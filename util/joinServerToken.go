package util

import (
	"github.com/go-resty/resty/v2"
	"log"
)

func JoinServerToken(inviteCode string, auth string) {
	client := resty.New()
	serverInviteURL := "https://discordapp.com/api/v6/invites/" + inviteCode
	client.SetProxy("http://user-***REMOVED***" + ":***REMOVED***" + "@us.smartproxy.com:10000")
	client.SetHeader("authorization", auth)
	client.SetHeader("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36")
	resp, err := client.R().Post(serverInviteURL)
	if err != nil {
		log.Print("error posting invite: ", err)

	}
	log.Print(resp.StatusCode())
	log.Print(resp.String())
	//log.Print("times attempted: ", attemptedJoins.Value())
	if resp.StatusCode() == 200 {
		//successfulJoins.Add(1)
		log.Print("successfully joined server, total joins is: ")
	}
}
