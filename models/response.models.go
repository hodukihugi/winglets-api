package models

// HTTPResponse represents a boilerplate of HTTP response payload
type HTTPResponse struct {
	Data          interface{}     `json:"data,omitempty"`
	Message       string          `json:"message,omitempty"`
	InvalidFields []string        `json:"invalid_fields,omitempty"`
	Pagination    *PaginationResp `json:"pagination,omitempty"`
}
