package service

import (
	"context"
	"time"

	"im.turms/server/internal/domain/common/constant"
	"im.turms/server/internal/domain/group/po"
	"im.turms/server/internal/domain/group/repository"
	"im.turms/server/internal/infra/exception"
	"im.turms/server/pkg/protocol"
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

// CheckNewGroupQuestion maps to checkNewGroupQuestion(NewGroupQuestion question)
func (s *GroupQuestionService) CheckNewGroupQuestion(answers []string, score int) error {
	if len(answers) == 0 {
		return exception.NewTurmsError(int32(constant.ResponseStatusCode_ILLEGAL_ARGUMENT), "The answers must not be empty")
	}
	if score < 0 {
		return exception.NewTurmsError(int32(constant.ResponseStatusCode_ILLEGAL_ARGUMENT), "The score must be greater than or equal to 0")
	}
	return nil
}

// CheckQuestionIdAndAnswer maps to checkQuestionIdAndAnswer(Map.Entry<Long, String> questionIdAndAnswer)
func (s *GroupQuestionService) CheckQuestionIdAndAnswer(questionID *int64, answer *string) error {
	if questionID == nil || answer == nil {
		return exception.NewTurmsError(int32(constant.ResponseStatusCode_ILLEGAL_ARGUMENT), "The question ID and answer must not be null")
	}
	return nil
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

// AuthAndDeleteGroupJoinQuestions performs batched deletion of join questions with authorization.
// @MappedFrom authAndDeleteGroupJoinQuestions(@NotNull Long userId, @NotNull Long groupId, @NotEmpty Set<Long> questionIds)
func (s *GroupQuestionService) AuthAndDeleteGroupJoinQuestions(ctx context.Context, requesterID int64, groupID int64, questionIDs []int64) error {
	if len(questionIDs) == 0 {
		return nil
	}
	isOwnerOrManager, err := s.groupMemberService.IsOwnerOrManager(ctx, groupID, requesterID)
	if err != nil {
		return err
	}
	if !isOwnerOrManager {
		return exception.NewTurmsError(int32(constant.ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_DELETE_GROUP_QUESTION), "Only owner or manager can delete questions")
	}

	for _, qid := range questionIDs {
		if err := s.questionRepo.Delete(ctx, qid); err != nil {
			return err
		}
	}
	return s.groupVersionService.UpdateJoinQuestionsVersion(ctx, groupID)
}

// AuthAndUpdateGroupJoinQuestion updates a join question with authorization, looking up groupId from the question itself.
// @MappedFrom authAndUpdateGroupJoinQuestion(@NotNull Long requesterId, @NotNull Long questionId, @Nullable String question, @Nullable Set<String> answers, @Nullable Integer score)
func (s *GroupQuestionService) AuthAndUpdateGroupJoinQuestion(ctx context.Context, requesterID int64, questionID int64, question *string, answers []string, score *int) error {
	// Look up the question's groupId for authorization
	groupID, err := s.questionRepo.FindGroupId(ctx, questionID)
	if err != nil {
		return err
	}
	if groupID == nil {
		return exception.NewTurmsError(int32(constant.ResponseStatusCode_UPDATE_INFO_OF_NONEXISTENT_GROUP), "Question not found")
	}
	return s.AuthAndUpdateQuestion(ctx, requesterID, *groupID, questionID, question, answers, score)
}

func (s *GroupQuestionService) UpdateJoinQuestion(ctx context.Context, questionID int64, groupID int64, newQuestion *string, newAnswers []string, newScore *int) error {
	return s.AuthAndUpdateQuestion(ctx, 0, groupID, questionID, newQuestion, newAnswers, newScore)
}

func (s *GroupQuestionService) CheckGroupJoinQuestionsAnswersAndJoin(ctx context.Context, requesterID int64, questionIdToAnswer map[int64]string) (*protocol.GroupJoinQuestionsAnswerResult, error) {
	if len(questionIdToAnswer) == 0 {
		return nil, exception.NewTurmsError(int32(constant.ResponseStatusCode_ILLEGAL_ARGUMENT), "The questions answers shouldn't be empty")
	}

	var questionIDs []int64
	for id := range questionIdToAnswer {
		questionIDs = append(questionIDs, id)
	}

	questions, err := s.questionRepo.FindQuestions(ctx, questionIDs, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	if len(questions) != len(questionIDs) {
		return nil, exception.NewTurmsError(int32(constant.ResponseStatusCode_ANSWER_GROUP_QUESTION_IS_DISABLED), "Question not found or disabled")
	}

	var groupID *int64
	totalScore := 0
	for _, q := range questions {
		if groupID == nil {
			groupID = &q.GroupID
		} else if *groupID != q.GroupID {
			return nil, exception.NewTurmsError(int32(constant.ResponseStatusCode_ILLEGAL_ARGUMENT), "The questions should belong to the same group")
		}

		userAnswer := questionIdToAnswer[q.ID]
		isCorrect := false
		for _, ans := range q.Answers {
			if ans == userAnswer {
				isCorrect = true
				break
			}
		}

		if isCorrect {
			totalScore += q.Score
		}
	}

	if groupID == nil {
		return nil, exception.NewTurmsError(int32(constant.ResponseStatusCode_ANSWER_GROUP_QUESTION_IS_DISABLED), "Question not found or disabled")
	}

	minimumScore, err := s.groupService.QueryGroupMinimumScoreIfActiveAndNotDeleted(ctx, *groupID)
	if err != nil {
		return nil, err
	}
	if minimumScore == nil {
		return nil, exception.NewTurmsError(int32(constant.ResponseStatusCode_UPDATE_INFO_OF_NONEXISTENT_GROUP), "Group does not exist or has been deleted")
	}

	joined := false
	if int32(totalScore) >= *minimumScore {
		// Add user to group
		err = s.groupMemberService.AddGroupMember(ctx, *groupID, requesterID, protocol.GroupMemberRole_MEMBER, nil, nil)
		if err != nil {
			if exception.IsCode(err, int32(constant.ResponseStatusCode_USER_ALREADY_GROUP_MEMBER)) {
				joined = true
			} else {
				return nil, err
			}
		} else {
			joined = true
		}
	}

	return &protocol.GroupJoinQuestionsAnswerResult{
		Score:       int32(totalScore),
		QuestionIds: questionIDs,
		Joined:      joined,
	}, nil
}

func (s *GroupQuestionService) QueryJoinQuestions(ctx context.Context, groupID int64) ([]po.GroupJoinQuestion, error) {
	return s.questionRepo.FindQuestionsByGroupID(ctx, groupID)
}

// Internal

func (s *GroupQuestionService) AuthAndQueryGroupJoinQuestionsWithVersion(ctx context.Context, requesterID int64, groupID int64, withAnswers bool, lastUpdatedDate *time.Time) (*po.GroupJoinQuestionsWithVersion, error) {
	groupTypeID, err := s.groupService.QueryGroupTypeIdIfActiveAndNotDeleted(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if groupTypeID == nil {
		return nil, exception.NewTurmsError(int32(constant.ResponseStatusCode_UPDATE_INFO_OF_NONEXISTENT_GROUP), "Group does not exist or has been deleted")
	}

	version, err := s.groupVersionService.QueryGroupJoinQuestionsVersion(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if lastUpdatedDate != nil && version != nil && !version.After(*lastUpdatedDate) {
		ts := version.UnixMilli()
		lastUpdatedDateInt := &ts
		return &po.GroupJoinQuestionsWithVersion{
			LastUpdatedDate: lastUpdatedDateInt,
		}, nil
	}

	questions, err := s.questionRepo.FindQuestionsByGroupID(ctx, groupID)
	if err != nil {
		return nil, err
	}

	questionPtrs := make([]*po.GroupJoinQuestion, len(questions))
	for i := range questions {
		if !withAnswers {
			questions[i].Answers = nil
		}
		questionPtrs[i] = &questions[i]
	}

	var lastUpdatedDateInt *int64
	if version != nil {
		ts := version.UnixMilli()
		lastUpdatedDateInt = &ts
	}
	return &po.GroupJoinQuestionsWithVersion{
		JoinQuestions:   questionPtrs,
		LastUpdatedDate: lastUpdatedDateInt,
	}, nil
}

func (s *GroupQuestionService) QueryGroupId(ctx context.Context, questionID int64) (*int64, error) {
	return s.questionRepo.FindGroupId(ctx, questionID)
}

func (s *GroupQuestionService) QueryJoinQuestionsWithAnswers(ctx context.Context, groupID int64) ([]po.GroupJoinQuestion, error) {
	return s.questionRepo.FindQuestionsByGroupID(ctx, groupID)
}

func (s *GroupQuestionService) CheckQuestionAnswerAndGetScore(ctx context.Context, questionId int64, answer string, groupID *int64) (int, error) {
	q, err := s.questionRepo.FindByID(ctx, questionId)
	if err != nil {
		return 0, err
	}
	if q == nil || (groupID != nil && q.GroupID != *groupID) {
		return 0, exception.NewTurmsError(int32(constant.ResponseStatusCode_ANSWER_GROUP_QUESTION_IS_DISABLED), "Question not found or disabled")
	}

	for _, ans := range q.Answers {
		if ans == answer {
			return q.Score, nil
		}
	}
	return 0, nil
}

func (s *GroupQuestionService) CountQuestions(ctx context.Context, ids []int64, groupIds []int64) (int64, error) {
	return s.questionRepo.CountQuestions(ctx, ids, groupIds)
}

func (s *GroupQuestionService) FindQuestions(ctx context.Context, ids []int64, groupIds []int64, page *int, size *int, withAnswers bool) ([]po.GroupJoinQuestion, error) {
	questions, err := s.questionRepo.FindQuestions(ctx, ids, groupIds, page, size)
	if err != nil {
		return nil, err
	}
	if !withAnswers {
		for i := range questions {
			questions[i].Answers = nil
		}
	}
	return questions, nil
}

func (s *GroupQuestionService) UpdateQuestions(ctx context.Context, ids []int64, groupID *int64, question *string, answers []string, score *int) error {
	var gidsToUpdate []int64
	if groupID != nil {
		gidsToUpdate = append(gidsToUpdate, *groupID)
	} else if len(ids) > 0 {
		qs, err := s.questionRepo.FindQuestions(ctx, ids, nil, nil, nil)
		if err != nil {
			return err
		}
		seen := make(map[int64]bool)
		for _, q := range qs {
			if !seen[q.GroupID] {
				seen[q.GroupID] = true
				gidsToUpdate = append(gidsToUpdate, q.GroupID)
			}
		}
	}

	err := s.questionRepo.UpdateQuestions(ctx, ids, groupID, question, answers, score)
	if err != nil {
		return err
	}

	for _, gid := range gidsToUpdate {
		_ = s.groupVersionService.UpdateJoinQuestionsVersion(ctx, gid)
	}
	return nil
}

// CreateGroupJoinQuestion creates a question directly (admin-level, no ownership check).
func (s *GroupQuestionService) CreateGroupJoinQuestion(ctx context.Context, id int64, groupID int64, question string, answers []string, score int) (*po.GroupJoinQuestion, error) {
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

// UpdateQuestionsNoVersion updates questions without updating the group version
// (matching Java behavior where updateGroupJoinQuestions does NOT update the version)
func (s *GroupQuestionService) UpdateQuestionsNoVersion(ctx context.Context, ids []int64, groupID *int64, question *string, answers []string, score *int) error {
	return s.questionRepo.UpdateQuestions(ctx, ids, groupID, question, answers, score)
}

// DeleteQuestions performs a batch delete of questions by IDs
func (s *GroupQuestionService) DeleteQuestions(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return s.questionRepo.DeleteByIds(ctx, ids)
}
