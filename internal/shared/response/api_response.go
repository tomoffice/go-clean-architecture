package response

type APIResponse[T any] struct {
	Data  T             `json:"data,omitempty"`
	Error *ErrorPayload `json:"errordefs,omitempty"`
	Meta  *MetaPayload  `json:"meta,omitempty"`
	BaseResponse
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
