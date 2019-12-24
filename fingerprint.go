package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"log"
)

func (user *User) GrabFingerprint() {
	log.Println("proxy", user.auth.proxy)
	log.Print("grabbing fingerprint")
	fingerPrintClient := resty.New()
	//fingerPrintClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	log.Print("created fingerPrintClient for fingerprint")
	proxy := user.auth.proxy
	fingerPrintClient.SetProxy(proxy)
	log.Println("set proxy")
	log.Println("proxy set: ", fingerPrintClient.IsProxySet())

	resp, err := fingerPrintClient.R().
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

	fingerprint := gjson.Get(resp.String(), "grabFingerprint")
	log.Println("grabFingerprint", fingerprint.String())
	user.auth.fingerprint = fingerprint.String()
}
