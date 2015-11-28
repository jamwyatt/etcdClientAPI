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
//  Mkdir - Function to make an etcd directory
//
// 	conn		ectdConnection, made with etcdMisc.MakeEtcdConnection()
// 	port		port to connect to
// 	key		etcd node directory
//
func Mkdir(conn etcdConnection, key string) (EtcdResponse, error) {

	var err error
	urlStr := fmt.Sprintf("%s://%s:%d/v2/keys%s", conn.Proto, conn.Host, conn.Port, key)
	data := url.Values{}
	data.Set("dir", "true")
	encoded := data.Encode()

	var request *http.Request
	request, err = http.NewRequest("PUT", urlStr, bytes.NewBufferString(encoded))
	if err != nil {
		return EtcdResponse{}, errors.New("http.NewRequest: " + err.Error())
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(encoded)))

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
		return r, errors.New("Mkdir etcd error: " + r.Message)
	}
	return r, nil // All good
}
