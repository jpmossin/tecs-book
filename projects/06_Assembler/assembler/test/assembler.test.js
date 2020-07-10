const runAssembler = require("../out/assembler").runAssembler;
const fs = require("fs");

function testPong() {
  try {
    runAssembler("test/Pong.asm");
    const generated = fs.readFileSync("test/Pong.hack").toString().trim();
    const expected = fs.readFileSync("test/Pong.expected.hack").toString().trim();
    if (!(generated === expected)) {
      throw "Pong failed! check diff..";
    }
  } finally {
    fs.unlinkSync("test/Pong.hack");
  }
}

testPong();
