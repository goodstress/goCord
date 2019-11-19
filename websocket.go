package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
)
//create dialer

func genOpenMsg() {

}

func (user *User) openSocket()  {
	user.CreateOpenMsg()

	var dialer = websocket.Dialer{
		Proxy: http.ProxyURL(&url.URL{
			Scheme: "http", // or "https" depending on your proxy
			Host: user.auth.proxy ,
			Path: "/",
		}),
	}
	c, _, err := dialer.Dial("wss://gateway.discord.gg/?encoding=json&v=6",nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	defer c.Close()

	done := make(chan struct{})
	c.WriteJSON(user.auth.OpenMsg)

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	if err != nil {
		log.Print("error occured::::")
		log.Print(err)
	}



}
