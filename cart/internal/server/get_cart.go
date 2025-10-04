package server

import (
	"context"

	pb "github.com/Coosis/go-ecommerce/cart/internal/pb/v1/cart"
	util "github.com/Coosis/go-ecommerce/cart/internal/util"
	sqlc "github.com/Coosis/go-ecommerce/cart/sqlc"
	log "github.com/sirupsen/logrus"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/metadata"
)

func (s *Server) GetActiveCartByUserId(
	ctx context.Context,
	req *pb.GetActiveCartByUserIdRequest,
) (*pb.GetActiveCartByUserIdResponse, error) {
	queries := sqlc.New(s.Pool)
	var userid_uuid pgtype.UUID
	if err := userid_uuid.Scan(req.GetUserId()); err != nil {
		return nil, err
	}
	sqlc_cart, err := queries.GetActiveCartByUserId(ctx, userid_uuid)
	if err != nil {
		return nil, err
	}
	sqlc_items, err := queries.GetCartItems(ctx, sqlc_cart.ID)
	if err != nil {
		return nil, err
	}
	items := []*pb.CartItem{}
	for _, item := range sqlc_items {
		items = append(items, &pb.CartItem{
			Id: item.ID.String(),
			ProductId: item.ProductID.String(),
			Qty: item.Qty,
			PriceCentsSnapshot: item.PriceCentsSnapshot,
			CreatedAt: item.CreatedAt.Time.String(),
			UpdatedAt: item.UpdatedAt.Time.String(),
		})
	}
	return &pb.GetActiveCartByUserIdResponse{
		Cart: &pb.Cart{
			Id: sqlc_cart.ID.String(),
			UserId: sqlc_cart.UserID.String(),
			Version: sqlc_cart.Version,
			Items: items,
			CreatedAt: sqlc_cart.CreatedAt.Time.String(),
			UpdatedAt: sqlc_cart.UpdatedAt.Time.String(),
		},
	}, nil
}

func (s *Server) GetActiveCartByCartId(
	ctx context.Context,
	req *pb.GetActiveCartByCartIdRequest,
) (*pb.GetActiveCartByCartIdResponse, error) {
	queries := sqlc.New(s.Pool)
	var cartid_uuid pgtype.UUID
	if err := cartid_uuid.Scan(req.GetCartId()); err != nil {
		return nil, err
	}
	sqlc_cart, err := queries.GetActiveCartByCartId(ctx, cartid_uuid)
	if err != nil {
		return nil, err
	}

	items, err := s.CartItems(ctx, cartid_uuid)
	if err != nil {
		return nil, err
	}

	cart, err := util.PbCartFromSqlcCart(&sqlc_cart)
	if err != nil {
		return nil, err
	}
	cart.Items = items

	return &pb.GetActiveCartByCartIdResponse{
		Cart: cart,
	}, nil
}

func(s *Server) GetMyCart(
	ctx context.Context,
	_ *pb.GetMyCartRequest,
) (*pb.GetMyCartResponse, error) {
	queries := sqlc.New(s.Pool)
	// user id from ctx(might be empty)
	user_id := ""
	// cart id from ctx
	cart_id := ""

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if val, exists := md[X_USER_ID]; exists && len(val) > 0 {
			user_id = val[0]
		}

		if val, exists := md[CART_ID]; exists && len(val) > 0 {
			cart_id = val[0]
		}
	}

	if cart_id == "" {
		return nil, nil
	}

	var cart_uuid pgtype.UUID
	if err := cart_uuid.Scan(cart_id); err != nil {
		return nil, err
	}

	var user_uuid pgtype.UUID
	if user_id != "" {
		if err := user_uuid.Scan(user_id); err != nil {
			return nil, err
		}
	}

	var sqlc_cart sqlc.Cart
	sqlc_cart, err := queries.GetCart(ctx, cart_uuid)
	if err != nil {
		return nil, err
	}

	// associate if not already
	if user_id != "" && sqlc_cart.UserID.Valid == false {
		log.Infof("Associating cart %s with user %s", cart_id, user_id)
		_ = queries.AssociateCartWithUser(ctx, sqlc.AssociateCartWithUserParams{
			ID: cart_uuid,
			UserID: user_uuid,
		})
	}

	// get cart items
	items, err := s.CartItems(ctx, cart_uuid)
	if err != nil {
		return nil, err
	}

	cart, err := util.PbCartFromSqlcCart(&sqlc_cart)
	if err != nil {
		return nil, err
	}
	cart.Items = items

	return &pb.GetMyCartResponse{
		Cart: cart,
	}, nil
}
