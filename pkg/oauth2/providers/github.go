package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"regexp"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// GitHub represents an GitHub Identity Provider
type GitHub struct {
	oauth2Config *oauth2.Config
	State        string
	ValidateURL  *url.URL
}

var (
	githubDefaultValidateURL = &url.URL{
		Scheme: "https",
		Host:   "api.github.com",
		Path:   "/",
	}
)

// NewGitHubProvider returns a new OAuth2 config for GitHub provider.
func NewGitHubProvider(clientID string, clientSecret string) *GitHub {
	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"user:email", "read:user"},
		Endpoint:     github.Endpoint,
	}
	return &GitHub{
		oauth2Config: conf,
		ValidateURL:  githubDefaultValidateURL,
	}
}

// GetAuthURL AuthCodeURL returns a URL to OAuth 2.0 provider's consent
// page that asks for permissions for the required scopes explicitly.
func (g *GitHub) GetAuthURL(state string) string {
	return g.oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// Exchange converts an authorization code into a token.
func (g *GitHub) Exchange(code string) (*oauth2.Token, error) {
	token, err := g.oauth2Config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	return token, err
}

// ValidateToken reports whether token is valid and have access to a given scope
func (g *GitHub) ValidateToken(token *oauth2.Token) (bool, error) {

	if !token.Valid() {
		return false, fmt.Errorf("Token is invalid, expire %v", token.Expiry)
	}

	client := g.oauth2Config.Client(context.TODO(), token)

	isEmailValid, err := g.validateEmails(client)
	if err != nil {
		return false, err
	}

	isUserValid, err := g.validateUser(client)
	if err != nil {
		return false, err
	}

	return (isEmailValid || isUserValid), nil
}

func (g *GitHub) validateEmails(client *http.Client) (bool, error) {
	endpoint := &url.URL{
		Scheme: g.ValidateURL.Scheme,
		Host:   g.ValidateURL.Host,
		Path:   path.Join(g.ValidateURL.Path, "/user/emails"),
	}

	resp, err := client.Get(endpoint.String())
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	bodyString := string(body)

	if resp.StatusCode == http.StatusUnauthorized {
		return false, fmt.Errorf("Cannot validate token: %v", bodyString)
	}

	var result []interface{}
	json.Unmarshal([]byte(body), &result)

	for _, data := range result {
		email := data.(map[string]interface{})
		if email["verified"].(bool) && email["primary"].(bool) {
			result := find(viper.GetStringSlice("server.oauth2-email"), email["email"].(string))
			if result {
				return true, nil
			}
		}
	}

	return false, nil
}

func (g *GitHub) validateUser(client *http.Client) (bool, error) {
	endpoint := &url.URL{
		Scheme: g.ValidateURL.Scheme,
		Host:   g.ValidateURL.Host,
		Path:   path.Join(g.ValidateURL.Path, "/user"),
	}

	resp, err := client.Get(endpoint.String())
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	bodyString := string(body)

	if resp.StatusCode == http.StatusUnauthorized {
		return false, fmt.Errorf("Cannot validate token: %v", bodyString)
	}

	var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)

	isUserValid := find(viper.GetStringSlice("server.oauth2-github-user"), result["login"].(string))

	return isUserValid, nil
}

func find(slice []string, val string) bool {
	for _, item := range slice {
		re := regexp.MustCompile(item)
		if re.MatchString(val) {
			return true
		}
	}
	return false
}
