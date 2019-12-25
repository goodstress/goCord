package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/pquerna/otp/totp"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"log"
)

func (user *User) enableTotp() {
	key, _ := totp.Generate(totp.GenerateOpts{
		Issuer:      "discord",
		AccountName: user.details.username,
	})
	log.Print("example code: ", key.String())
	log.Print("2fa key is: ", key)
	client := resty.New()
	client.SetProxy(user.auth.proxy)
	client.SetHeader("User-Agent", user.auth.userAgent)
	client.SetHeader("Authorization", user.auth.token)
	client.SetHeader("content-type", "application/json")
	defaultPayload := `{"code":"411977","secret":"MAM3LQ4RQTNU6TPB"}`
	changeCode, err := sjson.Set(defaultPayload, "code", key.String())
	if err != nil {
		log.Print("error changing code of default payload: ", err)
	}
	payload, err := sjson.Set(changeCode, "secret", key.Secret())
	if err != nil {
		log.Print("error changing secret of default payload: ", err)
	}
	resp, err := client.R().
		SetBody(payload).
		Post("https://discordapp.com/api/v6/users/@me/mfa/totp/enable")
	if err != nil {
		log.Println("error sending ENABLE TOTP request: ", err)
	}

	log.Println(resp.String())
	token := gjson.Get(resp.String(), "token").String()
	user.auth.token = token
}
