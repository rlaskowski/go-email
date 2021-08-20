package router

import "net/http"

type Response struct {
	Writer http.ResponseWriter
}

func NewResponse(w http.ResponseWriter) *Response {
	return &Response{Writer: w}
}

func (r *Response) Header() http.Header {
	return r.Writer.Header()
}
