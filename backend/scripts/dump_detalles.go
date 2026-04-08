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

	rows, _ := db.Query("SELECT column_name, data_type FROM information_schema.columns WHERE table_name='documentos_detalles' ORDER BY ordinal_position")
	f, _ := os.Create("detalles_schema.txt")
	defer f.Close()
	for rows.Next() {
		var c, t string
		rows.Scan(&c, &t)
		fmt.Fprintf(f, "- %s (%s)\n", c, t)
	}
}
