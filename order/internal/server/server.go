package server

import (
	pb "github.com/Coosis/go-ecommerce/order/internal/pb/v1/order"
	pbcatalog "github.com/Coosis/go-ecommerce/order/internal/pb/v1/catalog"
	pbinventory "github.com/Coosis/go-ecommerce/order/internal/pb/v1/inventory"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct{
	Pool *pgxpool.Pool

	CatalogClient pbcatalog.CatalogServiceClient
	InventoryClient pbinventory.InventoryServiceClient

	pb.UnimplementedOrderServiceServer
}
