package util

import (
	pb "github.com/Coosis/go-ecommerce/cart/internal/pb/v1/cart"
	sqlc "github.com/Coosis/go-ecommerce/cart/sqlc"
)

// convert sqlc cart to pb cart, ignoring items
func PbCartFromSqlcCart(c *sqlc.Cart) (*pb.Cart, error) {
	return &pb.Cart{
		Id:        c.ID.String(),
		UserId:    c.UserID.String(),
		Version:   c.Version,
		Items:     nil,
		CreatedAt: c.CreatedAt.Time.String(),
		UpdatedAt: c.UpdatedAt.Time.String(),
	}, nil
}

