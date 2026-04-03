package repository

// AdminRoleRepository maps to AdminRoleRepository in Java.
// @MappedFrom AdminRoleRepository
type AdminRoleRepository struct {
}

// @MappedFrom updateAdminRoles(Set<Long> roleIds, String newName, @Nullable Set<AdminPermission> permissions, @Nullable Integer rank)
func (r *AdminRoleRepository) UpdateAdminRoles() {
}

// @MappedFrom countAdminRoles(@Nullable Set<Long> ids, @Nullable Set<String> names, @Nullable Set<AdminPermission> includedPermissions, @Nullable Set<Integer> ranks)
func (r *AdminRoleRepository) CountAdminRoles() {
}

// @MappedFrom findAdminRoles(@Nullable Set<Long> roleIds, @Nullable Set<String> names, @Nullable Set<AdminPermission> includedPermissions, @Nullable Set<Integer> ranks, @Nullable Integer page, @Nullable Integer size)
func (r *AdminRoleRepository) FindAdminRoles() {
}

// @MappedFrom findAdminRolesByIdsAndRankGreaterThan(@NotNull Collection<Long> roleIds, @Nullable Integer rankGreaterThan)
func (r *AdminRoleRepository) FindAdminRolesByIdsAndRankGreaterThan() {
}

// @MappedFrom findHighestRankByRoleIds(Set<Long> roleIds)
func (r *AdminRoleRepository) FindHighestRankByRoleIds() {
}
