// Binary cocopilot uses the homonymous package to retrieve Copilot tokens. It
// only outputs the token to the standard out, allowing you to use the binary
// to directly feed an environment variable, for instance.
//
// It stores both the GitHub API token and the Copilot token in the OS-provided
// keyring, reusing those tokens as possible and renewing otherwise.
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	"github.com/goccy/go-json"
	"github.com/tmc/keyring"
	"github.com/wwmoraes/cocopilot"
	"golang.org/x/oauth2"
)

var (
	exitCode atomic.Int64

	// ErrUnknownResponseContentType occurs if an API returns an unknown content-type response
	ErrUnknownResponseContentType = errors.New("unknown response content type")
)

func fail(format string, args ...any) {
	exitCode.CompareAndSwap(0, int64(1))
	fmt.Fprintf(os.Stderr, format, args...)
	runtime.Goexit()
}

func main() {
	defer func() {
		os.Exit(int(exitCode.Load()))
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	githubToken, err := getGithubToken(ctx)
	if err != nil {
		fail("failed to get a valid Github API token: %s\n", err)
	}

	copilotKeyHandler := JSONKeyHandler[cocopilot.Token]{
		Service:  "https://api.githubcopilot.com",
		Username: cocopilot.ClientID,
	}

	copilotToken, err := copilotKeyHandler.Get()
	if err != nil && !errors.Is(err, keyring.ErrNotFound) {
		fail("failed to retrieve copilot token from keyring: %s\n", err)
	}

	if copilotToken == nil || !time.Unix(copilotToken.ExpiresAt, 0).After(time.Now()) {
		copilotToken, err = getCopilotToken(ctx, githubToken)
		if err != nil {
			fail("failed to get copilot token from API: %s\n", err)
		}

		err = copilotKeyHandler.Set(copilotToken)
		if err != nil {
			fail("failed to store copilot token to keyring: %s\n", err)
		}
	}

	fmt.Fprintln(os.Stdout, copilotToken.AccessToken)
}

func getGithubToken(ctx context.Context) (*oauth2.Token, error) {
	githubKeyHandler := JSONKeyHandler[oauth2.Token]{
		Service:  "https://github.com",
		Username: cocopilot.ClientID,
	}

	githubToken, err := githubKeyHandler.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s; re-authenticating device\n", err)
	}

	tokenSource := oauth2.ReuseTokenSource(githubToken, cocopilot.NewTokenSource(ctx))

	githubToken, err = tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve GitHub token from keyring: %w", err)
	}

	err = githubKeyHandler.Set(githubToken)
	if err != nil {
		return nil, fmt.Errorf("failed to store GitHub token to keyring: %w", err)
	}

	return githubToken, nil
}

func getCopilotToken(ctx context.Context, token *oauth2.Token) (*cocopilot.Token, error) {
	req, err := cocopilot.NewRequest(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %w", err)
	}
	defer res.Body.Close()

	if strings.SplitN(res.Header.Get("Content-Type"), ";", 2)[0] != "application/json" {
		return nil, fmt.Errorf(
			"%w (%s)",
			ErrUnknownResponseContentType,
			res.Header.Get("Content-Type"),
		)
	}

	var copilotResponse cocopilot.Response

	err = json.NewDecoder(res.Body).Decode(&copilotResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode copilot API response: %w", err)
	}

	copilotToken, err := copilotResponse.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to parse copilot response: %w", err)
	}

	return copilotToken, nil
}
