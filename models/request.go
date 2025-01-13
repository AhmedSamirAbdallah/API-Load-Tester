package models

type ApiRequest struct {
	URL     string
	Method  string
	Headers map[string]string
	Body    string
}
