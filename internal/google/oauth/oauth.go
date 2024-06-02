package oauth

import (
	"context"
	"errors"
	"golang.org/x/oauth2"
	"strings"
)

const GoogleAuthURL = "https://accounts.google.com/o/oauth2/auth"
const GoogleTokenURL = "https://oauth2.googleapis.com/token"

type OAuth interface {
	AuthCodeURL(state string, isForce bool) string
	Token(ctx context.Context, code string) (*oauth2.Token, error)
	IsAllowSync(scope []string) bool
}

type googleOAuth struct {
	config *oauth2.Config
}

type Config struct {
	ClientId     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	RedirectUri  string   `json:"redirect_uri"`
	Scopes       []string `json:"scopes"`
}

func New(c Config) (OAuth, error) {
	if c.ClientId == "" || c.ClientSecret == "" || c.RedirectUri == "" || c.Scopes == nil {
		return nil, errors.New("client_id, client_secret, redirect_uri, scopes must be set")
	}

	cfg := oauth2.Config{
		ClientID:     c.ClientId,
		ClientSecret: c.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  GoogleAuthURL,
			TokenURL: GoogleTokenURL,
		},
		RedirectURL: c.RedirectUri,
		Scopes:      c.Scopes,
	}

	return &googleOAuth{config: &cfg}, nil
}

func (g *googleOAuth) AuthCodeURL(state string, isForce bool) string {
	if isForce {
		return g.config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	}
	return g.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (g *googleOAuth) Token(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (g *googleOAuth) IsAllowSync(scope []string) bool {
	for i := range scope {
		if strings.Contains(scope[i], "calendar.events.owned.readonly") {
			return true
		}
	}

	return false
}
