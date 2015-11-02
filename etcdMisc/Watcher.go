package etcdMisc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

//
//  Watcher - Function to watch/report change on node/tree from etcd.
//
// 	client		http.Client that can control functionality, like Timeouts (nil is ok)
// 	tr		http.Transport that can set TLS client attributes (nil is ok)
// 	ctrl		channel that can be used to abort a long timeout (single write aborts)
// 	proto		"http" or "https"
// 	host		host to connect with
// 	port		port to connect to
// 	key		etcd node key/directory
// 	recursive	true = watch recursively
// 	waitIndex	index of the node to watch for (useful to avoid missing an event)
//
//
func Watcher(client *http.Client, tr *http.Transport, ctrl chan bool,
	proto string, host string, port int, key string, recursive bool, waitIndex ...int) (EtcdResponse, error) {

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

	url := fmt.Sprintf("%s://%s:%d/v2/keys/%s?wait=true&recursive=%t", proto, host, port, key, recursive)
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
		resp, err = client.Do(request)
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
		tr.CancelRequest(request)
		<-syncChannel
		close(syncChannel)
		return EtcdResponse{}, errors.New("Controller Directed Cancel")
	}

}

//
//  EventStream - a go routine that returns an event channel to receive continuous stream of watch events
//                Unlike Watcher, this starts with the first event to be received and then watches the
//		  next event in sequence.
//
// 	client		http.Client that can control functionality, like Timeouts (nil is ok)
// 	tr		http.Transport that can set TLS client attributes (nil is ok)
// 	ctrl		channel that can be used to abort a long timeout (single write aborts)
// 	proto		"http" or "https"
// 	host		host to connect with
// 	port		port to connect to
// 	key		etcd node key/directory
// 	recursive	true = watch recursively
//
//
func EventStream(client *http.Client, tr *http.Transport, ctrl chan bool,
	proto string, host string, port int, key string, recursive bool) chan EtcdResponse {

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
					resp, err = Watcher(client, tr, myCtrl, proto, host, port, key, recursive, index)
				} else {
					// First one takes the first response
					resp, err = Watcher(client, tr, myCtrl, proto, host, port, key, recursive)
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
