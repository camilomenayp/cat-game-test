package main

import (
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type config struct {
	cfg pixelgl.WindowConfig
}

var (
	frames = 0
	second = time.Tick(time.Second)
)

func getConfig() pixelgl.WindowConfig {
	cfg := pixelgl.WindowConfig{
		Title:  "Cat game test",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	return cfg
}
