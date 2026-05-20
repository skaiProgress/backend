package quizzes

import "time"

// QuizMeta is admin-facing quiz summary.
type QuizMeta struct {
	ID             string    `json:"id"`
	LessonID       string    `json:"lesson_id"`
	SourceFileName *string   `json:"source_file_name"`
	QuestionCount  int       `json:"question_count"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// AdminQuestion includes the correct answer.
type AdminQuestion struct {
	ID            string `json:"id"`
	OrderIndex    int    `json:"order_index"`
	QuestionText  string `json:"question_text"`
	OptionA       string `json:"option_a"`
	OptionB       string `json:"option_b"`
	OptionC       string `json:"option_c"`
	CorrectOption string `json:"correct_option"`
}

// AdminQuiz is full quiz for admin panel.
type AdminQuiz struct {
	QuizMeta
	Questions []AdminQuestion `json:"questions"`
}

// EmployeeQuestion hides correct answer.
type EmployeeQuestion struct {
	ID           string `json:"id"`
	OrderIndex   int    `json:"order_index"`
	QuestionText string `json:"question_text"`
	Options      []string `json:"options"`
}

// AttemptSummary is internal attempt data (includes score).
type AttemptSummary struct {
	Score       int       `json:"score"`
	MaxScore    int       `json:"max_score"`
	Passed      bool      `json:"passed"`
	CompletedAt time.Time `json:"completed_at"`
}

// EmployeeAttemptView is what employees see — no score, only status.
type EmployeeAttemptView struct {
	Submitted   bool      `json:"submitted"`
	Passed      bool      `json:"passed"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
}

// EmployeeQuizPayload for taking a quiz.
type EmployeeQuizPayload struct {
	LessonID  string             `json:"lesson_id"`
	Questions []EmployeeQuestion `json:"questions"`
	Attempt   *EmployeeAttemptView `json:"attempt,omitempty"`
}

// SubmitAnswer is one answer in a submission.
type SubmitAnswer struct {
	QuestionID string `json:"question_id"`
	Answer     string `json:"answer"`
}

// SubmitRequest body for employee quiz submit.
type SubmitRequest struct {
	Answers []SubmitAnswer `json:"answers"`
}

// SubmitResult response after grading (internal).
type SubmitResult struct {
	Score    int  `json:"score"`
	MaxScore int  `json:"max_score"`
	Passed   bool `json:"passed"`
}

// EmployeeSubmitResult is returned to employees — no numeric score.
type EmployeeSubmitResult struct {
	Submitted bool `json:"submitted"`
	Passed    bool `json:"passed"`
}
