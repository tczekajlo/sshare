package grpc

import (
	"context"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"sshare/pkg/logger"
	"sshare/pkg/ssh"
	"sshare/pkg/version"
	pb "sshare/protobuf"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/google/uuid"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/kyokomi/emoji"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/grpc/metadata"
)

func retryDial(sshareClient *SshareClient, address string) (*grpc.ClientConn, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		log.Fatalf("Failed to parse address %v", err)
	}

	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(100 * time.Millisecond)),
		grpc_retry.WithCodes(codes.DataLoss, codes.Aborted, codes.Unavailable),
	}

	opts := []grpc.DialOption{
		grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(retryOpts...)),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...)),
	}

	if sshareClient.TLS.AuthEnabled && sshareClient.oauth2Token != nil {
		// Set up the credentials for the connection.
		sshareClient.log.Debug("Set up the credentials for the connection")
		authPerRPC := oauth.NewOauthAccess(sshareClient.oauth2Token)
		opts = append(opts, grpc.WithPerRPCCredentials(authPerRPC))
	}

	if !viper.GetBool("client.tls-disabled") {
		// Create a certificate pool from the certificate authority
		certPool := x509.NewCertPool()
		// Append the client certificates from the CA
		if ok := certPool.AppendCertsFromPEM(sshareClient.TLS.CACert); !ok {
			log.Fatal("Failed to append client certs")
		}

		creds := credentials.NewClientTLSFromCert(certPool, "")
		opts = append(opts, grpc.WithTransportCredentials(creds))

		if port != strconv.Itoa(int(sshareClient.TLS.OAuth2ServerPort)) {
			address = fmt.Sprintf("%s:%d", host, sshareClient.TLS.TLSServerPort)
		}
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	return grpc.Dial(address, opts...)
}

// SshareClient stores data for sshare client
type SshareClient struct {
	log               *zap.SugaredLogger
	streamID          string
	backendConnection *pb.Connection
	instanceName      string
	waitSpinner       *spinner.Spinner
	sshKeys           *ssh.Keys
	sshPublicKey      string
	streamClient      pb.Create_BackendClient
	TLS               *pb.TLSResponse
	localPort         int32
	onlyTCP           bool
	sessionTimeout    int32
	oauth2Code        string
	oauth2Token       *oauth2.Token
}

func (s *SshareClient) generateSSHKeys() {
	s.waitSpinner.Suffix = " Creating key pair..."
	s.waitSpinner.Start()

	s.sshKeys.Init()

	sshPublicKey, err := s.sshKeys.GetPublicKey()
	if err != nil {
		s.log.Fatalw("Cannot get SSH public key", "error", err)
	}
	s.sshPublicKey = sshPublicKey

	s.log.Debugw("SSH public key",
		"public-key", sshPublicKey,
		"private-key", s.sshKeys.PrivateKeyPEM,
	)
}

func (s *SshareClient) spinnerNewMsg(msg string) {
	s.waitSpinner.Suffix = emoji.Sprint(msg)
	s.waitSpinner.Restart()
}

func (s *SshareClient) initStreamClient(client pb.CreateClient) {
	md := metadata.Pairs(
		"timestamp", time.Now().Format(time.StampNano),
		"stream-id", s.streamID)

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// Make RPC using the context with the metadata
	stream, err := client.Backend(ctx)
	if err != nil {
		log.Fatalf("Open stream error: %v", err)
	}

	s.streamClient = stream
}

func (s *SshareClient) runSender() {
	req := pb.BackendData{
		Name:         s.instanceName,
		StreamID:     s.streamID,
		SshPublicKey: s.sshPublicKey,
		HTTPOptions: &pb.HTTPOptions{
			CORSEnabled:   viper.GetBool("http-enable-cors"),
			HTTPSRedirect: viper.GetBool("https-redirect"),
		},
		OnlyTCP: s.onlyTCP,
		Connection: &pb.Connection{
			LocalPort: s.localPort,
		},
	}

	if err := s.streamClient.Send(&req); err != nil {
		log.Fatalf("Can not send: %v", err)
	}
	s.log.Debugw("Sent data", "name", req.Name, "stream-id", req.StreamID)
	s.streamClient.CloseSend()
}

func (s *SshareClient) runReceiver() {
	// Read the header when the header arrives.
	header, err := s.streamClient.Header()
	if err != nil {
		log.Fatalf("Failed to get header from stream: %v", err)
	}

	s.log.Debugw("Received metadata", "header", header)

	for {
		resp, err := s.streamClient.Recv()
		if err == io.EOF {
			// read done
			return
		}
		if err != nil {
			log.Fatalf("Can not receive: %v", err)
		}
		s.log.Debugw("Received data",
			"connection", resp.Connection,
			"ready", resp.Ready,
		)
		if resp.Ready {
			s.backendConnection = resp.Connection
			s.backendConnection.LocalPort = s.localPort
			s.sessionTimeout = resp.ClientSessionTimeout
		} else {
			s.waitSpinner.FinalMSG = "Something went wrong :(\n"
			s.waitSpinner.Restart()
			s.waitSpinner.Stop()
		}
	}
}

func (s *SshareClient) sessionTimeoutClose(sigs chan os.Signal) {
	if s.sessionTimeout != 0 {
		timeout := time.After(time.Duration(s.sessionTimeout) * time.Second)
		ticker := time.Tick(5 * time.Second)
		now := time.Now()
		end := now.Add(time.Duration(s.sessionTimeout) * time.Second)

		s.log.Debugw("Session timeout is set", "timeout", s.sessionTimeout)

		s.waitSpinner.UpdateCharSet(spinner.CharSets[21])
		s.waitSpinner.Start()

		for {
			select {
			// Got a timeout!
			case <-timeout:
				s.waitSpinner.Stop()

				fmt.Println(color.YellowString(emoji.Sprintf("Session timed out :clock: The server that you're connected to allows for a session no longer than %v",
					time.Duration(s.sessionTimeout)*time.Second)))

				sigs <- syscall.SIGTERM
				return
			// Got a tick
			case <-ticker:
				diff := end.Sub(time.Now())
				timeLeft := time.Time{}.Add(diff).Format("15:04:05")
				s.spinnerNewMsg(
					color.YellowString(
						emoji.Sprintf(" Your session is time-limited by the server that you're connected to: %v left",
							timeLeft),
					),
				)
			}
		}

	}
}

func handleSignals(sigs chan os.Signal, sshareClient *SshareClient) {
	<-sigs
	clean(sshareClient)
}

func clean(sshareClient *SshareClient) {
	conn, err := retryDial(sshareClient, viper.GetString("client.server-address"))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer func() {
		if e := conn.Close(); e != nil {
			log.Printf("Failed to close connection: %s", e)
		}
	}()
	// Make a greeter client and send an RPC.
	client := pb.NewDeleteClient(conn)
	cleanSend(client, sshareClient)

	fmt.Println()
	sshareClient.waitSpinner.Stop()
	emoji.Println("Bye :wave:")
	conn.Close()
	os.Exit(0)
}

func cleanSend(c pb.DeleteClient, sshareClient *SshareClient) {
	md := metadata.Pairs(
		"timestamp", time.Now().Format(time.StampNano),
		"stream-id", sshareClient.streamID)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ctxMetadata := metadata.NewOutgoingContext(ctx, md)

	c.Backend(
		ctxMetadata,
		&pb.BackendData{
			Name:     sshareClient.instanceName,
			StreamID: sshareClient.streamID,
			OnlyTCP:  sshareClient.onlyTCP,
		},
		grpc_retry.WithMax(3),
		grpc.WaitForReady(true),
	)
}

func tlsRequester(streamID string) *pb.TLSResponse {
	log := logger.GetInstance()
	// Set up a connection to the server.
	conn, err := grpc.Dial(viper.GetString("client.server-address"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTLSClient(conn)

	md := metadata.Pairs(
		"timestamp", time.Now().Format(time.StampNano),
		"stream-id", streamID)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ctxMetadata := metadata.NewOutgoingContext(ctx, md)
	r, err := c.Connection(ctxMetadata, &pb.TLSRequest{Send: true}, grpc.WaitForReady(true))
	if err != nil {
		log.Fatalf("Could not get data for TLS connection: %v", err)
	}
	log.Debugw("Received data from TLS Responder", "data", r)

	return r
}

func (s *SshareClient) checkLocalListener() {
	if s.localPort == 64000 {
		fmt.Printf("%s\n", color.RedString("Port 64000 is reserved for OAuth authentication, please use different port"))
		os.Exit(1)
	}

	local, err := net.Dial("tcp", fmt.Sprintf("0.0.0.0:%d", s.localPort))
	if err != nil {
		fmt.Printf("%v: %s\n", color.RedString("Cannot dial into local service"), err)
		os.Exit(1)
	}
	defer local.Close()
}

func (s *SshareClient) printAccessData() {
	boldColorCyan := color.New(color.FgCyan, color.Bold)
	if viper.GetBool("client.tcp") {
		fmt.Printf("%v: %s:%d %v 0.0.0.0:%d\n",
			boldColorCyan.Sprintf("Address"),
			s.backendConnection.SSHHost,
			s.backendConnection.LocalPort,
			color.YellowString("->"),
			s.backendConnection.LocalPort,
		)
	} else {
		if s.backendConnection.HTTPScheme {
			fmt.Printf("%v: https://%s %v 0.0.0.0:%d\n",
				boldColorCyan.Sprintf("Address"),
				s.backendConnection.Domain,
				color.YellowString("->"),
				s.backendConnection.LocalPort,
			)
			if viper.GetBool("https-redirect") {
				fmt.Printf("%v: %v\n", boldColorCyan.Sprintf("HTTPs redirect"), color.GreenString("enabled"))
				fmt.Printf("%v: http://%s %v https://%s\n",
					boldColorCyan.Sprintf("Address"),
					s.backendConnection.Domain,
					color.YellowString("->"),
					s.backendConnection.Domain,
				)
			} else {
				fmt.Printf("%v: http://%s %v 0.0.0.0:%d\n",
					boldColorCyan.Sprintf("Address"),
					s.backendConnection.Domain,
					color.YellowString("->"),
					s.backendConnection.LocalPort,
				)
			}
		} else {
			fmt.Printf("%v: http://%s %v 0.0.0.0:%d\n",
				boldColorCyan.Sprintf("Address"),
				s.backendConnection.Domain,
				color.YellowString("->"),
				s.backendConnection.LocalPort,
			)
		}
	}
	fmt.Println()
}

// RunClient runs gRPC client and establishes SSH tunnel
func RunClient(localPort int32) {
	sigs := make(chan os.Signal, 1)
	sshareClient := &SshareClient{
		log:          logger.GetInstance(),
		streamID:     uuid.New().String(),
		instanceName: uuid.New().String(),
		waitSpinner:  spinner.New(spinner.CharSets[11], 100*time.Millisecond),
		sshKeys:      &ssh.Keys{},
		localPort:    localPort,
		onlyTCP:      viper.GetBool("client.tcp"),
		TLS:          &pb.TLSResponse{},
	}

	emoji.Println(fmt.Sprintf("sshare %s :rocket:", version.VERSION))

	// Check local service
	sshareClient.checkLocalListener()

	// Generate SSH keys
	sshareClient.generateSSHKeys()

	if !viper.GetBool("client.tls-disabled") {
		sshareClient.spinnerNewMsg(" Requesting CA cert for securing connection...")
		sshareClient.TLS = tlsRequester(sshareClient.streamID)

		sshareClient.spinnerNewMsg(" Waiting for the authentication to be finished...")
		if err := sshareClient.oauth2Auth(); err != nil {
			sshareClient.waitSpinner.Stop()
			log.Fatalf("Cannot finish authentication: %v", err)
			os.Exit(1)
		}
	}

	sshareClient.spinnerNewMsg(" Preparing a secure tunnel...")

	// Set up a connection to the server.
	conn, err := retryDial(sshareClient, viper.GetString("client.server-address"))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func() {
		if e := conn.Close(); e != nil {
			log.Printf("Failed to close connection: %s", e)
		}
	}()

	// Create client
	client := pb.NewCreateClient(conn)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go handleSignals(sigs, sshareClient)

	sshareClient.initStreamClient(client)
	sshareClient.runSender()
	sshareClient.runReceiver()

	tunnelReady := make(chan bool)
	tunnel := &ssh.Tunnel{
		User:          "sshare",
		Log:           sshareClient.log,
		Connection:    sshareClient.backendConnection,
		WaitSpinner:   sshareClient.waitSpinner,
		PrivateKeySSH: sshareClient.sshKeys.PrivateKeyPEM,
		Ready:         tunnelReady,
	}

	go func() {
		<-tunnel.Ready
		fmt.Println()
		sshareClient.printAccessData()
		return
	}()

	go func() {
		<-tunnel.Ready
		sshareClient.sessionTimeoutClose(sigs)
		return
	}()

	if err := tunnel.ReverseTunnel(); err != nil {
		sshareClient.log.Errorw("Cannot establish the tunnel", "error", err)
		clean(sshareClient)
	}

}
