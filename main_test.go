package main

import (
	"net/http"
	"net/http/httptest"
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

func TestSecureFileHandler(t *testing.T) {
	// Create a dummy handler to wrap
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	handler := secureFileHandler(nextHandler)

	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedLoc    string // for redirects
	}{
		{"Root redirect", "/", http.StatusTemporaryRedirect, "/client.html"},
		{"Allowed HTML", "/client.html", http.StatusOK, ""},
		{"Allowed PNG", "/logo.png", http.StatusOK, ""},
		{"Allowed CSS", "/style.css", http.StatusOK, ""},
		{"Blocked Go Source", "/main.go", http.StatusForbidden, ""},
		{"Blocked Git Config", "/.git/config", http.StatusForbidden, ""},
		{"Blocked Sensitive Config", "/speakers.txt", http.StatusForbidden, ""},
		{"Blocked JSON Data", "/kjv.json", http.StatusForbidden, ""},
		{"Blocked Unknown", "/readme.md", http.StatusForbidden, ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tc.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("handler returned wrong status code for %s: got %v want %v",
					tc.path, rr.Code, tc.expectedStatus)
			}

			if tc.expectedLoc != "" {
				loc := rr.Header().Get("Location")
				if loc != tc.expectedLoc {
					t.Errorf("handler returned wrong redirect location for %s: got %v want %v",
						tc.path, loc, tc.expectedLoc)
				}
			}
		})
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
