package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Srgkharkov/recordingmeet/internal/auth"
	"github.com/Srgkharkov/recordingmeet/internal/handlers"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

func main() {
	if len(jwtKey) == 0 {
		log.Fatal("JWT_SECRET_KEY is not set")
	}

	http.Handle("/record", auth.JWTMiddleware(http.HandlerFunc(handlers.HandleRecordRequest)))
	http.Handle("/download", auth.JWTMiddleware(http.HandlerFunc(handlers.HandleDownloadRequest)))

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
