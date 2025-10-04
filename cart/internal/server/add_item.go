package server

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/metadata"

	pb "github.com/Coosis/go-ecommerce/cart/internal/pb/v1/cart"
	pbcatalog "github.com/Coosis/go-ecommerce/cart/internal/pb/v1/catalog"
	sqlc "github.com/Coosis/go-ecommerce/cart/sqlc"
	util "github.com/Coosis/go-ecommerce/cart/internal/util"
	log "github.com/sirupsen/logrus"
)

func(s *Server) AddItem(
	ctx context.Context,
	req *pb.AddItemRequest,
) (*pb.AddItemResponse, error) {
	metadata, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("no metadata")
	}
	cart_id, ok := metadata[CART_ID]
	if !ok || len(cart_id) == 0 {
		return nil, fmt.Errorf("no cart id provided")
	}
	log.Infof("Adding item to cart %s", cart_id[0])
	var cart_id_uuid pgtype.UUID
	if err := cart_id_uuid.Scan(cart_id[0]); err != nil {
		return nil, err
	}

	var product_id_uuid pgtype.UUID
	if err := product_id_uuid.Scan(req.GetProductId()); err != nil {
		return nil, err
	}

	r, err := s.CatalogClient.GetProduct(ctx, &pbcatalog.GetProductRequest{
		ProductId: req.GetProductId(),
	})
	if err != nil {
		return nil, err
	}

	var sku_id_uuid pgtype.UUID
	if err := sku_id_uuid.Scan(req.GetSkuId()); err != nil {
		return nil, err
	}

	queries := sqlc.New(s.Pool)
	queries.InsertItemToCart(ctx, sqlc.InsertItemToCartParams{
		CartID: cart_id_uuid,
		ProductID: product_id_uuid,
		SkuID: sku_id_uuid,
		Qty: req.GetQty(),
		PriceCentsSnapshot: r.GetPriceCents(),
	})

	sqlcCart, err := queries.GetCart(ctx, cart_id_uuid)
	if err != nil {
		return nil, err
	}

	cart, err := util.PbCartFromSqlcCart(&sqlcCart)
	if err != nil {
		return nil, err
	}

	return &pb.AddItemResponse{
		Cart: cart,
	}, nil
}
