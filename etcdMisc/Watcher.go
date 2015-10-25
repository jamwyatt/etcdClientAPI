package etcdMisc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Watcher(ctrl chan bool, key string, recursive bool, waitIndex ...int) (WatchResponse, error) {
	var err error
	tr := &http.Transport{
		DisableKeepAlives:     true, // No persistent connections
		ResponseHeaderTimeout: 0,    // No timeouts
	}
	client := &http.Client{
		Timeout:   0, // No timeout for watch
		Transport: tr,
	}

	url := fmt.Sprintf("http://localhost:4001/v2/keys/%s?wait=true&recursive=%t", key, recursive)
	if len(waitIndex) > 0 {
		url += fmt.Sprintf("&waitIndex=%d", waitIndex[0])
	}
	fmt.Println("URL: ", url)
	var request *http.Request
	request, err = http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Request Error:", err)
		return WatchResponse{}, err
	}

	syncChannel := make(chan WatchResponse)
	go func(s chan WatchResponse) {
		var resp *http.Response
		resp, err = client.Do(request)
		if err != nil {
			fmt.Println("GET Error:", err)
			s <- WatchResponse{}
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Body Error:", err)
			s <- WatchResponse{}
		}
		r := WatchResponse{}
		err = json.Unmarshal(body, &r)
		if err != nil {
			fmt.Println("Failed to decode response")
			s <- WatchResponse{}
		}
		s <- r
	}(syncChannel)

	select {
	case msg := <-syncChannel:
		return msg, err
	case <-ctrl:
		tr.CancelRequest(request)
		return WatchResponse{}, errors.New("Canceled")
	}

}
