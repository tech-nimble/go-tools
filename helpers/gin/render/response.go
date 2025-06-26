package render

import (
	"net/http"

	"github.com/gin-gonic/gin/render"
)

type Response struct {
	Status    string `json:"status"`
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
	Data      any    `json:"data,omitempty"`
}

func (r Response) Render(w http.ResponseWriter) (err error) {
	return render.JSON{Data: r}.Render(w)
}

func (r Response) WriteContentType(w http.ResponseWriter) {
	render.JSON{}.WriteContentType(w)
}

func NewResponse(data any, errorCode int, status, message string) Response {
	return Response{
		Status:    status,
		ErrorCode: errorCode,
		Message:   message,
		Data:      data,
	}
}

func NewSuccessResponse(data any, message string) Response {
	return NewResponse(data, 0, "success", message)
}
