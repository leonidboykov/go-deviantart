package deviantart

import (
	"github.com/dghubble/sling"
	"github.com/google/uuid"
)

type Gallery struct {
	FolderID          uuid.UUID   `json:"folderid"`
	Parent            uuid.UUID   `json:"parent,omitempty"`
	Name              string      `json:"name"`
	Description       string      `json:"folder"`
	Size              uint32      `json:"size,omitempty"`
	Thumb             *Deviation  `json:"thumb"`
	PremiumFolderData any         `json:"premium_folder_data,omitempty"` // premium_folder_data
	HasSubfolders     bool        `json:"has_subfolders"`
	Subfolders        []Gallery   `json:"subfolders,omitempty"`
	Deviations        []Deviation `json:"deviations,omitempty"`
}

type galleryService struct {
	sling   *sling.Sling
	Folders *foldersService[Gallery]
}

func newGalleryService(sling *sling.Sling) *galleryService {
	base := sling.Path("gallery/")
	return &galleryService{
		sling:   base,
		Folders: newFoldersService[Gallery](base.New()),
	}
}
