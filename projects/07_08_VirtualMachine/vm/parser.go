package vm

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type commandType = int

const (
	CPush commandType = iota
	CPop
	CLabel
	CGoto
	CIfGoto
	CFunction
	CReturn
	CCall

	// arithmetic commands.
	CAdd
	CSub
	CNeg
	CEq
	CGt
	CLt
	CAnd
	COr
	CNot
)

var commandNames = map[string]commandType{
	"push":     CPush,
	"pop":      CPop,
	"label":    CLabel,
	"goto":     CGoto,
	"if-goto":  CIfGoto,
	"function": CFunction,
	"return":   CReturn,
	"call":     CCall,
	"add":      CAdd,
	"sub":      CSub,
	"neg":      CNeg,
	"eq":       CEq,
	"gt":       CGt,
	"lt":       CLt,
	"and":      CAnd,
	"or":       COr,
	"not":      CNot,
}

type vmInstruction struct {
	raw         string
	commandType commandType

	// The meaning (and presence) of these args depend on the instruction,
	// e.g., "add" has no args, for "pop/push seg n" arg1=segment and arg2=index,
	// for "function f n" arg1=function name and arg2=number of arguments.
	arg1 *string
	arg2 *int
}

func parse(filePath string) []vmInstruction {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed opening file %s: %v", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var instructions []vmInstruction
	for scanner.Scan() {
		line := trimLine(scanner.Text())
		if len(line) > 0 {
			instructions = append(instructions, parseInstruction(line))
		}

	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Failed reading file: %v", err)
	}

	return instructions
}

func parseInstruction(line string) vmInstruction {
	parts := strings.Split(line, " ")
	cmdType, ok := commandNames[parts[0]]
	if !ok {
		log.Fatalf("Unrecognized command: %s\n", line)
	}

	var arg1 *string
	var arg2 *int
	// If there are segment and/or index args,
	//just include them regardless of the command to keep it simple
	if len(parts) >= 2 {
		arg1 = &parts[1]
	}
	if len(parts) >= 3 {
		idx, err := strconv.Atoi(parts[2])
		if err != nil {
			log.Fatal("Error parsing index arg for: " + line)
		}
		arg2 = &idx
	}

	return vmInstruction{
		raw:         line,
		commandType: cmdType,
		arg1:        arg1,
		arg2:        arg2,
	}

}

func trimLine(line string) string {
	commentIdx := strings.Index(line, "//")
	if commentIdx >= 0 {
		line = line[:commentIdx]
	}
	return strings.TrimSpace(line)
}
