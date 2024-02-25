package lru

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T){
	_,err := NewCache[string, string](0)

	if err == nil {
		t.Errorf("New failed, created a cache of invalid size 0: %v",err)
	}

	_,err = NewCache[string, string](1)
	if err == nil {
		t.Errorf("New failed, created a cache of invalid size 1: %v",err)
	}

	_,err = NewCache[string, string](2)
	if err != nil {
		t.Errorf("New failed, could not create a cache of valid size 2: %v",err)
	}
}

func TestRead(t *testing.T){
	var empty Cache[string,string]

	// test the empty cache
	blob, ok := empty.Read("Uninitialized cache")
	if ok || blob != "" {
		t.Errorf("Test unitialized cache failed")
	}

	// create a cache of size 4
	c , err := NewCache[string, string](4)
	
	if err != nil {
		t.Errorf("New failed, could not create a cache of valid size 4: %v",err)
	}


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
}

func TestWrite(t *testing.T){
	var empty Cache[string,string]

	// test the empty cache
	err := empty.Write("Uninitialized cache","Should not be here")
	
	if err == nil {
		t.Errorf("Test unitialized cache failed")
	}

	// create a cache of size 4
	c,err := NewCache[string, string](4)

	if err != nil {
		t.Errorf("Failed to create a cache of size 4: %v",err)
	}

	// cannot write to the empty key
	var emptyKey string

	err = c.Write(emptyKey,"Should not be here")
	
	if err == nil {
		t.Errorf("Test write to empty key succeeded and it should not have: %v",err)
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
		err = c.Write(v.addr, v.value)
		if (err != nil){
			t.Errorf("Test write failed: %v",err)
		}
	}

	if c.maxSize != 4 || c.oldest != "a" || c.newest != "d" || c.size != 4 {
		t.Errorf("Test failed (exp,res), maxSize:(%d,%d) size: (%d,%d) oldest(%s,%s) newest:(%s,%s)",c.maxSize,4,c.size,4, c.oldest, "a", c.newest, "d")
	}
		
	for i :=0; i < 4 ; i++ {
		if tests[i].value !=  c.nodes[tests[i].addr].blob {
			t.Errorf("Cache entry wrong: idx %d addr:%v exp %s got %s", i,tests[i].addr, tests[i].value, c.nodes[tests[i].addr].blob)
		}
	}
}

func TestReadWrite(t *testing.T){
	var empty Cache[string,string]
	var emptyString string

	// test the empty cache
	blob, ok := empty.Read("Uninitialized cache")
	
	if ok || blob != emptyString {
		t.Errorf("Test read empty cache failed : ok %t blob %v", ok, blob)
	}

	err := empty.Write("Uninitialized cache","Should not be here")
	
	if err == nil {
		t.Errorf("Test write unitialized cache failed")
	}

	const maxCacheSize = 10

	// create a cache
	c, err := NewCache[string, string](maxCacheSize)

	if err != nil {
		t.Errorf("Create cache failed: %v",err)
	}

	for i :=0 ; i < 100 ; i++ {
		addr := fmt.Sprint(i%(maxCacheSize+1))
		value := fmt.Sprint(i+1000)

		// no entry
		r, ok :=  c.Read(addr)
		if  ok {
			 t.Errorf("Test %d read 1 failed, addr:%v value:%v", i, addr, r)
		}

		// write the entry
		err = c.Write(addr, value)
		if (err != nil){
			t.Errorf("Test write failed: %v",err)
		}

		// verify
		r, ok =  c.Read(addr)
		if  !ok || r != value {
			 t.Errorf("Test %d read 2 failed, addr:%v value:%v got blob %v res %t", i, addr, value, r, ok)
		}

		// overwrite the entry
		err = c.Write(addr, value + " overwrite")
		if (err != nil){
			t.Errorf("Test write failed: %v",err)
		}

		// Do 10 reads for each write
		for j := 0; j < 10 ; j++ {
			// verify
			r, ok =  c.Read(addr)
			if  !ok || r != value + " overwrite" {
				t.Errorf("Test %d read 3 failed, addr:%v value:%v got blob %v res %t", i, addr, value + "overwrite", r, ok)
			}	
		}
	}
}

func TestPrint(t *testing.T){
	// test print of an invalid cache

	// create a cache of size 4
	c , err := NewCache[string, string](4)
	
	if err != nil {
		t.Errorf("New failed, could not create a cache of valid size 4: %v",err)
	}

	// populate the cache, unit test so will not use the write
	// make the cache circular
	c.oldest = "d"
	c.newest = "a"
	c.size = 4
	c.maxSize = 4
	c.nodes["a"] = &node[string, string]{older: "b", newer: "d", blob: "0"}
	c.nodes["b"] = &node[string, string]{older: "c", newer: "a", blob: "1"}
	c.nodes["c"] = &node[string, string]{older: "d", newer: "b", blob: "2"}
	c.nodes["d"] = &node[string, string]{older: "a", newer: "c", blob: "3"}

	// if this does not run forever, than so far so good
	c.Print()	
	
	//  invalid cache
	c.nodes["a"] = &node[string, string]{older: "b", newer: "", blob: "0"}
	c.nodes["b"] = &node[string, string]{older: "c", newer: "a", blob: "1"}
	c.nodes["c"] = &node[string, string]{older: "", newer: "", blob: "2"}
	c.nodes["d"] = &node[string, string]{older: "", newer: "c", blob: "3"}

	c.Print()

	// valid cache
	c.nodes["a"] = &node[string, string]{older: "b", newer: "", blob: "0"}
	c.nodes["b"] = &node[string, string]{older: "c", newer: "a", blob: "1"}
	c.nodes["c"] = &node[string, string]{older: "d", newer: "b", blob: "2"}
	c.nodes["d"] = &node[string, string]{older: "", newer: "c", blob: "3"}
	
	c.Print()
}

func BenchmarkCacheReadWrite(b *testing.B){
	const maxCacheSize = 10

	// create a cache
	c, err := NewCache[string, string](maxCacheSize)

	if err != nil {
		b.Errorf("Create cache failed: %v",err)
	}

	for i :=0 ; i < b.N ; i++ {
		addr := fmt.Sprint(i%(maxCacheSize+1))
		value := "Benchmarking"

		// no entry, new address
		r, ok :=  c.Read(addr)
		if  ok {
			 b.Errorf("Test %d read 1 failed, addr:%v value:%v", i, addr, r)
		}

		// write the entry
		c.Write(addr, value) 

		// verify
		r, ok =  c.Read(addr)
		if  !ok || r != value {
			 b.Errorf("Test %d read 2 failed, addr:%v value:%v got blob %v res %t", i, addr, value, r, ok)
		}

		// overwrite the entry
		c.Write(addr, value + " overwrite")

		// Do 10 reads for each write
		for j := 0; j < 10 ; j++ {
			// verify
			r, ok =  c.Read(addr)
			if  !ok || r != value + " overwrite" {
				b.Errorf("Test %d read 3 failed, addr:%v value:%v got blob %v res %t", i, addr, value + "overwrite", r, ok)
			}
		}
	}
}