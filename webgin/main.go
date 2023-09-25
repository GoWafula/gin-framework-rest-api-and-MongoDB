package main

import (
	"github.com/gin-gonic/gin"
	routes "github.com/kevin/webgin/routes"
)

func main() {
	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)

	router.Run("localhost:8080")

}
