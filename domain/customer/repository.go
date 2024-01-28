package customer

import (
	"errors"

	"github.com/google/uuid"
	"github.com/tedxted/go-impl-ddd/aggregate"
)

var (
	ErrCustomerNotFound    = errors.New("the customer was not fount in repository")
	ErrFailedToAddCustomer = errors.New("failed to add customer")
	ErrUpdateCustomer      = errors.New("failed to update teh customer")
)

type CustomerRepository interface {
	Get(uuid.UUID) (aggregate.Customer, error)
	Add(aggregate.Customer) error
	Update(aggregate.Customer) error
}
