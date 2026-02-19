package profileService_test

import (
	"context"
	"testing"
	"time"

	dto "github.com/YoungFlores/Case_Go/Profile/internal/profile/dto"
	"github.com/YoungFlores/Case_Go/Profile/internal/profile/models"
	"github.com/YoungFlores/Case_Go/Profile/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	service "github.com/YoungFlores/Case_Go/Profile/internal/profile/service"
)

func ptrString(s string) *string { return &s }
func ptrInt(i int) *int          { return &i }

func TestCreateProfileService(t *testing.T) {
	mockRepo := new(mocks.ProfileRepo)
	mockTx := new(mocks.Tx)
	catRepo := new(mocks.CategoryRepo)
	svc := service.NewProfileService(mockRepo, catRepo)

	ctx := context.Background()
	userID := int64(1)
	userInfo := models.UserIdentity{
		UserID: userID,
		Role:   models.User,
	}

	req := dto.CreateProfileRequest{
		Info: dto.ProfileInfoDTO{
			Avatar:      "https://avatar.com",
			Name:        "Маша",
			Surname:     "Залужная",
			Description: "Создатель орешника",
			City:        ptrString("Moscow"),
			Age:         ptrInt(21),
			Sex:         ptrInt(1),
			Profession:  ptrString("Проектировщик ракет"),
		},
		SocialLinks: []dto.SocialLinkDTO{
			{
				Type: "telegram",
				URL:  "https://t.me/MashaZalushnaya",
			},
		},
		Purposes: []dto.UserPurposeDTO{
			{Purpose: "Донбас"},
		},
	}

	expectedProfile := &models.Profile{
		ID:          0,
		UserID:      0,
		Avatar:      "https://avatar.com",
		IsActive:    true,
		Description: "Создатель орешника",
		Username:    "",
		Name:        "Маша",
		Surname:     "Залужная",
		Patronymic:  nil,
		City:        ptrString("Moscow"),
		Age:         ptrInt(21),
		Sex:         (*models.UserSex)(ptrInt(1)),
		Profession:  ptrString("Проектировщик ракет"),
		CaseCount:   0,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
	}

	expectedPurposes := []models.UserPurpose{
		{ID: 0, Purpose: "Донбас", UserID: userID},
	}

	expectedLinks := []models.UserSocialLink{
		{ID: 0, Type: "telegram", URL: "https://t.me/MashaZalushnaya", UserID: userID},
	}

	mockRepo.On("BeginTx", ctx).Return(mockTx, nil)
	mockRepo.On("WithTx", mockTx).Return(mockRepo)
	mockRepo.On("CreateProfile", ctx, mock.MatchedBy(func(p *models.Profile) bool {
		return p.Name == req.Info.Name
	})).Return(expectedProfile, nil)
	mockRepo.On("AddSocial", ctx, expectedLinks).Return(expectedLinks, nil)
	mockRepo.On("AddPurposes", ctx, expectedPurposes).Return(expectedPurposes, nil)

	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	res, err := svc.CreateProfileService(ctx, req, userInfo)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, expectedProfile.Name, res.UsrProfile.Name)

	mockRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}
