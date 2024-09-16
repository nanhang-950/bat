package main

import (
    "github.com/AllenDang/giu"
)

func main() {
    // Initialize the Giu application
    giu.NewMasterWindow("Giu Example", 400, 300, 0, func() {
        // Define the layout of the main window
        giu.SingleWindow().Layout(
            giu.Label("Hello, Giu!"), // Display a label
            giu.Button("Click Me").OnClick(func() {
                // Handle button click
                giu.OpenPopup("Popup")
            }),
            giu.ModalPopup("Popup").Layout(
                giu.Label("Button was clicked!"),
                giu.Button("Close").OnClick(func() {
                    giu.CloseCurrentPopup()
                }),
            ),
        )
    }).Run()
}
