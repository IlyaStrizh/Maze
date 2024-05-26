package controller

import (
	"image/color"
	"math"

	"Maze/internal/model"
)

type Controller struct {
	maze         model.Maze
	cave         model.Cave
	qLearning    model.QLearningAgent
	canvasWidth  int
	canvasHeight int
}

func NewController(x, y int) Controller {
	return Controller{
		maze:         model.NewMaze(),
		cave:         model.NewCave(),
		canvasWidth:  x,
		canvasHeight: y,
	}
}

func (c *Controller) InitQLearning() {
	c.qLearning = model.NewQLearningAgent(&c.maze, c.qLearning.GetEndPoint())
}

func (c *Controller) IsTrained() bool {
	return c.qLearning.IsTrained()
}

func (c *Controller) GetObjMaze() *model.Maze {
	return &c.maze
}

func (c *Controller) GetObjCave() *model.Cave {
	return &c.cave
}

func (c *Controller) GetObjQLearning() *model.QLearningAgent {
	return &c.qLearning
}

func (c *Controller) Load(obj interface{}, path string) error {
	return obj.(model.MazeCave).Load(path)
}

func (c *Controller) GenerateMaze(cols, rows int) error {
	return c.maze.GenerateMaze(cols, rows)
}

func (c *Controller) GenerateCave(cols, rows, N, chance int) error {
	return c.cave.GenerateCave(cols, rows, N, chance)
}

func (c *Controller) UpdateCave(birthLimit, deathLimit int) (bool, error) {
	return c.cave.UpdateCave(birthLimit, deathLimit)
}

func (c *Controller) Save(obj interface{}, path string) error {
	return obj.(model.MazeCave).Save(path)
}

func (c *Controller) IsEmpty(obj interface{}) bool {
	return obj.(model.MazeCave).IsEmpty()
}

func (c *Controller) Clear(obj interface{}) {
	obj.(model.MazeCave).Clear()
}

func (c *Controller) ScalePoints() (int, int) {
	cellWidth := float64(c.canvasWidth) / float64(c.maze.GetCols())
	cellHeight := float64(c.canvasHeight) / float64(c.maze.GetRows())

	cellH, cellW := 0, 0
	if c.maze.GetRows()%2 == 0 {
		cellH = c.canvasHeight>>1 + int(cellHeight/2)
	} else {
		cellH = c.canvasHeight >> 1
	}

	if c.maze.GetCols()%2 == 0 {
		cellW = c.canvasWidth>>1 + int(cellWidth/2)
	} else {
		cellW = c.canvasWidth >> 1
	}

	return cellW, cellH
}

func (c *Controller) MovePoint(scaledMaze [][]color.Color, x, y *int, orientation,
	vector bool, color color.RGBA) {
	var point *int
	cellHeight := float64(c.canvasHeight) / float64(c.maze.GetRows())

	if vector {
		point = y
	} else {
		point = x
	}

	if orientation {
		tmp := *point - int(math.Ceil(cellHeight))
		switch {
		case tmp < 0:
			*point = c.canvasHeight - int(math.Ceil(cellHeight/2))
		default:
			*point = tmp
		}
	} else {
		tmp := *point + int(math.Ceil(cellHeight))
		switch {
		case tmp > c.canvasHeight:
			*point = int(math.Ceil(cellHeight / 2))
		default:
			*point = tmp
		}
	}

	c.DrawPoint(scaledMaze, *x, *y, color)
}

func (c *Controller) DrawPoint(scaledMaze [][]color.Color, x, y int, color color.RGBA) {
	c.maze.DrawPoint(scaledMaze, x, y, c.canvasWidth, c.canvasHeight, color)
}

func (c *Controller) FindPath(scaledMaze [][]color.Color, x, y, x1, y1 int) {
	if !c.maze.IsEmpty() {
		cellHeight := float64(c.canvasHeight) / float64(c.maze.GetRows())

		start := model.NewPoint(int(float64(y)/cellHeight), int(float64(x)/cellHeight))
		end := model.NewPoint(int(float64(y1)/cellHeight), int(float64(x1)/cellHeight))

		c.maze.DrawPath(scaledMaze, c.maze.FindPath(start, end),
			c.canvasWidth, c.canvasHeight)
	}
}

func (c *Controller) CreatePixelMatrix(pixelMatrix *[][]color.Color) {
	*pixelMatrix = model.CreatePixelMatrix(c.canvasWidth, c.canvasHeight, 0, 10, 23) //color.Transparent
}

func (c *Controller) ScaleToPixelMatrix(scaledMaze [][]color.Color, obj interface{}) {
	if obj == c.GetObjMaze() || obj == c.GetObjQLearning() {
		obj.(model.MazeCave).ScaleToPixelMatrix(scaledMaze, c.canvasWidth, c.canvasHeight, color.RGBA{0, 255, 255, 250})
	} else {
		obj.(model.MazeCave).ScaleToPixelMatrix(scaledMaze, c.canvasWidth, c.canvasHeight, color.RGBA{90, 0, 240, 250})
	}
}

func (c *Controller) TrainAgent() {
	c.qLearning.TrainAgent()
}

func (c *Controller) InteractWithAgent(scaledMaze [][]color.Color, x, y int) {
	cellHeight := float64(c.canvasHeight) / float64(c.maze.GetRows())
	start := model.NewPoint(int(float64(y)/cellHeight), int(float64(x)/cellHeight))

	c.maze.DrawPath(scaledMaze, c.qLearning.InteractWithAgent(start), c.canvasWidth, c.canvasHeight)
}
