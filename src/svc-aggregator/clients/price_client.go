package clients

import (
  "fmt"
  "github.com/go-resty/resty"

  "svc-aggregator/models"
)


func GetPrice(itemID string, apiResponses chan models.ApiResponse) {
  r := models.ApiResponse{"price", nil}
  price, _ := resty.R().SetHeader("Accept", "application/json").Get(fmt.Sprintf("http://localhost:10002/price/item/%s", itemID))
  if price.StatusCode() != 200 {
    // An error occurred while fetching the JSON
    apiResponses <- r
  }
  r.Message = price.Body()
  apiResponses <- r
}
