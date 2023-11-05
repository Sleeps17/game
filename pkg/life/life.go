package life

import (
	"errors"
	"math/rand"
)

type World struct {
	Height int
	Wide   int
	Cells  [][]bool
}

var (
	x_move = []int{-1, -1, -1, 0, 0, 1, 1, 1}
	y_move = []int{0, 1, -1, 1, -1, 0, 1, -1}
)

func NewWorld(height, wide int) (*World, error) {

	if height <= 0 || wide <= 0 {
		return nil, errors.New("height and width cannot be less then 0")
	}

	cells := make([][]bool, height)
	for i := range cells {
		cells[i] = make([]bool, wide)
	}

	return &World{
		Height: height,
		Wide:   wide,
		Cells:  cells,
	}, nil
}

func (w *World) Neighbours(x, y int) int {
	cnt := 0

	for i := range x_move {
		if x+x_move[i] > 0 && x+x_move[i] < w.Height && y+y_move[i] > 0 && y+y_move[i] < w.Wide && w.Cells[x+x_move[i]][y+y_move[i]] {
			cnt++
		}
	}

	return cnt
}

func (w *World) Next(x, y int) bool {
	n := w.Neighbours(x, y)

	alive := w.Cells[x][y]
	if (n == 2 || n == 3) && alive {
		return true
	}

	if !alive && n == 3 {
		return true
	}

	return false
}

func NextState(oldWorld, newWorld *World) {
	for i := 0; i < oldWorld.Height; i++ {
		for j := 0; j < oldWorld.Wide; j++ {
			newWorld.Cells[i][j] = oldWorld.Next(i, j)
		}
	}
}

func (w *World) Seed(fill int) {
	for _, row := range w.Cells {
		for i := range row {
			if rand.Intn(10) < fill/10 {
				row[i] = true
			}
		}
	}
}

func (w *World) String() string {

	brownSquare := "\xF0\x9F\x9F\xAB"
	greenSquare := "\xF0\x9F\x9F\xA9"

	ans := ""

	for i := range w.Cells {
		for j := range w.Cells[i] {

			if w.Cells[i][j] {
				ans += greenSquare
			} else {
				ans += brownSquare
			}
		}
		ans += "\n"
	}

	return ans
}
