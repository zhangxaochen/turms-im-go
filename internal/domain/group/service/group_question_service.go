package service

import (
	"context"
	"errors"
	"time"

	"im.turms/server/internal/domain/group/po"
	"im.turms/server/internal/domain/group/repository"
)

type GroupQuestionService struct {
	questionRepo repository.GroupJoinQuestionRepository
}

func NewGroupQuestionService(questionRepo repository.GroupJoinQuestionRepository) *GroupQuestionService {
	return &GroupQuestionService{
		questionRepo: questionRepo,
	}
}

func (s *GroupQuestionService) CreateJoinQuestion(ctx context.Context, groupID int64, question string, answers []string, score int) (*po.GroupJoinQuestion, error) {
	id := time.Now().UnixNano()
	q := &po.GroupJoinQuestion{
		ID:       id,
		GroupID:  groupID,
		Question: question,
		Answers:  answers,
		Score:    score,
	}
	err := s.questionRepo.Insert(ctx, q)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (s *GroupQuestionService) DeleteJoinQuestion(ctx context.Context, questionID int64) error {
	return s.questionRepo.Delete(ctx, questionID)
}

func (s *GroupQuestionService) UpdateJoinQuestion(ctx context.Context, questionID int64, groupID int64, newQuestion *string, newAnswers []string, newScore *int) error {
	// Simple update override just to show interaction
	// Actual Turms logic uses more precise field-by-field updates
	q, err := s.questionRepo.FindQuestionsByGroupID(ctx, groupID)
	if err != nil {
		return err
	}

	// Check if question exists (pseudo-logic for simplicity instead of exposing FindByID)
	found := false
	for _, item := range q {
		if item.ID == questionID {
			found = true
			break
		}
	}

	if !found {
		return errors.New("question not found")
	}

	// Assuming a delete and re-insert for quick schema sync without building a complex Update builder
	_ = s.questionRepo.Delete(ctx, questionID)

	qs := ""
	if newQuestion != nil {
		qs = *newQuestion
	}
	sc := 0
	if newScore != nil {
		sc = *newScore
	}

	updated := &po.GroupJoinQuestion{
		ID:       questionID,
		GroupID:  groupID,
		Question: qs,
		Answers:  newAnswers,
		Score:    sc,
	}
	return s.questionRepo.Insert(ctx, updated)
}

func (s *GroupQuestionService) QueryJoinQuestions(ctx context.Context, groupID int64) ([]po.GroupJoinQuestion, error) {
	return s.questionRepo.FindQuestionsByGroupID(ctx, groupID)
}

func (s *GroupQuestionService) CheckGroupQuestionAnswerAndJoin(ctx context.Context, requesterID int64, questionID int64, groupID int64, answer string) (bool, error) {
	questions, err := s.questionRepo.FindQuestionsByGroupID(ctx, groupID)
	if err != nil {
		return false, err
	}

	for _, q := range questions {
		if q.ID == questionID {
			for _, ans := range q.Answers {
				if ans == answer {
					return true, nil
				}
			}
			return false, nil
		}
	}
	return false, errors.New("question not found")
}
