package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB // Mayúscula para que otros paquetes la vean

func ConnectDB() {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, name)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error en parámetros:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("No se pudo conectar a la DB:", err)
	}

	fmt.Println("¡Conexión exitosa a PostgreSQL!")

	// Inicializar tablas
	InitTables()
}

func InitTables() {
	// Crear tabla roles
	rolesTable := `
	CREATE TABLE IF NOT EXISTS roles (
		id SERIAL PRIMARY KEY,
		nombre VARCHAR(50) NOT NULL UNIQUE,
		descripcion TEXT
	);`

	// Crear tabla usuarios
	usuariosTable := `
	CREATE TABLE IF NOT EXISTS usuarios (
		id SERIAL PRIMARY KEY,
		nombre_completo VARCHAR(100) NOT NULL,
		username VARCHAR(50) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		rol_id INTEGER REFERENCES roles(id) ON DELETE SET NULL,
		activo BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// Ejecutar creación de tablas
	if _, err := DB.Exec(rolesTable); err != nil {
		log.Fatal("Error creando tabla roles:", err)
	}

	if _, err := DB.Exec(usuariosTable); err != nil {
		log.Fatal("Error creando tabla usuarios:", err)
	}

	// Insertar roles básicos
	insertRoles := `
	INSERT INTO roles (nombre, descripcion) VALUES
	('Administrador', 'Acceso completo al sistema'),
	('Usuario', 'Acceso limitado')
	ON CONFLICT (nombre) DO NOTHING;`

	if _, err := DB.Exec(insertRoles); err != nil {
		log.Println("Aviso: Error insertando roles básicos:", err)
	}

	fmt.Println("Tablas inicializadas correctamente.")
}
