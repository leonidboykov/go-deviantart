package deviantart

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"mime/multipart"
	"strings"

	"github.com/google/go-querystring/query"
)

type StashSubmitParams struct {
	// The title of the submission.
	Title string `url:"title,omitempty"`

	// Additional information about the submission provided by the author.
	Description string `url:"artist_comments,omitempty"`

	// An array of tags describing the submission. Letters, numbers and
	// underscore only.
	Tags []string `url:"tags,brackets,omitempty"`

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

type SubmitResponse struct {
	StatusResponse
	ItemID  int64  `json:"itemid"`
	Stack   string `json:"stack,omitempty"`
	StackID int64  `json:"stackid,omitempty"`
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
func (s *StashService) Submit(params *StashSubmitParams, files ...fs.File) (SubmitResponse, error) {
	var (
		success SubmitResponse
		failure Error
	)
	provider, err := newMultipartBodyProvider(params, files...)
	if err != nil {
		return SubmitResponse{}, fmt.Errorf("prepare submit data: %w", err)
	}
	_, err = s.sling.New().Post("submit").BodyProvider(provider).Receive(&success, &failure)
	if err := relevantError(err, failure); err != nil {
		return SubmitResponse{}, fmt.Errorf("unable to submit file to sta.sh: %w", err)
	}
	return success, nil
}

type multipartBodyProvider struct {
	reader      io.Reader
	contentType string
}

func newMultipartBodyProvider(params *StashSubmitParams, files ...fs.File) (*multipartBodyProvider, error) {
	values, err := query.Values(params)
	if err != nil {
		return nil, fmt.Errorf("encode params: %w", err)
	}
	if len(files) == 0 {
		// Fallback to form body provider.
		return &multipartBodyProvider{
			reader:      strings.NewReader(values.Encode()),
			contentType: "application/x-www-form-urlencoded",
		}, nil
	}
	buf := new(bytes.Buffer)
	mp := multipart.NewWriter(buf)
	for _, file := range files {
		fs, err := file.Stat()
		if err != nil {
			return nil, fmt.Errorf("get file stat: %w", err)
		}
		part, err := mp.CreateFormFile(fs.Name(), fs.Name()) // TODO: Is it correct?
		if err != nil {
			return nil, fmt.Errorf("create from file: %w", err)
		}
		io.Copy(part, file)
	}
	for key, vals := range values {
		for _, val := range vals {
			mp.WriteField(key, val)
		}
	}
	return &multipartBodyProvider{
		reader:      buf,
		contentType: mp.FormDataContentType(),
	}, nil
}

// ContentType returns the Content-Type of the body.
func (p *multipartBodyProvider) ContentType() string {
	return p.contentType
}

// Body returns the io.Reader body.
func (p *multipartBodyProvider) Body() (io.Reader, error) {
	return p.reader, nil
}
