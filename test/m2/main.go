// file: test/m2/main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/KarelKubat/flagnames"
)

func main() {
	c, err := parseSubCmdFlags(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	// What have we got?
	fmt.Println("Flags:")
	fmt.Println("  verbose =", c.verbose)
	fmt.Println("  id =     ", c.id)
	fmt.Println("  item =   ", c.item)
	fmt.Println("  prefix = ", c.prefix)
	for _, arg := range flag.Args() {
		fmt.Println("Positional argument:", arg)
	}
}

type subCmd struct {
	verbose bool
	id      int
	item    int
	prefix  string
}

func parseSubCmdFlags(args []string) (*subCmd, error) {
	cmdData := &subCmd{}
	// Create a dedicated flagset and define some options for it.
	fs := flag.NewFlagSet("myprog", flag.ContinueOnError)
	fs.BoolVar(&cmdData.verbose, "verbose", false, "increase verbosity")
	fs.IntVar(&cmdData.id, "id", 0, "ID to process")
	fs.IntVar(&cmdData.item, "item", 0, "item number to fetch")
	fs.StringVar(&cmdData.prefix, "prefix", "", "report prefix")

	// Patch up short flags into the known flags and parse.
	flagnames.PatchFlagSet(fs, &args)
	if err := fs.Parse(args[1:]); err != nil {
		return nil, err
	}

	return cmdData, nil
}
