package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/models"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/services/taskservice"
)

type CreateTaskRequest struct {
	Title string `json:"title" validate:"required"`
}

type CreateTaskResponse struct {
	Id string `json:"id"`
}

func CreateTask(ts *taskservice.TaskService) echo.HandlerFunc {
	return func(c echo.Context) error {

		var req CreateTaskRequest
		if err := c.Bind(&req); err != nil {
			return err
		}
		if err := c.Validate(&req); err != nil {
			return err
		}

		task, err := ts.CreateTask(c.Request().Context(), &models.CreateTask{
			Title: req.Title,
		})
		if err != nil {
			return err
		}
		resp := CreateTaskResponse{
			Id: task.Id.String(),
		}

		return c.JSON(200, resp)
	}
}
