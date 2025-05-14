package errorz

type ErrorDefinition struct {
	Code        string `json:"code"`
	HTTPStatus  int    `json:"http_status"`
	Message     string `json:"message"`
	Description string `json:"description"`
	Module      string `json:"module"`
	Retryable   bool   `json:"retryable"`  // Whether the operation can be retried
	Severity    string `json:"severity"`   // Severity level, e.g., "info", "warn", "error"
	ErrorType   string `json:"error_type"` // Type/category of error, e.g., "validation", "internal", "db"
}
