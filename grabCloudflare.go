package main

import (
	"crypto/tls"
	"github.com/go-resty/resty/v2"
	"log"
)

func (user *User) grabCloudflare() string {

	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetProxy("127.0.0.1:24000")
	client.SetHeader("User-Agent", user.auth.userAgent)
	resp, err := client.R().
		SetHeaders(map[string]string{
			"User-Agent":                user.auth.userAgent,
			"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
			"Accept-Language":           "en-US,en;q=0.5",
			"DNT":                       "1",
			"Connection":                "keep-alive",
			"Upgrade-Insecure-Requests": "1",
			"Pragma":                    "no-cache",
			"Cache-Control":             "no-cache",
			"TE":                        "Trailers",
		}).
		Get("https://discordapp.com")
	if err != nil {
		log.Println(err)
	}
	// resp.Cookies()
	//set cookies
	user.auth.cookies = resp.Cookies()


	return ""

}
