package deviantart

import (
	"fmt"

	"github.com/google/uuid"
)

type CreateJournalParams struct {
	// Journal title.
	Title string `url:"title"`

	// The `body` of the journal.
	Body string `url:"body,omitempty"`

	// Journal tags.
	Tags []string `url:"tags,brackets,omitempty"`

	// Cover deviation ID.
	CoverImageDeviationID uuid.UUID `url:"cover_image_deviation_id,omitempty"`

	// ID of the embeded deviation.
	EmbeddedImageDeviationID uuid.UUID `url:"embedded_image_deviation_id,omitempty"`

	// Submission is mature or not.
	IsMature bool `url:"is_mature,omitempty"`

	// Allow comments on the submission.
	AllowComments bool `url:"allow_comments,omitempty"`

	// License options.
	LicenseOptions LicenseOptions `url:"license_options,omitempty"`
}

// CreateJournal creates journal.
//
// To connect to this endpoint OAuth2 Access Token from the Authorization Code
// Grant is required.
//
// The following scopes are required to access this resource:
//
//   - user.manage
func (s *DeviationService) CreateJournal(params *CreateJournalParams) (uuid.UUID, error) {
	var (
		success map[string]uuid.UUID
		failure Error
	)
	_, err := s.sling.New().Post("journal/create/").BodyForm(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return uuid.UUID{}, fmt.Errorf("unable to create journal: %w", err)
	}
	return success["deviationid"], nil
}

type UpdateJournalParams struct {
	// Journal title.
	Title string `url:"title"`

	// Journal tags.
	Tags []string `url:"tags,brackets,omitempty"`

	// Cover deviation ID.
	CoverImageDeviationID uuid.UUID `url:"cover_image_deviation_id,omitempty"`

	// Reset cover deviation ID.
	ResetCoverImageDeviationID bool `url:"reset_cover_image_deviation_id,omitempty"`

	// Submission is mature or not.
	IsMature bool `url:"is_mature,omitempty"`

	// Allow comments on the submission.
	AllowComments bool `url:"allow_comments,omitempty"`

	// License options.
	LicenseOptions LicenseOptions `url:"license_options,omitempty"`
}

// UpdateJournal updates journal. All values left empty, except cover image
// deviation id, will have the corresponding fields cleared. To keep a field
// value send the old one. To clear cover image deviation id value - pass
// `reset_cover_image_deviation_id` param with value `true`.
//
// To connect to this endpoint OAuth2 Access Token from the Authorization Code
// Grant is required.
//
// The following scopes are required to access this resource:
//
//   - user.manage
func (s *DeviationService) UpdateJournal(deviationID uuid.UUID, params *UpdateJournalParams) (DeviationUpdateResponse, error) {
	var (
		success DeviationUpdateResponse
		failure Error
	)
	_, err := s.sling.New().Post("journal/update/").Path(deviationID.String()).BodyForm(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return DeviationUpdateResponse{}, fmt.Errorf("unable to update journal: %w", err)
	}
	return success, nil
}
