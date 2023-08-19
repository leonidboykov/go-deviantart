package deviantart

import (
	"fmt"

	"github.com/dghubble/sling"
	"github.com/google/uuid"
)

type User struct {
	UserID       uuid.UUID `json:"userid"`
	UserName     string    `json:"username"`
	UserIcon     string    `json:"usericon"`
	Type         string    `json:"type"`
	IsWatching   bool      `json:"is_watching,omitempty"`
	IsSubscribed bool      `json:"is_subscribed,omitempty"`
	Details      struct {
		Sex      string `json:"sex,omitempty"`
		Age      uint8  `json:"age,omitempty"`
		JoinDate string `json:"joindate"`
	} `json:"details,omitempty"`
	Geo struct {
		Country   string `json:"country"`
		CountryID uint8  `json:"countryid"`
		Timezone  string `json:"timezone"`
	} `json:"get,omitempty"`
	Profile struct {
		UserIsArtist     bool   `json:"user_is_artist"`
		ArtistLevel      string `json:"artist_level,omitempty"`
		ArtistSpeciality string `json:"artist_speciality,omitempty"`
		RealName         string `json:"real_name"`
		Tagline          string `json:"tagline"`
		Website          string `json:"website"`
		CoverPhoto       string `json:"cover_photo"`
	} `json:"profile,omitempty"`
	Stats struct {
		Watchers int32 `json:"watchers"`
		Friends  int32 `json:"friends"`
	} `json:"stats,omitempty"`
	Sidebar struct {
		Watched struct {
			HasNewContent bool `json:"has_new_content"`
			LinkSubnav    struct {
				ContentType string `json:"content_type"`
			} `json:"link_subnav"`
			IsPinned bool `json:"is_pinned"`
		} `json:"watched,omitempty"`
	} `json:"sidebar,omitempty"`
}

type UserService struct {
	sling   *sling.Sling
	Friends *friendsService
}

func newUserService(sling *sling.Sling) *UserService {
	base := sling.Path("user/")
	return &UserService{
		sling:   base,
		Friends: newFriendsService(base.New()),
	}
}

// DAmnToken retrieves the dAmn auth token required to connect to the dAmn servers.
//
// To connect to this endpoint OAuth2 Access Token from the Authorization Code
// Grant is required.
//
// The following scopes are required to access this resource:
//
//   - user
func (s *UserService) DAmnToken() (string, error) {
	var (
		success map[string]any
		failure Error
	)
	_, err := s.sling.New().Get("damntoken").Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return "", fmt.Errorf("unable to fetch dAmn token: %w", err)
	}
	return success["damntoken"].(string), nil
}

// Tiers fetches users tiers.
//
// To connect to this endpoint OAuth2 Access Token from the Authorization Code
// Grant is required.
//
// The following scopes are required to access this resource:
//
//   - user
//
// TODO: Scope is missing in docs.
func (s *UserService) Tiers(username string) ([]Deviation, error) {
	var (
		success singleResponse[Deviation]
		failure Error
	)
	_, err := s.sling.New().Get("tiers/").Path(username).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return nil, fmt.Errorf("unable to user tiers: %w", err)
	}
	return success.Results, nil
}

// Watchers gets the user's list of watchers.
//
// To connect to this endpoint OAuth2 Access Token from the Client Credentials
// Grant, or Authorization Code Grant is required.
//
// The following scopes are required to access this resource:
//
//   - browse
func (s *UserService) Watchers(username string, page *OffsetParams) (OffsetResponse[Friend], error) {
	var (
		success OffsetResponse[Friend]
		failure Error
	)
	_, err := s.sling.New().Get("watchers/").Path(username).QueryStruct(page).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return OffsetResponse[Friend]{}, fmt.Errorf("unable to get user watchers: %w", err)
	}
	return success, nil
}

// Whoami fetches user info of authenticated user.
//
// To connect to this endpoint OAuth2 Access Token from the Authorization Code
// Grant is required.
//
// The following scopes are required to access this resource:
//
//   - user
func (s *UserService) Whoami() (User, error) {
	var (
		success User
		failure Error
	)
	_, err := s.sling.New().Get("whoami").Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return User{}, fmt.Errorf("unable to fetch whoami: %w", err)
	}
	return success, nil
}

// Whois fetches user info for given usernames.
//
// To connect to this endpoint OAuth2 Access Token from the Authorization Code
// Grant is required.
//
// The following scopes are required to access this resource:
//
//   - browse
func (s *UserService) Whois(usernames ...string) ([]User, error) {
	type usernameParams struct {
		Usernames []string `url:"usernames"` // TODO: Implement square brackets with number inside.
	}
	var (
		success singleResponse[User]
		failure Error
	)
	params := &usernameParams{Usernames: usernames}
	_, err := s.sling.New().Post("whois/").BodyForm(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return nil, fmt.Errorf("unable to fetch whois: %w", err)
	}
	return success.Results, nil
}
