package app

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lowerc4s3/29.07.2025-zipload/batch"
	"github.com/lowerc4s3/29.07.2025-zipload/task"
)

type BatchService interface {
	DownloadAll(ctx context.Context, sources []string) ([]byte, error)
}

type TaskService interface {
	Create(ctx context.Context) (uuid.UUID, error)
	AppendFile(ctx context.Context, taskID uuid.UUID, source string) error
	Check(ctx context.Context, taskID uuid.UUID) (task.TaskInfo, error)
	PopArchive(ctx context.Context, taskID uuid.UUID) ([]byte, error)
}

func mapError(err error) error {
	switch {
	case errors.Is(err, batch.ErrForbiddenMIME):
		return HttpErrForbiddenMime
	case errors.Is(err, task.ErrTaskNotFound):
		return HttpErrTaskNotFound
	case errors.Is(err, batch.ErrDownload):
		return HttpErrDownload
	case errors.Is(err, task.ErrTooManyTasks):
		return HttpErrTooManyTasks
	case errors.Is(err, task.ErrTooManyTaskFiles):
		return HttpErrTooManyTaskFiles
	case errors.Is(err, batch.ErrTooManySources):
		return HttpErrTooManySources
	default:
		return err
	}
}

type Handler struct {
	batch BatchService
	task  TaskService
}

func NewHandler(batch BatchService, task TaskService) *Handler {
	return &Handler{
		batch: batch,
		task:  task,
	}
}

func (h *Handler) DownloadBatch(c echo.Context) error {
	request := new(batch.Batch)
	if err := c.Bind(request); err != nil {
		return HttpErrMalformedRequest.WithInternal(err)
	}

	archive, err := h.batch.DownloadAll(c.Request().Context(), request.Sources)
	status := http.StatusOK
	if err != nil {
		if !errors.Is(err, batch.ErrPartitialArchive) {
			return mapError(err)
		}

		// HACK: Using HTTP 206 code for partitial success
		// which is not idiomatic use of this status
		// but does the job of telling the client that
		// archive doesn't contain some files
		// while preserving same response structure
		status = http.StatusPartialContent
	}

	return c.Blob(status, "application/zip", archive)
}

func (h *Handler) CreateTask(c echo.Context) error {
	id, err := h.task.Create(c.Request().Context())
	if err != nil {
		return HttpErrTooManyTasks.WithInternal(err)
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"id": id.String(),
	})
}

func (h *Handler) AppendTask(c echo.Context) error {
	request := new(AppendTaskRequest)
	if err := c.Bind(request); err != nil {
		return HttpErrMalformedRequest.WithInternal(err)
	}

	id, err := uuid.Parse(request.ID)
	if err != nil {
		return HttpErrMalformedID.WithInternal(err)
	}

	if err := h.task.AppendFile(c.Request().Context(), id, request.Source); err != nil {
		return mapError(err)
	}

	return c.NoContent(http.StatusAccepted)
}

func (h *Handler) CheckTask(c echo.Context) error {
	rawID := c.Param("taskid")
	if rawID == "" {
		return HttpErrMalformedRequest.WithInternal(errors.New("ID parameter wasn't provided"))
	}

	id, err := uuid.Parse(rawID)
	if err != nil {
		return HttpErrMalformedID.WithInternal(err)
	}

	status, err := h.task.Check(c.Request().Context(), id)
	if err != nil {
		return mapError(err)
	}

	resp := TaskResponse{
		Files: status.Files,
	}
	// If task is finished, then provide a download link in response
	if status.Ready {
		downloadLink := c.Echo().Reverse("downloadTask", id.String())
		resp.Link = &downloadLink
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) DownloadTask(c echo.Context) error {
	rawID := c.Param("taskid")
	if rawID == "" {
		return HttpErrMalformedRequest.WithInternal(errors.New("ID parameter wasn't provided"))
	}

	id, err := uuid.Parse(rawID)
	if err != nil {
		return HttpErrMalformedID.WithInternal(err)
	}

	archive, err := h.task.PopArchive(c.Request().Context(), id)
	if err != nil {
		return mapError(err)
	}

	return c.Blob(http.StatusOK, "application/zip", archive)
}
