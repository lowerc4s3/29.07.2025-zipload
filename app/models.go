package app

import "github.com/lowerc4s3/29.07.2025-zipload/task"

type AppendTaskRequest struct {
	ID     string `json:"id"`
	Source string `json:"source"`
}

type TaskResponse struct {
	Files []task.DownloadInfo `json:"files"`
	Link  *string             `json:"link,omitempty"`
}
