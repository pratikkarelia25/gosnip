package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/pratikkarelia25/gosnip/internal/api"
	"github.com/pratikkarelia25/gosnip/internal/store"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" || dbUser == "" || dbPassword == "" || dbName == "" {
		fmt.Println("DB_HOST, DB_USER, DB_PASSWORD and DB_NAME must all be set")
		return
	}

	db, err := store.NewFromParts(dbHost, dbUser, dbPassword, dbName)
	if err != nil {
		fmt.Printf("An Error occured connecting to the database: %s", err)
		return
	}
	defer db.Close()

	if err := db.Migrate(); err != nil {
		fmt.Printf("An Error occured creating the table: %s", err)
		return
	}

	handler := api.NewHandler(db, frontendURL)

	if err := http.ListenAndServe(":"+port, handler.Routes()); err != nil {
		fmt.Printf("An Error occured starting the server: %s", err)
	}
}
