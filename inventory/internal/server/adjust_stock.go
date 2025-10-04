package server

import (
	"fmt"
	"context"

	pb "github.com/Coosis/go-ecommerce/inventory/internal/pb/v1/inventory"
	sqlc "github.com/Coosis/go-ecommerce/inventory/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

func(s *Server) AdjustStock(
	ctx context.Context,
	req *pb.AdjustStockRequest,
) (*pb.AdjustStockResponse, error) {
	queries := sqlc.New(s.Pool)
	var sku_uuid pgtype.UUID
	if err := sku_uuid.Scan(req.GetSkuId()); err != nil {
		return nil, fmt.Errorf("failed to scan sku ID: %w", err)
	}
	var warehouse_uuid pgtype.UUID
	if err := warehouse_uuid.Scan(req.GetWarehouseId()); err != nil {
		return nil, fmt.Errorf("failed to scan warehouse ID: %w", err)
	}
	queries.AdjustStockLevel(ctx, sqlc.AdjustStockLevelParams{
		SkuID: sku_uuid,
		WarehouseID: warehouse_uuid,
		Delta: req.Delta,
		Reason: req.Reason,
		CreatedBy: req.CreatedBy,
	});
	return &pb.AdjustStockResponse{}, nil
}
