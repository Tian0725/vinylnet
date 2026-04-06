package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tu-usuario/vinylnet/config" // Verifica que este sea el nombre en tu go.mod
)

// --- ESTRUCTURAS ---

// UsuarioResponse: Lo que enviamos al Frontend (GET)
type UsuarioResponse struct {
	ID             int    `json:"id"`
	NombreCompleto string `json:"nombre_completo"`
	Username       string `json:"username"`
	Rol            string `json:"rol"`
	Activo         bool   `json:"activo"`
	CreadoAt       string `json:"creado_at"`
}

// CreateUsuarioRequest: Lo que recibimos del Frontend (POST)
type CreateUsuarioRequest struct {
	NombreCompleto string `json:"nombre_completo" binding:"required"`
	Username       string `json:"username" binding:"required"`
	Password       string `json:"password" binding:"required"`
	RolID          int    `json:"rol_id" binding:"required"`
}

// --- FUNCIONES DEL CRUD ---

// 1. CONSULTAR (GET /usuarios)
func GetUsuarios(c *gin.Context) {
	query := `
        SELECT u.id, u.nombre_completo, u.username, r.nombre, u.activo, u.creado_at
        FROM usuarios u
        JOIN roles r ON u.rol_id = r.id 
        ORDER BY u.id ASC`

	rows, err := config.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al consultar usuarios"})
		return
	}
	defer rows.Close()

	var usuarios []UsuarioResponse
	for rows.Next() {
		var u UsuarioResponse
		if err := rows.Scan(&u.ID, &u.NombreCompleto, &u.Username, &u.Rol, &u.Activo, &u.CreadoAt); err != nil {
			continue
		}
		usuarios = append(usuarios, u)
	}
	c.JSON(http.StatusOK, usuarios)
}

// 2. CREAR (POST /usuarios)
func CreateUsuario(c *gin.Context) {
	var req CreateUsuarioRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos incompletos o inválidos"})
		return
	}

	query := `
        INSERT INTO usuarios (nombre_completo, username, password, rol_id, activo) 
        VALUES ($1, $2, $3, $4, TRUE) RETURNING id`

	var lastInsertId int
	err := config.DB.QueryRow(query, req.NombreCompleto, req.Username, req.Password, req.RolID).Scan(&lastInsertId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al insertar usuario en la DB"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Usuario creado con éxito",
		"id":      lastInsertId,
	})
}

// 3. BORRAR (DELETE /usuarios/:id)
func DeleteUsuario(c *gin.Context) {
	id := c.Param("id")

	query := `DELETE FROM usuarios WHERE id = $1`

	result, err := config.DB.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al intentar eliminar"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Usuario no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuario eliminado correctamente"})
}
