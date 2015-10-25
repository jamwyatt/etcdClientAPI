package main

import (
	"fmt"
	"github.com/jamwyatt/etcdClientAPI/etcdMisc"
	"os"
	"time"
)

// Test will watch the target twice and then timeout the last
// Drive the test by updating the target in etcd twice and then wait for the timeout.

func main() {
	if len(os.Args) < 1 {
		fmt.Println("usage: Missing watch key")
		os.Exit(-1)
	}

	ctrl := make(chan bool)

	// Watch for the next event
	r, err := etcdMisc.Watcher(ctrl, os.Args[1], true)
	if err != nil {
		fmt.Println("Failed")
	} else {
		fmt.Printf("WatchResponse: %s\n", r)
	}

	// Watch for the 'proper' next event, by number
	r, err = etcdMisc.Watcher(ctrl, os.Args[1], true, r.Node.ModifiedIndex+1)
	if err != nil {
		fmt.Println("Failed")
	} else {
		fmt.Printf("WatchResponse: %s\n", r)
	}

	// Watch for an event ad then cancel it.
	go func() {
		r, err := etcdMisc.Watcher(ctrl, os.Args[1], true)
		if err != nil {
			fmt.Println("Passed with expected fail", err)
		} else {
			fmt.Printf("Failed with WatchResponse: %s\n", r)
		}
	}()
	time.Sleep(time.Second * 2)
	ctrl <- true
	time.Sleep(time.Second * 1)

}
