/*
Slack bot for automatization purposes
*/
package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	ci "xr/slackbot/ci-gitlab"
	"xr/slackbot/slack"
)

const (
	SLACK_BOT_TOKEN  = "SLACK_BOT_TOKEN"
	CI_TRIGGER_TOKEN = "CI_TRIGGER_TOKEN"
)

var (
	BuildVersion = ""
)

func main() {
	// connect to slack via websocket with auth token
	ws, id, err := slack.Connect(os.Getenv(SLACK_BOT_TOKEN))
	if err != nil {
		// @TODO enable retries or cb in slack package
		log.Fatal("could not init slack connection", err)
	}

	fmt.Printf("bot ready, build version: %s\n", BuildVersion)
	// configure ci gitlab conf
	ci.Configure("https://gitlab-ci.regium.com", "https://gitlab-ci.regium.com/api/v1", os.Getenv(CI_TRIGGER_TOKEN))
	// start listening to messages
	for {
		// read each incoming message
		m, err := slack.GetMsg(ws)
		if err != nil {
			log.Println("get msg error ", err)
		}

		// see if we're mentioned
		if m.Type == "message" && strings.HasPrefix(m.Text, "<@"+id+">") {
			// if so try to parse if
			parts := strings.Fields(m.Text)
			if len(parts) == 4 && parts[1] == ci.CMDTEST {
				// looks good, get the quote and reply with the result
				go func(m slack.Msg) {
					// trigger test action
					resp, err := ci.Trigger(ci.CMDTEST, parts[2], parts[3])
					if err != nil {
						// set error as msg
						m.Text = "error happened " + err.Error()
						log.Println("could not trigger ci work", err)
					} else {
						// format response
						m.Text = resp.Msg()
					}
					// send msg to slack
					if err := slack.PostMsg(ws, m); err != nil {
						log.Println("could not send msg", err)
					}
				}(m)
				// NOTE: the Message object is copied, this is intentional
			} else {
				// huh?
				m.Text = fmt.Sprintf("cmd does not recognized\n")
				if err := slack.PostMsg(ws, m); err != nil {
					log.Println("could not send msg", err)
				}
			}
		}
	}
}
