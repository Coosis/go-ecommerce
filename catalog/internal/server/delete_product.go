package server

import (
	"context"

	pb "github.com/Coosis/go-ecommerce/catalog/internal/pb/v1/catalog"
	sqlc "github.com/Coosis/go-ecommerce/catalog/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

func(s *Server) DeleteProduct(
	ctx context.Context,
	req *pb.DeleteProductRequest,
) (*pb.DeleteProductResponse, error) {
	queries := sqlc.New(s.Pool)
	var pg_product_id pgtype.UUID
	if err := pg_product_id.Scan(req.GetProductId()); err != nil {
		return &pb.DeleteProductResponse{Success: false}, err
	}
	err := queries.DeleteProduct(ctx, pg_product_id)
	if err != nil {
		return &pb.DeleteProductResponse{Success: false}, err
	}
	return &pb.DeleteProductResponse{
		Success: true,
	}, nil
}
