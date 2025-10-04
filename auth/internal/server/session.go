package server

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	pb "github.com/Coosis/go-ecommerce/auth/internal/pb/v1/auth"
	sqlc "github.com/Coosis/go-ecommerce/auth/sqlc"
	log "github.com/sirupsen/logrus"
)

func(s *Server) CreateSession(user_id pgtype.UUID) (string, error) {
	session_id, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("failed to create new session ID: %w", err)
	}
	ctx := context.Background()
	err = s.VKclient.Do(ctx, s.VKclient.
		B().
		Set().
		Key("session:" + session_id.String()).
		Value(user_id.String()).
		Build(),
	).Error()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}

	// non-critical, we can ignore
	queries := sqlc.New(s.Pool)
	_, err = queries.UpdateUserLoginTime(ctx, user_id);
	if err != nil {
		log.Errorf("failed to update user login time: %v", err)
	}

	return session_id.String(), nil
}

func(s *Server) VerifySession(
    ctx context.Context,
	req *pb.VerifySessionRequest,
) (*pb.VerifySessionResponse, error) {
	val, err := s.VKclient.Do(ctx, s.VKclient.
		B().
		Get().
		Key("session:" + req.GetSessionId()).
		Build(),
	).ToString()
	if err != nil {
		return nil, fmt.Errorf("failed to verify session: %w", err)
	}

	if val == "" {
		return nil, fmt.Errorf("session not found")
	}

	return &pb.VerifySessionResponse{
		UserId: val,
	}, nil
}

func(s *Server) InvalidateSession(
	ctx context.Context,
	req *pb.InvalidateSessionRequest,
) (*pb.InvalidateSessionResponse, error) {
	err := s.VKclient.Do(ctx, s.VKclient.
		B().
		Del().
		Key("session:" + req.GetSessionId()).
		Build(),
	).Error()
	if err != nil {
		return nil, fmt.Errorf("failed to invalidate session: %w", err)
	}

	return &pb.InvalidateSessionResponse{
		Success: true,
	}, nil
}
