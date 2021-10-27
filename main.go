package main

import (
	"fmt"
	"os"
	"time"
)

const (
	rows                    = 25
	columns                 = 25
	escape                  = 27
	cellWidth               = 3
	cellHeight              = 3
	setWindowSizePattern    = "%c[8;60;160t"
	stepUpCursorPattern     = "%c[%c[1A"
	clearCurrentLinePattern = "%c[2K"
)

var Out = os.Stdout

func main() {
	// set appropriate size of terminal window
	_, _ = fmt.Fprintf(Out, setWindowSizePattern, escape)

	startGen := startGeneration()
	renderUI(createUI(startGen))

	for {
		nextGen := [columns][rows]int{}
		restart := false

		for i := 1; i < columns-1; i++ {
			for j := 1; j < rows-1; j++ {
				neighbours := make([][]int, 0)
				p := startGen[i-1 : i+2]
				for k := range p {
					neighbours = append(neighbours, p[k][j-1:j+2])
				}

				neighbourState := handleNeighbours(neighbours)
				// check if we are close to some (horizontal or vertical) border
				// and if yes start game from initial point
				if neighbourState == 1 && (i == (columns-2) || j == (rows-2)) {
					restart = true
				}
				nextGen[i][j] = neighbourState
			}
		}

		if restart {
			startGen = startGeneration()
		} else {
			startGen = nextGen
		}

		renderUI(destroyGrid() + createUI(nextGen))

		// small delay between UI rendering
		time.Sleep(100 * time.Millisecond)
	}
}

// start generation from initial point
func startGeneration() [columns][rows]int {
	pattern := [columns][rows]int{}
	pattern[1][2] = 1
	pattern[2][3] = 1
	pattern[3][1] = 1
	pattern[3][2] = 1
	pattern[3][3] = 1
	return pattern
}

// renderUI - write formatted string to Stdout to emulate UI
func renderUI(ui string) {
	_, _ = fmt.Fprint(Out, ui)
}

// createUI - creates formatted string that represents UI to be
// rendered for current generation
func createUI(grid [columns][rows]int) string {
	horizontal := createLine()
	s := horizontal

	for i := 0; i < rows; i++ {
		s = s + createCell(grid[i]) + horizontal
	}

	return s
}

// createLine - creates horizontal line for grid
func createLine() string {
	l := "+"
	pattern := "----+"
	i := 0
	for i < rows {
		l += pattern
		i++
	}
	l += "\n"
	return l
}

// createCell - creates vertical lines for grid along with substituting live cells with "><" sign.
// This func accepts array of integers. Indexes of elements correspond with cells numbers and int value could be 0 or 1.
// 0 - is a dead cell, 1 - alive cell.
func createCell(row [columns]int) string {
	c := []rune{'|'}
	for _, v := range row {
		if v == 1 {
			c = append(c, []rune{' ', '>', '<', ' ', '|'}...)
			continue
		}
		c = append(c, []rune{' ', ' ', ' ', ' ', '|'}...)
	}

	return string(c) + "\n"
}

// destroyGrid - clears previous grid letting next grid to be drawn in the same place.
func destroyGrid() string {
	i := 0
	// move cursor 1 line up
	stepUp := fmt.Sprintf(stepUpCursorPattern, escape, escape)
	// clear current line
	clearLine := fmt.Sprintf(clearCurrentLinePattern, escape)
	s := stepUp + clearLine
	for i < 2*rows {
		s = s + stepUp + clearLine
		i++
	}

	return s
}

// handleNeighbours - returns status of the middle cell depending on neighbours state.
// It processes sub-grid with dimension 3x3
// If alive returns 1 otherwise 0.
func handleNeighbours(neighbours [][]int) int {
	aliveNeighbours := 0
	for i := 0; i < cellWidth; i++ {
		for j := 0; j < cellHeight; j++ {
			if i == 1 && j == 1 {
				continue
			}
			if neighbours[i][j] == 1 {
				aliveNeighbours++
				continue
			}
		}
	}

	// neighbour with [1][1] indexes is our target cell we are trying to get state for
	if neighbours[1][1] == 1 && (aliveNeighbours == 2 || aliveNeighbours == 3) {
		return 1
	}
	if neighbours[1][1] == 0 && aliveNeighbours == 3 {
		return 1
	}
	return 0
}
