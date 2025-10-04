package main

import (
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/Coosis/go-ecommerce/payment/internal/pb/v1/payment"
	pborder "github.com/Coosis/go-ecommerce/payment/internal/pb/v1/order"
	internal "github.com/Coosis/go-ecommerce/payment/internal"
	. "github.com/Coosis/go-ecommerce/payment/internal/server"
)

func main() {
	order_conn, err := grpc.NewClient(
		internal.OrderURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	defer order_conn.Close()

	s := &Server{
		OrderClient: pborder.NewOrderServiceClient(order_conn),
	}
	lis, err := net.Listen("tcp", internal.PaymentPort)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPaymentServiceServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
