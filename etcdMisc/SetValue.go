package etcdMisc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

//
//  SetValue - Function to set a value on a node (pre-existing)
//
// 	client		http.Client that can control functionality, like Timeouts (nil is ok)
// 	tr		http.Transport that can set TLS client attributes (nil is ok)
// 	proto		"http" or "https"
// 	host		host to connect with
// 	port		port to connect to
// 	key		etcd node key/directory
// 	Value		string value to set ... yes, string, that's all etcd works with.
//	ttl		optional integer TTL for this key/value (expires after TTL)
//
func SetValue(client *http.Client, tr *http.Transport,
	proto string, host string, port int,
	key string, value string, ttl ...int) (EtcdResponse, error) {

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

	urlStr := fmt.Sprintf("%s://%s:%d/v2/keys/%s", proto, host, port, key)
	data := url.Values{}
	data.Set("value", value)
	if len(ttl) > 0 {
		data.Set("ttl", strconv.Itoa(ttl[0]))
	}
	encoded := data.Encode()

	var request *http.Request
	request, err = http.NewRequest("PUT", urlStr, bytes.NewBufferString(encoded))
	if err != nil {
		return EtcdResponse{}, errors.New("http.NewRequest: " + err.Error())
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(encoded)))

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
	return r, nil // All good
}
