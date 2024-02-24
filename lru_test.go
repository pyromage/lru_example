package main

import (
	"fmt"
	"testing"
)

func TestRead(t *testing.T){
	c := cache{
		nodes:   make(map[string]*node),
		oldest:  "",
		newest:  "",
		size:    0,
		maxSize: 4,
	}

	// test the empty cache
	blob, ok := c.Read("Empty cache")
	if ok || blob != "" {
		t.Errorf("Test empty cache failed, got blob:%v exp: %v res: %t exp %t", blob, "", ok, false)
	}

	// create a cache
	c.oldest = "d"
	c.newest = "a"
	c.size = 4
	c.maxSize = 4
	c.nodes["a"] = &node{older: "b", newer: "", blob: "0"}
	c.nodes["b"] = &node{older: "c", newer: "a", blob: "1"}
	c.nodes["c"] = &node{older: "d", newer: "b", blob: "2"}
	c.nodes["d"] = &node{older: "", newer: "c", blob: "3"}

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

	c.printCache()

}

func TestWrite(t *testing.T){
	c := cache{
		nodes:   make(map[string]*node),
		oldest:  "",
		newest:  "",
		size:    0,
		maxSize: 4,
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

	for k,v := range tests {
		if c.Write(v.addr, v.value) != v.exp {
			 t.Errorf("Test %d failed, addr:%v vale:%v result %t, expected %t", k, v.addr, v.value, !v.exp, v.exp )
		}
	}

	c.printCache()

}

func TestReadWrite(t *testing.T){
	const maxCacheSize = 10

	c := cache{
		nodes:   make(map[string]*node),
		oldest:  "",
		newest:  "",
		size:    0,
		maxSize: maxCacheSize,
	} 

	for i :=0 ; i < 100 ; i++ {
		addr := fmt.Sprint(i%(maxCacheSize+1))
		value := fmt.Sprint(i+1000)

		// no entry
		r, ok :=  c.Read(addr)
		if  ok {
			 t.Errorf("Test %d read failed, addr:%v value:%v", i, addr, r)
		}

		// write the entry
		if !c.Write(addr, value) {
			 t.Errorf("Test %d write failed, addr:%v value:%v", i, addr, value)
		}

		// verify
		r, ok =  c.Read(addr)
		if  !ok || r != value {
			 t.Errorf("Test %d read 2 failed, addr:%v value:%v got blob %v res %t", i, addr, value, r, ok)
		}

		// overwrite the entry
		if !c.Write(addr, value + " overwrite") {
			 t.Errorf("Test %d overwrite 2 failed, addr:%v value:%v", i, addr, value + "overwrite")
		}

		// verify
		r, ok =  c.Read(addr)
		if  !ok || r != value + " overwrite" {
			 t.Errorf("Test %d.read failed, addr:%v value:%v got blob %v res %t", i, addr, value + "overwrite", r, ok)
		}
	}

	c.printCache()

}