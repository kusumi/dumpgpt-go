dumpgpt-go ([v0.2.2](https://github.com/kusumi/dumpgpt-go/releases/tag/v0.2.2))
========

## About

+ Parse and dump GPT in ASCII text.

+ Go version of [https://github.com/kusumi/dumpgpt](https://github.com/kusumi/dumpgpt).

## Requirements

go 1.18 or above

## Build

    $ make

## Usage

    $ ./dumpgpt-go
    usage: dumpgpt-go: [<options>] <path>
      -h    Print usage and exit
      -noalt
            Do not dump secondary header and entries
      -symbol
            Print symbol name if possible
      -v    Print version and exit
      -verbose
            Enable verbose print
