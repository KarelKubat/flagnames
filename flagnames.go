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

	for i := 1; i < len(*actualArgs); i++ {
		arg := (*actualArgs)[i]
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
		parts := strings.Split(arg, "=")
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
