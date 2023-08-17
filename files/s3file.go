package files

import (
	"errors"
	"fmt"
	s3 "github.com/ipfs/boxo/s3connection"
	"io"
	"net/url"
	"os"
)

// S3File is an implementation of File which reads it
// from a S3 URI. A READ request will be performed
// against the source when calling Read().
type S3File struct {
	body          io.ReadCloser
	url           *url.URL
	contentLength int64
	s3conn        s3.S3Backend
}

// NewS3File creates a S3File with the given URL, which
// will be used to perform the GET request on Read().
func NewS3File(s3conn s3.S3Backend, url *url.URL) *S3File {
	return &S3File{
		url:    url,
		s3conn: s3conn,
	}
}

func (s3f *S3File) start() error {
	if s3f.body == nil {
		fileSize, err := s3f.s3conn.FileInfo(s3f.url)
		if err != nil {
			return err
		}
		body, size, err := s3f.s3conn.Download(s3f.url)
		if err != nil {
			return err
		}
		if size != fileSize {
			return fmt.Errorf("S3 file size incorrect")
		}

		s3f.body = body
		s3f.contentLength = size
		return nil
	}
	return nil
}

// Read reads the File from its web location. On the first
// call to Read, a GET request will be performed against the
// S3File's URL, using Go's default HTTP client. Any further
// reads will keep reading from the HTTP Request body.
func (s3f *S3File) Read(b []byte) (int, error) {
	if err := s3f.start(); err != nil {
		return 0, err
	}
	return s3f.body.Read(b)
}

// Close closes the WebFile (or the request body).
func (s3f *S3File) Close() error {
	if s3f.body == nil {
		return nil
	}
	return s3f.body.Close()
}

// TODO: implement
func (s3f *S3File) Seek(offset int64, whence int) (int64, error) {
	return 0, ErrNotSupported
}

func (s3f *S3File) Size() (int64, error) {
	if err := s3f.start(); err != nil {
		return 0, err
	}
	if s3f.contentLength < 0 {
		return -1, errors.New("Content-Length hearer was not set")
	}

	return s3f.contentLength, nil
}

func (s3f *S3File) AbsPath() string {
	return s3f.url.String()
}

func (s3f *S3File) Stat() os.FileInfo {
	return nil
}

var _ File = &S3File{}
var _ FileInfo = &S3File{}
