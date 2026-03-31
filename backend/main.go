package main

import (
	"database/sql"
	"fmt"
	"log" // Mejor usar log para errores fatales
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// 1. Cargar .env
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: No se encontró archivo .env, usando variables de entorno del sistema")
	}

	// 2. Obtener variables
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")

	// Validar que las variables esenciales no estén vacías
	if user == "" || pass == "" || name == "" {
		log.Fatal("Error: Faltan variables de configuración en el .env")
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, name)

	// 3. Abrir conexión
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error en los parámetros de conexión:", err)
	}
	defer db.Close()

	// 4. Verificación REAL de conexión
	err = db.Ping()
	if err != nil {
		log.Fatal("No se pudo conectar a la base de datos de Vinylnet:", err)
	}

	fmt.Println("¡Conexión exitosa a PostgreSQL!")

	// ... (Tu código anterior de conexión aquí) ...

	r := gin.Default()

	// 1. MIDDLEWARE DE CORS (Permite que Vue en puerto 5173/3000 entre)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 2. ESTRUCTURA PARA EL LOGIN
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// 3. RUTA DE LOGIN
	r.POST("/Login", func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "JSON inválido"})
			return
		}

		var nombreCompleto, nombreRol string
		query := `
            SELECT u.nombre_completo, r.nombre 
            FROM usuarios u 
            JOIN roles r ON u.rol_id = r.id 
            WHERE u.username = $1 AND u.password = $2 AND u.activo = TRUE`

		err := db.QueryRow(query, req.Username, req.Password).Scan(&nombreCompleto, &nombreRol)

		if err != nil {
			c.JSON(401, gin.H{"message": "Credenciales incorrectas"})
			return
		}

		c.JSON(200, gin.H{
			"usuario": nombreCompleto,
			"rol":     nombreRol,
		})
	})

	// 4. INICIAR EN EL PUERTO 8080
	r.Run("0.0.0.0:8080") // Aquí continuaría tu servidor Gin...
}
