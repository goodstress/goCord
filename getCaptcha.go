package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"log"
)

func getCaptcha() string {
	client := resty.New()
	resp, err := client.R().Get("localhost:8080")
	if err != nil {
		log.Println(err)
	}
	captcha := gjson.Get(resp.String(), "RequestKey")
	return captcha.String()
}

