package server

import (
	"context"

	pb "github.com/Coosis/go-ecommerce/order/internal/pb/v1/order"
	sqlc "github.com/Coosis/go-ecommerce/order/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

func(s *Server) GetOrder(
	ctx context.Context,
	req *pb.GetOrderRequest,
) (*pb.GetOrderResponse, error) {
	queries := sqlc.New(s.Pool)
	var order_id_uuid pgtype.UUID
	if err := order_id_uuid.Scan(req.OrderId); err != nil {
		return nil, err
	}
	order, err := queries.GetOrder(ctx, order_id_uuid)
	if err != nil {
		return nil, err
	}

	order_items, err := queries.GetOrderItems(ctx, order_id_uuid)
	if err != nil {
		return nil, err
	}

	var pb_order_items []*pb.OrderItem
	for _, item := range order_items {
		sku_name := ""
		if item.SkuCode.Valid {
			sku_name = item.SkuCode.String
		}
		var unit_price_cents int32
		unit_price_cents = int32(item.UnitPriceCents)

		pb_order_items = append(pb_order_items, &pb.OrderItem{
			ProductId: item.ProductID.String(),
			ProductName: &item.ProductName,
			SkuId: item.SkuID.String(),
			SkuName: &sku_name,
			Qty: item.Qty,
			PriceCentsSnapshot: &unit_price_cents,
		})
	}

	return &pb.GetOrderResponse{
		OrderId: order.ID.String(),
		UserId: order.UserID.String(),
		Email: order.Email.String,
		Phone: order.Phone.String,
		ShippingAddress: order.ShippingAddress,
		BillingAddress: order.BillingAddress.String,
		Items: pb_order_items,
		Status: string(order.Status),
		CreatedAt: order.CreatedAt.Time.String(),
		UpdatedAt: order.UpdatedAt.Time.String(),
	}, nil
}
