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

	// 2. Middleware CORS
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

	// 3. RUTAS DE LA API
	api := r.Group("/api")
	{
		// --- Módulo de Usuarios ---
		api.GET("/usuarios", controllers.GetUsuarios)
		api.POST("/usuarios", controllers.CreateUsuario)
		api.DELETE("/usuarios/:id", controllers.DeleteUsuario)

		// --- Módulo de Roles (NUEVO) ---
		api.GET("/roles", controllers.GetRoles)   // Para el select en el front
		api.POST("/roles", controllers.CreateRol) // Por si quieres crear roles nuevos
		api.PUT("/roles/:id", controllers.UpdateRol)
		api.DELETE("/roles/:id", controllers.DeleteRol) // Por si quieres editar descripciones

		// --- Autenticación ---
		api.POST("/login", controllers.Login)

		// --- Módulo de Facturas (VENTAS) ---
		api.GET("/invoices", controllers.GetInvoices)
		api.GET("/invoices/:id", controllers.GetInvoiceDetails)
		api.POST("/invoices", controllers.CreateInvoice)

		// RUTAS AUXILIARES
		api.GET("/clients", controllers.GetClients)
		api.GET("/products", controllers.GetProducts)
		api.GET("/rates", controllers.GetActiveRates)

		// CANONICAL SYSTEM (documentos)
		api.GET("/documentos", controllers.ListarDocumentos)
		api.GET("/documentos/buscar", controllers.BuscarDocumento)
		api.POST("/documentos/devolucion", controllers.CrearDevolucion)
		api.GET("/devoluciones/historial", controllers.ListarHistorialDevoluciones)
	}

	// 4. Encender servidor
	log.Println("Servidor VINYLNET corriendo en http://localhost:8080")
	r.Run(":8080")
}
