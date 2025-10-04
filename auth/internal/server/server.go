package server

import (
	"github.com/valkey-io/valkey-go"
	"github.com/jackc/pgx/v5/pgxpool"
	pb "github.com/Coosis/go-ecommerce/auth/internal/pb/v1/auth"
	internal "github.com/Coosis/go-ecommerce/auth/internal"
)

type Server struct{
	Pool *pgxpool.Pool
	VKclient valkey.Client
	AmqpClient *internal.Client

	pb.UnimplementedAuthServiceServer
}
