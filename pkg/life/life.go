package life

import (
	"errors"
	"math/rand"
)

const (
	BrownSquare = "\xF0\x9F\x9F\xAB"
	GreenSquare = "\xF0\x9F\x9F\xA9"
)

type World struct {
	Height int
	Wide   int
	Fill   int
	Cells  [][]bool
}

var (
	x_move = []int{-1, -1, -1, 0, 0, 1, 1, 1}
	y_move = []int{0, 1, -1, 1, -1, 0, 1, -1}
)

func NewWorld(height, wide, fill int) (*World, error) {

	if height <= 0 || wide <= 0 || fill <= 0 {
		return nil, errors.New("height, width and fill cannot be less then 0")
	}

	cells := make([][]bool, height)
	for i := range cells {
		cells[i] = make([]bool, wide)
	}

	return &World{
		Height: height,
		Wide:   wide,
		Fill:   fill,
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

	f := float64(0)
	for i := 0; i < oldWorld.Height; i++ {
		for j := 0; j < oldWorld.Wide; j++ {
			if newWorld.Cells[i][j] {
				f += 1.0
			}
		}
	}

	f /= float64(oldWorld.Height) * float64(oldWorld.Wide)
	f *= 100
	newWorld.Fill = int(f)
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

	ans := ""

	for i := range w.Cells {
		for j := range w.Cells[i] {

			if w.Cells[i][j] {
				ans += GreenSquare
			} else {
				ans += BrownSquare
			}
		}
		ans += "\n"
	}

	return ans
}
