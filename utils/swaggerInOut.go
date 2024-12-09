package utils

type LoginResponse struct {
	ID string `json:"id"`
	Token string `json:"token"`
}

type ResponseOK struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
type ErrorResponse struct {
	ErrorMsg string `json:"error_msg,omitempty"`
	Message  string `json:"message"`
}
