package deviantart

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/dghubble/sling"
)

// TODO: STASH
// 	/{stackid}
// 	/{stackid}/contents
// 	/delete
// 	/delta
// 	/item/{itemid}
// 	/move/{stackid}
// 	/position/{stackid}
// 	/publish
// 	/publish/categorytree
// 	/publish/userdata
// 	/space
// 	/submit
// 	/update/{stackid}

type StashObject struct {
	Src          string `json:"src"`
	Height       int    `json:"height"`
	Width        int    `json:"width"`
	Transparency bool   `json:"transparency"`
}

type stashService struct {
	sling *sling.Sling
}

func newStashService(sling *sling.Sling) *stashService {
	return &stashService{
		sling: sling.Path("stash/"),
	}
}

type StashMetadata struct {
	Title          string        `json:"title"`
	Path           string        `json:"path,omitempty"`
	Size           int           `json:"size,omitempty"`
	Description    string        `json:"description,omitempty"` // html
	ParentID       int           `json:"parentid,omitempty"`
	Thumb          *StashObject  `json:"thumb,omitempty"`
	ArtistComments string        `json:"artist_comments,omitempty"` //html
	OriginalURL    string        `json:"original_url,omitempty"`
	Category       string        `json:"category,omitempty"`
	CreationTime   int64         `json:"creation_time,omitempty"`
	Files          []StashObject `json:"files,omitempty"`
	Submission     *struct {
		FileSize      string `json:"file_size,omitempty"`
		Resolution    string `json:"resolution,omitempty"`
		SubmittedWith *struct {
			App string `json:"app,omitempty"`
			URL string `json:"url,omitempty"`
		} `json:"submitted_with,omitempty"`
	} `json:"submission,omitempty"`
	Stats *struct {
		Views          int `json:"views,omitempty"`
		ViewsToday     int `json:"views_today,omitempty"`
		Downloads      int `json:"downloads,omitempty"`
		DownloadsToday int `json:"downloads_today,omitempty"`
	} `json:"stats,omitempty"`
	Camera  any      `json:"camera,omitempty"`
	StackID int      `json:"stackid"`
	Tags    []string `json:"tags,omitempty"`
}

// Stack fetches a stash stack's metadata.
func (s *stashService) StackMetadata(stackID int64) (StashMetadata, error) {
	var (
		success StashMetadata
		failure Error
	)
	_, err := s.sling.New().Get(strconv.FormatInt(stackID, 10)).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return StashMetadata{}, fmt.Errorf("unable to fetch stack metadata: %w", err)
	}
	return success, nil
}

type StackContents struct {
	Results    []StashMetadata `json:"results"`
	HasMore    bool            `json:"has_more"`
	NextOffset int32           `json:"next_offset,omitempty"`
}

type StackContentsParams struct {
	Offset            uint16 `url:"offset,omitempty"`
	Limit             uint8  `url:"limit,omitempty"`
	IncludeSubmission bool   `url:"ext_submission,omitempty"`
	IncludeCamera     bool   `url:"ext_camera,omitempty"`
	IncludeStats      bool   `url:"ext_stats,omitempty"`
}

// RootStackID is an ID to list contents of a root stack.
const RootStackID = 0

// StackContents fetches stack contents.
//
// Requires Authorization Code grant.
func (s *stashService) StackContents(stackID int64, params *StackContentsParams) (StackContents, error) {
	var (
		success StackContents
		failure Error
	)
	_, err := s.sling.New().Get(strconv.FormatInt(stackID, 10)+"/contents").QueryStruct(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return StackContents{}, fmt.Errorf("unable to fetch stack contents: %w", err)
	}
	return success, nil
}

type deleteParams struct {
	ItemID int64 `json:"itemid"`
}

// Delete deletes a previously submitted file.
//
// Requires Authorization Code grant.
func (s *stashService) Delete(itemID int64) (bool, error) {
	var (
		// TODO: Check error_description.
		success map[string]any
		failure Error
	)
	_, err := s.sling.New().Post("delete").BodyForm(deleteParams{ItemID: itemID}).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return false, fmt.Errorf("unable to delete item: %w", err)
	}
	return success["success"].(bool), nil
}

// Delta retrieves contents of Sta.sh for a user.
//
// This endpoint is used to retrieve all data available in user's Sta.sh. The
// Sta.sh data is organized into stacks and items. A stack contains either a
// list of children stacks or a single item, that is, a stack can't contain more
// than one item or a mix of stacks and items. Stacks can be moved into another
// stack or positioned within a parent. An item represents a Sta.sh submission
// and is always contained within some stack.
//
// This endpoint is incremental. The first time you call it for a user, your app
// will receive the full list of that user's Sta.sh stacks and items. Your app
// should cache these results locally.
//
// Afterward, your app should then provide a cursor parameter for all /delta
// calls. This cursor tells us which data you have already received so that we
// can send you only new and modified stacks and items.
//
// Requires Authorization Code grant.
func (s *stashService) Delta() (any, error) {
	// TODO: Checkout this endpoint.
	return nil, errors.New("not implemented yet")
}

// Move moves the stack into the target stack.
//
// The response includes updated metadata of the target stack and changes in its
// contents. Changes include at least the stack being moved. If the move
// operation resulted in new stacks being created, they will be returned as
// well. New stack may be created if the target stack is a leaf stack, i.e. has
// no children stacks yet.
//
// Requires Authorization Code grant.
func (s *stashService) Move(stackID int64, targetID int64) (any, error) {
	// TODO: Checkout this endpoint.
	return nil, errors.New("not implemented yet")
}

func (s *stashService) Position(stackID int64, position int64) (any, error) {
	// TODO: Checkout this endpoint.
	return nil, errors.New("not implemented yet")
}
