package main

import (
	"fmt"

	flag "github.com/spf13/pflag"
)

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	fmt.Println("anti brute force attack prevention tool")
}
