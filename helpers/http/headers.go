package http

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetHeadersAsMap(h http.Header) map[string]string {
	headers := map[string]string{}

	for k, v := range h {
		headers[k] = strings.Join(v, ", ")
	}

	return headers
}

func GetIntHeader(ctx *gin.Context, key string) int {
	hStr := ctx.GetHeader(key)
	if hStr == "" {
		return 0
	}

	h, err := strconv.Atoi(hStr)
	if err != nil {
		return 0
	}

	return h
}
