// Defines the mapping from assembly instruction to the machine instruction number.
// The C-instruction has 3 parts, comp, dest, and jump, which are
// mapped independently per the tables provided in the book.

import {AInstruction, Instruction} from "./parse";

export function createInstructionWord(instruction: Instruction): number {
    if (instruction instanceof AInstruction) {
        return instruction.value
    } else {
        const dest = lookupMnemonic(instruction.dest, destMapping)
        const comp = lookupMnemonic(instruction.comp, compMapping)
        const jump = lookupMnemonic(instruction.jump, jumpMapping)
        return createCInstructionWord(dest, comp, jump)
    }
}

const destMapping: Record<string, number> = {
    "M": 1,
    "D": 2,
    "MD": 3,
    "A": 4,
    "AM": 5,
    "AD": 6,
    "AMD": 7
}

const jumpMapping: Record<string, number> = {
    'JGT': 1,
    'JEQ': 2,
    'JGE': 3,
    'JLT': 4,
    'JNE': 5,
    'JLE': 6,
    'JMP': 7,
}

const compMapping: Record<string, number> = {
    '0': 0b0101010,
    '1': 0b0111111,
    '-1': 0b0111010,
    'D': 0b0001100,
    'A': 0b0110000,
    'M': 0b1110000,
    '!D': 0b0001101,
    '!A': 0b0110001,
    '!M': 0b1110001,
    '-D': 0b0001111,
    '-A': 0b0110011,
    '-M': 0b1110011,
    'D+1': 0b0011111,
    'A+1': 0b0110111,
    'M+1': 0b1110111,
    'D-1': 0b0001110,
    'A-1': 0b0110010,
    'M-1': 0b1110010,
    'D+A': 0b0000010,
    'D+M': 0b1000010,
    'D-A': 0b0010011,
    'D-M': 0b1010011,
    'A-D': 0b0000111,
    'M-D': 0b1000111,
    'D&A': 0b0000000,
    'D&M': 0b1000000,
    'D|A': 0b0010101,
    'D|M': 0b1010101,
}

function lookupMnemonic(mnemonic: string | undefined, mapping: Record<string, number>): number {
    if (!mnemonic) {
        return 0
    }
    const bits = mapping[mnemonic];
    if (!bits) {
        throw "Invalid dest: " + mnemonic
    }
    return bits;

}

// Concatenates the instruction parts into a single 16 bit instruction word
function createCInstructionWord(dest: number, comp: number, jump: number): number {
    dest = dest || 0;
    jump = jump || 0;

    let word = 0b111 << 13; // three first bits are unused
    word += jump;
    word += dest << 3;
    word += comp << 6;

    return word;
}
