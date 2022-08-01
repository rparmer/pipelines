package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) > 0 && args[0] == "crd" {
		fmt.Println("Lookup using crd releaseRefs:")
		crd()
	} else {
		fmt.Println("Lookup using labels:")
		labels()
	}
}
