package deviantart

import (
	"fmt"
	"time"

	"github.com/dghubble/sling"
)

// TODO: BROWSE
// 	/tags
// 	/tags/search
// 	/topic
// 	/topics
// 	/toptopics
// 	/user/journals

type browseService struct {
	sling *sling.Sling
}

func newBrowseService(sling *sling.Sling) *browseService {
	return &browseService{
		sling: sling.Path("browse/"),
	}
}

type DeviationsResponse struct {
	HasMore        bool        `json:"has_more"`
	NextOffset     uint32      `json:"next_offset,omitempty"`
	EstimatedTotal uint32      `json:"estimated_total,omitempty"`
	Results        []Deviation `json:"results"`
	Session        *Session    `json:"session,omitempty"`
}

type DailyDeviationsParams struct {
	Date        time.Time `url:"date,omitempty" layout:"2006-01-02"`
	WithSession bool      `url:"with_session,omitempty"`
}

// DailyDeviations fetches daily deviations.
func (s *browseService) DailyDeviations(params *DailyDeviationsParams) (DeviationsResponse, error) {
	var (
		success DeviationsResponse
		failure Error
	)
	_, err := s.sling.New().Get("dailydeviations").QueryStruct(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return DeviationsResponse{}, fmt.Errorf("unable to fetch daily deviations: %w", err)
	}
	return success, nil
}

type DeviantsYouWatchParams struct {
	Limit       uint32 `url:"limit,omitempty"`
	Offset      uint32 `url:"offset,omitempty"`
	WithSession bool   `url:"with_session,omitempty"`
}

// DeviantsYouWatch fetches deviations of deviants you watch.
func (s *browseService) DeviantsYouWatch(params *DeviantsYouWatchParams) (DeviationsResponse, error) {
	var (
		success DeviationsResponse
		failure Error
	)
	_, err := s.sling.New().Get("deviantsyouwatch").QueryStruct(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return DeviationsResponse{}, fmt.Errorf("unable to fetch deviants for you: %w", err)
	}
	return success, nil
}

// TODO: MoreLikeThis/Preview is actual?

type SearchParams struct {
	// Search query term.
	Query string `url:"q,omitempty"`

	// The pagination offset.
	Offset int `url:"offset,omitempty"`

	// The pagination offset.
	Limit int `url:"offset,omitempty"`

	WithSession bool `url:"with_session,omitempty"`
}

// Newest fetches newest deviations.
func (s *browseService) Newest(params *SearchParams) (DeviationsResponse, error) {
	var (
		success DeviationsResponse
		failure Error
	)
	_, err := s.sling.New().Get("newest").QueryStruct(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return DeviationsResponse{}, fmt.Errorf("unable to fetch newest deviations: %w", err)
	}
	return DeviationsResponse{}, nil
}

// Popular fetches popular deviations.
func (s *browseService) Popular(params *SearchParams) (DeviationsResponse, error) {
	var (
		success DeviationsResponse
		failure Error
	)
	_, err := s.sling.New().Get("popular").QueryStruct(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return DeviationsResponse{}, fmt.Errorf("unable to fetch popular deviations: %w", err)
	}
	return success, nil
}

// TODO: what is it?
func (s *browseService) PostsDeviantsYouWatch(params *SearchParams) (DeviationsResponse, error) {
	var (
		success DeviationsResponse
		failure Error
	)
	_, err := s.sling.New().Get("posts/deviantsyouwatch").QueryStruct(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return DeviationsResponse{}, fmt.Errorf("unable to fetch deviants for you: %w", err)
	}
	return success, nil
}

// Recommended fetches recommended deviations.
func (s *browseService) Recommended(params *SearchParams) (DeviationsResponse, error) {
	var (
		success DeviationsResponse
		failure Error
	)
	_, err := s.sling.New().Get("recommended").QueryStruct(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return DeviationsResponse{}, fmt.Errorf("unable to fetch recommended deviations: %w", err)
	}
	return success, nil
}

// Tags fetches a tag.
// func (s *browseService) Tags()
