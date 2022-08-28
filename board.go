package main

import (
	"errors"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
)

type Cell struct {
	x0    int
	y0    int
	x1    int
	y1    int
	style tcell.Style
	text  string
}

type Board struct {
	colEmpty     int
	rowEmpty     int
	cells        [][]Cell
	solution     [][]string
	screen       tcell.Screen
}

func NewBoard(s tcell.Screen, labels [][]string, x, y int, filledStyle, emptyStyle tcell.Style) (b Board) {
	b.screen = s
	b.cells = make([][]Cell, 4)
	b.solution = make([][]string, 4)
	for i := 0; i < 4; i++ {
		b.cells[i] = make([]Cell, 4)
		b.solution[i] = make([]string, 4)
		for j := 0; j < 4; j++ {
			b.solution[i][j] = labels[i][j]
			xtmp := x + j*6
			ytmp := y + i*3
			b.cells[i][j] = Cell{
				x0:    xtmp,
				y0:    ytmp,
				x1:    xtmp + 5,
				y1:    ytmp + 2,
				text:  labels[i][j],
				style: filledStyle,
			}
		}
	}
	b.cells[3][3].style = emptyStyle
	b.colEmpty = 3
	b.rowEmpty = 3
	return
}

func (c Cell) Draw(s tcell.Screen) {
	drawBox(s, c.x0, c.y0, c.x1, c.y1, c.style, c.text)
}

func (b Board) Draw() {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			b.cells[i][j].Draw(b.screen)
		}
	}
}

func (b *Board) MoveUp() error {
	if b.rowEmpty < 3 {
		swapCell(b.screen, &b.cells[b.rowEmpty][b.colEmpty], &b.cells[b.rowEmpty+1][b.colEmpty])
		b.rowEmpty++
		return nil
	}
	return errors.New("Movimiento no válido")
}

func (b *Board) MoveDown() error {
	if b.rowEmpty > 0 {
		swapCell(b.screen, &b.cells[b.rowEmpty][b.colEmpty], &b.cells[b.rowEmpty-1][b.colEmpty])
		b.rowEmpty--
		return nil
	}
	return errors.New("Movimiento no válido")
}

func (b *Board) MoveLeft() error {
	if b.colEmpty < 3 {
		swapCell(b.screen, &b.cells[b.rowEmpty][b.colEmpty], &b.cells[b.rowEmpty][b.colEmpty+1])
		b.colEmpty++
		return nil
	}
	return errors.New("Movimiento no válido")
}

func (b *Board) MoveRight() error {
	if b.colEmpty > 0 {
		swapCell(b.screen, &b.cells[b.rowEmpty][b.colEmpty], &b.cells[b.rowEmpty][b.colEmpty-1])
		b.colEmpty--
		return nil
	}
	return errors.New("Movimiento no válido")
}

func (b *Board) Shuffle(nMoves int) {
	for i := 0; i < nMoves; {
		rand.Seed(time.Now().UnixNano())
		var err error
		switch rand.Intn(4) {
		case 0:
			err = b.MoveUp()
		case 1:
			err = b.MoveLeft()
		case 2:
			err = b.MoveDown()
		case 3:
			err = b.MoveRight()
		}
		if err == nil {
			i++
		}
	}
}

func (b Board) Solved() bool {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if b.cells[i][j].text != b.solution[i][j] {
				return false
			}
		}
	}
	return true
}

func swapCell(s tcell.Screen, c1, c2 *Cell) {
	tmpText := c1.text
	tmpStyle := c1.style
	c1.text = c2.text
	c1.style = c2.style
	c2.text = tmpText
	c2.style = tmpStyle
	c1.Draw(s)
	c2.Draw(s)
}

func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range text {
		if r == '\n' {
			row++
			col = x1
			continue
		}
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}

func drawBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	// Background
	for row := y1; row <= y2; row++ {
		for col := x1; col <= x2; col++ {
			s.SetContent(col, row, ' ', nil, style)
		}
	}

	// Borders
	for col := x1; col <= x2; col++ {
		s.SetContent(col, y1, tcell.RuneHLine, nil, style)
		s.SetContent(col, y2, tcell.RuneHLine, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		s.SetContent(x1, row, tcell.RuneVLine, nil, style)
		s.SetContent(x2, row, tcell.RuneVLine, nil, style)
	}

	// Corners if necessary
	if y1 != y2 && x1 != x2 {
		s.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
		s.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
		s.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
		s.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
	}

	drawText(s, x1+2, y1+1, x2-1, y2-1, style, text)
}

func drawMessage(s tcell.Screen, style tcell.Style, msg string) {
	drawBox(s, 30, 2, 75, 13, style, msg)
}
