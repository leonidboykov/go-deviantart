package deviantart

import (
	"context"
	"fmt"

	"github.com/dghubble/sling"
	"golang.org/x/oauth2"
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
	BasicScope      = "basic"
	BrowseScope     = "browse"
	CollectionScope = "collection"
	PublishScope    = "publish"
	StashScope      = "stash"
	UserScope       = "user"
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
		url := conf.AuthCodeURL("state")
		fmt.Printf("Visit the URL for the auth dialog: %v", url)

		srv := authserver.NewAuthServer(conf)
		code, err := srv.ListenToken()
		if err != nil {
			return err
		}
		// Use the authorization code that is pushed to the redirect
		// URL. Exchange will do the handshake to retrieve the
		// initial access token. The HTTP Client returned by
		// conf.Client will refresh the token as necessary.
		tok, err := conf.Exchange(context.Background(), code)
		if err != nil {
			return err
		}

		fmt.Println("Auth Token:", tok.AccessToken)

		s.Client(conf.Client(context.Background(), tok))
		return nil
	}
}
