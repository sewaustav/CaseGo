package llm_service

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/sewaustav/CaseGoCore/internal/cases/dto"
	"github.com/sewaustav/CaseGoCore/internal/cases/models"
)

const (
	gigachatAuthURL = "https://ngw.devices.sberbank.ru:9443/api/v2/oauth"
	gigachatAPIURL  = "https://gigachat.devices.sberbank.ru/api/v1/chat/completions"
	gigachatModel   = "GigaChat-Pro"
)

type GigaChatService struct {
	authKey    string
	httpClient *http.Client

	mu          sync.Mutex
	accessToken string
	tokenExpiry time.Time
}

func NewGigaChatService(authKey string) *GigaChatService {
	// GigaChat использует российский CA, поэтому пропускаем проверку TLS
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec
	}
	return &GigaChatService{
		authKey: authKey,
		httpClient: &http.Client{
			Transport: tr,
			Timeout:   60 * time.Second,
		},
	}
}

func (g *GigaChatService) getToken(ctx context.Context) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.accessToken != "" && time.Now().Before(g.tokenExpiry) {
		return g.accessToken, nil
	}

	body := url.Values{}
	body.Set("scope", "GIGACHAT_API_PERS")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, gigachatAuthURL,
		strings.NewReader(body.Encode()))
	if err != nil {
		return "", fmt.Errorf("gigachat auth: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+g.authKey)
	req.Header.Set("RqUID", randomUUID())
	req.Header.Set("Accept", "application/json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("gigachat auth request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("gigachat auth status %d: %s", resp.StatusCode, b)
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		ExpiresAt   int64  `json:"expires_at"` // unix ms
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("gigachat auth decode: %w", err)
	}

	g.accessToken = tokenResp.AccessToken
	g.tokenExpiry = time.UnixMilli(tokenResp.ExpiresAt).Add(-30 * time.Second)
	return g.accessToken, nil
}

// ── internal types ────────────────────────────────────────────────────────────

type gcMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type gcRequest struct {
	Model       string      `json:"model"`
	Messages    []gcMessage `json:"messages"`
	Temperature float64     `json:"temperature"`
}

type gcResponse struct {
	Choices []struct {
		Message gcMessage `json:"message"`
	} `json:"choices"`
	Usage struct {
		TotalTokens int `json:"total_tokens"`
	} `json:"usage"`
}

func (g *GigaChatService) chat(ctx context.Context, messages []gcMessage, temperature float64) (*gcResponse, error) {
	token, err := g.getToken(ctx)
	if err != nil {
		return nil, err
	}

	payload, err := json.Marshal(gcRequest{
		Model:       gigachatModel,
		Messages:    messages,
		Temperature: temperature,
	})
	if err != nil {
		return nil, fmt.Errorf("gigachat marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, gigachatAPIURL, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("gigachat build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gigachat request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("gigachat status %d: %s", resp.StatusCode, b)
	}

	var result gcResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("gigachat decode: %w", err)
	}
	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("gigachat: empty choices")
	}
	return &result, nil
}

// ── LLM interface ─────────────────────────────────────────────────────────────

func (g *GigaChatService) GenerateCase(ctx context.Context, description string) (*models.Case, error) {
	systemPrompt := `Ты — профессиональный сценарист обучающих симуляторов.
На основе вводных данных создай JSON объект.
Структура JSON:
{
  "topic": "краткое название кейса",
  "description": "подробное описание ситуации и роли игрока",
  "first_question": "первая фраза, которую скажет ИИ-персонаж для начала диалога"
}
Правила:
- Пиши на русском языке.
- Описание должно быть погружающим.
- Верни ТОЛЬКО JSON без лишнего текста.`

	resp, err := g.chat(ctx, []gcMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: description},
	}, 0.8)
	if err != nil {
		return nil, err
	}

	content := extractJSON(resp.Choices[0].Message.Content)
	var result models.Case
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("gigachat parse case: %w", err)
	}
	result.IsGenerated = true
	return &result, nil
}

func (g *GigaChatService) GenerateResponse(ctx context.Context, caseModel *models.Case, activeCase *models.Dialog, history []models.Interaction) (*dto.CaseDto, error) {
	currentStep := int32(len(history) + 1)

	systemPrompt := fmt.Sprintf(
		"Ты — ведущий симулятора мягких навыков. Кейс: %s. "+
			"Сейчас шаг %d из 10. На 9-м шаге сворачивай историю, на 10-м — завершай. "+
			"ВАЖНО: Отвечай кратко (до 2-3 предложений). "+
			"Отвечай ТОЛЬКО в формате JSON: {\"message\": \"...\", \"is_finished\": true/false}",
		caseModel.Description, currentStep,
	)

	messages := []gcMessage{
		{Role: "system", Content: systemPrompt},
	}
	for _, inter := range history {
		messages = append(messages, gcMessage{Role: "assistant", Content: inter.Question})
		if inter.Answer != "" {
			messages = append(messages, gcMessage{Role: "user", Content: inter.Answer})
		}
	}

	resp, err := g.chat(ctx, messages, 0.8)
	if err != nil {
		return nil, err
	}

	content := extractJSON(resp.Choices[0].Message.Content)
	var result struct {
		Message    string `json:"message"`
		IsFinished bool   `json:"is_finished"`
	}
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("gigachat parse response: %w", err)
	}

	return &dto.CaseDto{
		DialogID: activeCase.ID,
		Question: result.Message,
		Model:    "gigachat",
		Step:     &currentStep,
	}, nil
}

func (g *GigaChatService) AnalyzeCase(ctx context.Context, conv []models.Interaction) (*dto.Result, error) {
	var historyText string
	for _, m := range conv {
		historyText += fmt.Sprintf("Система: %s\nПользователь: %s\n\n", m.Question, m.Answer)
	}

	systemPrompt := `Ты — эксперт по анализу soft-skills. Проанализируй диалог и оцени навыки пользователя.
Оценка для каждого навыка — число от 0.00 до 1.00.
Верни ТОЛЬКО JSON:
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

	resp, err := g.chat(ctx, []gcMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: "Диалог для анализа:\n" + historyText},
	}, 0.3)
	if err != nil {
		return nil, err
	}

	content := extractJSON(resp.Choices[0].Message.Content)
	var analysis struct {
		Skills   dto.Skills `json:"skills_rating"`
		Feedback string     `json:"feedback"`
	}
	if err := json.Unmarshal([]byte(content), &analysis); err != nil {
		return nil, fmt.Errorf("gigachat parse analysis: %w", err)
	}

	return &dto.Result{
		SkillsRating: analysis.Skills,
		StepsCount:   int32(len(conv)),
		TokensUsed:   int32(resp.Usage.TotalTokens),
	}, nil
}

// ── helpers ───────────────────────────────────────────────────────────────────

// extractJSON вырезает первый JSON-объект из строки.
// Нужно потому что GigaChat иногда оборачивает JSON в текст.
func extractJSON(s string) string {
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start == -1 || end == -1 || end <= start {
		return s
	}
	return s[start : end+1]
}

func randomUUID() string {
	b := make([]byte, 16)
	rand.Read(b) //nolint:errcheck
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
