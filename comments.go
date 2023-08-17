package deviantart

import (
	"fmt"

	"github.com/dghubble/sling"
	"github.com/google/uuid"
)

type commentsService struct {
	sling *sling.Sling
}

func newCommentsService(sling *sling.Sling) *commentsService {
	return &commentsService{
		sling: sling.Path("comments/"),
	}
}

const (
	HiddenByOwner     = "hidden_by_owner"
	HiddenByAdmin     = "hidden_by_admin"
	HiddenByCommenter = "hidden_by_commenter"
	HiddenAsSpam      = "hidden_as_spam"
)

type Comment struct {
	CommentID   uuid.UUID   `json:"commentid"`
	ParentID    uuid.UUID   `json:"parentid"`
	Posted      string      `json:"posted"`
	Replies     int         `json:"replies"`
	Body        string      `json:"body"`
	IsLiked     bool        `json:"is_liked"`
	IsFeatured  bool        `json:"is_featured"`
	Likes       int         `json:"likes"`
	User        User        `json:"user,omitempty"`
	TextContent *EditorText `json:"text_content,omitempty"`

	// The hidden field will be null when the comment is not hidden and one of
	// the following values when it is:
	//   - `hidden_by_owner` - The comment was hidden by the owner of the item
	//   - `hidden_by_admin` - The comment was hidden by an administrator
	//   - `hidden_by_commenter` - The comment was by the comment owner
	//   - `hidden_as_spam` - The comment was hidden because it was marked spam
	Hidden string `json:"hidden"`
}

type EditorText struct {
	Excerpt string `json:"excerpt"`
	Body    struct {
		Type     string `json:"type"`
		Markup   string `json:"markup,omitempty"`
		Features string `json:"features"`
	} `json:"body"`
}

type CommentSiblingsParams struct {
	// Fetch the related containing item (deviation, profile user, or status).
	IncludeItem bool `url:"ext_item,omitempty"`

	// The pagination offset.
	Offset int `url:"offset,omitempty"`

	// The pagination limit.
	Limit int `url:"offset,omitempty"`
}

type CommentSiblings struct {
	HasMore    bool      `json:"has_more"`
	NextOffset int       `json:"next_offset,omitempty"`
	HasLess    bool      `json:"has_less"`
	PrevOffset int       `json:"prev_offset,omitempty"`
	Thread     []Comment `json:"thread"`
	Context    struct {
		Parent        *Comment   `json:"parent,omitempty"`
		ItemProfile   *User      `json:"item_profile,omitempty"`
		ItemDeviation *Deviation `json:"item_deviation,omitempty"`
		ItemStatus    *Status    `json:"item_status,omitempty"`
	} `json:"context,omitempty"`
}

// CommentSiblings fetches siblings of a comment.
//
// To connect to this endpoint, OAuth2 Access Token, from the Client Credentials
// Grant, or Authorization Code Grant is required.
//
// The following scopes are required to access this resource:
//
//   - browse
func (s *commentsService) CommentSiblings(commentID uuid.UUID, params *CommentSiblingsParams) (CommentSiblings, error) {
	var (
		success CommentSiblings
		failure Error
	)
	_, err := s.sling.New().Get(commentID.String()+"/").Path("siblings").Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return CommentSiblings{}, fmt.Errorf("unable to fetch comment siblings: %w", err)
	}
	return success, nil
}

type FetchCommentsParams struct {
	// The commentid you want to fetch.
	CommentID uuid.UUID `url:"commentid,omitempty"`

	// Depth to query replies until.
	MaxDepth uint8 `json:"maxdepth,omitempty"`

	// The pagination offset.
	Offset int `url:"offset,omitempty"`

	// The pagination limit.
	Limit int `url:"limit,omitempty"`
}

type CommentsResponse struct {
	HasMore    bool      `json:"has_more,omitempty"`
	NextOffset int       `json:"next_offset,omitempty"`
	HasLess    bool      `json:"has_less,omitempty"`
	PrevOffset bool      `json:"prev_offset,omitempty"`
	Total      int       `json:"total,omitempty"`
	Thread     []Comment `json:"thread,omitempty"`
}

// DeviationComments fetch comments posted on deviation.
//
// To connect to this endpoint, OAuth2 Access Token, from the Client Credentials
// Grant, or Authorization Code Grant is required.
//
// The following scopes are required to access this resource:
//
//   - browse
func (s *commentsService) DeviationComments(deviationID uuid.UUID, params *FetchCommentsParams) (CommentsResponse, error) {
	var (
		success CommentsResponse
		failure Error
	)
	_, err := s.sling.New().Get("deviation/").Path(deviationID.String()).QueryStruct(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return CommentsResponse{}, fmt.Errorf("unable to fetch deviation comments: %w", err)
	}
	return success, nil
}

// ProfileComments fetch comments posted on user profile.
//
// To connect to this endpoint, OAuth2 Access Token, from the Client Credentials
// Grant, or Authorization Code Grant is required.
//
// The following scopes are required to access this resource:
//
//   - browse
func (s *commentsService) ProfileComments(username string, params *FetchCommentsParams) (CommentsResponse, error) {
	var (
		success CommentsResponse
		failure Error
	)
	_, err := s.sling.New().Get("profile/").Path(username).QueryStruct(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return CommentsResponse{}, fmt.Errorf("unable to fetch profile comments: %w", err)
	}
	return success, nil
}

// StatusComments fetch comments posted on status.
//
// To connect to this endpoint, OAuth2 Access Token, from the Client Credentials
// Grant, or Authorization Code Grant is required.
//
// The following scopes are required to access this resource:
//
//   - browse
func (s *commentsService) StatusComments(statusID uuid.UUID, params *FetchCommentsParams) (CommentsResponse, error) {
	var (
		success CommentsResponse
		failure Error
	)
	_, err := s.sling.New().Get("status/").Path(statusID.String()).QueryStruct(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return CommentsResponse{}, fmt.Errorf("unable to fetch status comments: %w", err)
	}
	return success, nil
}

type CommentParams struct {
	// The Comment ID you are replying to.
	CommentID uuid.UUID `url:"commentid,omitempty"`

	// The comment text.
	Text string `url:"body"`
}

// CommentDeviation posts a comment on a deviation.
//
// To connect to this endpoint, OAuth2 Access Token, from the Authorization Code
// Grant is required.
//
// The following scopes are required to access this resource:
//
//   - browse
//   - comment.post
func (s *commentsService) CommentDeviation(deviationID uuid.UUID, params *CommentParams) (Comment, error) {
	var (
		success Comment
		failure Error
	)
	_, err := s.sling.New().Post("post/deviation/").Path(deviationID.String()).BodyForm(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return Comment{}, fmt.Errorf("unable to comment a deviation: %w", err)
	}
	return success, nil
}

// CommentProfile posts a comment on a users profile.
//
// To connect to this endpoint, OAuth2 Access Token, from the Authorization Code
// Grant is required.
//
// The following scopes are required to access this resource:
//
//   - browse
//   - comment.post
func (s *commentsService) CommentProfile(username string, params *CommentParams) (Comment, error) {
	var (
		success Comment
		failure Error
	)
	_, err := s.sling.New().Post("post/profile/").Path(username).BodyForm(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return Comment{}, fmt.Errorf("unable to comment a profile: %w", err)
	}
	return success, nil
}

// CommentStatus posts a comment on a status.
//
// To connect to this endpoint, OAuth2 Access Token, from the Authorization Code
// Grant is required.
//
// The following scopes are required to access this resource:
//
//   - browse
//   - comment.post
func (s *commentsService) CommentStatus(statusID uuid.UUID, params *CommentParams) (Comment, error) {
	var (
		success Comment
		failure Error
	)
	_, err := s.sling.New().Post("post/status/").Path(statusID.String()).BodyForm(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return Comment{}, fmt.Errorf("unable to comment a status: %w", err)
	}
	return success, nil
}
