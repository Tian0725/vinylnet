package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tu-usuario/vinylnet/config"
)

type ExchangeRate struct {
	ID           int     `json:"id"`
	CurrencyCode string  `json:"currency_code"`
	Rate         float64 `json:"rate"`
	Source       string  `json:"source"`
}

func GetActiveRates(c *gin.Context) {
	// Obtener las tasas de USD/VES y USD/COP (o las que estén activas)
	rows, err := config.DB.Query(`
		SELECT id, currency_code, rate, source 
		FROM exchange_rates 
		WHERE is_active = true 
		ORDER BY id DESC`)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var rates []ExchangeRate
	for rows.Next() {
		var r ExchangeRate
		rows.Scan(&r.ID, &r.CurrencyCode, &r.Rate, &r.Source)
		rates = append(rates, r)
	}

	c.JSON(http.StatusOK, rates)
}
