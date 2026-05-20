package ai

import (
	"encoding/json"
	"fmt"
	"strings"
)

const systemPrompt = `Ты эксперт по пожарно-технической безопасности (ПТМ) Республики Казахстан.
Твоя задача — проанализировать результаты теста сотрудника по пожарной безопасности
и дать чёткие, практичные рекомендации на русском языке.
Отвечай ТОЛЬКО валидным JSON без лишнего текста и markdown-обёрток.
Формат ответа:
{
  "weak_topics": ["тема1", "тема2"],
  "recommendation": "Конкретные рекомендации по обучению...",
  "risk_level": "low|medium|high",
  "summary": "Краткое резюме результатов в 1-2 предложениях"
}`

// Analyze calls OpenAI with real employee quiz data and returns structured analysis.
func Analyze(result EmployeeAnalysis) (*Analysis, error) {
	wrongQuestionsText := buildWrongQuestionsText(result.WrongQuestions)

	scorePercent := 0.0
	if result.TotalQuestions > 0 {
		scorePercent = result.Score
	}

	userPrompt := fmt.Sprintf(`Сотрудник: %s
Отдел: %s
Курс: %s
Результат теста: %.1f%% (%d из %d вопросов правильно)

Неправильно отвеченные вопросы:
%s

Определи уровень риска:
- high: менее 60%%
- medium: 60-79%%
- low: 80%% и выше

Выдели слабые темы на основе реальных вопросов выше.
Дай конкретные рекомендации по улучшению знаний.`,
		result.FullName,
		result.Department,
		result.CourseName,
		scorePercent,
		countCorrect(result),
		result.TotalQuestions,
		wrongQuestionsText,
	)

	raw, err := CallLLM(systemPrompt, userPrompt)
	if err != nil {
		return nil, fmt.Errorf("openai call: %w", err)
	}

	return parseAnalysis(raw)
}

func buildWrongQuestionsText(questions []WrongQuestion) string {
	if len(questions) == 0 {
		return "Все вопросы отвечены правильно."
	}

	var sb strings.Builder
	for i, q := range questions {
		sb.WriteString(fmt.Sprintf("%d. Вопрос: %s\n", i+1, q.QuestionText))
		sb.WriteString(fmt.Sprintf("   Ответ сотрудника: %s\n", q.EmployeeAnswer))
		sb.WriteString(fmt.Sprintf("   Правильный ответ: %s\n", q.CorrectAnswer))
		if q.Topic != "" {
			sb.WriteString(fmt.Sprintf("   Тема: %s\n", q.Topic))
		}
	}
	return sb.String()
}

func countCorrect(result EmployeeAnalysis) int {
	wrong := len(result.WrongQuestions)
	correct := result.TotalQuestions - wrong
	if correct < 0 {
		correct = 0
	}
	return correct
}

func parseAnalysis(raw string) (*Analysis, error) {
	// Strip ```json ... ``` wrappers if present
	cleaned := strings.TrimSpace(raw)
	if strings.HasPrefix(cleaned, "```") {
		lines := strings.SplitN(cleaned, "\n", 2)
		if len(lines) == 2 {
			cleaned = lines[1]
		}
		if idx := strings.LastIndex(cleaned, "```"); idx != -1 {
			cleaned = cleaned[:idx]
		}
		cleaned = strings.TrimSpace(cleaned)
	}

	var a Analysis
	if err := json.Unmarshal([]byte(cleaned), &a); err != nil {
		return nil, fmt.Errorf("parse ai response: %w (raw: %s)", err, raw)
	}

	// Normalise risk level
	switch strings.ToLower(a.RiskLevel) {
	case "low", "medium", "high":
		a.RiskLevel = strings.ToLower(a.RiskLevel)
	default:
		a.RiskLevel = "medium"
	}

	if a.WeakTopics == nil {
		a.WeakTopics = []string{}
	}

	return &a, nil
}
