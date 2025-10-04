package server

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/metadata"

	pb "github.com/Coosis/go-ecommerce/cart/internal/pb/v1/cart"
	pborder "github.com/Coosis/go-ecommerce/cart/internal/pb/v1/order"
	log "github.com/sirupsen/logrus"
)

func(s *Server) Checkout(
	ctx context.Context,
	req *pb.CheckoutRequest,
) (*pb.CheckoutResponse, error) {
	metadata, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("no metadata")
	}

	user_id, ok := metadata[X_USER_ID]
	if !ok || len(user_id) == 0 {
		return nil, fmt.Errorf("no user id provided")
	}

	cart_id, ok := metadata[CART_ID]
	if !ok || len(cart_id) == 0 {
		return nil, fmt.Errorf("no cart id provided")
	}
	var cart_id_uuid pgtype.UUID
	if err := cart_id_uuid.Scan(cart_id[0]); err != nil {
		return nil, err
	}

	pb_items, err := s.CartItemsForCheckout(ctx, cart_id_uuid)
	if err != nil {
		log.Errorf("failed to get cart items: %v", err)
		return nil, err
	}

	order_resp, err := s.OrderClient.CreateOrder(ctx, &pborder.CreateOrderRequest{
		UserId: user_id[0],
		Email: nil,
		Phone: nil,
		CartId: &cart_id[0],
		ShippingAddress: "",
		BillingAddress: "",

		Items: pb_items,
	})
	if err != nil {
		log.Errorf("failed to create order: %v", err)
		return nil, err
	}

	return &pb.CheckoutResponse{
		OrderId: order_resp.OrderId,
	}, nil
}
