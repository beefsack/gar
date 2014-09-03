package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/codegangsta/cli"
)

const (
	FlagSourceDir      = "source"
	FlagSourceDirShort = "s"
	FlagRootDir        = "root"
	FlagRootDirShort   = "r"
	FlagInclude        = "include"
	FlagIncludeShort   = "i"
	FlagExclude        = "exclude"
	FlagExcludeShort   = "e"
	FlagGzip           = "gzip"
	FlagGzipShort      = "z"
	FlagExtract        = "extract"
	FlagExtractShort   = "x"
)

func flagsToA(flags ...string) string {
	return strings.Join(flags, ", ")
}

func main() {
	app := cli.NewApp()
	app.Name = "gar"
	app.Usage = "Create Go application archives as standalone executables"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  flagsToA(FlagSourceDir, FlagSourceDirShort),
			Usage: "source location",
		},
		cli.StringFlag{
			Name:  flagsToA(FlagRootDir, FlagRootDirShort),
			Usage: "root directory to source files for archive",
		},
		cli.StringSliceFlag{
			Name:  flagsToA(FlagInclude, FlagIncludeShort),
			Usage: "list of regular expressions of files to include",
			Value: &cli.StringSlice{},
		},
		cli.StringSliceFlag{
			Name:  flagsToA(FlagExclude, FlagExcludeShort),
			Usage: `list of regular expressions of files to exclude, defaults to excluding Go source files and hidden files: "^\.", "/\.", "\.go$"`,
			Value: &cli.StringSlice{
				`^\.`,
				`/\.`,
				`\.go$`,
				fmt.Sprintf("%s$", regexp.QuoteMeta(Output)),
			},
		},
		cli.BoolFlag{
			Name:  flagsToA(FlagGzip, FlagGzipShort),
			Usage: "compress the archive using gzip internally, smaller archive size but longer load times",
		},
		cli.BoolFlag{
			Name:  flagsToA(FlagExtract, FlagExtractShort),
			Usage: "set the resulting archive to extract resources to a temporary directory instead of keeping in memory, better for archives with large resources",
		},
	}
	app.Action = build
	app.Run(os.Args)
}
