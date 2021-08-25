package router

import (
	"encoding/json"
	"net/http"
)

const (
	HeaderContentType   = "Content-Type"
	MIMEApplicationJson = "application/json"
)

type HandlerFunc func(Handler)

type Handler interface {
	FormValue(key string) string

	JSON(code int, i interface{})

	Reload(rw http.ResponseWriter, r *http.Request)

	Request() *http.Request

	Handler() HandlerFunc
}

type Handle struct {
	request  *http.Request
	response *Response
	handler  HandlerFunc
}

func NewHandle(rw http.ResponseWriter, r *http.Request) *Handle {
	return &Handle{
		request:  r,
		response: NewResponse(rw),
	}
}

func (h *Handle) FormValue(key string) string {
	return h.request.FormValue(key)
}

func (h *Handle) Handler() HandlerFunc {
	return h.handler
}

func (h *Handle) writeContentType(content string) {
	header := h.response.Header()
	header.Set("Content-Type", content)
}

func (h *Handle) json(code int, i interface{}) {
	h.writeContentType(MIMEApplicationJson)
	h.response.Status = code
	enc := json.NewEncoder(h.response.Writer)
	enc.Encode(i)
}

func (h *Handle) JSON(code int, i interface{}) {
	h.json(code, i)
}

func (h *Handle) Request() *http.Request {
	return h.request
}

func (h *Handle) SetRequest(req *http.Request) {
	h.request = req
}

func (h *Handle) Response() *Response {
	return h.response
}

func (h *Handle) SetResponse(res *Response) {
	h.response = res
}

func (h *Handle) Reload(rw http.ResponseWriter, r *http.Request) {
	h.response.Writer = rw
	h.request = r

}
