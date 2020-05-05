package game

import (
	"github.com/Spippolo/iron-snail/characters"
	"github.com/Spippolo/iron-snail/sprites"
	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

func Run() error {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Iron Snail")

	return ebiten.RunGame(&Game{
		character: characters.NewMarco(),
		weapons:   sprites.Weapons(),
	})
}
