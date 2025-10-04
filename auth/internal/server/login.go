package server

import (
	"context"
	"errors"
	"fmt"

	pb "github.com/Coosis/go-ecommerce/auth/internal/pb/v1/auth"
	sqlc "github.com/Coosis/go-ecommerce/auth/sqlc"
	pgx "github.com/jackc/pgx/v5"
)

func(s *Server) Login(
	ctx context.Context,
	req *pb.LoginClassicRequest,
) (*pb.Session, error) {
	queries := sqlc.New(s.Pool)
	user, err := queries.FindUserByContactMethod(ctx, sqlc.FindUserByContactMethodParams{
		Type: req.GetType(),
		Value: req.GetValue(),
		PasswordHash: req.GetPassword(),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("User not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	
	if user.IsVerified == false {
		return nil, fmt.Errorf("user is not verified")
	}

	session_id, err := s.CreateSession(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &pb.Session{
		SessionId: session_id,
	}, nil
}
