package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/AvivKermann/Chirpy/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

const port = "8080"
const filepathRoot = "."
const dbFilePath = "./database.json"

type apiConfig struct {
	fileServerHits int
	DB             *database.DB
	jwtSecret      string
}

func main() {
	router := chi.NewRouter()
	apiRouter := chi.NewRouter()
	adminRouter := chi.NewRouter()

	router.Mount("/api", apiRouter)
	router.Mount("/admin", adminRouter)

	godotenv.Load()

	db, err := database.NewDB(dbFilePath)
	if err != nil {
		log.Fatal("db cannot be loaded")
	}

	cfg := apiConfig{
		fileServerHits: 0,
		DB:             db,
		jwtSecret:      os.Getenv("JWT_SECRET"),
	}

	fsHandler := cfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	router.Handle("/app", fsHandler)
	router.Handle("/app/*", fsHandler)

	apiRouter.Get("/healthz", handlerHealthz)
	apiRouter.Get("/reset", cfg.handlerReset)
	apiRouter.Get("/chirps", cfg.handlerGetChirps)
	apiRouter.Get("/chirps/{chirpId}", cfg.handlerGetSingleChirp)
	apiRouter.Post("/chirps", cfg.handlerCreateChirp)
	apiRouter.Post("/users", cfg.handlerCreateUser)
	apiRouter.Post("/login", cfg.handlerLogin)
	apiRouter.Post("/refresh", cfg.handlerRefresh)
	apiRouter.Post("/revoke", cfg.handlerRevoke)
	apiRouter.Put("/users", cfg.handlerUpdateUser)

	adminRouter.Get("/metrics", cfg.handlerMetrics)

	fmt.Printf("started local host on http://localhost:%s\n", port)
	corsMux := middlewareCors(router)
	server := &http.Server{
		Handler: corsMux,
		Addr:    ":" + port,
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("server couldn't run")
	}
}
