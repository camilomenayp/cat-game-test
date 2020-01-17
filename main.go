package main

import (
	"fmt"
	_ "image/png"
	"time"

	"github.com/camilomenayp/cat-game-test/pkg/sprite"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type Movement = sprite.Movement

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
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	arrMov, err := initMovements()
	if err != nil {
		panic(err)
	}
	//last := time.Now()

	for !win.Closed() {
		//last = time.Now()
		//dt := time.Since(last).Seconds()

		win.Clear(colornames.Forestgreen)

		batch := pixel.NewBatch(&pixel.TrianglesData{}, *((arrMov)[0].GetSpriteSheet()))
		sprite := *(arrMov)[0].GetMovementSpriteFrame()
		sprite.Draw(batch, pixel.IM.Scaled(pixel.ZV, 0.2).Chained(pixel.IM.Moved(win.Bounds().Center())))
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
