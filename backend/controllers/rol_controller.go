package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/tu-usuario/vinylnet/config"
)

// Estructura para las respuestas (Coincide con tu DB)
type RolResponse struct {
	ID          int    `json:"id"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
}

// GET /roles - Listar todos los roles para el select de Vue
func GetRoles(c *gin.Context) {
	query := `SELECT id, nombre, descripcion FROM roles ORDER BY id ASC`

	rows, err := config.DB.Query(query)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error al consultar roles"})
		return
	}
	defer rows.Close()

	roles := []RolResponse{}

	for rows.Next() {
		var r RolResponse
		if err := rows.Scan(&r.ID, &r.Nombre, &r.Descripcion); err != nil {
			continue
		}
		roles = append(roles, r)
	}

	c.JSON(200, roles)
}

// POST /roles - Crear un nuevo rol
func CreateRol(c *gin.Context) {
	var req struct {
		Nombre      string `json:"nombre" binding:"required"`
		Descripcion string `json:"descripcion"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Datos inválidos o nombre faltante"})
		return
	}

	query := `INSERT INTO roles (nombre, descripcion) VALUES ($1, $2) RETURNING id`
	var id int
	err := config.DB.QueryRow(query, req.Nombre, req.Descripcion).Scan(&id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error al crear el rol"})
		return
	}

	c.JSON(201, gin.H{"message": "Rol creado", "id": id})
}

// PUT /roles/:id - Actualizar un rol (ESTA ES LA QUE TE FALTABA)
func UpdateRol(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Nombre      string `json:"nombre" binding:"required"`
		Descripcion string `json:"descripcion"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Datos inválidos"})
		return
	}

	query := `UPDATE roles SET nombre = $1, descripcion = $2 WHERE id = $3`
	_, err := config.DB.Exec(query, req.Nombre, req.Descripcion, id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error al actualizar el rol"})
		return
	}

	c.JSON(200, gin.H{"message": "Rol actualizado correctamente"})
}

// DELETE /roles/:id
func DeleteRol(c *gin.Context) {
	id := c.Param("id")

	// Ejecutamos la eliminación en la base de datos
	result, err := config.DB.Exec("DELETE FROM roles WHERE id = $1", id)
	if err != nil {
		c.JSON(500, gin.H{"error": "No se puede eliminar el rol (puede estar en uso)"})
		return
	}

	// Verificamos si realmente se eliminó algo
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(404, gin.H{"error": "El rol no existe"})
		return
	}

	c.JSON(200, gin.H{"message": "Rol eliminado correctamente"})
}
