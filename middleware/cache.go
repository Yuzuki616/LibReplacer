package middleware

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"net/http"
	"strings"
)

type responseCache struct {
	Status int
	Header http.Header
	Body   []byte
}

type CacheWriter struct {
	c       *cache.Cache
	status  int
	written bool
	body    *bytes.Buffer
	gin.ResponseWriter
}

func (w *CacheWriter) WriteHeader(code int) {
	w.status = code
	w.written = true
	w.ResponseWriter.WriteHeader(code)
}

func (w *CacheWriter) Status() int {
	return w.ResponseWriter.Status()
}

func (w *CacheWriter) Written() bool {
	return w.ResponseWriter.Written()
}

func (w *CacheWriter) Write(data []byte) (int, error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}

func (w *CacheWriter) WriteString(data string) (n int, err error) {
	w.body.WriteString(data)
	return w.ResponseWriter.WriteString(data)
}

func (w *CacheWriter) SetCache(key string) {
	if w.written && w.status == 200 {
		w.c.SetDefault(key, &responseCache{
			Status: w.status,
			Header: w.Header(),
			Body:   w.body.Bytes(),
		})
	}
}

func genHashString(data string) string {
	s := sha1.Sum([]byte(data))
	return hex.EncodeToString(s[:])
}

var blackList = []string{
	"played",
	"isfavorite",
	"isresumable",
	"likes",
	"userdata",
}

func (m *Middleware) CacheResponse(c *gin.Context) {
	for _, v := range blackList {
		if strings.Contains(strings.ToLower(c.Request.URL.RawQuery), v) {
			return
		}
	}
	q, _, _ := strings.Cut(c.Request.URL.RawQuery, "&X-")
	s := genHashString(q)
	if b, e := m.cache.Get(s); e {
		res := b.(*responseCache)
		c.Writer.WriteHeader(res.Status)
		for k, v := range res.Header {
			for i := range v {
				c.Writer.Header().Set(k, v[i])
			}
		}
		c.Writer.Write(res.Body)
		c.Abort()
		return
	}
	w := &CacheWriter{c: m.cache, ResponseWriter: c.Writer, body: bytes.NewBufferString("")}
	c.Writer = w
	c.Next()
	w.SetCache(s)
}
