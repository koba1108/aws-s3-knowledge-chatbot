package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime/types"
	"github.com/rs/cors"
)

type ChatRequest struct {
	Message        string `json:"message"`
	SessionID      string `json:"session_id,omitempty"`
	KnowledgeBaseID string `json:"knowledge_base_id,omitempty"`
}

type ChatResponse struct {
	Response  string              `json:"response"`
	SessionID string              `json:"session_id"`
	Sources   []RetrievedSource   `json:"sources,omitempty"`
	Error     string              `json:"error,omitempty"`
}

type RetrievedSource struct {
	Content  string            `json:"content"`
	Location map[string]string `json:"location"`
}

type Server struct {
	bedrockClient   *bedrockagentruntime.Client
	knowledgeBaseID string
	modelID         string
}

func NewServer() (*Server, error) {
	ctx := context.Background()
	
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS config: %w", err)
	}

	// Get configuration from environment
	knowledgeBaseID := os.Getenv("KNOWLEDGE_BASE_ID")
	if knowledgeBaseID == "" {
		log.Println("Warning: KNOWLEDGE_BASE_ID not set, using default")
		knowledgeBaseID = "default-kb-id"
	}

	modelID := os.Getenv("MODEL_ID")
	if modelID == "" {
		modelID = "anthropic.claude-3-sonnet-20240229-v1:0"
	}

	return &Server{
		bedrockClient:   bedrockagentruntime.NewFromConfig(cfg),
		knowledgeBaseID: knowledgeBaseID,
		modelID:         modelID,
	}, nil
}

func (s *Server) handleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Message == "" {
		http.Error(w, "Message is required", http.StatusBadRequest)
		return
	}

	// Use provided knowledge base ID or default
	kbID := req.KnowledgeBaseID
	if kbID == "" {
		kbID = s.knowledgeBaseID
	}

	// Generate session ID if not provided
	sessionID := req.SessionID
	if sessionID == "" {
		sessionID = fmt.Sprintf("session-%d", time.Now().Unix())
	}

	// Call Bedrock Agent Runtime to retrieve and generate response
	ctx := context.Background()
	input := &bedrockagentruntime.RetrieveAndGenerateInput{
		Input: &types.RetrieveAndGenerateInput{
			Text: &req.Message,
		},
		RetrieveAndGenerateConfiguration: &types.RetrieveAndGenerateConfiguration{
			Type: types.RetrieveAndGenerateTypeKnowledgeBase,
			KnowledgeBaseConfiguration: &types.KnowledgeBaseRetrieveAndGenerateConfiguration{
				KnowledgeBaseId: &kbID,
				ModelArn:        &s.modelID,
			},
		},
		SessionId: &sessionID,
	}

	result, err := s.bedrockClient.RetrieveAndGenerate(ctx, input)
	if err != nil {
		log.Printf("Error calling Bedrock: %v", err)
		resp := ChatResponse{
			Error:     fmt.Sprintf("Failed to get response: %v", err),
			SessionID: sessionID,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Extract response text
	var responseText string
	if result.Output != nil && result.Output.Text != nil {
		responseText = *result.Output.Text
	}

	// Extract sources
	var sources []RetrievedSource
	if result.Citations != nil {
		for _, citation := range result.Citations {
			for _, ref := range citation.RetrievedReferences {
				source := RetrievedSource{
					Location: make(map[string]string),
				}
				if ref.Content != nil && ref.Content.Text != nil {
					source.Content = *ref.Content.Text
				}
				if ref.Location != nil && ref.Location.S3Location != nil {
					source.Location["uri"] = *ref.Location.S3Location.Uri
				}
				sources = append(sources, source)
			}
		}
	}

	resp := ChatResponse{
		Response:  responseText,
		SessionID: *result.SessionId,
		Sources:   sources,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func main() {
	server, err := NewServer()
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/chat", server.handleChat)
	mux.HandleFunc("/api/health", server.handleHealth)

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("Knowledge Base ID: %s", server.knowledgeBaseID)
	log.Printf("Model ID: %s", server.modelID)
	
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
