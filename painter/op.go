package painter

import (
	"image"
	"image/color"

	"golang.org/x/exp/shiny/screen"
)

// TextureOperation defines an interface for operations that modify a texture.
type TextureOperation interface {
	// Apply performs the operation and returns true if the texture is ready to be displayed.
	Apply(t screen.Texture) (ready bool)
}

// CompositeOperation combines multiple operations into one.
type CompositeOperation []TextureOperation

func (co CompositeOperation) Apply(t screen.Texture) (ready bool) {
	for _, op := range co {
		ready = op.Apply(t) || ready
	}
	return
}

// MarkUpdated signals that the texture should be considered ready for display.
var MarkUpdated = markUpdated{}

type markUpdated struct{}

func (mu markUpdated) Apply(t screen.Texture) bool { return true }

// TextureFunc wraps a texture update function into a TextureOperation.
type TextureFunc func(t screen.Texture)

func (tf TextureFunc) Apply(t screen.Texture) bool {
	tf(t)
	return false
}

// FillTexture creates a TextureFunc that fills the texture with the specified color.
func FillTexture(fillColor color.Color) TextureFunc {
	return func(t screen.Texture) {
		t.Fill(t.Bounds(), fillColor, screen.Src)
	}
}

// DrawRectangle creates a TextureFunc that draws a rectangle with the specified coordinates and color.
func DrawRectangle(x1, y1, x2, y2 int, rectColor color.Color) TextureFunc {
	return func(t screen.Texture) {
		t.Fill(image.Rect(x1, y1, x2, y2), rectColor, screen.Src)
	}
}

// Shape represents a drawable shape with a center position.
type Shape struct {
	CenterX int
	CenterY int
}

// DrawShape creates a TextureFunc that draws the shape at its center position.
func (s *Shape) DrawShape() TextureFunc {
	return func(t screen.Texture) {
		t.Fill(image.Rect(s.CenterX-150, s.CenterY-100, s.CenterX+150, s.CenterY), color.RGBA{255, 255, 0, 1}, screen.Src)
		t.Fill(image.Rect(s.CenterX-50, s.CenterY, s.CenterX+50, s.CenterY+100), color.RGBA{255, 255, 0, 1}, screen.Src)
	}
}

// Move changes the center position of the shape by the specified offsets.
func (s *Shape) Move(xOffset, yOffset int) {
	s.CenterX += xOffset
	s.CenterY += yOffset
}
