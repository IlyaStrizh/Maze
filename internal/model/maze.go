package model

import (
	"bufio"
	"errors"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
)

type MazeCell struct {
	rightBorder bool
	downBorder  bool
}

type Maze struct {
	rows, cols int
	maze       [][]MazeCell
}

type Point struct {
	x, y int
}

func NewMaze() Maze {
	return Maze{}
}

func NewPoint(x, y int) Point {
	return Point{
		x: x,
		y: y,
	}
}

func (p *Point) GetX() int {
	return p.x
}

func (p *Point) GetY() int {
	return p.y
}

func (m *Maze) GenerateMaze(cols, rows int) error {
	if cols <= 0 || rows <= 0 {
		return errors.New("incorrect values cols or rows")
	}

	m.cols = cols
	m.rows = rows
	nextSetID := 1
	rowSets := make([]int, m.cols)
	m.maze = make([][]MazeCell, m.rows)

	for i := 0; i < m.rows; i++ {
		m.maze[i] = make([]MazeCell, m.cols)
		m.addUniqueSet(rowSets, &nextSetID, i)
		m.addRightBorders(rowSets, i)
		m.addDownBorders(rowSets, i)
	}

	return nil
}

func (m *Maze) addRightBorders(rowSets []int, pos int) {
	for i := 0; i < m.cols; i++ {
		if i == m.cols-1 || rowSets[i] == rowSets[i+1] ||
			rand.Intn(math.MaxInt64) > math.MaxInt64>>1 {
			// Установка правой границы
			m.maze[pos][i].rightBorder = true
		} else {
			oldSetID := rowSets[i]
			for j := 0; j < m.cols; j++ {
				// Если ячейка принадлежит предыдущему множеству
				if rowSets[j] == oldSetID {
					// Обновление ID множества для этой ячейки
					rowSets[j] = rowSets[i+1]
				}
			}
		}
	}
	if pos == m.rows-1 {
		for i := 0; i < m.cols-1; i++ {
			// Если текущая ячейка не принадлежит тому же множеству, что и следующая ячейка
			if rowSets[i] != rowSets[i+1] {
				// Сброс правой границы чтобы не было зацикленных областей
				m.maze[pos][i].rightBorder = false
				oldSetID := rowSets[i]
				for j := 0; j < m.cols; j++ {
					if rowSets[j] == oldSetID {
						rowSets[j] = rowSets[i+1]
					}
				}
			}
		}
	}
}

func (m *Maze) addDownBorders(rowSets []int, pos int) {
	// Карта для подсчета размеров множеств
	setSizes := make(map[int]int)

	for _, setID := range rowSets {
		setSizes[setID]++
	}

	for i := 0; i < m.cols; i++ {
		if rand.Intn(math.MaxInt64) > math.MaxInt64>>1 && setSizes[rowSets[i]] > 1 {
			// Установка нижней границы
			m.maze[pos][i].downBorder = true
			// Уменьшение размера множества
			setSizes[rowSets[i]]--
		} else if pos == m.rows-1 {
			m.maze[pos][i].downBorder = true
		}
	}
}

func (m *Maze) addUniqueSet(rowSets []int, nextSetID *int, pos int) {
	// Добавление уникальных множеств для ячеек в заданной строке
	for i := 0; i < m.cols; i++ {
		if pos > 0 && m.CheckDown(pos-1, i) || rowSets[i] == 0 {
			rowSets[i] = *nextSetID
			*nextSetID++
		}
	}
}

func (m *Maze) Load(path string) error {
	var err error
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	if m.rows, m.cols, err = readRowsCols(scanner); err != nil {
		return err
	}

	m.maze = make([][]MazeCell, m.rows)

	// Заполняем вертикальные границы
	m.readVertical(scanner)
	// Пропускаем пустую строку
	scanner.Scan()
	// Заполняем горизонтальные границы
	m.readHorizontal(scanner)

	return nil
}

func readRowsCols(scanner *bufio.Scanner) (rows int, cols int, err error) {
	scanner.Scan()

	data := strings.Fields(scanner.Text())
	if len(data) < 2 {
		return 0, 0, errors.New("incorrect data in file")
	}

	if rows, err = strconv.Atoi(data[0]); err != nil || rows <= 0 {
		return 0, 0, errors.New("incorrect data in file")
	}
	if cols, err = strconv.Atoi(data[1]); err != nil || cols <= 0 {
		return 0, 0, errors.New("incorrect data in file")
	}

	return rows, cols, nil
}

func (m *Maze) readVertical(scanner *bufio.Scanner) {
	for i := 0; i < m.rows; i++ {
		scanner.Scan()
		m.maze[i] = make([]MazeCell, m.cols)
		data := strings.Fields(scanner.Text())
		if len(data) < m.cols {
			continue
		}
		for j := 0; j < m.cols; j++ {
			val, err := strconv.Atoi(string(data[j]))
			if err != nil {
				continue
			}
			if val == 1 {
				m.maze[i][j].rightBorder = true
			}
		}
	}
}

func (m *Maze) readHorizontal(scanner *bufio.Scanner) {
	for i := 0; i < m.rows && scanner.Scan(); i++ {
		data := strings.Fields(scanner.Text())
		if len(data) < m.cols {
			continue
		}
		for j := 0; j < m.cols; j++ {
			val, err := strconv.Atoi(string(data[j]))
			if err != nil {
				continue
			}
			if val == 1 {
				m.maze[i][j].downBorder = true
			}
		}
	}
}

func (m *Maze) Save(path string) error {
	if m.IsEmpty() {
		return errors.New("the maze is empty")
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	str := strconv.Itoa(m.rows) + " " + strconv.Itoa(m.cols) + "\n"
	if _, err = writer.Write([]byte(str)); err != nil {
		return err
	}

	if err = m.writeVertical(writer); err != nil {
		return err
	}
	if err = m.writeHorizontal(writer); err != nil {
		return err
	}

	return nil
}

func (m *Maze) writeVertical(writer *bufio.Writer) error {
	for i := range m.maze {
		builded := &strings.Builder{}
		for j := range m.maze[i] {
			var s string = "0"
			if m.maze[i][j].rightBorder {
				s = "1"
			}
			builded.WriteString(s)
			if j != m.rows-1 {
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

func (m *Maze) writeHorizontal(writer *bufio.Writer) error {
	if _, err := writer.Write([]byte("\n")); err != nil {
		return err
	}

	for i := range m.maze {
		builded := &strings.Builder{}
		for j := range m.maze[i] {
			var s string = "0"
			if m.maze[i][j].downBorder {
				s = "1"
			}
			builded.WriteString(s)
			if j != m.rows-1 {
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

func (m *Maze) Clear() {
	m.cols, m.rows, m.maze = 0, 0, nil
}

func (m *Maze) IsEmpty() bool {
	return m.cols == 0 || m.rows == 0
}

func (m *Maze) CheckRight(i, j int) bool {
	return m.maze[i][j].rightBorder
}

func (m *Maze) CheckDown(i, j int) bool {
	return m.maze[i][j].downBorder
}

func (m *Maze) GetRows() int {
	return m.rows
}

func (m *Maze) GetCols() int {
	return m.cols
}

func (m *Maze) FindPath(start, end Point) []Point {
	start = Point{clamp(start.x, 0, m.rows-1), clamp(start.y, 0, m.cols-1)}
	end = Point{clamp(end.x, 0, m.rows-1), clamp(end.y, 0, m.cols-1)}

	visited := make([][]bool, m.rows)
	for i := range visited {
		visited[i] = make([]bool, m.cols)
	}

	var path []Point
	visited[start.x][start.y] = true
	path = append(path, start)

	for len(path) > 0 {
		current := path[len(path)-1]
		if current == end {
			break
		}
		if !m.IsCanMoveDown(visited, &path, current) &&
			!m.IsCanMoveLeft(visited, &path, current) &&
			!m.IsCanMoveRight(visited, &path, current) &&
			!m.IsCanMoveUp(visited, &path, current) {
			path = path[:len(path)-1]
		}
	}
	return path
}

func (m *Maze) IsCanMoveLeft(visited [][]bool, path *[]Point, current Point) bool {
	result := true
	if current.x-1 < 0 || visited[current.x-1][current.y] || m.maze[current.x-1][current.y].downBorder {
		result = false
	} else {
		*path = append(*path, Point{current.x - 1, current.y})
		visited[current.x-1][current.y] = true
	}
	return result
}

func (m *Maze) IsCanMoveRight(visited [][]bool, path *[]Point, current Point) bool {
	result := true
	if current.x+1 >= m.rows || visited[current.x+1][current.y] || m.maze[current.x][current.y].downBorder {
		result = false
	} else {
		*path = append(*path, Point{current.x + 1, current.y})
		visited[current.x+1][current.y] = true
	}
	return result
}

func (m *Maze) IsCanMoveUp(visited [][]bool, path *[]Point, current Point) bool {
	result := true
	if current.y-1 < 0 || visited[current.x][current.y-1] || m.maze[current.x][current.y-1].rightBorder {
		result = false
	} else {
		*path = append(*path, Point{current.x, current.y - 1})
		visited[current.x][current.y-1] = true
	}
	return result
}

func (m *Maze) IsCanMoveDown(visited [][]bool, path *[]Point, current Point) bool {
	result := true
	if current.y+1 >= m.cols || visited[current.x][current.y+1] || m.maze[current.x][current.y].rightBorder {
		result = false
	} else {
		*path = append(*path, Point{current.x, current.y + 1})
		visited[current.x][current.y+1] = true
	}
	return result
}

func clamp(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func CreatePixelMatrix(rows, cols int, R, G, B uint8) [][]color.Color {
	pixelMatrix := make([][]color.Color, rows)
	for l := range pixelMatrix {
		pixelMatrix[l] = make([]color.Color, cols)
		for rows := range pixelMatrix[l] {
			pixelMatrix[l][rows] = color.RGBA{R, G, B, 255}
		}
	}

	return pixelMatrix
}

func (m *Maze) DrawPoint(scaledMaze [][]color.Color, x, y, width, height int, color color.RGBA) {
	drawLine(scaledMaze, x-2, y+2, x+2, y+2, width, height, true, color)
	drawLine(scaledMaze, x+2, y+2, x+2, y-2, width, height, true, color)
	drawLine(scaledMaze, x+2, y-2, x-3, y-2, width, height, true, color)
	drawLine(scaledMaze, x-3, y-3, x-3, y+3, width, height, true, color)
}

func (m *Maze) DrawPath(scaledMaze [][]color.Color, path []Point, width, height int) {
	if !m.IsEmpty() {
		cellWidth := float64(width) / float64(m.GetCols())
		cellHeight := float64(height) / float64(m.GetRows())

		for i := 1; i < len(path); i++ {
			x0 := int(float64(path[i-1].y)*cellWidth + cellWidth/2)
			y0 := int(float64(path[i-1].x)*cellHeight + cellHeight/2)
			x1 := int(float64(path[i].y)*cellWidth + cellWidth/2)
			y1 := int(float64(path[i].x)*cellHeight + cellHeight/2)
			drawLine(scaledMaze, x0, y0, x1, y1, width, height, true, color.RGBA{0, 255, 0, 255})
		}
	}
}

func (m *Maze) ScaleToPixelMatrix(scaledMaze [][]color.Color, width, height int, color color.RGBA) {
	if !m.IsEmpty() {
		borderX, borderY := width-1, height-1

		cellWidth := float64(width) / float64(m.GetCols())
		cellHeight := float64(height) / float64(m.GetRows())

		wg := new(sync.WaitGroup)

		wg.Add(1)
		go func() {
			defer wg.Done()
			drawLine(scaledMaze, 0, 1, borderX, 1, width, height, true, color)             //color.RGBA{0, 255, 255, 255})
			drawLine(scaledMaze, borderX, 0, borderX, borderY, width, height, true, color) //color.RGBA{0, 255, 255, 255})
			drawLine(scaledMaze, borderX, borderY, 0, borderY, width, height, true, color) //color.RGBA{0, 255, 255, 255})
			drawLine(scaledMaze, 1, borderY, 1, 0, width, height, true, color)             //color.RGBA{0, 255, 255, 255})
		}()
		for i := 0; i < m.GetRows(); i++ {
			for j := 0; j < m.GetCols(); j++ {
				x0 := int(math.Floor(float64(cellWidth) * float64(j)))
				y0 := int(math.Floor(float64(cellHeight) * float64(i)))
				x1 := int(math.Ceil(float64(cellWidth) * float64(j+1)))
				y1 := int(math.Ceil(float64(cellHeight) * float64(i+1)))

				x1 = int(math.Min(float64(x1), float64(width-1)))
				y1 = int(math.Min(float64(y1), float64(height-1)))

				if m.CheckRight(i, j) {
					drawLine(scaledMaze, x1, y0, x1, y1, width, height, true, color) //color.RGBA{0, 255, 255, 255}) // vertical
				}

				if m.CheckDown(i, j) {
					drawLine(scaledMaze, x0, y1, x1, y1, width, height, true, color) //color.RGBA{0, 255, 255, 255}) // horizontal
				}
			}
		}

		wg.Wait()
	}
}

func drawLine(img [][]color.Color, x0, y0, x1, y1, width, height int,
	doubleLine bool, color color.RGBA) {
	// Алгоритм Брезенхэма
	// Проверяем, что начальная и конечная точки лежат в пределах изображения
	if x0 < 0 || x0 >= width || y0 < 0 || y0 >= height ||
		x1 < 0 || x1 >= width || y1 < 0 || y1 >= height {
		return
	}

	// Определяем длину и направление линии
	dx := x1 - x0
	dy := y1 - y0
	stepX := 1
	stepY := 1
	if dx < 0 {
		stepX = -1
		dx = -dx
	}
	if dy < 0 {
		stepY = -1
		dy = -dy
	}

	// Определяем, ось x или ось y является основной
	if dx > dy {
		// Основная ось - ось x
		y := y0
		eps := 0
		for x := x0; x != x1; x += stepX {
			if doubleLine {
				img[y-1][x] = color
			}
			img[y][x] = color
			eps += dy
			if (eps << 1) >= dx {
				y += stepY
				eps -= dx
			}
		}
	} else {
		// Основная ось - ось y
		x := x0
		eps := 0
		for y := y0; y != y1; y += stepY {
			if doubleLine {
				img[y][x-1] = color
			}
			img[y][x] = color
			eps += dx
			if (eps << 1) >= dy {
				x += stepX
				eps -= dy
			}
		}
	}
}

func CreateImage(screen [][]color.Color, name string) error {
	x := len(screen)
	y := len(screen[x-1])

	file, err := os.Create(name)
	if err != nil {
		return err
	}

	img := image.NewRGBA(image.Rect(0, 0, x, y))
	for i := range screen {
		for j := range screen[i] {
			img.Set(j, i, screen[i][j])
		}
	}

	err = png.Encode(file, img)
	if err != nil {
		return err
	}
	return nil
}
