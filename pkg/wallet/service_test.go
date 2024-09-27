package wallet

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/Yessentemir256/wallet/pkg/types"
	"github.com/google/uuid"
)

func TestService_FindAccountByID_success(t *testing.T) {
	// Создаем сервис с тестовым аккаунтом
	testAccount := &types.Account{ID: 123, Balance: 1000}
	s := &Service{
		accounts: []*types.Account{testAccount},
	}

	// Тест: Ищем существующий аккаунт
	foundAccount, err := s.FindAccountByID(testAccount.ID)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	if foundAccount != testAccount {
		t.Errorf("Expected account %v, but got %v", testAccount, foundAccount)
	}
}

func TestService_FindAccountByID_notFound(t *testing.T) {
	// Создаем сервис без аккаунтов
	s := &Service{
		accounts: []*types.Account{},
	}
	var nonExistingAccountID int64 = 456

	// Тест: Ищем несуществующий аккаунт
	_, err := s.FindAccountByID(nonExistingAccountID)

	expectedError := errors.New("account not found")

	if err == nil || err.Error() != expectedError.Error() {
		t.Errorf("Expected 'account not found' error, but got %v", err)
	}
}

func TestReject(t *testing.T) {
	s := &Service{
		payments: []*types.Payment{
			{ID: "1", AccountID: 1, Amount: 100, Category: "Test", Status: types.PaymentStatusOk},
		},
		accounts: []*types.Account{
			{ID: 1, Balance: 0},
		},
	}

	err := s.Reject("1")
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// Check if payment status is updated
	if s.payments[0].Status != types.PaymentStatusFail {
		t.Errorf("Payment status not updated")
	}

	// Check if funds are added back to the account
	if s.accounts[0].Balance != 100 {
		t.Errorf("Funds not added back to account")
	}
}

func TestFindPaymentByID(t *testing.T) {
	s := &Service{
		payments: []*types.Payment{
			{ID: "1", AccountID: 1, Amount: 100, Category: "Test", Status: types.PaymentStatusOk},
		},
	}

	payment, err := s.FindPaymentByID("1")
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	if payment.ID != "1" {
		t.Errorf("Incorrect payment found")
	}
}

func (s *testService) addAccountWithBalance(phone types.Phone, balance types.Money) (*types.Account, error) {
	// регистрируем там пользователя
	account, err := s.RegisterAccount(phone)
	if err != nil {
		return nil, fmt.Errorf("can't register account, error = %v", err)
	}

	// пополняем его счет
	err = s.Deposit(account.ID, balance)
	if err != nil {
		return nil, fmt.Errorf("can't deposit account, error = %v", err)
	}

	return account, nil
}

func TestService_FindPaymentByID_success(t *testing.T) {
	// создаем сервис
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// пробуем найти платёж
	payment := payments[0]
	got, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("FindPaymentByID(): error = %v", err)
		return
	}

	// сравниваем платежи
	if !reflect.DeepEqual(payment, got) {
		t.Errorf("FindPaymentByID(): wrong payment returned = %v", err)
		return
	}
}

func TestService_FindPaymentByID_fail(t *testing.T) {
	// создаем сервис
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// пробуем найти несуществующий платёж
	_, err = s.FindPaymentByID(uuid.New().String())
	if err == nil {
		t.Error("FindPaymentByID(): must return error, returned nil")
		return
	}

	if err != ErrPaymentNotFound {
		t.Errorf("FindPaymentByID(): must return ErrPaymentNotFound, returned = %v", err)
		return
	}
}

type testService struct {
	*Service // embedding (встраивание)
}

func newTestService() *testService {
	return &testService{Service: &Service{}} //функция-конструктор
}

type testAccount struct {
	phone    types.Phone
	balance  types.Money
	payments []struct {
		amount   types.Money
		category types.PaymentCategory
	}
}

var defaultTestAccount = testAccount{
	phone:   "+992000000001",
	balance: 10_000_00,
	payments: []struct {
		amount   types.Money
		category types.PaymentCategory
	}{
		{amount: 1_000_00, category: "auto"},
	},
}

func (s *testService) addAccount(data testAccount) (*types.Account, []*types.Payment, error) {
	// регистрируем там пользователя
	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, fmt.Errorf("can't register account, error = %v", err)
	}

	// пополняем его счет
	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("can't register account, error = %v", err)
	}

	// выполняем платежи
	// можем создать слайс нужной длины, послкольку знаем размер
	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments {
		// тогда здесь работаем просто через index, а не через append
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, fmt.Errorf("can't make payment, error = %v", err)
		}
	}

	return account, payments, nil
}

func TestService_Reject_success(t *testing.T) {
	// создаем сервис
	s := newTestService()                                // это наша функция конструктор, которая вышла из embedding
	_, payments, err := s.addAccount(defaultTestAccount) // добавление пользователя с помощью метода который принадлежит testService
	if err != nil {
		t.Error(err)
		return
	}

	// пробуем отменить платеж
	payment := payments[0] // выбираем конкретно платеж который в итоге хотим итменить
	err = s.Reject(payment.ID)
	if err != nil {
		t.Errorf("Reject(): error = %v", err)
		return
	}

	savedPayment, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("Reject(): can't find payment by id, error = %v", err)
		return
	}
	if savedPayment.Status != types.PaymentStatusFail {
		t.Errorf("Reject(): status didn't changed, payment = %v", savedPayment)
		return
	}

	savedAccount, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		t.Errorf("Reject(): can't find account by id, error = %v", err)
		return
	}
	if savedAccount.Balance != defaultTestAccount.balance {
		t.Errorf("Reject(): balance didn't changed, account = %v", savedAccount)
		return
	}
}

func TestService_Repeat_success(t *testing.T) {
	// создаем сервис
	s := newTestService()                                // это наша функция конструктор, которая вышла из embedding
	_, payments, err := s.addAccount(defaultTestAccount) // добавление пользователя с помощью метода который принадлежит testService
	if err != nil {
		t.Error(err)
		return
	}

	// пробуем повторить платеж
	payment := payments[0] // выбираем конкретно платеж который в итоге хотим повторить
	paymentRepeated, err := s.Repeat(payment.ID)
	if err != nil {
		t.Errorf("Repeat(): error = %v", err)
		return
	}

	// проверка на то что id разные
	if payment.ID == paymentRepeated.ID {
		t.Errorf("Repeat(): ID is not different, paymentID = %v", payment.ID)
		return
	}

	// проверка на то что суммы одинковые
	if payment.Amount != paymentRepeated.Amount {
		t.Errorf("Repeat(): amount is not equal, paymentID = %v", payment.ID)
		return
	}

}

func TestService_FavoritePayment_success(t *testing.T) {
	// создаем сервис
	s := newTestService()                                // это наша функция конструктор, которая вышла из embedding
	_, payments, err := s.addAccount(defaultTestAccount) // добавление пользователя с помощью метода который принадлежит testService
	if err != nil {
		t.Error(err)
		return
	}

	// пробуем
	payment := payments[0] // выбираем конкретно платеж который в итоге хотим сделать избранным.
	var name string
	favorite, err := s.FavoritePayment(payment.ID, name)
	if err != nil {
		t.Errorf("FavoritePayment(): error = %v", err)
		return
	}

	// проверка на то что id разные.
	if payment.ID != favorite.ID {
		t.Errorf("FavoritePayment(): ID is equal, paymentID = %v", payment.ID)
		return
	}

	// проверка на то что суммы одинковые
	if payment.Amount != favorite.Amount {
		t.Errorf("Repeat(): amount is not equal, paymentID = %v", payment.ID)
		return
	}
}
