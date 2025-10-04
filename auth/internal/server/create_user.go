package server

import (
	"context"
	"errors"
	"fmt"

	pb "github.com/Coosis/go-ecommerce/auth/internal/pb/v1/auth"
	sqlc "github.com/Coosis/go-ecommerce/auth/sqlc"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	UNIQUE_VIOLATION = "23505"
)

func(s *Server) CreateUserClassic(
	ctx context.Context,
	req *pb.CreateUserClassicRequest,
) (*pb.CreateUserClassicResponse, error) {
	queries := sqlc.New(s.Pool)
	user, err := queries.CreateUserClassic(ctx, sqlc.CreateUserClassicParams{
		PasswordHash: req.GetPassword(),
		Type: req.GetType(),
		Value: req.GetValue(),
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == UNIQUE_VIOLATION{
				existingUser, err := queries.FindUserByContactValue(ctx, sqlc.FindUserByContactValueParams{
					Type: req.GetType(),
					Value: req.GetValue(),
				})
				if err != nil {
					return nil, fmt.Errorf("internal error: %v", err)
				}

				challenge_id, err := s.CreateChallenge(ctx, existingUser.ID.String(), req.GetValue())
				if err != nil {
					return nil, fmt.Errorf("failed to create challenge: %v", err)
				}

				return &pb.CreateUserClassicResponse{
					UserId: user.ID.String(),
					Challenge: challenge_id,
				}, nil
			}
		}
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	challenge_id, err := s.CreateChallenge(ctx, user.ID.String(), req.GetValue())
	if err != nil {
		return nil, fmt.Errorf("failed to create challenge: %v", err)
	}

	return &pb.CreateUserClassicResponse{
		UserId: user.ID.String(),
		Challenge: challenge_id,
	}, nil
}

