package repository

// ConversationSettingsRepository maps to ConversationSettingsRepository.java
// @MappedFrom ConversationSettingsRepository
type ConversationSettingsRepository struct {
}

// @MappedFrom findByIdAndSettingNames(Long ownerId, @Nullable Collection<String> settingNames, @Nullable Date lastUpdatedDateStart)
func (r *ConversationSettingsRepository) FindByIdAndSettingNames() {
}

// @MappedFrom findByIdAndSettingNames(Collection<ConversationSettings.Key> keys, @Nullable Collection<String> settingNames, @Nullable Date lastUpdatedDateStart)
func (r *ConversationSettingsRepository) FindByIdAndSettingNamesWithKeys() {
}

// @MappedFrom findSettingFields(Long ownerId, Long targetId, Collection<String> includedFields)
func (r *ConversationSettingsRepository) FindSettingFields() {
}

// @MappedFrom deleteByOwnerIds(Collection<Long> ownerIds, @Nullable ClientSession clientSession)
func (r *ConversationSettingsRepository) DeleteByOwnerIds() {
}
