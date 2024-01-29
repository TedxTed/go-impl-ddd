package services

import "github.com/tedxted/go-impl-ddd/domain/customer"

type OrderConfiguration func(os *OrderService) error

type OrderService struct {
	customers customer.CustomerRepository
}

func NewOrderService(cfgs ...OrderConfiguration) (*OrderService, error) {
	os := &OrderService{}

	for _, cfg := range cfgs {
		err := cfg(os)

		if err != nil {
			return nil, err
		}
	}
	return os, nil
}
