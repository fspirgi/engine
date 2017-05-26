package dirwalk

import (
	"os"
	"container/list"
	"path/filepath"
	"strings"
	"unicode"
	"fmt"
)


func Find(dir string, name string) (*list.List,error) {
	files := list.New()
	var walkTest filepath.WalkFunc = func(path string, info os.FileInfo, err error) error {
		if (strings.Contains(filepath.Base(path),name)) {
			files.PushBack(path)
		}
		return err
	}
	filepath.Walk(dir,walkTest)
	return files,nil
}

// finds executable files in <dir>. These are in the form [A-Z][0-9]{4}_.*
func FindExecFiles(dir string) (*list.List,error) {
	files := list.New()
	var walkTest filepath.WalkFunc = func(path string, info os.FileInfo, err error) error {
		validcnt := 0
		for i,cp := range filepath.Base(path) {
			if i == 0 && unicode.IsUpper(cp) {
				// here we have the stream letter! use it
				validcnt++
			} else if i == 1 && unicode.IsDigit(cp) {
				// here is the parallel flag
				validcnt++
			} else if i > 1 && i < 5 && unicode.IsDigit(cp) {
				// here is the parallel flag
				validcnt++
			} else if i == 5 && cp != '_' {
				validcnt++
			} else if i > 5 {
				break
			}
		}
		if validcnt == 5 {
			files.PushBack(path)
		}
		return err
	}
	err := filepath.Walk(dir,walkTest)
	return files,err
}

// arranges executable files in a suitable order
// returns a list of lists of slices
// each element of the list has to be executed one after the other
// each slice of the element is executed in parallel
// each element of the slice has to be excuted one after the other
// go to next element if all processes executed sucessfully, fail otherwise at the end of each top level element
//
// update 26. november: make all lists, two parallel flags x and y [A-Z]xy[0-9]{2}
func FindAndOrderExecFiles(dir string) (*list.List,error) {
	// resulting list
	rlist := list.New()
	files,ok := FindExecFiles(dir)

	// last id and last parallel flag
	var lastId string
	var lastPf int
	var curId string
	var curPf int

	if ok != nil {
		return rlist,ok
	}
	// this should probably be included in FindExecFiles!
	for elem := files.Front(); elem != nil; elem = elem.Next() {
		for i,cp := range filepath.Base(elem.Value.(string)) {
			if i == 0 {
				curId = string(cp)
			} else if i == 1 {
				fmt.Sscanf(string(cp),"%d",&curPf)
			} else {
				break
			}
		}
		if curId != lastId {
			// new top level list element
			elist := list.New()
			// esl := make([]string,0,0)
			// esl = append(esl,elem.Value.(string))
			esl := list.New()
			esl.PushBack(elem.Value.(string))
			elist.PushBack(esl)
			rlist.PushBack(elist)
		} else if (curPf != lastPf) {
			// new slice in last rlist element
			// esl := make([]string,0,0)
			// esl = append(esl,elem.Value.(string))
			esl := list.New()
			esl.PushBack(elem.Value.(string))
			elist := rlist.Back()
			elist.Value.(*list.List).PushBack(esl)
		} else {
			// new entry in last slice
			elist := rlist.Back()
			esl := elist.Value.(*list.List).Back()
			// esl.Value = append(esl.Value.([]string),elem.Value.(string))
			esl.Value.(*list.List).PushBack(elem.Value.(string))
		}
		lastId = curId
		lastPf = curPf
	}
	return rlist,ok
}

func DisplayExecTree (tree *list.List ) {
	var t interface{}
	for telem := tree.Front(); telem != nil; telem = telem.Next() {
		t = telem.Value
		switch et := t.(type) {
		default:
			fmt.Printf("Unexpected type %T\n",et)
		case string:
			fmt.Printf("%s\n",telem.Value.(string))
		case *list.List:
			fmt.Printf("\n")
			DisplayExecTree(telem.Value.(*list.List))
		}
	}
}


// function that finds all files named x in dir y UP the tree, e.g. /a/b/c/d/etc/x, /a/b/c/etc/x ...
// order should be top down afterwards
// func FindUp(startdir string, finddir string, findfile string) (*list.List,error)
func FindUp(startdir string, finddir string, findfile string) (*list.List,error) {
	dirs,err := FindUpDir(startdir,finddir)
	rcs := list.New()
	// this is just the list of possible dirs, not tested against actual file system
	for {
		if cdir := filepath.Dir(startdir); cdir != startdir {
			dirs.PushBack(filepath.Join(cdir,finddir))
			startdir = cdir
		} else {
			break
		}
	}
	// now we look whether we find any files name findfile
	for elem := dirs.Front(); elem != nil; elem = elem.Next() {
		if found,ok := Find(elem.Value.(string),findfile); ok == nil {
			rcs.PushFrontList(found)
		}
	}
	return rcs,err
}

// FindUpDir(string,string) (*list.List,error)
// generates just the list of possible dirs, not tested against actual file system
func FindUpDir(startdir string, finddir string) (*list.List,error) {
        // generate a list of dirs to be scanned
        startdir,err := filepath.Abs(startdir)
        dirs := list.New()
        for {
                if cdir := filepath.Dir(startdir); cdir != startdir {
                        dirs.PushBack(filepath.Join(cdir,finddir))
                        startdir = cdir
                } else {
                        break
                }
        }
	return dirs,err
}
