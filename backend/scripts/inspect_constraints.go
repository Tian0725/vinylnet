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

	fmt.Println("--- Checking unique constraints on 'documentos' ---")
	rows, _ := db.Query(`
		SELECT conname, pg_get_constraintdef(c.oid)
		FROM pg_constraint c
		JOIN pg_namespace n ON n.oid = c.connamespace
		WHERE n.nspname = 'public' AND c.conrelid = 'documentos'::regclass
	`)
	for rows.Next() {
		var name, def string
		rows.Scan(&name, &def)
		fmt.Printf("- %s: %s\n", name, def)
	}
}
