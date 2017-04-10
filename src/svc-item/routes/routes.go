package routes

import (
    "net/http"
    "log"
    "os"
    "github.com/go-zoo/bone"
    "encoding/json"

    "svc-item/controllers/item"
    "svc-item/controllers/health"
)

var endpoints = map[string]string {
    "health": "/health",
    "itemLookup": "/item/lookup/:id",
}

func Build(mux *bone.Mux) {
    mux.NotFoundFunc(handleNotFound)

    mux.Get(endpoints["health"], logger(health.Check))

    itemController := item.New()
    mux.Get(endpoints["itemLookup"], logger(itemController.Lookup))
    
    for _,value := range endpoints {
        log.Printf("Registering endpoint %s", value)
    }
}

type ErrorMessage struct {
    Type    string
    Message string    
}

func handleNotFound(rw http.ResponseWriter, req *http.Request) {
    message := &ErrorMessage{"NotFound", "Route Not Found"}
    encoder := json.NewEncoder(rw)
    
    rw.Header().Set("Content-Type", "application/json")

    err := encoder.Encode(message)
    if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
    }
}

func logger(next http.HandlerFunc) http.Handler {
    logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)

    return http.HandlerFunc(
        func (rw http.ResponseWriter, req *http.Request) {
            logger.Printf("%s %s %s", req.RemoteAddr, req.Method, req.RequestURI)
            next.ServeHTTP(rw, req)
        })
}
