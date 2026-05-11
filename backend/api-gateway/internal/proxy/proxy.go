package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func To(target string) gin.HandlerFunc {
	targetURL, _ := url.Parse(target)
	rp := httputil.NewSingleHostReverseProxy(targetURL)

	rp.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Del("Access-Control-Allow-Origin")
		return nil
	}

	return func(c *gin.Context) {
		c.Request.URL.Host = targetURL.Host
		c.Request.URL.Scheme = targetURL.Scheme
		c.Request.Header.Set("X-Forwarded-Host", c.Request.Host)
		rp.ServeHTTP(c.Writer, c.Request)
	}
}
