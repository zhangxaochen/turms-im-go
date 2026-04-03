const fs = require('fs');

const path = 'docs/refactor_progress_report.md';
let content = fs.readFileSync(path, 'utf8');

// The mappings we want to apply to [ ] items
const mappings = [
    ['authenticate(@NotNull Long userId, @Nullable String rawPassword)', 'internal/domain/gateway/session/user_service.go:Authenticate(ctx context.Context, userID int64, rawPassword string)'],
    ['getConflictedDeviceTypes(@NotNull @ValidDeviceType DeviceType deviceType)', 'internal/domain/gateway/session/manager/user_simultaneous_login_service.go:GetConflictedDeviceTypes(deviceType protocol.DeviceType)'],
    ['isForbiddenDeviceType(DeviceType deviceType)', 'internal/domain/gateway/session/manager/user_simultaneous_login_service.go:IsForbiddenDeviceType(deviceType protocol.DeviceType)'],
    ['shouldDisconnectLoggingInDeviceIfConflicts()', 'internal/domain/gateway/session/manager/user_simultaneous_login_service.go:ShouldDisconnectLoggingInDeviceIfConflicts()'],
    ['getWsAddress()', 'internal/infra/address/service_address_manager.go:GetWsAddress()'],
    ['getTcpAddress()', 'internal/infra/address/service_address_manager.go:GetTcpAddress()'],
    ['getUdpAddress()', 'internal/infra/address/service_address_manager.go:GetUdpAddress()'],
    ['connect()', 'internal/infra/ldap/ldap_client.go:Connect()'],
    ['bind(boolean useFastBind, String dn, String password)', 'internal/infra/ldap/ldap_client.go:Bind(useFastBind bool, dn string, password string)'],
    ['modify(String dn, List<ModifyOperationChange> changes)', 'internal/infra/ldap/ldap_client.go:Modify(dn string, changes []any)'],
    
    // BerBuffer
    ['skipTag()', 'internal/infra/ldap/asn1/ber_buffer.go:SkipTag()'],
    ['skipTagAndLength()', 'internal/infra/ldap/asn1/ber_buffer.go:SkipTagAndLength()'],
    ['skipTagAndLengthAndValue()', 'internal/infra/ldap/asn1/ber_buffer.go:SkipTagAndLengthAndValue()'],
    ['readTag()', 'internal/infra/ldap/asn1/ber_buffer.go:ReadTag()'],
    ['peekAndCheckTag(int tag)', 'internal/infra/ldap/asn1/ber_buffer.go:PeekAndCheckTag(tag int)'],
    ['skipLength()', 'internal/infra/ldap/asn1/ber_buffer.go:SkipLength()'],
    ['skipLengthAndValue()', 'internal/infra/ldap/asn1/ber_buffer.go:SkipLengthAndValue()'],
    ['writeLength(int length)', 'internal/infra/ldap/asn1/ber_buffer.go:WriteLength(length int)'],
    ['readLength()', 'internal/infra/ldap/asn1/ber_buffer.go:ReadLength()'],
    ['tryReadLengthIfReadable()', 'internal/infra/ldap/asn1/ber_buffer.go:TryReadLengthIfReadable()'],
    ['beginSequence()', 'internal/infra/ldap/asn1/ber_buffer.go:BeginSequence()'],
    ['beginSequence(int tag)', 'internal/infra/ldap/asn1/ber_buffer.go:BeginSequenceWithTag(tag int)'],
    ['endSequence()', 'internal/infra/ldap/asn1/ber_buffer.go:EndSequence()'],
    ['writeBoolean(boolean value)', 'internal/infra/ldap/asn1/ber_buffer.go:WriteBoolean(value bool)'],
    ['writeBoolean(int tag, boolean value)', 'internal/infra/ldap/asn1/ber_buffer.go:WriteBooleanWithTag(tag int, value bool)'],
    ['readBoolean()', 'internal/infra/ldap/asn1/ber_buffer.go:ReadBoolean()'],
    ['writeInteger(int value)', 'internal/infra/ldap/asn1/ber_buffer.go:WriteInteger(value int)'],
    ['writeInteger(int tag, int value)', 'internal/infra/ldap/asn1/ber_buffer.go:WriteIntegerWithTag(tag int, value int)'],
    ['readInteger()', 'internal/infra/ldap/asn1/ber_buffer.go:ReadInteger()'],
    ['readIntWithTag(int tag)', 'internal/infra/ldap/asn1/ber_buffer.go:ReadIntWithTag(tag int)'],
    ['writeOctetString(String value)', 'internal/infra/ldap/asn1/ber_buffer.go:WriteOctetString(value string)'],
    ['writeOctetString(byte[] value)', 'internal/infra/ldap/asn1/ber_buffer.go:WriteOctetStringBytes(value []byte)'],
    ['writeOctetString(int tag, byte[] value)', 'internal/infra/ldap/asn1/ber_buffer.go:WriteOctetStringBytesWithTag(tag int, value []byte)'],
    ['writeOctetString(byte[] value, int start, int length)', 'internal/infra/ldap/asn1/ber_buffer.go:WriteOctetStringBytesRange(value []byte, start int, length int)'],
    ['writeOctetString(int tag, byte[] value, int start, int length)', 'internal/infra/ldap/asn1/ber_buffer.go:WriteOctetStringBytesRangeWithTag(tag int, value []byte, start int, length int)'],
    ['writeOctetString(int tag, String value)', 'internal/infra/ldap/asn1/ber_buffer.go:WriteOctetStringWithTag(tag int, value string)'],
    ['writeOctetStrings(List<String> values)', 'internal/infra/ldap/asn1/ber_buffer.go:WriteOctetStrings(values []string)'],
    ['readOctetString()', 'internal/infra/ldap/asn1/ber_buffer.go:ReadOctetString()'],
    ['readOctetStringWithTag(int tag)', 'internal/infra/ldap/asn1/ber_buffer.go:ReadOctetStringWithTag(tag int)'],
    ['readOctetStringWithLength(int length)', 'internal/infra/ldap/asn1/ber_buffer.go:ReadOctetStringWithLength(length int)'],
    ['writeEnumeration(int value)', 'internal/infra/ldap/asn1/ber_buffer.go:WriteEnumeration(value int)'],
    ['readEnumeration()', 'internal/infra/ldap/asn1/ber_buffer.go:ReadEnumeration()'],
    ['getBytes()', 'internal/infra/ldap/asn1/ber_buffer.go:GetBytes()'],
    ['skipBytes(int length)', 'internal/infra/ldap/asn1/ber_buffer.go:SkipBytes(length int)'],
    ['refCnt()', 'internal/infra/ldap/asn1/ber_buffer.go:RefCnt()'],
    ['retain()', 'internal/infra/ldap/asn1/ber_buffer.go:Retain()'],
    ['retain(int increment)', 'internal/infra/ldap/asn1/ber_buffer.go:RetainIncrement(increment int)'],
    ['touch()', 'internal/infra/ldap/asn1/ber_buffer.go:Touch()'],
    ['touch(Object hint)', 'internal/infra/ldap/asn1/ber_buffer.go:TouchWithHint(hint any)'],
    ['release()', 'internal/infra/ldap/asn1/ber_buffer.go:Release()'],
    ['release(int decrement)', 'internal/infra/ldap/asn1/ber_buffer.go:ReleaseDecrement(decrement int)'],
    ['isReadable(int length)', 'internal/infra/ldap/asn1/ber_buffer.go:IsReadableLen(length int)'],
    ['isReadable()', 'internal/infra/ldap/asn1/ber_buffer.go:IsReadable()'],
    ['isReadableWithEnd(int end)', 'internal/infra/ldap/asn1/ber_buffer.go:IsReadableWithEnd(end int)'],
    ['readerIndex()', 'internal/infra/ldap/asn1/ber_buffer.go:ReaderIndex()'],
    
    // elements
    ['decode(BerBuffer buffer)', 'internal/infra/ldap/element/elements.go:Decode(buffer *asn1.BerBuffer)'],
    ['estimateSize()', 'internal/infra/ldap/element/elements.go:EstimateSize()'],
    ['writeTo(BerBuffer buffer)', 'internal/infra/ldap/element/elements.go:WriteTo(buffer *asn1.BerBuffer)'],
    ['isSuccess()', 'internal/infra/ldap/element/elements.go:IsSuccess()'],
    ['write(BerBuffer buffer, String filter)', 'internal/infra/ldap/element/elements.go:Write(buffer *asn1.BerBuffer, filter string)'],
    ['isComplete()', 'internal/infra/ldap/element/elements.go:IsComplete()'],

    // ApiLoggingContext
    ['shouldLogHeartbeatRequest()', 'internal/infra/logging/api_logging_context.go:ShouldLogHeartbeatRequest()'],

    // Proto Utilities
    ['SimpleTurmsNotification(long requesterId, Integer closeStatus, TurmsRequest.KindCase relayedRequestType)', 'internal/infra/proto/proto_parser.go:NewSimpleTurmsNotification(requesterID int64, closeStatus *int32, relayedRequestType *protocol.TurmsRequest_Kind)'],
    ['SimpleTurmsRequest(long requestId, TurmsRequest.KindCase type, CreateSessionRequest createSessionRequest)', 'internal/infra/proto/proto_parser.go:NewSimpleTurmsRequest(requestID int64, reqType *protocol.TurmsRequest_Kind, createSessionReq *protocol.CreateSessionRequest)'],
    ['parseSimpleNotification(CodedInputStream turmsRequestInputStream)', 'internal/infra/proto/proto_parser.go:ParseSimpleNotification(turmsRequestInputStream []byte)'],
    ['parseSimpleRequest(CodedInputStream turmsRequestInputStream)', 'internal/infra/proto/proto_parser.go:ParseSimpleRequest(turmsRequestInputStream []byte)'],

    // Common requests
    ['turmsRequest()', 'internal/domain/common/dto/client_request.go:TurmsRequest()'],
    ['userId()', 'internal/domain/common/dto/client_request.go:UserId()'],
    ['deviceType()', 'internal/domain/common/dto/client_request.go:DeviceType()'],
    ['clientIp()', 'internal/domain/common/dto/client_request.go:ClientIp()'],
    ['requestId()', 'internal/domain/common/dto/client_request.go:RequestId()'],
    ['equals(Object obj)', 'internal/domain/common/dto/client_request.go:Equals(obj interface{})'],
    ['hashCode()', 'internal/domain/common/dto/client_request.go:HashCode()'],
    
    // RequestHandlerResult overrides
    ['RequestHandlerResult(ResponseStatusCode code, @Nullable String reason, @Nullable TurmsNotification.Data response, List<Notification> notifications)', 'internal/domain/common/dto/request_handler_result.go:NewRequestHandlerResult(...)'],

    // Admin controllers
    ['checkLoginNameAndPassword()', 'internal/domain/admin/access/admin/controller/admin_controllers.go:CheckLoginNameAndPassword()'],
    ['addAdmin(RequestContext requestContext, @RequestBody AddAdminDTO addAdminDTO)', 'internal/domain/admin/access/admin/controller/admin_controllers.go:AddAdmin()'],
    ['queryAdmins(@QueryParam(required = false)', 'internal/domain/admin/access/admin/controller/admin_controllers.go:QueryAdmins()'],
    ['updateAdmins(RequestContext requestContext, Set<Long> ids, @RequestBody UpdateAdminDTO updateAdminDTO)', 'internal/domain/admin/access/admin/controller/admin_controllers.go:UpdateAdmins()'],
    ['deleteAdmins(RequestContext requestContext, Set<Long> ids)', 'internal/domain/admin/access/admin/controller/admin_controllers.go:DeleteAdmins()'],
    ['queryAdminPermissions()', 'internal/domain/admin/access/admin/controller/admin_controllers.go:QueryAdminPermissions()'],
    ['addAdminRole(RequestContext requestContext, @RequestBody AddAdminRoleDTO addAdminRoleDTO)', 'internal/domain/admin/access/admin/controller/admin_controllers.go:AddAdminRole()'],
    ['queryAdminRoles(@QueryParam(required = false)', 'internal/domain/admin/access/admin/controller/admin_controllers.go:QueryAdminRoles()'],
    ['updateAdminRole(RequestContext requestContext, Set<Long> ids, @RequestBody UpdateAdminRoleDTO updateAdminRoleDTO)', 'internal/domain/admin/access/admin/controller/admin_controllers.go:UpdateAdminRole()'],
    ['deleteAdminRoles(RequestContext requestContext, Set<Long> ids)', 'internal/domain/admin/access/admin/controller/admin_controllers.go:DeleteAdminRoles()'],

    // Admin DTOs
    ['AddAdminDTO(String loginName, @SensitiveProperty(SensitiveProperty.Access.ALLOW_DESERIALIZATION)', 'internal/domain/admin/access/admin/dto/admin_dtos.go:AddAdminDTO'],
    ['AddAdminRoleDTO(Long id, String name, Set<String> permissions, Integer rank)', 'internal/domain/admin/access/admin/dto/admin_dtos.go:AddAdminRoleDTO'],
    ['UpdateAdminDTO(@SensitiveProperty(SensitiveProperty.Access.ALLOW_DESERIALIZATION)', 'internal/domain/admin/access/admin/dto/admin_dtos.go:UpdateAdminDTO'],
    ['UpdateAdminRoleDTO(String name, Set<String> permissions, Integer rank)', 'internal/domain/admin/access/admin/dto/admin_dtos.go:UpdateAdminRoleDTO'],
    ['PermissionDTO(String group, AdminPermission permission)', 'internal/domain/admin/access/admin/dto/admin_dtos.go:PermissionDTO'],

    // Admin Repo
    ['updateAdmins(Set<Long> ids, @Nullable byte[] password, @Nullable String displayName, @Nullable Set<Long> roleIds)', 'internal/domain/admin/repository/admin_repository.go:UpdateAdmins()'],
    ['countAdmins(@Nullable Set<Long> ids, @Nullable Set<Long> roleIds)', 'internal/domain/admin/repository/admin_repository.go:CountAdmins()'],
    ['findAdmins(@Nullable Collection<Long> ids, @Nullable Collection<String> loginNames, @Nullable Collection<Long> roleIds, @Nullable Integer page, @Nullable Integer size)', 'internal/domain/admin/repository/admin_repository.go:FindAdmins()'],

    // Admin Role Repo
    ['updateAdminRoles(Set<Long> roleIds, String newName, @Nullable Set<AdminPermission> permissions, @Nullable Integer rank)', 'internal/domain/admin/repository/admin_role_repository.go:UpdateAdminRoles()'],
    ['countAdminRoles(@Nullable Set<Long> ids, @Nullable Set<String> names, @Nullable Set<AdminPermission> includedPermissions, @Nullable Set<Integer> ranks)', 'internal/domain/admin/repository/admin_role_repository.go:CountAdminRoles()'],
    ['findAdminRoles(@Nullable Set<Long> roleIds, @Nullable Set<String> names, @Nullable Set<AdminPermission> includedPermissions, @Nullable Set<Integer> ranks, @Nullable Integer page, @Nullable Integer size)', 'internal/domain/admin/repository/admin_role_repository.go:FindAdminRoles()'],
    ['findAdminRolesByIdsAndRankGreaterThan(@NotNull Collection<Long> roleIds, @Nullable Integer rankGreaterThan)', 'internal/domain/admin/repository/admin_role_repository.go:FindAdminRolesByIdsAndRankGreaterThan()'],
    ['findHighestRankByRoleIds(Set<Long> roleIds)', 'internal/domain/admin/repository/admin_role_repository.go:FindHighestRankByRoleIds()'],

    // Admin Role Service
    ['authAndAddAdminRole(@NotNull Long requesterId, @NotNull Long roleId, @NotNull @NoWhitespace @Size( min = MIN_ROLE_NAME_LIMIT, max = MAX_ROLE_NAME_LIMIT)', 'internal/domain/admin/service/admin_services.go:AuthAndAddAdminRole()'],
    ['addAdminRole(@NotNull Long roleId, @NotNull @NoWhitespace @Size( min = MIN_ROLE_NAME_LIMIT, max = MAX_ROLE_NAME_LIMIT)', 'internal/domain/admin/service/admin_services.go:AddAdminRole()'],
    ['authAndDeleteAdminRoles(@NotNull Long requesterId, @NotEmpty Set<Long> roleIds)', 'internal/domain/admin/service/admin_services.go:AuthAndDeleteAdminRoles()'],
    ['deleteAdminRoles(@NotEmpty Set<Long> roleIds)', 'internal/domain/admin/service/admin_services.go:DeleteAdminRoles()'],
    ['authAndUpdateAdminRoles(@NotNull Long requesterId, @NotEmpty Set<Long> roleIds, @Nullable @NoWhitespace @Size( min = MIN_ROLE_NAME_LIMIT, max = MAX_ROLE_NAME_LIMIT)', 'internal/domain/admin/service/admin_services.go:AuthAndUpdateAdminRoles()'],
    ['updateAdminRole(@NotEmpty Set<Long> roleIds, @Nullable @NoWhitespace @Size( min = MIN_ROLE_NAME_LIMIT, max = MAX_ROLE_NAME_LIMIT)', 'internal/domain/admin/service/admin_services.go:UpdateAdminRole()'],
    ['queryAdminRoles(@Nullable Set<Long> ids, @Nullable Set<String> names, @Nullable Set<AdminPermission> includedPermissions, @Nullable Set<Integer> ranks, @Nullable Integer page, @Nullable Integer size)', 'internal/domain/admin/service/admin_services.go:QueryAdminRoles()'],
    ['queryAndCacheRolesByRoleIdsAndRankGreaterThan(@NotNull Collection<Long> roleIds, @NotNull Integer rankGreaterThan)', 'internal/domain/admin/service/admin_services.go:QueryAndCacheRolesByRoleIdsAndRankGreaterThan()'],
    ['countAdminRoles(@Nullable Set<Long> ids, @Nullable Set<String> names, @Nullable Set<AdminPermission> includedPermissions, @Nullable Set<Integer> ranks)', 'internal/domain/admin/service/admin_services.go:CountAdminRoles()'],
    ['queryHighestRankByAdminId(@NotNull Long adminId)', 'internal/domain/admin/service/admin_services.go:QueryHighestRankByAdminId()'],
    ['queryHighestRankByRoleIds(@NotNull Set<Long> roleIds)', 'internal/domain/admin/service/admin_services.go:QueryHighestRankByRoleIds()'],
    ['isAdminRankHigherThanRank(@NotNull Long adminId, @NotNull Integer rank)', 'internal/domain/admin/service/admin_services.go:IsAdminRankHigherThanRank()'],
    ['queryPermissions(@NotNull Long adminId)', 'internal/domain/admin/service/admin_services.go:QueryPermissions()'],

    // Admin Service
    ['queryRoleIdsByAdminIds(@NotEmpty Set<Long> adminIds)', 'internal/domain/admin/service/admin_services.go:QueryRoleIdsByAdminIds()']
];

for (const [javaSig, goPath] of mappings) {
    const searchString = `- [ ] \`${javaSig}\``;
    const replacementString = `- [x] \`${javaSig}\` -> [${goPath}](../${goPath.split(':')[0]})`;
    
    // Some lines might already been partially matched, we should be careful.
    content = content.replace(searchString, replacementString);
}

fs.writeFileSync(path, content, 'utf8');
console.log('Done mapping checks in ' + path);
