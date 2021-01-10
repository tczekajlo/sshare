package oauth2

import (
	"sshare/pkg/oauth2/providers"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

// Provider represents an upstream identity provider implementation
type Provider interface {
	GetAuthURL(state string) string
	Exchange(code string) (*oauth2.Token, error)
	ValidateToken(token *oauth2.Token) (bool, error)
}

// New provides a new Provider based on the --oauth2-provider flag
func New() Provider {
	switch "github" {
	case "github":
		return providers.NewGitHubProvider(
			viper.GetString("server.oauth2-client-id"),
			viper.GetString("server.oauth2-client-secret"),
		)
	default:
		return nil
	}
}
