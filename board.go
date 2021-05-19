package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/husnimun/blocks/assets"
)

const (
	boardWidth      = 10
	boardHeight     = 20
	maxTick         = 60
	maxLandingCount = 60
	boardXZero      = 5
	boardYZero      = 5
	nextX           = 137
	nextY           = 19
)

var imageFrame *ebiten.Image

type Board struct {
	blocks       [boardHeight][boardWidth]BlockType
	piece        *Piece
	nextPiece    *Piece
	pieceX       int
	pieceY       int
	pieceAngle   Angle
	tick         int
	landingCount int
	gameover     bool
}

func init() {
	f, err := assets.Assets.Open("images/frame.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	imageFrame = ebiten.NewImageFromImage(img)
}

func (b *Board) isBlocked(x, y int) bool {
	if x < 0 || x >= boardWidth {
		return true
	}

	if y < 0 {
		return false
	}

	if y >= boardHeight {
		return true
	}

	return b.blocks[y][x] != BlockTypeNone
}

func (b *Board) collides(x, y int, angle Angle) bool {
	size := len(b.piece.blocks)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if b.piece.isBlock(i, j, angle) && b.isBlocked(x+j, y+i) {
				return true
			}
		}
	}
	return false
}

func (b *Board) movePieceRight() {
	if !b.collides(b.pieceX+1, b.pieceY, b.pieceAngle) {
		b.pieceX++
	}
}

func (b *Board) movePieceLeft() {
	if !b.collides(b.pieceX-1, b.pieceY, b.pieceAngle) {
		b.pieceX--
	}
}

func (b *Board) rotatePiece() {
	angle := b.pieceAngle.RotateRight()
	if !b.collides(b.pieceX, b.pieceY, angle) {
		b.pieceAngle = angle
	}
}

func (b *Board) dropPiece() {
	if !b.collides(b.pieceX, b.pieceY+1, b.pieceAngle) {
		b.pieceY++
	}
}

func (b *Board) pieceDropAble() bool {
	return !b.collides(b.pieceX, b.pieceY+1, b.pieceAngle)
}

func (b *Board) absorbPiece() {
	size := len(b.piece.blocks)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if b.piece.isBlock(i, j, b.pieceAngle) {
				b.blocks[b.pieceY+i][b.pieceX+j] = b.piece.blockType
			}
		}
	}
}

func (b *Board) flush() {
	for i := 0; i < boardHeight; i++ {
		if b.flushableLine(i) {
			b.flushLine(i)
		}
	}
}

func (b *Board) flushLine(i int) {
	for i2 := i; i2 > 0; i2-- {
		for j := 0; j < boardWidth; j++ {
			b.blocks[i2][j] = b.blocks[i2-1][j]
		}
	}

	for j := 0; j < boardWidth; j++ {
		b.blocks[0][j] = BlockTypeNone
	}
}

func (b *Board) flushableLine(i int) bool {
	for j := 0; j < boardWidth; j++ {
		if b.blocks[i][j] == BlockTypeNone {
			return false
		}
	}

	return true
}

func (b *Board) initPiece(p *Piece) {
	b.piece = p
	b.pieceX, b.pieceY = p.InitialPosition()
	b.pieceAngle = Angle0
}

func (b *Board) changeToNextPiece() {
	b.initPiece(b.nextPiece)
	b.nextPiece = RandomPiece()
	b.landingCount = 0
	if !b.pieceDropAble() {
		b.gameover = true
	}
}

func (b *Board) Update() {
	if b.gameover {
		return
	}

	if b.piece == nil {
		b.initPiece(RandomPiece())
	}

	if b.nextPiece == nil {
		b.nextPiece = RandomPiece()
	}

	if d := inpututil.KeyPressDuration(ebiten.KeyRight); d == 1 || (d > 10 && d%2 == 0) {
		b.movePieceRight()
	} else if d := inpututil.KeyPressDuration(ebiten.KeyLeft); d == 1 || (d > 10 && d%2 == 0) {
		b.movePieceLeft()
	} else if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		b.rotatePiece()
	} else if d := inpututil.KeyPressDuration(ebiten.KeyDown); d == 1 || (d > 10 && d%2 == 0) {
		b.dropPiece()
	}

	b.tick++
	if b.tick >= maxTick {
		b.tick = 0
		b.dropPiece()
	}

	if !b.pieceDropAble() {
		if d := inpututil.KeyPressDuration(ebiten.KeyDown); d > 0 {
			b.landingCount += 10
		} else {
			b.landingCount++
		}

		if b.landingCount >= maxLandingCount {
			b.absorbPiece()
			b.changeToNextPiece()
			b.flush()
		}
	}
}

func (b *Board) drawFrame(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(imageFrame, op)
}

func (b *Board) drawNext(screen *ebiten.Image) {
	if b.nextPiece != nil {
		b.nextPiece.DrawCenter(screen, nextX, nextY, Angle0)
	}
}

func (b *Board) Draw(screen *ebiten.Image) {
	b.drawFrame(screen)

	b.piece.Draw(screen, boardXZero+b.pieceX*blockWidth, boardYZero+b.pieceY*blockHeight, b.pieceAngle)

	for i := 0; i < boardHeight; i++ {
		for j := 0; j < boardWidth; j++ {
			if b.blocks[i][j] != BlockTypeNone {
				drawBlock(screen, b.blocks[i][j], boardXZero+j*blockWidth, boardYZero+i*blockHeight)
			}
		}
	}

	b.drawNext(screen)
}
