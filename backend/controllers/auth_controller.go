package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/tu-usuario/vinylnet/config" // Ruta completa según tu go.mod
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
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

	// Usamos config.DB que definimos en db.go
	err := config.DB.QueryRow(query, req.Username, req.Password).Scan(&nombreCompleto, &nombreRol)

	if err != nil {
		c.JSON(401, gin.H{"message": "Credenciales incorrectas"})
		return
	}

	c.JSON(200, gin.H{
		"usuario": nombreCompleto,
		"rol":     nombreRol,
	})
}
