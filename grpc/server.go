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

	"sshare/driver"
	"sshare/logger"
	"sshare/types"
	"sshare/version"

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
}

func (s *tlsServer) Connection(ctx context.Context, data *pb.TLSRequest) (*pb.TLSResponse, error) {
	// Read metadata from client.
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		s.log.Errorw("failed to get metadata")
	}
	streamID := md.Get("stream-id")[0]

	response := &pb.TLSResponse{}
	if data.Send {
		response.TLSServerPort = viper.GetInt32("server.tls-port")
		response.CACert = s.CACert

		s.log.Infow("received data", "data", data, "stream-id", streamID)
		s.log.Infow("sent data", "data", response, "stream-id", streamID)

		return response, nil
	}
	return response, nil
}

func (s *deleteServer) Backend(ctx context.Context, data *pb.BackendData) (*pb.BackendReply, error) {
	// Read metadata from client.
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		s.log.Errorw("failed to get metadata")
	}
	streamID := md.Get("stream-id")[0]

	if err := s.driver.Delete(data); err != nil {
		s.log.Errorw("cannot delete the backend",
			"name", data.Name,
			"stream-id", streamID,
		)
		return &pb.BackendReply{Deleted: false}, nil
	}
	s.log.Infow("the backed has been deleted",
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
		s.log.Errorw("failed to get metadata")
		return status.Errorf(codes.DataLoss, "failed to get metadata")
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
			s.log.Infow("exit (close stream from server side)", "stream-id", streamID)
			return nil
		}
		if err != nil {
			s.log.Errorw("receive error", "error", err, "stream-id", streamID)
			return status.Errorf(codes.DataLoss, "failed to prepare a tunnel (stream-id: %s)", streamID)
		}
		s.log.Debugw("received data", "data", req)

		// Prepare backend
		if err := s.driver.Create(req); err != nil {
			if errDel := s.driver.Delete(req); errDel != nil {
				s.log.Errorw("cannot delete backend",
					"name", req.Name,
					"stream-id", streamID,
					"error", errDel,
				)
			} else {
				s.log.Infow("the backend has been deleted",
					"name", req.Name,
					"stream-id", streamID,
				)
			}
			return status.Errorf(codes.Internal, "failed to prepare a tunnel, cannot create a backend (stream-id: %s)", streamID)
		}

		// Check if backend is ready
		ready, err := s.driver.IsReady(req)
		if err != nil {
			s.log.Errorw("the backend is not ready", "error", err)
			if errDel := s.driver.Delete(req); errDel != nil {
				s.log.Errorw("cannot delete backend",
					"name", req.Name,
					"stream-id", streamID,
					"error", errDel,
				)
			} else {
				s.log.Infow("the backend has been deleted",
					"name", req.Name,
					"stream-id", streamID,
				)
			}
			return status.Errorf(codes.Internal, "failed to prepare a tunnel, the backend is not ready (stream-id: %s)", streamID)
		}

		connData, err := s.driver.GetConnectionData(req)
		if err != nil {
			s.log.Errorw("cannot get connection data", "error", err)
			return status.Errorf(codes.Internal, "failed to prepare a tunnel, couldn't get connection data (stream-id: %s)", streamID)
		}

		connData.RemotePort = req.Connection.LocalPort
		resp := pb.BackendReply{
			Ready:                ready,
			Connection:           connData,
			ClientSessionTimeout: viper.GetInt32("server.client-session-timeout"),
		}
		if err := stream.Send(&resp); err != nil {
			s.log.Errorw("send error", "error", err, "stream-id", streamID)
		}

		s.log.Infow("sent reply",
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
			server.log.Fatalf("cannot read CA certificate %v", err)
		}
		server.CACert = caData

		address := fmt.Sprintf("%s:%d", viper.GetString("server.address"), viper.GetInt32("server.port"))
		lis, err := net.Listen("tcp", address)
		if err != nil {
			server.log.Fatalw("failed to listen", "error", err)
		}

		server.log.Infow("TLS Responder",
			"address", address)

		grpcServer := grpc.NewServer()

		pb.RegisterTLSServer(grpcServer, server)

		if err := grpcServer.Serve(lis); err != nil {
			server.log.Fatalw("failed to serve", "error", err)
		}
	} else {
		server.log.Infow("TLS is disabled. Skipping TLS Responder")
	}
}

// RunServer runs sshare server
func RunServer() {
	log := logger.GetInstance()
	port := viper.GetInt32("server.port")

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

	var streamInterceptors []grpc.StreamServerInterceptor = []grpc.StreamServerInterceptor{
		grpc_prometheus.StreamServerInterceptor,
	}
	var unaryInterceptors []grpc.UnaryServerInterceptor = []grpc.UnaryServerInterceptor{
		grpc_prometheus.UnaryServerInterceptor,
	}

	if viper.GetString("server.auth-token") != "" {
		streamInterceptors = append(streamInterceptors, streamEnsureValidTokenInterceptor)
		unaryInterceptors = append(unaryInterceptors, unaryEnsureValidTokenInterceptor)

		log.Debug("authorization token is set")
	}

	var opts []grpc.ServerOption = []grpc.ServerOption{
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
		opts = append(opts, grpc.Creds(creds))
	}

	address := fmt.Sprintf("%s:%d", viper.GetString("server.address"), port)

	log.Infow("sshare gRPC server",
		"version", version.VERSION,
		"address", address)

	go tlsResponder(tlsServer)

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
		log.Fatalw("failed to listen", "error", err)
	}

	grpcServer := grpc.NewServer(opts...)

	pb.RegisterCreateServer(grpcServer, createServer)
	pb.RegisterDeleteServer(grpcServer, deleteServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalw("failed to serve", "error", err)
	}

}
