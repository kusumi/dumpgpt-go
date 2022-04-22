package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	version        [3]int = [3]int{0, 2, 0}
	dumpOptVerbose bool
	dumpOptSymbol  bool
	dumpOptNoalt   bool
)

func printVersion() {
	fmt.Printf("%d.%d.%d\n",
		version[0],
		version[1],
		version[2])
}

func usage(progname string) {
	fmt.Fprintln(os.Stderr, "usage: "+progname+": "+
		"[--verbose] "+
		"[--symbol] "+
		"[--noalt] "+
		"[-v] "+
		"[-h] "+
		"[-u] "+
		"<gpt_image_path>")
}

func main() {
	assertDs()

	progname := os.Args[0]

	if !isLe() {
		fmt.Fprintln(os.Stderr, "big-endian arch unsupported")
		os.Exit(1)
	}

	opt_verbose := flag.Bool("verbose", false, "")
	opt_symbol := flag.Bool("symbol", false, "")
	opt_noalt := flag.Bool("noalt", false, "")
	opt_version := flag.Bool("v", false, "")
	opt_help_h := flag.Bool("h", false, "")
	opt_help_u := flag.Bool("u", false, "")

	flag.Parse()
	args := flag.Args()
	dumpOptVerbose = *opt_verbose
	dumpOptSymbol = *opt_symbol
	dumpOptNoalt = *opt_noalt

	if *opt_version {
		printVersion()
		os.Exit(1)
	}

	if *opt_help_h || *opt_help_u {
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
		panic(err)
	}
	defer fp.Close()

	err = dumpGpt(fp)
	if err != nil {
		panic(err)
	}
}
