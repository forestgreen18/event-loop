package ui

import (
	"image"
	"image/color"
	"log"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/draw"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

type Visualizer struct {
	Title         string
	Debug         bool
	OnScreenReady func(s screen.Screen)

	w    screen.Window
	tx   chan screen.Texture
	done chan struct{}

	sz  size.Event
	pos image.Rectangle
	mousePos image.Point
}

func (pw *Visualizer) Main() {
	pw.tx = make(chan screen.Texture)
	pw.done = make(chan struct{})

	pw.pos.Max.X = 200
	pw.pos.Max.Y = 200
	pw.pos = image.Rect(300, 300, 500, 500) // Center the shape
	pw.mousePos = image.Point{400, 400}
	driver.Main(pw.run)
}

func (pw *Visualizer) UpdateTexture(t screen.Texture) {

	pw.tx <- t
}

func (pw *Visualizer) run(s screen.Screen) {
	w, err := s.NewWindow(&screen.NewWindowOptions{
		Title: pw.Title,
		Height: 800,
		Width: 800,
	})
	if err != nil {
		log.Fatal("Failed to initialize the app window:", err)
	}
	defer func() {
		w.Release()
		close(pw.done)
	}()

	if pw.OnScreenReady != nil {
		pw.OnScreenReady(s)
	}

	pw.w = w

	events := make(chan any)
	go func() {
		for {
			e := w.NextEvent()
			if pw.Debug {
				log.Printf("new event: %v", e)
			}
			if detectTerminate(e) {
				close(events)
				break
			}
			events <- e
		}
	}()

	var t screen.Texture

	for {
		select {
		case e, ok := <-events:
			if !ok {
				return
			}
			pw.handleEvent(e, t)

		case t = <-pw.tx:
			w.Send(paint.Event{})
		}
	}
}

func detectTerminate(e any) bool {
	switch e := e.(type) {
	case lifecycle.Event:
		if e.To == lifecycle.StageDead {
			return true // Window destroy initiated.
		}
	case key.Event:
		if e.Code == key.CodeEscape {
			return true // Esc pressed.
		}
	}
	return false
}

func (pw *Visualizer) handleEvent(e any, t screen.Texture) {
	switch e := e.(type) {
	// ... other cases ...

	case size.Event:
		pw.sz = e
	case paint.Event:
		if t != nil {
			// Use the texture received from the update.
			pw.w.Copy(pw.sz.Bounds().Min, t, t.Bounds(), draw.Src, nil)
		} else {
			// If there is no texture, draw the default UI.
			pw.drawDefaultUI()
		}
		pw.w.Publish() // Publish the window contents to the screen.

	case mouse.Event:
		if e.Button == mouse.ButtonRight {
			pw.mousePos = image.Point{
				X: int(e.X),
				Y: int(e.Y),}
			pw.w.Send(paint.Event{})
		}


	}
}
func (pw *Visualizer) drawDefaultUI() {


	x, y := pw.mousePos.X, pw.mousePos.Y


	pw.w.Fill(pw.sz.Bounds(), color.RGBA{G: 255, A: 255}, draw.Src)

	pw.pos = image.Rect(
		x-100, y-100,
		x+100, y+100,
	)
	pw.drawShape(pw.w, pw.pos)
}


func (pw *Visualizer) drawShape(w screen.Window, pos image.Rectangle) {



	w.Fill(image.Rect(
		pos.Min.X, pos.Min.Y+80,
		pos.Max.X, pos.Min.Y+120,
	),  color.RGBA{B: 255, A: 255}, draw.Src)
	w.Fill(image.Rect(
		pos.Min.X+80, pos.Min.Y,
		pos.Min.X+120, pos.Max.Y,
	), color.RGBA{B: 255, A: 255}, draw.Src)
}
