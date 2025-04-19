package middlewares

import (
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
)

func Logging(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		defer func() {
			duration := time.Since(start)
			slog.Info(
				"request",
				slog.String("method", c.Request().Method),
				slog.String("path", c.Path()),
				slog.Int("status", c.Response().Status),
				slog.Int("size", int(c.Response().Size)),
				slog.Duration("latency", duration),
			)
		}()

		return next(c)
	}
}
