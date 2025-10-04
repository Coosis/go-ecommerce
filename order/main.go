package main

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/Coosis/go-ecommerce/order/internal/pb/v1/order"
	pbcatalog "github.com/Coosis/go-ecommerce/order/internal/pb/v1/catalog"
	pbinventory "github.com/Coosis/go-ecommerce/order/internal/pb/v1/inventory"
	internal "github.com/Coosis/go-ecommerce/order/internal"
	log "github.com/sirupsen/logrus"

	. "github.com/Coosis/go-ecommerce/order/internal/server"
)

func main() {
	inventoryClient, err := grpc.NewClient(
		internal.InventoryURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}

	catalogClient, err := grpc.NewClient(
		internal.CatalogURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}

	ok, err := MigrateWithBackoff(context.Background(), internal.OrderPostgresURL)
	if err != nil {
		panic("failed to run migrations: " + err.Error())
	}
	if !ok {
		panic("failed to run migrations")
	}

	pool, err := ConnectPoolWithBackoff(context.Background(), internal.OrderPostgresURL)
	if err != nil {
		panic(err)
	}
	log.Info("Connected to Postgres!")

	s := &Server{
		Pool: pool,
		CatalogClient: pbcatalog.NewCatalogServiceClient(catalogClient),
		InventoryClient: pbinventory.NewInventoryServiceClient(inventoryClient),
	}

	lis, err := net.Listen("tcp", internal.OrderPort)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
