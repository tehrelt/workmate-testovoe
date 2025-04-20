package handlers

import (
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/lib/tracer"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/models"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/services/taskservice"
	"go.opentelemetry.io/otel"
)

type ListTasksRequest struct {
	Status string `query:"status"`
	Page   uint64 `query:"page"`
	Limit  uint64 `query:"limit"`
}

type task struct {
	Id        string     `json:"id"`
	Title     string     `json:"title"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

type ListTasksResponse struct {
	Tasks       []task `json:"tasks"`
	Total       uint64 `json:"total"`
	HasNextPage bool   `json:"hasNextPage"`
}

func ListTasks(ts *taskservice.TaskService) echo.HandlerFunc {
	return func(c echo.Context) error {

		var req ListTasksRequest
		if err := c.Bind(&req); err != nil {
			return err
		}
		if err := c.Validate(&req); err != nil {
			return err
		}

		filter := &models.TaskFilter{}

		if req.Status != "" {
			filter = filter.SetStatus(models.TaskStatus(req.Status))
		}

		if req.Page != 0 {
			filter = filter.SetFrom((req.Page - 1) * req.Limit)
		}

		if req.Limit != 0 {
			filter = filter.SetLimit(req.Limit)
		}

		slog.Debug("query params", slog.Any("filter", filter))

		taskch, total, err := ts.Tasks(c.Request().Context(), filter)
		if err != nil {
			return err
		}

		_, span := otel.Tracer(tracer.TracerKey).Start(c.Request().Context(), "building a response")
		defer span.End()
		var tasks []task
		for t := range taskch {
			tasks = append(tasks, task{
				Id:        t.Id.String(),
				Title:     t.Title,
				Status:    t.Status.String(),
				CreatedAt: t.CreatedAt,
				UpdatedAt: t.UpdatedAt,
			})
		}

		resp := ListTasksResponse{
			Tasks:       tasks,
			Total:       total,
			HasNextPage: total > (req.Page)*req.Limit,
		}

		return c.JSON(200, resp)
	}
}
