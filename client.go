package deviantart

import (
	"github.com/dghubble/sling"
)

const deviantArtURL = "https://www.deviantart.com/api/v1/oauth2/"

// Client provices access to DeviantArt API endpoint.
type Client struct {
	Stash *stashService
}

// TODO: Add http.Client to args.
func NewClient(auth Authenticator) (*Client, error) {
	sling := sling.New().Base(deviantArtURL)
	if err := auth(sling); err != nil {
		return nil, err
	}

	c := Client{
		Stash: newStashService(sling.New()),
	}
	return &c, nil
}
