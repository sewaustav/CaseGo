package models

import "time"

type Case struct {
	ID            int64     `json:"id" db:"id"`
	Topic         string    `json:"topic" db:"topic"`
	Category      int32     `json:"category" db:"category"`
	IsGenerated   bool      `json:"is_generated" db:"is_generated"`
	Description   string    `json:"description" db:"description"`
	FirstQuestion string    `json:"first_question" db:"first_question"`
	Creator       int64     `json:"creator" db:"creator"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

type Dialog struct {
	ID        int64      `json:"id" db:"id"`
	CaseID    int64      `json:"case_id" db:"case_id"`
	UserID    int64      `json:"user_id" db:"user_id"`
	ModelName *string    `json:"model_name" db:"model_name"`
	StartedAt time.Time  `json:"started_at" db:"started_at"`
	EndedAt   *time.Time `json:"ended_at" db:"ended_at"`
}

type Interaction struct {
	ID         int64     `json:"id" db:"id"`
	DialogID   int64     `json:"dialog_id" db:"dialog_id"`
	Step       int32     `json:"step" db:"step"`
	Question   string    `json:"question" db:"question"`
	Answer     string    `json:"answer" db:"answer"`
	TokensUsed int32     `json:"tokens_used" db:"tokens_used"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
