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

type Physics struct {
	gravity   float64
	speed     float64
	jumpSpeed float64
	vel       pixel.Vec
	ground    bool
}

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
	actual_frame int32
}

type Entity struct {
	animation   *Animation
	entity_type int32
	entity_name string
	physics     Physics
	rect        pixel.Rect
	pos_x       float64
	pos_y       float64
}

func (ent *Entity) UpdateMovement(mov Movement) {
	if mov.action_name != ent.animation.movement.action_name {
		ent.animation.movement = mov
		ent.animation.counter = 0
		ent.animation.actual_frame = 0
		ent.physics.vel.X = 0
		ent.physics.vel.Y = 0
	}
}
func (ent *Entity) UpdatePhysics(vector pixel.Vec) {

	x, y := vector.XY()
	ent.physics.vel.X = x
	ent.physics.vel.Y = y

}

func (mov Movement) GetFrame(index int32) *pixel.Sprite {
	for _, v := range mov.sprites {
		if v.Sprite_number == index+1 {
			return v.Sprite
		}
	}
	return nil
}
func (ent *Entity) update(dt float64) {
	ent.animation.counter += dt
	if ent.animation.counter >= ent.animation.rate {
		ent.animation.actual_frame++
		ent.animation.sprite = ent.animation.movement.GetFrame(ent.animation.actual_frame % (ent.animation.movement.total_sprites - 1))
		ent.animation.counter = 0
		ent.pos_x = ent.pos_x + ent.physics.vel.X*ent.physics.speed
	}
}

func (ent *Entity) draw(t pixel.Target) {
	if ent.animation.sprite == nil {
		ent.animation.sprite = ent.animation.movement.GetFrame(0)
	}
	ent.animation.sprite.Draw(t, pixel.IM.Scaled(pixel.ZV, 0.2).Chained(pixel.IM.Moved(pixel.V(ent.pos_x, 300))))
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
		movement:     arrMov["Idle"],
		rate:         1.0 / 10.0,
		counter:      0.0,
		actual_frame: 0,
	}
	characterEntity := Entity{animation: anim, entity_name: "cat", entity_type: 1, pos_x: 200, pos_y: 300, physics: Physics{gravity: 1, speed: 1}}

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

		if win.JustPressed(pixelgl.KeySpace) {
			characterEntity.UpdateMovement(arrMov["Jump"])
		}

		if win.Pressed(pixelgl.KeyRight) {
			characterEntity.UpdateMovement(arrMov["Walk"])
			characterEntity.UpdatePhysics(pixel.V(1, 0))
		}

		if win.JustReleased(pixelgl.KeyRight) {
			characterEntity.UpdateMovement(arrMov["Idle"])
			characterEntity.UpdatePhysics(pixel.V(0, 0))
		}

		canvas.Clear(colornames.Black)
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
