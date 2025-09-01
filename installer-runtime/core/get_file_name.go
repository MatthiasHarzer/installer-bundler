package core

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
)

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func getFileNameFromHeader(header http.Header, fileURL string) (string, error) {
	var filename string
	contentDisposition := header.Get("Content-Disposition")
	if contentDisposition != "" {
		_, err := fmt.Sscanf(contentDisposition, "attachment; filename=%q", &filename)
		if err == nil {
			return filename, nil
		}
	}

	// Fallback to extracting the filename from the URL
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return "", err
	}
	pathSegments := strings.Split(parsedURL.Path, "/")
	if len(pathSegments) > 0 {
		filename = pathSegments[len(pathSegments)-1]
	}

	hasExtension := strings.Contains(filename, ".")
	if hasExtension {
		return filename, nil
	}

	contentType := header.Get("Content-Type")
	mime := mimetype.Lookup(contentType)

	if mime != nil {
		return fmt.Sprintf("%s_%d%s", randomString(10), time.Now().UnixMilli(), mime.Extension()), nil
	}

	return "", fmt.Errorf("could not determine filename from URL or headers")
}
