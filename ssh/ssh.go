package ssh

import (
	"fmt"
	"log"
	"net"
	"time"

	pb "sshare/protobuf"

	"sshare/proxy"

	"github.com/briandowns/spinner"
	"github.com/kyokomi/emoji"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
)

// Tunnel stores data needed to run SSH tunnel
type Tunnel struct {
	Connection    *pb.Connection
	User          string
	Log           *zap.SugaredLogger
	WaitSpinner   *spinner.Spinner
	PrivateKeySSH string
	Ready         chan bool
}

// ReverseTunnel creates a SSH tunnel for forwarding traffic from local service to remote server
func (t *Tunnel) ReverseTunnel() error {
	sshConfig := &ssh.ClientConfig{
		User: t.User,
		Auth: []ssh.AuthMethod{
			privateKeySSH(t.PrivateKeySSH),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	proxy := &proxy.Proxy{
		Done:       make(chan struct{}),
		Connection: t.Connection,
		Log:        t.Log,
	}

	// Connect to SSH remote server
	serverConn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", t.Connection.SSHHost, t.Connection.SSHPort), sshConfig)
	if err != nil {
		return fmt.Errorf("Cannot dial into remote server: %s", err)
	}
	defer serverConn.Close()

	// Listen on remote server port
	listener, err := serverConn.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", t.Connection.RemotePort))
	if err != nil {
		return fmt.Errorf("Cannot listen on remote server: %s", err)
	}
	defer listener.Close()

	ready, err := t.checkIfRemoteServerIsReady()
	if err != nil {
		t.WaitSpinner.Stop()
		return fmt.Errorf("Remote listener is not ready: %s", err)
	}

	if ready {
		t.WaitSpinner.Stop()
		close(t.Ready)
		proxy.Run(listener)
	}

	return nil
}

func (t *Tunnel) checkIfRemoteServerIsReady() (bool, error) {
	timeout := time.After(120 * time.Second)
	ticker := time.Tick(1000 * time.Millisecond)

	t.WaitSpinner.Suffix = emoji.Sprint(" Waiting for the tunnel to be ready")
	t.WaitSpinner.Restart()

	t.Log.Debug("Checking if remote server is ready...")

	// Keep trying until we're time out or get a result or get an error
	for {
		select {
		// Got a timeout! fail with a timeout error
		case <-timeout:
			return false, fmt.Errorf("timed out")
		// Got a tick, we should check on checkSomething()
		case <-ticker:
			conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", t.Connection.SSHHost, t.Connection.SSHPort))
			if err != nil {
				t.Log.Debug("Cannot dial into remote server, not ready", err)
			} else {
				t.Log.Debug("Remote server is ready")
				return true, nil
			}
			defer conn.Close()
		}
	}
}

func privateKeySSH(key string) ssh.AuthMethod {
	keyParsed, err := ssh.ParsePrivateKey([]byte(key))
	if err != nil {
		log.Fatalln(fmt.Sprintf("Cannot parse SSH public key: %s", err))
		return nil
	}
	return ssh.PublicKeys(keyParsed)
}
