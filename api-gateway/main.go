package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	internal "github.com/Coosis/go-ecommerce/api-gateway/internal"
	pbauth "github.com/Coosis/go-ecommerce/api-gateway/internal/pb/v1/auth"
	pbcatalog "github.com/Coosis/go-ecommerce/api-gateway/internal/pb/v1/catalog"
	pbcart "github.com/Coosis/go-ecommerce/api-gateway/internal/pb/v1/cart"
	pbinventory "github.com/Coosis/go-ecommerce/api-gateway/internal/pb/v1/inventory"
	pbpayment "github.com/Coosis/go-ecommerce/api-gateway/internal/pb/v1/payment"
	middleware "github.com/Coosis/go-ecommerce/api-gateway/internal/middleware"
	util "github.com/Coosis/go-ecommerce/api-gateway/util"
	log "github.com/sirupsen/logrus"

	. "github.com/Coosis/go-ecommerce/api-gateway/internal/oauth"
)

func main() {
	log.SetLevel(log.DebugLevel)
	gwmux := runtime.NewServeMux(
		runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
			md := metadata.MD{}
			if c, err := req.Cookie(middleware.SESSION_ID); err == nil && c.Value != "" {
				md.Append("session_id", c.Value)
			}

			if c, err := req.Cookie(util.CART_ID); err == nil && c.Value != "" {
				md.Append("cart_id", c.Value)
			}

			if in, ok := metadata.FromIncomingContext(ctx); ok {
				if vals := in.Get(middleware.X_USER_ID); len(vals) > 0 && vals[0] != "" {
					md.Append(middleware.X_USER_ID, vals[0]) // "x-user-id"
				}
			}
			return md
		}),
	)

	uri := fmt.Sprintf("%s://localhost:%s", internal.Protocol, internal.Port)

	allowed := map[string]bool{
		"http://localhost:9765": true,
		"http://localhost:8080": true,
	}

	// gorilla_mux := mux.NewRouter().StrictSlash(true)
	gorilla_mux := mux.NewRouter()
	gorilla_mux.Use(middleware.CORS(allowed))

	// auth grpc gateway setup
	auth_conn, err := grpc.NewClient(
		internal.AuthURI,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	defer auth_conn.Close()
	pbauth.RegisterAuthServiceHandler(context.Background(), gwmux, auth_conn)

	auth_client := pbauth.NewAuthServiceClient(auth_conn)
	soft_userid_resolver := middleware.SoftUserIDResolverInterceptor(auth_client)

	// catalog grpc gateway setup
	catalog_conn, err := grpc.NewClient(
		internal.CatalogURI,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(soft_userid_resolver),
	)
	if err != nil {
		panic(err)
	}
	defer catalog_conn.Close()
	pbcatalog.RegisterCatalogServiceHandler(context.Background(), gwmux, catalog_conn)

	// cart service
	cart_conn, err := grpc.NewClient(
		internal.CartURI,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(soft_userid_resolver),
	)
	if err != nil {
		panic(err)
	}
	defer cart_conn.Close()
	pbcart.RegisterCartServiceHandler(context.Background(), gwmux, cart_conn)
	gorilla_mux.PathPrefix("/v1/cart").Handler(
		middleware.SoftUserIDResolver(
			middleware.EnsureCartID(gwmux, pbcart.NewCartServiceClient(cart_conn)),
			// gwmux,
			auth_client,
		),
	)

	// inventory service
	inventory_conn, err := grpc.NewClient(
		internal.InventoryURI,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(soft_userid_resolver),
	)
	if err != nil {
		panic(err)
	}
	defer inventory_conn.Close()
	pbinventory.RegisterInventoryServiceHandler(context.Background(), gwmux, inventory_conn)

	// payment service
	payment_conn, err := grpc.NewClient(
		internal.PaymentURI,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(soft_userid_resolver),
	)
	if err != nil {
		panic(err)
	}
	defer payment_conn.Close()
	pbpayment.RegisterPaymentServiceHandler(context.Background(), gwmux, payment_conn)

	oauth_uri := fmt.Sprintf("%s/v1/auth/oauth", uri)
	NewOAuthMux(
		gorilla_mux.PathPrefix("/v1/auth/oauth/").Subrouter(),

		oauth_uri,
		auth_conn,
	)

	PROTECTED := "/protected"
	gorilla_mux.PathPrefix(PROTECTED).Handler(middleware.AuthMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log.Debug("hit!")
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				log.Error("No metadata found in context")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			sessionID := md.Get("session_id")
			if len(sessionID) == 0 {
				log.Error("No session ID found in metadata")
				http.Redirect(w, r, fmt.Sprintf("%s/github", oauth_uri), http.StatusFound)
				return
			}
			w.Write([]byte("Protected resource accessed!" + "\nSession ID: " + sessionID[0] + "\n"))
		}),
		PROTECTED,
	))

	gorilla_mux.PathPrefix("/").Handler(gwmux)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s", internal.Port),
		Handler: gorilla_mux,
	}

	log.Infof("Starting server on %s", uri)

	server.ListenAndServe()
}
