package deviantart

import (
	"fmt"

	"github.com/dghubble/sling"
	"github.com/google/uuid"
)

type galleryService struct {
	sling   *sling.Sling
	Folders *foldersService
}

func newGalleryService(sling *sling.Sling) *galleryService {
	base := sling.Path("gallery/")
	return &galleryService{
		sling:   base,
		Folders: newFoldersService(base.New()),
	}
}

type GalleryFolderParams struct {
	// The user to query, defaults to current user.
	Username string `url:"username,omitempty"`

	// Sort results by either newest or popular (when querying all folders
	// only).
	Mode string `url:"mode,omitempty"` // values(newest, popular) default: popular

	// The pagination offset.
	Offset int `url:"offset,omitempty"`

	// The pagination limit.
	Limit int `url:"limit,omitempty"`

	WithSession bool `url:"with_session,omitempty"`
}

// Folder fetches gallery folder contents.
func (s *galleryService) Folder(folderID uuid.UUID, params *GalleryFolderParams) (FolderContent, error) {
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

// All gets the "all" view of a users gallery.
func (s *galleryService) All(params *CollectionsFolderParams) (FolderContent, error) {
	var (
		success FolderContent
		failure Error
	)
	_, err := s.sling.New().Get("all").QueryStruct(params).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return FolderContent{}, fmt.Errorf("unable to fetch all content: %w", err)
	}
	return success, nil
}
