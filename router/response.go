package router

import "net/http"

type Response struct {
	Status int
	Writer http.ResponseWriter
}

func NewResponse(w http.ResponseWriter) *Response {
	return &Response{Writer: w}
}

func (r *Response) Header() http.Header {
	return r.Writer.Header()
}

func (r *Response) WriteHeader(code int) {
	r.Writer.WriteHeader(code)
}

func (r *Response) Write(b []byte) {
	if r.Status == 0 {
		r.Status = http.StatusOK
	}
	r.WriteHeader(r.Status)

	r.Writer.Write(b)
}
