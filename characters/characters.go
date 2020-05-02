package characters

import (
	"image"
	"image/color"

	"github.com/Spippolo/iron-snail/sprites"
	"github.com/Spippolo/iron-snail/utils"
	"github.com/hajimehoshi/ebiten"
	log "github.com/sirupsen/logrus"
)

type Character interface {
	Update(int) error
	Draw() *ebiten.Image
	MakeAction(a Action) error
	CurrentAction() Action
	SetDirection(d Direction) error
	CurrentDirection() Direction
}

type Direction int

const (
	Right Direction = iota
	Left
)

type Action int

const (
	Stand Action = iota
	Walk
	Shoot
	Knife
	KnifeUp
	Jump
	JumpShoot
	Crouch
	CrouchWalk
)

type Marco struct {
	tick         int
	currentFrame int // number of frame in the current animation
	sprite       *sprites.Sprite
	action       Action
	direction    Direction
}

func NewMarco() *Marco {
	return &Marco{
		sprite: sprites.Marco(),
	}
}

func (c *Marco) CurrentDirection() Direction {
	return c.direction
}

func (c *Marco) SetDirection(d Direction) error {
	// TODO: optional: validate direction
	c.direction = d
	return nil
}

func (c *Marco) CurrentAction() Action {
	return c.action
}

func (c *Marco) MakeAction(action Action) error {
	log.Debugf("Action %d", action)
	// TODO: optional: validate action
	c.action = action
	// Reset the animation
	// TODO: some actions can't be reset, like shooting, and must be completed before resetting
	c.currentFrame = 0
	return nil
}

func (c *Marco) Update(tick int) error {
	c.tick = tick
	c.currentFrame++
	return nil
}

func (c *Marco) Draw() *ebiten.Image {
	legsImage, legsOptions, _, _ := c.drawLegs()
	legsW, legsH := legsImage.Size()

	bodyImage, bodyOptions, _, YOffset := c.drawBody()
	bodyW, bodyH := bodyImage.Size()
	frameWidth := utils.Max(legsW, bodyW)
	frameHeight := bodyH - YOffset + legsH

	bodyOptions.GeoM.Translate(0, float64(-bodyH+YOffset))

	legsOptions.GeoM.Translate(0, float64(frameHeight-legsH))
	bodyOptions.GeoM.Translate(0, float64(frameHeight-legsH))

	marco, err := ebiten.NewImage(frameWidth, frameHeight, ebiten.FilterNearest)
	if err != nil {
		log.Fatal(err)
	}
	marco.Fill(color.White)

	if err := marco.DrawImage(legsImage, legsOptions); err != nil {
		log.Fatal(err)
	}

	if err := marco.DrawImage(bodyImage, bodyOptions); err != nil {
		log.Fatal(err)
	}

	return marco
}

func (c *Marco) drawBody() (*ebiten.Image, *ebiten.DrawImageOptions, int, int) {
	a := c.CurrentAction()
	var part sprites.BodyPart
	if a == Stand {
		part = sprites.BodyStandingPart
	} else if a == Shoot {
		part = sprites.BodyShootingPart
	} else if a == Knife {
		part = sprites.BodyKnifePart
	} else if a == KnifeUp {
		part = sprites.BodyKnifeUpPart
	}
	return c.drawPart(part)
}

func (c *Marco) drawLegs() (*ebiten.Image, *ebiten.DrawImageOptions, int, int) {
	return c.drawPart(sprites.LegsStandingPart)
}

func (c *Marco) drawPart(part sprites.BodyPart) (*ebiten.Image, *ebiten.DrawImageOptions, int, int) {
	s := c.sprite.Desc[part]
	// Number of frames for this part
	frameNum := len(s.Tiles)
	frame := (c.currentFrame / (*c.sprite.Desc[part]).Speed) % frameNum
	t := s.Tiles[frame]
	options := &ebiten.DrawImageOptions{}
	return c.sprite.Image.SubImage(image.Rect(t.X0, t.Y0, t.X0+t.W, t.Y0+t.H)).(*ebiten.Image), options, t.XOffset, t.YOffset
}
