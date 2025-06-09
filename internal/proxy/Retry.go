package proxy

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

type RetryTransport struct {
	Transport    http.RoundTripper
	MaxRetries   int
	RetryBackoff time.Duration
}

func (t *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var (
		res *http.Response
		err error
		attempt int
		shouldRetry bool
		bodyBytes []byte
	)

	if req.Body != nil {
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		req.Body.Close()
	}

	transport := t.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	// Configurar reintentos
	maxRetries := t.MaxRetries
	if maxRetries <= 0 {
		maxRetries = 3
	}

	retryBackoff := t.RetryBackoff
	if retryBackoff <= 0 {
		retryBackoff = 100 * time.Millisecond 
	}

	for attempt = 0; attempt <= maxRetries; attempt++ {
		if bodyBytes != nil {
			req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}
		res, err = transport.RoundTrip(req)
		
		shouldRetry = false
		
		if err != nil {
			shouldRetry = true
			log.Printf("Error de red en el intento %d: %v", attempt+1, err)
		} else if res != nil {
			switch res.StatusCode {
			case http.StatusServiceUnavailable, // 503
				 http.StatusGatewayTimeout,     // 504
				 http.StatusBadGateway,         // 502
				 http.StatusTooManyRequests:    // 429
				shouldRetry = true
				log.Printf("Código de estado %d en el intento %d", res.StatusCode, attempt+1)
				res.Body.Close()
			}
		}

		
		if !shouldRetry || attempt >= maxRetries {
			break
		}


		backoff := retryBackoff * time.Duration(1<<uint(attempt))
		log.Printf("Reintentando en %v...", backoff)
		time.Sleep(backoff)
	}

	return res, err
}

func ConfigureRetries(proxy *httputil.ReverseProxy, maxRetries int, initialBackoff time.Duration) {
	originalTransport := proxy.Transport
	if originalTransport == nil {
		originalTransport = http.DefaultTransport
	}

	proxy.Transport = &RetryTransport{
		Transport:    originalTransport,
		MaxRetries:   maxRetries,
		RetryBackoff: initialBackoff,
	}
}