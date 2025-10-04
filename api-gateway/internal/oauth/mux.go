package oauth

import (
	"fmt"
	"os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	grpc "google.golang.org/grpc"
)

// NOTE:
// for any provider integration, you SHOULD be able to handle an option `redirect` query parameter
// and after session creation, redirect the user to the specified URL if exists

const (
	// github
	GITHUB_OAUTH_URL = "https://github.com/login/oauth/authorize"
	GITHUB_STATE = "sdklfjslkdjfehguytuk4rty"
	GITHUB_CLIENT_ID = "GITHUB_CLIENT_ID"
	GITHUB_CLIENT_SECRET = "GITHUB_CLIENT_SECRET"
	GITHUB_EXCHANGE_URL = "https://github.com/login/oauth/access_token"
	GITHUB_USER_URL = "https://api.github.com/user"
)

// !!Note that prefix does not include trailing slash
func NewOAuthMux(
	mux *mux.Router,
    prefix string,
	authService *grpc.ClientConn,
) {
	router := mux.StrictSlash(true)

	// github
	github_client_id := os.Getenv(GITHUB_CLIENT_ID)
	redirect_uri := fmt.Sprintf("%s/github/callback", prefix)
	log.Debugf("Redirect URI for GitHub OAuth: %s", redirect_uri)
	github_client_secret := os.Getenv(GITHUB_CLIENT_SECRET)
	if github_client_id == "" {
		log.Fatal("GITHUB_CLIENT_ID environment variable is not set")
	}
	if github_client_secret == "" {
		log.Fatal("GITHUB_CLIENT_SECRET environment variable is not set")
	}
	router.HandleFunc("/github", GithubOAuthHandler(
		github_client_id, 
		redirect_uri,
		GITHUB_STATE,
	))

	router.HandleFunc("/github/callback", GithubOAuthCallbackHandler(
		github_client_id,
		github_client_secret,
		redirect_uri,
		GITHUB_EXCHANGE_URL,
		GITHUB_USER_URL,
		GITHUB_STATE,
		authService,
	))
}
