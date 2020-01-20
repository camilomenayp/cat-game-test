package sprite

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"image"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/faiface/pixel"
)

type Size struct {
	sizex float64
	sizey float64
}

type SpriteSheet struct {
	Sprites_sheet *pixel.Picture
	Sprites       *map[string][]Sprite
}

type Sprite struct {
	Sprite_number int32
	max_x         int32
	max_y         int32
	height        int32
	width         int32
	Sprite        *pixel.Sprite
	name          string
	total_sprites int32
}

func LoadAllSprites(object_name string, csvPath string, spritePath string) SpriteSheet {
	sprInfo := loadCsv("images/png/" + object_name + "/" + csvPath)
	pic, _ := loadPicture("images/png/" + object_name + "/" + spritePath)
	//width := int32(pic.Bounds().W())
	height := int32(pic.Bounds().H())
	for i, v := range sprInfo {
		for j, w := range v {
			rect := pixel.R(float64(w.max_x), float64(height-w.max_y-w.height), float64(w.max_x+w.width), float64(height-w.max_y))
			fmt.Println("[Sprite] Building action: *", i, "(", w.Sprite_number, ")*")
			fmt.Println("[Sprite] ", rect.String())
			sprInfo[i][j].Sprite = pixel.NewSprite(pic, rect)

			sprInfo[i][j].total_sprites = int32(len(v))
		}

	}
	return SpriteSheet{&pic, &sprInfo}
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func loadCsv(path string) map[string][]Sprite {
	csvFile, _ := os.Open(path)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	sprites := make(map[string][]Sprite)
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		ini := strings.IndexRune(line[0], '(')
		end := strings.IndexRune(line[0], ')')

		number, _ := strconv.ParseInt(line[0][ini+1:end], 10, 32)
		max_x, _ := strconv.ParseInt(line[1], 10, 32)
		max_y, _ := strconv.ParseInt(line[2], 10, 32)
		height, _ := strconv.ParseInt(line[4], 10, 32)
		width, _ := strconv.ParseInt(line[3], 10, 32)

		spr := Sprite{
			Sprite_number: int32(number),
			name:          strings.Split(line[0], " ")[0],
			max_x:         int32(max_x),
			max_y:         int32(max_y),
			height:        int32(height),
			width:         int32(width),
		}
		if val, ok := sprites[strings.Split(line[0], " ")[0]]; ok {
			sprites[strings.Split(line[0], " ")[0]] = append(val, spr)
		} else {
			sprites[strings.Split(line[0], " ")[0]] = append(sprites[strings.Split(line[0], " ")[0]], spr)
		}
	}
	return sprites
}
