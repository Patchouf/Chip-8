package emulator

import "math/rand"

// décodage d'un opcode et exécute l'instruction correspondante.
func (c *Cpu) decode(opcode uint16) {
	// Diviser l'opcode en parties individuelles pour faciliter le décodage
	// opcodeN := byte(opcode>>12) & 0x000F // 4 premiers bits
	opcodeX := byte(opcode>>8) & 0x000F // Bits 8 à 11
	opcodeY := byte(opcode>>4) & 0x000F // Bits 4 à 7
	opcodeNNN := opcode & 0x0FFF        // Bits 0 à 11
	opcodeNN := byte(opcode & 0x00FF)   // Bits 0 à 7
	opcodeN4 := byte(opcode & 0x000F)   // 4 derniers bits

	switch opcode & 0xF000 {
	case 0x0000:
		switch opcode {
		case 0x00E0:
			c.op00E0()
		case 0x00EE:
			c.op00EE()
		default:
			// Gérer les opcodes 0NNN ici (non standard)
		}
	case 0x1000:
		c.op1nnn(uint16(opcodeNNN)) // PTET ERREUR = opcodeN a la place
	case 0x2000:
		c.op2nnn(uint16(opcodeNNN))
	case 0x3000:
		c.op3nnn(opcodeX, opcodeNN)
	case 0x4000:
		c.op4nnn(opcodeX, opcodeNN)
	case 0x5000:
		c.op5nnn(opcodeX, opcodeY)
	case 0x6000:
		c.op6nnn(opcodeX, opcodeNN)
	case 0x7000:
		c.op7nnn(opcodeX, opcodeNN)
	case 0x8000:
		// Gérer les opcodes 8XY0 à 8XYE
		switch opcode & 0x000F {
		case 0x0000:
			c.op8nn0(opcodeX, opcodeY)
		case 0x0001:
			c.op8nn1(opcodeX, opcodeY)
		case 0x0002:
			c.op8nn2(opcodeX, opcodeY)
		case 0x0003:
			c.op8nn3(opcodeX, opcodeY)
		case 0x0004:
			c.op8nn4(opcodeX, opcodeY)
		case 0x0005:
			c.op8nn5(opcodeX, opcodeY)
		case 0x0006:
			c.op8nn6(opcodeX, opcodeY)
		case 0x0007:
			c.op8nn7(opcodeX, opcodeY)
		case 0xE:
			c.op8nnE(opcodeX, opcodeY)
		default:
			// Gérer les opcodes 8XY0 à 8XYE ici (non standard)
		}
	case 0x9000:
		c.op9nn0(opcodeX, opcodeY)
	case 0xA000:
		c.opAnnn(uint16(opcodeNNN))
	case 0xB000:
		c.opBnnn(opcodeNNN)
	case 0xC000:
		// Opcode CXNN - Génération d'un nombre aléatoire (0 à 255) =
		//Set Vx = random byte AND kk. The interpreter generates a random number from 0 to 255, which is then
		//ANDed with the value kk. The results are stored in Vx. See instruction 8xy2 for more information on AND.
		c.Registre[opcodeN4] = byte(rand.Int()*256) & byte(opcodeNNN)
	case 0xD000:
		c.opDxyn(opcodeX, opcodeY, opcodeN4)
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
			c.opFx07(opcodeX)
		case 0x000A:
			// Opcode FX0A - Attente de touche
		case 0x0005:
			switch opcode & 0x00F0 {
			case 0x0010:
				c.opFx15(opcodeX)
			case 0x0050:
				c.opFx55(opcodeX)
			case 0x0060:
				c.opFx65(opcodeX)
			default:
				// Gérer les opcodes FX07 à FX65 ici (non standard)
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
