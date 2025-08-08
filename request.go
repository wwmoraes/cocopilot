package cocopilot

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

//
//nolint:lll // user agents are really long ¯\_(ツ)_/¯
const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Code/1.102.2 Chrome/134.0.6998.205 Electron/35.6.0 Safari/537.36"

// ErrInvalidAPIToken occurs if the API returns an empty access token
var ErrInvalidAPIToken = errors.New("invalid API token")

// NewRequest generates a HTTP request that when sent gets a response that
// contains a session token usable by the Copilot APIs.
//
// Headers are as per VSCode, minus browser-only Sec-Fetch-* ones. It is not
// clear whether Origin and User-Agent make a functional difference; for the
// sake of work's "security wizards" BS we send them anyway.
func NewRequest(ctx context.Context, token *oauth2.Token) (*http.Request, error) {
	if token == nil || token.AccessToken == "" {
		return nil, ErrInvalidAPIToken
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://api.github.com/copilot_internal/v2/token",
		http.NoBody,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Origin", "vscode-file://vscode-app")
	req.Header.Set(
		"User-Agent",
		userAgent,
	)

	token.SetAuthHeader(req)

	return req, nil
}
