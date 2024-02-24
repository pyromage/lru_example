package lru

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T){
	var c Cache[string,string]

	if c.New(0) == true || c.New(1) == true {
		t.Errorf("New failed, tried to allocate a cache that is too small")
	}

	if c.New(2) == false {
		t.Errorf("New failed to create a cache of size 2")
	}

	if c.New(3) == true {
		t.Errorf("New failed, cache is already initialized")
	}

}

func TestRead(t *testing.T){
	var c Cache[string,string]

	// test the empty cache
	blob, ok := c.Read("Uninitialized cache")
	if ok || blob != "" {
		t.Errorf("Test unitialized cache failed")
	}

	// create a cache of size 4
	c.New(4)

	// populate the cache, unit test so will not use the write
	c.oldest = "d"
	c.newest = "a"
	c.size = 4
	c.maxSize = 4
	c.nodes["a"] = &node[string, string]{older: "b", newer: "", blob: "0"}
	c.nodes["b"] = &node[string, string]{older: "c", newer: "a", blob: "1"}
	c.nodes["c"] = &node[string, string]{older: "d", newer: "b", blob: "2"}
	c.nodes["d"] = &node[string, string]{older: "", newer: "c", blob: "3"}

	type test struct {
		addr string
		exp string
		res bool
	}

	tests := []test{
		{"a","0",true},
		{"b","1",true},
		{"c","2",true},
		{"d","3",true},
		{"invalid","",false},
	}

	for k,v := range tests {
		blob, ok := c.Read(v.addr)
		if blob != v.exp || ok != v.res {
			 t.Errorf("Test %d failed, addr:%v blob:%v exp: %v res: %t exp %t",k, v.addr, blob, v.exp, ok, v.res)
		}
	}

	c.Print()

}

func TestWrite(t *testing.T){
	var c Cache[string,string]

	// test the empty cache
	ok := c.Write("Uninitialized cache","Should not be here")
	
	if ok {
		t.Errorf("Test unitialized cache failed")
	}

	// create a cache of size 4
	c.New(4)

	// cannot write to the empty key
	var emptyKey string

	ok = c.Write(emptyKey,"Should not be here")
	
	if ok {
		t.Errorf("Test write to empty key succeeded")
	}

	type test struct {
		addr string
		value string
		exp bool
	}

	tests := []test{
		{"a","1",true },
		{"b","2",true },
		{"c","3",true },
		{"d","4",true },
	}

	// without read all we can really do is look for panics and do a final compare
	for _,v := range tests {
		c.Write(v.addr, v.value)
	}

	if c.maxSize != 4 || c.oldest != "a" || c.newest != "d" || c.size != 4 {
		t.Errorf("Test failed (exp,res), maxSize:(%d,%d) size: (%d,%d) oldest(%s,%s) newest:(%s,%s)",c.maxSize,4,c.size,4, c.oldest, "a", c.newest, "d")
	}
		
	for i :=0; i < 4 ; i++ {
		if tests[i].value !=  c.nodes[tests[i].addr].blob {
			t.Errorf("Cache entry wrong: idx %d addr:%v exp %s got %s", i,tests[i].addr, tests[i].value, c.nodes[tests[i].addr].blob)
		}
	}

	c.Print()

}

func TestReadWrite(t *testing.T){
	var c Cache[string,string]

	// test the empty cache
	blob, ok := c.Read("Uninitialized cache")
	
	if ok || blob != "" {
		t.Errorf("Test read unitialized cache failed")
	}

	ok = c.Write("Uninitialized cache","Should not be here")
	
	if ok {
		t.Errorf("Test write unitialized cache failed")
	}

	const maxCacheSize = 10

	// create a cache
	c.New(maxCacheSize)

	for i :=0 ; i < 100 ; i++ {
		addr := fmt.Sprint(i%(maxCacheSize+1))
		value := fmt.Sprint(i+1000)

		// no entry
		r, ok :=  c.Read(addr)
		if  ok {
			 t.Errorf("Test %d read 1 failed, addr:%v value:%v", i, addr, r)
		}

		// write the entry
		c.Write(addr, value) 

		// verify
		r, ok =  c.Read(addr)
		if  !ok || r != value {
			 t.Errorf("Test %d read 2 failed, addr:%v value:%v got blob %v res %t", i, addr, value, r, ok)
		}

		// overwrite the entry
		c.Write(addr, value + " overwrite")

		// verify
		r, ok =  c.Read(addr)
		if  !ok || r != value + " overwrite" {
			 t.Errorf("Test %d read 3 failed, addr:%v value:%v got blob %v res %t", i, addr, value + "overwrite", r, ok)
		}
	}

	c.Print()

}