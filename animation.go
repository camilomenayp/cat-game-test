package main

import "github.com/faiface/pixel"

type Animation struct {
	movement     Movement
	rate         float64
	counter      float64
	sprite       *pixel.Sprite
	actual_frame int32
	finished     bool
}

type Movement struct {
	action_name        string
	type_name          string
	total_sprites      int32
	sprites            []Sprite
	sprites_sheet      *pixel.Picture
	can_be_interrupted bool
}

type Physics struct {
	gravity   float64
	speed     float64
	jumpSpeed float64
	vel       pixel.Vec
	ground    bool
}

func createMovement(arrMovs *map[string]Movement, arrSprites *map[string][]Sprite, action string, pic *pixel.Picture, canBeInterrupted bool) {
	mov := (*arrSprites)[action]
	(*arrMovs)[action] = Movement{action_name: action, type_name: "cat", total_sprites: int32(len(mov)), sprites: mov, sprites_sheet: pic, can_be_interrupted: canBeInterrupted}
}

func initMovements() (map[string]Movement, error) {
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

func (mov Movement) GetFrame(index int32) *pixel.Sprite {
	for _, v := range mov.sprites {
		if v.Sprite_number == index+1 {
			return v.Sprite
		}
	}
	return nil
}
