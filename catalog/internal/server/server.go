package server

import (
	pb "github.com/Coosis/go-ecommerce/catalog/internal/pb/v1/catalog"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	Pool *pgxpool.Pool

	pb.UnimplementedCatalogServiceServer
}
