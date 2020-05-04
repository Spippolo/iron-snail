package characters

import (
	"image"
	"image/color"
	"math"

	"github.com/Spippolo/iron-snail/common"
	"github.com/Spippolo/iron-snail/sprites"
	"github.com/Spippolo/iron-snail/utils"
	"github.com/Spippolo/iron-snail/weapons"
	"github.com/hajimehoshi/ebiten"
	log "github.com/sirupsen/logrus"
)

type Character interface {
	Update()
	Draw() (*ebiten.Image, [2]int)
	MakeAction(a Action) error
	CurrentWeapon() weapons.Weapon
	SetWeapon(weapons.Weapon)
	CurrentAction() Action
	GetDirection() common.Direction
	SetDirection(d common.Direction) error
	CurrentDirection() common.Direction
}

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
	currentFrame int // number of frame in the current animation
	sprite       *sprites.Sprite
	action       Action
	direction    common.Direction
	weapon       weapons.Weapon
}

func NewMarco() *Marco {
	return &Marco{
		sprite: sprites.Marco(),
		weapon: weapons.Gun,
	}
}

func (c *Marco) CurrentWeapon() weapons.Weapon {
	return c.weapon
}

func (c *Marco) SetWeapon(w weapons.Weapon) {
	c.weapon = w
}

func (c *Marco) CurrentDirection() common.Direction {
	return c.direction
}

func (c *Marco) SetDirection(d common.Direction) error {
	// TODO: optional: validate direction
	c.direction = d
	return nil
}

func (c *Marco) GetDirection() common.Direction {
	return c.direction
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

func (c *Marco) Update() {
	c.currentFrame++
}

func (c *Marco) Draw() (*ebiten.Image, [2]int) {
	legsImage, legsOptions, legsJoint := c.drawLegs()
	legsW, legsH := legsImage.Size()

	bodyImage, bodyOptions, bodyJoint := c.drawBody()
	bodyW, _ := bodyImage.Size()

	// body and legs can move right/left during an animation, creating difference of size
	// so we take into account that difference to calculate to total width of the image and its
	// horizontal translation
	jointDiff := math.Abs(float64(bodyJoint[0] - legsJoint[0]))

	legsOptions.GeoM.Translate(jointDiff, float64(bodyJoint[1]-legsJoint[1]))

	frameWidth := utils.Max(legsW+int(jointDiff), bodyW+int(jointDiff))
	frameHeight := bodyJoint[1] + legsH - legsJoint[1]
	marco, err := ebiten.NewImage(frameWidth, frameHeight, ebiten.FilterNearest)
	utils.CheckErr(err, "Cannot create image for Marco")
	marco.Fill(color.White)

	if err := marco.DrawImage(legsImage, legsOptions); err != nil {
		log.Fatal(err)
	}

	if err := marco.DrawImage(bodyImage, bodyOptions); err != nil {
		log.Fatal(err)
	}

	return marco, bodyJoint
}

func (c *Marco) drawBody() (*ebiten.Image, *ebiten.DrawImageOptions, [2]int) {
	a := c.CurrentAction()
	var part sprites.SpriteName
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

func (c *Marco) drawLegs() (*ebiten.Image, *ebiten.DrawImageOptions, [2]int) {
	return c.drawPart(sprites.LegsStandingPart)
}

func (c *Marco) drawPart(part sprites.SpriteName) (*ebiten.Image, *ebiten.DrawImageOptions, [2]int) {
	s := c.sprite.Desc[part]
	// Number of frames for this part
	frame := (c.currentFrame / s.Speed) % s.Frames
	t := s.Tiles[frame]
	return c.sprite.Image.SubImage(image.Rect(t.X0, t.Y0, t.X0+t.W, t.Y0+t.H)).(*ebiten.Image), &ebiten.DrawImageOptions{}, t.Joint
}
