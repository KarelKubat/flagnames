# flagnames

Package `flagnames` is meant to be used with the standard Go package [flag](https://pkg.go.dev/flag). When long flags are defined, `flagnames` allow these to be abbreviated. So you get short flags as well, "for free".

Example:

```go
import (
        "flag"
        "github.com/KarelKubat/flagnames"
)

var (
        verboseFlag = flag.Bool("verbose", false, "increase verbosity")
        idFlag      = flag.Int("id", 0, "ID to process")
        itemFlag    = flag.Int("item", 0, "item number to fetch")
        prefixFlag  = flag.String("prefix", "", "report prefix")
)

func main() {
        // Patch up the short flags into the known flags. That's the only
        // code change you'll need.
        flagnames.Patch()

        flag.Parse()
        // ... handle the flags, perform whatever the program should do
}
```

This allows for the following invocations:

- `myprog -verbose -id=1 -item=2 -prefix=myprefix`: what you'd expect
- `myprog --verbose --id=1 --item=2 --prefix=myprefix`: same, but with `--`
- `myprog -v -p=myprefix`: `-verbose` and `-prefix` can be abbreviated to their shortest form (`--` also works)
- `myprog -id=1 -it=2`: the shortest form of `-id` and `-item` is 2 characters, abbreviating to `-i` won't work since it's ambiguous
- `myprog -p myprefix a b c -p`: the first `-p` is expanded to `--prefix=myprefix`, but the one beyond `a b c` is left as-is; it is a positional argument given that flags stop at `a`
- `myprog -p myprefix -- -p`: the first `-p` is again expanded, the second one not as `--` indicates end-of-flags

When `flagnames` can't resolve shortened flags to their longer form, then no expansion happens - and `flag.Parse()` will fail:

- `myprog -i 1`: will print that flag `-i` is given but not defined (`-i` could mean `-id` or `-item`, and `flagnames` can't resolve it)

The standard flag `-help` is also automatically handled:

- `myprog -h` (or `-he`, `-hel`, `-help`): will call the usual `flag.Usage()` function, like with the standard `myprog -help`

The order of actions is important:
1. First the flags need to be defined
1. Then `flagnames.Patch()` is called (or `flagnames.PatchFlagSet()` for a specific `flag.FlagSet`)
1. Finally `flag.Parse()` is called.

## Synopsis for global flags

```go
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
```

## Synopsis for flag.FlagSet usage

```go
// file: test/m2/main.go
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

```

## Debugging

If the `flagnames` doesn't behave the way you'd expect it to behave, then prior to calling `flagnames.Patch()`, set `flagnames.Debug = true`. That will generate debug messages, stating what's going on and why. If the behavior is a bug, then send me that list along with your invocation and what you would expect `flagnames` to do.

Or better yet, fix the bug and send me a pull request :)

