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
	CIf
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
	"if":       CIf,
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

func cmdName(cmd commandType) string {
	for name, t := range commandNames {
		if t == cmd {
			return name
		}
	}
	return "UNKNOWN:" + strconv.Itoa(cmd);
}

type vmInstruction struct {
	commandType commandType

	// The segment and index args are null when an instruction does not include them
	segment *string
	index   *int
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

	var segment *string
	var index *int
	// If there are segment and/or index args,
	//just include them regardless of the command to keep it simple
	if len(parts) >= 2 {
		segment = &parts[1]
	}
	if len(parts) >= 3 {
		idx, err := strconv.Atoi(parts[2])
		if err != nil {
			log.Fatal("Error parsing index arg for: " + line)
		}
		index = &idx
	}

	return vmInstruction{
		commandType: cmdType,
		segment:     segment,
		index:       index,
	}

}

func trimLine(line string) string {
	commentIdx := strings.Index(line, "//")
	if commentIdx >= 0 {
		line = line[:commentIdx]
	}
	return strings.TrimSpace(line)
}
