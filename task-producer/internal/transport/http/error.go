package http

type ValidationError struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details"`
}

func newValidationError(code int, message string, details map[string]string) *ValidationError {
	return &ValidationError{
		Code:    code,
		Message: message,
		Details: details,
	}
}
