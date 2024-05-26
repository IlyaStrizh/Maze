package main

import (
	maze "Maze/internal/view"

	"fyne.io/fyne/v2/app"
)

func main() {
	myApp := app.New()
	maze.NewView(myApp)

	myApp.Run()
}
