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
//  Function to watch/report change on node/tree from etcd.
//
//	ctrl		channel that can be used to abort a long timeout (single write aborts)
// 	key		etcd node key/directory
// 	recursive	true = watch recursively
// 	waitIndex	index of the node to watch for (useful to avoid missing an event)
//
//
func (conn EtcdConnection) Watcher(ctrl chan bool, key string, recursive bool, waitIndex ...int) (EtcdResponse, error) {

	var err error
	url := fmt.Sprintf("%s://%s:%d/v2/keys%s?wait=true&recursive=%t", conn.Proto, conn.Host, conn.Port, key, recursive)
	if len(waitIndex) > 0 {
		url += fmt.Sprintf("&waitIndex=%d", waitIndex[0])
	}
	var request *http.Request
	request, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return EtcdResponse{err: errors.New("http.NewRequest: " + err.Error())}, err
	}

	syncChannel := make(chan EtcdResponse)
	go func(s chan EtcdResponse) {
		var resp *http.Response
		resp, err = conn.Client.Do(request)
		if err != nil {
			s <- EtcdResponse{err: errors.New("http.client.Do: " + err.Error())}
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			s <- EtcdResponse{err: errors.New("ioutil.ReadAll: " + err.Error())}
			return
		}
		r := EtcdResponse{}
		err = json.Unmarshal(body, &r)
		if err != nil {
			s <- EtcdResponse{err: errors.New("json.Unmarshal: " + err.Error())}
			return
		}
		s <- r
	}(syncChannel)

	select {
	case msg := <-syncChannel:
		return msg, err
	case <-ctrl:
		conn.Transport.CancelRequest(request)
		<-syncChannel
		close(syncChannel)
		return EtcdResponse{}, errors.New("Controller Directed Cancel")
	}

}

//
//  A routine that returns an event channel to receive continuous stream of watch events
//  Unlike Watcher, this starts with the first event to be received and then watches the
//  next event in sequence.
//
//	ctrl		channel that can be used to abort a long timeout (single write aborts)
// 	key		etcd node key/directory
// 	recursive	true = watch recursively
//
//
func (conn EtcdConnection) EventStream(ctrl chan bool, key string, recursive bool) chan EtcdResponse {

	index := -1
	response := make(chan EtcdResponse) // returned to caller
	go func() {
		myCtrl := make(chan bool)
		insideSync := make(chan EtcdResponse)
		for {
			go func() {
				var resp EtcdResponse
				var err error
				if index > 0 {
					// Index matching to avoid loss
					resp, err = conn.Watcher(myCtrl, key, recursive, index)
				} else {
					// First one takes the first response
					resp, err = conn.Watcher(myCtrl, key, recursive)
				}
				if err != nil {
					insideSync <- EtcdResponse{err: err}
				} else {
					index = resp.Node.ModifiedIndex + 1
					insideSync <- resp
				}
			}()
			var msg EtcdResponse
			select {
			case msg = <-insideSync:
				// Pass the message to original caller
				response <- msg
			case <-ctrl:
				// Outside shutdown ... pass it along
				myCtrl <- true
				// ending this thread, but wait for final message (and send it)
				response <- <-insideSync
				close(insideSync)
				close(myCtrl)
				close(response)
				return
			}
		}
	}()
	// Return channel for event stream back to caller
	return response
}
