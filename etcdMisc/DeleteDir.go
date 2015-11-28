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
// 	conn		ectdConnection, made with etcdMisc.MakeEtcdConnection()
// 	key		etcd node key/directory
//	recursive	true for recursive delete of everything
//
func DeleteDir(conn etcdConnection, key string, recursive ...bool) (EtcdResponse, error) {

	var err error
	urlStr := fmt.Sprintf("%s://%s:%d/v2/keys%s?dir=true", conn.Proto, conn.Host, conn.Port, key)
	if len(recursive) > 0 && recursive[0] {
		urlStr += "&recursive=true"
	}

	var request *http.Request
	request, err = http.NewRequest("DELETE", urlStr, nil)
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
		return r, errors.New("DeleteDir etcd error: " + r.Message)
	}
	return r, nil // All good
}
