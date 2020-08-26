package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	"bytes"
	"io/ioutil"

	"github.com/go-chi/chi"
	"github.com/go-redis/redis"
	"github.com/gorilla/handlers"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/speps/go-hashids"
)

type HandlerWrapper struct {
	db DbInterface
}

// http handler wrapper
func NewHandlerWrapper() (*HandlerWrapper, error) {
	db, err := NewRedisDb()
	if err != nil {
		return nil, err
	}
	return &HandlerWrapper{db: db}, nil
}

// expect a json object as param to return as response
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	msg := payload
	response, err := json.Marshal(msg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

// expect a json string as param to return as response
func respondJSONStr(w http.ResponseWriter, status int, payload string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(payload))
}

// error response
func respondError(w http.ResponseWriter, code int, err error) {
	respondJSON(w, code, map[string]string{"error": err.Error()})
}

func (a *HandlerWrapper) WrapHandler(h http.Handler) http.Handler {
	authFunc := func(w http.ResponseWriter, r *http.Request) {
		dumpRequest(r)
		h.ServeHTTP(w, r)
	}
	return allowCORS(http.HandlerFunc(authFunc))

}

func dumpRequest(r *http.Request) {
	httputil.DumpRequest(r, true)
}

// cors
func allowCORS(h http.Handler) http.Handler {
	options := []handlers.CORSOption{
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"POST", "GET", "PUT"}),
		handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "ClientId", "Env", "Access-Control-Allow-Origin"}),
	}

	return handlers.CORS(options...)(h)
}

func (h *HandlerWrapper) Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the endpoint generation service!")
}

// @Summary Create a Unique Endpoint
// @Description Create an Endpoint which is unique and accepts GET/POST
// @ID GenerateUniqueEndpoint
// @Produce json
// @Tags generate
// @Success 201 {object} main.EndpointOutput "Endpoint created: url is returned"
// @Failure 500 "There was as error during processing your request. "
// @Router /generate [put]
func (handler *HandlerWrapper) GenerateUniqueEndpoint(w http.ResponseWriter, r *http.Request) {
	var url Endpoint
	// always generate a new endpoint
	id, err := handler.GenerateUniqueId()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	url.ID = id
	url.Url = "http://" + getEnv("APP_URL", "app") + ":" + getEnv("APP_PORT", "8080") + "/" + url.ID
	err = handler.db.Save("url:"+url.ID, url.Url)
	if err != nil {
		err := errors.New("There was as error during processing your request. " + err.Error())
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	m := EndpointOutput{url.Url, "To use from outside the docker environment use: http://localhost:8080/" + url.ID}
	respondJSON(w, http.StatusCreated, m)
}

// generate unique id: IS NANOSECOND sufficient to eliminate possibility of duplicates??
func (handler *HandlerWrapper) GenerateUniqueId() (string, error) {
	id := ""
	hd := hashids.NewData()
	h, err := hashids.NewWithData(hd)
	if err != nil {
		return id, err
	}
	now := time.Now()
	id, err = h.Encode([]int{int(now.UnixNano())})
	return id, err
}

// @Summary Unique Endpoint for POST
// @Description Unique Endpoint to submit a payload
// @ID SubmitContent
// @Accept json
// @Produce json
// @Tags submit
// @Success 202 "submitted"
// @Failure 404 "Endpoint not found"
// @Failure 400 "Bad request body"
// @Failure 500 "Internal server error"
// @Router / [post]
func (handler *HandlerWrapper) SubmitContent(w http.ResponseWriter, r *http.Request) {
	// validate if this is a generated endpoint
	urlId := chi.URLParam(r, "id")
	if len(urlId) <= 0 {
		err := errors.New("Endpoint not found")
		respondError(w, http.StatusNotFound, err)
		return
	}
	cmd, err := handler.db.Retrieve("url:" + urlId)
	log.Debug(cmd)
	if err == redis.Nil {
		err := errors.New("Endpoint not found")
		respondError(w, http.StatusNotFound, err)
		return
	}
	if err != nil && err != redis.Nil {
		err := errors.New("There was as error during processing your request. " + err.Error())
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	// store the payload
	value, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err := errors.New("Body content not accepted: " + err.Error())
		respondError(w, http.StatusBadRequest, err)
		return
	}
	err = handler.db.AddToList("list:"+urlId, bytes.NewBuffer(value).String())
	if err != nil {
		err := errors.New("There was as error during processing your request. " + err.Error())
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	respondJSON(w, http.StatusAccepted, map[string]string{"status": "submitted"})
}

// @Summary Unique Endpoint for GET
// @Description Unique Endpoint to get the last submitted payload
// @ID GetLastContent
// @Produce json
// @Tags retrieve
// @Success 302 "Found"
// @Success 404 "No payload was submitted"
// @Failure 404 "Endpoint not found"
// @Failure 500 "Internal server error"
// @Router / [get]
func (handler *HandlerWrapper) GetLastContent(w http.ResponseWriter, r *http.Request) {
	// validate endpoint
	urlId := chi.URLParam(r, "id")
	if len(urlId) <= 0 {
		err := errors.New("Endpoint not found")
		respondError(w, http.StatusNotFound, err)
		return
	}
	cmd, err := handler.db.Retrieve("url:" + urlId)
	log.Debug(cmd)
	if err == redis.Nil {
		err := errors.New("Endpoint not found")
		respondError(w, http.StatusNotFound, err)
		return
	}
	if err != nil && err != redis.Nil {
		err := errors.New("There was as error during processing your request. " + err.Error())
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	// get the last payload
	val, err := handler.db.ReadFromList("list:"+urlId, 0, 0)
	if err != nil {
		err := errors.New("There was as error during processing your request. " + err.Error())
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	if len(val) > 0 {
		respondJSONStr(w, http.StatusFound, val[0])
		return
	} else {
		err := errors.New("No payload was submitted to this endpoint")
		respondError(w, http.StatusNotFound, err)
		return
	}
}

// @Summary Unique Endpoint health
// @Description Unique Endpoint healthcheck
// @ID Health
// @Produce json
// @Tags health
// @Success 200 "Status:healthy"
// @Failure 404 "Status:unhealthy"
// @Failure 500 "Status:unhealthy"
// @Router /health [get]
func (handler *HandlerWrapper) Health(w http.ResponseWriter, r *http.Request) {
	urlId := chi.URLParam(r, "id")
	if len(urlId) <= 0 {
		err := errors.New("Endpoint not found")
		respondError(w, http.StatusNotFound, err)
		return
	}
	cmd, err := handler.db.Retrieve("url:" + urlId)
	log.Debug(cmd)
	if err == redis.Nil {
		err := errors.New("Endpoint not found")
		respondError(w, http.StatusNotFound, err)
		return
	}
	if err != nil && err != redis.Nil {
		err := errors.New("There was as error during processing your request. " + err.Error())
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}
