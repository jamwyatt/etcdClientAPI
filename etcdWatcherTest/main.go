package main

/*

Simple test tool for the Watch tools in etcdMisc

*/

import (
	"fmt"
	"github.com/jamwyatt/etcdClientAPI/etcdMisc"
	"os"
)

// Test will watch the target a total of 4 times.
// NOTE: Outside action required to trigger event on etcd target (4X)

func main() {
	if len(os.Args) < 1 {
		fmt.Println("usage: Missing watch key")
		os.Exit(-1)
	}

	// Get a single event for the watched key (blocking call)
	r, err := etcdMisc.Watcher(make(chan bool), os.Args[1], true)
	if err != nil {
		fmt.Println("Failed")
	} else {
		fmt.Printf("WatchResponse: %s\n", r)
	}

	// Non-blocking event stream. Shutdown via bool channel.
	// Create event stream from key on cmd line, recursive == true
	ctrl := make(chan bool)
	events := etcdMisc.EventStream(ctrl, os.Args[1], true)

	// Receive exactly three events
	for i := 0; i < 3; i++ {
		msg := <-events
		fmt.Println("Event: ", msg)
	}
	// Shutdown event stream
	ctrl <- true
	close(ctrl)

}
