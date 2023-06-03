package deviantart

import (
	"fmt"

	"github.com/dghubble/sling"
	"github.com/google/uuid"
)

type collectionsService struct {
	sling   *sling.Sling
	Folders *foldersService
}

func newCollectionsService(sling *sling.Sling) *collectionsService {
	base := sling.Path("collections/")
	return &collectionsService{
		sling:   base,
		Folders: newFoldersService(base.New()),
	}
}

type CollectionsFolderParams struct {
	// The user who owns the folder, defaults to current user.
	Username string `url:"username,omitempty"`

	// The pagination offset.
	Offset int `url:"offset,omitempty"`

	// The pagination limit.
	Limit int `url:"limit,omitempty"`

	WithSession bool `url:"with_session,omitempty"`
}

type FolderContent struct {
	Name string `json:"name,omitempty"`
	OffsetResponse[Deviation]
}

// Folder fetches collection folder contents.
func (s *collectionsService) Folder(folderID uuid.UUID, params *CollectionsFolderParams) (FolderContent, error) {
	var (
		success FolderContent
		failure Error
	)
	_, err := s.sling.New().Get(folderID.String()).QueryStruct(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return FolderContent{}, fmt.Errorf("unable to fetch folder: %w", err)
	}
	return success, nil
}

// All fetches all deviations in user's collection.
func (s *collectionsService) All(params *CollectionsFolderParams) (OffsetResponse[Deviation], error) {
	var (
		success OffsetResponse[Deviation]
		failure Error
	)
	_, err := s.sling.New().Get("all").QueryStruct(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return OffsetResponse[Deviation]{}, fmt.Errorf("unable to fetch all content: %w", err)
	}
	return success, nil
}

type FaveParams struct {
	// ID of the Deviation to favourite.
	DeviationID uuid.UUID `url:"deviationid"`

	// Optional `UUID` of the Collection folder to add the favourite into.
	FolderID []uuid.UUID `url:"folderid,optional"`
}

// Fave adds deviation to favourites.
//
// You can add deviation to multiple collections at once. If you omit `folderID`
// parameter, it will be added to Featured collection.
//
// Returns the total number of times this deviation was favourited after the
// fave event.
//
// Users can fave their own deviations, when this happens the fave is not
// counted but the item is added to the requested folder.
func (s *collectionsService) Fave(params *FaveParams) (int, error) {
	var (
		success map[string]any
		failure Error
	)
	_, err := s.sling.New().Post("fave").BodyForm(params).Receive(&success, &failure)
	if err != nil {
		return 0, fmt.Errorf("unable to fave the deviation: %w", err)
	}
	return success["favourites"].(int), nil
}

// Unfave removes deviation from favourites.
//
// You can remove deviation from multiple collections at once. If you omit
// folderid parameter, it will be removed from Featured collection.
//
// Returns the total number of times this deviation was favourited after the
// unfave event.
//
// If a user has faved their own deviation, unfave can be used to remove the
// deviation from a given folder. Favorite counts are not affected if the
// deviation is owned by the user.
func (s *collectionsService) Unfave(params *FaveParams) (int, error) {
	var (
		success map[string]any
		failure Error
	)
	_, err := s.sling.New().Post("unfave").BodyForm(params).Receive(&success, &failure)
	if err != nil {
		return 0, fmt.Errorf("unable to unfave the deviation: %w", err)
	}
	return success["favourites"].(int), nil
}
