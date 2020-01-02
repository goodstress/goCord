package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"log"
	"time"
)

func (user *User) GrabFingerprint() {
	log.Println("proxy", user.auth.proxy)
	log.Print("grabbing fingerprint")
	client := resty.New()
	//client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	log.Print("created client for fingerprint")
	proxy := user.auth.proxy
	client.SetProxy(proxy)
	log.Println("set proxy")
	log.Println("proxy set: ", client.IsProxySet())
	client.SetTimeout(10 * time.Second)
	resp, err := client.R().
		SetHeaders(map[string]string{
			"User-Agent":      user.auth.userAgent,
			"Accept":          "*/*",
			"Accept-Language": "en-US,en;q=0.5",
			// "X-Track": Auth.superProp,
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

	fingerprint := gjson.Get(resp.String(), "fingerprint")
	log.Println("grabFingerprint", fingerprint.String())
	user.auth.fingerprint = fingerprint.String()
}
