package main

import (
	"log"

	"github.com/Spippolo/iron-snail/characters"
	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

type Game struct {
	count     int
	character characters.Character
}

func (g *Game) Update(screen *ebiten.Image) error {
	g.count++
	g.character.Update(g.count)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(screenWidth/2, screenHeight/2)

	if err := screen.DrawImage(g.character.Draw(), op); err != nil {
		log.Fatal(err)
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Marco character standing")
	if err := ebiten.RunGame(&Game{
		character: characters.NewMarco(),
	}); err != nil {
		log.Fatal(err)
	}
}
