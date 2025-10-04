package server

import (
	pb "github.com/Coosis/go-ecommerce/inventory/internal/pb/v1/inventory"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct{
	Pool *pgxpool.Pool

	pb.UnimplementedInventoryServiceServer
}
