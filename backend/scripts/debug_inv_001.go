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
	// Root .env
	if err := godotenv.Load("../.env"); err != nil {
		fmt.Println("No se cargo .env de root")
	}

	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, name)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("--- Buscando Documento INV-00001 ---")
	var id int
	var numero string
	var tipo int
	err = db.QueryRow("SELECT id, numero_documento, tipo_documento_id FROM documentos WHERE numero_documento = 'INV-00001'").Scan(&id, &numero, &tipo)
	if err != nil {
		fmt.Printf("Error o No Encontrado: %v\n", err)
		return
	}
	fmt.Printf("ID: %d, Numero: %s, Tipo: %d\n", id, numero, tipo)

	fmt.Println("\n--- Buscando Detalles ---")
	rows, err := db.Query("SELECT id, producto_id, cantidad, total FROM documentos_detalles WHERE documento_id = $1", id)
	if err != nil {
		fmt.Printf("Error en detalles: %v\n", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var detID, prodID int
		var cant, total float64
		rows.Scan(&detID, &prodID, &cant, &total)
		fmt.Printf("Detalle ID: %d, Producto: %d, Cantidad: %v, Total: %v\n", detID, prodID, cant, total)
		count++
	}
	fmt.Printf("Total detalles encontrados: %d\n", count)
}
