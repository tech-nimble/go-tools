package gin

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func QueryInt(ctx *gin.Context, key string) int {
	i, err := strconv.Atoi(ctx.Query(key))
	if err != nil {
		return 0
	}

	return i
}

func ParamInt(ctx *gin.Context, key string) int {
	i, err := strconv.Atoi(ctx.Param(key))
	if err != nil {
		return 0
	}

	return i
}
