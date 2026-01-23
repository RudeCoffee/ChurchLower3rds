package main

import (
	"strings"
	"testing"
)

// Ensure loadBibleData is called once for benchmarks
func init() {
	// loadBibleData relies on "kjv.json" being in the current directory.
	// We assume the test is run from the root directory.
	loadBibleData()
}

func TestSearchBible(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		wantText string // Substring expected in first result
	}{
		{"Case Insensitive", "LOVE", "love"},
		{"Book Name", "Genesis", "In the beginning God created the heaven and the earth"},
		{"Specific Verse", "John 3:16", "For God so loved the world"},
		{"Partial Word", "begat", "begat"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := searchBible(tt.query)
			if len(results) == 0 {
				t.Fatalf("searchBible(%q) returned no results", tt.query)
			}

			firstResult := results[0]

			// Check if the expected text is in the *Text* of the first result
			// We check Text (original case) but using case-insensitive contains for validation
			if !strings.Contains(strings.ToLower(firstResult.Text), strings.ToLower(tt.wantText)) {
				t.Errorf("searchBible(%q) first result = %q, want text containing %q", tt.query, firstResult.Text, tt.wantText)
			}
		})
	}
}

func BenchmarkSearchBible(b *testing.B) {
	query := "love"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = searchBible(query)
	}
}

func BenchmarkSearchBibleReference(b *testing.B) {
	query := "John 3:16"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = searchBible(query)
	}
}

func BenchmarkSearchBibleComplex(b *testing.B) {
	query := "kingdom of god"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = searchBible(query)
	}
}
