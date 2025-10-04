package server

import (
	"context"

	pb "github.com/Coosis/go-ecommerce/inventory/internal/pb/v1/inventory"
	sqlc "github.com/Coosis/go-ecommerce/inventory/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

func(s *Server) Release(
	ctx context.Context,
	req *pb.ReleaseRequest,
) (*pb.ReleaseResponse, error) {
	var reserve_uuid pgtype.UUID
	if err := reserve_uuid.Scan(req.ReservationId); err != nil {
		return nil, err
	}

	queries := sqlc.New(s.Pool)
	_, err := queries.ReleaseReservation(ctx, reserve_uuid)
	if err != nil {
		return nil, err
	}
	return &pb.ReleaseResponse{}, nil
}
