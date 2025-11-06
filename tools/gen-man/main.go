package main

import (
	"fmt"
	"os"
)

func main() {
	manPage := generateManPage()
	fmt.Print(manPage)
	os.Exit(0)
}
