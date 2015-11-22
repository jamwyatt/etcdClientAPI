# etcd client API example

A simple application that can interface with most features of ectd. This will operate based off of
https://github.com/coreos/etcd/blob/master/Documentation/api.md

This is not intended to be earth shattering of any sort, but a place to experiment with golang and etcd.

So far, there is an event watching system that allows for single blocking event watch and an Event stream watcher. The stream
requeries the same key using the 'next' index to ensure that the stream never drops an event. Both cases, the API allows
for recursive watching of a tree or a single node.

Aside from event 'watching', the general functions supported:

Set/Get for a value
Mkdir/DeletDir for a directory (recursive delete supported)
Recursive Get

Get supports recursive 'get' and the response datastructure supports a recursive structure and related print.

Also, there are error values in the main response structure to report Etcd errors (Cause/ErrorCode/Message).

Note that a leading '/' is required in any key or path. The responses from etcd will contain mostly the complete set of etcd supported elements. Use the 'Cause' field to check of errors.

Get responses will contain the 'dir' boolean as defined by etcd. This will indicate a directory or a key.

