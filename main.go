package main

import (
	"fmt"
	"log"
	"unicode"

	"os"

	"github.com/gdamore/tcell/v2"
)

func main() {
	screen()
}

func screen() {
	defStyle := tcell.StyleDefault.Foreground(tcell.ColorReset).Background(tcell.ColorReset)
	boxStyle := tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorBlue)
	okStyle := tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorGreen)
	errStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorRed)

	shuffleMoves := 0
	userMoves := 0

	// Initialize screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	s.SetStyle(defStyle)
	s.Clear()

	labels := [][]string{
		{"1", "2", "3", "4"},
		{"5", "6", "7", "8"},
		{"9", "10", "11", "12"},
		{"13", "14", "15", " "},
	}

	board := NewBoard(s, labels, 4, 2, boxStyle, defStyle)

	message := "Este programa permite practicar el juego del 15.\n\n"
	message += "Puedes pulsar la tecla <Esc> o <Ctrl+C> para salir del juego en\ncualquier momento.\n\n"
	message += fmt.Sprintf("Ingresa la cantidad de movimientos a realizar para mezclar el\ntablero: %d", shuffleMoves)

	drawBox(s, 4, 2, 75, 13, boxStyle, message)

	quit := func() {
		s.Fini()
		os.Exit(0)
	}

	enterPress := false
	for !enterPress {
		s.Show()
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				quit()
			} else if ev.Key() == tcell.KeyEnter {
				enterPress = true
			} else if unicode.IsDigit(ev.Rune()) {
				shuffleMoves = shuffleMoves*10 + int(ev.Rune()) - 48
			} else if ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 {
				shuffleMoves /= 10
			}
			message = "Este programa permite practicar el juego del 15.\n\n"
			message += "Puedes pulsar la tecla <Esc> o <Ctrl+C> para salir del juego en\ncualquier momento\n\n"
			message += fmt.Sprintf("Ingresa la cantidad de movimientos a realizar para mezclar el\ntablero: %d", shuffleMoves)
			drawBox(s, 4, 2, 75, 13, boxStyle, message)
		}
	}

	s.Clear()
	board.Shuffle(shuffleMoves)
	board.Draw()
	message = fmt.Sprintf("\nMovimientos de mezcla: %d", shuffleMoves)
	info(s, message, nil, boxStyle, errStyle)

	checkBoard := func() {
		if board.Solved() {
			message = fmt.Sprintf("\nMovimientos de mezcla: %d\n\n", shuffleMoves)
			message += fmt.Sprintf("¡Haz ganado!\nRealizaste: %d movimientos\n\n", userMoves)
			message += "Presiona cualquier tecla para salir."
			info(s, message, nil, okStyle, errStyle)
		}
	}
	checkBoard()

	for {
		s.Show()
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			message = fmt.Sprintf("\nMovimientos de mezcla: %d\n\n", shuffleMoves)
			var err error
			if board.Solved() {
				quit()
			} else if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				quit()
			} else if ev.Key() == tcell.KeyCtrlL {
				s.Sync()
			} else if ev.Key() == tcell.KeyUp || ev.Rune() == 'k' {
				err = board.MoveUp()
				message += "Movimiento hacia arriba"
			} else if ev.Key() == tcell.KeyDown || ev.Rune() == 'j' {
				err = board.MoveDown()
				message += "Movimiento hacia abajo"
			} else if ev.Key() == tcell.KeyLeft || ev.Rune() == 'h' {
				err = board.MoveLeft()
				message += "Movimiento hacia la izquierda"
			} else if ev.Key() == tcell.KeyRight || ev.Rune() == 'l' {
				err = board.MoveRight()
				message += "Movimiento hacia la derecha"
			}
			if err == nil {
				userMoves++
			}
			message += fmt.Sprintf("\n\nMovimientos realizados: %d", userMoves)
			info(s, message, err, boxStyle, errStyle)
		}
		checkBoard()
	}
}

func info(s tcell.Screen, msg string, err error, okStyle, errorStyle tcell.Style) {
	message := fmt.Sprintln("Utilize las flechas del teclado o h,j,k,l para mover las celdas del tablero.")
	if err != nil {
		drawMessage(s, errorStyle, message+"\n"+err.Error())
	} else {
		message += msg
		drawMessage(s, okStyle, message)
	}
}
