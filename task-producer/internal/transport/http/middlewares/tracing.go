package middlewares

import (
	"fmt"
	"log/slog"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func Tracing(serviceName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {

			req := c.Request()
			tracer := otel.Tracer(serviceName)
			ctx, span := tracer.Start(req.Context(), fmt.Sprintf("%s %s", c.Request().Method, c.Path()))
			defer span.End()
			slog.Debug("tracing request",
				slog.String("span_id", span.SpanContext().SpanID().String()),
				slog.String("trace_id", span.SpanContext().TraceID().String()),
			)

			span.SetAttributes(
				attribute.String("http.method", c.Request().Method),
				attribute.String("http.host", c.Request().Host),
				attribute.String("http.path", c.Path()),
				attribute.String("htpp.real_ip", c.RealIP()),
			)

			c.SetRequest(req.WithContext(ctx))

			defer func() {
				if err != nil {
					span.RecordError(err)
					span.SetStatus(codes.Error, err.Error())
				}
			}()

			err = next(c)
			return
		}
	}
}
