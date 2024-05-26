package model

import (
	"errors"
	"image/color"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type Action int

const (
	Left Action = iota
	Right
	Up
	Down
)

var ACTIONS = []Action{Left, Right, Up, Down}

const (
	NUM_EPISODES = 1500
	EPSILON      = 0.1
)

type QLearningAgent struct {
	maze      *Maze
	qTable    map[Point]map[Action]float64
	visited   [][]bool
	discount  float64
	alpha     float64
	currState Point
	endPoint  Point
	isTrained bool
}

func NewQLearningAgent(m *Maze, p Point) QLearningAgent {
	return QLearningAgent{
		maze:     m,
		qTable:   make(map[Point]map[Action]float64),
		discount: 1,
		alpha:    0.9,
		endPoint: p,
	}
}

func (q *QLearningAgent) GetEndPoint() Point {
	return q.endPoint
}

func (q *QLearningAgent) initQTable() {
	for x := 0; x < q.maze.GetRows(); x++ {
		for y := 0; y < q.maze.GetCols(); y++ {
			q.qTable[Point{x, y}] = make(map[Action]float64)
			for _, action := range ACTIONS {
				if q.isValidCell(Point{x, y}, action) {
					q.qTable[Point{x, y}][action] = 0
				}
			}
		}
	}
}

func (q *QLearningAgent) TrainAgent() {
	if len(q.qTable) == 0 {
		q.initQTable()
	}

	for episode := 0; episode < NUM_EPISODES; episode++ {
		q.visited = make([][]bool, q.maze.GetRows())
		for i := range q.visited {
			q.visited[i] = make([]bool, q.maze.GetCols())
		}

		startPoint := Point{rand.Intn(q.maze.GetRows() - 1), rand.Intn(q.maze.GetCols() - 1)}
		q.currState = startPoint
		for count := 0; !q.isTerminalCell(q.currState); count++ {
			action := q.epsilonGreedy(q.currState)
			q.evalQFunction(q.currState, action)
			q.visited[q.currState.x][q.currState.y] = true
			q.currState = q.getCellAfterAction(q.currState, action)

			if q.currState == startPoint && count > q.maze.GetRows()*q.maze.GetCols()*100 {
				break
			}
		}
	}
	q.isTrained = true
}

func (q *QLearningAgent) epsilonGreedy(state Point) Action {
	if rand.Float64() < EPSILON {
		var validActions []Action
		for _, action := range ACTIONS {
			if q.isValidCell(state, action) {
				validActions = append(validActions, action)
			}
		}
		if len(validActions) == 0 {
			return Down // или любое другое действие
		}
		return validActions[rand.Intn(len(validActions))]
	} else {
		var bestAction Action
		bestQValue := -1e9
		for action, qValue := range q.qTable[state] {
			if qValue > bestQValue {
				bestAction = action
				bestQValue = qValue
			}
		}
		return bestAction
	}
}

func (q *QLearningAgent) evalQFunction(coord Point, action Action) {
	nextCell := q.getCellAfterAction(coord, action)
	reward := q.getCellValue(coord, action)
	maxQSPrime := -1e9
	for _, action2 := range ACTIONS {
		// проверка в посещенных и валиднось перехода для оценки веса
		if q.isValidCell(nextCell, action2) && !q.visited[nextCell.x][nextCell.y] {
			maxQSPrime = math.Max(maxQSPrime, q.qTable[nextCell][action2])
		}
	}

	q.qTable[coord] = make(map[Action]float64)
	q.qTable[coord][action] += q.alpha * (float64(reward) + q.discount*float64(maxQSPrime) - q.qTable[coord][action])
}

func (q *QLearningAgent) isValidCell(p Point, a Action) bool {
	switch {
	case a == Left && p.x > 0:
		return !q.maze.CheckDown(p.x-1, p.y)
	case a == Right && p.x < q.maze.GetRows():
		return !q.maze.CheckDown(p.x, p.y)
	case a == Up && p.y > 0:
		return !q.maze.CheckRight(p.x, p.y-1)
	case a == Down && p.y < q.maze.GetRows():
		return !q.maze.CheckRight(p.x, p.y)
	default:
		return false
	}
}

func (q *QLearningAgent) isTerminalCell(p Point) bool {
	return p == q.endPoint
}

func (q *QLearningAgent) getCellAfterAction(p Point, a Action) Point {
	switch {
	case a == Left && p.x > 0:
		return Point{p.x - 1, p.y}
	case a == Right && p.x < q.maze.GetRows()-1:
		return Point{p.x + 1, p.y}
	case a == Up && p.y > 0:
		return Point{p.x, p.y - 1}
	case a == Down && p.y < q.maze.GetCols()-1:
		return Point{p.x, p.y + 1}
	default:
		return p
	}
}

func (q *QLearningAgent) getCellValue(p Point, a Action) float64 {
	var reward float64

	if q.isValidCell(p, a) && !q.visited[p.x][p.y] {
		reward = 0.5
	} else {
		reward = -10 // Штраф за недействительные ходы
	}

	return reward
}

func (q *QLearningAgent) InteractWithAgent(startState Point) []Point {
	var result []Point
	q.currState = startState
	result = append(result, q.currState)

	for count := 0; !q.isTerminalCell(q.currState); count++ {
		bestAction := q.getBestAction(q.currState)
		nextState := q.getCellAfterAction(q.currState, bestAction)

		// Обновляем текущее состояние и добавляем его в решение
		q.currState = nextState
		result = append(result, q.currState)

		if count > q.maze.GetRows()*q.maze.GetCols()*100 {
			break
		}
	}

	if !q.isTerminalCell(q.currState) {
		result = []Point{}
	}

	return result
}

func (q *QLearningAgent) IsTrained() bool {
	return q.isTrained
}

func (q *QLearningAgent) getBestAction(state Point) Action {
	var bestAction Action
	bestQValue := -1e9

	for action, qValue := range q.qTable[state] {
		if qValue > bestQValue {
			bestAction = action
			bestQValue = qValue
		}
	}

	return bestAction
}

func (q *QLearningAgent) Load(path string) error {
	err := q.maze.Load(path)
	if err != nil {
		return err
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	lastLine := lines[len(lines)-2]

	data := strings.Fields(lastLine)
	if len(data) < 2 {
		return errors.New("invalid data in file")
	} else {
		x, _ := strconv.Atoi(data[0])
		y, _ := strconv.Atoi(data[1])
		q.endPoint = NewPoint(x, y)
	}
	q.initQTable()

	return nil
}

func (q *QLearningAgent) Clear() {
	q.maze = nil
	q.endPoint = Point{}
}

func (q *QLearningAgent) GetCols() int {
	return q.maze.GetCols()
}

func (q *QLearningAgent) GetRows() int {
	return q.maze.GetRows()
}

func (q *QLearningAgent) IsEmpty() bool {
	return q.maze == nil && q.endPoint == Point{}
}

func (q *QLearningAgent) Save(path string) error {
	return q.maze.Save(path)
}

func (q *QLearningAgent) ScaleToPixelMatrix(scaledMaze [][]color.Color,
	width, height int, color color.RGBA) {
	q.maze.ScaleToPixelMatrix(scaledMaze, width, height, color)
}
