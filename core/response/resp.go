package response

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func NewResponse(code int, msg string, data interface{}) *Response {
	if code == 0 {
		code = 200
	}
	if msg == "" {
		msg = Success
	}
	if data == nil {
		data = map[string]interface{}{}
	}
	return &Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}
