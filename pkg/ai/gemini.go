package ai

import (
	"context"
	"errors"
	"log"
	"os"
	"sync"

	"google.golang.org/genai"
)

var (
	client     *genai.Client
	clientOnce sync.Once
	clientErr  error
)

// InitClient initializes the Gemini client singleton.
// It assumes GEMINI_API_KEY is set in the environment variables.
func InitClient() {
	clientOnce.Do(func() {
		apiKey := os.Getenv("GEMINI_API_KEY")
		if apiKey == "" {
			clientErr = errors.New("GEMINI_API_KEY environment variable is not set")
			log.Println("Warning:", clientErr)
			return
		}

		ctx := context.Background()
		client, clientErr = genai.NewClient(ctx, nil)
		if clientErr != nil {
			log.Printf("Failed to create Gemini client: %v\n", clientErr)
		} else {
			log.Println("Gemini client initialized successfully")
		}
	})
}

// GenerateContent generates text content using the Gemini Flash 2.5 model.
func GenerateContent(ctx context.Context, prompt string) (string, error) {
	// Ensure client is initialized
	if client == nil {
		InitClient()
	}
	if clientErr != nil {
		return "", clientErr
	}
	if client == nil {
		return "", errors.New("gemini client is not initialized")
	}

	// Use gemini-2.5-flash as requested
	modelName := "gemini-2.0-flash"
	// Note: User requested "flash 2.5", but standard model names are usually 1.5-flash, 2.0-flash-exp, etc.
	// However, the user explicitly provided code with "gemini-2.5-flash".
	// I will use exactly what was in the snippet, assuming user knows the specific model alias/version available.
	// Wait, common usage is currently 1.5 or 2.0. If 2.5 doesn't exist it will fail.
	// But I must follow the user's snippet if they claim it works or is the target.
	// Snippet: "gemini-2.5-flash"
	modelName = "gemini-2.5-flash"

	result, err := client.Models.GenerateContent(
		ctx,
		modelName,
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return "", err
	}

	return result.Text(), nil
}
