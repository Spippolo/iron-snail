package game

import (
	"fmt"
	"image/color"
	"log"

	"github.com/Spippolo/iron-snail/characters"
	"github.com/Spippolo/iron-snail/common"
	"github.com/Spippolo/iron-snail/sprites"
	"github.com/Spippolo/iron-snail/utils"
	"github.com/Spippolo/iron-snail/weapons"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
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
	if action != g.character.CurrentAction() {
		// An action can fail, for example if the past action isn't over yet
		if g.character.MakeAction(action) {
			if action == characters.Shoot {
				b := weapons.NewBullet(g.character.CurrentWeapon(), g.character.GetDirection())
				g.bullets = append(g.bullets, b)
			}
		}
	}
	g.character.Update()
	for _, b := range g.bullets {
		b.Update()
		// TODO: Delete bullet if out of the screen
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
	_, h := i.Size()
	p := g.character.CurrentPosition()
	op.GeoM.Translate(p.X-float64(center[0]), screenHeight-p.Y-float64(h))
	if err := screen.DrawImage(i, op); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
