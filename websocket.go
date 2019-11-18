package main

import (
	"github.com/gorilla/websocket"
	"net/http"
)
//create dialer
var dialer = websocket.Dialer{
	Proxy: http.ProxyFromEnvironment,
}
func genOpenMsg() {

}

func openSocket(ch chan)  {

}
