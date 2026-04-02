# Turms Refactoring Progress Report

## Modules

### turms-gateway

> [简述功能]

#### Configurations

- **application-demo.yaml** (`resources/application-demo.yaml`): [简述功能]
- **application-dev.yaml** (`resources/application-dev.yaml`): [简述功能]
- **application-test.yaml** (`resources/application-test.yaml`): [简述功能]
- **application.yaml** (`resources/application.yaml`): [简述功能]

#### Java source tracking

- **TurmsGatewayApplication.java** (`java/im/turms/gateway/TurmsGatewayApplication.java`)
> [简述功能]

  - [ ] `main`

- **ClientRequestDispatcher.java** (`java/im/turms/gateway/access/client/common/ClientRequestDispatcher.java`)
> [简述功能]

  - [ ] `handleRequest`
  - [ ] `handleRequest0`
  - [ ] `handleServiceRequest`

- **IpRequestThrottler.java** (`java/im/turms/gateway/access/client/common/IpRequestThrottler.java`)
> [简述功能]

  - [x] `tryAcquireToken` -> `internal/domain/gateway/access/client/common/ip_request_throttler.go`

- **NotificationFactory.java** (`java/im/turms/gateway/access/client/common/NotificationFactory.java`)
> [简述功能]

  - [ ] `init`
  - [x] `create` -> `internal/domain/gateway/access/client/common/notification_factory.go`
  - [x] `create` -> `internal/domain/gateway/access/client/common/notification_factory.go`
  - [x] `create` -> `internal/domain/gateway/access/client/common/notification_factory.go`
  - [x] `createBuffer` -> `internal/domain/gateway/access/client/common/notification_factory.go`
  - [x] `sessionClosed` -> `internal/domain/gateway/access/client/common/notification_factory.go`

- **RequestHandlerResult.java** (`java/im/turms/gateway/access/client/common/RequestHandlerResult.java`)
> [简述功能]

  - [ ] `RequestHandlerResult`

- **UserSession.java** (`java/im/turms/gateway/access/client/common/UserSession.java`)
> [简述功能]

  - [ ] `setConnection`
  - [x] `setLastHeartbeatRequestTimestampToNow` -> `internal/domain/gateway/session/connection.go`
  - [x] `setLastRequestTimestampToNow` -> `internal/domain/gateway/session/connection.go`
  - [x] `close` -> `internal/storage/redis/redis.go`
  - [x] `isOpen` -> `internal/domain/gateway/session/connection.go`
  - [ ] `isConnected`
  - [ ] `supportsSwitchingToUdp`
  - [x] `sendNotification` -> `internal/domain/gateway/access/router/router.go`
  - [x] `sendNotification` -> `internal/domain/gateway/access/router/router.go`
  - [ ] `acquireDeleteSessionRequestLoggingLock`
  - [ ] `hasPermission`
  - [ ] `toString`

- **UserSessionWrapper.java** (`java/im/turms/gateway/access/client/common/UserSessionWrapper.java`)
> [简述功能]

  - [ ] `getIp`
  - [ ] `getIpStr`
  - [ ] `setUserSession`
  - [ ] `hasUserSession`

- **Policy.java** (`java/im/turms/gateway/access/client/common/authorization/policy/Policy.java`)
> [简述功能]

  - [ ] `Policy`

- **PolicyDeserializer.java** (`java/im/turms/gateway/access/client/common/authorization/policy/PolicyDeserializer.java`)
> [简述功能]

  - [ ] `parse`

- **PolicyStatement.java** (`java/im/turms/gateway/access/client/common/authorization/policy/PolicyStatement.java`)
> [简述功能]

  - [ ] `PolicyStatement`

- **ServiceAvailabilityHandler.java** (`java/im/turms/gateway/access/client/common/channel/ServiceAvailabilityHandler.java`)
> [简述功能]

  - [ ] `channelRegistered`
  - [ ] `exceptionCaught`

- **NetConnection.java** (`java/im/turms/gateway/access/client/common/connection/NetConnection.java`)
> [简述功能]

  - [ ] `getAddress`
  - [ ] `send`
  - [x] `close` -> `internal/storage/redis/redis.go`
  - [x] `close` -> `internal/storage/redis/redis.go`
  - [ ] `switchToUdp`
  - [ ] `tryNotifyClientToRecover`

- **ExtendedHAProxyMessageReader.java** (`java/im/turms/gateway/access/client/tcp/ExtendedHAProxyMessageReader.java`)
> [简述功能]

  - [ ] `channelRead`

- **HAProxyUtil.java** (`java/im/turms/gateway/access/client/tcp/HAProxyUtil.java`)
> [简述功能]

  - [ ] `addProxyProtocolHandlers`
  - [ ] `addProxyProtocolDetectorHandler`

- **TcpConnection.java** (`java/im/turms/gateway/access/client/tcp/TcpConnection.java`)
> [简述功能]

  - [ ] `getAddress`
  - [ ] `send`
  - [x] `close` -> `internal/storage/redis/redis.go`
  - [x] `close` -> `internal/storage/redis/redis.go`

- **TcpServerFactory.java** (`java/im/turms/gateway/access/client/tcp/TcpServerFactory.java`)
> [简述功能]

  - [x] `create` -> `internal/domain/gateway/access/client/common/notification_factory.go`

- **TcpUserSessionAssembler.java** (`java/im/turms/gateway/access/client/tcp/TcpUserSessionAssembler.java`)
> [简述功能]

  - [ ] `getHost`
  - [ ] `getPort`

- **UdpRequestDispatcher.java** (`java/im/turms/gateway/access/client/udp/UdpRequestDispatcher.java`)
> [简述功能]

  - [ ] `sendSignal`

- **UdpSignalResponseBufferPool.java** (`java/im/turms/gateway/access/client/udp/UdpSignalResponseBufferPool.java`)
> [简述功能]

  - [x] `get` -> `internal/domain/gateway/session/sharded_map.go`
  - [x] `get` -> `internal/domain/gateway/session/sharded_map.go`

- **UdpNotification.java** (`java/im/turms/gateway/access/client/udp/dto/UdpNotification.java`)
> [简述功能]

  - [ ] `UdpNotification`

- **UdpRequestType.java** (`java/im/turms/gateway/access/client/udp/dto/UdpRequestType.java`)
> [简述功能]

  - [ ] `parse`
  - [ ] `getNumber`

- **UdpSignalRequest.java** (`java/im/turms/gateway/access/client/udp/dto/UdpSignalRequest.java`)
> [简述功能]

  - [ ] `UdpSignalRequest`

- **HttpForwardedHeaderHandler.java** (`java/im/turms/gateway/access/client/websocket/HttpForwardedHeaderHandler.java`)
> [简述功能]

  - [ ] `apply`

- **WebSocketConnection.java** (`java/im/turms/gateway/access/client/websocket/WebSocketConnection.java`)
> [简述功能]

  - [ ] `getAddress`
  - [ ] `send`
  - [x] `close` -> `internal/storage/redis/redis.go`
  - [x] `close` -> `internal/storage/redis/redis.go`

- **WebSocketServerFactory.java** (`java/im/turms/gateway/access/client/websocket/WebSocketServerFactory.java`)
> [简述功能]

  - [x] `create` -> `internal/domain/gateway/access/client/common/notification_factory.go`

- **NotificationService.java** (`java/im/turms/gateway/domain/notification/service/NotificationService.java`)
> [简述功能]

  - [ ] `sendNotificationToLocalClients`

- **StatisticsService.java** (`java/im/turms/gateway/domain/observation/service/StatisticsService.java`)
> [简述功能]

  - [ ] `countLocalOnlineUsers`

- **ServiceRequestService.java** (`java/im/turms/gateway/domain/servicerequest/service/ServiceRequestService.java`)
> [简述功能]

  - [ ] `handleServiceRequest`

- **SessionController.java** (`java/im/turms/gateway/domain/session/access/admin/controller/SessionController.java`)
> [简述功能]

  - [ ] `deleteSessions`

- **SessionClientController.java** (`java/im/turms/gateway/domain/session/access/client/controller/SessionClientController.java`)
> [简述功能]

  - [ ] `handleDeleteSessionRequest`
  - [ ] `handleCreateSessionRequest`

- **UserLoginInfo.java** (`java/im/turms/gateway/domain/session/bo/UserLoginInfo.java`)
> [简述功能]

  - [ ] `UserLoginInfo`

- **UserPermissionInfo.java** (`java/im/turms/gateway/domain/session/bo/UserPermissionInfo.java`)
> [简述功能]

  - [ ] `UserPermissionInfo`

- **HeartbeatManager.java** (`java/im/turms/gateway/domain/session/manager/HeartbeatManager.java`)
> [简述功能]

  - [ ] `setCloseIdleSessionAfterSeconds`
  - [ ] `setClientHeartbeatIntervalSeconds`
  - [ ] `destroy`
  - [ ] `estimatedSize`
  - [ ] `next`

- **UserSessionsManager.java** (`java/im/turms/gateway/domain/session/manager/UserSessionsManager.java`)
> [简述功能]

  - [ ] `addSessionIfAbsent`
  - [ ] `closeSession`
  - [ ] `pushSessionNotification`
  - [x] `getSession` -> `internal/domain/gateway/session/sharded_map.go`
  - [ ] `countSessions`
  - [ ] `getLoggedInDeviceTypes`

- **UserRepository.java** (`java/im/turms/gateway/domain/session/repository/UserRepository.java`)
> [简述功能]

  - [ ] `findPassword`
  - [x] `isActiveAndNotDeleted` -> `internal/domain/user/repository/user_repository.go`

- **HttpSessionIdentityAccessManager.java** (`java/im/turms/gateway/domain/session/service/HttpSessionIdentityAccessManager.java`)
> [简述功能]

  - [ ] `verifyAndGrant`

- **JwtSessionIdentityAccessManager.java** (`java/im/turms/gateway/domain/session/service/JwtSessionIdentityAccessManager.java`)
> [简述功能]

  - [ ] `verifyAndGrant`

- **LdapSessionIdentityAccessManager.java** (`java/im/turms/gateway/domain/session/service/LdapSessionIdentityAccessManager.java`)
> [简述功能]

  - [ ] `verifyAndGrant`

- **NoopSessionIdentityAccessManager.java** (`java/im/turms/gateway/domain/session/service/NoopSessionIdentityAccessManager.java`)
> [简述功能]

  - [ ] `verifyAndGrant`

- **PasswordSessionIdentityAccessManager.java** (`java/im/turms/gateway/domain/session/service/PasswordSessionIdentityAccessManager.java`)
> [简述功能]

  - [ ] `verifyAndGrant`
  - [ ] `updateGlobalProperties`

- **SessionIdentityAccessManager.java** (`java/im/turms/gateway/domain/session/service/SessionIdentityAccessManager.java`)
> [简述功能]

  - [ ] `verifyAndGrant`

- **SessionService.java** (`java/im/turms/gateway/domain/session/service/SessionService.java`)
> [简述功能]

  - [ ] `destroy`
  - [ ] `handleHeartbeatUpdateRequest`
  - [ ] `handleLoginRequest`
  - [ ] `closeLocalSessions`
  - [ ] `closeLocalSessions`
  - [ ] `closeLocalSession`
  - [ ] `closeLocalSession`
  - [ ] `closeLocalSession`
  - [ ] `closeLocalSessions`
  - [ ] `authAndCloseLocalSession`
  - [ ] `closeAllLocalSessions`
  - [ ] `closeLocalSession`
  - [ ] `closeLocalSession`
  - [ ] `getSessions`
  - [ ] `authAndUpdateHeartbeatTimestamp`
  - [ ] `tryRegisterOnlineUser`
  - [ ] `getUserSessionsManager`
  - [ ] `getLocalUserSession`
  - [ ] `getLocalUserSession`
  - [ ] `countLocalOnlineUsers`
  - [ ] `onSessionEstablished`
  - [ ] `addOnSessionClosedListeners`
  - [ ] `invokeGoOnlineHandlers`

- **UserService.java** (`java/im/turms/gateway/domain/session/service/UserService.java`)
> [简述功能]

  - [ ] `authenticate`
  - [x] `isActiveAndNotDeleted` -> `internal/domain/user/repository/user_repository.go`

- **UserSimultaneousLoginService.java** (`java/im/turms/gateway/domain/session/service/UserSimultaneousLoginService.java`)
> [简述功能]

  - [ ] `getConflictedDeviceTypes`
  - [ ] `isForbiddenDeviceType`
  - [ ] `shouldDisconnectLoggingInDeviceIfConflicts`

- **ServiceAddressManager.java** (`java/im/turms/gateway/infra/address/ServiceAddressManager.java`)
> [简述功能]

  - [ ] `getWsAddress`
  - [ ] `getTcpAddress`
  - [ ] `getUdpAddress`

- **LdapClient.java** (`java/im/turms/gateway/infra/ldap/LdapClient.java`)
> [简述功能]

  - [ ] `isConnected`
  - [ ] `connect`
  - [ ] `bind`
  - [x] `search` -> `internal/storage/elasticsearch/elasticsearch_client.go`
  - [ ] `modify`

- **BerBuffer.java** (`java/im/turms/gateway/infra/ldap/asn1/BerBuffer.java`)
> [简述功能]

  - [ ] `skipTag`
  - [ ] `skipTagAndLength`
  - [ ] `skipTagAndLengthAndValue`
  - [ ] `readTag`
  - [ ] `peekAndCheckTag`
  - [ ] `skipLength`
  - [ ] `skipLengthAndValue`
  - [ ] `writeLength`
  - [ ] `readLength`
  - [ ] `tryReadLengthIfReadable`
  - [ ] `beginSequence`
  - [ ] `beginSequence`
  - [ ] `endSequence`
  - [ ] `writeBoolean`
  - [ ] `writeBoolean`
  - [ ] `readBoolean`
  - [ ] `writeInteger`
  - [ ] `writeInteger`
  - [ ] `readInteger`
  - [ ] `readIntWithTag`
  - [ ] `writeOctetString`
  - [ ] `writeOctetString`
  - [ ] `writeOctetString`
  - [ ] `writeOctetString`
  - [ ] `writeOctetString`
  - [ ] `writeOctetString`
  - [ ] `writeOctetStrings`
  - [ ] `readOctetString`
  - [ ] `readOctetStringWithTag`
  - [ ] `readOctetStringWithLength`
  - [ ] `writeEnumeration`
  - [ ] `readEnumeration`
  - [ ] `getBytes`
  - [ ] `skipBytes`
  - [x] `close` -> `internal/storage/redis/redis.go`
  - [ ] `refCnt`
  - [ ] `retain`
  - [ ] `retain`
  - [ ] `touch`
  - [ ] `touch`
  - [ ] `release`
  - [ ] `release`
  - [ ] `isReadable`
  - [ ] `isReadable`
  - [ ] `isReadableWithEnd`
  - [ ] `readerIndex`

- **Attribute.java** (`java/im/turms/gateway/infra/ldap/element/common/Attribute.java`)
> [简述功能]

  - [x] `isEmpty` -> `internal/domain/gateway/session/sharded_map.go`
  - [ ] `decode`

- **LdapMessage.java** (`java/im/turms/gateway/infra/ldap/element/common/LdapMessage.java`)
> [简述功能]

  - [ ] `estimateSize`
  - [ ] `writeTo`

- **LdapResult.java** (`java/im/turms/gateway/infra/ldap/element/common/LdapResult.java`)
> [简述功能]

  - [ ] `isSuccess`

- **Control.java** (`java/im/turms/gateway/infra/ldap/element/common/control/Control.java`)
> [简述功能]

  - [ ] `decode`

- **BindRequest.java** (`java/im/turms/gateway/infra/ldap/element/operation/bind/BindRequest.java`)
> [简述功能]

  - [ ] `estimateSize`
  - [ ] `writeTo`

- **BindResponse.java** (`java/im/turms/gateway/infra/ldap/element/operation/bind/BindResponse.java`)
> [简述功能]

  - [ ] `decode`

- **ModifyRequest.java** (`java/im/turms/gateway/infra/ldap/element/operation/modify/ModifyRequest.java`)
> [简述功能]

  - [ ] `estimateSize`
  - [ ] `writeTo`

- **ModifyResponse.java** (`java/im/turms/gateway/infra/ldap/element/operation/modify/ModifyResponse.java`)
> [简述功能]

  - [ ] `decode`

- **Filter.java** (`java/im/turms/gateway/infra/ldap/element/operation/search/Filter.java`)
> [简述功能]

  - [ ] `write`

- **SearchRequest.java** (`java/im/turms/gateway/infra/ldap/element/operation/search/SearchRequest.java`)
> [简述功能]

  - [ ] `estimateSize`
  - [ ] `writeTo`

- **SearchResult.java** (`java/im/turms/gateway/infra/ldap/element/operation/search/SearchResult.java`)
> [简述功能]

  - [ ] `decode`
  - [ ] `isComplete`

- **ApiLoggingContext.java** (`java/im/turms/gateway/infra/logging/ApiLoggingContext.java`)
> [简述功能]

  - [ ] `shouldLogHeartbeatRequest`
  - [x] `shouldLogRequest` -> `internal/infra/logging/api_logging_context.go`
  - [x] `shouldLogNotification` -> `internal/infra/logging/api_logging_context.go`

- **ClientApiLogging.java** (`java/im/turms/gateway/infra/logging/ClientApiLogging.java`)
> [简述功能]

  - [x] `log` -> `internal/infra/logging/client_api_logging.go`
  - [x] `log` -> `internal/infra/logging/client_api_logging.go`
  - [x] `log` -> `internal/infra/logging/client_api_logging.go`

- **NotificationLoggingManager.java** (`java/im/turms/gateway/infra/logging/NotificationLoggingManager.java`)
> [简述功能]

  - [x] `log` -> `internal/infra/logging/client_api_logging.go`

- **SimpleTurmsNotification.java** (`java/im/turms/gateway/infra/proto/SimpleTurmsNotification.java`)
> [简述功能]

  - [ ] `SimpleTurmsNotification`

- **SimpleTurmsRequest.java** (`java/im/turms/gateway/infra/proto/SimpleTurmsRequest.java`)
> [简述功能]

  - [ ] `SimpleTurmsRequest`
  - [ ] `toString`

- **TurmsNotificationParser.java** (`java/im/turms/gateway/infra/proto/TurmsNotificationParser.java`)
> [简述功能]

  - [ ] `parseSimpleNotification`

- **TurmsRequestParser.java** (`java/im/turms/gateway/infra/proto/TurmsRequestParser.java`)
> [简述功能]

  - [ ] `parseSimpleRequest`

- **MongoConfig.java** (`java/im/turms/gateway/storage/mongo/MongoConfig.java`)
> [简述功能]

  - [x] `adminMongoClient` -> `internal/storage/mongo/mongo_config.go`
  - [x] `userMongoClient` -> `internal/storage/mongo/mongo_config.go`
  - [ ] `mongoDataGenerator`

### turms-service

> [简述功能]

#### Configurations

- **application-demo.yaml** (`resources/application-demo.yaml`): [简述功能]
- **application-dev.yaml** (`resources/application-dev.yaml`): [简述功能]
- **application-test.yaml** (`resources/application-test.yaml`): [简述功能]
- **application.yaml** (`resources/application.yaml`): [简述功能]

#### Java source tracking

- **TurmsServiceApplication.java** (`java/im/turms/service/TurmsServiceApplication.java`)
> [简述功能]

  - [ ] `main`

- **ServiceRequestDispatcher.java** (`java/im/turms/service/access/servicerequest/dispatcher/ServiceRequestDispatcher.java`)
> [简述功能]

  - [x] `dispatch` -> `internal/domain/common/infra/cluster/rpc/router.go`

- **ClientRequest.java** (`java/im/turms/service/access/servicerequest/dto/ClientRequest.java`)
> [简述功能]

  - [ ] `toString`
  - [ ] `turmsRequest`
  - [ ] `userId`
  - [ ] `deviceType`
  - [ ] `clientIp`
  - [ ] `requestId`
  - [ ] `equals`
  - [ ] `hashCode`

- **RequestHandlerResult.java** (`java/im/turms/service/access/servicerequest/dto/RequestHandlerResult.java`)
> [简述功能]

  - [ ] `RequestHandlerResult`
  - [ ] `toString`
  - [ ] `of`
  - [ ] `of`
  - [ ] `of`
  - [ ] `of`
  - [ ] `of`
  - [ ] `of`
  - [ ] `of`
  - [ ] `of`
  - [ ] `of`
  - [ ] `of`
  - [ ] `of`
  - [ ] `of`
  - [ ] `of`
  - [ ] `of`
  - [ ] `of`
  - [ ] `ofDataLong`
  - [ ] `ofDataLong`
  - [ ] `ofDataLong`
  - [ ] `ofDataLong`
  - [ ] `ofDataLong`
  - [ ] `ofDataLongs`
  - [ ] `Notification`
  - [ ] `of`
  - [ ] `of`
  - [ ] `of`
  - [ ] `toString`

- **AdminController.java** (`java/im/turms/service/domain/admin/access/admin/controller/AdminController.java`)
> [简述功能]

  - [ ] `checkLoginNameAndPassword`
  - [ ] `addAdmin`
  - [ ] `queryAdmins`
  - [ ] `queryAdmins`
  - [ ] `updateAdmins`
  - [ ] `deleteAdmins`

- **AdminPermissionController.java** (`java/im/turms/service/domain/admin/access/admin/controller/AdminPermissionController.java`)
> [简述功能]

  - [ ] `queryAdminPermissions`

- **AdminRoleController.java** (`java/im/turms/service/domain/admin/access/admin/controller/AdminRoleController.java`)
> [简述功能]

  - [ ] `addAdminRole`
  - [ ] `queryAdminRoles`
  - [ ] `queryAdminRoles`
  - [ ] `updateAdminRole`
  - [ ] `deleteAdminRoles`

- **AddAdminDTO.java** (`java/im/turms/service/domain/admin/access/admin/dto/request/AddAdminDTO.java`)
> [简述功能]

  - [ ] `AddAdminDTO`
  - [ ] `toString`

- **AddAdminRoleDTO.java** (`java/im/turms/service/domain/admin/access/admin/dto/request/AddAdminRoleDTO.java`)
> [简述功能]

  - [ ] `AddAdminRoleDTO`

- **UpdateAdminDTO.java** (`java/im/turms/service/domain/admin/access/admin/dto/request/UpdateAdminDTO.java`)
> [简述功能]

  - [ ] `UpdateAdminDTO`
  - [ ] `toString`

- **UpdateAdminRoleDTO.java** (`java/im/turms/service/domain/admin/access/admin/dto/request/UpdateAdminRoleDTO.java`)
> [简述功能]

  - [ ] `UpdateAdminRoleDTO`

- **PermissionDTO.java** (`java/im/turms/service/domain/admin/access/admin/dto/response/PermissionDTO.java`)
> [简述功能]

  - [ ] `PermissionDTO`

- **AdminRepository.java** (`java/im/turms/service/domain/admin/repository/AdminRepository.java`)
> [简述功能]

  - [ ] `updateAdmins`
  - [ ] `countAdmins`
  - [ ] `findAdmins`

- **AdminRoleRepository.java** (`java/im/turms/service/domain/admin/repository/AdminRoleRepository.java`)
> [简述功能]

  - [ ] `updateAdminRoles`
  - [ ] `countAdminRoles`
  - [ ] `findAdminRoles`
  - [ ] `findAdminRolesByIdsAndRankGreaterThan`
  - [ ] `findHighestRankByRoleIds`

- **AdminRoleService.java** (`java/im/turms/service/domain/admin/service/AdminRoleService.java`)
> [简述功能]

  - [ ] `authAndAddAdminRole`
  - [ ] `addAdminRole`
  - [ ] `authAndDeleteAdminRoles`
  - [ ] `deleteAdminRoles`
  - [ ] `authAndUpdateAdminRoles`
  - [ ] `updateAdminRole`
  - [ ] `queryAdminRoles`
  - [ ] `queryAndCacheRolesByRoleIdsAndRankGreaterThan`
  - [ ] `countAdminRoles`
  - [ ] `queryHighestRankByAdminId`
  - [ ] `queryHighestRankByRoleIds`
  - [ ] `isAdminRankHigherThanRank`
  - [ ] `queryPermissions`

- **AdminService.java** (`java/im/turms/service/domain/admin/service/AdminService.java`)
> [简述功能]

  - [ ] `queryRoleIdsByAdminIds`
  - [ ] `authAndAddAdmin`
  - [ ] `addAdmin`
  - [ ] `queryAdmins`
  - [ ] `authAndDeleteAdmins`
  - [ ] `authAndUpdateAdmins`
  - [ ] `updateAdmins`
  - [ ] `countAdmins`
  - [ ] `errorRequesterNotExist`

- **IpBlocklistController.java** (`java/im/turms/service/domain/blocklist/access/admin/controller/IpBlocklistController.java`)
> [简述功能]

  - [ ] `addBlockedIps`
  - [ ] `queryBlockedIps`
  - [ ] `queryBlockedIps`
  - [ ] `deleteBlockedIps`

- **UserBlocklistController.java** (`java/im/turms/service/domain/blocklist/access/admin/controller/UserBlocklistController.java`)
> [简述功能]

  - [ ] `addBlockedUserIds`
  - [x] `queryBlockedUsers` -> `internal/domain/group/service/group_blocklist_service.go`
  - [x] `queryBlockedUsers` -> `internal/domain/group/service/group_blocklist_service.go`
  - [ ] `deleteBlockedUserIds`

- **AddBlockedIpsDTO.java** (`java/im/turms/service/domain/blocklist/access/admin/dto/request/AddBlockedIpsDTO.java`)
> [简述功能]

  - [ ] `AddBlockedIpsDTO`

- **AddBlockedUserIdsDTO.java** (`java/im/turms/service/domain/blocklist/access/admin/dto/request/AddBlockedUserIdsDTO.java`)
> [简述功能]

  - [ ] `AddBlockedUserIdsDTO`

- **BlockedIpDTO.java** (`java/im/turms/service/domain/blocklist/access/admin/dto/response/BlockedIpDTO.java`)
> [简述功能]

  - [ ] `BlockedIpDTO`

- **BlockedUserDTO.java** (`java/im/turms/service/domain/blocklist/access/admin/dto/response/BlockedUserDTO.java`)
> [简述功能]

  - [ ] `BlockedUserDTO`

- **BlockedClientSerializer.java** (`java/im/turms/service/domain/blocklist/codec/BlockedClientSerializer.java`)
> [简述功能]

  - [x] `serialize` -> `internal/storage/elasticsearch/model/elasticsearch_model.go`

- **MemberController.java** (`java/im/turms/service/domain/cluster/access/admin/controller/MemberController.java`)
> [简述功能]

  - [ ] `queryMembers`
  - [ ] `removeMembers`
  - [ ] `addMember`
  - [ ] `updateMember`
  - [ ] `queryLeader`
  - [ ] `electNewLeader`

- **SettingController.java** (`java/im/turms/service/domain/cluster/access/admin/controller/SettingController.java`)
> [简述功能]

  - [ ] `queryClusterSettings`
  - [ ] `updateClusterSettings`
  - [ ] `queryClusterConfigMetadata`

- **AddMemberDTO.java** (`java/im/turms/service/domain/cluster/access/admin/dto/request/AddMemberDTO.java`)
> [简述功能]

  - [ ] `AddMemberDTO`

- **UpdateMemberDTO.java** (`java/im/turms/service/domain/cluster/access/admin/dto/request/UpdateMemberDTO.java`)
> [简述功能]

  - [ ] `UpdateMemberDTO`

- **SettingsDTO.java** (`java/im/turms/service/domain/cluster/access/admin/dto/response/SettingsDTO.java`)
> [简述功能]

  - [ ] `SettingsDTO`

- **BaseController.java** (`java/im/turms/service/domain/common/access/admin/controller/BaseController.java`)
> [简述功能]

  - [ ] `getPageSize`
  - [ ] `queryBetweenDate`
  - [ ] `queryBetweenDate`
  - [ ] `checkAndQueryBetweenDate`
  - [ ] `checkAndQueryBetweenDate`

- **StatisticsRecordDTO.java** (`java/im/turms/service/domain/common/access/admin/dto/response/StatisticsRecordDTO.java`)
> [简述功能]

  - [ ] `StatisticsRecordDTO`

- **ServicePermission.java** (`java/im/turms/service/domain/common/permission/ServicePermission.java`)
> [简述功能]

  - [ ] `ServicePermission`
  - [x] `get` -> `internal/domain/gateway/session/sharded_map.go`
  - [x] `get` -> `internal/domain/gateway/session/sharded_map.go`

- **ExpirableEntityRepository.java** (`java/im/turms/service/domain/common/repository/ExpirableEntityRepository.java`)
> [简述功能]

  - [ ] `isExpired`
  - [ ] `getEntityExpirationDate`
  - [x] `deleteExpiredData` -> `internal/domain/user/repository/user_friend_request_repository.go`
  - [x] `findMany` -> `internal/domain/user/repository/user_repository.go`
  - [x] `findMany` -> `internal/domain/user/repository/user_repository.go`

- **ExpirableEntityService.java** (`java/im/turms/service/domain/common/service/ExpirableEntityService.java`)
> [简述功能]

  - [ ] `getEntityExpirationDate`

- **UserDefinedAttributesService.java** (`java/im/turms/service/domain/common/service/UserDefinedAttributesService.java`)
> [简述功能]

  - [ ] `updateGlobalProperties`
  - [ ] `parseAttributesForUpsert`

- **ExpirableRequestInspector.java** (`java/im/turms/service/domain/common/util/ExpirableRequestInspector.java`)
> [简述功能]

  - [ ] `isProcessedByResponder`

- **DataValidator.java** (`java/im/turms/service/domain/common/validation/DataValidator.java`)
> [简述功能]

  - [ ] `validRequestStatus`
  - [ ] `validResponseAction`
  - [ ] `validDeviceType`
  - [ ] `validProfileAccess`
  - [ ] `validRelationshipKey`
  - [ ] `validRelationshipGroupKey`
  - [ ] `validGroupMemberKey`
  - [ ] `validGroupMemberRole`
  - [ ] `validGroupBlockedUserKey`
  - [ ] `validNewGroupQuestion`
  - [ ] `validGroupQuestionIdAndAnswer`

- **CancelMeetingResult.java** (`java/im/turms/service/domain/conference/bo/CancelMeetingResult.java`)
> [简述功能]

  - [ ] `CancelMeetingResult`

- **UpdateMeetingInvitationResult.java** (`java/im/turms/service/domain/conference/bo/UpdateMeetingInvitationResult.java`)
> [简述功能]

  - [ ] `UpdateMeetingInvitationResult`

- **UpdateMeetingResult.java** (`java/im/turms/service/domain/conference/bo/UpdateMeetingResult.java`)
> [简述功能]

  - [ ] `UpdateMeetingResult`

- **ConferenceServiceController.java** (`java/im/turms/service/domain/conference/controller/ConferenceServiceController.java`)
> [简述功能]

  - [ ] `handleCreateMeetingRequest`
  - [ ] `handleDeleteMeetingRequest`
  - [ ] `handleUpdateMeetingRequest`
  - [ ] `handleQueryMeetingsRequest`
  - [ ] `handleUpdateMeetingInvitationRequest`

- **MeetingRepository.java** (`java/im/turms/service/domain/conference/repository/MeetingRepository.java`)
> [简述功能]

  - [ ] `updateEndDate`
  - [ ] `updateCancelDateIfNotCanceled`
  - [ ] `updateMeeting`
  - [ ] `find`
  - [ ] `find`

- **ConferenceService.java** (`java/im/turms/service/domain/conference/service/ConferenceService.java`)
> [简述功能]

  - [ ] `onExtensionStarted`
  - [ ] `authAndCancelMeeting`
  - [ ] `queryMeetingParticipants`
  - [ ] `authAndUpdateMeeting`
  - [ ] `authAndUpdateMeetingInvitation`
  - [ ] `authAndQueryMeetings`

- **ConversationController.java** (`java/im/turms/service/domain/conversation/access/admin/controller/ConversationController.java`)
> [简述功能]

  - [ ] `queryConversations`
  - [ ] `deleteConversations`
  - [ ] `updateConversations`

- **UpdateConversationDTO.java** (`java/im/turms/service/domain/conversation/access/admin/dto/request/UpdateConversationDTO.java`)
> [简述功能]

  - [ ] `UpdateConversationDTO`

- **ConversationsDTO.java** (`java/im/turms/service/domain/conversation/access/admin/dto/response/ConversationsDTO.java`)
> [简述功能]

  - [ ] `ConversationsDTO`

- **ConversationServiceController.java** (`java/im/turms/service/domain/conversation/access/servicerequest/controller/ConversationServiceController.java`)
> [简述功能]

  - [ ] `handleQueryConversationsRequest`
  - [ ] `handleUpdateTypingStatusRequest`
  - [ ] `handleUpdateConversationRequest`

- **ConversationSettingsServiceController.java** (`java/im/turms/service/domain/conversation/access/servicerequest/controller/ConversationSettingsServiceController.java`)
> [简述功能]

  - [ ] `handleUpdateConversationSettingsRequest`
  - [ ] `handleDeleteConversationSettingsRequest`
  - [ ] `handleQueryConversationSettingsRequest`

- **ConversationSettingsRepository.java** (`java/im/turms/service/domain/conversation/repository/ConversationSettingsRepository.java`)
> [简述功能]

  - [x] `upsertSettings` -> `internal/domain/user/service/user_settings_service.go`
  - [x] `unsetSettings` -> `internal/domain/user/service/user_settings_service.go`
  - [x] `findByIdAndSettingNames` -> `internal/domain/user/repository/user_settings_repository.go`
  - [x] `findByIdAndSettingNames` -> `internal/domain/user/repository/user_settings_repository.go`
  - [ ] `findSettingFields`
  - [ ] `deleteByOwnerIds`

- **GroupConversationRepository.java** (`java/im/turms/service/domain/conversation/repository/GroupConversationRepository.java`)
> [简述功能]

  - [x] `upsert` -> `internal/domain/user/repository/user_relationship_repository.go`
  - [x] `upsert` -> `internal/domain/user/repository/user_relationship_repository.go`
  - [ ] `deleteMemberConversations`

- **PrivateConversationRepository.java** (`java/im/turms/service/domain/conversation/repository/PrivateConversationRepository.java`)
> [简述功能]

  - [x] `upsert` -> `internal/domain/user/repository/user_relationship_repository.go`
  - [ ] `deleteConversationsByOwnerIds`
  - [ ] `findConversations`

- **ConversationService.java** (`java/im/turms/service/domain/conversation/service/ConversationService.java`)
> [简述功能]

  - [ ] `authAndUpsertGroupConversationReadDate`
  - [ ] `authAndUpsertPrivateConversationReadDate`
  - [ ] `upsertGroupConversationReadDate`
  - [ ] `upsertGroupConversationsReadDate`
  - [ ] `upsertPrivateConversationReadDate`
  - [ ] `upsertPrivateConversationsReadDate`
  - [x] `queryGroupConversations` -> `internal/domain/conversation/service/conversation_service.go`
  - [ ] `queryPrivateConversationsByOwnerIds`
  - [x] `queryPrivateConversations` -> `internal/domain/conversation/service/conversation_service.go`
  - [x] `queryPrivateConversations` -> `internal/domain/conversation/service/conversation_service.go`
  - [ ] `deletePrivateConversations`
  - [ ] `deletePrivateConversations`
  - [ ] `deleteGroupConversations`
  - [ ] `deleteGroupMemberConversations`
  - [ ] `authAndUpdateTypingStatus`

- **ConversationSettingsService.java** (`java/im/turms/service/domain/conversation/service/ConversationSettingsService.java`)
> [简述功能]

  - [ ] `upsertPrivateConversationSettings`
  - [ ] `upsertGroupConversationSettings`
  - [x] `deleteSettings` -> `internal/domain/user/service/user_settings_service.go`
  - [x] `unsetSettings` -> `internal/domain/user/service/user_settings_service.go`
  - [x] `querySettings` -> `internal/domain/user/service/user_settings_service.go`

- **GroupBlocklistController.java** (`java/im/turms/service/domain/group/access/admin/controller/GroupBlocklistController.java`)
> [简述功能]

  - [ ] `addGroupBlockedUser`
  - [ ] `queryGroupBlockedUsers`
  - [ ] `queryGroupBlockedUsers`
  - [ ] `updateGroupBlockedUsers`
  - [ ] `deleteGroupBlockedUsers`

- **GroupController.java** (`java/im/turms/service/domain/group/access/admin/controller/GroupController.java`)
> [简述功能]

  - [ ] `addGroup`
  - [ ] `queryGroups`
  - [ ] `queryGroups`
  - [x] `countGroups` -> `internal/domain/user/repository/user_relationship_group_member_repository.go`
  - [ ] `updateGroups`
  - [ ] `deleteGroups`

- **GroupInvitationController.java** (`java/im/turms/service/domain/group/access/admin/controller/GroupInvitationController.java`)
> [简述功能]

  - [ ] `addGroupInvitation`
  - [ ] `queryGroupInvitations`
  - [ ] `queryGroupInvitations`
  - [ ] `updateGroupInvitations`
  - [ ] `deleteGroupInvitations`

- **GroupJoinRequestController.java** (`java/im/turms/service/domain/group/access/admin/controller/GroupJoinRequestController.java`)
> [简述功能]

  - [ ] `addGroupJoinRequest`
  - [ ] `queryGroupJoinRequests`
  - [ ] `queryGroupJoinRequests`
  - [ ] `updateGroupJoinRequests`
  - [ ] `deleteGroupJoinRequests`

- **GroupMemberController.java** (`java/im/turms/service/domain/group/access/admin/controller/GroupMemberController.java`)
> [简述功能]

  - [x] `addGroupMember` -> `internal/domain/group/service/group_member_service.go`
  - [ ] `queryGroupMembers`
  - [ ] `queryGroupMembers`
  - [ ] `updateGroupMembers`
  - [ ] `deleteGroupMembers`

- **GroupQuestionController.java** (`java/im/turms/service/domain/group/access/admin/controller/GroupQuestionController.java`)
> [简述功能]

  - [ ] `queryGroupJoinQuestions`
  - [ ] `queryGroupJoinQuestions`
  - [ ] `addGroupJoinQuestion`
  - [ ] `updateGroupJoinQuestions`
  - [ ] `deleteGroupJoinQuestions`

- **GroupTypeController.java** (`java/im/turms/service/domain/group/access/admin/controller/GroupTypeController.java`)
> [简述功能]

  - [ ] `addGroupType`
  - [ ] `queryGroupTypes`
  - [ ] `queryGroupTypes`
  - [x] `updateGroupType` -> `internal/domain/group/repository/group_type_repository.go`
  - [ ] `deleteGroupType`

- **AddGroupBlockedUserDTO.java** (`java/im/turms/service/domain/group/access/admin/dto/request/AddGroupBlockedUserDTO.java`)
> [简述功能]

  - [ ] `AddGroupBlockedUserDTO`

- **AddGroupDTO.java** (`java/im/turms/service/domain/group/access/admin/dto/request/AddGroupDTO.java`)
> [简述功能]

  - [ ] `AddGroupDTO`

- **AddGroupInvitationDTO.java** (`java/im/turms/service/domain/group/access/admin/dto/request/AddGroupInvitationDTO.java`)
> [简述功能]

  - [ ] `AddGroupInvitationDTO`

- **AddGroupJoinQuestionDTO.java** (`java/im/turms/service/domain/group/access/admin/dto/request/AddGroupJoinQuestionDTO.java`)
> [简述功能]

  - [ ] `AddGroupJoinQuestionDTO`

- **AddGroupJoinRequestDTO.java** (`java/im/turms/service/domain/group/access/admin/dto/request/AddGroupJoinRequestDTO.java`)
> [简述功能]

  - [ ] `AddGroupJoinRequestDTO`

- **AddGroupMemberDTO.java** (`java/im/turms/service/domain/group/access/admin/dto/request/AddGroupMemberDTO.java`)
> [简述功能]

  - [ ] `AddGroupMemberDTO`

- **AddGroupTypeDTO.java** (`java/im/turms/service/domain/group/access/admin/dto/request/AddGroupTypeDTO.java`)
> [简述功能]

  - [ ] `AddGroupTypeDTO`

- **UpdateGroupBlockedUserDTO.java** (`java/im/turms/service/domain/group/access/admin/dto/request/UpdateGroupBlockedUserDTO.java`)
> [简述功能]

  - [ ] `UpdateGroupBlockedUserDTO`

- **UpdateGroupDTO.java** (`java/im/turms/service/domain/group/access/admin/dto/request/UpdateGroupDTO.java`)
> [简述功能]

  - [ ] `UpdateGroupDTO`

- **UpdateGroupInvitationDTO.java** (`java/im/turms/service/domain/group/access/admin/dto/request/UpdateGroupInvitationDTO.java`)
> [简述功能]

  - [ ] `UpdateGroupInvitationDTO`

- **UpdateGroupJoinQuestionDTO.java** (`java/im/turms/service/domain/group/access/admin/dto/request/UpdateGroupJoinQuestionDTO.java`)
> [简述功能]

  - [ ] `UpdateGroupJoinQuestionDTO`

- **UpdateGroupJoinRequestDTO.java** (`java/im/turms/service/domain/group/access/admin/dto/request/UpdateGroupJoinRequestDTO.java`)
> [简述功能]

  - [ ] `UpdateGroupJoinRequestDTO`

- **UpdateGroupMemberDTO.java** (`java/im/turms/service/domain/group/access/admin/dto/request/UpdateGroupMemberDTO.java`)
> [简述功能]

  - [ ] `UpdateGroupMemberDTO`

- **UpdateGroupTypeDTO.java** (`java/im/turms/service/domain/group/access/admin/dto/request/UpdateGroupTypeDTO.java`)
> [简述功能]

  - [ ] `UpdateGroupTypeDTO`

- **GroupStatisticsDTO.java** (`java/im/turms/service/domain/group/access/admin/dto/response/GroupStatisticsDTO.java`)
> [简述功能]

  - [ ] `GroupStatisticsDTO`

- **GroupServiceController.java** (`java/im/turms/service/domain/group/access/servicerequest/controller/GroupServiceController.java`)
> [简述功能]

  - [ ] `handleCreateGroupRequest`
  - [ ] `handleDeleteGroupRequest`
  - [ ] `handleQueryGroupsRequest`
  - [ ] `handleQueryJoinedGroupIdsRequest`
  - [ ] `handleQueryJoinedGroupsRequest`
  - [ ] `handleUpdateGroupRequest`
  - [ ] `handleCreateGroupBlockedUserRequest`
  - [ ] `handleDeleteGroupBlockedUserRequest`
  - [ ] `handleQueryGroupBlockedUserIdsRequest`
  - [ ] `handleQueryGroupBlockedUsersInfosRequest`
  - [ ] `handleCheckGroupQuestionAnswerRequest`
  - [ ] `handleCreateGroupInvitationRequestRequest`
  - [ ] `handleCreateGroupJoinRequestRequest`
  - [ ] `handleCreateGroupQuestionsRequest`
  - [ ] `handleDeleteGroupInvitationRequest`
  - [ ] `handleUpdateGroupInvitationRequest`
  - [ ] `handleDeleteGroupJoinRequestRequest`
  - [ ] `handleUpdateGroupJoinRequestRequest`
  - [ ] `handleDeleteGroupJoinQuestionsRequest`
  - [ ] `handleQueryGroupInvitationsRequest`
  - [ ] `handleQueryGroupJoinRequestsRequest`
  - [ ] `handleQueryGroupJoinQuestionsRequest`
  - [ ] `handleUpdateGroupJoinQuestionRequest`
  - [ ] `handleCreateGroupMembersRequest`
  - [ ] `handleDeleteGroupMembersRequest`
  - [ ] `handleQueryGroupMembersRequest`
  - [ ] `handleUpdateGroupMemberRequest`

- **CheckGroupQuestionAnswerResult.java** (`java/im/turms/service/domain/group/bo/CheckGroupQuestionAnswerResult.java`)
> [简述功能]

  - [ ] `CheckGroupQuestionAnswerResult`

- **GroupInvitationStrategy.java** (`java/im/turms/service/domain/group/bo/GroupInvitationStrategy.java`)
> [简述功能]

  - [x] `requiresApproval` -> `internal/domain/group/constant/group_strategy.go`

- **HandleHandleGroupInvitationResult.java** (`java/im/turms/service/domain/group/bo/HandleHandleGroupInvitationResult.java`)
> [简述功能]

  - [ ] `HandleHandleGroupInvitationResult`

- **HandleHandleGroupJoinRequestResult.java** (`java/im/turms/service/domain/group/bo/HandleHandleGroupJoinRequestResult.java`)
> [简述功能]

  - [ ] `HandleHandleGroupJoinRequestResult`

- **NewGroupQuestion.java** (`java/im/turms/service/domain/group/bo/NewGroupQuestion.java`)
> [简述功能]

  - [ ] `NewGroupQuestion`

- **GroupBlocklistRepository.java** (`java/im/turms/service/domain/group/repository/GroupBlocklistRepository.java`)
> [简述功能]

  - [ ] `updateBlockedUsers`
  - [x] `count` -> `internal/domain/user/repository/user_repository.go`
  - [ ] `findBlockedUserIds`
  - [ ] `findBlockedUsers`

- **GroupInvitationRepository.java** (`java/im/turms/service/domain/group/repository/GroupInvitationRepository.java`)
> [简述功能]

  - [x] `getEntityExpireAfterSeconds` -> `internal/domain/user/repository/user_friend_request_repository.go`
  - [x] `updateStatusIfPending` -> `internal/domain/user/repository/user_friend_request_repository.go`
  - [ ] `updateInvitations`
  - [x] `count` -> `internal/domain/user/repository/user_repository.go`
  - [ ] `findGroupIdAndInviteeIdAndStatus`
  - [ ] `findGroupIdAndInviterIdAndInviteeIdAndStatus`
  - [x] `findInvitationsByInviteeId` -> `internal/domain/group/repository/group_invitation_repository.go`
  - [ ] `findInvitationsByInviterId`
  - [x] `findInvitationsByGroupId` -> `internal/domain/group/repository/group_invitation_repository.go`
  - [ ] `findInviteeIdAndGroupIdAndCreationDateAndStatus`
  - [ ] `findInvitations`

- **GroupJoinRequestRepository.java** (`java/im/turms/service/domain/group/repository/GroupJoinRequestRepository.java`)
> [简述功能]

  - [x] `getEntityExpireAfterSeconds` -> `internal/domain/user/repository/user_friend_request_repository.go`
  - [x] `updateStatusIfPending` -> `internal/domain/user/repository/user_friend_request_repository.go`
  - [ ] `updateRequests`
  - [ ] `countRequests`
  - [ ] `findGroupId`
  - [ ] `findRequesterIdAndStatusAndGroupId`
  - [x] `findRequestsByGroupId` -> `internal/domain/group/repository/group_join_request_repository.go`
  - [x] `findRequestsByRequesterId` -> `internal/domain/group/repository/group_join_request_repository.go`
  - [ ] `findRequests`

- **GroupMemberRepository.java** (`java/im/turms/service/domain/group/repository/GroupMemberRepository.java`)
> [简述功能]

  - [ ] `deleteAllGroupMembers`
  - [ ] `updateGroupMembers`
  - [x] `countMembers` -> `internal/domain/user/repository/user_relationship_group_member_repository.go`
  - [x] `countMembers` -> `internal/domain/user/repository/user_relationship_group_member_repository.go`
  - [ ] `findGroupManagersAndOwnerId`
  - [ ] `findGroupMembers`
  - [ ] `findGroupMembers`
  - [ ] `findGroupsMembers`
  - [x] `findGroupMemberIds` -> `internal/domain/group/repository/group_member_repository.go`
  - [x] `findGroupMemberIds` -> `internal/domain/group/repository/group_member_repository.go`
  - [ ] `findGroupMemberKeyAndRoleParis`
  - [x] `findGroupMemberRole` -> `internal/domain/group/repository/group_member_repository.go`
  - [ ] `findMemberIdsByGroupId`
  - [x] `findUserJoinedGroupIds` -> `internal/domain/group/repository/group_member_repository.go`
  - [ ] `findUsersJoinedGroupIds`
  - [x] `isMemberMuted` -> `internal/domain/group/service/group_member_service.go`

- **GroupQuestionRepository.java** (`java/im/turms/service/domain/group/repository/GroupQuestionRepository.java`)
> [简述功能]

  - [ ] `updateQuestion`
  - [ ] `updateQuestions`
  - [ ] `countQuestions`
  - [ ] `checkQuestionAnswerAndGetScore`
  - [ ] `findGroupId`
  - [ ] `findQuestions`

- **GroupRepository.java** (`java/im/turms/service/domain/group/repository/GroupRepository.java`)
> [简述功能]

  - [ ] `updateGroupsDeletionDate`
  - [ ] `updateGroups`
  - [ ] `countCreatedGroups`
  - [ ] `countDeletedGroups`
  - [x] `countGroups` -> `internal/domain/user/repository/user_relationship_group_member_repository.go`
  - [x] `countOwnedGroups` -> `internal/domain/group/repository/group_repository.go`
  - [x] `countOwnedGroups` -> `internal/domain/group/repository/group_repository.go`
  - [x] `findGroups` -> `internal/domain/group/repository/group_repository.go`
  - [ ] `findNotDeletedGroups`
  - [x] `findAllNames` -> `internal/domain/user/repository/user_repository.go`
  - [ ] `findTypeId`
  - [ ] `findTypeIdAndGroupId`
  - [ ] `findTypeIdIfActiveAndNotDeleted`
  - [ ] `findMinimumScore`
  - [ ] `findOwnerId`
  - [ ] `isGroupMuted`
  - [ ] `isGroupActiveAndNotDeleted`

- **GroupTypeRepository.java** (`java/im/turms/service/domain/group/repository/GroupTypeRepository.java`)
> [简述功能]

  - [ ] `updateTypes`

- **GroupVersionRepository.java** (`java/im/turms/service/domain/group/repository/GroupVersionRepository.java`)
> [简述功能]

  - [ ] `updateVersions`
  - [ ] `updateVersions`
  - [x] `updateVersion` -> `internal/domain/group/repository/group_version_repository.go`
  - [x] `updateVersion` -> `internal/domain/group/repository/group_version_repository.go`
  - [ ] `findBlocklist`
  - [ ] `findInvitations`
  - [ ] `findJoinRequests`
  - [ ] `findJoinQuestions`
  - [ ] `findMembers`

- **GroupBlocklistService.java** (`java/im/turms/service/domain/group/service/GroupBlocklistService.java`)
> [简述功能]

  - [ ] `authAndBlockUser`
  - [x] `unblockUser` -> `internal/domain/group/service/group_blocklist_service.go`
  - [ ] `findBlockedUserIds`
  - [x] `isBlocked` -> `internal/domain/user/service/user_relationship_service.go`
  - [ ] `queryGroupBlockedUserIds`
  - [x] `queryBlockedUsers` -> `internal/domain/group/service/group_blocklist_service.go`
  - [ ] `countBlockedUsers`
  - [ ] `queryGroupBlockedUserIdsWithVersion`
  - [ ] `queryGroupBlockedUserInfosWithVersion`
  - [ ] `addBlockedUser`
  - [ ] `updateBlockedUsers`
  - [ ] `deleteBlockedUsers`

- **GroupInvitationService.java** (`java/im/turms/service/domain/group/service/GroupInvitationService.java`)
> [简述功能]

  - [ ] `authAndCreateGroupInvitation`
  - [ ] `createGroupInvitation`
  - [ ] `queryGroupIdAndInviterIdAndInviteeIdAndStatus`
  - [ ] `queryGroupIdAndInviteeIdAndStatus`
  - [ ] `authAndRecallPendingGroupInvitation`
  - [ ] `queryGroupInvitationsByInviteeId`
  - [ ] `queryGroupInvitationsByInviterId`
  - [ ] `queryGroupInvitationsByGroupId`
  - [ ] `queryUserGroupInvitationsWithVersion`
  - [ ] `authAndQueryGroupInvitationsWithVersion`
  - [ ] `queryInviteeIdAndGroupIdAndCreationDateAndStatusByInvitationId`
  - [ ] `queryInvitations`
  - [ ] `countInvitations`
  - [ ] `deleteInvitations`
  - [ ] `authAndHandleInvitation`
  - [ ] `updatePendingInvitationStatus`
  - [ ] `updateInvitations`

- **GroupJoinRequestService.java** (`java/im/turms/service/domain/group/service/GroupJoinRequestService.java`)
> [简述功能]

  - [ ] `authAndCreateGroupJoinRequest`
  - [ ] `authAndRecallPendingGroupJoinRequest`
  - [ ] `authAndQueryGroupJoinRequestsWithVersion`
  - [ ] `queryGroupJoinRequestsByGroupId`
  - [ ] `queryGroupJoinRequestsByRequesterId`
  - [ ] `queryGroupId`
  - [ ] `queryJoinRequests`
  - [ ] `countJoinRequests`
  - [ ] `deleteJoinRequests`
  - [ ] `authAndHandleJoinRequest`
  - [ ] `updatePendingJoinRequestStatus`
  - [ ] `updateJoinRequests`
  - [ ] `createGroupJoinRequest`

- **GroupMemberService.java** (`java/im/turms/service/domain/group/service/GroupMemberService.java`)
> [简述功能]

  - [x] `addGroupMember` -> `internal/domain/group/service/group_member_service.go`
  - [ ] `addGroupMembers`
  - [ ] `authAndAddGroupMembers`
  - [ ] `authAndDeleteGroupMembers`
  - [ ] `deleteGroupMember`
  - [ ] `deleteGroupMembers`
  - [ ] `updateGroupMember`
  - [ ] `updateGroupMembers`
  - [ ] `updateGroupMembers`
  - [x] `isGroupMember` -> `internal/domain/group/service/group_member_service.go`
  - [x] `isGroupMember` -> `internal/domain/group/service/group_member_service.go`
  - [ ] `findExistentMemberGroupIds`
  - [ ] `isAllowedToInviteUser`
  - [ ] `isAllowedToBeInvited`
  - [ ] `isAllowedToSendMessage`
  - [x] `isMemberMuted` -> `internal/domain/group/service/group_member_service.go`
  - [ ] `queryGroupMemberKeyAndRolePairs`
  - [ ] `queryGroupMemberRole`
  - [ ] `isOwner`
  - [ ] `isOwnerOrManager`
  - [ ] `isOwnerOrManagerOrMember`
  - [ ] `queryUserJoinedGroupIds`
  - [ ] `queryUsersJoinedGroupIds`
  - [ ] `queryMemberIdsInUsersJoinedGroups`
  - [ ] `queryGroupMemberIds`
  - [ ] `queryGroupMemberIds`
  - [ ] `queryGroupMembers`
  - [x] `countMembers` -> `internal/domain/user/repository/user_relationship_group_member_repository.go`
  - [ ] `deleteGroupMembers`
  - [ ] `queryGroupMembers`
  - [ ] `queryGroupMembers`
  - [ ] `authAndQueryGroupMembers`
  - [ ] `authAndQueryGroupMembersWithVersion`
  - [ ] `authAndUpdateGroupMember`
  - [ ] `deleteAllGroupMembers`
  - [ ] `queryGroupManagersAndOwnerId`

- **GroupQuestionService.java** (`java/im/turms/service/domain/group/service/GroupQuestionService.java`)
> [简述功能]

  - [ ] `checkGroupQuestionAnswerAndGetScore`
  - [ ] `authAndCheckGroupQuestionAnswerAndJoin`
  - [ ] `authAndCreateGroupJoinQuestions`
  - [ ] `createGroupJoinQuestions`
  - [ ] `queryGroupId`
  - [ ] `authAndDeleteGroupJoinQuestions`
  - [ ] `queryGroupJoinQuestions`
  - [ ] `countGroupJoinQuestions`
  - [ ] `deleteGroupJoinQuestions`
  - [ ] `authAndQueryGroupJoinQuestionsWithVersion`
  - [ ] `authAndUpdateGroupJoinQuestion`
  - [ ] `updateGroupJoinQuestions`

- **GroupService.java** (`java/im/turms/service/domain/group/service/GroupService.java`)
> [简述功能]

  - [x] `createGroup` -> `internal/domain/group/service/group_service.go`
  - [ ] `authAndDeleteGroup`
  - [ ] `authAndCreateGroup`
  - [ ] `deleteGroupsAndGroupMembers`
  - [ ] `queryGroups`
  - [ ] `queryGroupTypeIfActiveAndNotDeleted`
  - [ ] `queryGroupTypeIfActiveAndNotDeleted`
  - [ ] `queryGroupTypeId`
  - [ ] `queryGroupTypeIdIfActiveAndNotDeleted`
  - [ ] `queryGroupMinimumScore`
  - [ ] `authAndTransferGroupOwnership`
  - [ ] `queryGroupOwnerId`
  - [ ] `checkAndTransferGroupOwnership`
  - [ ] `checkAndTransferGroupOwnership`
  - [ ] `updateGroupInformation`
  - [ ] `updateGroupsInformation`
  - [ ] `authAndUpdateGroupInformation`
  - [ ] `authAndQueryGroups`
  - [ ] `queryJoinedGroups`
  - [ ] `queryJoinedGroupIdsWithVersion`
  - [ ] `queryJoinedGroupsWithVersion`
  - [ ] `isAllowedToCreateGroupAndHaveGroupType`
  - [ ] `isAllowedToCreateGroup`
  - [ ] `isAllowedCreateGroupWithGroupType`
  - [ ] `isAllowedUpdateGroupToGroupType`
  - [x] `countOwnedGroups` -> `internal/domain/group/repository/group_repository.go`
  - [x] `countOwnedGroups` -> `internal/domain/group/repository/group_repository.go`
  - [ ] `countCreatedGroups`
  - [x] `countGroups` -> `internal/domain/user/repository/user_relationship_group_member_repository.go`
  - [ ] `countDeletedGroups`
  - [x] `count` -> `internal/domain/user/repository/user_repository.go`
  - [ ] `isGroupMuted`
  - [ ] `isGroupActiveAndNotDeleted`

- **GroupTypeService.java** (`java/im/turms/service/domain/group/service/GroupTypeService.java`)
> [简述功能]

  - [ ] `initGroupTypes`
  - [ ] `queryGroupTypes`
  - [ ] `addGroupType`
  - [ ] `updateGroupTypes`
  - [ ] `deleteGroupTypes`
  - [ ] `queryGroupType`
  - [ ] `queryGroupTypes`
  - [ ] `groupTypeExists`
  - [ ] `countGroupTypes`

- **GroupVersionService.java** (`java/im/turms/service/domain/group/service/GroupVersionService.java`)
> [简述功能]

  - [ ] `queryMembersVersion`
  - [ ] `queryBlocklistVersion`
  - [x] `queryGroupJoinRequestsVersion` -> `internal/domain/user/service/user_version_service.go`
  - [ ] `queryGroupJoinQuestionsVersion`
  - [ ] `queryGroupInvitationsVersion`
  - [x] `updateVersion` -> `internal/domain/group/repository/group_version_repository.go`
  - [x] `updateMembersVersion` -> `internal/domain/group/service/group_version_service.go`
  - [x] `updateMembersVersion` -> `internal/domain/group/service/group_version_service.go`
  - [x] `updateMembersVersion` -> `internal/domain/group/service/group_version_service.go`
  - [x] `updateBlocklistVersion` -> `internal/domain/group/service/group_version_service.go`
  - [x] `updateJoinRequestsVersion` -> `internal/domain/group/repository/group_version_repository.go`
  - [x] `updateJoinQuestionsVersion` -> `internal/domain/group/repository/group_version_repository.go`
  - [ ] `updateGroupInvitationsVersion`
  - [x] `updateSpecificVersion` -> `internal/domain/user/service/user_version_service.go`
  - [x] `updateSpecificVersion` -> `internal/domain/user/service/user_version_service.go`
  - [x] `updateSpecificVersion` -> `internal/domain/user/service/user_version_service.go`
  - [x] `upsert` -> `internal/domain/user/repository/user_relationship_repository.go`
  - [x] `delete` -> `internal/domain/user/service/user_version_service.go`

- **MessageController.java** (`java/im/turms/service/domain/message/access/admin/controller/MessageController.java`)
> [简述功能]

  - [ ] `createMessages`
  - [x] `queryMessages` -> `internal/domain/message/service/message_service.go`
  - [x] `queryMessages` -> `internal/domain/message/service/message_service.go`
  - [ ] `countMessages`
  - [ ] `updateMessages`
  - [ ] `deleteMessages`

- **CreateMessageDTO.java** (`java/im/turms/service/domain/message/access/admin/dto/request/CreateMessageDTO.java`)
> [简述功能]

  - [ ] `CreateMessageDTO`

- **UpdateMessageDTO.java** (`java/im/turms/service/domain/message/access/admin/dto/request/UpdateMessageDTO.java`)
> [简述功能]

  - [ ] `UpdateMessageDTO`

- **MessageStatisticsDTO.java** (`java/im/turms/service/domain/message/access/admin/dto/response/MessageStatisticsDTO.java`)
> [简述功能]

  - [ ] `MessageStatisticsDTO`

- **MessageServiceController.java** (`java/im/turms/service/domain/message/access/servicerequest/controller/MessageServiceController.java`)
> [简述功能]

  - [x] `handleCreateMessageRequest` -> `internal/domain/message/controller/message_controller.go`
  - [ ] `handleQueryMessagesRequest`
  - [ ] `handleUpdateMessageRequest`
  - [ ] `handleCreateMessageReactionsRequest`
  - [ ] `handleDeleteMessageReactionsRequest`

- **MessageAndRecipientIds.java** (`java/im/turms/service/domain/message/bo/MessageAndRecipientIds.java`)
> [简述功能]

  - [ ] `MessageAndRecipientIds`

- **Message.java** (`java/im/turms/service/domain/message/po/Message.java`)
> [简述功能]

  - [ ] `groupId`

- **MessageRepository.java** (`java/im/turms/service/domain/message/repository/MessageRepository.java`)
> [简述功能]

  - [ ] `updateMessages`
  - [ ] `updateMessagesDeletionDate`
  - [ ] `existsBySenderIdAndTargetId`
  - [ ] `countMessages`
  - [ ] `countUsersWhoSentMessage`
  - [ ] `countGroupsThatSentMessages`
  - [ ] `countSentMessages`
  - [ ] `findDeliveryDate`
  - [ ] `findExpiredMessageIds`
  - [ ] `findMessageGroupId`
  - [ ] `findMessageSenderIdAndTargetIdAndIsGroupMessage`
  - [ ] `findMessages`
  - [ ] `findIsGroupMessageAndTargetId`
  - [ ] `findIsGroupMessageAndTargetIdAndDeliveryDate`
  - [ ] `getGroupConversationId`
  - [ ] `getPrivateConversationId`

- **MessageService.java** (`java/im/turms/service/domain/message/service/MessageService.java`)
> [简述功能]

  - [ ] `isMessageRecipientOrSender`
  - [ ] `authAndQueryCompleteMessages`
  - [ ] `queryMessage`
  - [x] `queryMessages` -> `internal/domain/message/service/message_service.go`
  - [ ] `saveMessage`
  - [ ] `queryExpiredMessageIds`
  - [ ] `deleteExpiredMessages`
  - [ ] `deleteMessages`
  - [ ] `updateMessages`
  - [ ] `hasPrivateMessage`
  - [ ] `countMessages`
  - [ ] `countUsersWhoSentMessage`
  - [ ] `countGroupsThatSentMessages`
  - [ ] `countSentMessages`
  - [ ] `countSentMessagesOnAverage`
  - [ ] `authAndUpdateMessage`
  - [ ] `queryMessageRecipients`
  - [x] `authAndSaveMessage` -> `internal/domain/message/service/message_service.go`
  - [ ] `saveMessage`
  - [ ] `authAndCloneAndSaveMessage`
  - [ ] `cloneAndSaveMessage`
  - [x] `authAndSaveAndSendMessage` -> `internal/domain/message/service/message_service.go`
  - [ ] `saveAndSendMessage`
  - [ ] `saveAndSendMessage`
  - [ ] `deleteGroupMessageSequenceIds`
  - [ ] `deletePrivateMessageSequenceIds`
  - [ ] `fetchGroupMessageSequenceId`
  - [ ] `fetchPrivateMessageSequenceId`

- **StatisticsService.java** (`java/im/turms/service/domain/observation/service/StatisticsService.java`)
> [简述功能]

  - [ ] `countOnlineUsersByNodes`
  - [x] `countOnlineUsers` -> `internal/domain/gateway/session/sharded_map.go`

- **StorageServiceController.java** (`java/im/turms/service/domain/storage/access/servicerequest/controller/StorageServiceController.java`)
> [简述功能]

  - [ ] `handleDeleteResourceRequest`
  - [ ] `handleQueryResourceUploadInfoRequest`
  - [ ] `handleQueryResourceDownloadInfoRequest`
  - [ ] `handleUpdateMessageAttachmentInfoRequest`
  - [ ] `handleQueryMessageAttachmentInfosRequest`

- **StorageResourceInfo.java** (`java/im/turms/service/domain/storage/bo/StorageResourceInfo.java`)
> [简述功能]

  - [ ] `StorageResourceInfo`

- **StorageService.java** (`java/im/turms/service/domain/storage/service/StorageService.java`)
> [简述功能]

  - [x] `deleteResource` -> `internal/domain/storage/service/storage_service.go`
  - [x] `queryResourceUploadInfo` -> `internal/domain/storage/service/storage_service.go`
  - [x] `queryResourceDownloadInfo` -> `internal/domain/storage/service/storage_service.go`
  - [ ] `shareMessageAttachmentWithUser`
  - [ ] `shareMessageAttachmentWithGroup`
  - [ ] `unshareMessageAttachmentWithUser`
  - [ ] `unshareMessageAttachmentWithGroup`
  - [ ] `queryMessageAttachmentInfosUploadedByRequester`
  - [ ] `queryMessageAttachmentInfosInPrivateConversations`
  - [ ] `queryMessageAttachmentInfosInGroupConversations`

- **UserController.java** (`java/im/turms/service/domain/user/access/admin/controller/UserController.java`)
> [简述功能]

  - [x] `addUser` -> `internal/domain/user/service/user_service.go`
  - [x] `queryUsers` -> `internal/domain/user/service/user_service.go`
  - [x] `queryUsers` -> `internal/domain/user/service/user_service.go`
  - [x] `countUsers` -> `internal/domain/user/service/user_service.go`
  - [x] `updateUser` -> `internal/domain/user/service/user_service.go`
  - [x] `deleteUsers` -> `internal/domain/user/service/user_service.go`

- **UserOnlineInfoController.java** (`java/im/turms/service/domain/user/access/admin/controller/UserOnlineInfoController.java`)
> [简述功能]

  - [x] `countOnlineUsers` -> `internal/domain/gateway/session/sharded_map.go`
  - [x] `queryUserSessions` -> `internal/domain/user/service/onlineuser/session_service.go`
  - [ ] `queryUserStatuses`
  - [ ] `queryUsersNearby`
  - [ ] `queryUserLocations`
  - [ ] `updateUserOnlineStatus`

- **UserRoleController.java** (`java/im/turms/service/domain/user/access/admin/controller/UserRoleController.java`)
> [简述功能]

  - [x] `addUserRole` -> `internal/domain/user/service/user_role_service.go`
  - [x] `queryUserRoles` -> `internal/domain/user/service/user_role_service.go`
  - [ ] `queryUserRoleGroups`
  - [ ] `updateUserRole`
  - [ ] `deleteUserRole`

- **UserFriendRequestController.java** (`java/im/turms/service/domain/user/access/admin/controller/relationship/UserFriendRequestController.java`)
> [简述功能]

  - [x] `createFriendRequest` -> `internal/domain/user/service/user_friend_request_service.go`
  - [x] `queryFriendRequests` -> `internal/domain/user/service/user_friend_request_service.go`
  - [x] `queryFriendRequests` -> `internal/domain/user/service/user_friend_request_service.go`
  - [x] `updateFriendRequests` -> `internal/domain/user/service/user_friend_request_service.go`
  - [x] `deleteFriendRequests` -> `internal/domain/user/service/user_friend_request_service.go`

- **UserRelationshipController.java** (`java/im/turms/service/domain/user/access/admin/controller/relationship/UserRelationshipController.java`)
> [简述功能]

  - [ ] `addRelationship`
  - [x] `queryRelationships` -> `internal/domain/user/service/user_relationship_service.go`
  - [x] `queryRelationships` -> `internal/domain/user/service/user_relationship_service.go`
  - [ ] `updateRelationships`
  - [ ] `deleteRelationships`

- **UserRelationshipGroupController.java** (`java/im/turms/service/domain/user/access/admin/controller/relationship/UserRelationshipGroupController.java`)
> [简述功能]

  - [ ] `addRelationshipGroup`
  - [x] `deleteRelationshipGroups` -> `internal/domain/user/service/user_relationship_group_service.go`
  - [x] `updateRelationshipGroups` -> `internal/domain/user/service/user_relationship_group_service.go`
  - [x] `queryRelationshipGroups` -> `internal/domain/user/service/user_relationship_group_service.go`
  - [x] `queryRelationshipGroups` -> `internal/domain/user/service/user_relationship_group_service.go`

- **AddFriendRequestDTO.java** (`java/im/turms/service/domain/user/access/admin/dto/request/AddFriendRequestDTO.java`)
> [简述功能]

  - [ ] `AddFriendRequestDTO`

- **AddRelationshipDTO.java** (`java/im/turms/service/domain/user/access/admin/dto/request/AddRelationshipDTO.java`)
> [简述功能]

  - [ ] `AddRelationshipDTO`

- **AddRelationshipGroupDTO.java** (`java/im/turms/service/domain/user/access/admin/dto/request/AddRelationshipGroupDTO.java`)
> [简述功能]

  - [ ] `AddRelationshipGroupDTO`

- **AddUserDTO.java** (`java/im/turms/service/domain/user/access/admin/dto/request/AddUserDTO.java`)
> [简述功能]

  - [ ] `AddUserDTO`
  - [ ] `toString`

- **AddUserRoleDTO.java** (`java/im/turms/service/domain/user/access/admin/dto/request/AddUserRoleDTO.java`)
> [简述功能]

  - [ ] `AddUserRoleDTO`

- **UpdateFriendRequestDTO.java** (`java/im/turms/service/domain/user/access/admin/dto/request/UpdateFriendRequestDTO.java`)
> [简述功能]

  - [ ] `UpdateFriendRequestDTO`

- **UpdateOnlineStatusDTO.java** (`java/im/turms/service/domain/user/access/admin/dto/request/UpdateOnlineStatusDTO.java`)
> [简述功能]

  - [ ] `UpdateOnlineStatusDTO`

- **UpdateRelationshipDTO.java** (`java/im/turms/service/domain/user/access/admin/dto/request/UpdateRelationshipDTO.java`)
> [简述功能]

  - [ ] `UpdateRelationshipDTO`

- **UpdateRelationshipGroupDTO.java** (`java/im/turms/service/domain/user/access/admin/dto/request/UpdateRelationshipGroupDTO.java`)
> [简述功能]

  - [ ] `UpdateRelationshipGroupDTO`

- **UpdateUserDTO.java** (`java/im/turms/service/domain/user/access/admin/dto/request/UpdateUserDTO.java`)
> [简述功能]

  - [ ] `UpdateUserDTO`
  - [ ] `toString`

- **UpdateUserRoleDTO.java** (`java/im/turms/service/domain/user/access/admin/dto/request/UpdateUserRoleDTO.java`)
> [简述功能]

  - [ ] `UpdateUserRoleDTO`

- **OnlineUserCountDTO.java** (`java/im/turms/service/domain/user/access/admin/dto/response/OnlineUserCountDTO.java`)
> [简述功能]

  - [ ] `OnlineUserCountDTO`

- **UserFriendRequestDTO.java** (`java/im/turms/service/domain/user/access/admin/dto/response/UserFriendRequestDTO.java`)
> [简述功能]

  - [ ] `UserFriendRequestDTO`

- **UserLocationDTO.java** (`java/im/turms/service/domain/user/access/admin/dto/response/UserLocationDTO.java`)
> [简述功能]

  - [ ] `UserLocationDTO`

- **UserRelationshipDTO.java** (`java/im/turms/service/domain/user/access/admin/dto/response/UserRelationshipDTO.java`)
> [简述功能]

  - [ ] `UserRelationshipDTO`
  - [ ] `fromDomain`
  - [ ] `fromDomain`
  - [ ] `Key`

- **UserStatisticsDTO.java** (`java/im/turms/service/domain/user/access/admin/dto/response/UserStatisticsDTO.java`)
> [简述功能]

  - [ ] `UserStatisticsDTO`

- **UserStatusDTO.java** (`java/im/turms/service/domain/user/access/admin/dto/response/UserStatusDTO.java`)
> [简述功能]

  - [ ] `UserStatusDTO`

- **UserRelationshipServiceController.java** (`java/im/turms/service/domain/user/access/servicerequest/controller/UserRelationshipServiceController.java`)
> [简述功能]

  - [ ] `handleCreateFriendRequestRequest`
  - [ ] `handleCreateRelationshipGroupRequest`
  - [ ] `handleCreateRelationshipRequest`
  - [ ] `handleDeleteFriendRequestRequest`
  - [ ] `handleDeleteRelationshipGroupRequest`
  - [ ] `handleDeleteRelationshipRequest`
  - [ ] `handleQueryFriendRequestsRequest`
  - [ ] `handleQueryRelatedUserIdsRequest`
  - [ ] `handleQueryRelationshipGroupsRequest`
  - [ ] `handleQueryRelationshipsRequest`
  - [ ] `handleUpdateFriendRequestRequest`
  - [ ] `handleUpdateRelationshipGroupRequest`
  - [ ] `handleUpdateRelationshipRequest`

- **UserServiceController.java** (`java/im/turms/service/domain/user/access/servicerequest/controller/UserServiceController.java`)
> [简述功能]

  - [ ] `handleQueryUserProfilesRequest`
  - [ ] `handleQueryNearbyUsersRequest`
  - [ ] `handleQueryUserOnlineStatusesRequest`
  - [ ] `handleUpdateUserLocationRequest`
  - [ ] `handleUpdateUserOnlineStatusRequest`
  - [ ] `handleUpdateUserRequest`

- **UserSettingsServiceController.java** (`java/im/turms/service/domain/user/access/servicerequest/controller/UserSettingsServiceController.java`)
> [简述功能]

  - [ ] `handleDeleteUserSettingsRequest`
  - [ ] `handleUpdateUserSettingsRequest`
  - [ ] `handleQueryUserSettingsRequest`

- **HandleFriendRequestResult.java** (`java/im/turms/service/domain/user/bo/HandleFriendRequestResult.java`)
> [简述功能]

  - [ ] `HandleFriendRequestResult`

- **UpsertRelationshipResult.java** (`java/im/turms/service/domain/user/bo/UpsertRelationshipResult.java`)
> [简述功能]

  - [ ] `UpsertRelationshipResult`

- **UserFriendRequestRepository.java** (`java/im/turms/service/domain/user/repository/UserFriendRequestRepository.java`)
> [简述功能]

  - [x] `getEntityExpireAfterSeconds` -> `internal/domain/user/repository/user_friend_request_repository.go`
  - [x] `updateFriendRequests` -> `internal/domain/user/service/user_friend_request_service.go`
  - [x] `updateStatusIfPending` -> `internal/domain/user/repository/user_friend_request_repository.go`
  - [x] `countFriendRequests` -> `internal/domain/user/service/user_friend_request_service.go`
  - [x] `findFriendRequests` -> `internal/domain/user/repository/user_friend_request_repository.go`
  - [x] `findFriendRequestsByRecipientId` -> `internal/domain/user/repository/user_friend_request_repository.go`
  - [x] `findFriendRequestsByRequesterId` -> `internal/domain/user/repository/user_friend_request_repository.go`
  - [x] `findRecipientId` -> `internal/domain/user/repository/user_friend_request_repository.go`
  - [x] `findRequesterIdAndRecipientIdAndStatus` -> `internal/domain/user/repository/user_friend_request_repository.go`
  - [x] `findRequesterIdAndRecipientIdAndCreationDateAndStatus` -> `internal/domain/user/repository/user_friend_request_repository.go`
  - [x] `hasPendingFriendRequest` -> `internal/domain/user/service/user_friend_request_service.go`
  - [x] `hasPendingOrDeclinedOrIgnoredOrExpiredRequest` -> `internal/domain/user/repository/user_friend_request_repository.go`

- **UserRelationshipGroupMemberRepository.java** (`java/im/turms/service/domain/user/repository/UserRelationshipGroupMemberRepository.java`)
> [简述功能]

  - [x] `deleteAllRelatedUserFromRelationshipGroup` -> `internal/domain/user/repository/user_relationship_group_member_repository.go`
  - [x] `deleteRelatedUserFromRelationshipGroup` -> `internal/domain/user/service/user_relationship_group_service.go`
  - [x] `deleteRelatedUsersFromAllRelationshipGroups` -> `internal/domain/user/service/user_relationship_group_service.go`
  - [x] `countGroups` -> `internal/domain/user/repository/user_relationship_group_member_repository.go`
  - [x] `countMembers` -> `internal/domain/user/repository/user_relationship_group_member_repository.go`
  - [x] `findGroupIndexes` -> `internal/domain/user/repository/user_relationship_group_member_repository.go`
  - [x] `findRelationshipGroupMemberIds` -> `internal/domain/user/repository/user_relationship_group_member_repository.go`
  - [x] `findRelationshipGroupMemberIds` -> `internal/domain/user/repository/user_relationship_group_member_repository.go`
  - [x] `findRelationshipGroupMembers` -> `internal/domain/user/repository/user_relationship_group_member_repository.go`

- **UserRelationshipGroupRepository.java** (`java/im/turms/service/domain/user/repository/UserRelationshipGroupRepository.java`)
> [简述功能]

  - [x] `deleteAllRelationshipGroups` -> `internal/domain/user/service/user_relationship_group_service.go`
  - [x] `updateRelationshipGroupName` -> `internal/domain/user/service/user_relationship_group_service.go`
  - [x] `updateRelationshipGroups` -> `internal/domain/user/service/user_relationship_group_service.go`
  - [x] `countRelationshipGroups` -> `internal/domain/user/service/user_relationship_group_service.go`
  - [x] `findRelationshipGroups` -> `internal/domain/user/repository/user_relationship_group_repository.go`
  - [x] `findRelationshipGroupsInfos` -> `internal/domain/user/repository/user_relationship_group_repository.go`

- **UserRelationshipRepository.java** (`java/im/turms/service/domain/user/repository/UserRelationshipRepository.java`)
> [简述功能]

  - [x] `deleteAllRelationships` -> `internal/domain/user/service/user_relationship_service.go`
  - [x] `updateUserOneSidedRelationships` -> `internal/domain/user/service/user_relationship_service.go`
  - [x] `countRelationships` -> `internal/domain/user/service/user_relationship_service.go`
  - [x] `findRelatedUserIds` -> `internal/domain/user/repository/user_relationship_repository.go`
  - [x] `findRelationships` -> `internal/domain/user/repository/user_relationship_repository.go`
  - [x] `findRelationships` -> `internal/domain/user/repository/user_relationship_repository.go`
  - [x] `hasRelationshipAndNotBlocked` -> `internal/domain/user/service/user_relationship_service.go`
  - [x] `isBlocked` -> `internal/domain/user/service/user_relationship_service.go`

- **UserRepository.java** (`java/im/turms/service/domain/user/repository/UserRepository.java`)
> [简述功能]

  - [x] `updateUsers` -> `internal/domain/user/service/user_service.go`
  - [x] `updateUsersDeletionDate` -> `internal/domain/user/repository/user_repository.go`
  - [x] `checkIfUserExists` -> `internal/domain/user/service/user_service.go`
  - [x] `countRegisteredUsers` -> `internal/domain/user/service/user_service.go`
  - [x] `countDeletedUsers` -> `internal/domain/user/service/user_service.go`
  - [x] `countUsers` -> `internal/domain/user/service/user_service.go`
  - [x] `countUsers` -> `internal/domain/user/service/user_service.go`
  - [x] `findName` -> `internal/domain/user/repository/user_repository.go`
  - [x] `findAllNames` -> `internal/domain/user/repository/user_repository.go`
  - [x] `findProfileAccessIfNotDeleted` -> `internal/domain/user/repository/user_repository.go`
  - [x] `findUsers` -> `internal/domain/user/repository/user_repository.go`
  - [x] `findNotDeletedUserProfiles` -> `internal/domain/user/repository/user_repository.go`
  - [x] `findUsersProfile` -> `internal/domain/user/repository/user_repository.go`
  - [x] `findUserRoleId` -> `internal/domain/user/repository/user_repository.go`
  - [x] `isActiveAndNotDeleted` -> `internal/domain/user/repository/user_repository.go`

- **UserRoleRepository.java** (`java/im/turms/service/domain/user/repository/UserRoleRepository.java`)
> [简述功能]

  - [x] `updateUserRoles` -> `internal/domain/user/service/user_role_service.go`

- **UserSettingsRepository.java** (`java/im/turms/service/domain/user/repository/UserSettingsRepository.java`)
> [简述功能]

  - [x] `upsertSettings` -> `internal/domain/user/service/user_settings_service.go`
  - [x] `unsetSettings` -> `internal/domain/user/service/user_settings_service.go`
  - [x] `findByIdAndSettingNames` -> `internal/domain/user/repository/user_settings_repository.go`

- **UserVersionRepository.java** (`java/im/turms/service/domain/user/repository/UserVersionRepository.java`)
> [简述功能]

  - [x] `updateSpecificVersion` -> `internal/domain/user/service/user_version_service.go`
  - [x] `updateSpecificVersion` -> `internal/domain/user/service/user_version_service.go`
  - [x] `updateSpecificVersion` -> `internal/domain/user/service/user_version_service.go`
  - [ ] `findGroupJoinRequests`
  - [ ] `findJoinedGroup`
  - [ ] `findReceivedGroupInvitations`
  - [x] `findRelationships` -> `internal/domain/user/repository/user_relationship_repository.go`
  - [x] `findRelationshipGroups` -> `internal/domain/user/repository/user_relationship_group_repository.go`
  - [ ] `findSentGroupInvitations`
  - [ ] `findSentFriendRequests`
  - [ ] `findReceivedFriendRequests`

- **UserFriendRequestService.java** (`java/im/turms/service/domain/user/service/UserFriendRequestService.java`)
> [简述功能]

  - [x] `removeAllExpiredFriendRequests` -> `internal/domain/user/service/user_friend_request_service.go`
  - [x] `hasPendingFriendRequest` -> `internal/domain/user/service/user_friend_request_service.go`
  - [x] `createFriendRequest` -> `internal/domain/user/service/user_friend_request_service.go`
  - [x] `authAndCreateFriendRequest` -> `internal/domain/user/service/user_friend_request_service.go`
  - [x] `authAndRecallFriendRequest` -> `internal/domain/user/service/user_friend_request_service.go`
  - [x] `updatePendingFriendRequestStatus` -> `internal/domain/user/service/user_friend_request_service.go`
  - [x] `updateFriendRequests` -> `internal/domain/user/service/user_friend_request_service.go`
  - [x] `queryRecipientId` -> `internal/domain/user/service/user_friend_request_service.go`
  - [x] `queryRequesterIdAndRecipientIdAndStatus` -> `internal/domain/user/service/user_friend_request_service.go`
  - [x] `queryRequesterIdAndRecipientIdAndCreationDateAndStatus` -> `internal/domain/user/service/user_friend_request_service.go`
  - [x] `authAndHandleFriendRequest` -> `internal/domain/user/service/user_friend_request_service.go`
  - [x] `queryFriendRequestsWithVersion` -> `internal/domain/user/service/user_friend_request_service.go`
  - [x] `queryFriendRequestsByRecipientId` -> `internal/domain/user/service/user_friend_request_service.go`
  - [x] `queryFriendRequestsByRequesterId` -> `internal/domain/user/service/user_friend_request_service.go`
  - [x] `deleteFriendRequests` -> `internal/domain/user/service/user_friend_request_service.go`
  - [x] `queryFriendRequests` -> `internal/domain/user/service/user_friend_request_service.go`
  - [x] `countFriendRequests` -> `internal/domain/user/service/user_friend_request_service.go`

- **UserRelationshipGroupService.java** (`java/im/turms/service/domain/user/service/UserRelationshipGroupService.java`)
> [简述功能]

  - [x] `createRelationshipGroup` -> `internal/domain/user/service/user_relationship_group_service.go`
  - [x] `queryRelationshipGroupsInfos` -> `internal/domain/user/service/user_relationship_group_service.go`
  - [x] `queryRelationshipGroupsInfosWithVersion` -> `internal/domain/user/service/user_relationship_group_service.go`
  - [x] `queryGroupIndexes` -> `internal/domain/user/service/user_relationship_group_service.go`
  - [x] `queryRelationshipGroupMemberIds` -> `internal/domain/user/service/user_relationship_group_service.go`
  - [x] `queryRelationshipGroupMemberIds` -> `internal/domain/user/service/user_relationship_group_service.go`
  - [x] `updateRelationshipGroupName` -> `internal/domain/user/service/user_relationship_group_service.go`
  - [x] `upsertRelationshipGroupMember` -> `internal/domain/user/service/user_relationship_group_service.go`
  - [x] `updateRelationshipGroups` -> `internal/domain/user/service/user_relationship_group_service.go`
  - [x] `addRelatedUserToRelationshipGroups` -> `internal/domain/user/service/user_relationship_group_service.go`
  - [x] `deleteRelationshipGroupAndMoveMembersToNewGroup` -> `internal/domain/user/service/user_relationship_group_service.go`
  - [x] `deleteAllRelationshipGroups` -> `internal/domain/user/service/user_relationship_group_service.go`
  - [x] `deleteRelatedUserFromRelationshipGroup` -> `internal/domain/user/service/user_relationship_group_service.go`
  - [x] `deleteRelatedUserFromAllRelationshipGroups` -> `internal/domain/user/service/user_relationship_group_service.go`
  - [x] `deleteRelatedUsersFromAllRelationshipGroups` -> `internal/domain/user/service/user_relationship_group_service.go`
  - [x] `moveRelatedUserToNewGroup` -> `internal/domain/user/service/user_relationship_group_service.go`
  - [x] `deleteRelationshipGroups` -> `internal/domain/user/service/user_relationship_group_service.go`
  - [x] `deleteRelationshipGroups` -> `internal/domain/user/service/user_relationship_group_service.go`
  - [x] `queryRelationshipGroups` -> `internal/domain/user/service/user_relationship_group_service.go`
  - [x] `countRelationshipGroups` -> `internal/domain/user/service/user_relationship_group_service.go`
  - [x] `countRelationshipGroups` -> `internal/domain/user/service/user_relationship_group_service.go`
  - [x] `countRelationshipGroupMembers` -> `internal/domain/user/service/user_relationship_group_service.go`

- **UserRelationshipService.java** (`java/im/turms/service/domain/user/service/UserRelationshipService.java`)
> [简述功能]

  - [x] `deleteAllRelationships` -> `internal/domain/user/service/user_relationship_service.go`
  - [x] `deleteOneSidedRelationships` -> `internal/domain/user/service/user_relationship_service.go`
  - [x] `deleteOneSidedRelationship` -> `internal/domain/user/service/user_relationship_service.go`
  - [x] `tryDeleteTwoSidedRelationships` -> `internal/domain/user/service/user_relationship_service.go`
  - [x] `queryRelatedUserIdsWithVersion` -> `internal/domain/user/service/user_relationship_service.go`
  - [x] `queryRelationshipsWithVersion` -> `internal/domain/user/service/user_relationship_service.go`
  - [x] `queryRelatedUserIds` -> `internal/domain/user/service/user_relationship_service.go`
  - [x] `queryRelatedUserIds` -> `internal/domain/user/service/user_relationship_service.go`
  - [x] `queryRelationships` -> `internal/domain/user/service/user_relationship_service.go`
  - [x] `queryMembersRelationships` -> `internal/domain/user/service/user_relationship_service.go`
  - [x] `countRelationships` -> `internal/domain/user/service/user_relationship_service.go`
  - [x] `countRelationships` -> `internal/domain/user/service/user_relationship_service.go`
  - [x] `friendTwoUsers` -> `internal/domain/user/service/user_relationship_service.go`
  - [x] `upsertOneSidedRelationship` -> `internal/domain/user/service/user_relationship_service.go`
  - [x] `isBlocked` -> `internal/domain/user/service/user_relationship_service.go`
  - [x] `isNotBlocked` -> `internal/domain/user/service/user_relationship_service.go`
  - [x] `hasRelationshipAndNotBlocked` -> `internal/domain/user/service/user_relationship_service.go`
  - [x] `hasRelationshipAndNotBlocked` -> `internal/domain/user/service/user_relationship_service.go`
  - [x] `updateUserOneSidedRelationships` -> `internal/domain/user/service/user_relationship_service.go`
  - [x] `hasOneSidedRelationship` -> `internal/domain/user/service/user_relationship_service.go`

- **UserRoleService.java** (`java/im/turms/service/domain/user/service/UserRoleService.java`)
> [简述功能]

  - [x] `queryUserRoles` -> `internal/domain/user/service/user_role_service.go`
  - [x] `addUserRole` -> `internal/domain/user/service/user_role_service.go`
  - [x] `updateUserRoles` -> `internal/domain/user/service/user_role_service.go`
  - [x] `deleteUserRoles` -> `internal/domain/user/service/user_role_service.go`
  - [x] `queryUserRoleById` -> `internal/domain/user/service/user_role_service.go`
  - [x] `queryStoredOrDefaultUserRoleByUserId` -> `internal/domain/user/service/user_role_service.go`
  - [x] `countUserRoles` -> `internal/domain/user/service/user_role_service.go`

- **UserService.java** (`java/im/turms/service/domain/user/service/UserService.java`)
> [简述功能]

  - [x] `isAllowedToSendMessageToTarget` -> `internal/domain/user/service/user_service.go`
  - [x] `createUser` -> `internal/domain/user/service/user_service.go`
  - [x] `addUser` -> `internal/domain/user/service/user_service.go`
  - [x] `isAllowToQueryUserProfile` -> `internal/domain/user/service/user_service.go`
  - [x] `authAndQueryUsersProfile` -> `internal/domain/user/service/user_service.go`
  - [x] `queryUserName` -> `internal/domain/user/service/user_service.go`
  - [x] `queryUsersProfile` -> `internal/domain/user/service/user_service.go`
  - [x] `queryUserRoleIdByUserId` -> `internal/domain/user/service/user_service.go`
  - [x] `deleteUsers` -> `internal/domain/user/service/user_service.go`
  - [x] `checkIfUserExists` -> `internal/domain/user/service/user_service.go`
  - [x] `updateUser` -> `internal/domain/user/service/user_service.go`
  - [x] `queryUsers` -> `internal/domain/user/service/user_service.go`
  - [x] `countRegisteredUsers` -> `internal/domain/user/service/user_service.go`
  - [x] `countDeletedUsers` -> `internal/domain/user/service/user_service.go`
  - [x] `countUsers` -> `internal/domain/user/service/user_service.go`
  - [x] `countUsers` -> `internal/domain/user/service/user_service.go`
  - [x] `updateUsers` -> `internal/domain/user/service/user_service.go`

- **UserSettingsService.java** (`java/im/turms/service/domain/user/service/UserSettingsService.java`)
> [简述功能]

  - [x] `upsertSettings` -> `internal/domain/user/service/user_settings_service.go`
  - [x] `deleteSettings` -> `internal/domain/user/service/user_settings_service.go`
  - [x] `unsetSettings` -> `internal/domain/user/service/user_settings_service.go`
  - [x] `querySettings` -> `internal/domain/user/service/user_settings_service.go`

- **UserVersionService.java** (`java/im/turms/service/domain/user/service/UserVersionService.java`)
> [简述功能]

  - [x] `queryRelationshipsLastUpdatedDate` -> `internal/domain/user/service/user_version_service.go`
  - [x] `querySentGroupInvitationsLastUpdatedDate` -> `internal/domain/user/service/user_version_service.go`
  - [x] `queryReceivedGroupInvitationsLastUpdatedDate` -> `internal/domain/user/service/user_version_service.go`
  - [x] `queryGroupJoinRequestsVersion` -> `internal/domain/user/service/user_version_service.go`
  - [x] `queryRelationshipGroupsLastUpdatedDate` -> `internal/domain/user/service/user_version_service.go`
  - [x] `queryJoinedGroupVersion` -> `internal/domain/user/service/user_version_service.go`
  - [x] `querySentFriendRequestsVersion` -> `internal/domain/user/service/user_version_service.go`
  - [x] `queryReceivedFriendRequestsVersion` -> `internal/domain/user/service/user_version_service.go`
  - [x] `upsertEmptyUserVersion` -> `internal/domain/user/service/user_version_service.go`
  - [x] `updateRelationshipsVersion` -> `internal/domain/user/service/user_version_service.go`
  - [x] `updateRelationshipsVersion` -> `internal/domain/user/service/user_version_service.go`
  - [x] `updateSentFriendRequestsVersion` -> `internal/domain/user/service/user_version_service.go`
  - [x] `updateReceivedFriendRequestsVersion` -> `internal/domain/user/service/user_version_service.go`
  - [x] `updateRelationshipGroupsVersion` -> `internal/domain/user/service/user_version_service.go`
  - [x] `updateRelationshipGroupsVersion` -> `internal/domain/user/service/user_version_service.go`
  - [x] `updateRelationshipGroupsMembersVersion` -> `internal/domain/user/service/user_version_service.go`
  - [x] `updateRelationshipGroupsMembersVersion` -> `internal/domain/user/service/user_version_service.go`
  - [x] `updateSentGroupInvitationsVersion` -> `internal/domain/user/service/user_version_service.go`
  - [x] `updateReceivedGroupInvitationsVersion` -> `internal/domain/user/service/user_version_service.go`
  - [x] `updateSentGroupJoinRequestsVersion` -> `internal/domain/user/service/user_version_service.go`
  - [x] `updateJoinedGroupsVersion` -> `internal/domain/user/service/user_version_service.go`
  - [x] `updateSpecificVersion` -> `internal/domain/user/service/user_version_service.go`
  - [x] `updateSpecificVersion` -> `internal/domain/user/service/user_version_service.go`
  - [x] `updateSpecificVersion` -> `internal/domain/user/service/user_version_service.go`
  - [x] `delete` -> `internal/domain/user/service/user_version_service.go`

- **NearbyUserService.java** (`java/im/turms/service/domain/user/service/onlineuser/NearbyUserService.java`)
> [简述功能]

  - [x] `queryNearbyUsers` -> `internal/domain/user/service/onlineuser/nearby_user_service.go`

- **SessionService.java** (`java/im/turms/service/domain/user/service/onlineuser/SessionService.java`)
> [简述功能]

  - [x] `disconnect` -> `internal/domain/user/service/onlineuser/session_service.go`
  - [x] `disconnect` -> `internal/domain/user/service/onlineuser/session_service.go`
  - [x] `disconnect` -> `internal/domain/user/service/onlineuser/session_service.go`
  - [x] `disconnect` -> `internal/domain/user/service/onlineuser/session_service.go`
  - [x] `disconnect` -> `internal/domain/user/service/onlineuser/session_service.go`
  - [x] `queryUserSessions` -> `internal/domain/user/service/onlineuser/session_service.go`

- **LocaleUtil.java** (`java/im/turms/service/infra/locale/LocaleUtil.java`)
> [简述功能]

  - [x] `isAvailableLanguage` -> `internal/infra/locale/locale_util.go`

- **ApiLoggingContext.java** (`java/im/turms/service/infra/logging/ApiLoggingContext.java`)
> [简述功能]

  - [x] `shouldLogRequest` -> `internal/infra/logging/api_logging_context.go`
  - [x] `shouldLogNotification` -> `internal/infra/logging/api_logging_context.go`

- **ClientApiLogging.java** (`java/im/turms/service/infra/logging/ClientApiLogging.java`)
> [简述功能]

  - [x] `log` -> `internal/infra/logging/client_api_logging.go`

- **AcceptMeetingInvitationResult.java** (`java/im/turms/service/infra/plugin/extension/model/AcceptMeetingInvitationResult.java`)
> [简述功能]

  - [ ] `AcceptMeetingInvitationResult`

- **CreateMeetingOptions.java** (`java/im/turms/service/infra/plugin/extension/model/CreateMeetingOptions.java`)
> [简述功能]

  - [ ] `CreateMeetingOptions`

- **CreateMeetingResult.java** (`java/im/turms/service/infra/plugin/extension/model/CreateMeetingResult.java`)
> [简述功能]

  - [ ] `CreateMeetingResult`

- **ProtoModelConvertor.java** (`java/im/turms/service/infra/proto/ProtoModelConvertor.java`)
> [简述功能]

  - [x] `toList` -> `internal/infra/proto/proto_model_convertor.go`
  - [x] `value2proto` -> `internal/infra/proto/proto_model_convertor.go`

- **DefaultLanguageSettings.java** (`java/im/turms/service/storage/elasticsearch/DefaultLanguageSettings.java`)
> [简述功能]

  - [x] `getSetting` -> `internal/storage/elasticsearch/default_language_settings.go`

- **ElasticsearchClient.java** (`java/im/turms/service/storage/elasticsearch/ElasticsearchClient.java`)
> [简述功能]

  - [x] `healthcheck` -> `internal/storage/elasticsearch/elasticsearch_client.go`
  - [x] `putIndex` -> `internal/storage/elasticsearch/elasticsearch_client.go`
  - [x] `putDoc` -> `internal/storage/elasticsearch/elasticsearch_client.go`
  - [x] `deleteDoc` -> `internal/storage/elasticsearch/elasticsearch_client.go`
  - [x] `deleteByQuery` -> `internal/storage/elasticsearch/elasticsearch_client.go`
  - [x] `updateByQuery` -> `internal/storage/elasticsearch/elasticsearch_client.go`
  - [x] `search` -> `internal/storage/elasticsearch/elasticsearch_client.go`
  - [x] `bulk` -> `internal/storage/elasticsearch/elasticsearch_client.go`
  - [x] `deletePit` -> `internal/storage/elasticsearch/elasticsearch_client.go`

- **ElasticsearchManager.java** (`java/im/turms/service/storage/elasticsearch/ElasticsearchManager.java`)
> [简述功能]

  - [x] `putUserDoc` -> `internal/storage/elasticsearch/elasticsearch_manager.go`
  - [x] `putUserDocs` -> `internal/storage/elasticsearch/elasticsearch_manager.go`
  - [x] `deleteUserDoc` -> `internal/storage/elasticsearch/elasticsearch_manager.go`
  - [x] `deleteUserDocs` -> `internal/storage/elasticsearch/elasticsearch_manager.go`
  - [x] `searchUserDocs` -> `internal/storage/elasticsearch/elasticsearch_manager.go`
  - [x] `putGroupDoc` -> `internal/storage/elasticsearch/elasticsearch_manager.go`
  - [x] `putGroupDocs` -> `internal/storage/elasticsearch/elasticsearch_manager.go`
  - [x] `deleteGroupDocs` -> `internal/storage/elasticsearch/elasticsearch_manager.go`
  - [x] `deleteAllGroupDocs` -> `internal/storage/elasticsearch/elasticsearch_manager.go`
  - [x] `searchGroupDocs` -> `internal/storage/elasticsearch/elasticsearch_manager.go`
  - [x] `deletePitForUserDocs` -> `internal/storage/elasticsearch/elasticsearch_manager.go`

- **IndexTextFieldSetting.java** (`java/im/turms/service/storage/elasticsearch/IndexTextFieldSetting.java`)
> [简述功能]

  - [ ] `IndexTextFieldSetting`

- **BulkRequest.java** (`java/im/turms/service/storage/elasticsearch/model/BulkRequest.java`)
> [简述功能]

  - [ ] `BulkRequest`
  - [x] `serialize` -> `internal/storage/elasticsearch/model/elasticsearch_model.go`

- **BulkResponse.java** (`java/im/turms/service/storage/elasticsearch/model/BulkResponse.java`)
> [简述功能]

  - [ ] `BulkResponse`

- **BulkResponseItem.java** (`java/im/turms/service/storage/elasticsearch/model/BulkResponseItem.java`)
> [简述功能]

  - [ ] `BulkResponseItem`

- **ClosePointInTimeRequest.java** (`java/im/turms/service/storage/elasticsearch/model/ClosePointInTimeRequest.java`)
> [简述功能]

  - [ ] `ClosePointInTimeRequest`

- **CreateIndexRequest.java** (`java/im/turms/service/storage/elasticsearch/model/CreateIndexRequest.java`)
> [简述功能]

  - [ ] `CreateIndexRequest`

- **DeleteByQueryRequest.java** (`java/im/turms/service/storage/elasticsearch/model/DeleteByQueryRequest.java`)
> [简述功能]

  - [ ] `DeleteByQueryRequest`

- **DeleteByQueryResponse.java** (`java/im/turms/service/storage/elasticsearch/model/DeleteByQueryResponse.java`)
> [简述功能]

  - [ ] `DeleteByQueryResponse`

- **DeleteResponse.java** (`java/im/turms/service/storage/elasticsearch/model/DeleteResponse.java`)
> [简述功能]

  - [ ] `DeleteResponse`

- **ErrorCause.java** (`java/im/turms/service/storage/elasticsearch/model/ErrorCause.java`)
> [简述功能]

  - [ ] `ErrorCause`

- **ErrorResponse.java** (`java/im/turms/service/storage/elasticsearch/model/ErrorResponse.java`)
> [简述功能]

  - [ ] `ErrorResponse`

- **FieldCollapse.java** (`java/im/turms/service/storage/elasticsearch/model/FieldCollapse.java`)
> [简述功能]

  - [ ] `FieldCollapse`

- **HealthResponse.java** (`java/im/turms/service/storage/elasticsearch/model/HealthResponse.java`)
> [简述功能]

  - [ ] `HealthResponse`

- **Highlight.java** (`java/im/turms/service/storage/elasticsearch/model/Highlight.java`)
> [简述功能]

  - [ ] `Highlight`

- **IndexSettings.java** (`java/im/turms/service/storage/elasticsearch/model/IndexSettings.java`)
> [简述功能]

  - [ ] `IndexSettings`

- **IndexSettingsAnalysis.java** (`java/im/turms/service/storage/elasticsearch/model/IndexSettingsAnalysis.java`)
> [简述功能]

  - [ ] `IndexSettingsAnalysis`
  - [x] `merge` -> `internal/storage/elasticsearch/model/elasticsearch_model.go`

- **PointInTimeReference.java** (`java/im/turms/service/storage/elasticsearch/model/PointInTimeReference.java`)
> [简述功能]

  - [ ] `PointInTimeReference`

- **Property.java** (`java/im/turms/service/storage/elasticsearch/model/Property.java`)
> [简述功能]

  - [ ] `Property`

- **Script.java** (`java/im/turms/service/storage/elasticsearch/model/Script.java`)
> [简述功能]

  - [ ] `Script`

- **SearchRequest.java** (`java/im/turms/service/storage/elasticsearch/model/SearchRequest.java`)
> [简述功能]

  - [ ] `SearchRequest`

- **ShardFailure.java** (`java/im/turms/service/storage/elasticsearch/model/ShardFailure.java`)
> [简述功能]

  - [ ] `ShardFailure`

- **ShardStatistics.java** (`java/im/turms/service/storage/elasticsearch/model/ShardStatistics.java`)
> [简述功能]

  - [ ] `ShardStatistics`

- **TypeMapping.java** (`java/im/turms/service/storage/elasticsearch/model/TypeMapping.java`)
> [简述功能]

  - [ ] `TypeMapping`

- **UpdateByQueryRequest.java** (`java/im/turms/service/storage/elasticsearch/model/UpdateByQueryRequest.java`)
> [简述功能]

  - [ ] `UpdateByQueryRequest`

- **UpdateByQueryResponse.java** (`java/im/turms/service/storage/elasticsearch/model/UpdateByQueryResponse.java`)
> [简述功能]

  - [ ] `UpdateByQueryResponse`

- **MongoCollectionMigrator.java** (`java/im/turms/service/storage/mongo/MongoCollectionMigrator.java`)
> [简述功能]

  - [x] `migrate` -> `internal/storage/mongo/mongo_collection_migrator.go`

- **MongoConfig.java** (`java/im/turms/service/storage/mongo/MongoConfig.java`)
> [简述功能]

  - [x] `adminMongoClient` -> `internal/storage/mongo/mongo_config.go`
  - [x] `userMongoClient` -> `internal/storage/mongo/mongo_config.go`
  - [x] `groupMongoClient` -> `internal/storage/mongo/mongo_config.go`
  - [x] `conversationMongoClient` -> `internal/storage/mongo/mongo_config.go`
  - [x] `messageMongoClient` -> `internal/storage/mongo/mongo_config.go`
  - [x] `conferenceMongoClient` -> `internal/storage/mongo/mongo_config.go`

- **MongoFakeDataGenerator.java** (`java/im/turms/service/storage/mongo/MongoFakeDataGenerator.java`)
> [简述功能]

  - [x] `populateCollectionsWithFakeData` -> `internal/storage/mongo/mongo_fake_data_generator.go`

- **RedisConfig.java** (`java/im/turms/service/storage/redis/RedisConfig.java`)
> [简述功能]

  - [x] `newSequenceIdRedisClientManager` -> `internal/storage/redis/redis_config.go`
  - [x] `sequenceIdRedisClientManager` -> `internal/storage/redis/redis_config.go`

