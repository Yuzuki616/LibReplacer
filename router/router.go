package router

import (
	"fmt"
	"github.com/Yuzuki616/LibReplacer/conf"
	"github.com/Yuzuki616/LibReplacer/handler"
	"github.com/Yuzuki616/LibReplacer/library"
	"github.com/gin-gonic/gin"
)

type Router struct {
	handler *handler.Handler
	library *library.Library
	engine  *gin.Engine
	config  *conf.Conf
}

func New(c *conf.Conf, l *library.Library) *Router {
	gin.SetMode(gin.ReleaseMode)
	return &Router{
		handler: handler.New(c, l),
		engine:  gin.New(),
		library: l,
		config:  c,
	}
}

func (r *Router) Start() error {
	r.engine.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			return fmt.Sprintf("| %s | %s [%s:%d] %s \n",
				param.ClientIP,
				param.TimeStamp.Format("2006-01-02 15:04:05"),
				param.Method,
				param.StatusCode,
				param.Path)
		},
	}), gin.Recovery())
	r.loadRoute()
	if r.config.EnableSsl {
		return r.engine.RunTLS(r.config.Addr, r.config.CertConfig.CertFile, r.config.CertConfig.KeyFile)
	}
	return r.engine.Run(r.config.Addr)
}
