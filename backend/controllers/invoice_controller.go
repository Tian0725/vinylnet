package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tu-usuario/vinylnet/config"
)

type Invoice struct {
	ID              int             `json:"id"`
	ClientID        int             `json:"client_id"`
	UserID          int             `json:"user_id"`
	Number          string          `json:"number"`
	IssueDate       string          `json:"issue_date"`
	DueDate         string          `json:"due_date"`
	Status          string          `json:"status"`
	Total           float64         `json:"total"`
	Subtotal        float64         `json:"subtotal"`
	TaxAmount       float64         `json:"tax_amount"`
	Currency        string          `json:"currency"`
	Notes           string          `json:"notes"`
	PrintFormat     string          `json:"print_format"`
	IGTFAmount      float64         `json:"igtf_amount"`
	IGTFPercent     float64         `json:"igtf_percent"`
	TaxPercent      float64         `json:"tax_percent"`
	ClientName      sql.NullString  `json:"client_name"` // To join with clients
}

type CreateInvoiceRequest struct {
	ClientID    int           `json:"client_id"`
	UserID      int           `json:"user_id"`
	Currency    string        `json:"currency"`
	ApplyIVA    bool          `json:"apply_iva"`
	ApplyIGTF   bool          `json:"apply_igtf"`
	Notes       string        `json:"notes"`
	Items       []ItemRequest `json:"items"`
}

type ItemRequest struct {
	ProductID int     `json:"product_id"`
	Quantity  float64 `json:"quantity"`
	Price     float64 `json:"price"`
}

// GetInvoices retrieves all invoices or filtered by status
func GetInvoices(c *gin.Context) {
	status := c.Query("status")
	printFormat := c.Query("print_format")

	query := `
		SELECT i.id, i.client_id, i.user_id, i.number, i.issue_date, i.due_date, i.status, 
		       i.total, i.subtotal, i.tax_amount, i.currency, i.notes, i.print_format,
		       i.igtf_amount, i.igtf_percent, i.tax_percent,
		       cl.name as client_name
		FROM invoices i
		LEFT JOIN clients cl ON i.client_id = cl.id
		WHERE 1=1
	`
	args := []interface{}{}

	if status != "" {
		query += " AND i.status = $1"
		args = append(args, status)
	}
	
	if printFormat != "" {
		if status != "" {
			query += " AND i.print_format = $2"
		} else {
			query += " AND i.print_format = $1"
		}
		args = append(args, printFormat)
	}

	query += " ORDER BY i.issue_date DESC"

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var invoices []Invoice
	for rows.Next() {
		var inv Invoice
		if err := rows.Scan(
			&inv.ID, &inv.ClientID, &inv.UserID, &inv.Number, &inv.IssueDate, &inv.DueDate, &inv.Status,
			&inv.Total, &inv.Subtotal, &inv.TaxAmount, &inv.Currency, &inv.Notes, &inv.PrintFormat,
			&inv.IGTFAmount, &inv.IGTFPercent, &inv.TaxPercent,
			&inv.ClientName,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		invoices = append(invoices, inv)
	}

	c.JSON(http.StatusOK, invoices)
}

// GetInvoiceDetails retrieves a single invoice with its items
func GetInvoiceDetails(c *gin.Context) {
	id := c.Param("id")

	// 1. Get header
	var inv Invoice
	err := config.DB.QueryRow(`
		SELECT i.id, i.client_id, i.number, i.status, i.total, cl.name
		FROM invoices i
		LEFT JOIN clients cl ON i.client_id = cl.id
		WHERE i.id = $1`, id).Scan(&inv.ID, &inv.ClientID, &inv.Number, &inv.Status, &inv.Total, &inv.ClientName)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Factura no encontrada"})
		return
	}

	// 2. Get items (simplified for now)
	rows, err := config.DB.Query("SELECT id, product_id, quantity, price, subtotal FROM invoice_items WHERE invoice_id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error cargando items"})
		return
	}
	defer rows.Close()

	type Item struct {
		ID        int     `json:"id"`
		ProductID int     `json:"product_id"`
		Quantity  float64 `json:"quantity"`
		Price     float64 `json:"price"`
		Subtotal  float64 `json:"subtotal"`
	}
	var items []Item
	for rows.Next() {
		var item Item
		rows.Scan(&item.ID, &item.ProductID, &item.Quantity, &item.Price, &item.Subtotal)
		items = append(items, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"invoice": inv,
		"items":   items,
	})
}

// CreateInvoice handles the creation of a new sale with transaction
func CreateInvoice(c *gin.Context) {
	var req CreateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Iniciar Transacción
	tx, err := config.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error de base de datos"})
		return
	}
	defer tx.Rollback()

	// 2. Generar Correlativo automático (INV-0001...)
	var lastNumber string
	var nextNumber string
	err = tx.QueryRow("SELECT number FROM invoices ORDER BY id DESC LIMIT 1 FOR UPDATE").Scan(&lastNumber)
	if err != nil {
		nextNumber = "INV-0001"
	} else {
		var currentID int
		fmt.Sscanf(lastNumber, "INV-%04d", &currentID)
		nextNumber = fmt.Sprintf("INV-%04d", currentID+1)
	}

	// 3. Calcular Totales
	var subtotal float64
	for _, item := range req.Items {
		subtotal += item.Quantity * item.Price
	}

	var taxAmount float64
	var taxPercent float64
	if req.ApplyIVA {
		taxPercent = 16.0
		taxAmount = subtotal * 0.16
	}

	var igtfAmount float64
	var igtfPercent float64
	if req.ApplyIGTF && req.Currency != "VES" {
		igtfPercent = 3.0
		igtfAmount = (subtotal + taxAmount) * 0.03
	}

	total := subtotal + taxAmount + igtfAmount

	// 4. Insertar Cabecera
	var invoiceID int
	err = tx.QueryRow(`
		INSERT INTO invoices (
			client_id, user_id, number, issue_date, status, total, subtotal, 
			tax_amount, tax_percent, igtf_amount, igtf_percent, currency, notes
		) VALUES ($1, $2, $3, NOW(), 'FINAL', $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id`,
		req.ClientID, req.UserID, nextNumber, total, subtotal, 
		taxAmount, taxPercent, igtfAmount, igtfPercent, req.Currency, req.Notes,
	).Scan(&invoiceID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar cabecera: " + err.Error()})
		return
	}

	// 5. Insertar Items
	for _, item := range req.Items {
		itemSubtotal := item.Quantity * item.Price
		_, err = tx.Exec(`
			INSERT INTO invoice_items (invoice_id, product_id, quantity, price, subtotal)
			VALUES ($1, $2, $3, $4, $5)`,
			invoiceID, item.ProductID, item.Quantity, item.Price, itemSubtotal,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar item: " + err.Error()})
			return
		}
	}

	// 6. Commit
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al confirmar transacción"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Venta creada con éxito",
		"invoice_id": invoiceID,
		"number": nextNumber,
	})
}
