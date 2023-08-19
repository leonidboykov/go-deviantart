package deviantart

import (
	"fmt"
	"strconv"
)

type StashItem struct {
	ItemID   int64    `json:"itemid"`
	HTML     string   `json:"html,omitempty"`
	CSS      string   `json:"css,omitempty"`
	CSSFonts []string `json:"css_fonts,omitempty"`
	StashMetadata
}

type ItemParams struct {
	IncludeSubmission bool `url:"ext_submission,omitempty"`
	IncludeCamera     bool `url:"ext_camera,omitempty"`
	IncludeStats      bool `url:"ext_stats,omitempty"`
}

// Item fetches stash item's metadata.
//
// To connect to this endpoint OAuth2 Access Token from the Client Credentials
// Grant, or Authorization Code Grant is required.
//
// The following scopes are required to access this resource:
//
//   - stash
func (s *StashService) Item(itemID int64, params *ItemParams) (StashItem, error) {
	var (
		success StashItem
		failure Error
	)
	_, err := s.sling.New().Post("item/").Path(strconv.FormatInt(itemID, 10)).QueryStruct(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return StashItem{}, fmt.Errorf("unable to fetch item: %w", err)
	}
	return success, nil
}
