package router

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"sync"

	"github.com/bmizerany/pat"
	"github.com/rlaskowski/go-email/config"
	"github.com/rlaskowski/go-email/controller"
	"github.com/rlaskowski/go-email/model"
	"github.com/rlaskowski/go-email/queue"
	"github.com/rlaskowski/go-email/registries"
)

type HttpServer struct {
	server        *http.Server
	router        *pat.PatternServeMux
	context       context.Context
	cancel        context.CancelFunc
	registries    *registries.Registries
	multipartPool sync.Pool
}

type Router struct {
	method string
	host   string
	name   http.HandlerFunc
}

func NewHttpServer(registries *registries.Registries) *HttpServer {
	ctx, cancel := context.WithCancel(context.Background())

	h := &HttpServer{
		context:    ctx,
		cancel:     cancel,
		router:     pat.New(),
		registries: registries,
	}

	h.server = &http.Server{
		Addr:           fmt.Sprintf(":%d", config.HttpServerPort),
		ReadTimeout:    config.HttpServerReadTimeout,
		WriteTimeout:   config.HttpServerWriteTimeout,
		MaxHeaderBytes: config.HttpMaxHeaderSize,
		Handler:        h,
	}

	h.multipartPool.New = func() interface{} {
		return new(controller.MutlipartController)
	}

	return h
}

func (h *HttpServer) Start() error {
	go func() {
		h.configureEndpoints()

		log.Printf("Starting REST API on http://localhost:%d", config.HttpServerPort)

		if err := h.server.ListenAndServe(); err != nil {
			log.Fatalf("Caught error while starting server: %s", err.Error())
		}
	}()

	return nil
}

func (h *HttpServer) Stop() error {
	h.cancel()

	log.Print("Stopping REST API")

	return h.server.Close()
}

func (h *HttpServer) configureEndpoints() {
	h.Post("/send/file", h.SendWithFile)
	h.Post("/send", h.Send)
}

func (h *HttpServer) Get(path string, handler http.HandlerFunc) {
	h.router.Get(path, handler)
}

func (h *HttpServer) Post(path string, handler http.HandlerFunc) {
	h.router.Post(path, handler)
}

func (h *HttpServer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(rw, r)
}

func (h *HttpServer) SendWithFile(rw http.ResponseWriter, r *http.Request) {
	multipartReader, err := r.MultipartReader()
	//c := r.Context()

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
	}

	multipartController := h.acquireMultipart(multipartReader)

	file, err := multipartController.File()
	if err != nil {
		h.json(rw, map[string]string{
			"error_message": err.Error(),
		})
	}

	h.json(rw, map[string]string{
		"file_name": file.FileName(),
		"form_name": file.FormName(),
	})

}

func (h *HttpServer) Send(rw http.ResponseWriter, r *http.Request) {
	messageForm := r.FormValue("message")

	message := new(model.Message)

	if err := json.Unmarshal([]byte(messageForm), message); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(err.Error()))
	}

	queue, err := h.registries.QueueFactory.GetOrCreate(queue.EmailQueueType)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(err.Error()))
	}

	queue.Publish(message)

	rw.Write([]byte(messageForm))
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

	m.MutlipartReader = reader

	return m
}

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
