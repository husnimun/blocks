package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/husnimun/blocks/assets"
)

var imageBlocks *ebiten.Image

const (
	blockHeight = 10
	blockWidth  = 10
)

func init() {
	f, err := assets.Assets.Open("images/blocks.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	imageBlocks = ebiten.NewImageFromImage(img)
}

type BlockType int

const (
	BlockTypeNone BlockType = iota
	BlockType1
	BlockType2
	BlockType3
	BlockType4
	BlockType5
	BlockType6
	BlockType7
	BlockTypeMax = BlockType7
)

func drawBlock(screen *ebiten.Image, blockType BlockType, x, y int) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))

	srcX := int((blockType - 1)) * blockWidth
	screen.DrawImage(imageBlocks.SubImage(image.Rect(srcX, 0, srcX+blockWidth, blockHeight)).(*ebiten.Image), op)
}
