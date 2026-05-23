package profileService_test

import (
	"context"
	"errors"
	"testing"
	"time"

	dto "github.com/YoungFlores/Case_Go/Profile/internal/profile/dto"
	"github.com/YoungFlores/Case_Go/Profile/internal/profile/models"
	repoerr "github.com/YoungFlores/Case_Go/Profile/internal/profile/repository/errors"
	service "github.com/YoungFlores/Case_Go/Profile/internal/profile/service"
	"github.com/YoungFlores/Case_Go/Profile/mocks"
	apperrors "github.com/YoungFlores/Case_Go/Profile/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func ptrString(s string) *string { return &s }
func ptrInt(i int) *int          { return &i }
func ptrInt16(i int16) *int16    { return &i }

// ==================== CreateProfileService Tests ====================

func TestCreateProfileService(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		mockTx := new(mocks.Tx)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userID := int64(1)
		userInfo := models.UserIdentity{
			UserID: userID,
			Role:   models.User,
		}

		req := dto.CreateProfileRequest{
			Info: dto.ProfileInfoDTO{
				Avatar:      "https://avatar.com",
				Username:    "masha_zaluzhnaya",
				Name:        "Маша",
				Surname:     "Залужная",
				Description: "Создатель орешника",
				City:        ptrString("Moscow"),
				Age:         ptrInt(21),
				Sex:         ptrInt(1),
				Profession:  ptrString("Проектировщик ракет"),
			},
			SocialLinks: []dto.SocialLinkDTO{
				{Type: "telegram", URL: "https://t.me/MashaZalushnaya"},
			},
			Purposes: []dto.UserPurposeDTO{
				{Purpose: "Донбас"},
			},
		}

		expectedProfile := &models.Profile{
			ID:          0,
			UserID:      userID,
			Avatar:      "https://avatar.com",
			IsActive:    true,
			Description: "Создатель орешника",
			Username:    "masha_zaluzhnaya",
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

		mockRepo.On("Begin", mock.Anything).Return(mockTx, nil)
		mockRepo.On("WithTx", mockTx).Return(mockRepo)
		mockRepo.On("CreateProfile", ctx, mock.MatchedBy(func(p *models.Profile) bool {
			return p.Name == req.Info.Name &&
				p.Username == req.Info.Username &&
				p.Avatar == req.Info.Avatar &&
				p.UserID == userID &&
				p.IsActive == true
		})).Return(expectedProfile, nil)
		mockRepo.On("AddSocial", ctx, expectedLinks).Return(expectedLinks, nil)
		mockRepo.On("AddPurposes", ctx, expectedPurposes).Return(expectedPurposes, nil)

		mockTx.On("Commit").Return(nil)
		mockTx.On("Rollback").Return(nil)

		res, err := svc.CreateProfileService(ctx, req, userInfo)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, expectedProfile.Name, res.UsrProfile.Name)
		assert.Equal(t, expectedProfile.Username, res.UsrProfile.Username)
		assert.Equal(t, expectedProfile.Avatar, res.UsrProfile.Avatar)
		assert.Equal(t, userID, res.UsrProfile.UserID)
		assert.Equal(t, 1, len(res.UsrPurposes))
		assert.Equal(t, "Донбас", res.UsrPurposes[0].Purpose)
		assert.Equal(t, 1, len(res.UsrSocials))
		assert.Equal(t, "telegram", res.UsrSocials[0].Type)

		mockRepo.AssertExpectations(t)
		mockTx.AssertExpectations(t)
	})

	t.Run("success without sex", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		mockTx := new(mocks.Tx)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userID := int64(1)
		userInfo := models.UserIdentity{UserID: userID, Role: models.User}

		req := dto.CreateProfileRequest{
			Info: dto.ProfileInfoDTO{
				Avatar:   "https://avatar.com",
				Username: "test_user",
				Name:     "Test",
				Surname:  "User",
			},
			Purposes: []dto.UserPurposeDTO{{Purpose: "Purpose"}},
		}

		expectedProfile := &models.Profile{ID: 1, UserID: userID, Name: "Test", Username: "test_user"}
		expectedPurposes := []models.UserPurpose{{ID: 1, Purpose: "Purpose", UserID: userID}}

		mockRepo.On("Begin", mock.Anything).Return(mockTx, nil)
		mockRepo.On("WithTx", mockTx).Return(mockRepo)
		mockRepo.On("CreateProfile", ctx, mock.MatchedBy(func(p *models.Profile) bool {
			return p.Sex == nil && p.Name == "Test"
		})).Return(expectedProfile, nil)
		mockRepo.On("AddSocial", ctx, mock.Anything).Return([]models.UserSocialLink{}, nil)
		mockRepo.On("AddPurposes", ctx, mock.MatchedBy(func(purposes []models.UserPurpose) bool {
			return len(purposes) == 1 && purposes[0].Purpose == "Purpose" && purposes[0].UserID == userID
		})).Return(expectedPurposes, nil)
		mockTx.On("Commit").Return(nil)
		mockTx.On("Rollback").Return(nil)

		res, err := svc.CreateProfileService(ctx, req, userInfo)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "Test", res.UsrProfile.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("begin tx error", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}
		req := dto.CreateProfileRequest{
			Info:     dto.ProfileInfoDTO{Avatar: "a", Username: "u", Name: "n", Surname: "s"},
			Purposes: []dto.UserPurposeDTO{{Purpose: "p"}},
		}

		mockRepo.On("Begin", mock.Anything).Return(nil, errors.New("connection error"))

		res, err := svc.CreateProfileService(ctx, req, userInfo)

		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("create profile error - conflict", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		mockTx := new(mocks.Tx)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}
		req := dto.CreateProfileRequest{
			Info:     dto.ProfileInfoDTO{Avatar: "a", Username: "taken", Name: "n", Surname: "s"},
			Purposes: []dto.UserPurposeDTO{{Purpose: "p"}},
		}

		mockRepo.On("Begin", mock.Anything).Return(mockTx, nil)
		mockRepo.On("WithTx", mockTx).Return(mockRepo)
		mockRepo.On("CreateProfile", ctx, mock.Anything).Return(nil, &repoerr.RepoError{Field: "username", Err: repoerr.ErrConflict})
		mockTx.On("Rollback").Return(nil)

		res, err := svc.CreateProfileService(ctx, req, userInfo)

		assert.Error(t, err)
		assert.Nil(t, res)
		var repoErr *repoerr.RepoError
		assert.True(t, errors.As(err, &repoErr))
		assert.Equal(t, "username", repoErr.Field)
		mockRepo.AssertExpectations(t)
	})

	t.Run("add social error", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		mockTx := new(mocks.Tx)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}
		req := dto.CreateProfileRequest{
			Info:        dto.ProfileInfoDTO{Avatar: "a", Username: "u", Name: "n", Surname: "s"},
			SocialLinks: []dto.SocialLinkDTO{{Type: "tg", URL: "https://t.me/test"}},
			Purposes:    []dto.UserPurposeDTO{{Purpose: "p"}},
		}

		mockRepo.On("Begin", mock.Anything).Return(mockTx, nil)
		mockRepo.On("WithTx", mockTx).Return(mockRepo)
		mockRepo.On("CreateProfile", ctx, mock.Anything).Return(&models.Profile{ID: 1}, nil)
		mockRepo.On("AddSocial", ctx, mock.Anything).Return(nil, errors.New("db error"))
		mockTx.On("Rollback").Return(nil)

		res, err := svc.CreateProfileService(ctx, req, userInfo)

		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("add purposes error", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		mockTx := new(mocks.Tx)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}
		req := dto.CreateProfileRequest{
			Info:     dto.ProfileInfoDTO{Avatar: "a", Username: "u", Name: "n", Surname: "s"},
			Purposes: []dto.UserPurposeDTO{{Purpose: "p"}},
		}

		mockRepo.On("Begin", mock.Anything).Return(mockTx, nil)
		mockRepo.On("WithTx", mockTx).Return(mockRepo)
		mockRepo.On("CreateProfile", ctx, mock.Anything).Return(&models.Profile{ID: 1}, nil)
		mockRepo.On("AddSocial", ctx, mock.Anything).Return([]models.UserSocialLink{}, nil)
		mockRepo.On("AddPurposes", ctx, mock.Anything).Return(nil, errors.New("db error"))
		mockTx.On("Rollback").Return(nil)

		res, err := svc.CreateProfileService(ctx, req, userInfo)

		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("commit error", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		mockTx := new(mocks.Tx)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}
		req := dto.CreateProfileRequest{
			Info:     dto.ProfileInfoDTO{Avatar: "a", Username: "u", Name: "n", Surname: "s"},
			Purposes: []dto.UserPurposeDTO{{Purpose: "p"}},
		}

		mockRepo.On("Begin", mock.Anything).Return(mockTx, nil)
		mockRepo.On("WithTx", mockTx).Return(mockRepo)
		mockRepo.On("CreateProfile", ctx, mock.Anything).Return(&models.Profile{ID: 1}, nil)
		mockRepo.On("AddSocial", ctx, mock.Anything).Return([]models.UserSocialLink{}, nil)
		mockRepo.On("AddPurposes", ctx, mock.Anything).Return([]models.UserPurpose{}, nil)
		mockTx.On("Commit").Return(errors.New("commit error"))
		mockTx.On("Rollback").Return(nil)

		res, err := svc.CreateProfileService(ctx, req, userInfo)

		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})
}

// ==================== UpdateProfileService Tests ====================

func TestUpdateProfileService(t *testing.T) {
	t.Run("success with sex", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userID := int64(1)
		userInfo := models.UserIdentity{UserID: userID, Role: models.User}
		req := dto.ProfileInfoDTO{
			Avatar:      "https://new-avatar.com",
			Username:    "new_username",
			Name:        "NewName",
			Surname:     "NewSurname",
			Description: "New description",
			City:        ptrString("NewCity"),
			Age:         ptrInt(25),
			Sex:         ptrInt(1),
			Profession:  ptrString("Developer"),
		}

		expectedProfile := &models.Profile{
			ID:          1,
			UserID:      userID,
			Avatar:      "https://new-avatar.com",
			Username:    "new_username",
			Name:        "NewName",
			Surname:     "NewSurname",
			Description: "New description",
			City:        ptrString("NewCity"),
			Age:         ptrInt(25),
			Profession:  ptrString("Developer"),
		}

		mockRepo.On("UpdateProfile", ctx, mock.MatchedBy(func(p *models.Profile) bool {
			return p.UserID == userID &&
				p.Username == "new_username" &&
				p.Name == "NewName" &&
				p.Avatar == "https://new-avatar.com" &&
				p.Sex != nil &&
				*p.Sex == models.UserSex(1)
		})).Return(expectedProfile, nil)

		res, err := svc.UpdateProfileService(ctx, userInfo, req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "new_username", res.Username)
		assert.Equal(t, "NewName", res.Name)
		assert.Equal(t, userID, res.UserID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success without sex", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}
		req := dto.ProfileInfoDTO{
			Avatar:   "https://avatar.com",
			Username: "username",
			Name:     "Name",
			Surname:  "Surname",
		}

		mockRepo.On("UpdateProfile", ctx, mock.MatchedBy(func(p *models.Profile) bool {
			return p.Sex == nil
		})).Return(&models.Profile{ID: 1, Name: "Name"}, nil)

		res, err := svc.UpdateProfileService(ctx, userInfo, req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - conflict", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}
		req := dto.ProfileInfoDTO{
			Avatar:   "https://avatar.com",
			Username: "taken_username",
			Name:     "Name",
			Surname:  "Surname",
		}

		mockRepo.On("UpdateProfile", ctx, mock.Anything).Return(nil, &repoerr.RepoError{Field: "username", Err: repoerr.ErrConflict})

		res, err := svc.UpdateProfileService(ctx, userInfo, req)

		assert.Error(t, err)
		assert.Nil(t, res)
		var repoErr *repoerr.RepoError
		assert.True(t, errors.As(err, &repoErr))
		mockRepo.AssertExpectations(t)
	})
}

// ==================== PatchProfileService Tests ====================

func TestPatchProfileService(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userID := int64(1)
		userInfo := models.UserIdentity{UserID: userID, Role: models.User}
		newName := "UpdatedName"
		newAge := 30
		req := dto.UpdateProfilePartialDTO{
			Name: &newName,
			Age:  &newAge,
		}

		expectedProfile := &models.Profile{ID: 1, UserID: userID, Name: newName, Age: &newAge}

		mockRepo.On("PatchProfile", ctx, userID, req).Return(expectedProfile, nil)

		res, err := svc.PatchProfileService(ctx, userInfo, req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, newName, res.Name)
		assert.Equal(t, &newAge, res.Age)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}
		req := dto.UpdateProfilePartialDTO{}

		mockRepo.On("PatchProfile", ctx, int64(1), req).Return(nil, errors.New("db error"))

		res, err := svc.PatchProfileService(ctx, userInfo, req)

		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})
}

// ==================== GetUserProfileService Tests ====================

func TestGetUserProfileService(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userID := int64(1)
		userInfo := models.UserIdentity{UserID: userID, Role: models.User}

		profile := &models.Profile{
			ID:       1,
			UserID:   userID,
			Name:     "Test",
			Username: "test_user",
			IsActive: true,
		}
		purposes := []models.UserPurpose{
			{ID: 1, UserID: userID, Purpose: "Purpose 1"},
			{ID: 2, UserID: userID, Purpose: "Purpose 2"},
		}
		links := []models.UserSocialLink{
			{ID: 1, UserID: userID, Type: "telegram", URL: "https://t.me/test"},
		}

		mockRepo.On("GetUserProfile", ctx, userID).Return(profile, nil)
		mockRepo.On("GetUserPurposes", ctx, userID).Return(purposes, nil)
		mockRepo.On("GetUserSocials", ctx, userID).Return(links, nil)

		res, err := svc.GetUserProfileService(ctx, userInfo)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "Test", res.UsrProfile.Name)
		assert.Equal(t, userID, res.UsrProfile.UserID)
		assert.Equal(t, 2, len(res.UsrPurposes))
		assert.Equal(t, "Purpose 1", res.UsrPurposes[0].Purpose)
		assert.Equal(t, 1, len(res.UsrSocials))
		assert.Equal(t, "telegram", res.UsrSocials[0].Type)
		mockRepo.AssertExpectations(t)
	})

	t.Run("profile not found", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}

		mockRepo.On("GetUserProfile", ctx, int64(1)).Return(nil, errors.New("not found"))

		res, err := svc.GetUserProfileService(ctx, userInfo)

		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("user not active", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}

		profile := &models.Profile{ID: 1, UserID: 1, IsActive: false}
		mockRepo.On("GetUserProfile", ctx, int64(1)).Return(profile, nil)

		res, err := svc.GetUserProfileService(ctx, userInfo)

		assert.Error(t, err)
		assert.Equal(t, apperrors.ErrIsNotActive, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("get purposes error", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}

		profile := &models.Profile{ID: 1, UserID: 1, IsActive: true}
		mockRepo.On("GetUserProfile", ctx, int64(1)).Return(profile, nil)
		mockRepo.On("GetUserPurposes", ctx, int64(1)).Return(nil, errors.New("db error"))

		res, err := svc.GetUserProfileService(ctx, userInfo)

		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("get socials error", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}

		profile := &models.Profile{ID: 1, UserID: 1, IsActive: true}
		mockRepo.On("GetUserProfile", ctx, int64(1)).Return(profile, nil)
		mockRepo.On("GetUserPurposes", ctx, int64(1)).Return([]models.UserPurpose{}, nil)
		mockRepo.On("GetUserSocials", ctx, int64(1)).Return(nil, errors.New("db error"))

		res, err := svc.GetUserProfileService(ctx, userInfo)

		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})
}

// ==================== GetUserProfileByIDService Tests ====================

func TestGetUserProfileByIDService(t *testing.T) {
	t.Run("success - own profile", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userID := int64(1)
		userInfo := models.UserIdentity{UserID: userID, Role: models.User}
		profileID := int64(5)

		profile := &models.Profile{ID: profileID, UserID: userID, Name: "Test"}
		purposes := []models.UserPurpose{{ID: 1, UserID: userID, Purpose: "Test purpose"}}
		links := []models.UserSocialLink{{ID: 1, UserID: userID, Type: "github", URL: "https://github.com/test"}}

		mockRepo.On("GetProfileByID", ctx, profileID).Return(profile, nil)
		mockRepo.On("GetUserPurposes", ctx, userID).Return(purposes, nil)
		mockRepo.On("GetUserSocials", ctx, userID).Return(links, nil)

		res, err := svc.GetUserProfileByIDService(ctx, userInfo, profileID)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, profileID, res.UsrProfile.ID)
		assert.Equal(t, userID, res.UsrProfile.UserID)
		assert.Equal(t, 1, len(res.UsrPurposes))
		assert.Equal(t, 1, len(res.UsrSocials))
		mockRepo.AssertExpectations(t)
	})

	t.Run("success - admin access to other user profile", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		adminID := int64(1)
		otherUserID := int64(2)
		userInfo := models.UserIdentity{UserID: adminID, Role: models.Admin}
		profileID := int64(5)

		profile := &models.Profile{ID: profileID, UserID: otherUserID, Name: "Other User"}
		mockRepo.On("GetProfileByID", ctx, profileID).Return(profile, nil)
		mockRepo.On("GetUserPurposes", ctx, otherUserID).Return([]models.UserPurpose{}, nil)
		mockRepo.On("GetUserSocials", ctx, otherUserID).Return([]models.UserSocialLink{}, nil)

		res, err := svc.GetUserProfileByIDService(ctx, userInfo, profileID)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, otherUserID, res.UsrProfile.UserID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("forbidden - not owner and not admin", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}
		profileID := int64(5)

		profile := &models.Profile{ID: profileID, UserID: 2, Name: "Other User"}
		mockRepo.On("GetProfileByID", ctx, profileID).Return(profile, nil)

		res, err := svc.GetUserProfileByIDService(ctx, userInfo, profileID)

		assert.Error(t, err)
		assert.Equal(t, apperrors.ErrForbidden, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("profile not found", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}

		mockRepo.On("GetProfileByID", ctx, int64(5)).Return(nil, errors.New("not found"))

		res, err := svc.GetUserProfileByIDService(ctx, userInfo, 5)

		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("get purposes error", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}

		profile := &models.Profile{ID: 5, UserID: 1}
		mockRepo.On("GetProfileByID", ctx, int64(5)).Return(profile, nil)
		mockRepo.On("GetUserPurposes", ctx, int64(1)).Return(nil, errors.New("db error"))

		res, err := svc.GetUserProfileByIDService(ctx, userInfo, 5)

		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("get socials error", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}

		profile := &models.Profile{ID: 5, UserID: 1}
		mockRepo.On("GetProfileByID", ctx, int64(5)).Return(profile, nil)
		mockRepo.On("GetUserPurposes", ctx, int64(1)).Return([]models.UserPurpose{}, nil)
		mockRepo.On("GetUserSocials", ctx, int64(1)).Return(nil, errors.New("db error"))

		res, err := svc.GetUserProfileByIDService(ctx, userInfo, 5)

		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})
}

// ==================== UpdatePurposeService Tests ====================

func TestUpdatePurposeService(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		mockTx := new(mocks.Tx)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userID := int64(1)
		userInfo := models.UserIdentity{UserID: userID, Role: models.User}
		req := dto.UserPurposeDTO{Purpose: "Updated Purpose"}
		purposeID := int64(5)

		expectedPurposes := []models.UserPurpose{
			{ID: purposeID, UserID: userID, Purpose: "Updated Purpose"},
		}

		mockRepo.On("Begin", mock.Anything).Return(mockTx, nil)
		mockRepo.On("WithTx", mockTx).Return(mockRepo)
		mockRepo.On("EditPurpose", ctx, mock.MatchedBy(func(p *models.UserPurpose) bool {
			return p.ID == purposeID && p.Purpose == "Updated Purpose" && p.UserID == userID
		})).Return(expectedPurposes, nil)
		mockTx.On("Commit").Return(nil)
		mockTx.On("Rollback").Return(nil)

		res, err := svc.UpdatePurposeService(ctx, req, userInfo, purposeID)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 1, len(res))
		assert.Equal(t, "Updated Purpose", res[0].Purpose)
		assert.Equal(t, purposeID, res[0].ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("begin tx error", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}
		req := dto.UserPurposeDTO{Purpose: "Purpose"}

		mockRepo.On("Begin", mock.Anything).Return(nil, errors.New("connection error"))

		res, err := svc.UpdatePurposeService(ctx, req, userInfo, 5)

		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("edit purpose error", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		mockTx := new(mocks.Tx)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}
		req := dto.UserPurposeDTO{Purpose: "Purpose"}

		mockRepo.On("Begin", mock.Anything).Return(mockTx, nil)
		mockRepo.On("WithTx", mockTx).Return(mockRepo)
		mockRepo.On("EditPurpose", ctx, mock.Anything).Return(nil, errors.New("db error"))
		mockTx.On("Rollback").Return(nil)

		res, err := svc.UpdatePurposeService(ctx, req, userInfo, 5)

		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("commit error", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		mockTx := new(mocks.Tx)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}
		req := dto.UserPurposeDTO{Purpose: "Purpose"}

		mockRepo.On("Begin", mock.Anything).Return(mockTx, nil)
		mockRepo.On("WithTx", mockTx).Return(mockRepo)
		mockRepo.On("EditPurpose", ctx, mock.Anything).Return([]models.UserPurpose{}, nil)
		mockTx.On("Commit").Return(errors.New("commit error"))
		mockTx.On("Rollback").Return(nil)

		res, err := svc.UpdatePurposeService(ctx, req, userInfo, 5)

		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})
}

// ==================== UpdateSocialLinkService Tests ====================

func TestUpdateSocialLinkService(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		mockTx := new(mocks.Tx)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userID := int64(1)
		userInfo := models.UserIdentity{UserID: userID, Role: models.User}
		req := dto.SocialLinkDTO{Type: "twitter", URL: "https://twitter.com/test"}
		linkID := int64(5)

		expectedLinks := []models.UserSocialLink{
			{ID: linkID, UserID: userID, Type: "twitter", URL: "https://twitter.com/test"},
		}

		mockRepo.On("Begin", mock.Anything).Return(mockTx, nil)
		mockRepo.On("WithTx", mockTx).Return(mockRepo)
		mockRepo.On("EditSocial", ctx, mock.MatchedBy(func(l *models.UserSocialLink) bool {
			return l.ID == linkID && l.Type == "twitter" && l.URL == "https://twitter.com/test" && l.UserID == userID
		})).Return(expectedLinks, nil)
		mockTx.On("Commit").Return(nil)
		mockTx.On("Rollback").Return(nil)

		res, err := svc.UpdateSocialLinkService(ctx, req, userInfo, linkID)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 1, len(res))
		assert.Equal(t, "twitter", res[0].Type)
		assert.Equal(t, linkID, res[0].ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("begin tx error", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}
		req := dto.SocialLinkDTO{Type: "tg", URL: "https://t.me/test"}

		mockRepo.On("Begin", mock.Anything).Return(nil, errors.New("connection error"))

		res, err := svc.UpdateSocialLinkService(ctx, req, userInfo, 5)

		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("edit social error", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		mockTx := new(mocks.Tx)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}
		req := dto.SocialLinkDTO{Type: "tg", URL: "https://t.me/test"}

		mockRepo.On("Begin", mock.Anything).Return(mockTx, nil)
		mockRepo.On("WithTx", mockTx).Return(mockRepo)
		mockRepo.On("EditSocial", ctx, mock.Anything).Return(nil, errors.New("db error"))
		mockTx.On("Rollback").Return(nil)

		res, err := svc.UpdateSocialLinkService(ctx, req, userInfo, 5)

		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("commit error", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		mockTx := new(mocks.Tx)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}
		req := dto.SocialLinkDTO{Type: "tg", URL: "https://t.me/test"}

		mockRepo.On("Begin", mock.Anything).Return(mockTx, nil)
		mockRepo.On("WithTx", mockTx).Return(mockRepo)
		mockRepo.On("EditSocial", ctx, mock.Anything).Return([]models.UserSocialLink{}, nil)
		mockTx.On("Commit").Return(errors.New("commit error"))
		mockTx.On("Rollback").Return(nil)

		res, err := svc.UpdateSocialLinkService(ctx, req, userInfo, 5)

		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})
}

// ==================== AddPurposesService Tests ====================

func TestAddPurposesService(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userID := int64(1)
		userInfo := models.UserIdentity{UserID: userID, Role: models.User}
		req := []dto.UserPurposeDTO{
			{Purpose: "Purpose 1"},
			{Purpose: "Purpose 2"},
		}

		expectedPurposes := []models.UserPurpose{
			{ID: 1, UserID: userID, Purpose: "Purpose 1"},
			{ID: 2, UserID: userID, Purpose: "Purpose 2"},
		}

		mockRepo.On("AddPurposes", ctx, mock.MatchedBy(func(purposes []models.UserPurpose) bool {
			return len(purposes) == 2 &&
				purposes[0].Purpose == "Purpose 1" &&
				purposes[0].UserID == userID &&
				purposes[1].Purpose == "Purpose 2" &&
				purposes[1].UserID == userID
		})).Return(expectedPurposes, nil)

		res, err := svc.AddPurposesService(ctx, req, userInfo)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 2, len(res))
		assert.Equal(t, "Purpose 1", res[0].Purpose)
		assert.Equal(t, "Purpose 2", res[1].Purpose)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}
		req := []dto.UserPurposeDTO{{Purpose: "Purpose"}}

		mockRepo.On("AddPurposes", ctx, mock.Anything).Return(nil, errors.New("db error"))

		res, err := svc.AddPurposesService(ctx, req, userInfo)

		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})
}

// ==================== AddSocialLinksService Tests ====================

func TestAddSocialLinksService(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userID := int64(1)
		userInfo := models.UserIdentity{UserID: userID, Role: models.User}
		req := []dto.SocialLinkDTO{
			{Type: "telegram", URL: "https://t.me/test"},
			{Type: "github", URL: "https://github.com/test"},
		}

		expectedLinks := []models.UserSocialLink{
			{ID: 1, UserID: userID, Type: "telegram", URL: "https://t.me/test"},
			{ID: 2, UserID: userID, Type: "github", URL: "https://github.com/test"},
		}

		mockRepo.On("AddSocial", ctx, mock.MatchedBy(func(links []models.UserSocialLink) bool {
			return len(links) == 2 &&
				links[0].Type == "telegram" &&
				links[0].URL == "https://t.me/test" &&
				links[0].UserID == userID &&
				links[1].Type == "github" &&
				links[1].UserID == userID
		})).Return(expectedLinks, nil)

		res, err := svc.AddSocialLinksService(ctx, req, userInfo)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 2, len(res))
		assert.Equal(t, "telegram", res[0].Type)
		assert.Equal(t, "github", res[1].Type)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}
		req := []dto.SocialLinkDTO{{Type: "tg", URL: "https://t.me/test"}}

		mockRepo.On("AddSocial", ctx, mock.Anything).Return(nil, errors.New("db error"))

		res, err := svc.AddSocialLinksService(ctx, req, userInfo)

		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})
}

// ==================== DeletePurposeService Tests ====================

func TestDeletePurposeService(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userID := int64(1)
		userInfo := models.UserIdentity{UserID: userID, Role: models.User}
		purposeID := int64(5)

		mockRepo.On("GetUserByProfileID", ctx, purposeID, userID).Return(userID, nil)
		mockRepo.On("DeletePurpose", ctx, purposeID).Return(nil)

		err := svc.DeletePurposeService(ctx, purposeID, userInfo)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("get owner error", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}

		mockRepo.On("GetUserByProfileID", ctx, int64(5), int64(1)).Return(int64(0), errors.New("not found"))

		err := svc.DeletePurposeService(ctx, 5, userInfo)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("forbidden - not owner", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}

		mockRepo.On("GetUserByProfileID", ctx, int64(5), int64(1)).Return(int64(2), nil)

		err := svc.DeletePurposeService(ctx, 5, userInfo)

		assert.Error(t, err)
		assert.Equal(t, apperrors.ErrForbidden, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("delete error", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}

		mockRepo.On("GetUserByProfileID", ctx, int64(5), int64(1)).Return(int64(1), nil)
		mockRepo.On("DeletePurpose", ctx, int64(5)).Return(errors.New("db error"))

		err := svc.DeletePurposeService(ctx, 5, userInfo)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

// ==================== DeleteLinkService Tests ====================

func TestDeleteLinkService(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userID := int64(1)
		userInfo := models.UserIdentity{UserID: userID, Role: models.User}
		linkID := int64(5)

		mockRepo.On("GetUserByProfileID", ctx, linkID, userID).Return(userID, nil)
		mockRepo.On("DeleteSocial", ctx, linkID).Return(nil)

		err := svc.DeleteLinkService(ctx, linkID, userInfo)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("get owner error", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}

		mockRepo.On("GetUserByProfileID", ctx, int64(5), int64(1)).Return(int64(0), errors.New("not found"))

		err := svc.DeleteLinkService(ctx, 5, userInfo)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("forbidden - not owner", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}

		mockRepo.On("GetUserByProfileID", ctx, int64(5), int64(1)).Return(int64(2), nil)

		err := svc.DeleteLinkService(ctx, 5, userInfo)

		assert.Error(t, err)
		assert.Equal(t, apperrors.ErrForbidden, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("delete error", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}

		mockRepo.On("GetUserByProfileID", ctx, int64(5), int64(1)).Return(int64(1), nil)
		mockRepo.On("DeleteSocial", ctx, int64(5)).Return(errors.New("db error"))

		err := svc.DeleteLinkService(ctx, 5, userInfo)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

// ==================== DeleteProfileService Tests ====================

func TestDeleteProfileService(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userID := int64(1)
		userInfo := models.UserIdentity{UserID: userID, Role: models.User}

		mockRepo.On("DeleteProfile", ctx, userID).Return(nil)

		err := svc.DeleteProfileService(ctx, userInfo)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}

		mockRepo.On("DeleteProfile", ctx, int64(1)).Return(errors.New("db error"))

		err := svc.DeleteProfileService(ctx, userInfo)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

// ==================== DeleteProfileWithoutRecoveryService Tests ====================

func TestDeleteProfileWithoutRecoveryService(t *testing.T) {
	t.Run("success - admin", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.Admin}
		targetUserID := int64(5)

		mockRepo.On("DeleteProfileWithoutRecovery", ctx, targetUserID).Return(nil)

		err := svc.DeleteProfileWithoutRecoveryService(ctx, userInfo, targetUserID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("forbidden - not admin", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}

		err := svc.DeleteProfileWithoutRecoveryService(ctx, userInfo, 5)

		assert.Error(t, err)
		assert.Equal(t, apperrors.ErrForbidden, err)
	})

	t.Run("delete error", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.Admin}

		mockRepo.On("DeleteProfileWithoutRecovery", ctx, int64(5)).Return(errors.New("db error"))

		err := svc.DeleteProfileWithoutRecoveryService(ctx, userInfo, 5)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

// ==================== AddProfessionService Tests ====================

func TestAddProfessionService(t *testing.T) {
	t.Run("success without parent", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		mockTx := new(mocks.Tx)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userID := int64(1)
		userInfo := models.UserIdentity{UserID: userID, Role: models.User}
		req := []dto.ProfessionDTO{
			{ProfessionID: 1},
		}

		expectedProfession := &models.UserProfession{ID: 1, UserID: userID, ProfessionID: 1}

		mockRepo.On("Begin", mock.Anything).Return(mockTx, nil)
		mockRepo.On("WithTx", mockTx).Return(mockRepo)
		mockRepo.On("AddProfession", ctx, mock.MatchedBy(func(p *models.UserProfession) bool {
			return p.UserID == userID && p.ProfessionID == 1
		})).Return(expectedProfession, nil)
		catRepo.On("GetParentOfCategory", ctx, int16(1)).Return(nil, nil)
		mockTx.On("Commit").Return(nil)
		mockTx.On("Rollback").Return(nil)

		res, err := svc.AddProfessionService(ctx, userInfo, req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 1, len(res))
		assert.Equal(t, int16(1), res[0].ProfessionID)
		mockRepo.AssertExpectations(t)
		catRepo.AssertExpectations(t)
	})

	t.Run("success with parent", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		mockTx := new(mocks.Tx)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userID := int64(1)
		userInfo := models.UserIdentity{UserID: userID, Role: models.User}
		req := []dto.ProfessionDTO{
			{ProfessionID: 5},
		}

		childProfession := &models.UserProfession{ID: 1, UserID: userID, ProfessionID: 5}
		parentProfession := &models.UserProfession{ID: 2, UserID: userID, ProfessionID: 1}

		mockRepo.On("Begin", mock.Anything).Return(mockTx, nil)
		mockRepo.On("WithTx", mockTx).Return(mockRepo)
		mockRepo.On("AddProfession", ctx, mock.MatchedBy(func(p *models.UserProfession) bool {
			return p.ProfessionID == 5
		})).Return(childProfession, nil).Once()
		catRepo.On("GetParentOfCategory", ctx, int16(5)).Return(ptrInt16(1), nil)
		mockRepo.On("AddProfession", ctx, mock.MatchedBy(func(p *models.UserProfession) bool {
			return p.ProfessionID == 1
		})).Return(parentProfession, nil).Once()
		mockTx.On("Commit").Return(nil)
		mockTx.On("Rollback").Return(nil)

		res, err := svc.AddProfessionService(ctx, userInfo, req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 2, len(res))
		mockRepo.AssertExpectations(t)
		catRepo.AssertExpectations(t)
	})

	t.Run("begin tx error", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}
		req := []dto.ProfessionDTO{{ProfessionID: 1}}

		mockRepo.On("Begin", mock.Anything).Return(nil, errors.New("connection error"))

		res, err := svc.AddProfessionService(ctx, userInfo, req)

		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("add profession error", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		mockTx := new(mocks.Tx)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}
		req := []dto.ProfessionDTO{{ProfessionID: 1}}

		mockRepo.On("Begin", mock.Anything).Return(mockTx, nil)
		mockRepo.On("WithTx", mockTx).Return(mockRepo)
		mockRepo.On("AddProfession", ctx, mock.Anything).Return(nil, errors.New("db error"))
		mockTx.On("Rollback").Return(nil)

		res, err := svc.AddProfessionService(ctx, userInfo, req)

		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("get parent error", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		mockTx := new(mocks.Tx)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}
		req := []dto.ProfessionDTO{{ProfessionID: 1}}

		mockRepo.On("Begin", mock.Anything).Return(mockTx, nil)
		mockRepo.On("WithTx", mockTx).Return(mockRepo)
		mockRepo.On("AddProfession", ctx, mock.Anything).Return(&models.UserProfession{ID: 1}, nil)
		catRepo.On("GetParentOfCategory", ctx, int16(1)).Return(nil, errors.New("db error"))
		mockTx.On("Rollback").Return(nil)

		res, err := svc.AddProfessionService(ctx, userInfo, req)

		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
		catRepo.AssertExpectations(t)
	})

	t.Run("add parent profession error", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		mockTx := new(mocks.Tx)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}
		req := []dto.ProfessionDTO{{ProfessionID: 5}}

		mockRepo.On("Begin", mock.Anything).Return(mockTx, nil)
		mockRepo.On("WithTx", mockTx).Return(mockRepo)
		mockRepo.On("AddProfession", ctx, mock.MatchedBy(func(p *models.UserProfession) bool {
			return p.ProfessionID == 5
		})).Return(&models.UserProfession{ID: 1, ProfessionID: 5}, nil).Once()
		catRepo.On("GetParentOfCategory", ctx, int16(5)).Return(ptrInt16(1), nil)
		mockRepo.On("AddProfession", ctx, mock.MatchedBy(func(p *models.UserProfession) bool {
			return p.ProfessionID == 1
		})).Return(nil, errors.New("db error")).Once()
		mockTx.On("Rollback").Return(nil)

		res, err := svc.AddProfessionService(ctx, userInfo, req)

		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
		catRepo.AssertExpectations(t)
	})

	t.Run("commit error", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		mockTx := new(mocks.Tx)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}
		req := []dto.ProfessionDTO{{ProfessionID: 1}}

		mockRepo.On("Begin", mock.Anything).Return(mockTx, nil)
		mockRepo.On("WithTx", mockTx).Return(mockRepo)
		mockRepo.On("AddProfession", ctx, mock.Anything).Return(&models.UserProfession{ID: 1}, nil)
		catRepo.On("GetParentOfCategory", ctx, int16(1)).Return(nil, nil)
		mockTx.On("Commit").Return(errors.New("commit error"))
		mockTx.On("Rollback").Return(nil)

		res, err := svc.AddProfessionService(ctx, userInfo, req)

		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})
}

// ==================== EditProfessionCategoryService Tests ====================

func TestEditProfessionCategoryService(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userID := int64(1)
		userInfo := models.UserIdentity{UserID: userID, Role: models.User}
		profession := &models.UserProfession{ID: 5, ProfessionID: 3}

		expectedProfession := &models.UserProfession{ID: 5, UserID: userID, ProfessionID: 3}

		mockRepo.On("GetProfileIDByProfessionID", ctx, int64(5)).Return(userID, nil)
		mockRepo.On("EditProfession", ctx, profession).Return(expectedProfession, nil)

		res, err := svc.EditProfessionCategoryService(ctx, userInfo, profession)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, int64(5), res.ID)
		assert.Equal(t, int16(3), res.ProfessionID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("get owner error", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}
		profession := &models.UserProfession{ID: 5, ProfessionID: 3}

		mockRepo.On("GetProfileIDByProfessionID", ctx, int64(5)).Return(int64(0), errors.New("not found"))

		res, err := svc.EditProfessionCategoryService(ctx, userInfo, profession)

		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("forbidden - not owner", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}
		profession := &models.UserProfession{ID: 5, ProfessionID: 3}

		mockRepo.On("GetProfileIDByProfessionID", ctx, int64(5)).Return(int64(2), nil)

		res, err := svc.EditProfessionCategoryService(ctx, userInfo, profession)

		assert.Error(t, err)
		assert.Equal(t, apperrors.ErrForbidden, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("edit profession error", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}
		profession := &models.UserProfession{ID: 5, ProfessionID: 3}

		mockRepo.On("GetProfileIDByProfessionID", ctx, int64(5)).Return(int64(1), nil)
		mockRepo.On("EditProfession", ctx, profession).Return(nil, errors.New("db error"))

		res, err := svc.EditProfessionCategoryService(ctx, userInfo, profession)

		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})
}

// ==================== DeleteProfessionService Tests ====================

func TestDeleteProfessionService(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userID := int64(1)
		userInfo := models.UserIdentity{UserID: userID, Role: models.User}
		professionID := int64(5)

		mockRepo.On("GetProfileIDByProfessionID", ctx, professionID).Return(userID, nil)
		mockRepo.On("DeleteProfession", ctx, professionID).Return(nil)

		err := svc.DeleteProfessionService(ctx, professionID, userInfo)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("get owner error", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}

		mockRepo.On("GetProfileIDByProfessionID", ctx, int64(5)).Return(int64(0), errors.New("not found"))

		err := svc.DeleteProfessionService(ctx, 5, userInfo)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("forbidden - not owner", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}

		mockRepo.On("GetProfileIDByProfessionID", ctx, int64(5)).Return(int64(2), nil)

		err := svc.DeleteProfessionService(ctx, 5, userInfo)

		assert.Error(t, err)
		assert.Equal(t, apperrors.ErrForbidden, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("delete error", func(t *testing.T) {
		mockRepo := new(mocks.ProfileRepo)
		catRepo := new(mocks.CategoryRepo)
		levelRepo := new(mocks.LevelRepo)
		svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

		ctx := context.Background()
		userInfo := models.UserIdentity{UserID: 1, Role: models.User}

		mockRepo.On("GetProfileIDByProfessionID", ctx, int64(5)).Return(int64(1), nil)
		mockRepo.On("DeleteProfession", ctx, int64(5)).Return(errors.New("db error"))

		err := svc.DeleteProfessionService(ctx, 5, userInfo)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

// ==================== NewProfileService Tests ====================

func TestNewProfileService(t *testing.T) {
	mockRepo := new(mocks.ProfileRepo)
	catRepo := new(mocks.CategoryRepo)
	levelRepo := new(mocks.LevelRepo)

	svc := service.NewProfileService(mockRepo, catRepo, levelRepo)

	assert.NotNil(t, svc)
}
