package server

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	pb "github.com/Coosis/go-ecommerce/catalog/internal/pb/v1/catalog"
	sqlc "github.com/Coosis/go-ecommerce/catalog/sqlc"
)

func (s *Server) CreateProduct(
	ctx context.Context,
	req *pb.CreateProductRequest,
) (*pb.CreateProductResponse, error) {
	queries := sqlc.New(s.Pool)
	var pg_description pgtype.Text
	if err := pg_description.Scan(req.GetDescription()); err != nil {
		return nil, err
	}
	p, err := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Name: req.GetName(),
		Slug: req.GetSlug(),
		Description: pg_description,
		PriceCents: req.GetPriceCents(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreateProductResponse{
		ProductId: p.ID.String(),
	}, nil
}
