package main

import (
	"fmt"
	"net/http"

	"github.com/pratikkarelia25/gosnip/internal/api"
	"github.com/pratikkarelia25/gosnip/internal/store"
)

func main() {
	db, err := store.New("database/gosnip.db")
	if err != nil {
		fmt.Printf("An Error occured connecting to the database: %s", err)
		return
	}
	defer db.Close()

	if err := db.Migrate(); err != nil {
		fmt.Printf("An Error occured creating the table: %s", err)
		return
	}

	handler := api.NewHandler(db)

	if err := http.ListenAndServe(":8080", handler.Routes()); err != nil {
		fmt.Printf("An Error occured starting the server: %s", err)
	}
}
