package main

import (
	"crypto/tls"
	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	"time"

	//"github.com/tidwall/gjson"
	"log"
	"net/http"
	"net/url"
	//"time"
)

//create dialer

//func genOpenMsg() {
//
//}

func (user *User) openSocket(smsNeeded chan string) {
	log.Print("in websocket function")
	user.CreateOpenMsg()
	log.Print("open msg created")
	host := user.auth.hostname
	log.Print(host)
	var proxyDialer = websocket.Dialer{
		Proxy: http.ProxyURL(&url.URL{

			Scheme: "http", // or "https" depending on your proxy,
			User:   url.UserPassword(user.auth.user, "***REMOVED***"),
			Host:   host,
		}),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	log.Print("setup dialer")
	c, _, err := proxyDialer.Dial("wss://gateway.discord.gg/?encoding=json&v=6", nil)
	if err != nil {
		log.Print("dial err: ", err)
	}
	log.Print("dialed")

	err = c.WriteMessage(websocket.TextMessage, user.auth.OpenMsg)
	if err != nil {
		log.Print("write error: ", err)
	}
	log.Print("sent open msg")

	err = c.WriteMessage(websocket.TextMessage, []byte(`{"op":4,"d":{"guild_id":null,"channel_id":null,"self_mute":true,"self_deaf":false,"self_video":false}}`))
	if err != nil {
		log.Print("secondary open msg error: ", err)
	}
	log.Print("sent secondary open msg")

	//defer c.Close()

	done := make(chan struct{})
	//_ = c.WriteMessage(websocket.TextMessage, user.auth.OpenMsg)

	heartbeat := time.NewTicker(39000 * time.Millisecond)

	go func() {
		defer close(done)
		//smsTicker := time.NewTicker(20*time.Second)
		//sendOpen := time.NewTicker(5*time.Second)
		for {
			//select {
			//case <-smsTicker.C:
			//	smsNeeded<-"not"
			//
			//}
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
			}
			log.Printf("recv: %s", string(message))

			requiredAction := gjson.Get(string(message), "d.required_action").String()

			log.Print("requiredAction: ", requiredAction)
			if requiredAction == "REQUIRE_VERIFIED_PHONE" {
				log.Print("phone required")
				smsNeeded <- "yes"
				log.Print("sent yes to channel")

				//c.Close()

			}
			verifiedTrue := gjson.Get(string(message), "d.verified").String()
			log.Print("verifiedTrue: ", verifiedTrue)
			if verifiedTrue == "true" {
				log.Print("verified is needed")
				smsNeeded <- "verified"
			}
			//log.Printf("recv: %s", message)

		}

	}()

	go func() {
		defer close(done)
		//smsTicker := time.NewTicker(20*time.Second)
		//sendOpen := time.NewTicker(5*time.Second)
		for {
			//select {
			//case <-smsTicker.C:
			//	smsNeeded<-"not"
			//
			//}
			select {
			case <-heartbeat.C:
				err = c.WriteMessage(websocket.TextMessage, []byte(`{"op":1,"d":37}`))
				if err != nil {
					log.Print("writeHeartBeat error: ", err)
				}
			}
			//requiredAction := gjson.Get(string(message), "d.required_action").String()
			//log.Print("requiredAction: ", requiredAction)
			//if requiredAction == "REQUIRE_VERIFIED_PHONE" {
			//log.Print("Error phone required")
			//c.Close()
			//
			//
			//}
			//log.Printf("recv: %s", message)

		}

	}()

}
