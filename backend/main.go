package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"backend/storage"

	"github.com/joho/godotenv"

	_ "modernc.org/sqlite"
)

func databaseInit() {
	db, err := sql.Open("sqlite", "./data/data.db")
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

	_, err := os.Stat("./data/data.db")
	if os.IsNotExist(err) {
		databaseInit()
	} else if err == nil {
		fmt.Println("Trying to Load DB")
	} else {
		log.Fatal("Error: ", err)
	}

	db, err := sql.Open("sqlite", "./data/data.db")
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

	log.Println("server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", NewRouter(store)))
}
