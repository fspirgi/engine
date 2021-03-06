package frcfile

// Package to read configuration files of the form
// key = value
// comments are #
// and empty lines are ignored
// returns a pointer to a map with key/value pairs (as strings)
//
// Exported function:
// ReadRc(path string) (*map[string]string,error)
///
// rcfile format:
// key = values
//
// TODO make it a struct with methods
// type Config struct {
//	config map[string]string
// }
//

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fspirgi/engine/dirwalk"
)

// package variable: alread read rc files
type rrcf struct {
	rcfile map[string]bool
	sync.Mutex
}

func newRrcf() (retval *rrcf) {
	retval = new(rrcf)
	retval.rcfile = make(map[string]bool)
	return
}

// var arrcs = make(map[string]bool)
var arrcs = newRrcf()

// ReadRc(path string) (map[string]string,error)
// Reads an rc file and return a key/value map or error
func ReadRc(path string) (map[string]string, error) {
	result := make(map[string]string)
	var err error

	file, err := os.Open(path)
	if err == nil {
		defer file.Close()
		lineNum := 0
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		// read file line by line
		for scanner.Scan() {
			lineNum++
			fields := strings.Fields(scanner.Text())
			length := len(fields)
			isComment := false
			for i := range fields {
				if strings.Contains(fields[i], "#") {
					isComment = true
				}
				if isComment {
					fields[i] = ""
					length--
				}
			}
			if length > 2 {
				if fields[1] != "=" {
					log.Println("WARNING (", path, "): Not a valid entry [", strings.Join(fields[:], " "), "] on line", lineNum)
					continue
				}
				result[fields[0]] = strings.Join(fields[2:], " ")
				continue
			}
			if length <= 2 && length > 0 {
				log.Println("WARNING (", path, "): Not a valid entry [", strings.Join(fields[:], " "), "] on line", lineNum)
			} else {
				continue
			}
		}
		// now expand variables (the ones with $ in front)
		for key, val := range result {
			for {
				if sidx := strings.Index(val, "$"); sidx >= 0 {
					var nval string
					// put the $ at position 0
					vtla := val[sidx:]
					splitted := strings.Fields(vtla)
					nkey := strings.Trim(splitted[0], " $")
					rval := strings.Trim(splitted[0], " ")
					if nkey == key {
						log.Println("WARNING: Endless recursion detected with: ", key)
						nval = strings.Replace(val, rval, "<ERR>", -1)
					} else {
						nval = strings.Replace(val, rval, result[nkey], -1)
					}
					result[key] = nval
					val = nval
				} else {
					break
				}
			}
		}
	}
	return result, err
}

// func ReadAll(*list.List) (map[string]string,error)
// reads every rcfile in list and returns a map containing all entries
func ReadAll(rclist *list.List) (map[string]string, error) {
	result := make(map[string]string)
	var nok error

	// test whether we have already read the file

	for elem := rclist.Front(); elem != nil; elem = elem.Next() {
		if arrcs.Lock(); arrcs.rcfile[elem.Value.(string)] {
			arrcs.Unlock()
			continue
		} else {
			if rcvals, err := ReadRc(elem.Value.(string)); err == nil {
				nok = err
				// merge that into the result map
				for key, val := range rcvals {
					result[key] = val
				}
			}
			arrcs.rcfile[elem.Value.(string)] = true
			arrcs.Unlock()
		}
	}
	return result, nok
}

// FindAndRead(startdir,findfile,finddir string) (map[string]string,error)
func FindAndRead(startdir, finddir, findfile string) (map[string]string, error) {
	result := make(map[string]string)
	var err error
	if flist, ok := dirwalk.FindUp(startdir, finddir, findfile); ok == nil {
		if rl, ok := ReadAll(flist); ok == nil {
			result = rl
		} else {
			log.Println("WARNING: Read rc files:", ok)
		}
	} else {
		log.Println("WARNING: Find rc files:", ok)
		err = ok
	}
	return result, err
}

// FindAndReadRc() (map[string]string,error)
func FindAndReadRc(findfile, path string) (map[string]string, error) {
	//	startdir,err := os.Getwd()
	startdir, err := filepath.Abs(path)
	// this will make sure we also find ../etc
	if filepath.Base(startdir) != "bin" {
		startdir = filepath.Join(startdir, "bin")
	}
	if err != nil {
		log.Println("WARNING: Cannot determine current starting directory for rcfile search:", err)
		return nil, err
	}
	finddir := "etc"
	return FindAndRead(startdir, finddir, findfile)
}

// finds and reads rcfiles and pushes the values into the environment
func PopulateEnvWithRc(findfile, path string) error {
	var ret error
	if env, ok := FindAndReadRc(findfile, path); ok == nil {
		for key, val := range env {
			if ok := os.Setenv(key, val); ok != nil {
				ret = ok
			}
		}
	} else {
		ret = ok
	}
	return ret
}

// Pushes the entries into the environment
func PushEnv(conf map[string]string) error {
	var ret error
	for key, value := range conf {
		nok := os.Setenv(key, value)
		if nok != nil {
			ret = nok
		}
	}
	return ret
}

// DisplayEnv displays the environment, currently unused
func DisplayEnv() {
	out, err := exec.Command("env").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", out)
}

// SetupPath sets up the path
func SetupPath(startdir string) error {
	var tpath []string
	var ret error
	// save previous path
	cpath := os.Getenv("PATH")
	if dlist, ok := dirwalk.FindUpDir(startdir, "bin"); ok == nil {
		for elem := dlist.Front(); elem != nil; elem = elem.Next() {
			tpath = append(tpath, elem.Value.(string))
		}
		tpath = append(tpath, cpath)
		path := strings.Join(tpath, string(os.PathListSeparator))
		if nok := os.Setenv("PATH", path); nok != nil {
			ret = nok
		}
	} else {
		ret = ok
	}
	return ret
}

// func FindAndReadDefault(name string) (map[string]string,error)
func FindAndReadDefault(path string) (map[string]string, error) {
	findfile := filepath.Base(os.Args[0]) + "rc"
	return FindAndReadRc(findfile, path)
}

// find general rc files
func PopulateEnvDefault(path string) error {
	findfile := filepath.Base(os.Args[0]) + "rc"
	return PopulateEnvWithRc(findfile, path)
}

// find all rcfile suitable for a filename
func PopulateEnvFile(filename, path string) error {
	basename := filepath.Base(filename)
	dirname := filepath.Base(filepath.Dir(filename))
	var seen string
	var rerr error
	rcfiles := make([]string, 0, 0)

	rcfiles = append(rcfiles, dirname+"rc")
	for i := 0; i < 5; i++ {
		seen = seen + string(basename[i])
		rcfiles = append(rcfiles, seen+"rc")
	}
	for _, file := range rcfiles {
		if ok := PopulateEnvWithRc(file, path); ok != nil {
			rerr = ok
		}
	}
	return rerr
}
