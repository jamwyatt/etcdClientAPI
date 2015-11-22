package etcdMisc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

//
//  DeleteDir - Function to delete a directory (optional recursion)
//
// 	client		http.Client that can control functionality, like Timeouts (nil is ok)
// 	tr		http.Transport that can set TLS client attributes (nil is ok)
// 	proto		"http" or "https"
// 	host		host to connect with
// 	port		port to connect to
// 	key		etcd node key/directory
//	recursive	true for recursive delete of everything
//
func DeleteDir(client *http.Client, tr *http.Transport,
	proto string, host string, port int,
	key string, recursive ...bool) (EtcdResponse, error) {

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

	urlStr := fmt.Sprintf("%s://%s:%d/v2/keys%s?dir=true", proto, host, port, key)
	if len(recursive) > 0 && recursive[0] {
		urlStr += "&recursive=true"
	}

	var request *http.Request
	request, err = http.NewRequest("DELETE", urlStr, nil)
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
		return EtcdResponse{}, errors.New("DeleteDir etcd error: " + r.Message)
	}
	return r, nil // All good
}
