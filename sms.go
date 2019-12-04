package main

import (
	"context"
	"encoding/json"
	"github.com/bxcodec/faker/v3"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"log"
	"strings"
	"time"
)

func (user *User) smsVerification() {
	config := smsApi{
		apiKey:  "***REMOVED***",
		service: "ds",
		country: 0,
	}

	number := getNumber(config)
	//todo: submit number to discord.
	sendToDiscord(number)
	//todo: notify sms-activate.ru that number is ready.

	//todo: get code from sms msg

}

func sendToDiscord(number phoneObject) {

}

type smsApi struct {
	apiKey, service string
	country         int
}

func getNumber(config smsApi) phoneObject {
	client := resty.New()
	URLString := "http://sms-activate.ru/stubs/handler_api.php?api_key=" + config.apiKey + "&action=getNumber&service=" + config.service + "&country=" + string(config.country)
	resp, err := client.R().
		Post(URLString)
	if err != nil {
		log.Print("Error getting phone number: ", err)
	}
	object := phoneObject{
		phoneNumber: strings.TrimLeft(resp.String(), ":"),
		numberId:    strings.TrimRight(resp.String(), ":"),
	}
	log.Print("Phone Number: ", object.phoneNumber)
	log.Print("Number ID: ", object.numberId)
	return object
}

type phoneObject struct {
	phoneNumber, numberId string
}
