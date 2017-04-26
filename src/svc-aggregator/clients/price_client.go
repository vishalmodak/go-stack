package clients

import (
  "fmt"
  "encoding/json"
)


func GetPrice(itemID string) (string, error) {
  price, err := GetContent(fmt.Sprintf("http://localhost:10002/price/item/%s", itemID))
  if err != nil {
    // An error occurred while fetching the JSON
    return "", err
  }
  var priceJSON string
  err = json.Unmarshal(price, &priceJSON)
  if err != nil {
    // An error occurred while converting our JSON to an object
    return "", err
  }
  return priceJSON, nil;
}
