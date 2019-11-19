package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"golang.org/x/net/http2"
	b64 "encoding/base64"
	"math/rand"
	"github.com/bxcodec/faker/v3"
	"github.com/go-resty/resty/v2"
	"github.com/mssola/user_agent"
	"log"
	"net/http"
	"time"
	"mvdan.cc/xurls/v2"
	// "go.zoe.im/surferua"
	// 	    "github.com/mileusna/useragent"
	"github.com/tidwall/gjson"
)

func main() {
	fmt.Println("vim-go")



	//create user
	createUser()

}

//noinspection SpellCheckingInspection
func randomStickyIP() string{
	rand.Seed(time.Now().UnixNano())

	min := 10001
	max := 29999
	randomPort := rand.Intn(max - min + 1) + min
	ipString := "***REMOVED***:***REMOVED***" + "@us.smartproxy.com:" + string(randomPort)
	return ipString
}
func createUser() {

	user := new(User)
	//set proxy
	user.auth.proxy = randomStickyIP()
	//noinspection SpellCheckingInspection
	ua := user_agent.New("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36")
	user.init()

	log.Println(ua.Bot())
	// userp := &user

}

func (user *User) init() {


	user.grabCloudflare()
	user.grabFingerprint()
	//setUsername
	user.details.username = faker.Username()
	//create superProp
	user.createSuperProp()
	//user.createXTrack()
	user.details.password = faker.Password()

}

func (user *User) createSuperProp() {
	var prop SuperProp
	ua := user_agent.New(user.auth.userAgent)
	prop.Browser, prop.BrowserVersion = ua.Browser()
	prop.BrowserUserAgent = user.auth.userAgent
	prop.OS = ua.OSInfo().Name
	prop.OSVersion = ua.OSInfo().Version
	//current build number
	prop.ClientBuildNumber = 49868
	prop.ClientEventSource = nil
	prop.Device, prop.Referrer, prop.ReferringDomain, prop.ReferringDomainCurrent, prop.ReferrerCurrent = "", "", "", "", ""
	prop.ReleaseChannel = "stable"
	//set superProp
	marshalledProp, err := prop.Marshal()
	if err != nil {
		log.Print(err)
	}
	sEnc := b64.StdEncoding.EncodeToString(marshalledProp)
	log.Print("Encoded super prop: ", sEnc)
	user.auth.SuperProp = sEnc

}
//func (user *User) createXTrack() {
//	var prop SuperProp
//	ua := user_agent.New(user.auth.userAgent)
//	prop.Browser, prop.BrowserVersion = ua.Browser()
//	prop.BrowserUserAgent = user.auth.userAgent
//	prop.OS = ua.OSInfo().Name
//	prop.OSVersion = ua.OSInfo().Version
//	prop.ClientBuildNumber = 9999
//	prop.ClientEventSource = nil
//	prop.Device, prop.Referrer, prop.ReferringDomain, prop.ReferringDomainCurrent, prop.ReferrerCurrent = "", "", "", "", ""
//	prop.ReleaseChannel = "stable"
//	//set superProp
//	user.auth.Xtrack = prop
//}
func (user *User) genUserAgent() {
	//noinspection ALL
	agent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36"
	user.auth.userAgent = agent
	log.Println("Set useragent")
}

// func superProp(agent string) string {
// 	ua.Parse(agent)
// }
// func xTrack() string {

// }

func (user *User) register() {
	captcha := getCaptcha()
	realRegister := RegPayload{Fingerprint: user.auth.fingerprint, Email: user.details.email, Username: user.details.username, Password: user.details.password, Invite: nil, Consent: true, GiftCodeSkuID: nil, CAPTCHAKey: captcha}
	registerURL := "https://discordapp.com/api/v6/auth/register"
	client := user.client
	resp, err := client.R().
		SetBody(realRegister).
		SetHeaders(map[string]string{
			"Accept":          "*/*",
			"Accept-Language": "en-US,en;q=0.5",
			"DNT":             "1",
			"Connection":      "keep-alive",
			"Referer":         "https://discordapp.com/",
			"TE":              "Trailers",
		}).
		Post(registerURL)
	if err != nil {
		log.Println(err)
	}
	token := gjson.Get(resp.String(), "token").String()
	user.auth.token = token
	log.Println("set token in user to: ", token)

}
type emailString struct {
	email string
}
func (user *User) GenEmail() string {
	client := user.client
	client.SetProxy(user.auth.proxy)
	_, err := client.R().
		SetBody(`{"min_name_length": 10,"max_name_length": 10}`).
		SetHeaders(map[string]string{
		"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:69.0) Gecko/20100101 Firefox/69.0",
		"Accept": "application/json, text/plain, */*",
		  "Accept-Language": "en-US,en;q=0.5",
		  "Content-Type": "application/json;charset=utf-8",
		  "Origin": "https://temp-mail.io",
		  "Connection": "keep-alive",
		  "Referer": "https://temp-mail.io/en",
		  "TE": "Trailers",
		}).Options("https://api.internal.temp-mail.io/api/v2/email/new")
	if err != nil {
		log.Print(err)
	}
	secondRequest, err := client.R().
		SetBody(`{"min_name_length": 10,"max_name_length": 10}`).
		SetHeaders(map[string]string{
		"Accept": "application/json, text/plain, */*",
	    "Referer": "https://temp-mail.io/en",
	    "Origin": "https://temp-mail.io",
	    "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.90 Safari/537.36",
	    "DNT": "1",
	    "Sec-Fetch-Mode": "cors",
	    "Content-Type": "application/json;charset=UTF-8",
		}).Post("https://api.internal.temp-mail.io/api/v2/email/new")
	var newEmail emailString
	json.Unmarshal(secondRequest.Body(), &newEmail)
	log.Print(user.details.email)
	user.details.email = newEmail.email
	return user.details.email






}


func (user *User) getVerifyString() string{
checkEmail := "https://api.internal.temp-mail.io/api/v2/email/replaceThis/messages"
fullCheckUrl := strings.Replace(checkEmail, "replaceThis", user.details.email, 1)
client := user.client
client.SetProxy(user.auth.proxy)

resp, err := client.R().
	Get(fullCheckUrl)
	if err != nil {
		log.Print(err)
	}
	parsedEmail, err := UnmarshalEmails(resp.Body())
	rxRelaxed := xurls.Relaxed()
	verifyUrl := rxRelaxed.FindString(parsedEmail[0].BodyText)
	justVerifyKey := strings.Replace(verifyUrl, "https://discordapp.com/verify?token=", "", 1)
	log.Print(justVerifyKey)
	return justVerifyKey

}

func (user *User) confirmEmail() {
	verifyString := user.getVerifyString()
	initialPayload := new(EmailVerify)
	initialPayload.Token = verifyString
	initialPayload.CAPTCHAKey = nil
	client := user.client
	client.SetProxy(user.auth.proxy)
	refferer := "https://discordapp.com/verify?token=" +verifyString
	initialMarshalled, err := initialPayload.Marshal()
	if err != nil {
		log.Print("error occurred")
		log.Print(err)
	}

		verifyNoCaptcha, err := client.R().
		SetHeaders(map[string]string{
		"authority": "discordapp.com",
		"pragma": "no-cache",
		"cache-control": "no-cache",
		"x-super-properties": user.auth.SuperProp,
		"x-fingerprint": user.auth.fingerprint,
		"accept-language": "en-US",
		"user-agent": user.auth.userAgent,
		"content-type": "application/json",
		"authorization": "undefined",
		"dnt": "1",
		"origin": "https://discordapp.com",
		"accept": "*/*",
		"sec-fetch-site": "same-origin",
		"sec-fetch-mode": "cors",
		"referer": refferer,
		"accept-encoding": "gzip, deflate, br",
	}).SetCookies(user.auth.cookies).
		SetBody(initialMarshalled).
		Post("https://discordapp.com/api/v6/auth/verify")
	if err != nil {
		log.Print("error occurred")
		log.Print(verifyNoCaptcha.Body())
		log.Print(err)
	}
time.Sleep(5*time.Second)
	payloadWithCaptcha := new(EmailVerify)
	captcha := getCaptcha()
	payloadWithCaptcha.CAPTCHAKey = captcha
	payloadWithCaptcha.Token = verifyString
	captchaMarshalled, err := payloadWithCaptcha.Marshal()
	if err != nil {
		log.Print("error occurred")
		log.Print(err)
	}
	verifyWithCaptcha, err := client.R().
		SetHeaders(map[string]string{
			"authority": "discordapp.com",
			"pragma": "no-cache",
			"cache-control": "no-cache",
			"x-super-properties": user.auth.SuperProp,
			"x-fingerprint": user.auth.fingerprint,
			"accept-language": "en-US",
			"user-agent": user.auth.userAgent,
			"content-type": "application/json",
			"authorization": "undefined",
			"dnt": "1",
			"origin": "https://discordapp.com",
			"accept": "*/*",
			"sec-fetch-site": "same-origin",
			"sec-fetch-mode": "cors",
			"referer": refferer,
			"accept-encoding": "gzip, deflate, br",
		}).SetCookies(user.auth.cookies).
		SetBody(captchaMarshalled).
		Post("https://discordapp.com/api/v6/auth/verify")
	if err != nil {
		log.Print("error occurred")
		log.Print(err)
	}
	token, err := UnmarshalVerifyResponse(verifyWithCaptcha.Body())

	if err != nil {
		log.Print("error occurred")
		log.Print(err)
	}
	//set user token
	user.auth.token	= token.Token





}
// Generated by https://quicktype.io

// RegPayload The payload used to register the account.
type RegPayload struct {
	Fingerprint   string      `json:"fingerprint"`
	Email         string      `json:"email"`
	Username      string      `json:"username"`
	Password      string      `json:"password"`
	Invite        interface{} `json:"invite"`
	Consent       bool        `json:"consent"`
	GiftCodeSkuID interface{} `json:"gift_code_sku_id"`
	CAPTCHAKey    string      `json:"captcha_key"`
}

type auth struct {
	fingerprint, cfuid, userAgent, token, proxy, SuperProp string
	cookies []*http.Cookie
	OpenMsg OpenMsg
}

type userDetails struct {
	username, password, email string
}

// User Struct that defines the user
type User struct {
	details userDetails
	auth    auth
	client  *resty.Client
}
