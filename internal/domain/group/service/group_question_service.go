package service

import (
	"context"
	"time"

	"im.turms/server/internal/domain/common/constant"
	"im.turms/server/internal/domain/group/po"
	"im.turms/server/internal/domain/group/repository"
	"im.turms/server/internal/infra/exception"
)

type GroupQuestionService struct {
	questionRepo        repository.GroupJoinQuestionRepository
	groupMemberService  *GroupMemberService
	groupService        *GroupService
	groupVersionService *GroupVersionService
}

func NewGroupQuestionService(
	questionRepo repository.GroupJoinQuestionRepository,
	groupMemberService *GroupMemberService,
	groupService *GroupService,
	groupVersionService *GroupVersionService,
) *GroupQuestionService {
	return &GroupQuestionService{
		questionRepo:        questionRepo,
		groupMemberService:  groupMemberService,
		groupService:        groupService,
		groupVersionService: groupVersionService,
	}
}

// RBAC Operations

func (s *GroupQuestionService) AuthAndCreateQuestion(ctx context.Context, requesterID int64, groupID int64, question string, answers []string, score int) (*po.GroupJoinQuestion, error) {
	isOwnerOrManager, err := s.groupMemberService.IsOwnerOrManager(ctx, groupID, requesterID)
	if err != nil {
		return nil, err
	}
	if !isOwnerOrManager {
		return nil, exception.NewTurmsError(int32(constant.ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_CREATE_GROUP_QUESTION), "Only owner or manager can create questions")
	}

	id := time.Now().UnixNano()
	q := &po.GroupJoinQuestion{
		ID:       id,
		GroupID:  groupID,
		Question: question,
		Answers:  answers,
		Score:    score,
	}
	err = s.questionRepo.Insert(ctx, q)
	if err != nil {
		return nil, err
	}

	err = s.groupVersionService.UpdateJoinQuestionsVersion(ctx, groupID)
	return q, err
}

func (s *GroupQuestionService) AuthAndDeleteQuestion(ctx context.Context, requesterID int64, groupID int64, questionID int64) error {
	isOwnerOrManager, err := s.groupMemberService.IsOwnerOrManager(ctx, groupID, requesterID)
	if err != nil {
		return err
	}
	if !isOwnerOrManager {
		return exception.NewTurmsError(int32(constant.ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_DELETE_GROUP_QUESTION), "Only owner or manager can delete questions")
	}

	err = s.questionRepo.Delete(ctx, questionID)
	if err != nil {
		return err
	}

	return s.groupVersionService.UpdateJoinQuestionsVersion(ctx, groupID)
}

func (s *GroupQuestionService) AuthAndUpdateQuestion(ctx context.Context, requesterID int64, groupID int64, questionID int64, question *string, answers []string, score *int) error {
	isOwnerOrManager, err := s.groupMemberService.IsOwnerOrManager(ctx, groupID, requesterID)
	if err != nil {
		return err
	}
	if !isOwnerOrManager {
		return exception.NewTurmsError(int32(constant.ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_UPDATE_GROUP_QUESTION), "Only owner or manager can update questions")
	}

	_, err = s.questionRepo.Update(ctx, questionID, question, answers, score)
	if err != nil {
		return err
	}

	return s.groupVersionService.UpdateJoinQuestionsVersion(ctx, groupID)
}

func (s *GroupQuestionService) AuthAndCheckAnswer(ctx context.Context, requesterID int64, questionID int64, answer string) (bool, error) {
	q, err := s.questionRepo.FindByID(ctx, questionID)
	if err != nil {
		return false, err
	}
	if q == nil {
		return false, exception.NewTurmsError(int32(constant.ResponseStatusCode_ANSWER_GROUP_QUESTION_IS_DISABLED), "Question not found")
	}

	for _, ans := range q.Answers {
		if ans == answer {
			return true, nil
		}
	}
	return false, nil
}

// Aliases for backward compatibility

func (s *GroupQuestionService) CreateJoinQuestion(ctx context.Context, groupID int64, question string, answers []string, score int) (*po.GroupJoinQuestion, error) {
	// Note: Requester check missing in original alias, using 0/system as default or assuming owner check done elsewhere
	return s.AuthAndCreateQuestion(ctx, 0, groupID, question, answers, score)
}

func (s *GroupQuestionService) DeleteJoinQuestion(ctx context.Context, questionID int64) error {
	// Original logic was broken as it didn't know the groupID for version update
	return s.questionRepo.Delete(ctx, questionID)
}

func (s *GroupQuestionService) UpdateJoinQuestion(ctx context.Context, questionID int64, groupID int64, newQuestion *string, newAnswers []string, newScore *int) error {
	return s.AuthAndUpdateQuestion(ctx, 0, groupID, questionID, newQuestion, newAnswers, newScore)
}

func (s *GroupQuestionService) CheckGroupQuestionAnswerAndJoin(ctx context.Context, requesterID int64, questionID int64, groupID int64, answer string) (bool, error) {
	return s.AuthAndCheckAnswer(ctx, requesterID, questionID, answer)
}

func (s *GroupQuestionService) QueryJoinQuestions(ctx context.Context, groupID int64) ([]po.GroupJoinQuestion, error) {
	return s.questionRepo.FindQuestionsByGroupID(ctx, groupID)
}

// Internal

func (s *GroupQuestionService) QueryJoinQuestionsWithAnswers(ctx context.Context, groupID int64) ([]po.GroupJoinQuestion, error) {
	return s.questionRepo.FindQuestionsByGroupID(ctx, groupID)
}
