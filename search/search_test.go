package search

import (
	"testing"
)

func TestSearch(t *testing.T) {
	verses := []Verse{
		{Book: "John", Chapter: 3, Verse: 16, Text: "For God so loved the world, that he gave his only begotten Son, that whosoever believeth in him should not perish, but have everlasting life."},
		{Book: "Genesis", Chapter: 1, Verse: 1, Text: "In the beginning God created the heaven and the earth."},
	}

	engine := NewEngine(verses)

	tests := []struct {
		transcript string
		wantBook   string
		wantVerse  int
	}{
		{"god loved the world", "John", 16},
		{"beginning created heaven", "Genesis", 1},
		{"for god so loved world", "John", 16},
	}

	for _, tt := range tests {
		suggestions := engine.Search(tt.transcript)
		if len(suggestions) == 0 {
			t.Errorf("Search(%q) returned no suggestions", tt.transcript)
			continue
		}
		best := suggestions[0]
		if best.Verse.Book != tt.wantBook || best.Verse.Verse != tt.wantVerse {
			t.Errorf("Search(%q) = %v %d:%d, want %s %d", tt.transcript, best.Verse.Book, best.Verse.Chapter, best.Verse.Verse, tt.wantBook, tt.wantVerse)
		}
	}
}
