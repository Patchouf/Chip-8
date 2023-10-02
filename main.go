package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"main.go/opcodes"
)

type Game struct {
	cpu opcodes.Cpu
}

// update du jeu
func (g *Game) Update() error {
	g.cpu.Update()
	time.Sleep(time.Millisecond)
	return nil
}

// dessin des pixels du jeu
func (g *Game) Draw(screen *ebiten.Image) {
	// ebitenutil.DebugPrint(screen,"Hello")
	for x, row := range g.cpu.Gfx {
		for y, pixel := range row {
			if pixel == 1 {
				screen.Set(x, y, color.RGBA{R: 255, G: 205, B: 1, A: 255})
			} else {
				screen.Set(x, y, color.RGBA{R: 153, G: 102, B: 1, A: 255})
			}
		}
	}
}

// fonction pour set l'écran
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 64, 32
}

// fonction pour start la game, ouverture du screen
func main() {
	filename := os.Args[1]
	rombytes := readROM(filename)
	// fmt.Println(rombytes)
	// PrintROM(rombytes)

	var game Game
	opcodes.InitCpu(&game.cpu, rombytes)
	fmt.Println(game.cpu.Memory)

	ebiten.SetWindowSize(640, 320)
	ebiten.SetWindowTitle("Chip8 Emulator :  " + filename[6:])
	ebiten.RunGame(&game)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}

}

// fonction pour lire le fichier rom
func readROM(filename string) []byte {
	dat, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return dat
}

// fonction pour print le rom
func PrintROM(rom []byte) {
	for i, byt := range rom {
		if i%2 == 0 {
			fmt.Printf("0x%03x: ", 0x200+i)
		}
		fmt.Printf("%02x", byt)
		if i%2 == 1 {
			fmt.Print("\n")
		}
	}
}
