package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/go-zoo/bone"

	"svc-item/routes"
)

func main() {
	var serverPort string
	flag.StringVar(&serverPort, "port", ":10001", "HTTP port")
	flag.Parse()
	if !strings.HasPrefix(serverPort, ":") {
		serverPort = ":" + serverPort
	}

	mux := bone.New()
	routes.Build(mux)

	log.Printf("svc-item started sucessfully on port %s ....", serverPort)
	log.Fatal(http.ListenAndServe(serverPort, mux))
}
