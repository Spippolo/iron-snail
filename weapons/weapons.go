package weapons

import (
	"image"

	"github.com/Spippolo/iron-snail/common"
	"github.com/Spippolo/iron-snail/sprites"
	"github.com/hajimehoshi/ebiten"
	log "github.com/sirupsen/logrus"
)

type Weapon int

const (
	Gun Weapon = iota
	HeavyMachineGun
	RocketLauncher
	ShotGun
	FlameShot
	LaserGun
)

type Bullet interface {
	Update() // Update bullets position
	Draw() (*ebiten.Image, *ebiten.DrawImageOptions)
}

type GunBullet struct {
	x            float64
	y            float64
	sprite       *sprites.Sprite
	currentFrame int
	direction    common.Direction
}

func NewBullet(w Weapon, d common.Direction) Bullet {
	if w == Gun {
		return newGunBullet(d)
	}
	return nil
}

func newGunBullet(d common.Direction) *GunBullet {
	return &GunBullet{
		direction: d,
		sprite:    sprites.Weapons(),
	}
}

func (g *GunBullet) Update() {
	g.currentFrame++

	x := g.sprite.Desc[sprites.BulletBase].Speed
	if g.direction == common.West {
		x *= -1
	}
	g.x += float64(x)
}

func (g *GunBullet) Draw() (*ebiten.Image, *ebiten.DrawImageOptions) {
	s := g.sprite.Desc[sprites.BulletBase]
	// Number of frames for this part
	if s == nil {
		log.Fatal("nil sprite for GunBullet")
	}
	frame := (g.currentFrame / s.Speed) % s.Frames
	t := s.Tiles[frame]

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.x, 0) // TODO Y is 0 for horizontal shoots, but not for other shoots

	return g.sprite.Image.SubImage(image.Rect(t.X0, t.Y0, t.X0+t.W, t.Y0+t.H)).(*ebiten.Image), op
}
