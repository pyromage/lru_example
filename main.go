package main

import (
	"fmt"
	"log"

	"github.com/pyromage/lru_example/lru"
)

func main(){

	// just some basic testing and console logging
	testCache,err := lru.NewCache[string, string](3)
	
	if err != nil {
		log.Fatalf("Faled to create cache: %v",err)
	}

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
