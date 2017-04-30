package main

import (
  "github.com/pressly/chi"
  "github.com/pressly/chi/middleware"
  "time"
  "net/http"
)

func Router() chi.Router {
  r := chi.NewRouter()

  // A good base middleware stack
  //Injects a request ID into the context of each request.
  r.Use(middleware.RequestID)
  //Sets a http.Request's RemoteAddr to either X-Forwarded-For or X-Real-IP.
  r.Use(middleware.RealIP)
  //Logs the start and end of each request with the elapsed processing time.
  r.Use(middleware.Logger)
  //Gracefully absorb panics and prints the stack trace.
  r.Use(middleware.Recoverer)

  // Set a timeout value on the request context (ctx), that will signal
  // through ctx.Done() that the request has timed out and further
  // processing should be stopped.
  r.Use(middleware.Timeout(60 * time.Second))

  r.Get("/", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("welcome"))
  })

  r.Route("/shopping", func(r chi.Router) {
    r.Route("/:itemID", func(r chi.Router) {
      r.Use(ShoppingCtx)
      r.Get("/", getShoppingInfo)
    })
  })

  return r;
}
