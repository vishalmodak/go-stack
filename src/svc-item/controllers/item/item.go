package item

import (
    "encoding/json"
    "net/http"
    "github.com/go-zoo/bone"
)

type Controller struct {

}

func New() *Controller {
    return &Controller{}    
}

type Item struct {
    ItemId      string
    Title       string
    Description string
    Category    string
    Brand       string
}

func (c *Controller) Lookup(rw http.ResponseWriter, req *http.Request) {
    itemId := bone.GetValue(req, "id")

    record := &Item{ itemId, "Hersheys Dark Chocolate Kisses, 12-Ounce Bags", "Special Dark Mildly Sweet Chocolate", "Food & Beverage", "Hersheys"}    

    encoder := json.NewEncoder(rw)

    rw.Header().Set("Content-Type", "application/json")
    err := encoder.Encode(record)
    if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
    }
}
