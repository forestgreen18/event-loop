package main

import (
	"net/http"

	"github.com/roman-mazur/architecture-lab-3/painter"
	"github.com/roman-mazur/architecture-lab-3/painter/lang"
	"github.com/roman-mazur/architecture-lab-3/ui"
)

func main() {
	var (
		pv ui.Visualizer // The visualizer creates a window and draws in it.

		// Needed for part 2.
		eventLoop painter.EventLoop // Event loop for processing operations.
		processor lang.CommandProcessor // Command processor.
		artboard  lang.ArtboardState // Artboard state.
	)

	//pv.Debug = true
	pv.Title = "Simple painter"

	pv.OnScreenReady = eventLoop.Initiate
	eventLoop.Receiver = &pv

	// Initialize the command processor with the artboard state.
	processor = *lang.NewCommandProcessor(&artboard)

	go func() {
		http.Handle("/", lang.CommandHttpHandler(&eventLoop, &processor))
		_ = http.ListenAndServe("localhost:17000", nil)
	}()

	pv.Main()
	eventLoop.Terminate()
}
