package deviantart

import (
	"fmt"

	"github.com/dghubble/sling"
)

type friendsService struct {
	sling *sling.Sling
}

func newFriendsService(sling *sling.Sling) *friendsService {
	return &friendsService{
		sling: sling.Path("friends/"),
	}
}

type Friend struct {
	User       *User     `json:"user"`
	IsWatching bool      `json:"is_watching"`
	WatchesYou bool      `json:"watches_you"`
	LastVisit  string    `json:"lastvisit,omitempty"` // TODO: Parse time.
	Watch      UserWatch `json:"watch"`
}

type UserWatch struct {
	Friend       bool `json:"friend" url:"friend"`
	Deviations   bool `json:"deviations" url:"deviations"`
	Journals     bool `json:"journals" url:"journals"`
	ForumThreads bool `json:"forum_threads" url:"forum_threads"`
	Critiques    bool `json:"critiques" url:"critiques"`
	Scraps       bool `json:"scraps" url:"scraps"`
	Activity     bool `json:"activity" url:"activity"`
	Collections  bool `json:"collections" url:"collections"`
}

// Get gets the users list of friends.
//
// To connect to this endpoint OAuth2 Access Token from the Client Credentials
// Grant, or Authorization Code Grant is required.
//
// The following scopes are required to access this resource:
//
//   - browse
func (s *friendsService) Get(username string, page *OffsetParams) (OffsetResponse[Friend], error) {
	var (
		success OffsetResponse[Friend]
		failure Error
	)
	_, err := s.sling.New().Get(username).QueryStruct(page).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return OffsetResponse[Friend]{}, fmt.Errorf("unable to get friends: %w", err)
	}
	return success, nil
}

type FriendsSearchParams struct {
	Username string `url:"username,omitempty"`
	Search   string `url:"search,omitempty"`
	Query    string `url:"query,omitempty"`
}

// Search searches friends by username.
//
// To connect to this endpoint OAuth2 Access Token from the Client Credentials
// Grant, or Authorization Code Grant is required.
//
// The following scopes are required to access this resource:
//
//   - browse
func (s *friendsService) Search(params *FriendsSearchParams) ([]User, error) {
	var (
		success map[string]any
		failure Error
	)
	_, err := s.sling.New().Get("search").QueryStruct(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return nil, fmt.Errorf("unable to search friends: %w", err)
	}
	return success["results"].([]User), nil
}

// Watch watches a user.
//
// Requires Authorization Code grant with the following scopes access this resource:
//
//   - browse
//   - user.manage
//
// TODO: Make better comments for access types and scopes (make it in a single paragraph/sentence).
func (s *friendsService) Watch(username string, params *UserWatch) (bool, error) {
	type watch struct {
		Watch UserWatch `url:"watch"`
	}
	var (
		success map[string]any
		failure Error
	)
	_, err := s.sling.New().Post("watch/").Path(username).BodyForm(&watch{Watch: *params}).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return false, fmt.Errorf("unable to watch a user: %w", err)
	}
	return success["success"].(bool), nil
}

// Unwatch unwatches a user.
//
// Requires Authorization Code grant with the following scopes access this resource:
//
//   - browse
//   - user.manage
//
// TODO: Make better comments for access types and scopes (make it in a single paragraph/sentence).
func (s *friendsService) Unwatch(username string) (bool, error) {
	var (
		success map[string]any
		failure Error
	)
	_, err := s.sling.New().Get("unwatch/").Path(username).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return false, fmt.Errorf("unable to unwatch a user: %w", err)
	}
	return success["success"].(bool), nil
}

// Watching checks if user is being watched by the given user.
//
// Requires Authorization Code grant with the following scopes access this resource:
//
//   - browse
//   - user
//
// TODO: Make better comments for access types and scopes (make it in a single paragraph/sentence).
func (s *friendsService) Watching(username string) (bool, error) {
	var (
		success map[string]any
		failure Error
	)
	_, err := s.sling.New().Get("watching/").Path(username).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return false, fmt.Errorf("unable to check watching for a user: %w", err)
	}
	return success["watching"].(bool), nil
}
