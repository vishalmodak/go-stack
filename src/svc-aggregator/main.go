package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"flag"
	"github.com/go-kit/kit/endpoint"
	consulsd "github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	"github.com/hashicorp/consul/api"
	"os"
	"svc-aggregator/clients"
	"svc-aggregator/models"
	"time"
)

func main() {
	var (
		httpAddr     = flag.String("http.addr", ":11000", "Address for HTTP (JSON) server")
		consulAddr   = flag.String("consul.addr", "", "Consul agent address")
		retryMax     = flag.Int("retry.max", 3, "per-request retries to different instances")
		retryTimeout = flag.Duration("retry.timeout", 500*time.Millisecond, "per-request timeout, including retries")
	)
	flag.Parse()

	// Service discovery domain. In this example we use Consul.
	var client consulsd.Client
	{
		consulConfig := api.DefaultConfig()
		if len(*consulAddr) > 0 {
			consulConfig.Address = *consulAddr
		}
		consulClient, err := api.NewClient(consulConfig)
		if err != nil {
			//logger.Log("err", err)
			os.Exit(1)
		}
		client = consulsd.NewClient(consulClient)
	}

	r := Router()

	routes := r.Routes()
	for i := 0; i < len(routes); i++ {
		fmt.Printf("%s\n", routes[i].Pattern)
		subRoutes := routes[i].SubRoutes.Routes()
		for j := 0; j < len(subRoutes); j++ {
			fmt.Println("%s\n", subRoutes[j].Pattern)
		}
	}

	{
		var (
			tags        = []string{}
			passingOnly = true
			uppercase   endpoint.Endpoint
			count       endpoint.Endpoint
		)
		{
			factory := stringsvcFactory(ctx, "GET", "/shopping/sync")
			subscriber := consulsd.NewSubscriber(client, factory, logger, "stringsvc", tags, passingOnly)
			balancer := lb.NewRoundRobin(subscriber)
			retry := lb.Retry(*retryMax, *retryTimeout, balancer)
			uppercase = retry
		}
		{
			factory := stringsvcFactory(ctx, "GET", "/shopping/async")
			subscriber := consulsd.NewSubscriber(client, factory, logger, "stringsvc", tags, passingOnly)
			balancer := lb.NewRoundRobin(subscriber)
			retry := lb.Retry(*retryMax, *retryTimeout, balancer)
			count = retry
		}
	}

	// service := chi.ServerBaseContext(r, serverCtx)
	fmt.Printf("svc-aggregator started sucessfully on port %d .....\n", 10000)
	http.ListenAndServe(*httpAddr, r)
}

func getShoppingInfoSync(w http.ResponseWriter, r *http.Request) {
	itemID, ok := r.Context().Value("itemID").(string)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	shopping := syncClientRequests(itemID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&shopping); err != nil {
		panic(err)
	}

}

func syncClientRequests(itemID string) models.ShoppingInfo {
	var shopping = models.ShoppingInfo{[]byte(`{}`), []byte(`{}`)}
	apiResponses := make(chan models.ApiResponse, 2)
	clients.GetItem(itemID, apiResponses)
	clients.GetPrice(itemID, apiResponses)
	close(apiResponses)
	for response := range apiResponses {
		if response.Client == "item" {
			shopping.Item = response.Message
		}
		if response.Client == "price" {
			shopping.Price = response.Message
		}
	}
	return shopping
}

func getShoppingInfoAsync(w http.ResponseWriter, r *http.Request) {
	itemID, ok := r.Context().Value("itemID").(string)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	shopping := asyncClientRequests(itemID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&shopping); err != nil {
		panic(err)
	}

}

func asyncClientRequests(itemID string) models.ShoppingInfo {
	var shopping = models.ShoppingInfo{[]byte(`{}`), []byte(`{}`)}
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
		if response.Client == "item" {
			shopping.Item = response.Message
		}
		if response.Client == "price" {
			shopping.Price = response.Message
		}
	}
	return shopping
}
