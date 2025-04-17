package response

import "module-clean/internal/shared/enum"

type BaseResponse struct {
	UnixTimeStamp    int64          `json:"unix_timestamp"`
	RFC3339TimeStamp string         `json:"rfc3339_timestamp"`
	Status           enum.APIStatus `json:"status"`
}
