package deviantart

import (
	"fmt"
	"time"

	"github.com/dghubble/sling"
	"github.com/google/uuid"
)

type BrowseService struct {
	sling *sling.Sling
}

func newBrowseService(sling *sling.Sling) *BrowseService {
	return &BrowseService{
		sling: sling.Path("browse/"),
	}
}

// DailyDeviations fetches daily deviations.
//
// To connect to this endpoint OAuth2 Access Token from the [ClientCredentials],
// or [AuthorizationCode] with a [BrowseScope] scope is required.
//
// TODO: The endpoint returns the `has_more` field, but there is no offset or
// cursor pagination information. This case requires further investigation.
func (s *BrowseService) DailyDeviations(date time.Time) (OffsetResponse[Deviation], error) {
	type dateParams struct {
		Date time.Time `url:"date,omitempty" layout:"2006-01-02"`
	}
	var (
		success OffsetResponse[Deviation]
		failure Error
	)
	_, err := s.sling.New().Get("dailydeviations").QueryStruct(&dateParams{Date: date}).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return OffsetResponse[Deviation]{}, fmt.Errorf("unable to fetch daily deviations: %w", err)
	}
	return success, nil
}

// DeviantsYouWatch fetches deviations of deviants you watch.
//
// To connect to this endpoint OAuth2 Access Token from the [AuthorizationCode]
// with a [BrowseScope] is required.
func (s *BrowseService) DeviantsYouWatch(page *OffsetParams) (OffsetResponse[Deviation], error) {
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

// Home browses homepage.
//
// To connect to this endpoint OAuth2 Access Token from the [AuthorizationCode]
// with a [BrowseScope] is required.
func (s *BrowseService) Home(query string, page *OffsetParams) (OffsetResponse[Deviation], error) {
	type searchParams struct {
		// Search query term.
		//
		// Estimated total results count would be available on EstimatedTotal field.
		Query string `url:"q,omitempty"`
	}
	var (
		success OffsetResponse[Deviation]
		failure Error
	)
	params := &searchParams{Query: query}
	_, err := s.sling.New().Get("home").QueryStruct(params).QueryStruct(page).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return OffsetResponse[Deviation]{}, fmt.Errorf("unable to fetch home deviations: %w", err)
	}
	return success, nil
}

type MoreLikeThisPreviewResponse struct {
	Seed                 uuid.UUID   `json:"seed"`
	Author               User        `json:"user"`
	MoreFromArtist       []Deviation `json:"more_from_artist"`
	MoreFromDeviantArt   []Deviation `json:"more_from_da"`
	SuggestedCollections []struct {
		Collection Folder      `json:"collection"` //Gallection
		Deviations []Deviation `json:"deviations"`
	} `json:"suggested_collections,omitempty"`
}

// MoreLikeThisPreview fetches More Like This preview result for a seed deviation.
//
// To connect to this endpoint OAuth2 Access Token from the Client Credentials
// Grant, or Authorization Code Grant is required.
//
// The following scopes are required to access this resource:
//
//   - browse
//   - browse.mlt
func (s *BrowseService) MoreLikeThisPreview(seed uuid.UUID) (MoreLikeThisPreviewResponse, error) {
	type seedParams struct {
		Seed string `url:"seed"`
	}
	var (
		success MoreLikeThisPreviewResponse
		failure Error
	)
	params := &seedParams{Seed: seed.String()}
	_, err := s.sling.New().Get("morelikethis/preview").QueryStruct(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return MoreLikeThisPreviewResponse{}, fmt.Errorf("unable to fetch more like this: %w", err)
	}
	return success, nil
}

// Tags fetches a tag.
//
// To connect to this endpoint OAuth2 Access Token from the [ClientCredentials],
// or [AuthorizationCode] with a [BrowseScope] scope is required.
//
// NOTE: This endpoint supports cursor- and offset-base pagination.
// But for simplicity, I'll stick to cursor params for now.
func (s *BrowseService) Tags(tag string, page *CursorParams) (CursorResponse[Deviation], error) {
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
// The `tag` parameter should not contain spaces. If it does, spaces will be
// stripped and remainder will be treated as a single tag.
//
// To connect to this endpoint OAuth2 Access Token from the [ClientCredentials],
// or [AuthorizationCode] with a [BrowseScope] scope is required.
func (s *BrowseService) TagsSearch(tag string) ([]string, error) {
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
// To connect to this endpoint OAuth2 Access Token from the [ClientCredentials],
// or [AuthorizationCode] with a [BrowseScope] scope is required.
func (s *BrowseService) Topic(topic string, page *CursorParams) (CursorResponse[Deviation], error) {
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
//
// To connect to this endpoint OAuth2 Access Token from the [ClientCredentials],
// or [AuthorizationCode] with a [BrowseScope] scope is required.
func (s *BrowseService) Topics(page *CursorParams) (CursorResponse[Topic], error) {
	var (
		success CursorResponse[Topic]
		failure Error
	)
	_, err := s.sling.New().Get("topics").QueryStruct(page).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return CursorResponse[Topic]{}, fmt.Errorf("unable to fetch topics: %w", err)
	}
	return success, nil
}

// Topics fetches top topics with example deviation for each one.
//
// To connect to this endpoint OAuth2 Access Token from the [ClientCredentials],
// or [AuthorizationCode] with a [BrowseScope] scope is required.
func (s *BrowseService) TopTopics(page *CursorParams) (CursorResponse[Topic], error) {
	var (
		success CursorResponse[Topic]
		failure Error
	)
	_, err := s.sling.New().Get("toptopics").QueryStruct(page).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return CursorResponse[Topic]{}, fmt.Errorf("unable to fetch topics: %w", err)
	}
	return success, nil
}
