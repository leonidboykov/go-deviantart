package deviantart

import "fmt"

// Placebo call confirms access_token is valid.
func (c *Client) Placebo() error {
	var (
		success StatusResponse
		failure Error
	)
	_, err := c.base.New().Get("placebo/").Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return fmt.Errorf("unable to validate access_token: %w", err)
	}
	return nil
}
