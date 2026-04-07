package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load("../backend/.env"); err != nil {
		log.Fatal("No .env found")
	}

	connStr := "host=" + os.Getenv("DB_HOST") + " port=" + os.Getenv("DB_PORT") + " user=" + os.Getenv("DB_USER") + " password=" + os.Getenv("DB_PASSWORD") + " dbname=" + os.Getenv("DB_NAME") + " sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Alter table to add IGTF if missing
	_, err = db.Exec("ALTER TABLE invoices ADD COLUMN IF NOT EXISTS igtf_amount NUMERIC(15,2) DEFAULT 0")
	if err != nil {
		log.Fatal("Error adding igtf_amount:", err)
	}
	_, err = db.Exec("ALTER TABLE invoices ADD COLUMN IF NOT EXISTS igtf_percent NUMERIC(5,2) DEFAULT 0")
	if err != nil {
		log.Fatal("Error adding igtf_percent:", err)
	}

	log.Println("Database schema updated for IGTF.")
}
