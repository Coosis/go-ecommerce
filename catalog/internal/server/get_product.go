package server

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"

	pb "github.com/Coosis/go-ecommerce/catalog/internal/pb/v1/catalog"
	sqlc "github.com/Coosis/go-ecommerce/catalog/sqlc"
)

func(s *Server) GetProduct(
	ctx context.Context,
	req *pb.GetProductRequest,
) (*pb.GetProductResponse, error) {
	queries := sqlc.New(s.Pool)
	// casting to pgx uuid type
	var pg_product_id pgtype.UUID
	if err := pg_product_id.Scan(req.ProductId); err != nil {
		return nil, fmt.Errorf("failed to scan product ID: %w", err)
	}
	product, err := queries.GetProduct(ctx, pg_product_id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return &pb.GetProductResponse{
		ProductId: req.ProductId,
		Name: product.Name,
		Slug: product.Slug,
		Description: product.Description.String,
		PriceCents: product.PriceCents,
		PriceVersion: product.PriceVersion,
	}, nil
}
