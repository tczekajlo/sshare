// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: sshare.proto

package sshare

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type TLSRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Send bool `protobuf:"varint,1,opt,name=send,proto3" json:"send,omitempty"`
}

func (x *TLSRequest) Reset() {
	*x = TLSRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sshare_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TLSRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TLSRequest) ProtoMessage() {}

func (x *TLSRequest) ProtoReflect() protoreflect.Message {
	mi := &file_sshare_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TLSRequest.ProtoReflect.Descriptor instead.
func (*TLSRequest) Descriptor() ([]byte, []int) {
	return file_sshare_proto_rawDescGZIP(), []int{0}
}

func (x *TLSRequest) GetSend() bool {
	if x != nil {
		return x.Send
	}
	return false
}

type TLSResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CACert        []byte `protobuf:"bytes,1,opt,name=CACert,proto3" json:"CACert,omitempty"`
	TLSServerPort int32  `protobuf:"varint,2,opt,name=TLSServerPort,proto3" json:"TLSServerPort,omitempty"`
}

func (x *TLSResponse) Reset() {
	*x = TLSResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sshare_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TLSResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TLSResponse) ProtoMessage() {}

func (x *TLSResponse) ProtoReflect() protoreflect.Message {
	mi := &file_sshare_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TLSResponse.ProtoReflect.Descriptor instead.
func (*TLSResponse) Descriptor() ([]byte, []int) {
	return file_sshare_proto_rawDescGZIP(), []int{1}
}

func (x *TLSResponse) GetCACert() []byte {
	if x != nil {
		return x.CACert
	}
	return nil
}

func (x *TLSResponse) GetTLSServerPort() int32 {
	if x != nil {
		return x.TLSServerPort
	}
	return 0
}

type BackendData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StreamID     string       `protobuf:"bytes,1,opt,name=streamID,proto3" json:"streamID,omitempty"`
	Name         string       `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	SshPublicKey string       `protobuf:"bytes,3,opt,name=sshPublicKey,proto3" json:"sshPublicKey,omitempty"`
	HTTPOptions  *HTTPOptions `protobuf:"bytes,4,opt,name=HTTPOptions,proto3" json:"HTTPOptions,omitempty"`
	OnlyTCP      bool         `protobuf:"varint,5,opt,name=onlyTCP,proto3" json:"onlyTCP,omitempty"`
	Connection   *Connection  `protobuf:"bytes,6,opt,name=connection,proto3" json:"connection,omitempty"`
}

func (x *BackendData) Reset() {
	*x = BackendData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sshare_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BackendData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BackendData) ProtoMessage() {}

func (x *BackendData) ProtoReflect() protoreflect.Message {
	mi := &file_sshare_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BackendData.ProtoReflect.Descriptor instead.
func (*BackendData) Descriptor() ([]byte, []int) {
	return file_sshare_proto_rawDescGZIP(), []int{2}
}

func (x *BackendData) GetStreamID() string {
	if x != nil {
		return x.StreamID
	}
	return ""
}

func (x *BackendData) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *BackendData) GetSshPublicKey() string {
	if x != nil {
		return x.SshPublicKey
	}
	return ""
}

func (x *BackendData) GetHTTPOptions() *HTTPOptions {
	if x != nil {
		return x.HTTPOptions
	}
	return nil
}

func (x *BackendData) GetOnlyTCP() bool {
	if x != nil {
		return x.OnlyTCP
	}
	return false
}

func (x *BackendData) GetConnection() *Connection {
	if x != nil {
		return x.Connection
	}
	return nil
}

type BackendReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error      string      `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
	Ready      bool        `protobuf:"varint,2,opt,name=ready,proto3" json:"ready,omitempty"`
	Connection *Connection `protobuf:"bytes,3,opt,name=connection,proto3" json:"connection,omitempty"`
	Deleted    bool        `protobuf:"varint,4,opt,name=deleted,proto3" json:"deleted,omitempty"`
}

func (x *BackendReply) Reset() {
	*x = BackendReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sshare_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BackendReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BackendReply) ProtoMessage() {}

func (x *BackendReply) ProtoReflect() protoreflect.Message {
	mi := &file_sshare_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BackendReply.ProtoReflect.Descriptor instead.
func (*BackendReply) Descriptor() ([]byte, []int) {
	return file_sshare_proto_rawDescGZIP(), []int{3}
}

func (x *BackendReply) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

func (x *BackendReply) GetReady() bool {
	if x != nil {
		return x.Ready
	}
	return false
}

func (x *BackendReply) GetConnection() *Connection {
	if x != nil {
		return x.Connection
	}
	return nil
}

func (x *BackendReply) GetDeleted() bool {
	if x != nil {
		return x.Deleted
	}
	return false
}

type HTTPOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CORSEnabled   bool `protobuf:"varint,1,opt,name=CORSEnabled,proto3" json:"CORSEnabled,omitempty"`
	HTTPSRedirect bool `protobuf:"varint,2,opt,name=HTTPSRedirect,proto3" json:"HTTPSRedirect,omitempty"`
}

func (x *HTTPOptions) Reset() {
	*x = HTTPOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sshare_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HTTPOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HTTPOptions) ProtoMessage() {}

func (x *HTTPOptions) ProtoReflect() protoreflect.Message {
	mi := &file_sshare_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HTTPOptions.ProtoReflect.Descriptor instead.
func (*HTTPOptions) Descriptor() ([]byte, []int) {
	return file_sshare_proto_rawDescGZIP(), []int{4}
}

func (x *HTTPOptions) GetCORSEnabled() bool {
	if x != nil {
		return x.CORSEnabled
	}
	return false
}

func (x *HTTPOptions) GetHTTPSRedirect() bool {
	if x != nil {
		return x.HTTPSRedirect
	}
	return false
}

type Connection struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SSHHost    string `protobuf:"bytes,1,opt,name=SSHHost,proto3" json:"SSHHost,omitempty"`
	SSHPort    int32  `protobuf:"varint,2,opt,name=SSHPort,proto3" json:"SSHPort,omitempty"`
	Domain     string `protobuf:"bytes,3,opt,name=Domain,proto3" json:"Domain,omitempty"`
	HTTPScheme bool   `protobuf:"varint,4,opt,name=HTTPScheme,proto3" json:"HTTPScheme,omitempty"`
	RemotePort int32  `protobuf:"varint,5,opt,name=RemotePort,proto3" json:"RemotePort,omitempty"`
	LocalPort  int32  `protobuf:"varint,6,opt,name=LocalPort,proto3" json:"LocalPort,omitempty"`
}

func (x *Connection) Reset() {
	*x = Connection{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sshare_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Connection) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Connection) ProtoMessage() {}

func (x *Connection) ProtoReflect() protoreflect.Message {
	mi := &file_sshare_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Connection.ProtoReflect.Descriptor instead.
func (*Connection) Descriptor() ([]byte, []int) {
	return file_sshare_proto_rawDescGZIP(), []int{5}
}

func (x *Connection) GetSSHHost() string {
	if x != nil {
		return x.SSHHost
	}
	return ""
}

func (x *Connection) GetSSHPort() int32 {
	if x != nil {
		return x.SSHPort
	}
	return 0
}

func (x *Connection) GetDomain() string {
	if x != nil {
		return x.Domain
	}
	return ""
}

func (x *Connection) GetHTTPScheme() bool {
	if x != nil {
		return x.HTTPScheme
	}
	return false
}

func (x *Connection) GetRemotePort() int32 {
	if x != nil {
		return x.RemotePort
	}
	return 0
}

func (x *Connection) GetLocalPort() int32 {
	if x != nil {
		return x.LocalPort
	}
	return 0
}

var File_sshare_proto protoreflect.FileDescriptor

var file_sshare_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x73, 0x73, 0x68, 0x61, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x22, 0x20, 0x0a, 0x0a, 0x54, 0x4c, 0x53, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x65, 0x6e, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x73, 0x65, 0x6e, 0x64, 0x22, 0x4b, 0x0a, 0x0b, 0x54, 0x4c,
	0x53, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x43, 0x41, 0x43,
	0x65, 0x72, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x43, 0x41, 0x43, 0x65, 0x72,
	0x74, 0x12, 0x24, 0x0a, 0x0d, 0x54, 0x4c, 0x53, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x50, 0x6f,
	0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0d, 0x54, 0x4c, 0x53, 0x53, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x50, 0x6f, 0x72, 0x74, 0x22, 0xea, 0x01, 0x0a, 0x0b, 0x42, 0x61, 0x63, 0x6b,
	0x65, 0x6e, 0x64, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x49, 0x44, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x73, 0x73, 0x68, 0x50, 0x75,
	0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x73,
	0x73, 0x68, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x12, 0x37, 0x0a, 0x0b, 0x48,
	0x54, 0x54, 0x50, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x48, 0x54, 0x54, 0x50,
	0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x0b, 0x48, 0x54, 0x54, 0x50, 0x4f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6f, 0x6e, 0x6c, 0x79, 0x54, 0x43, 0x50, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x6f, 0x6e, 0x6c, 0x79, 0x54, 0x43, 0x50, 0x12, 0x34,
	0x0a, 0x0a, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x43, 0x6f,
	0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x22, 0x8a, 0x01, 0x0a, 0x0c, 0x42, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64,
	0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x72,
	0x65, 0x61, 0x64, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x72, 0x65, 0x61, 0x64,
	0x79, 0x12, 0x34, 0x0a, 0x0a, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x63, 0x6f, 0x6e,
	0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x64, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x64, 0x22, 0x55, 0x0a, 0x0b, 0x48, 0x54, 0x54, 0x50, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x12, 0x20, 0x0a, 0x0b, 0x43, 0x4f, 0x52, 0x53, 0x45, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x43, 0x4f, 0x52, 0x53, 0x45, 0x6e, 0x61, 0x62, 0x6c,
	0x65, 0x64, 0x12, 0x24, 0x0a, 0x0d, 0x48, 0x54, 0x54, 0x50, 0x53, 0x52, 0x65, 0x64, 0x69, 0x72,
	0x65, 0x63, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0d, 0x48, 0x54, 0x54, 0x50, 0x53,
	0x52, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x22, 0xb6, 0x01, 0x0a, 0x0a, 0x43, 0x6f, 0x6e,
	0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x53, 0x53, 0x48, 0x48, 0x6f,
	0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x53, 0x53, 0x48, 0x48, 0x6f, 0x73,
	0x74, 0x12, 0x18, 0x0a, 0x07, 0x53, 0x53, 0x48, 0x50, 0x6f, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x07, 0x53, 0x53, 0x48, 0x50, 0x6f, 0x72, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x44,
	0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x44, 0x6f, 0x6d,
	0x61, 0x69, 0x6e, 0x12, 0x1e, 0x0a, 0x0a, 0x48, 0x54, 0x54, 0x50, 0x53, 0x63, 0x68, 0x65, 0x6d,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x48, 0x54, 0x54, 0x50, 0x53, 0x63, 0x68,
	0x65, 0x6d, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x52, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x50, 0x6f, 0x72,
	0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x52, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x50,
	0x6f, 0x72, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x50, 0x6f, 0x72, 0x74,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x50, 0x6f, 0x72,
	0x74, 0x32, 0x48, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x3e, 0x0a, 0x07, 0x42,
	0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x12, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x42, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x16, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x42, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64,
	0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x32, 0x44, 0x0a, 0x06, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x3a, 0x0a, 0x07, 0x42, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64,
	0x12, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x42, 0x61, 0x63, 0x6b,
	0x65, 0x6e, 0x64, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x42, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22,
	0x00, 0x32, 0x42, 0x0a, 0x03, 0x54, 0x4c, 0x53, 0x12, 0x3b, 0x0a, 0x0a, 0x43, 0x6f, 0x6e, 0x6e,
	0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x4c, 0x53, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x4c, 0x53, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x3b, 0x73, 0x73, 0x68, 0x61, 0x72,
	0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_sshare_proto_rawDescOnce sync.Once
	file_sshare_proto_rawDescData = file_sshare_proto_rawDesc
)

func file_sshare_proto_rawDescGZIP() []byte {
	file_sshare_proto_rawDescOnce.Do(func() {
		file_sshare_proto_rawDescData = protoimpl.X.CompressGZIP(file_sshare_proto_rawDescData)
	})
	return file_sshare_proto_rawDescData
}

var file_sshare_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_sshare_proto_goTypes = []interface{}{
	(*TLSRequest)(nil),   // 0: protobuf.TLSRequest
	(*TLSResponse)(nil),  // 1: protobuf.TLSResponse
	(*BackendData)(nil),  // 2: protobuf.BackendData
	(*BackendReply)(nil), // 3: protobuf.BackendReply
	(*HTTPOptions)(nil),  // 4: protobuf.HTTPOptions
	(*Connection)(nil),   // 5: protobuf.Connection
}
var file_sshare_proto_depIdxs = []int32{
	4, // 0: protobuf.BackendData.HTTPOptions:type_name -> protobuf.HTTPOptions
	5, // 1: protobuf.BackendData.connection:type_name -> protobuf.Connection
	5, // 2: protobuf.BackendReply.connection:type_name -> protobuf.Connection
	2, // 3: protobuf.Create.Backend:input_type -> protobuf.BackendData
	2, // 4: protobuf.Delete.Backend:input_type -> protobuf.BackendData
	0, // 5: protobuf.TLS.Connection:input_type -> protobuf.TLSRequest
	3, // 6: protobuf.Create.Backend:output_type -> protobuf.BackendReply
	3, // 7: protobuf.Delete.Backend:output_type -> protobuf.BackendReply
	1, // 8: protobuf.TLS.Connection:output_type -> protobuf.TLSResponse
	6, // [6:9] is the sub-list for method output_type
	3, // [3:6] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_sshare_proto_init() }
func file_sshare_proto_init() {
	if File_sshare_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_sshare_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TLSRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_sshare_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TLSResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_sshare_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BackendData); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_sshare_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BackendReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_sshare_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HTTPOptions); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_sshare_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Connection); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_sshare_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   3,
		},
		GoTypes:           file_sshare_proto_goTypes,
		DependencyIndexes: file_sshare_proto_depIdxs,
		MessageInfos:      file_sshare_proto_msgTypes,
	}.Build()
	File_sshare_proto = out.File
	file_sshare_proto_rawDesc = nil
	file_sshare_proto_goTypes = nil
	file_sshare_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// CreateClient is the client API for Create service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type CreateClient interface {
	Backend(ctx context.Context, opts ...grpc.CallOption) (Create_BackendClient, error)
}

type createClient struct {
	cc grpc.ClientConnInterface
}

func NewCreateClient(cc grpc.ClientConnInterface) CreateClient {
	return &createClient{cc}
}

func (c *createClient) Backend(ctx context.Context, opts ...grpc.CallOption) (Create_BackendClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Create_serviceDesc.Streams[0], "/protobuf.Create/Backend", opts...)
	if err != nil {
		return nil, err
	}
	x := &createBackendClient{stream}
	return x, nil
}

type Create_BackendClient interface {
	Send(*BackendData) error
	Recv() (*BackendReply, error)
	grpc.ClientStream
}

type createBackendClient struct {
	grpc.ClientStream
}

func (x *createBackendClient) Send(m *BackendData) error {
	return x.ClientStream.SendMsg(m)
}

func (x *createBackendClient) Recv() (*BackendReply, error) {
	m := new(BackendReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// CreateServer is the server API for Create service.
type CreateServer interface {
	Backend(Create_BackendServer) error
}

// UnimplementedCreateServer can be embedded to have forward compatible implementations.
type UnimplementedCreateServer struct {
}

func (*UnimplementedCreateServer) Backend(Create_BackendServer) error {
	return status.Errorf(codes.Unimplemented, "method Backend not implemented")
}

func RegisterCreateServer(s *grpc.Server, srv CreateServer) {
	s.RegisterService(&_Create_serviceDesc, srv)
}

func _Create_Backend_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(CreateServer).Backend(&createBackendServer{stream})
}

type Create_BackendServer interface {
	Send(*BackendReply) error
	Recv() (*BackendData, error)
	grpc.ServerStream
}

type createBackendServer struct {
	grpc.ServerStream
}

func (x *createBackendServer) Send(m *BackendReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *createBackendServer) Recv() (*BackendData, error) {
	m := new(BackendData)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Create_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protobuf.Create",
	HandlerType: (*CreateServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Backend",
			Handler:       _Create_Backend_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "sshare.proto",
}

// DeleteClient is the client API for Delete service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type DeleteClient interface {
	Backend(ctx context.Context, in *BackendData, opts ...grpc.CallOption) (*BackendReply, error)
}

type deleteClient struct {
	cc grpc.ClientConnInterface
}

func NewDeleteClient(cc grpc.ClientConnInterface) DeleteClient {
	return &deleteClient{cc}
}

func (c *deleteClient) Backend(ctx context.Context, in *BackendData, opts ...grpc.CallOption) (*BackendReply, error) {
	out := new(BackendReply)
	err := c.cc.Invoke(ctx, "/protobuf.Delete/Backend", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteServer is the server API for Delete service.
type DeleteServer interface {
	Backend(context.Context, *BackendData) (*BackendReply, error)
}

// UnimplementedDeleteServer can be embedded to have forward compatible implementations.
type UnimplementedDeleteServer struct {
}

func (*UnimplementedDeleteServer) Backend(context.Context, *BackendData) (*BackendReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Backend not implemented")
}

func RegisterDeleteServer(s *grpc.Server, srv DeleteServer) {
	s.RegisterService(&_Delete_serviceDesc, srv)
}

func _Delete_Backend_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BackendData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeleteServer).Backend(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.Delete/Backend",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeleteServer).Backend(ctx, req.(*BackendData))
	}
	return interceptor(ctx, in, info, handler)
}

var _Delete_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protobuf.Delete",
	HandlerType: (*DeleteServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Backend",
			Handler:    _Delete_Backend_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sshare.proto",
}

// TLSClient is the client API for TLS service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TLSClient interface {
	Connection(ctx context.Context, in *TLSRequest, opts ...grpc.CallOption) (*TLSResponse, error)
}

type tLSClient struct {
	cc grpc.ClientConnInterface
}

func NewTLSClient(cc grpc.ClientConnInterface) TLSClient {
	return &tLSClient{cc}
}

func (c *tLSClient) Connection(ctx context.Context, in *TLSRequest, opts ...grpc.CallOption) (*TLSResponse, error) {
	out := new(TLSResponse)
	err := c.cc.Invoke(ctx, "/protobuf.TLS/Connection", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TLSServer is the server API for TLS service.
type TLSServer interface {
	Connection(context.Context, *TLSRequest) (*TLSResponse, error)
}

// UnimplementedTLSServer can be embedded to have forward compatible implementations.
type UnimplementedTLSServer struct {
}

func (*UnimplementedTLSServer) Connection(context.Context, *TLSRequest) (*TLSResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Connection not implemented")
}

func RegisterTLSServer(s *grpc.Server, srv TLSServer) {
	s.RegisterService(&_TLS_serviceDesc, srv)
}

func _TLS_Connection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TLSRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TLSServer).Connection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.TLS/Connection",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TLSServer).Connection(ctx, req.(*TLSRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _TLS_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protobuf.TLS",
	HandlerType: (*TLSServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Connection",
			Handler:    _TLS_Connection_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sshare.proto",
}
