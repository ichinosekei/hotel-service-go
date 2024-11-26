// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.3
// source: hotelier-service/api/protobufs/hotelier_service.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	HotelierService_GetRoomPrice_FullMethodName = "/hotelier.HotelierService/GetRoomPrice"
)

// HotelierServiceClient is the client API for HotelierService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HotelierServiceClient interface {
	GetRoomPrice(ctx context.Context, in *RoomRequest, opts ...grpc.CallOption) (*RoomResponse, error)
}

type hotelierServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewHotelierServiceClient(cc grpc.ClientConnInterface) HotelierServiceClient {
	return &hotelierServiceClient{cc}
}

func (c *hotelierServiceClient) GetRoomPrice(ctx context.Context, in *RoomRequest, opts ...grpc.CallOption) (*RoomResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RoomResponse)
	err := c.cc.Invoke(ctx, HotelierService_GetRoomPrice_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HotelierServiceServer is the server API for HotelierService service.
// All implementations must embed UnimplementedHotelierServiceServer
// for forward compatibility.
type HotelierServiceServer interface {
	GetRoomPrice(context.Context, *RoomRequest) (*RoomResponse, error)
	mustEmbedUnimplementedHotelierServiceServer()
}

// UnimplementedHotelierServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedHotelierServiceServer struct{}

func (UnimplementedHotelierServiceServer) GetRoomPrice(context.Context, *RoomRequest) (*RoomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRoomPrice not implemented")
}
func (UnimplementedHotelierServiceServer) mustEmbedUnimplementedHotelierServiceServer() {}
func (UnimplementedHotelierServiceServer) testEmbeddedByValue()                         {}

// UnsafeHotelierServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HotelierServiceServer will
// result in compilation errors.
type UnsafeHotelierServiceServer interface {
	mustEmbedUnimplementedHotelierServiceServer()
}

func RegisterHotelierServiceServer(s grpc.ServiceRegistrar, srv HotelierServiceServer) {
	// If the following call pancis, it indicates UnimplementedHotelierServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&HotelierService_ServiceDesc, srv)
}

func _HotelierService_GetRoomPrice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RoomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HotelierServiceServer).GetRoomPrice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: HotelierService_GetRoomPrice_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HotelierServiceServer).GetRoomPrice(ctx, req.(*RoomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// HotelierService_ServiceDesc is the grpc.ServiceDesc for HotelierService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var HotelierService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "hotelier.HotelierService",
	HandlerType: (*HotelierServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetRoomPrice",
			Handler:    _HotelierService_GetRoomPrice_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "hotelier-service/api/protobufs/hotelier_service.proto",
}