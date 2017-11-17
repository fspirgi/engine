package main

import (
	"flag"
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
	var stream string
	var path string

	flag.StringVar(&path, "path", path, "Start directory")
	flag.StringVar(&stream, "stream", ".*", "Stream regexp")
	flag.Parse()

	if path == "" {
		var err error
		path, err = os.Getwd()
		if err != nil {
			// something's very wrong here
			log.Fatalf("Can not determine current working directory: %s", err)
		}

	}

	log.Println("Using path: ", path)

	// regex handling for the stream value
	// default is .* (everything)
	rStream := regexp.MustCompile(stream)

	frcfile.SetupPath(path)
	frcfile.PopulateEnvDefault(path)
	if err := frcfile.PopulateEnvWithRc("testrc", path); err != nil {
		log.Fatal(err)
	}

	flist, err := dirwalk.FindAndOrderExecFiles(path, rStream)
	if err != nil {
		log.Fatal(err)
	}

	exectree.ExecuteToplevel(flist, path)

}
