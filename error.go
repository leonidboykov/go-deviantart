package deviantart

type Error struct {
	Text        string            `json:"error"`
	Code        int               `json:"code"`
	Description string            `json:"error_description"`
	Details     map[string]string `json:"error_details"`
}

func (e Error) Error() string {
	// TODO: Append details.
	return e.Description
}

func relevantError(httpError error, apiError Error) error {
	if httpError != nil {
		return httpError
	} else if apiError.Text != "" {
		return apiError
	}
	return nil
}
