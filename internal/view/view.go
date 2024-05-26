package view

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"Maze/internal/controller"
)

const (
	// Размеры окна отрисовки в пикселях
	canvasWidth  = 500
	canvasHeight = 500
)

type Points struct {
	endX          int
	endY          int
	startX        int
	startY        int
	pointX        int
	pointY        int
	setWalkPoints bool
}

type View struct {
	mainWin          fyne.Window
	hBox             *fyne.Container
	pixelMatrix      [][]color.Color
	controller       controller.Controller
	walkPoints       Points
	colsGenerateMaze int
	rowsGenerateMaze int
	colsGenerateCave int
	rowsGenerateCave int
	n                int
	chance           int
	birth            int
	death            int
	milliseconds     int
	qLearningFlag    bool
}

func NewView(myApp fyne.App) *View {
	menu := &View{
		mainWin: myApp.NewWindow("Maze"),
		controller: controller.NewController(
			canvasWidth,
			canvasHeight,
		),
		colsGenerateMaze: 10,
		rowsGenerateMaze: 10,
		colsGenerateCave: 10,
		rowsGenerateCave: 10,
		n:                100,
		chance:           50,
		birth:            4,
		death:            3,
		milliseconds:     500,
	}
	menu.controller.CreatePixelMatrix(&menu.pixelMatrix)
	menu.createBoxMaze()

	menu.mainWin.CenterOnScreen()
	menu.mainWin.SetContent(menu.hBox)
	menu.mainWin.Resize(fyne.NewSize(720, 500))
	menu.mainWin.SetFixedSize(true)
	menu.mainWin.Show()

	return menu
}

func (v *View) createBoxMaze() {
	content := v.createImage()
	content.SetMinSize(fyne.NewSize(500, 500))

	v.hBox = container.NewHBox(
		content,
		v.createVBoxMaze(),
	)
}

func (v *View) createBoxCave() {
	content := v.createImage()
	content.SetMinSize(fyne.NewSize(500, 500))

	v.hBox = container.NewHBox(
		content,
		v.createVBoxCave(),
	)
}

func (v *View) createHBox1() *fyne.Container {

	mazeButton := widget.NewButton("Maze", func() { v.mazeButton() })
	caveButton := widget.NewButton("Cave", func() { v.caveButton() })

	return container.NewGridWithColumns(
		2,
		mazeButton,
		caveButton,
	)
}

func (v *View) createHBox3Maze() *fyne.Container {
	upButton := widget.NewButtonWithIcon("", theme.MoveUpIcon(), func() {
		v.moveButton(true, true)
	})

	return container.NewGridWithColumns(
		3,
		createEmptyLabel(),
		upButton,
		createEmptyLabel(),
	)
}

func (v *View) createHBox4Maze() *fyne.Container {
	leftButton := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		v.moveButton(true, false)
	})
	okButton := widget.NewButton("OK", func() { v.okButton() })
	rightButton := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
		v.moveButton(false, false)
	})

	return container.NewGridWithColumns(
		3,
		leftButton,
		okButton,
		rightButton,
	)

}

func (v *View) createHbox5Maze() *fyne.Container {
	downButton := widget.NewButtonWithIcon("", theme.MoveDownIcon(), func() {
		v.moveButton(false, true)
	})

	return container.NewGridWithColumns(
		3,
		createEmptyLabel(),
		downButton,
		createEmptyLabel(),
	)

}

func (v *View) createVBoxMaze() *fyne.Container {
	loadMazeButton := widget.NewButton("Load maze", func() {
		v.loadButton(v.controller.GetObjMaze())
	})
	saveMazeButton := widget.NewButton("Save maze", func() {
		v.saveButton(v.controller.GetObjMaze())
	})
	clearMazeButton := widget.NewButton("Clear maze", func() {
		v.clearMazeButton(v.controller.GetObjMaze())
	})
	widthLabel, widthEntry := widget.NewLabel("width:"), widget.NewEntry()
	heightLabel, heightEntry := widget.NewLabel("height:"), widget.NewEntry()
	generateMazeButton := widget.NewButton("Generate maze", func() {
		v.generateMazeButton(widthEntry, heightEntry)
	})
	qLearningButton := widget.NewButton("QLearning", v.qLearningButton)

	widthEntry.SetText(strconv.Itoa(v.colsGenerateMaze))
	heightEntry.SetText(strconv.Itoa(v.rowsGenerateMaze))

	hBox2 := container.NewHBox(
		widthLabel,
		widthEntry,
		createDelim(),
		heightLabel,
		heightEntry,
	)

	findPathButton := widget.NewButton("Find path", v.findPathButton)

	return container.NewGridWithRows(
		14,
		v.createHBox1(),
		createEmptyLabel(),
		loadMazeButton,
		clearMazeButton,
		saveMazeButton,
		createEmptyLabel(),
		hBox2,
		generateMazeButton,
		createEmptyLabel(),
		v.createHBox3Maze(),
		v.createHBox4Maze(),
		v.createHbox5Maze(),
		findPathButton,
		qLearningButton,
	)
}

func (v *View) createVBoxCave() *fyne.Container {
	loadCaveButton := widget.NewButton("Load cave", func() {
		v.loadButton(v.controller.GetObjCave())
	})
	saveCaveButton := widget.NewButton("Save cave", func() {
		v.saveButton(v.controller.GetObjCave())
	})
	clearCaveButton := widget.NewButton("Clear cave", func() {
		v.clearCaveButton(v.controller.GetObjCave())
	})
	widthLabel, widthEntry := widget.NewLabel("width:"), widget.NewEntry()
	heightLabel, heightEntry := widget.NewLabel("height:"), widget.NewEntry()
	nLabel, nEntry := widget.NewLabel("                  N:"), widget.NewEntry()
	chanceLabel, chanceEntry := widget.NewLabel("        chance:"), widget.NewEntry()
	generateCaveButton := widget.NewButton("Generate cave", func() {
		v.generateCaveButton(widthEntry, heightEntry, nEntry, chanceEntry)
	})
	birthLabel, birthEntry := widget.NewLabel("birth:"), widget.NewEntry()
	deathhLabel, deathEntry := widget.NewLabel("death:"), widget.NewEntry()
	updateButton := widget.NewButton("Update cave", func() {
		v.updateCaveButton(birthEntry, deathEntry)
	})
	millisecondsLabel, millisecondsEntry := widget.NewLabel("milliseconds:"), widget.NewEntry()
	autoUpdateButton := widget.NewButton("Auto update", func() {
		v.autoUpdateButton(millisecondsEntry)
	})

	nEntry.SetText(strconv.Itoa(v.n))
	chanceEntry.SetText(strconv.Itoa(v.chance))
	widthEntry.SetText(strconv.Itoa(v.colsGenerateCave))
	heightEntry.SetText(strconv.Itoa(v.rowsGenerateCave))
	birthEntry.SetText(strconv.Itoa(v.birth))
	deathEntry.SetText(strconv.Itoa(v.death))
	millisecondsEntry.SetText(strconv.Itoa(v.milliseconds))

	hBox2 := container.NewHBox(
		nLabel,
		nEntry,
	)

	hBox3 := container.NewHBox(
		chanceLabel,
		chanceEntry,
	)

	hBox4 := container.NewHBox(
		widthLabel,
		widthEntry,
		createDelim(),
		heightLabel,
		heightEntry,
	)

	hBox5 := container.NewHBox(
		birthLabel,
		birthEntry,
		createDelim(),
		deathhLabel,
		deathEntry,
	)

	hBox6 := container.NewHBox(
		millisecondsLabel,
		millisecondsEntry,
	)

	return container.NewGridWithRows(
		14,
		v.createHBox1(),
		createEmptyLabel(),
		loadCaveButton,
		clearCaveButton,
		saveCaveButton,
		hBox2,
		hBox3,
		hBox4,
		generateCaveButton,
		createEmptyLabel(),
		hBox5,
		updateButton,
		hBox6,
		autoUpdateButton,
	)
}

func createEmptyLabel() *widget.Label {
	return widget.NewLabel("")
}

func createDelim() *widget.Separator {
	return widget.NewSeparator()
}

func (v *View) createImage() *canvas.Image {
	img := image.NewRGBA(image.Rect(0, 0, canvasWidth, canvasHeight))
	for y := 0; y < canvasWidth; y++ {
		for x := 0; x < canvasHeight; x++ {
			img.Set(x, y, v.pixelMatrix[y][x])
		}
	}

	return canvas.NewImageFromImage(img)
}

func (v *View) loadButton(obj interface{}) {
	dialog.NewFileOpen(func(writer fyne.URIReadCloser, err error) {
		if err != nil {
			v.errorWindow("FAILED TO OPEN FILE", err)
			return
		}
		if writer == nil {
			return
		}
		v.controller.CreatePixelMatrix(&v.pixelMatrix)
		if err = v.controller.Load(obj, writer.URI().Path()); err == nil {
			v.controller.ScaleToPixelMatrix(v.pixelMatrix, obj)
		}

		if obj == v.controller.GetObjMaze() || obj == v.controller.GetObjQLearning() {
			v.createBoxMaze()
		} else {
			v.createBoxCave()
		}
		v.walkPoints = Points{}

		v.mainWin.SetContent(v.hBox)
		writer.Close()
	}, v.mainWin).Show()
}

func (v *View) saveButton(obj interface{}) {
	if !v.controller.IsEmpty(obj) {
		dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil {
				v.errorWindow("FAILED TO CREATE FILE", err)
				return
			}
			if writer == nil {
				return
			}
			if err = v.controller.Save(obj, writer.URI().Path()); err != nil {
				v.errorWindow("FAILED TO SAVING", err)
			}

			writer.Close()
		}, v.mainWin).Show()
	} else {
		v.errorWindow("FAILED TO SAVING", errors.New("the object is empty"))
	}
}

func (v *View) errorWindow(text string, err error) {
	errWindow := fyne.CurrentApp().NewWindow(" ")
	content := widget.NewLabel(fmt.Sprintf("%s:\n\n \"%s\"\n", text, err))

	container := container.New(layout.NewCenterLayout(), content)
	errWindow.SetContent(container)
	errWindow.Resize(fyne.NewSize(250, 100))
	errWindow.CenterOnScreen()
	errWindow.Show()
}

func (v *View) generateMazeButton(widthEntry, heightEntry *widget.Entry) {
	var err error
	v.qLearningFlag = false
	v.colsGenerateMaze, err = strconv.Atoi(widthEntry.Text)
	if err != nil || v.colsGenerateMaze <= 0 {
		v.colsGenerateMaze, v.rowsGenerateMaze = 10, 10
	} else {
		v.rowsGenerateMaze, err = strconv.Atoi(heightEntry.Text)
		if err != nil || v.rowsGenerateMaze <= 0 {
			v.colsGenerateMaze, v.rowsGenerateMaze = 10, 10
		}
	}

	if v.colsGenerateMaze > 50 || v.rowsGenerateMaze > 50 {
		v.colsGenerateMaze, v.rowsGenerateMaze = 50, 50
	}

	v.walkPoints = Points{}
	v.controller.CreatePixelMatrix(&v.pixelMatrix)
	v.controller.GenerateMaze(v.colsGenerateMaze, v.rowsGenerateMaze)
	v.controller.ScaleToPixelMatrix(v.pixelMatrix, v.controller.GetObjMaze())
	v.createBoxMaze()

	v.mainWin.SetContent(v.hBox)
}

func (v *View) generateCaveButton(widthEntry, heightEntry, nEntry, chanceEntry *widget.Entry) {
	var err error
	v.colsGenerateCave, err = strconv.Atoi(widthEntry.Text)
	if err != nil || v.colsGenerateCave <= 0 {
		v.colsGenerateCave, v.rowsGenerateCave = 10, 10
	} else {
		v.rowsGenerateCave, err = strconv.Atoi(heightEntry.Text)
		if err != nil || v.rowsGenerateCave <= 0 {
			v.colsGenerateCave, v.rowsGenerateCave = 10, 10
		}
	}

	if v.colsGenerateCave > 50 || v.rowsGenerateCave > 50 {
		v.colsGenerateCave, v.rowsGenerateCave = 50, 50
	}

	v.n, err = strconv.Atoi(nEntry.Text)
	if err != nil || v.n <= 0 {
		v.n = 100
	}

	v.chance, err = strconv.Atoi(chanceEntry.Text)
	if err != nil || v.chance < 0 || v.chance > v.n {
		v.chance = v.n / 2
	}

	v.controller.CreatePixelMatrix(&v.pixelMatrix)
	v.controller.GenerateCave(v.colsGenerateCave, v.rowsGenerateCave, v.n, v.chance)
	v.controller.ScaleToPixelMatrix(v.pixelMatrix, v.controller.GetObjCave())
	v.createBoxCave()

	v.mainWin.SetContent(v.hBox)
}

func (v *View) updateCaveButton(birthEntry, deathEntry *widget.Entry) {
	if v.controller.IsEmpty(v.controller.GetObjCave()) {
		v.controller.GenerateCave(v.colsGenerateCave, v.rowsGenerateCave, v.n, v.chance)
	} else {
		var err error

		v.birth, err = strconv.Atoi(birthEntry.Text)
		if err != nil || v.birth < 0 || v.birth > 7 {
			v.birth = 4
		}

		v.death, _ = strconv.Atoi(deathEntry.Text)
		if err != nil || v.death < 0 || v.death > 7 {
			v.death = 3
		}
		v.controller.UpdateCave(v.birth, v.death)
	}

	v.controller.CreatePixelMatrix(&v.pixelMatrix)
	v.controller.ScaleToPixelMatrix(v.pixelMatrix, v.controller.GetObjCave())
	v.createBoxCave()

	v.mainWin.SetContent(v.hBox)
}

func (v *View) autoUpdateButton(millisecondsEntry *widget.Entry) {
	var err error

	v.milliseconds, err = strconv.Atoi(millisecondsEntry.Text)
	if err != nil || v.milliseconds < 0 || v.milliseconds > 3000 {
		v.milliseconds = 500
	}

	if v.controller.IsEmpty(v.controller.GetObjCave()) {
		v.controller.GenerateCave(v.colsGenerateCave, v.rowsGenerateCave, v.n, v.chance)
	}

	for {
		res, _ := v.controller.UpdateCave(v.birth, v.death)
		v.controller.CreatePixelMatrix(&v.pixelMatrix)
		v.controller.ScaleToPixelMatrix(v.pixelMatrix, v.controller.GetObjCave())
		v.createBoxCave()

		v.mainWin.SetContent(v.hBox)

		if !res {
			break
		}
		time.Sleep(time.Millisecond * time.Duration(v.milliseconds))
	}
}

func (v *View) clearMazeButton(obj interface{}) {
	v.qLearningFlag = false
	v.controller.Clear(v.controller.GetObjQLearning())
	v.controller.Clear(obj)
	v.mazeButton()
}

func (v *View) clearCaveButton(obj interface{}) {
	v.controller.Clear(obj)
	v.caveButton()
}

func (v *View) moveButton(orientation, vector bool) {
	if !v.controller.IsEmpty(v.controller.GetObjMaze()) {
		pointX, pointY := v.walkPoints.pointX, v.walkPoints.pointY

		if v.walkPoints.pointX == 0 && v.walkPoints.pointY == 0 {
			pointX, pointY = v.controller.ScalePoints()
		}

		v.controller.CreatePixelMatrix(&v.pixelMatrix)
		v.controller.ScaleToPixelMatrix(v.pixelMatrix, v.controller.GetObjMaze())
		if v.walkPoints.setWalkPoints {
			v.controller.DrawPoint(v.pixelMatrix, v.walkPoints.startX,
				v.walkPoints.startY, color.RGBA{255, 255, 0, 255})
			v.controller.MovePoint(v.pixelMatrix, &pointX, &pointY,
				orientation, vector, color.RGBA{255, 0, 255, 255})
		} else {
			v.walkPoints.endX, v.walkPoints.endY = 0, 0
			v.controller.MovePoint(v.pixelMatrix, &pointX, &pointY,
				orientation, vector, color.RGBA{255, 255, 0, 255})
		}
		v.walkPoints.pointX, v.walkPoints.pointY = pointX, pointY
		v.createBoxMaze()

		v.mainWin.SetContent(v.hBox)
	}
}

func (v *View) okButton() {
	if !v.controller.IsEmpty(v.controller.GetObjMaze()) {

		if !v.walkPoints.setWalkPoints {
			if v.walkPoints.pointX != 0 || v.walkPoints.pointY != 0 {
				v.walkPoints.startX, v.walkPoints.startY = v.walkPoints.pointX, v.walkPoints.pointY
			} else if v.walkPoints.startX != 0 && v.walkPoints.startY != 0 {
				v.controller.CreatePixelMatrix(&v.pixelMatrix)
				v.controller.ScaleToPixelMatrix(v.pixelMatrix, v.controller.GetObjMaze())
				v.controller.DrawPoint(v.pixelMatrix, v.walkPoints.startX,
					v.walkPoints.startY, color.RGBA{255, 255, 0, 255})
				v.createBoxMaze()
				v.mainWin.SetContent(v.hBox)
			}
			v.walkPoints.setWalkPoints = true
		} else {
			if v.walkPoints.pointX != 0 || v.walkPoints.pointY != 0 {
				v.walkPoints.endX, v.walkPoints.endY = v.walkPoints.pointX, v.walkPoints.pointY
			}
			v.walkPoints.setWalkPoints = false
		}
		v.walkPoints.pointX, v.walkPoints.pointY = 0, 0
	}
}

func (v *View) findPathButton() {
	if !v.controller.IsEmpty(v.controller.GetObjMaze()) {
		if v.walkPoints.startX != 0 && v.walkPoints.startY != 0 &&
			v.walkPoints.endX != 0 && v.walkPoints.endY != 0 {
			v.controller.FindPath(v.pixelMatrix, v.walkPoints.startX, v.walkPoints.startY,
				v.walkPoints.endX, v.walkPoints.endY)
		}
		v.createBoxMaze()
		v.mainWin.SetContent(v.hBox)
	}
}

func (v *View) mazeButton() {
	v.walkPoints = Points{}
	v.controller.CreatePixelMatrix(&v.pixelMatrix)
	v.controller.ScaleToPixelMatrix(v.pixelMatrix, v.controller.GetObjMaze())
	v.createBoxMaze()
	v.mainWin.SetContent(v.hBox)
}

func (v *View) caveButton() {
	v.controller.CreatePixelMatrix(&v.pixelMatrix)
	v.controller.ScaleToPixelMatrix(v.pixelMatrix, v.controller.GetObjCave())
	v.createBoxCave()
	v.mainWin.SetContent(v.hBox)
}

func (v *View) qLearningButton() {
	switch {
	case !v.qLearningFlag:
		v.initQLearning()

	case v.walkPoints.startX != 0 && v.walkPoints.startY != 0:
		if !v.controller.IsTrained() {
			v.controller.TrainAgent()
			v.controller.InteractWithAgent(v.pixelMatrix, v.walkPoints.startX, v.walkPoints.startY)
		} else {
			v.controller.InteractWithAgent(v.pixelMatrix, v.walkPoints.startX, v.walkPoints.startY)
		}
		v.walkPoints = Points{}
		v.createBoxMaze()
		v.mainWin.SetContent(v.hBox)
	default:
		v.initQLearning()
	}
}

func (v *View) initQLearning() {
	v.controller.InitQLearning()
	v.loadButton(v.controller.GetObjQLearning())
	v.qLearningFlag = true
}
