package cocopilot

import "fmt"

// APIError contains details about an error response.
//
//nolint:tagliatelle // we use a terms as per GitHub implementation
type APIError struct {
	URL            string `json:"url,omitempty"`
	Message        string `json:"message,omitempty"`
	Title          string `json:"title,omitempty"`
	NotificationID string `json:"notification_id,omitempty"`
}

func (err *APIError) Error() string {
	if err == nil {
		return ""
	}

	return fmt.Sprintf("%s; see %s", err.Message, err.URL)
}
