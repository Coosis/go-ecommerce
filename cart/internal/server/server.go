package server

import (
	"context"

	pb "github.com/Coosis/go-ecommerce/cart/internal/pb/v1/cart"
	pbcatalog "github.com/Coosis/go-ecommerce/cart/internal/pb/v1/catalog"
	pborder "github.com/Coosis/go-ecommerce/cart/internal/pb/v1/order"
	sqlc "github.com/Coosis/go-ecommerce/cart/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	CART_ID = "cart_id"
	X_USER_ID = "x-user-id"
)

type Server struct {
	Pool sqlc.DBTX

	CatalogClient pbcatalog.CatalogServiceClient
	OrderClient pborder.OrderServiceClient

	pb.UnimplementedCartServiceServer
}

func NewServer(
	pool sqlc.DBTX,
	catalogClient pbcatalog.CatalogServiceClient,
	orderClient pborder.OrderServiceClient,
) *Server {
	return &Server{
		Pool: pool,
		CatalogClient: catalogClient,
		OrderClient: orderClient,
	}
}

func(s *Server) CartItems(
	ctx context.Context,
	cart_id pgtype.UUID,
) ([]*pb.CartItem, error) {
	queries := sqlc.New(s.Pool)
	sqlcitems, err := queries.GetCartItems(ctx, cart_id)
	if err != nil {
		return nil, err
	}
	items := []*pb.CartItem{}
	for _, item := range sqlcitems {
		items = append(items, &pb.CartItem{
			Id: item.ID.String(),
			ProductId: item.ProductID.String(),
			SkuId: item.SkuID.String(),
			Qty: item.Qty,
			PriceCentsSnapshot: item.PriceCentsSnapshot,
			CreatedAt: item.CreatedAt.Time.String(),
			UpdatedAt: item.UpdatedAt.Time.String(),
		})
	}
	return items, nil
}

func(s *Server) CartItemsForCheckout(
	ctx context.Context,
	cart_id pgtype.UUID,
) ([]*pborder.OrderItem, error) {
	items, err := s.CartItems(ctx, cart_id)
	if err != nil {
		return nil, err
	}

	var orderItems []*pborder.OrderItem
	for _, item := range items {
		orderItems = append(orderItems, &pborder.OrderItem{
			ProductId: item.ProductId,
			SkuId: item.SkuId,
			Qty: item.Qty,
		})
	}
	return orderItems, nil
}
