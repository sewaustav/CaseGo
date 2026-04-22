package dto

import "time"

type CaseInitialDto struct {
	CaseID int64 `json:"case_id"`
}

type InteractionDto struct {
	DialogID int64  `json:"dialog_id"`
	Step     int32  `json:"step"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type CaseDto struct {
	DialogID int64  `json:"dialog_id"`
	Question string `json:"question"`
	Model    string `json:"model"`
	Step     *int32 `json:"step"`
}

type NewCaseDto struct {
	Topic         *string `json:"topic" validate:"excluded_with=Prompt"`
	Category      *int32  `json:"category" validate:"excluded_with=Prompt"`
	Description   *string `json:"description" validate:"excluded_with=Prompt"`
	FirstQuestion *string `json:"first_question" validate:"excluded_with=Prompt"`

	Prompt *string `json:"prompt" validate:"required_without_all=Topic Category Description FirstQuestion"`
}

type UserSettingsDto struct {
	Topic    *string `json:"topic"`
	Category *int32  `json:"category"`
	Model    *string `json:"model"`
}

type StartDialogResponse struct {
	DialogID      int64  `json:"dialog_id"`
	CaseID        int64  `json:"case_id"`
	FirstQuestion string `json:"first_question"`
	Step          int    `json:"step"`
}

type GetCasesDto struct {
	Limit    int     `json:"limit" form:"limit"`
	Page     int     `json:"page" form:"page"`
	Topic    *string `json:"topic" form:"topic"`
	Category *int32  `json:"category" form:"category"`
}

type Skills struct {
	Assertiveness        float64 `json:"assertiveness"`
	Empathy              float64 `json:"empathy"`
	ClarityCommunication float64 `json:"clarity_communication"`
	Resistance           float64 `json:"resistance"`
	Eloquence            float64 `json:"eloquence"`
	Initiative           float64 `json:"initiative"`
}

type Result struct {
	CaseID       int64     `json:"case_id"`
	UserID       int64     `json:"user_id"`
	DialogID     int64     `json:"dialog_id"`
	StepsCount   int32     `json:"steps_count"`
	TokensUsed   int32     `json:"tokens_used"`
	SkillsRating Skills    `json:"skills_rating"`
	FinishedAt   time.Time `json:"finished_at"`
}
