package etcdMisc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Transport as a parameter, to allow for TLS support
func Watcher(client *http.Client, tr *http.Transport, ctrl chan bool,
	host string, port int, key string, recursive bool, waitIndex ...int) (WatchResponse, error) {

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

	url := fmt.Sprintf("http://%s:%d/v2/keys/%s?wait=true&recursive=%t", host, port, key, recursive)
	if len(waitIndex) > 0 {
		url += fmt.Sprintf("&waitIndex=%d", waitIndex[0])
	}
	// fmt.Println("URL: ", url)
	var request *http.Request
	request, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return WatchResponse{err: errors.New("http.NewRequest: " + err.Error())}, err
	}

	syncChannel := make(chan WatchResponse)
	go func(s chan WatchResponse) {
		var resp *http.Response
		resp, err = client.Do(request)
		if err != nil {
			s <- WatchResponse{err: errors.New("http.client.Do: " + err.Error())}
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			s <- WatchResponse{err: errors.New("ioutil.ReadAll: " + err.Error())}
			return
		}
		r := WatchResponse{}
		err = json.Unmarshal(body, &r)
		if err != nil {
			s <- WatchResponse{err: errors.New("json.Unmarshal: " + err.Error())}
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
		return WatchResponse{}, errors.New("Canceled")
	}

}

// Transport as a parameter. Allows for TLS support
func EventStream(client *http.Client, tr *http.Transport, ctrl chan bool,
	host string, port int, key string, recursive bool) chan WatchResponse {

	index := -1
	response := make(chan WatchResponse) // returned to caller
	go func() {
		myCtrl := make(chan bool)
		insideSync := make(chan WatchResponse)
		for {
			go func() {
				var resp WatchResponse
				var err error
				if index > 0 {
					// Index matching to avoid loss
					resp, err = Watcher(client, tr, myCtrl, host, port, key, true, index)
				} else {
					// First one takes the first response
					resp, err = Watcher(client, tr, myCtrl, host, port, key, true)
				}
				if err != nil {
					insideSync <- WatchResponse{err: err}
				} else {
					index = resp.Node.ModifiedIndex + 1
					insideSync <- resp
				}
			}()
			var msg WatchResponse
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
