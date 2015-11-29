// Simple test tool for the Watch tools in etcdMisc
package main

/*
etcdClientAPI is a simple golang library to interface with etcd's API
Copyright (C) 2015 J. Robert Wyatt

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
*/

/*

Simple test tool for the Watch tools in etcdMisc

- Note: This test expects http://loadhost:4001 to be an etcd server that it
        can do anything it wants with.

*/

import (
	"fmt"
	"github.com/jamwyatt/etcdClientAPI/etcdMisc"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {

	connection, err := etcdMisc.MakeEtcdConnection(&http.Client{Timeout: 0}, nil, "http", "localhost", 4001)
	fmt.Printf("Connection: %s\n", connection)

	// ---------------------------------------------------------------------------------------------------
	// Protective Delete of a directory, might be non-existant
	r, err := connection.DeleteDir("/junk", true)
	if err != nil {
		fmt.Println("Failed to delete directory:", err)
	} else {
		fmt.Printf("Deleted directory: %s\n", r)
	}

	// Start making a directory
	r, err = connection.Mkdir("/junk")
	if err != nil {
		fmt.Println("Failed to create directory:", err)
		os.Exit(-1)
	}
	fmt.Printf("Created directory successfully: %s\n", r)

	// Add to the previous directory
	r, err = connection.Mkdir("/junk/one")
	if err != nil {
		fmt.Println("Failed to create directory:", err)
		os.Exit(-1)
	}
	fmt.Printf("Created directory successfully: %s\n", r)

	// ---------------------------------------------------------------------------------------------------

	// Set followed by a 'Get' to verify
	r, err = connection.SetValue("/junk/blob1", "Hello")
	if err != nil {
		fmt.Println("Failed to set etcd value:", err)
		os.Exit(-1)
	}
	fmt.Printf("SetValue with string(\"Hello\") Response: %s\n", r)

	r, err = connection.GetValue("/junk/blob1", false, false)
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
	r, err = connection.DeleteKey("/junk/blob1")
	if err != nil {
		fmt.Println("Failed to delete etcd key:", err)
		os.Exit(-1)
	}
	fmt.Printf("DeleteKey Response: %s\n", r)

	// ---------------------------------------------------------------------------------------------------

	// Delete non-existant Key
	r, err = connection.DeleteKey("/junk/doesNotExist")
	if err == nil {
		fmt.Printf("Found no error deleting a missing etcd key: %s/%s\n", err, r)
		os.Exit(-1)
	}
	fmt.Printf("DeleteKey errored as expected deleting a missing key: %s\n", err)

	// ---------------------------------------------------------------------------------------------------

	// Set with a TTL
	r, err = connection.SetValue("/junk/blob3", "Hello", 2)
	if err != nil {
		fmt.Println("Failed to set etcd value:", err)
		os.Exit(-1)
	}
	fmt.Printf("SetValue with string(\"Hello\", TTL=1) Response: %s\n", r)

	// Verify it was set
	r, err = connection.GetValue("/junk/blob3", false, false)
	if err != nil {
		fmt.Println("Failed to get etcd value:", err)
		os.Exit(-1)
	}
	fmt.Printf("GetValue Response: %s\n", r)
	if r.Node.Value != "Hello" {
		fmt.Printf("FAILED: get should return 'Hello'\n")
		os.Exit(-1)
	}

	// wait for the TTL above to expire and verify it is gone
	fmt.Printf("Sleep 3 seconds ..... to wait for key to expire\n")
	time.Sleep(3 * time.Second)
	r, err = connection.GetValue("/junk/blob3", false, false)
	if err == nil {
		fmt.Println("Failed to NOT get expired etcd value:", err)
		os.Exit(-1)
	}
	fmt.Println("Received expected error GETting expired key (missing key)")

	// ---------------------------------------------------------------------------------------------------

	// Background function to make 4 changes that will be picked up by the
	// coming watchers that require at least 4 changes to avoid timeouts.
	go func(keyRoot string) {
		time.Sleep(1 * time.Second)
		for t := 0; t < 4; t++ {
			newVal := keyRoot + strconv.Itoa(t)
			fmt.Printf("SetValue background Value for watcher test: %s\n", newVal)
			// Set followed by a 'Get' to verify
			_, err := connection.SetValue(newVal, "miscSillyValue")
			if err != nil {
				fmt.Println("Failed to set etcd value:", err)
				os.Exit(-1)
			}
			time.Sleep(time.Second / 2)
		}
	}("/junk/foobarOne")

	// ---------------------------------------------------------------------------------------------------

	// Get a single event for the watched key (blocking call) (consumes one of the background events)
	// Use the bool channel to abort this query if needed.
	ctrl := make(chan bool)
	r, err = connection.Watcher(ctrl, "/", true)
	if err != nil {
		fmt.Println("Failed: ", err)
	} else {
		fmt.Printf("EtcdResponse: %s\n", r)
	}

	// ---------------------------------------------------------------------------------------------------

	// Non-blocking event stream. Shutdown via bool channel. Consumes 3 of the background events.
	// Create event stream from key on cmd line, recursive == true
	events := connection.EventStream(ctrl, "/", true)

	// Receive exactly three events
	for i := 0; i < 3; i++ {
		timer := time.NewTimer(time.Second * 3)
		select {
		case msg := <-events:
			if msg.GetError() != nil {
				fmt.Println("Event ERR:", msg.GetError())
				// Shutdown the event stream and restart ... most likely out of sync or etcd is dead
				ctrl <- true
				events = connection.EventStream(ctrl, "/junk", true)
			} else {
				fmt.Println("Event: ", msg)
			}
		case <-timer.C:
			fmt.Println("TIMEOUT: FAILED TO RECEIVE WATCHER EVENT from etcd")
			os.Exit(-1)
		}
	}
	// Shutdown event stream
	ctrl <- true
	close(ctrl)

	// ---------------------------------------------------------------------------------------------------

	// Recursive Get
	r, err = connection.GetValue("/", true, false)
	if err != nil {
		fmt.Println("Failed recursive GET:", err)
		os.Exit(-1)
	}
	fmt.Printf("Recursive GET results:\n%s", r)

	// -------------------------------------------------------------------------------------------------------

	// Final cleanup ... must exist
	r, err = connection.DeleteDir("/junk", true)
	if err != nil {
		fmt.Println("Failed to delete directory:", err)
		os.Exit(-1)
	}
	fmt.Printf("Successfully cleaned up by Deleting directory: %s\n", r)

}
