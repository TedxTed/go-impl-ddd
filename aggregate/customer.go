package aggregate

import (
	"errors"

	"github.com/google/uuid"
	"github.com/tedxted/go-impl-ddd/entity"
	"github.com/tedxted/go-impl-ddd/valueobject"
)

var (
	ErrInvalidPerson = errors.New("a customer has to have a valid name")
)

type Customer struct {
	// person is the root entity of customer
	// which mean person.ID is the main identifier for the customer

	person   *entity.Person
	products []*entity.Item

	transaction []valueobject.Transaction
}

// factory pattern
// New Customer is a factory to create aggregate
// it will aggregate that the name is not empty
func NewCustomer(name string) (Customer, error) {
	if name == "" {
		return Customer{}, ErrInvalidPerson
	}

	person := &entity.Person{
		Name: name,
		ID:   uuid.New(),
	}

	return Customer{
		person:      person,
		products:    make([]*entity.Item, 0),
		transaction: make([]valueobject.Transaction, 0),
	}, nil

}

// person is private in the Customer , so we need to add a public function to get the person id
func (c *Customer) GetID() uuid.UUID {
	return c.person.ID
}

func (c *Customer) SetID(id uuid.UUID) {
	if c.person == nil {
		c.person = &entity.Person{}
	}
	c.person.ID = id
}

func (c *Customer) SetName(name string) {
	if c.person == nil {
		c.person = &entity.Person{}
	}

	c.person.Name = name
}

func (c *Customer) GetName() string {
	return c.person.Name
}
