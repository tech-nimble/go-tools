package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tech-nimble/go-tools/helpers/errors"
	"github.com/tech-nimble/go-tools/helpers/gin/render"
)

var errContentType = errors.DomainErrWithDefaultCode("Невалидный Content-Type")

var imageContentTypes = []string{"image/png", "image/jpg", "image/jpeg"}
var allowedContentTypes = append([]string{"application/json", "application/vnd.api+json"}, imageContentTypes...)
var allowedContentTypesV2 = append([]string{"application/x-www-form-urlencoded", "multipart/form-data"}, imageContentTypes...)

func ContentType() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		contentType := ctx.GetHeader("Content-Type")
		if contentType != "" && !contains(allowedContentTypes, contentType) {
			jsonApiError := render.NewJSONAPIErrors().AddError(errContentType)
			ctx.Render(
				jsonApiError.GetHttpStatus(),
				jsonApiError,
			)

			ctx.Abort()
			return
		}
	}
}

func ContentTypeV2() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		contentType := ctx.GetHeader("Content-Type")
		if contentType != "" && !contains(allowedContentTypesV2, contentType) {
			ctx.Render(
				http.StatusOK,
				render.NewResponseWithError(errContentType),
			)

			ctx.Abort()
			return
		}
	}
}

func contains(s []string, v string) bool {
	for _, e := range s {
		sv := strings.Split(strings.Trim(v, " ;"), ";")[0]
		if e == sv {
			return true
		}
	}

	return false
}
