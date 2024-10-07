package types

import "strconv"

// Money представляет собой денежную сумму в минимальных единицах (центы, копейки, дирамы и.т.д).
type Money int64

func (m Money) String() string {
	return strconv.FormatInt(int64(m), 10)
}

// PaymentCategory представляет собой категорию, в которой был совершен платеж (авто, аптеки, рестораны и.т.д.).
type PaymentCategory string

// PaymentStatus представляет собой статус платежа.
type PaymentStatus string

// Предопределенные статусы платежей.
const (
	PaymentStatusOk         PaymentStatus = "OK"
	PaymentStatusFail       PaymentStatus = "FAIL"
	PaymentStatusInProgress PaymentStatus = "INPROGRESS"
)

// Payment представляет информацию о платеже.
type Payment struct {
	ID        string
	AccountID int64
	Amount    Money
	Category  PaymentCategory
	Status    PaymentStatus
}

type Phone string

func (p Phone) String() string {
	return string(p)
}

// Account представляет информацию о счёте пользователя.
type Account struct {
	ID      int64
	Phone   Phone
	Balance Money
}

// Favorite представляет информацию о избранном платеже.
type Favorite struct {
	ID        string
	AccountID int64
	Name      string
	Amount    Money
	Category  PaymentCategory
}
