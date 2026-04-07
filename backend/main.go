package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tu-usuario/vinylnet/config"      // Tu conexión
	"github.com/tu-usuario/vinylnet/controllers" // Tus funciones
)

func main() {
	// 1. Cargar configuración
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: No se encontró .env")
	}
	config.ConnectDB()

	r := gin.Default()

	// 2. Middleware CORS (Cópialo igual que lo tenías)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 3. RUTAS DEL CRUD DE USUARIOS
	// Estas rutas llaman a las funciones que definiste en user_controller.go
	api := r.Group("/api")
	{
		api.GET("/usuarios", controllers.GetUsuarios)          // Consultar
		api.POST("/usuarios", controllers.CreateUsuario)       // Crear (Insert)
		api.DELETE("/usuarios/:id", controllers.DeleteUsuario) // Borrar (Delete)

		// El login que ya tenías
		api.POST("/login", controllers.Login)

		// RUTAS DE ROLES
		api.GET("/roles", controllers.GetRoles)          // Consultar
		api.POST("/roles", controllers.CreateRole)       // Crear
		api.DELETE("/roles/:id", controllers.DeleteRole) // Borrar

		// RUTAS DE FACTURAS (VENTAS)
		api.GET("/invoices", controllers.GetInvoices)
		api.GET("/invoices/:id", controllers.GetInvoiceDetails)
		api.POST("/invoices", controllers.CreateInvoice)

		// RUTAS AUXILIARES
		api.GET("/clients", controllers.GetClients)
		api.GET("/products", controllers.GetProducts)
		api.GET("/rates", controllers.GetActiveRates)
	}

	// 4. Encender servidor
	r.Run(":8080")
}
