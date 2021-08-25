package router

type Router struct {
	handlers map[string]HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		handlers: make(map[string]HandlerFunc),
	}
}

func (r *Router) Add(method, path string, handler HandlerFunc) {
	if _, ok := r.handlers[method+path]; !ok {
		r.handlers[method+path] = handler
	}
}

func (r *Router) FindHandle(method, path string) HandlerFunc {
	if hf, ok := r.handlers[method+path]; ok {
		return hf
	}
	return nil
}
