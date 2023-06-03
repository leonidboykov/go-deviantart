package deviantart

import (
	"fmt"

	"github.com/google/uuid"
)

type UserStatus struct {
	StatusID      uuid.UUID `json:"statusid,omitempty"`
	Body          string    `json:"html,omitempty"`
	Timestamp     string    `json:"ts,omitempty"` // TODO: Parse to time.Time.
	URL           string    `json:"url,omitempty"`
	CommentsCount int       `json:"comments_count,omitempty"`
	IsShare       bool      `json:"is_share,omitempty"`
	IsDeleted     bool      `json:"is_deleted,omitempty"`
	Author        *User     `json:"user,omitempty"`
	Items         []struct {
		Type      string      `json:"type"`
		Status    *UserStatus `json:"status,omitempty"`
		Deviation *Deviation  `json:"deviation,omitempty"`
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
func (s *userService) Status(statusID uuid.UUID) (UserStatus, error) {
	var (
		success UserStatus
		failure Error
	)
	_, err := s.sling.New().Get("statuses/").Path(statusID.String()).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return UserStatus{}, fmt.Errorf("unable to fetch status: %w", err)
	}
	return success, nil
}

// Statuses fetches user statuses.
func (s *userService) Statuses(username string, page *OffsetParams) (OffsetResponse[UserStatus], error) {
	var (
		success OffsetResponse[UserStatus]
		failure Error
	)
	params := &usernameParams{Username: username}
	_, err := s.sling.New().Get("statuses").QueryStruct(params).QueryStruct(page).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return OffsetResponse[UserStatus]{}, fmt.Errorf("unable to fetch user statuses: %w", err)
	}
	return success, nil
}

type PostStatusParams struct {
	// The body of the status.
	Text string `url:"body,omitempty"`

	// The ID of the object you wish to share.
	ID uuid.UUID `url:"id,omitempty"` // TODO: Is it UUID?

	ParentID uuid.UUID
}

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
