// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    verifyResponse, err := UnmarshalVerifyResponse(bytes)
//    bytes, err = verifyResponse.Marshal()

package main

import "encoding/json"

func UnmarshalVerifyResponse(data []byte) (VerifyResponse, error) {
	var r VerifyResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *VerifyResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type VerifyResponse struct {
	Token  string `json:"token"`
	UserID string `json:"user_id"`
}
