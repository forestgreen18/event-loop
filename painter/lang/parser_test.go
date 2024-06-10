package lang

import (
	"bytes"
	"testing"
)

func TestCommandProcessor_ProcessCommands(t *testing.T) {
	artboard := NewArtboardState()
	processor := NewCommandProcessor(artboard)

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Empty input",
			input:   "",
			wantErr: false,
		},
		{
			name:    "Valid white command",
			input:   "white",
			wantErr: false,
		},
		{
			name:    "Valid green command",
			input:   "green",
			wantErr: false,
		},
		{
			name:    "Invalid command",
			input:   "invalid",
			wantErr: true,
		},
		{
			name:    "Valid bgrect command",
			input:   "bgrect 0.1 0.1 0.5 0.5",
			wantErr: false,
		},
		{
			name:    "Invalid bgrect command",
			input:   "bgrect 0.1",
			wantErr: true,
		},
		{
			name:    "Valid figure command",
			input:   "figure 0.3 0.3",
			wantErr: false,
		},
		{
			name:    "Invalid figure command",
			input:   "figure 0.3",
			wantErr: true,
		},
		{
			name:    "Valid move command",
			input:   "move 0.2 0.2",
			wantErr: false,
		},
		{
			name:    "Invalid move command",
			input:   "move 0.2",
			wantErr: true,
		},
		{
			name:    "Valid update command",
			input:   "update",
			wantErr: false,
		},
		{
			name:    "Valid reset command",
			input:   "reset",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := processor.ProcessCommands(bytes.NewBufferString(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessCommands() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConvertToCoordinates(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    []int
		wantErr bool
	}{
		{
			name:    "Valid coordinates",
			args:    []string{"0.1", "0.2", "0.3", "0.4"},
			want:    []int{80, 160, 240, 320},
			wantErr: false,
		},
		{
			name:    "Invalid coordinates",
			args:    []string{"invalid", "0.2", "0.3", "0.4"},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertToCoordinates(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertToCoordinates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !compareSlices(got, tt.want) {
				t.Errorf("ConvertToCoordinates() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function to compare two slices.
func compareSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
