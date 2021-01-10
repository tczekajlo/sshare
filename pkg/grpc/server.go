package grpc

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	pb "sshare/protobuf"

	"sshare/pkg/driver"
	"sshare/pkg/logger"
	"sshare/pkg/oauth2"
	"sshare/pkg/types"
	"sshare/pkg/version"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type createServer struct {
	driver types.DriverAdapter
	log    *zap.SugaredLogger
	pb.UnimplementedCreateServer
}

type deleteServer struct {
	driver types.DriverAdapter
	log    *zap.SugaredLogger
	pb.UnimplementedDeleteServer
}

type tlsServer struct {
	log    *zap.SugaredLogger
	CACert []byte
	pb.UnimplementedTLSServer
	oauth2 oauth2.Provider
}

func (s *tlsServer) Connection(ctx context.Context, data *pb.TLSRequest) (*pb.TLSResponse, error) {
	// Read metadata from client.
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		s.log.Errorw("Failed to get metadata")
	}
	streamID := md.Get("stream-id")[0]

	response := &pb.TLSResponse{}
	if data.Send {
		response.TLSServerPort = viper.GetInt32("server.tls-port")
		response.CACert = s.CACert
		response.AuthEnabled = viper.GetBool("server.auth-enabled")

		if viper.GetBool("server.auth-enabled") {
			response.AuthURL = s.oauth2.GetAuthURL(streamID)
			response.OAuth2ServerPort = viper.GetInt32("server.oauth2-port")
		}

		s.log.Infow("Received data", "data", data, "stream-id", streamID)
		s.log.Infow("Sent data", "data", response, "stream-id", streamID)

		return response, nil
	}
	return response, nil
}

func (s *deleteServer) Backend(ctx context.Context, data *pb.BackendData) (*pb.BackendReply, error) {
	// Read metadata from client.
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		s.log.Errorw("Failed to get metadata")
	}
	streamID := md.Get("stream-id")[0]

	if err := s.driver.Delete(data); err != nil {
		s.log.Errorw("Cannot delete the backend",
			"name", data.Name,
			"stream-id", streamID,
		)
		return &pb.BackendReply{Deleted: false}, nil
	}
	s.log.Infow("The backed has been deleted",
		"name", data.Name,
		"stream-id", streamID,
	)
	return &pb.BackendReply{Deleted: true}, nil
}

func (s *createServer) Backend(stream pb.Create_BackendServer) error {
	ctx := stream.Context()

	// Create trailer in defer to record function return time.
	defer func() {
		trailer := metadata.Pairs("timestamp", time.Now().Format(time.StampNano))
		stream.SetTrailer(trailer)
	}()

	// Read metadata from client.
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		s.log.Errorw("Failed to get metadata")
		return status.Errorf(codes.DataLoss, "Failed to get metadata")
	}

	// Create and send header.
	streamID := md.Get("stream-id")[0]
	header := metadata.New(map[string]string{"stream-id": streamID, "timestamp": time.Now().Format(time.StampNano)})
	stream.SendHeader(header)

	for {
		// Receive data from stream
		req, err := stream.Recv()
		if err == io.EOF {
			// return will close stream from server side
			s.log.Infow("Exit (close stream from server side)", "stream-id", streamID)
			return nil
		}
		if err != nil {
			s.log.Errorw("Receive error", "error", err, "stream-id", streamID)
			return status.Errorf(codes.DataLoss, "Failed to prepare a tunnel (stream-id: %s)", streamID)
		}
		s.log.Debugw("Received data", "data", req)

		// Prepare backend
		if err := s.driver.Create(req); err != nil {
			if errDel := s.driver.Delete(req); errDel != nil {
				s.log.Errorw("Cannot delete backend",
					"name", req.Name,
					"stream-id", streamID,
					"error", errDel,
				)
			} else {
				s.log.Infow("The backend has been deleted",
					"name", req.Name,
					"stream-id", streamID,
				)
			}
			return status.Errorf(codes.Internal, "Failed to prepare a tunnel, cannot create a backend (stream-id: %s)", streamID)
		}

		// Check if backend is ready
		ready, err := s.driver.IsReady(req)
		if err != nil {
			s.log.Errorw("the backend is not ready", "error", err)
			if errDel := s.driver.Delete(req); errDel != nil {
				s.log.Errorw("Cannot delete backend",
					"name", req.Name,
					"stream-id", streamID,
					"error", errDel,
				)
			} else {
				s.log.Infow("The backend has been deleted",
					"name", req.Name,
					"stream-id", streamID,
				)
			}
			return status.Errorf(codes.Internal, "Failed to prepare a tunnel, the backend is not ready (stream-id: %s)", streamID)
		}

		connData, err := s.driver.GetConnectionData(req)
		if err != nil {
			s.log.Errorw("Cannot get connection data", "error", err)
			return status.Errorf(codes.Internal, "Failed to prepare a tunnel, couldn't get connection data (stream-id: %s)", streamID)
		}

		connData.RemotePort = req.Connection.LocalPort
		resp := pb.BackendReply{
			Ready:                ready,
			Connection:           connData,
			ClientSessionTimeout: viper.GetInt32("server.client-session-timeout"),
		}
		if err := stream.Send(&resp); err != nil {
			s.log.Errorw("Send error", "error", err, "stream-id", streamID)
		}

		s.log.Infow("Sent reply",
			"stream-id", streamID,
			"ready", resp.Ready,
			"name", req.Name,
			"ssh-host", connData.SSHHost,
			"ssh-port", connData.SSHPort,
			"https-enabled", connData.HTTPScheme,
			"domain", connData.Domain,
			"client-session-timeout", resp.ClientSessionTimeout,
		)
	}
}

func tlsResponder(server *tlsServer) {
	if viper.GetBool("server.tls-enabled") {
		// read CA certs
		caData, err := ioutil.ReadFile(viper.GetString("server.tls-ca"))
		if err != nil {
			server.log.Fatalf("Cannot read CA certificate %v", err)
		}
		server.CACert = caData

		address := fmt.Sprintf("%s:%d", viper.GetString("server.address"), viper.GetInt32("server.port"))
		lis, err := net.Listen("tcp", address)
		if err != nil {
			server.log.Fatalw("Failed to listen", "error", err)
		}

		server.log.Infow("TLS Responder",
			"address", address)

		grpcServer := grpc.NewServer()

		pb.RegisterTLSServer(grpcServer, server)

		if err := grpcServer.Serve(lis); err != nil {
			server.log.Fatalw("Failed to serve", "error", err)
		}
	} else {
		server.log.Infow("TLS is disabled. Skipping TLS Responder")
	}
}

// RunServer runs sshare server
func RunServer() {
	log := logger.GetInstance()
	port := viper.GetInt32("server.port")
	var serverCreds grpc.ServerOption

	driverInstance := driver.Driver{}
	driver := driverInstance.New()

	createServer := &createServer{
		driver: driver,
		log:    log,
	}

	deleteServer := &deleteServer{
		driver: driver,
		log:    log,
	}

	tlsServer := &tlsServer{
		log: log,
	}

	oauth2Server := &oauth2Server{
		log: log,
	}

	streamInterceptors := []grpc.StreamServerInterceptor{
		grpc_prometheus.StreamServerInterceptor,
	}
	unaryInterceptors := []grpc.UnaryServerInterceptor{
		grpc_prometheus.UnaryServerInterceptor,
	}

	if viper.GetBool("server.auth-enabled") {
		streamInterceptors = append(streamInterceptors, streamEnsureValidTokenInterceptor)
		unaryInterceptors = append(unaryInterceptors, unaryEnsureValidTokenInterceptor)

		log.Debug("Authorization is enabled")
	}

	opts := []grpc.ServerOption{
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				streamInterceptors...,
			),
		),
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				unaryInterceptors...,
			),
		),
	}

	if viper.GetBool("server.tls-enabled") {
		port = viper.GetInt32("server.tls-port")
		creds, err := credentials.NewServerTLSFromFile(
			viper.GetString("server.tls-cert"),
			viper.GetString("server.tls-key"),
		)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		serverCreds = grpc.Creds(creds)
		opts = append(opts, serverCreds)
	}

	address := fmt.Sprintf("%s:%d", viper.GetString("server.address"), port)

	log.Infow("sshare gRPC server",
		"version", version.VERSION,
		"address", address)

	go tlsResponder(tlsServer)

	// check if TLS is enabled
	if !viper.GetBool("server.tls-enabled") && viper.GetBool("server.auth-enabled") {
		log.Error("Authentication can't be enabled without TLS")
	} else if viper.GetBool("server.auth-enabled") {
		oauth2Provider := oauth2.New()
		tlsServer.oauth2 = oauth2Provider
		oauth2Server.oauth2Provider = oauth2Provider

		log.Info("OAuth2 enabled")
		go oauth2Exchanger(oauth2Server, serverCreds)
	}

	// Run prometheus metrics
	go func() {
		metricsAddress := fmt.Sprintf(":%d", viper.GetInt32("server.metrics-port"))
		metricsEndpoint := "/metrics"
		log.Infow("Running Prometheus metrics",
			"address", metricsAddress,
			"endpoint", metricsEndpoint,
		)
		http.Handle(metricsEndpoint, promhttp.Handler())
		http.ListenAndServe(metricsAddress, nil)
	}()

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalw("Failed to listen", "error", err)
	}

	grpcServer := grpc.NewServer(opts...)

	pb.RegisterCreateServer(grpcServer, createServer)
	pb.RegisterDeleteServer(grpcServer, deleteServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalw("Failed to serve", "error", err)
	}

}
