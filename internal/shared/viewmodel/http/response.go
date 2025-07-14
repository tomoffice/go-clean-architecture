package http

import (
	"github.com/tomoffice/go-clean-architecture/internal/shared/enum"
	"time"
)

type BaseHTTPResponse struct {
	UnixTimeStamp    int64          `json:"unix_timestamp"`
	RFC3339TimeStamp string         `json:"rfc3339_timestamp"`
	Status           enum.APIStatus `json:"status"`
}

func NewBaseHTTPResponse(status enum.APIStatus) BaseHTTPResponse {
	now := time.Now()
	return BaseHTTPResponse{
		UnixTimeStamp:    now.Unix(),
		RFC3339TimeStamp: now.Format(time.RFC3339),
		Status:           status,
	}
}

type HTTPResponse[T any] struct {
	Data  T             `json:"data,omitempty"`
	Error *ErrorPayload `json:"errordefs,omitempty"`
	Meta  *MetaPayload  `json:"meta,omitempty"`
	BaseHTTPResponse
}
type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type MetaPayload struct {
	Total  int `json:"total"`
	Page   int `json:"page"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
