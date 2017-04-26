package clients

import (
    "fmt"
    "io/ioutil"
    "net/http"
)


func GetItem(itemID string) ([]byte, error) {
  item, err := GetContent(fmt.Sprintf("http://localhost:10001/item/lookup/%s", itemID))
  if err != nil {
    // An error occurred while fetching the JSON
    return nil, err
  }
  // var itemJSON
  // err = json.Unmarshal(item, &itemJSON)
  // if err != nil {
  //   // An error occurred while converting our JSON to an object
  //   log.Fatal(err)
  // }
  // log.Print(item)
  return item, nil;
}


func GetContent(url string) ([]byte, error) {
  // Build the request
  req, err := http.NewRequest("GET", url, nil)
  if err != nil {
    return nil, err
  }
  // Send the request via a client
  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    return nil, err
  }
  // Defer the closing of the body
  defer resp.Body.Close()
  // Read the content into a byte array
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return nil, err
  }
  // At this point we're done - simply return the bytes
  return body, nil
}
