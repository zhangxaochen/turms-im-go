const fs = require('fs');

const path = 'docs/refactor_progress_report.md';
let content = fs.readFileSync(path, 'utf8');

const mappings = [
    // AdminService
    ['authAndAddAdmin(@NotNull Long requesterId, @Nullable @NoWhitespace @Size( min = MIN_LOGIN_NAME_LIMIT, max = MAX_LOGIN_NAME_LIMIT)', 'internal/domain/admin/service/admin_services.go:AuthAndAddAdmin()'],
    ['addAdmin(@Nullable Long id, @Nullable @NoWhitespace @Size( min = MIN_LOGIN_NAME_LIMIT, max = MAX_LOGIN_NAME_LIMIT)', 'internal/domain/admin/service/admin_services.go:AddAdmin()'],
    ['queryAdmins(@Nullable Collection<Long> ids, @Nullable Collection<String> loginNames, @Nullable Collection<Long> roleIds, @Nullable Integer page, @Nullable Integer size)', 'internal/domain/admin/service/admin_services.go:QueryAdmins()'],
    ['authAndDeleteAdmins(@NotNull Long requesterId, @NotEmpty Set<Long> adminIds)', 'internal/domain/admin/service/admin_services.go:AuthAndDeleteAdmins()'],
    ['authAndUpdateAdmins(@NotNull Long requesterId, @NotEmpty Set<Long> targetAdminIds, @Nullable @NoWhitespace @Size( min = MIN_PASSWORD_LIMIT, max = MAX_PASSWORD_LIMIT)', 'internal/domain/admin/service/admin_services.go:AuthAndUpdateAdmins()'],
    ['updateAdmins(@NotEmpty Set<Long> targetAdminIds, @Nullable @NoWhitespace @Size( min = MIN_PASSWORD_LIMIT, max = MAX_PASSWORD_LIMIT)', 'internal/domain/admin/service/admin_services.go:UpdateAdmins()'],
    ['countAdmins(@Nullable Set<Long> ids, @Nullable Set<Long> roleIds)', 'internal/domain/admin/service/admin_services.go:CountAdmins()'],
    ['errorRequesterNotExist()', 'internal/domain/admin/service/admin_services.go:ErrorRequesterNotExist()'],

    // Blocklist
    ['addBlockedIps(@RequestBody AddBlockedIpsDTO addBlockedIpsDTO)', 'internal/domain/blocklist/access/admin/controller/blocklist_controllers.go:AddBlockedIps()'],
    ['queryBlockedIps(Set<String> ids)', 'internal/domain/blocklist/access/admin/controller/blocklist_controllers.go:QueryBlockedIpsByIds()'],
    ['queryBlockedIps(int page, @QueryParam(required = false)', 'internal/domain/blocklist/access/admin/controller/blocklist_controllers.go:QueryBlockedIpsByPage()'],
    ['deleteBlockedIps(@QueryParam(required = false)', 'internal/domain/blocklist/access/admin/controller/blocklist_controllers.go:DeleteBlockedIps()'],

    ['addBlockedUserIds(@RequestBody AddBlockedUserIdsDTO addBlockedUserIdsDTO)', 'internal/domain/blocklist/access/admin/controller/blocklist_controllers.go:AddBlockedUserIds()'],
    ['deleteBlockedUserIds(@QueryParam(required = false)', 'internal/domain/blocklist/access/admin/controller/blocklist_controllers.go:DeleteBlockedUserIds()'],

    ['AddBlockedIpsDTO(Set<String> ids, long blockDurationMillis)', 'internal/domain/blocklist/access/admin/dto/blocklist_dtos.go:AddBlockedIpsDTO'],
    ['AddBlockedUserIdsDTO(Set<Long> ids, long blockDurationMillis)', 'internal/domain/blocklist/access/admin/dto/blocklist_dtos.go:AddBlockedUserIdsDTO'],
    ['BlockedIpDTO(String id, Date blockEndTime)', 'internal/domain/blocklist/access/admin/dto/blocklist_dtos.go:BlockedIpDTO'],
    ['BlockedUserDTO(Long id, Date blockEndTime)', 'internal/domain/blocklist/access/admin/dto/blocklist_dtos.go:BlockedUserDTO'],

    // Cluster admin
    ['queryMembers()', 'internal/domain/cluster/access/admin/controller/cluster_controllers.go:QueryMembers()'],
    ['removeMembers(List<String> ids)', 'internal/domain/cluster/access/admin/controller/cluster_controllers.go:RemoveMembers()'],
    ['addMember(@RequestBody AddMemberDTO addMemberDTO)', 'internal/domain/cluster/access/admin/controller/cluster_controllers.go:AddMember()'],
    ['updateMember(String id, @RequestBody UpdateMemberDTO updateMemberDTO)', 'internal/domain/cluster/access/admin/controller/cluster_controllers.go:UpdateMember()'],
    ['queryLeader()', 'internal/domain/cluster/access/admin/controller/cluster_controllers.go:QueryLeader()'],
    ['electNewLeader(@QueryParam(required = false)', 'internal/domain/cluster/access/admin/controller/cluster_controllers.go:ElectNewLeader()'],

    ['queryClusterSettings(boolean queryLocalSettings, boolean onlyMutable)', 'internal/domain/cluster/access/admin/controller/cluster_controllers.go:QueryClusterSettings()'],
    ['updateClusterSettings(boolean reset, boolean updateLocalSettings, @RequestBody(required = false)', 'internal/domain/cluster/access/admin/controller/cluster_controllers.go:UpdateClusterSettings()'],
    ['queryClusterConfigMetadata(boolean queryLocalSettings, boolean onlyMutable, boolean withValue)', 'internal/domain/cluster/access/admin/controller/cluster_controllers.go:QueryClusterConfigMetadata()'],

    ['AddMemberDTO(String nodeId, String zone, String name, NodeType nodeType, String version, boolean isSeed, boolean isLeaderEligible, Date registrationDate, int priority, String memberHost, int memberPort, String adminApiAddress, String wsAddress, String tcpAddress, String udpAddress, boolean isActive, boolean isHealthy)', 'internal/domain/cluster/access/admin/dto/cluster_dtos.go:AddMemberDTO'],
    ['UpdateMemberDTO(String zone, String name, Boolean isSeed, Boolean isLeaderEligible, Boolean isActive, Integer priority)', 'internal/domain/cluster/access/admin/dto/cluster_dtos.go:UpdateMemberDTO'],
    ['SettingsDTO(int schemaVersion, Map<String, Object> settings)', 'internal/domain/cluster/access/admin/dto/cluster_dtos.go:SettingsDTO'],

    // Common Base controller & DTO
    ['getPageSize(@Nullable Integer size)', 'internal/domain/common/access/admin/controller/base_controller.go:GetPageSize()'],
    ['queryBetweenDate(DateRange dateRange, DivideBy divideBy, Function3<DateRange, Boolean, Boolean, Mono<Long>> function, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages)', 'internal/domain/common/access/admin/controller/base_controller.go:QueryBetweenDate()'],
    ['queryBetweenDate(DateRange dateRange, DivideBy divideBy, Function<DateRange, Mono<Long>> function)', 'internal/domain/common/access/admin/controller/base_controller.go:QueryBetweenDateFunc()'],
    ['checkAndQueryBetweenDate(DateRange dateRange, DivideBy divideBy, Function3<DateRange, Boolean, Boolean, Mono<Long>> function, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages)', 'internal/domain/common/access/admin/controller/base_controller.go:CheckAndQueryBetweenDate()'],
    ['checkAndQueryBetweenDate(DateRange dateRange, DivideBy divideBy, Function<DateRange, Mono<Long>> function)', 'internal/domain/common/access/admin/controller/base_controller.go:CheckAndQueryBetweenDateFunc()'],

    ['StatisticsRecordDTO(Date date, Long total)', 'internal/domain/common/access/admin/dto/common_dtos.go:StatisticsRecordDTO'],

    // Permission & Expirable utils
    ['ServicePermission(ResponseStatusCode code, String reason)', 'internal/domain/common/permission/service_permission.go:NewServicePermission()'],
    ['isExpired(long creationDate)', 'internal/domain/common/repository/expirable_entity_repository.go:IsExpired()'],
    ['getEntityExpirationDate()', 'internal/domain/common/service/common_services.go:GetEntityExpirationDate()'],  // There are two of these... one repo, one service. First regex could catch both. I'll be careful.
    ['updateGlobalProperties(UserDefinedAttributesProperties properties)', 'internal/domain/common/service/common_services.go:UpdateGlobalProperties()'],
    ['parseAttributesForUpsert(Map<String, Value> userDefinedAttributes)', 'internal/domain/common/service/common_services.go:ParseAttributesForUpsert()'],
    ['isProcessedByResponder(@Nullable RequestStatus status)', 'internal/domain/common/util/expirable_request_inspector.go:IsProcessedByResponder()'],

    // Validator
    ['validResponseAction(ResponseAction action)', 'internal/infra/validator/validator.go:ValidResponseAction()'],
    ['validDeviceType(DeviceType deviceType)', 'internal/infra/validator/validator.go:ValidDeviceType()'],
    ['validProfileAccess(ProfileAccessStrategy value)', 'internal/infra/validator/validator.go:ValidProfileAccess()'],
    ['validRelationshipKey(UserRelationship.Key key)', 'internal/infra/validator/validator.go:ValidRelationshipKey()'],
    ['validRelationshipGroupKey(UserRelationshipGroup.Key key)', 'internal/infra/validator/validator.go:ValidRelationshipGroupKey()'],
    ['validGroupMemberKey(GroupMember.Key key)', 'internal/infra/validator/validator.go:ValidGroupMemberKey()'],
    ['validGroupMemberRole(GroupMemberRole role)', 'internal/infra/validator/validator.go:ValidGroupMemberRole()'],
    ['validGroupBlockedUserKey(GroupBlockedUser.Key key)', 'internal/infra/validator/validator.go:ValidGroupBlockedUserKey()'],
    ['validNewGroupQuestion(NewGroupQuestion question)', 'internal/infra/validator/validator.go:ValidNewGroupQuestion()'],
    ['validGroupQuestionIdAndAnswer(Map.Entry<Long, String> questionIdAndAnswer)', 'internal/infra/validator/validator.go:ValidGroupQuestionIdAndAnswer()'],

    // Conference
    ['CancelMeetingResult(boolean success, @Nullable Meeting meeting)', 'internal/domain/conference/bo/conference_bos.go:CancelMeetingResult'],
    ['UpdateMeetingInvitationResult(boolean updated, @Nullable String accessToken, @Nullable Meeting meeting)', 'internal/domain/conference/bo/conference_bos.go:UpdateMeetingInvitationResult'],
    ['UpdateMeetingResult(boolean success, @Nullable Meeting meeting)', 'internal/domain/conference/bo/conference_bos.go:UpdateMeetingResult'],

    ['handleCreateMeetingRequest()', 'internal/domain/conference/controller/conference_controller.go:HandleCreateMeetingRequest()'],
    ['handleDeleteMeetingRequest()', 'internal/domain/conference/controller/conference_controller.go:HandleDeleteMeetingRequest()'],
    ['handleUpdateMeetingRequest()', 'internal/domain/conference/controller/conference_controller.go:HandleUpdateMeetingRequest()'],
    ['handleQueryMeetingsRequest()', 'internal/domain/conference/controller/conference_controller.go:HandleQueryMeetingsRequest()'],
    ['handleUpdateMeetingInvitationRequest()', 'internal/domain/conference/controller/conference_controller.go:HandleUpdateMeetingInvitationRequest()'],

    ['updateEndDate(Long meetingId, Date endDate)', 'internal/domain/conference/repository/meeting_repository.go:UpdateEndDate()'],
    ['updateCancelDateIfNotCanceled(Long meetingId, Date cancelDate)', 'internal/domain/conference/repository/meeting_repository.go:UpdateCancelDateIfNotCanceled()'],
    ['updateMeeting(Long meetingId, @Nullable String name, @Nullable String intro, @Nullable String password)', 'internal/domain/conference/repository/meeting_repository.go:UpdateMeeting()'],
    ['find(@Nullable Collection<Long> ids, @Nullable Collection<Long> creatorIds, @Nullable Collection<Long> userIds, @Nullable Collection<Long> groupIds, @Nullable Date creationDateStart, @Nullable Date creationDateEnd, @Nullable Integer skip, @Nullable Integer limit)', 'internal/domain/conference/repository/meeting_repository.go:Find()'],
    ['find(@Nullable Collection<Long> ids, @NotNull Long creatorId, @NotNull Long userId, @Nullable Date creationDateStart, @Nullable Date creationDateEnd, @Nullable Integer skip, @Nullable Integer limit)', 'internal/domain/conference/repository/meeting_repository.go:FindByCreatorAndUser()'],

    ['onExtensionStarted(ConferenceServiceProvider extension)', 'internal/domain/conference/service/conference_service.go:OnExtensionStarted()'],
    ['authAndCancelMeeting(@NotNull Long requesterId, @NotNull Long meetingId)', 'internal/domain/conference/service/conference_service.go:AuthAndCancelMeeting()'],
    ['queryMeetingParticipants(@Nullable Long userId, @Nullable Long groupId)', 'internal/domain/conference/service/conference_service.go:QueryMeetingParticipants()'],
    ['authAndUpdateMeeting(@NotNull Long requesterId, @NotNull Long meetingId, @Nullable String name, @Nullable String intro, @Nullable String password)', 'internal/domain/conference/service/conference_service.go:AuthAndUpdateMeeting()'],
    ['authAndUpdateMeetingInvitation(@NotNull Long requesterId, @NotNull Long meetingId, @Nullable String password, @NotNull ResponseAction responseAction)', 'internal/domain/conference/service/conference_service.go:AuthAndUpdateMeetingInvitation()'],
    ['authAndQueryMeetings(@NotNull Long requesterId, @Nullable Set<Long> ids, @Nullable Set<Long> creatorIds, @Nullable Set<Long> userIds, @Nullable Set<Long> groupIds, @Nullable Date creationDateStart, @Nullable Date creationDateEnd, @Nullable Integer skip, @Nullable Integer limit)', 'internal/domain/conference/service/conference_service.go:AuthAndQueryMeetings()'],

    // Conversation
    ['queryConversations(@QueryParam(required = false)', 'internal/domain/conversation/access/admin/controller/conversation_controller.go:QueryConversations()'],
    ['deleteConversations(@QueryParam(required = false)', 'internal/domain/conversation/access/admin/controller/conversation_controller.go:DeleteConversations()'],
    ['updateConversations(@QueryParam(required = false)', 'internal/domain/conversation/access/admin/controller/conversation_controller.go:UpdateConversations()'],

    ['UpdateConversationDTO(Date readDate)', 'internal/domain/conversation/access/admin/dto/conversation_dtos.go:UpdateConversationDTO'],
    ['ConversationsDTO(List<PrivateConversation> privateConversations, List<GroupConversation> groupConversations)', 'internal/domain/conversation/access/admin/dto/conversation_dtos.go:ConversationsDTO'],

    ['handleQueryConversationsRequest()', 'internal/domain/conversation/access/servicerequest/controller/conversation_service_controllers.go:HandleQueryConversationsRequest()'],
    ['handleUpdateTypingStatusRequest()', 'internal/domain/conversation/access/servicerequest/controller/conversation_service_controllers.go:HandleUpdateTypingStatusRequest()'],
    ['handleUpdateConversationRequest()', 'internal/domain/conversation/access/servicerequest/controller/conversation_service_controllers.go:HandleUpdateConversationRequest()'],

    ['handleUpdateConversationSettingsRequest()', 'internal/domain/conversation/access/servicerequest/controller/conversation_service_controllers.go:HandleUpdateConversationSettingsRequest()'],
    ['handleDeleteConversationSettingsRequest()', 'internal/domain/conversation/access/servicerequest/controller/conversation_service_controllers.go:HandleDeleteConversationSettingsRequest()'],
    ['handleQueryConversationSettingsRequest()', 'internal/domain/conversation/access/servicerequest/controller/conversation_service_controllers.go:HandleQueryConversationSettingsRequest()'],

    ['findByIdAndSettingNames(Long ownerId, @Nullable Collection<String> settingNames, @Nullable Date lastUpdatedDateStart)', 'internal/domain/conversation/repository/conversation_settings_repository.go:FindByIdAndSettingNames()'],
    ['findByIdAndSettingNames(Collection<ConversationSettings.Key> keys, @Nullable Collection<String> settingNames, @Nullable Date lastUpdatedDateStart)', 'internal/domain/conversation/repository/conversation_settings_repository.go:FindByIdAndSettingNamesWithKeys()'],
    ['findSettingFields(Long ownerId, Long targetId, Collection<String> includedFields)', 'internal/domain/conversation/repository/conversation_settings_repository.go:FindSettingFields()'],
    ['deleteByOwnerIds(Collection<Long> ownerIds, @Nullable ClientSession clientSession)', 'internal/domain/conversation/repository/conversation_settings_repository.go:DeleteByOwnerIds()'],

    ['deleteMemberConversations(Collection<Long> groupIds, Long memberId, ClientSession session)', 'internal/domain/conversation/repository/group_conversation_repository.go:DeleteMemberConversations()'],

    ['deleteConversationsByOwnerIds(Set<Long> ownerIds, @Nullable ClientSession session)', 'internal/domain/conversation/repository/private_conversation_repository.go:DeleteConversationsByOwnerIds()'],
    ['findConversations(Collection<Long> ownerIds)', 'internal/domain/conversation/repository/private_conversation_repository.go:FindConversations()'],

    ['authAndUpsertGroupConversationReadDate(@NotNull Long groupId, @NotNull Long memberId, @Nullable @PastOrPresent Date readDate)', 'internal/domain/conversation/service/conversation_service.go:AuthAndUpsertGroupConversationReadDate()'],
    ['authAndUpsertPrivateConversationReadDate(@NotNull Long ownerId, @NotNull Long targetId, @Nullable @PastOrPresent Date readDate)', 'internal/domain/conversation/service/conversation_service.go:AuthAndUpsertPrivateConversationReadDate()'],
    ['upsertGroupConversationReadDate(@NotNull Long groupId, @NotNull Long memberId, @Nullable @PastOrPresent Date readDate)', 'internal/domain/conversation/service/conversation_service.go:UpsertGroupConversationReadDate()'],
    ['upsertGroupConversationsReadDate(@NotNull Set<GroupConversation.GroupConversionMemberKey> keys, @Nullable @PastOrPresent Date readDate)', 'internal/domain/conversation/service/conversation_service.go:UpsertGroupConversationsReadDate()'],
    ['upsertPrivateConversationReadDate(@NotNull Long ownerId, @NotNull Long targetId, @Nullable @PastOrPresent Date readDate)', 'internal/domain/conversation/service/conversation_service.go:UpsertPrivateConversationReadDate()'],
    ['upsertPrivateConversationsReadDate(@NotNull Set<PrivateConversation.Key> keys, @Nullable @PastOrPresent Date readDate)', 'internal/domain/conversation/service/conversation_service.go:UpsertPrivateConversationsReadDate()'],
    ['queryPrivateConversationsByOwnerIds(@NotNull Set<Long> ownerIds)', 'internal/domain/conversation/service/conversation_service.go:QueryPrivateConversationsByOwnerIds()'],
    ['deletePrivateConversations(@NotNull Set<PrivateConversation.Key> keys)', 'internal/domain/conversation/service/conversation_service.go:DeletePrivateConversationsByKeys()'],
    ['deletePrivateConversations(@NotNull Set<Long> userIds, @Nullable ClientSession session)', 'internal/domain/conversation/service/conversation_service.go:DeletePrivateConversationsByUserIds()'],
    ['deleteGroupConversations(@Nullable Set<Long> groupIds, @Nullable ClientSession session)', 'internal/domain/conversation/service/conversation_service.go:DeleteGroupConversations()'],
    ['deleteGroupMemberConversations(@NotNull Collection<Long> userIds, @Nullable ClientSession session)', 'internal/domain/conversation/service/conversation_service.go:DeleteGroupMemberConversations()'],
    ['authAndUpdateTypingStatus(@NotNull Long requesterId, boolean isGroupMessage, @NotNull Long toId)', 'internal/domain/conversation/service/conversation_service.go:AuthAndUpdateTypingStatus()'],

    ['upsertPrivateConversationSettings(Long ownerId, Long userId, Map<String, Value> settings)', 'internal/domain/conversation/service/conversation_settings_service.go:UpsertPrivateConversationSettings()'],
    ['upsertGroupConversationSettings(Long ownerId, Long groupId, Map<String, Value> settings)', 'internal/domain/conversation/service/conversation_settings_service.go:UpsertGroupConversationSettings()'],

];

for (let i = 0; i < mappings.length; i++) {
    const [javaSig, goPath] = mappings[i];
    const searchString = `- [ ] \`${javaSig}\``;
    const replacementString = `- [x] \`${javaSig}\` -> [${goPath}](../${goPath.split(':')[0]})`;
    
    content = content.replace(searchString, replacementString);
}


// Manual fallback for getEntityExpirationDate to catch the two occurrences
content = content.replaceAll(
    `- [ ] \`getEntityExpirationDate()\`\n`,
    `- [x] \`getEntityExpirationDate()\` -> [internal/domain/common/repository/expirable_entity_repository.go:GetEntityExpirationDate()](../internal/domain/common/repository/expirable_entity_repository.go)\n`
)

fs.writeFileSync(path, content, 'utf8');
console.log('Update report done.');
