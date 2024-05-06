package models

// HTTPResponse represents a boilerplate of HTTP response payload
type HTTPResponse struct {
	Message       string          `json:"message,omitempty"`
	Data          interface{}     `json:"data,omitempty"`
	InvalidFields []string        `json:"invalid_fields,omitempty"`
	Pagination    *PaginationResp `json:"pagination,omitempty"`
}
