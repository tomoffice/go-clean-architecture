package response

type APIResponse[T any] struct {
	Data  T             `json:"data,omitempty"`
	Error *ErrorPayload `json:"error,omitempty"`
	Meta  *MetaPayload  `json:"meta,omitempty"`
	BaseResponse
}

type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type MetaPayload struct {
	Total  int `json:"total,omitempty"`
	Page   int `json:"page,omitempty"`
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}
