# go etcd client API example (version 0.5)

A simple golang library that can interface with key features of ectd. This operates based off of
https://github.com/coreos/etcd/blob/master/Documentation/api.md

This is not intended to be earth shattering of any sort, but a place to experiment with golang and etcd. I will attempt to make it reasonably stable and useful. There is an event watching system that allows for single blocking event watch and an Event stream watcher. The stream requeries the same key using the 'next' index to ensure that the stream never drops an event. Both cases, the API allows
for recursive watching of a tree or a single node.

Everything in terms of connecting to an instance of etcd is handled within the EtcdConnection structure. Do not make one yourself, instead call the factory function **EtcdConnection.makeEtcdConnection()**. When making this instance, you must provide an http.Client and  and optional http.Transport. Using these objects, you can control details about https connections as well as setting performance related attributes, like timeout and keepalives. The most interesting of those is the use of timeout. When you are using **EtcdConnection.Watcher()** or **EtcdConnection.EvenetStream()**, timeout is interesting.

Aside from event 'watching/streaming', other general functions are:

* Set/Get for a value
* Mkdir/DeletDir for a directory (recursive delete supported)
* Recursive Get

**Note** that a leading '/' is required in any etcd key or path. 

The response datastructure (**EtcdResponse**) is recursive and supports the **String()** function. Error fields are located in the main  structure (Cause/ErrorCode/Message). These error fields and non-error fields within the response structure, line up with etcd responses (they start with a capital letter in go, but etcd starts with lowercase).

```
type EtcdResponse struct {
        Action   string
        Node     Node
        PrevNode Node
        Cause     string
        ErrorCode int
        Message   string
        err error
}
```

When using '**Watcher()/EventStream()**', you should consider using a separate connection that has a timeout of 0 for the http.Client. Each instance of the **EtcdConnection** is a connection to etcd.

Get responses will contain the 'dir' boolean as defined by etcd. This will indicate a directory or a key.

