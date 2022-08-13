package route

import "github.com/gin-gonic/gin"

func PublicRoutes(g *gin.RouterGroup) {
	g.POST("/signup", handler.SignUp)
	g.GET("/login", handler.Login)
}

func PrivateRoutes(g *gin.RouterGroup) {

}
