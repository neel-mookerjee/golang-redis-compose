package main

import (
	log "github.com/sirupsen/logrus"
	"os"

)

// @title Sample Golang RESTFul APIs
// @version 0.1
// @description Provide a service endpoint that will generate a unique URI

//

// @contact.name Arghanil
// @contact.email ARghanil@gmail.com

// @BasePath
// @Host http://localhost:8080/

func main() {
	// handler
	h, err := NewHandlerWrapper()
	if err != nil {
		log.Fatal(err)
	}

	// router
	httpRouter := NewHttpRouter(h)
	httpRouter.Register()
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
