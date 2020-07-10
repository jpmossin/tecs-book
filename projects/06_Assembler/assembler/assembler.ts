import {parse} from "./parse"
import {createInstructionWord} from "./codegen"
import {writeFileSync} from "fs";


// Parses a .asm file and writes the corresponding .hack file.
// Since a single assembly file is typically only a few MBs
// we keep everything sync and simple :)
export function runAssembler(inFile: string) {
    if (!inFile.endsWith(".asm")) {
        console.log("Unsupported file type");
        return;
    }

    const asmInstructions = parse(inFile);
    const hackInstructions = asmInstructions
        .map(createInstructionWord)
        .map(word => word.toString(2))
        .map(word => padding[16 - word.length] + word);

    const outFile = inFile.substring(0, inFile.length - 3) + "hack";
    writeFileSync(outFile, hackInstructions.join("\n"))
}

const padding = [''];
for (let i = 1; i < 16; i++) {
    padding[i] = padding[i-1] + '0'
}
