package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"backend/handlers"
	"backend/middlewares"
	"backend/storage"

	"github.com/joho/godotenv"

	_ "modernc.org/sqlite"
)

func databaseInit() {
	db, err := sql.Open("sqlite", "./data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	schema, err := os.ReadFile("../migration/005_user_root.sql")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(string(schema)); err != nil {
		log.Fatal(err)
	}

	log.Println("Schema applied successfully")
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on system environment")
	}

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
	adminUsername := os.Getenv("ADMIN_USERNAME")
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminUsername != "" && adminPassword != "" {
		if err := store.SeedAdmin(context.Background(), adminUsername, adminPassword); err != nil {
			log.Fatal(err)
		}
	}

	userHandlers := handlers.NewUserHandlers(store)
	authHandlers := handlers.NewAuthHandlers(store)
	foldHandlers := handlers.NewFolderHandlers(store)
	fileHanlder := handlers.NewFileHandler(store)
	shareHandler := handlers.NewShareHandlers(store)
	authMW := middlewares.NewAuthMiddlewares(store)

	mux := http.NewServeMux()
	// Open Endpoints
	mux.HandleFunc("POST /users", userHandlers.Register)
	mux.HandleFunc("POST /login", authHandlers.Login)
	// Auth Endpoints
	mux.HandleFunc("GET /users/me", authMW.RequireAuth(userHandlers.GetMe))
	mux.HandleFunc("POST /logout", authHandlers.Logout)

	mux.HandleFunc("POST /folders", authMW.RequireAuth(foldHandlers.CreateFolder))
	mux.HandleFunc("GET /folders/{id}", authMW.RequireAuth(foldHandlers.GetFolder))
	mux.HandleFunc("GET /folders/{id}/content", authMW.RequireAuth(foldHandlers.GetFolderContents))
	mux.HandleFunc("PATCH /folders/{id}", authMW.RequireAuth(foldHandlers.UpdateFolder))
	mux.HandleFunc("DELETE /folders/{id}", authMW.RequireAuth(foldHandlers.DeleteFolder))

	mux.HandleFunc("POST /files", authMW.RequireAuth(fileHanlder.CreateFile))
	mux.HandleFunc("POST /files/{id}/chunk", authMW.RequireAuth(fileHanlder.UploadChunk))
	mux.HandleFunc("GET /files/{id}", authMW.RequireAuth(fileHanlder.GetFile))
	mux.HandleFunc("GET /files/{id}/content", authMW.RequireAuth(fileHanlder.GetFileContent))
	mux.HandleFunc("PATCH /files/{id}", authMW.RequireAuth(fileHanlder.UpdateFile))
	mux.HandleFunc("DELETE /files/{id}", authMW.RequireAuth(fileHanlder.DeleteFile))

	mux.HandleFunc("POST /shares", authMW.RequireAuth(shareHandler.CreateShare))
	mux.HandleFunc("GET /shares/incoming", authMW.RequireAuth(shareHandler.ViewIncomingShares))
	mux.HandleFunc("GET /shares/outgoing", authMW.RequireAuth(shareHandler.ViewOutgoingShares))
	mux.HandleFunc("PATCH /shares/{id}", authMW.RequireAuth(shareHandler.UpdatePermission))
	mux.HandleFunc("DELETE /shares/{id}", authMW.RequireAuth(shareHandler.DeleteShare))
	// Admin Endpoints
	mux.HandleFunc("GET /users/requests", authMW.RequireAuth(authMW.RequireAdmin(userHandlers.ListPendingRequests)))
	mux.HandleFunc("POST /users/requests/{id}/aprove", authMW.RequireAuth(authMW.RequireAdmin(userHandlers.AproveRequest)))
	mux.HandleFunc("POST /users/requests/{id}/reject", authMW.RequireAuth(authMW.RequireAdmin(userHandlers.RejectRequest)))

	log.Println("server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", middlewares.CORS(mux)))
}
