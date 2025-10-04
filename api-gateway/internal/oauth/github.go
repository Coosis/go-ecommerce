// uses GITHUB_CLIENT_ID env variable
package oauth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"

	pb "github.com/Coosis/go-ecommerce/api-gateway/internal/pb/v1/auth"
	log "github.com/sirupsen/logrus"
)

type githubState struct {
	Nounce string `json:"nounce"`
	SiteRedirect string `json:"site_redirect"`
}

func encodeState(n, sr string) (string, error) {
	state := &githubState {
		Nounce: n,
		SiteRedirect: sr,
	}
	data, err := json.Marshal(state)
	if err != nil {
		return "", fmt.Errorf("failed to encode state: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(data), nil
}

func decodeState(encodedState string) (*githubState, error) {
	data, err := base64.RawURLEncoding.DecodeString(encodedState)
	if err != nil {
		return nil, fmt.Errorf("failed to decode state: %w", err)
	}

	var state githubState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to unmarshal state: %w", err)
	}
	return &state, nil
}

func GithubOAuthHandler(
	client_id string,
	redirect_uri string,
	state string,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		site_redirect := r.URL.Query().Get("site_redirect")
		params := url.Values{}
		params.Add("client_id", client_id)
		params.Add("redirect_uri", redirect_uri)
		state_with_redirect, err := encodeState(state, site_redirect)
		if err != nil {
			log.Errorf("Failed to encode state: %v", err)
			http.Error(w, fmt.Sprintf("Failed to encode state: %v", err), http.StatusInternalServerError)
			return
		}
		params.Add("state", state_with_redirect)
		dest := fmt.Sprintf("%s?%s", GITHUB_OAUTH_URL, params.Encode())
		http.Redirect(w, r, dest, http.StatusFound)
	}
}

func GithubOAuthCallbackHandler(
	client_id string,
	client_secret string,
	redirect_uri string,
	exchange_url string,
	user_info_url string,
	state string,

	conn *grpc.ClientConn,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		newstate := r.URL.Query().Get("state")
		if code == "" {
			log.Error("Missing code in request")
			http.Error(w, "Missing code", http.StatusBadRequest)
			return
		}

		github_state, err := decodeState(newstate)
		if err != nil {
			log.Errorf("Failed to decode state: %v", err)
			http.Error(w, fmt.Sprintf("Failed to decode state: %v", err), http.StatusBadRequest)
			return
		}
		site_redirect := "/" // actual redirect to our own site, not with github
		if github_state.SiteRedirect != "" {
			site_redirect = github_state.SiteRedirect
		}

		if github_state.Nounce != state {
			http.Error(w, "Invalid state", http.StatusBadRequest)
			return
		}

		params := url.Values{}
		params.Add("client_id", client_id)
		params.Add("client_secret", client_secret)
		params.Add("code", code)
		params.Add("redirect_uri", redirect_uri)
		body := strings.NewReader(params.Encode())

		req, err := http.NewRequest("POST", exchange_url, body)
		if err != nil {
			log.Errorf("Failed to create request: %v", err)
			http.Error(w, "Failed to create request", http.StatusInternalServerError)
			return
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Accept", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Errorf("Failed to exchange code for token: %v", err)
			http.Error(w, "Failed to exchange code for token", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			log.Errorf("Failed to exchange code for token: status code %d", resp.StatusCode)
			http.Error(w, "Failed to exchange code for token", resp.StatusCode)
			return
		}

		type GithubExchangeResponse struct {
			AccessToken string `json:"access_token"`
			TokenType   string `json:"token_type"`
		}

		var exchangeResponse GithubExchangeResponse
		if err := json.NewDecoder(resp.Body).Decode(&exchangeResponse); err != nil {
			log.Errorf("Failed to decode response: %v", err)
			http.Error(w, "Failed to decode response", http.StatusInternalServerError)
			return
		}

		log.Infof("Access Token: %s", exchangeResponse.AccessToken)

		info, err := GithubFetchInfo(exchangeResponse.AccessToken, user_info_url)
		if err != nil {
			log.Errorf("Failed to fetch user info: %v", err)
			http.Error(w, fmt.Sprintf("Failed to fetch user info: %v", err), http.StatusInternalServerError)
			return
		}

		authClient := pb.NewAuthServiceClient(conn)
		uid := strconv.Itoa(info.ID)
		session, err := authClient.GetOAuthSession(r.Context(), &pb.OAuthSessionRequest{
			Provider: "github",
			UserId: uid,
		})
		if err != nil {
			log.Errorf("Failed to get OAuth session: %v", err)
			http.Error(w, fmt.Sprintf("Failed to get OAuth session: %v", err), http.StatusInternalServerError)
			return
		}

		// set session and redirect to home page
		http.SetCookie(w, &http.Cookie{
			Path:    "/",
			Name:     "session_id",
			Value:    session.SessionId,
			Expires:  time.Now().Add(24 * time.Hour),
		})

		http.Redirect(w, r, site_redirect, http.StatusFound)
	}
}

type GithubUserInfo struct {
	ID int `json:"id"`
	Login string `json:"login"`
}

func GithubFetchInfo(
	token string,
	userInfoURL string,
) (*GithubUserInfo, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create user request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer " + token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: status code %d", resp.StatusCode)
	}

	var userInfo GithubUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &userInfo, nil
}

// http://localhost:9765/v1/auth/oauth/github
