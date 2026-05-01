package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"esimply/internal/domain"
)

const (
	geminiAPIURL = "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent"
	systemPrompt = "You are a travel-savvy eSIM assistant for eSimply — enthusiastic about travel, genuinely helpful, and always moving the conversation forward. When users mention a destination, affirm their choice and offer a useful tip, data recommendation, or interesting insight. Suggest next steps when the conversation stalls (e.g. \"For a week in Italy, a 5 GB plan usually covers maps, messaging, and daily photos — want me to explain what to look for?\"). Never say you don't know; instead, offer your best guidance and invite them to share more details. Keep replies concise — 2 to 3 sentences max — and skip bullet lists unless the user asks for a breakdown."
)

type Handler struct {
	usecase      domain.CountryUsecase
	geminiAPIKey string
}

func NewHandler(uc domain.CountryUsecase) *Handler {
	return &Handler{
		usecase:      uc,
		geminiAPIKey: os.Getenv("GEMINI_API_KEY"),
	}
}

func (h *Handler) GetCountries(w http.ResponseWriter, r *http.Request) {
	countries, err := h.usecase.GetAllCountries()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errResp("internal server error"))
		return
	}
	writeJSON(w, http.StatusOK, countries)
}

func (h *Handler) GetPlans(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	plans, err := h.usecase.GetPlansByCountryID(id)
	if errors.Is(err, domain.ErrNotFound) {
		writeJSON(w, http.StatusNotFound, errResp("country not found"))
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errResp("internal server error"))
		return
	}
	writeJSON(w, http.StatusOK, plans)
}

type chatRequest struct {
	Messages []chatMessage `json:"messages"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Response string `json:"response"`
}

// Gemini API types

type geminiRequest struct {
	SystemInstruction geminiContent   `json:"systemInstruction"`
	Contents          []geminiContent `json:"contents"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
	Role  string       `json:"role,omitempty"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiResponse struct {
	Candidates []struct {
		Content geminiContent `json:"content"`
	} `json:"candidates"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (h *Handler) AIChat(w http.ResponseWriter, r *http.Request) {
	if h.geminiAPIKey == "" {
		writeJSON(w, http.StatusInternalServerError, errResp("AI service not configured"))
		return
	}

	var req chatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.Messages) == 0 {
		writeJSON(w, http.StatusBadRequest, errResp("invalid request body"))
		return
	}

	contents := make([]geminiContent, len(req.Messages))
	for i, m := range req.Messages {
		role := m.Role
		if role == "assistant" {
			role = "model" // Gemini uses "model" instead of "assistant"
		}
		contents[i] = geminiContent{
			Role:  role,
			Parts: []geminiPart{{Text: m.Content}},
		}
	}

	payload := geminiRequest{
		SystemInstruction: geminiContent{Parts: []geminiPart{{Text: systemPrompt}}},
		Contents:          contents,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errResp("failed to prepare request"))
		return
	}

	url := fmt.Sprintf("%s?key=%s", geminiAPIURL, h.geminiAPIKey)
	httpReq, err := http.NewRequestWithContext(r.Context(), http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errResp("failed to create request"))
		return
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, errResp("failed to reach AI service"))
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, errResp("failed to read AI response"))
		return
	}

	var gr geminiResponse
	if err := json.Unmarshal(respBody, &gr); err != nil {
		writeJSON(w, http.StatusBadGateway, errResp("failed to parse AI response"))
		return
	}

	if resp.StatusCode != http.StatusOK {
		msg := "AI service error"
		if gr.Error != nil {
			msg = gr.Error.Message
		}
		log.Printf("Gemini error (HTTP %d): %s", resp.StatusCode, msg)
		writeJSON(w, http.StatusBadGateway, errResp(msg))
		return
	}

	if len(gr.Candidates) == 0 || len(gr.Candidates[0].Content.Parts) == 0 {
		writeJSON(w, http.StatusBadGateway, errResp("empty AI response"))
		return
	}

	writeJSON(w, http.StatusOK, chatResponse{Response: gr.Candidates[0].Content.Parts[0].Text})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func errResp(msg string) map[string]string {
	return map[string]string{"error": msg}
}
