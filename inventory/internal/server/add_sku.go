package server

import (
	"context"

	pb "github.com/Coosis/go-ecommerce/inventory/internal/pb/v1/inventory"
	sqlc "github.com/Coosis/go-ecommerce/inventory/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

func(s *Server) AddSku(
	ctx context.Context,
	req *pb.AddSkuRequest,
) (*pb.AddSkuResponse, error) {
	queries := sqlc.New(s.Pool)
	var product_id_uuid pgtype.UUID
	if err := product_id_uuid.Scan(req.ProductId); err != nil {
		return nil, err
	}
	resp, err := queries.CreateSku(ctx, sqlc.CreateSkuParams{
		ProductID: product_id_uuid,
		Code: req.Code,
	})
	if err != nil {
		return nil, err
	}
	return &pb.AddSkuResponse{ SkuId: resp.ID.String() }, nil
}
