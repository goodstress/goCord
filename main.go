package main

import (
	"crypto/tls"
	b64 "encoding/base64"
	"github.com/bxcodec/faker/v3"
	"github.com/corpix/uarand"
	"github.com/go-resty/resty/v2"
	"github.com/mssola/user_agent"
	"github.com/thanhpk/randstr"
	"log"
	"math/rand"
	"mvdan.cc/xurls/v2"
	"net/http"
	//"strconv"
	"strings"
	"sync"
	"time"
	// "go.zoe.im/surferua"
	// 	    "github.com/mileusna/useragent"
	"github.com/tidwall/gjson"
)

var wg sync.WaitGroup

func main() {
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 100

	//for i := 1; i <= 10; i++ {
	//	fmt.Println(i)
	//}

	//create user
	log.Print("creating user")
	wg.Add(1)
	go createUser()
	wg.Wait()
	log.Println("terminated")
}

//noinspection SpellCheckingInspection
func (user *User) randomStickyIP() {
	rand.Seed(time.Now().UnixNano())

	//min := 10001
	//max := 29999
	//randomPort := rand.Intn(max - min + 1) + min
	//log.Print("random port: ", randomPort)
	randSession := randstr.String(16)
	ipString := "http://user-***REMOVED***-country-us-session-" + randSession + ":***REMOVED***" + "@gate.smartproxy.com:7000"

	user.auth.proxy = ipString
	user.auth.hostname = "gate.smartproxy.com:7000"
	user.auth.user = "user-***REMOVED***-country-us-session-" + randSession

}
func createUser() {
	log.Print("ran createUser")

	user := new(User)
	//set proxy
	user.randomStickyIP()

	log.Print("Proxy: ", user.auth.proxy)
	user.genUserAgent()
	//user.auth.userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36"
	//noinspection SpellCheckingInspection
	user.init()

	// userp := &user

}

func (user *User) init() {
	log.Print("Ran init")

	user.details.username = faker.Username()
	user.details.password = faker.Password()
	user.GrabCloudflare()

	user.GrabFingerprint()
	//s := spinner.New(spinner.CharSets[38], 100*time.Millisecond) // Build our new spinner
	//s.Prefix = "Waiting 90 seconds for fingerprint: "
	//s.Start() // Start the spinner
	//time.Sleep(90 * time.Second)
	//s.Stop()
	//setUsername
	//create superProp
	user.createSuperProp()
	//user.createXTrack()
	user.GenEmail()
	var registerWaitGroup sync.WaitGroup

	user.register(registerWaitGroup)
	registerWaitGroup.Wait()
	//var waitNoSmsGroup sync.WaitGroup
	var writeAccount sync.WaitGroup

	smsNeeded := make(chan string)

	log.Print("email process created")
	go user.openSocket(smsNeeded)
	writeAccount.Add(2)
	time.Sleep(20)
	go user.confirmEmail(writeAccount)
	//waitNoSmsGroup.Add(1)
	needSms := <-smsNeeded
	var checked bool

	checked = false
	log.Print("reached if statement")
	if (checked == false) && (needSms == "yes") {
		log.Print("ran SMS verification from IF statement")
		//waitNoSmsGroup.Done()
		writeAccount.Add(1)
		go user.smsVerification(writeAccount)
		//log.Print("need phone verification, continuing process.")
		checked = true
		//writeAccount.Done()

	}
	if needSms == "verified" {
		log.Print("detected verified message")
		log.Print("writing account")
		go user.writeAccount()
		log.Print("finished writing account")
	}
	log.Print("after if statement")
	log.Print("socket go process created")
	//time.Sleep(10*time.Second)

	//
	//log.Print("ran confirm email")
	//var smsConfirmed sync.WaitGroup
	//
	//smsConfirmed.Wait()
	//log.Print("sms done, running email.")
	time.Sleep(2 * time.Second)

	log.Print(user.auth.token)
	writeAccount.Wait()
	log.Print("complete")
	log.Print(user.auth.token)

	//user.writeAccount()

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
	agent := uarand.GetRandom()
	//agent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36"
	user.auth.userAgent = agent
	log.Println("Set useragent to: ", agent)
}

// func superProp(agent string) string {
// 	ua.Parse(agent)
// }
// func xTrack() string {

// }

func (user *User) register(complete sync.WaitGroup) {
	defer complete.Done()
	complete.Add(1)
	captcha := user.NewKey()
	log.Print("captcha: ", captcha)
	realRegister := RegPayload{Fingerprint: user.auth.fingerprint, Email: user.details.email, Username: user.details.username, Password: user.details.password, Invite: nil, Consent: true, GiftCodeSkuID: nil, CAPTCHAKey: captcha}
	registerURL := "https://discordapp.com/api/v6/auth/register"
	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetProxy(user.auth.proxy)
	//todo: handle incorrect captcha
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
	log.Print(resp.String())
	token := gjson.Get(resp.String(), "token").String()
	user.auth.token = token
	log.Println("set token in user to: ", token)
	//complete.Done()

}

//type emailString struct {
//	email string
//}
func (user *User) GenEmail() string {
	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetProxy(user.auth.proxy)
	_, err := client.R().
		SetBody(`{"min_name_length": 10,"max_name_length": 10}`).
		SetHeaders(map[string]string{
			"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:69.0) Gecko/20100101 Firefox/69.0",
			"Accept":          "application/json, text/plain, */*",
			"Accept-Language": "en-US,en;q=0.5",
			"Content-Type":    "application/json;charset=utf-8",
			"Origin":          "https://temp-mail.io",
			"Connection":      "keep-alive",
			"Referer":         "https://temp-mail.io/en",
			"TE":              "Trailers",
		}).Options("https://api.internal.temp-mail.io/api/v2/email/new")
	if err != nil {
		log.Print(err)
	}
	secondRequest, err := client.R().
		SetBody(`{"min_name_length": 10,"max_name_length": 10}`).
		SetHeaders(map[string]string{
			"Accept":         "application/json, text/plain, */*",
			"Referer":        "https://temp-mail.io/en",
			"Origin":         "https://temp-mail.io",
			"User-Agent":     "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.90 Safari/537.36",
			"DNT":            "1",
			"Sec-Fetch-Mode": "cors",
			"Content-Type":   "application/json;charset=UTF-8",
		}).Post("https://api.internal.temp-mail.io/api/v2/email/new")
	if err != nil {
		log.Print(err)
	}
	log.Print(secondRequest.String())
	user.details.email = gjson.Get(secondRequest.String(), "email").String()
	log.Print(user.details.email)

	return user.details.email

}

func (user *User) getVerifyString() string {
	checkEmail := "https://api.internal.temp-mail.io/api/v2/email/replaceThis/messages"
	fullCheckUrl := strings.Replace(checkEmail, "replaceThis", user.details.email, 1)
	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetProxy(user.auth.proxy)
	time.Sleep(4 * time.Second)
	resp, err := client.R().
		Get(fullCheckUrl)
	if err != nil {
		log.Print(err)
	}
	parsedEmail, err := UnmarshalEmails(resp.Body())
	rxRelaxed := xurls.Relaxed()
	verifyUrl := rxRelaxed.FindString(parsedEmail[0].BodyText)
	//resp, err := client.R().
	//	Get(verifyUrl)
	//if err != nil {
	//	log.Print("error getting real VerifyURL: ", err)
	//}
	realVerifyResponse, getVerifyURLError := client.R().
		Get(verifyUrl)
	if getVerifyURLError != nil {
		log.Fatalf("error getting redirect: => %v", getVerifyURLError)
	}

	// Your magic function. The Request in the Response is the last URL the
	// client tried to access.
	realVerifyUrl := realVerifyResponse.RawResponse.Request.URL.String()

	justVerifyKey := strings.Replace(realVerifyUrl, "https://discordapp.com/verify?token=", "", 1)
	log.Print("email verification key: ", justVerifyKey)
	return justVerifyKey

}

func (user *User) confirmEmail(confirmedWait sync.WaitGroup) {
	defer confirmedWait.Done()
	//confirmedWait.Add(1)
	time.Sleep(20 * time.Second)
	verifyString := user.getVerifyString()
	initialPayload := new(EmailVerify)
	initialPayload.Token = verifyString
	initialPayload.CAPTCHAKey = nil
	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetProxy(user.auth.proxy)
	referrer := "https://discordapp.com/verify?token=" + verifyString
	//initialMarshalled, err := initialPayload.Marshal()
	//if err != nil {
	//	log.Print("error occurred")
	//	log.Print(err)
	//}
	//
	//	verifyNoCaptcha, err := client.R().
	//	SetHeaders(map[string]string{
	//	"authority":          "discordapp.com",
	//	"pragma":             "no-cache",
	//	"cache-control":      "no-cache",
	//	"x-super-properties": user.auth.SuperProp,
	//	"x-fingerprint":      user.auth.fingerprint,
	//	"accept-language":    "en-US",
	//	"user-agent":         user.auth.userAgent,
	//	"content-type":       "application/json",
	//	"authorization":      "undefined",
	//	"dnt":                "1",
	//	"origin":             "https://discordapp.com",
	//	"accept":             "*/*",
	//	"sec-fetch-site":     "same-origin",
	//	"sec-fetch-mode":     "cors",
	//	"referer":            referrer,
	//	"accept-encoding":    "gzip, deflate, br",
	//}).SetCookies(user.auth.cookies).
	//	SetBody(initialMarshalled).
	//	Post("https://discordapp.com/api/v6/auth/verify")
	//if err != nil {
	//	log.Print("error occurred")
	//	log.Print(verifyNoCaptcha.String())
	//	log.Print(err)
	//}
	time.Sleep(5 * time.Second)
	payloadWithCaptcha := new(EmailVerify)
	captcha := user.NewKey()
	payloadWithCaptcha.CAPTCHAKey = captcha
	payloadWithCaptcha.Token = verifyString
	captchaMarshalled, err := payloadWithCaptcha.Marshal()
	if err != nil {
		log.Print("error occurred")
		log.Print(err)
	}
	verifyWithCaptcha, err := client.R().
		SetHeaders(map[string]string{
			"authority":          "discordapp.com",
			"pragma":             "no-cache",
			"cache-control":      "no-cache",
			"x-super-properties": user.auth.SuperProp,
			"x-fingerprint":      user.auth.fingerprint,
			"accept-language":    "en-US",
			"user-agent":         user.auth.userAgent,
			"content-type":       "application/json",
			"authorization":      "undefined",
			"dnt":                "1",
			"origin":             "https://discordapp.com",
			"accept":             "*/*",
			"sec-fetch-site":     "same-origin",
			"sec-fetch-mode":     "cors",
			"referer":            referrer,
			"accept-encoding":    "gzip, deflate, br",
		}).SetCookies(user.auth.cookies).
		SetBody(captchaMarshalled).
		Post("https://discordapp.com/api/v6/auth/verify")
	if err != nil {
		log.Print("error occurred")
		log.Print(err)
	}

	//set user token
	user.auth.token = gjson.Get(verifyWithCaptcha.String(), "token").String()
	log.Print(verifyWithCaptcha.String())

	//confirmedWait.Done()

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
	fingerprint, cfuid, userAgent, token, proxy, SuperProp, hostname, user string
	cookies                                                                []*http.Cookie
	OpenMsg                                                                []byte
}
type PhoneNumber struct {
	phoneNumber, numberId string
}
type userDetails struct {
	username, password, email string
}

// User Struct that defines the user
type User struct {
	details     userDetails
	auth        auth
	PhoneNumber PhoneNumber
	smsApi      smsApi
}
