package main

import (
	"context"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/Coosis/go-ecommerce/cart/internal/pb/v1/cart"
	pbcatalog "github.com/Coosis/go-ecommerce/cart/internal/pb/v1/catalog"
	pborder "github.com/Coosis/go-ecommerce/cart/internal/pb/v1/order"
	internal "github.com/Coosis/go-ecommerce/cart/internal"
	log "github.com/sirupsen/logrus"
	. "github.com/Coosis/go-ecommerce/cart/internal/server"
)

func main() {
	ok, err := MigrateWithBackoff(context.Background(), internal.CartPostgresURL)
	if !ok || err != nil {
		panic("failed to run migrations: " + err.Error())
	}

	catalog_client, err := grpc.NewClient(
		internal.CatalogURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}

	order_url := os.Getenv(internal.OrderURL)
	if order_url == "" {
		panic(internal.OrderURL + " environment variable is not set")
	}
	order_client, err := grpc.NewClient(
		order_url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}

	pool, err := ConnectPoolWithBackoff(context.Background(), internal.CartPostgresURL)
	s := &Server{
		Pool: pool,
		CatalogClient: pbcatalog.NewCatalogServiceClient(catalog_client),
		OrderClient: pborder.NewOrderServiceClient(order_client),
	}
	log.Info("Connected to Postgres!")

	lis, err := net.Listen("tcp", internal.CartPort)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCartServiceServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
