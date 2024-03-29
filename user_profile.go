package deviantart

import (
	"fmt"

	"github.com/google/uuid"
)

type ModuleCoverDeviation struct {
	CoverDeviation          Deviation `json:"cover_deviation,omitempty"`
	CoverDeviationIDOffsetY uint32    `json:"cover_deviationid_offset_y"`
	ImageWidth              uint32    `json:"image_width,omitempty"`
	ImageHeight             uint32    `json:"image_height,omitempty"`
	CropX                   uint32    `json:"crop_x,omitempty"`
	CropY                   uint32    `json:"crop_y,omitempty"`
	CropWidth               uint32    `json:"crop_width,omitempty"`
	CropHeight              uint32    `json:"crop_height,omitempty"`
}

type Profile struct {
	User             User                  `json:"user"`
	IsWatching       bool                  `json:"is_watching"`
	ProfileURL       string                `json:"profile_url"`
	UserIsArtist     bool                  `json:"user_is_artist"`
	ArtistLevel      string                `json:"artist_level,omitempty"`
	ArtistSpeciality string                `json:"artist_specialty,omitempty"`
	RealName         string                `json:"real_name"`
	Tagline          string                `json:"tagline"`
	CountryID        uint8                 `json:"countryid"`
	Country          string                `json:"country"`
	Website          string                `json:"website"`
	Bio              string                `json:"bio"`
	CoverPhoto       string                `json:"cover_photo,omitempty"`
	CoverDeviation   *ModuleCoverDeviation `json:"cover_deviation,omitempty"`
	LastStatus       *Status               `json:"last_status,omitempty"`
	Stats            struct {
		UserDeviations   uint32 `json:"user_deviations"`
		UserFavourites   uint32 `json:"user_favourites"`
		UserComments     uint32 `json:"user_comments"`
		ProfilePageViews uint32 `json:"profile_pageviews"`
		ProfileComments  uint32 `json:"profile_comments"`
	} `json:"stats"`
	Collections []Folder `json:"collections,omitempty"`
	Galleries   []struct {
		FolderID uuid.UUID `json:"folderid"`
		Parent   uuid.UUID `json:"parent,omitempty"`
		Name     string    `json:"name"`
	} `json:"galleries,omitempty"`
}

type GetProfileParams struct {
	// Include collection folder info.
	IncludeCollections bool `url:"ext_collections,omitempty"`

	// Include gallery folder info.
	IncludeGalleries bool `url:"ext_galleries,omitempty"`
}

// Profile gets user profile information.
//
// To connect to this endpoint OAuth2 Access Token from the Client Credentials
// Grant, or Authorization Code Grant is required.
//
// The following scopes are required to access this resource:
//
//   - browse
func (s *UserService) Profile(username string, params *GetProfileParams) (Profile, error) {
	var (
		success Profile
		failure Error
	)
	_, err := s.sling.New().Post("profile/").Path(username).QueryStruct(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return Profile{}, fmt.Errorf("unable to fetch profile: %w", err)
	}
	return success, nil
}

// Posts returns all journals & status updates for a given user in a single feed.
//
// To connect to this endpoint OAuth2 Access Token from the Client Credentials
// Grant, or Authorization Code Grant is required.
//
// The following scopes are required to access this resource:
//
//   - browse
func (s *UserService) Posts(username string, page *CursorParams) (CursorResponse[Deviation], error) {
	var (
		success CursorResponse[Deviation]
		failure Error
	)
	params := &usernameParams{Username: username}
	_, err := s.sling.New().Get("profile/posts").QueryStruct(params).QueryStruct(page).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return CursorResponse[Deviation]{}, fmt.Errorf("unable to fetch profile: %w", err)
	}
	return success, nil
}

const (
	ArtistLevelNone         = "None"
	ArtistLevelStudent      = "Student"
	ArtistLevelHobbyist     = "Hobbyist"
	ArtistLevelProfessional = "Professional"
)

const (
	ArtistSpecialityNone          = "None"
	ArtistSpecialityArtisanCrafts = "Artisan Crafts"
	ArtistLevelDesignInterfaces   = "Design & Interfaces"
	ArtistLevelDigitalArt         = "Digital Art"
	ArtistLevelFilmAnimation      = "Film & Animation"
	ArtistLevelLiterature         = "Literature"
	ArtistLevelPhotography        = "Photography"
	ArtistLevelTraditionalArt     = "Traditional Art"
	ArtistLevelOther              = "Other"
	ArtistLevelVaried             = "Varied"
)

type UserInfoParams struct {
	// Is the user an artist?
	UserIsArtist bool `url:"user_is_artist,omitempty"`

	// If the user is an artist, what level are they.
	ArtistLevel string `url:"artist_level,omitempty"`

	// If the user is an artist, what is their specialty.
	ArtistSpeciality string `url:"artist_specialty,omitempty"`

	// The users location.
	CountryID int `url:"countryid,omitempty"`

	// The users personal website.
	Website      string `url:"website,omitempty"`
	WebsiteLabel string `url:"website_label,omitempty"`

	// The users tagline.
	Tagline string `url:"tagline,omitempty"`

	ShowBadges  bool     `url:"show_badges,omitempty"`
	Interests   []string `url:"interests,omitempty"`    // TODO: Check positional params.
	SocialLinks []string `url:"social_links,omitempty"` // TODO: Check positional params.
}

// UpdateProfile updates the users profile information.
//
// Check [Countries] to get a list of countries and their IDs.
func (s *UserService) UpdateProfile(params *UserInfoParams) (bool, error) {
	var (
		success map[string]any
		failure Error
	)
	_, err := s.sling.New().Post("profile/update").BodyForm(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return false, fmt.Errorf("unable to update user profile: %w", err)
	}
	return success["success"].(bool), nil
}
