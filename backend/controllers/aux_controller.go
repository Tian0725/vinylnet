package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tu-usuario/vinylnet/config"
)

type Client struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	RIF  string `json:"rif"`
}

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock float64 `json:"stock"`
}

func GetClients(c *gin.Context) {
	rows, err := config.DB.Query("SELECT id, name, rif FROM clients WHERE is_active = true")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var clis []Client
	for rows.Next() {
		var cli Client
		rows.Scan(&cli.ID, &cli.Name, &cli.RIF)
		clis = append(clis, cli)
	}
	c.JSON(http.StatusOK, clis)
}

func GetProducts(c *gin.Context) {
	rows, err := config.DB.Query("SELECT id, name, price, stock FROM products")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var prods []Product
	for rows.Next() {
		var prod Product
		rows.Scan(&prod.ID, &prod.Name, &prod.Price, &prod.Stock)
		prods = append(prods, prod)
	}
	c.JSON(http.StatusOK, prods)
}
