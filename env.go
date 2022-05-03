package main

import (
	"image"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

type Environment interface {
	Draw(win pixel.Target, at pixel.Vec)
}

type ImageEnv struct {
	sprite *pixel.Sprite
}

func (s *ImageEnv) Draw(win pixel.Target, at pixel.Vec) {
	s.sprite.Draw(win, pixel.IM.Moved(s.sprite.Frame().Max.Scaled(.5)).Scaled(pixel.ZV, ImageScale).Moved(at))
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
