package flagnames

import (
	"flag"
	"fmt"
	"testing"
)

func TestPatchFlagSet(t *testing.T) {
	for _, test := range []struct {
		args        []string
		wantError   bool
		wantVerbose bool
		wantId      int
		wantItem    int
		wantPrefix  string
		wantNArg    int
	}{
		{
			// No flags, no action
			args:        []string{},
			wantError:   false,
			wantVerbose: false,
			wantId:      0,
			wantItem:    0,
			wantPrefix:  "",
			wantNArg:    0,
		},
		{
			// Just 3 args
			args:        []string{"a", "b", "c"},
			wantError:   false,
			wantVerbose: false,
			wantId:      0,
			wantItem:    0,
			wantPrefix:  "",
			wantNArg:    3,
		},
		{
			// -v --> -verbose, -p=myprefix --> -prefix
			args:        []string{"-v", "-p=myprefix", "a", "b", "c"},
			wantError:   false,
			wantVerbose: true,
			wantId:      0,
			wantItem:    0,
			wantPrefix:  "myprefix",
			wantNArg:    3,
		},
		{
			// With --
			args:        []string{"--v", "--p=myprefix", "a", "b", "c"},
			wantError:   false,
			wantVerbose: true,
			wantId:      0,
			wantItem:    0,
			wantPrefix:  "myprefix",
			wantNArg:    3,
		},
		{
			// Same with separated flag/value
			args:        []string{"-v", "-p", "myprefix", "a", "b", "c"},
			wantError:   false,
			wantVerbose: true,
			wantId:      0,
			wantItem:    0,
			wantPrefix:  "myprefix",
			wantNArg:    3,
		},
		{
			// With --
			args:        []string{"--v", "--p", "myprefix", "a", "b", "c"},
			wantError:   false,
			wantVerbose: true,
			wantId:      0,
			wantItem:    0,
			wantPrefix:  "myprefix",
			wantNArg:    3,
		},
		{
			// -id and -it are not ambiguous
			args:        []string{"-v", "-p", "myprefix", "-id=19", "-it=62", "a", "b", "c"},
			wantError:   false,
			wantVerbose: true,
			wantId:      19,
			wantItem:    62,
			wantPrefix:  "myprefix",
			wantNArg:    3,
		},
		{
			// Same with separate flag/value for the int flags
			args:        []string{"-v", "-p=myprefix", "-id", "19", "-it", "62", "a", "b", "c"},
			wantError:   false,
			wantVerbose: true,
			wantId:      19,
			wantItem:    62,
			wantPrefix:  "myprefix",
			wantNArg:    3,
		},
		{
			// Value with = doesn't get clobbered
			args:        []string{"-v", "-p=a=b=c=d", "-id", "19", "-it", "62", "a", "b", "c"},
			wantError:   false,
			wantVerbose: true,
			wantId:      19,
			wantItem:    62,
			wantPrefix:  "a=b=c=d",
			wantNArg:    3,
		},
		{
			// With --
			args:        []string{"--v", "--p=a=b=c=d", "--id", "19", "--it", "62", "a", "b", "c"},
			wantError:   false,
			wantVerbose: true,
			wantId:      19,
			wantItem:    62,
			wantPrefix:  "a=b=c=d",
			wantNArg:    3,
		},
		{
			// Repeated flags are ok, last one is picked up (same as flag package)
			args:        []string{"-v", "-p=a=b=c=d", "-id", "19", "-it", "62", "-p=prefix", "a", "b", "c"},
			wantError:   false,
			wantVerbose: true,
			wantId:      19,
			wantItem:    62,
			wantPrefix:  "prefix",
			wantNArg:    3,
		},
		{
			// With --
			args:        []string{"--v", "--p=a=b=c=d", "--id", "19", "--it", "62", "-p=prefix", "a", "b", "c"},
			wantError:   false,
			wantVerbose: true,
			wantId:      19,
			wantItem:    62,
			wantPrefix:  "prefix",
			wantNArg:    3,
		},
		{
			// Flag-like args aren't consumed.
			args:        []string{"-v", "a", "-v"},
			wantError:   false,
			wantVerbose: true,
			wantId:      0,
			wantItem:    0,
			wantPrefix:  "",
			wantNArg:    2,
		},
		{
			// Ambiguous flags are not patched, flag.Parse will barf
			args:        []string{"-v", "-i", "19", "-it", "62", "-p=prefix", "a", "b", "c"},
			wantError:   true,
			wantVerbose: true,
			wantId:      19,
			wantItem:    62,
			wantPrefix:  "prefix",
			wantNArg:    3,
		},
		{
			// With --
			args:        []string{"--v", "--i", "19", "--it", "62", "--p=prefix", "a", "b", "c"},
			wantError:   true,
			wantVerbose: true,
			wantId:      19,
			wantItem:    62,
			wantPrefix:  "prefix",
			wantNArg:    3,
		},
	} {
		var verbose bool
		var id int
		var item int
		var prefix string

		// Create a dedicated flagset and define some options for it.
		fs := flag.NewFlagSet("myprog", flag.ContinueOnError)
		fs.BoolVar(&verbose, "verbose", false, "increase verbosity")
		fs.IntVar(&id, "id", 0, "ID to process")
		fs.IntVar(&item, "item", 0, "item number to fetch")
		fs.StringVar(&prefix, "prefix", "", "report prefix")

		// Patch up short flags into the known flags and parse.
		originalArgs := fmt.Sprintf("%v", test.args)
		PatchFlagSet(fs, &test.args)
		err := fs.Parse(test.args)
		gotError := err != nil

		if gotError != test.wantError {
			t.Errorf("parseSubCmdFlags(%v) = %q, gotError=%v, wantError=%v (modified args: %v)", originalArgs, err.Error(), gotError, test.wantError, test.args)
			continue
		}
		if gotError {
			continue
		}
		if verbose != test.wantVerbose {
			t.Errorf("parseSubCmdFlags(%v): verbose=%v, want %v (modified args: %v)", originalArgs, verbose, test.wantVerbose, test.args)
		}
		if id != test.wantId {
			t.Errorf("parseSubCmdFlags(%v): id=%v, want %v (modified args: %v)", originalArgs, id, test.wantId, test.args)
		}
		if item != test.wantItem {
			t.Errorf("parseSubCmdFlags(%v): item=%v, want %v (modified args: %v)", originalArgs, item, test.wantItem, test.args)
		}
		if prefix != test.wantPrefix {
			t.Errorf("parseSubCmdFlags(%v): prefix=%v, want %v (modified args: %v)", originalArgs, prefix, test.wantPrefix, test.args)
		}
		if fs.NArg() != test.wantNArg {
			t.Errorf("parseSubCmdFlags(%v): NArg=%v, want %v (modified args: %v)", originalArgs, fs.NArg(), test.wantNArg, test.args)
		}
	}
}
