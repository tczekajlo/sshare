package proxy

import (
	"fmt"
	"io"
	"net"
	"sync"

	pb "sshare/protobuf"

	"go.uber.org/zap"
)

// Proxy stores data for proxy
type Proxy struct {
	Done       chan struct{}
	Connection *pb.Connection
	Log        *zap.SugaredLogger
}

// Run runs proxy
func (p *Proxy) Run(listener net.Listener) {
	for {
		select {
		case <-p.Done:
			return
		default:
			connection, err := listener.Accept()
			if err == nil {
				go p.forward(connection)
			} else {
				p.Log.Errorw("error accepting connection", "error", err)
			}
		}
	}
}

func (p *Proxy) forward(connection net.Conn) {
	p.Log.Debug("Forwarding")
	defer p.Log.Debug("Done forwarding")
	defer connection.Close()

	local, err := net.Dial("tcp", fmt.Sprintf("0.0.0.0:%d", p.Connection.LocalPort))
	if err != nil {
		p.Log.Errorf("Error dialing local host: %v", err)
		return
	}

	defer local.Close()
	wg := &sync.WaitGroup{}

	wg.Add(2)
	go p.copy(local, connection, wg)
	go p.copy(connection, local, wg)
	wg.Wait()
}

func (p *Proxy) copy(from, to net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	select {
	case <-p.Done:
		return
	default:
		if _, err := io.Copy(to, from); err != nil {
			p.Log.Errorf("cannot copy connection: %v", err)
			p.Stop()
			return
		}
	}
}

// Stop stops proxy
func (p *Proxy) Stop() {
	p.Log.Debug("stopping proxy")
	if p.Done == nil {
		return
	}
	close(p.Done)
	p.Done = nil
}
