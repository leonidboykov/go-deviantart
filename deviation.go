package deviantart

import (
	"fmt"

	"github.com/dghubble/sling"
	"github.com/google/uuid"
)

type deviationService struct {
	sling *sling.Sling
}

func newDeviationService(sling *sling.Sling) *deviationService {
	return &deviationService{
		sling: sling.Path("deviation/"),
	}
}

// deviationIDParam is a wrapper for single deviation ID.
type deviationIDParam struct {
	DeviationID uuid.UUID `url:"deviationid"`
}

type Deviation struct {
	DeviationID  uuid.UUID `json:"deviationid"`
	PrintID      uuid.UUID `json:"printid,omitempty"`
	URL          string    `json:"url,omitempty"`
	Title        string    `json:"title,omitempty"`
	Category     string    `json:"category,omitempty"`
	CategoryPath string    `json:"category_path,omitempty"`
	IsFavourited bool      `json:"is_favourited,omitempty"`
	IsDeleted    bool      `json:"is_deleted"`
	IsPublished  bool      `json:"is_published,omitempty"`
	IsBlocked    bool      `json:"is_blocked,omitempty"`
	Author       User      `json:"author,omitempty"`
	Stats        struct {
		Comments   uint32 `json:"comments"`
		Favourites uint32 `json:"favourites"`
	} `json:"stats,omitempty"`
	PublishedTime  string        `json:"published_time,omitempty"`
	AllowsComments bool          `json:"allows_comments,omitempty"`
	Tier           DeviationTier `json:"tier,omitempty"`
	Preview        StashFile     `json:"preview,omitempty"`
	Content        struct {
		StashFile
		FileSize uint32 `json:"filesize"`
	} `json:"content,omitempty"`
	Thumbs []StashFile `json:"thumbs"`
	Videos []struct {
		Src      string `json:"src"`
		Quality  string `json:"quality"`
		FileSize uint32 `json:"filesize"`
		Duration uint32 `json:"duration"`
	} `json:"videos,omitempty"`
	Flash          []fileInfo `json:"flash,omitempty"`
	DailyDeviation struct {
		Body      string `json:"body"`
		Time      string `json:"time"`
		Giver     User   `json:"giver"`
		Suggester User   `json:"suggester,omitempty"`
	} `json:"daily_deviation,omitempty"`
	PremiumFolderData *PremiumFolderData `json:"premium_folder_data,omitempty"`
	TextContent       *EditorText        `json:"text_content,omitempty"`
	IsPinned          bool               `json:"is_pinned,omitempty"`
	CoverImage        *Deviation         `json:"cover_image,omitempty"`
	TierAccess        string             `json:"tier_access,omitempty"`
	PrimaryTier       *Deviation         `json:"primary_tier,omitempty"`
	Excerpt           string             `json:"excerpt,omitempty"`
	IsMature          bool               `json:"is_mature,omitempty"`
	IsDownloadable    bool               `json:"is_downloadable,omitempty"`
	DownloadFileSize  uint32             `json:"download_filesize,omitempty"`
	MotionBook        struct {
		EmbedURL string `json:"embed_url,omitempty"`
	} `json:"motion_book,omitempty"`
	SuggestedReasons []any `json:"suggested_reasons,omitempty"`
}

type PremiumFolderData struct {
	Type           string    `json:"type"`
	HasAccess      bool      `json:"has_access"`
	GalleryID      uuid.UUID `json:"gallery_id"`
	PointsPrice    int       `json:"points_price,omitempty"`
	DollarPrice    float64   `json:"dollar_price,omitempty"` // TODO: DeviationTier has the same string field.
	NumSubscribers int       `json:"num_subscribers,omitempty"`
	SubproductID   int       `json:"subproductid,omitempty"` // TODO: Is it really an integer field and not an UUID?
}

type DeviationTier struct {
	State            string `json:"state,omitempty"` // TODO: enum[draft,active,pending_deletion,deleted]
	IsUserSubscribed bool   `json:"is_user_subscribed,omitempty"`
	CanUserSubscribe bool   `json:"can_user_subscribe,omitempty"`
	SubproductID     uint64 `json:"subproductid,omitempty"`
	DollarPrice      string `json:"dollar_price,omitempty"` // TODO: PremiumFolderData has the same float field.
	Settings         struct {
		AccessSettings string `json:"access_settings"` // TODO: enum[all,future_only,limited_past_and_future]
	} `json:"settings,omitempty"`
	Stats struct {
		Subscribers uint32 `json:"subscribers,omitempty"`
		Deviations  uint32 `json:"deviations,omitempty"`
		Posts       uint32 `json:"posts,omitempty"`
		Total       uint32 `json:"total,omitempty"`
	} `json:"stats"`
	Benefits []string `json:"benefits"`
}

type DeviationUpdateResponse struct {
	StatusResponse
	URL         string    `json:"url"`
	DeviationID uuid.UUID `json:"deviationid"`
}

// Deviation fetches a deviation.
//
// To connect to this endpoint OAuth2 Access Token from the Client Credentials
// Grant, or Authorization Code Grant is required.
//
// The following scopes are required to access this resource:
//
//   - browse
func (s *deviationService) Deviation(deviationID uuid.UUID) (Deviation, error) {
	var (
		success Deviation
		failure Error
	)
	_, err := s.sling.New().Get(deviationID.String()).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return Deviation{}, fmt.Errorf("unable to fetch deviation: %w", err)
	}
	return success, nil
}

type Content struct {
	HTML     string   `json:"html,omitempty"`
	CSS      string   `json:"css,omitempty"`
	CSSFonts []string `json:"css_fonts,omitempty"`
}

// Content fetches a full data that is not included in the main deviation
// object.
//
// The endpoint works with journals and literatures. Deviation objects returned
// from API contain only excerpt of a journal, use this endpoint to load full
// content. Any custom CSS rules and fonts applied to journal are also returned.
//
// To connect to this endpoint OAuth2 Access Token from the Client Credentials
// Grant, or Authorization Code Grant is required.
//
// The following scopes are required to access this resource:
//
//   - browse
func (s *deviationService) Content(deviationID uuid.UUID) (Content, error) {
	var (
		success Content
		failure Error
	)
	params := &deviationIDParam{DeviationID: deviationID}
	_, err := s.sling.New().Get("content/").QueryStruct(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return Content{}, fmt.Errorf("unable to fetch deviation content: %w", err)
	}
	return success, nil
}

type DownloadResponse struct {
	fileInfo
	FileName string `json:"filename"`
	FileSize uint32 `json:"filesize"`
}

// Download fetches the original file download (if allowed).
//
// To connect to this endpoint OAuth2 Access Token from the Client Credentials
// Grant, or Authorization Code Grant is required.
//
// The following scopes are required to access this resource:
//
//   - browse
func (s *deviationService) Download(deviationID uuid.UUID) (DownloadResponse, error) {
	var (
		success DownloadResponse
		failure Error
	)
	_, err := s.sling.New().Get(deviationID.String()).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return DownloadResponse{}, fmt.Errorf("unable to fetch download data: %w", err)
	}
	return success, nil
}

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

// Edit edits deviation. Note: null/empty values will have the corresponding
// fields cleared. To keep a field value send the old one.
//
// To connect to this endpoint OAuth2 Access Token from the Authorization Code
// Grant is required.
//
// The following scopes are required to access this resource:
//
//   - stash
//   - publish
func (s *deviationService) Edit(deviationID uuid.UUID, params *EditDeviationParams) (DeviationUpdateResponse, error) {
	var (
		success DeviationUpdateResponse
		failure Error
	)
	_, err := s.sling.New().Get("edit/").Path(deviationID.String()).BodyForm(params).Receive(success, failure)
	if err := relevantError(err, failure); err != nil {
		return DeviationUpdateResponse{}, fmt.Errorf("unable to edit deviation: %w", err)
	}
	return success, nil
}

type EmbeddedContentParams struct {
	// The deviation ID of container deviation.
	DeviationID uuid.UUID `url:"deviationid"`

	// ID of embedded deviation to use as an offset.
	OffsetDeviationID uuid.UUID `url:"offset_deviationid,omitempty"`
}

// EmbeddedContent fetch a content embedded in a deviation.
//
// Journal and literature deviations support embedding of deviations inside
// them.
//
// To connect to this endpoint OAuth2 Access Token from the Client Credentials
// Grant, or Authorization Code Grant is required.
//
// The following scopes are required to access this resource:
//
//   - browse
func (s *deviationService) EmbeddedContent(params *EmbeddedContentParams, page *OffsetParams) (OffsetResponse[Deviation], error) {
	var (
		success OffsetResponse[Deviation]
		failure Error
	)
	_, err := s.sling.New().Get("embeddedcontent/").QueryStruct(params).QueryStruct(page).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return OffsetResponse[Deviation]{}, fmt.Errorf("fetch content embedded in a deviation: %w", err)
	}
	return success, nil
}

type DeviationMetadata struct {
	DeviationID          uuid.UUID            `json:"deviationid"`
	PrintID              uuid.UUID            `json:"uuid,omitempty"`
	Author               *User                `json:"author"`
	IsWatching           bool                 `json:"is_watching"`
	Title                string               `json:"title"`
	Description          string               `json:"description"`
	License              string               `json:"license"`
	AllowsComments       bool                 `json:"allow_comments"`
	Tags                 []DeviationTag       `json:"tags"`
	IsFavourited         bool                 `json:"is_favourited"`
	IsMature             bool                 `json:"is_mature"`
	MatureLevel          string               `json:"mature_level,omitempty"`
	MatureClassification []string             `json:"mature_classification,omitempty"`
	Submission           *DeviationSubmission `json:"submission,omitempty"`
	Stats                *DeviationStats      `json:"stats,omitempty"`
	Camera               map[string]string    `json:"camera,omitempty"`
	Collections          []Folder             `json:"collections,omitempty"`
	Galleries            []Folder             `json:"galleries,omitempty"`
	CanPostComments      bool                 `json:"can_post_comments,omitempty"`
}

type DeviationTag struct {
	Name      string `json:"tag_name"`
	Sponsored bool   `json:"sponsored"`
	Sponsor   bool   `json:"sponsor"`
}

type DeviationSubmission struct {
	CreationTime  string `json:"creation_time"`
	Category      string `json:"category"`
	FileSize      string `json:"file_size,omitempty"`
	Resolution    string `json:"resolution,omitempty"`
	SubmittedWith struct {
		App string `json:"app"`
		URL string `json:"url"`
	} `json:"submitted_with"`
}

type DeviationStats struct {
	Views          int `json:"views"`
	ViewsToday     int `json:"views_today,omitempty"`
	Favourites     int `json:"favourites"`
	Comments       int `json:"comments"`
	Downloads      int `json:"downloads"`
	DownloadsToday int `json:"downloads_today,omitempty"`
}

type MetadataResponse struct {
	Metatada []DeviationMetadata `json:"metadata"`
}

type MetadataParams struct {
	// The deviation IDs you want metadata for.
	DeviationIDs []uuid.UUID `url:"deviationids"`

	IncludeSubmission bool `url:"ext_submission,omitempty"`
	IncludeCamera     bool `url:"ext_camera,omitempty"`
	IncludeStats      bool `url:"ext_stats,omitempty"`
	IncludeCollection bool `url:"ext_collection,omitempty"`
	IncludeGallery    bool `url:"ext_gallery,omitempty"`
}

// Metadata fetches a deviation metadata for a set of deviations.
//
// This endpoint is limited to 50 deviations per query when fetching the base
// data and 10 when fetching extended data.
//
// To connect to this endpoint OAuth2 Access Token from the Client Credentials
// Grant, or Authorization Code Grant is required.
//
// The following scopes are required to access this resource:
//
//   - browse
func (s *deviationService) Metadata(params *MetadataParams) (MetadataResponse, error) {
	var (
		success MetadataResponse
		failure Error
	)
	_, err := s.sling.New().Get("metadata").QueryStruct(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return MetadataResponse{}, fmt.Errorf("unable to fetch deviation metadata: %w", err)
	}
	return success, nil
}

type FaveInfo struct {
	User *User `json:"user"`
	Time int64 `json:"time"`
}

// WhoFaved fetches a list of users who faved the deviation.
//
// To connect to this endpoint OAuth2 Access Token from the Client Credentials
// Grant, or Authorization Code Grant is required.
//
// The following scopes are required to access this resource:
//
//   - browse
func (s *deviationService) WhoFaved(deviationID uuid.UUID, page *OffsetParams) (OffsetResponse[FaveInfo], error) {
	var (
		success OffsetResponse[FaveInfo]
		failure Error
	)
	params := &deviationIDParam{DeviationID: deviationID}
	_, err := s.sling.New().Get("whofaved").QueryStruct(params).QueryStruct(page).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return OffsetResponse[FaveInfo]{}, fmt.Errorf("unable to fetch whofaved a deviation: %w", err)
	}
	return success, nil
}
