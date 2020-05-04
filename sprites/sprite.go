package sprites

import (
	"encoding/json"
	"io/ioutil"

	"github.com/Spippolo/iron-snail/utils"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	log "github.com/sirupsen/logrus"
)

type SpriteName string

const (
	LegsStandingPart SpriteName = "legs-standing-base"
	LegsRunningPart  SpriteName = "legs-running-base"
	BodyStandingPart SpriteName = "body-standing-base"
	BodyShootingPart SpriteName = "body-shooting-base"
	BodyKnifePart    SpriteName = "body-knife-base"
	BodyKnifeUpPart  SpriteName = "body-knife-up"
	BulletBase       SpriteName = "bullet-base"
)

type Sprite struct {
	Image *ebiten.Image
	Desc  map[SpriteName]*SpriteDesc
}

type SpriteDesc struct {
	Speed  int
	Frames int
	Tiles  []*TileDesc
}

type TileDesc struct {
	X0    int
	Y0    int
	W     int
	H     int
	Joint [2]int // The point (x,y) in the image that must be places over the joint in another image to build a grouped image. The (0,0) point is the first top-left pixel
}

var weapons *Sprite
var marco *Sprite

func Weapons() *Sprite {
	if weapons == nil {
		log.Debug("Generating weapons")
		weapons = buildSprite("./assets/weapons.png", "./assets/weapons.json")
	}
	return weapons
}

func Marco() *Sprite {
	if marco == nil {
		log.Debug("Generating Marco")
		marco = buildSprite("./assets/marco.gif", "./assets/marco.json")
	}
	return marco
}

func buildSprite(imgPath, specPath string) *Sprite {
	img, _, err := ebitenutil.NewImageFromFile(imgPath, ebiten.FilterDefault)
	utils.CheckErr(err, "Error reading image file")
	j, err := ioutil.ReadFile(specPath)
	utils.CheckErr(err, "Error reading json file")

	var sprite map[SpriteName]*SpriteDesc
	err = json.Unmarshal(j, &sprite)
	utils.CheckErr(err, "Cannot unmarshal sprite")

	// Set sprite frame numbers so that we can avoid calculating it each frame
	for _, s := range sprite {
		s.Frames = len(s.Tiles)
	}

	return &Sprite{
		Image: img,
		Desc:  sprite,
	}
}
