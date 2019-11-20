package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"log"
	"time"
)

func getCaptcha() string {

	check := time.NewTicker(5*time.Second)
	for {
		select {
		case <-check.C:
			response := gitIt()
			captcha := gjson.Get(response, "RequestKey")
			if len(captcha.String()) < 5 {
				log.Print("error occurred, captcha was nil")

			} else {
				return captcha.String()
			}

		}
	}
}

func gitIt() string {
	client := resty.New()
	client.RemoveProxy()
	resp, err := client.R().Get("http://0.0.0.0:8080")
	if err != nil {
		log.Println(err)

	}
	return resp.String()

}


