package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/manhdaovan/mytail/solutions"
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
	dieIf(err)

	err = solutions.Tail(args.filePaths, args.numLine)
	dieIf(err)
}

func dieIf(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
