package models

import "time"

type CaseProfile struct {
	ID                   int64   `json:"id" db:"id"`
	UserID               int64   `json:"user_id" db:"user_id"`
	TotalCases           int64   `json:"total_cases" db:"total_cases"`
	Assertiveness        float64 `json:"assertiveness" db:"assertiveness"`
	Empathy              float64 `json:"empathy" db:"empathy"`
	ClarityCommunication float64 `json:"clarity_communication" db:"clarity_communication"`
	Resistance           float64 `json:"resistance" db:"resistance"`
	Eloquence            float64 `json:"eloquence" db:"eloquence"`
	Initiative           float64 `json:"initiative" db:"initiative"`
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
	Date                 time.Time `json:"date" db:"date"`
}
