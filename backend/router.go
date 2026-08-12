package main

import (
	"net/http"
	"time"

	"backend/handlers"
	"backend/middlewares"
	"backend/storage"
)

func NewRouter(store *storage.Store) http.Handler {
	userHandlers := handlers.NewUserHandlers(store)
	authHandlers := handlers.NewAuthHandlers(store)
	foldHandlers := handlers.NewFolderHandlers(store)
	fileHandler := handlers.NewFileHandler(store)
	shareHandler := handlers.NewShareHandlers(store)
	authMW := middlewares.NewAuthMiddlewares(store)

	loginLimiter := middlewares.NewRateLimiter(5, time.Minute)
	registerLimiter := middlewares.NewRateLimiter(5, time.Minute)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", registerLimiter.Limit(userHandlers.Register))
	mux.HandleFunc("POST /login", loginLimiter.Limit(authHandlers.Login))
	mux.HandleFunc("GET /users/me", authMW.RequireAuth(userHandlers.GetMe))
	mux.HandleFunc("GET /users/summary", authMW.RequireAuth(userHandlers.ListSharableUsers))
	mux.HandleFunc("POST /logout", authHandlers.Logout)

	mux.HandleFunc("POST /folders", authMW.RequireAuth(foldHandlers.CreateFolder))
	mux.HandleFunc("GET /folders/{id}", authMW.RequireAuth(foldHandlers.GetFolder))
	mux.HandleFunc("GET /folders/{id}/content", authMW.RequireAuth(foldHandlers.GetFolderContents))
	mux.HandleFunc("PATCH /folders/{id}", authMW.RequireAuth(foldHandlers.UpdateFolder))
	mux.HandleFunc("DELETE /folders/{id}", authMW.RequireAuth(foldHandlers.DeleteFolder))

	mux.HandleFunc("POST /files", authMW.RequireAuth(fileHandler.CreateFile))
	mux.HandleFunc("POST /files/{id}/chunk", authMW.RequireAuth(fileHandler.UploadChunk))
	mux.HandleFunc("GET /files/{id}", authMW.RequireAuth(fileHandler.GetFile))
	mux.HandleFunc("GET /files/{id}/content", authMW.RequireAuth(fileHandler.GetFileContent))
	mux.HandleFunc("PATCH /files/{id}", authMW.RequireAuth(fileHandler.UpdateFile))
	mux.HandleFunc("DELETE /files/{id}", authMW.RequireAuth(fileHandler.DeleteFile))

	mux.HandleFunc("POST /shares", authMW.RequireAuth(shareHandler.CreateShare))
	mux.HandleFunc("GET /shares/incoming", authMW.RequireAuth(shareHandler.ViewIncomingShares))
	mux.HandleFunc("GET /shares/outgoing", authMW.RequireAuth(shareHandler.ViewOutgoingShares))
	mux.HandleFunc("PATCH /shares/{id}", authMW.RequireAuth(shareHandler.UpdatePermission))
	mux.HandleFunc("DELETE /shares/{id}", authMW.RequireAuth(shareHandler.DeleteShare))

	mux.HandleFunc("GET /users/requests", authMW.RequireAuth(authMW.RequireAdmin(userHandlers.ListPendingRequests)))
	mux.HandleFunc("GET /users", authMW.RequireAuth(authMW.RequireAdmin(userHandlers.ListUser)))
	mux.HandleFunc("POST /users/requests/{id}/aprove", authMW.RequireAuth(authMW.RequireAdmin(userHandlers.AproveRequest)))
	mux.HandleFunc("POST /users/requests/{id}/reject", authMW.RequireAuth(authMW.RequireAdmin(userHandlers.RejectRequest)))

	return middlewares.CORS(mux)
}
