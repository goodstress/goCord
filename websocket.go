package main

import (
	"crypto/tls"
	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
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
	log.Print("open msg complete")
	var dialer = websocket.Dialer{

		Proxy: http.ProxyURL(&url.URL{

			Scheme: "http", // or "https" depending on your proxy,
			User: url.UserPassword("***REMOVED***", "***REMOVED***"),
			Host: user.auth.hostname+user.auth.port,
		}),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		HandshakeTimeout: 30*time.Second,
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
	}
	log.Print("setup dialer")
	c, resp, err := dialer.Dial("wss://gateway.discord.gg/?encoding=json&v=6",nil)
	if err == websocket.ErrBadHandshake {
		log.Printf("handshake failed with status %d", resp.StatusCode)
	}
	log.Print("dialed")
	if err != nil {
		log.Print(err)
		log.Fatal("dial:", err)
	}

	defer c.Close()

	done := make(chan struct{})
	_ = c.WriteMessage(websocket.TextMessage, user.auth.OpenMsg)
log.Print("sent open msg")

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
			requiredAction := gjson.Get(string(message), "d.required_action").String()
			log.Print("requiredAction: ", requiredAction)
			if requiredAction == "REQUIRE_VERIFIED_PHONE" {
			log.Print("Error phone required")
			c.Close()


			}
			log.Printf("recv: %s", message)
		}
	}()

	if err != nil {
		log.Print("error occured::::")
		log.Print(err)
	}



}
