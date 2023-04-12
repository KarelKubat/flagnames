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
	// We need at least SOME args: the program main and one other.
	if len(*actualArgs) < 2 {
		return
	}
	// Gather up the names of defined flags.
	definedFlags := []string{}
	fs.VisitAll(func(f *flag.Flag) {
		definedFlags = append(definedFlags, f.Name)
	})
	definedFlags = append(definedFlags, "help")

	newArgs := []string{}
	parsingFlags := true

	for i, arg := range *actualArgs {
		// Skip the program name.
		if i == 0 {
			continue
		}

		// Stop examining flags when:
		// - We see a solitary -- or -
		// - We see a positional argument
		// Add the the new args if we're already not examining flags.
		if arg == "--" || arg == "-" {
			parsingFlags = false
		}
		if !parsingFlags {
			newArgs = append(newArgs, arg)
			continue
		}

		// We're at a flag. Build up a list of possible hits. What full-length flag can this maybe abbreviated flag mean?
		parts := strings.SplitN(arg, "=", 2)
		hits := []int{}
		givenFlag := strings.TrimLeft(parts[0], "-")
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
