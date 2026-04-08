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

	fmt.Println("--- Checking for Return related tables ---")
	rows, _ := db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema='public' AND table_name LIKE '%return%' OR table_name LIKE '%devolucion%'")
	for rows.Next() {
		var t string
		rows.Scan(&t)
		fmt.Println("- " + t)
	}
	rows.Close()

	fmt.Println("\n--- Checking invoices columns for references ---")
	rows, _ = db.Query("SELECT column_name FROM information_schema.columns WHERE table_name='invoices' AND (column_name LIKE '%parent%' OR column_name LIKE '%original%' OR column_name LIKE '%ref%')")
	for rows.Next() {
		var c string
		rows.Scan(&c)
		fmt.Println("- " + c)
	}
}
