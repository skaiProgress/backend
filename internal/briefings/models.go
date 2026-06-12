package briefings

import "time"

// BriefingKind enumerates all possible briefing types.
type BriefingKind string

const (
	KindIntroductory BriefingKind = "introductory"
	KindPrimary      BriefingKind = "primary"
	KindRepeat       BriefingKind = "repeat"
	KindUnscheduled  BriefingKind = "unscheduled"
	KindTargeted     BriefingKind = "targeted"
)

// BriefingKindLabel returns the human-readable Russian label.
func (k BriefingKind) Label() string {
	switch k {
	case KindIntroductory:
		return "Вводный"
	case KindPrimary:
		return "Первичный"
	case KindRepeat:
		return "Повторный"
	case KindUnscheduled:
		return "Внеплановый"
	case KindTargeted:
		return "Целевой"
	default:
		return string(k)
	}
}

// OrgEvent is a calendar event belonging to an organization.
type OrgEvent struct {
	ID             string       `json:"id"`
	OrganizationID string       `json:"organization_id"`
	EmployeeID     *string      `json:"employee_id,omitempty"`
	CourseID       *string      `json:"course_id,omitempty"`
	Title          string       `json:"title"`
	EventType      string       `json:"type"`
	BriefingKind   *string      `json:"briefing_kind,omitempty"`
	StartsAt       time.Time    `json:"starts_at"`
	EndsAt         *time.Time   `json:"ends_at,omitempty"`
	Time           string       `json:"time"`
	Location       string       `json:"location"`
	Participants   *int         `json:"participants,omitempty"`
	CreatedAt      time.Time    `json:"created_at"`
}

// BriefingRecord is a journal entry for a fire-safety briefing.
type BriefingRecord struct {
	ID                  string     `json:"id"`
	OrganizationID      string     `json:"organization_id"`
	EventID             *string    `json:"event_id,omitempty"`
	EmployeeID          string     `json:"employee_id"`
	EmployeeName        string     `json:"employee_name"`
	Position            string     `json:"position"`
	BriefingKind        string     `json:"briefing_kind"`
	InstructorName      string     `json:"instructor_name"`
	InstructorID        *string    `json:"instructor_id,omitempty"`
	DateConducted       string     `json:"date_conducted"`
	EmployeeSigned      bool       `json:"employee_signed"`
	EmployeeSignedAt    *time.Time `json:"employee_signed_at,omitempty"`
	InstructorSigned    bool       `json:"instructor_signed"`
	InstructorSignedAt  *time.Time `json:"instructor_signed_at,omitempty"`
	RowNumber           int        `json:"row_number"`
	CreatedAt           time.Time  `json:"created_at"`
}

// CreateEventInput is used internally to schedule a new calendar event.
type CreateEventInput struct {
	OrganizationID string
	EmployeeID     *string
	CourseID       *string
	Title          string
	EventType      string
	BriefingKind   *BriefingKind
	StartsAt       time.Time
	EndsAt         *time.Time
	Location       string
	Participants   *int
	CreatedBy      *string
}

// EmployeeBriefing is what an employee sees as a pending briefing.
type EmployeeBriefing struct {
	EventID      string     `json:"event_id"`
	Title        string     `json:"title"`
	BriefingKind string     `json:"briefing_kind"`
	StartsAt     time.Time  `json:"starts_at"`
	EndsAt       *time.Time `json:"ends_at,omitempty"`
	Location     string     `json:"location"`
	RecordID     *string    `json:"record_id,omitempty"`
	Confirmed    bool       `json:"confirmed"`
}

// EmployeeBriefingDetail is the payload for the dedicated briefing page.
type EmployeeBriefingDetail struct {
	EventID      string     `json:"event_id"`
	Title        string     `json:"title"`
	BriefingKind string     `json:"briefing_kind"`
	Location     string     `json:"location"`
	StartsAt     time.Time  `json:"starts_at"`
	EndsAt       *time.Time `json:"ends_at,omitempty"`
	VideoURL     string     `json:"video_url"`
	Confirmed    bool       `json:"confirmed"`
	WindowActive bool       `json:"window_active"`
	NotStarted   bool       `json:"not_started"`
	Expired      bool       `json:"expired"`
}

// BriefingVideo is one uploaded video for a course/kind.
type BriefingVideo struct {
	BriefingKind string `json:"briefing_kind"`
	VideoURL     string `json:"video_url"`
}

// BriefingCourse is a course an org-admin can use to create briefing links.
type BriefingCourse struct {
	CourseID string   `json:"course_id"`
	Title    string   `json:"title"`
	Kinds    []string `json:"kinds"` // briefing kinds that have a video
}

// UpdateEventRequest allows org-admin to change event date/time.
type UpdateEventRequest struct {
	StartsAt string `json:"starts_at"` // RFC3339
}

// CreateBriefingEventRequest schedules any of the 5 briefing kinds tied to a course video.
type CreateBriefingEventRequest struct {
	EmployeeID   string `json:"employee_id"`
	CourseID     string `json:"course_id"`
	BriefingKind string `json:"briefing_kind"` // any of the 5 kinds
	StartsAt     string `json:"starts_at"`     // RFC3339, window start
	EndsAt       string `json:"ends_at"`       // RFC3339, window end
	Location     string `json:"location"`
}
