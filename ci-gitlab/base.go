/*
Package cigitlab is layer for interaction with Gitlab CI API
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
	// CmdTest command to trigger test action
	CmdTest = "test"
	// CmdDeploy command to trigger deploy action
	CmdDeploy = "deploy"
)

var (
	// domain
	apiHost = ""
	// part or url
	apiURL   = ""
	apiToken = ""
	once     sync.Once
	// project ids in gitlab
	// @TODO make configurable map created from conf/ENV
	projIDs = map[string]string{
		"slackbot": "5",
	}
	// ErrProjID if there is no association for between project name and id
	ErrProjID = errors.New("project id unknow")
	// ErrReq if post request to CI server return not expected status code
	ErrReq = errors.New("wrong request status from CI gitlab")
	// ErrWrongCMD undefined command provided
	ErrWrongCMD = errors.New("undefined command provided")
)

// Resp is esponse object from ci server
// currently only commit unmarshaled from response
type Resp struct {
	Commit `json:"commit"`
}

// Commit object from response object
type Commit struct {
	Ref    string
	Sha    string
	ProjID int    `json:"project_id"`
	GitMsg string `json:"git_commit_message"`
}

// Msg returns formated string from resp
// @TODO better to remove from here
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

// Configure set host, api urls and token of CI
//@TODO start goroute per project
func Configure(host, api, token string) {
	// protect from race
	once.Do(func() {
		apiHost = host
		apiURL = api
		apiToken = token
	})
}

// Trigger particular api in gitlab ci for particular action
func Trigger(cmd, proj, ref string) (resp Resp, err error) {
	switch cmd {
	// react to test command
	case CmdTest:
		resp, err = req(proj, ref)
		if err != nil {
			return
		}
	}

	return resp, ErrWrongCMD
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
	fullURL := apiURL + "/projects/" + id + "/refs/" + ref + "/trigger"
	v := url.Values{}
	v.Set("token", apiToken)
	resp, err := http.PostForm(fullURL, v)
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
