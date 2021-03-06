package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	service "github.com/ocelotconsulting/go-ocelot/api"
	"github.com/ocelotconsulting/go-ocelot/middleware"
	"github.com/ocelotconsulting/go-ocelot/proxy"
	"github.com/ocelotconsulting/go-ocelot/routes"
)

type ports struct {
	serverPort    string
	serverTLSPort string
}

func main() {
	start(os.Args)
}

func start(args []string) {
	config := &ports{
		serverPort:    "0.0.0.0:8080",
		serverTLSPort: "0.0.0.0:8443",
	}

	redisURL := flag.String("redisURL", "redis:6379", "redis url, 'redis:6379'")

	flag.Parse()

	fmt.Println(fmt.Sprintf("running on HTTP: %s, TLS: %s", config.serverPort, config.serverTLSPort))

	//  Start Route Synchronizer
	repo := routes.New(10, *redisURL)
	repo.Start()

	proxy := proxy.New(repo)

	api := service.New(repo)

	mux := http.NewServeMux()
	mux.Handle("/api/", api.Mux())
	mux.HandleFunc("/", proxy.ServeHTTP)

	loggedHandler := middleware.LoggedHandler(mux)
	headeredHandler := middleware.HeaderedHandler(loggedHandler)
	corsHandler := middleware.CORSHandler(headeredHandler)

	//  Start HTTP
	go func() {
		errHTTP := http.ListenAndServe(config.serverPort, corsHandler)
		if errHTTP != nil {
			log.Fatal("HTTP Serving Error: ", errHTTP)
		}
	}()

	// Start TLS
	errTLS := http.ListenAndServeTLS(config.serverTLSPort, "cert.pem", "key.pem", corsHandler)
	if errTLS != nil {
		log.Fatal("TLS Serving Error: ", errTLS)
	}
}
