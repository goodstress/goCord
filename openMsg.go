package main

import (
	"encoding/json"
	//"errors"
	"github.com/mssola/user_agent"
	//"go/types"
	"log"
)

func UnmarshalOpenMsg(data []byte) (OpenMsg, error) {
	var r OpenMsg
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *OpenMsg) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
func (user *User) CreateOpenMsg() {
	log.Print("creating opening msg")
	ua := user_agent.New(user.auth.userAgent)
	defaultString := `{"op":2,"d":{"token":"yourTokenHere","properties":{"os":"Windows","browser":"Chrome","device":"","browser_user_agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.97 Safari/537.36","browser_version":"78.0.3904.97","os_version":"10","referrer":"https://discordapp.com/widget?id=518014517138030603&theme=dark","referring_domain":"discordapp.com","utm_source":"Discord Widget","utm_medium":"Logo","referrer_current":"","referring_domain_current":"","release_channel":"stable","client_build_number":49868,"client_event_source":null},"presence":{"status":"online","since":0,"activities":[],"afk":false},"compress":false}}`
	msg, _ := UnmarshalOpenMsg([]byte(defaultString))
	msg.D.Token = user.auth.token
	log.Print("set token")
	msg.D.Properties.BrowserUserAgent = user.auth.userAgent
	msg.D.Properties.Browser, msg.D.Properties.BrowserVersion = ua.Browser()

	msg.D.Properties.OS = ua.OSInfo().Name
	msg.D.Properties.OSVersion = ua.OSInfo().Version
	msg.D.Properties.ClientEventSource = "null"

	msg.D.Properties.Browser, msg.D.Properties.BrowserVersion = ua.Browser()
	user.auth.OpenMsg, _ = msg.Marshal()
	//if err != nil {
	//	log.Print("marshal error: ", err)
	//}
	//user.Auth.OpenMsg = []byte(string(user.Auth.OpenMsg) + "}")
	log.Print("marshalled msg")
	log.Print(string(user.auth.OpenMsg))
	//return msg.Marshal()
}

type OpenMsg struct {
	Op int64 `json:"op"`
	D  D     `json:"d"`
}

type D struct {
	Token      string     `json:"token"`
	Properties Properties `json:"properties"`
	Presence   Presence   `json:"presence"`
	Compress   bool       `json:"compress"`
}

type Presence struct {
	Status     string        `json:"status"`
	Since      int64         `json:"since"`
	Activities []interface{} `json:"activities"`
	Afk        bool          `json:"afk"`
}

type Properties struct {
	OS                     string      `json:"os"`
	Browser                string      `json:"browser"`
	Device                 string      `json:"device"`
	BrowserUserAgent       string      `json:"browser_user_agent"`
	BrowserVersion         string      `json:"browser_version"`
	OSVersion              string      `json:"os_version"`
	Referrer               string      `json:"referrer"`
	ReferringDomain        string      `json:"referring_domain"`
	UtmSource              string      `json:"utm_source"`
	UtmMedium              string      `json:"utm_medium"`
	ReferrerCurrent        string      `json:"referrer_current"`
	ReferringDomainCurrent string      `json:"referring_domain_current"`
	ReleaseChannel         string      `json:"release_channel"`
	ClientBuildNumber      int64       `json:"client_build_number"`
	ClientEventSource      interface{} `json:"client_event_source"`
}
