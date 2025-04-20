package middlewares

import (
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/tehrelt/workmate-testovoe/pkg/sl"
)

func Logging(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		start := time.Now()
		defer func() {
			duration := time.Since(start)
			if err != nil {
				slog.Error("request failed", sl.Err(err))
				return
			}
			slog.Info(
				"request",
				slog.String("method", c.Request().Method),
				slog.String("path", c.Path()),
				slog.Int("status", c.Response().Status),
				slog.Int("size", int(c.Response().Size)),
				slog.Duration("latency", duration),
			)
		}()

		err = next(c)
		return err
	}
}
