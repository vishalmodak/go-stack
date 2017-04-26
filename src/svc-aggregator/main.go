package main

import (
	"net/http"
  "time"
  "log"
  "fmt"
	"github.com/pressly/chi"
  "github.com/pressly/chi/middleware"
  "context"

  "svc-aggregator/clients"
)

func main() {
	r := chi.NewRouter()

  // Setup the logger
	// logger := logrus.New()
	// logger.Formatter = &logrus.TextFormatter{}
  //
	// lg.RedirectStdlogOutput(logger)
	// lg.DefaultLogger = logger
  //
  // lg.Infoln("Welcome")
  //
  // serverCtx := context.Background()
	// serverCtx = lg.WithLoggerContext(serverCtx, logger)
	// lg.Log(serverCtx).Infof("Booting up server, %s", "v1.0")

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

  // service := chi.ServerBaseContext(r, serverCtx)
	http.ListenAndServe(":10000", r)
}

type ShoppingInfo struct {
  item []byte `json:"item"`
  price []byte `json:"price"`
}

func ShoppingCtx(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    itemID := chi.URLParam(r, "itemID")

    item, err := clients.GetItem(itemID)
    if err != nil {
      log.Fatal(err)
      http.Error(w, http.StatusText(404), 404)
      return
    }
    // priceJSON, err := clients.GetPrice(itemID)
    // if err != nil {
    //   http.Error(w, http.StatusText(404), 404)
    //   return
    // }
    shopping := ShoppingInfo{item, nil}
    ctx := context.WithValue(r.Context(), "shopping", shopping)
    next.ServeHTTP(w, r.WithContext(ctx))
  })
}

func getShoppingInfo(w http.ResponseWriter, r *http.Request) {
  shopping, ok := r.Context().Value("shopping").(ShoppingInfo)
  log.Println(fmt.Sprintf("%s", shopping))
  if !ok {
    http.Error(w, http.StatusText(422), 422)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  w.Write([]byte(fmt.Sprintf("%s", shopping)))
	// if err := json.NewEncoder(w).Encode(shopping); err != nil {
	// 	panic(err)
	// }
}
