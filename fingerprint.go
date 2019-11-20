package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"log"
)

func (user *User) grabFingerprint()  {
	log.Print("grabbing fingerprint")
	client := new(resty.Client)
	client.SetProxy(user.auth.proxy)
	resp, err := client.R().
		SetHeaders(map[string]string{
			"User-Agent":      user.auth.userAgent,
			"Accept":          "*/*",
			"Accept-Language": "en-US,en;q=0.5",
			// "X-Track": auth.superProp,
			"DNT":        "1",
			"Connection": "keep-alive",
			"Referer":    "https://discordapp.com/",
			"TE":         "Trailers",
		}).
		Get("https://discordapp.com/api/v6/experiments")
	if err != nil {
		log.Print("fingerprint error:::")
		log.Println(err)
	}

	fingerprint := gjson.Get(resp.String(), "grabFingerprint")
	log.Println("grabFingerprint", fingerprint.String())
	user.auth.fingerprint = fingerprint.String()
}
