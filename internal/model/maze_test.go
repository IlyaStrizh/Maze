package model_test

import (
	"Maze/internal/model"
	"reflect"
	"testing"

	"image/color"

	colorT "github.com/fatih/color"
)

func TestMaze(t *testing.T) {
	const imageSize = 500

	t.Run("TestLoad&Save", func(t *testing.T) {
		maze := model.NewMaze()
		result := model.NewMaze()

		maze.Load("../../txt/maze_10.txt")
		maze.Save("Load&SaveMaze.txt")
		result.Load("Load&SaveMaze.txt")

		if !reflect.DeepEqual(result, maze) {
			errorMsg := colorT.MagentaString(`loadMaze: "%v", safeMaze: "%v"`, maze, result)
			t.Errorf(errorMsg)
		}
	})

	t.Run("TestGenerate", func(t *testing.T) {
		rows, cols := 33, 33
		maze := model.NewMaze()

		maze.GenerateMaze(cols, rows)
		if !reflect.DeepEqual(cols, maze.GetCols()) {
			errorMsg := colorT.MagentaString(`cols: "%v", maze.GetCols: "%v"`, cols, maze.GetCols())
			t.Errorf(errorMsg)
		}

		if !reflect.DeepEqual(rows, maze.GetRows()) {
			errorMsg := colorT.MagentaString(`rows: "%v", maze.GetRows: "%v"`, rows, maze.GetRows())
			t.Errorf(errorMsg)
		}

		// Проверка отсутствия изолированных областей
		visited := make([][]bool, maze.GetRows())
		for i := range visited {
			visited[i] = make([]bool, maze.GetCols())
		}
		dfs(&maze, 0, 0, visited)

		for row := 0; row < maze.GetRows(); row++ {
			for col := 0; col < maze.GetCols(); col++ {
				if !visited[row][col] {
					errorMsg := colorT.MagentaString("Maze has isolated areas at (%d, %d)", row, col)
					t.Errorf(errorMsg)
				}
			}
		}
	})

	t.Run("FindPath", func(t *testing.T) {
		maze := model.NewMaze()
		maze.Load("../../txt/maze_10.txt")

		start := model.NewPoint(0, 0)
		end := model.NewPoint(9, 9)
		screen := model.CreatePixelMatrix(imageSize, imageSize, 0, 0, 0)
		result := []model.Point{model.NewPoint(0, 0), model.NewPoint(0, 1),
			model.NewPoint(0, 2), model.NewPoint(1, 2), model.NewPoint(2, 2),
			model.NewPoint(2, 1), model.NewPoint(3, 1), model.NewPoint(3, 2),
			model.NewPoint(3, 3), model.NewPoint(4, 3), model.NewPoint(4, 4),
			model.NewPoint(3, 4), model.NewPoint(3, 5), model.NewPoint(2, 5),
			model.NewPoint(2, 4), model.NewPoint(2, 3), model.NewPoint(1, 3),
			model.NewPoint(0, 3), model.NewPoint(0, 4), model.NewPoint(1, 4),
			model.NewPoint(1, 5), model.NewPoint(1, 6), model.NewPoint(1, 7),
			model.NewPoint(2, 7), model.NewPoint(3, 7), model.NewPoint(3, 8),
			model.NewPoint(2, 8), model.NewPoint(1, 8), model.NewPoint(0, 8),
			model.NewPoint(0, 9), model.NewPoint(1, 9), model.NewPoint(2, 9),
			model.NewPoint(3, 9), model.NewPoint(4, 9), model.NewPoint(5, 9),
			model.NewPoint(6, 9), model.NewPoint(6, 8), model.NewPoint(5, 8),
			model.NewPoint(5, 7), model.NewPoint(6, 7), model.NewPoint(6, 6),
			model.NewPoint(6, 5), model.NewPoint(7, 5), model.NewPoint(7, 6),
			model.NewPoint(8, 6), model.NewPoint(9, 6), model.NewPoint(9, 7),
			model.NewPoint(9, 8), model.NewPoint(9, 9)}

		maze.ScaleToPixelMatrix(screen, imageSize, imageSize, color.RGBA{0, 255, 255, 255})
		path := maze.FindPath(start, end)

		if !reflect.DeepEqual(result, path) {
			errorMsg := colorT.MagentaString(`result: "%v", path: "%v"`, result, path)
			t.Errorf(errorMsg)
		}

		maze.DrawPath(screen, path, imageSize, imageSize)
		maze.DrawPoint(screen, 25, 25, imageSize, imageSize, color.RGBA{255, 255, 0, 255})
		maze.DrawPoint(screen, 475, 475, imageSize, imageSize, color.RGBA{255, 0, 255, 255})

		model.CreateImage(screen, "testMaze.png")

		maze.Clear()

		if !reflect.DeepEqual(true, maze.IsEmpty()) {
			errorMsg := colorT.MagentaString(`result: "false", need: "true"`)
			t.Errorf(errorMsg)
		}
	})

	t.Run("QLearning", func(t *testing.T) {
		maze := model.NewMaze()
		qLearning := model.NewQLearningAgent(&maze, model.NewPoint(0, 0))
		err := qLearning.Load("../../txt/qLearningMaze.txt")
		if err != nil {
			t.Errorf(err.Error())
		}

		if qLearning.GetCols() != 10 || qLearning.GetRows() != 10 {
			t.Errorf("GetCols() != 10 || GetRows() != 10")
		}

		if !reflect.DeepEqual(model.NewPoint(9, 9), qLearning.GetEndPoint()) {
			errorMsg := colorT.MagentaString(`endPoint: "%v", want: "%v"`, qLearning.GetEndPoint(), model.NewPoint(9, 9))
			t.Errorf(errorMsg)
		}

		result := []model.Point{
			model.NewPoint(3, 3), model.NewPoint(4, 3), model.NewPoint(4, 4),
			model.NewPoint(3, 4), model.NewPoint(3, 5), model.NewPoint(2, 5),
			model.NewPoint(2, 4), model.NewPoint(2, 3), model.NewPoint(1, 3),
			model.NewPoint(0, 3), model.NewPoint(0, 4), model.NewPoint(1, 4),
			model.NewPoint(1, 5), model.NewPoint(1, 6), model.NewPoint(1, 7),
			model.NewPoint(2, 7), model.NewPoint(3, 7), model.NewPoint(3, 8),
			model.NewPoint(2, 8), model.NewPoint(1, 8), model.NewPoint(0, 8),
			model.NewPoint(0, 9), model.NewPoint(1, 9), model.NewPoint(2, 9),
			model.NewPoint(3, 9), model.NewPoint(4, 9), model.NewPoint(5, 9),
			model.NewPoint(6, 9), model.NewPoint(6, 8), model.NewPoint(5, 8),
			model.NewPoint(5, 7), model.NewPoint(6, 7), model.NewPoint(6, 6),
			model.NewPoint(6, 5), model.NewPoint(7, 5), model.NewPoint(7, 6),
			model.NewPoint(8, 6), model.NewPoint(9, 6), model.NewPoint(9, 7),
			model.NewPoint(9, 8), model.NewPoint(9, 9),
		}

		qLearning.TrainAgent()

		if !qLearning.IsTrained() {
			t.Errorf("is not trained")
		}

		path := qLearning.InteractWithAgent(model.NewPoint(3, 3))
		if !reflect.DeepEqual(result, path) {
			errorMsg := colorT.MagentaString(`result: "%v", path: "%v"`, result, path)
			t.Errorf(errorMsg)
		}

		screen := model.CreatePixelMatrix(imageSize, imageSize, 0, 0, 0)
		qLearning.ScaleToPixelMatrix(screen, imageSize, imageSize, color.RGBA{0, 255, 255, 255})
		maze.DrawPath(screen, path, imageSize, imageSize)
		maze.DrawPoint(screen, 175, 175, imageSize, imageSize, color.RGBA{255, 255, 0, 255})
		model.CreateImage(screen, "qLearning.png")

		qLearning.Clear()

		if !qLearning.IsEmpty() {
			t.Errorf("qLearning.Clear() don`t work")
		}
	})
}

func TestErrorsMaze(t *testing.T) {
	maze := model.NewMaze()
	err := maze.GenerateMaze(-3, 3)

	if !reflect.DeepEqual(err.Error(), "incorrect values cols or rows") {
		errorMsg := colorT.MagentaString(`result: "%v", need: "incorrect values cols or rows"`, err)
		t.Errorf(errorMsg)
	}

	maze.Clear()
	err = maze.Save("errMaze.txt")

	if !reflect.DeepEqual(err.Error(), "the maze is empty") {
		errorMsg := colorT.MagentaString(`result: "%v", need: "the maze is empty"`, err)
		t.Errorf(errorMsg)
	}
}

func dfs(maze *model.Maze, row, col int, visited [][]bool) {
	if row < 0 || row >= maze.GetRows() || col < 0 || col >= maze.GetCols() || visited[row][col] {
		return
	}
	visited[row][col] = true
	if !maze.CheckRight(row, col) {
		dfs(maze, row, col+1, visited)
	}
	if !maze.CheckDown(row, col) {
		dfs(maze, row+1, col, visited)
	}
	if col > 0 && !maze.CheckRight(row, col-1) {
		dfs(maze, row, col-1, visited)
	}
	if row > 0 && !maze.CheckDown(row-1, col) {
		dfs(maze, row-1, col, visited)
	}
}
