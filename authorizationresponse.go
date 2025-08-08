package cocopilot

import (
	"errors"
	"strings"
)

// AuthorizationResponse contains the parameters set by the authorization server
// on the request back to the redirect URI. This includes both the success and
// error cases, as both are sent over a HTTP 302 Found request.
//
// See https://datatracker.ietf.org/doc/html/rfc6749#section-4.1.2
type AuthorizationResponse struct {
	Code             string
	State            string
	ErrorCode        string
	ErrorDescription string
	ErrorURI         string
}

// Error formats the response error if any, otherwise returns nil.
func (response *AuthorizationResponse) Error() error {
	if response.ErrorCode == "" {
		return nil
	}

	errorMessage := strings.Builder{}
	errorMessage.WriteString(response.ErrorCode)

	if response.ErrorDescription != "" {
		errorMessage.WriteString(": " + response.ErrorDescription)
	}

	if response.ErrorURI != "" {
		errorMessage.WriteString(" [" + response.ErrorURI + "]")
	}

	//nolint:err113 // dynamic errors, we won't map those at all
	return errors.New(errorMessage.String())
}
