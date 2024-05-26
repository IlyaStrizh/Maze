package model

import "image/color"

type MazeCave interface {
	ScaleToPixelMatrix(scaledMaze [][]color.Color,
		width, height int, color color.RGBA)
	Load(path string) error
	Save(path string) error
	IsEmpty() bool
	GetRows() int
	GetCols() int
	Clear()
}
