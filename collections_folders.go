package deviantart

import (
	"fmt"

	"github.com/dghubble/sling"
	"github.com/google/uuid"
)

type Folder struct {
	FolderID uuid.UUID `json:"folderid"`
	Name     string    `json:"name"`
}

type foldersService struct {
	sling *sling.Sling
}

func newFoldersService(sling *sling.Sling) *foldersService {
	return &foldersService{
		sling: sling,
	}
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

	// The pagination offset.
	Offset int `url:"offset,omitempty"`

	// The pagination limit.
	Limit int `url:"limit,omitempty"`
}

// Folders fetches collection folders.
func (s *foldersService) Folders(params *FolderParams) (FolderContent, error) {
	var (
		success FolderContent
		failure Error
	)
	_, err := s.sling.New().Get("folders").QueryStruct(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return FolderContent{}, fmt.Errorf("unable to fetch folders: %w", err)
	}
	return success, nil
}

type CopyDeviationsParams struct {
	TargetFolderID uuid.UUID   `url:"target_folderid,omitempty"`
	DeviationIDs   []uuid.UUID `url:"deviationids,omitempty"`
}

// CopyDeviations copies a list of deviations to a folder destination.
//
// Requires Authorization Code grant.
func (s *foldersService) CopyDeviations(params *CopyDeviationsParams) error {
	var (
		success map[string]any
		failure Error
	)
	_, err := s.sling.New().Get("folders/copy_deviations").QueryStruct(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return fmt.Errorf("unable to copy deviations: %w", err)
	}
	return nil
}

type createFolderParams struct {
	Folder string `url:"folder"`
}

// Creates new collection folder.
//
// Requires Authorization Code grant.
func (s *foldersService) Create(folderName string) (Folder, error) {
	var (
		success Folder
		failure Error
	)
	_, err := s.sling.New().Post("folders/create").BodyForm(&createFolderParams{Folder: folderName}).Receive(&success, &failure)
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
	DeviationIDS []uuid.UUID `url:"deviationids"`
}

// MoveDeviations moves a list of deviations to a folder destination.
func (s *foldersService) MoveDeviations(params *MoveDeviationsParams) error {
	var (
		failure Error
	)
	_, err := s.sling.New().Post("folder/move_destination").BodyForm(params).Receive(nil, &failure)
	if err := relevantError(err, failure); err != nil {
		return fmt.Errorf("unable to move deviations: %w", err)
	}
	return nil
}

// Remove deletes collection folder.
//
// Requires Authorization Code grant.
func (s *foldersService) Remove(folderID uuid.UUID) error {
	var (
		failure Error
	)
	_, err := s.sling.New().Get("remove/").Path(folderID.String()).Receive(nil, &failure)
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
func (s *foldersService) RemoveDeviations(params *RemoveDeviationsParams) error {
	var (
		failure Error
	)
	_, err := s.sling.New().Post("/remove_deviations").BodyForm(params).Receive(nil, &failure)
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
func (s *foldersService) Update(params *UpdateFoldersParams) error {
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
func (s *foldersService) UpdateDeviationOrder(params *UpdateDeviationOrderParams) error {
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
func (s *foldersService) UpdateOrder(params *UpdateOrderParams) error {
	var (
		failure Error
	)
	_, err := s.sling.New().Post("update_order").BodyForm(params).Receive(nil, &failure)
	if err := relevantError(err, failure); err != nil {
		return fmt.Errorf("unable to update folders order: %w", err)
	}
	return nil
}
