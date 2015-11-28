package etcdMisc

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type etcdConnection struct {
	Client    *http.Client
	Transport *http.Transport
	Proto     string
	Host      string
	Port      int
}

func (c etcdConnection) String() string {
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
func MakeEtcdConnection(client *http.Client, trans *http.Transport, proto string, host string, port ...int) (etcdConnection, error) {

	if client == nil {
		return etcdConnection{}, errors.New("Missing http.client")
	}
	if host == "" {
		return etcdConnection{}, errors.New("Missing host")
	}
	if proto == "" {
		return etcdConnection{}, errors.New("Missing proto")
	}
	lowerProto := strings.ToLower(proto)
	if lowerProto != "http" && lowerProto != "https" {
		return etcdConnection{}, errors.New("Unsupported proto: " + proto)
	}

	connection := etcdConnection{Client: client}
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
