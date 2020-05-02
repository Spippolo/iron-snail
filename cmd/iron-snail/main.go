package main

import (
	"image/color"

	"github.com/Spippolo/iron-snail/characters"
	"github.com/Spippolo/iron-snail/utils"
	"github.com/hajimehoshi/ebiten"
	log "github.com/sirupsen/logrus"
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
	var action characters.Action
	var direction characters.Direction
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		action = characters.Jump
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
		action = characters.Crouch
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		direction = characters.Right
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		direction = characters.Left
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		action = characters.Shoot
	}
	err := g.character.Update(g.count)
	utils.CheckErr(err, "Update failed")
	if action != g.character.CurrentAction() {
		err := g.character.MakeAction(action)
		utils.CheckErr(err, "Action failed")
	}
	if direction != g.character.CurrentDirection() {
		err := g.character.SetDirection(direction)
		utils.CheckErr(err, "Set direction failed")
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(screenWidth/2, screenHeight/2)
	// In case we're drawing multiple characters, draw them one near the other (and not over)
	var xOffset int
	op.GeoM.Translate(float64(xOffset), 0)
	i := g.character.Draw(screen)
	if err := screen.DrawImage(i, op); err != nil {
		log.Fatal(err)
	}
	xOffset, _ = i.Size()

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Iron Snail")
	if err := ebiten.RunGame(&Game{
		character: characters.NewMarco(),
	}); err != nil {
		log.Fatal(err)
	}
}
