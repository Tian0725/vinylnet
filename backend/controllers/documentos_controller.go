package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tu-usuario/vinylnet/config"
)

type Documento struct {
	ID                int             `json:"id"`
	NumeroDocumento   string          `json:"numero_documento"`
	TipoDocumentoID   int             `json:"tipo_documento_id"`
	ClienteID         int             `json:"cliente_id"`
	SucursalID        int             `json:"sucursal_id"`
	UsuarioID         int             `json:"usuario_id"`
	FechaEmision      string          `json:"fecha_emision"`
	TotalNeto         float64         `json:"total_neto"`
	Subtotal          float64         `json:"subtotal"`
	MontoIVA          float64         `json:"monto_iva"`
	MontoIGTF         float64         `json:"monto_igtf"`
	NumeroReferencia  sql.NullString  `json:"numero_referencia"`
	DocReferenciaID   sql.NullInt64   `json:"documento_referencia_id"`
	ClienteNombre     sql.NullString  `json:"cliente_nombre"`
	TipoDescripcion   sql.NullString  `json:"tipo_descripcion"`
}

type DetalleDocumento struct {
	ID             int     `json:"id"`
	ProductoID     int     `json:"producto_id"`
	Cantidad       float64 `json:"cantidad"`
	PrecioUnitario float64 `json:"precio_unitario"`
	Subtotal       float64 `json:"subtotal"`
	Total          float64 `json:"total"`
	ProductoNombre string  `json:"producto_nombre"`
}

// BuscarDocumento busca una factura por número para devoluciones
func BuscarDocumento(c *gin.Context) {
	numero := c.Query("numero")
	if numero == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Número de documento requerido"})
		return
	}

	var d Documento
	query := `
		SELECT d.id, d.numero_documento, d.tipo_documento_id, d.cliente_id, d.fecha_emision, d.total_neto, 
		       cl.nombre_razon_social as cliente_nombre, td.descripcion as tipo_descripcion
		FROM documentos d
		LEFT JOIN clientes cl ON d.cliente_id = cl.id
		LEFT JOIN tipos_documentos td ON d.tipo_documento_id = td.id
		WHERE d.numero_documento = $1 AND d.tipo_documento_id IN (1, 5)
		LIMIT 1
	`
	err := config.DB.QueryRow(query, numero).Scan(
		&d.ID, &d.NumeroDocumento, &d.TipoDocumentoID, &d.ClienteID, &d.FechaEmision, &d.TotalNeto,
		&d.ClienteNombre, &d.TipoDescripcion,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Factura no encontrada"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Obtener detalles
	rows, err := config.DB.Query(`
		SELECT det.id, det.producto_id, det.cantidad, det.precio_unitario, det.subtotal, det.total, p.descripcion
		FROM documentos_detalles det
		JOIN productos p ON det.producto_id = p.id
		WHERE det.documento_id = $1`, d.ID)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al cargar detalles"})
		return
	}
	defer rows.Close()

	var detalles []DetalleDocumento
	for rows.Next() {
		var det DetalleDocumento
		rows.Scan(&det.ID, &det.ProductoID, &det.Cantidad, &det.PrecioUnitario, &det.Subtotal, &det.Total, &det.ProductoNombre)
		detalles = append(detalles, det)
	}

	c.JSON(http.StatusOK, gin.H{
		"documento": d,
		"detalles":  detalles,
	})
}

// CrearDevolucion genera una nota de crédito vinculada
func CrearDevolucion(c *gin.Context) {
	var req struct {
		OriginalID int     `json:"original_id"`
		Motivo     string  `json:"motivo"`
		Detalles   []struct {
			ProductoID int     `json:"producto_id"`
			Cantidad   float64 `json:"cantidad"`
		} `json:"detalles"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := config.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error de base de datos"})
		return
	}
	defer tx.Rollback()

	// 1. Obtener datos de la factura original
	var orig Documento
	err = tx.QueryRow(`SELECT numero_documento, tipo_documento_id, cliente_id, sucursal_id, usuario_id FROM documentos WHERE id = $1`, req.OriginalID).Scan(
		&orig.NumeroDocumento, &orig.TipoDocumentoID, &orig.ClienteID, &orig.SucursalID, &orig.UsuarioID,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Factura original no encontrada"})
		return
	}

	// 2. Determinar tipo de devolución (6=DEV para 1=FAC, 7=DEVN para 5=FAN)
	tipoDevID := 6
	prefix := "DEV"
	if orig.TipoDocumentoID == 5 {
		tipoDevID = 7
		prefix = "DEVN"
	}

	// 3. Generar correlativo de devolución
	var lastNum string
	var nextNum string
	tx.QueryRow("SELECT numero_documento FROM documentos WHERE tipo_documento_id = $1 ORDER BY id DESC LIMIT 1", tipoDevID).Scan(&lastNum)
	if lastNum == "" {
		nextNum = fmt.Sprintf("%s-0001", prefix)
	} else {
		var currentID int
		fmt.Sscanf(lastNum, prefix+"-%04d", &currentID)
		nextNum = fmt.Sprintf("%s-%04d", prefix, currentID+1)
	}

	// 4. Insertar Cabecera de Devolución (Valores negativos o marcados)
	// Para simplicidad, guardamos los montos como positivos pero con el tipo_documento_id correcto
	var devID int
	err = tx.QueryRow(`
		INSERT INTO documentos (
			numero_documento, tipo_documento_id, cliente_id, sucursal_id, usuario_id, 
			numero_referencia, documento_referencia_id, fecha_emision, total_neto
		) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), 0) RETURNING id`,
		nextNum, tipoDevID, orig.ClienteID, orig.SucursalID, orig.UsuarioID, orig.NumeroDocumento, req.OriginalID,
	).Scan(&devID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear devolución: " + err.Error()})
		return
	}

	// 5. Insertar detalles y calcular total
	var totalDev float64
	for _, it := range req.Detalles {
		var precio float64
		tx.QueryRow("SELECT precio_unitario FROM documentos_detalles WHERE documento_id = $1 AND producto_id = $2", req.OriginalID, it.ProductoID).Scan(&precio)
		
		subtotal := it.Cantidad * precio
		totalDev += subtotal

		_, err = tx.Exec(`
			INSERT INTO documentos_detalles (documento_id, producto_id, cantidad, precio_unitario, subtotal, total)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			devID, it.ProductoID, it.Cantidad, precio, subtotal, subtotal,
		)
	}

	// 6. Actualizar total en cabecera
	tx.Exec("UPDATE documentos SET total_neto = $1, subtotal = $1 WHERE id = $2", totalDev, devID)

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al confirmar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Devolución creada con éxito", "numero": nextNum})
}

// ListarHistorialDevoluciones devuelve las devoluciones realizadas
func ListarHistorialDevoluciones(c *gin.Context) {
	rows, err := config.DB.Query(`
		SELECT d.id, d.numero_documento, d.numero_referencia, d.fecha_emision, d.total_neto, td.descripcion as tipo_descripcion
		FROM documentos d
		LEFT JOIN tipos_documentos td ON d.tipo_documento_id = td.id
		WHERE d.tipo_documento_id IN (6, 7)
		ORDER BY d.fecha_emision DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var docs []Documento
	for rows.Next() {
		var d Documento
		rows.Scan(&d.ID, &d.NumeroDocumento, &d.NumeroReferencia, &d.FechaEmision, &d.TotalNeto, &d.TipoDescripcion)
		docs = append(docs, d)
	}

	c.JSON(http.StatusOK, docs)
}

// ListarDocumentos devuelve todos los documentos (FAC, FAN)
func ListarDocumentos(c *gin.Context) {
	tipo := c.Query("tipo") // Opcional: filtrar por tipo
	query := `
		SELECT d.id, d.numero_documento, d.fecha_emision, d.total_neto, td.descripcion as tipo_descripcion, cl.nombre_razon_social as cliente_nombre
		FROM documentos d
		LEFT JOIN tipos_documentos td ON d.tipo_documento_id = td.id
		LEFT JOIN clientes cl ON d.cliente_id = cl.id
		WHERE d.tipo_documento_id IN (1, 5)
	`
	args := []interface{}{}
	if tipo != "" {
		query += " AND d.tipo_documento_id = $1"
		args = append(args, tipo)
	}
	query += " ORDER BY d.fecha_emision DESC"

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var docs []Documento
	for rows.Next() {
		var d Documento
		rows.Scan(&d.ID, &d.NumeroDocumento, &d.FechaEmision, &d.TotalNeto, &d.TipoDescripcion, &d.ClienteNombre)
		docs = append(docs, d)
	}

	c.JSON(http.StatusOK, docs)
}
