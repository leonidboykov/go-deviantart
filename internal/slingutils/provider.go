package slingutils

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"mime/multipart"
	"strings"

	"github.com/dghubble/sling"
	"github.com/google/go-querystring/query"
)

var _ sling.BodyProvider = new(multipartBodyProvider)

type multipartBodyProvider struct {
	reader      io.Reader
	contentType string
}

func NewMultipartProvider(formData any, file fs.File) (*multipartBodyProvider, error) {
	values, err := query.Values(formData)
	if err != nil {
		return nil, fmt.Errorf("encode params: %w", err)
	}
	if file == nil {
		return &multipartBodyProvider{
			reader:      strings.NewReader(values.Encode()),
			contentType: "application/x-www-form-urlencoded",
		}, nil
	}

	buf := new(bytes.Buffer)
	mp := multipart.NewWriter(buf)
	fs, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("get file stat: %w", err)
	}
	part, err := mp.CreateFormFile(fs.Name(), fs.Name()) // TODO: Is it correct?
	if err != nil {
		return nil, fmt.Errorf("create from file: %w", err)
	}
	io.Copy(part, file)
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
