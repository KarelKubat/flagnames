// file: test/m1/main.go
package main

import (
	"flag"
	"fmt"

	"github.com/KarelKubat/flagnames"
)

var (
	verboseFlag = flag.Bool("verbose", false, "increase verbosity")
	idFlag      = flag.Int("id", 0, "ID to process")
	itemFlag    = flag.Int("item", 0, "item number to fetch")
	prefixFlag  = flag.String("prefix", "", "report prefix")
)

func main() {
	// Trace what's happening.
	flagnames.Debug = true
	// Patch up the short flags into the known flags and parse.
	flagnames.Patch()
	flag.Parse()

	// What have we got?
	fmt.Println("Flags:")
	fmt.Println("  verbose =", *verboseFlag)
	fmt.Println("  id =     ", *idFlag)
	fmt.Println("  item =   ", *itemFlag)
	fmt.Println("  prefix = ", *prefixFlag)
	for _, arg := range flag.Args() {
		fmt.Println("Positional argument:", arg)
	}
}
