package task

import "encoding/json"

type TaskResponse struct {
	Files []DownloadInfo `json:"files"`
	Link  *string        `json:"link"`
}

type TaskInfo struct {
	Files []DownloadInfo
	Ready bool
}

type DownloadInfo struct {
	URL    string         `json:"url"`
	Status DownloadStatus `json:"status"`
}

type DownloadStatus uint8

const (
	StatusDownloading DownloadStatus = iota
	StatusDone
	StatusFailed
)

func (s DownloadStatus) String() string {
	switch s {
	case StatusDownloading:
		return "downloading"
	case StatusDone:
		return "done"
	case StatusFailed:
		return "failed"
	default:
		panic("FileTaskStatus internal number is out of bounds")
	}
}

func (s DownloadStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}
