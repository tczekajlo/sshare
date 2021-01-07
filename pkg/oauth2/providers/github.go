package providers

import (
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// GitHub represents an GitHub Identity Provider
type GitHub struct {
	oauth2Config *oauth2.Config
	State        string
}

func NewGitHubProvider() *GitHub {
	conf := &oauth2.Config{
		ClientID:     "f79797c93fd0d6b4cadb",
		ClientSecret: "c58b42235fc6b098bcb963adc10f4856e2d7b5b4",
		Scopes:       []string{"user:email", "read:user"},
		Endpoint:     github.Endpoint,
	}
	return &GitHub{
		oauth2Config: conf,
	}
}

func (g *GitHub) GetAuthURL() string {
	g.State = "ssharestateoauth2"
	return g.oauth2Config.AuthCodeURL(g.State, oauth2.AccessTypeOffline)
}

func (g *GitHub) ServerCallbackEndpoint() {
	http.HandleFunc("/callback", g.handleCallback)
	server := &http.Server{Addr: ":8080"}
	server.ListenAndServe()

}

func (g *GitHub) handleCallback(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("#%v", r)
	//content, _ := getUserInfo(, r.FormValue("code"))
	token, err := g.oauth2Config.Exchange(oauth2.NoContext, r.FormValue("code"))
	if err != nil {
		fmt.Println("code exchange failed:", err.Error())
	}
	fmt.Println(token)
	//if s.TLS.AuthState != r.FormValue("state") {
	//	s.log.Error("Invalid oauth state")
	//	os.Exit(1)
	//}

	fmt.Fprintf(w, "Success, please go back to the terminal\n")
	return
}
