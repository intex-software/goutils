package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/MadAppGang/httplog"
)

// ! TLS Configuration
func Proxy(target, addr string) {
	proxyURL, err := url.Parse(target)
	if err != nil {
		panic(err)
	}

	logger := httplog.LoggerWithConfig(httplog.LoggerConfig{
		RouterName:  "Proxy",
		Formatter:   httplog.FullFormatterWithRequestAndResponseHeadersAndBody,
		CaptureBody: true,
	})

	proxy := NewSingleHostReverseProxy(proxyURL)
	proxyHandler := logger(proxy)

	// HTTP-Server mit Proxy-Handler starten
	fmt.Println("Proxy-Server gestartet auf", addr, "und leitet auf", proxyURL.String())
	http.ListenAndServe(addr, proxyHandler)
}

func NewSingleHostReverseProxy(target *url.URL) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		rewriteRequestURL(req, target)
	}
	return &httputil.ReverseProxy{Director: director}
}

func rewriteRequestURL(req *http.Request, target *url.URL) {
	targetQuery := target.RawQuery
	req.URL.Scheme = target.Scheme
	req.URL.Host = target.Host
	req.Host = target.Host
	req.URL.Path, req.URL.RawPath = joinURLPath(target, req.URL)
	if targetQuery == "" || req.URL.RawQuery == "" {
		req.URL.RawQuery = targetQuery + req.URL.RawQuery
	} else {
		req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
	}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func joinURLPath(a, b *url.URL) (path, rawpath string) {
	if a.RawPath == "" && b.RawPath == "" {
		return singleJoiningSlash(a.Path, b.Path), ""
	}
	// Same as singleJoiningSlash, but uses EscapedPath to determine
	// whether a slash should be added
	apath := a.EscapedPath()
	bpath := b.EscapedPath()

	aslash := strings.HasSuffix(apath, "/")
	bslash := strings.HasPrefix(bpath, "/")

	switch {
	case aslash && bslash:
		return a.Path + b.Path[1:], apath + bpath[1:]
	case !aslash && !bslash:
		return a.Path + "/" + b.Path, apath + "/" + bpath
	}
	return a.Path + b.Path, apath + bpath
}
