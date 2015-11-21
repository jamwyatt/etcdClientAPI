package main

/*

Simple test tool for the Watch tools in etcdMisc

*/

import (
	"fmt"
	"github.com/jamwyatt/etcdClientAPI/etcdMisc"
	"net/http"
	"os"
	"time"
)

// Test will watch the target a total of 4 times.
// NOTE: Outside action required to trigger event on etcd target (4X)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: Missing watch key")
		os.Exit(-1)
	}

	client := &http.Client{
		Timeout: 0, // No time out on a Watch
	}

	// Set followed by a 'Get' to verify
	r, err := etcdMisc.SetValue(client, nil, "http", "localhost", 4001, "chickens/blob1", "Hello")
	if err != nil {
		fmt.Println("Failed to set etcd value:", err)
		os.Exit(-1)
	}
	fmt.Printf("SetValue with string(\"Hello\") Response: %s\n", r)

	r, err = etcdMisc.GetValue(client, nil, "http", "localhost", 4001, "chickens/blob1")
	if err != nil {
		fmt.Println("Failed to get etcd value:", err)
		os.Exit(-1)
	}
	fmt.Printf("GetValue Response: %s\n", r)
	if r.Node.Value != "Hello" {
		fmt.Printf("FAILED: get should return 'Hello'\n")
		os.Exit(-1)
	}

	// Delete Key
	r, err = etcdMisc.DeleteKey(client, nil, "http", "localhost", 4001, "chickens/blob1")
	if err != nil {
		fmt.Println("Failed to delete etcd key:", err)
		os.Exit(-1)
	}
	fmt.Printf("DeleteKey Response: %s\n", r)

	// Delete non-existant Key
	r, err = etcdMisc.DeleteKey(client, nil, "http", "localhost", 4001, "chickens/doesNotExist")
	if err == nil {
		fmt.Printf("Found no error deleting a missing etcd key: %s/%s\n", err, r)
		os.Exit(-1)
	}
	fmt.Printf("DeleteKey errored as expected deleting a missing key: %s\n", err)

	// Set with a TTL
	r, err = etcdMisc.SetValue(client, nil, "http", "localhost", 4001, "chickens/blob3", "Hello", 2)
	if err != nil {
		fmt.Println("Failed to set etcd value:", err)
		os.Exit(-1)
	}
	fmt.Printf("SetValue with string(\"Hello\", TTL=2) Response: %s\n", r)
	r, err = etcdMisc.GetValue(client, nil, "http", "localhost", 4001, "chickens/blob3")
	if err != nil {
		fmt.Println("Failed to get etcd value:", err)
		os.Exit(-1)
	}
	fmt.Printf("GetValue Response: %s\n", r)
	if r.Node.Value != "Hello" {
		fmt.Printf("FAILED: get should return 'Hello'\n")
		os.Exit(-1)
	}

	fmt.Printf("Sleep 3 seconds ..... should have expired\n")
	time.Sleep(3 * time.Second)
	r, err = etcdMisc.GetValue(client, nil, "http", "localhost", 4001, "chickens/blob3")
	if err == nil {
		fmt.Println("Failed to NOT get expired etcd value:", err)
		os.Exit(-1)
	}
	fmt.Println("Received expected error GETting expired key (missing key)")

	// err = etcdMisc.SetValue(client, nil, "http", "localhost", 4001, "/chickens/blob", true)
	// err = etcdMisc.SetValue(client, nil, "http", "localhost", 4001, "/chickens/blob", "hello")
	// Get a single event for the watched key (blocking call)
	// Use the bool channel to abort this query if needed.
	r, err = etcdMisc.Watcher(client, nil, make(chan bool), "http", "localhost", 4001, os.Args[1], true)
	if err != nil {
		fmt.Println("Failed: ", err)
	} else {
		fmt.Printf("EtcdResponse: %s\n", r)
	}

	tr := &http.Transport{
		DisableKeepAlives: false, // Allow connection reuse
	}
	// Non-blocking event stream. Shutdown via bool channel.
	// Create event stream from key on cmd line, recursive == true
	ctrl := make(chan bool)
	events := etcdMisc.EventStream(client, tr, ctrl, "http", "localhost", 4001, os.Args[1], true)

	// Receive exactly three events
	for i := 0; i < 3; i++ {
		msg := <-events
		if msg.GetError() != nil {
			fmt.Println("Event ERR:", msg.GetError())
			// Shutdown the event stream and restart ... most likely out of sync or etcd is dead
			ctrl <- true
			events = etcdMisc.EventStream(client, tr, ctrl, "http", "localhost", 4001, os.Args[1], true)
		} else {
			fmt.Println("Event: ", msg)
		}
	}
	// Shutdown event stream
	ctrl <- true
	close(ctrl)

}
