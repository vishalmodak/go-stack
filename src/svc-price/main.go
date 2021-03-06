package main

import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Global middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.GET("/price/item/:id", func(c *gin.Context) {
		itemID := c.Param("id")

		itemPrice := Price {
			ItemId: itemID,
			LastUpdate: time.Now(),
			ListPrice : ListPrice {
				MinPrice: 5.1,
				MaxPrice: 25.5,
				Price: 20.0,
				PriceType: "Regular",
			},
			OfferPrice: OfferPrice {
				MinPrice: 5.0,
				MaxPrice: 15.0,
				Price: 15.0,
				PriceType: "Regular",
				StartDate: time.Now(),
				EndDate: time.Now().Add(5 * time.Minute),
			},
		}

		c.JSON(http.StatusOK, itemPrice)
	})
	router.Run(":10002")
}
