package main

import (
	"jpmossin.com/hackvm/vm"
	"os"
)

func main() {
	path := os.Args[1]
	vm.TranslateFiles(path)
}
