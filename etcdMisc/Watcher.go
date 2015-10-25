package etcdMisc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func Watcher(key string, recursive bool, waitIndex ...int) (WatchResponse, error) {
	client := &http.Client{
		Timeout: 0, // No timeout for watch
	}

	url := fmt.Sprintf("http://localhost:4001/v2/keys/%s?wait=true&recursive=%t", key, recursive)
	if len(waitIndex) > 0 {
		url += fmt.Sprintf("&waitIndex=%d", waitIndex[0])
	}
	fmt.Println("URL: ", url)

	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("GET Error:", err)
		os.Exit(-1)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Body Error:", err)
		os.Exit(-1)
	}
	// fmt.Printf("body: %v\n", string(body))

	r := WatchResponse{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		fmt.Println("Failed to decode response")
		return WatchResponse{}, err
	}
	return r, nil
}
