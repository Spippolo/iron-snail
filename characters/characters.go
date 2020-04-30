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

type tile struct {
	x0 int
	y0 int
	w  int
	h  int
}

type Character interface {
	Update(int)
	Draw() *ebiten.Image
}

type part string

const (
	legsStandingPart part = "legs-standing"
	legsRunningPart       = "legs-running"
	bodyStandingPart      = "body-standing"
	bodyShootingPart      = "body-shooting"
)

type Marco struct {
	tick        int
	speed       int // slow down the animation. 1 is fast
	legsYOffset int // TODO: Bad design. This is the amount of pixels the legs must be moved down to be in the correct position compared to the body

	parts map[part][]*tile
}

func init() {
	var err error
	sprite, _, err = ebitenutil.NewImageFromFile("./assets/11226.gif", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
}

func NewMarco() *Marco {
	m := Marco{
		speed: 8,
		parts: make(map[part][]*tile),
	}

	m.parts[legsStandingPart] = m.buildParts(legsStandingPart)
	m.parts[bodyStandingPart] = m.buildParts(bodyStandingPart)
	return &m
}

func (c *Marco) buildParts(pName part) []*tile {
	t := []*tile{}
	switch pName {
	case legsStandingPart:
		x0 := 155
		y0 := 1536
		w := 24
		h := 16
		for i := 0; i < 4; i++ {
			t = append(t, &tile{x0: x0 + w*i, y0: y0, w: w, h: h})
		}
		return t
	case bodyStandingPart:
		x0 := 12
		y0 := 8
		w := 33
		h := 29
		c.legsYOffset = h
		for i := 0; i < 4; i++ {
			t = append(t, &tile{x0: x0 + w*i, y0: y0, w: w, h: h})
		}
		return t
	default:
		log.Fatalf("Unknown part %s", pName)
	}
	return nil
}

func (c *Marco) Update(tick int) {
	c.tick = tick
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (c *Marco) Draw() *ebiten.Image {
	legsImage, legsOptions := c.drawLegs(c.legsYOffset)
	legsW, legsH := legsImage.Size()

	bodyImage, bodyOptions := c.drawBody()
	bodyW, bodyH := bodyImage.Size()

	frameHeight := bodyH + legsH - int(float64(legsH)/2)

	marco, err := ebiten.NewImage(max(legsW, bodyW), frameHeight, ebiten.FilterNearest)
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

func (c *Marco) drawBody() (*ebiten.Image, *ebiten.DrawImageOptions) {
	// Number of frames for this part
	frameNum := len(c.parts[bodyStandingPart])
	frame := (c.tick / c.speed) % frameNum
	t := c.parts[bodyStandingPart][frame]
	options := &ebiten.DrawImageOptions{}
	return sprite.SubImage(image.Rect(t.x0, t.y0, t.x0+t.w, t.y0+t.h)).(*ebiten.Image), options
}

func (c *Marco) drawLegs(offset int) (*ebiten.Image, *ebiten.DrawImageOptions) {
	// Number of frames for this part
	frameNum := len(c.parts[legsStandingPart])
	frame := (c.tick / c.speed) % frameNum
	t := c.parts[legsStandingPart][frame]

	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(0, float64(offset)-float64(t.h)/2)
	return sprite.SubImage(image.Rect(t.x0, t.y0, t.x0+t.w, t.y0+t.h)).(*ebiten.Image), options
}
