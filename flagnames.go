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
	// Avoid parsing empty args. This should not occur irl.
	if len(*actualArgs) == 0 {
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

	for _, arg := range *actualArgs {
		// Anything that doesn't start with a hyphen can be a positional arg, or a flag value.
		// We add it to the reworked args and continue - incase this was a flag arg and more flags
		// follow.
		if !strings.HasPrefix(arg, "-") {
			newArgs = append(newArgs, arg)
			continue
		}

		// Stop examining flags when we see a solitary -- or -.
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
			newFlag := fmt.Sprintf("--%v", definedFlags[hits[0]])
			if len(parts) > 1 {
				newFlag += fmt.Sprintf("=%v", parts[1])
			}
			newArgs = append(newArgs, newFlag)
		} else {
			newArgs = append(newArgs, arg)
		}
		fmt.Println("hits:", hits, "givenflag:", givenFlag, "parts:", parts, "newargs now:", newArgs)
	}

	// Reset the args to the resolved flags.
	*actualArgs = newArgs
	fmt.Println("final newargs:", newArgs)
}

// Patch patches the default (global) flags, witch is the flag.CommandLine.
func Patch() {
	if len(os.Args) > 1 {
		beyondArgs := os.Args[1:]
		PatchFlagSet(flag.CommandLine, &beyondArgs)
		os.Args = beyondArgs
	}
}
