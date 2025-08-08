package cocopilot

import (
	"context"
	"fmt"
	"os"

	"github.com/pkg/browser"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/authhandler"
)

var _ oauth2.TokenSource = (*TokenSource)(nil)

// TokenSource implements an [oauth2.TokenSource] to refresh Github API
// tokens. It uses a device authorization code flow to retrieve the first one.
//
//nolint:containedctx // unfortunately [oauth2.TokenSource.Token] doesn't take a context
type TokenSource struct {
	Context         context.Context
	Config          *oauth2.Config
	Handler         authhandler.AuthorizationHandler
	AuthCodeOptions []oauth2.AuthCodeOption
}

// NewTokenSource initializes a GitHub token source to retrieve/renew its core
// API token. It uses scopes and redirect URLs as per the upstream Visual Studio
// Code implementation.
func NewTokenSource(ctx context.Context) *TokenSource {
	return &TokenSource{
		Context: ctx,
		Config: &oauth2.Config{
			ClientID:     ClientID,
			ClientSecret: "",
			// unfortunately we cannot use OpenID Connect discovery with something
			// like [github.com/coreos/go-oidc/v3/oidc.NewProvider] as GitHub's
			// well-known endpoint does not provide any of the needed URLs 👌
			Endpoint: oauth2.Endpoint{
				AuthURL:       "https://github.com/login/oauth/authorize",
				DeviceAuthURL: "https://github.com/login/device/code",
				TokenURL:      "https://github.com/login/oauth/access_token",
				AuthStyle:     oauth2.AuthStyleAutoDetect,
			},
			RedirectURL: "https://vscode.dev/redirect",
			Scopes:      []string{"user:email"},
		},
		AuthCodeOptions: []oauth2.AuthCodeOption{
			oauth2.SetAuthURLParam("prompt", "select_account"),
		},
		//nolint:exhaustruct // result is an internal field handled at runtime
		Handler: (&GithubDeviceAuthGrantFlowHandler{
			Context: ctx,
		}).AuthorizationHandler,
	}
}

// Token initializes and handles a device authorization code flow, polling the
// server to exchange the code for a token.
func (source *TokenSource) Token() (*oauth2.Token, error) {
	deviceAuth, err := source.Config.DeviceAuth(source.Context)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate device authorization grant flow: %w", err)
	}

	fmt.Fprintln(
		os.Stderr,
		"Copy this code and paste on the device verification page:",
		deviceAuth.UserCode,
	)

	err = browser.OpenURL(deviceAuth.VerificationURI)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to open browser:", err)
		fmt.Fprintln(
			os.Stderr,
			"open it manually and navigate to this URL to proceed:",
			deviceAuth.VerificationURI,
		)
	}

	token, err := source.Config.DeviceAccessToken(source.Context, deviceAuth)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange device code for a token: %w", err)
	}

	return token, nil
}
