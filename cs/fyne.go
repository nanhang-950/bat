package main

import (
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
)

func main() {
    myApp := app.New()
    myWindow := myApp.NewWindow("Hello")

    myWindow.SetContent(container.NewVBox(
        widget.NewLabel("Hello Fyne!"),
        widget.NewButton("Quit", func() {
            myApp.Quit()
        }),
    ))

    myWindow.ShowAndRun()
}
