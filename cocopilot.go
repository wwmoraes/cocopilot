// Package cocopilot provides bindings to retrieve a GitHub Copilot token to
// interact with their LLM models. It uses the same public client as Visual
// Studio Code to authenticate using OAuth to retrieve a device token, and then
// use that token to get the Copilot one.
package cocopilot

const (
	// ClientID contains the VSCode client ID, as we need to impersonate it to get
	// a token
	ClientID = "01ab8ac9400c4e429b23"
)
