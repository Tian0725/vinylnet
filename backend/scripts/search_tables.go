package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load("../backend/.env")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, name)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("--- Searching for tables with 'tipo' or 'documento' in their name ---")
	rows, _ := db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema='public' AND (table_name LIKE '%tipo%' OR table_name LIKE '%documento%')")
	for rows.Next() {
		var t string
		rows.Scan(&t)
		fmt.Println("- " + t)
	}
}
