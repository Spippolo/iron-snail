package characters

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	log "github.com/sirupsen/logrus"
)

var (
	sprite *ebiten.Image
)

type Character interface {
	Update(int)
	Draw() *ebiten.Image
}

type Marco struct {
	tick     int
	speed    int // slow down the animation. 1 is fast
	frameNum int

	bottomWidth  int
	bottomHeight int
	bottom0X     int
	bottom0Y     int

	topWidth  int
	topHeight int
	top0X     int
	top0Y     int
}

func init() {
	var err error
	sprite, _, err = ebitenutil.NewImageFromFile("./assets/11226.gif", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
}

func NewMarco() *Marco {
	return &Marco{
		speed:        8,
		frameNum:     4,
		bottomWidth:  24,
		bottomHeight: 16,
		bottom0X:     155,
		bottom0Y:     1536,

		top0X:     12,
		top0Y:     8,
		topWidth:  33,
		topHeight: 29,
	}
}

func (c *Marco) frameWidth() int {
	if c.topWidth > c.bottomWidth {
		return c.topWidth
	}
	return c.bottomWidth
}

func (c *Marco) frameHeight() int {
	return c.topHeight + c.bottomHeight - int(float64(c.bottomHeight)/2)
}

func (c *Marco) Update(tick int) {
	c.tick = tick
}

func (c *Marco) Draw() *ebiten.Image {
	marco, err := ebiten.NewImage(c.frameWidth(), c.frameHeight(), ebiten.FilterNearest)
	if err != nil {
		log.Fatal(err)
	}
	marco.Fill(color.White)

	frame := (c.tick / c.speed) % c.frameNum
	if err := marco.DrawImage(c.drawLegs(frame)); err != nil {
		log.Fatal(err)
	}

	if err := marco.DrawImage(c.drawBody(frame)); err != nil {
		log.Fatal(err)
	}

	return marco
}

func (c *Marco) drawBody(frame int) (*ebiten.Image, *ebiten.DrawImageOptions) {
	options := &ebiten.DrawImageOptions{}
	topX, topY := c.top0X+frame*c.topWidth, c.top0Y
	return sprite.SubImage(image.Rect(topX, topY, topX+c.topWidth, topY+c.topHeight)).(*ebiten.Image), options
}

func (c *Marco) drawLegs(frame int) (*ebiten.Image, *ebiten.DrawImageOptions) {
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(0, float64(c.topHeight)-float64(c.bottomHeight)/2)
	bottomX, bottomY := c.bottom0X+frame*c.bottomWidth, c.bottom0Y
	return sprite.SubImage(image.Rect(bottomX, bottomY, bottomX+c.bottomWidth, bottomY+c.bottomHeight)).(*ebiten.Image), options
}
