package main

import (
	"flag"
	"fmt"
	"os"
	"path"
)

var (
	version    [3]int = [3]int{0, 2, 2}
	optVerbose bool
	optSymbol  bool
	optNoAlt   bool
)

func printVersion() {
	fmt.Printf("%d.%d.%d\n", version[0], version[1], version[2])
}

func usage(progname string) {
	fmt.Fprintln(os.Stderr, "usage: "+progname+": [<options>] <path>")
	flag.PrintDefaults()
}

func main() {
	assertDs()

	progname := path.Base(os.Args[0])

	if !isLe() {
		fmt.Println("big-endian arch unsupported")
		os.Exit(1)
	}

	optVerboseAddr := flag.Bool("verbose", false, "Enable verbose print")
	optSymbolAddr := flag.Bool("symbol", false, "Print symbol name if possible")
	optNoAltAddr := flag.Bool("noalt", false, "Do not dump secondary header and entries")
	optVersion := flag.Bool("v", false, "Print version and exit")
	optHelp := flag.Bool("h", false, "Print usage and exit")

	flag.Parse()
	args := flag.Args()
	optVerbose = *optVerboseAddr
	optSymbol = *optSymbolAddr
	optNoAlt = *optNoAltAddr

	if *optVersion {
		printVersion()
		os.Exit(1)
	}

	if *optHelp {
		usage(progname)
		os.Exit(1)
	}

	if len(args) < 1 {
		usage(progname)
		os.Exit(1)
	}

	device := args[0]
	fmt.Println(device)
	fmt.Println("")

	fp, err := os.Open(device)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer fp.Close()

	err = dumpGpt(fp)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
