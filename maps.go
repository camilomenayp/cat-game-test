package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/faiface/pixel"
)

type Tiles struct {
	tile_number int32
	sprite      *pixel.Sprite
	isGround    bool
	pos_x       float64
	pos_y       float64
	width       float64
	height      float64
}

type MapDefinition struct {
	tiles  *[][]Tiles
	size_x int32
	size_y int32
}

var TILE_WIDTH = 32.0
var TILE_HEIGHT = 32.0
var GROUND_TILES = [...]int32{1, 2, 3, 7, 11, 13, 14, 15}

func (md *MapDefinition) draw(t pixel.Target, dt float64) {
	for _, tileLine := range *md.tiles {
		for _, tile := range tileLine {
			if tile.tile_number != -1 {
				tile.sprite.Draw(t, pixel.IM.Moved(pixel.V(tile.pos_x, tile.pos_y)))
			}
		}
	}

}
func GetMap(mapName string) MapDefinition {
	return loadMapCsv(mapName)
}

func loadMapCsv(mapName string) MapDefinition {
	arrTiles := make([][]Tiles, 0)
	csvFile, _ := os.Open("maps/" + mapName)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	tilesSpriteSheet := spriteManager.GetSpriteSheet("tiles")
	i := 0
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		for j, v := range line {
			n, err := strconv.ParseInt(v, 10, 32)
			if err != nil {
				panic(err)
			}
			if i+1 > len(arrTiles) {
				arrTiles = append(arrTiles, make([]Tiles, 0))
			}
			isGround := false
			if checkIfGround(int32(n)) {
				isGround = true
			}
			var sprite *pixel.Sprite
			if v != "-1" {
				sprite = (*tilesSpriteSheet.Sprites)[v][0].Sprite
			}

			tile := Tiles{
				tile_number: int32(n),
				isGround:    isGround,
				pos_x:       float64(j) * TILE_WIDTH,
				pos_y:       768 - float64(i)*TILE_HEIGHT - TILE_HEIGHT,
				width:       TILE_WIDTH,
				height:      TILE_HEIGHT,
				sprite:      sprite,
			}
			arrTiles[i] = append(arrTiles[i], tile)

		}
		i++
	}

	return MapDefinition{tiles: &arrTiles, size_x: int32(len(arrTiles[0])), size_y: int32(len(arrTiles))}
}

func checkIfGround(n int32) bool {
	for _, v := range GROUND_TILES {
		if v == n {
			return true
		}
	}
	return false
}
