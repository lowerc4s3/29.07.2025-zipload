package downloader

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path"
	"sync"
)

type DownloadResult struct {
	Ok  *MIMEFile
	Err error
}

type MIMEFile struct {
	Content []byte
	Name    string
	MIME    string
}

func DownloadBatch(ctx context.Context, sources []string) <-chan DownloadResult {
	out := make(chan DownloadResult)
	var wg sync.WaitGroup

	for _, url := range sources {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			file, err := Download(ctx, url)
			if err != nil {
				out <- DownloadResult{Err: err}
			} else {
				out <- DownloadResult{Ok: file}
			}
		}(url)
	}

	// Close channel when all sources were processed
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func Download(ctx context.Context, source string) (*MIMEFile, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", source, nil)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http response: %w", err)
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("response read: %w", err)
	}

	return &MIMEFile{
		Content: content,
		Name:    getName(resp),
		MIME:    getMIME(resp, content),
	}, nil
}

func getName(resp *http.Response) string {
	// Try to get filename from Content-Disposition header
	_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Disposision"))
	if err == nil {
		if headerName, ok := params["filename"]; ok {
			return headerName
		}
	}

	// Fallback to filename from URL
	return path.Base(resp.Request.URL.Path)
}

func getMIME(resp *http.Response, content []byte) string {
	// Try to get MIME from Content-Type header
	headerMIME, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err == nil {
		return headerMIME
	}

	// Try to get MIME from URL extension if there is one
	extMIME := mime.TypeByExtension(path.Ext(resp.Request.URL.Path))
	if extMIME != "" {
		return extMIME
	}

	// Fallback to content scanning based detection
	return http.DetectContentType(content)
}
