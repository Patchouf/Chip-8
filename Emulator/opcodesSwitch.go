package emulator

import "math/rand"

// décodage d'un opcode et exécute l'instruction correspondante.
func (c *Cpu) decode(opcode uint16) {
	// Diviser l'opcode en parties individuelles PROBLEME
	// opcodeN := byte(opcode>>12) & 0x000F // 4 premiers bits
	opcodeX := byte(opcode>>8) & 0x000F // Bits 8 à 11
	opcodeY := byte(opcode>>4) & 0x000F // Bits 4 à 7
	opcodeNNN := opcode & 0x0FFF        // Bits 0 à 11
	opcodeNN := byte(opcode & 0x00FF)   // Bits 0 à 7
	opcodeN4 := byte(opcode & 0x000F)   // 4 derniers bits
	//x := byte(opcodeX & )

	// Utilisez un switch pour gérer chaque opcode
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

		// Opcode 1NNN - Saut
		// Jump to location nnn. The interpreter sets the program counter to nnn.

		c.op1nnn(uint16(opcodeNNN)) // PTET ERREUR = opcodeN a la place

	case 0x2000:

		// Opcode 2NNN - Appel de sous-routine =
		// Call subroutine at nnn. The interpreter increments the stack pointer, then puts the current PC on the top
		//of the stack. The PC is then set to nnn.
		c.op2nnn(uint16(opcodeNNN))

	case 0x3000:

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

			c.op8nn1(opcodeX, opcodeY)

		case 0x0002:

			// Opcode 8XY2 - Opération ET (bitwise AND) =
			//Set Vx = Vx AND Vy. Performs a bitwise AND on the values of Vx and Vy, then stores the result in Vx.
			//A bitwise AND compares the corresponding bits from two values, and if both bits are 1, then the same bit
			//in the result is also 1. Otherwise, it is 0.

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

		// Opcode ANNN - Chargement de l'index (I) =
		// Set I = nnn. The value of register I is set to nnn.

		c.opAnnn(uint16(opcodeNNN)) // PTET ERREUR  PTET ERREUR = opcodeN a la place

	case 0xB000:

		// Opcode BNNN - Saut avec offset =
		//Jump to location nnn + V0. The program counter is set to nnn plus the value of V0

		c.opBnnn(opcodeNNN) // PTET ERREUR  PTET ERREUR = opcodeN a la place

	case 0xC000:

		// Opcode CXNN - Génération d'un nombre aléatoire (0 à 255) =
		//Set Vx = random byte AND kk. The interpreter generates a random number from 0 to 255, which is then
		//ANDed with the value kk. The results are stored in Vx. See instruction 8xy2 for more information on AND.

		c.Registre[opcodeN4] = byte(rand.Int()*256) & byte(opcodeNNN)

	case 0xD000:
		// Opcode DXYN - Dessin à l'écran
		c.opDxyn(opcodeX, opcodeY, opcodeN4)
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
