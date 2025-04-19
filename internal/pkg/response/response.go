package response

type AppResponse struct {
	Data    any    `json:"data,omitempty"`
	Msg     string `json:"msg,omitempty"`
	Success bool   `json:"success,omitempty"`
}

func NewAppResponse(msg string, data any) *AppResponse {
	return &AppResponse{
		Msg:     msg,
		Data:    data,
		Success: true,
	}
}
