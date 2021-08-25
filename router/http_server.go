package router

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/rlaskowski/go-email/config"
	"github.com/rlaskowski/go-email/controller"
	"github.com/rlaskowski/go-email/registry"
)

type HttpServer struct {
	server        *http.Server
	router        *Router //*pat.PatternServeMux
	context       context.Context
	cancel        context.CancelFunc
	registry      registry.Registry
	serviceConfig config.ServiceConfig
	multipartPool sync.Pool
	handlePool    sync.Pool
}

func NewHttpServer(registry registry.Registry) *HttpServer {
	ctx, cancel := context.WithCancel(context.Background())

	sc := registry.ServiceConfig()

	h := &HttpServer{
		context:       ctx,
		cancel:        cancel,
		router:        NewRouter(), //pat.New(),
		registry:      registry,
		serviceConfig: sc,
	}

	h.server = &http.Server{
		Addr:           fmt.Sprintf(":%d", h.serviceConfig.HttpServerPort),
		ReadTimeout:    h.serviceConfig.HttpServerReadTimeout,
		WriteTimeout:   h.serviceConfig.HttpServerWriteTimeout,
		MaxHeaderBytes: h.serviceConfig.HttpMaxHeaderSize,
		Handler:        h,
	}

	h.multipartPool.New = func() interface{} {
		return new(controller.MutlipartController)
	}

	h.handlePool.New = func() interface{} {
		return NewHandle(nil, nil)
	}

	return h
}

func (h *HttpServer) Start() error {
	go func() {
		h.configureEndpoints()

		log.Printf("Starting REST API on %d port", h.serviceConfig.HttpServerPort)

		if err := h.server.ListenAndServe(); err != nil {
			log.Fatalf("Caught error while starting server: %s", err.Error())
		}
	}()

	return nil
}

func (h *HttpServer) Stop() error {
	log.Print("Stopping REST API")

	if err := h.server.Close(); err != nil {
		return err
	}

	h.cancel()

	return nil
}

func (h *HttpServer) configureEndpoints() {
	//h.Post("/file/send", h.SendWithFile)
	//h.Post("/send", h.Send)
	h.Get("/receive/list", h.ReceiveList)
}

func (h *HttpServer) add(method, path string, handler HandlerFunc) {
	h.router.Add(method, path, handler)
	//handler := http.HandlerFunc(nil)
	//h.router.Add(http.MethodGet, path, handler)
}

func (h *HttpServer) Get(path string, handlerFunc HandlerFunc) {
	h.add(http.MethodGet, path, handlerFunc)
}

func (h *HttpServer) Post(path string, handlerFunc HandlerFunc) {
	h.add(http.MethodPost, path, handlerFunc)
}

func (h *HttpServer) Put(path string, handlerFunc HandlerFunc) {
	h.add(http.MethodPut, path, handlerFunc)
}

func (h *HttpServer) Delete(path string, handlerFunc HandlerFunc) {
	h.add(http.MethodDelete, path, handlerFunc)
}

func (h *HttpServer) Options(path string, handlerFunc HandlerFunc) {
	h.add(http.MethodOptions, path, handlerFunc)
}

func (h *HttpServer) Head(path string, handlerFunc HandlerFunc) {
	h.add(http.MethodHead, path, handlerFunc)
}

func (h *HttpServer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	handle := h.handlePool.Get().(*Handle)
	defer h.handlePool.Put(handle)

	handle.Reload(rw, r)

	handlerFunc := h.router.FindHandle(r.Method, r.URL.Path)
	handlerFunc(handle)
}

func (h *HttpServer) ReceiveList(handler Handler) {
	/* que, err := h.registries.QueueFactory.GetOrCreate(queue.EmailQueueType)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(err.Error()))
		return
	}

	list, _ := que.Subscribe(queue.SubjectReceiving)
	h.json(rw, list) */
	a := struct {
		Name     string
		Lastname string
	}{
		"Rafal",
		"Laskowski",
	}

	handler.JSON(http.StatusOK, a)
}

func (h *HttpServer) SendWithFile(rw http.ResponseWriter, r *http.Request) {
	/*var result error

	multipartReader, err := r.MultipartReader()
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(err.Error()))
		return
	}
	//c := r.Context()

	multipartController := h.acquireMultipart(multipartReader)

	message, err := multipartController.Message()
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(err.Error()))
		return
	}

	file, err := multipartController.File()
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(err.Error()))
	}

	que, err := h.registries.QueueFactory.GetOrCreate(queue.EmailQueueType)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(err.Error()))
		return
	}

	 if err := que.Publish(queue.SubjectSending, message, file); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(err.Error()))
		return
	}

	if result != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(err.Error()))
	} else {
		h.json(rw, map[string]string{
			"result": "Message published successfully",
		})
		rw.WriteHeader(http.StatusOK)
	}*/

}

func (h *HttpServer) Send(rw http.ResponseWriter, r *http.Request) {
	/*var result error

	messageForm := r.FormValue("message")

	message := new(model.Message)

	if err := json.Unmarshal([]byte(messageForm), message); err != nil {
		result = h.multiFailure(err.Error())
	}

	 que, err := h.registries.QueueFactory.GetOrCreate(queue.EmailQueueType)
	if err != nil {
		result = h.multiFailure(err.Error())
	}

	if err := que.Publish(queue.SubjectSending, message); err != nil {
		result = h.multiFailure(err.Error())
	}

	if result != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(err.Error()))
	} else {
		h.json(rw, map[string]string{
			"result": "Message published successfully",
		})
		rw.WriteHeader(http.StatusOK)
	}*/

}

/* func (h *HttpServer) storeFile(multipartController *controller.MutlipartController) (string, error) {
	file, err := multipartController.File()
	if err != nil {
		return "", err
	}

	fileStore := store.NewFileStore(config.FileStorePath)

	fileuuid, err := fileStore.Store(file)
	if err != nil {
		return "", err
	}

	return fileuuid, nil
} */

/* func (h *HttpServer) multiFailure(err string) error {
	return fmt.Errorf("%s\r\n", err)
}

func (h *HttpServer) json(rw http.ResponseWriter, i interface{}) {
	rw.Header().Add("Content-Type", "application/json")

	marshal, err := json.Marshal(i)
	if err != nil {
		rw.Write([]byte(err.Error()))
	}
	rw.Write(marshal)
}

func (h *HttpServer) acquireMultipart(reader *multipart.Reader) *controller.MutlipartController {
	m := h.multipartPool.Get().(*controller.MutlipartController)
	defer h.multipartPool.Put(m)

	m.Reader = reader

	return m
} */

/* func (h *HttpServer) BME280(rw http.ResponseWriter, r *http.Request) {
	driver, err := h.registries.RaspiDriver.BME280Driver()
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
	} else {
		h.json(rw, driver.Stat())
	}

}

func (h *HttpServer) Drivers(rw http.ResponseWriter, r *http.Request) {
	h.json(rw, h.registries.DriverRepository.FindAll())
}

func (h *HttpServer) DriversByGroup(rw http.ResponseWriter, r *http.Request) {
	group := r.URL.Query().Get(":group")
	h.json(rw, h.registries.DriverRepository.FindByGroup(group))
}


*/
