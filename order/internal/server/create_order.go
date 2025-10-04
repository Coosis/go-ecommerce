package server

import (
	"context"
	"fmt"

	pbcatalog "github.com/Coosis/go-ecommerce/order/internal/pb/v1/catalog"
	pbinventory "github.com/Coosis/go-ecommerce/order/internal/pb/v1/inventory"
	pb "github.com/Coosis/go-ecommerce/order/internal/pb/v1/order"
	sqlc "github.com/Coosis/go-ecommerce/order/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	log "github.com/sirupsen/logrus"

	"crypto/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

var entropy = ulid.Monotonic(rand.Reader, 0)
func NewOrderNumber() string {
  return ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
}

func(s *Server) CreateOrder(
	ctx context.Context,
	req *pb.CreateOrderRequest,
) (*pb.CreateOrderResponse, error) {
	queries := sqlc.New(s.Pool)
	var user_id_uuid pgtype.UUID
	if err := user_id_uuid.Scan(req.UserId); err != nil {
		return nil, err
	}

	var cart_id_uuid pgtype.UUID
	if err := cart_id_uuid.Scan(*req.CartId); err != nil {
		return nil, err
	}

	subtotal := 0
	var items []sqlc.AddOrderItemParams
	for _, item := range req.Items {
		catalog_resp, err := s.CatalogClient.GetProduct(ctx, &pbcatalog.GetProductRequest{
			ProductId: item.ProductId,
		})
		if err != nil {
			return nil, err
		}

		inventory_resp, err := s.InventoryClient.GetSkuCode(ctx, &pbinventory.GetSkuCodeRequest{
			SkuId: item.SkuId,
		})
		if err != nil {
			return nil, err
		}

		var skucode pgtype.Text
		if err := skucode.Scan(inventory_resp.Code); err != nil {
			return nil, fmt.Errorf("invalid sku code: %v", err)
		}

		subtotal += int(item.Qty) * int(catalog_resp.PriceCents)
		var product_id_uuid pgtype.UUID
		if err := product_id_uuid.Scan(item.ProductId); err != nil {
			return nil, fmt.Errorf("invalid product id: %v", err)
		}
		var sku_id_uuid pgtype.UUID
		if err := sku_id_uuid.Scan(item.SkuId); err != nil {
			return nil, fmt.Errorf("invalid sku id: %v", err)
		}
		items = append(items, sqlc.AddOrderItemParams{
			OrderID: pgtype.UUID{}, // will be set later
			ProductID: product_id_uuid,
			ProductName: catalog_resp.Name,
			SkuID: sku_id_uuid,
			SkuCode: skucode,
			Qty: item.Qty,
			UnitPriceCents: int64(catalog_resp.PriceCents),
			DiscountCents: 0,
			TaxRateBp: 0,
			TotalLineCents: 0,
			PriceVersion: pgtype.Int8{Valid: true, Int64: int64(catalog_resp.PriceVersion)},
			Metadata: pgtype.Text{Valid: false},
		})
		log.Infof("Added item: %+v", items[len(items)-1])
	}
	// TODO
	order, err := queries.CreateOrder(ctx, sqlc.CreateOrderParams{
		OrderNumber:   NewOrderNumber(),
		UserID:        user_id_uuid,
		Email:         pgtype.Text{Valid: false},
		Phone:         pgtype.Text{Valid: false},
		CartID:        cart_id_uuid,
		Currency:      "USD",
		SubtotalCents: int64(subtotal),
		DiscountCents: 0,
		TaxCents:      0,
		ShippingCents: 0,
		TotalCents:    int64(subtotal),
		PaymentIntentID: pgtype.Text{Valid: false},
		Notes:         pgtype.Text{Valid: false},
		ShippingAddr: "",
		BillingAddr:  pgtype.Text{Valid: false},
	})
	if err != nil {
		log.Error("failed to create order: ", err)
		return nil, err
	}

	ok := true
	for _, item := range items {
		item.OrderID = order.ID
		queries.AddOrderItem(ctx, item)
	}

	if !ok {
		_, err := queries.CancelOrder(ctx, order.ID)
		if err != nil {
			log.Error("failed to cancel order: ", err)
			return nil, err
		}
		return &pb.CreateOrderResponse{
			OrderId: order.ID.String(),
			Status: "FAILED",
		}, nil
	}

	return &pb.CreateOrderResponse{
		OrderId: order.ID.String(),
		Status: "CREATED",
	}, nil
}
