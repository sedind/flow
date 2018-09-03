package googleauth

import (
	"encoding/json"
	"io/ioutil"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// User is a retrieved and authentiacted user.
type User struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
	Locale        string `json:"locale"`
	HostedDomain  string `json:"hd"`
}

// GoogleAuth provides operations needed for google authentication
type GoogleAuth struct {
	id     string
	secret string
	cfg    *oauth2.Config
}

// New creates instance of google auth for provided cid and secret
func New(clientID string, clientSecret string, redirectURL string) *GoogleAuth {
	return &GoogleAuth{
		id:     clientID,
		secret: clientSecret,
		cfg: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}

// LoginURL returns a URL to OAuth 2.0 provider's consent page
// that asks for permissions for the required scopes explicitly.
func (g *GoogleAuth) LoginURL(state string) string {
	return g.cfg.AuthCodeURL(state)
}

// Authenticate user by provided code and return user info
func (g *GoogleAuth) Authenticate(code string) (User, error) {
	user := User{}
	tok, err := g.cfg.Exchange(oauth2.NoContext, code)
	if err != nil {
		return user, err
	}

	client := g.cfg.Client(oauth2.NoContext, tok)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return user, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return user, err
	}

	err = json.Unmarshal(data, &user)
	return user, err
}
