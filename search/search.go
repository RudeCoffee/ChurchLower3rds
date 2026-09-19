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
	verseMap      map[string]map[int]map[int]*Verse // book -> chapter -> verse -> *Verse
	invertedIndex map[string][]int                  // word -> []verseIndex
	wordFreq      map[string]int                    // word -> count across all verses
	totalVerses   int
	mu            sync.RWMutex
}

// NewEngine creates a new search engine and indexes the provided verses
func NewEngine(verses []Verse) *Engine {
	e := &Engine{
		verses:        verses,
		verseMap:      make(map[string]map[int]map[int]*Verse),
		invertedIndex: make(map[string][]int),
		wordFreq:      make(map[string]int),
		totalVerses:   len(verses),
	}
	e.buildIndex()
	return e
}

// buildIndex creates the inverted index and verseMap
func (e *Engine) buildIndex() {
	for i := range e.verses {
		v := &e.verses[i]

		// Build verse lookup map
		if e.verseMap[v.Book] == nil {
			e.verseMap[v.Book] = make(map[int]map[int]*Verse)
		}
		if e.verseMap[v.Book][v.Chapter] == nil {
			e.verseMap[v.Book][v.Chapter] = make(map[int]*Verse)
		}
		e.verseMap[v.Book][v.Chapter][v.Verse] = v

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

// GetVerseByRef returns a verse pointer by book, chapter, and verse number
func (e *Engine) GetVerseByRef(book string, chapter int, verseNum int) *Verse {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if chMap, ok := e.verseMap[book]; ok {
		if vMap, ok := chMap[chapter]; ok {
			if v, ok := vMap[verseNum]; ok {
				return v
			}
		}
	}
	return nil
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

// Search finds the best matching verses for the given transcript
func (e *Engine) Search(transcript string) []Suggestion {
	e.mu.RLock()
	defer e.mu.RUnlock()

	// 1. Check direct spoken scripture reference (e.g. "John three sixteen")
	if ref := ParseSpokenReference(transcript); ref != nil {
		if chMap, ok := e.verseMap[ref.Book]; ok {
			if vMap, ok := chMap[ref.Chapter]; ok {
				if v, ok := vMap[ref.Verse]; ok {
					return []Suggestion{
						{
							Verse:      *v,
							Score:      1.0,
							Confidence: "High",
						},
					}
				}
			}
		}
	}

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

	if len(queryWords) < 2 {
		return nil // Reduced from 3 to 2 for improved desktop reactivity
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
	threshold := 3.0
	if e.totalVerses < 100 {
		threshold = 0.1
	}

	var topCandidates []int
	for idx, score := range candidates {
		if score > threshold {
			topCandidates = append(topCandidates, idx)
		}
	}

	// Sort by initial score
	sort.Slice(topCandidates, func(i, j int) bool {
		return candidates[topCandidates[i]] > candidates[topCandidates[j]]
	})

	// Desktop expanded candidate limit from 20 to 50
	if len(topCandidates) > 50 {
		topCandidates = topCandidates[:50]
	}

	// Refined scoring: Sequence Alignment / Phrase Matching
	var suggestions []Suggestion
	for _, idx := range topCandidates {
		verse := e.verses[idx]
		verseWords := tokenize(verse.Text)

		// Calculate match ratio
		matchScore := calculateSequenceScore(queryWords, verseWords)

		if matchScore > 0.25 { // Lowered threshold to 25% for broader suggestion coverage
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
func calculateSequenceScore(query []string, verse []string) float64 {
	verseMap := make(map[string][]int)
	for i, w := range verse {
		verseMap[w] = append(verseMap[w], i)
	}

	matchedCount := 0
	lastPos := -1

	for _, qw := range query {
		if positions, ok := verseMap[qw]; ok {
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
				matchedCount++
			}
		}
	}

	return float64(matchedCount) / float64(len(query))
}
