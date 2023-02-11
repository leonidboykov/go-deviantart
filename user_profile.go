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
	User             User                 `json:"user"`
	IsWatching       bool                 `json:"is_watching"`
	ProfileURL       string               `json:"profile_url"`
	UserIsArtist     bool                 `json:"user_is_artist"`
	ArtistLevel      string               `json:"artist_level,omitempty"`
	ArtistSpeciality string               `json:"artist_specialty,omitempty"`
	RealName         string               `json:"real_name"`
	Tagline          string               `json:"tagline"`
	CountryID        uint8                `json:"countryid"`
	Country          string               `json:"country"`
	Website          string               `json:"website"`
	Bio              string               `json:"bio"`
	CoverPhoto       string               `json:"cover_photo,omitempty"`
	CoverDeviation   ModuleCoverDeviation `json:"cover_deviation,omitempty"`
	LastStatus       any                  `json:"last_status,omitempty"` // TODO: Status object.
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
	Session Session `json:"session,omitempty"`
}

type ProfileParams struct {
	// Include collection folder info.
	IncludeCollections bool `url:"ext_collections,omitempty"`

	// Include gallery folder info.
	IncludeGalleries bool `url:"ext_galleries,omitempty"`

	// Session data is not always needed for this endpoint.
	WithSession bool `url:"with_session,omitempty"`
}

// Profile gets user profile information.
func (s *userService) Profile(username string, params *ProfileParams) (Profile, error) {
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
