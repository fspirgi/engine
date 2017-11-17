package exectree

import (
	"container/list"
	"fmt"
	"log"
	"os/exec"
	"sync"

	"github.com/fspirgi/engine/frcfile"
)

// esl execute (see esl list in dirwalk)
// esl contains pathnames of files to execute sequentially
func executeFiles(files *list.List, path string) error {
	var rerr error
	var status string
	for filename := files.Front(); filename != nil; filename = filename.Next() {
		if ok := frcfile.PopulateEnvFile(filename.Value.(string), path); ok != nil {
			rerr = ok
		}
		log.Println("START", filename.Value.(string))
		out, err := exec.Command(filename.Value.(string)).Output()
		fmt.Println(string(out))
		if err != nil {
			rerr = err
			status = "FAILED: " + err.Error()
		} else {
			status = "OK"
		}
		log.Printf("STOP %s: %s", filename.Value.(string), status)
	}
	return rerr
}

// elist execute (see elist in dirwalk)
// elist contains esl element to be executed in parallel
func executePar(elist *list.List, path string) error {
	var rerr error
	var waiter = new(sync.WaitGroup)
	for esl := elist.Front(); esl != nil; esl = esl.Next() {
		waiter.Add(1)
		go func(tesl *list.List) {
			err := executeFiles(tesl, path)
			if err != nil {
				rerr = err
			}
			waiter.Done()
		}(esl.Value.(*list.List))
	}
	waiter.Wait()

	return rerr
}

// ExecuteToplevel rlist execute
// just execute one after the other...
// stop on error
func ExecuteToplevel(rlist *list.List, path string) {
	log.Println("PROCESSING STARTED")
	for elist := rlist.Front(); elist != nil; elist = elist.Next() {
		if err := executePar(elist.Value.(*list.List), path); err != nil {
			log.Fatalf("Execution aborted: %s", err)
		}
	}
	log.Println("PROCESSING SUCCESSFULLY ENDED")
}
