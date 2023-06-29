package middleware

import (
	"a21hc3NpZ25tZW50/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Auth() gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context) {
		cookie, err := ctx.Cookie("session_token")
		if err != nil {
			contentType := ctx.Request.Header.Get("Content-Type")
			if contentType == "application/json"{
				ctx.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Cookie not found"})
				return
			} else {
				ctx.Redirect(http.StatusSeeOther, "/user/login")
			}
			ctx.Abort()
			return
		}

		claims := &model.Claims{}
		
		token, err := jwt.ParseWithClaims(cookie, claims, func(token *jwt.Token) (interface{}, error) {
			return model.JwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				ctx.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Invalid signature"})
				return
			}
			ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "token error"})
			return
		}

		if !token.Valid {
			ctx.JSON(http.StatusUnauthorized, "token tidak valid")
			return
		}

		ctx.Set("email", claims.Email)
		// ctx.JSON(http.StatusCreated, model.SuccessResponse{Message: "auth valid"})
		ctx.Next()
	})
}
