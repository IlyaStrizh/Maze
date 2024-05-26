package model

import (
	"bufio"
	"errors"
	"image/color"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type Cave struct {
	rows, cols int
	cave       [][]bool
}

func NewCave() Cave {
	return Cave{}
}

func (c *Cave) GenerateCave(cols, rows, N, chance int) error {
	if cols <= 0 || rows <= 0 || N <= 0 || chance < 0 {
		return errors.New("incorrect values of data")
	}

	c.cols = cols
	c.rows = rows
	c.cave = make([][]bool, c.rows)

	if N > 1 {
		N -= 1
	}

	for i := 0; i < c.rows; i++ {
		c.cave[i] = make([]bool, c.cols)
		for j := 0; j < c.cols; j++ {
			if rand.Intn(N) < chance {
				c.cave[i][j] = true
			}
		}
	}

	return nil
}

func (c *Cave) UpdateCave(birthLimit, deathLimit int) (bool, error) {
	var res bool

	if c.IsEmpty() || birthLimit < 0 || deathLimit < 0 {
		return res, errors.New("the cave is empty")
	}

	newCave := NewCave()
	newCave.rows = c.GetRows()
	newCave.cols = c.GetCols()
	newCave.cave = make([][]bool, c.GetRows())

	for i := 0; i < c.rows; i++ {
		newCave.cave[i] = make([]bool, c.GetCols())
		for j := 0; j < c.cols; j++ {
			count := c.neighborsCount(i, j)
			switch {
			case c.cave[i][j] && count < deathLimit:
				newCave.cave[i][j] = false
				res = true
			case !c.cave[i][j] && count > birthLimit:
				newCave.cave[i][j] = true
				res = true
			default:
				newCave.cave[i][j] = c.cave[i][j]
			}
		}
	}
	*c = newCave

	return res, nil
}

func (c *Cave) neighborsCount(i, j int) int {
	var res int
	offsets := [8][2]int{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}}

	for _, offset := range offsets {
		row, col := i+offset[0], j+offset[1]
		if row < 0 || row >= c.rows || col < 0 || col >= c.cols ||
			c.cave[row][col] {
			res++
		}
	}

	return res
}

func (c *Cave) Load(path string) error {
	var err error
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	if c.rows, c.cols, err = readRowsCols(scanner); err != nil {
		return err
	}

	c.cave = make([][]bool, c.rows)

	c.readField(scanner)

	return nil
}

func (c *Cave) readField(scanner *bufio.Scanner) {
	for i := 0; i < c.rows; i++ {
		scanner.Scan()
		c.cave[i] = make([]bool, c.cols)
		data := strings.Fields(scanner.Text())
		if len(data) < c.cols {
			continue
		}
		for j := 0; j < c.cols; j++ {
			val, err := strconv.Atoi(string(data[j]))
			if err != nil {
				continue
			}
			if val == 1 {
				c.cave[i][j] = true
			}
		}
	}
}

func (c *Cave) Save(path string) error {
	if c.IsEmpty() {
		return errors.New("the cave is empty")
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	str := strconv.Itoa(c.rows) + " " + strconv.Itoa(c.cols) + "\n"
	if _, err = writer.Write([]byte(str)); err != nil {
		return err
	}

	if err = c.writeField(writer); err != nil {
		return err
	}

	return nil
}

func (c *Cave) writeField(writer *bufio.Writer) error {
	for i := range c.cave {
		builded := &strings.Builder{}
		for j := range c.cave[i] {
			var s string = "0"
			if c.cave[i][j] {
				s = "1"
			}
			builded.WriteString(s)
			if j != c.rows-1 {
				builded.WriteString(" ")
			}
		}
		builded.WriteString("\n")
		if _, err := writer.Write([]byte(builded.String())); err != nil {
			return err
		}
	}

	return nil
}

func (c *Cave) Clear() {
	c.rows, c.cols, c.cave = 0, 0, nil
}

func (c *Cave) IsEmpty() bool {
	return c.rows == 0 || c.cols == 0
}

func (c *Cave) GetRows() int {
	return c.rows
}

func (c *Cave) GetCols() int {
	return c.cols
}

func (c *Cave) ScaleToPixelMatrix(scaledMaze [][]color.Color,
	width, height int, color color.RGBA) {
	if !c.IsEmpty() {
		cellWidth := float64(width) / float64(c.GetCols())
		cellHeight := float64(height) / float64(c.GetRows())

		for i := 0; i < c.GetRows(); i++ {
			for j := 0; j < c.GetCols(); j++ {
				x0 := int(math.Floor(float64(cellWidth) * float64(j)))
				y0 := int(math.Floor(float64(cellHeight) * float64(i)))
				x1 := int(math.Ceil(float64(cellWidth) * float64(j+1)))
				y1 := int(math.Ceil(float64(cellHeight) * float64(i+1)))

				x1 = int(math.Min(float64(x1), float64(width-1)))
				y1 = int(math.Min(float64(y1), float64(height-1)))

				if c.cave[i][j] {
					for k := y0; k < y1; k++ {
						drawLine(scaledMaze, x0, k, x1, k, width, height, false, color)
					}
				}
			}
		}
	}
}
