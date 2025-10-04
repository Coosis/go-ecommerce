package server

import (
	"context"

	pb "github.com/Coosis/go-ecommerce/catalog/internal/pb/v1/catalog"
	sqlc "github.com/Coosis/go-ecommerce/catalog/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

func(s *Server) ListProducts(
	ctx context.Context,
	req *pb.ListProductsRequest,
) (*pb.ListProductsResponse, error) {
	queries := sqlc.New(s.Pool)
	var pg_slug pgtype.Text
	if req.CategorySlug == nil {
		pg_slug = pgtype.Text{Valid: false}
	} else {
		pg_slug = pgtype.Text{Valid: true, String: req.GetCategorySlug()}
	}
	var pg_min_price pgtype.Int4
	if req.MinPriceCents == nil {
		pg_min_price = pgtype.Int4{Valid: false}
	} else {
		pg_min_price = pgtype.Int4{Valid: true, Int32: req.GetMinPriceCents()}
	}

	var pg_max_price pgtype.Int4
	if req.MaxPriceCents == nil {
		pg_max_price = pgtype.Int4{Valid: false}
	} else {
		pg_max_price = pgtype.Int4{Valid: true, Int32: req.GetMaxPriceCents()}
	}
	ps, err := queries.ListProducts(ctx, sqlc.ListProductsParams{
		MaxPriceCents: pg_max_price,
		MinPriceCents: pg_min_price,
		Slug: pg_slug,
		PageNumber: req.GetPage(),
		PageSize: req.GetPerpage(),
	})
	if err != nil {
		return nil, err
	}
	var products []*pb.GetProductResponse
	for _, p := range ps {
		products = append(products, &pb.GetProductResponse{
			ProductId: p.ID.String(),
			Name: p.Name,
			Slug: p.Slug,
			Description: p.Description.String,
			PriceCents: p.PriceCents,
		})
	}
	return &pb.ListProductsResponse{
		Products: products,
		TotalCount: int32(len(products)),
	}, nil
}
