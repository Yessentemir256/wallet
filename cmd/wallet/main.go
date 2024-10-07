package main

import (
	"github.com/Yessentemir256/wallet/pkg/types"  // Импортируйте ваш пакет types
	"github.com/Yessentemir256/wallet/pkg/wallet" // Импортируйте ваш пакет wallet
	"log"
)

func main() {
	s := &wallet.Service{} // Создание экземпляра Service

	// Заполнение s данными
	phone := types.Phone("123456789")
	_, err := s.RegisterAccount(phone)
	if err != nil {
		log.Fatal(err)
	}

	// Добавление денег на счет
	err = s.Deposit(1, 1000)
	if err != nil {
		log.Fatal(err)
	}

	// Вызов метода ExportToFile для тестирования
	err = s.ExportToFile("test.txt")
	if err != nil {
		log.Fatal(err)
	}
}
