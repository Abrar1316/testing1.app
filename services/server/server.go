package main

import (
	"log"
	"net/http"

	"github.com/TFTPL/AWS-Cost-Calculator/services/webbff"
	"github.com/gorilla/mux"
)

func main() {

	// create router
	router := mux.NewRouter().StrictSlash(true)

	// add handlers to router
	webbff.InitHandlers(router)

	// start server
	log.Println("Starting the HTTP server on port 8070")
	log.Fatal(http.ListenAndServe(":8070", router))
}
