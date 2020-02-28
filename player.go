package eyes

import (
	"image"
	"image/color"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

func LoadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func LoadSession(path string) (Session, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	session, err := DecodeSession(file)
	if err != nil {
		return nil, err
	}
	return session, nil
}

type Player struct {
	Picture  pixel.Picture
	Session  Session
	Time     time.Duration
	Position int

	spr *pixel.Sprite
	imd *imdraw.IMDraw
}

func LoadPlayer(picPath, sessionPath string) (*Player, error) {
	pic, err := LoadPicture(picPath)
	if err != nil {
		return nil, err
	}
	session, err := LoadSession(sessionPath)
	if err != nil {
		return nil, err
	}
	return &Player{
		Picture:  pic,
		Session:  session,
		Time:     0,
		Position: 0,
		spr:      pixel.NewSprite(pic, pic.Bounds()),
		imd:      imdraw.New(nil),
	}, nil
}

func (p *Player) Finished() bool {
	return p.Position >= len(p.Session)-1
}

func (p *Player) Update(dt time.Duration) {
	p.Time += dt
	for p.Position < len(p.Session)-1 && p.Time > p.Session[p.Position].Time {
		p.Position++
	}
	for p.Position > 0 && p.Time <= p.Session[p.Position].Time {
		p.Position--
	}
}

func (p *Player) LeftEye() pixel.Vec {
	return p.Session[p.Position].LeftPos
}

func (p *Player) RightEye() pixel.Vec {
	return p.Session[p.Position].RightPos
}

func (p *Player) Draw(t pixel.Target, picture bool) {
	p.imd.Clear()

	// picture
	if picture {
		p.spr.Draw(t, pixel.IM.Moved(p.Picture.Bounds().Center()))
	}

	// eyes
	p.drawEye(p.LeftEye(), colornames.Blue)
	p.drawEye(p.RightEye(), colornames.Red)

	p.imd.Draw(t)
}

func (p *Player) drawEye(pos pixel.Vec, clr color.Color) {
	p.imd.Color = colornames.White
	p.imd.Push(pos)
	p.imd.Circle(4, 0)
	p.imd.Color = clr
	p.imd.Push(pos)
	p.imd.Circle(4, 2)
}
