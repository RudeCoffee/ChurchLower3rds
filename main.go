package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

type Message struct {
	Type    string `json:"type"`
	Book    string `json:"book,omitempty"`
	Verse   string `json:"verse,omitempty"`
	Speaker string `json:"speaker,omitempty"`
	Show    bool   `json:"show,omitempty"`
}

type BibleData struct {
	Books []BibleBook `json:"books"`
}

type BibleBook struct {
	Name     string        `json:"name"`
	Chapters []BibleChapter `json:"chapters"`
}

type BibleChapter struct {
	Chapter int          `json:"chapter"`
	Verses  []BibleVerse `json:"verses"`
	Name    string       `json:"name"`
}

type BibleVerse struct {
	Chapter int    `json:"chapter"`
	Text    string `json:"text"`
	Verse   int    `json:"verse"`
	Name    string `json:"name"`
}

type SearchRequest struct {
	Type  string `json:"type"`
	Query string `json:"query"`
}

type SearchResponse struct {
	Type    string       `json:"type"`
	Results []BibleVerse `json:"results"`
}

type BooksResponse struct {
	Type  string   `json:"type"`
	Books []string `json:"books"`
}

type ChaptersResponse struct {
	Type     string `json:"type"`
	Book     string `json:"book"`
	Chapters []int  `json:"chapters"`
}

type VersesResponse struct {
	Type    string `json:"type"`
	Book    string `json:"book"`
	Chapter int    `json:"chapter"`
	Verses  []int  `json:"verses"`
}

var obsClient *websocket.Conn
var controlClients []*websocket.Conn
var bibleData BibleData

func loadBibleData() {
	data, err := ioutil.ReadFile("kjv.json")
	if err != nil {
		log.Printf("Warning: Could not load KJV Bible data: %v", err)
		log.Println("Bible search will not be available")
		return
	}

	err = json.Unmarshal(data, &bibleData)
	if err != nil {
		log.Printf("Error parsing Bible data: %v", err)
		return
	}

	log.Printf("Loaded %d Bible books", len(bibleData.Books))
}

func searchBible(query string) []BibleVerse {
	var results []BibleVerse
	query = strings.ToLower(query)

	// Search by book and chapter/verse (e.g., "john 3:16")
	parts := strings.Fields(query)
	if len(parts) >= 2 {
		bookName := parts[0]
		reference := parts[1]

		// Parse chapter:verse
		if strings.Contains(reference, ":") {
			refParts := strings.Split(reference, ":")
			if len(refParts) == 2 {
				chapter, err1 := strconv.Atoi(refParts[0])
				verse, err2 := strconv.Atoi(refParts[1])

				if err1 == nil && err2 == nil {
					// Find matching book and verse
					for _, book := range bibleData.Books {
						if strings.Contains(strings.ToLower(book.Name), bookName) {
							// Find the specific chapter
							for _, chapterData := range book.Chapters {
								if chapterData.Chapter == chapter {
									// Find the specific verse
									for _, verseData := range chapterData.Verses {
										if verseData.Verse == verse {
											results = append(results, BibleVerse{
												Chapter: verseData.Chapter,
												Verse:   verseData.Verse,
												Text:    verseData.Text,
												Name:    verseData.Name,
											})
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// If no specific reference found, search text content
	if len(results) == 0 {
		for _, book := range bibleData.Books {
			for _, chapter := range book.Chapters {
				for _, verse := range chapter.Verses {
					if strings.Contains(strings.ToLower(verse.Text), query) ||
						strings.Contains(strings.ToLower(book.Name), query) {
						results = append(results, BibleVerse{
							Chapter: verse.Chapter,
							Verse:   verse.Verse,
							Text:    verse.Text,
							Name:    verse.Name,
						})
						if len(results) >= 20 { // Limit results
							break
						}
					}
				}
				if len(results) >= 20 {
					break
				}
			}
			if len(results) >= 20 {
				break
			}
		}
	}

	return results
}

func getBookNames(filter string) []string {
	var books []string
	filter = strings.ToLower(filter)

	for _, book := range bibleData.Books {
		if filter == "" || strings.Contains(strings.ToLower(book.Name), filter) {
			books = append(books, book.Name)
		}
	}
	return books
}

func getChapterNumbers(bookName string) []int {
	var chapters []int

	for _, book := range bibleData.Books {
		if strings.EqualFold(book.Name, bookName) {
			for _, chapter := range book.Chapters {
				chapters = append(chapters, chapter.Chapter)
			}
			break
		}
	}
	return chapters
}

func getVerseNumbers(bookName string, chapterNum int) []int {
	var verses []int

	for _, book := range bibleData.Books {
		if strings.EqualFold(book.Name, bookName) {
			for _, chapter := range book.Chapters {
				if chapter.Chapter == chapterNum {
					for _, verse := range chapter.Verses {
						verses = append(verses, verse.Verse)
					}
					break
				}
			}
			break
		}
	}
	return verses
}

func getVerse(bookName string, chapterNum int, verseNum int) *BibleVerse {
	for _, book := range bibleData.Books {
		if strings.EqualFold(book.Name, bookName) {
			for _, chapter := range book.Chapters {
				if chapter.Chapter == chapterNum {
					for _, verse := range chapter.Verses {
						if verse.Verse == verseNum {
							return &BibleVerse{
								Chapter: verse.Chapter,
								Verse:   verse.Verse,
								Text:    verse.Text,
								Name:    verse.Name,
							}
						}
					}
					break
				}
			}
			break
		}
	}
	return nil
}

func handleOBSWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("OBS upgrade failed: ", err)
		return
	}
	defer conn.Close()

	obsClient = conn
	log.Println("OBS client connected")

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("OBS client disconnected:", err)
			obsClient = nil
			break
		}
	}
}

func handleControlWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Control upgrade failed: ", err)
		return
	}
	defer conn.Close()

	controlClients = append(controlClients, conn)
	log.Println("Control client connected")

	for {
		var rawMsg json.RawMessage
		err := conn.ReadJSON(&rawMsg)
		if err != nil {
			log.Println("Control client disconnected:", err)
			// Remove client from slice
			for i, client := range controlClients {
				if client == conn {
					controlClients = append(controlClients[:i], controlClients[i+1:]...)
					break
				}
			}
			break
		}

		// Parse message to check type
		var msgType struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(rawMsg, &msgType); err != nil {
			continue
		}

		switch msgType.Type {
		case "search":
			var searchReq SearchRequest
			if err := json.Unmarshal(rawMsg, &searchReq); err == nil {
				results := searchBible(searchReq.Query)
				response := SearchResponse{
					Type:    "search_results",
					Results: results,
				}
				conn.WriteJSON(response)
			}

		case "get_books":
			var req struct {
				Type   string `json:"type"`
				Filter string `json:"filter,omitempty"`
			}
			if err := json.Unmarshal(rawMsg, &req); err == nil {
				books := getBookNames(req.Filter)
				response := BooksResponse{
					Type:  "books_response",
					Books: books,
				}
				conn.WriteJSON(response)
			}

		case "get_chapters":
			var req struct {
				Type string `json:"type"`
				Book string `json:"book"`
			}
			if err := json.Unmarshal(rawMsg, &req); err == nil {
				chapters := getChapterNumbers(req.Book)
				response := ChaptersResponse{
					Type:     "chapters_response",
					Book:     req.Book,
					Chapters: chapters,
				}
				conn.WriteJSON(response)
			}

		case "get_verses":
			var req struct {
				Type    string `json:"type"`
				Book    string `json:"book"`
				Chapter int    `json:"chapter"`
			}
			if err := json.Unmarshal(rawMsg, &req); err == nil {
				verses := getVerseNumbers(req.Book, req.Chapter)
				response := VersesResponse{
					Type:    "verses_response",
					Book:    req.Book,
					Chapter: req.Chapter,
					Verses:  verses,
				}
				conn.WriteJSON(response)
			}

		case "get_verse":
			var req struct {
				Type    string `json:"type"`
				Book    string `json:"book"`
				Chapter int    `json:"chapter"`
				Verse   int    `json:"verse"`
			}
			if err := json.Unmarshal(rawMsg, &req); err == nil {
				verse := getVerse(req.Book, req.Chapter, req.Verse)
				if verse != nil {
					// Send verse to OBS
					obsMsg := Message{
						Type:  "bible",
						Book:  verse.Name, // Use the verse name which includes book reference
						Verse: verse.Text,
						Show:  true,
					}
					if obsClient != nil {
						obsClient.WriteJSON(obsMsg)
					}
				}
			}

		case "preview_verse":
			var req struct {
				Type    string `json:"type"`
				Book    string `json:"book"`
				Chapter int    `json:"chapter"`
				Verse   int    `json:"verse"`
			}
			if err := json.Unmarshal(rawMsg, &req); err == nil {
				verse := getVerse(req.Book, req.Chapter, req.Verse)
				if verse != nil {
					// Send verse data back to client for preview
					response := struct {
						Type  string     `json:"type"`
						Verse BibleVerse `json:"verse"`
					}{
						Type:  "verse_preview",
						Verse: *verse,
					}
					conn.WriteJSON(response)
				}
			}

		default:
			// Handle regular message (speaker, show/hide, etc.)
			var msg Message
			if err := json.Unmarshal(rawMsg, &msg); err == nil {
				// Forward message to OBS client
				if obsClient != nil {
					err := obsClient.WriteJSON(msg)
					if err != nil {
						log.Println("Error sending to OBS:", err)
					}
				}
			}
		}
	}
}

func main() {
	// Load Bible data
	loadBibleData()

	// Serve static files
	http.Handle("/", http.FileServer(http.Dir("./")))

	// WebSocket endpoints
	http.HandleFunc("/ws/obs", handleOBSWebSocket)
	http.HandleFunc("/ws/control", handleControlWebSocket)

	fmt.Println("Server starting on :8080")
	fmt.Println("OBS Browser Source URL: http://localhost:8080/obs.html")
	fmt.Println("Control Interface URL: http://localhost:8080/client.html")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
