package main

import (
	"github.com/Yessentemir256/wallet/pkg/types"
	"github.com/Yessentemir256/wallet/pkg/wallet"
	"log"
)

func main() {
	s := &wallet.Service{}

	phone := types.Phone("123456789")
	_, err := s.RegisterAccount(phone)
	if err != nil {
		log.Fatal(err)
	}

	err = s.Deposit(1, 1000)
	if err != nil {
		log.Fatal(err)
	}

	err = s.ExportToFile("test.txt")
	if err != nil {
		log.Fatal(err)
	}

	err = s.ImportFromFile("test.txt")
	if err != nil {
		log.Fatal(err)
	}
}
