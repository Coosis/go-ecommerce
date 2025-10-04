package server

import (
	"context"

	pb "github.com/Coosis/go-ecommerce/inventory/internal/pb/v1/inventory"
	sqlc "github.com/Coosis/go-ecommerce/inventory/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

func(s *Server) GetSkuCode(
	ctx context.Context,
	req *pb.GetSkuCodeRequest,
) (*pb.GetSkuCodeResponse, error) {
	var sku_id_uuid pgtype.UUID
	if err := sku_id_uuid.Scan(req.SkuId); err != nil {
		return nil, err
	}

	queries := sqlc.New(s.Pool)
	code, err := queries.GetSkuCode(ctx, sku_id_uuid)
	if err != nil {
		return nil, err
	}
	return &pb.GetSkuCodeResponse{
		Code: code,
	}, nil
}

func (s *Server) GetAllSkus(
	ctx context.Context,
	req *pb.GetAllSkusRequest,
) (*pb.GetAllSkusResponse, error) {
	queries := sqlc.New(s.Pool)
	var product_id_uuid pgtype.UUID
	if err := product_id_uuid.Scan(req.ProductId); err != nil {
		return nil, err
	}
	skus, err := queries.GetAllSkus(ctx, product_id_uuid)
	if err != nil {
		return nil, err
	}

	var pb_skus []string
	for _, sku := range skus {
		pb_skus = append(pb_skus, sku.ID.String())
	}

	return &pb.GetAllSkusResponse{
		SkuIds: pb_skus,
	}, nil
}
