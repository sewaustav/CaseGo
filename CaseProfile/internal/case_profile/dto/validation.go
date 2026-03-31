package dto

import "time"

type Result struct {
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
