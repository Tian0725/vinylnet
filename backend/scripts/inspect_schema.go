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
	if err := godotenv.Load("../.env"); err != nil {
		fmt.Println("No se cargo .env de root")
	}

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

	fmt.Println("--- Schema: documentos ---")
	rows, _ := db.Query("SELECT column_name FROM information_schema.columns WHERE table_name='documentos'")
	for rows.Next() {
		var c string
		rows.Scan(&c)
		fmt.Printf("- %s\n", c)
	}

	fmt.Println("\n--- Schema: documentos_detalles ---")
	rows2, _ := db.Query("SELECT column_name FROM information_schema.columns WHERE table_name='documentos_detalles'")
	for rows2.Next() {
		var c string
		rows2.Scan(&c)
		fmt.Printf("- %s\n", c)
	}
}
