package grpc

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sshare/pkg/metrics"
	"strings"
	"time"

	"github.com/pkg/browser"
	"github.com/spf13/viper"
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

func (s *SshareClient) handleCallback(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("#%v", r)
	//content, _ := getUserInfo(, r.FormValue("code"))

	if s.TLS.AuthState != r.FormValue("state") {
		s.log.Error("Invalid oauth state")
		os.Exit(1)
	}

	fmt.Fprintf(w, "Success, please go back to the terminal\n")
	close(callbackDone)
	return
}

func (s *SshareClient) oauth2Auth() {
	if s.TLS.AuthEnabled {
		//ctx := context.Background()

		//http.HandleFunc("/callback", s.handleCallback)
		//server := &http.Server{Addr: ":8080"}
		//go server.ListenAndServe()

		fmt.Println("Server requires authentication, please use the following link:", s.TLS.AuthURL)
		browser.OpenURL(s.TLS.AuthURL)
		time.Sleep(1000 * time.Second)
		//<-callbackDone
		//server.Shutdown(ctx)
	}

}

func fetchToken() *oauth2.Token {
	return &oauth2.Token{
		AccessToken: viper.GetString("client.token"),
	}
}

// validToken validates the authorization.
func validToken(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	return token == viper.GetString("server.auth-token")
}

func unaryEnsureValidTokenInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}
	// The keys within metadata.MD are normalized to lowercase.
	// See: https://godoc.org/google.golang.org/grpc/metadata#New
	if !validToken(md["authorization"]) {
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
	if !validToken(md["authorization"]) {
		metrics.AuthenticationTotal.WithLabelValues("invalid_token", "stream").Inc()

		return errInvalidToken
	}

	metrics.AuthenticationTotal.WithLabelValues("success", "stream").Inc()

	return handler(srv, ss)
}
