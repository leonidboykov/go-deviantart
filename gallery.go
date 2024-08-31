package deviantart

import (
	"github.com/dghubble/sling"
	"github.com/google/uuid"
)

type Gallery struct {
	FolderID    uuid.UUID `json:"folderid"`
	Parent      uuid.UUID `json:"parent,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description"`

	// Content count per each gallery folder. This field is only presented if
	// calculate_size param is true.
	Size              int64              `json:"size,omitempty"`
	Thumb             *Deviation         `json:"thumb,omitempty"`
	PremiumFolderData *PremiumFolderData `json:"premium_folder_data,omitempty"`
	HasSubfolders     bool               `json:"has_subfolders"`
	Subfolders        []Gallery          `json:"subfolders,omitempty"`
	Deviations        []Deviation        `json:"deviations,omitempty"`
}

type GalleryService struct {
	sling *sling.Sling
	*FoldersService[Gallery]
}

func newGalleryService(sling *sling.Sling) *GalleryService {
	base := sling.Path("gallery/")
	return &GalleryService{
		sling:          base,
		FoldersService: newFoldersService[Gallery](base.New()),
	}
}
