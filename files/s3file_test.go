package files

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type mockS3Client struct {
	s3iface.S3API
}

func (m *mockS3Client) AbortMultipartUpload(input *s3.AbortMultipartUploadInput) (*s3.AbortMultipartUploadOutput, error) {
	// mock response/functionality
	return nil, nil
}

func TestS3File(t *testing.T) {
	const content = "Hello world!"
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, content)
	}))
	defer s.Close()

	u, err := url.Parse(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	wf := NewWebFile(u)
	body, err := io.ReadAll(wf)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != content {
		t.Fatalf("expected %q but got %q", content, string(body))
	}
}

func TestS3File_notFound(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "File not found.", http.StatusNotFound)
	}))
	defer s.Close()

	u, err := url.Parse(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	wf := NewWebFile(u)
	_, err = io.ReadAll(wf)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestS3FileSize(t *testing.T) {
	body := "Hello world!"
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, body)
	}))
	defer s.Close()

	u, err := url.Parse(s.URL)
	if err != nil {
		t.Fatal(err)
	}

	// Read size before reading file.

	wf1 := NewWebFile(u)
	if size, err := wf1.Size(); err != nil {
		t.Error(err)
	} else if int(size) != len(body) {
		t.Errorf("expected size to be %d, got %d", len(body), size)
	}

	actual, err := io.ReadAll(wf1)
	if err != nil {
		t.Fatal(err)
	}
	if string(actual) != body {
		t.Fatal("should have read the web file")
	}

	wf1.Close()

	// Read size after reading file.

	wf2 := NewWebFile(u)
	actual, err = io.ReadAll(wf2)
	if err != nil {
		t.Fatal(err)
	}
	if string(actual) != body {
		t.Fatal("should have read the web file")
	}

	if size, err := wf2.Size(); err != nil {
		t.Error(err)
	} else if int(size) != len(body) {
		t.Errorf("expected size to be %d, got %d", len(body), size)
	}
}
