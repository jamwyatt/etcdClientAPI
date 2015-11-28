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
// 	conn		ectdConnection, made with etcdMisc.MakeEtcdConnection()
// 	key		etcd node key/directory
//
//
func DeleteKey(conn etcdConnection, key string) (EtcdResponse, error) {

	var err error
	url := fmt.Sprintf("%s://%s:%d/v2/keys%s", conn.Proto, conn.Host, conn.Port, key)
	var request *http.Request
	request, err = http.NewRequest("DELETE", url, nil)
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
		return r, errors.New("etcd error: " + r.Message)
	}
	return r, nil // All good
}
