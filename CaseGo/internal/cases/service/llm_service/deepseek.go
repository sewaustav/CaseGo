package llm_service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sashabaranov/go-openai"
	"github.com/sewaustav/CaseGoCore/internal/cases/dto"
	"github.com/sewaustav/CaseGoCore/internal/cases/models"
)

func (l *LLMService) GenerateCase(ctx context.Context, description string) (*models.Case, error) {
	resp, err := l.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: "deepseek-chat",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are a professional case generator. Return ONLY JSON.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: description,
			},
		},
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		},
	})

	if err != nil {
		return nil, fmt.Errorf("deepseek error: %w", err)
	}

	var result models.Case
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &result); err != nil {
		return nil, fmt.Errorf("json unmarshal error: %w", err)
	}

	return &result, nil
}

func (l *LLMService) GenerateResponse(ctx context.Context, history []models.Interaction) (*dto.CaseDto, error) {
	return &dto.CaseDto{}, nil
}

func (l *LLMService) AnalyzeCase(ctx context.Context, conv []models.Interaction) (*dto.Result, error) {
	return &dto.Result{}, nil
}
