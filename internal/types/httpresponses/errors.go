package httpresponses

type ErrorCode struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

type ErrorResponse struct {
	ErrorCode ErrorCode `json:"error_code"`
}
