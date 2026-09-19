package search

import (
	"strconv"
	"strings"
	"unicode"
)

// VerseReference represents a parsed Bible reference
type VerseReference struct {
	Book    string
	Chapter int
	Verse   int
}

var wordToNumber = map[string]int{
	"zero": 0, "one": 1, "first": 1, "1st": 1, "two": 2, "second": 2, "2nd": 2,
	"three": 3, "third": 3, "3rd": 3, "four": 4, "fourth": 4, "4th": 4,
	"five": 5, "fifth": 5, "5th": 5, "six": 6, "sixth": 6, "6th": 6,
	"seven": 7, "seventh": 7, "7th": 7, "eight": 8, "eighth": 8, "8th": 8,
	"nine": 9, "ninth": 9, "9th": 9, "ten": 10, "tenth": 10, "10th": 10,
	"eleven": 11, "eleventh": 11, "twelve": 12, "twelfth": 12,
	"thirteen": 13, "thirteenth": 13, "fourteen": 14, "fourteenth": 14,
	"fifteen": 15, "fifteenth": 15, "sixteen": 16, "sixteenth": 16,
	"seventeen": 17, "seventeenth": 17, "eighteen": 18, "eighteenth": 18,
	"nineteen": 19, "nineteenth": 19, "twenty": 20, "twentieth": 20,
	"thirty": 30, "thirtieth": 30, "forty": 40, "fortieth": 40,
	"fifty": 50, "fiftieth": 50, "sixty": 60, "sixtieth": 60,
	"seventy": 70, "seventieth": 70, "eighty": 80, "eightieth": 80,
	"ninety": 90, "ninetieth": 90, "hundred": 100,
}

// Canonical book names map matching standard KJV book names
var bookAliases = map[string]string{
	"genesis": "Genesis", "exodus": "Exodus", "leviticus": "Leviticus", "numbers": "Numbers",
	"deuteronomy": "Deuteronomy", "joshua": "Joshua", "judges": "Judges", "ruth": "Ruth",
	"1 samuel": "1 Samuel", "first samuel": "1 Samuel", "1st samuel": "1 Samuel", "one samuel": "1 Samuel",
	"2 samuel": "2 Samuel", "second samuel": "2 Samuel", "2nd samuel": "2 Samuel", "two samuel": "2 Samuel",
	"1 kings": "1 Kings", "first kings": "1 Kings", "1st kings": "1 Kings", "one kings": "1 Kings",
	"2 kings": "2 Kings", "second kings": "2 Kings", "2nd kings": "2 Kings", "two kings": "2 Kings",
	"1 chronicles": "1 Chronicles", "first chronicles": "1 Chronicles", "1st chronicles": "1 Chronicles", "one chronicles": "1 Chronicles",
	"2 chronicles": "2 Chronicles", "second chronicles": "2 Chronicles", "2nd chronicles": "2 Chronicles", "two chronicles": "2 Chronicles",
	"ezra": "Ezra", "nehemiah": "Nehemiah", "esther": "Esther", "job": "Job",
	"psalm": "Psalms", "psalms": "Psalms", "proverbs": "Proverbs", "ecclesiastes": "Ecclesiastes",
	"song of solomon": "Song of Solomon", "song of songs": "Song of Solomon", "canticles": "Song of Solomon",
	"isaiah": "Isaiah", "jeremiah": "Jeremiah", "lamentations": "Lamentations", "ezekiel": "Ezekiel", "daniel": "Daniel",
	"hosea": "Hosea", "joel": "Joel", "amos": "Amos", "obadiah": "Obadiah", "jonah": "Jonah",
	"micah": "Micah", "nahum": "Nahum", "habakkuk": "Habakkuk", "zephaniah": "Zephaniah",
	"haggai": "Haggai", "zechariah": "Zechariah", "malachi": "Malachi",
	"matthew": "Matthew", "mark": "Mark", "luke": "Luke", "john": "John",
	"acts": "Acts", "romans": "Romans",
	"1 corinthians": "1 Corinthians", "first corinthians": "1 Corinthians", "1st corinthians": "1 Corinthians", "one corinthians": "1 Corinthians",
	"2 corinthians": "2 Corinthians", "second corinthians": "2 Corinthians", "2nd corinthians": "2 Corinthians", "two corinthians": "2 Corinthians",
	"galatians": "Galatians", "ephesians": "Ephesians", "philippians": "Philippians", "colossians": "Colossians",
	"1 thessalonians": "1 Thessalonians", "first thessalonians": "1 Thessalonians", "1st thessalonians": "1 Thessalonians", "one thessalonians": "1 Thessalonians",
	"2 thessalonians": "2 Thessalonians", "second thessalonians": "2 Thessalonians", "2nd thessalonians": "2 Thessalonians", "two thessalonians": "2 Thessalonians",
	"1 timothy": "1 Timothy", "first timothy": "1 Timothy", "1st timothy": "1 Timothy", "one timothy": "1 Timothy",
	"2 timothy": "2 Timothy", "second timothy": "2 Timothy", "2nd timothy": "2 Timothy", "two timothy": "2 Timothy",
	"titus": "Titus", "philemon": "Philemon", "hebrews": "Hebrews", "james": "James",
	"1 peter": "1 Peter", "first peter": "1 Peter", "1st peter": "1 Peter", "one peter": "1 Peter",
	"2 peter": "2 Peter", "second peter": "2 Peter", "2nd peter": "2 Peter", "two peter": "2 Peter",
	"1 john": "1 John", "first john": "1 John", "1st john": "1 John", "one john": "1 John",
	"2 john": "2 John", "second john": "2 John", "2nd john": "2 John", "two john": "2 John",
	"3 john": "3 John", "third john": "3 John", "3rd john": "3 John", "three john": "3 John",
	"jude": "Jude", "revelation": "Revelation", "revelations": "Revelation",
}

// ParseSpokenReference attempts to extract a Bible book, chapter, and verse from spoken text
func ParseSpokenReference(text string) *VerseReference {
	cleanText := strings.ToLower(text)
	// Replace non-alphanumeric except spaces
	f := func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) {
			return r
		}
		return ' '
	}
	cleanText = strings.Map(f, cleanText)
	words := strings.Fields(cleanText)
	if len(words) == 0 {
		return nil
	}

	// Try matching book names starting from different positions in the transcript
	for i := 0; i < len(words); i++ {
		// Try multi-word book names (e.g., "first corinthians", "song of solomon", "john")
		for wordCount := 3; wordCount >= 1; wordCount-- {
			if i+wordCount > len(words) {
				continue
			}
			bookCandidate := strings.Join(words[i:i+wordCount], " ")
			canonicalBook, found := bookAliases[bookCandidate]
			if !found {
				continue
			}

			// Found book! Now parse numbers following the book name
			remainingWords := words[i+wordCount:]

			// Remove filler words like "chapter", "verse", "and"
			filtered := filterFillerWords(remainingWords)
			ch, v := parseChapterAndVerse(filtered)

			if ch > 0 {
				if v == 0 {
					v = 1 // Default to verse 1 if chapter is given without verse
				}
				return &VerseReference{
					Book:    canonicalBook,
					Chapter: ch,
					Verse:   v,
				}
			}
		}
	}

	return nil
}

func filterFillerWords(words []string) []string {
	var res []string
	for _, w := range words {
		if w == "chapter" || w == "verse" || w == "and" || w == "verses" || w == "in" || w == "on" || w == "the" || w == "of" {
			continue
		}
		res = append(res, w)
	}
	return res
}

// parseChapterAndVerse converts remaining number tokens into chapter and verse integers
func parseChapterAndVerse(tokens []string) (int, int) {
	if len(tokens) == 0 {
		return 0, 0
	}

	nums := parseAllNumbers(tokens)
	if len(nums) == 0 {
		return 0, 0
	}

	if len(nums) == 1 {
		return nums[0], 0
	}

	return nums[0], nums[1]
}

// parseAllNumbers converts a sequence of tokens into integer numbers
func parseAllNumbers(tokens []string) []int {
	var numbers []int
	currentVal := 0
	inNumber := false

	for _, token := range tokens {
		// Check if it's digit string
		if val, err := strconv.Atoi(token); err == nil {
			if inNumber {
				numbers = append(numbers, currentVal)
				currentVal = 0
				inNumber = false
			}
			numbers = append(numbers, val)
			continue
		}

		// Check if it's a number word
		if val, ok := wordToNumber[token]; ok {
			inNumber = true
			if val == 100 {
				if currentVal == 0 {
					currentVal = 100
				} else {
					currentVal *= 100
				}
			} else if currentVal > 0 && currentVal%10 == 0 && val < 10 {
				// e.g. twenty (20) + three (3) -> 23
				currentVal += val
			} else if currentVal > 0 && val >= 10 && val < 100 {
				// e.g. previous number complete, start new
				numbers = append(numbers, currentVal)
				currentVal = val
			} else {
				if currentVal > 0 {
					numbers = append(numbers, currentVal)
				}
				currentVal = val
			}
		} else {
			if inNumber {
				numbers = append(numbers, currentVal)
				currentVal = 0
				inNumber = false
			}
		}
	}

	if inNumber && currentVal > 0 {
		numbers = append(numbers, currentVal)
	}

	return numbers
}
