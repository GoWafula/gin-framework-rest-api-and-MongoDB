package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/kevin/webgin/controllers"
)

func UserRoutes(router *gin.Engine) {
	router.GET("/albums", controllers.GetAlbums)
	router.POST("/albums", controllers.PostAlbums)
	router.GET("/albums/:id", controllers.GetAlbumByID)
	router.PUT("/albums/:id", controllers.UpdateAlbum)
	router.DELETE("/albums/:id", controllers.DeleteAlbum)
}
