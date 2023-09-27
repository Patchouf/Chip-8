package opcodes

import (
	"errors"
	"math/rand"
)

type Cpu struct {
	memory      [4096]byte
	registre    [16]byte
	i           uint16
	pc          uint16
	Gfx         [64 * 32]byte
	delay_timer byte
	stack       [16]uint16
	sp          byte
	key         [16]byte
}

type Object struct {
	clavier [16]byte
}

// Fonction stackPop
func (c *Cpu) stackPop() (uint16, error) {
	if c.sp == 0 {
		return 0, errors.New("pile vide")
	}
	c.sp--
	address := c.stack[c.sp]

	return address, nil
}

// Fonction stackPush
func (c *Cpu) stackPush(address uint16) {
	// Vérifiez que le pointeur de pile (SP) est dans la plage valide (0-15).
	if c.sp >= 15 {
		return
	}
	c.stack[c.sp] = address
	c.sp++
}

// Fonction uint16 to uint8
func (c *Cpu) uint16ToUint8(n uint16) (uint8, uint8) {
	return uint8(n >> 8), uint8(n & 0x00FF)
}

// Fonction uint8 to uint4
func (c *Cpu) uint8ToUint4(n uint8) (uint8, uint8) {
	return uint8(n >> 4), uint8(n & 0x0F)
}

// Decode décode un opcode et exécute l'instruction correspondante.
func (c *Cpu) decode(opcode uint16) {
	// Diviser l'opcode en parties individuelles
	opcodeN := (opcode >> 12) & 0x000F // 4 premiers bits
	// opcodeX := (opcode >> 8) & 0x000F   // Bits 8 à 11
	// opcodeY := (opcode >> 4) & 0x000F   // Bits 4 à 7
	opcodeNNN := opcode & 0x0FFF // Bits 0 à 11
	// opcodeKK := uint8(opcode & 0x00FF) // Bits 0 à 7
	opcodeN4 := uint8(opcode & 0x000F) // 4 derniers bits

	// Utilisez un switch pour gérer chaque opcode
	switch opcodeN {
	case 0x0:
		switch opcode {
		case 0x00E0:
			// Opcode 00E0 - Effacer l'écran
			for i := range c.Gfx {
				c.Gfx[i] = 0
			}
		case 0x00EE:
			// Opcode 00EE - Retour de sous-routine
			c.stackPop()
		default:
			// Gérer les opcodes 0NNN ici (non standard)
		}
	case 0x1:
		// Opcode 1NNN - Saut
		c.pc = opcodeNNN
	case 0x2:
		// Opcode 2NNN - Appel de sous-routine
		c.stackPush(c.pc)
	case 0x3:
		// Opcode 3XNN - Saut conditionnel (égal)
		if c.registre[opcodeN4] == byte(opcodeNNN) {
			c.pc += 2
		}
	case 0x4:
		// Opcode 4XNN - Saut conditionnel (différent)
		if c.registre[opcodeN4] != byte(opcodeNNN) {
			c.pc += 2
		}
	case 0x5:
		// Opcode 5XY0 - Saut conditionnel (égalité de registres)
		if c.registre[opcodeN4] == c.registre[opcodeN4] {
			c.pc += 2
		}
	case 0x6:
		// Opcode 6XNN - Chargement de valeur constante
		c.registre[opcodeN4] = byte(opcodeNNN)
	case 0x7:
		// Opcode 7XNN - Ajout de valeur constante
		c.registre[opcodeN4] += byte(opcodeNNN)
	case 0x8:
		// Gérer les opcodes 8XY0 à 8XYE
		switch opcodeN4 {
		case 0x0:
			// Opcode 8XY0 - Copie de registre
			c.registre[opcodeN4] = c.registre[opcodeN4]
		case 0x1:
			// Opcode 8XY1 - Opération OU (bitwise OR)
			c.registre[opcodeN4] |= byte(opcodeNNN)
		case 0x2:
			// Opcode 8XY2 - Opération ET (bitwise AND)
			c.registre[opcodeN4] &= byte(opcodeNNN)
		case 0x3:
			// Opcode 8XY3 - Opération XOR (bitwise XOR)
			c.registre[opcodeN4] ^= byte(opcodeNNN)
		case 0x4:
			// Opcode 8XY4 - Ajout avec retenue
		case 0x5:
			// Opcode 8XY5 - Soustraction avec retenue
		case 0x6:
			// Opcode 8XY6 - Décalage à droite
			c.registre[opcodeN4] >>= 1
		case 0x7:
			// Opcode 8XY7 - Soustraction inversée avec retenue
		case 0xE:
			// Opcode 8XYE - Décalage à gauche
			c.registre[opcodeN4] <<= 1
		default:
			// Gérer les opcodes 8XY0 à 8XYE ici (non standard)
		}
	case 0x9:
		// Opcode 9XY0 - Saut conditionnel (différents registres)
		if c.registre[opcodeN4] != c.registre[opcodeN4] {
			c.pc += 2
		}
	case 0xA:
		// Opcode ANNN - Chargement de l'index (I)
		c.i = opcodeNNN
	case 0xB:
		// Opcode BNNN - Saut avec offset
		c.pc = opcodeNNN + uint16(c.registre[0])
	case 0xC:
		// Opcode CXNN - Génération d'un nombre aléatoire (0 à 255)
		c.registre[opcodeN4] = byte(rand.Int()*256) & byte(opcodeNNN)
	case 0xD:
		// Opcode DXYN - Dessin à l'écran

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
