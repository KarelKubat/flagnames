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
- `myprog -v -p=myprefix`: -verbose and -prefix can be abbreviated to their shortest form (`--` also works)
- `myprog -id=1 -it=2`: the shortest form of `id` and `item` is 2 characters

When `flagnames` can't resolve shortened flags to their longer form, then nothing happens - and `flag.Parse()` will fail:

- `myprog -i 1`: will print that flag `-i` is given but not defined (`-i` could mean `-id` or `-item`, and `flagnames` can't resolve it)

The standard flag `-help` is also automatically handled:

- `myprog -h` (or -he, -hel, -help): will call the usual `flag.Usage()` function, like with `myprog -help`

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
```
