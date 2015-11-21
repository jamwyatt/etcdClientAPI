# etcd client API example

A simple application that can interface with most features of ectd. This will operate based off of
https://github.com/coreos/etcd/blob/master/Documentation/api.md

This is not intended to be earth shattering of any sort, but a place to experiment with golang.

So far, there is an event watching system that allows for single blocking event watch and an Event stream watcher. The stream
requeries the same key using the 'next' index to ensure that the stream never drops an event. Both cases, the APi allows
for recursive watching of a tree or a single node.

Now supports basic Set/Get/Delete of a value in etcd. The node must exist. Next is node management and listing.

Get supports recursive 'get' and the response datastructure supports a recursive structure and related print.

Also, there are error values in the main response structure to report Etcd errors (Cause/ErrorCode/Message).

