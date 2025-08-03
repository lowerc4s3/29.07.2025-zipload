package batch

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/lowerc4s3/29.07.2025-zipload/downloader"
)

var (
	ErrTooManySources   = errors.New("amount of provided sources exceeds max allowed amount")
	ErrForbiddenMIME    = errors.New("got resource with forbidden type")
	ErrPartitialArchive = errors.New("some files couldn't be processed")
	ErrDownload         = errors.New("all files couldn't be processed")
)

type BatchSerivce struct {
	allowedMIMEs map[string]struct{}
	maxSources   int
}

func NewBatchService(allowedMIMEs []string, maxSources int) *BatchSerivce {
	mimemap := make(map[string]struct{}, len(allowedMIMEs))
	for _, mime := range allowedMIMEs {
		mimemap[mime] = struct{}{}
	}
	return &BatchSerivce{
		allowedMIMEs: mimemap,
		maxSources:   maxSources,
	}
}

func (s *BatchSerivce) DownloadAll(ctx context.Context, sources []string) ([]byte, error) {
	// Abort if too many sources were requested
	if len(sources) > s.maxSources {
		return nil, ErrTooManySources
	}

	archiver := downloader.NewArchiver()
	inCh := make(chan *downloader.MIMEFile) // File channel
	outCh := make(chan []byte)              // Final archive channel
	errCh := make(chan error)               // Error channel

	// Check if channel was closed manually
	inClosed := false
	defer func() {
		if !inClosed {
			close(inCh)
		}
	}()

	// Spawn archiver on a different goroutine to dynamically process files
	go func() {
		defer close(errCh)
		defer close(outCh)

		// Archive files on demand
		for file := range inCh {
			if err := archiver.AddFile(file.Name, file.Content); err != nil {
				errCh <- err
				return
			}
		}

		// Finalize archive
		result, err := archiver.FinishArchive()
		if err != nil {
			errCh <- err
		} else {
			outCh <- result
		}
	}()

	filesCounter := 0                      // Track amount of successful downloads
	errs := make([]error, 0, len(sources)) // Store all errors in the slice

	// Create another context to stop download if returning early
	downloadCtx, downloadCancel := context.WithCancel(ctx)
	defer downloadCancel()
	for result := range downloader.DownloadBatch(downloadCtx, sources) {
		select {
		case err := <-errCh:
			// If archiver failed, abort downloads and return...
			return nil, fmt.Errorf("archiving file: %w", err)
		default:
			// ...otherwise proceed with downloads
		}

		// If download failed, skip it
		if result.Err != nil {
			errs = append(errs, result.Err)
			continue
		}

		// If downloaded file's type is forbidden, skip it
		if _, ok := s.allowedMIMEs[result.Ok.MIME]; !ok {
			errs = append(errs, fmt.Errorf("%w: %v", ErrForbiddenMIME, result.Ok.MIME))
			continue
		}

		filesCounter++
		inCh <- result.Ok
	}
	close(inCh)
	inClosed = true

	// Return error if all downloads were failed
	if len(errs) != 0 && filesCounter == 0 {
		if slices.ContainsFunc(errs, func(e error) bool { return errors.Is(e, ErrForbiddenMIME) }) {
			return nil, ErrForbiddenMIME
		}
		return nil, ErrDownload
	}

	select {
	case err := <-errCh:
		return nil, fmt.Errorf("archive finishing: %w", err)
	case archive := <-outCh:
		if len(errs) != 0 {
			return archive, ErrPartitialArchive
		}
		return archive, nil
	}
}
