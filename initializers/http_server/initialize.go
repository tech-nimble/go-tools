package http_server

import (
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/envy"
)

const (
	HTTPServerAddrEnv     = "HTTP_SERVER_ADDR"
	DefaultHTTPServerAddr = ":8080"
)

type HTTPServerConfig struct {
	HTTPServerAddr string
	Router         *gin.Engine
}

type Server interface {
	ListenAndServe() (err error)
}

func InitializeHTTPServerConfig(router *gin.Engine) *HTTPServerConfig {
	return &HTTPServerConfig{
		HTTPServerAddr: envy.Get(HTTPServerAddrEnv, DefaultHTTPServerAddr),
		Router:         router,
	}
}

func InitializeHTTPServer(cfg *HTTPServerConfig) (Server, error) {
	// create http server
	// Пакет endless позволяет завершить работу http сервер правильно, не обрывая запросы
	// Сперва мы обработаем поступившие запросы и только потом закроем соединение и завершим работу приложения
	srv := endless.NewServer(cfg.HTTPServerAddr, cfg.Router)

	return srv, nil
}
