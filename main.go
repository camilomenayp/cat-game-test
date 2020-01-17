package main

import (
	"fmt"
	_ "image/png"
	"math"
	"time"

	"github.com/camilomenayp/cat-game-test/pkg/sprite"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type Movement = sprite.Movement

type Animation struct {
	movement Movement
	rate     float64
	counter  float64
	frame    *pixel.Sprite
}

func (an *Animation) update(dt float64) {
	an.counter += dt
	i := int32(math.Floor(an.counter / an.rate))
	an.frame = an.movement.GetMovementSpriteFrame(i % (an.movement.GetTotalSprites() - 1))
}
func initMovements() ([]Movement, error) {
	arrMovements := make([]Movement, 0)
	/*
		arrMovements = append(arrMovements, Movement{"Dead", 10, nil})
		arrMovements = append(arrMovements, Movement{"Fall", 8, nil})
		arrMovements = append(arrMovements, Movement{"Hurt", 10, nil})
		arrMovements = append(arrMovements, Movement{"Idle", 10, nil})
		arrMovements = append(arrMovements, Movement{"Jump", 8, nil})
		arrMovements = append(arrMovements, Movement{"Run", 8, nil})
		arrMovements = append(arrMovements, Movement{"Slide", 10, nil})
		arrMovements = append(arrMovements, Movement{"Walk", 10, nil}) */

	mov := Movement{}
	mov.Init("Walk", "cat", 10)
	err := mov.SetSpritesFrames()
	if err != nil {
		return nil, (err)
	}
	arrMovements = append(arrMovements, mov)
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
		movement: arrMov[0],
		rate:     1.0 / 10,
		counter:  0.0,
	}
	if err != nil {
		panic(err)
	}
	last := time.Now()

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		win.Clear(colornames.Forestgreen)

		batch := pixel.NewBatch(&pixel.TrianglesData{}, *anim.movement.GetSpriteSheet())
		anim.update(dt)
		anim.frame.Draw(batch, pixel.IM.Scaled(pixel.ZV, 0.2).Chained(pixel.IM.Moved(win.Bounds().Center())))
		batch.Draw(win)
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
