package response

type AppResponse struct {
	Msg     string `json:"msg,omitempty"`
	Data    any    `json:"data,omitempty"`
	Success bool   `json:"success,omitempty"`
}

func NewAppResponse(msg string, data any) *AppResponse {
	return &AppResponse{
		Msg:     msg,
		Data:    data,
		Success: true,
	}
}
