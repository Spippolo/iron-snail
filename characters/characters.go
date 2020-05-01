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
	Update(int)
	Draw(screen *ebiten.Image) *ebiten.Image
}

type Marco struct {
	tick   int
	speed  int // slow down the animation. 1 is fast
	sprite *sprites.Sprite
}

func NewMarco() *Marco {
	return &Marco{
		speed:  15,
		sprite: sprites.Marco(),
	}
}

func (c *Marco) Update(tick int) {
	c.tick = tick
}

func (c *Marco) Draw(screen *ebiten.Image) *ebiten.Image {

	legsImage, legsOptions, _ := c.drawLegs()
	legsW, legsH := legsImage.Size()

	bodyImage, bodyOptions, YOffset := c.drawBody()
	bodyW, bodyH := bodyImage.Size()
	bodyH -= YOffset
	frameHeight := bodyH + legsH
	// Move Marco at the bottom of the image
	legsOptions.GeoM.Translate(0, float64(frameHeight-legsH))
	bodyOptions.GeoM.Translate(0, float64(frameHeight-2*legsH-YOffset))

	marco, err := ebiten.NewImage(utils.Max(legsW, bodyW), frameHeight, ebiten.FilterNearest)
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

func (c *Marco) drawBody() (*ebiten.Image, *ebiten.DrawImageOptions, int) {
	return c.drawPart(sprites.BodyStandingPart)
}

func (c *Marco) drawLegs() (*ebiten.Image, *ebiten.DrawImageOptions, int) {
	return c.drawPart(sprites.LegsStandingPart)
}

func (c *Marco) drawPart(part sprites.BodyPart) (*ebiten.Image, *ebiten.DrawImageOptions, int) {
	s := c.sprite.Desc[part]
	// Number of frames for this part
	frameNum := len(s.Tiles)
	frame := (c.tick / c.speed) % frameNum
	t := s.Tiles[frame]
	options := &ebiten.DrawImageOptions{}
	return c.sprite.Image.SubImage(image.Rect(t.X0, t.Y0, t.X0+t.W, t.Y0+t.H)).(*ebiten.Image), options, t.YOffset
}
