package deviantart

import (
	"fmt"

	"github.com/google/uuid"
)

type Status struct {
	StatusID      uuid.UUID `json:"statusid,omitempty"`
	Body          string    `json:"html,omitempty"`
	Timestamp     string    `json:"ts,omitempty"` // TODO: Parse to time.Time.
	URL           string    `json:"url,omitempty"`
	CommentsCount int       `json:"comments_count,omitempty"`
	IsShare       bool      `json:"is_share,omitempty"`
	IsDeleted     bool      `json:"is_deleted,omitempty"`
	Author        *User     `json:"user,omitempty"`
	Items         []struct {
		Type      string     `json:"type"`
		Status    *Status    `json:"status,omitempty"`
		Deviation *Deviation `json:"deviation,omitempty"`
	} `json:"items,omitempty"`
	TextContent *struct {
		Excerpt string `json:"excerpt"`
		Body    struct {
			Type     string `json:"type"`
			Markup   string `json:"markup,omitempty"`
			Features string `json:"features"`
		} `json:"body"`
	} `json:"text_content,omitempty"`
}

// Status fetches the status.
//
// To connect to this endpoint OAuth2 Access Token from the Client Credentials
// Grant, or Authorization Code Grant is required.
//
// The following scopes are required to access this resource:
//
//   - browse
func (s *userService) Status(statusID uuid.UUID) (Status, error) {
	var (
		success Status
		failure Error
	)
	_, err := s.sling.New().Get("statuses/").Path(statusID.String()).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return Status{}, fmt.Errorf("unable to fetch status: %w", err)
	}
	return success, nil
}

// Statuses fetches user statuses.
//
// To connect to this endpoint OAuth2 Access Token from the Client Credentials
// Grant, or Authorization Code Grant is required.
//
// The following scopes are required to access this resource:
//
//   - browse
func (s *userService) Statuses(username string, page *OffsetParams) (OffsetResponse[Status], error) {
	var (
		success OffsetResponse[Status]
		failure Error
	)
	params := &usernameParams{Username: username}
	_, err := s.sling.New().Get("statuses").QueryStruct(params).QueryStruct(page).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return OffsetResponse[Status]{}, fmt.Errorf("unable to fetch user statuses: %w", err)
	}
	return success, nil
}

type PostStatusParams struct {
	// The body of the status.
	Text string `url:"body,omitempty"`

	// The ID of the object you wish to share.
	ID uuid.UUID `url:"id,omitempty"`

	ParentID uuid.UUID
}

// PostStatus postes a status.
//
// When posting a status, it is possible to share another status or deviation
// along with the status text. To do so, pass UUID of an object being shared
// (status or deviation) in `id` parameter. Note that it is not possible to
// share a status which already shares something. Sometimes the object you want
// to share is contained within the status. To share such object pass UUID of
// the containing status in `parentid` parameter (in addition to the `id`).
//
// To connect to this endpoint OAuth2 Access Token from the Authorization Code
// Grant is required.
//
// The following scopes are required to access this resource:
//
//   - browse
//   - user.manage
func (s *userService) PostStatus(params *PostStatusParams) (uuid.UUID, error) {
	type response struct {
		StatusID uuid.UUID `json:"statusid"`
	}
	var (
		success response
		failure Error
	)
	_, err := s.sling.New().Post("statuses/post").BodyForm(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return uuid.Nil, fmt.Errorf("unable to post status: %w", err)
	}
	return success.StatusID, nil
}
