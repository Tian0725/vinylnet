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

	nf, _ := os.Create("documentos_full_info.txt")
	defer nf.Close()

	rows, err := db.Query("SELECT column_name, is_nullable, column_default FROM information_schema.columns WHERE table_name='documentos' ORDER BY ordinal_position")
	if err != nil {
		log.Fatalf("Error querying documentos: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var c, n string
		var d sql.NullString
		rows.Scan(&c, &n, &d)
		fmt.Fprintf(nf, "- %s | NULLABLE: %s | DEFAULT: %v\n", c, n, d.String)
	}
}
