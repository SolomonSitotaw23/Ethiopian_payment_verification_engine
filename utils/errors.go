package utils

type AppError struct {
	Status  int    `json:"status"`
	Message string `json:"error"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(message string, status int) *AppError {
	return &AppError{
		Status:  status,
		Message: message,
	}
}

func NewValidationError(message string) *AppError {
	return NewAppError(message, 400)
}

func NewNotFoundError(message string) *AppError {
	return NewAppError(message, 404)
}

func NewConnectionTimeoutError(message string) *AppError {
	if message == "" {
		message = "Request timed out"
	}
	return NewAppError(message, 504)
}

func NewUpstreamServiceError(message string) *AppError {
	if message == "" {
		message = "Upstream service unavailable"
	}
	return NewAppError(message, 502)
}

func GetHTTPStatus(err error) int {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Status
	}
	return 500
}
