package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
	"github.com/swaggo/http-swagger"
	"io/ioutil"
	"encoding/json"
)

type HttpRouterInterface interface {
	Register()
}

type HttpRouter struct {
	handlerWrapper *HandlerWrapper
}

func NewHttpRouter(handlerWrapper *HandlerWrapper) HttpRouterInterface {
	return &HttpRouter{handlerWrapper: handlerWrapper}
}

func (r *HttpRouter) Register() {
	// create the router
	router := chi.NewRouter()

	router.Get("/swagger/*", httpSwagger.WrapHandler)

	router.Put("/generate", r.handlerWrapper.GenerateUniqueEndpoint)
	router.Post("/{id}", r.handlerWrapper.SubmitContent)
	router.Get("/{id}", r.handlerWrapper.GetLastContent)
	router.Get("/{id}/health", r.handlerWrapper.Health)
	//router.Get("/", r.handlerWrapper.Index)
	router.Get("/", GetDetail)

	httpServer := &http.Server{
		Handler:      router,
		Addr:         ":" + getEnv("APP_PORT", "8080"),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Info("Starting Server...")
	log.Fatal(httpServer.ListenAndServe())

}

func GetDetail(w http.ResponseWriter, r *http.Request) {

	b,_ := ioutil.ReadFile("/app/docs/swagger/swagger.json");

	rawIn := json.RawMessage(string(b))
	var objmap map[string]*json.RawMessage
	err := json.Unmarshal(rawIn, &objmap)
	if err != nil {
		log.Error(err)
	}
	log.Info(objmap)

	json.NewEncoder(w).Encode(objmap)
}
