package utils

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/ledongthuc/pdf"
)

// ExtractTextFromPDF extracts plain text from a PDF file at the given path or URL.
func ExtractTextFromPDF(pathOrURL string) (string, error) {
	var readerAt io.ReaderAt
	var size int64

	// Check if it is a URL
	if strings.HasPrefix(pathOrURL, "http://") || strings.HasPrefix(pathOrURL, "https://") {
		// Download to temp file
		resp, err := http.Get(pathOrURL)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		tmpFile, err := os.CreateTemp("", "pdf-*.pdf")
		if err != nil {
			return "", err
		}
		defer os.Remove(tmpFile.Name()) // Clean up
		defer tmpFile.Close()

		size, err = io.Copy(tmpFile, resp.Body)
		if err != nil {
			return "", err
		}
		readerAt = tmpFile
	} else {
		// Assume local file
		// If path starts with "/", remove it to make it relative to project root if needed
		// But usually absolute path is safer or relative to CWD.
		// The app runs from project root.
		// Also handle cases where path starts with "storage/" vs "/storage/"
		cleanPath := strings.TrimPrefix(pathOrURL, "/")
		if _, err := os.Stat(cleanPath); os.IsNotExist(err) {
			// try absolute or without trimming?
			// Let's rely on standard os.Open
			cleanPath = pathOrURL
			// If starts with /, and we are on windows?
			// On windows file path d:\... is absolute.
			// If web path /storage/..., it translates to ./storage/...
			if strings.HasPrefix(cleanPath, "/storage") || strings.HasPrefix(cleanPath, "storage") {
				// Ensure we look in current directory
				if strings.HasPrefix(cleanPath, "/") {
					cleanPath = "." + cleanPath
				}
			}
		}

		f, err := os.Open(cleanPath)
		if err != nil {
			return "", err
		}
		defer f.Close()

		fs, err := f.Stat()
		if err != nil {
			return "", err
		}
		size = fs.Size()
		readerAt = f
	}

	r, err := pdf.NewReader(readerAt, size)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	for pageIndex := 1; pageIndex <= r.NumPage(); pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		// Extract text
		text, err := p.GetPlainText(nil)
		if err != nil {
			continue
		}
		buf.WriteString(text)
		buf.WriteString("\n")
	}

	return buf.String(), nil
}
