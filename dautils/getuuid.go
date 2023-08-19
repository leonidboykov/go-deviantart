package dautils

import (
	"errors"
	"io"
	"net/http"
	"regexp"

	"github.com/google/uuid"
)

var re = regexp.MustCompile(`content="DeviantArt:\/\/deviation\/(?P<uuid>[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})`)

func GetDeviationUUIDByURL(url string) (uuid.UUID, error) {
	resp, err := http.Get(url)
	if err != nil {
		return uuid.UUID{}, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return uuid.UUID{}, err
	}
	match := re.FindSubmatch(data)

	if len(match) != 2 {
		return uuid.UUID{}, errors.New("parsing issues, len(match) != 2")
	}

	return uuid.ParseBytes(match[1])
}
