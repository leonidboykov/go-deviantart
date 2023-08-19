package deviantart

import (
	"fmt"

	"github.com/dghubble/sling"
	"github.com/google/uuid"
)

// TODO: Embed to Gallery and Collection?
type Folder struct {
	// FolderID uuid.UUID `json:"folderid"` // TODO: Remove it?
	FolderID int64  `json:"folderid"`
	Name     string `json:"name"`
	Owner    *User  `json:"owner,omitempty"` // TODO: Do we need this field?
}

type FoldersService[T Collection | Gallery] struct {
	sling *sling.Sling
}

func newFoldersService[T Collection | Gallery](sling *sling.Sling) *FoldersService[T] {
	return &FoldersService[T]{
		sling: sling,
	}
}

type FolderParams struct {
	// The user who owns the folder, defaults to current user.
	Username string `url:"username,omitempty"`

	// Sort results by either newest or popular (when querying all folders
	// only).
	// This field is supported only by galleries.
	SortMode string `url:"mode,omitempty"` // values(newest, popular) default: popular
}

type FoldersParams struct {
	// The user to list folders for, if omitted the authenticated user is used.
	Username string `url:"username,omitempty"`

	// The option to include the content count per each collection folder.
	CalculateSize bool `url:"calculate_size,omitempty"`

	// Include first 5 deviations from the folder.
	IncludePreload bool `url:"ext_preload,omitempty"`

	// Filters collections with no deviations if true.
	FilterEmptyFolder bool `url:"filter_empty_folder,omitempty"`
}

type FolderContent struct {
	OffsetResponse[Deviation]
	Name string `json:"name"`
}

// Folder fetches folder contents.
//
// To connect to this endpoint OAuth2 Access Token from the Client Credentials
// Grant, or Authorization Code Grant is required.
//
// The following scopes are required to access this resource:
//
//   - browse
func (s *FoldersService[T]) Folder(folderID uuid.UUID, params *FolderParams, page *OffsetParams) (FolderContent, error) {
	var (
		success FolderContent
		failure Error
	)
	_, err := s.sling.New().Get(folderID.String()).QueryStruct(params).QueryStruct(page).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return FolderContent{}, fmt.Errorf("unable to fetch folder: %w", err)
	}
	return success, nil
}

type usernameParams struct {
	Username string `url:"username,omitempty"`
}

// All gets the "all" view of a users collection/gallery.
//
// To connect to this endpoint OAuth2 Access Token from the Client Credentials
// Grant, or Authorization Code Grant is required.
//
// The following scopes are required to access this resource:
//
//   - browse
func (s *FoldersService[T]) All(username string) (OffsetResponse[Deviation], error) {
	var (
		success OffsetResponse[Deviation]
		failure Error
	)
	params := &usernameParams{Username: username}
	_, err := s.sling.New().Get("all").QueryStruct(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return OffsetResponse[Deviation]{}, fmt.Errorf("unable to fetch all content: %w", err)
	}
	return success, nil
}

// Folders fetches collection folders.
//
// To connect to this endpoint OAuth2 Access Token from the Client Credentials
// Grant, or Authorization Code Grant is required.
//
// The following scopes are required to access this resource:
//
//   - browse
//
// TODO: Support `ext_preload` for collections.
func (s *FoldersService[T]) Folders(params *FoldersParams, page *OffsetParams) (OffsetResponse[T], error) {
	var (
		success OffsetResponse[T]
		failure Error
	)
	_, err := s.sling.New().Get("folders").QueryStruct(params).QueryStruct(page).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return OffsetResponse[T]{}, fmt.Errorf("unable to fetch folders: %w", err)
	}
	return success, nil
}

type CopyDeviationsParams struct {
	TargetFolderID uuid.UUID   `url:"target_folderid,omitempty"`
	DeviationIDs   []uuid.UUID `url:"deviationids,omitempty"`
}

// CopyDeviations copies a list of deviations to a folder destination.
//
// To connect to this endpoint OAuth2 Access Token from the Authorization Code
// Grant is required.
//
// The following scopes are required to access this resource:
//
//   - browse
//   - collection or gallery
func (s *FoldersService[T]) CopyDeviations(params *CopyDeviationsParams) error {
	var (
		failure Error
	)
	_, err := s.sling.New().Get("folders/copy_deviations").QueryStruct(params).Receive(nil, &failure)
	if err := relevantError(err, failure); err != nil {
		return fmt.Errorf("unable to copy deviations: %w", err)
	}
	return nil
}

type CreateFolderParams struct {
	// The name of the folder to create.
	Folder string `url:"folder"`

	// Folder description.
	// This field is supported only by galleries.
	Description string `url:"description,omitempty"`

	// The UUID of the parent gallery if this is a subgallery.
	// This field is supported only by galleries.
	ParentFolderID uuid.UUID `url:"parent_folderid,omitempty"`
}

// Creates new collection folder.
//
// To connect to this endpoint OAuth2 Access Token from the Authorization Code
// Grant is required.
//
// The following scopes are required to access this resource:
//
//   - browse
//   - collection or gallery
func (s *FoldersService[T]) Create(params *CreateFolderParams) (Folder, error) {
	var (
		success Folder
		failure Error
	)
	_, err := s.sling.New().Post("folders/create").BodyForm(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return Folder{}, fmt.Errorf("unable to create folder: %w", err)
	}
	return success, nil
}

type MoveDeviationsParams struct {
	// The UUID of the folder to copy to.
	SourceFolderID uuid.UUID `url:"source_folderid"`

	// The UUID of the folder to copy to.
	TargetFolderID uuid.UUID `url:"target_folderid"`

	// The UUIDs of the deviations.
	DeviationIDs []uuid.UUID `url:"deviationids"`
}

// MoveDeviations moves a list of deviations to a folder destination.
//
// To connect to this endpoint OAuth2 Access Token from the Authorization Code
// Grant is required.
//
// The following scopes are required to access this resource:
//
//   - browse
//   - collection or gallery
func (s *FoldersService[T]) MoveDeviations(params *MoveDeviationsParams) error {
	var (
		failure Error
	)
	_, err := s.sling.New().Post("folders/move_destination").BodyForm(params).Receive(nil, &failure)
	if err := relevantError(err, failure); err != nil {
		return fmt.Errorf("unable to move deviations: %w", err)
	}
	return nil
}

// Remove deletes collection folder.
//
// To connect to this endpoint OAuth2 Access Token from the Authorization Code
// Grant is required.
//
// The following scopes are required to access this resource:
//
//   - browse
//   - collection or gallery
func (s *FoldersService[T]) Remove(folderID uuid.UUID) error {
	var (
		failure Error
	)
	_, err := s.sling.New().Get("folders/remove/").Path(folderID.String()).Receive(nil, &failure)
	if err := relevantError(err, failure); err != nil {
		return fmt.Errorf("unable to remove folder: %w", err)
	}
	return nil
}

type RemoveDeviationsParams struct {
	// The UUID of the folder to remove.
	FolderID uuid.UUID `url:"folderid"`

	// The UUIDs of the deviations.
	DeviationIDs []uuid.UUID `url:"deviationids"`
}

// RemoveDeviations removes a list of deviations from a gallery folder.
//
// To connect to this endpoint OAuth2 Access Token from the Authorization Code
// Grant is required.
//
// The following scopes are required to access this resource:
//
//   - browse
//   - collection or gallery
func (s *FoldersService[T]) RemoveDeviations(params *RemoveDeviationsParams) error {
	var (
		failure Error
	)
	_, err := s.sling.New().Post("folders/remove_deviations").BodyForm(params).Receive(nil, &failure)
	if err := relevantError(err, failure); err != nil {
		return fmt.Errorf("unable to remove deviations: %w", err)
	}
	return nil
}

type UpdateFoldersParams struct {
	// The UUID of the folder to rename.
	FolderID uuid.UUID `url:"folderid"`

	// Folder new name.
	Name string `url:"name,omitempty"`

	// Folder description.
	Description string `url:"description,omitempty"`

	// Folder thumb.
	CoverDeviationID uuid.UUID `url:"cover_deviationid,omitempty"`
}

// Update updates folder.
//
// To connect to this endpoint OAuth2 Access Token from the Authorization Code
// Grant is required.
//
// The following scopes are required to access this resource:
//
//   - browse
//   - collection or gallery
func (s *FoldersService[T]) Update(params *UpdateFoldersParams) error {
	var (
		failure Error
	)
	_, err := s.sling.New().Get("update").QueryStruct(params).Receive(nil, &failure)
	if err := relevantError(err, failure); err != nil {
		return fmt.Errorf("unable to update folders: %w", err)
	}
	return nil
}

type UpdateDeviationOrderParams struct {
	// The UUID of the gallery folder.
	FolderID uuid.UUID `url:"folderid"`

	// The UUID of the deviation.
	DeviationID uuid.UUID `url:"deviationid"`

	// The new position.
	Position int `url:"position"`
}

// UpdateDeviationOrder updates order of deviation in folder.
//
// To connect to this endpoint OAuth2 Access Token from the Authorization Code
// Grant is required.
//
// The following scopes are required to access this resource:
//
//   - browse
//   - collection or gallery
func (s *FoldersService[T]) UpdateDeviationOrder(params *UpdateDeviationOrderParams) error {
	var (
		failure Error
	)
	_, err := s.sling.New().Post("update_deviation_order").BodyForm(params).Receive(nil, &failure)
	if err := relevantError(err, failure); err != nil {
		return fmt.Errorf("unable to update deviations order: %w", err)
	}
	return nil
}

type UpdateOrderParams struct {
	// The UUID of the folder to reposition.
	FolderID uuid.UUID `url:"folderid"`

	// The new position.
	Position int `url:"position"`
}

// UpdateOrder rearranges the position of folders.
//
// To connect to this endpoint OAuth2 Access Token from the Authorization Code
// Grant is required.
//
// The following scopes are required to access this resource:
//
//   - browse
//   - collection or gallery
func (s *FoldersService[T]) UpdateOrder(params *UpdateOrderParams) error {
	var (
		failure Error
	)
	_, err := s.sling.New().Post("update_order").BodyForm(params).Receive(nil, &failure)
	if err := relevantError(err, failure); err != nil {
		return fmt.Errorf("unable to update folders order: %w", err)
	}
	return nil
}
