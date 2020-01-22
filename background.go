package main

import (
	"math"

	"github.com/faiface/pixel"
)

type ScrollingBackground struct {
	width        float64
	height       float64
	displacement float64
	speed        float64
	backgrounds  [2]*pixel.Sprite
	positions    [2]pixel.Vec
}

func CreateBackground(path string) *ScrollingBackground {
	pic, _ := loadPicture(path)
	return NewScrollingBackground(pic, pic.Bounds().W(), pic.Bounds().H())
}
func NewScrollingBackground(pic pixel.Picture, width, height float64) *ScrollingBackground {
	sb := &ScrollingBackground{
		width:  width,
		height: height,
		backgrounds: [2]*pixel.Sprite{
			pixel.NewSprite(pic, pixel.R(0, 0, width, height)),
			pixel.NewSprite(pic, pixel.R(0, 0, width, height)),
		},
	}

	sb.positionImages()
	return sb
}

func (sb *ScrollingBackground) positionImages() {

	sb.positions = [2]pixel.Vec{
		pixel.V(sb.width/2, sb.height/2),
		pixel.V(sb.width+(sb.width/2), sb.height/2),
	}

}

func (sb *ScrollingBackground) Update(win pixel.Target, dt float64, vel pixel.Vec) {
	if math.Abs(sb.displacement) >= sb.width {
		sb.displacement = 0
		sb.positions[0], sb.positions[1] = sb.positions[1], sb.positions[0]
	}
	d := pixel.V(sb.displacement, 0)
	sb.backgrounds[0].Draw(win, pixel.IM.Moved(sb.positions[0].Add(d)))
	sb.backgrounds[1].Draw(win, pixel.IM.Moved(sb.positions[1].Add(d)))
	sb.displacement += -vel.X * dt
}
