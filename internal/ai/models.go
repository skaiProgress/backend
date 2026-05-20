package ai

// EmployeeAnalysis holds the raw data needed to build an AI prompt.
type EmployeeAnalysis struct {
	EmployeeID     string
	FullName       string
	Department     string
	CourseName     string
	Score          float64
	TotalQuestions int
	WrongQuestions []WrongQuestion
}

// WrongQuestion represents one incorrectly answered question.
type WrongQuestion struct {
	QuestionText   string
	EmployeeAnswer string
	CorrectAnswer  string
	Topic          string
}

// Analysis is the parsed output from Gemini.
type Analysis struct {
	WeakTopics     []string `json:"weak_topics"`
	Recommendation string   `json:"recommendation"`
	RiskLevel      string   `json:"risk_level"` // low | medium | high
	Summary        string   `json:"summary"`
}

// OrgStats aggregates analysis results for the whole organization.
type OrgStats struct {
	TotalAnalyzed    int              `json:"total_analyzed"`
	RiskDistribution map[string]int   `json:"risk_distribution"`
	TopWeakTopics    []string         `json:"top_weak_topics"`
	AvgScore         float64          `json:"avg_score"`
	Employees        []EmployeeResult `json:"employees"`
}

// EmployeeResult is one row in the OrgStats employee list.
type EmployeeResult struct {
	EmployeeID     string   `json:"employee_id"`
	FullName       string   `json:"full_name"`
	Department     string   `json:"department"`
	CourseName     string   `json:"course_name"`
	Score          float64  `json:"score"`
	RiskLevel      string   `json:"risk_level"`
	WeakTopics     []string `json:"weak_topics"`
	Recommendation string   `json:"recommendation"`
	Summary        string   `json:"summary"`
	AnalyzedAt     string   `json:"analyzed_at"`
}
