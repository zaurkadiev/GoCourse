package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {

	var pathTale string
	for pos, char := range path {
		if char == '~' {
			pathTale = path[pos:]
			path = path[:pos]
			break
		}
	}

	file, err := os.Open(path)

	if err != nil {
		return err
	}

	// get files in directory
	fileInf, err1 := file.Readdir(100)

	if err1 != nil {
		return err
	}

	// sorting files
	sort.Slice(fileInf, func(i, j int) bool {
		return fileInf[i].Name() < fileInf[j].Name()
	})

	// delete files if printFiles = false
	if !printFiles {
		for i := 0; i < len(fileInf); i++ {
			if !fileInf[i].IsDir() {
				fileInf = append(fileInf[:i], fileInf[i+1:]...)
				i--
			}
		}
	}

	// files/dirs in folder
	for idx, val := range fileInf {
		var gr string
		level := strings.Count(path, "/")
		var isLast bool

		for i := 0; i < level; i++ {
			if s := strings.Split(pathTale, "~"); contains(s, strconv.Itoa(i)) {
				gr += "\t"
			} else {
				gr += "│\t"
			}
		}

		switch {
		case idx == len(fileInf)-1: // last item in directory
			gr += "└───"
			isLast = true
		default:
			gr += "├───"
		}

		// edit graphics

		if val.IsDir() {
			str := fmt.Sprintf("%v%v\n", gr, val.Name())
			fmt.Fprintf(out, str)
		} else {
			si := strconv.FormatInt(val.Size(), 10) + "b"
			if size := val.Size(); size == 0 {
				si = "empty"
			}
			str := fmt.Sprintf("%v%v (%v)\n", gr, val.Name(), si)
			fmt.Fprintf(out, str)
		}

		// add or not & to the path
		if isLast {
			dirTree(out, path+"/"+val.Name()+pathTale+"~"+strconv.Itoa(level), printFiles)
		} else {
			dirTree(out, path+"/"+val.Name()+pathTale, printFiles)
		}
	}
	return nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
