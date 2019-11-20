package main

import (
	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)
//create dialer

func genOpenMsg() {

}

func (user *User) openSocket(waitNoSmsGroup sync.WaitGroup, smsNeeded chan string)  {
	_, _ = user.CreateOpenMsg()

	var dialer = websocket.Dialer{
		Proxy: http.ProxyURL(&url.URL{
			Scheme: "http", // or "https" depending on your proxy
			Host: user.auth.proxy ,
			Path: "/",
		}),
	}
	waitNoSmsGroup.Add(1)
	c, _, err := dialer.Dial("wss://gateway.discord.gg/?encoding=json&v=6",nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	defer c.Close()

	done := make(chan struct{})
	err = c.WriteJSON(user.auth.OpenMsg)
	if err != nil {
		log.Print("error occurred: ", err)
	}
	smsTicker := time.NewTicker(20*time.Second)
	for {
		select {
		case <-smsTicker.C:
		smsNeeded<-"not"
		}
	}
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			requiredAction := gjson.Get(string(message), "d.required_action").String()
			log.Print("requiredAction: ", requiredAction)
			if requiredAction == "REQUIRE_VERIFIED_PHONE" {
			log.Print("Error phone required")
			c.Close()
				waitNoSmsGroup.Done()


			}
			log.Printf("recv: %s", message)
		}
	}()

	if err != nil {
		log.Print("error occured::::")
		log.Print(err)
	}



}
