package clients

import (
    "fmt"
    "github.com/go-resty/resty"

    "svc-aggregator/models"
)


func GetItem(itemID string, apiResponses chan models.ApiResponse) {
  r := models.ApiResponse{"item", nil}
  item, _ := resty.R().SetHeader("Accept", "application/json").Get(fmt.Sprintf("http://localhost:10001/item/lookup/%s", itemID))
  if item.StatusCode() != 200 {
    // An error occurred while fetching the JSON
    apiResponses <- r
  }
  r.Message = item.Body()
  apiResponses <- r
}
