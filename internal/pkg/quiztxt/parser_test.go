package quiztxt

import "testing"

const sample = `Q: Первый вопрос?
A) Вариант 1
B) Вариант 2
C) Вариант 3
ANSWER: B

Q: Второй вопрос?
A) Да
B) Нет
C) Не знаю
ANSWER: A

Q: Третий?
A) One
B) Two
C) Three
ANSWER: C

Q: Четвёртый?
A) A1
B) B1
C) C1
ANSWER: A

Q: Пятый?
A) X
B) Y
C) Z
ANSWER: B
`

func TestParseValid(t *testing.T) {
	qs, err := Parse(sample)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(qs) != 5 {
		t.Fatalf("expected 5 questions, got %d", len(qs))
	}
	if qs[0].CorrectOption != "B" {
		t.Fatalf("expected B, got %s", qs[0].CorrectOption)
	}
}

func TestParseWrongCount(t *testing.T) {
	_, err := Parse(`Q: Only one?
A) 1
B) 2
C) 3
ANSWER: A`)
	if err == nil {
		t.Fatal("expected error for single question")
	}
}
