package handlers

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/services/taskservice"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/storage"
)

type GetTaskRequest struct {
	Id string `param:"id" validate:"required"`
}

type GetTaskResponse struct {
	Id        string     `json:"id"`
	Title     string     `json:"title"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

func GetTask(ts *taskservice.TaskService) echo.HandlerFunc {
	return func(c echo.Context) error {

		var req GetTaskRequest
		if err := c.Bind(&req); err != nil {
			return err
		}
		if err := c.Validate(&req); err != nil {
			return err
		}

		id, err := uuid.Parse(req.Id)
		if err != nil {
			return echo.NewHTTPError(400, "bad id")
		}

		task, err := ts.Task(c.Request().Context(), id)
		if err != nil {
			if errors.Is(err, storage.ErrTaskNotFound) {
				return echo.NewHTTPError(404, "task not found")
			}
			return err
		}
		resp := GetTaskResponse{
			Id:        task.Id.String(),
			Title:     task.Title,
			Status:    task.Status.String(),
			CreatedAt: task.CreatedAt,
			UpdatedAt: task.UpdatedAt,
		}

		return c.JSON(200, resp)
	}
}
