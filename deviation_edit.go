package deviantart

import (
	"fmt"

	"github.com/google/uuid"
)

type EditDeviationParams struct {
	// Title.
	Title string `url:"title,omitempty"`

	// Submission is mature or not.
	IsMature bool `url:"is_mature"`

	// The mature level of the submission, required for mature submissions.
	MatureLevel string `url:"mature_level,omitempty"`

	// The mature classification of the submission.
	MatureClassification []string `url:"mature_classification,brackets,omitempty"`

	// Allow comments on the submission. Default: true.
	AllowComments bool `url:"allow_comments,omitempty"`

	// License options.
	LicenseOptions LicenseOptions `url:"license_options,omitempty"`

	// UUIDs of gallery folders to publish this submission to.
	GalleryIDs []string `url:"galleryids,omitempty"`

	// Offer original file as a free download.
	AllowFreeDownload bool `url:"allow_free_download,omitempty"`

	// Add watermark. Available only if display_resolution is present.
	AddWatermark bool `url:"add_watermark,omitempty"`
}

func (s *deviationService) Edit(deviationID uuid.UUID, params *EditDeviationParams) (DeviationUpdateResponse, error) {
	var (
		success DeviationUpdateResponse
		failure Error
	)
	_, err := s.sling.New().Path("edit/").Post(deviationID.String()).BodyForm(params).Receive(success, failure)
	if err := relevantError(err, failure); err != nil {
		return DeviationUpdateResponse{}, fmt.Errorf("unable to edit deviation: %w", err)
	}
	return success, nil
}
