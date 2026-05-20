package quiztxt

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const RequiredQuestions = 5

var (
	questionPrefix = regexp.MustCompile(`^(?:Q:|(\d+)\.)\s*(.+)$`)
	optionPrefix   = regexp.MustCompile(`^([ABC])[\)\.]\s*(.+)$`)
	answerPrefix   = regexp.MustCompile(`^ANSWER:\s*([ABC])\s*$`)
)

// Question is one parsed quiz question.
type Question struct {
	OrderIndex    int
	QuestionText  string
	OptionA       string
	OptionB       string
	OptionC       string
	CorrectOption string
}

// Parse reads a .txt quiz file with exactly 5 questions.
func Parse(content string) ([]Question, error) {
	blocks := splitBlocks(content)
	if len(blocks) == 0 {
		return nil, errors.New("файл пустой")
	}
	if len(blocks) != RequiredQuestions {
		return nil, fmt.Errorf("нужно ровно %d вопросов, найдено %d", RequiredQuestions, len(blocks))
	}

	out := make([]Question, 0, RequiredQuestions)
	for i, block := range blocks {
		q, err := parseBlock(block, i+1)
		if err != nil {
			return nil, fmt.Errorf("вопрос %d: %w", i+1, err)
		}
		out = append(out, q)
	}
	return out, nil
}

func splitBlocks(content string) []string {
	lines := strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")
	blocks := make([][]string, 0)
	current := make([]string, 0)

	flush := func() {
		if len(current) == 0 {
			return
		}
		blocks = append(blocks, current)
		current = make([]string, 0)
	}

	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			flush()
			continue
		}
		current = append(current, line)
	}
	flush()

	out := make([]string, 0, len(blocks))
	for _, b := range blocks {
		out = append(out, strings.Join(b, "\n"))
	}
	return out
}

func parseBlock(block string, order int) (Question, error) {
	lines := strings.Split(block, "\n")
	if len(lines) < 5 {
		return Question{}, errors.New("ожидаются строка вопроса, 3 варианта и ANSWER")
	}

	qMatch := questionPrefix.FindStringSubmatch(lines[0])
	if qMatch == nil {
		return Question{}, errors.New("первая строка должна начинаться с Q: или номера")
	}

	var opts [3]string
	found := map[string]string{}
	correct := ""

	for _, line := range lines[1:] {
		trimmed := strings.TrimSpace(line)
		if m := answerPrefix.FindStringSubmatch(strings.ToUpper(trimmed)); m != nil {
			correct = m[1]
			continue
		}
		if m := optionPrefix.FindStringSubmatch(trimmed); m != nil {
			key := strings.ToUpper(m[1])
			if _, exists := found[key]; exists {
				return Question{}, fmt.Errorf("дубликат варианта %s", key)
			}
			found[key] = strings.TrimSpace(m[2])
		}
	}

	if found["A"] == "" || found["B"] == "" || found["C"] == "" {
		return Question{}, errors.New("нужны варианты A, B и C")
	}
	opts[0], opts[1], opts[2] = found["A"], found["B"], found["C"]

	if correct == "" {
		return Question{}, errors.New("строка ANSWER: A|B|C обязательна")
	}
	if correct != "A" && correct != "B" && correct != "C" {
		return Question{}, errors.New("правильный ответ должен быть A, B или C")
	}

	return Question{
		OrderIndex:    order,
		QuestionText:  strings.TrimSpace(qMatch[2]),
		OptionA:       opts[0],
		OptionB:       opts[1],
		OptionC:       opts[2],
		CorrectOption: correct,
	}, nil
}
