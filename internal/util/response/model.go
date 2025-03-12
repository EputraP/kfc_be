package response

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func NewResponse(code int, message string, data interface{}) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
}
