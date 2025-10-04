package middleware

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pbauth "github.com/Coosis/go-ecommerce/api-gateway/internal/pb/v1/auth"
)

const (
	SESSION_ID = "session_id"
	X_USER_ID     = "x-user-id"

	VERIFY_SESSION_TIMEOUT = 3 * 1000 // milliseconds
)

func AuthMiddleware(
	handler http.Handler,
	site_redirect string,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// exclude paths that starts with /v1/auth
		println("Request URL:", r.URL.Path)
		if strings.HasPrefix(r.URL.Path, "/gateway/v1/auth") {
			handler.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie(SESSION_ID)
		if err != nil || cookie == nil || cookie.Value == "" {
			uri := "http://localhost:9765/v1/auth/oauth/github"
			if site_redirect != "" {
				uri += "?site_redirect=" + url.QueryEscape(site_redirect)
			}

			log.Debugf("Redirecting to auth URI: %s", uri)
			http.Redirect(w, r, uri, http.StatusFound)
			return
		}

		ctx := metadata.NewIncomingContext(r.Context(), metadata.Pairs(SESSION_ID, cookie.Value))

		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

func SoftUserIDResolver(
	handler http.Handler,
	client pbauth.AuthServiceClient,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(SESSION_ID)
		if err == nil && cookie != nil && cookie.Value != "" {
			resp, err := client.VerifySession(r.Context(), &pbauth.VerifySessionRequest{
				SessionId: cookie.Value,
			})
			if err != nil {
				log.Errorf("Failed to verify session: %v", err)
				handler.ServeHTTP(w, r)
				return
			}
			user_id := resp.GetUserId()
			log.Debugf("Resolved user id: %s", user_id)
			ctx := metadata.NewIncomingContext(
				r.Context(),
				metadata.Pairs(X_USER_ID, user_id),
			)
			handler.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func SoftUserIDResolverInterceptor(
	client pbauth.AuthServiceClient,
) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		// skip for auth
		if strings.Contains(method, "AuthService") {
			log.Debug("Skipping auth for AuthService")
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		var session_id string
		md, ok := metadata.FromOutgoingContext(ctx)
		if ok && md != nil {
			vals := md.Get(SESSION_ID)
			if len(vals) > 0 { session_id = vals[0] }
		}
		if session_id == "" {
			if ctx_md, ok := metadata.FromIncomingContext(ctx); ok {
				vals := ctx_md.Get(SESSION_ID)
				if len(vals) > 0 { session_id = vals[0] }
			}
		}

		if session_id != "" {
			verifyCtx, cancel := context.WithTimeout(ctx, VERIFY_SESSION_TIMEOUT * time.Millisecond)
			defer cancel()

			resp, err := client.VerifySession(verifyCtx, &pbauth.VerifySessionRequest{
				SessionId: session_id,
			})
			if err == nil && resp.GetUserId() != "" {
				log.Debugf("Resolved user id: %s using interceptor", resp.GetUserId())
				ctx = metadata.AppendToOutgoingContext(ctx, X_USER_ID, resp.GetUserId())
			}
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
