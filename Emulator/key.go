package emulator

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type Clavier struct {
	IsPressed [16]bool
}

func RefreshKeys(chip8 *Clavier, cpu *Cpu) {
	chip8.IsPressed = [16]bool{
		ebiten.IsKeyPressed(ebiten.Key1), //1
		ebiten.IsKeyPressed(ebiten.Key2), //2
		ebiten.IsKeyPressed(ebiten.Key3), //3
		ebiten.IsKeyPressed(ebiten.Key4), //4
		ebiten.IsKeyPressed(ebiten.KeyQ), //A
		ebiten.IsKeyPressed(ebiten.KeyW), //Z
		ebiten.IsKeyPressed(ebiten.KeyE), //E
		ebiten.IsKeyPressed(ebiten.KeyR), //R
		ebiten.IsKeyPressed(ebiten.KeyA), //Q
		ebiten.IsKeyPressed(ebiten.KeyS), //S
		ebiten.IsKeyPressed(ebiten.KeyD), //D
		ebiten.IsKeyPressed(ebiten.KeyF), //F
		ebiten.IsKeyPressed(ebiten.KeyZ), //W
		ebiten.IsKeyPressed(ebiten.KeyX), //X
		ebiten.IsKeyPressed(ebiten.KeyC), //C
		ebiten.IsKeyPressed(ebiten.KeyV), //V
	}

	// set up tab key
	for i := 0 ; i < len(chip8.IsPressed); i++ {
		if chip8.IsPressed[i] {
			cpu.Key[i] = 1
		} else {
			cpu.Key[i] = 0
		}
	}
	fmt.Printf("%x \n", cpu.Key)

	// test fonctionnement
	// for index := 0; index < len(chip8.IsPressed); index++ {
	// 	if chip8.IsPressed[index] {
	// 		fmt.Println(index)
	// 		cpu.Key[index] = 255
	// 		for ix := range cpu.Key {
	// 			if cpu.Key[ix] == 255 {
	// 				fmt.Println("oura", +ix)
	// 			}
	// 		}
	// 	}
	// }
}




func (clavier *Clavier) GetKey(key byte) bool {
	fmt.Println(clavier.IsPressed[key])
	return clavier.IsPressed[key]
}
func (clavier *Clavier) Update(keyIndex byte, pressed bool) {
	clavier.IsPressed[keyIndex] = pressed
}
