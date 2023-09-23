package deviantart

import (
	"fmt"
	"strconv"

	"github.com/dghubble/sling"
	"github.com/dustin/go-humanize"
)

type fileInfo struct {
	Source string `json:"src"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type StashFile struct {
	fileInfo
	Transparency bool `json:"transparency"`
}

type StashService struct {
	sling *sling.Sling
}

func newStashService(sling *sling.Sling) *StashService {
	return &StashService{
		sling: sling.Path("stash/"),
	}
}

type StashMetadata struct {
	Title          string            `json:"title"`
	Path           string            `json:"path,omitempty"`
	Size           int               `json:"size,omitempty"`
	Description    string            `json:"description,omitempty"` // html
	ParentID       int               `json:"parentid,omitempty"`
	Thumb          *StashFile        `json:"thumb,omitempty"`
	ArtistComments string            `json:"artist_comments,omitempty"` //html
	OriginalURL    string            `json:"original_url,omitempty"`
	Category       string            `json:"category,omitempty"`
	CreationTime   int64             `json:"creation_time,omitempty"`
	Files          []StashFile       `json:"files,omitempty"`
	Submission     *StashSubmission  `json:"submission,omitempty"`
	Stats          *StashStats       `json:"stats,omitempty"`
	Camera         map[string]string `json:"camera,omitempty"`
	StackID        int               `json:"stackid"`
	Tags           []string          `json:"tags,omitempty"`
}

type StashSubmission struct {
	FileSize      string `json:"file_size,omitempty"`
	Resolution    string `json:"resolution,omitempty"`
	SubmittedWith *struct {
		App string `json:"app,omitempty"`
		URL string `json:"url,omitempty"`
	} `json:"submitted_with,omitempty"`
}

type StashStats struct {
	Views          int `json:"views,omitempty"`
	ViewsToday     int `json:"views_today,omitempty"`
	Downloads      int `json:"downloads,omitempty"`
	DownloadsToday int `json:"downloads_today,omitempty"`
}

// Stack fetches a stash stack's metadata.
//
// To connect to this endpoint OAuth2 Access Token from the Client Credentials
// Grant, or Authorization Code Grant is required.
//
// The following scopes are required to access this resource:
//
//   - stash
func (s *StashService) Stack(stackID int64) (StashMetadata, error) {
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

type StackContentsParams struct {
	Offset uint16 `url:"offset,omitempty"`
	Limit  uint8  `url:"limit,omitempty"`

	IncludeSubmission bool `url:"ext_submission,omitempty"`
	IncludeCamera     bool `url:"ext_camera,omitempty"`
	IncludeStats      bool `url:"ext_stats,omitempty"`
}

// RootStackID is an ID to list contents of a root stack.
const RootStackID = 0

// StackContents fetches stack contents.
//
// To connect to this endpoint OAuth2 Access Token from the Authorization Code
// Grant is required.
//
// The following scopes are required to access this resource:
//
//   - stash
func (s *StashService) StackContents(stackID int64, params *StackContentsParams) (OffsetResponse[StashMetadata], error) {
	var (
		success OffsetResponse[StashMetadata]
		failure Error
	)
	_, err := s.sling.New().Get(strconv.FormatInt(stackID, 10)+"/contents").QueryStruct(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return OffsetResponse[StashMetadata]{}, fmt.Errorf("unable to fetch stack contents: %w", err)
	}
	return success, nil
}

type deleteParams struct {
	ItemID int64 `json:"itemid"`
}

// Delete deletes a previously submitted file.
//
// To connect to this endpoint OAuth2 Access Token from the Authorization Code
// Grant is required.
//
// The following scopes are required to access this resource:
//
//   - stash
func (s *StashService) Delete(itemID int64) (bool, error) {
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

type StashDeltaResponse struct {
	Cursor     string `json:"cursor"`
	HasMore    bool   `json:"has_more"`
	NextOffset int    `json:"next_offset"`
	Reset      bool   `json:"reset"`
	Entries    []struct {
		ItemID   int64         `json:"itemid,omitempty"`
		StackID  int64         `json:"stackid,omitempty"`
		Metadata StashMetadata `json:"stash_metadata"`
		Position int           `json:"position,omitempty"`
	} `json:"entries,omitempty"`
}

type StashDeltaParams struct {
	// The cursor hash provided to your app in the last delta call.
	Cursor string `url:"cursor,omitempty"`

	// The pagination offset
	Offset int `url:"offset,omitempty"`

	// The pagination limit
	Limit int `url:"limit,omitempty"`

	// Include extended submission information
	IncludeSubmission bool `url:"ext_submission"`

	// Include camera EXIF information
	IncludeCamera bool `url:"ext_camera"`

	// Include extended statistics
	IncludeStats bool `url:"ext_stats"`
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
// To connect to this endpoint OAuth2 Access Token from the Authorization Code
// Grant is required.
//
// The following scopes are required to access this resource:
//
//   - stash
func (s *StashService) Delta(params *StashDeltaParams) (StashDeltaResponse, error) {
	var (
		success StashDeltaResponse
		failure Error
	)
	_, err := s.sling.New().Get("delta").QueryStruct(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return StashDeltaResponse{}, fmt.Errorf("unable to fetch felta: %w", err)
	}
	return success, nil
}

type StashMoveResponse struct {
	Target  StashMetadata   `json:"target"`
	Changes []StashMetadata `json:"changes"`
}

type moveParams struct {
	TargetID int64 `url:"targetid,omitempty"`
}

// Move moves the stack into the target stack.
//
// The response includes updated metadata of the target stack and changes in its
// contents. Changes include at least the stack being moved. If the move
// operation resulted in new stacks being created, they will be returned as
// well. New stack may be created if the target stack is a leaf stack, i.e. has
// no children stacks yet.
//
// To connect to this endpoint OAuth2 Access Token from the Authorization Code
// Grant is required.
//
// The following scopes are required to access this resource:
//
//   - stash
func (s *StashService) Move(stackID, targetID int64) (StashMoveResponse, error) {
	var (
		success StashMoveResponse
		failure Error
	)
	params := &moveParams{TargetID: targetID}
	_, err := s.sling.New().Post("move/").Path(strconv.FormatInt(stackID, 10)).BodyForm(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return StashMoveResponse{}, fmt.Errorf("unable to move stash: %w", err)
	}
	return success, nil
}

// Position changes the position of a stack within its parent.
//
// To connect to this endpoint OAuth2 Access Token from the Authorization Code
// Grant is required.
//
// The following scopes are required to access this resource:
//
//   - stash
func (s *StashService) Position(stackID, position int64) (bool, error) {
	type positionParams struct {
		Position int64 `url:"position"`
	}
	var (
		success map[string]any
		failure Error
	)
	stackPath := strconv.FormatInt(stackID, 10)
	params := &positionParams{Position: position}
	_, err := s.sling.New().Post("position/").Path(stackPath).BodyForm(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return false, fmt.Errorf("unable to change stash position: %w", err)
	}
	return success["success"].(bool), nil
}

type StashUserdata struct {
	Features   []string `json:"features"`
	Agreements []string `json:"agreements"`
}

// Userdata fetches users data about features and agreements.
//
// To connect to this endpoint OAuth2 Access Token from the Authorization Code
// Grant is required.
//
// The following scopes are required to access this resource:
//
//   - stash
//   - publish
func (s *StashService) Userdata() (StashUserdata, error) {
	var (
		success StashUserdata
		failure Error
	)
	_, err := s.sling.New().Get("publish/userdata").Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return StashUserdata{}, fmt.Errorf("unable to fetch userdata: %w", err)
	}
	return success, nil
}

type StashSpace struct {
	AvailableSpace uint64 `json:"available_space"`
	TotalSpace     uint64 `json:"total_space"`
}

// String represents available space in human readable form.
func (s StashSpace) String() string {
	return fmt.Sprintf("%s of %s", humanize.IBytes(s.AvailableSpace), humanize.IBytes(s.TotalSpace))
}

// Space returns how much sta.sh space (expressed in bytes) a user has available
// for new uploads.
//
// To connect to this endpoint OAuth2 Access Token from the Authorization Code
// Grant is required.
//
// The following scopes are required to access this resource:
//
//   - stash
func (s *StashService) Space() (StashSpace, error) {
	var (
		success StashSpace
		failure Error
	)
	_, err := s.sling.New().Get("space").Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return StashSpace{}, fmt.Errorf("unable to fetch sta.sh space: %w", err)
	}
	return success, nil
}

type StashUpdateParams struct {
	Title       string `url:"title,omitempty"`
	Description string `url:"description,omitempty"`
}

// Update updates the stash stack's details.
//
// To connect to this endpoint OAuth2 Access Token from the Authorization Code
// Grant is required.
//
// The following scopes are required to access this resource:
//
//   - stash
func (s *StashService) Update(stackID int64, params *StashUpdateParams) (bool, error) {
	var (
		success map[string]any
		failure Error
	)
	stackPath := strconv.FormatInt(stackID, 10)
	_, err := s.sling.New().Post("update/").Path(stackPath).BodyForm(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return false, fmt.Errorf("unable to update stack: %w", err)
	}
	return success["success"].(bool), nil
}
