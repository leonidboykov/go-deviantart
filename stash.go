package deviantart

import (
	"github.com/dghubble/sling"
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
