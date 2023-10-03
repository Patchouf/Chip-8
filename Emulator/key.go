package emulator

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func RefreshKeys(chip8 *Clavier) {
	chip8.IsPressed = [16]bool{
		ebiten.IsKeyPressed(ebiten.Key1),
		ebiten.IsKeyPressed(ebiten.Key2),
		ebiten.IsKeyPressed(ebiten.Key3),
		ebiten.IsKeyPressed(ebiten.Key4),
		ebiten.IsKeyPressed(ebiten.KeyA),
		ebiten.IsKeyPressed(ebiten.KeyZ),
		ebiten.IsKeyPressed(ebiten.KeyE),
		ebiten.IsKeyPressed(ebiten.KeyR),
		ebiten.IsKeyPressed(ebiten.KeyQ),
		ebiten.IsKeyPressed(ebiten.KeyS),
		ebiten.IsKeyPressed(ebiten.KeyD),
		ebiten.IsKeyPressed(ebiten.KeyF),
		ebiten.IsKeyPressed(ebiten.KeyW),
		ebiten.IsKeyPressed(ebiten.KeyX),
		ebiten.IsKeyPressed(ebiten.KeyC),
		ebiten.IsKeyPressed(ebiten.KeyV),
	}
}

type Clavier struct {
	IsPressed [16]bool
}

func (clavier *Clavier) GetKey(key byte) bool {
	return clavier.IsPressed[key]
}

func (clavier *Clavier) Update(keyIndex byte, pressed bool){
	clavier.IsPressed[keyIndex] = pressed
}