package server

import (
	"context"

	pb "github.com/Coosis/go-ecommerce/catalog/internal/pb/v1/catalog"
	sqlc "github.com/Coosis/go-ecommerce/catalog/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

func(s *Server) UpdateProduct(
	ctx context.Context,
	req *pb.UpdateProductRequest,
) (*pb.UpdateProductResponse, error) {
	queries := sqlc.New(s.Pool)
	var pg_product_id pgtype.UUID
	if err := pg_product_id.Scan(req.GetDescription()); err != nil {
		return nil, err
	}
	p, err := queries.UpdateProduct(ctx, sqlc.UpdateProductParams{
		ID: pg_product_id,
		Name: req.GetName(),
		Slug: req.GetSlug(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.UpdateProductResponse{
		ProductId: p.ID.String(),
	}, nil
}
