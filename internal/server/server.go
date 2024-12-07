package server

import (
	"net"

	"google.golang.org/grpc"

	"fincraft-finance/api/finance"
)

// RunGRPCServer запускает gRPC сервер на указанном порту
func RunGRPCServer(port string, handler finance.FinanceServiceServer) error {
	grpcServer := grpc.NewServer()

	finance.RegisterFinanceServiceServer(grpcServer, handler)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	return grpcServer.Serve(lis)
}
