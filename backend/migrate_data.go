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
	godotenv.Load(".env")
	user, pass, name, host, port := os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT")
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, name)
	db, err := sql.Open("postgres", connStr)
	if err != nil { log.Fatal(err) }
	defer db.Close()

	// 1. Nuclear Cleanup: Reiniciar tabla
	fmt.Println("Borrando tablas dependientes y tabla de productos...")
	_, err = db.Exec(`DROP TABLE IF EXISTS invoice_items CASCADE`)
	if err != nil { log.Printf("Aviso: No se pudo borrar invoice_items: %v", err) }
	
	_, err = db.Exec(`DROP TABLE IF EXISTS products CASCADE`)
	if err != nil { log.Fatal("Error fatal borrando productos:", err) }

	// 2. Recrear Tabla (Basada en la definición extraída y necesidades del usuario)
	fmt.Println("Recreando tabla products...")
	createTableSQL := `
	CREATE TABLE products (
		id BIGINT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		unit VARCHAR(50) DEFAULT 'CM',
		price NUMERIC(10,2) NOT NULL DEFAULT 0,
		stock NUMERIC(10,2) NOT NULL DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	_, err = db.Exec(createTableSQL)
	if err != nil { log.Fatal("Error recreando tabla:", err) }

	// 3. Insertar los 9 productos de la imagen
	fmt.Println("Insertando nuevos productos...")
	products := []struct {
		ID   int
		Name string
		Price float64
		Stock float64
	}{
		{1001, "6010 DE 61 CM (BLANCO BRILLANTE)", 2.22, 7.52},
		{1002, "6011 DE 61 CM (NEGRO BRILLANTE)", 2.22, 45.99},
		{1003, "6631 DE 61 CM (AMARILLO LIMON)", 2.22, 93.93},
		{1004, "6025 DE 61 CM (AMARILLO CLARO)", 2.22, 0.89},
		{1005, "6021 DE 61 CM (AMARILLO BANDER)", 2.22, 36.76},
		{1006, "6814 DE 61 CM (AMARILLO CATTER)", 2.22, 6.92},
		{1007, "6038 DE 61 CM (NARANJA)", 2.22, -0.27},
		{1008, "6037 DE 61 CM (FUCSIA BRILLANTE)", 2.22, 1.05},
		{1009, "6033 DE 61 CM (ROJO PASION)", 2.22, 1.22},
	}

	for _, p := range products {
		_, err = db.Exec(`INSERT INTO products (id, name, price, stock) VALUES ($1, $2, $3, $4)`, 
			p.ID, p.Name, p.Price, p.Stock)
		if err != nil {
			log.Printf("Error insertando producto %d: %v", p.ID, err)
		}
	}

	fmt.Println("¡Migración y carga de datos completada!")
}
