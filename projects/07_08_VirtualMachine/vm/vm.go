package vm

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Translates the given .vm file, or all contained .vm files for a directory path,
// writing to output to a single .hack file.
func TranslateFiles(dirOrFilePath string) {
	files, err := getVMFiles(dirOrFilePath)
	if err != nil {
		panic(err)
	}
	if len(files) == 0 {
		panic("No .vm files specified")
	}

	var outfile string
	if len(files) == 1 {
		outfile = strings.ReplaceAll(files[0], ".vm", ".asm")
	} else {
		dirPath, _ := filepath.Abs(dirOrFilePath)
		_, dirName := filepath.Split(dirPath)
		outfile = dirPath + "/" + dirName + ".asm"
	}
	translator := createTranslator(outfile)
	defer translator.done()

	for _, filePath := range files {
		instructions := parse(filePath)
		fmt.Printf("Read %d instructions from %s\n", len(instructions), filePath)
		nameParts := strings.Split(filePath, "/")
		err = translator.translateVMInstructions(instructions, nameParts[len(nameParts)-1])
		if err != nil {
			panic("Failed translating code: " + err.Error())
		}
	}
	fmt.Println("Done: " + outfile)
}

// Returns all .vm files in the given directory path, or the
// path itself if it is a .vm file
func getVMFiles(path string) ([]string, error) {
	fi, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	if !fi.IsDir() {
		if !strings.HasSuffix(path, ".vm") {
			return nil, errors.New("invalid file type")
		}
		return []string{path}, nil
	}

	var vmFiles []string
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	files, err := ioutil.ReadDir(path)
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".vm") {
			vmFiles = append(vmFiles, path+f.Name())
		}
	}
	return vmFiles, nil
}
