package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/tu-usuario/vinylnet/config"
)

type UsuarioResponse struct {
	ID             int    `json:"id"`
	NombreCompleto string `json:"nombre_completo"`
	Username       string `json:"username"`
	Rol            string `json:"rol"`
	Activo         bool   `json:"activo"`
	CreadoAt       string `json:"creado_at"`
}

type UserRequest struct {
	NombreCompleto string `json:"nombre_completo" binding:"required"`
	Username       string `json:"username" binding:"required"`
	Password       string `json:"password"`
	RolID          int    `json:"rol_id" binding:"required"`
}

// GET /usuarios
func GetUsuarios(c *gin.Context) {
	query := `
        SELECT u.id, u.nombre_completo, u.username, r.nombre, u.activo 
        FROM usuarios u
        JOIN roles r ON u.rol_id = r.id 
        ORDER BY u.id ASC`

	rows, err := config.DB.Query(query)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var usuarios []UsuarioResponse
	for rows.Next() {
		var u UsuarioResponse
		if err := rows.Scan(&u.ID, &u.NombreCompleto, &u.Username, &u.Rol, &u.Activo); err != nil {
			continue
		}
		usuarios = append(usuarios, u)
	}
	c.JSON(200, usuarios)
}

// POST /usuarios (Crear)
func CreateUsuario(c *gin.Context) {
	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Datos inválidos"})
		return
	}
	query := `INSERT INTO usuarios (nombre_completo, username, password, rol_id, activo) 
	          VALUES ($1, $2, $3, $4, TRUE) RETURNING id`
	var id int
	err := config.DB.QueryRow(query, req.NombreCompleto, req.Username, req.Password, req.RolID).Scan(&id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error al guardar"})
		return
	}
	c.JSON(201, gin.H{"message": "Creado", "id": id})
}

// DELETE /usuarios/:id
func DeleteUsuario(c *gin.Context) {
	id := c.Param("id")
	_, err := config.DB.Exec("DELETE FROM usuarios WHERE id = $1", id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error al eliminar"})
		return
	}
	c.JSON(200, gin.H{"message": "Eliminado"})
}
