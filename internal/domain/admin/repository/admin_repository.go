package repository

// AdminRepository maps to AdminRepository in Java.
// @MappedFrom AdminRepository
type AdminRepository struct {
}

// @MappedFrom updateAdmins(Set<Long> ids, @Nullable byte[] password, @Nullable String displayName, @Nullable Set<Long> roleIds)
func (r *AdminRepository) UpdateAdmins() {
}

// @MappedFrom countAdmins(@Nullable Set<Long> ids, @Nullable Set<Long> roleIds)
func (r *AdminRepository) CountAdmins() {
}

// @MappedFrom findAdmins(@Nullable Collection<Long> ids, @Nullable Collection<String> loginNames, @Nullable Collection<Long> roleIds, @Nullable Integer page, @Nullable Integer size)
func (r *AdminRepository) FindAdmins() {
}
