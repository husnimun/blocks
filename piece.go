package main

import (
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Angle int

const (
	Angle0 Angle = iota
	Angle90
	Angle180
	Angle270
)

func (a Angle) RotateRight() Angle {
	if a == Angle270 {
		return Angle0
	}
	return a + 1
}

type Piece struct {
	blockType BlockType
	blocks    [][]bool
}

var Pieces map[BlockType]*Piece

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	const (
		t = true
		f = false
	)
	Pieces = map[BlockType]*Piece{
		BlockType1: {
			blockType: BlockType1,
			blocks: [][]bool{
				{f, f, f, f},
				{t, t, t, t},
				{f, f, f, f},
				{f, f, f, f},
			},
		},
		BlockType2: {
			blockType: BlockType2,
			blocks: [][]bool{
				{t, t},
				{t, t},
			},
		},
		BlockType3: {
			blockType: BlockType3,
			blocks: [][]bool{
				{f, t, f},
				{t, t, t},
				{f, f, f},
			},
		},
		BlockType4: {
			blockType: BlockType4,
			blocks: [][]bool{
				{t, f, f},
				{t, t, t},
				{f, f, f},
			},
		},
		BlockType5: {
			blockType: BlockType5,
			blocks: [][]bool{
				{f, f, t},
				{t, t, t},
				{f, f, f},
			},
		},
		BlockType6: {
			blockType: BlockType6,
			blocks: [][]bool{
				{f, t, t},
				{t, t, f},
				{f, f, f},
			},
		},
		BlockType7: {
			blockType: BlockType7,
			blocks: [][]bool{
				{t, t, f},
				{f, t, t},
				{f, f, f},
			},
		},
	}
}

func RandomPiece() *Piece {
	n := rand.Intn(int(BlockTypeMax) + 1)
	return Pieces[BlockType(n)]
}

func (p *Piece) InitialPosition() (x, y int) {
	size := len(p.blocks)
	x = (boardWidth - size) / 2
	y = 0

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if p.blocks[i][j] {
				return x, y
			}
		}
		y--
	}

	return x, y
}

func (p *Piece) isBlock(i, j int, angle Angle) bool {
	size := len(p.blocks)
	i2, j2 := i, j
	switch angle {
	case Angle0:
	case Angle90:
		i2 = size - 1 - j
		j2 = i
	case Angle180:
		i2 = size - 1 - i
		j2 = size - 1 - j
	case Angle270:
		i2 = j
		j2 = size - 1 - i
	}
	return p.blocks[i2][j2]
}

func (p *Piece) Draw(screen *ebiten.Image, x, y int, angle Angle) {
	for i := range p.blocks {
		for j := range p.blocks[i] {
			if p.isBlock(i, j, angle) {
				drawBlock(screen, p.blockType, x+j*blockWidth, y+i*blockHeight)
			}
		}
	}
}

func (p *Piece) DrawCenter(screen *ebiten.Image, x, y int, angle Angle) {
	size := len(p.blocks)

	firstY := size - 1
	lastY := 0
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if p.isBlock(i, j, angle) && i < firstY {
				firstY = i
			}
			if p.isBlock(i, j, angle) && i > lastY {
				lastY = i
			}
		}

	}

	firstY *= blockHeight
	lastY = (lastY + 1) * blockHeight

	midX := size * blockWidth / 2
	midY := (lastY-firstY)/2 + firstY

	x -= midX
	y -= midY

	p.Draw(screen, x, y, angle)
}
