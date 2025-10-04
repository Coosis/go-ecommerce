package server

import (
	"context"
	"errors"
	"fmt"

	pb "github.com/Coosis/go-ecommerce/auth/internal/pb/v1/auth"
	sqlc "github.com/Coosis/go-ecommerce/auth/sqlc"
	pgx "github.com/jackc/pgx/v5"
)

func(s *Server) GetOAuthSession(
	ctx context.Context,
	req *pb.OAuthSessionRequest,
) (*pb.Session, error) {
	queries := sqlc.New(s.Pool)
	user, err := queries.FindUserByOAuth(ctx, sqlc.FindUserByOAuthParams{
		Provider: req.GetProvider(),
		ProviderUserID: req.GetUserId(),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// create a new user
			_, err := queries.CreateUserOAuth(ctx, sqlc.CreateUserOAuthParams{
				Provider: req.GetProvider(),
				ProviderUserID: req.GetUserId(),
			})
			if err != nil {
				return nil, fmt.Errorf("failed to create user by OAuth: %w", err)
			}

			return s.GetOAuthSession(ctx, req)
		} else {
			return nil, fmt.Errorf("failed to find user by OAuth: %w", err)
		}
	}

	sessionID, err := s.CreateSession(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &pb.Session{
		SessionId: sessionID,
	}, nil
}
