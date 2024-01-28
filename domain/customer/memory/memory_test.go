package memory

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/tedxted/go-impl-ddd/aggregate"
	"github.com/tedxted/go-impl-ddd/domain/customer"
)

func TestMemory_GetCustom(t *testing.T) {
	type testCase struct {
		name        string
		id          uuid.UUID
		expectedErr error
	}

	cust, err := aggregate.NewCustomer("ted")

	if err != nil {
		t.Fatal(err)
	}

	id := cust.GetID()

	repo := MemoryRepository{
		customers: map[uuid.UUID]aggregate.Customer{
			id: cust,
		},
	}

	testCases := []testCase{
		{
			name:        "no customer by id",
			id:          uuid.MustParse("e628bc45-854e-4e7e-9316-1876ac570d90"),
			expectedErr: customer.ErrCustomerNotFound,
		}, {
			name:        "customer by id ",
			id:          id,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := repo.Get(tc.id)

			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}
