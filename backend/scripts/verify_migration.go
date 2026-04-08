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

	fmt.Println("--- Verifying 'documentos' migration ---")
	rows, _ := db.Query("SELECT id, numero_documento, tipo_documento_id, total_neto FROM documentos")
	for rows.Next() {
		var id, tID int
		var num string
		var total float64
		rows.Scan(&id, &num, &tID, &total)
		fmt.Printf("- ID: %d | Num: %s | Tipo: %d | Total: %.2f\n", id, num, tID, total)
	}
}
