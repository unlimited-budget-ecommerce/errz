package errorz

type ErrorDefinition struct {
	Domain      string   `json:"domain"`
	Code        string   `json:"code"`
	Msg         string   `json:"msg"`
	Cause       string   `json:"cause"`
	HTTPStatus  int      `json:"http_status"`
	Category    string   `json:"category"`
	Severity    string   `json:"severity"`
	Solution    string   `json:"solution,omitempty"`
	IsRetryable bool     `json:"is_retryable"`
	Tags        []string `json:"tags,omitempty"`
}
