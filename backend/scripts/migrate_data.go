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

	fmt.Println("--- Migration: Moving invoices to documentos ---")
	
	// Query current invoices
	rows, err := db.Query("SELECT id, client_id, user_id, number, issue_date, total, subtotal, tax_amount, tax_percent, igtf_amount, igtf_percent, currency, notes, print_format FROM invoices")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var oldID, clientID, userID int
		var number, currency, notes, printFormat string
		var issueDate string
		var total, subtotal, taxAmount, taxPercent, igtfAmount, igtfPercent float64

		rows.Scan(&oldID, &clientID, &userID, &number, &issueDate, &total, &subtotal, &taxAmount, &taxPercent, &igtfAmount, &igtfPercent, &currency, &notes, &printFormat)

		// Determine tipo_documento_id
		tipoID := 1 // FAC
		if printFormat == "FAN" {
			tipoID = 5
		}

		// Insert into documentos
		var newID int
		err = db.QueryRow(`
			INSERT INTO documentos (
				numero_documento, tipo_documento_id, cliente_id, usuario_id, 
				fecha_emision, total_neto, subtotal, monto_iva, monto_igtf, base_imponible
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`,
			number, tipoID, clientID, userID, issueDate, total, subtotal, taxAmount, igtfAmount, subtotal,
		).Scan(&newID)

		if err != nil {
			fmt.Printf("Error migrating invoice %s: %v\n", number, err)
			continue
		}

		// Migrate items
		itemRows, _ := db.Query("SELECT product_id, quantity, price, subtotal FROM invoice_items WHERE invoice_id = $1", oldID)
		for itemRows.Next() {
			var prodID int
			var qty, price, itemSub float64
			itemRows.Scan(&prodID, &qty, &price, &itemSub)

			_, err = db.Exec(`
				INSERT INTO documentos_detalles (
					documento_id, producto_id, cantidad, precio_unitario, subtotal, total
				) VALUES ($1, $2, $3, $4, $5, $6)`,
				newID, prodID, qty, price, itemSub, itemSub,
			)
			if err != nil {
				fmt.Printf("  Error migrating item for %s: %v\n", number, err)
			}
		}
		itemRows.Close()
		fmt.Printf("Migrated invoice %s -> documento ID %d\n", number, newID)
	}
}
