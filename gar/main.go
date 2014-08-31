package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "gar"
	app.Usage = "Create Go application archives as standalone executables"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "src",
			Usage: "source location",
		},
		cli.StringFlag{
			Name:  "root",
			Usage: "root directory to source files for archive",
		},
		cli.StringSliceFlag{
			Name:  "include",
			Usage: "list of regular expressions of files to include",
			Value: &cli.StringSlice{},
		},
		cli.StringSliceFlag{
			Name:  "exclude",
			Usage: `list of regular expressions of files to exclude, defaults to excluding Go source files and hidden files: "^\.", "/\.", "\.go$"`,
			Value: &cli.StringSlice{
				`^\.`,
				`/\.`,
				`\.go$`,
				fmt.Sprintf("%s$", regexp.QuoteMeta(Output)),
			},
		},
		cli.BoolFlag{
			Name:  "gzip",
			Usage: "compress the archive using gzip internally, smaller archive size but longer load times",
		},
	}
	app.Action = build
	app.Run(os.Args)
}
