package server

import (
	"context"

	pb "github.com/Coosis/go-ecommerce/cart/internal/pb/v1/cart"
	sqlc "github.com/Coosis/go-ecommerce/cart/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Server) RotateAndMergeCart(
	ctx context.Context,
	req *pb.RotateAndMergeCartRequest,
) (*pb.RotateAndMergeCartResponse, error) {
	queries := sqlc.New(s.Pool)

	var userid_uuid pgtype.UUID
	if err := userid_uuid.Scan(req.UserId); err != nil {
		return nil, err
	}

	var oldcart_uuid pgtype.UUID
	if req.OldCartId == "" {
		oldcart_uuid = pgtype.UUID{Valid: false}
	} else if err := oldcart_uuid.Scan(req.OldCartId); err != nil {
		oldcart_uuid = pgtype.UUID{Valid: false}
	}

	sqlc_cart, err := queries.RotateAndMergeCartForUser(ctx, sqlc.RotateAndMergeCartForUserParams{
		UserID:    userid_uuid,
		OldCartID: oldcart_uuid,
	})
	if err != nil {
		return nil, err
	}

	items, err := queries.GetCartItems(ctx, sqlc_cart.ID)
	if err != nil {
		return nil, err
	}

	var cart_items []*pb.CartItem
	for _, item := range items {
		cart_items = append(cart_items, &pb.CartItem{
			Id:   item.ID.String(),
			ProductId: item.ProductID.String(),
			Qty:  item.Qty,
			PriceCentsSnapshot: item.PriceCentsSnapshot,
			CreatedAt: item.CreatedAt.Time.String(),
			UpdatedAt: item.UpdatedAt.Time.String(),
		})
	}

	return &pb.RotateAndMergeCartResponse{
		Cart: &pb.Cart{
			Id:        sqlc_cart.ID.String(),
			UserId:    sqlc_cart.UserID.String(),
			Version:   sqlc_cart.Version,
			Items:     cart_items,
			CreatedAt: sqlc_cart.CreatedAt.Time.String(),
			UpdatedAt: sqlc_cart.UpdatedAt.Time.String(),
		},
	}, nil
}
