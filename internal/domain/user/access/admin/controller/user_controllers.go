package controller

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
func (c *UserOnlineInfoController) QueryUserLocations() {
	// TODO: implement
}

// @MappedFrom updateUserOnlineStatus(Set<Long> ids, @QueryParam(required = false)
func (c *UserOnlineInfoController) UpdateUserOnlineStatus() {
	// TODO: implement
}

// UserRoleController maps to UserRoleController.java
// @MappedFrom UserRoleController
type UserRoleController struct {
}

// @MappedFrom addUserRole(@RequestBody AddUserRoleDTO addUserRoleDTO)
func (c *UserRoleController) AddUserRole() {
	// TODO: implement
}

// @MappedFrom queryUserRoles(@QueryParam(required = false)
func (c *UserRoleController) QueryUserRoles() {
	// TODO: implement
}

// @MappedFrom queryUserRoleGroups(int page, @QueryParam(required = false)
func (c *UserRoleController) QueryUserRoleGroups() {
	// TODO: implement
}

// @MappedFrom updateUserRole(Set<Long> ids, @RequestBody UpdateUserRoleDTO updateUserRoleDTO)
func (c *UserRoleController) UpdateUserRole() {
	// TODO: implement
}

// @MappedFrom deleteUserRole(Set<Long> ids)
func (c *UserRoleController) DeleteUserRole() {
	// TODO: implement
}

// UserFriendRequestController maps to UserFriendRequestController.java
// @MappedFrom UserFriendRequestController
type UserFriendRequestController struct {
}

// @MappedFrom createFriendRequest(@RequestBody AddFriendRequestDTO addFriendRequestDTO)
func (c *UserFriendRequestController) CreateFriendRequest() {
	// TODO: implement
}

// @MappedFrom queryFriendRequests(@QueryParam(required = false)
func (c *UserFriendRequestController) QueryFriendRequests() {
	// TODO: implement
}

// @MappedFrom updateFriendRequests(Set<Long> ids, @RequestBody UpdateFriendRequestDTO updateFriendRequestDTO)
func (c *UserFriendRequestController) UpdateFriendRequests() {
	// TODO: implement
}

// @MappedFrom deleteFriendRequests(@QueryParam(required = false)
func (c *UserFriendRequestController) DeleteFriendRequests() {
	// TODO: implement
}

// UserRelationshipController maps to UserRelationshipController.java
// @MappedFrom UserRelationshipController
type UserRelationshipController struct {
}

// @MappedFrom addRelationship(@RequestBody AddRelationshipDTO addRelationshipDTO)
func (c *UserRelationshipController) AddRelationship() {
	// TODO: implement
}

// @MappedFrom queryRelationships(@QueryParam(required = false)
func (c *UserRelationshipController) QueryRelationships() {
	// TODO: implement
}

// @MappedFrom updateRelationships(List<UserRelationship.Key> keys, @RequestBody UpdateRelationshipDTO updateRelationshipDTO)
func (c *UserRelationshipController) UpdateRelationships() {
	// TODO: implement
}

// @MappedFrom deleteRelationships(List<UserRelationship.Key> keys)
func (c *UserRelationshipController) DeleteRelationships() {
	// TODO: implement
}

// UserRelationshipGroupController maps to UserRelationshipGroupController.java
// @MappedFrom UserRelationshipGroupController
type UserRelationshipGroupController struct {
}

// @MappedFrom addRelationshipGroup(@RequestBody AddRelationshipGroupDTO addRelationshipGroupDTO)
func (c *UserRelationshipGroupController) AddRelationshipGroup() {
	// TODO: implement
}

// @MappedFrom deleteRelationshipGroups(@QueryParam(required = false)
func (c *UserRelationshipGroupController) DeleteRelationshipGroups() {
	// TODO: implement
}

// @MappedFrom updateRelationshipGroups(List<UserRelationshipGroup.Key> keys, @RequestBody UpdateRelationshipGroupDTO updateRelationshipGroupDTO)
func (c *UserRelationshipGroupController) UpdateRelationshipGroups() {
	// TODO: implement
}

// @MappedFrom queryRelationshipGroups(@QueryParam(required = false)
func (c *UserRelationshipGroupController) QueryRelationshipGroups() {
	// TODO: implement
}
