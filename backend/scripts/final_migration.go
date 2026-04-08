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

	fmt.Println("--- Final Resilient Migration ---")

	// 1. Get/Create Sucursal
	var sucursalID int
	db.QueryRow("SELECT id FROM sucursales WHERE nombre = 'Sede Principal' LIMIT 1").Scan(&sucursalID)
	if sucursalID == 0 {
		db.QueryRow("INSERT INTO sucursales (nombre, es_principal, codigo_unico) VALUES ('Sede Principal', true, 'SEDE-001') RETURNING id").Scan(&sucursalID)
	}
	fmt.Printf("Using Sucursal ID: %d\n", sucursalID)

	// 2. Migrate Clients (Ignore errors)
	db.Exec("INSERT INTO clientes (id, nombre_razon_social, rif_cedula) SELECT id, name, rif FROM clients ON CONFLICT (id) DO NOTHING")
	
	// 3. Migrate Products (Ignore errors)
	db.Exec("INSERT INTO productos (id, descripcion, codigo) SELECT id, name, 'PROD-' || id FROM products ON CONFLICT (id) DO NOTHING")
	db.Exec(`INSERT INTO productos_sucursal (producto_id, sucursal_id, existencia_total, precio_venta_1) 
	         SELECT id, $1, stock, price FROM products 
			 ON CONFLICT (producto_id, sucursal_id) DO NOTHING`, sucursalID)

	// 4. Migrate Invoices
	rows, err := db.Query("SELECT id, client_id, user_id, number, issue_date, total, subtotal, tax_amount, igtf_amount, print_format FROM invoices")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var oldID, clientID, userID int
		var number, printFormat, issueDate string
		var total, subtotal, taxAmount, igtfAmount float64

		rows.Scan(&oldID, &clientID, &userID, &number, &issueDate, &total, &subtotal, &taxAmount, &igtfAmount, &printFormat)

		tipoID := 1 // FAC
		if printFormat == "FAN" {
			tipoID = 5
		}

		// Check if already exists
		var existingID int
		db.QueryRow("SELECT id FROM documentos WHERE numero_documento = $1 AND tipo_documento_id = $2", number, tipoID).Scan(&existingID)
		if existingID > 0 {
			fmt.Printf("Documento %s already exists, skipping.\n", number)
			continue
		}

		var newID int
		err = db.QueryRow(`
			INSERT INTO documentos (
				numero_documento, tipo_documento_id, cliente_id, sucursal_id, usuario_id, 
				fecha_emision, total_neto, subtotal, monto_iva, monto_igtf, base_imponible
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) 
			RETURNING id`,
			number, tipoID, clientID, sucursalID, userID, issueDate, total, subtotal, taxAmount, igtfAmount, subtotal,
		).Scan(&newID)

		if err != nil {
			fmt.Printf("Error migrating %s: %v\n", number, err)
			continue
		}

		// Items
		db.Exec(`INSERT INTO documentos_detalles (documento_id, producto_id, cantidad, precio_unitario, subtotal, total)
		         SELECT $1, product_id, quantity, price, subtotal, subtotal FROM invoice_items WHERE invoice_id = $2`, newID, oldID)

		fmt.Printf("Documento %s centralizado.\n", number)
	}

	fmt.Println("--- Migration Finished ---")
}
