package deviantart

import (
	"errors"
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
	Preview        StashObject   `json:"preview,omitempty"`
	Content        struct {
		StashObject
		FileSize uint32 `json:"filesize"`
	} `json:"content,omitempty"`
	Thumbs []StashObject `json:"thumbs"`
	Videos []struct {
		Src      string `json:"src"`
		Quality  string `json:"quality"`
		FileSize uint32 `json:"filesize"`
		Duration uint32 `json:"duration"`
	} `json:"videos,omitempty"`
	Flash []struct {
		Src    string `json:"src"`
		Height uint32 `json:"height"`
		Width  uint32 `json:"width"`
	} `json:"flash,omitempty"`
	DailyDeviation struct {
		Body      string `json:"body"`
		Time      string `json:"time"`
		Giver     User   `json:"giver"`
		Suggester User   `json:"suggester,omitempty"`
	} `json:"daily_deviation,omitempty"`
	PremiumFolderData any        `json:"premium_folder_data,omitempty"` // TODO: premium folder data.
	TextContent       any        `json:"text_content,omitempty"`        // TODO: editor object.
	IsPinned          bool       `json:"is_pinned,omitempty"`
	CoverImage        *Deviation `json:"cover_image,omitempty"`
	TierAccess        string     `json:"tier_access,omitempty"`
	PrimaryTier       *Deviation `json:"primary_tier,omitempty"`
	Excerpt           string     `json:"excerpt,omitempty"`
	IsMature          bool       `json:"is_mature,omitempty"`
	IsDownloadable    bool       `json:"is_downloadable,omitempty"`
	DownloadFileSize  uint32     `json:"download_filesize,omitempty"`
	MotionBook        struct {
		EmbedURL string `json:"embed_url,omitempty"`
	} `json:"motion_book,omitempty"`
	SuggestedReasons []any `json:"suggested_reasons,omitempty"`
}

type DeviationTier struct {
	State            string `json:"state,omitempty"` // TODO: enum[draft,active,pending_deletion,deleted]
	IsUserSubscribed bool   `json:"is_user_subscribed,omitempty"`
	CanUserSubscribe bool   `json:"can_user_subscribe,omitempty"`
	SubproductID     uint64 `json:"subproductid,omitempty"`
	DollarPrice      string `json:"dollar_price,omitempty"`
	Settings         struct {
		AccessSettings string `json:"access_settings"` // TODO: enum[all,future_only,limited_past_and_future]
	} `json:"settings,omitempty"`
	Stats struct {
		Subscribers uint32 `json:"subscribers,omitempty"`
		Deviations  uint32 `json:"deviations,omitempty"`
		Posts       uint32 `json:"posts,omitempty"`
		Total       uint32 `json:"total,omitempty"`
	} `json:"stats"`
	Benefits []any `json:"benefits"`
}

type DeviationUpdateResponse struct {
	Status      string    `json:"status"`
	URL         string    `json:"url"`
	DeviationID uuid.UUID `json:"deviationid"`
}

// Deviation fetches a deviation.
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

func (s *deviationService) Content(deviationID uuid.UUID) (any, error) {
	// TODO: Implement /deviation/content endpoint.
	return nil, errors.New("not implemented yet")
}

func (s *deviationService) EmbeddedContent(deviationID uuid.UUID) (any, error) {
	// TODO: Implement /deviation/embeddedcontent endpoint.
	return nil, errors.New("not implemented yet")
}

func (s *deviationService) Metadata() (any, error) {
	// TODO: Implement /metadata.
	return nil, errors.New("not implemented yet")
}

func (s *deviationService) WhoFaved() (any, error) {
	// TODO: Implement /whofaved.
	return nil, errors.New("not implemented yet")
}
