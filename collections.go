package deviantart

import (
	"fmt"

	"github.com/dghubble/sling"
	"github.com/google/uuid"
)

type Collection struct {
	FolderID    uuid.UUID   `json:"folderid"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Size        uint32      `json:"size,omitempty"`
	Thumb       *Deviation  `json:"thumb,omitempty"`
	Deviations  []Deviation `json:"deviations,omitempty"`
}

type CollectionsService struct {
	sling *sling.Sling
	*FoldersService[Collection]
}

func newCollectionsService(sling *sling.Sling) *CollectionsService {
	base := sling.Path("collections/")
	return &CollectionsService{
		sling:          base,
		FoldersService: newFoldersService[Collection](base.New()),
	}
}

type faveParams struct {
	// ID of the Deviation to favourite.
	DeviationID uuid.UUID `url:"deviationid"`

	// Optional `UUID` of the Collection folder to add the favourite into.
	FolderIDs []uuid.UUID `url:"folderid,omitempty"`
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
//
// To connect to this endpoint OAuth2 Access Token from the Authorization Code
// Grant is required.
//
// The following scopes are required to access this resource:
//
//   - browse
//   - collection
func (s *CollectionsService) Fave(deviationID uuid.UUID, folderIDs ...uuid.UUID) (int, error) {
	var (
		success map[string]any
		failure Error
	)
	_, err := s.sling.New().Post("fave").BodyForm(&faveParams{DeviationID: deviationID, FolderIDs: folderIDs}).Receive(&success, &failure)
	if err != nil {
		return 0, fmt.Errorf("unable to fave the deviation: %w", err)
	}
	return success["favourites"].(int), nil
}

// Unfave removes deviation from favourites.
//
// You can remove deviation from multiple collections at once. If you omit
// `folderID` parameter, it will be removed from Featured collection.
//
// Returns the total number of times this deviation was favourited after the
// unfave event.
//
// If a user has faved their own deviation, unfave can be used to remove the
// deviation from a given folder. Favorite counts are not affected if the
// deviation is owned by the user.
//
// To connect to this endpoint OAuth2 Access Token from the Authorization Code
// Grant is required.
//
// The following scopes are required to access this resource:
//
//   - browse
//   - collection
func (s *CollectionsService) Unfave(deviationID uuid.UUID, folderIDs ...uuid.UUID) (int, error) {
	var (
		success map[string]any
		failure Error
	)
	_, err := s.sling.New().Post("unfave").BodyForm(&faveParams{DeviationID: deviationID, FolderIDs: folderIDs}).Receive(&success, &failure)
	if err != nil {
		return 0, fmt.Errorf("unable to unfave the deviation: %w", err)
	}
	return success["favourites"].(int), nil
}
