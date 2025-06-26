package sentry

import (
	"context"

	"github.com/getsentry/sentry-go"
	ginsentry "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

func GetHubFromContext(ctx context.Context) *sentry.Hub {
	if gCtx, ok := ctx.(*gin.Context); ok {
		return ginsentry.GetHubFromContext(gCtx)
	}

	return sentry.GetHubFromContext(ctx)
}
