package main

import (
	"context"
	"net"

	"google.golang.org/grpc"

	pb "github.com/Coosis/go-ecommerce/inventory/internal/pb/v1/inventory"
	internal "github.com/Coosis/go-ecommerce/inventory/internal"
	log "github.com/sirupsen/logrus"
	. "github.com/Coosis/go-ecommerce/inventory/internal/server"
)


func main() {
	ok, err := MigrateWithBackoff(context.Background(), internal.InventoryPostgresURL)
	if err != nil {
		panic("failed to run migrations: " + err.Error())
	}
	if !ok {
		panic("failed to run migrations")
	}

	pool, err := ConnectPoolWithBackoff(context.Background(), internal.InventoryPostgresURL)
	if err != nil {
		panic(err)
	}
	log.Info("Connected to Postgres!")

	s := &Server{
		Pool: pool,
	}

	lis, err := net.Listen("tcp", internal.InventoryPort)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterInventoryServiceServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
