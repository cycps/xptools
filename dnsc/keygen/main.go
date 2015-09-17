package main

import (
	"fmt"
	"github.com/cycps/xptools/dnsc"
	"os"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: keygen <dns name>")
		os.Exit(1)
	}

	dnsc.Keygen(os.Args[1])
}
