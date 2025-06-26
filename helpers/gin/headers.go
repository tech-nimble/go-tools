package gin

import (
	"github.com/gin-gonic/gin"
)

const (
	XUserAgentHeader         = "X-User-Agent"
	UserAgentHeader          = "User-Agent"
	UserAgentForwardedHeader = "X-Forwarded-Via"
	OriginHeader             = "Origin"
)

func GetOrigin(ctx *gin.Context) string {
	return ctx.GetHeader(OriginHeader)
}

func GetUserAgent(ctx *gin.Context) string {
	userAgent := ctx.GetHeader(XUserAgentHeader)

	if userAgent == "" {
		userAgent = ctx.GetHeader(UserAgentForwardedHeader)
	}

	if userAgent == "" {
		userAgent = ctx.GetHeader(UserAgentHeader)
	}

	return userAgent
}
