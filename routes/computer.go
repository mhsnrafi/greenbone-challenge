package routes

import (
	"github.com/gin-gonic/gin"
	"greenbone-task/controllers"
	"greenbone-task/middlewares"
)

func Computer(router *gin.RouterGroup) {
	auth := router.Group("/")
	{
		auth.POST(
			"/computers",
			middlewares.JWTMiddleware(),
			controllers.CreateComputer,
		)
		auth.GET(
			"/computers",
			middlewares.JWTMiddleware(),
			controllers.GetAllComputers,
		)
		auth.GET(
			"/computers/:computer_id",
			middlewares.JWTMiddleware(),
			controllers.GetComputerByID,
		)
		auth.PUT(
			"/computers/:computer_id/:employee_abbrev",
			middlewares.JWTMiddleware(),
			controllers.UpdateComputer,
		)
		auth.DELETE(
			"/computers/:computer_id",
			middlewares.JWTMiddleware(),
			controllers.DeleteComputer,
		)
	}
}
