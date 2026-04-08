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

	fmt.Println("--- Migration: Adding reference columns to 'documentos' ---")
	_, err = db.Exec("ALTER TABLE documentos ADD COLUMN IF NOT EXISTS numero_referencia character varying(20)")
	if err != nil {
		fmt.Println("Error adding numero_referencia:", err)
	} else {
		fmt.Println("numero_referencia column OK")
	}

	_, err = db.Exec("ALTER TABLE documentos ADD COLUMN IF NOT EXISTS documento_referencia_id integer")
	if err != nil {
		fmt.Println("Error adding documento_referencia_id:", err)
	} else {
		fmt.Println("documento_referencia_id column OK")
	}
}
