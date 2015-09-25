/*
Slack sevice communication package
auth, connecting, sending and receiving
*/

package slack

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync/atomic"

	"golang.org/x/net/websocket"
)

const (
	slackAPIURL    = "https://api.slack.com/"
	slackAPIRtmURL = "https://slack.com/api/rtm.start"
)

var (
	// ErrReqFail if request status is not 200
	ErrReqFail = errors.New("request to slack not returned 200")
	// ErrReqErr if request object has error
	ErrReqErr = errors.New("request object return error")
)

// RespRtmStart and RespSelf these two structures represent the response of the Slack API rtm.start.
// Only some fields are included. The rest are ignored by json.Unmarshal.
type RespRtmStart struct {
	Ok    bool     `json:"ok"`
	Error string   `json:"error"`
	URL   string   `json:"url"`
	Self  RespSelf `json:"self"`
}

// RespSelf object part of response from slack api
type RespSelf struct {
	ID string `json:"id"`
}

// start does a rtm.start, and returns a websocket URL and user ID. The
// websocket URL can be used to initiate an RTM session.
func start(token string) (wsurl, id string, err error) {
	url := fmt.Sprintf(slackAPIRtmURL+"?token=%s", token)
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = ErrReqFail
		log.Println(resp.StatusCode)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var respObj RespRtmStart
	err = json.Unmarshal(body, &respObj)
	if err != nil {
		return
	}

	if !respObj.Ok {
		err = ErrReqErr
		log.Println(respObj.Error)
		return
	}

	wsurl = respObj.URL
	id = respObj.Self.ID
	return
}

// Msg these are the messages read off and written into the websocket. Since this
// struct serves as both read and write, we include the "Id" field which is
// required only for writing.
type Msg struct {
	ID      uint64 `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

// GetMsg get slack msg object unmarshaled from active websocket connection or return error
func GetMsg(ws *websocket.Conn) (m Msg, err error) {
	err = websocket.JSON.Receive(ws, &m)
	return
}

// counter will increment msg id
// should be protected by concurrent use
// with atomic or mutex
var counter uint64

// PostMsg send slack message object to active websocket connection or return error
func PostMsg(ws *websocket.Conn, m Msg) error {
	m.ID = atomic.AddUint64(&counter, 1)
	return websocket.JSON.Send(ws, m)
}

// Connect establish websocket-based Real Time API session and return the websocket
// and the ID of the (bot-)user whom the token belongs to.
func Connect(token string) (ws *websocket.Conn, id string, err error) {
	wsurl, id, err := start(token)
	if err != nil {
		return
	}
	ws, err = websocket.Dial(wsurl, "", slackAPIURL)
	if err != nil {
		return
	}

	return
}
