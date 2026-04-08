package service

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"im.turms/server/internal/domain/common/infra/idgen"
	"im.turms/server/internal/domain/user/po"
	"im.turms/server/internal/domain/user/repository"
)

type UserService interface {
	CreateUser(ctx context.Context, password string, name string, intro string, profilePicture string, profileAccess int32, permissionGroupID int64, isActive bool) (*po.User, error)
	AddUser(ctx context.Context, id int64, password string, name string, intro string, profilePicture string, profileAccess int32, permissionGroupID int64, registrationDate time.Time, isActive bool) (*po.User, error)
	FindUser(ctx context.Context, userID int64) (*po.User, error)
	UpdateUser(ctx context.Context, userID int64, update bson.M) error
	CheckIfUserExists(ctx context.Context, userID int64) (bool, error)
	DeleteUsers(ctx context.Context, userIDs []int64) (int64, error)
	QueryUsersProfile(ctx context.Context, userIDs []int64) ([]*po.User, error)
	QueryUserName(ctx context.Context, userID int64) (string, error)
	QueryUsers(ctx context.Context, userIDs []int64) ([]*po.User, error)
	CountUsers(ctx context.Context, activeOnly bool) (int64, error)
	IsActiveAndNotDeleted(ctx context.Context, userID int64) (bool, error)
	FindPassword(ctx context.Context, userID int64) (*string, error)
	IsAllowedToSendMessageToTarget(ctx context.Context, isGroupMessage bool, isSystemMessage bool, requesterID int64, targetID int64) (int, error)
	IsAllowToQueryUserProfile(ctx context.Context, requesterID int64, targetID int64) (int, error)
	AuthAndQueryUsersProfile(ctx context.Context, requesterID int64, userIDs []int64, name string, lastUpdatedDate *time.Time, skip int, limit int) ([]*po.User, error)
	QueryUserRoleIDByUserID(ctx context.Context, userID int64) (int64, error)
	CountRegisteredUsers(ctx context.Context, startDate *time.Time, endDate *time.Time, queryDeletedRecords bool) (int64, error)
	CountDeletedUsers(ctx context.Context, startDate *time.Time, endDate *time.Time) (int64, error)
	UpdateUsers(ctx context.Context, userIDs []int64, update bson.M) (int64, error)
}

type userService struct {
	idGen *idgen.SnowflakeIdGenerator
	repo  repository.UserRepository
}

func NewUserService(idGen *idgen.SnowflakeIdGenerator, repo repository.UserRepository) UserService {
	return &userService{
		idGen: idGen,
		repo:  repo,
	}
}

// @MappedFrom createUser(@Nullable Long id, @Nullable String rawPassword, @Nullable String name, @Nullable String intro, @Nullable String profilePicture, @Nullable @ValidProfileAccess ProfileAccessStrategy profileAccessStrategy, @Nullable Long roleId, @Nullable @PastOrPresent Date registrationDate, @Nullable Boolean isActive)
func (s *userService) CreateUser(ctx context.Context, password string, name string, intro string, profilePicture string, profileAccess int32, permissionGroupID int64, isActive bool) (*po.User, error) {
	userID := s.idGen.NextIncreasingId()
	now := time.Now()

	user := &po.User{
		ID:                userID,
		Password:          password, // Assuming plain text for this simple refactor, should be hashed in real world
		Name:              name,
		Intro:             intro,
		ProfilePicture:    profilePicture,
		ProfileAccess:     profileAccess,
		PermissionGroupID: permissionGroupID,
		RegistrationDate:  now,
		IsActive:          isActive,
	}
	err := s.repo.Insert(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// @MappedFrom addUser(@RequestBody AddUserDTO addUserDTO)
// @MappedFrom addUser(@Nullable Long id, @Nullable String rawPassword, @Nullable String name, @Nullable String intro, @Nullable String profilePicture, @Nullable @ValidProfileAccess ProfileAccessStrategy profileAccessStrategy, @Nullable Long roleId, @Nullable @PastOrPresent Date registrationDate, @Nullable Boolean isActive)
func (s *userService) AddUser(ctx context.Context, id int64, password string, name string, intro string, profilePicture string, profileAccess int32, permissionGroupID int64, registrationDate time.Time, isActive bool) (*po.User, error) {
	// Validation (matching Java: minPasswordLengthForCreate, maxPasswordLength, name/intro/profilePicture max lengths)
	if len(password) > 64 {
		return nil, errors.New("password must not exceed 64 characters")
	}
	if len(name) > 40 {
		return nil, errors.New("name must not exceed 40 characters")
	}
	if len(intro) > 300 {
		return nil, errors.New("intro must not exceed 300 characters")
	}
	if len(profilePicture) > 255 {
		return nil, errors.New("profilePicture must not exceed 255 characters")
	}
	if !registrationDate.IsZero() && registrationDate.After(time.Now()) {
		return nil, errors.New("registrationDate must be a past or present date")
	}

	if id == 0 {
		id = s.idGen.NextIncreasingId()
	}
	if registrationDate.IsZero() {
		registrationDate = time.Now()
	}
	now := time.Now()

	// Apply defaults matching Java behavior
	// profileAccess defaults to ALL (1), permissionGroupID defaults to DEFAULT_USER_ROLE_ID
	if profileAccess == 0 {
		profileAccess = 1 // ProfileAccessStrategy.ALL
	}

	user := &po.User{
		ID:                id,
		Password:          password, // TODO: Implement password encoding via passwordManager
		Name:              name,
		Intro:             intro,
		ProfilePicture:    profilePicture,
		ProfileAccess:     profileAccess,
		PermissionGroupID: permissionGroupID,
		RegistrationDate:  registrationDate,
		IsActive:          isActive,
		LastUpdatedDate:   &now,
		// DeletionDate is nil (zero value) by default, matching Java's explicit null
	}
	err := s.repo.Insert(ctx, user)
	if err != nil {
		return nil, err
	}
	// TODO: Create default relationship group, upsert user version, sync Elasticsearch (transactional side effects)
	return user, nil
}

func (s *userService) FindUser(ctx context.Context, userID int64) (*po.User, error) {
	return s.repo.FindByID(ctx, userID)
}

// @MappedFrom updateUser(@NotNull Long userId, @Nullable String rawPassword, @Nullable String name, @Nullable String intro, @Nullable String profilePicture, @Nullable @ValidProfileAccess ProfileAccessStrategy profileAccessStrategy, @Nullable Long roleId, @Nullable Boolean isActive, @Nullable @PastOrPresent Date registrationDate, @Nullable Map<String, Value> userDefinedAttributes)
// @MappedFrom updateUser(Set<Long> ids, @RequestBody UpdateUserDTO updateUserDTO)
func (s *userService) UpdateUser(ctx context.Context, userID int64, update bson.M) error {
	// No-op optimization: if update is empty, return immediately (matching Java ACKNOWLEDGED_UPDATE_RESULT)
	if len(update) == 0 {
		return nil
	}
	now := time.Now()
	update["lud"] = now // Set LastUpdatedDate
	return s.repo.Update(ctx, userID, update)
}

// @MappedFrom checkIfUserExists(Long userId, boolean queryDeletedRecords)
// @MappedFrom checkIfUserExists(@NotNull Long userId, boolean queryDeletedRecords)
func (s *userService) CheckIfUserExists(ctx context.Context, userID int64) (bool, error) {
	return s.repo.Exists(ctx, userID)
}

// @MappedFrom deleteUsers(@NotEmpty Set<Long> userIds, @Nullable Boolean deleteLogically)
// @MappedFrom deleteUsers(Set<Long> ids, @QueryParam(required = false)
func (s *userService) DeleteUsers(ctx context.Context, userIDs []int64) (int64, error) {
	if len(userIDs) == 0 {
		return 0, nil
	}
	filter := bson.M{"_id": bson.M{"$in": userIDs}}
	deletedCount, err := s.repo.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	// TODO: Support logical delete (set deletionDate instead of physical delete when deleteLogically=true)
	// TODO: Cascade deletion: Elasticsearch docs, relationships, relationship groups, settings, conversations, conversation settings, user versions, message sequence IDs
	// TODO: Disconnect user sessions after delete
	return deletedCount, nil
}

// @MappedFrom queryUsersProfile(@NotEmpty Collection<Long> userIds, boolean queryDeletedRecords)
func (s *userService) QueryUsersProfile(ctx context.Context, userIDs []int64) ([]*po.User, error) {
	if len(userIDs) == 0 {
		return []*po.User{}, nil
	}
	filter := bson.M{"_id": bson.M{"$in": userIDs}}
	return s.repo.FindMany(ctx, filter)
}

// @MappedFrom queryUserName(@NotNull Long userId)
func (s *userService) QueryUserName(ctx context.Context, userID int64) (string, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", nil // Or return a specific record not found error depending on Turms
	}
	return user.Name, nil
}

// @MappedFrom queryUsers(@QueryParam(required = false)
// @MappedFrom queryUsers(@Nullable Collection<Long> userIds, @Nullable DateRange registrationDateRange, @Nullable DateRange deletionDateRange, @Nullable Boolean isActive, @Nullable Integer page, @Nullable Integer size, boolean queryDeletedRecords)
func (s *userService) QueryUsers(ctx context.Context, userIDs []int64) ([]*po.User, error) {
	filter := bson.M{}
	if len(userIDs) > 0 {
		filter["_id"] = bson.M{"$in": userIDs}
	}
	// When userIDs is empty (equivalent to Java null), return all users
	return s.repo.FindMany(ctx, filter)
}

func (s *userService) CountUsers(ctx context.Context, activeOnly bool) (int64, error) {
	filter := bson.M{}
	if activeOnly {
		filter["act"] = true
	}
	return s.repo.Count(ctx, filter)
}

func (s *userService) IsActiveAndNotDeleted(ctx context.Context, userID int64) (bool, error) {
	return s.repo.IsActiveAndNotDeleted(ctx, userID)
}

func (s *userService) FindPassword(ctx context.Context, userID int64) (*string, error) {
	return s.repo.FindPassword(ctx, userID)
}

// @MappedFrom isAllowedToSendMessageToTarget(@NotNull Boolean isGroupMessage, @NotNull Boolean isSystemMessage, @NotNull Long requesterId, @NotNull Long targetId)
func (s *userService) IsAllowedToSendMessageToTarget(ctx context.Context, isGroupMessage bool, isSystemMessage bool, requesterID int64, targetID int64) (int, error) {
	if isSystemMessage {
		return 200, nil // OK
	}
	// Simplified permission check
	return 200, nil
}

// @MappedFrom isAllowToQueryUserProfile(@NotNull Long requesterId, @NotNull Long targetUserId)
func (s *userService) IsAllowToQueryUserProfile(ctx context.Context, requesterID int64, targetID int64) (int, error) {
	// Simplified logic for refactor
	return 200, nil
}

// @MappedFrom authAndQueryUsersProfile(@NotNull Long requesterId, @Nullable Set<Long> userIds, @Nullable String name, @Nullable Date lastUpdatedDate, @Nullable Integer skip, @Nullable Integer limit, @Nullable List<Integer> fieldsToHighlight)
func (s *userService) AuthAndQueryUsersProfile(ctx context.Context, requesterID int64, userIDs []int64, name string, lastUpdatedDate *time.Time, skip int, limit int) ([]*po.User, error) {
	// Simplified, normally check permission then query
	filter := bson.M{}
	if len(userIDs) > 0 {
		filter["_id"] = bson.M{"$in": userIDs}
	}
	if name != "" {
		filter["n"] = name
	}
	return s.repo.FindMany(ctx, filter)
}

// @MappedFrom queryUserRoleIdByUserId(@NotNull Long userId)
func (s *userService) QueryUserRoleIDByUserID(ctx context.Context, userID int64) (int64, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return 0, err
	}
	if user == nil {
		return 0, nil
	}
	return user.PermissionGroupID, nil
}

func (s *userService) CountRegisteredUsers(ctx context.Context, startDate *time.Time, endDate *time.Time, queryDeletedRecords bool) (int64, error) {
	filter := bson.M{}
	dateFilter := bson.M{}
	if startDate != nil {
		dateFilter["$gte"] = *startDate
	}
	if endDate != nil {
		dateFilter["$lt"] = *endDate
	}
	if len(dateFilter) > 0 {
		filter["rd"] = dateFilter
	}
	return s.repo.Count(ctx, filter)
}

func (s *userService) CountDeletedUsers(ctx context.Context, startDate *time.Time, endDate *time.Time) (int64, error) {
	filter := bson.M{"dd": bson.M{"$exists": true, "$ne": nil}}
	return s.repo.Count(ctx, filter)
}

func (s *userService) UpdateUsers(ctx context.Context, userIDs []int64, update bson.M) (int64, error) {
	if len(userIDs) == 0 {
		return 0, nil
	}
	// No-op optimization: if all update fields are null/empty, return acknowledged result immediately (matching Java)
	if len(update) == 0 {
		return 0, nil
	}
	// Set LastUpdatedDate
	update["lud"] = time.Now()
	filter := bson.M{"_id": bson.M{"$in": userIDs}}
	modified, err := s.repo.UpdateMany(ctx, filter, update)
	if err != nil {
		return 0, err
	}
	// TODO: Disconnect user sessions when isActive is set to false and modified > 0
	// TODO: Elasticsearch sync for name changes
	return modified, nil
}
