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

	f, _ := os.Create("not_null_columns.txt")
	defer f.Close()

	tables := []string{"sucursales", "documentos", "documentos_detalles", "clientes", "productos", "productos_sucursal"}
	for _, t := range tables {
		fmt.Fprintf(f, "\nTable: %s\n", t)
		rows, err := db.Query(fmt.Sprintf("SELECT column_name FROM information_schema.columns WHERE table_name='%s' AND is_nullable='NO'", t))
		if err != nil {
			fmt.Printf("Error querying table %s: %v\n", t, err)
			continue
		}
		for rows.Next() {
			var c string
			rows.Scan(&c)
			fmt.Fprintf(f, "- %s\n", c)
		}
		rows.Close()
	}
}
