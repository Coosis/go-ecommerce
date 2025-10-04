package server

import (
	"context"

	pb "github.com/Coosis/go-ecommerce/inventory/internal/pb/v1/inventory"
	sqlc "github.com/Coosis/go-ecommerce/inventory/sqlc"
)

func(s *Server) AddWarehouse(
	ctx context.Context,
	req *pb.AddWarehouseRequest,
) (*pb.AddWarehouseResponse, error) {
	queries := sqlc.New(s.Pool)
	resp, err := queries.CreateWarehouse(ctx, sqlc.CreateWarehouseParams{
		Code: req.Code,
		Name: req.Name,
	})
	if err != nil {
		return nil, err
	}
	return &pb.AddWarehouseResponse{ WarehouseId: resp.ID.String() }, nil
}
