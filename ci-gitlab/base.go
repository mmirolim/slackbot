package cigitlab

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
)

const (
	TEST   = "test"
	DEPLOY = "deploy"
)

var (
	apiUrl   = ""
	apiToken = ""
	once     sync.Once
	// project ids in gitlab
	// @TODO make configurable map created from conf/ENV
	projIDs = map[string]string{
		"slack-bot": "5",
	}
	// Errors
	ErrProjID = errors.New("project id unknow")
	ErrReq    = errors.New("request status is not 200")
)

type Resp struct {
	Commit
}

type Commit struct {
	Ref    string
	Sha    string
	ProjID string `json:"project_id"`
	GitMsg string `json:"git_commit_message"`
}

func (c *Commit) Stringer() string {
	return fmt.Sprintf("project id: %s commit msg: %s ref: %s sha: %s", c.ProjID, c.GitMsg, c.Sha)
}

//@TODO start goroute per project
func Configure(url, token string) {
	// protect from race
	once.Do(func() {
		apiUrl = url
		apiToken = token
	})
}

// trigger particular api in gitlab ci for particular action
func Trigger(cmd, proj, ref string) (resp Resp, err error) {
	switch cmd {
	// react to test command
	case TEST:
		resp, err = req(proj, ref)
		if err != nil {
			return
		}
	}

	return
}

// construct correct url for gitlab ci trigger
// and make http post request with project id, ref and token
func req(proj, ref string) (Resp, error) {
	var ciResp Resp
	// check if we have project ids
	id, ok := projIDs[proj]
	if !ok {
		return ciResp, ErrProjID
	}
	// construct project and branch specific url
	fullUrl := apiUrl + "/projects/" + id + "/refs/" + ref + "/trigger"
	// send post form with token
	v := url.Values{}
	v.Set("token", apiToken)
	log.Println("values", v, "fullurl", fullUrl)
	resp, err := http.PostForm(fullUrl, v)

	//	if resp.StatusCode != http.StatusOK {
	//		return ciResp, ErrReq
	//	}

	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		return ciResp, err
	}
	// unmarshal object
	err = json.Unmarshal(data, &ciResp)
	log.Printf("resp %#v, body %#v", resp, ciResp)
	return ciResp, err
}
