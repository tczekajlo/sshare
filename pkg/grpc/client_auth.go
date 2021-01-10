package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	pb "sshare/protobuf"
	"time"

	"github.com/pkg/browser"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"k8s.io/client-go/util/homedir"
)

var tokenFile string = filepath.Join(homedir.HomeDir(), ".sshare", "auth.json")

func (s *SshareClient) oauth2TokenRequester() *pb.OAuth2Response {
	// Set up a connection to the server.
	host, _, _ := net.SplitHostPort(viper.GetString("client.server-address"))
	address := fmt.Sprintf("%s:%d", host, s.TLS.OAuth2ServerPort)

	conn, err := retryDial(s, address)
	if err != nil {
		s.log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewOAuth2Client(conn)

	md := metadata.Pairs(
		"timestamp", time.Now().Format(time.StampNano),
		"stream-id", s.streamID)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ctxMetadata := metadata.NewOutgoingContext(ctx, md)
	r, err := c.Exchange(ctxMetadata, &pb.OAuth2Request{Code: s.oauth2Code}, grpc.WaitForReady(true))
	if err != nil {
		s.log.Fatalf("Could not connect to OAuth2 Exchanger: %v", err)
	}
	s.log.Debugw("Received data from OAuth2 Exchanger", "data", r)

	return r
}

func (s *SshareClient) handleCallback(w http.ResponseWriter, r *http.Request) {
	s.oauth2Code = r.FormValue("code")

	if s.streamID != r.FormValue("state") {
		fmt.Fprintf(w, "Something went wrong, please go back to the terminal for more details\n")
		s.log.Error("Invalid OAuth state")
		os.Exit(1)
	}

	fmt.Fprintf(w, "Success, please go back to the terminal\n")
	close(callbackDone)
	return
}

func (s *SshareClient) oauth2Auth() error {
	if s.TLS.AuthEnabled {

		if !viper.GetBool("client.login") {
			err := s.readStoredToken()
			if err == nil {
				s.log.Debug("Using stored token")
				return nil
			}
			s.log.Debugw("Cannot read stored token", "file", tokenFile, "error", err)
		}

		ctx := context.Background()

		http.HandleFunc("/callback", s.handleCallback)
		server := &http.Server{Addr: ":64000"}
		go server.ListenAndServe()

		fmt.Println("Server requires authentication, please use the following link:", s.TLS.AuthURL)
		browser.OpenURL(s.TLS.AuthURL)

		<-callbackDone
		server.Shutdown(ctx)

		// Exchange code
		response := s.oauth2TokenRequester()
		if err := s.saveToken(response); err != nil {
			return err
		}
	}
	return nil
}

func (s *SshareClient) saveToken(token *pb.OAuth2Response) error {
	jsonString, err := json.Marshal(token)
	if err != nil {
		return err
	}

	_, err = os.Stat(tokenFile)
	if os.IsNotExist(err) {
		os.MkdirAll(filepath.Join(homedir.HomeDir(), ".sshare"), os.ModePerm)
	}

	if err := ioutil.WriteFile(tokenFile, jsonString, os.ModePerm); err != nil {
		return err
	}

	s.oauth2Token = &oauth2.Token{
		AccessToken: token.Token,
	}

	return nil
}

func (s *SshareClient) readStoredToken() error {

	jsonFile, err := os.Open(tokenFile)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	// Unmarshall JSON
	data, _ := ioutil.ReadAll(jsonFile)

	var auth *pb.OAuth2Response

	err = json.Unmarshal(data, &auth)
	if err != nil {
		return err
	}

	s.oauth2Token = &oauth2.Token{
		AccessToken: auth.Token,
	}

	return nil
}
