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
