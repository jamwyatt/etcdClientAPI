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
// 	conn		ectdConnection, made with etcdMisc.MakeEtcdConnection()
// 	key		etcd node key/directory
// 	Value		string value to set ... yes, string, that's all etcd works with.
//	ttl		optional integer TTL for this key/value (expires after TTL)
//
func SetValue(conn etcdConnection, key string, value string, ttl ...int) (EtcdResponse, error) {

	var err error
	urlStr := fmt.Sprintf("%s://%s:%d/v2/keys%s", conn.Proto, conn.Host, conn.Port, key)
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
		return r, errors.New("Set etcd error: " + r.Message)
	}
	return r, nil // All good
}
