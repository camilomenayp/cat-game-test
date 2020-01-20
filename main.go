package main

import (
	"fmt"
	_ "image/png"
	"math"
	"time"

	"github.com/camilomenayp/cat-game-test/pkg/sprite"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type Movement struct {
	action_name   string
	type_name     string
	total_sprites int32
	sprites       []sprite.Sprite
	sprites_sheet *pixel.Picture
}

type Animation struct {
	movement     Movement
	rate         float64
	counter      float64
	sprite       *pixel.Sprite
	rect         pixel.Rect
	actual_frame int32
}

func (mov Movement) GetFrame(index int32) *pixel.Sprite {
	for _, v := range mov.sprites {
		if v.Sprite_number == index+1 {
			return v.Sprite
		}
	}
	return nil
}
func (an *Animation) update(dt float64) {
	an.counter += dt
	if an.counter >= an.rate {
		an.actual_frame++
		an.sprite = an.movement.GetFrame(an.actual_frame % (an.movement.total_sprites - 1))
		an.counter = 0
	}
}

func (an *Animation) draw(t pixel.Target) {
	if an.sprite == nil {
		an.sprite = an.movement.GetFrame(0)
	}
	an.sprite.Draw(t, pixel.IM.Scaled(pixel.ZV, 0.2).Chained(pixel.IM.Moved(an.rect.Center())))
}

var characterSpriteSheet sprite.SpriteSheet

func createMovement(arrMovs *map[string]Movement, arrSprites *map[string][]sprite.Sprite, action string, pic *pixel.Picture) {
	mov := (*arrSprites)[action]
	(*arrMovs)[action] = Movement{action_name: action, type_name: "cat", total_sprites: int32(len(mov)), sprites: mov, sprites_sheet: pic}
}

func initMovements() (map[string]Movement, error) {
	characterSpriteSheet = sprite.LoadAllSprites("cat", "sprites.txt", "spritesheet.png")
	arr := *characterSpriteSheet.Sprites
	arrMovements := make(map[string]Movement, 0)
	createMovement(&arrMovements, &arr, "Jump", characterSpriteSheet.Sprites_sheet)
	createMovement(&arrMovements, &arr, "Walk", characterSpriteSheet.Sprites_sheet)
	createMovement(&arrMovements, &arr, "Fall", characterSpriteSheet.Sprites_sheet)
	createMovement(&arrMovements, &arr, "Hurt", characterSpriteSheet.Sprites_sheet)
	createMovement(&arrMovements, &arr, "Idle", characterSpriteSheet.Sprites_sheet)
	createMovement(&arrMovements, &arr, "Slide", characterSpriteSheet.Sprites_sheet)
	createMovement(&arrMovements, &arr, "Run", characterSpriteSheet.Sprites_sheet)
	createMovement(&arrMovements, &arr, "Dead", characterSpriteSheet.Sprites_sheet)

	return arrMovements, nil
}

func run() {
	var (
		frames = 0
		second = time.Tick(time.Second)
	)

	cfg := pixelgl.WindowConfig{
		Title:  "Cat game test",
		Bounds: pixel.R(0, 0, 1024, 768),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	arrMov, err := initMovements()
	anim := &Animation{
		movement:     arrMov["Jump"],
		rate:         1.0 / 10.0,
		counter:      0.0,
		rect:         pixel.R(0, 0, 542, 474),
		actual_frame: 0,
	}
	if err != nil {
		panic(err)
	}
	last := time.Now()
	canvas := pixelgl.NewCanvas(pixel.R(0, 0, 1024, 768))
	imd := imdraw.New(*anim.movement.sprites_sheet)
	imd.Precision = 32
	frameTime := float64(0)
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		frameTime += dt

		canvas.Clear(colornames.Black)
		imd.Clear()
		anim.update(dt)
		anim.draw(imd)
		imd.Draw(canvas)

		win.SetMatrix(pixel.IM.Scaled(pixel.ZV,
			math.Min(
				win.Bounds().W()/canvas.Bounds().W(),
				win.Bounds().H()/canvas.Bounds().H(),
			),
		).Moved(win.Bounds().Center()))
		canvas.Draw(win, pixel.IM)
		win.Update()
		frames++

		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
	}

}

func main() {
	pixelgl.Run(run)
}
