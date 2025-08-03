package app

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lowerc4s3/29.07.2025-zipload/batch"
	"github.com/lowerc4s3/29.07.2025-zipload/models"
)

type BatchService interface {
	DownloadAll(ctx context.Context, sources []string) ([]byte, error)
}

type TaskService interface {
	Create() (uuid.UUID, error)
	AppendFile(taskID uuid.UUID) error
	Check(taskID uuid.UUID) (models.TaskStatus, error)
	PopArchive(taskID uuid.UUID) ([]byte, error)
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
		return ErrMalformedRequest.WithInternal(err)
	}

	archive, err := h.batch.DownloadAll(c.Request().Context(), request.Sources)
	status := http.StatusOK
	if err != nil {
		if errors.Is(err, batch.ErrForbiddenMIME) {
			return ErrForbiddenMime
		} else if !errors.Is(err, batch.ErrPartitialArchive) {
			return ErrDownload
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
	id, err := h.task.Create()
	if err != nil {
		return ErrTooManyTasks.WithInternal(err)
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"id": id.String(),
	})
}

func (h *Handler) AppendTask(c echo.Context) error {
	request := new(models.AppendTaskRequest)
	if err := c.Bind(request); err != nil {
		return ErrMalformedRequest.WithInternal(err)
	}

	id, err := uuid.Parse(request.ID)
	if err != nil {
		return ErrMalformedID.WithInternal(err)
	}

	if err := h.task.AppendFile(id); err != nil {
		return err // TODO: Check returned error type
	}

	return c.NoContent(http.StatusAccepted)
}

func (h *Handler) CheckTask(c echo.Context) error {
	rawID := c.Param("taskid")
	if rawID == "" {
		return ErrMalformedRequest.WithInternal(errors.New("ID parameter wasn't provided"))
	}

	id, err := uuid.Parse(rawID)
	if err != nil {
		return ErrMalformedID.WithInternal(err)
	}

	status, err := h.task.Check(id)
	if err != nil {
		return err // TODO: Check returned error type
	}

	return c.JSON(http.StatusOK, status)
}

func (h *Handler) DownloadTask(c echo.Context) error {
	rawID := c.Param("taskid")
	if rawID == "" {
		return ErrMalformedRequest.WithInternal(errors.New("ID parameter wasn't provided"))
	}

	id, err := uuid.Parse(rawID)
	if err != nil {
		return ErrMalformedID.WithInternal(err)
	}

	archive, err := h.task.PopArchive(id)
	if err != nil {
		return err // TODO:
	}

	return c.Blob(http.StatusOK, "application/zip", archive)
}
