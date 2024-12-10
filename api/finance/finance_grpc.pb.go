// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.3
// source: finance/finance.proto

package finance

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	FinanceService_AddIncome_FullMethodName          = "/finance.FinanceService/AddIncome"
	FinanceService_GetIncomeForPeriod_FullMethodName = "/finance.FinanceService/GetIncomeForPeriod"
)

// FinanceServiceClient is the client API for FinanceService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FinanceServiceClient interface {
	AddIncome(ctx context.Context, in *AddIncomeRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetIncomeForPeriod(ctx context.Context, in *GetIncomeForPeriodRequest, opts ...grpc.CallOption) (*GetIncomeForPeriodResponse, error)
}

type financeServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFinanceServiceClient(cc grpc.ClientConnInterface) FinanceServiceClient {
	return &financeServiceClient{cc}
}

func (c *financeServiceClient) AddIncome(ctx context.Context, in *AddIncomeRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, FinanceService_AddIncome_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *financeServiceClient) GetIncomeForPeriod(ctx context.Context, in *GetIncomeForPeriodRequest, opts ...grpc.CallOption) (*GetIncomeForPeriodResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetIncomeForPeriodResponse)
	err := c.cc.Invoke(ctx, FinanceService_GetIncomeForPeriod_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FinanceServiceServer is the server API for FinanceService service.
// All implementations must embed UnimplementedFinanceServiceServer
// for forward compatibility.
type FinanceServiceServer interface {
	AddIncome(context.Context, *AddIncomeRequest) (*emptypb.Empty, error)
	GetIncomeForPeriod(context.Context, *GetIncomeForPeriodRequest) (*GetIncomeForPeriodResponse, error)
	mustEmbedUnimplementedFinanceServiceServer()
}

// UnimplementedFinanceServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedFinanceServiceServer struct{}

func (UnimplementedFinanceServiceServer) AddIncome(context.Context, *AddIncomeRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddIncome not implemented")
}
func (UnimplementedFinanceServiceServer) GetIncomeForPeriod(context.Context, *GetIncomeForPeriodRequest) (*GetIncomeForPeriodResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetIncomeForPeriod not implemented")
}
func (UnimplementedFinanceServiceServer) mustEmbedUnimplementedFinanceServiceServer() {}
func (UnimplementedFinanceServiceServer) testEmbeddedByValue()                        {}

// UnsafeFinanceServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FinanceServiceServer will
// result in compilation errors.
type UnsafeFinanceServiceServer interface {
	mustEmbedUnimplementedFinanceServiceServer()
}

func RegisterFinanceServiceServer(s grpc.ServiceRegistrar, srv FinanceServiceServer) {
	// If the following call pancis, it indicates UnimplementedFinanceServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&FinanceService_ServiceDesc, srv)
}

func _FinanceService_AddIncome_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddIncomeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FinanceServiceServer).AddIncome(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FinanceService_AddIncome_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceServiceServer).AddIncome(ctx, req.(*AddIncomeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FinanceService_GetIncomeForPeriod_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetIncomeForPeriodRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FinanceServiceServer).GetIncomeForPeriod(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FinanceService_GetIncomeForPeriod_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceServiceServer).GetIncomeForPeriod(ctx, req.(*GetIncomeForPeriodRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FinanceService_ServiceDesc is the grpc.ServiceDesc for FinanceService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FinanceService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "finance.FinanceService",
	HandlerType: (*FinanceServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddIncome",
			Handler:    _FinanceService_AddIncome_Handler,
		},
		{
			MethodName: "GetIncomeForPeriod",
			Handler:    _FinanceService_GetIncomeForPeriod_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "finance/finance.proto",
}
