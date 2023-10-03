package emulator

import (
	"math/rand"
)

func (c *Cpu) StackPush(address uint16) {
	// Vérifiez que le pointeur de pile (SP) est dans la plage valide (0-15).
	c.Stack[c.Sp] = address
	c.Sp++
}
func (c *Cpu) stackPop() uint16 {
	c.Sp--
	address := c.Stack[c.Sp]

	return address
}

// Opcode 00E0 - Effacer l'écran =
// Clear the display.
func (c *Cpu) op00E0() {
	for x := 0; x < 64; x++ {
		for y := 0; y < 32; y++ {
			c.Gfx[x][y] = 0
		}
	}
}

// stock les données
func (c *Cpu) op6XNN(opcodeX, opcodeNNN byte) {
	c.Registre[opcodeX] = opcodeNNN
}

// Opcode 00EE - Retour de sous-routine =
// Return from a subroutine.The interpreter sets the program counter to the address at the top of the stack,
// then subtracts 1 from the stack pointer
func (c *Cpu) op00EE() {
	c.Pc = c.stackPop()

}

func (c *Cpu) op1nnn(address uint16) {
	c.Pc = address - 2
}

func (c *Cpu) op2nnn(address uint16) {
	// Vérifiez que le pointeur de pile (SP) est dans la plage valide (0-15).
	c.StackPush(c.Pc)
	c.Pc = address - 2
}

func (c *Cpu) op3nnn(opcodeX, opcodeNNN byte) {
	if c.Registre[opcodeX] == opcodeNNN {
		c.Pc += 2
	}
}

func (c *Cpu) op4nnn(opcodeX, opcodeNN byte) {
	if c.Registre[opcodeX] != opcodeNN {
		c.Pc += 2
	}
}

func (c *Cpu) op5nnn(opcodeX, opcodeY byte) {
	if c.Registre[opcodeX] == c.Registre[opcodeY] {
		c.Pc += 2
	}
}

func (c *Cpu) op6nnn(opcodeX, opcodeNN byte) {
	c.Registre[opcodeX] = opcodeNN
}

func (c *Cpu) op7nnn(opcodeX, opcodeNN byte) {
	c.Registre[opcodeX] = c.Registre[opcodeX] + opcodeNN

}

func (c *Cpu) op8nn0(opcodeX, opcodeY byte) {
	c.Registre[opcodeX] = c.Registre[opcodeY]

}

// Opcode 8XY1 - Opération OU (bitwise OR) =
// Set Vx = Vx OR Vy. Performs a bitwise OR on the values of Vx and Vy, then stores the result in Vx. A
// bit wise OR compares the corresponding bits from two values, and if either bit is 1, then the same bit in the
// result is also 1. Otherwise, it is 0.
func (c *Cpu) op8nn1(opcodeX, opcodeY byte) {
	c.Registre[opcodeX] |= c.Registre[opcodeY]

}

func (c *Cpu) op8nn2(opcodeX, opcodeY byte) {
	c.Registre[opcodeX] &= c.Registre[opcodeY]

}

func (c *Cpu) op8nn3(opcodeX, opcodeY byte) {
	c.Registre[opcodeX] ^= c.Registre[opcodeY]

}

func (c *Cpu) op8nn4(opcodeX, opcodeY byte) {

	final := c.Registre[opcodeX] + c.Registre[opcodeY]

	if final > 255 {
		c.Registre[0xF] = 1
		c.Registre[opcodeX] = 255
	} else {
		c.Registre[0xF] = 0
		c.Registre[opcodeX] = final
	}
}

func (c *Cpu) op8nn5(opcodeX, opcodeY byte) {

	if c.Registre[opcodeX] > c.Registre[opcodeY] {
		c.Registre[0xF] = 1
	} else {
		c.Registre[0xF] = 0
	}

	c.Registre[opcodeX] -= c.Registre[opcodeY]
}

func (c *Cpu) op8nn6(opcodeX, opcodeY byte) {

	if c.Registre[opcodeX]&0xF == 1 {
		c.Registre[0xF] = 1
	} else {
		c.Registre[0xF] = 0
	}

	c.Registre[opcodeX] /= 2
}

func (c *Cpu) op8nn7(opcodeX, opcodeY byte) {

	if c.Registre[opcodeY] > c.Registre[opcodeX] {
		c.Registre[0xF] = 1
	} else {
		c.Registre[0xF] = 0
	}
	c.Registre[opcodeX] = c.Registre[opcodeY] - c.Registre[opcodeX]
}

func (c *Cpu) op8nnE(opcodeX, opcodeY byte) {
	if c.Registre[opcodeX]&0x000F == 1 {
		c.Registre[0xF] = 1
	} else {
		c.Registre[0xF] = 0
	}
	c.Registre[opcodeX] *= 2
}

func (c *Cpu) op9nn0(opcodeX, opcodeY byte) {
	if c.Registre[opcodeX] != c.Registre[opcodeY] {
		c.Pc += 2
	}
}

func (c *Cpu) opAnnn(address uint16) { // verifier si nnn = opcodennn ou 0
	c.I = address
}

func (c *Cpu) opBnnn(address uint16) {
	c.Pc = address + uint16(c.Registre[0])
}

func (c *Cpu) opCxkk(opcodeX, opcodeNN byte) {

	c.Registre[opcodeX] = byte(rand.Int()*256) & opcodeNN
}

func (clavier *Clavier) GetKey(key byte) bool {
	return clavier.IsPressed[key]
}

// dessine les pixels
func (c *Cpu) opDxyn(opcodeX, opcodeY, opcodeN byte) {
	xval := c.Registre[opcodeX]
	yval := c.Registre[opcodeY]
	c.Registre[0xF] = 0
	var i byte = 0
	for ; i < opcodeN; i++ {
		row := c.Memory[c.I+uint16(i)]
		if erased := c.DrawSprite(xval, yval+i, row); erased {
			c.Registre[0xF] = 1
		}

	}

}
