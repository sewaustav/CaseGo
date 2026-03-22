package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sewaustav/CaseGoCore/internal/cases/dto"
	"github.com/sewaustav/CaseGoCore/internal/cases/models"
)

func (s *CaseGoCoreService) StartDialogService(ctx context.Context, caseID int64, user models.UserIdentity) (*models.Case, error) {
	caseModel, err := s.caseGoRepo.GetCaseByID(ctx, caseID)
	if err != nil {
		return nil, err
	}

	_, err = s.dialogRepo.StartDialog(ctx, &models.Dialog{
		UserID: user.UserID,
		CaseID: caseID,
	})

	if err != nil {
		return nil, err
	}

	return caseModel, nil
}

func (s *CaseGoCoreService) HandleInteractionService(ctx context.Context, interaction *dto.InteractionDto, user models.UserIdentity) (*dto.CaseDto, error) {
	dialog, err := s.dialogRepo.GetDialogByID(ctx, interaction.DialogID)
	if err != nil {
		return nil, err
	}

	if dialog.UserID != user.UserID {
		return nil, errors.New("user not authorized to interact with this dialog")
	}

	interactionModel := &models.Interaction{
		DialogID:  interaction.DialogID,
		Step:      interaction.Step,
		Question:  interaction.Question,
		Answer:    interaction.Answer,
		CreatedAt: time.Now(),
	}

	history, err := s.redisClient.GetFullHistory(ctx, interaction.DialogID)
	if err != nil {
		return nil, err
	}

	history = append(history, *interactionModel)

	llmResponse, err := s.llmService.GenerateResponse(ctx, history)
	if err != nil {
		return nil, err
	}

	if err := s.redisClient.Push(ctx, interactionModel); err != nil {
		return nil, err
	}
	return &dto.CaseDto{
		DialogID: interaction.DialogID,
		Question: interaction.Question,
		Model:    llmResponse.Model,
		Step:     new(interaction.Step + 1),
	}, nil
}

func (s *CaseGoCoreService) CompleteDialogService(ctx context.Context, dialogID int64, user models.UserIdentity) (*dto.Result, error) {
	analysis := make(chan *dto.Result, 1)
	errChan := make(chan error, 1)

	dialog, err := s.dialogRepo.GetDialogByID(ctx, dialogID)
	if err != nil {
		return nil, err
	}
	if dialog.UserID != user.UserID {
		return nil, errors.New("user is not owner of dialog")
	}

	tx, err := s.interactionRepo.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	txRepo := s.interactionRepo.WithTx(tx)

	dialogHistory, err := s.redisClient.GetFullHistory(ctx, dialogID)
	if err != nil {
		return nil, err
	}

	go func() {
		result, err := s.llmService.AnalyzeCase(ctx, dialogHistory)
		if err != nil {
			errChan <- err
			return
		}
		analysis <- result
	}()

	for _, interaction := range dialogHistory {
		if err := txRepo.PutInteraction(ctx, &interaction); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	if err := s.redisClient.Clear(ctx, dialogID); err != nil {
		return nil, err
	}

	select {
	case result := <-analysis:
		return result, nil
	case err := <-errChan:
		return nil, fmt.Errorf("analysis failed: %w", err)
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (s *CaseGoCoreService) GetUsersDialogsService(ctx context.Context, user models.UserIdentity, userID int64, limit, offset int) ([]models.Conversation, error) {
	if userID != user.UserID && user.Role != models.Admin {
		return nil, errors.New("only admin or owner can get user dialogs")
	}

	dialogs, err := s.dialogRepo.GetUserDialogs(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	if dialogs == nil || len(dialogs) == 0 {
		if user.Role != models.Admin {
			return nil, errors.New("only admin can get user dialogs")
		}
		return nil, nil
	}

	var conversations []models.Conversation

	for _, dialog := range dialogs {
		history, err := s.redisClient.GetFullHistory(ctx, dialog.ID)
		if err != nil {
			return nil, err
		}
		if history == nil {
			history, err = s.interactionRepo.GetInteractionsByDialogID(ctx, dialog.ID)
		}
		conversations = append(conversations, models.Conversation{
			Dialog:       dialog,
			Interactions: history,
		})
	}
	return conversations, nil

}

func (s *CaseGoCoreService) GetUserDialogByIDService(ctx context.Context, user models.UserIdentity, dialogID int64) (*models.Conversation, error) {
	dialog, err := s.dialogRepo.GetDialogByID(ctx, dialogID)
	if err != nil {
		return nil, err
	}
	if dialog.UserID != user.UserID && user.Role != models.Admin {
		return nil, errors.New("only admin or owner can get user dialogs")
	}

	history, err := s.redisClient.GetFullHistory(ctx, dialogID)
	if err != nil {
		return nil, err
	}
	if history == nil {
		history, err = s.interactionRepo.GetInteractionsByDialogID(ctx, dialogID)
		if err != nil {
			return nil, err
		}
	}
	return &models.Conversation{
		Dialog:       *dialog,
		Interactions: history,
	}, nil

}
