package main

import (
	"image"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type Environment interface {
	Draw(win *pixelgl.Window, at pixel.Vec)
}

type ImageEnv struct {
	sprite *pixel.Sprite
}

func (s *ImageEnv) Draw(win *pixelgl.Window, at pixel.Vec) {

	middleX := win.Bounds().W() / 2
	middleY := win.Bounds().H() / 2
	screenMiddle := pixel.V(middleX, middleY)

	spriteMiddle := s.sprite.Picture().Bounds().Size().Scaled(.5).Scaled(ImageScale)

	unitMiddle := pixel.V(8*ImageScale, 8*ImageScale)

	s.sprite.Draw(win, pixel.IM.Scaled(pixel.V(-1, -1), ImageScale).Moved(at.Add(screenMiddle).Add(spriteMiddle).Sub(unitMiddle)))

}
func NewImageEnvFromFile(path string) (*ImageEnv, error) {
	img, err := loadImageImage(path)
	if err != nil {
		return nil, err
	}
	pic := pixel.PictureDataFromImage(img)
	rect := img.Bounds()
	sp := pixel.NewSprite(pic, pixel.R(float64(rect.Min.X), float64(rect.Min.Y), float64(rect.Max.X), float64(rect.Max.Y)))

	ie := &ImageEnv{
		sprite: sp,
	}
	return ie, nil
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

type GridEnv struct {
	imd *imdraw.IMDraw
}

func (g *GridEnv) Draw(win pixel.Target, at pixel.Vec) {

	g.imd.Draw(win)
}

func (g *GridEnv) MakeGrid() {
	grid := imdraw.New(nil)

	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			if (x%2 == 0) != (y%2 == 0) {
				grid.Color = colornames.Red
			} else {
				grid.Color = colornames.Orange
			}
			grid.Push(pixel.V(float64(x*16)*ImageScale, float64(y*16)*ImageScale))
			grid.Push(pixel.V(float64((x+1)*16)*ImageScale, float64(y*16)*ImageScale))
			grid.Push(pixel.V(float64((x+1)*16)*ImageScale, float64((y+1)*16)*ImageScale))
			grid.Push(pixel.V(float64(x*16)*ImageScale, float64((y+1)*16)*ImageScale))
			grid.Polygon(0)
		}
	}
	g.imd = grid
}
