package deviantart

import (
	"errors"
	"fmt"

	"github.com/dghubble/sling"
	"github.com/google/uuid"
)

// TODO: COMMENTS
// 	/{commentid}/siblings
// 	/deviation/{deviationid}
// 	/post/deviation/{deviationid}
// 	/post/profile/{username}
// 	/post/status/{statusid}
// 	/profile/{username}
// 	/status/{statusid}

type commentsService struct {
	sling *sling.Sling
}

func newCommentsService(sling *sling.Sling) *commentsService {
	return &commentsService{
		sling: sling.Path("comments/"),
	}
}

type Comment struct {
	CommentID   uuid.UUID `json:"commentid"`
	ParentID    uuid.UUID `json:"parentid"`
	Posted      string    `json:"posted"`
	Replies     int       `json:"replies"`
	Body        string    `json:"body"`
	IsLiked     bool      `json:"is_liked"`
	IsFeatured  bool      `json:"is_featured"`
	Likes       int       `json:"likes"`
	User        User      `json:"user,omitempty"`
	TextContent any       `json:"text_content,omitempty"`

	// The hidden field will be null when the comment is not hidden and one of
	// the following values when it is:
	// - `hidden_by_owner` - The comment was hidden by the owner of the item
	// - `hidden_by_admin` - The comment was hidden by an administrator
	// - `hidden_by_commenter` - The comment was by the comment owner
	// - `hidden_as_spam` - The comment was hidden because it was marked spam
	Hidden string `json:"hidden"`
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
		ItemStatus    any        `json:"item_status,omitempty"` // TODO: Status object.
	} `json:"context,omitempty"`
}

// CommentSiblings fetches siblings of a comment.
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

type DeviationCommentsParams struct {
	// The commentid you want to fetch.
	CommentID uuid.UUID `url:"commentid,omitempty"`

	// Depth to query replies until.
	MaxDepth uint8 `json:"maxdepth,omitempty"`

	// The pagination offset.
	Offset int `url:"offset,omitempty"`

	// The pagination limit.
	Limit int `url:"limit,omitempty"`
}

type DeviationComments struct {
	HasMore    bool      `json:"has_more,omitempty"`
	NextOffset int       `json:"next_offset,omitempty"`
	HasLess    bool      `json:"has_less,omitempty"`
	PrevOffset bool      `json:"prev_offset,omitempty"`
	Total      int       `json:"total,omitempty"`
	Thread     []Comment `json:"thread,omitempty"`
}

// DeviationComments fetch comments posted on deviation.
func (s *commentsService) DeviationComments(deviationID uuid.UUID, params *DeviationCommentsParams) (DeviationComments, error) {
	var (
		success DeviationComments
		failure Error
	)
	_, err := s.sling.New().Get("deviation/").Path(deviationID.String()).QueryStruct(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return DeviationComments{}, fmt.Errorf("unable to fetch deviant comments: %w", err)
	}
	return success, nil
}

func (s *commentsService) PostDeviationComments(deviationID uuid.UUID) error {
	// TODO: Implement /comments/post/deviation/{deviationid}.
	return errors.New("not implemented yet")
}
