package grpc

import (
	"context"
	"fmt"
	"net"
	pb "sshare/protobuf"
	"strings"

	"sshare/pkg/metrics"
	oauth2Provider "sshare/pkg/oauth2"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	errMissingMetadata           = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken              = status.Errorf(codes.Unauthenticated, "invalid token")
	callbackDone       chan bool = make(chan bool)
)

type oauth2Server struct {
	log            *zap.SugaredLogger
	oauth2Provider oauth2Provider.Provider
	pb.UnimplementedOAuth2Server
}

func (o *oauth2Server) Exchange(ctx context.Context, data *pb.OAuth2Request) (*pb.OAuth2Response, error) {
	// Read metadata from client.
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		o.log.Errorw("Failed to get metadata")
	}
	streamID := md.Get("stream-id")[0]

	token, err := o.oauth2Provider.Exchange(data.Code)
	if err != nil {
		return nil, fmt.Errorf("Code exchange failed: %s", err.Error())
	}

	response := &pb.OAuth2Response{
		Token: token.AccessToken,
	}

	o.log.Infow("Received data", "data", data, "stream-id", streamID)
	o.log.Infow("Sent data", "data", response, "stream-id", streamID)

	return response, nil
}

func oauth2Exchanger(server *oauth2Server, opts ...grpc.ServerOption) {
	// Run OAuth2 Exchanger only if TLS is enabled
	if viper.GetBool("server.tls-enabled") {
		address := fmt.Sprintf("%s:%d", viper.GetString("server.address"), viper.GetInt32("server.oauth2-port"))
		lis, err := net.Listen("tcp", address)
		if err != nil {
			server.log.Fatalw("Failed to listen", "error", err)
		}

		server.log.Infow("OAuth2 Exchanger",
			"address", address)

		grpcServer := grpc.NewServer(opts...)

		// Register reflection service on gRPC server.
		pb.RegisterOAuth2Server(grpcServer, server)

		if err := grpcServer.Serve(lis); err != nil {
			server.log.Fatalw("Failed to serve", "error", err)
		}
	} else {
		server.log.Infow("TLS is disabled. Skipping OAuth2 Exchanger")
	}
}

// validToken validates the authorization.
func validToken(authorization []string) (bool, error) {
	if len(authorization) < 1 {
		return false, nil
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")

	oauth2Provider := oauth2Provider.New()
	isValid, err := oauth2Provider.ValidateToken(&oauth2.Token{
		AccessToken: token,
	})
	return isValid, err
}

func unaryEnsureValidTokenInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}
	// The keys within metadata.MD are normalized to lowercase.
	// See: https://godoc.org/google.golang.org/grpc/metadata#New
	isValid, err := validToken(md["authorization"])
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "%s, try to use the --login flag", err)
	}
	if !isValid {
		metrics.AuthenticationTotal.WithLabelValues("invalid_token", "unary").Inc()

		return nil, errInvalidToken
	}
	// Continue execution of handler after ensuring a valid token.
	metrics.AuthenticationTotal.WithLabelValues("success", "unary").Inc()

	return handler(ctx, req)
}

func streamEnsureValidTokenInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	// authentication (token verification)
	md, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		return errMissingMetadata
	}
	isValid, err := validToken(md["authorization"])
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "%s, try to use the --login flag", err)
	}
	if !isValid {
		metrics.AuthenticationTotal.WithLabelValues("invalid_token", "stream").Inc()

		return errInvalidToken
	}

	metrics.AuthenticationTotal.WithLabelValues("success", "stream").Inc()

	return handler(srv, ss)
}
