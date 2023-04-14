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
	if err := parseSubCmdFlags(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func parseSubCmdFlags(args []string) error {
	// Create a dedicated flagset and define some options for it.
	fs := flag.NewFlagSet("myprog", flag.ContinueOnError)

	var verboseFlag bool
	fs.BoolVar(&verboseFlag, "verbose", false, "increase verbosity")

	var IDFlag int
	fs.IntVar(&IDFlag, "id", 0, "ID to process")

	var itemFlag int
	fs.IntVar(&itemFlag, "item", 0, "item number to fetch")

	var prefixFlag string
	fs.StringVar(&prefixFlag, "prefix", "", "report prefix")

	// Patch up short flags into the known flags and parse.
	flagnames.Debug = true
	flagnames.PatchFlagSet(fs, &args)
	if err := fs.Parse(args); err != nil {
		return err
	}

	fmt.Println("verbose =", verboseFlag)
	fmt.Println("id      =", IDFlag)
	fmt.Println("item    =", itemFlag)
	fmt.Println("prefix  =", prefixFlag)

	return nil
}
