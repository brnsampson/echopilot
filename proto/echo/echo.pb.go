// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/echo/echo.proto

package echo

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// The request message containing a user defined string.
type EchoRequest struct {
	Content              string   `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EchoRequest) Reset()         { *m = EchoRequest{} }
func (m *EchoRequest) String() string { return proto.CompactTextString(m) }
func (*EchoRequest) ProtoMessage()    {}
func (*EchoRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6e0b19995c129c8b, []int{0}
}

func (m *EchoRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EchoRequest.Unmarshal(m, b)
}
func (m *EchoRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EchoRequest.Marshal(b, m, deterministic)
}
func (m *EchoRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EchoRequest.Merge(m, src)
}
func (m *EchoRequest) XXX_Size() int {
	return xxx_messageInfo_EchoRequest.Size(m)
}
func (m *EchoRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_EchoRequest.DiscardUnknown(m)
}

var xxx_messageInfo_EchoRequest proto.InternalMessageInfo

func (m *EchoRequest) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

// The response message containing the content of the original requests
type EchoReply struct {
	Content              string   `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EchoReply) Reset()         { *m = EchoReply{} }
func (m *EchoReply) String() string { return proto.CompactTextString(m) }
func (*EchoReply) ProtoMessage()    {}
func (*EchoReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_6e0b19995c129c8b, []int{1}
}

func (m *EchoReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EchoReply.Unmarshal(m, b)
}
func (m *EchoReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EchoReply.Marshal(b, m, deterministic)
}
func (m *EchoReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EchoReply.Merge(m, src)
}
func (m *EchoReply) XXX_Size() int {
	return xxx_messageInfo_EchoReply.Size(m)
}
func (m *EchoReply) XXX_DiscardUnknown() {
	xxx_messageInfo_EchoReply.DiscardUnknown(m)
}

var xxx_messageInfo_EchoReply proto.InternalMessageInfo

func (m *EchoReply) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

func init() {
	proto.RegisterType((*EchoRequest)(nil), "echo.EchoRequest")
	proto.RegisterType((*EchoReply)(nil), "echo.EchoReply")
}

func init() { proto.RegisterFile("proto/echo/echo.proto", fileDescriptor_6e0b19995c129c8b) }

var fileDescriptor_6e0b19995c129c8b = []byte{
	// 199 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2d, 0x28, 0xca, 0x2f,
	0xc9, 0xd7, 0x4f, 0x4d, 0xce, 0x80, 0x10, 0x7a, 0x60, 0xbe, 0x10, 0x0b, 0x88, 0x2d, 0x25, 0x93,
	0x9e, 0x9f, 0x9f, 0x9e, 0x93, 0xaa, 0x9f, 0x58, 0x90, 0xa9, 0x9f, 0x98, 0x97, 0x97, 0x5f, 0x92,
	0x58, 0x92, 0x99, 0x9f, 0x57, 0x0c, 0x51, 0xa3, 0xa4, 0xce, 0xc5, 0xed, 0x9a, 0x9c, 0x91, 0x1f,
	0x94, 0x5a, 0x58, 0x9a, 0x5a, 0x5c, 0x22, 0x24, 0xc1, 0xc5, 0x9e, 0x9c, 0x9f, 0x57, 0x92, 0x9a,
	0x57, 0x22, 0xc1, 0xa8, 0xc0, 0xa8, 0xc1, 0x19, 0x04, 0xe3, 0x2a, 0xa9, 0x72, 0x71, 0x42, 0x14,
	0x16, 0xe4, 0x54, 0xe2, 0x56, 0x66, 0xe4, 0xcd, 0xc5, 0x02, 0x52, 0x26, 0xe4, 0x0c, 0xa5, 0x05,
	0xf5, 0xc0, 0x0e, 0x42, 0xb2, 0x43, 0x8a, 0x1f, 0x59, 0xa8, 0x20, 0xa7, 0x52, 0x49, 0xbc, 0xe9,
	0xf2, 0x93, 0xc9, 0x4c, 0x82, 0x4a, 0x3c, 0xfa, 0x65, 0x86, 0x60, 0x67, 0x82, 0xe4, 0xad, 0x18,
	0xb5, 0x9c, 0x74, 0xa2, 0xb4, 0xd2, 0x33, 0x4b, 0x32, 0x4a, 0x93, 0xf4, 0x92, 0xf3, 0x73, 0xf5,
	0x93, 0x8a, 0xf2, 0x8a, 0x13, 0x73, 0x0b, 0x8a, 0xf3, 0xf3, 0xc0, 0x0a, 0x0a, 0x32, 0x73, 0xf2,
	0x4b, 0xf4, 0x11, 0x3e, 0x4f, 0x62, 0x03, 0xb3, 0x8d, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x46,
	0x4a, 0xb4, 0xe6, 0x0e, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// EchoClient is the client API for Echo service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type EchoClient interface {
	// Echos a string back to sender
	Echo(ctx context.Context, in *EchoRequest, opts ...grpc.CallOption) (*EchoReply, error)
}

type echoClient struct {
	cc *grpc.ClientConn
}

func NewEchoClient(cc *grpc.ClientConn) EchoClient {
	return &echoClient{cc}
}

func (c *echoClient) Echo(ctx context.Context, in *EchoRequest, opts ...grpc.CallOption) (*EchoReply, error) {
	out := new(EchoReply)
	err := c.cc.Invoke(ctx, "/echo.Echo/Echo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EchoServer is the server API for Echo service.
type EchoServer interface {
	// Echos a string back to sender
	Echo(context.Context, *EchoRequest) (*EchoReply, error)
}

// UnimplementedEchoServer can be embedded to have forward compatible implementations.
type UnimplementedEchoServer struct {
}

func (*UnimplementedEchoServer) Echo(ctx context.Context, req *EchoRequest) (*EchoReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Echo not implemented")
}

func RegisterEchoServer(s *grpc.Server, srv EchoServer) {
	s.RegisterService(&_Echo_serviceDesc, srv)
}

func _Echo_Echo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EchoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EchoServer).Echo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/echo.Echo/Echo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EchoServer).Echo(ctx, req.(*EchoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Echo_serviceDesc = grpc.ServiceDesc{
	ServiceName: "echo.Echo",
	HandlerType: (*EchoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Echo",
			Handler:    _Echo_Echo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/echo/echo.proto",
}
