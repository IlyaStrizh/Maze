package model_test

import (
	"Maze/internal/model"
	"image/color"
	"reflect"
	"testing"

	colorT "github.com/fatih/color"
)

const imageSize = 500

func TestLoadSaveCave(t *testing.T) {
	cave := model.NewCave()
	result := model.NewCave()

	cave.Load("../../txt/cave_10.txt")
	cave.Save("Load&SaveCave.txt")
	result.Load("Load&SaveCave.txt")

	if !reflect.DeepEqual(result, cave) {
		errorMsg := colorT.MagentaString(`loadCave: "%v", safeCave: "%v"`, cave, result)
		t.Errorf(errorMsg)
	}

	cave.Clear()
	if !reflect.DeepEqual(true, cave.IsEmpty()) {
		errorMsg := colorT.MagentaString(`result: "false", need: "true"`)
		t.Errorf(errorMsg)
	}
}

func TestGenerateUpdateCave(t *testing.T) {
	rows, cols := 10, 10
	cave := model.NewCave()

	cave.GenerateCave(rows, cols, 100, 50)
	if !reflect.DeepEqual(rows, cave.GetRows()) {
		errorMsg := colorT.MagentaString(`rows: "%v", need: "%v"`, cave.GetRows(), rows)
		t.Errorf(errorMsg)
	}

	if !reflect.DeepEqual(rows, cave.GetCols()) {
		errorMsg := colorT.MagentaString(`cols: "%v", need: "%v"`, cave.GetCols(), cols)
		t.Errorf(errorMsg)
	}

	cave.Load("../../txt/cave_10.txt")
	for {
		res, _ := cave.UpdateCave(4, 3)
		if !res {
			break
		}
	}
	result := model.NewCave()
	result.Load("../../txt/result.txt")

	if !reflect.DeepEqual(cave, result) {
		errorMsg := colorT.MagentaString(`result: "%v", need: "%v"`, cave, result)
		t.Errorf(errorMsg)
	}

	screen := model.CreatePixelMatrix(imageSize, imageSize, 255, 255, 255)
	cave.ScaleToPixelMatrix(screen, imageSize, imageSize, color.RGBA{0, 0, 0, 255})

	model.CreateImage(screen, "testCave.png")
}

func TestErrorsCave(t *testing.T) {
	cave := model.NewCave()
	err := cave.Load(".txt")

	if !reflect.DeepEqual(err.Error(), "open .txt: no such file or directory") {
		errorMsg := colorT.MagentaString(`result: "%v", need: "open .txt: no such file or directory"`, err)
		t.Errorf(errorMsg)
	}

	cave.Load("../../txt/cave_10.txt")
	cave.Clear()
	err = cave.Save("errCave.txt")

	if !reflect.DeepEqual(err.Error(), "the cave is empty") {
		errorMsg := colorT.MagentaString(`result: "%v", need: "the cave is empty"`, err)
		t.Errorf(errorMsg)
	}

	err = cave.GenerateCave(-3, 3, 100, 50)
	if !reflect.DeepEqual(err.Error(), "incorrect values of data") {
		errorMsg := colorT.MagentaString(`result: "%v", need: "incorrect values of data"`, err)
		t.Errorf(errorMsg)
	}

	_, err = cave.UpdateCave(-5, -3)
	if !reflect.DeepEqual(err.Error(), "the cave is empty") {
		errorMsg := colorT.MagentaString(`result: "%v", need: "the cave is empty"`, err)
		t.Errorf(errorMsg)
	}
}
