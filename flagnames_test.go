package flagnames

import (
	"flag"
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
			args:        []string{"myprog"},
			wantError:   false,
			wantVerbose: false,
			wantId:      0,
			wantItem:    0,
			wantPrefix:  "",
			wantNArg:    0,
		},
		{
			// Just 3 args
			args:        []string{"myprog", "a", "b", "c"},
			wantError:   false,
			wantVerbose: false,
			wantId:      0,
			wantItem:    0,
			wantPrefix:  "",
			wantNArg:    3,
		},
		{
			// -v --> -verbose, -p=myprefix --> -prefix
			args:        []string{"myprog", "-v", "-p=myprefix", "a", "b", "c"},
			wantError:   false,
			wantVerbose: true,
			wantId:      0,
			wantItem:    0,
			wantPrefix:  "myprefix",
			wantNArg:    3,
		},
		{
			// Same with separated flag/value
			args:        []string{"myprog", "-v", "-p", "myprefix", "a", "b", "c"},
			wantError:   false,
			wantVerbose: true,
			wantId:      0,
			wantItem:    0,
			wantPrefix:  "myprefix",
			wantNArg:    3,
		},
		{
			// -id and -it are not ambiguous
			args:        []string{"myprog", "-v", "-p", "myprefix", "-id=19", "-it=62", "a", "b", "c"},
			wantError:   false,
			wantVerbose: true,
			wantId:      19,
			wantItem:    62,
			wantPrefix:  "myprefix",
			wantNArg:    3,
		},
		{
			// Same with separate flag/value for the int flags
			args:        []string{"myprog", "-v", "-p=myprefix", "-id", "19", "-it", "62", "a", "b", "c"},
			wantError:   false,
			wantVerbose: true,
			wantId:      19,
			wantItem:    62,
			wantPrefix:  "myprefix",
			wantNArg:    3,
		},
		{
			// Value with = doesn't get clobbered
			args:        []string{"myprog", "-v", "-p=a=b=c=d", "-id", "19", "-it", "62", "a", "b", "c"},
			wantError:   false,
			wantVerbose: true,
			wantId:      19,
			wantItem:    62,
			wantPrefix:  "a=b=c=d",
			wantNArg:    3,
		},
		{
			// Repeated flags are ok, last one is picked up (same as flag package)
			args:        []string{"myprog", "-v", "-p=a=b=c=d", "-id", "19", "-it", "62", "-p=prefix", "a", "b", "c"},
			wantError:   false,
			wantVerbose: true,
			wantId:      19,
			wantItem:    62,
			wantPrefix:  "prefix",
			wantNArg:    3,
		},
		{
			// Ambiguous flags are not patched, flag.Parse will barf
			args:        []string{"myprog", "-v", "-i", "19", "-it", "62", "-p=prefix", "a", "b", "c"},
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
		PatchFlagSet(fs, &test.args)
		err := fs.Parse(test.args[1:])
		gotError := err != nil

		if gotError != test.wantError {
			t.Errorf("parseSubCmdFlags(%v): gotError=%v, wantError=%v", test.args, gotError, test.wantError)
			continue
		}
		if gotError {
			continue
		}
		if verbose != test.wantVerbose {
			t.Errorf("parseSubCmdFlags(%v): verbose=%v, want %v", test.args, verbose, test.wantVerbose)
		}
		if id != test.wantId {
			t.Errorf("parseSubCmdFlags(%v): id=%v, want %v", test.args, id, test.wantId)
		}
		if item != test.wantItem {
			t.Errorf("parseSubCmdFlags(%v): item=%v, want %v", test.args, item, test.wantItem)
		}
		if prefix != test.wantPrefix {
			t.Errorf("parseSubCmdFlags(%v): prefix=%v, want %v", test.args, prefix, test.wantPrefix)
		}
		if fs.NArg() != test.wantNArg {
			t.Errorf("parseSubCmdFlags(%v): NArg=%v, want %v", test.args, fs.NArg(), test.wantNArg)
		}
	}
}
