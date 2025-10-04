package main

import (
	"context"
	"net"

	"google.golang.org/grpc"

	pb "github.com/Coosis/go-ecommerce/auth/internal/pb/v1/auth"
	internal "github.com/Coosis/go-ecommerce/auth/internal"
	log "github.com/sirupsen/logrus"

	. "github.com/Coosis/go-ecommerce/auth/internal/server"
)

func main() {
	// valkey client setup
	client, err := ConnectValkeyWithBackoff(context.Background(), internal.ValkeyURL)
	// client, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{valkey_url}})
	if err != nil {
		panic(err)
	}
	defer client.Close()
	log.Info("Connected to Valkey!")

	// postgres connection setup
	ok, err := MigrateWithBackoff(context.Background(), internal.AuthPostgresUrl)
	if !ok || err != nil {
		panic("failed to run migrations: " + err.Error())
	}
	pool, err := ConnectPoolWithBackoff(
		context.Background(),
		internal.AuthPostgresUrl,
	)
	if err != nil {
		panic(err)
	}
	defer pool.Close()
	log.Info("Connected to Postgres!")

	amqp_client, err := internal.NewClient(context.Background(), internal.AmqpUrl)
	if err != nil {
		panic("failed to connect to RabbitMQ, " + err.Error())
	}
	log.Info("Connected to RabbitMQ!")

	ser := &Server{
		Pool: pool,
		VKclient: client,
		AmqpClient: amqp_client,
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, ser)

	lis, err := net.Listen("tcp", internal.AuthPort)
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
