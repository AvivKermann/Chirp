package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const port = "8080"
const filepathRoot = "."

type apiConfig struct {
	fileServerHits int
}

func main() {
	router := chi.NewRouter()
	apiRouter := chi.NewRouter()
	adminRouter := chi.NewRouter()
	router.Mount("/api", apiRouter)
	router.Mount("/admin", adminRouter)

	cfg := apiConfig{
		fileServerHits: 0,
	}

	fsHandler := cfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	router.Handle("/app", fsHandler)
	router.Handle("/app/*", fsHandler)
	apiRouter.Get("/healthz", handlerHealthz)
	adminRouter.Get("/metrics", cfg.handlerMetrics)
	apiRouter.Get("/reset", cfg.handlerReset)

	fmt.Printf("started local host on http://localhost:%s\n", port)
	corsMux := middlewareCors(router)
	server := &http.Server{
		Handler: corsMux,
		Addr:    ":" + port,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("server couldn't run")
	}
}
