package main

import "fmt"


type node[K comparable, V any] struct {
	newer K
	older K
	blob V
}


type cache[K comparable,V any] struct {
	nodes map[K]*node[K,V] 
	oldest K
	newest K
	size int
	maxSize int
}

func main(){

	var testCache cache[string, string]

	// just some basic testing and console logging
	testCache.New(3)	

	fmt.Println("0")
	testCache.Read("a")
	testCache.printCache()

	fmt.Println("1")
	testCache.Write("a","1")
	testCache.printCache()

	fmt.Println("2")
	testCache.Write("a","2")
	testCache.printCache()

	fmt.Println("3")
	testCache.Write("b","3")
	testCache.printCache()

	fmt.Println("4")
	testCache.Write("b","4")
	testCache.printCache()

	fmt.Println("5")
	testCache.Write("a","5")
	testCache.printCache()

}

func (c *cache[K,V])New(size int) bool {
	// already exists
	if c.maxSize != 0 {
		return false
	}
	
	// too small, must be larger than 2
	if (size <= 1) {
		return false
	}
	
	var empty K

	c.nodes = 	make(map[K]*node[K,V])
	c.oldest = empty
	c.newest = empty
	c.size = 0
	c.maxSize = size

	return true
}

func (c *cache[K,V])Read(addr K) (V, bool) {
	var nilKey K
	var nilValue V


	// empty cache or not initialized or trying to cache the nil value
	if c.newest == nilKey || c.oldest == nilKey || c.maxSize <= 1 || addr == nilKey {
		return nilValue, false	
	}
	
	val, ok := c.nodes[addr]

	// value does not exist
	if !ok {
		return nilValue, false
	} 

	// update the lru for the address
	c.updateLRU(addr, false)

	return val.blob,ok
}

// true if inserted, false if failed
func (c *cache[K,V])Write(addr K, value V) bool {
	var nilKey K

	// do not address cache of size 1...
	// adds too much code for no value
	if (c.maxSize <= 1) {
		return false
	}

	// check if exists(overwrite) or new(allocate)
	_, exists := c.nodes[addr]

	// if addr is not already there and size is too large, need to free up space
	if !exists && c.size >= c.maxSize {

		// delete the current oldest after the updates
		defer delete(c.nodes, c.oldest)

		c.oldest = c.nodes[c.oldest].newer
		c.nodes[c.oldest].older = nilKey
		c.size--
		}

	// is new node, need to allocate memory and update cache size
	if (!exists) {
		c.nodes[addr] = &node[K,V]{
			blob: value,
			newer: nilKey,
			older: nilKey,
		}
		c.size++
	} else {
		// update existing
		c.nodes[addr].blob = value
	}

	// update the lru
	c.updateLRU(addr, !exists)

	return true
}

// Internal(private)used in package only functions

// stick the node in list at the most recent position, it is assumed
// the list is not past max size
func (c *cache[K,V]) updateLRU(current K, isNew bool) {
	if c.newest == current {
		// nothing to do
		return
	}

	if !isNew && c.oldest == c.newest {
		// nothing to do, single node
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

	c.nodes[c.newest].newer = current
	c.nodes[current].older = c.newest
	c.newest = current
}

// print in order of new to old
func (c *cache[K,V])printCache(){
	var nilKey K
	idx := 0

	fmt.Printf("Cache max: %d sz: %d newest: %v oldest: %v\n", c.maxSize, c.size, c.newest, c.oldest)
	
	fmt.Println("Raw")
	for addr, blob := range c.nodes {
		fmt.Printf("  Node %v body %v newer %v older %v\n",addr,blob,c.nodes[addr].newer,c.nodes[addr].older)
	}

	fmt.Println("Traverse")
	for addr := c.newest; addr != nilKey; addr = c.nodes[addr].older {
		fmt.Printf("  Node: %v blob: %v newer: %v older: %v\n",addr,c.nodes[addr].blob,c.nodes[addr].newer,c.nodes[addr].older)
		idx++

		if idx > c.maxSize {
			fmt.Println("Error in structure of cache")
			panic("Error in structure of cache")
		}
	}
}