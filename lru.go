package main

import "fmt"


type node struct {
	newer string
	older string
	blob string
}


type cache struct {
	nodes map[string]*node 
	oldest string
	newest string
	size int
	maxSize int
}

func main(){

	testCache := cache{
		nodes:   make(map[string]*node),
		oldest:  "",
		newest:  "",
		size:    0,
		maxSize: 3,
	}

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


func (c *cache)Read(addr string) (string, bool) {
	
	// empty cache
	if c.newest == "" || c.oldest == "" || c.maxSize <= 1 {
		return "", false	
	}
	
	val, ok := c.nodes[addr]

	// value does not exist
	if !ok {
		return "", false
	} 

	// update the lru for the address
	c.update(addr, false)

	return val.blob,ok
}

// true if inserted, false if failed
func (c *cache)Write(addr string, value string) bool {

	if (c.maxSize <= 1) {
		return false
	}

	// check if exists(overwrite) or new
	_,exists := c.nodes[addr]

	// if addr is not already there and too large, need to free up space
	if !exists && (c.size >= c.maxSize) {
		if !c.remove() {
			return false
		}
	}

	// is new node, need to allocate space
	if (!exists) {
		c.nodes[addr] = &node{
			blob: value,
			newer: "",
			older: "",
		}
		c.size++
	} else {
		// update existing
		c.nodes[addr].blob = value
	}

	// update the lru
	c.update(addr, !exists)
	 
	return true
}

// remove the lru, delete the node
func (c *cache)remove() bool{
	// is the cache empty
	if ( c.newest == "" && c.oldest == "" ) || c.size == 0 {
		return false
	}

	// delete the current oldest after the updates
	defer delete(c.nodes, c.oldest)

	c.oldest = c.nodes[c.oldest].newer
	c.nodes[c.oldest].older = ""
	c.size--

	return true
	}

// stick the node in list at the most recent position, it is assumed
// the list is not past max size
func (c *cache)update(current string, new bool) {

	if c.newest == current{
		// nothing to do
		return
	}

	if !new && c.oldest == c.newest{
		// nothing to do, single node
		return
	}

	if c.newest == "" && c.oldest == "" {
		// cache is empty
		c.newest = current
		c.oldest = current
		return
	}

	// unlink first if not a new node, ie overwrite
	if new {
		// already unlinked, add it to the head
//		c.addToNewest(current)
	} else {
		if (c.oldest == current) {
			// is it at the tail
			newer := c.nodes[current].newer
			c.oldest = newer
			c.nodes[newer].older = ""
		} else {
			// not at the head and not at the tail			
			older := c.nodes[current].older
			newer := c.nodes[current].newer
			c.nodes[older].newer = newer
			c.nodes[newer].older = older
		}
//		c.addToNewest(current)
	}

	c.nodes[c.newest].newer = current
	c.nodes[current].older = c.newest
	c.newest = current

}

// add a new (unlinked) node to the head
//func (c *cache)addToNewest(current string){

//	if c.newest == "" || c.oldest == "" {
//		// cache is empty
//		c.newest = current
//		c.oldest = current
//		return
//	}

//	c.nodes[c.newest].newer = current
//	c.nodes[current].older = c.newest
//	c.newest = current

//}

// print in order of new to old
func (c *cache)printCache(){

	idx := 0

	fmt.Printf("Cache max: %d sz: %d newest: %s oldest: %s\n", c.maxSize, c.size, c.newest, c.oldest)
	
	fmt.Println("Raw")
	for addr, blob := range c.nodes {
		fmt.Printf("  Node %v body %v newer %v older %v\n",addr,blob,c.nodes[addr].newer,c.nodes[addr].older)
	}

	fmt.Println("Traverse")
	for addr := c.newest; addr != ""; addr = c.nodes[addr].older {
		fmt.Printf("  Node: %v blob: %v newer: %v older: %v\n",addr,c.nodes[addr].blob,c.nodes[addr].newer,c.nodes[addr].older)
		idx++

		if idx > c.maxSize {
			fmt.Println("Error in structure of cache")
			panic("Error in structure of cache")
		}
	}
}