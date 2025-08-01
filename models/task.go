package models

import (
	"encoding/json"
)

// NOTE: Using UUIDs as using plain uint IDs is not really safe
type DownloadTask struct {
	Files []FileTask `json:"files"`
	Link  *string    `json:"link"`
}

type FileTask struct {
	Name            string     `json:"string"`
	BytesDownloaded uint       `json:"bytes_downloaded"`
	BytesLeft       uint       `json:"bytes_left"`
	Status          TaskStatus `json:"status"`
}

type TaskStatus uint8

const (
	Downloading TaskStatus = iota
	Done
	Failed
)

func (s TaskStatus) String() string {
	switch s {
	case Downloading:
		return "downloading"
	case Done:
		return "done"
	case Failed:
		return "failed"
	default:
		panic("FileTaskStatus internal number is out of bounds")
	}
}

func (s TaskStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}
