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
	} else if len(args) > 0 && args[0] == "cm" {
		fmt.Println("Lookup using configmap:")
		cm()
	} else if len(args) > 0 && args[0] == "labels" {
		fmt.Println("Lookup using labels:")
		labels()
	} else {
		fmt.Println("Starting controller:")
		controller()
	}
}
