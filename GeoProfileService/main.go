package main

import (
	"GeoProfileService/router"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

// setupGlobalMiddleware will setup CORS
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	handleCORS := cors.Default().Handler
	return handleCORS(handler)
}

// our main function
func main() {

	e := godotenv.Load() //Load .env file
	if e != nil {
		log.Print(e)
	}
	certfile := os.Getenv("certfile")
	keyfile := os.Getenv("keyfile")
	// create router and start listen on port 8000
	router := router.NewRouter()
	if certfile != "" && keyfile != "" {
		log.Fatal(http.ListenAndServeTLS(":6001", certfile, keyfile, setupGlobalMiddleware(router)))
	} else {
		log.Fatal(http.ListenAndServe(":6001", setupGlobalMiddleware(router)))
	}
}
