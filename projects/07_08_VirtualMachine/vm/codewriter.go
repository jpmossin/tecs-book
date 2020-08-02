// Provides translation from VM code to hack assembly.
// The implementation is per the Standard VM Mapping of the book.
package vm

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

// One translator instance is used for creating one output asm file
// from one or more input vm files.
type translator struct {
	writer *bufio.Writer

	currentFunction string
}

func createTranslator(outfile string) *translator {
	outFile, err := os.Create(outfile)
	if err != nil {
		log.Fatal("Failed opening output file: " + err.Error())
	}
	return &translator{
		writer: bufio.NewWriter(outFile),
	}
}

func (t *translator) translateVMInstructions(instructions []vmInstruction, fileName string) error {
	t.writeASM(bootstrapCode())

	for _, instruction := range instructions {
		var asm string
		cmd := instruction.commandType
		switch cmd {
		case CAdd, CSub, CNeg, CEq, CGt, CLt, CAnd, COr, CNot:
			asm = translateArithmetic(instruction)
		case CPush:
			asm = translatePush(instruction, fileName)
		case CPop:
			asm = translatePop(instruction, fileName)
		case CLabel:
			asm = translateLabel(instruction, t.currentFunction)
		case CGoto:
			asm = translateGoto(instruction, t.currentFunction)
		case CIfGoto:
			asm = translateIfGoto(instruction, t.currentFunction)
		case CCall:
			asm = translateCall(instruction)
		case CFunction:
			t.currentFunction = *instruction.arg1
			asm = translateFunction(instruction)
		case CReturn:
			asm = translateReturn(instruction)
		default:
			log.Fatalln("Unhandled instruction: " + instruction.raw)
		}
		t.writeASM("// " + instruction.raw + "\n")
		t.writeASM(asm)
	}

	return nil
}

func translateReturn(instruction vmInstruction) string {
	// Restore SP for the caller, with return value on top of stack
	asm := popStack + "D=M\n@R14\nM=D\n"  // tmp=pop()
	asm += "@ARG\nD=M\n@SP\nM=D\n"  // sp=arg
	asm += "@R14\nD=M\n" + pushDToStack  // push(tmp)

	asm += "@LCL\nD=M\n@R14\nM=D\n" // store current LCL so that we can get return address below

	// Restore segment registers for the caller; THAT=LCL-1, THIS=LCL-2 etc
	for _, segmentReg := range []string{"THAT", "THIS", "ARG", "LCL"} {
		asm += "@LCL\n"
		asm += "AM=M-1\n"
		asm += "D=M\n"
		asm += "@" + segmentReg + "\nM=D\n"
	}

	// Get return address (old LCL -5) and jump
	asm += "@R14\nD=M\n@5\nA=D-A\nA=M\n0;JMP\n"

	return asm
}

func translateCall(instruction vmInstruction) string {
	// push return address
	returnLabel := nextDynamicLabel()
	asm := "@" + returnLabel + "\nD=A\n" + pushDToStack

	// push segment addresses
	for _, segmentReg := range []string{"LCL", "ARG", "THIS", "THAT"} {
		asm += "@" + segmentReg + "\nA=M\nD=A\n" + pushDToStack
	}

	// Reposition ARG; ARG=SP-5-n (where 5 = the stuff pushed above)
	numArgs := *instruction.arg2
	asm += loadStackAdr +
		"D=A\n@5\nD=D-A\n@" + strconv.Itoa(numArgs) + "\nD=D-A\n" +
		"@ARG\nM=D\n"
	// Reposition LCL;, LCL=SP
	asm += loadStackAdr + "D=A\n@LCL\nM=D\n"

	// goto function
	funcName := *instruction.arg1
	asm += "@" + funcName + "\n"
	asm += "0;JMP\n"

	asm += label(returnLabel)
	return asm
}

func translateFunction(instruction vmInstruction) string {
	funcName := *instruction.arg1
	numLocals := *instruction.arg2

	asm := label(funcName)
	asm += "D=0\n"
	for i := 0; i < numLocals; i++ {
		asm += pushDToStack
	}
	return asm
}

func translateIfGoto(instruction vmInstruction, function string) string {
	target := function + "$" + *instruction.arg1
	return popStack +
		"D=M\n" +
		"@" + target + "\n" +
		"D;JNE\n"
}

func translateGoto(instruction vmInstruction, function string) string {
	target := function + "$" + *instruction.arg1
	asm := "@" + target + "\n"
	asm += "0;JMP\n"
	return asm
}

func translateLabel(instruction vmInstruction, function string) string {
	l := *instruction.arg1
	return label(function + "$" + l)
}

func (t *translator) writeASM(asm string) {
	_, err := t.writer.WriteString(asm)
	if err != nil {
		log.Fatal("Failed writing asm: " + err.Error())
	}
}

func bootstrapCode() string {
	baseAddress := map[string]string{
		"0": "256",  // SP
		//"1": "300",  // LCL
		//"2": "400",  // ARG
		//"3": "3000", // THIS
		//"4": "3010", // THAT
	}
	var asm string
	for reg, base := range baseAddress {
		asm += "@" + base + "\nD=A\n" +
			"@" + reg + "\nM=D\n"
	}

	// todo call og ikke bare jump?
	asm += "@Sys.init\n0;JMP\n"

	return asm
}

func (t *translator) done() {
	err := t.writer.Flush()
	if err != nil {
		log.Fatalln("Failed closing out file, file may be incomplete")
	}
}

func translatePush(iPush vmInstruction, fileName string) string {
	seg := *iPush.arg1
	idx := *iPush.arg2
	idxS := strconv.Itoa(idx)
	var asm string
	switch seg {
	// Load the value to push into D
	case "constant":
		asm = "@" + idxS + "\n" +
			"D=A\n"
	case "local", "argument", "this", "that":
		asm = readHeapSegmentToD(segmentBaseRegister(seg), idx)
	case "pointer", "temp":
		adr := fixedRamSegmentAdr(seg, *iPush.arg2)
		asm = "@" + strconv.Itoa(adr) + "\n" +
			"D=M\n"
	case "static":
		asm = "@" + fileName + "." + idxS + "\n" +
			"D=M\n"
	default:
		panic("Invalid segment for push: " + seg)
	}

	asm += pushDToStack
	return asm
}

func translatePop(iPop vmInstruction, fileName string) string {
	seg := *iPop.arg1
	idx := *iPop.arg2

	var asm string
	switch seg {
	case "local", "argument", "this", "that":
		// First compute adr to write to and save it to a temp reg
		// then pop stack and write to computed adr
		if idx > 0 {
			asm += "@" + strconv.Itoa(idx) + "\n" +
				"D=A\n"
		}
		asm += "@" + segmentBaseRegister(seg) + "\n"
		if idx == 0 {
			asm += "D=M\n"
		} else {
			asm += "D=D+M\n"
		}
		asm += "@R13\n" +
			"M=D\n" + // R13 = adr to write to
			popStack +
			"D=M\n" + // D = value to write
			"@R13\n" +
			"A=M\n" +
			"M=D\n"

	case "pointer", "temp":
		adr := fixedRamSegmentAdr(seg, idx)
		asm += popStack +
			"D=M\n" +
			"@" + strconv.Itoa(adr) + "\n" +
			"M=D\n"

	case "static":
		asm += popStack +
			"D=M\n" +
			"@" + fileName + "." + strconv.Itoa(idx) + "\n" +
			"M=D\n"

	default:
		panic("Invalid segment for pop: " + seg)
	}

	return asm
}

func translateArithmetic(instruction vmInstruction) string {
	op := instruction.commandType
	twoOperands := op == CAdd || op == CSub || op == CEq || op == CGt || op == CLt || op == CAnd || op == COr

	asm := popStack
	if twoOperands {
		asm += "D=M\n" +
			popStack
	}

	switch op {
	case CAdd:
		asm += "M=M+D\n"
	case CSub:
		asm += "M=M-D\n"
	case CEq:
		asm += logicalCmd("JEQ")
	case CLt:
		asm += logicalCmd("JLT")
	case CGt:
		asm += logicalCmd("JGT")
	case CAnd:
		asm += "M=M&D\n"
	case COr:
		asm += "M=M|D\n"
	case CNeg:
		asm += "M=-M\n"
	case CNot:
		asm += "M=!M\n"
	default:
		log.Fatalln("Unhandled command type: " + instruction.raw)
	}

	asm += incSP
	return asm
}

func readHeapSegmentToD(segmentRegister string, index int) string {
	var asm string
	if index > 0 {
		asm = "@" + strconv.Itoa(index) + "\n" +
			"D=A\n"
	}
	asm += "@" + segmentRegister + "\n"
	asm += "A=M\n"
	if index > 0 {
		asm += "A=D+A\n"
	}
	asm += "D=M\n"
	return asm
}

func logicalCmd(jumpCmd string) string {
	// todo: can we skip the jump and just do bit stuff?
	// (No, because we need true=-1 and not just some nonzero value?)
	jumpLabel := nextDynamicLabel()
	doneLabel := nextDynamicLabel()
	return "D=M-D\n" +
		"@" + jumpLabel + "\n" +
		"D;" + jumpCmd + "\n" +
		loadStackAdr +
		"M=0\n" +
		"@" + doneLabel + "\n" +
		"0;JMP\n" +
		label(jumpLabel) +
		loadStackAdr +
		"M=-1;\n" +
		label(doneLabel)
}

var dynamicLabelIdx = 0

func nextDynamicLabel() string {
	dynamicLabelIdx += 1
	return "$DynLabel$" + strconv.Itoa(dynamicLabelIdx)
}

func label(l string) string {
	return "(" + l + ")\n"
}

func segmentBaseRegister(segment string) string {
	switch segment {
	case "local":
		return "LCL"
	case "argument":
		return "ARG"
	case "this":
		return "THIS"
	case "that":
		return "THAT"
	default:
		panic("Not a segment base: " + segment)
	}
}

func fixedRamSegmentAdr(segment string, index int) int {
	switch segment {
	case "pointer":
		return 3 + index
	case "temp":
		return 5 + index
	default:
		panic("Not a fixed ram segment: " + segment)
	}
}

const incSP = "@SP\nM=M+1\n"

const pushDToStack = loadStackAdr +
	"M=D\n" +
	incSP

const loadStackAdr = "@SP\n" +
	"A=M\n"
const popStack = "@SP\nAM=M-1\n"
