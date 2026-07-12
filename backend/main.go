package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"backend/handlers"
	"backend/middlewares"
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
	authMW := middlewares.NewAuthMiddlewares(store)

	mux := http.NewServeMux()
	// Open Endpoints
	mux.HandleFunc("POST /users", userHandlers.Register)
	mux.HandleFunc("POST /login", authHandlers.Login)
	// Auth Endpoints
	mux.HandleFunc("GET /users/me", authMW.RequireAuth(userHandlers.GetMe))
	mux.HandleFunc("POST /logout", authHandlers.Logout)
	// Admin Endpoints
	mux.HandleFunc("GET /users/requests", authMW.RequireAuth(authMW.RequireAdmin(userHandlers.ListPendingRequests)))
	mux.HandleFunc("POST /users/requests/{id}/aprove", authMW.RequireAuth(authMW.RequireAdmin(userHandlers.AproveRequest)))
	mux.HandleFunc("POST /users/requests/{id}/reject", authMW.RequireAuth(authMW.RequireAdmin(userHandlers.RejectRequest)))

	log.Println("server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
