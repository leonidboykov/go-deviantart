package deviantart

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	MatureLevelStrict   = "strict"
	MatureLevelModerate = "moderate"
)

const (
	MatureClassificationNudity   = "nudity"
	MatureClassificationSexual   = "sexual"
	MatureClassificationGore     = "gore"
	MatureClassificationLanguage = "language"
	MatureClassificationIdeology = "ideology"
)

type DisplayResolution uint8

const (
	DisplayResolutionOriginal DisplayResolution = iota
	DisplayResolution400px
	DisplayResolution600px
	DisplayResolution800px
	DisplayResolution900px
	DisplayResolution1024px
	DisplayResolution1280px
	DisplayResolution1600px
	DisplayResolution1920px
)

const (
	SharingOptionsAllow              = "allow"
	SharingOptionsHideShareButtons   = "hide_share_buttons"
	SharingOptionsHideAndMembersOnly = "hide_and_members_only"
)

type StashPublishParams struct {
	// The ID of the stash item to publish.
	ItemID int64 `url:"itemid"`

	// Submission is mature or not.
	IsMature bool `url:"is_mature"`

	// The mature level of the submission, required for mature submissions.
	MatureLevel string `url:"mature_level,omitempty"`

	// The mature classification of the submission.
	MatureClassification []string `url:"mature_classification,brackets,omitempty"`

	// Agree to submission policy.
	AgreeSubmission bool `url:"agree_submission"`

	// Agree to terms of service.
	AgreeToS bool `url:"agree_tos"`

	// Feature the submission. Default: true.
	Feature bool `url:"feature,omitempty"`

	// Allow comments on the submission. Default: true.
	AllowComments bool `url:"allow_comments,omitempty"`

	// Request a critique, only available to some users see `/publish/userdata`.
	RequestCritique bool `url:"request_critique,omitempty"`

	// Resize image to. Cannot be bigger than image and is ignored for non images.
	DisplayResolution DisplayResolution `url:"display_resolution,omitempty"`

	// Sharing options.
	SharingOptions string `url:"sharing,omitempty"`

	// License options.
	LicenseOptions LicenseOptions `url:"license_options,omitempty"`

	// UUIDs of gallery folders to publish this submission to.
	GalleryIDs []string `url:"galleryids,omitempty"`

	// Offer original file as a free download.
	AllowFreeDownload bool `url:"allow_free_download,omitempty"`

	// Add watermark. Available only if display_resolution is present.
	AddWatermark bool `url:"add_watermark,omitempty"`
}

type StashPublishResponse struct {
	StatusResponse
	URL         string    `json:"url"`
	DeviationID uuid.UUID `json:"deviationid"`
}

// Publish a Sta.sh item to deviantART.
//
// To connect to this endpoint OAuth2 Access Token from the Authorization Code
// Grant is required.
//
// The following scopes are required to access this resource:
//
//   - stash
//   - publish
func (s *stashService) Publish(params StashPublishParams) (StashPublishResponse, error) {
	var (
		success StashPublishResponse
		failure Error
	)
	_, err := s.sling.New().Post("publish").QueryStruct(&params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return StashPublishResponse{}, fmt.Errorf("unable to publish item: %w", err)
	}
	return success, nil
}
