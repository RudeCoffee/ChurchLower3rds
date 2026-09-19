package search

import (
	"reflect"
	"testing"
)

func TestParseSpokenReference(t *testing.T) {
	tests := []struct {
		input string
		want  *VerseReference
	}{
		{
			input: "turn with me to John three sixteen",
			want:  &VerseReference{Book: "John", Chapter: 3, Verse: 16},
		},
		{
			input: "let us read First Corinthians chapter thirteen verse four",
			want:  &VerseReference{Book: "1 Corinthians", Chapter: 13, Verse: 4},
		},
		{
			input: "open your bibles to Genesis chapter one verse one",
			want:  &VerseReference{Book: "Genesis", Chapter: 1, Verse: 1},
		},
		{
			input: "second samuel twenty three verse two",
			want:  &VerseReference{Book: "2 Samuel", Chapter: 23, Verse: 2},
		},
		{
			input: "psalm twenty three verse one",
			want:  &VerseReference{Book: "Psalms", Chapter: 23, Verse: 1},
		},
		{
			input: "song of solomon chapter two verse five",
			want:  &VerseReference{Book: "Song of Solomon", Chapter: 2, Verse: 5},
		},
		{
			input: "just talking about life and faith",
			want:  nil,
		},
	}

	for _, tt := range tests {
		got := ParseSpokenReference(tt.input)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("ParseSpokenReference(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}
