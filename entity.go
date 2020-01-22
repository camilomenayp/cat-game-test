package main

import (
	"fmt"

	"github.com/faiface/pixel"
)

type Entity struct {
	animation   *Animation
	entity_type string
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
		if ent.pos_y+ent.physics.vel.Y*dt >= 214.0 && ent.animation.movement.action_name == "Jump" {
			ent.pos_y = ent.pos_y + ent.physics.vel.Y*dt
			ent.physics.vel.Y = ent.physics.vel.Y + ent.physics.gravity*dt
		} else if ent.animation.movement.action_name == "Jump" {
			ent.pos_y = 214
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
