package errorz

import "fmt"

type AppError struct {
	Code        string
	Message     string
	HTTPStatus  int
	Description string
	Module      string
	Retryable   bool
	Severity    string
	ErrorType   string
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}
