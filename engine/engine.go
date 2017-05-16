package main

import (
	"fmt"
	"os"
	"github.com/fspirgi/frcfile"
	"github.com/fspirgi/dirwalk"
	"github.com/fspirgi/exectree"
	"log"
)

// todo
// password store handling
//
// later
// execution server (to provide remote execution)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <directory>\n",os.Args[0])
		os.Exit(0)
	}

//	result,err := frcfile.FindAndReadRc("testrc")
//	if err != nil {
//		fmt.Println("Error occured: ", err)
//		os.Exit(1)
//	}
//
//	defrc,err := frcfile.FindAndReadDefault()
//	if err != nil {
//		fmt.Println("Error occured: ", err)
//		os.Exit(1)
//	}
//	for key, val := range defrc {
//		fmt.Println("defrc",key,val)
//		result[key] = val
//	}
//	for key, val := range result {
//		fmt.Println(key,val)
//	}

	frcfile.SetupPath(os.Args[1])
	frcfile.PopulateEnvDefault()
	if err := frcfile.PopulateEnvWithRc("testrc"); err != nil {
		log.Fatal(err)
	}

	flist,err := dirwalk.FindAndOrderExecFiles(os.Args[1])
	if err != nil {
		fmt.Println("Error occured: ", err)
		os.Exit(1)
	}

	exectree.ExecuteToplevel(flist)

}

