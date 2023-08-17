package deviantart

import (
	"context"

	"github.com/dghubble/sling"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/authhandler"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/leonidboykov/deviantart/internal/authserver"
)

// Authenticator describes authentication pipeline.
type Authenticator func(s *sling.Sling) error

// ClientCredentials allows gives access to "public" endpoints and do not
// require user authorization. Use this method to access read-only endpoints.
func ClientCredentials(clientID, clientSecret string) Authenticator {
	conf := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     "https://www.deviantart.com/oauth2/token",
	}
	return func(s *sling.Sling) error {
		s.Client(conf.Client(context.Background()))
		return nil
	}
}

const (
	BasicScope       = "basic"
	BrowseScope      = "browse"
	BrowseMLTScope   = "browse.mlt"
	CollectionScope  = "collection"
	CommentPostScope = "comment.post"
	PublishScope     = "publish"
	StashScope       = "stash"
	UserScope        = "user"
)

// AuthorizationCode grant is the most common OAuth2 grant type and gives access
// to aspects of a users account. Use this method if you need to upload images.
func AuthorizationCode(clientID, clientSecret string, scopes []string, callbackURL string) Authenticator {
	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.deviantart.com/oauth2/authorize",
			TokenURL: "https://www.deviantart.com/oauth2/token",
		},
		RedirectURL: callbackURL,
		Scopes:      scopes,
	}
	return func(s *sling.Sling) error {
		tok, err := authhandler.TokenSource(
			context.Background(),
			conf,
			"state", // TODO: This is unsecure.
			authserver.AuthHandler(callbackURL),
		).Token()
		if err != nil {
			return nil
		}

		s.Client(conf.Client(context.Background(), tok))
		return nil
	}
}
