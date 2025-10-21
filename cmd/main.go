package main

import (
	"fmt"
	"os"
)

func main() {
	if err := Serve(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
