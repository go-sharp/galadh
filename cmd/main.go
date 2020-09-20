package main

import (
	"fmt"

	"github.com/go-sharp/galadh"
)

func main() {
	fmt.Println("Hello World")
	g := galadh.NewGaladh()

	g.PrintTree(".")
}
