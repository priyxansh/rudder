package main

import (
	"fmt"
	"os"
)

var Version = "0.1.0-dev"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Printf("rudder version %s\n", Version)
		return
	}
	fmt.Printf("rudder daemon %s (CLI skeleton)\n", Version)
}
