// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package neutrinorpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// NeutrinoKitClient is the client API for NeutrinoKit service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NeutrinoKitClient interface {
	// lncli: `neutrino status`
	// Status returns the status of the light client neutrino instance,
	// along with height and hash of the best block, and a list of connected
	// peers.
	Status(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusResponse, error)
	// lncli: `neutrino addpeer`
	// AddPeer adds a new peer that has already been connected to the server.
	AddPeer(ctx context.Context, in *AddPeerRequest, opts ...grpc.CallOption) (*AddPeerResponse, error)
	// lncli: `neutrino disconnectpeer`
	// DisconnectPeer disconnects a peer by target address. Both outbound and
	// inbound nodes will be searched for the target node. An error message will
	// be returned if the peer was not found.
	DisconnectPeer(ctx context.Context, in *DisconnectPeerRequest, opts ...grpc.CallOption) (*DisconnectPeerResponse, error)
	// lncli: `neutrino isbanned`
	// IsBanned returns true if the peer is banned, otherwise false.
	IsBanned(ctx context.Context, in *IsBannedRequest, opts ...grpc.CallOption) (*IsBannedResponse, error)
	// lncli: `neutrino getblockheader`
	// GetBlockHeader returns a block header with a particular block hash.
	GetBlockHeader(ctx context.Context, in *GetBlockHeaderRequest, opts ...grpc.CallOption) (*GetBlockHeaderResponse, error)
	// GetBlock returns a block with a particular block hash.
	GetBlock(ctx context.Context, in *GetBlockRequest, opts ...grpc.CallOption) (*GetBlockResponse, error)
	// lncli: `neutrino getcfilter`
	// GetCFilter returns a compact filter from a block.
	GetCFilter(ctx context.Context, in *GetCFilterRequest, opts ...grpc.CallOption) (*GetCFilterResponse, error)
	// Deprecated: Do not use.
	//
	// Deprecated, use chainrpc.GetBlockHash instead.
	// GetBlockHash returns the header hash of a block at a given height.
	GetBlockHash(ctx context.Context, in *GetBlockHashRequest, opts ...grpc.CallOption) (*GetBlockHashResponse, error)
	// UnbanPeer unbans a previously banned peer and connects to it.
	UnbanPeer(ctx context.Context, in *UnbanPeerRequest, opts ...grpc.CallOption) (*UnbanPeerResponse, error)
}

type neutrinoKitClient struct {
	cc grpc.ClientConnInterface
}

func NewNeutrinoKitClient(cc grpc.ClientConnInterface) NeutrinoKitClient {
	return &neutrinoKitClient{cc}
}

func (c *neutrinoKitClient) Status(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusResponse, error) {
	out := new(StatusResponse)
	err := c.cc.Invoke(ctx, "/neutrinorpc.NeutrinoKit/Status", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *neutrinoKitClient) AddPeer(ctx context.Context, in *AddPeerRequest, opts ...grpc.CallOption) (*AddPeerResponse, error) {
	out := new(AddPeerResponse)
	err := c.cc.Invoke(ctx, "/neutrinorpc.NeutrinoKit/AddPeer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *neutrinoKitClient) DisconnectPeer(ctx context.Context, in *DisconnectPeerRequest, opts ...grpc.CallOption) (*DisconnectPeerResponse, error) {
	out := new(DisconnectPeerResponse)
	err := c.cc.Invoke(ctx, "/neutrinorpc.NeutrinoKit/DisconnectPeer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *neutrinoKitClient) IsBanned(ctx context.Context, in *IsBannedRequest, opts ...grpc.CallOption) (*IsBannedResponse, error) {
	out := new(IsBannedResponse)
	err := c.cc.Invoke(ctx, "/neutrinorpc.NeutrinoKit/IsBanned", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *neutrinoKitClient) GetBlockHeader(ctx context.Context, in *GetBlockHeaderRequest, opts ...grpc.CallOption) (*GetBlockHeaderResponse, error) {
	out := new(GetBlockHeaderResponse)
	err := c.cc.Invoke(ctx, "/neutrinorpc.NeutrinoKit/GetBlockHeader", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *neutrinoKitClient) GetBlock(ctx context.Context, in *GetBlockRequest, opts ...grpc.CallOption) (*GetBlockResponse, error) {
	out := new(GetBlockResponse)
	err := c.cc.Invoke(ctx, "/neutrinorpc.NeutrinoKit/GetBlock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *neutrinoKitClient) GetCFilter(ctx context.Context, in *GetCFilterRequest, opts ...grpc.CallOption) (*GetCFilterResponse, error) {
	out := new(GetCFilterResponse)
	err := c.cc.Invoke(ctx, "/neutrinorpc.NeutrinoKit/GetCFilter", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Deprecated: Do not use.
func (c *neutrinoKitClient) GetBlockHash(ctx context.Context, in *GetBlockHashRequest, opts ...grpc.CallOption) (*GetBlockHashResponse, error) {
	out := new(GetBlockHashResponse)
	err := c.cc.Invoke(ctx, "/neutrinorpc.NeutrinoKit/GetBlockHash", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *neutrinoKitClient) UnbanPeer(ctx context.Context, in *UnbanPeerRequest, opts ...grpc.CallOption) (*UnbanPeerResponse, error) {
	out := new(UnbanPeerResponse)
	err := c.cc.Invoke(ctx, "/neutrinorpc.NeutrinoKit/UnbanPeer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NeutrinoKitServer is the server API for NeutrinoKit service.
// All implementations must embed UnimplementedNeutrinoKitServer
// for forward compatibility
type NeutrinoKitServer interface {
	// lncli: `neutrino status`
	// Status returns the status of the light client neutrino instance,
	// along with height and hash of the best block, and a list of connected
	// peers.
	Status(context.Context, *StatusRequest) (*StatusResponse, error)
	// lncli: `neutrino addpeer`
	// AddPeer adds a new peer that has already been connected to the server.
	AddPeer(context.Context, *AddPeerRequest) (*AddPeerResponse, error)
	// lncli: `neutrino disconnectpeer`
	// DisconnectPeer disconnects a peer by target address. Both outbound and
	// inbound nodes will be searched for the target node. An error message will
	// be returned if the peer was not found.
	DisconnectPeer(context.Context, *DisconnectPeerRequest) (*DisconnectPeerResponse, error)
	// lncli: `neutrino isbanned`
	// IsBanned returns true if the peer is banned, otherwise false.
	IsBanned(context.Context, *IsBannedRequest) (*IsBannedResponse, error)
	// lncli: `neutrino getblockheader`
	// GetBlockHeader returns a block header with a particular block hash.
	GetBlockHeader(context.Context, *GetBlockHeaderRequest) (*GetBlockHeaderResponse, error)
	// GetBlock returns a block with a particular block hash.
	GetBlock(context.Context, *GetBlockRequest) (*GetBlockResponse, error)
	// lncli: `neutrino getcfilter`
	// GetCFilter returns a compact filter from a block.
	GetCFilter(context.Context, *GetCFilterRequest) (*GetCFilterResponse, error)
	// Deprecated: Do not use.
	//
	// Deprecated, use chainrpc.GetBlockHash instead.
	// GetBlockHash returns the header hash of a block at a given height.
	GetBlockHash(context.Context, *GetBlockHashRequest) (*GetBlockHashResponse, error)
	// UnbanPeer unbans a previously banned peer and connects to it.
	UnbanPeer(context.Context, *UnbanPeerRequest) (*UnbanPeerResponse, error)
	mustEmbedUnimplementedNeutrinoKitServer()
}

// UnimplementedNeutrinoKitServer must be embedded to have forward compatible implementations.
type UnimplementedNeutrinoKitServer struct {
}

func (UnimplementedNeutrinoKitServer) Status(context.Context, *StatusRequest) (*StatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}
func (UnimplementedNeutrinoKitServer) AddPeer(context.Context, *AddPeerRequest) (*AddPeerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddPeer not implemented")
}
func (UnimplementedNeutrinoKitServer) DisconnectPeer(context.Context, *DisconnectPeerRequest) (*DisconnectPeerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DisconnectPeer not implemented")
}
func (UnimplementedNeutrinoKitServer) IsBanned(context.Context, *IsBannedRequest) (*IsBannedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsBanned not implemented")
}
func (UnimplementedNeutrinoKitServer) GetBlockHeader(context.Context, *GetBlockHeaderRequest) (*GetBlockHeaderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBlockHeader not implemented")
}
func (UnimplementedNeutrinoKitServer) GetBlock(context.Context, *GetBlockRequest) (*GetBlockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBlock not implemented")
}
func (UnimplementedNeutrinoKitServer) GetCFilter(context.Context, *GetCFilterRequest) (*GetCFilterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCFilter not implemented")
}
func (UnimplementedNeutrinoKitServer) GetBlockHash(context.Context, *GetBlockHashRequest) (*GetBlockHashResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBlockHash not implemented")
}
func (UnimplementedNeutrinoKitServer) UnbanPeer(context.Context, *UnbanPeerRequest) (*UnbanPeerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnbanPeer not implemented")
}
func (UnimplementedNeutrinoKitServer) mustEmbedUnimplementedNeutrinoKitServer() {}

// UnsafeNeutrinoKitServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NeutrinoKitServer will
// result in compilation errors.
type UnsafeNeutrinoKitServer interface {
	mustEmbedUnimplementedNeutrinoKitServer()
}

func RegisterNeutrinoKitServer(s grpc.ServiceRegistrar, srv NeutrinoKitServer) {
	s.RegisterService(&NeutrinoKit_ServiceDesc, srv)
}

func _NeutrinoKit_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NeutrinoKitServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/neutrinorpc.NeutrinoKit/Status",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NeutrinoKitServer).Status(ctx, req.(*StatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NeutrinoKit_AddPeer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddPeerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NeutrinoKitServer).AddPeer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/neutrinorpc.NeutrinoKit/AddPeer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NeutrinoKitServer).AddPeer(ctx, req.(*AddPeerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NeutrinoKit_DisconnectPeer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DisconnectPeerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NeutrinoKitServer).DisconnectPeer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/neutrinorpc.NeutrinoKit/DisconnectPeer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NeutrinoKitServer).DisconnectPeer(ctx, req.(*DisconnectPeerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NeutrinoKit_IsBanned_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IsBannedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NeutrinoKitServer).IsBanned(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/neutrinorpc.NeutrinoKit/IsBanned",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NeutrinoKitServer).IsBanned(ctx, req.(*IsBannedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NeutrinoKit_GetBlockHeader_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBlockHeaderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NeutrinoKitServer).GetBlockHeader(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/neutrinorpc.NeutrinoKit/GetBlockHeader",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NeutrinoKitServer).GetBlockHeader(ctx, req.(*GetBlockHeaderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NeutrinoKit_GetBlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBlockRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NeutrinoKitServer).GetBlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/neutrinorpc.NeutrinoKit/GetBlock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NeutrinoKitServer).GetBlock(ctx, req.(*GetBlockRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NeutrinoKit_GetCFilter_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCFilterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NeutrinoKitServer).GetCFilter(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/neutrinorpc.NeutrinoKit/GetCFilter",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NeutrinoKitServer).GetCFilter(ctx, req.(*GetCFilterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NeutrinoKit_GetBlockHash_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBlockHashRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NeutrinoKitServer).GetBlockHash(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/neutrinorpc.NeutrinoKit/GetBlockHash",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NeutrinoKitServer).GetBlockHash(ctx, req.(*GetBlockHashRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NeutrinoKit_UnbanPeer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnbanPeerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NeutrinoKitServer).UnbanPeer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/neutrinorpc.NeutrinoKit/UnbanPeer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NeutrinoKitServer).UnbanPeer(ctx, req.(*UnbanPeerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// NeutrinoKit_ServiceDesc is the grpc.ServiceDesc for NeutrinoKit service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NeutrinoKit_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "neutrinorpc.NeutrinoKit",
	HandlerType: (*NeutrinoKitServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Status",
			Handler:    _NeutrinoKit_Status_Handler,
		},
		{
			MethodName: "AddPeer",
			Handler:    _NeutrinoKit_AddPeer_Handler,
		},
		{
			MethodName: "DisconnectPeer",
			Handler:    _NeutrinoKit_DisconnectPeer_Handler,
		},
		{
			MethodName: "IsBanned",
			Handler:    _NeutrinoKit_IsBanned_Handler,
		},
		{
			MethodName: "GetBlockHeader",
			Handler:    _NeutrinoKit_GetBlockHeader_Handler,
		},
		{
			MethodName: "GetBlock",
			Handler:    _NeutrinoKit_GetBlock_Handler,
		},
		{
			MethodName: "GetCFilter",
			Handler:    _NeutrinoKit_GetCFilter_Handler,
		},
		{
			MethodName: "GetBlockHash",
			Handler:    _NeutrinoKit_GetBlockHash_Handler,
		},
		{
			MethodName: "UnbanPeer",
			Handler:    _NeutrinoKit_UnbanPeer_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "neutrinorpc/neutrino.proto",
}
