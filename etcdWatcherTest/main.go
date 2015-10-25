package main

import (
	"fmt"
	"github.com/jamwyatt/etcdClientAPI/etcdMisc"
	"os"
)

func main() {
	if len(os.Args) < 1 {
		fmt.Println("usage: Missing watch key")
		os.Exit(-1)
	}

	// Watch for the next event
	r, err := etcdMisc.Watcher(os.Args[1], true)
	if err != nil {
		fmt.Println("Failed")
	} else {
		fmt.Printf("WatchResponse: %s\n", r)
	}

	// Watch for the 'proper' next event, by number
	r, err = etcdMisc.Watcher(os.Args[1], true, r.Node.ModifiedIndex+1)
	if err != nil {
		fmt.Println("Failed")
	} else {
		fmt.Printf("WatchResponse: %s\n", r)
	}

}
