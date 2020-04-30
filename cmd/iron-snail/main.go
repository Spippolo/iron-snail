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
	count      int
	characters []characters.Character
}

func (g *Game) Update(screen *ebiten.Image) error {
	g.count++
	for _, c := range g.characters {
		c.Update(g.count)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(screenWidth/2, screenHeight/2)
	// In case we're drawing multiple characters, draw them one near the other (and not over)
	var xOffset int
	for _, c := range g.characters {
		op.GeoM.Translate(float64(xOffset), 0)
		i := c.Draw()
		if err := screen.DrawImage(i, op); err != nil {
			log.Fatal(err)
		}
		xOffset, _ = i.Size()
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Iron Snail")
	if err := ebiten.RunGame(&Game{
		characters: []characters.Character{
			characters.NewMarco(),
		},
	}); err != nil {
		log.Fatal(err)
	}
}
