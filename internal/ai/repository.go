package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository provides DB operations for AI analytics.
type Repository interface {
	GetQuizResultForAnalysis(ctx context.Context, employeeID, quizResultID string) (*EmployeeAnalysis, error)
	SaveAnalysis(ctx context.Context, employeeID, orgID, quizResultID, courseName string, score float64, a *Analysis) error
	GetOrgStats(ctx context.Context, orgID string) (*OrgStats, error)
	ListPendingAnalysis(ctx context.Context, orgID string) ([]PendingItem, error)
	ListPendingForEmployeeCourse(ctx context.Context, employeeID, courseID string) ([]PendingItem, error)
	HasAnalysis(ctx context.Context, employeeID, quizResultID string) (bool, error)
}

// PendingItem is an employee quiz attempt without AI analysis yet.
type PendingItem struct {
	EmployeeID    string
	QuizResultID  string
	OrganizationID string
}

type postgresRepository struct {
	pool *pgxpool.Pool
}

// NewRepository creates an AI repository.
func NewRepository(pool *pgxpool.Pool) Repository {
	return &postgresRepository{pool: pool}
}

// GetQuizResultForAnalysis fetches everything needed to build EmployeeAnalysis.
func (r *postgresRepository) GetQuizResultForAnalysis(ctx context.Context, employeeID, quizResultID string) (*EmployeeAnalysis, error) {
	// 1. Get attempt basics
	const attemptQ = `
		SELECT a.score, a.max_score, a.answers,
		       COALESCE(p.full_name, p.email, ''), 
		       COALESCE(p.department, ''),
		       COALESCE(c.title, '')
		FROM public.lesson_quiz_attempts a
		JOIN public.profiles p ON p.id = a.user_id
		JOIN public.lessons l ON l.id = a.lesson_id
		JOIN public.courses c ON c.id = l.course_id
		WHERE a.id = $1::uuid
		  AND a.user_id = $2::uuid
	`
	var score, maxScore int
	var answersJSON []byte
	var fullName, department, courseName string

	err := r.pool.QueryRow(ctx, attemptQ, quizResultID, employeeID).Scan(
		&score, &maxScore, &answersJSON,
		&fullName, &department, &courseName,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, pgx.ErrNoRows
	}
	if err != nil {
		return nil, fmt.Errorf("get attempt: %w", err)
	}

	// 2. Parse submitted answers: [{question_id, answer}]
	var submitted []struct {
		QuestionID string `json:"question_id"`
		Answer     string `json:"answer"`
	}
	if err := json.Unmarshal(answersJSON, &submitted); err != nil {
		return nil, fmt.Errorf("parse answers json: %w", err)
	}

	// Build a map question_id → submitted answer
	answerMap := make(map[string]string, len(submitted))
	for _, s := range submitted {
		answerMap[s.QuestionID] = s.Answer
	}

	// 3. Load the questions with correct answers
	const questionsQ = `
		SELECT qq.id::text, qq.question_text,
		       qq.option_a, qq.option_b, qq.option_c, qq.correct_option
		FROM public.lesson_quiz_attempts a
		JOIN public.lesson_quizzes lq ON lq.id = a.quiz_id
		JOIN public.lesson_quiz_questions qq ON qq.quiz_id = lq.id
		WHERE a.id = $1::uuid
		ORDER BY qq.order_index
	`
	rows, err := r.pool.Query(ctx, questionsQ, quizResultID)
	if err != nil {
		return nil, fmt.Errorf("get questions: %w", err)
	}
	defer rows.Close()

	type questionRow struct {
		ID            string
		QuestionText  string
		OptionA       string
		OptionB       string
		OptionC       string
		CorrectOption string
	}

	var allQuestions []questionRow
	for rows.Next() {
		var q questionRow
		if err := rows.Scan(&q.ID, &q.QuestionText, &q.OptionA, &q.OptionB, &q.OptionC, &q.CorrectOption); err != nil {
			return nil, err
		}
		allQuestions = append(allQuestions, q)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// 4. Build wrong questions list
	optionText := func(q questionRow, opt string) string {
		switch opt {
		case "A":
			return q.OptionA
		case "B":
			return q.OptionB
		case "C":
			return q.OptionC
		}
		return opt
	}

	var wrong []WrongQuestion
	for _, q := range allQuestions {
		submitted := answerMap[q.ID]
		if submitted != q.CorrectOption {
			wrong = append(wrong, WrongQuestion{
				QuestionText:   q.QuestionText,
				EmployeeAnswer: optionText(q, submitted),
				CorrectAnswer:  optionText(q, q.CorrectOption),
			})
		}
	}

	totalQ := len(allQuestions)
	if totalQ == 0 {
		totalQ = maxScore
	}
	scorePercent := 0.0
	if totalQ > 0 {
		scorePercent = float64(score) / float64(totalQ) * 100
	}

	return &EmployeeAnalysis{
		EmployeeID:     employeeID,
		FullName:       fullName,
		Department:     department,
		CourseName:     courseName,
		Score:          scorePercent,
		TotalQuestions: totalQ,
		WrongQuestions: wrong,
	}, nil
}

// SaveAnalysis persists Gemini analysis to ai_analysis table.
func (r *postgresRepository) SaveAnalysis(
	ctx context.Context,
	employeeID, orgID, quizResultID, courseName string,
	score float64,
	a *Analysis,
) error {
	const q = `
		INSERT INTO public.ai_analysis (
			employee_id, organization_id, quiz_result_id,
			course_name, score, weak_topics,
			recommendation, risk_level, summary
		) VALUES (
			$1::uuid, $2::uuid, $3::uuid,
			$4, $5, $6,
			$7, $8, $9
		)
	`
	_, err := r.pool.Exec(ctx, q,
		employeeID, orgID, quizResultID,
		courseName, score, a.WeakTopics,
		a.Recommendation, a.RiskLevel, a.Summary,
	)
	return err
}

// HasAnalysis checks if analysis already exists for this quiz attempt.
func (r *postgresRepository) HasAnalysis(ctx context.Context, employeeID, quizResultID string) (bool, error) {
	const q = `
		SELECT EXISTS(
			SELECT 1 FROM public.ai_analysis
			WHERE employee_id = $1::uuid AND quiz_result_id = $2::uuid
		)
	`
	var exists bool
	err := r.pool.QueryRow(ctx, q, employeeID, quizResultID).Scan(&exists)
	return exists, err
}

// GetOrgStats returns aggregated AI analytics for an organization.
func (r *postgresRepository) GetOrgStats(ctx context.Context, orgID string) (*OrgStats, error) {
	const q = `
		SELECT 
			aa.employee_id::text,
			COALESCE(p.full_name, p.email, '') AS full_name,
			COALESCE(p.department, '') AS department,
			aa.course_name,
			aa.score,
			aa.risk_level,
			aa.weak_topics,
			aa.recommendation,
			aa.summary,
			aa.created_at::text
		FROM public.ai_analysis aa
		JOIN public.profiles p ON p.id = aa.employee_id
		WHERE aa.organization_id = $1::uuid
		ORDER BY aa.created_at DESC
	`
	rows, err := r.pool.Query(ctx, q, orgID)
	if err != nil {
		return nil, fmt.Errorf("get org stats: %w", err)
	}
	defer rows.Close()

	riskDist := map[string]int{"low": 0, "medium": 0, "high": 0}
	topicCount := map[string]int{}
	totalScore := 0.0
	var employees []EmployeeResult

	for rows.Next() {
		var e EmployeeResult
		var weakTopics []string
		if err := rows.Scan(
			&e.EmployeeID, &e.FullName, &e.Department,
			&e.CourseName, &e.Score, &e.RiskLevel,
			&weakTopics,
			&e.Recommendation, &e.Summary, &e.AnalyzedAt,
		); err != nil {
			return nil, err
		}
		if weakTopics == nil {
			weakTopics = []string{}
		}
		e.WeakTopics = weakTopics
		employees = append(employees, e)

		if _, ok := riskDist[e.RiskLevel]; ok {
			riskDist[e.RiskLevel]++
		}
		totalScore += e.Score
		for _, t := range weakTopics {
			topicCount[t]++
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	avgScore := 0.0
	if len(employees) > 0 {
		avgScore = totalScore / float64(len(employees))
	}

	topTopics := topN(topicCount, 5)
	if employees == nil {
		employees = []EmployeeResult{}
	}
	if topTopics == nil {
		topTopics = []string{}
	}

	return &OrgStats{
		TotalAnalyzed:    len(employees),
		RiskDistribution: riskDist,
		TopWeakTopics:    topTopics,
		AvgScore:         avgScore,
		Employees:        employees,
	}, nil
}

// ListPendingAnalysis finds quiz attempts without an AI analysis for an org.
func (r *postgresRepository) ListPendingAnalysis(ctx context.Context, orgID string) ([]PendingItem, error) {
	const q = `
		SELECT DISTINCT ON (a.user_id)
			a.user_id::text,
			a.id::text,
			p.organization_id::text
		FROM public.lesson_quiz_attempts a
		JOIN public.profiles p ON p.id = a.user_id
		WHERE p.organization_id = $1::uuid
		  AND NOT EXISTS (
		      SELECT 1 FROM public.ai_analysis aa
		      WHERE aa.employee_id = a.user_id
		        AND aa.quiz_result_id = a.id
		  )
		ORDER BY a.user_id, a.completed_at DESC
	`
	rows, err := r.pool.Query(ctx, q, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []PendingItem
	for rows.Next() {
		var it PendingItem
		if err := rows.Scan(&it.EmployeeID, &it.QuizResultID, &it.OrganizationID); err != nil {
			return nil, err
		}
		items = append(items, it)
	}
	return items, rows.Err()
}

// ListPendingForEmployeeCourse returns quiz attempts for a course without AI analysis yet.
func (r *postgresRepository) ListPendingForEmployeeCourse(
	ctx context.Context,
	employeeID, courseID string,
) ([]PendingItem, error) {
	const q = `
		SELECT DISTINCT ON (a.lesson_id)
			a.user_id::text, a.id::text, p.organization_id::text
		FROM public.lesson_quiz_attempts a
		JOIN public.lessons l ON l.id = a.lesson_id
		JOIN public.profiles p ON p.id = a.user_id
		WHERE a.user_id = $1::uuid
		  AND l.course_id = $2::uuid
		  AND p.organization_id IS NOT NULL
		  AND NOT EXISTS (
		      SELECT 1 FROM public.ai_analysis aa
		      WHERE aa.employee_id = a.user_id
		        AND aa.quiz_result_id = a.id
		  )
		ORDER BY a.lesson_id, a.completed_at DESC
	`
	rows, err := r.pool.Query(ctx, q, employeeID, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []PendingItem
	for rows.Next() {
		var it PendingItem
		if err := rows.Scan(&it.EmployeeID, &it.QuizResultID, &it.OrganizationID); err != nil {
			return nil, err
		}
		items = append(items, it)
	}
	return items, rows.Err()
}

func topN(m map[string]int, n int) []string {
	type kv struct {
		Key   string
		Value int
	}
	var pairs []kv
	for k, v := range m {
		pairs = append(pairs, kv{k, v})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Value > pairs[j].Value
	})
	result := make([]string, 0, n)
	for i, p := range pairs {
		if i >= n {
			break
		}
		result = append(result, p.Key)
	}
	return result
}
