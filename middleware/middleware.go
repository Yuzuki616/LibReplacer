package middleware

import (
	"github.com/Yuzuki616/LibReplacer/handler"
	"github.com/Yuzuki616/LibReplacer/library"
	"github.com/patrickmn/go-cache"
	"time"
)

type Middleware struct {
	library *library.Library
	handler *handler.Handler
	cache   *cache.Cache
}

func New(l *library.Library, h *handler.Handler) *Middleware {
	return &Middleware{
		library: l,
		handler: h,
		cache:   cache.New(time.Minute*10, time.Hour),
	}
}
