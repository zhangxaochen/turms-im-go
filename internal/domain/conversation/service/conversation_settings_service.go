package service

// ConversationSettingsService maps to ConversationSettingsService.java
// @MappedFrom ConversationSettingsService
type ConversationSettingsService struct {
}

// @MappedFrom upsertPrivateConversationSettings(Long ownerId, Long userId, Map<String, Value> settings)
func (s *ConversationSettingsService) UpsertPrivateConversationSettings() {
}

// @MappedFrom upsertGroupConversationSettings(Long ownerId, Long groupId, Map<String, Value> settings)
func (s *ConversationSettingsService) UpsertGroupConversationSettings() {
}
