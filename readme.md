Example from an interview question to build a LRU in memory cache.

This implementation is in Go and uses a map of pointers to nodes, where the nodes maintain a double linked list with ordering constraint most recently used to least recently used. The cache has a fixed size that must be set when it is created, but the cache could easily be made resizable.

``` 
  newest -> [Key 1 ]      [Key 2 ]      [Key 3 ]      [Key 4 ]      [Key 5 ] <- oldest
            [Node *] <--> [Node *] <--> [Node *] <--> [Node *] <--> [Node *] 
            [Blob  ]      [Blob  ]      [Blob  ]      [Blob  ]      [Blob  ]
```

Methods
* New (size)
* Read (value)
* Write (address, value)
* Print

Complexity
* Read is O(1) - hash (map key) + unlink node + link at newest (head)
* Write is O(1) - hash (map key) + update & unlink node || create node  + link at newest (head)
* New is O(1) - allocate empty map, set maximum size 
* Print is O(n) - obviously to traverse the structure

Constraints:
* Does not currently support multi-threaded access, will probably update that.
* Key must be of a comparable type, blobs can be any type (uses Go generics)
* No disk or other backstore for the writes, this is just a toy example

Algorithm:

Essentially there are three non-trivial cases
1. Read hit
2. Write hit
3. Write miss

When there is a read or write hit, then node is unlinked and moved to most recent (unless it is already the most recent). Nodes are quickly found by key using a map (hash).

When there is a write miss then we need to either:
1. Add to the cache if there is space
2. Remove the least used node to create space and then add the new as the most recent

The common function is to take an unlinked node and update the most recent or head of the list.  Thus the write hit just adds the unlink node, then uses update most recent.
