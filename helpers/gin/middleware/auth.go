package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tech-nimble/go-tools/helpers/gin/credentials"
)

const UnauthorizedErrCode = http.StatusUnauthorized

func AuthOnly() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cred := credentials.GetCredentialsFromCtx(ctx)
		if cred.UserID == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
