package wallet

import (
	"errors"
	"github.com/Yessentemir256/wallet/pkg/types"
	"testing"
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
