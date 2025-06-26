package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	ginHelpers "github.com/tech-nimble/go-tools/helpers/gin"
	"github.com/tech-nimble/go-tools/helpers/gin/credentials"
	"github.com/tech-nimble/go-tools/helpers/http"
)

const (
	MerchantHeader      = "X-Merchant"
	UserHeader          = "X-User"
	ClientIDHeader      = "X-Client-Id"
	IPHeader            = "X-User-Ip"
	FingerprintHeader   = "X-Fingerprint"
	SessionIDHeader     = "X-Session-Id"
	AuthorizationHeader = "Authorization"
)

func ParseHeaders(ctx *gin.Context) {
	userAgent := ginHelpers.GetUserAgent(ctx)

	cred := credentials.NewCredentials(
		http.GetIntHeader(ctx, MerchantHeader),
		http.GetIntHeader(ctx, UserHeader),
		http.GetIntHeader(ctx, ClientIDHeader),
		ctx.GetHeader(IPHeader),
		ctx.GetHeader(FingerprintHeader),
		userAgent,
		ctx.GetHeader(SessionIDHeader),
		strings.ReplaceAll(ctx.GetHeader(AuthorizationHeader), "Bearer ", ""),
	)

	ctx.Set("credentials", cred)
}
