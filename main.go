package main

import (
	"fmt"
	"log"
	"net/http"
)

const port = "8080"

func main() {
	app := http.NewServeMux()
	corsMux := middlewareCors(app)

	myServer := &http.Server{
		Handler: corsMux,
		Addr:    ":" + port,
	}

	app.Handle("/", http.FileServer(http.Dir(".")))

	fmt.Printf("started local host on http://localhost:%s\n", port)
	err := myServer.ListenAndServe()
	if err != nil {
		log.Fatal("server couldn't run")
	}
}
