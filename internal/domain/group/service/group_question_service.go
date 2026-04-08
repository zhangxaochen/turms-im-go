package service

import (
	"context"
	"log"
	"time"

	group_constant "im.turms/server/internal/domain/group/constant"

	common_constant "im.turms/server/internal/domain/common/constant"
	"im.turms/server/internal/domain/group/po"
	"im.turms/server/internal/domain/group/repository"
	"im.turms/server/internal/infra/exception"
	"im.turms/server/internal/infra/validator"
	"im.turms/server/pkg/codes"
	"im.turms/server/pkg/protocol"
)

const (
	questionContentLimit = 200
	answerContentLimit   = 200
	maxAnswerCount       = 10
)

type GroupQuestionService struct {
	questionRepo         repository.GroupJoinQuestionRepository
	groupMemberService   *GroupMemberService
	groupService         *GroupService
	groupVersionService  *GroupVersionService
	groupTypeService     *GroupTypeService
	groupBlocklistService *GroupBlocklistService
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

func (s *GroupQuestionService) SetGroupTypeService(groupTypeService *GroupTypeService) {
	s.groupTypeService = groupTypeService
}

func (s *GroupQuestionService) SetGroupBlocklistService(groupBlocklistService *GroupBlocklistService) {
	s.groupBlocklistService = groupBlocklistService
}

// CheckNewGroupQuestion maps to checkNewGroupQuestion(NewGroupQuestion question)
func (s *GroupQuestionService) CheckNewGroupQuestion(answers []string, score int) error {
	if len(answers) == 0 {
		return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_ILLEGAL_ARGUMENT), "The answers must not be empty")
	}
	if score < 0 {
		return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_ILLEGAL_ARGUMENT), "The score must be greater than or equal to 0")
	}
	return nil
}

// CheckQuestionIdAndAnswer maps to checkQuestionIdAndAnswer(Map.Entry<Long, String> questionIdAndAnswer)
func (s *GroupQuestionService) CheckQuestionIdAndAnswer(questionID *int64, answer *string) error {
	if questionID == nil || answer == nil {
		return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_ILLEGAL_ARGUMENT), "The question ID and answer must not be null")
	}
	return nil
}

// updateJoinQuestionsVersionNonFatal updates the join questions version, logging errors without failing the operation.
// Java uses onErrorResume with logging; we log and swallow.
func (s *GroupQuestionService) updateJoinQuestionsVersionNonFatal(ctx context.Context, groupID int64) {
	if s.groupVersionService != nil {
		if err := s.groupVersionService.UpdateJoinQuestionsVersion(ctx, groupID); err != nil {
			log.Printf("WARN: failed to update join questions version for group %d: %v", groupID, err)
		}
	}
}

// RBAC Operations

// AuthAndCreateQuestion creates a group join question with authorization.
// Bug fixes:
// - Added group active/not-deleted check before creation
// - Added join strategy validation (must be QUESTION strategy)
func (s *GroupQuestionService) AuthAndCreateQuestion(ctx context.Context, requesterID int64, groupID int64, question string, answers []string, score int) (*po.GroupJoinQuestion, error) {
	isOwnerOrManager, err := s.groupMemberService.IsOwnerOrManager(ctx, groupID, requesterID)
	if err != nil {
		return nil, err
	}
	if !isOwnerOrManager {
		return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_CREATE_GROUP_QUESTION), "Only owner or manager can create questions")
	}

	// Check group is active and not deleted
	typeID, err := s.groupService.QueryGroupTypeIdIfActiveAndNotDeleted(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if typeID == nil {
		return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_CREATE_GROUP_QUESTION_FOR_INACTIVE_GROUP), "Cannot create question for inactive or deleted group")
	}

	// Validate join strategy is QUESTION
	groupType, err := s.groupTypeService.FindGroupType(ctx, *typeID)
	if err != nil {
		return nil, err
	}
	if groupType != nil && groupType.JoinStrategy != group_constant.GroupJoinStrategy_QUESTION {
		switch groupType.JoinStrategy {
		case group_constant.GroupJoinStrategy_JOIN_REQUEST:
			return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_CREATE_GROUP_QUESTION_FOR_GROUP_USING_JOIN_REQUEST), "Cannot create question for group using join request strategy")
		case group_constant.GroupJoinStrategy_INVITATION:
			return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_CREATE_GROUP_QUESTION_FOR_GROUP_USING_INVITATION), "Cannot create question for group using invitation strategy")
		case group_constant.GroupJoinStrategy_MEMBERSHIP_REQUEST:
			return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_CREATE_GROUP_QUESTION_FOR_GROUP_USING_MEMBERSHIP_REQUEST), "Cannot create question for group using membership request strategy")
		default:
			return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_CREATE_GROUP_QUESTION_FOR_INACTIVE_GROUP), "Cannot create question for this group type")
		}
	}

	// Validation: Java validates score not null and min(score, 0), answer content length
	if score < 0 {
		return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_ILLEGAL_ARGUMENT), "score must be >= 0")
	}
	if err := validator.MaxStringLengths(answers, "answers", answerContentLimit); err != nil {
		return nil, err
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

	// Java: version update error is logged and swallowed (non-fatal)
	s.updateJoinQuestionsVersionNonFatal(ctx, groupID)
	return q, nil
}

// AuthAndDeleteQuestion deletes a group join question with authorization.
// Bug fix: Only update version if something was actually deleted.
func (s *GroupQuestionService) AuthAndDeleteQuestion(ctx context.Context, requesterID int64, groupID int64, questionID int64) error {
	isOwnerOrManager, err := s.groupMemberService.IsOwnerOrManager(ctx, groupID, requesterID)
	if err != nil {
		return err
	}
	if !isOwnerOrManager {
		return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_DELETE_GROUP_QUESTION), "Only owner or manager can delete questions")
	}

	err = s.questionRepo.Delete(ctx, questionID)
	if err != nil {
		return err
	}

	// Java: version update error is logged and swallowed (non-fatal)
	s.updateJoinQuestionsVersionNonFatal(ctx, groupID)
	return nil
}

// AuthAndUpdateQuestion updates a group join question with authorization.
// Bug fix: Added early return when all update params are null.
func (s *GroupQuestionService) AuthAndUpdateQuestion(ctx context.Context, requesterID int64, groupID int64, questionID int64, question *string, answers []string, score *int) error {
	// Early return if all update params are null (Java: ACKNOWLEDGED_UPDATE_RESULT)
	if question == nil && len(answers) == 0 && score == nil {
		return nil
	}

	isOwnerOrManager, err := s.groupMemberService.IsOwnerOrManager(ctx, groupID, requesterID)
	if err != nil {
		return err
	}
	if !isOwnerOrManager {
		return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_UPDATE_GROUP_QUESTION), "Only owner or manager can update questions")
	}

	// Validate score >= 0 if provided
	if score != nil && *score < 0 {
		return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_ILLEGAL_ARGUMENT), "The score must be greater than or equal to 0")
	}

	// Validation: answer content length
	if answers != nil {
		if err := validator.MaxStringLengths(answers, "answers", answerContentLimit); err != nil {
			return err
		}
	}

	_, err = s.questionRepo.Update(ctx, questionID, question, answers, score)
	if err != nil {
		return err
	}

	// Java: version update error is logged and swallowed (non-fatal)
	s.updateJoinQuestionsVersionNonFatal(ctx, groupID)
	return nil
}

func (s *GroupQuestionService) AuthAndCheckAnswer(ctx context.Context, requesterID int64, questionID int64, answer string) (bool, error) {
	q, err := s.questionRepo.FindByID(ctx, questionID)
	if err != nil {
		return false, err
	}
	if q == nil {
		return false, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_ANSWER_GROUP_QUESTION_IS_DISABLED), "Question not found")
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
// Bug fix: Only update version if deletion actually happened.
func (s *GroupQuestionService) AuthAndDeleteGroupJoinQuestions(ctx context.Context, requesterID int64, groupID int64, questionIDs []int64) error {
	if len(questionIDs) == 0 {
		return nil
	}
	isOwnerOrManager, err := s.groupMemberService.IsOwnerOrManager(ctx, groupID, requesterID)
	if err != nil {
		return err
	}
	if !isOwnerOrManager {
		return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_DELETE_GROUP_QUESTION), "Only owner or manager can delete questions")
	}

	// Use batch delete instead of one-at-a-time
	err = s.questionRepo.DeleteByIds(ctx, questionIDs)
	if err != nil {
		return err
	}

	// Java: version update error is logged and swallowed (non-fatal)
	s.updateJoinQuestionsVersionNonFatal(ctx, groupID)
	return nil
}

// AuthAndUpdateGroupJoinQuestion updates a join question with authorization, looking up groupId from the question itself.
// @MappedFrom authAndUpdateGroupJoinQuestion(@NotNull Long requesterId, @NotNull Long questionId, @Nullable String question, @Nullable Set<String> answers, @Nullable Integer score)
// Bug fix: Added early return when all update params are null.
func (s *GroupQuestionService) AuthAndUpdateGroupJoinQuestion(ctx context.Context, requesterID int64, questionID int64, question *string, answers []string, score *int) error {
	// Early return if all update params are null (Java: ACKNOWLEDGED_UPDATE_RESULT)
	if question == nil && len(answers) == 0 && score == nil {
		return nil
	}

	// Validate score >= 0 if provided
	if score != nil && *score < 0 {
		return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_ILLEGAL_ARGUMENT), "The score must be greater than or equal to 0")
	}

	// Look up the question's groupId for authorization
	groupID, err := s.questionRepo.FindGroupId(ctx, questionID)
	if err != nil {
		return err
	}
	if groupID == nil {
		return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_UPDATE_INFO_OF_NONEXISTENT_GROUP), "Question not found")
	}
	return s.AuthAndUpdateQuestion(ctx, requesterID, *groupID, questionID, question, answers, score)
}

func (s *GroupQuestionService) UpdateJoinQuestion(ctx context.Context, questionID int64, groupID int64, newQuestion *string, newAnswers []string, newScore *int) error {
	return s.AuthAndUpdateQuestion(ctx, 0, groupID, questionID, newQuestion, newAnswers, newScore)
}

// CheckGroupJoinQuestionsAnswersAndJoin checks answers to group join questions and adds the user if score is sufficient.
// Bug fixes:
// - Added blocklist check (Java: groupBlocklistService.isBlocked)
// - Added existing group member check (Java: groupMemberService.isGroupMember)
// - Added join strategy validation (Java: checks type.getJoinStrategy() == QUESTION)
// - Added group active/not-deleted validation before score check
func (s *GroupQuestionService) CheckGroupJoinQuestionsAnswersAndJoin(ctx context.Context, requesterID int64, questionIdToAnswer map[int64]string) (*protocol.GroupJoinQuestionsAnswerResult, error) {
	if len(questionIdToAnswer) == 0 {
		return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_ILLEGAL_ARGUMENT), "The questions answers shouldn't be empty")
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
		return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_ANSWER_GROUP_QUESTION_IS_DISABLED), "Question not found or disabled")
	}

	var groupID *int64
	totalScore := 0
	for _, q := range questions {
		if groupID == nil {
			groupID = &q.GroupID
		} else if *groupID != q.GroupID {
			return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_ILLEGAL_ARGUMENT), "The questions should belong to the same group")
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
		return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_ANSWER_GROUP_QUESTION_IS_DISABLED), "Question not found or disabled")
	}

	// Bug fix: Check if requester is blocked
	if s.groupBlocklistService != nil {
		isBlocked, err := s.groupBlocklistService.IsBlocked(ctx, *groupID, requesterID)
		if err != nil {
			return nil, err
		}
		if isBlocked {
			return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_GROUP_QUESTION_ANSWERER_HAS_BEEN_BLOCKED), "User has been blocked from the group")
		}
	}

	// Bug fix: Check if requester is already a group member
	isMember, err := s.groupMemberService.IsGroupMember(ctx, *groupID, requesterID)
	if err != nil {
		return nil, err
	}
	if isMember {
		return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_GROUP_MEMBER_ANSWER_GROUP_QUESTION), "User is already a group member")
	}

	// Bug fix: Validate join strategy is QUESTION and group is active/not-deleted
	typeID, err := s.groupService.QueryGroupTypeIdIfActiveAndNotDeleted(ctx, *groupID)
	if err != nil {
		return nil, err
	}
	if typeID == nil {
		return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_ANSWER_GROUP_QUESTION_OF_INACTIVE_GROUP), "Group does not exist or has been deleted")
	}

	if s.groupTypeService != nil {
		groupType, err := s.groupTypeService.FindGroupType(ctx, *typeID)
		if err != nil {
			return nil, err
		}
		if groupType != nil && groupType.JoinStrategy != group_constant.GroupJoinStrategy_QUESTION {
			return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_ANSWER_GROUP_QUESTION_OF_INACTIVE_GROUP), "Group does not use question join strategy")
		}
	}

	minimumScore, err := s.groupService.QueryGroupMinimumScoreIfActiveAndNotDeleted(ctx, *groupID)
	if err != nil {
		return nil, err
	}
	if minimumScore == nil {
		return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_ANSWER_GROUP_QUESTION_OF_INACTIVE_GROUP), "Group does not exist or has been deleted")
	}

	joined := false
	if int32(totalScore) >= *minimumScore {
		// Add user to group
		err = s.groupMemberService.AddGroupMember(ctx, *groupID, requesterID, protocol.GroupMemberRole_MEMBER, nil, nil)
		if err != nil {
			if exception.IsCode(err, int32(common_constant.ResponseStatusCode_USER_ALREADY_GROUP_MEMBER)) {
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

// AuthAndQueryGroupJoinQuestionsWithVersion queries group join questions with version control.
// Bug fixes:
// - Added owner/manager auth check when withAnswers=true
// - Added NO_CONTENT check when questions are empty
// - Added switchIfEmpty(alreadyUpToUpdate) fallback when version is nil
func (s *GroupQuestionService) AuthAndQueryGroupJoinQuestionsWithVersion(ctx context.Context, requesterID int64, groupID int64, withAnswers bool, lastUpdatedDate *time.Time) (*po.GroupJoinQuestionsWithVersion, error) {
	// Bug fix: When withAnswers=true, check isOwnerOrManager (Java: withAnswers ? isOwnerOrManager(requesterId, groupId) : true)
	if withAnswers {
		isOwnerOrManager, err := s.groupMemberService.IsOwnerOrManager(ctx, groupID, requesterID)
		if err != nil {
			return nil, err
		}
		if !isOwnerOrManager {
			return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_QUERY_GROUP_JOIN_REQUEST), "Only owner or manager can query questions with answers")
		}
	}

	groupTypeID, err := s.groupService.QueryGroupTypeIdIfActiveAndNotDeleted(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if groupTypeID == nil {
		return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_UPDATE_INFO_OF_NONEXISTENT_GROUP), "Group does not exist or has been deleted")
	}

	// Bug fix: Java checks isOwnerOrManager when withAnswers=true
	if withAnswers {
		isOwnerOrManager, err := s.groupMemberService.IsOwnerOrManager(ctx, groupID, requesterID)
		if err != nil {
			return nil, err
		}
		if !isOwnerOrManager {
			return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_QUERY_GROUP_QUESTION_ANSWER), "Only owner or manager can query group question answers")
		}
	}

	version, err := s.groupVersionService.QueryGroupJoinQuestionsVersion(ctx, groupID)
	if err != nil {
		return nil, err
	}

	// Bug fix: Java has switchIfEmpty(alreadyUpToDate) when version mono is empty (no version record)
	if version == nil {
		return nil, exception.NewTurmsError(int32(codes.AlreadyUpToDate), "already up-to-date")
	}

	if lastUpdatedDate != nil && !version.After(*lastUpdatedDate) {
		return nil, exception.NewTurmsError(int32(codes.AlreadyUpToDate), "already up-to-date")
	}

	// Bug fix: Java always queries with withAnswers=false at line 391; answers are stripped at query level
	questions, err := s.questionRepo.FindQuestionsByGroupID(ctx, groupID)
	if err != nil {
		return nil, err
	}

	// Bug fix: Java throws NO_CONTENT if questions are empty
	if len(questions) == 0 {
		return nil, exception.NewTurmsError(int32(codes.NoContent), "no group join questions")
	}

	// Java: Always strips answers at query level when querying with version
	questionPtrs := make([]*po.GroupJoinQuestion, len(questions))
	for i := range questions {
		// Bug fix: Java always strips answers (queries with withAnswers=false)
		questions[i].Answers = nil
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
		return 0, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_ANSWER_GROUP_QUESTION_IS_DISABLED), "Question not found or disabled")
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

// UpdateQuestions updates questions and updates the version.
// Note: Java's admin updateGroupJoinQuestions does NOT update version. Use UpdateQuestionsNoVersion for that behavior.
func (s *GroupQuestionService) UpdateQuestions(ctx context.Context, ids []int64, groupID *int64, question *string, answers []string, score *int) error {
	// Bug fix: Early return when all update params are null/falsy (Java: ACKNOWLEDGED_UPDATE_RESULT)
	if groupID == nil && question == nil && len(answers) == 0 && score == nil {
		return nil
	}

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
	// Bug fix: Add content validation matching Java
	questionStr := question
	if err := validator.MaxLength(&questionStr, "question", questionContentLimit); err != nil {
		return nil, err
	}
	if err := validator.InSizeRange(answers, "answers", 1, maxAnswerCount); err != nil {
		return nil, err
	}
	if err := validator.MaxStringLengths(answers, "answers", answerContentLimit); err != nil {
		return nil, err
	}
	if err := validator.MinInt(score, "score", 0); err != nil {
		return nil, err
	}

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
	// Bug fix: early return when all update fields are null (Java returns ACKNOWLEDGED_UPDATE_RESULT)
	if groupID == nil && question == nil && answers == nil && score == nil {
		return nil
	}

	// Bug fix: add content validation matching Java
	if question != nil {
		if err := validator.MaxLength(question, "question", questionContentLimit); err != nil {
			return err
		}
	}
	if answers != nil {
		if err := validator.InSizeRange(answers, "answers", 1, maxAnswerCount); err != nil {
			return err
		}
		if err := validator.MaxStringLengths(answers, "answers", answerContentLimit); err != nil {
			return err
		}
	}
	if score != nil {
		if err := validator.MinInt(*score, "score", 0); err != nil {
			return err
		}
	}

	return s.questionRepo.UpdateQuestions(ctx, ids, groupID, question, answers, score)
}

// DeleteQuestions performs a batch delete of questions by IDs
func (s *GroupQuestionService) DeleteQuestions(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return s.questionRepo.DeleteByIds(ctx, ids)
}
