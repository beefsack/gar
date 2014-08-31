package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/beefsack/gar"
	"github.com/codegangsta/cli"
)

const (
	Output = "out.gar"
)

func build(c *cli.Context) {
	// Flags
	src := c.String("src")
	if src == "" {
		src = "."
	}
	if !strings.ContainsRune("*./", rune(src[0])) {
		log.Fatalf("If src is relative it must start with ./")
	}
	root := c.String("root")
	if root == "" {
		root = "."
	}
	includes := []*regexp.Regexp{}
	for _, i := range c.StringSlice("include") {
		log.Printf("Parsing include flag %s", i)
		r, err := regexp.Compile(i)
		if err != nil {
			log.Fatalf("Unable to parse %s, %v", i, err)
		}
		includes = append(includes, r)
	}
	excludes := []*regexp.Regexp{}
	for _, e := range c.StringSlice("exclude") {
		log.Printf("Parsing exclude flag %s", e)
		r, err := regexp.Compile(e)
		if err != nil {
			log.Fatalf("Unable to parse %s, %v", e, err)
		}
		excludes = append(excludes, r)
	}
	gz := c.Bool("gzip")
	// Build
	goArgs := []string{"build", "-o", Output, src}
	cmd := exec.Command("go", goArgs...)
	cmd.Stdout = os.Stderr
	cmd.Stdout = os.Stdout
	log.Println("Building...")
	if err := cmd.Run(); err != nil {
		log.Fatalf("Build failed, %v", err.Error())
	}
	log.Printf("Build complete, created %s.", Output)
	file, err := os.OpenFile("out.gar", os.O_RDWR|os.O_APPEND, 0)
	if err != nil {
		log.Fatalf("Error opening %s, %v", err.Error())
	}
	defer file.Close()
	// Archive
	w := gar.NewWriter(file)
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if len(includes) > 0 {
			found := false
			for _, i := range includes {
				if i.MatchString(rel) {
					found = true
					break
				}
			}
			if !found {
				return nil
			}
		}
		for _, e := range excludes {
			if e.MatchString(rel) {
				return nil
			}
		}
		log.Printf("Adding %s", rel)
		if err := w.WriteFileAtPath(path, rel); err != nil {
			log.Fatalf("Error adding %s, %v", rel, err)
		}
		return nil
	})
	compression := gar.CompressionNone
	if gz {
		compression = gar.CompressionGzip
	}
	if err := w.Close(compression); err != nil {
		log.Fatalf("Error writing file, %v", err)
	}
}
