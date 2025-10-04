package util

import "net/http"

const CART_ID = "cart_id"

func SetCartCookie(w http.ResponseWriter, r *http.Request, cart_id string) {
	http.SetCookie(w, &http.Cookie{
		Name:     CART_ID,
		Value:    cart_id,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})
	r.AddCookie(&http.Cookie{Name: CART_ID, Value: cart_id})
}
