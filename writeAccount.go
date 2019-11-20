package main

import (
	"io/ioutil"
	"log"
	"os"
)

func (user *User) writeAccount() {
	jsonFile, err := os.Open("accounts.json")
	if err != nil {
		log.Print(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var accounts Accounts
	accounts, err = UnmarshalAccounts(byteValue)
	if err != nil {
		log.Print(err)
	}
	currentAccount := Account{Token: user.auth.token}
	accounts = append(accounts, currentAccount)
	marshalledAccounts, _ := accounts.Marshal()
	_, _ = jsonFile.Write(marshalledAccounts)
	_ = jsonFile.Close()
wg.Done()
}