# DDD in go

## what is ddd 
- 領域驅動設計（Domain-Driven Design，簡稱DDD）是一種依據其所屬領域進行軟體結構化和模型化的方法。這表示在撰寫軟體之前，首先需要考慮到領域。領域是軟體意圖處理的主題或問題。軟體應該撰寫成能夠反映領域的結構
- DDD主張工程團隊必須與領域物料專家（Subject Matter Experts，SME）進行會議，這些專家是該領域內的專家。這麼做的原因是因為SME掌握了關於領域的知識，而這些知識應該在軟體中得以體現

## Entities and Value Objects
- 以可變和不可變狀態解釋Entities與Value Objects對象

### Entities
- Entities是一個具有標識符的結構體，可以改變狀態，改變狀態意味著實體的值可以改變。
- 一個具有唯一標識符的結構體，具有可以改變的狀態

```go
// Package entities 包含所有跨子域共享的實體
package entity

import (
	"github.com/google/uuid"
)

// Person 是一個代表所有域中的人的實體
type Person struct {
	// ID 是實體的標識符，ID 在所有子域中共享
	ID uuid.UUID
	// Name 是人的名字
	Name string
	// Age 是人的年齡
	Age int
}
```
```go
package entity

import "github.com/google/uuid"

// Item 代表所有子域的一個項目
type Item struct {
	ID          uuid.UUID 
	Name        string    
	Description string    
}

```

### Value Objects
- 有時我們會有一些結構體是不可變的，不需要唯一標識符，這些結構體被稱為Value Objects
- 沒有標識符且創建後值持久的結構體(structs without an identifier and persistent values after creation)
- Value Objects通常在域內部找到，用於描述該域中的某些方面。我們現在將創建一個值對象，Transaction，一旦交易執行，就不能改變狀態(Value objects are often found inside domains and used to describe certain aspects in that domain)

```go
package valueobject

import (
	"time"
)

// Transaction is a payment between two parties
type Transaction struct {
	// all values lowercase since they are immutable
	amount    int
	from      uuid.UUID
	to        uuid.UUID
	createdAt time.Time
}
```

### Aggregates — Combined Entities and Value Objects
- Aggregates是一組Entities和Value Objects的組合。因此，在我們的案例中，我們可以開始創建一個新的聚合，即 Customer
- 業務邏輯將應用於 Customer 聚合，而不是每個Entities持有邏輯。聚合不允許直接訪問底層Entities
- 在現實生活中，為了正確表示數據，通常需要多個Entities，例如，一個 Customer。它是一個 Person，但他/她可以持有 Products，並進行交易
- 聚合中的一個重要規則是，它們應該只有一個實體作為根實體。這意味著根實體的引用也用於引用聚合。對於我們的 customer 聚合，這意味著 Person 的 ID 是唯一標識符

```go
// Package aggregates 包含結合許多實體成為一個完整對象的聚合
package aggregate

import (
	"github.com/percybolmer/ddd-go/entity"
	"github.com/percybolmer/ddd-go/valueobject"
)

// Customer 是一個聚合，結合了表示客戶所需的所有實體
type Customer struct {
	// person 是客戶的根實體
	// 這意味著 person.ID 是這個聚合的主要標識符
	person *entity.Person 
	// 一個客戶可以持有許多產品
	products []*entity.Item 
	// 一個客戶可以進行許多交易
	transactions []valueobject.Transaction 
}
```
- 將所有Entities設為指針，這是因為Entities可以改變狀態，我希望這可以在運行時訪問它的所有實例中反映出來。雖然值對象則作為非指針持有，因為它們不能改變狀態(The value objects are held as nonpointers though since they cannot change state)

### Factories — Encapsulate complex logic
- 工廠模式用於封裝複雜的創建邏輯，使調用者無需了解實例的具體實現細節
- 在 Domain-Driven Design（領域驅動設計）中，建議使用工廠模式創建複雜的聚合、倉庫和服務

```go
package aggregate

import (
    "errors"
    "github.com/google/uuid"
    "github.com/percybolmer/ddd-go/entity"
    "github.com/percybolmer/ddd-go/valueobject"
)

type Customer struct {
    person        *entity.Person
    products      []*entity.Item
    transactions  []valueobject.Transaction
}

var ErrInvalidPerson = errors.New("一個客戶必須有一個有效的人物")

func NewCustomer(name string) (Customer, error) {
    if name == "" {
        return Customer{}, ErrInvalidPerson
    }

    person := &entity.Person{
        Name: name,
        ID:   uuid.New(),
    }

    return Customer{
        person:       person,
        products:     make([]*entity.Item, 0),
        transactions: make([]valueobject.Transaction, 0),
    }, nil
}
```

## Repositories — The Repository Pattern
- It is a pattern that relies on hiding the implementation of the storage/database solution behind an interface
- This allows us to define a set of methods that has to be present, and if they are present it is qualified to be used as a repository.
- The advantage of this design pattern is that it allows us to exchange the solution without breaking anything.

```go
// Package Customer holds all the domain logic for the customer domain.
package customer

import (
	"github.com/google/uuid"
	"github.com/percybolmer/ddd-go/aggregate"
)
var (
	// ErrCustomerNotFound is returned when a customer is not found.
	ErrCustomerNotFound = errors.New("the customer was not found in the repository")
	// ErrFailedToAddCustomer is returned when the customer could not be added to the repository.
	ErrFailedToAddCustomer = errors.New("failed to add the customer to the repository")
	// ErrUpdateCustomer is returned when the customer could not be updated in the repository.
	ErrUpdateCustomer = errors.New("failed to update the customer in the repository")
)
// CustomerRepository is a interface that defines the rules around what a customer repository
// Has to be able to perform
type CustomerRepository interface {
	Get(uuid.UUID) (aggregate.Customer, error)
	Add(aggregate.Customer) error
	Update(aggregate.Customer) error
}
```

## Services — Connecting the Business Logic
- 把所有鬆散耦合的倉庫綁定到一個業務邏輯中，以滿足某個特定領域的需求 A service will tie all loosely coupled repositories into a business logic that fulfills the needs of a certain domain
- 在酒館的例子中，我們可能會有一個訂單服務，負責將倉庫連接起來以執行訂單。因此，服務將訪問客戶倉庫和產品倉庫