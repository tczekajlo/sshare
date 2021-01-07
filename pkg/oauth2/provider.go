package oauth2

import (
	"sshare/pkg/oauth2/providers"
)

// Provider represents an upstream identity provider implementation
type Provider interface {
	GetAuthURL() string
	ServerCallbackEndpoint()
}

// New provides a new Provider based on the --oauth2-provider flag
func New() Provider {
	switch "github" {
	case "github":
		return providers.NewGitHubProvider()
	default:
		return nil
	}
}
