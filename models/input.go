package models

type AppendTaskRequest struct {
	ID     string `json:"id"`
	Source string `json:"source"`
}
