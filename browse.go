package deviantart

import (
	"fmt"
	"time"

	"github.com/dghubble/sling"
	"github.com/google/uuid"
)

type browseService struct {
	sling *sling.Sling
}

func newBrowseService(sling *sling.Sling) *browseService {
	return &browseService{
		sling: sling.Path("browse/"),
	}
}

type DailyDeviationsParams struct {
	Date        time.Time `url:"date,omitempty" layout:"2006-01-02"`
	WithSession bool      `url:"with_session,omitempty"` // TODO: Move WithSession to parameters.
}

// DailyDeviations fetches daily deviations.
//
// The following scopes are required to access this resource:
//
//   - browse
//
// TODO: The endpoint returns the `has_more` field, but there is no offset or
// cursor pagination information. This case requires further investigation.
func (s *browseService) DailyDeviations(params *DailyDeviationsParams) (OffsetResponse[Deviation], error) {
	var (
		success OffsetResponse[Deviation]
		failure Error
	)
	_, err := s.sling.New().Get("dailydeviations").QueryStruct(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return OffsetResponse[Deviation]{}, fmt.Errorf("unable to fetch daily deviations: %w", err)
	}
	return success, nil
}

// DeviantsYouWatch fetches deviations of deviants you watch.
//
// The following scopes are required to access this resource:
//
//   - browse
func (s *browseService) DeviantsYouWatch(page *OffsetParams) (OffsetResponse[Deviation], error) {
	var (
		success OffsetResponse[Deviation]
		failure Error
	)
	_, err := s.sling.New().Get("deviantsyouwatch").QueryStruct(page).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return OffsetResponse[Deviation]{}, fmt.Errorf("unable to fetch deviants for you: %w", err)
	}
	return success, nil
}

type MoreLikeThisPreviewResponse struct {
	Seed                 uuid.UUID   `json:"seed"`
	Author               User        `json:"user"`
	MoreFromArtist       []Deviation `json:"more_from_artist"`
	MoreFromDA           []Deviation `json:"more_from_da"`
	SuggestedCollections []struct {
		Collection Folder      `json:"collection"` //Gallection
		Deviations []Deviation `json:"deviations"`
	} `json:"suggested_collections,omitempty"`
}

// MoreLikeThisPreview fetches More Like This preview result for a seed deviation.
//
// The following scopes are required to access this resource:
//
//   - browse
//   - browse.mlt
func (s *browseService) MoreLikeThisPreview(seed uuid.UUID) (MoreLikeThisPreviewResponse, error) {
	type seedParams struct {
		Seed uuid.UUID `url:"seed"`
	}
	var (
		success MoreLikeThisPreviewResponse
		failure Error
	)
	_, err := s.sling.New().Get("morelikethis/preview").QueryStruct(&seedParams{Seed: seed}).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return MoreLikeThisPreviewResponse{}, fmt.Errorf("unable to fetch more like this: %w", err)
	}
	return success, nil
}

type SearchParams struct {
	// Search query term.
	//
	// Estimated total results count would be available on EstimatedTotal field.
	Query string `url:"q,omitempty"`
}

// Newest fetches newest deviations.
//
// The following scopes are required to access this resource:
//
//   - browse
func (s *browseService) Newest(params *SearchParams, page *OffsetParams) (OffsetResponse[Deviation], error) {
	var (
		success OffsetResponse[Deviation]
		failure Error
	)
	_, err := s.sling.New().Get("newest").QueryStruct(params).QueryStruct(page).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return OffsetResponse[Deviation]{}, fmt.Errorf("unable to fetch newest deviations: %w", err)
	}
	return success, nil
}

var (
	TimeRangeNow   = "now"
	TimeRangeWeek  = "1week"
	TimeRangeMonth = "1month"
	TimeRangeAll   = "alltime"
)

type PopularParams struct {
	// Search query term.
	//
	// Estimated total results count would be available on EstimatedTotal field.
	Query string `url:"q,omitempty"`

	// The timerange.
	//
	// TODO: Valid values are: values(now, 1week, 1month, alltime).
	TimeRange string `url:"timerange,omitempty"`
}

// Popular fetches popular deviations.
//
// The following scopes are required to access this resource:
//
//   - browse
//
// BUG: Query does not work properly.
// See: https://github.com/wix-incubator/DeviantArt-API/issues/206.
func (s *browseService) Popular(params *PopularParams, page *OffsetParams) (OffsetResponse[Deviation], error) {
	var (
		success OffsetResponse[Deviation]
		failure Error
	)
	_, err := s.sling.New().Get("popular").QueryStruct(params).QueryStruct(page).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return OffsetResponse[Deviation]{}, fmt.Errorf("unable to fetch popular deviations: %w", err)
	}
	return success, nil
}

type JournalStatus struct {
	Journal *Deviation `json:"journal"`
	Status  *Status    `json:"status"`
}

// PostsDeviantsYouWatch returns deviants you watch.
//
// The following scopes are required to access this resource:
//
//   - browse
func (s *browseService) PostsDeviantsYouWatch(page *OffsetParams) (OffsetResponse[JournalStatus], error) {
	var (
		success OffsetResponse[JournalStatus]
		failure Error
	)
	_, err := s.sling.New().Get("posts/deviantsyouwatch").QueryStruct(page).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return OffsetResponse[JournalStatus]{}, fmt.Errorf("unable to fetch deviants for you: %w", err)
	}
	return success, nil
}

// Recommended fetches recommended deviations.
//
// The following scopes are required to access this resource:
//
//   - browse
//
// TODO: Documentation specifies the `suggested_reasons` field but is absend in
// all responses. This case requires further investigation.
func (s *browseService) Recommended(params *SearchParams) (OffsetResponse[Deviation], error) {
	var (
		success OffsetResponse[Deviation]
		failure Error
	)
	_, err := s.sling.New().Get("recommended").QueryStruct(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return OffsetResponse[Deviation]{}, fmt.Errorf("unable to fetch recommended deviations: %w", err)
	}
	return success, nil
}

// Tags fetches a tag.
//
// NOTE: This endpoint supports cursor- and offset-base pagination.
// But for simplicity, I'll stick to cursor params for now.
func (s *browseService) Tags(tag string, page *CursorParams) (CursorResponse[Deviation], error) {
	type tagParams struct {
		Tag string `url:"tag"`
	}
	var (
		success CursorResponse[Deviation]
		failure Error
	)
	_, err := s.sling.New().Get("tags").QueryStruct(&tagParams{Tag: tag}).QueryStruct(page).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return CursorResponse[Deviation]{}, fmt.Errorf("unable to fetch tags: %w", err)
	}
	return success, nil
}

// TagsSearch autocompletes tags.
//
// The `tag_name“ parameter should not contain spaces. If it does, spaces will
// be stripped and remainder will be treated as a single tag.
func (s *browseService) TagsSearch(tag string) ([]string, error) {
	type tagName struct {
		Name string `json:"tag_name" url:"tag_name"`
	}
	var (
		success singleResponse[tagName]
		failure Error
	)
	_, err := s.sling.New().Get("tags/search").QueryStruct(&tagName{Name: tag}).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return nil, fmt.Errorf("unable to search tags: %w", err)
	}

	tags := make([]string, 0, len(success.Results))
	for _, tag := range success.Results {
		tags = append(tags, tag.Name)
	}
	return tags, nil
}

// Topic fetches topic deviations.
//
// The following scopes are required to access this resource:
//
//   - browse
func (s *browseService) Topic(topic string, page *CursorParams) (CursorResponse[Deviation], error) {
	type topicParams struct {
		Topic string `url:"topic"`
	}
	var (
		success CursorResponse[Deviation]
		failure Error
	)
	_, err := s.sling.New().Get("topic").QueryStruct(&topicParams{Topic: topic}).QueryStruct(page).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return CursorResponse[Deviation]{}, fmt.Errorf("unable to fetch topic: %w", err)
	}
	return success, nil
}

type Topic struct {
	Name              string      `json:"name"`
	CanonicalName     string      `json:"canonical_name"`
	ExampleDeviations []Deviation `json:"example_deviations,omitempty"`
}

// Topics fetches topics and deviations from each topic.
func (s *browseService) Topics(page *CursorParams) (CursorResponse[Topic], error) {
	var (
		success CursorResponse[Topic]
		failure Error
	)
	_, err := s.sling.New().Get("topics").QueryStruct(page).Receive(success, failure)
	if err := relevantError(err, failure); err != nil {
		return CursorResponse[Topic]{}, fmt.Errorf("unable to fetch topics: %w", err)
	}
	return success, nil
}

type TopTopic struct {
	Name              string      `json:"name"`
	CanonicalName     string      `json:"canonical_name"`
	ExampleDeviations []Deviation `json:"example_deviations,omitempty"`
}

// Topics fetches top topics with example deviation for each one.
func (s *browseService) TopTopics(page *CursorParams) (CursorResponse[TopTopic], error) {
	var (
		success CursorResponse[TopTopic]
		failure Error
	)
	_, err := s.sling.New().Get("toptopics").QueryStruct(page).Receive(success, failure)
	if err := relevantError(err, failure); err != nil {
		return CursorResponse[TopTopic]{}, fmt.Errorf("unable to fetch topics: %w", err)
	}
	return success, nil
}

type UserJournalsParams struct {
	// The username of the user to fetch journals for.
	Username string `url:"username"`

	// Fetch only featured or not.
	Featured bool `url:"featured,omitempty"`
}

// UserJournals browses journals of a user.
func (s *browseService) UserJournals(params *UserJournalsParams, page *OffsetParams) (OffsetResponse[Deviation], error) {
	var (
		success OffsetResponse[Deviation]
		failure Error
	)
	_, err := s.sling.New().Get("user/journals").QueryStruct(params).QueryStruct(page).Receive(success, failure)
	if err := relevantError(err, failure); err != nil {
		return OffsetResponse[Deviation]{}, fmt.Errorf("unable to browse user journals: %w", err)
	}
	return success, nil
}
