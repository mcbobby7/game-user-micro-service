package routes

import (
	controller "github.com/Gameware/controllers"
	"github.com/Gameware/middleware"
	"github.com/gin-gonic/gin"
)

func TournamentRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.POST("/tournament/save", controller.SaveTournament())
	incomingRoutes.GET("/tournaments", controller.GetTournaments())
	incomingRoutes.GET("/tournament/:id", controller.GetTournament())
}