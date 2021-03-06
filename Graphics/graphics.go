package graphics

import (
	"encoding/json"
	"image"
	"io"
	"os"

	"image/color"
	_ "image/jpeg"
	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type SpriteGroup struct {
	Sprites map[string]*FlippableSprite
}
type FlippableSprite struct {
	s    *pixel.Sprite
	flip bool
}

func (f *FlippableSprite) DrawScreenPosition(win pixel.Target, screenLocation pixel.Vec, ImageScale float64) {
	f.s.Draw(win, pixel.IM.Scaled(pixel.ZV, ImageScale).Moved(screenLocation))
}
func (f *FlippableSprite) DrawWorldPosition(win *pixelgl.Window, worldLocation pixel.Vec, ImageScale float64) {
	scaleVec := pixel.V(1, 1)
	if f.flip {
		scaleVec.X = -1
	}

	middleX := win.Bounds().W() / 2
	middleY := win.Bounds().H() / 2
	middle := pixel.V(middleX, middleY)
	//Calculate window position
	//TODO revisit this. x should be centered. bottom of sprite is actually where it is so draw it there. Not actually sure
	//x is centered.
	//y should be centered on the bottom 16. Here middle is 0
	//if 20 pixels tall, middle of that is ten pixels up. difference is 2. Add 2
	//if 32 pixels tall, middle is 16, difference is 8
	middleSprite := f.s.Frame().H() / 2
	offset := middleSprite - 6
	x := (worldLocation.X * 16) * ImageScale
	y := (worldLocation.Y*16 + offset) * ImageScale

	f.s.Draw(win, pixel.IM.ScaledXY(pixel.V(0, 0), scaleVec.Scaled(ImageScale)).Moved(pixel.V(x, y).Add(middle)))
}

type AnimationSheet struct {
	Sprites map[string]SpriteLocation
}

type SpriteLocation struct {
	Minx, Miny, Maxx, Maxy float32
	Flip                   bool
}

func (sl SpriteLocation) Rect() pixel.Rect {
	return pixel.R(float64(sl.Minx), float64(sl.Miny), float64(sl.Maxx), float64(sl.Maxy))
}

func LoadSprite(imagePath, dataPath string) (*SpriteGroup, error) {
	img, err := loadImageImage(imagePath)
	if err != nil {
		return nil, err
	}
	pic := pixel.PictureDataFromImage(img)
	makeTransparent(pic, pic.Color(pixel.V(0, 0)))

	dataFile, err := os.Open(dataPath)
	if err != nil {
		return nil, err
	}
	dataBytes, err := io.ReadAll(dataFile)
	if err != nil {
		return nil, err
	}

	var sg SpriteGroup = SpriteGroup{
		Sprites: map[string]*FlippableSprite{},
	}
	var as AnimationSheet = AnimationSheet{
		Sprites: map[string]SpriteLocation{},
	}
	err = json.Unmarshal(dataBytes, &as)
	if err != nil {
		return nil, err
	}

	for k, v := range as.Sprites {
		sg.Sprites[k] = &FlippableSprite{
			s:    pixel.NewSprite(pic, v.Rect()),
			flip: v.Flip,
		}
	}

	return &sg, nil
}

func loadImageImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func makeTransparent(img *pixel.PictureData, transColor pixel.RGBA) {
	for y := 0; y < int(img.Bounds().H()); y++ {
		for x := 0; x < int(img.Bounds().W()); x++ {
			if img.Color(pixel.V(float64(x), float64(y))) == transColor {
				img.Pix[img.Index(pixel.V(float64(x), float64(y)))] = color.RGBA{0, 0, 0, 0}
			}
		}
	}
}
