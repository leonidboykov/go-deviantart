package deviantart

import (
	"github.com/dghubble/sling"
)

const deviantArtURL = "https://www.deviantart.com/api/v1/oauth2/"

type StatusResponse struct {
	Status string `json:"status"`
}

// Client provices access to DeviantArt API endpoint.
type Client struct {
	base        *sling.Sling
	Browse      *browseService
	Collections *collectionsService
	Deviation   *deviationService
	Gallery     *galleryService
	Stash       *stashService
	User        *userService
}

// TODO: Add http.Client to args.
func NewClient(auth Authenticator) (*Client, error) {
	sling := sling.New().Base(deviantArtURL)
	if err := auth(sling); err != nil {
		return nil, err
	}

	c := Client{
		base:        sling,
		Browse:      newBrowseService(sling.New()),
		Collections: newCollectionsService(sling.New()),
		Deviation:   newDeviationService(sling.New()),
		Gallery:     newGalleryService(sling.New()),
		Stash:       newStashService(sling.New()),
		User:        newUserService(sling.New()),
	}
	return &c, nil
}
