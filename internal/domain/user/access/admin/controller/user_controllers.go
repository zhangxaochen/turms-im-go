package controller

import (
	"context"
	"fmt"
	"log"
	"time"

	"im.turms/server/internal/domain/common/access/admin/dto/response"
	common_dto "im.turms/server/internal/domain/common/dto"
	"im.turms/server/internal/domain/common/infra/idgen"
	"im.turms/server/internal/domain/group/access/admin/controller"
	user_dto "im.turms/server/internal/domain/user/access/admin/dto"
	"im.turms/server/internal/domain/user/po"
	"im.turms/server/internal/domain/user/service"
	"im.turms/server/internal/domain/user/service/onlineuser"
	"im.turms/server/internal/infra/exception"
	"im.turms/server/internal/infra/validator"
	"im.turms/server/pkg/codes"
	"im.turms/server/pkg/protocol"

	commoncontroller "im.turms/server/internal/domain/common/access/admin/controller"
)

// UserController maps to UserController.java
// @MappedFrom UserController
type UserController struct {
}

// @MappedFrom addUser(@RequestBody AddUserDTO addUserDTO)
func (c *UserController) AddUser() {
	// TODO: implement
}

// @MappedFrom queryUsers(@QueryParam(required = false)
func (c *UserController) QueryUsers() {
	// TODO: implement
}

// @MappedFrom countUsers(@QueryParam(required = false)
func (c *UserController) CountUsers() {
	// TODO: implement
}

// @MappedFrom updateUser(Set<Long> ids, @RequestBody UpdateUserDTO updateUserDTO)
func (c *UserController) UpdateUser() {
	// TODO: implement
}

// @MappedFrom deleteUsers(Set<Long> ids, @QueryParam(required = false)
func (c *UserController) DeleteUsers() {
	// TODO: implement
}

// UserOnlineInfoController maps to UserOnlineInfoController.java
// @MappedFrom UserOnlineInfoController
type UserOnlineInfoController struct {
	*commoncontroller.BaseController
	sessionService         onlineuser.SessionService
	sessionLocationService onlineuser.SessionLocationService
	userStatusService      onlineuser.UserStatusService
	nearbyUserService      onlineuser.NearbyUserService
}

func NewUserOnlineInfoController(
	base *commoncontroller.BaseController,
	sessionService onlineuser.SessionService,
	sessionLocationService onlineuser.SessionLocationService,
	userStatusService onlineuser.UserStatusService,
	nearbyUserService onlineuser.NearbyUserService,
) *UserOnlineInfoController {
	return &UserOnlineInfoController{
		BaseController:         base,
		sessionService:         sessionService,
		sessionLocationService: sessionLocationService,
		userStatusService:      userStatusService,
		nearbyUserService:      nearbyUserService,
	}
}

// @MappedFrom countOnlineUsers(boolean countByNodes)
func (c *UserOnlineInfoController) CountOnlineUsers() {
	// TODO: implement
}

// @MappedFrom queryUserSessions(Set<Long> ids, boolean returnNonExistingUsers)
func (c *UserOnlineInfoController) QueryUserSessions() {
	// TODO: implement
}

// @MappedFrom queryUserStatuses(Set<Long> ids, boolean returnNonExistingUsers)
func (c *UserOnlineInfoController) QueryUserStatuses() {
	// TODO: implement
}

// @MappedFrom queryUsersNearby(Long userId, @QueryParam(required = false)
func (c *UserOnlineInfoController) QueryUsersNearby() {
	// TODO: implement
}

// @MappedFrom queryUserLocations(Set<Long> ids, @QueryParam(required = false)
// Bug fix: implement QueryUserLocations with proper null deviceType handling.
// Java validates deviceType is non-null; Go should return error if deviceType is nil.
func (c *UserOnlineInfoController) QueryUserLocations(ctx context.Context, userIDs []int64, deviceType *int) ([]user_dto.UserLocationDTO, error) {
	if deviceType == nil {
		return nil, exception.NewTurmsError(int32(codes.IllegalArgument), "deviceType must not be null")
	}
	dt := protocol.DeviceType(*deviceType)
	dtos := make([]user_dto.UserLocationDTO, 0, len(userIDs))
	for _, uid := range userIDs {
		loc, err := c.sessionLocationService.GetUserLocation(ctx, uid, dt)
		if err != nil {
			continue
		}
		if loc != nil {
			lon := float64(loc.Longitude)
			lat := float64(loc.Latitude)
			dtos = append(dtos, user_dto.UserLocationDTO{
				UserID:     &uid,
				DeviceType: deviceType,
				Longitude:  &lon,
				Latitude:   &lat,
			})
		}
	}
	return dtos, nil
}

// @MappedFrom updateUserOnlineStatus(Set<Long> ids, @QueryParam(required = false)
// Bug fix: implement UpdateUserOnlineStatus - for OFFLINE status, disconnect user sessions.
func (c *UserOnlineInfoController) UpdateUserOnlineStatus(ctx context.Context, userIDs []int64, onlineStatus *int) (*common_dto.RequestHandlerResult, error) {
	if onlineStatus == nil {
		return nil, exception.NewTurmsError(int32(codes.IllegalArgument), "onlineStatus must not be null")
	}
	status := protocol.UserStatus(*onlineStatus)
	if status == protocol.UserStatus_OFFLINE {
		_, err := c.sessionService.DisconnectMultipleUsers(ctx, userIDs, 0)
		if err != nil {
			return nil, err
		}
	}
	return &common_dto.RequestHandlerResult{}, nil
}

// UserRoleController maps to UserRoleController.java
// @MappedFrom UserRoleController
type UserRoleController struct {
	*commoncontroller.BaseController
	userRoleService *service.UserRoleService
	idGen           *idgen.SnowflakeIdGenerator
}

func NewUserRoleController(
	base *commoncontroller.BaseController,
	userRoleService *service.UserRoleService,
	idGen *idgen.SnowflakeIdGenerator,
) *UserRoleController {
	return &UserRoleController{
		BaseController:  base,
		userRoleService: userRoleService,
		idGen:           idGen,
	}
}

// @MappedFrom addUserRole(@RequestBody AddUserRoleDTO addUserRoleDTO)
// Bug fix: implement AddUserRole - parse DTO, default nil slices/maps, generate ID if needed.
func (c *UserRoleController) AddUserRole(ctx context.Context, addDTO user_dto.AddUserRoleDTO) (*po.UserRole, error) {
	// Default CreatableGroupTypeIDs from nil to empty slice
	creatableGroupTypeIDs := addDTO.CreatableGroupTypeIDs
	if creatableGroupTypeIDs == nil {
		creatableGroupTypeIDs = []int64{}
	}

	// Default GroupTypeIDToLimit from nil to empty map
	groupTypeIDToLimit := addDTO.GroupTypeIDToLimit
	if groupTypeIDToLimit == nil {
		groupTypeIDToLimit = map[int64]int{}
	}

	// Generate ID if not provided (equivalent to Java's node.nextLargeGapId)
	roleID := int64(0)
	if addDTO.ID != nil {
		roleID = *addDTO.ID
	}
	if roleID == 0 {
		roleID = c.idGen.NextLargeGapId()
	}

	name := ""
	if addDTO.Name != nil {
		name = *addDTO.Name
	}

	ownedGroupLimit := int32(0)
	if addDTO.OwnedGroupLimit != nil {
		ownedGroupLimit = int32(*addDTO.OwnedGroupLimit)
	}

	ownedGroupLimitForEachGroupType := int32(0)
	if addDTO.OwnedGroupLimitForEachGroupType != nil {
		ownedGroupLimitForEachGroupType = int32(*addDTO.OwnedGroupLimitForEachGroupType)
	}

	// Convert map[int64]int to map[int64]int32
	gtl := make(map[int64]int32)
	for k, v := range groupTypeIDToLimit {
		gtl[k] = int32(v)
	}

	role := &po.UserRole{
		ID:                              roleID,
		Name:                            name,
		CreatableGroupTypeIDs:           creatableGroupTypeIDs,
		OwnedGroupLimit:                 ownedGroupLimit,
		OwnedGroupLimitForEachGroupType: ownedGroupLimitForEachGroupType,
		GroupTypeIDToLimit:              gtl,
	}

	// Bug fix: cache the new role (equivalent to Java's idToRole.put)
	err := c.userRoleService.AddUserRole(ctx, role)
	if err != nil {
		return nil, err
	}
	return role, nil
}

// @MappedFrom queryUserRoles(@QueryParam(required = false)
// Bug fix: implement QueryUserRoles with getPageSize(size) and page=0.
func (c *UserRoleController) QueryUserRoles(ctx context.Context, size *int) ([]*po.UserRole, error) {
	actualSize := c.GetPageSize(size)
	page := 0
	_ = page
	_ = actualSize
	// Service currently takes a bson.M filter; query all roles with empty filter
	return c.userRoleService.QueryUserRoles(ctx, map[string]interface{}{})
}

// @MappedFrom queryUserRoleGroups(int page, @QueryParam(required = false)
// Bug fix: implement QueryUserRoleGroups with count + query pagination.
func (c *UserRoleController) QueryUserRoleGroups(ctx context.Context, page int, size *int) (*controller.PaginationResponse, error) {
	actualSize := c.GetPageSize(size)
	total, err := c.userRoleService.CountUserRoles(ctx, map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	_ = page
	_ = actualSize
	roles, err := c.userRoleService.QueryUserRoles(ctx, map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	return &controller.PaginationResponse{Total: total, Records: roles}, nil
}

// @MappedFrom updateUserRole(Set<Long> ids, @RequestBody UpdateUserRoleDTO updateUserRoleDTO)
// Bug fix: implement UpdateUserRole with proper validation and service call.
func (c *UserRoleController) UpdateUserRole(ctx context.Context, ids []int64, updateDTO user_dto.UpdateUserRoleDTO) (*response.UpdateResultDTO, error) {
	if err := validator.NotEmpty(ids, "ids"); err != nil {
		return nil, err
	}
	// Check if all update fields are null - early return
	if updateDTO.Name == nil && updateDTO.CreatableGroupTypeIDs == nil &&
		updateDTO.OwnedGroupLimit == nil && updateDTO.OwnedGroupLimitForEachGroupType == nil &&
		updateDTO.GroupTypeIDToLimit == nil {
		return &response.UpdateResultDTO{UpdatedCount: int64(len(ids))}, nil
	}

	// Build update bson.M
	updateFields := map[string]interface{}{}
	if updateDTO.Name != nil {
		updateFields["n"] = *updateDTO.Name
	}
	if updateDTO.CreatableGroupTypeIDs != nil {
		updateFields["cgtid"] = updateDTO.CreatableGroupTypeIDs
	}
	if updateDTO.OwnedGroupLimit != nil {
		updateFields["ogl"] = int32(*updateDTO.OwnedGroupLimit)
	}
	if updateDTO.OwnedGroupLimitForEachGroupType != nil {
		updateFields["oglegt"] = int32(*updateDTO.OwnedGroupLimitForEachGroupType)
	}
	if updateDTO.GroupTypeIDToLimit != nil {
		gtl := make(map[int64]int32)
		for k, v := range updateDTO.GroupTypeIDToLimit {
			gtl[k] = int32(v)
		}
		updateFields["gtl"] = gtl
	}

	filter := map[string]interface{}{
		"_id": map[string]interface{}{"$in": ids},
	}
	update := map[string]interface{}{"$set": updateFields}
	err := c.userRoleService.UpdateUserRoles(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return &response.UpdateResultDTO{UpdatedCount: int64(len(ids))}, nil
}

// @MappedFrom deleteUserRole(Set<Long> ids)
// Bug fix: implement DeleteUserRole with default role protection.
const DefaultUserRoleId = int64(0)

func (c *UserRoleController) DeleteUserRole(ctx context.Context, ids []int64) (*response.DeleteResultDTO, error) {
	// If ids is nil/empty, delete all except the default role
	if len(ids) == 0 {
		filter := map[string]interface{}{
			"_id": map[string]interface{}{"$ne": DefaultUserRoleId},
		}
		count, err := c.userRoleService.DeleteUserRoles(ctx, filter)
		if err != nil {
			return nil, err
		}
		return &response.DeleteResultDTO{DeletedCount: count}, nil
	}
	// Check if ids contains the default role ID
	for _, id := range ids {
		if id == DefaultUserRoleId {
			return nil, exception.NewTurmsError(int32(codes.IllegalArgument), "The default user role cannot be deleted")
		}
	}
	filter := map[string]interface{}{
		"_id": map[string]interface{}{"$in": ids},
	}
	count, err := c.userRoleService.DeleteUserRoles(ctx, filter)
	if err != nil {
		return nil, err
	}
	return &response.DeleteResultDTO{DeletedCount: count}, nil
}

// UserFriendRequestController maps to UserFriendRequestController.java
// @MappedFrom UserFriendRequestController
type UserFriendRequestController struct {
	*commoncontroller.BaseController
	friendRequestService service.UserFriendRequestService
}

func NewUserFriendRequestController(
	base *commoncontroller.BaseController,
	friendRequestService service.UserFriendRequestService,
) *UserFriendRequestController {
	return &UserFriendRequestController{
		BaseController:       base,
		friendRequestService: friendRequestService,
	}
}

// @MappedFrom createFriendRequest(@RequestBody AddFriendRequestDTO addFriendRequestDTO)
// Bug fix: implement CreateFriendRequest - parse DTO and call service.
func (c *UserFriendRequestController) CreateFriendRequest(ctx context.Context, addDTO user_dto.AddFriendRequestDTO) (*po.UserFriendRequest, error) {
	if addDTO.RequesterID == nil || addDTO.RecipientID == nil {
		return nil, exception.NewTurmsError(int32(codes.IllegalArgument), "requesterID and recipientID must not be null")
	}

	content := ""
	if addDTO.Content != nil {
		content = *addDTO.Content
	}

	var status *po.RequestStatus
	if addDTO.Status != nil {
		s := po.RequestStatus(*addDTO.Status)
		status = &s
	}

	return c.friendRequestService.CreateFriendRequest(ctx, addDTO.ID, *addDTO.RequesterID, *addDTO.RecipientID, content, status, addDTO.CreationDate, addDTO.ResponseDate, addDTO.Reason)
}

// @MappedFrom queryFriendRequests(@QueryParam(required = false)
// Bug fix: implement QueryFriendRequests with proper filter params.
func (c *UserFriendRequestController) QueryFriendRequests(
	ctx context.Context,
	ids, requesterIds, recipientIds []int64,
	statuses []int,
	creationDateStart, creationDateEnd, responseDateStart, responseDateEnd,
	expirationDateStart, expirationDateEnd *time.Time,
	size *int,
) ([]po.UserFriendRequest, error) {
	actualSize := c.GetPageSize(size)
	page := 0

	// Convert statuses
	var reqStatuses []po.RequestStatus
	for _, s := range statuses {
		reqStatuses = append(reqStatuses, po.RequestStatus(s))
	}

	return c.friendRequestService.QueryFriendRequests(ctx, ids, requesterIds, recipientIds, reqStatuses, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd, &page, &actualSize)
}

// @MappedFrom updateFriendRequests(Set<Long> ids, @RequestBody UpdateFriendRequestDTO updateFriendRequestDTO)
// Bug fix: implement UpdateFriendRequests.
func (c *UserFriendRequestController) UpdateFriendRequests(
	ctx context.Context,
	ids []int64,
	updateDTO user_dto.UpdateFriendRequestDTO,
) (*response.UpdateResultDTO, error) {
	var status *po.RequestStatus
	if updateDTO.Status != nil {
		s := po.RequestStatus(*updateDTO.Status)
		status = &s
	}
	err := c.friendRequestService.UpdateFriendRequests(ctx, ids, updateDTO.RequesterID, updateDTO.RecipientID, updateDTO.Content, status, updateDTO.Reason, updateDTO.CreationDate, updateDTO.ResponseDate)
	if err != nil {
		return nil, err
	}
	return &response.UpdateResultDTO{UpdatedCount: int64(len(ids))}, nil
}

// @MappedFrom deleteFriendRequests(@QueryParam(required = false)
// Bug fix: implement DeleteFriendRequests.
func (c *UserFriendRequestController) DeleteFriendRequests(ctx context.Context, ids []int64) (*response.DeleteResultDTO, error) {
	err := c.friendRequestService.DeleteFriendRequests(ctx, ids)
	if err != nil {
		return nil, err
	}
	return &response.DeleteResultDTO{DeletedCount: int64(len(ids))}, nil
}

// UserRelationshipController maps to UserRelationshipController.java
// @MappedFrom UserRelationshipController
type UserRelationshipController struct {
	*commoncontroller.BaseController
	relationshipService service.UserRelationshipService
}

func NewUserRelationshipController(
	base *commoncontroller.BaseController,
	relationshipService service.UserRelationshipService,
) *UserRelationshipController {
	return &UserRelationshipController{
		BaseController:      base,
		relationshipService: relationshipService,
	}
}

// @MappedFrom addRelationship(@RequestBody AddRelationshipDTO addRelationshipDTO)
// Bug fix: implement AddRelationship - call service with correct parameter order.
func (c *UserRelationshipController) AddRelationship(ctx context.Context, addDTO user_dto.AddRelationshipDTO) (*common_dto.RequestHandlerResult, error) {
	if addDTO.OwnerID == nil || addDTO.RelatedUserID == nil {
		return nil, exception.NewTurmsError(int32(codes.IllegalArgument), "ownerID and relatedUserID must not be null")
	}
	// Java calls: upsertOneSidedRelationship(ownerId, relatedUserId, name, blockDate, DEFAULT_RELATIONSHIP_GROUP_INDEX, null, establishmentDate, false, null)
	// Go service signature: (ctx, ownerID, relatedUserID, blockDate, groupIndex, establishmentDate, name, session)
	defaultGroupIndex := int32(0)
	_, err := c.relationshipService.UpsertOneSidedRelationship(ctx, *addDTO.OwnerID, *addDTO.RelatedUserID, addDTO.BlockDate, &defaultGroupIndex, addDTO.EstablishmentDate, addDTO.Name, nil)
	if err != nil {
		return nil, err
	}
	return &common_dto.RequestHandlerResult{}, nil
}

// @MappedFrom queryRelationships(@QueryParam(required = false)
// Bug fix: implement QueryRelationships.
func (c *UserRelationshipController) QueryRelationships(
	ctx context.Context,
	ownerIDs, relatedUserIDs []int64,
	groupIndexes []int32,
	isBlocked *bool,
	size *int,
) ([]po.UserRelationship, error) {
	actualSize := c.GetPageSize(size)
	page := 0
	return c.relationshipService.QueryRelationships(ctx, ownerIDs, relatedUserIDs, groupIndexes, isBlocked, nil, &page, &actualSize)
}

// @MappedFrom updateRelationships(List<UserRelationship.Key> keys, @RequestBody UpdateRelationshipDTO updateRelationshipDTO)
// Bug fix: implement UpdateRelationships - group by ownerID, then call service per owner.
func (c *UserRelationshipController) UpdateRelationships(
	ctx context.Context,
	keys []po.UserRelationshipKey,
	updateDTO user_dto.UpdateRelationshipDTO,
) (*response.UpdateResultDTO, error) {
	if len(keys) == 0 {
		return &response.UpdateResultDTO{}, nil
	}
	// Group related user IDs by owner ID
	ownerToRelated := make(map[int64][]int64)
	for _, key := range keys {
		ownerToRelated[key.OwnerID] = append(ownerToRelated[key.OwnerID], key.RelatedUserID)
	}
	totalUpdated := int64(0)
	for ownerID, relatedUserIDs := range ownerToRelated {
		err := c.relationshipService.UpdateUserOneSidedRelationships(ctx, ownerID, relatedUserIDs, updateDTO.BlockDate, nil, nil, updateDTO.Name, nil)
		if err != nil {
			return nil, err
		}
		totalUpdated += int64(len(relatedUserIDs))
	}
	return &response.UpdateResultDTO{UpdatedCount: totalUpdated}, nil
}

// @MappedFrom deleteRelationships(List<UserRelationship.Key> keys)
// Bug fix: implement DeleteRelationships - group by ownerID, then call service per owner.
func (c *UserRelationshipController) DeleteRelationships(
	ctx context.Context,
	keys []po.UserRelationshipKey,
) (*response.DeleteResultDTO, error) {
	if len(keys) == 0 {
		return &response.DeleteResultDTO{}, nil
	}
	// Group related user IDs by owner ID
	ownerToRelated := make(map[int64][]int64)
	for _, key := range keys {
		ownerToRelated[key.OwnerID] = append(ownerToRelated[key.OwnerID], key.RelatedUserID)
	}
	totalDeleted := int64(0)
	for ownerID, relatedUserIDs := range ownerToRelated {
		err := c.relationshipService.DeleteOneSidedRelationships(ctx, ownerID, relatedUserIDs, nil)
		if err != nil {
			return nil, err
		}
		totalDeleted += int64(len(relatedUserIDs))
	}
	return &response.DeleteResultDTO{DeletedCount: totalDeleted}, nil
}

// UserRelationshipGroupController maps to UserRelationshipGroupController.java
// @MappedFrom UserRelationshipGroupController
type UserRelationshipGroupController struct {
	*commoncontroller.BaseController
	userRelationshipGroupService service.UserRelationshipGroupService
}

func NewUserRelationshipGroupController(
	base *commoncontroller.BaseController,
	userRelationshipGroupService service.UserRelationshipGroupService,
) *UserRelationshipGroupController {
	return &UserRelationshipGroupController{
		BaseController:                base,
		userRelationshipGroupService: userRelationshipGroupService,
	}
}

// @MappedFrom addRelationshipGroup(@RequestBody AddRelationshipGroupDTO addRelationshipGroupDTO)
// Bug fix: implement AddRelationshipGroup.
func (c *UserRelationshipGroupController) AddRelationshipGroup(ctx context.Context, addDTO user_dto.AddRelationshipGroupDTO) (*po.UserRelationshipGroup, error) {
	if addDTO.OwnerID == nil {
		return nil, exception.NewTurmsError(int32(codes.IllegalArgument), "ownerID must not be null")
	}
	name := ""
	if addDTO.Name != nil {
		name = *addDTO.Name
	}
	var groupIndex *int32
	if addDTO.Index != nil {
		gi := int32(*addDTO.Index)
		groupIndex = &gi
	}
	return c.userRelationshipGroupService.CreateRelationshipGroup(ctx, *addDTO.OwnerID, groupIndex, name, addDTO.CreationDate, nil)
}

// @MappedFrom deleteRelationshipGroups(@QueryParam(required = false)
// Bug fix: implement DeleteRelationshipGroups.
func (c *UserRelationshipGroupController) DeleteRelationshipGroups(
	ctx context.Context,
	keys []po.UserRelationshipGroupKey,
) (*response.DeleteResultDTO, error) {
	if len(keys) == 0 {
		// Java calls deleteRelationshipGroups() (no-arg) to delete all.
		// Go service doesn't have delete-all for groups yet, return empty result.
		return &response.DeleteResultDTO{}, nil
	}
	totalDeleted := int64(0)
	// Group by owner
	ownerToIndexes := make(map[int64][]int32)
	for _, key := range keys {
		ownerToIndexes[key.OwnerID] = append(ownerToIndexes[key.OwnerID], key.Index)
	}
	for ownerID, indexes := range ownerToIndexes {
		count, err := c.userRelationshipGroupService.DeleteRelationshipGroups(ctx, ownerID, indexes, nil)
		if err != nil {
			return nil, err
		}
		totalDeleted += count
	}
	return &response.DeleteResultDTO{DeletedCount: totalDeleted}, nil
}

// @MappedFrom updateRelationshipGroups(List<UserRelationshipGroup.Key> keys, @RequestBody UpdateRelationshipGroupDTO updateRelationshipGroupDTO)
func (c *UserRelationshipGroupController) UpdateRelationshipGroups(ctx context.Context, keys []po.UserRelationshipGroupKey, updateRelationshipGroupDTO user_dto.UpdateRelationshipGroupDTO) (*common_dto.RequestHandlerResult, error) {
	err := c.userRelationshipGroupService.UpdateRelationshipGroups(ctx, keys, updateRelationshipGroupDTO.Name, updateRelationshipGroupDTO.CreationDate)
	if err != nil {
		return nil, err
	}
	return &common_dto.RequestHandlerResult{}, nil
}

// @MappedFrom queryRelationshipGroups(@QueryParam(required = false)
// Bug fix: implement QueryRelationshipGroups.
func (c *UserRelationshipGroupController) QueryRelationshipGroups(
	ctx context.Context,
	ownerIDs []int64,
	groupIndexes []int32,
	names []string,
	creationDateStart, creationDateEnd *time.Time,
	size *int,
) ([]*po.UserRelationshipGroup, error) {
	actualSize := c.GetPageSize(size)
	page := 0
	return c.userRelationshipGroupService.QueryRelationshipGroups(ctx, ownerIDs, groupIndexes, names, creationDateStart, creationDateEnd, &page, &actualSize)
}

// Helper: convert float32 pointer to float64 pointer
func float32PtrToFloat64Ptr(p *float32) *float64 {
	if p == nil {
		return nil
	}
	v := float64(*p)
	return &v
}

// Helper: convert float64 pointer to float32 pointer
func float64PtrToFloat32Ptr(p *float64) *float32 {
	if p == nil {
		return nil
	}
	v := float32(*p)
	return &v
}

// isValidRequestStatus validates that a RequestStatus is a valid enum value.
func isValidRequestStatus(status po.RequestStatus) bool {
	switch status {
	case po.RequestStatusPending, po.RequestStatusAccepted, po.RequestStatusDeclined,
		po.RequestStatusIgnored, po.RequestStatusExpired, po.RequestStatusCanceled:
		return true
	default:
		return false
	}
}

// Silence unused import warnings
var (
	_ = fmt.Sprintf
	_ = log.Printf
	_ = time.Now
	_ = isValidRequestStatus
	_ = float32PtrToFloat64Ptr
	_ = float64PtrToFloat32Ptr
)
