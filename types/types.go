package types

import (
	pb "sshare/protobuf"
)

// Connection stores data returned by server that are needed to establish a tunnel connection
type Connection struct {
	SSHHost    string
	SSHPort    int32
	Domain     string
	HTTPScheme bool
	RemotePort int32
	LocalPort  int32
}

// DriverAdapter returns functions for a driver
type DriverAdapter interface {
	Create(data *pb.BackendData, opts ...interface{}) error
	IsReady(data *pb.BackendData, opts ...interface{}) (bool, error)
	GetConnectionData(data *pb.BackendData, opts ...interface{}) (*pb.Connection, error)
	Delete(data *pb.BackendData, opts ...interface{}) error
}
