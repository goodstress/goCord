// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    emails, err := UnmarshalEmails(bytes)
//    bytes, err = emails.Marshal()

package main

import "encoding/json"

type Emails []Email

func UnmarshalEmails(data []byte) (Emails, error) {
	var r Emails
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Emails) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Email struct {
	ID          string        `json:"id"`
	From        string        `json:"from"`
	To          string        `json:"to"`
	Cc          interface{}   `json:"cc"`
	Subject     string        `json:"subject"`
	BodyText    string        `json:"body_text"`
	BodyHTML    string        `json:"body_html"`
	CreatedAt   string        `json:"created_at"`
	Attachments []interface{} `json:"attachments"`
}
