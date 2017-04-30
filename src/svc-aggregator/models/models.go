package models

import (
  "encoding/json"
)

type ShoppingInfo struct {
  Item json.RawMessage `json:"item"`
  Price json.RawMessage `json:"price"`
}

type ApiResponse struct {
  Client string
  Message json.RawMessage
}
