// Simple etcd client library to interface with etcd's HTTP API
package etcdMisc

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

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type EtcdConnection struct {
	Client    *http.Client
	Transport *http.Transport
	Proto     string
	Host      string
	Port      int
}

func (c EtcdConnection) String() string {
	return c.Proto + "://" + c.Host + ":" + strconv.Itoa(c.Port)
}

//
//  makeEtcdConnection - Build a connection for use with the other etcdMisc functions
//
//      client          http.Client that can control functionality, like Timeouts (nil is ok)
//			Note: timeouts should be 0 when using 'EventStream' and 'Watcher'
//      tr              http.Transport that can set TLS client attributes (nil is ok)
//      proto           "http" or "https"
//      host            host to connect with
//      port            Optional, defaults to 80/443 depending on proto
//
func MakeEtcdConnection(client *http.Client, trans *http.Transport, proto string, host string, port ...int) (EtcdConnection, error) {

	if client == nil {
		return EtcdConnection{}, errors.New("Missing http.client")
	}
	if host == "" {
		return EtcdConnection{}, errors.New("Missing host")
	}
	if proto == "" {
		return EtcdConnection{}, errors.New("Missing proto")
	}
	lowerProto := strings.ToLower(proto)
	if lowerProto != "http" && lowerProto != "https" {
		return EtcdConnection{}, errors.New("Unsupported proto: " + proto)
	}

	connection := EtcdConnection{Client: client}
	if trans == nil {
		connection.Transport = &http.Transport{}
	}
	connection.Proto = lowerProto
	connection.Host = host
	if len(port) > 0 {
		connection.Port = port[0]
	} else if lowerProto == "http" {
		connection.Port = 80
	} else {
		connection.Port = 443
	}
	return connection, nil
}
