package main

import (
	"etcdClientAPI/etcdMisc"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 1 {
		fmt.Println("usage: Missing watch key")
		os.Exit(-1)
	}

	// Watch for the next event
	r, err := etcdMisc.WatchEvent(os.Args[1], true)
	if err != nil {
		fmt.Println("Failed")
	} else {
		fmt.Printf("WatchResponse: %s\n", r)
	}

	// Watch for the 'proper' next event, by number
	r, err = etcdMisc.WatchEvent(os.Args[1], true, r.Node.ModifiedIndex+1)
	if err != nil {
		fmt.Println("Failed")
	} else {
		fmt.Printf("WatchResponse: %s\n", r)
	}

}
