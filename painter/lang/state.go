package lang

import (
	"image/color"

	"github.com/roman-mazur/architecture-lab-3/painter"
)

type ArtboardState struct {
	Background painter.TextureFunc
	Rectangle  painter.TextureFunc
	Shapes     []*painter.Shape
}

func NewArtboardState() *ArtboardState {
	return &ArtboardState{}
}

func (as *ArtboardState) ConfigureBackground(op painter.TextureFunc) {
	as.Background = op
}

func (as *ArtboardState) DefineRectangle(op painter.TextureFunc) {
	as.Rectangle = op
}

func (as *ArtboardState) PlaceShape(s *painter.Shape) {
	as.Shapes = append(as.Shapes, s)
}

func (as *ArtboardState) ClearArtboard() {
	as.Background = painter.TextureFunc(painter.FillTexture(color.Black))
	as.Rectangle = painter.DrawRectangle(0, 0, 0, 0, color.RGBA{255, 0, 0, 255})
	as.Shapes = nil
}

func (as *ArtboardState) RefreshArtboard() []painter.TextureOperation {
	var ops []painter.TextureOperation

	if as.Background != nil {
		ops = append(ops, as.Background)
	}

	if as.Rectangle != nil {
		ops = append(ops, as.Rectangle)
	}

	for _, shape := range as.Shapes {
		ops = append(ops, shape.DrawShape())
	}

	ops = append(ops, painter.MarkUpdated)

	return ops
}

func (as *ArtboardState) RepositionShapes(dx, dy int) {
	for _, shape := range as.Shapes {
		shape.Move(dx, dy)
	}
}
