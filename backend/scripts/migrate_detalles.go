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
	if err := godotenv.Load("../.env"); err != nil {
		fmt.Println("No se cargo .env de root")
	}

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

	fmt.Println("--- Aplicando Migración de Esquema ---")

	// 1. Agregar columnas si no existen
	_, err = db.Exec(`
		ALTER TABLE documentos_detalles 
		ADD COLUMN IF NOT EXISTS subtotal NUMERIC(15,2) DEFAULT 0,
		ADD COLUMN IF NOT EXISTS total NUMERIC(15,2) DEFAULT 0;
	`)
	if err != nil {
		log.Fatalf("Error al agregar columnas: %v", err)
	}

	fmt.Println("Columnas agregadas (o ya existentes).")

	// 2. Actualizar registros existentes con cálculos básicos
	_, err = db.Exec(`
		UPDATE documentos_detalles 
		SET 
			subtotal = cantidad * precio_unitario,
			total = (cantidad * precio_unitario) + COALESCE(monto_iva_renglon, 0) - COALESCE(monto_descuento_renglon, 0)
		WHERE subtotal = 0 AND total = 0;
	`)
	if err != nil {
		log.Fatalf("Error al actualizar registros: %v", err)
	}

	fmt.Println("Registros actualizados exitosamente.")
}
