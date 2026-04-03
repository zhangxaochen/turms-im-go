package service

// AdminRoleService maps to AdminRoleService in Java.
// @MappedFrom AdminRoleService
type AdminRoleService struct {
}

// @MappedFrom authAndAddAdminRole(@NotNull Long requesterId, @NotNull Long roleId, @NotNull @NoWhitespace @Size( min = MIN_ROLE_NAME_LIMIT, max = MAX_ROLE_NAME_LIMIT)
func (s *AdminRoleService) AuthAndAddAdminRole() {
}

// @MappedFrom addAdminRole(@NotNull Long roleId, @NotNull @NoWhitespace @Size( min = MIN_ROLE_NAME_LIMIT, max = MAX_ROLE_NAME_LIMIT)
func (s *AdminRoleService) AddAdminRole() {
}

// @MappedFrom authAndDeleteAdminRoles(@NotNull Long requesterId, @NotEmpty Set<Long> roleIds)
func (s *AdminRoleService) AuthAndDeleteAdminRoles() {
}

// @MappedFrom deleteAdminRoles(@NotEmpty Set<Long> roleIds)
func (s *AdminRoleService) DeleteAdminRoles() {
}

// @MappedFrom authAndUpdateAdminRoles(@NotNull Long requesterId, @NotEmpty Set<Long> roleIds, @Nullable @NoWhitespace @Size( min = MIN_ROLE_NAME_LIMIT, max = MAX_ROLE_NAME_LIMIT)
func (s *AdminRoleService) AuthAndUpdateAdminRoles() {
}

// @MappedFrom updateAdminRole(@NotEmpty Set<Long> roleIds, @Nullable @NoWhitespace @Size( min = MIN_ROLE_NAME_LIMIT, max = MAX_ROLE_NAME_LIMIT)
func (s *AdminRoleService) UpdateAdminRole() {
}

// @MappedFrom queryAdminRoles(@Nullable Set<Long> ids, @Nullable Set<String> names, @Nullable Set<AdminPermission> includedPermissions, @Nullable Set<Integer> ranks, @Nullable Integer page, @Nullable Integer size)
func (s *AdminRoleService) QueryAdminRoles() {
}

// @MappedFrom queryAndCacheRolesByRoleIdsAndRankGreaterThan(@NotNull Collection<Long> roleIds, @NotNull Integer rankGreaterThan)
func (s *AdminRoleService) QueryAndCacheRolesByRoleIdsAndRankGreaterThan() {
}

// @MappedFrom countAdminRoles(@Nullable Set<Long> ids, @Nullable Set<String> names, @Nullable Set<AdminPermission> includedPermissions, @Nullable Set<Integer> ranks)
func (s *AdminRoleService) CountAdminRoles() {
}

// @MappedFrom queryHighestRankByAdminId(@NotNull Long adminId)
func (s *AdminRoleService) QueryHighestRankByAdminId() {
}

// @MappedFrom queryHighestRankByRoleIds(@NotNull Set<Long> roleIds)
func (s *AdminRoleService) QueryHighestRankByRoleIds() {
}

// @MappedFrom isAdminRankHigherThanRank(@NotNull Long adminId, @NotNull Integer rank)
func (s *AdminRoleService) IsAdminRankHigherThanRank() {
}

// @MappedFrom queryPermissions(@NotNull Long adminId)
func (s *AdminRoleService) QueryPermissions() {
}

// AdminService maps to AdminService in Java.
// @MappedFrom AdminService
type AdminService struct {
}

// @MappedFrom queryRoleIdsByAdminIds(@NotEmpty Set<Long> adminIds)
func (s *AdminService) QueryRoleIdsByAdminIds() {
}

// @MappedFrom authAndAddAdmin(@NotNull Long requesterId, @Nullable @NoWhitespace @Size( min = MIN_LOGIN_NAME_LIMIT, max = MAX_LOGIN_NAME_LIMIT)
func (s *AdminService) AuthAndAddAdmin() {
}

// @MappedFrom addAdmin(@Nullable Long id, @Nullable @NoWhitespace @Size( min = MIN_LOGIN_NAME_LIMIT, max = MAX_LOGIN_NAME_LIMIT)
func (s *AdminService) AddAdmin() {
}

// @MappedFrom queryAdmins(@Nullable Collection<Long> ids, @Nullable Collection<String> loginNames, @Nullable Collection<Long> roleIds, @Nullable Integer page, @Nullable Integer size)
func (s *AdminService) QueryAdmins() {
}

// @MappedFrom authAndDeleteAdmins(@NotNull Long requesterId, @NotEmpty Set<Long> adminIds)
func (s *AdminService) AuthAndDeleteAdmins() {
}

// @MappedFrom authAndUpdateAdmins(@NotNull Long requesterId, @NotEmpty Set<Long> targetAdminIds, @Nullable @NoWhitespace @Size( min = MIN_PASSWORD_LIMIT, max = MAX_PASSWORD_LIMIT)
func (s *AdminService) AuthAndUpdateAdmins() {
}

// @MappedFrom updateAdmins(@NotEmpty Set<Long> targetAdminIds, @Nullable @NoWhitespace @Size( min = MIN_PASSWORD_LIMIT, max = MAX_PASSWORD_LIMIT)
func (s *AdminService) UpdateAdmins() {
}

// @MappedFrom countAdmins(@Nullable Set<Long> ids, @Nullable Set<Long> roleIds)
func (s *AdminService) CountAdmins() {
}

// @MappedFrom errorRequesterNotExist()
func (s *AdminService) ErrorRequesterNotExist() {
}
