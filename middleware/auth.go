package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Auth(ctx *gin.Context) {
	username := ctx.Request.Header[http.CanonicalHeaderKey("username")][0]
	session := sessions.Default(ctx)
	user := session.Get(username)
	fmt.Println(user)

	if user == nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx.Next()
}
