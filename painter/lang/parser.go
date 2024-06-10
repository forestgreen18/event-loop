package lang

import (
	"bufio"
	"errors"
	"image/color"
	"io"
	"strconv"
	"strings"

	"github.com/roman-mazur/architecture-lab-3/painter"
)

type CommandProcessor struct {
	Artboard *ArtboardState
}

func NewCommandProcessor(artboard *ArtboardState) *CommandProcessor {
	return &CommandProcessor{Artboard: artboard}
}

func (cp *CommandProcessor) ProcessCommands(input io.Reader) ([]painter.TextureOperation, error) {
	var textureOps []painter.TextureOperation

	commandReader := bufio.NewScanner(input)

	for commandReader.Scan() {
		commands := strings.Split(commandReader.Text(), ",")
		if len(commands) == 0 {
			continue
		}

		for _, cmd := range commands {
			cmdParts := strings.Fields(cmd)

			switch cmdParts[0] {
			case "white":
				cp.Artboard.ConfigureBackground(painter.FillTexture(color.White))
			case "green":
				cp.Artboard.ConfigureBackground(painter.FillTexture(color.RGBA{R: 0, G: 128, B: 0, A: 255}))
			case "bgrect":
				if len(cmdParts) != 5 {
					return nil, errors.New("bgrect command expects four arguments")
				}

				coords, err := convertToCoordinates(cmdParts[1:])
				if err != nil {
					return nil, err
				}

				cp.Artboard.DefineRectangle(painter.DrawRectangle(coords[0], coords[1], coords[2], coords[3], color.RGBA{255, 0, 0, 255}))
			case "figure":
				if len(cmdParts) != 3 {
					return nil, errors.New("figure command expects two arguments")
				}

				center, err := convertToCoordinates(cmdParts[1:])
				if err != nil {
					return nil, err
				}

				cp.Artboard.PlaceShape(&painter.Shape{
					CenterX: center[0],
					CenterY: center[1],
				})
			case "move":
				if len(cmdParts) != 3 {
					return nil, errors.New("move command expects two arguments")
				}

				delta, err := convertToCoordinates(cmdParts[1:])
				if err != nil {
					return nil, err
				}

				cp.Artboard.RepositionShapes(delta[0], delta[1])
			case "update":
				textureOps = append(textureOps, cp.Artboard.RefreshArtboard()...)
			case "reset":
				cp.Artboard.ClearArtboard()
			default:
				return nil, errors.New("unrecognized command")
			}
		}
	}

	if err := commandReader.Err(); err != nil {
		return nil, err
	}


	return textureOps, nil
}

func convertToCoordinates(args []string) ([]int, error) {
	coordinates := make([]int, len(args))
	for i, arg := range args {
		value, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			return nil, errors.New("error parsing coordinate")
		}
		coordinates[i] = int(value * 800) // Assuming 800 is the scale factor for the artboard dimensions
	}
	return coordinates, nil
}
