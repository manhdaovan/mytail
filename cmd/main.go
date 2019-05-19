package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/manhdaovan/mytail/pkg/mytail"
)

type cmdArgs struct {
	filePaths []string
	numLine   uint64
}

func (args cmdArgs) validate() error {
	if len(args.filePaths) == 0 {
		return fmt.Errorf("no file path set")
	}

	for _, fp := range args.filePaths {
		if _, err := os.Stat(fp); err != nil {
			return err
		}
	}

	return nil
}

func parseArgs() cmdArgs {
	args := cmdArgs{}
	flag.Uint64Var(&args.numLine, "n", 10, "number last line")
	flag.Parse()
	args.filePaths = flag.Args()

	return args
}

func main() {
	args := parseArgs()
	err := args.validate()
	if err != nil {
		printHelp(err)
		return
	}

	err = mytail.Tail(args.filePaths, args.numLine)
	if err != nil {
		printHelp(err)
	}
}

func dieIf(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func printHelp(err error) {
	help := `
usage: ./mytail [-n number] file [file2...]

The options are as follows:
    -n number: Number lines will be tailed.
`
	if err != nil {
		fmt.Printf("Error: %+v\n", err)
	}
	fmt.Print(help)
}
