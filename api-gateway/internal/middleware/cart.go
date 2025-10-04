package middleware

import (
	"net/http"

	pb "github.com/Coosis/go-ecommerce/api-gateway/internal/pb/v1/cart"
	log "github.com/sirupsen/logrus"
	util "github.com/Coosis/go-ecommerce/api-gateway/util"
	"google.golang.org/grpc/metadata"
	"github.com/google/uuid"
)

// ensure a cart id exists in the context
// if user is logged in: 
//  - if no cart id
//     - if cart with user id exists, use it
//     - if no cart with user id exists, create one and use it
//  - if cart with cart id exists, use it
//  - if no cart with cart id exists, create one and make it active cart, rotate old carts
// if user is not logged in:
//  - if no cart id, create one and use it
//  - if cart with cart id exists, use it
//  - if no cart with cart id exists, create one

// create cart id if not present
func EnsureCartID(
	next http.Handler,
	cart pb.CartServiceClient,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cart_id := ""
		if c, err := r.Cookie(util.CART_ID); err == nil && c != nil && c.Value != "" {
			cart_id = c.Value
		}

		var user_id string
		if md, ok := metadata.FromIncomingContext(r.Context()); ok {
			vals, exists := md[X_USER_ID]
			if exists && len(vals) > 0 {
				user_id = vals[0]
			}
		}
		if user_id != "" {
			// have user id
			log.Infof("User ID found in context: %s, rotating...", user_id)
			rotate_res, err := cart.RotateAndMergeCart(r.Context(), &pb.RotateAndMergeCartRequest{
				UserId: user_id,
				OldCartId: cart_id,
			})
			if err != nil {
				log.Errorf("Error rotating and merging cart: %v", err)
			}
			if rotate_res != nil && rotate_res.Cart.Id != "" {
				if rotate_res.Cart.Id != cart_id {
					util.SetCartCookie(w, r, rotate_res.Cart.Id)
					log.Infof("Rotated and merged cart, new cart ID: %s", cart_id)
				} else {
					log.Infof("Rotated and merged cart, cart ID unchanged: %s", cart_id)
				}
			} else {
				log.Warnf("Rotate and merge cart returned nil or empty cart ID")
			}
		} else {
			// no user id
			log.Infof("No user ID found in context, ensuring cart exists for cart ID: %s", cart_id)
			if cart_id == "" {
				new_cart_id := uuid.New().String()
				util.SetCartCookie(w, r, new_cart_id)
			}
			_, err := cart.GetActiveCartByCartId(r.Context(), &pb.GetActiveCartByCartIdRequest{
				CartId: cart_id,
			})
			if err != nil {
				log.Errorf("Error getting active cart by cart id: %v", err)
			}
		}
		next.ServeHTTP(w, r)
	})
}
