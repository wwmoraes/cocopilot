package cocopilot

import (
	"fmt"
)

// Response contains the properties of the non-standard token response that
// Copilot uses.
//
//nolint:tagliatelle // we use a terms as per GitHub implementation
type Response struct {
	ErrorDetails *APIError `json:"error_details,omitempty"`
	Message      string    `json:"message,omitempty"`
	Token
}

// Parse extracts a Copilot token and error from a response.
func (res *Response) Parse() (*Token, error) {
	if res.Message != "" || res.ErrorDetails != nil {
		return nil, fmt.Errorf("%s: %w", res.Message, res.ErrorDetails)
	}

	return &res.Token, nil
}
