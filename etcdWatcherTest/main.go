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

	fmt.Printf("Time.Second: %d\n", time.Second)

	client := &http.Client{
		Timeout: 0, // No time out on a Watch
	}

	tr := &http.Transport{
		DisableKeepAlives: false, // Allow connection reuse
	}

	// Get a single event for the watched key (blocking call)
	// Use the bool channel to abort this query if needed.
	r, err := etcdMisc.Watcher(client, tr, make(chan bool), "localhost", 4001, os.Args[1], true)
	if err != nil {
		fmt.Println("Failed: ", err)
	} else {
		fmt.Printf("WatchResponse: %s\n", r)
	}

	// Non-blocking event stream. Shutdown via bool channel.
	// Create event stream from key on cmd line, recursive == true
	ctrl := make(chan bool)
	events := etcdMisc.EventStream(client, tr, ctrl, "localhost", 4001, os.Args[1], true)

	// Receive exactly three events
	for i := 0; i < 3; i++ {
		msg := <-events
		if msg.GetError() != nil {
			fmt.Println("Event ERR:", msg.GetError())
			// Shutdown the event stream and restart ... most likely out of sync or etcd is dead
			ctrl <- true
			events = etcdMisc.EventStream(client, tr, ctrl, "localhost", 4001, os.Args[1], true)
		} else {
			fmt.Println("Event: ", msg)
		}
	}
	// Shutdown event stream
	ctrl <- true
	close(ctrl)

}
