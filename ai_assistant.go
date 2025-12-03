package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

// AIClient handles interactions with the local AI model
type AIClient struct {
	client *openai.Client
	model  string
}

// AISuggestion represents the decision made by the AI
type AISuggestion struct {
	Book        string `json:"book"`
	Chapter     int    `json:"chapter"`
	Verse       int    `json:"verse"`
	Confidence  string `json:"confidence"` // "high", "medium", "low"
	Reasoning   string `json:"reasoning,omitempty"`
	Action      string `json:"action"` // "display", "next", "none"
}

// NewAIClient creates a new client connected to the local endpoint
func NewAIClient(endpoint string, model string) *AIClient {
	config := openai.DefaultConfig("dummy-key") // Key not needed for local Ollama usually
	config.BaseURL = endpoint                   // e.g., "http://localhost:11434/v1"

	return &AIClient{
		client: openai.NewClientWithConfig(config),
		model:  model,
	}
}

// AnalyzeTranscript sends the transcript to the AI to decide on the next verse
func (ai *AIClient) AnalyzeTranscript(transcript string, currentBook string, currentChapter int, currentVerse int) (*AISuggestion, error) {
	// Construct the system prompt
	systemPrompt := `You are a helpful assistant for a church projection system.
Your task is to analyze the live transcript of a sermon and decide if a Bible verse should be displayed.
You will be given the current verse being displayed (if any) and the recent transcript.

Rules:
1. If the speaker explicitly mentions a bible reference (e.g., "Let's turn to John 3:16"), identify it.
2. If the speaker is reading a passage and moves to the next verse text, identify the next verse.
3. If the speaker is just talking/preaching without reading specific scripture, return action "none".
4. Return your response in JSON format ONLY.
5. The JSON structure should be: {"book": "string", "chapter": int, "verse": int, "action": "display"|"next"|"none", "reasoning": "string"}
6. "display" action is for a new jump (e.g., jumping to a new book). "next" is for advancing to the next verse in sequence.
`

	userMessage := fmt.Sprintf(`Current Context: %s %d:%d
Transcript: "%s"

Analyze the transcript and determine the correct verse to display.`,
		currentBook, currentChapter, currentVerse, transcript)

	resp, err := ai.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: ai.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userMessage,
				},
			},
			Temperature: 0.1, // Low temperature for deterministic results
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: openai.ChatCompletionResponseFormatTypeJSONObject,
			},
		},
	)

	if err != nil {
		return nil, err
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from AI")
	}

	content := resp.Choices[0].Message.Content

	// Parse JSON response
	var suggestion AISuggestion
	err = json.Unmarshal([]byte(content), &suggestion)
	if err != nil {
		// Try to find JSON in the response if it's wrapped in markdown
		start := strings.Index(content, "{")
		end := strings.LastIndex(content, "}")
		if start != -1 && end != -1 {
			jsonPart := content[start : end+1]
			if err2 := json.Unmarshal([]byte(jsonPart), &suggestion); err2 != nil {
				return nil, fmt.Errorf("failed to parse AI response: %v", err)
			}
		} else {
			return nil, fmt.Errorf("failed to parse AI response: %v", err)
		}
	}

	return &suggestion, nil
}

// CheckAvailability checks if the local AI server is reachable
func (ai *AIClient) CheckAvailability() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Try a lightweight call, e.g., listing models
	_, err := ai.client.ListModels(ctx)
	if err != nil {
		log.Printf("AI Server check failed: %v", err)
		return false
	}
	return true
}
