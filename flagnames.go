// Package flagnames resolves abbreviated flags to their full form, so that the flag package may understand them.
package flagnames

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// PatchFlagSet patches a flag.FlagSet to a known list of flags.
func PatchFlagSet(fs *flag.FlagSet, actualArgs *[]string) {
	// Gather up the names of defined flags.
	definedFlags := []string{}
	fs.VisitAll(func(f *flag.Flag) {
		definedFlags = append(definedFlags, f.Name)
	})
	definedFlags = append(definedFlags, "help")

	newArgs := []string{}
	parsingFlags := true

	for _, arg := range *actualArgs {
		// Stop examining flags once we've seen --.
		if arg == "--" {
			parsingFlags = false
			newArgs = append(newArgs, arg)
			continue
		}
		// Just add the arg to the new set of args if we've stopped parsing, or the arg isn't a flag.
		if !parsingFlags || arg[0] != '-' {
			newArgs = append(newArgs, arg)
			continue
		}
		// Build up a list of possible hits. What full-length flag can this maybe abbreviated flag mean?
		parts := strings.Split(arg, "=")
		hits := []int{}
		givenFlag := parts[0]
		for givenFlag[0] == '-' {
			givenFlag = givenFlag[1:]
		}
		for index, knownFlag := range definedFlags {
			if strings.HasPrefix(knownFlag, givenFlag) {
				hits = append(hits, index)
			}
		}
		// If we have exactly one match for the given flag, then modify it. Else use whatever was there.
		if len(hits) == 1 {
			newParts := []string{fmt.Sprintf("-%v", definedFlags[hits[0]])}
			for _, p := range parts[1:] {
				newParts = append(newParts, p)
			}
			newArgs = append(newArgs, strings.Join(newParts, "="))
		} else {
			newArgs = append(newArgs, arg)
		}
	}

	// Reset the args to the resolved flags.
	*actualArgs = newArgs
}

// Patch patches the default (global) flags, witch is the flag.CommandLine.
func Patch() {
	PatchFlagSet(flag.CommandLine, &os.Args)
}
