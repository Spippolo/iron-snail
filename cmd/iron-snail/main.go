package main

import (
	"fmt"
	"image/color"

	"github.com/Spippolo/iron-snail/characters"
	"github.com/Spippolo/iron-snail/common"
	"github.com/Spippolo/iron-snail/sprites"
	"github.com/Spippolo/iron-snail/utils"
	"github.com/Spippolo/iron-snail/weapons"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	log "github.com/sirupsen/logrus"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

type Game struct {
	count     int
	character characters.Character
	bullets   []weapons.Bullet
	weapons   *sprites.Sprite
}

func (g *Game) Update(screen *ebiten.Image) error {
	g.count++
	// TODO: collisions
	var action characters.Action
	var direction common.Direction
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		action = characters.Jump
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
		action = characters.Crouch
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		direction = common.East
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		direction = common.West
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		action = characters.Shoot
	} else if ebiten.IsKeyPressed(ebiten.KeyE) {
		action = characters.Knife
	} else if ebiten.IsKeyPressed(ebiten.KeyW) {
		action = characters.KnifeUp
	}
	g.character.Update()
	for _, b := range g.bullets {
		b.Update()
		// TODO: Delete bullet if out of the screen
	}
	if action != g.character.CurrentAction() {
		err := g.character.MakeAction(action)
		if action == characters.Shoot {
			b := weapons.NewBullet(g.character.CurrentWeapon(), g.character.GetDirection())
			g.bullets = append(g.bullets, b)
		}
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
	g.drawCharacter(screen)
	g.drawBullets(screen)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %f", ebiten.CurrentTPS()))
}

func (g *Game) drawBullets(screen *ebiten.Image) {
	for _, b := range g.bullets {
		i, op := b.Draw()
		op.GeoM.Translate(screenWidth/2, screenHeight/2)
		if err := screen.DrawImage(i, op); err != nil {
			log.Fatal(err)
		}
	}
}

func (g *Game) drawCharacter(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	i, center := g.character.Draw()
	op.GeoM.Translate(screenWidth/2-float64(center[0]), screenHeight/2-float64(center[1]))
	if err := screen.DrawImage(i, op); err != nil {
		log.Fatal(err)
	}
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
		weapons:   sprites.Weapons(),
	}); err != nil {
		log.Fatal(err)
	}
}
