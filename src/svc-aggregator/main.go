package main

import (
	"net/http"
  "context"
  "encoding/json"
  "sync"
  "github.com/pressly/chi"

  "svc-aggregator/clients"
  "svc-aggregator/models"
)

func main() {
  r := Router()
  // service := chi.ServerBaseContext(r, serverCtx)
	http.ListenAndServe(":10000", r)
}

func ShoppingCtx(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    itemID := chi.URLParam(r, "itemID")

    shopping := asyncClientRequests(itemID)

    ctx := context.WithValue(r.Context(), "shopping", shopping)
    next.ServeHTTP(w, r.WithContext(ctx))
  })
}

func asyncClientRequests(itemID string) models.ShoppingInfo {
  var shopping = models.ShoppingInfo{ []byte(`{}`), []byte(`{}`)}
  apiResponses := make(chan models.ApiResponse, 2)
  var wg sync.WaitGroup
  wg.Add(2)

  go func() {
    defer wg.Done()
    clients.GetItem(itemID, apiResponses)
  }()

  go func() {
    defer wg.Done()
    clients.GetPrice(itemID, apiResponses)
  }()

  wg.Wait()
  close(apiResponses)

  for response := range apiResponses {
    if (response.Client == "item") {
      shopping.Item = response.Message
    }
    if (response.Client == "price") {
      shopping.Price = response.Message
    }
  }
  return shopping
}

func getShoppingInfo(w http.ResponseWriter, r *http.Request) {
  shopping, ok := r.Context().Value("shopping").(models.ShoppingInfo)
  if !ok {
    http.Error(w, http.StatusText(422), 422)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  if err := json.NewEncoder(w).Encode(&shopping); err != nil {
    panic(err)
  }

}
