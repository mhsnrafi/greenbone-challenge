package middlewares

import (
	"github.com/gin-gonic/gin"
	"greenbone-task/models"
	db "greenbone-task/models/db"
	"greenbone-task/services"
	"net/http"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Bearer-Token")
		tokenModel, err := services.VerifyToken(token, db.TokenTypeAccess)
		if err != nil {
			models.SendErrorResponse(c, http.StatusUnauthorized, err.Error())
			return
		}

		c.Set("userIdHex", tokenModel.ID)
		c.Set("userId", tokenModel.ID)

		c.Next()
	}
}
