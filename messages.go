package deviantart

import (
	"fmt"

	"github.com/dghubble/sling"
	"github.com/google/uuid"
)

type messagesService struct {
	sling *sling.Sling
}

func newMessagesService(sling *sling.Sling) *messagesService {
	return &messagesService{
		sling: sling.Path("messages/"),
	}
}

type Message struct {
	MessageID  string `json:"messageid"`
	Type       string `json:"type"`
	Orphaned   bool   `json:"orphaned"`
	TS         string `json:"ts,omitempty"`
	StackID    string `json:"stackid,omitempty"`
	StackCount int    `json:"stack_count,omitempty"`
	IsNew      bool   `json:"is_new"`
	Originator *User  `json:"originator,omitempty"`
	Subject    *struct {
		Profile    *User      `json:"profile,omitempty"`
		Deviation  *Deviation `json:"deviation,omitempty"`
		Status     any        `json:"status,omitempty"` // TODO: Status object
		Comment    *Comment   `json:"comment,omitempty"`
		Collection *Folder    `json:"collection"`
		Gallery    *Folder    `json:"gallery"`
	} `json:"subject,omitempty"`
	HTML       string     `json:"html,omitempty"`
	Profile    *User      `json:"profile,omitempty"`
	Deviation  *Deviation `json:"deviation,omitempty"`
	Status     any        `json:"status,omitempty"` // TODO: Status object
	Comment    *Comment   `json:"comment,omitempty"`
	Collection *Folder    `json:"collection,omitempty"`
}

type DeleteMessageParams struct {
	// The folder to delete the message from, defaults to inbox.
	FolderID uuid.UUID `url:"folderid,omitempty"`

	// The message to delete.
	MessageID string `url:"messageid,omitempty"`

	// The stack to delete.
	StackID string `url:"stackid,omitempty"`
}

// Delete deletes a message or a message stack.
func (s *messagesService) Delete(params *DeleteMessageParams) error {
	var (
		failure Error
	)
	_, err := s.sling.New().Post("delete").BodyForm(params).Receive(nil, &failure)
	if err := relevantError(err, failure); err != nil {
		return fmt.Errorf("unable to delete message: %w", err)
	}
	return nil
}

type MessagesFeedParams struct {
	// The folder to fetch messages from, defaults to inbox.
	FolderID uuid.UUID `url:"folderid,omitempty"`

	// True to use stacked mode, false to use flat mode.
	Stack bool `url:"stack,omitempty"`

	Cursor      string `url:"cursor,omitempty"`
	WithSession bool   `url:"with_session,omitempty"`
}

type MessagesFeed struct {
	HasMore bool      `json:"has_more"`
	Cursor  string    `json:"cursor"`
	Results []Message `json:"results"`
	Session *Session  `json:"session,omitempty"`
}

// Feed fetches feed of all messages.
//
// Messages can be fetched in a stacked (default) or flat mode. In the stacked
// mode similar messages will be grouped together and the most recent one will
// be returned. stackid can be used to fetch the rest of the messages in the
// stack.
func (s *messagesService) Feed(params *MessagesFeedParams) (MessagesFeed, error) {
	var (
		success MessagesFeed
		failure Error
	)
	_, err := s.sling.New().Get("feed").QueryStruct(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return MessagesFeed{}, fmt.Errorf("unable to fetch message feed: %w", err)
	}
	return success, nil
}

type MessagesFeedbackParams struct {
	// Type of feedback messages to fetch.
	Type string `url:"type"`

	// The folder to fetch messages from, defaults to inbox.
	FolderID uuid.UUID `jsurlon:"folderid,omitempty"`

	// True to use stacked mode, false to use flat mode.
	Stack bool `url:"stack,omitempty"`

	// The pagination offset.
	Offset int `url:"offset,omitempty"`

	// The pagination limit.
	Limit int `url:"limit,omitempty"`

	WithSession bool `url:"with_session,omitempty"`
}

// Feedback fetches feedback messages.
//
// Messages can be fetched in a stacked (default) or flat mode. In the stacked
// mode similar messages will be grouped together and the most recent one will
// be returned. stackid can be used to fetch the rest of the messages in the
// stack.
func (s *messagesService) Feedback(params *MessagesFeedbackParams) (MessagesFeed, error) {
	var (
		success MessagesFeed
		failure Error
	)
	_, err := s.sling.New().Get("feedback").QueryStruct(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return MessagesFeed{}, fmt.Errorf("unable to fetch message feedback: %w", err)
	}
	return success, nil
}

type StackFeedbackParams struct {
	// The pagination offset.
	Offset int `url:"offset,omitempty"`

	// The pagination limit.
	Limit int `url:"limit,omitempty"`

	WithSession bool `url:"with_session,omitempty"`
}

// Fetch messages in a stack.
func (s *messagesService) StackFeedback(stackID string, params *StackFeedbackParams) (MessagesFeed, error) {
	var (
		success MessagesFeed
		failure Error
	)
	_, err := s.sling.New().Get("feedback").Path(stackID).QueryStruct(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return MessagesFeed{}, fmt.Errorf("unable to fetch stack feedback: %w", err)
	}
	return success, nil
}

type MessagesMentionsParams struct {
	// The folder to fetch messages from, defaults to inbox.
	FolderID uuid.UUID `url:"folderid,omitempty"`

	// True to use stacked mode, false to use flat mode.
	Stack bool `url:"stack,omitempty"`

	// The pagination offset.
	Offset int `url:"offset,omitempty"`

	// The pagination limit.
	Limit int `url:"limit,omitempty"`

	WithSession bool `url:"with_session,omitempty"`
}

type MessagesMentions struct {
	HasMore    bool      `json:"has_more"`
	NextOffset string    `json:"cursor"`
	Results    []Message `json:"results"`
	Session    *Session  `json:"session,omitempty"`
}

// Mentions fetches mention messages.
//
// Messages can be fetched in a stacked (default) or flat mode. In the stacked
// mode similar messages will be grouped together and the most recent one will
// be returned. stackid can be used to fetch the rest of the messages in the
// stack.
func (s *messagesService) Mentions(params *MessagesMentionsParams) (MessagesMentions, error) {
	var (
		success MessagesMentions
		failure Error
	)
	_, err := s.sling.New().Get("mentions").QueryStruct(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return MessagesMentions{}, fmt.Errorf("unable to fetch message mentions: %w", err)
	}
	return success, nil
}

// StackMentions fetches messages in a stack.
func (s *messagesService) StackMentions(stackID string, params *StackFeedbackParams) (MessagesMentions, error) {
	var (
		success MessagesMentions
		failure Error
	)
	_, err := s.sling.New().Get("messages/").Path(stackID).QueryStruct(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return MessagesMentions{}, fmt.Errorf("unable to fetch stack mentions: %w", err)
	}
	return success, nil
}