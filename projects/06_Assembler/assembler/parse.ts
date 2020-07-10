import * as fs from "fs";
import {SymbolTable} from "./symtable";

export type Instruction = CInstruction | AInstruction

export class CInstruction {
    constructor(public dest: string | undefined, public comp: string, public jump: string | undefined) {
    }
}

export class AInstruction {
    constructor(public value: number) {
    }
}

/**
 * Maps the given asm file to the corresponding list
 * of parsed instructions.
 */
export function parse(filePath: string): Instruction[] {
    const lines = fs.readFileSync(filePath, {encoding: 'utf8', flag: 'r'})
        .split("\n")
        .map(line => line.includes("//") ? line.substring(0, line.indexOf("//")) : line)
        .map(line => line.trim())
        .filter(Boolean)


    // First do one pass through the file to define all label symbols
    const symbols = new SymbolTable();
    let instructionNum = 0;
    for (const line of lines) {
        if (line.startsWith("(")) {
            const symbol = line.substring(1, line.length - 1);
            symbols.addLabelSymbol(symbol, instructionNum);
        } else {
            instructionNum++;
        }
    }

    return lines
        .filter(line => !line.startsWith("("))
        .map(line => parseInstruction(line, symbols));
}

function parseInstruction(line: string, symbols: SymbolTable): Instruction {
    if (line.startsWith("@")) {
        const value = line.substring(1);
        if (value[0] >= '0' && value[0] <= '9') {
            return new AInstruction(parseInt(value));
        } else {
            return new AInstruction(symbols.getOrCreate(value));
        }
    } else {
        const cParts = line.split(";");
        const jump = cParts.length > 1 ? cParts[1] : undefined;
        const destComp = cParts[0].split("=");
        let dest, comp;
        if (destComp.length == 1) {
            comp = destComp[0]
        } else {
            dest = destComp[0];
            comp = destComp[1];
        }
        return new CInstruction(dest, comp, jump);
    }
}

