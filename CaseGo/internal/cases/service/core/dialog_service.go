package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/sewaustav/CaseGoCore/apperrors"
	"github.com/sewaustav/CaseGoCore/internal/cases/dto"
	"github.com/sewaustav/CaseGoCore/internal/cases/models"
)

func (s *CaseGoCoreService) StartDialogService(ctx context.Context, caseID int64, user models.UserIdentity) (*models.Case, error) {
	caseModel, err := s.caseGoRepo.GetCaseByID(ctx, caseID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.NewNotFound("case not found", err)
		}
		return nil, apperrors.NewInternal("failed to get case", err)
	}

	_, err = s.dialogRepo.StartDialog(ctx, &models.Dialog{
		UserID: user.UserID,
		CaseID: caseID,
	})
	if err != nil {
		return nil, apperrors.NewInternal("failed to start dialog", err)
	}

	return caseModel, nil
}

func (s *CaseGoCoreService) HandleInteractionService(ctx context.Context, interaction *dto.InteractionDto, user models.UserIdentity) (*dto.CaseDto, error) {
	if interaction == nil {
		return nil, apperrors.NewBadRequest("interaction data is required", nil)
	}

	dialog, err := s.dialogRepo.GetDialogByID(ctx, interaction.DialogID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.NewNotFound("dialog not found", err)
		}
		return nil, apperrors.NewInternal("failed to get dialog", err)
	}

	if dialog.UserID != user.UserID {
		return nil, apperrors.NewForbidden("user is not authorized to interact with this dialog", nil)
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
		return nil, apperrors.NewInternal("failed to get dialog history from cache", err)
	}

	history = append(history, *interactionModel)

	llmResponse, err := s.llmService.GenerateResponse(ctx, history)
	if err != nil {
		return nil, apperrors.NewInternal("failed to generate LLM response", err)
	}

	if err := s.redisClient.Push(ctx, interactionModel); err != nil {
		return nil, apperrors.NewInternal("failed to save interaction to cache", err)
	}

	return &dto.CaseDto{
		DialogID: interaction.DialogID,
		Question: interaction.Question,
		Model:    llmResponse.Model,
		Step:     new(interaction.Step + 1),
	}, nil
}

func (s *CaseGoCoreService) CompleteDialogService(ctx context.Context, dialogID int64, user models.UserIdentity) (*dto.Result, error) {
	dialog, err := s.dialogRepo.GetDialogByID(ctx, dialogID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.NewNotFound("dialog not found", err)
		}
		return nil, apperrors.NewInternal("failed to get dialog", err)
	}

	if dialog.UserID != user.UserID {
		return nil, apperrors.NewForbidden("user is not owner of dialog", nil)
	}

	tx, err := s.interactionRepo.Begin(ctx)
	if err != nil {
		return nil, apperrors.NewInternal("failed to begin transaction", err)
	}
	defer tx.Rollback()

	txRepo := s.interactionRepo.WithTx(tx)

	dialogHistory, err := s.redisClient.GetFullHistory(ctx, dialogID)
	if err != nil {
		return nil, apperrors.NewInternal("failed to get dialog history from cache", err)
	}

	analysis := make(chan *dto.Result, 1)
	errChan := make(chan error, 1)

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
			return nil, apperrors.NewInternal("failed to save interaction", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, apperrors.NewInternal("failed to commit transaction", err)
	}

	if err := s.redisClient.Clear(ctx, dialogID); err != nil {
		return nil, apperrors.NewInternal("failed to clear dialog cache", err)
	}

	select {
	case result := <-analysis:
		go func(r *dto.Result) {
			grpcResult := &models.Result{
				UserID:               user.UserID,
				CaseID:               dialog.CaseID,
				DialogID:             dialogID,
				StepsCount:           r.StepsCount,
				FinishedAt:           time.Now(),
				Assertiveness:        r.SkillsRating.Assertiveness,
				Empathy:              r.SkillsRating.Empathy,
				ClarityCommunication: r.SkillsRating.ClarityCommunication,
				Resistance:           r.SkillsRating.Resistance,
				Eloquence:            r.SkillsRating.Eloquence,
				Initiative:           r.SkillsRating.Initiative,
			}

			if err := s.grpcHandler.SendResults(context.Background(), *grpcResult); err != nil {
				log.Printf("failed to send grpc results: %v", err)
			}
		}(result)

		return result, nil

	case err := <-errChan:
		return nil, apperrors.NewInternal("LLM analysis failed", fmt.Errorf("analyze case: %w", err))

	case <-ctx.Done():
		return nil, apperrors.NewInternal("request timeout", ctx.Err())
	}
}

func (s *CaseGoCoreService) GetUsersDialogsService(ctx context.Context, user models.UserIdentity, userID int64, limit, offset int) ([]models.Conversation, error) {
	if userID != user.UserID && user.Role != models.Admin {
		return nil, apperrors.NewForbidden("only admin or owner can view user dialogs", nil)
	}

	dialogs, err := s.dialogRepo.GetUserDialogs(ctx, userID, limit, offset)
	if err != nil {
		return nil, apperrors.NewInternal("failed to get user dialogs", err)
	}

	if len(dialogs) == 0 {
		return nil, nil
	}

	var conversations []models.Conversation
	for _, dialog := range dialogs {
		history, err := s.redisClient.GetFullHistory(ctx, dialog.ID)
		if err != nil {
			return nil, apperrors.NewInternal("failed to get dialog history from cache", err)
		}

		if history == nil {
			history, err = s.interactionRepo.GetInteractionsByDialogID(ctx, dialog.ID)
			if err != nil {
				return nil, apperrors.NewInternal("failed to get dialog interactions", err)
			}
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
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.NewNotFound("dialog not found", err)
		}
		return nil, apperrors.NewInternal("failed to get dialog", err)
	}

	if dialog.UserID != user.UserID && user.Role != models.Admin {
		return nil, apperrors.NewForbidden("only admin or owner can view this dialog", nil)
	}

	history, err := s.redisClient.GetFullHistory(ctx, dialogID)
	if err != nil {
		return nil, apperrors.NewInternal("failed to get dialog history from cache", err)
	}

	if history == nil {
		history, err = s.interactionRepo.GetInteractionsByDialogID(ctx, dialogID)
		if err != nil {
			return nil, apperrors.NewInternal("failed to get dialog interactions", err)
		}
	}

	return &models.Conversation{
		Dialog:       *dialog,
		Interactions: history,
	}, nil
}
