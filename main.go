package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

func getOnlyDirectories(data []os.FileInfo) []os.FileInfo {
	var dirData []os.FileInfo

	for _, element := range data {
		if element.IsDir() {
			dirData = append(dirData, element)
		}
	}

	return dirData
}

func getFileString(element os.FileInfo, prefixes string, isLast bool) string {
	var sizeStr string

	if element.Size() == 0 {
		sizeStr = "empty"
	} else {
		sizeStr = strconv.FormatInt(element.Size(), 10) + "b"
	}

	if isLast {
		return prefixes + "└───" + element.Name() + " (" + sizeStr + ")"
	} else {
		return prefixes + "├───" + element.Name() + " (" + sizeStr + ")"
	}
}

func getDirString(element os.FileInfo, prefixes string, isLast bool) string {
	if isLast {
		return prefixes + "└───" + element.Name()
	} else {
		return prefixes + "├───" + element.Name()
	}
}

func printTree(out io.Writer, path string, printFiles bool, prefixes string) (err error) {
	data, err := ioutil.ReadDir(path)

	if !printFiles {
		data = getOnlyDirectories(data)
	}

	if err != nil {
		return
	}

	if len(data) == 0 {
		return
	}

	for i := 0; i < len(data); i++ {
		element := data[i]

		// This is the last element.
		if i == len(data)-1 {
			if err = processLastElement(element, out, path, prefixes, printFiles); err != nil {
				return
			}

			break
		}

		if element.IsDir() {
			_, _ = fmt.Fprintln(out, getDirString(element, prefixes, false))

			err = printTree(out, path+"/"+element.Name(), printFiles, prefixes+"│\t")

			if err != nil {
				return
			}
		} else {
			_, _ = fmt.Fprintln(out, getFileString(element, prefixes, false))
		}
	}

	return
}

func processLastElement(element os.FileInfo, out io.Writer, path string, prefixes string, printFiles bool) (err error) {
	if element.IsDir() {
		_, _ = fmt.Fprintln(out, getDirString(element, prefixes, true))
		err = printTree(out, path+"/"+element.Name(), printFiles, prefixes+"\t")
	} else {
		_, _ = fmt.Fprintln(out, getFileString(element, prefixes, true))
	}

	return
}

func dirTree(out io.Writer, path string, printFiles bool) (err error) {
	err = printTree(out, path, printFiles, "")

	if err != nil {
		return
	}

	return
}

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
