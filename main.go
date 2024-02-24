package main

import (
	"fmt"
	"github.com/pyromage/lru_example/lru"
)

func main(){

	var testCache lru.Cache[string, string]

	// just some basic testing and console logging
	testCache.New(3)	

	fmt.Println("0")
	testCache.Read("a")
	testCache.Print()

	fmt.Println("1")
	testCache.Write("a","1")
	testCache.Print()

	fmt.Println("2")
	testCache.Write("a","2")
	testCache.Print()

	fmt.Println("3")
	testCache.Write("b","3")
	testCache.Print()

	fmt.Println("4")
	testCache.Write("b","4")
	testCache.Print()

	fmt.Println("5")
	testCache.Write("a","5")
	testCache.Print()

}
