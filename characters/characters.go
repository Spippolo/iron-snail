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
	speed    int
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
	// var marco *ebiten.Image
	marco, err := ebiten.NewImage(c.frameWidth(), c.frameHeight(), ebiten.FilterNearest)
	if err != nil {
		log.Fatal(err)
	}
	marco.Fill(color.White)
	// x, y := marco.Size()
	// log.Infof("Size: %d %d", x, y)

	topOp := &ebiten.DrawImageOptions{}
	// topOp.GeoM.Translate(-float64(c.topWidth)/2, -float64(c.topHeight)/2)
	// topOp.GeoM.Translate(screenWidth/2, screenHeight/2)

	bottomOp := &ebiten.DrawImageOptions{}
	// bottomOp.GeoM.Translate(-float64(c.topWidth)/2, -float64(c.topHeight)/2)
	bottomOp.GeoM.Translate(0, float64(c.topHeight)-float64(c.bottomHeight)/2)
	// bottomOp.GeoM.Translate(screenWidth/2, screenHeight/2)

	// slow down the animation. 1 is fast
	i := (c.tick / c.speed) % c.frameNum

	bottomX, bottomY := c.bottom0X+i*c.bottomWidth, c.bottom0Y
	if err := marco.DrawImage(
		sprite.SubImage(
			image.Rect(bottomX, bottomY, bottomX+c.bottomWidth, bottomY+c.bottomHeight),
		).(*ebiten.Image),
		bottomOp); err != nil {
		log.Fatal(err)
	}

	topX, topY := c.top0X+i*c.topWidth, c.top0Y
	if err := marco.DrawImage(sprite.SubImage(image.Rect(topX, topY, topX+c.topWidth, topY+c.topHeight)).(*ebiten.Image), topOp); err != nil {
		log.Fatal(err)
	}

	return marco
}
