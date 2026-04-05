package models

import "time"

type UserRole int

const (
	Admin UserRole = iota
	User
	Creator
)

type UserIdentity struct {
	UserID int64
	Role   UserRole
}

type CaseProfile struct {
	ID                   int64     `json:"id" db:"id"`
	UserID               int64     `json:"user_id" db:"user_id"`
	TotalCases           int64     `json:"total_cases" db:"total_cases"`
	Assertiveness        float32   `json:"assertiveness" db:"assertiveness"`
	Empathy              float32   `json:"empathy" db:"empathy"`
	ClarityCommunication float32   `json:"clarity_communication" db:"clarity_communication"`
	Resistance           float32   `json:"resistance" db:"resistance"`
	Eloquence            float32   `json:"eloquence" db:"eloquence"`
	Initiative           float32   `json:"initiative" db:"initiative"`
	ChangedAt            time.Time `json:"changed_at" db:"changed_at"`
}

type CaseProfileHistory struct {
	ID                   int64     `json:"id" db:"id"`
	UserID               int64     `json:"user_id" db:"user_id"`
	Assertiveness        float64   `json:"assertiveness" db:"assertiveness"`
	Empathy              float64   `json:"empathy" db:"empathy"`
	ClarityCommunication float64   `json:"clarity_communication" db:"clarity_communication"`
	Resistance           float64   `json:"resistance" db:"resistance"`
	Eloquence            float64   `json:"eloquence" db:"eloquence"`
	Initiative           float64   `json:"initiative" db:"initiative"`
	Date                 time.Time `json:"actual_date" db:"actual_date"`
}

type CaseResult struct {
	ID                   int64     `json:"id" db:"id"`
	CaseID               int64     `json:"case_id" db:"case_id"`
	UserID               int64     `json:"user_id" db:"user_id"`
	DialogID             int64     `json:"dialog_id" db:"dialog_id"`
	StepsCount           int32     `json:"steps_count" db:"steps_count"`
	TokensUsed           int32     `json:"tokens_used" db:"tokens_used"`
	FinishedAt           time.Time `json:"finished_at" db:"finished_at"`
	Assertiveness        float32   `json:"assertiveness" db:"assertiveness"`
	Empathy              float32   `json:"empathy" db:"empathy"`
	ClarityCommunication float32   `json:"clarity_communication" db:"clarity_communication"`
	Resistance           float32   `json:"resistance" db:"resistance"`
	Eloquence            float32   `json:"eloquence" db:"eloquence"`
	Initiative           float32   `json:"initiative" db:"initiative"`
}
