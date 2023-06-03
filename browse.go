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
	Query string `url:"q,omitempty"`
}

// Newest fetches newest deviations.
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

type PopularParams struct {
	// Search query term.
	Query string `url:"q,omitempty"`

	// The timerange.
	//
	// TODO: Valid values are: values(now, 1week, 1month, alltime).
	TimeRange string `url:"timerange,omitempty"`
}

// Popular fetches popular deviations.
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

// TODO: what is it?
func (s *browseService) PostsDeviantsYouWatch(page *OffsetParams) (OffsetResponse[Deviation], error) {
	var (
		success OffsetResponse[Deviation]
		failure Error
	)
	_, err := s.sling.New().Get("posts/deviantsyouwatch").QueryStruct(page).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return OffsetResponse[Deviation]{}, fmt.Errorf("unable to fetch deviants for you: %w", err)
	}
	return success, nil
}

// Recommended fetches recommended deviations.
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

type TagsParams struct {
	// The tag to browse.
	Tag string `url:"tag"`
}

// Tags fetches a tag.
//
// NOTE: This endpoing supports cursor- and offset-base pagination.
// But for simplicity, I'll stick to cursor params for now.
func (s *browseService) Tags(params *TagsParams, page *CursorParams) (CursorResponse[Deviation], error) {
	var (
		success CursorResponse[Deviation]
		failure Error
	)
	_, err := s.sling.New().Get("tags").QueryStruct(params).QueryStruct(page).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return CursorResponse[Deviation]{}, fmt.Errorf("unable to fetch tags: %w", err)
	}
	return success, nil
}

type TagsSearchResponse struct {
	Results []struct {
		Name string `json:"tag_name"`
	} `json:"results"`
}

type TagsSearchParams struct {
	TagName     string `url:"tag_name"`
	WithSession bool   `url:"with_session,omitempty"`
}

// TagsSearch autocompletes tags.
func (s *browseService) TagsSearch(params *TagsSearchParams) ([]string, error) {
	var (
		success TagsSearchResponse
		failure Error
	)
	_, err := s.sling.New().Get("tags/search").QueryStruct(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return nil, fmt.Errorf("unable to search tags: %w", err)
	}

	tags := make([]string, 0, len(success.Results))
	for _, tag := range success.Results {
		tags = append(tags, tag.Name)
	}
	return tags, nil
}

type TopicParams struct {
	// Topic name.
	Topic string `url:"topic"`
}

// Topic fetches topic deviations.
func (s *browseService) Topic(params *TopicParams, page *CursorParams) (CursorResponse[Deviation], error) {
	var (
		success CursorResponse[Deviation]
		failure Error
	)
	_, err := s.sling.New().Get("topic").QueryStruct(params).QueryStruct(page).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return CursorResponse[Deviation]{}, fmt.Errorf("unable to fetch topic: %w", err)
	}
	return success, nil
}

type Topic struct {
	Name              string      `json:"name"`
	CanonicalName     string      `json:"canonical_name"`
	ExampleDeviations []Deviation `json:"example_deviations,omitempty"`
	Deviations        []Deviation `json:"deviations,omitempty"`
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
	Featured bool `url:"featured"`
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
