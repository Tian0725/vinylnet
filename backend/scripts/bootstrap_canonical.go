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

	fmt.Println("--- Bootstrapping Canonical System ---")

	// 1. Create Default Sucursal if not exists
	var sucursalID int
	err = db.QueryRow("SELECT id FROM sucursales WHERE nombre = 'Sede Principal' LIMIT 1").Scan(&sucursalID)
	if err != nil {
		fmt.Println("Sede Principal not found, creating...")
		// Added codigo_unico which is NOT NULL
		err = db.QueryRow("INSERT INTO sucursales (nombre, es_principal, codigo_unico) VALUES ('Sede Principal', true, 'SEDE-001') RETURNING id").Scan(&sucursalID)
		if err != nil {
			log.Fatalf("Fatal error creating sucursal: %v", err)
		}
	}
	fmt.Printf("Sucursal ID: %d\n", sucursalID)

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	// 2. Migrate Clients -> clientes
	rows, _ := tx.Query("SELECT id, name, rif FROM clients")
	for rows.Next() {
		var id int
		var n, r string
		rows.Scan(&id, &n, &r)
		_, err = tx.Exec("INSERT INTO clientes (id, nombre_razon_social, rif_cedula) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET nombre_razon_social = EXCLUDED.nombre_razon_social", id, n, r)
		if err != nil {
			fmt.Printf("Error migrating client %d: %v\n", id, err)
		}
	}
	rows.Close()
	fmt.Println("Clients migrated to 'clientes'")

	// 3. Migrate Products -> productos & productos_sucursal
	rows, _ = tx.Query("SELECT id, name, price, stock FROM products")
	for rows.Next() {
		var id int
		var n string
		var p, s float64
		rows.Scan(&id, &n, &p, &s)
		
		// Insert into productos
		_, err = tx.Exec("INSERT INTO productos (id, descripcion, codigo) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET descripcion = EXCLUDED.descripcion", id, n, fmt.Sprintf("PROD-%d", id))
		if err != nil {
			fmt.Printf("Error migrating product %d: %v\n", id, err)
			continue
		}

		// Insert into productos_sucursal
		_, err = tx.Exec("INSERT INTO productos_sucursal (producto_id, sucursal_id, existencia_total, precio_venta_1) VALUES ($1, $2, $3, $4) ON CONFLICT (producto_id, sucursal_id) DO UPDATE SET precio_venta_1 = EXCLUDED.precio_venta_1", id, sucursalID, s, p)
		if err != nil {
			// In case direct INSERT fails due to missing unique constraint, we check if we should add it? 
			// I'll assume it exists or I'll just try to insert.
			fmt.Printf("Error migrating product_sucursal %d: %v\n", id, err)
		}
	}
	rows.Close()
	fmt.Println("Products migrated to 'productos' and 'productos_sucursal'")

	// 4. Migrate Invoices -> documentos & documentos_detalles
	fmt.Println("Querying invoices...")
	invoiceRows, err := tx.Query("SELECT id, client_id, user_id, number, issue_date, total, subtotal, tax_amount, tax_percent, igtf_amount, igtf_percent, currency, notes, print_format FROM invoices")
	if err == nil {
		count := 0
		for invoiceRows.Next() {
			count++
			var oldID, clientID, userID int
			var number, currency, notes, printFormat string
			var issueDate string
			var total, subtotal, taxAmount, taxPercent, igtfAmount, igtfPercent float64

			err = invoiceRows.Scan(&oldID, &clientID, &userID, &number, &issueDate, &total, &subtotal, &taxAmount, &taxPercent, &igtfAmount, &igtfPercent, &currency, &notes, &printFormat)
			if err != nil {
				fmt.Printf("Error scanning invoice row: %v\n", err)
				continue
			}
			fmt.Printf("Processing invoice: %s\n", number)

			tipoID := 1 // FAC
			if printFormat == "FAN" {
				tipoID = 5
			}

			var newID int
			err = tx.QueryRow(`
				INSERT INTO documentos (
					numero_documento, tipo_documento_id, cliente_id, sucursal_id, usuario_id, 
					fecha_emision, total_neto, subtotal, monto_iva, monto_igtf, base_imponible
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`,
				number, tipoID, clientID, sucursalID, userID, issueDate, total, subtotal, taxAmount, igtfAmount, subtotal,
			).Scan(&newID)

			if err != nil {
				fmt.Printf("Error migrating document %s: %v\n", number, err)
				continue
			}

			// Items
			itemRows, _ := tx.Query("SELECT product_id, quantity, price, subtotal FROM invoice_items WHERE invoice_id = $1", oldID)
			if itemRows != nil {
				for itemRows.Next() {
					var prodID int
					var qty, price, itemSub float64
					itemRows.Scan(&prodID, &qty, &price, &itemSub)

					_, err = tx.Exec(`
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
			}
		}
		invoiceRows.Close()
	} else {
		fmt.Printf("Error querying invoices: %v\n", err)
	}

	err = tx.Commit()
	if err != nil {
		fmt.Printf("Commit failed: %v\n", err)
	} else {
		fmt.Println("--- Bootstrap Complete and Committed ---")
	}
}
