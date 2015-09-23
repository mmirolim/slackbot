/*
Gitlab CI integration package
*/
package cigitlab

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

const (
	// cmd for ci
	CMDTEST   = "test"
	CMDDEPLOY = "deploy"
)

var (
	// domain
	apiHost = ""
	// part or url
	apiUrl   = ""
	apiToken = ""
	once     sync.Once
	// project ids in gitlab
	// @TODO make configurable map created from conf/ENV
	projIDs = map[string]string{
		"slackbot": "5",
	}
	// Errors
	ErrProjID = errors.New("project id unknow")
	ErrReq    = errors.New("request status is not 200")
)

type Resp struct {
	Commit `json:"commit"`
}

type Commit struct {
	Ref    string
	Sha    string
	ProjID int    `json:"project_id"`
	GitMsg string `json:"git_commit_message"`
}

// implement Stringer interface
func (r *Resp) Msg() string {
	var obj string
	// check if not empty
	if len(r.Sha) < 7 {
		return obj
	}
	pid := strconv.Itoa(r.ProjID)

	// ci specific
	link := apiHost + "/projects/" + pid + "/refs/" + r.Ref + "/commits/" + r.Sha
	return fmt.Sprintf("project-id: %d commit-msg: %s ref: %s %s", r.ProjID, r.GitMsg, r.Ref, link)
}

//@TODO start goroute per project
func Configure(host, api, token string) {
	// protect from race
	once.Do(func() {
		apiHost = host
		apiUrl = api
		apiToken = token
	})
}

// trigger particular api in gitlab ci for particular action
func Trigger(cmd, proj, ref string) (resp Resp, err error) {
	switch cmd {
	// react to test command
	case CMDTEST:
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
	v := url.Values{}
	v.Set("token", "2df1de069095cfda2edde54d57ebbe")
	resp, err := http.PostForm(fullUrl, v)
	// if ok http result is 201
	if resp.StatusCode != http.StatusCreated {
		return ciResp, ErrReq
	}

	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		return ciResp, err
	}
	// unmarshal object
	err = json.Unmarshal(data, &ciResp)
	return ciResp, err
}
