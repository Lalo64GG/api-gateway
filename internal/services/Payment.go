package services

import (
	"net/http"

	"github.com/Lalo64GG/api-gateway/internal/proxy"
	"github.com/go-chi/chi/v5"
)

type PaymentService struct {
	proxy *proxy.ServiceProxy
}

func NewPaymentService(targetURL string) (*PaymentService, error) {
	proxy, err := proxy.NewServiceProxy("payment",targetURL)

	if err != nil {
		return nil, err
	}

	return &PaymentService{
		proxy: proxy,
	}, nil
}


func (s *PaymentService) Routes() http.Handler{
	r := chi.NewRouter() 

	r.Post("/", s.proxy.Handler())

	r.Get("/{id}", s.proxy.Handler())

	r.Post("/process/{id}", s.proxy.Handler())


	return r
}