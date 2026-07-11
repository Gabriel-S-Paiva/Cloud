package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"backend/handlers"
	"backend/storage"

	_ "modernc.org/sqlite"
)

func databaseInit() {
	db, err := sql.Open("sqlite", "./data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	schema, err := os.ReadFile("../migration/003_auth.sql")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(string(schema)); err != nil {
		log.Fatal(err)
	}

	log.Println("Schema applied successfully")
}

func main() {
	_, err := os.Stat("./data.db")
	if os.IsNotExist(err) {
		databaseInit()
	} else if err == nil {
		fmt.Println("Trying to Load DB")
	} else {
		log.Fatal("Error: ", err)
	}

	db, err := sql.Open("sqlite", "./data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := storage.NewStore(db)
	userHandlers := handlers.NewUserHandlers(store)
	authHandlers := handlers.NewAuthHandlers(store)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", userHandlers.Register)
	mux.HandleFunc("GET /users/me", userHandlers.GetMe)
	mux.HandleFunc("GET /users/requests", userHandlers.ListPendingRequests)
	mux.HandleFunc("POST /users/requests/{id}/aprove", userHandlers.AproveRequest)
	mux.HandleFunc("POST /users/requests/{id}/reject", userHandlers.RejectRequest)
	mux.HandleFunc("POST /login", authHandlers.Login)

	log.Println("server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
