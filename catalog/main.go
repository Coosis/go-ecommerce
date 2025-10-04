package main

import (
	"context"
	"net"

	"google.golang.org/grpc"

	pb "github.com/Coosis/go-ecommerce/catalog/internal/pb/v1/catalog"
	internal "github.com/Coosis/go-ecommerce/catalog/internal"
	log "github.com/sirupsen/logrus"
	. "github.com/Coosis/go-ecommerce/catalog/internal/server"
)

func main() {
	ok, err := MigrateWithBackoff(context.Background(), internal.CATALOG_POSTGRES_URL)
	if !ok || err != nil {
		panic("failed to run migrations: " + err.Error())
	}
	pool, err := ConnectPoolWithBackoff(context.Background(), internal.CATALOG_POSTGRES_URL)
	if err != nil {
		panic(err)
	}
	defer pool.Close()
	log.Info("Connected to Postgres!")

	s := &Server{Pool: pool}

	grpcServer := grpc.NewServer()
	pb.RegisterCatalogServiceServer(grpcServer, s)

	lis, err := net.Listen("tcp", internal.CATALOG_PORT)
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}

}
