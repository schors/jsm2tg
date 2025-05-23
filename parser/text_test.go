package parser

import (
	"testing"
)

func TestDetectRightBlock(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		block   string
		want    bool
		wantIdx int
	}{
		{
			name:  "empty string",
			input: "",
			block: "color",
			want:  false,
		},
		{
			name:  "no left block",
			input: "This is a test.",
			block: "color",
			want:  false,
		},
		{
			name:    "right block present",
			input:   "{color}",
			block:   "color",
			want:    true,
			wantIdx: 7,
		},
		{
			name:    "left block in center",
			input:   "{color}This is a test.",
			block:   "color",
			want:    true,
			wantIdx: 7,
		},
	}

	for _, tt := range tests {
		got, i := DetectRightBlock(tt.input, tt.block)
		if got != tt.want {
			t.Errorf("%s: detectRightBlock(%q) = %v, want %v", tt.name, tt.input, got, tt.want)
		}

		if i != tt.wantIdx {
			t.Errorf("%s: detectRightBlock(%q) index = %d, want %d", tt.name, tt.input, i, tt.wantIdx)
		}
	}
}

func TestDetectLeftBlock(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		block   string
		want    bool
		wantIdx int
	}{
		{
			name:  "empty string",
			input: "",
			block: "color",
			want:  false,
		},
		{
			name:  "no left block",
			input: "This is a test.",
			block: "color",
			want:  false,
		},
		{
			name:    "left block present",
			input:   "{color:red}This is a test.",
			block:   "color",
			want:    true,
			wantIdx: 11,
		},
		{
			name:    "left block in center",
			input:   "{color:red}\nThis is a test.",
			block:   "color",
			want:    true,
			wantIdx: 11,
		},
	}

	for _, tt := range tests {
		got, i := DetectLeftBlock(tt.input, tt.block)
		if got != tt.want {
			t.Errorf("%s: detectLeftBlock(%q) = %v, want %v", tt.name, tt.input, got, tt.want)
		}

		if i != tt.wantIdx {
			t.Errorf("%s: detectLeftBlock(%q) index = %d, want %d", tt.name, tt.input, i, tt.wantIdx)
		}
	}
}
