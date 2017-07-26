package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"regexp"

	"github.com/fspirgi/engine/dirwalk"
	"github.com/fspirgi/engine/exectree"
	"github.com/fspirgi/engine/frcfile"
)

// todo
// password store handling
//
// later
// execution server (to provide remote execution)

func main() {

	var path string
	var stream string
	flag.StringVar(&path, "path", path, "Start directory")
	flag.StringVar(&stream, "stream", ".*", "Stream regexp")
	flag.Parse()

	if path == "" {
		fmt.Printf("Usage: %s --path <directory>\n", os.Args[0])
		os.Exit(0)
	}

	// regex handling for the stream value
	// default is .* (everything)
	rStream := regexp.MustCompile(stream)

	frcfile.SetupPath(path)
	frcfile.PopulateEnvDefault()
	if err := frcfile.PopulateEnvWithRc("testrc"); err != nil {
		log.Fatal(err)
	}

	flist, err := dirwalk.FindAndOrderExecFiles(path, rStream)
	if err != nil {
		log.Fatal(err)
	}

	exectree.ExecuteToplevel(flist)

}
