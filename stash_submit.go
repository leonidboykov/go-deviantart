package deviantart

import "fmt"

type StashSubmitParams struct {
	// The title of the submission.
	Title string `url:"title,omitempty"`

	// Additional information about the submission provided by the author.
	Description string `url:"artist_comments,omitempty"`

	// An array of tags describing the submission. Letters, numbers and
	// underscore only.
	Tags []string `url:"tags,omitempty"`

	// A link to the original, in case the artwork has already been posted
	// elsewhere. This field can be restricted with a whitelist by editing your
	// deviantART app.
	OriginalURL string `url:"original_url,omitempty"`

	// Is the submission being worked on currently. You can use this flag to
	// warn users that the item is being edited and may change if they reload.
	// Note this does NOT provide any type of locking.
	IsDirty bool `url:"is_dirty,omitempty"`

	// The id of an existing Sta.sh submission. This can be used to overwrite
	// files and /metadata of existing submissions. If you make a new API call
	// containing files, the files that were previously associated with the
	// artwork will be replaced by the new ones.
	ItemID int64 `url:"itemid,omitempty"`

	// The name of the stack to create and place the new submission in. Applies
	// to new submissions only. (Ignored if `stackid` is set).
	Stack string `url:"stack,omitempty"`

	// The id of the stack to create and place the new submission in. Applies to
	// new submissions only.
	StackID int64 `url:"stackid,omitempty"`
}

// Submit submits files to Sta.sh or modify existing files.
//
// It can receive files in any format. Some formats like JPG, PNG, GIF, HTML or
// plain text can be viewed directly on Sta.sh and DeviantArt. Other file types
// are made available for download and may have a preview image.
//
// To connect to this endpoint OAuth2 Access Token from the Authorization Code
// Grant is required.
//
// The following scopes are required to access this resource:
//
//   - stash
func (s *stashService) Submit(params *StashSubmitParams) error {
	var (
		success map[string]any
		failure map[string]any
	)
	// TODO: Upload file.
	_, err := s.sling.New().Post("submit").BodyForm(params).Receive(&success, &failure)
	fmt.Println("=== SUCCESS ===")
	fmt.Println(success)
	fmt.Println("")
	fmt.Println("=== FAILURE ===")
	fmt.Println(failure)
	fmt.Println("")
	fmt.Println("===  ERROR  ===")
	fmt.Println(err)

	req, _ := s.sling.New().Post("submit").BodyForm(params).BodyForm(nil).Request()
	fmt.Println(req.URL.String())
	return err
}
