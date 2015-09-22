/*
Slack sevice communication package
auth, connecting, sending and receiving
*/

package slack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync/atomic"

	"golang.org/x/net/websocket"
)

const (
	SLACK_API_URL     = "https://api.slack.com/"
	SLACK_API_RTM_URL = "https://slack.com/api/rtm.start"
)

// These two structures represent the response of the Slack API rtm.start.
// Only some fields are included. The rest are ignored by json.Unmarshal.

type RespRtmStart struct {
	Ok    bool     `json:"ok"`
	Error string   `json:"error"`
	Url   string   `json:"url"`
	Self  RespSelf `json:"self"`
}

type RespSelf struct {
	Id string `json:"id"`
}

// start does a rtm.start, and returns a websocket URL and user ID. The
// websocket URL can be used to initiate an RTM session.
func start(token string) (wsurl, id string, err error) {
	url := fmt.Sprintf(SLACK_API_RTM_URL+"?token=%s", token)
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("API request failed with code %d", resp.StatusCode)
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
		err = fmt.Errorf("Slack error: %s", respObj.Error)
		return
	}

	wsurl = respObj.Url
	id = respObj.Self.Id
	return
}

// These are the messages read off and written into the websocket. Since this
// struct serves as both read and write, we include the "Id" field which is
// required only for writing.

type Msg struct {
	Id      uint64 `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

// get slack msg object unmarshaled from active websocket connection or return error
func GetMsg(ws *websocket.Conn) (m Msg, err error) {
	err = websocket.JSON.Receive(ws, &m)
	return
}

// counter will increment msg id
// should be protected by concurrent use
// with atomic or mutex
var counter uint64

// send slack message object to active websocket connection or return error
func PostMsg(ws *websocket.Conn, m Msg) error {
	m.Id = atomic.AddUint64(&counter, 1)
	return websocket.JSON.Send(ws, m)
}

// Starts a websocket-based Real Time API session and return the websocket
// and the ID of the (bot-)user whom the token belongs to.
func Connect(token string) (*websocket.Conn, string, error) {
	wsurl, id, err := start(token)
	if err != nil {
		return nil, "", err
	}

	ws, err := websocket.Dial(wsurl, "", SLACK_API_URL)
	if err != nil {
		return nil, "", err
	}

	return ws, id, err
}
