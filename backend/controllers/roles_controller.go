package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tu-usuario/vinylnet/config"
)

type RoleResponse struct {
	ID          int    `json:"id"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
}

type CreateRoleRequest struct {
	Nombre      string `json:"nombre" binding:"required"`
	Descripcion string `json:"descripcion"`
}

// 1. CONSULTAR ROLES (GET /roles)
func GetRoles(c *gin.Context) {
	query := `SELECT id, nombre, descripcion FROM roles ORDER BY id ASC`

	rows, err := config.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al consultar roles"})
		return
	}
	defer rows.Close()

	var roles []RoleResponse
	for rows.Next() {
		var r RoleResponse
		if err := rows.Scan(&r.ID, &r.Nombre, &r.Descripcion); err != nil {
			continue
		}
		roles = append(roles, r)
	}
	c.JSON(http.StatusOK, roles)
}

// 2. CREAR ROL (POST /roles)
func CreateRole(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nombre es obligatorio"})
		return
	}

	query := `INSERT INTO roles (nombre, descripcion) VALUES ($1, $2) RETURNING id`
	var lastInsertId int
	err := config.DB.QueryRow(query, req.Nombre, req.Descripcion).Scan(&lastInsertId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear rol"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Rol creado con éxito", "id": lastInsertId})
}

// 3. ELIMINAR ROL (DELETE /roles/:id)
func DeleteRole(c *gin.Context) {
	id := c.Param("id")

	// Verificar si el rol está en uso (opcional, pero la DB lo impedirá por claves foráneas)
	query := `DELETE FROM roles WHERE id = $1`
	result, err := config.DB.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al intentar eliminar rol"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Rol no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rol eliminado correctamente"})
}
