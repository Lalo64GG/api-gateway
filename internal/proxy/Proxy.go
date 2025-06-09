package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)


type ServiceProxy struct {
	Name string
	TargetURL *url.URL
	Proxy *httputil.ReverseProxy
}

func NewServiceProxy(name, targetURLStr string) (*ServiceProxy, error){
	targetURL, err := url.Parse(targetURLStr)

	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	originalDirector := proxy.Director
	proxy.Director = func (req *http.Request){
		originalDirector(req)


		req.Header.Set("X-Proxy", "API-Gateway")
		req.Header.Set("X-Forwarded-Host", req.Host)
		req.Host = targetURL.Host
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error){
		log.Printf("Error proxying to %s: %v", name, err)
		http.Error(w, "Service Unavailable", http.StatusBadGateway)
	}

	proxy.Transport = &http.Transport{
		ResponseHeaderTimeout: 10 * time.Second,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
	}

	return &ServiceProxy{
		Name:      name,
		TargetURL: targetURL,
		Proxy:     proxy,
	}, nil
}

func (sp *ServiceProxy) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Proxying request to %s: %s %s", sp.Name, r.Method, r.URL.Path)

		sp.Proxy.ServeHTTP(w, r)
	}
}

func StripPrefix(prefix string, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, prefix) {
			r.URL.Path = strings.TrimPrefix(r.URL.Path, prefix)
			r.RequestURI = r.URL.Path
		}
		handler.ServeHTTP(w, r)
	})
}