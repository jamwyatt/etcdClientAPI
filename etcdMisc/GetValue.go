package etcdMisc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

//
// GetValue  - Function to get the value of a pre-existing key
//
// 	conn		ectdConnection, made with etcdMisc.MakeEtcdConnection()
// 	port		port to connect to
// 	key		etcd node key/directory
// 	recurse		Recursive get, useful on directories, but the response is recursive too
// 	sort		apply etcd sorting?
//
//
func GetValue(conn etcdConnection, key string, recurse bool, sort bool) (EtcdResponse, error) {

	var err error
	url := fmt.Sprintf("%s://%s:%d/v2/keys%s?recursive=%t&sorted=%t", conn.Proto, conn.Host, conn.Port, key, recurse, sort)
	var request *http.Request
	request, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return EtcdResponse{}, errors.New("http.NewRequest: " + err.Error())
	}

	var resp *http.Response
	resp, err = conn.Client.Do(request)
	if err != nil {
		return EtcdResponse{}, errors.New("http.client.Do: " + err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return EtcdResponse{}, errors.New("ioutil.ReadAll: " + err.Error())
	}
	r := EtcdResponse{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return EtcdResponse{}, errors.New("json.Unmarshal: " + err.Error())
	}
	// Check for etcd error
	if r.Cause != "" {
		return r, errors.New("GET etcd error: " + r.Message)
	}

	return r, nil // All good
}
