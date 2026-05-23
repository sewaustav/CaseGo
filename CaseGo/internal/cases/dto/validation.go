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
	Xp            *int32  `json:"xp"`

	Prompt *string `json:"prompt" validate:"required_without_all=Topic Category Description FirstQuestion"`
}

type UserSettingsDto struct {
	Topic    *string `json:"topic"`
	Category *int32  `json:"category"`
	Model    *string `json:"model"`
}

type StatsResponse struct {
	TotalCases   int `json:"total_cases"`
	TotalDialogs int `json:"total_dialogs"`
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

type ResultResponse struct {
	Result *Result          `json:"result"`
	Level  *CaseLevelResult `json:"level"`
}

type SubscriptionStatusDto struct {
	Status    int32     `db:"subscription" json:"subscription"`
	ExpiredAt time.Time `db:"expired_at" json:"expired_at"`
}

type LevelXpResult struct {
	Xp   int32     `json:"xp"`
	Date time.Time `json:"date"`
}

type UserLevelInfo struct {
	Level          int32     `json:"level"`
	Xp             int32     `json:"xp"`
	Streak         int32     `json:"streak"`
	LastActiveDate time.Time `json:"last_active_date"`
	IsLevelUp      bool      `json:"is_level_up"`
}

type CaseLevelResult struct {
	Level    int32 `json:"level"`
	Xp       int32 `json:"xp"`
	XpEarned int32 `json:"xp_earned"`
	LevelUp  bool  `json:"level_up"`
}
