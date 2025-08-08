package cocopilot

// Token contains a custom token representation from GitHub. It falls short of
// fitting in a regular [oauth2.Token] by using a distinct property name for the
// access token, hence why it is a custom type.
//
//nolint:tagliatelle // we use a terms as per GitHub implementation
type Token struct {
	// AccessToken contains a session/cookie string, where it sports key-values
	// that looks like this:
	//   tid=...;exp=...;sku=...;...
	//
	// Copilot APIs require this whole value as-is to work.
	AccessToken string `json:"token,omitempty"`
	ExpiresAt   int64  `json:"expires_at,omitempty"`
	RefreshIn   int64  `json:"refresh_in,omitempty"`
}
