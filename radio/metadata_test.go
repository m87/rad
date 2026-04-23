package radio

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMetadata(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Metadata
	}{
		{
			name:  "artist and title",
			input: "StreamTitle='Artist Name - Song Title';",
			want:  Metadata{Artist: "Artist Name", Title: "Song Title"},
		},
		{
			name:  "no separator",
			input: "StreamTitle='Just A Title';",
			want:  Metadata{},
		},
		{
			name:  "empty stream title",
			input: "StreamTitle='';",
			want:  Metadata{},
		},
		{
			name:  "empty string",
			input: "",
			want:  Metadata{},
		},
		{
			name:  "multiple dashes in title",
			input: "StreamTitle='Artist - Title - Remix';",
			want:  Metadata{Artist: "Artist", Title: "Title - Remix"},
		},
		{
			name:  "extra fields ignored",
			input: "StreamTitle='Foo - Bar';StreamUrl='http://example.com';",
			want:  Metadata{Artist: "Foo", Title: "Bar"},
		},
		{
			name:  "no StreamTitle key",
			input: "StreamUrl='http://example.com';",
			want:  Metadata{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseMetadata(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
