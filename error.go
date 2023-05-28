package deviantart

import "fmt"

type Error struct {
	StatusResponse

	// The error type.
	Type string `json:"error"`

	// The error message.
	Description string `json:"error_description"`

	// An additional endpoint specific error code.
	Details map[string]string `json:"error_details"`

	// An additional endpoint specific error code.
	Code int `json:"error_code"`
}

func (e Error) Error() string {
	// TODO: Append details.
	return fmt.Sprintf("%s: %s", e.Type, e.Description)
}

func relevantError(httpError error, apiError Error) error {
	if httpError != nil {
		return httpError
	} else if apiError.Type != "" {
		return apiError
	}
	return nil
}
