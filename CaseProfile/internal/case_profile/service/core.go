package service

import (
	"time"

	"github.com/sewaustav/CaseGoProfile/internal/case_profile/models"
)

func (s CaseResultService) updateProfile(actualResults *models.CaseProfile, newResults *models.CaseResult) (*models.CaseProfile, error) {
	return &models.CaseProfile{
		UserID:               actualResults.UserID,
		TotalCases:           actualResults.TotalCases + 1,
		Assertiveness:        s.calcNewRating(actualResults.Assertiveness, newResults.Assertiveness),
		Empathy:              s.calcNewRating(actualResults.Empathy, newResults.Empathy),
		ClarityCommunication: s.calcNewRating(actualResults.ClarityCommunication, newResults.ClarityCommunication),
		Resistance:           s.calcNewRating(actualResults.Resistance, newResults.Resistance),
		Eloquence:            s.calcNewRating(actualResults.Eloquence, newResults.Eloquence),
		Initiative:           s.calcNewRating(actualResults.Initiative, newResults.Initiative),
		ChangedAt:            time.Now(),
	}, nil
}

func (s CaseResultService) calcNewRating(oldRating, newRating float32) float32 {
	return oldRating*(1-s.coefficient) + newRating*s.coefficient
}
