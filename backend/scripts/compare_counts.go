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

	var c1, c2 int
	db.QueryRow("SELECT COUNT(*) FROM clients").Scan(&c1)
	db.QueryRow("SELECT COUNT(*) FROM clientes").Scan(&c2)
	fmt.Printf("clients: %d | clientes: %d\n", c1, c2)

	db.QueryRow("SELECT COUNT(*) FROM products").Scan(&c1)
	db.QueryRow("SELECT COUNT(*) FROM productos").Scan(&c2)
	fmt.Printf("products: %d | productos: %d\n", c1, c2)
}
