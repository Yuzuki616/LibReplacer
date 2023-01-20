package handler

import (
	"github.com/Yuzuki616/LibReplacer/conf"
	"github.com/Yuzuki616/LibReplacer/library"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Handler struct {
	proxy     *httputil.ReverseProxy
	library   *library.Library
	blackList []string
}

func New(c *conf.Conf, l *library.Library) *Handler {
	target, err := url.Parse(c.EmbyUrl)
	if err != nil {
		log.Panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Director = func(req *http.Request) {
		req.Host = target.Host
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
	}
	return &Handler{
		proxy:     proxy,
		library:   l,
		blackList: c.BlackList,
	}
}

func (h *Handler) ReverseProxy(c *gin.Context) {
	for _, v := range h.blackList {
		if strings.Contains(c.Request.URL.Path, v) {
			c.Status(403)
			return
		}
	}
	h.proxy.ServeHTTP(c.Writer, c.Request)
}

func (h *Handler) Items(c *gin.Context) {
	q := c.Request.URL.Query()
	q.Set("Ids", strings.Join(h.library.ListItem(c.Query("ParentId")), ","))
	q.Del("ParentId")
	//q.Del("IncludeItemTypes")
	c.Request.URL.RawQuery = q.Encode()
	h.ReverseProxy(c)
}

func (h *Handler) Latest(c *gin.Context) {
	q := c.Request.URL.Query()
	q.Set("Ids", strings.Join(h.library.ListItem(c.Query("ParentId")), ","))
	q.Set("SortBy", "DateCreated")
	q.Del("ParentId")
	//q.Del("IncludeItemTypes")
	c.Request.URL.RawQuery = q.Encode()
	c.Request.URL.Path = strings.Replace(c.Request.URL.Path, "/Latest", "", 1)
	h.ReverseProxy(c)
}

func (h *Handler) Resume(c *gin.Context) {
	q := c.Request.URL.Query()
	q.Set("Ids", strings.Join(h.library.ListItem(c.Query("ParentId")), ","))
	q.Set("Filters", "IsResumable")
	q.Del("ParentId")
	//q.Del("IncludeItemTypes")
	c.Request.URL.RawQuery = q.Encode()
	c.Request.URL.Path = strings.Replace(c.Request.URL.Path, "/Resume", "", 1)
	log.Println(c.Request.URL.RawQuery)
	h.ReverseProxy(c)
}
func (h *Handler) Genres(c *gin.Context) {
	q := c.Request.URL.Query()
	q.Set("Ids", strings.Join(h.library.ListItem(c.Query("ParentId")), ","))
	q.Del("ParentId")
	q.Del("IncludeItemTypes")
	c.Request.URL.RawQuery = q.Encode()
	h.ReverseProxy(c)
}
