package sprite

import (
	"image"
	"os"

	"github.com/faiface/pixel"
)

type Size struct {
	sizex float64
	sizey float64
}

type SpriteSheet struct {
	sprite_sheet     *pixel.Picture
	sprites_position []SpritePosition
}

type SpritePosition struct {
	sprite_number int32
	pos_x         int32
	pos_y         int32
	size_x        int32
	size_y        int32
	sprite_rect   pixel.Rect
}

type Movement struct {
	action_name          string
	type_name            string
	total_sprites        int32
	sprite_sheet         *SpriteSheet
	actual_sprite_number int32
}

func (mov *Movement) SetSpritesFrames() error {

	sprite, err := getSpriteSheet(mov.type_name, mov.action_name+"-sheet.png", mov.total_sprites)
	if err != nil {
		return err
	}
	mov.sprite_sheet = sprite

	return nil
}

func (mov *Movement) Init(action_name string, type_name string, total_sprites int32) {
	mov.total_sprites = total_sprites
	mov.action_name = action_name
	mov.type_name = type_name
	mov.actual_sprite_number = -1

}

func (mov *Movement) GetSpriteSheet() *pixel.Picture {
	return mov.sprite_sheet.sprite_sheet
}

func (mov *Movement) GetMovementSpriteFrame() *pixel.Sprite {
	n := mov.actual_sprite_number
	if n == -1 || n+1 > mov.total_sprites-1 {
		n = 0
	} else {
		n++
	}
	sprite := pixel.NewSprite(*mov.sprite_sheet.sprite_sheet, mov.sprite_sheet.sprites_position[n].sprite_rect)
	mov.actual_sprite_number = n
	return sprite
}
func getSpriteSheet(folder string, path string, n_sprites int32) (*SpriteSheet, error) {
	pic, err := loadPicture("images/png/" + folder + "/" + path)
	if err != nil {
		return nil, err
	}

	height := pic.Bounds().H()
	width := pic.Bounds().W()

	size_x := int32(width) / n_sprites
	size_y := int32(height)
	sprite_pos := make([]SpritePosition, 0)
	for i := int32(0); i < n_sprites; i++ {
		sprite_pos = append(sprite_pos, SpritePosition{i + 1, i * size_x, 0 * size_y, size_x, size_y, pixel.R(float64(i*size_x), 0, float64(i*size_x+size_x), float64(size_y))})
	}

	frame := &SpriteSheet{&pic, sprite_pos}
	return frame, nil
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
