// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Fill.asm

// Runs an infinite loop that listens to the keyboard input.
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel;
// the screen should remain fully black as long as the key is pressed.
// When no key is pressed, the program clears the screen, i.e. writes
// "white" in every pixel;
// the screen should remain fully clear as long as no key is pressed.


(START_LOOP)
	// init counter. We go through each of the 8K words of the screen memory
	@8192
	D=A
	@R1
	M=D

(NEXT_PIXEL)
	// Check keypress for writing white or black
	@R0
	M=0
	@KBD
	D=M
	@WRITE
	D;JEQ
	@R0
	M=-1
	@WRITE
	0;JMP


(WRITE)
	// compute adr to write: screen start adr + current index
	@SCREEN
	D=A
	@R1
	D=D+M
	@R2
	M=D

	// write pixel
	@R0
	D=M
	@R2
	A=M
	M=D

	// count down and loop
	@R1
	MD=M-1
	@START_LOOP
	D;JLT
	@NEXT_PIXEL
	0;JMP

