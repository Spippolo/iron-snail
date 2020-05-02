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
)

type Sprite struct {
	Image *ebiten.Image
	Desc  map[BodyPart]*SpriteDesc
}

type SpriteDesc struct {
	Speed int
	Tiles []*TileDesc
}

type TileDesc struct {
	X0      int
	Y0      int
	W       int
	H       int
	YOffset int `json:"y_offset"` // Tells how much th image below this one must be moved up to be visually correct. Think to the legs compared to a body, they must be attached to the body but the body image can have some space below the belt
}

func Marco() *Sprite {
	img, _, err := ebitenutil.NewImageFromFile("./assets/marco.gif", ebiten.FilterDefault)
	utils.CheckErr(err, "Error reading image file")
	j, err := ioutil.ReadFile("./assets/marco.json")
	utils.CheckErr(err, "Error reading json file")

	var s map[BodyPart]*SpriteDesc
	err = json.Unmarshal(j, &s)
	utils.CheckErr(err, "Cannot unmarshal sprite")

	return &Sprite{
		Image: img,
		Desc:  s,
	}
}
