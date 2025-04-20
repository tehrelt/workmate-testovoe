package http

import (
	"fmt"
	"log/slog"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type httpValidator struct {
	validator  *validator.Validate
	translator ut.Translator
}

func newHttpValidator(v *validator.Validate, t ut.Translator) *httpValidator {
	return &httpValidator{
		validator:  v,
		translator: t,
	}
}

func (hv *httpValidator) Validate(v any) error {
	if err := hv.validator.Struct(v); err != nil {
		if e, ok := err.(validator.ValidationErrors); ok {
			slog.Debug("validate error", slog.Any("type", fmt.Sprintf("%t", err)))
			translated := e.Translate(hv.translator)
			for _, err := range translated {
				slog.Debug("translated error", slog.Any("error", err))
			}

			validationError := echo.NewHTTPError(400, newValidationError(400, "Validation failed", translated))
			return validationError
		}
		return err
	}
	return nil
}
