package opcodes

import (
	"errors"
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

func InitCpu(cpu *Cpu, rombytes []byte) {
	cpu.loadROM(rombytes)
	cpu.Pc = 0x200 - 2
}

func (cpu *Cpu) Update() {
	cpu.Pc += 2
	op1 := cpu.Memory[cpu.Pc]
	op2 := cpu.Memory[cpu.Pc]
	opcode := cpu.uint8ToUint16(op1, op2)
	cpu.decode(opcode)
}

func (cpu *Cpu) loadROM(rombytes []byte) {
	cpu.Romlength = uint16(len(rombytes))
	for i, byt := range rombytes {
		cpu.Memory[0x200+i] = byt
	}
}

// Fonction stackPop
func (c *Cpu) stackPop() (uint16, error) {
	if c.Sp == 0 {
		return 0, errors.New("pile vide")
	}
	c.Sp--
	address := c.Stack[c.Sp]

	return address, nil
}

// Fonction stackPush
func (c *Cpu) StackPush(address uint16) {
	// Vérifiez que le pointeur de pile (SP) est dans la plage valide (0-15).
	if c.Sp >= 15 {
		return
	}
	c.Stack[c.Sp] = address
	c.Sp++
}

// Fonction uint16 to uint8
func (c *Cpu) Uint16ToUint8(n uint16) (uint8, uint8) {
	return uint8(n >> 8), uint8(n & 0x00FF)
}

// unit 8 to uint16
func (c *Cpu) uint8ToUint16(n1 uint8, n2 uint8) uint16 {
	return uint16(n1)<<8 | uint16(n2)
}

// Fonction uint8 to uint4
func (c *Cpu) Uint8ToUint4(n uint8) (uint8, uint8) {
	return uint8(n >> 4), uint8(n & 0x0F)
}

func (c *Cpu) DrawSprite(X, Y, height byte) bool {
	ScreenWidth := uint16(c.Registre[X])
	ScreenHeight := uint16(c.Registre[Y])

	c.Registre[0xF] = 0

	// Parcourez les lignes du sprite.
	for row := byte(0); row < height; row++ {
		spriteByte := c.Memory[c.I+uint16(row)]

		for bit := byte(0); bit < 8; bit++ {

			if (spriteByte & (0x80 >> bit)) != 0 {

				x := int(ScreenWidth) + int(bit)
				y := int(ScreenHeight) + int(row)

				if x < 64 && y < 32 {

					index := y*64 + x

					if c.Gfx[index][0] == byte(1) {

						c.Registre[0xF] = 1
					}
					for i := 0; i < len(c.Gfx[index]); i++ {
						c.Gfx[index][i] ^= byte(1)
					}
				}
			}
		}
	}

	return c.Registre[0xF] == 1
}

// func (c *Cpu) DrawSprite(x byte, y byte, row byte) bool {
// 	erased := false
// 	yIndex := y % 64

// 	for i := x; i < x+8; i++ {
// 		xIndex := i % 32

// 		wasSet := c.Gfx[xIndex][yIndex] == 1
// 		value := row >> (x + 8 - i - 1) & 0x01

// 		c.Gfx[xIndex][yIndex] ^= value

// 		if wasSet && c.Gfx[xIndex][yIndex] == 0 {
// 			erased = true
// 		}
// 	}

// 	return erased
// }

// Decode décode un opcode et exécute l'instruction correspondante.
func (c *Cpu) decode(opcode uint16) {
	// Diviser l'opcode en parties individuelles PROBLEME
	opcodeN := byte(opcode>>12) & 0x000F // 4 premiers bits
	opcodeX := byte(opcode>>8) & 0x000F  // Bits 8 à 11
	opcodeY := byte(opcode>>4) & 0x000F  // Bits 4 à 7
	opcodeNNN := opcode & 0x0FFF         // Bits 0 à 11
	// opcodeKK := uint8(opcode & 0x00FF) // Bits 0 à 7
	opcodeN4 := uint8(opcode & 0x000F) // 4 derniers bits
	//x := byte(opcodeX & )

	// Utilisez un switch pour gérer chaque opcode
	switch opcodeN {
	case 0x0:
		switch opcode {
		case 0x00E0:
			// Opcode 00E0 - Effacer l'écran
			c.op00E0()
		case 0x00EE:
			// Opcode 00EE - Retour de sous-routine
			c.stackPop()
		default:
			// Gérer les opcodes 0NNN ici (non standard)
		}
	case 0x1:
		// Opcode 1NNN - Saut
		c.op1nnn(uint16(opcodeNNN)) // PTET ERREUR = opcodeN a la place
	case 0x2:
		// Opcode 2NNN - Appel de sous-routine
		c.StackPush(c.Pc)
	case 0x3:
		// Opcode 3XNN - Saut conditionnel (égal)
		if c.Registre[opcodeN4] == byte(opcodeNNN) {
			c.Pc += 2
		}
	case 0x4:
		// Opcode 4XNN - Saut conditionnel (différent)
		if c.Registre[opcodeN4] != byte(opcodeNNN) {
			c.Pc += 2
		}
	case 0x5:
		// Opcode 5XY0 - Saut conditionnel (égalité de registres)
		if c.Registre[opcodeN4] == c.Registre[opcodeN4] {
			c.Pc += 2
		}
	case 0x6:
		// Opcode 6XNN - Chargement de valeur constante
		c.op6XNN(opcodeX, opcodeN)
		// c.Registre[opcodeN4] = byte(opcodeNNN)
	case 0x7:
		// Opcode 7XNN - Ajout de valeur constante
		c.Registre[opcodeN4] += byte(opcodeNNN)
	case 0x8:
		// Gérer les opcodes 8XY0 à 8XYE
		switch opcodeN4 {
		case 0x0:
			// Opcode 8XY0 - Copie de Registre
			c.Registre[opcodeN4] = c.Registre[opcodeN4]
		case 0x1:
			// Opcode 8XY1 - Opération OU (bitwise OR)
			c.Registre[opcodeN4] |= byte(opcodeNNN)
		case 0x2:
			// Opcode 8XY2 - Opération ET (bitwise AND)
			c.Registre[opcodeN4] &= byte(opcodeNNN)
		case 0x3:
			// Opcode 8XY3 - Opération XOR (bitwise XOR)
			c.Registre[opcodeN4] ^= byte(opcodeNNN)
		case 0x4:
			// Opcode 8XY4 - Ajout avec retenue
		case 0x5:
			// Opcode 8XY5 - Soustraction avec retenue
		case 0x6:
			// Opcode 8XY6 - Décalage à droite
			c.Registre[opcodeN4] >>= 1
		case 0x7:
			// Opcode 8XY7 - Soustraction inversée avec retenue
		case 0xE:
			// Opcode 8XYE - Décalage à gauche
			c.Registre[opcodeN4] <<= 1
		default:
			// Gérer les opcodes 8XY0 à 8XYE ici (non standard)
		}
	case 0x9:
		// Opcode 9XY0 - Saut conditionnel (différents registres)
		if c.Registre[opcodeN4] != c.Registre[opcodeN4] {
			c.Pc += 2
		}
	case 0xA:
		c.opAnnn(uint16(opcodeNNN)) // PTET ERREUR  PTET ERREUR = opcodeN a la place
		// Opcode ANNN - Chargement de l'index (I)
		// c.I = opcodeNNN
	case 0xB:
		// Opcode BNNN - Saut avec offset
		// c.Pc = opcodeNNN + uint16(c.Registre[0])
	case 0xC:
		// Opcode CXNN - Génération d'un nombre aléatoire (0 à 255)
		c.Registre[opcodeN4] = byte(rand.Int()*256) & byte(opcodeNNN)
	case 0xD:
		// Opcode DXYN - Dessin à l'écran
		c.opDxyn(opcodeX, opcodeY, opcodeN)
		// Gérer l'opcode DXYN ici
		c.DrawSprite(opcodeN4, opcodeN4, opcodeN4)
		break

	case 0xE:
		// Gérer les opcodes EX9E et EXA1
		switch opcodeNNN {
		case 0x9E:
			// Opcode EX9E - Saut si touche pressée
		case 0xA1:
			// Opcode EXA1 - Saut si touche non pressée
		default:
			// Gérer les opcodes EX9E et EXA1 ici (non standard)
		}
	case 0xF:
		switch opcodeNNN {
		case 0x07:
			// Opcode FX07 - Chargement du retard
		case 0x0A:
			// Opcode FX0A - Attente de touche
		case 0x15:
			// Opcode FX15 - Réglage du retard
		case 0x18:
			// Opcode FX18 - Réglage du son
		case 0x1E:
			// Opcode FX1E - Ajout de l'index (I)
		case 0x29:
			// Opcode FX29 - Chargement de l'emplacement du caractère
		case 0x33:
			// Opcode FX33 - Chargement des chiffres décimaux
		case 0x55:
			// Opcode FX55 - Sauvegarde des registres
		case 0x65:
			// Opcode FX65 - Chargement des registres
		default:
			// Gérer les opcodes FX07 à FX65 ici (non standard)
		}
	default:
		// Gérer les opcodes non pris en charge ou inconnus ici
	}
}
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
func (c *Cpu) op6XNN(opcodeX, opcodeNNN byte) {
	c.Registre[opcodeX] = opcodeNNN
}
func (c *Cpu) op1nnn(address uint16) {
	c.Pc = address
}

func (c *Cpu) opAnnn(address uint16) {
	c.I = address
}
