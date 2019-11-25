package main

import (
	"crypto/tls"
	"github.com/gorilla/websocket"
	//"github.com/tidwall/gjson"
	"log"
	"net/http"
	"net/url"
	"time"
)
//create dialer

func genOpenMsg() {

}

func (user *User) openSocket( smsNeeded chan string)  {
	log.Print("in websocket function")
	user.CreateOpenMsg()
	host := user.auth.hostname + user.auth.port
	log.Print(host)
	log.Print("open msg complete")
	var dialer = websocket.Dialer{

		Proxy: http.ProxyURL(&url.URL{

			Scheme: "http", // or "https" depending on your proxy,
			User: url.UserPassword("***REMOVED***", "***REMOVED***"),
			Host: host,
		}),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		HandshakeTimeout: 30*time.Second,
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
	}
	log.Print("setup dialer")
	c, _, err := dialer.Dial("wss://gateway.discord.gg/?encoding=json&v=6",nil)
	if err != nil {
		log.Print("dial err: ", err)
	}
	log.Print("dialed")

	err = c.WriteMessage(websocket.TextMessage, user.auth.OpenMsg)
	if err != nil {
		log.Print("write error: ", err)
	}
	log.Print("sent open msg")


	defer c.Close()

	done := make(chan struct{})
	//_ = c.WriteMessage(websocket.TextMessage, user.auth.OpenMsg)


	go func() {
		defer close(done)
		smsTicker := time.NewTicker(20*time.Second)
		for {
			select {
			case <-smsTicker.C:
				smsNeeded<-"not"
			}
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
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

	if err != nil {
		log.Print("error occured::::")
		log.Print(err)
	}



}
