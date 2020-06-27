const predefinedSymbols: Record<string, number> = {
    "SP": 0,
    "LCL": 1,
    "ARG": 2,
    "THIS": 3,
    "THAT": 4,
    "SCREEN": 16384,
    "KBD": 24576,
}
for (let i = 0; i < 16; i++) {
    predefinedSymbols["R" + i] = i
}

// Keeps track of all label- and variable symbols.
export class SymbolTable {

    private variableAdr = 16;
    private readonly symbols: Record<string, number>;

    constructor() {
        this.symbols = Object.assign(predefinedSymbols);
    }

    getOrCreate(symbol: string) {
        if (!(symbol in this.symbols)) {
            // First time this variable is used: define it.
            this.symbols[symbol] = this.variableAdr++;
        }
        return this.symbols[symbol];
    }

    addLabelSymbol(symbol: string, address: number) {
        this.symbols[symbol] = address
    }

}
