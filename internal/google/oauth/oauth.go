package oauth

import (
	"context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"strings"
)

type OAuth interface {
	AuthCodeURL(state string, isForce bool) string
	Token(ctx context.Context, code string) (*oauth2.Token, error)
	IsAllowSync(scope []string) bool
}

type googleOAuth struct {
	config *oauth2.Config
}

func New(cfgData []byte, scopes []string) (OAuth, error) {
	cfg, err := google.ConfigFromJSON(cfgData, scopes...)
	if err != nil {
		return nil, err
	}

	return &googleOAuth{cfg}, nil
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
