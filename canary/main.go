package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"fmt"
)

type HealthStatus struct {
	Url    string `json:"url,omitempty"`
	Status string `json:"status,omitempty"`
}

func main() {
	// need a repo for all the endpoints
	repo, err := NewEndpointRepo()
	if err != nil {
		log.Fatalln(err)
		return
	}
	// infinite
	for {
		// get all endpoints from the repo
		endpoints, err := repo.GetEndpoints()
		if err != nil {
			log.Fatalln(err)
			return
		}
		// check the endpoints
		CheckEndpoints(endpoints)
		// wait
		time.Sleep(30 * time.Second)
	}
}

func CheckEndpoints(endpoints []string) {
	for _, url := range endpoints {
		MakeRequest(url)
	}
}

func MakeRequest(url string) {
	// see if the endpoint is up
	resp, err := http.Get(url)
	if err != nil {
		// report error
		log.Error(HealthStatus{url, "unhealthy"})
		log.Error(err)
		return
	}
	if resp.StatusCode == http.StatusOK {
		// see if the endpoint returns a body
		_, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			// report error
			log.Error(HealthStatus{url, "unhealthy"})
			log.Error(err)
			return
		}
		// log success
		log.Info(HealthStatus{url, "healthy"})
	} else {
		log.Error(HealthStatus{url, "unhealthy"})
		log.Error(fmt.Sprintf("Response code: %d", resp.StatusCode))
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
