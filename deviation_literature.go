package deviantart

import (
	"fmt"

	"github.com/google/uuid"
)

type CreateLiteratureParams struct {
	// Literature title.
	Title string `url:"title"`

	// The `body` of the literature.
	Body string `url:"body,omitempty"`

	// Literature description.
	Description string `url:"description,omitempty"`

	// Literature tags.
	Tags []string `url:"tags,omitempty"`

	// UUIDs of gallery folders to publish this submission to.
	GalleryIDs []uuid.UUID `url:"galleryids,omitempty"`

	// Submission is mature or not.
	IsMature bool `url:"is_mature"`

	// The mature level of the submission, required for mature submissions.
	MatureLevel string `url:"mature_level,omitempty"`

	// The mature classification of the submission.
	MatureClassification []string `url:"mature_classification,brackets,omitempty"`

	// Allow comments on the submission
	AllowComments bool `url:"allow_comments,omitempty"`

	// License options.
	LicenseOptions []LicenseOptions `url:"license_options"`

	// ID of the embeded deviation.
	EmbeddedImageDeviationID string `url:"embedded_image_deviation_id,omitempty"`
}

// CreateLiterature creates literature.
//
// To connect to this endpoint OAuth2 Access Token from the Authorization Code
// Grant is required.
//
// The following scopes are required to access this resource:
//
//   - user.manage
func (s *DeviationService) CreateLiterature(params *CreateLiteratureParams) (uuid.UUID, error) {
	var (
		success map[string]uuid.UUID
		failure Error
	)
	_, err := s.sling.New().Post("literature/create/").BodyForm(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return uuid.UUID{}, fmt.Errorf("unable to create literature: %w", err)
	}
	return success["deviationid"], nil
}

type UpdateLiteratureParams struct {
	// Literature title.
	Title string `url:"title"`

	// Literature tags.
	Tags []string `url:"tags,omitempty"`

	// UUIDs of gallery folders to publish this submission to.
	GalleryIDs []uuid.UUID `url:"galleryids,omitempty"`

	// Submission is mature or not.
	IsMature bool `url:"is_mature"`

	// The mature level of the submission, required for mature submissions.
	MatureLevel string `url:"mature_level,omitempty"`

	// The mature classification of the submission.
	MatureClassification []string `url:"mature_classification,brackets,omitempty"`

	// Allow comments on the submission
	AllowComments bool `url:"allow_comments,omitempty"`

	// License options.
	LicenseOptions []LicenseOptions `url:"license_options"`
}

// UpdateLiterature updates literature. Note: null/empty values will have the
// corresponding fields cleared. To keep a field value send the old one.
//
// To connect to this endpoint OAuth2 Access Token from the Authorization Code
// Grant is required.
//
// The following scopes are required to access this resource:
//
//   - user.manage
func (s *DeviationService) UpdateLiterature(deviationID uuid.UUID, params *UpdateLiteratureParams) (DeviationUpdateResponse, error) {
	var (
		success DeviationUpdateResponse
		failure Error
	)
	_, err := s.sling.New().Post("literature/update/").Path(deviationID.String()).BodyForm(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return DeviationUpdateResponse{}, fmt.Errorf("unable to update literature: %w", err)
	}
	return success, nil
}
