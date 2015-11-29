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
// Function to get the value of a pre-existing key
//
// 	port		port to connect to
// 	key		etcd node key/directory
// 	recurse		Recursive get, useful on directories, but the response is recursive too
// 	sort		apply etcd sorting?
//
//
func (conn EtcdConnection) GetValue(key string, recurse bool, sort bool) (EtcdResponse, error) {

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
