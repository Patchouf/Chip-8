package opcodes

import (
	"fmt"
	"math/rand"
)

type Cpu struct {
	Memory      [4096]byte
	Registre    [16]byte
	I           uint16
	Pc          uint16
	Gfx         [64][32]byte
	Delay_timer byte
	Stack       [16]uint16
	Sp          byte
	Key         [16]byte
	Romlength   uint16
}

type Object struct {
	Clavier [16]byte
}

func (cpu *Cpu) initialiseFont() {
	//0x000-0x1FF - Chip 8 interpreter (contains font set in emu)
	//0x050-0x0A0 - Used for the built in 4x5 pixel font set (0-F)
	//0 0x050
	cpu.Memory[0x050] = 0xF0
	cpu.Memory[0x051] = 0x90
	cpu.Memory[0x052] = 0x90
	cpu.Memory[0x053] = 0x90
	cpu.Memory[0x054] = 0xF0
	// 1
	cpu.Memory[0x055] = 0x20
	cpu.Memory[0x056] = 0x60
	cpu.Memory[0x057] = 0x20
	cpu.Memory[0x058] = 0x20
	cpu.Memory[0x059] = 0x70
	// 2
	cpu.Memory[0x05A] = 0xF0
	cpu.Memory[0x05B] = 0x10
	cpu.Memory[0x05C] = 0xF0
	cpu.Memory[0x05D] = 0x80
	cpu.Memory[0x05E] = 0xF0
	// 3
	cpu.Memory[0x05F] = 0xF0
	cpu.Memory[0x060] = 0x10
	cpu.Memory[0x061] = 0xF0
	cpu.Memory[0x062] = 0x10
	cpu.Memory[0x063] = 0xF0
	// 4
	cpu.Memory[0x064] = 0x90
	cpu.Memory[0x065] = 0x90
	cpu.Memory[0x066] = 0xF0
	cpu.Memory[0x067] = 0x10
	cpu.Memory[0x068] = 0x10
	// 5
	cpu.Memory[0x069] = 0xF0
	cpu.Memory[0x06A] = 0x80
	cpu.Memory[0x06B] = 0xF0
	cpu.Memory[0x06C] = 0x10
	cpu.Memory[0x06D] = 0xF0
	// 6
	cpu.Memory[0x06E] = 0xF0
	cpu.Memory[0x06F] = 0x80
	cpu.Memory[0x070] = 0xF0
	cpu.Memory[0x071] = 0x90
	cpu.Memory[0x072] = 0xF0
	// 7
	cpu.Memory[0x073] = 0xF0
	cpu.Memory[0x074] = 0x10
	cpu.Memory[0x075] = 0x20
	cpu.Memory[0x076] = 0x40
	cpu.Memory[0x077] = 0x40
	// 8
	cpu.Memory[0x078] = 0xF0
	cpu.Memory[0x079] = 0x90
	cpu.Memory[0x07A] = 0xF0
	cpu.Memory[0x07B] = 0x90
	cpu.Memory[0x07C] = 0xF0
	// 9
	cpu.Memory[0x07D] = 0xF0
	cpu.Memory[0x07E] = 0x90
	cpu.Memory[0x07F] = 0xF0
	cpu.Memory[0x080] = 0x10
	cpu.Memory[0x081] = 0xF0
	// A
	cpu.Memory[0x082] = 0xF0
	cpu.Memory[0x083] = 0x90
	cpu.Memory[0x084] = 0xF0
	cpu.Memory[0x085] = 0x90
	cpu.Memory[0x086] = 0x90
	// B
	cpu.Memory[0x087] = 0xE0
	cpu.Memory[0x088] = 0x90
	cpu.Memory[0x089] = 0xE0
	cpu.Memory[0x08A] = 0x90
	cpu.Memory[0x08B] = 0xE0
	// C
	cpu.Memory[0x08C] = 0xF0
	cpu.Memory[0x08D] = 0x80
	cpu.Memory[0x08E] = 0x80
	cpu.Memory[0x08F] = 0x80
	cpu.Memory[0x090] = 0xF0
	// D
	cpu.Memory[0x091] = 0xE0
	cpu.Memory[0x092] = 0x90
	cpu.Memory[0x093] = 0x90
	cpu.Memory[0x094] = 0x90
	cpu.Memory[0x095] = 0xE0
	// E
	cpu.Memory[0x096] = 0xF0
	cpu.Memory[0x097] = 0x80
	cpu.Memory[0x098] = 0xF0
	cpu.Memory[0x099] = 0x80
	cpu.Memory[0x09A] = 0xF0
	// F
	cpu.Memory[0x09B] = 0xF0
	cpu.Memory[0x09C] = 0x80
	cpu.Memory[0x09D] = 0xF0
	cpu.Memory[0x09E] = 0x80
	cpu.Memory[0x09F] = 0x80
}

// Initialisation du cpu
func InitCpu(cpu *Cpu, rombytes []byte) {
	cpu.initialiseFont()
	cpu.loadROM(rombytes)
	cpu.Pc = 0x200 - 2
	// fmt.Println(cpu.Pc)

}

// Update du cpu
func (cpu *Cpu) Update() {
	cpu.Pc += 2
	op1 := cpu.Memory[cpu.Pc]
	op2 := cpu.Memory[cpu.Pc+1]
	// fmt.Println(op1, " ", op2)
	opcode := cpu.uint8ToUint16(op1, op2)
	fmt.Printf("%02x", opcode)
	fmt.Println()
	cpu.decode(opcode)
}

// chargement du rom
func (cpu *Cpu) loadROM(rombytes []byte) {
	cpu.Romlength = uint16(len(rombytes))
	for i, byt := range rombytes {
		cpu.Memory[0x200+i] = byt
	}
}

// Fonction stackPush

// Fonction uint16 to uint8
func (c *Cpu) Uint16ToUint8(n uint16) (uint8, uint8) {
	return uint8(n >> 8), uint8(n & 0x00FF)
}

// unit 8 to uint16
func (c *Cpu) uint8ToUint16(n1 uint8, n2 uint8) uint16 {
	return uint16(uint16(n1)<<8 | uint16(n2))
}

// Fonction uint8 to uint4
func (c *Cpu) Uint8ToUint4(n uint8) (uint8, uint8) {
	return uint8(n >> 4), uint8(n & 0x0F)
}

// func (c *Cpu) DrawSprite(X, Y, height byte) bool {
// 	ScreenWidth := uint16(c.Registre[X])
// 	ScreenHeight := uint16(c.Registre[Y])

// 	c.Registre[0xF] = 0

// 	// Parcourez les lignes du sprite.
// 	for row := byte(0); row < height; row++ {
// 		spriteByte := c.Memory[c.I+uint16(row)]

// 		for bit := byte(0); bit < 8; bit++ {

// 			if (spriteByte & (0x80 >> bit)) != 0 {

// 				x := int(ScreenWidth) + int(bit)
// 				y := int(ScreenHeight) + int(row)

// 				if x < 64 && y < 32 {

// 					index := y*64 + x

// 					if c.Gfx[index][0] == byte(1) {

// 						c.Registre[0xF] = 1
// 					}
// 					for i := 0; i < len(c.Gfx[index]); i++ {
// 						c.Gfx[index][i] ^= byte(1)
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return c.Registre[0xF] == 1
// }

// DrawSprite dessine un sprite à l'écran et renvoie true si un pixel a été effacé
func (c *Cpu) DrawSprite(x byte, y byte, row byte) bool {
	erased := false
	yIndex := y % 32

	for i := x; i < x+8; i++ {
		xIndex := i % 64

		wasSet := c.Gfx[xIndex][yIndex] == 1
		value := row >> (x + 8 - i - 1) & 0x01

		c.Gfx[xIndex][yIndex] ^= value

		if wasSet && c.Gfx[xIndex][yIndex] == 0 {
			erased = true
		}
	}

	return erased
}

// décodage d'un opcode et exécute l'instruction correspondante.
func (c *Cpu) decode(opcode uint16) {
	// Diviser l'opcode en parties individuelles PROBLEME
	opcodeN := byte(opcode>>12) & 0x000F // 4 premiers bits
	opcodeX := byte(opcode>>8) & 0x000F  // Bits 8 à 11
	opcodeY := byte(opcode>>4) & 0x000F  // Bits 4 à 7
	opcodeNNN := opcode & 0x0FFF         // Bits 0 à 11
	opcodeNN := byte(opcode & 0x00FF)    // Bits 0 à 7
	opcodeN4 := byte(opcode & 0x000F)    // 4 derniers bits
	//x := byte(opcodeX & )

	// Utilisez un switch pour gérer chaque opcode
	switch opcode & 0xF000 {
	case 0x0000:
		switch opcode {
		case 0x00E0:

			// Opcode 00E0 - Effacer l'écran =
			// Clear the display.

			c.op00E0()

		case 0x00EE:

			// Opcode 00EE - Retour de sous-routine =
			// Return from a subroutine.The interpreter sets the program counter to the address at the top of the stack,
			// then subtracts 1 from the stack pointer

			c.op00EE()

		default:
			// Gérer les opcodes 0NNN ici (non standard)
		}
	case 0x1000:

		// Opcode 1NNN - Saut =
		// Jump to location nnn. The interpreter sets the program counter to nnn.

		c.op1nnn(uint16(opcodeNNN)) // PTET ERREUR = opcodeN a la place

	case 0x2000:

		// Opcode 2NNN - Appel de sous-routine =
		// Call subroutine at nnn. The interpreter increments the stack pointer, then puts the current PC on the top
		//of the stack. The PC is then set to nnn.

		c.op2nnn(opcode & 0x0FFF)

	case 0x3:

		// Opcode 3XNN - Saut conditionnel (égal) =
		// Skip next instruction if Vx = kk. The interpreter compares register Vx to kk, and if they are equal,
		//increments the program counter by 2.

		c.op3nnn(opcodeX, opcodeNN)

	case 0x4000:

		// Opcode 4XNN - Saut conditionnel (différent)
		// Skip next instruction if Vx != kk. The interpreter compares register Vx to kk, and if they are not equal,
		// increments the program counter by 2.

		c.op4nnn(opcodeX, opcodeNN)

	case 0x5000:

		// Opcode 5XY0 - Saut conditionnel (égalité de registres)
		// Skip next instruction if Vx = Vy. The interpreter compares register Vx to register Vy, and if they are equal,
		// increments the program counter by 2.

		c.op5nnn(opcodeX, opcodeY)

	case 0x6000:

		// Opcode 6XNN - Chargement de valeur constante
		// Set Vx = kk. The interpreter puts the value kk into register Vx

		c.op6nnn(opcodeX, opcodeNN)

	case 0x7000:

		// Opcode 7XNN - Ajout de valeur constante =
		// Set Vx = Vx + kk. Adds the value kk to the value of register Vx, then stores the result in Vx.

		c.op7nnn(opcodeX, opcodeNN)

	case 0x8000:
		// Gérer les opcodes 8XY0 à 8XYE
		switch opcode & 0x000F {

		case 0x0000:

			// Opcode 8XY0 - Copie de Registre =
			//Set Vx = Vy. Stores the value of register Vy in register Vx.

			c.op8nn0(opcodeX, opcodeY)

		case 0x0001:

			// Opcode 8XY1 - Opération OU (bitwise OR) =
			//Set Vx = Vx OR Vy. Performs a bitwise OR on the values of Vx and Vy, then stores the result in Vx. A
			//bit wise OR compares the corresponding bits from two values, and if either bit is 1, then the same bit in the
			// result is also 1. Otherwise, it is 0.

			c.op8nn1(opcodeX, opcodeY)

		case 0x0002:

			// Opcode 8XY2 - Opération ET (bitwise AND) =
			//Set Vx = Vx AND Vy. Performs a bitwise AND on the values of Vx and Vy, then stores the result in Vx.
			//A bitwise AND compares the corresponding bits from two values, and if both bits are 1, then the same bit
			//in the result is also 1. Otherwise, it is 0.

			c.op8nn2(opcodeX, opcodeY)

		case 0x0003:

			// Opcode 8XY3 - Opération XOR (bitwise XOR) =
			//Set Vx = Vx XOR Vy. Performs a bitwise exclusive OR on the values of Vx and Vy, then stores the result
			// in Vx. An exclusive OR compares the corresponding bits from two values, and if the bits are not both the
			// same, then the corresponding bit in the result is set to 1. Otherwise, it is 0.

			c.op8nn3(opcodeX, opcodeY)

		case 0x0004:

			// Opcode 8XY4 - Ajout avec retenue =
			//Set Vx = Vx + Vy, set VF = carry. The values of Vx and Vy are added together. If the result is greater
			//than 8 bits (i.e., ¿ 255,) VF is set to 1, otherwise 0. Only the lowest 8 bits of the result are kept, and stored
			//in Vx.

			//  Vx += Vy

			c.op8nn4(opcodeX, opcodeY)

		case 0x0005:

			// Opcode 8XY5 - Soustraction avec retenue
			// Vx -= Vy
			//Set Vx = Vx - Vy, set VF = NOT borrow. If Vx ¿ Vy, then VF is set to 1, otherwise 0. Then Vy is
			//subtracted from Vx, and the results stored in Vx.

			c.op8nn5(opcodeX, opcodeY)

		case 0x0006:

			// Opcode 8XY6 - Décalage à droite
			//Set Vx = Vx SHR 1. If the least-significant bit of Vx is 1, then VF is set to 1, otherwise 0. Then Vx is
			//divided 	by 2

			c.op8nn6(opcodeX, opcodeY)

		case 0x0007:

			// Opcode 8XY7 - Soustraction inversée avec retenue =
			//Set Vx = Vy - Vx, set VF = NOT borrow. If Vy ¿ Vx, then VF is set to 1, otherwise 0. Then Vx is
			//subtracted from Vy, and the results stored in Vx.

			c.op8nn7(opcodeX, opcodeY)

		case 0xE:

			// Opcode 8XYE - Décalage à gauche =
			//Set Vx = Vx SHL 1. If the most-significant bit of Vx is 1, then VF is set to 1, otherwise to 0. Then Vx is
			//multiplied by 2

			c.op8nnE(opcodeX, opcodeY)

		default:
			// Gérer les opcodes 8XY0 à 8XYE ici (non standard)
		}
	case 0x9000:

		// Opcode 9XY0 - Saut conditionnel (différents registres)=
		//Skip next instruction if Vx != Vy. The values of Vx and Vy are compared, and if they are not equal, the
		//program counter is increased by 2

		c.op9nn0(opcodeX, opcodeY)

	case 0xA000:

		// Opcode ANNN - Chargement de l'index (I) =
		// Set I = nnn. The value of register I is set to nnn.

		c.opAnnn(uint16(opcodeNNN)) // PTET ERREUR  PTET ERREUR = opcodeN a la place

	case 0xB000:

		// Opcode BNNN - Saut avec offset =
		//Jump to location nnn + V0. The program counter is set to nnn plus the value of V0

		c.opBnnn(uint16(opcodeNNN))

	case 0xC000:

		// Opcode CXNN - Génération d'un nombre aléatoire (0 à 255) =
		//Set Vx = random byte AND kk. The interpreter generates a random number from 0 to 255, which is then
		//ANDed with the value kk. The results are stored in Vx. See instruction 8xy2 for more information on AND.

		c.Registre[opcodeN4] = byte(rand.Int()*256) & byte(opcodeNNN)

	case 0xD000:
		// Opcode DXYN - Dessin à l'écran
		c.opDxyn(opcodeX, opcodeY, opcodeN)
		// Gérer l'opcode DXYN ici
		// c.DrawSprite(opcodeN4, opcodeN4, opcodeN4)

	case 0xE000:
		// Gérer les opcodes EX9E et EXA1
		switch opcode & 0x000F {
		case 0x000E:
			// Opcode EX9E - Saut si touche pressée
		case 0x0001:
			// Opcode EXA1 - Saut si touche non pressée
		default:
			// Gérer les opcodes EX9E et EXA1 ici (non standard)
		}
	case 0xF000:
		switch opcode & 0x000F {
		case 0x0007:
			// Opcode FX07 - Chargement du retard
			c.Delay_timer = c.Registre[opcodeX]
		case 0x000A:
			// Opcode FX0A - Attente de touche
		case 0x0005:
			switch opcode & 0x00F0 {
			case 0x0010:
				// Opcode FX15 - Réglage du retard
			case 0x0050:
				// Opcode FX55 - Sauvegarde des registres
			case 0x0060:
				// Opcode FX65 - Chargement des registres
			}

		case 0x0008:
			// Opcode FX18 - Réglage du son
		case 0x000E:
			// Opcode FX1E - Ajout de l'index (I)
		case 0x0009:
			// Opcode FX29 - Chargement de l'emplacement du caractère
		case 0x0003:
			// Opcode FX33 - Chargement des chiffres décimaux

		default:
			// Gérer les opcodes FX07 à FX65 ici (non standard)
		}
	default:
		// Gérer les opcodes non pris en charge ou inconnus ici
	}
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
func (c *Cpu) op00EE() {
	c.Pc = c.Stack[c.Sp]
	c.Sp--
}

func (c *Cpu) op1nnn(address uint16) {
	c.Pc = address - 2
}
func (c *Cpu) op2nnn(address uint16) {
	// Vérifiez que le pointeur de pile (SP) est dans la plage valide (0-15).
	if c.Sp >= 15 {
		return
	}
	c.Sp++
	c.Pc = c.Stack[c.Sp]
	c.Pc = address
}

func (c *Cpu) op3nnn(opcodeX, opcodeNN byte) {
	if c.Registre[opcodeX] == opcodeNN {
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

func (c *Cpu) op8nn1(opcodeX, opcodeY byte) {
	c.Registre[opcodeX] |= c.Registre[opcodeY] //??????????

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
