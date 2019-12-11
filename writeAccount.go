package main

import (
	"io/ioutil"
	"log"
	"os"
)

func (user *User) writeAccount() {
	defer wg.Done()
	jsonFile, err := os.OpenFile("accounts.json", os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Print(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var accounts Accounts
	accounts, err = UnmarshalAccounts(byteValue)
	if err != nil {
		log.Print("error on UnmarshalAccounts to accounts: ", err)
	}
	currentAccount := Account{Token: user.auth.token}
	accounts = append(accounts, currentAccount)
	marshalledAccounts, _ := accounts.Marshal()
	err = ioutil.WriteFile("accounts.json", marshalledAccounts, os.ModePerm)
	if err != nil {
		log.Print("error writing accounts: ", err)
	}
	//_, err = jsonFile.Write(marshalledAccounts)
	//if err != nil {
	//	log.Print("error writing marshalledAccounts: ", err)
	//}
	err = jsonFile.Close()
	if err != nil {
		log.Print("error closing jsonFile: ", err)
	}
	//wg.Done()
}
