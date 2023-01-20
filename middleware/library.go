package middleware

import "github.com/gin-gonic/gin"

func (m *Middleware) LibraryCheck(c *gin.Context) {
	libId := c.Query("ParentId")
	if !m.library.IsExist(libId) {
		m.handler.ReverseProxy(c)
		c.Abort()
	}
}
