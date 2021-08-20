package router

import (
	"net/http"
)

type HandlerFunc func(Handler)

type Handler interface {
	FormValue(key string) string

	Reload()

	Handler() http.HandlerFunc
}

type Handle struct {
	request  *http.Request
	response *Response
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

func (h *Handle) Handler() http.HandlerFunc {
	handler := func(rw http.ResponseWriter, r *http.Request) {
		rw = h.response.Writer
		r = h.request
	}

	return handler
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
