package task

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrTaskNotFound     = errors.New("no task with such ID exists")
	ErrTooManyTasks     = errors.New("tasks max amount was reached")
	ErrTooManyTaskFiles = errors.New("task's max file amount was reached")
)

type TaskService struct {
	tasks        map[uuid.UUID]*Task
	maxTaskFiles int
	maxTasks     int
}

func NewTaskService(maxTaskFiles int, maxTasks int) *TaskService {
	return &TaskService{
		tasks:        make(map[uuid.UUID]*Task, maxTasks),
		maxTaskFiles: maxTaskFiles,
		maxTasks:     maxTasks,
	}
}

func (s *TaskService) Create(ctx context.Context) (uuid.UUID, error) {
	if len(s.tasks) >= s.maxTasks {
		return uuid.UUID{}, ErrTooManyTasks
	}
	id := uuid.New()
	s.tasks[id] = NewTask()
	return id, nil
}

func (s *TaskService) AppendFile(ctx context.Context, taskID uuid.UUID, source string) error {
	task, ok := s.tasks[taskID]
	if !ok {
		return ErrTaskNotFound
	}
	if task.FilesAmount() >= s.maxTaskFiles {
		return ErrTooManyTaskFiles
	}
	task.AddDownload(source)
	return nil
}

func (s *TaskService) Check(ctx context.Context, taskID uuid.UUID) (TaskInfo, error) {
	task, ok := s.tasks[taskID]
	if !ok {
		return TaskInfo{}, ErrTaskNotFound
	}
	status := task.Status()
	if len(status.Files) < s.maxTaskFiles {
		status.Ready = false
	}
	return status, nil
}

func (s *TaskService) PopArchive(ctx context.Context, taskID uuid.UUID) ([]byte, error) {
	task, ok := s.tasks[taskID]
	if !ok {
		return nil, ErrTaskNotFound
	}
	delete(s.tasks, taskID)

	archive, err := task.Finish(ctx)
	if err != nil {
		return nil, err
	}
	return archive, nil
}
