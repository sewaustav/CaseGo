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
	// Детальный промпт, чтобы модель понимала структуру
	systemPrompt := `Ты — профессиональный сценарист обучающих симуляторов. 
На основе вводных данных создай JSON объект.
Структура JSON:
{
  "topic": "краткое название кейса",
  "description": "подробное описание ситуации и роли игрока",
  "first_question": "первая фраза, которую скажет ИИ-персонаж для начала диалога"
}
Правила: 
- Пиши на языке пользователя.
- Описание должно быть погружающим.
- Верни ТОЛЬКО JSON.`

	resp, err := l.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: "deepseek-chat",
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
			{Role: openai.ChatMessageRoleUser, Content: description},
		},
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		},
		Temperature: 0.8,
	})

	if err != nil {
		return nil, fmt.Errorf("deepseek error: %w", err)
	}

	var result models.Case
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &result); err != nil {
		return nil, fmt.Errorf("json unmarshal error: %w", err)
	}

	result.IsGenerated = true
	return &result, nil
}

func (l *LLMService) GenerateResponse(ctx context.Context, caseModel *models.Case, activeCase *models.Dialog, history []models.Interaction) (*dto.CaseDto, error) {

	currentStep := int32(len(history) + 1)

	systemInstruction := fmt.Sprintf(
		"Ты — ведущий симулятора. Кейс: %s. "+
			"Сейчас шаг %d из 10. На 9-м шаге сворачивай историю, на 10-м — завершай. "+
			"ВАЖНО: Отвечай кратко (до 2-3 предложений). "+
			"Отвечай ТОЛЬКО в формате JSON: {\"message\": \"...\", \"is_finished\": bool}",
		caseModel.Description, currentStep,
	)

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemInstruction,
		},
	}

	for _, inter := range history {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: inter.Question,
		})
		if inter.Answer != "" {
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: inter.Answer,
			})
		}
	}
	resp, err := l.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    "deepseek-chat",
		Messages: messages,
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		},
		Temperature: 0.8,
	})

	if err != nil {
		return nil, fmt.Errorf("api error: %v", err)
	}

	var result struct {
		Message    string `json:"message"`
		IsFinished bool   `json:"is_finished"`
	}

	content := resp.Choices[0].Message.Content
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("failed to parse model response: %v", err)
	}

	return &dto.CaseDto{
		DialogID: activeCase.ID,
		Question: result.Message,
		Model:    "deepseek",
		Step:     &currentStep,
	}, nil
}

func (l *LLMService) AnalyzeCase(ctx context.Context, conv []models.Interaction) (*dto.Result, error) {
	var historyText string
	for _, m := range conv {
		historyText += fmt.Sprintf("Система: %s\nПользователь: %s\n\n", m.Question, m.Answer)
	}

	// Промпт под твою структуру dto.Skills
	systemPrompt := `Ты — эксперт по анализу soft-skills. Проанализируй диалог и оцени навыки пользователя.
Оценка для каждого навыка — это число от 0.00 до 1.00.
Верни ТОЛЬКО JSON с полями:
{
  "skills_rating": {
    "assertiveness": 0.00,
    "empathy": 0.00,
    "clarity_communication": 0.00,
    "resistance": 0.00,
    "eloquence": 0.00,
    "initiative": 0.00
  },
  "feedback": "краткий текстовый разбор"
}
Будь строгим, но объективным аналитиком.`

	resp, err := l.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: "deepseek-chat",
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
			{Role: openai.ChatMessageRoleUser, Content: "Диалог для анализа:\n" + historyText},
		},
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		},
		Temperature: 0.3,
	})

	if err != nil {
		return nil, err
	}

	var analysis struct {
		Skills   dto.Skills `json:"skills_rating"`
		Feedback string     `json:"feedback"`
	}

	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &analysis); err != nil {
		return nil, err
	}

	return &dto.Result{
		SkillsRating: analysis.Skills,
		StepsCount:   int32(len(conv)),
		TokensUsed:   int32(resp.Usage.TotalTokens),
	}, nil
}
