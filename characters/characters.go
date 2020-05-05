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
	MakeAction(a Action) bool
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
	tick            int
	sprite          *sprites.Sprite
	action          Action
	direction       common.Direction
	weapon          weapons.Weapon
	bodyPart        sprites.SpriteName // Current body part to be drawn
	legsPart        sprites.SpriteName // Current legs part to be drawn
	lastActionFrame bool
}

func NewMarco() *Marco {
	return &Marco{
		sprite:   sprites.Marco(),
		weapon:   weapons.Gun,
		bodyPart: sprites.BodyStandingPart,
		legsPart: sprites.LegsStandingPart,
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

// Can actually perform an action if it's standing. It's not possible to interrupt an ongoing action
func (c *Marco) canPerform(action Action) bool {
	return c.action == Stand
}

func (c *Marco) MakeAction(action Action) bool {
	if !c.canPerform(action) {
		return false
	}
	log.Debugf("Action %d", action)

	if action == Stand {
		c.bodyPart = sprites.BodyStandingPart
	} else if action == Shoot {
		c.bodyPart = sprites.BodyShootingPart
	} else if action == Knife {
		c.bodyPart = sprites.BodyKnifePart
	} else if action == KnifeUp {
		c.bodyPart = sprites.BodyKnifeUpPart
	} else {
		log.Fatalf("Unknown action %v", action)
	}
	c.action = action

	c.legsPart = sprites.LegsStandingPart

	// Reset the animation
	c.tick = 0
	return true
}

func (c *Marco) Update() {
	c.tick++

	// The previous Update calculated it was the last tick for the action, so we need to
	// reset Marco to its rest position
	if c.lastActionFrame {
		c.resetAction()
	}
}

func (c *Marco) resetAction() {
	c.tick = 0
	c.lastActionFrame = false
	c.action = Stand
	c.bodyPart = sprites.BodyStandingPart
	c.legsPart = sprites.LegsStandingPart
}

func (c *Marco) Draw() (*ebiten.Image, [2]int) {
	nextLegsframe := (c.tick / c.sprite.Desc[c.legsPart].Speed) % c.sprite.Desc[c.legsPart].Frames
	l := c.sprite.Desc[c.legsPart].Tiles[nextLegsframe]
	legsImage := c.sprite.Image.SubImage(image.Rect(l.X0, l.Y0, l.X0+l.W, l.Y0+l.H)).(*ebiten.Image)
	legsOptions := &ebiten.DrawImageOptions{}
	legsJoint := l.Joint
	legsW, legsH := legsImage.Size()

	nextBodyframe := (c.tick / c.sprite.Desc[c.bodyPart].Speed) % c.sprite.Desc[c.bodyPart].Frames
	b := c.sprite.Desc[c.bodyPart].Tiles[nextBodyframe]
	bodyImage := c.sprite.Image.SubImage(image.Rect(b.X0, b.Y0, b.X0+b.W, b.Y0+b.H)).(*ebiten.Image)
	bodyOptions := &ebiten.DrawImageOptions{}
	bodyJoint := b.Joint
	bodyW, _ := bodyImage.Size()

	if nextBodyframe == c.sprite.Desc[c.bodyPart].Frames-1 {
		c.lastActionFrame = true
	}

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
