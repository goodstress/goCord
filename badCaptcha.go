package main

import (
	"github.com/go-resty/resty/v2"
	"log"
)

func (user *User) badCaptcha() {
	client := resty.New()
	client.SetQueryParam("key", configuration.apikey)
	client.SetQueryParam("action", "reportbad")
	client.SetQueryParam("id", user.auth.Captcha.captchaID)
	client.SetQueryParam("json", "1")
	resp, err := client.R().Post("https://2captcha.com/res.php")
	if err != nil {
		log.Println(err)
	}
	log.Println(resp.String())
}

func (user *User) goodCaptcha() {
	client := resty.New()
	client.SetQueryParam("key", configuration.apikey)
	client.SetQueryParam("action", "reportgood")
	client.SetQueryParam("id", user.auth.Captcha.captchaID)
	client.SetQueryParam("json", "1")
	resp, err := client.R().Post("https://2captcha.com/res.php")
	if err != nil {
		log.Println(err)
	}
	log.Println(resp.String())
}
