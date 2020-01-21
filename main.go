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
	action_name        string
	type_name          string
	total_sprites      int32
	sprites            []sprite.Sprite
	sprites_sheet      *pixel.Picture
	can_be_interrupted bool
}

type Animation struct {
	movement     Movement
	rate         float64
	counter      float64
	sprite       *pixel.Sprite
	actual_frame int32
	finished     bool
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

func (ent *Entity) UpdateMovement(mov Movement, vector pixel.Vec) {

	if (ent.animation.finished && ent.animation.movement.action_name != "Idle") || (ent.animation.movement.can_be_interrupted && mov.action_name != ent.animation.movement.action_name) {
		fmt.Println("[Game] [Entity:", ent.entity_name, "]", "Starting action:", mov.action_name)
		ent.animation.movement = mov
		ent.animation.counter = 0
		ent.animation.actual_frame = 0
		ent.animation.finished = false
		ent.UpdatePhysics(vector)
	}
}
func (ent *Entity) UpdatePhysics(vector pixel.Vec) {

	if vector != ent.physics.vel {
		x, y := vector.XY()
		ent.physics.vel.X = x
		ent.physics.vel.Y = y
		fmt.Println("[Game] [Entity:", ent.entity_name, "]", "Physics vector updated to:(", x, y, ")")
	}
}

func (mov Movement) GetFrame(index int32) *pixel.Sprite {
	for _, v := range mov.sprites {
		if v.Sprite_number == index+1 {
			return v.Sprite
		}
	}
	return nil
}

func (ent *Entity) checkMovementStatus() {

	if ent.animation.actual_frame == ent.animation.movement.total_sprites-1 && ent.animation.movement.action_name != "Jump" {
		ent.animation.finished = true
	}
}
func (ent *Entity) update(dt float64) {
	ent.animation.counter += dt
	x := ent.pos_x
	y := ent.pos_y
	if ent.animation.counter >= ent.animation.rate {
		ent.pos_x = ent.pos_x + ent.physics.vel.X*dt
		if ent.pos_y+ent.physics.vel.Y*dt >= 300.0 && ent.animation.movement.action_name == "Jump" {
			ent.pos_y = ent.pos_y + ent.physics.vel.Y*dt
			ent.physics.vel.Y = ent.physics.vel.Y + ent.physics.gravity*dt
		} else if ent.animation.movement.action_name == "Jump" {
			ent.pos_y = 300
			ent.animation.finished = true
		}

		ent.animation.actual_frame++
		ent.animation.sprite = ent.animation.movement.GetFrame(ent.animation.actual_frame % (ent.animation.movement.total_sprites - 1))
		ent.animation.counter = 0

		if x != ent.pos_x || y != ent.pos_y {
			fmt.Println("[Game] [Entity:", ent.entity_name, "]", "Position updated to:(", ent.pos_x, ent.pos_y, ")")
			fmt.Println("[Game] [Entity:", ent.entity_name, "]", "Speed vector:(", ent.physics.vel.X, ent.physics.vel.Y, ")")

		}
	}
}

func (ent *Entity) draw(t pixel.Target) {
	if ent.animation.sprite == nil {
		ent.animation.sprite = ent.animation.movement.GetFrame(0)
	}
	ent.animation.sprite.Draw(t, pixel.IM.Scaled(pixel.ZV, 0.2).Chained(pixel.IM.Moved(pixel.V(ent.pos_x, ent.pos_y))))
}

var characterSpriteSheet sprite.SpriteSheet

func createMovement(arrMovs *map[string]Movement, arrSprites *map[string][]sprite.Sprite, action string, pic *pixel.Picture, canBeInterrupted bool) {
	mov := (*arrSprites)[action]
	(*arrMovs)[action] = Movement{action_name: action, type_name: "cat", total_sprites: int32(len(mov)), sprites: mov, sprites_sheet: pic, can_be_interrupted: canBeInterrupted}
}

func initMovements() (map[string]Movement, error) {
	characterSpriteSheet = sprite.LoadAllSprites("cat", "sprites.txt", "spritesheet.png")
	arr := *characterSpriteSheet.Sprites
	arrMovements := make(map[string]Movement, 0)
	createMovement(&arrMovements, &arr, "Jump", characterSpriteSheet.Sprites_sheet, false)
	createMovement(&arrMovements, &arr, "Walk", characterSpriteSheet.Sprites_sheet, true)
	createMovement(&arrMovements, &arr, "Fall", characterSpriteSheet.Sprites_sheet, false)
	createMovement(&arrMovements, &arr, "Hurt", characterSpriteSheet.Sprites_sheet, false)
	createMovement(&arrMovements, &arr, "Idle", characterSpriteSheet.Sprites_sheet, true)
	createMovement(&arrMovements, &arr, "Slide", characterSpriteSheet.Sprites_sheet, false)
	createMovement(&arrMovements, &arr, "Run", characterSpriteSheet.Sprites_sheet, true)
	createMovement(&arrMovements, &arr, "Dead", characterSpriteSheet.Sprites_sheet, false)

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
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	arrMov, err := initMovements()
	anim := &Animation{
		movement:     arrMov["Idle"],
		rate:         1.0 / 15.0,
		counter:      0.0,
		actual_frame: 0,
	}
	characterEntity := Entity{animation: anim, entity_name: "cat", entity_type: 1, pos_x: 200, pos_y: 300, physics: Physics{gravity: -3000, speed: 1, jumpSpeed: 700}}

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
