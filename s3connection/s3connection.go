package s3connection

import (
	"io"
	"net/url"
)

type S3Backend interface {
	FileInfo(url.URL) (int64, error)
	Download(url.URL) io.ReadCloser
}
