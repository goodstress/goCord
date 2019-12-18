package main

import (
	"github.com/go-resty/resty/v2"
	"log"
	"strings"
	"sync"
	"time"
)

func (user *User) smsVerification(wg *sync.WaitGroup) {
	defer wg.Done()
	//verification.Add(1)
	time.Sleep(20 * time.Second)
	user.smsApi = smsApi{
		apiKey:  "***REMOVED***",
		service: "ds",
		country: 0,
	}

	user.getNumber()
	//submit number to discord.
	log.Print("sending phone to discord")
	user.sendPhoneToDiscord()
	//notify sms-activate.ru that number is ready.
	log.Print("notifying SMS service number is ready")
	user.notifyReady()
	//get code from sms msg
	log.Print("attempting to grab sms code")
	code := user.getCode()

	user.sendCodeToDiscord(code)
}

type smsApi struct {
	apiKey, service string
	country         int
}

func (user *User) getNumber() PhoneNumber {
	client := resty.New()
	URLString := "http://sms-activate.ru/stubs/handler_api.php?api_key=" + user.smsApi.apiKey + "&action=getNumber&service=" + user.smsApi.service + "&country=0"
	resp, err := client.R().
		Post(URLString)
	if err != nil {
		log.Print("Error getting phone number: ", err)
	}
	phoneSlice := strings.Split(resp.String(), ":")
	object := PhoneNumber{
		phoneNumber: phoneSlice[2],
		numberId:    phoneSlice[1],
	}
	user.PhoneNumber = object
	log.Print("Phone Number: ", user.PhoneNumber.phoneNumber)
	log.Print("Number ID: ", user.PhoneNumber.numberId)
	return object

}
func (user *User) sendPhoneToDiscord() {
	client := resty.New()
	client.SetProxy(user.auth.proxy)
	payload := `{"phone": ` + `"+` + user.PhoneNumber.phoneNumber + `"}`
	_, err := client.R().
		SetBody(payload).
		SetHeaders(map[string]string{
			"User-Agent":         user.auth.userAgent,
			"Accept":             "*/*",
			"Accept-Language":    "en-US,en;q=0.5",
			"Authorization":      user.auth.token,
			"Content-Type":       "application/json;charset=utf-8",
			"Connection":         "keep-alive",
			"X-Super-Properties": user.auth.SuperProp,
			"Referer":            "https://discordapp.com/channels/@me",
			"TE":                 "Trailers",
		}).
		Post("https://discordapp.com/api/v6/users/@me/phone")
	if err != nil {
		log.Print(err)
	}
	user.notifyReady()

}
func (user *User) notifyReady() {
	URLString := "http://sms-activate.ru/stubs/handler_api.php?api_key=" + user.smsApi.apiKey + "&action=setStatus&status=1&id=" + user.PhoneNumber.numberId
	client := resty.New()
	_, err := client.R().
		Post(URLString)
	if err != nil {
		log.Print("error notifying SMS api of being ready: ", err)
	}
}

func (user *User) getCode() string {
	URLString := "http://sms-activate.ru/stubs/handler_api.php?api_key=" + user.smsApi.apiKey + "&action=getStatus&status=1&id=" + user.PhoneNumber.numberId
	client := resty.New()
	checkCode := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-checkCode.C:
			resp, err := client.R().
				Post(URLString)
			if err != nil {
				log.Print("error grabbing phone code: ", err)
			}
			if strings.HasPrefix(resp.String(), "STATUS_WAIT_CODE") {
				log.Print("sms code not ready, waiting another 10 seconds.")

			} else {
				log.Print(resp.String())

				regularCode := strings.Replace(resp.String(), "STATUS_OK:", "", 1)

				return regularCode
			}

		}
	}

}

func (user *User) sendCodeToDiscord(code string) {
	client := resty.New()
	client.SetProxy(user.auth.proxy)
	payload := `{"code": ` + `"` + code + `"}`
	_, err := client.R().
		SetBody(payload).
		SetHeaders(map[string]string{
			"User-Agent":         user.auth.userAgent,
			"Accept":             "*/*",
			"Accept-Language":    "en-US,en;q=0.5",
			"Authorization":      user.auth.token,
			"Content-Type":       "application/json;charset=utf-8",
			"Connection":         "keep-alive",
			"X-Super-Properties": user.auth.SuperProp,
			"Referer":            "https://discordapp.com/channels/@me",
			"TE":                 "Trailers",
		}).
		Post("https://discordapp.com/api/v6/users/@me/phone/verify")
	if err != nil {
		log.Print(err)
	}

}
