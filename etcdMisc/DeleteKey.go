package etcdMisc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

//
// DeleteKey  - Function to delete a pre-existing key
//
// 	client		http.Client that can control functionality, like Timeouts (nil is ok)
// 	tr		http.Transport that can set TLS client attributes (nil is ok)
// 	proto		"http" or "https"
// 	host		host to connect with
// 	port		port to connect to
// 	key		etcd node key/directory
//
//
func DeleteKey(client *http.Client, tr *http.Transport, proto string, host string, port int, key string) (EtcdResponse, error) {

	var err error
	if client == nil {
		client = &http.Client{
			Timeout: 0,
		}
	}
	if tr == nil {
		tr = &http.Transport{}
		client.Transport = tr
	}

	url := fmt.Sprintf("%s://%s:%d/v2/keys%s", proto, host, port, key)
	var request *http.Request
	request, err = http.NewRequest("DELETE", url, nil)
	if err != nil {
		return EtcdResponse{}, errors.New("http.NewRequest: " + err.Error())
	}

	var resp *http.Response
	resp, err = client.Do(request)
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
		return EtcdResponse{}, errors.New("etcd error: " + r.Message)
	}
	return r, nil // All good
}
