package task

import (
	"context"
	"sync"

	"github.com/lowerc4s3/29.07.2025-zipload/downloader"
)

type downloadsMap struct {
	data map[string]DownloadStatus
	mu   sync.Mutex
}

type Task struct {
	statuses downloadsMap
	archiver *downloader.GoArchiver
}

func NewTask() *Task {
	return &Task{
		statuses: downloadsMap{make(map[string]DownloadStatus, 0), sync.Mutex{}},
		archiver: downloader.NewGoArchiver(),
	}
}

func (t *Task) AddDownload(source string) {
	t.statuses.mu.Lock()
	defer t.statuses.mu.Unlock()
	t.statuses.data[source] = StatusDownloading

	go func() {
		download, err := downloader.Download(context.Background(), source)
		if err != nil {
			t.statuses.mu.Lock()
			defer t.statuses.mu.Unlock()
			t.statuses.data[source] = StatusFailed
			return
		}
		t.archiver.AddFile(download.Name, download.Content)

		t.statuses.mu.Lock()
		defer t.statuses.mu.Unlock()
		t.statuses.data[source] = StatusDone
	}()
}

func (t *Task) Status() TaskInfo {
	t.statuses.mu.Lock()
	defer t.statuses.mu.Unlock()

	ready := true
	files := make([]DownloadInfo, 0, len(t.statuses.data))
	for source := range t.statuses.data {
		status := t.statuses.data[source]
		if status == StatusDownloading {
			ready = false
		}
		files = append(files, DownloadInfo{URL: source, Status: status})
	}
	return TaskInfo{files, ready}
}

func (t *Task) FilesAmount() int {
	t.statuses.mu.Lock()
	defer t.statuses.mu.Unlock()
	return len(t.statuses.data)
}

func (t *Task) Finish(ctx context.Context) ([]byte, error) {
	return t.archiver.Finish(ctx)
}
