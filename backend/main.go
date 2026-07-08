package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "modernc.org/sqlite"
)

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

func databaseInit() {
	db, err := sql.Open("sqlite", "./data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	schema, err := os.ReadFile("../migration/001_init.sql")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(string(schema)); err != nil {
		log.Fatal(err)
	}

	log.Println("Schema applied successfully")
}

func main() {
	//databaseInit()

	http.HandleFunc("/hello", helloWorld)

	log.Println("Server Started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
