package search

import (
	"math"
	"sort"
	"strings"
	"sync"
	"unicode"
)

// Verse represents a bible verse for searching
type Verse struct {
	Book    string `json:"book"`
	Chapter int    `json:"chapter"`
	Verse   int    `json:"verse"`
	Text    string `json:"text"`
}

// Suggestion represents a search result
type Suggestion struct {
	Verse      Verse   `json:"verse"`
	Score      float64 `json:"score"`
	Confidence string  `json:"confidence"`
}

// Engine implements the smart search logic
type Engine struct {
	verses        []Verse
	invertedIndex map[string][]int // word -> []verseIndex
	wordFreq      map[string]int   // word -> count across all verses
	totalVerses   int
	mu            sync.RWMutex
}

// NewEngine creates a new search engine and indexes the provided verses
func NewEngine(verses []Verse) *Engine {
	e := &Engine{
		verses:        verses,
		invertedIndex: make(map[string][]int),
		wordFreq:      make(map[string]int),
		totalVerses:   len(verses),
	}
	e.buildIndex()
	return e
}

// buildIndex creates the inverted index
func (e *Engine) buildIndex() {
	for i, v := range e.verses {
		words := tokenize(v.Text)
		uniqueWords := make(map[string]bool)

		for _, w := range words {
			if !uniqueWords[w] {
				e.invertedIndex[w] = append(e.invertedIndex[w], i)
				e.wordFreq[w]++
				uniqueWords[w] = true
			}
		}
	}
}

// tokenize splits text into lowercase words, removing punctuation
func tokenize(text string) []string {
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}
	fields := strings.FieldsFunc(text, f)
	var words []string
	for _, field := range fields {
		if len(field) > 2 { // Skip very short words to reduce noise
			words = append(words, strings.ToLower(field))
		}
	}
	return words
}

// commonStopWords is a set of words to ignore in search
var commonStopWords = map[string]bool{
	"the": true, "and": true, "that": true, "have": true, "for": true,
	"not": true, "with": true, "you": true, "this": true, "but": true,
	"his": true, "from": true, "they": true, "we": true, "say": true,
	"her": true, "she": true, "will": true, "an": true, "my": true,
	"one": true, "all": true, "would": true, "there": true, "their": true,
	"what": true, "so": true, "up": true, "out": true, "if": true,
	"about": true, "who": true, "get": true, "which": true, "go": true,
	"me": true, "when": true, "make": true, "can": true, "like": true,
	"time": true, "no": true, "just": true, "him": true, "know": true,
	"take": true, "people": true, "into": true, "year": true, "your": true,
	"good": true, "some": true, "could": true, "them": true, "see": true,
	"other": true, "than": true, "then": true, "now": true, "look": true,
	"only": true, "come": true, "its": true, "over": true, "think": true,
	"also": true, "back": true, "after": true, "use": true, "two": true,
	"how": true, "our": true, "work": true, "first": true, "well": true,
	"way": true, "even": true, "new": true, "want": true, "because": true,
	"any": true, "these": true, "give": true, "day": true, "most": true,
	"us": true,
}

// Search finds the best matching verse for the given transcript
func (e *Engine) Search(transcript string) []Suggestion {
	e.mu.RLock()
	defer e.mu.RUnlock()

	words := tokenize(transcript)
	if len(words) == 0 {
		return nil
	}

	// Filter stop words
	var queryWords []string
	for _, w := range words {
		if !commonStopWords[w] {
			queryWords = append(queryWords, w)
		}
	}

	if len(queryWords) < 3 {
		return nil // Need enough context
	}

	// Candidate scoring
	candidates := make(map[int]float64)

	for _, w := range queryWords {
		if idxs, ok := e.invertedIndex[w]; ok {
			// TF-IDF like score
			// IDF: log(totalVerses / wordFreq)
			idf := math.Log(float64(e.totalVerses) / float64(e.wordFreq[w]))

			for _, idx := range idxs {
				candidates[idx] += idf
			}
		}
	}

	// Filter low scores
	threshold := 5.0
	if e.totalVerses < 100 {
		threshold = 0.1
	}

	var topCandidates []int
	for idx, score := range candidates {
		if score > threshold { // Threshold
			topCandidates = append(topCandidates, idx)
		}
	}

	// Sort by initial score
	sort.Slice(topCandidates, func(i, j int) bool {
		return candidates[topCandidates[i]] > candidates[topCandidates[j]]
	})

	if len(topCandidates) > 20 {
		topCandidates = topCandidates[:20]
	}

	// Refined scoring: Sequence Alignment / Phrase Matching
	var suggestions []Suggestion
	for _, idx := range topCandidates {
		verse := e.verses[idx]
		verseWords := tokenize(verse.Text)

		// Calculate match ratio (how many query words are in verse, in order-ish)
		matchScore := calculateSequenceScore(queryWords, verseWords)

		if matchScore > 0.3 { // 30% match
			confidence := "Low"
			if matchScore > 0.6 {
				confidence = "High"
			} else if matchScore > 0.4 {
				confidence = "Medium"
			}

			suggestions = append(suggestions, Suggestion{
				Verse:      verse,
				Score:      matchScore,
				Confidence: confidence,
			})
		}
	}

	// Sort by final score
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Score > suggestions[j].Score
	})

	if len(suggestions) > 3 {
		suggestions = suggestions[:3]
	}

	return suggestions
}

// calculateSequenceScore calculates how well the query matches the verse
// A simple approach: longest common subsequence or just matched words / total query words
func calculateSequenceScore(query []string, verse []string) float64 {
	// Map verse words to positions
	verseMap := make(map[string][]int)
	for i, w := range verse {
		verseMap[w] = append(verseMap[w], i)
	}

	matchedCount := 0
	lastPos := -1

	for _, qw := range query {
		if positions, ok := verseMap[qw]; ok {
			// Find the first position after lastPos
			found := false
			for _, pos := range positions {
				if pos > lastPos {
					lastPos = pos
					matchedCount++
					found = true
					break
				}
			}
			if !found {
				// Word exists but not in order, still count it but reset position tracking?
				// Or just count it as a match but penalize?
				// For simple fuzzy search, just existence is good, order is better.
				// Let's just count existence for now but give bonus for order?
				// Actually, let's stick to simple existence for robustness against paraphrasing.
				matchedCount++
			}
		}
	}

	return float64(matchedCount) / float64(len(query))
}
