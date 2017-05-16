package exectree

import (
	"os/exec"
	"fmt"
	"container/list"
	"log"
	"github.com/fspirgi/engine/frcfile"
)


// esl execute (see esl list in dirwalk)
// esl contains pathnames of files to execute sequentially

func executeFiles(files *list.List) error {
	var rerr error
	var status string
	for filename := files.Front(); filename != nil; filename = filename.Next() {
		if ok := frcfile.PopulateEnvFile(filename.Value.(string)); ok != nil {
			rerr = ok
		}
		log.Println("START",filename.Value.(string))
		out,err := exec.Command(filename.Value.(string)).Output()
		fmt.Println(string(out))
		if err != nil {
			rerr = err
			status = "FAILED: " + err.Error()
		} else {
			status = "OK"
		}
		log.Printf("STOP %s: %s",filename.Value.(string), status)
	}
	return rerr
}

// elist execute (see elist in dirwalk)
// elist contains esl element to be executed in parallel

func executePar(elist *list.List) error {
	var rerr error
	channels := make([]chan bool,0,0)
	ccnt := 0

	for esl := elist.Front(); esl != nil; esl = esl.Next() {
		channels = append(channels,make(chan bool,1))
		go func(tesl *list.List, ret chan bool) {
			err := executeFiles(tesl)
			if err != nil {
				rerr = err
			}
			ret <- true
		}(esl.Value.(*list.List), channels[ccnt])
		ccnt++
	}
	// this just waits for the first then the second etc, it doesn't care about which could be faster
	// this is probably a poor mans implementation
	for _,channel := range channels {
		_ = <-channel
	}
	return rerr
}

// rlist execute
// just execute one after the other...
// stop on error

func ExecuteToplevel(rlist *list.List)  {
	log.Println("PROCESSING STARTED")
	for elist := rlist.Front(); elist != nil; elist = elist.Next() {
		if err := executePar(elist.Value.(*list.List)); err != nil {
			panic(fmt.Sprintf("Execution aborted: %s", err))
		}
	}
	log.Println("PROCESSING SUCCESSFULLY ENDED")
}
