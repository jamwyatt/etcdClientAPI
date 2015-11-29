// Simple etcd client library to interface with etcd's HTTP API
package etcdMisc

/*
etcdClientAPI is a simple golang library to interface with etcd's API
Copyright (C) 2015 J. Robert Wyatt

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
*/

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

//
//  Function to delete a directory (optional recursion)
//
// 	key		etcd node key/directory
//	recursive	true for recursive delete of everything
//
func (conn EtcdConnection) DeleteDir(key string, recursive ...bool) (EtcdResponse, error) {

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
