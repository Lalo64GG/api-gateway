package services

import (
	"github.com/Lalo64GG/api-gateway/internal/config"
)

type Registry struct {
	Payment *PaymentService
}


func NewRegistry(cfg *config.Config) (*Registry, error) {

	paymentService, err := NewPaymentService(cfg.PaymentServiceURL)

	if err != nil {
		return nil, err
	}

	return &Registry{
		Payment: paymentService,
	}, nil
}