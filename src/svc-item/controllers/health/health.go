package health

import (
    "encoding/json"
    "net/http"
)

type Status struct {
    Status string
}

func Check(rw http.ResponseWriter, req *http.Request) {
    message := &Status{"UP"}
    encoder := json.NewEncoder(rw)
   
    err := encoder.Encode(message)
    if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
    }
}
