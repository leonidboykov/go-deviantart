package deviantart

import (
	"fmt"

	"github.com/google/uuid"
)

type DownloadResponse struct {
	Src      string `json:"src"`
	FileName string `json:"filename"`
	Width    uint32 `json:"width"`
	Height   uint32 `json:"height"`
	FileSize uint32 `json:"filesize"`
}

// Download fetches the original file download (if allowed).
func (s *deviationService) Download(deviationID uuid.UUID) (DownloadResponse, error) {
	var (
		success DownloadResponse
		failure Error
	)
	_, err := s.sling.New().Get(deviationID.String()).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return DownloadResponse{}, fmt.Errorf("unable to fetch download data: %w", err)
	}
	return success, nil
}
