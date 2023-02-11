package deviantart

import "fmt"

type Error struct {
	StatusResponse
	Text        string            `json:"error"`
	Code        int               `json:"code"`
	Description string            `json:"error_description"`
	Details     map[string]string `json:"error_details"`
}

func (e Error) Error() string {
	// TODO: Append details.
	return fmt.Sprintf("%s: %s", e.Text, e.Description)
}

func relevantError(httpError error, apiError Error) error {
	if httpError != nil {
		return httpError
	} else if apiError.Text != "" {
		return apiError
	}
	return nil
}
