package response

import (
	"module-clean/internal/shared/enum"
	"time"
)

type BaseResponse struct {
	UnixTimeStamp    int64          `json:"unix_timestamp"`
	RFC3339TimeStamp string         `json:"rfc3339_timestamp"`
	Status           enum.APIStatus `json:"status"`
}

func NewBaseResponse(status enum.APIStatus) BaseResponse {
	now := time.Now()
	return BaseResponse{
		UnixTimeStamp:    now.Unix(),
		RFC3339TimeStamp: now.Format(time.RFC3339),
		Status:           status,
	}
}
