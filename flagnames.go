// Package flagnames resolves abbreviated flags to their full form, so that the flag package may understand them.
package flagnames

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	Debug bool
)

// dbg is for outputting what happens.
func dbg(f string, a ...interface{}) {
	if Debug {
		fmt.Print("flagnames: ")
		fmt.Printf(f, a...)
		fmt.Println()
	}
}

// PatchFlagSet patches a flag.FlagSet to a known list of flags.
func PatchFlagSet(fs *flag.FlagSet, actualArgs *[]string) {
	// Avoid parsing empty args. This should not occur irl.
	if len(*actualArgs) == 0 {
		dbg("actual args %v are empty, nothing to do", *actualArgs)
		return
	}
	dbg("patching: %v", *actualArgs)

	// Gather up the names of defined flags and ensure that "help" is defined as well.
	// The value bool states whether that flag is a flag.BoolVar. If so, its next argument
	// may be a 'false' or a 'true' and is candidate for consumption. If not, e.g. in the
	// case of a flag.StringVar, its next argument can ALWAYS be consumed.
	definedFlags := map[string]bool{
		"help": true,
	}
	fs.VisitAll(func(f *flag.Flag) {
		// If the default value as string is 'true' or 'false', then it's a boolean flag that may or
		// may not be followed by its value-argument.
		// Otherwise it's a flag that is always followed by its value-argument.
		isBool := f.DefValue == "true" || f.DefValue == "false"
		definedFlags[f.Name] = isBool
		dbg("defined flag %q, is that a bool flag: %v", f.Name, isBool)
	})

	newArgs := []string{}
	parsingFlags := true

	for i := 0; i < len(*actualArgs); i++ {
		arg := (*actualArgs)[i]
		dbg("looking at %v %q, new list so far: %v", i, arg, newArgs)

		// Stop if we're not at a flag or if we see a solitary -- (end-of-flag marker).
		if parsingFlags && (!strings.HasPrefix(arg, "-") || arg == "--") {
			dbg("end of flags seen at %q, stopping flags parsing", arg)
			parsingFlags = false
		}
		if !parsingFlags {
			newArgs = append(newArgs, arg)
			dbg("not parsing flags, taking %q as is", arg)
			continue
		}

		// We're at a flag. Build up a list of possible hits. What full-length flag can this maybe abbreviated flag mean?
		parts := strings.SplitN(arg, "=", 2)
		longCandidates := []string{}
		givenFlag := strings.TrimLeft(parts[0], "-")

		var name string
		for name = range definedFlags {
			if strings.HasPrefix(name, givenFlag) {
				longCandidates = append(longCandidates, name)
				dbg("candidate for %q: %q, is that a bool flag: %v", givenFlag, name, definedFlags[name])
			}
		}
		if len(longCandidates) == 0 {
			dbg("there is no candidate for %q, taking as-is and stopping further parsing", arg)
			newArgs = append(newArgs, arg)
			continue
		}
		isBoolFlag := definedFlags[longCandidates[0]]

		// If we more than 1 candidate for this short flag, then leave it as-is. `flag.Parse()` will complain.
		if len(longCandidates) > 1 {
			newArgs = append(newArgs, arg)
			dbg("there are multiple candidates for %q, further handling not possible", name)
			parsingFlags = false
			continue
		}
		newFlag := fmt.Sprintf("--%v", longCandidates[0])
		if len(parts) > 1 {
			// The flag is in the format --whatever=something, one string. No need to consume the next commandline argument.
			dbg("given flag %q already contains the value, taking that for %q", givenFlag, newFlag)
			newFlag += fmt.Sprintf("=%v", parts[1])
			newArgs = append(newArgs, newFlag)
			continue
		}

		// This was a solitary --flag and not --flag=value.
		if i == len(*actualArgs)-1 {
			dbg("there are no more args to use as value for %q", newFlag)
			newArgs = append(newArgs, newFlag)
			continue
		}

		// The flag is in the format --whatever and we have more args.
		nextArg := (*actualArgs)[i+1]
		dbg("considering %q as value for %q", nextArg, newFlag)
		switch {
		case strings.HasPrefix(nextArg, "-"):
			dbg("candidate value %q starts with a hyphen, not taking it", nextArg)
		case isBoolFlag && nextArg != "true" && nextArg != "false":
			dbg("candidate value %q does not fit bool flag %q", nextArg, newFlag)
		default:
			newFlag += fmt.Sprintf("=%v", nextArg)
			dbg("using candidate value %q as %q", nextArg, newFlag)
			i++
		}
		newArgs = append(newArgs, newFlag)
	}

	// Reset the args to the resolved flags.
	dbg("final list: %v", newArgs)
	*actualArgs = newArgs
}

// Patch patches the default (global) flags, witch is the flag.CommandLine.
func Patch() {
	if len(os.Args) > 1 {
		beyondArgs := os.Args[1:]
		PatchFlagSet(flag.CommandLine, &beyondArgs)
		os.Args = []string{os.Args[0]}
		os.Args = append(os.Args, beyondArgs...)
	}
}
