package sprites

import (
	"encoding/json"
	"io/ioutil"

	"github.com/Spippolo/iron-snail/utils"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type BodyPart string

const (
	LegsStandingPart BodyPart = "legs-standing-base"
	LegsRunningPart  BodyPart = "legs-running-base"
	BodyStandingPart BodyPart = "body-standing-base"
	BodyShootingPart BodyPart = "body-shooting-base"
	BodyKnifePart    BodyPart = "body-knife-base"
	BodyKnifeUpPart  BodyPart = "body-knife-up"
)

type Sprite struct {
	Image *ebiten.Image
	Desc  map[BodyPart]*SpriteDesc
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

func Marco() *Sprite {
	img, _, err := ebitenutil.NewImageFromFile("./assets/marco.gif", ebiten.FilterDefault)
	utils.CheckErr(err, "Error reading image file")
	j, err := ioutil.ReadFile("./assets/marco.json")
	utils.CheckErr(err, "Error reading json file")

	var sprite map[BodyPart]*SpriteDesc
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
