package server

import (
	"fincraft-finance/api/finance"
	"google.golang.org/grpc"
	"net"
)

func RunGRPCServer(port string, handler finance.FinanceServiceServer) error {
	grpcServer := grpc.NewServer()

	finance.RegisterFinanceServiceServer(grpcServer, handler)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	return grpcServer.Serve(lis)
}
