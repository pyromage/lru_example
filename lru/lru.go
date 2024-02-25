package lru

import (
	"errors"
	"fmt"
)

type node[K comparable, V any] struct {
	newer K
	older K
	blob V
}

type Cache[K comparable,V any] struct {
	nodes map[K]*node[K, V] 
	oldest K
	newest K
	maxSize int
}

func NewCache[K comparable, V any](size int) (*Cache[K,V], error) {
	// too small, must be larger than 2
	if size <= 1 {
		return nil, errors.New("maximum size must be larger than 1")
	}
	
	var empty K

	return &Cache[K,V]{
		nodes: 	make(map[K]*node[K,V]),
		oldest: empty,
		newest: empty,
		maxSize: size,
	}, nil
}

func (c *Cache[K,V])Read(addr K) (V, bool) {
	var nilKey K
	var nilValue V

	if c.newest == nilKey || c.oldest == nilKey || c.maxSize <= 1 || addr == nilKey {
		return nilValue, false	
	}
	
	val, ok := c.nodes[addr]
	if !ok {
		return nilValue, false
	} 

	c.moveToHead(addr, false)

	return val.blob,ok
}

// true if inserted, false if failed to insert
func (c *Cache[K,V])Write(addr K, value V) error {
	var nilKey K

	if (c.maxSize <= 1) {
		return errors.New("cache size must be larger than 1")
	}

	if (addr == nilKey){
		return errors.New("cannot write to the empty key")

	}

	// is this a cache hit or miss
	_, exists := c.nodes[addr]

	// if this is a miss but cache is full then free up space
	if !exists && len(c.nodes) >= c.maxSize {

		oldTmp := c.oldest
		c.oldest = c.nodes[c.oldest].newer
		c.nodes[c.oldest].older = nilKey
		delete(c.nodes, oldTmp)
	}

	// a cache miss, need new node and update cache size
	if !exists {
		c.nodes[addr] = &node[K,V]{
			blob: value,
			newer: nilKey,
			older: nilKey,
		}
	} else {
		// cache hit, just update the blob
		c.nodes[addr].blob = value
	}

	c.moveToHead(addr, !exists)

	return nil
}

// Internal(private)used in package only functions

// stick the node in list at the most recent position, it is assumed
// the list is not past max size
func (c *Cache[K,V]) moveToHead(current K, isNew bool) {
	if c.newest == current {
		// nothing to do
		return
	}

	var nilKey K

	if c.newest == nilKey && c.oldest == nilKey {
		// cache is empty
		c.newest = current
		c.oldest = current
		return
	}

	// unlink first if not a new node, i.e., overwrite
	if !isNew {
		if c.oldest == current {
			// is it at the tail
			c.oldest = c.nodes[current].newer
			c.nodes[c.oldest].older = nilKey
		} else {
			// not at the head and not at the tail
			older := c.nodes[current].older
			newer := c.nodes[current].newer
			c.nodes[older].newer = newer
			c.nodes[newer].older = older
		}
	}

	// node is unlinked, add to the most recent end
	c.nodes[c.newest].newer = current
	c.nodes[current].older = c.newest
	c.newest = current
}

// print in order of new to old
func (c *Cache[K,V])Print(){
	var nilKey K

	// Header
	fmt.Printf("Cache max: %d sz: %d newest: %v oldest: %v\n", c.maxSize, len(c.nodes), c.newest, c.oldest)
	
	// Unordered form, range over the whole map
	fmt.Println("Raw")
	for addr, blob := range c.nodes {
		fmt.Printf("  Node %v body %v newer %v older %v\n",addr,blob,c.nodes[addr].newer,c.nodes[addr].older)
	}

	// Traverse from most recent to least recent
	fmt.Println("Traverse from most to least recently used")
	idx := 0
	for addr := c.newest; addr != nilKey; addr = c.nodes[addr].older {
		fmt.Printf("  Node: %v blob: %v newer: %v older: %v\n",addr,c.nodes[addr].blob,c.nodes[addr].newer,c.nodes[addr].older)
		idx++

		if idx > c.maxSize {
			fmt.Println("Error in structure of cache")
			return		
		}
	}
}