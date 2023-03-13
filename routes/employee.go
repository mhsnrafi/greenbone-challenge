package routes

import (
	"github.com/gin-gonic/gin"
	"greenbone-task/controllers"
	"greenbone-task/middlewares"
)

func Employee(router *gin.RouterGroup) {
	auth := router.Group("/api")
	{

		auth.POST(
			"/employees/",
			middlewares.JWTMiddleware(),
			controllers.CreateEmployee,
		)
		auth.GET(
			"/employees/computers/:employee_abbrev",
			middlewares.JWTMiddleware(),
			controllers.GetEmployeeComputers,
		)
		auth.DELETE(
			"/employees/computers/:computer_id/:employee_abbrev",
			middlewares.JWTMiddleware(),
			controllers.DeleteEmployeeComputer,
		)
	}
}
