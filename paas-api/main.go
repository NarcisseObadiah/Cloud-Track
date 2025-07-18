package main


import (
	"github.com/gin-gonic/gin"
	"log"
	"paas-api/handlers"
)
 
func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context){
		c.JSON(200, gin.H{"message":"Pass API is running!"})
	})

	r.POST("/databases", handlers.CreateDatabase)
	log.Println("API listening on port 8080")
	r.Run(":8080")
}
