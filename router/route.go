package router

import (
	"LibReplacer/middleware"
)

func (r *Router) loadRoute() {
	m := middleware.New(r.library, r.handler)
	mod := r.engine.Group("/emby/", m.LibraryCheck)
	cacheQueryMod := mod.Group("", m.CacheResponse)
	cacheQueryMod.GET("Users/:userid/Items", r.handler.Items)
	cacheQueryMod.GET("Users/:userid/Items/Latest", r.handler.Latest)
	mod.GET("Users/:userid/Items/Resume", r.handler.Resume)
	cacheQueryMod.GET("Genres", r.handler.Genres)
	r.engine.NoRoute(r.handler.ReverseProxy)
}
