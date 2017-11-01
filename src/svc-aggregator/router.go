package main

import (
	"context"
	"encoding/json"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"net/http"
	"time"
	"fmt"
)

type Status struct {
	Status string
}

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

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(&Status{"UP"}); err != nil {
			panic(err)
		}
	})

	r.Route("/shopping", func(r chi.Router) {
		r.Route("/sync/:itemID", func(r chi.Router) {
			r.Use(ShoppingCtx)
			r.Get("/", getShoppingInfoSync)
		})
		r.Route("/async/:itemID", func(r chi.Router) {
			r.Use(ShoppingCtx)
			r.Get("/", getShoppingInfoAsync)
		})
	})

	routes := r.Routes()
	for i := 0; i < len(routes); i++ {
		fmt.Printf("%s\n", routes[i].Pattern)
		subRoutes := routes[i].SubRoutes.Routes()
		for j := 0; j < len(subRoutes); j++ {
			fmt.Println("%s\n", subRoutes[j].Pattern)
		}
	}

	return r
}

func ShoppingCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		itemID := chi.URLParam(r, "itemID")

		ctx := context.WithValue(r.Context(), "itemID", itemID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
