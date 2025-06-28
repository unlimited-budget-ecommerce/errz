//go:generate go run ./cmd/gen_errors/gen.go
package errz

type ErrorDefinition struct {
	Domain string `json:"domain"`
	Code   string `json:"code"`
	Msg    string `json:"msg"`
	Cause  string `json:"cause"`
}
