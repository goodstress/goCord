// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    accounts, err := UnmarshalAccounts(bytes)
//    bytes, err = accounts.Marshal()

package main

import "encoding/json"

type Accounts []Account

func UnmarshalAccounts(data []byte) (Accounts, error) {
	var r Accounts
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Accounts) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Account struct {
	Token string `json:"token"`
}
