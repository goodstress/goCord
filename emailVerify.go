// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    emailVerify, err := UnmarshalEmailVerify(bytes)
//    bytes, err = emailVerify.Marshal()

package main

import "encoding/json"

func UnmarshalEmailVerify(data []byte) (EmailVerify, error) {
	var r EmailVerify
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *EmailVerify) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type EmailVerify struct {
	Token      string `json:"token"`
	CAPTCHAKey interface{} `json:"captcha_key"`
}
