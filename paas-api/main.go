package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"paas-api/handlers"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "PaaS API is running!"})
	})

	r.POST("/databases", handlers.CreateDatabase)

	log.Println("API listening on port 8080")
	r.Run(":8080")
}
