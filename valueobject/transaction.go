package valueobject

import (
	"time"

	"github.com/google/uuid"
)

// transaction is a valueobject because it has no identifier ,and is immitible
type Transaction struct {
	amount    int
	from      uuid.UUID
	to        uuid.UUID
	createdAt time.Time
}
