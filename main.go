package main

import (
	"fmt"
	_ "image/png"
	"math"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var characterSpriteSheet SpriteSheet
var spriteManager SpriteManager

func run() {

	cfg := getConfig()
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	spriteManager = SpriteManager{}
	spriteManager.LoadSpriteSheet("cat", "sprites.txt", "spritesheet.png")
	spriteManager.LoadSpriteSheet("tiles", "sprites.txt", "spritesheet.png")
	mapDefinition := GetMap("level-test-1.csv")
	characterSpriteSheet = spriteManager.GetSpriteSheet("cat")

	arrMov, err := initMovements()
	anim := &Animation{
		movement:     arrMov["Idle"],
		rate:         1.0 / 15.0,
		counter:      0.0,
		actual_frame: 0,
	}
	characterEntity := Entity{animation: anim, entity_name: "cat", entity_type: "player", pos_x: 200, pos_y: 214, physics: Physics{gravity: -4000, speed: 1, jumpSpeed: 700}}
	background := CreateBackground("images/png/bg/bg.png")

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

		characterEntity.checkMovementStatus()
		if win.JustPressed(pixelgl.KeySpace) {
			characterEntity.UpdateMovement(arrMov["Jump"], pixel.V(150, characterEntity.physics.jumpSpeed))
		} else if win.Pressed(pixelgl.KeyRight) {
			characterEntity.UpdateMovement(arrMov["Walk"], pixel.V(150, 0))
		} else {
			characterEntity.UpdateMovement(arrMov["Idle"], pixel.V(0, 0))
		}

		canvas.Clear(colornames.Black)
		background.Update(canvas, dt, pixel.V(characterEntity.physics.vel.X, 0))
		mapDefinition.draw(canvas, dt)
		imd.Clear()
		characterEntity.update(dt)
		characterEntity.draw(imd)
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
