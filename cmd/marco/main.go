package main

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	screenWidth  = 320
	screenHeight = 240

	frameNum = 4

	bottomWidth  = 24
	bottomHeight = 16
	bottom0X     = 155
	bottom0Y     = 1536

	top0X     = 12
	top0Y     = 8
	topWidth  = 33
	topHeight = 29

	frameWidth  = topWidth + bottomWidth
	frameHeight = topHeight + bottomHeight
)

var (
	runnerImage *ebiten.Image
)

type Game struct {
	count int
}

func (g *Game) Update(screen *ebiten.Image) error {
	g.count++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	topOp := &ebiten.DrawImageOptions{}
	topOp.GeoM.Translate(-float64(topWidth)/2, -float64(topHeight)/2)
	topOp.GeoM.Translate(screenWidth/2, screenHeight/2)

	bottomOp := &ebiten.DrawImageOptions{}
	bottomOp.GeoM.Translate(-float64(topWidth)/2, -float64(topHeight)/2)
	bottomOp.GeoM.Translate(0, topHeight-float64(bottomHeight)/2)
	bottomOp.GeoM.Translate(screenWidth/2, screenHeight/2)

	// /5 to slow down the animation
	i := (g.count / 8) % frameNum

	bottomX, bottomY := bottom0X+i*bottomWidth, bottom0Y
	screen.DrawImage(runnerImage.SubImage(image.Rect(bottomX, bottomY, bottomX+bottomWidth, bottomY+bottomHeight)).(*ebiten.Image), bottomOp)

	topX, topY := top0X+i*topWidth, top0Y
	screen.DrawImage(runnerImage.SubImage(image.Rect(topX, topY, topX+topWidth, topY+topHeight)).(*ebiten.Image), topOp)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	var err error
	runnerImage, _, err = ebitenutil.NewImageFromFile("./assets/11226.gif", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Marco character standing")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
