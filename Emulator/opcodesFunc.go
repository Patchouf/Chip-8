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

// Opcode 1NNN - Saut
// Jump to location nnn. The interpreter sets the program counter to nnn.
func (c *Cpu) op1nnn(address uint16) {
	c.Pc = address - 2
}

// Opcode 2NNN - Appel de sous-routine =
// Call subroutine at nnn. The interpreter increments the stack pointer, then puts the current PC on the top
// of the stack. The PC is then set to nnn.
func (c *Cpu) op2nnn(address uint16) {
	// Vérifie que le pointeur de pile (SP) est dans la plage valide (0-15).
	c.StackPush(c.Pc)
	c.Pc = address - 2
}

// Opcode 3XNN - Saut conditionnel (égal) =
// Skip next instruction if Vx = kk. The interpreter compares register Vx to kk, and if they are equal,
// increments the program counter by 2.
func (c *Cpu) op3nnn(opcodeX, opcodeNNN byte) {
	if c.Registre[opcodeX] == opcodeNNN {
		c.Pc += 2
	}
}

// Opcode 4XNN - Saut conditionnel (différent)
// Skip next instruction if Vx != kk. The interpreter compares register Vx to kk, and if they are not equal,
// increments the program counter by 2.
func (c *Cpu) op4nnn(opcodeX, opcodeNN byte) {
	if c.Registre[opcodeX] != opcodeNN {
		c.Pc += 2
	}
}

// Opcode 5XY0 - Saut conditionnel (égalité de registres)
// Skip next instruction if Vx = Vy. The interpreter compares register Vx to register Vy, and if they are equal,
// increments the program counter by 2.
func (c *Cpu) op5nnn(opcodeX, opcodeY byte) {
	if c.Registre[opcodeX] == c.Registre[opcodeY] {
		c.Pc += 2
	}
}

// Opcode 6XNN - Chargement de valeur constante
// Set Vx = kk. The interpreter puts the value kk into register Vx
func (c *Cpu) op6nnn(opcodeX, opcodeNN byte) {
	c.Registre[opcodeX] = opcodeNN
}

// Opcode 7XNN - Ajout de valeur constante =
// Set Vx = Vx + kk. Adds the value kk to the value of register Vx, then stores the result in Vx.
func (c *Cpu) op7nnn(opcodeX, opcodeNN byte) {
	c.Registre[opcodeX] = c.Registre[opcodeX] + opcodeNN

}

// Opcode 8XY0 - Copie de Registre =
// Set Vx = Vy. Stores the value of register Vy in register Vx.
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

// Opcode 8XY2 - Opération ET (bitwise AND) =
// Set Vx = Vx AND Vy. Performs a bitwise AND on the values of Vx and Vy, then stores the result in Vx.
// A bitwise AND compares the corresponding bits from two values, and if both bits are 1, then the same bit
// in the result is also 1. Otherwise, it is 0.
func (c *Cpu) op8nn2(opcodeX, opcodeY byte) {
	c.Registre[opcodeX] &= c.Registre[opcodeY]

}

// Opcode 8XY3 - Opération XOR (bitwise XOR) =
// Set Vx = Vx XOR Vy. Performs a bitwise exclusive OR on the values of Vx and Vy, then stores the result
// in Vx. An exclusive OR compares the corresponding bits from two values, and if the bits are not both the
// same, then the corresponding bit in the result is set to 1. Otherwise, it is 0.
func (c *Cpu) op8nn3(opcodeX, opcodeY byte) {
	c.Registre[opcodeX] ^= c.Registre[opcodeY]

}

// Opcode 8XY4 - Ajout avec retenue =
//Set Vx = Vx + Vy, set VF = carry. The values of Vx and Vy are added together. If the result is greater
//than 8 bits (i.e., ¿ 255,) VF is set to 1, otherwise 0. Only the lowest 8 bits of the result are kept, and stored
//in Vx.

// Vx += Vy
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

// Opcode 8XY5 - Soustraction avec retenue
// Vx -= Vy
// Set Vx = Vx - Vy, set VF = NOT borrow. If Vx ¿ Vy, then VF is set to 1, otherwise 0. Then Vy is
// subtracted from Vx, and the results stored in Vx.
func (c *Cpu) op8nn5(opcodeX, opcodeY byte) {

	if c.Registre[opcodeX] > c.Registre[opcodeY] {
		c.Registre[0xF] = 1
	} else {
		c.Registre[0xF] = 0
	}

	c.Registre[opcodeX] -= c.Registre[opcodeY]
}

// Opcode 8XY6 - Décalage à droite
// Set Vx = Vx SHR 1. If the least-significant bit of Vx is 1, then VF is set to 1, otherwise 0. Then Vx is
// divided 	by 2
func (c *Cpu) op8nn6(opcodeX, opcodeY byte) {

	if c.Registre[opcodeX]&0xF == 1 {
		c.Registre[0xF] = 1
	} else {
		c.Registre[0xF] = 0
	}

	c.Registre[opcodeX] /= 2
}

// Opcode 8XY7 - Soustraction inversée avec retenue =
// Set Vx = Vy - Vx, set VF = NOT borrow. If Vy ¿ Vx, then VF is set to 1, otherwise 0. Then Vx is
// subtracted from Vy, and the results stored in Vx.
func (c *Cpu) op8nn7(opcodeX, opcodeY byte) {

	if c.Registre[opcodeY] > c.Registre[opcodeX] {
		c.Registre[0xF] = 1
	} else {
		c.Registre[0xF] = 0
	}
	c.Registre[opcodeX] = c.Registre[opcodeY] - c.Registre[opcodeX]
}

// Opcode 8XYE - Décalage à gauche =
// Set Vx = Vx SHL 1. If the most-significant bit of Vx is 1, then VF is set to 1, otherwise to 0. Then Vx is
// multiplied by 2
func (c *Cpu) op8nnE(opcodeX, opcodeY byte) {
	if c.Registre[opcodeX]&0x000F == 1 {
		c.Registre[0xF] = 1
	} else {
		c.Registre[0xF] = 0
	}
	c.Registre[opcodeX] *= 2
}

// Opcode 9XY0 - Saut conditionnel (différents registres)=
// Skip next instruction if Vx != Vy. The values of Vx and Vy are compared, and if they are not equal, the
// program counter is increased by 2
func (c *Cpu) op9nn0(opcodeX, opcodeY byte) {
	if c.Registre[opcodeX] != c.Registre[opcodeY] {
		c.Pc += 2
	}
}

// Opcode ANNN - Chargement de l'index (I) =
// Set I = nnn. The value of register I is set to nnn.
func (c *Cpu) opAnnn(address uint16) { // verifier si nnn = opcodennn ou 0
	c.I = address
}

// Opcode BNNN - Saut avec offset =
// Jump to location nnn + V0. The program counter is set to nnn plus the value of V0
func (c *Cpu) opBnnn(address uint16) {
	c.Pc = address + uint16(c.Registre[0])
}

func (c *Cpu) opCxkk(opcodeX, opcodeNN byte) {
	c.Registre[opcodeX] = byte(rand.Int()*256) & opcodeNN
}

// Opcode DXYN - Dessin à l'écran
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

//

// Opcode FX15 - Réglage du retard
func (c *Cpu) opFx15(opcodeX byte) {
	c.Delay_timer = c.Registre[opcodeX]
}

// Opcode FX07 - Chargement du retard
func (c *Cpu) opFx07(opcodeX byte) {
	c.Registre[opcodeX] = c.Delay_timer
}

// Opcode FX55 - Sauvegarde des registres
func (c *Cpu) opFx55(opcodeX byte) {
	for i := byte(0); i <= opcodeX; i++ {
		c.Memory[c.I+uint16(i)] = c.Registre[i]
	}
}

//Fills V0 to VX with values from memory starting at address I. I is then set to I + x + 1.

func (c *Cpu) opFx65(opcodeX byte) {
	for i := byte(0); i <= opcodeX; i++ {
		c.Registre[i] = c.Memory[c.I+uint16(i)]
	}
}

//Set I = I + Vx. The values of I and Vx are added, and the results are stored in I.

func (c *Cpu) opFx1E(opcodeX byte) {
	c.I += uint16(c.Registre[opcodeX])
}
