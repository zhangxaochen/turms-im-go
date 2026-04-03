# Turms Refactoring Progress Report

## Modules

### turms-gateway

> [简述功能]

#### Configurations

- **application-demo.yaml** ([resources/application-demo.yaml](../turms-orig/turms-gateway/src/main/resources/application-demo.yaml)): [简述功能]
- **application-dev.yaml** ([resources/application-dev.yaml](../turms-orig/turms-gateway/src/main/resources/application-dev.yaml)): [简述功能]
- **application-test.yaml** ([resources/application-test.yaml](../turms-orig/turms-gateway/src/main/resources/application-test.yaml)): [简述功能]
- **application.yaml** ([resources/application.yaml](../turms-orig/turms-gateway/src/main/resources/application.yaml)): [简述功能]

#### Java source tracking

- **TurmsGatewayApplication.java** ([java/im/turms/gateway/TurmsGatewayApplication.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/TurmsGatewayApplication.java))
> [简述功能]

  - [x] `main(String[] args)` -> [cmd/turms-gateway/main.go:main()](../cmd/turms-gateway/main.go)

- **ClientRequestDispatcher.java** ([java/im/turms/gateway/access/client/common/ClientRequestDispatcher.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/common/ClientRequestDispatcher.java))
> [简述功能]

  - [x] `handleRequest(UserSessionWrapper sessionWrapper, ByteBuf serviceRequestBuffer)` -> [internal/domain/gateway/access/client/common/client_request_dispatcher.go:HandleRequest(ctx context.Context, sessionWrapper *UserSessionWrapper, serviceRequestBuffer []byte)](../internal/domain/gateway/access/client/common/client_request_dispatcher.go)
  - [x] `handleRequest0(UserSessionWrapper sessionWrapper, ByteBuf serviceRequestBuffer)` -> [internal/domain/gateway/access/client/common/client_request_dispatcher.go:HandleRequest0(ctx context.Context, sessionWrapper *UserSessionWrapper, serviceRequestBuffer []byte)](../internal/domain/gateway/access/client/common/client_request_dispatcher.go)
  - [x] `handleServiceRequest(UserSessionWrapper sessionWrapper, SimpleTurmsRequest request, ByteBuf serviceRequestBuffer, TracingContext tracingContext)` -> [internal/domain/gateway/access/client/common/client_request_dispatcher.go:HandleServiceRequest(ctx context.Context, sessionWrapper *UserSessionWrapper, request *protocol.TurmsRequest, serviceRequestBuffer []byte)](../internal/domain/gateway/access/client/common/client_request_dispatcher.go)

- **IpRequestThrottler.java** ([java/im/turms/gateway/access/client/common/IpRequestThrottler.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/common/IpRequestThrottler.java))
> [简述功能]

  - [x] `tryAcquireToken(ByteArrayWrapper ip, long timestamp)` -> [internal/domain/gateway/access/client/common/ip_request_throttler.go:TryAcquireToken(ip string)](../internal/domain/gateway/access/client/common/ip_request_throttler.go)

- **NotificationFactory.java** ([java/im/turms/gateway/access/client/common/NotificationFactory.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/common/NotificationFactory.java))
> [简述功能]

  - [x] `init(TurmsPropertiesManager propertiesManager)` -> [internal/domain/gateway/access/client/common/notification_factory.go:NewNotificationFactory(props *config.GatewayProperties)](../internal/domain/gateway/access/client/common/notification_factory.go)
  - [x] `create(ResponseStatusCode code, long requestId)` -> [internal/domain/gateway/access/client/common/notification_factory.go:Create(requestID *int64, code constant.ResponseStatusCode)](../internal/domain/gateway/access/client/common/notification_factory.go)
  - [x] `create(ResponseStatusCode code, @Nullable String reason, long requestId)` -> [internal/domain/gateway/access/client/common/notification_factory.go:CreateWithReason(requestID *int64, code constant.ResponseStatusCode, reason string)](../internal/domain/gateway/access/client/common/notification_factory.go)
  - [x] `create(ThrowableInfo info, long requestId)` -> [internal/domain/gateway/access/client/common/notification_factory.go:CreateFromError(err error, requestID *int64)](../internal/domain/gateway/access/client/common/notification_factory.go)
  - [x] `createBuffer(CloseReason closeReason)` -> [internal/domain/gateway/access/client/common/notification_factory.go:CreateBuffer(requestID *int64, code constant.ResponseStatusCode, reason string)](../internal/domain/gateway/access/client/common/notification_factory.go)
  - [x] `sessionClosed(long requestId)` -> [internal/domain/gateway/access/client/common/notification_factory.go:SessionClosed(requestID *int64)](../internal/domain/gateway/access/client/common/notification_factory.go)

- **RequestHandlerResult.java** ([java/im/turms/gateway/access/client/common/RequestHandlerResult.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/common/RequestHandlerResult.java))
> [简述功能]

  - [x] `RequestHandlerResult(ResponseStatusCode code, String reason)` -> [internal/domain/gateway/access/client/common/request_handler_result.go:NewRequestHandlerResult(code constant.ResponseStatusCode, reason string)](../internal/domain/gateway/access/client/common/request_handler_result.go)

- **UserSession.java** ([java/im/turms/gateway/access/client/common/UserSession.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/common/UserSession.java))
> [简述功能]

  - [x] `setConnection(NetConnection connection, ByteArrayWrapper ip)` -> [internal/domain/gateway/session/connection.go:SetConnection(connection Connection, ip string)](../internal/domain/gateway/session/connection.go)
  - [x] `setLastHeartbeatRequestTimestampToNow()` -> [internal/domain/gateway/session/connection.go:SetLastHeartbeatRequestTimestampToNow()](../internal/domain/gateway/session/connection.go)
  - [x] `setLastRequestTimestampToNow()` -> [internal/domain/gateway/session/connection.go:SetLastRequestTimestampToNow()](../internal/domain/gateway/session/connection.go)
  - [x] `close(@NotNull CloseReason closeReason)` -> [internal/domain/gateway/session/connection.go:Close(closeReason any)](../internal/domain/gateway/session/connection.go)
  - [x] `isOpen()` -> [internal/domain/gateway/session/connection.go:IsOpen()](../internal/domain/gateway/session/connection.go)
  - [x] `isConnected()` -> [internal/domain/gateway/session/connection.go:IsConnected()](../internal/domain/gateway/session/connection.go)
  - [x] `supportsSwitchingToUdp()` -> [internal/domain/gateway/session/connection.go:SupportsSwitchingToUdp()](../internal/domain/gateway/session/connection.go)
  - [x] `sendNotification(ByteBuf byteBuf)` -> [internal/domain/gateway/access/router/router.go:sendNotification(s *session.UserSession, requestID *int64, code int32, reason string)](../internal/domain/gateway/access/router/router.go)
  - [x] `sendNotification(ByteBuf byteBuf, TracingContext tracingContext)` -> [internal/domain/gateway/access/router/router.go:sendNotification(s *session.UserSession, requestID *int64, code int32, reason string)](../internal/domain/gateway/access/router/router.go)
  - [x] `acquireDeleteSessionRequestLoggingLock()` -> [internal/domain/gateway/session/connection.go:AcquireDeleteSessionRequestLoggingLock()](../internal/domain/gateway/session/connection.go)
  - [x] `hasPermission(TurmsRequest.KindCase requestType)` -> [internal/domain/gateway/session/connection.go:HasPermission(requestType any)](../internal/domain/gateway/session/connection.go)
  - [x] `toString()` -> [internal/domain/gateway/session/connection.go:ToString()](../internal/domain/gateway/session/connection.go)

- **UserSessionWrapper.java** ([java/im/turms/gateway/access/client/common/UserSessionWrapper.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/common/UserSessionWrapper.java))
> [简述功能]

  - [x] `getIp()` -> [internal/domain/gateway/access/client/common/user_session_wrapper.go:GetIP()](../internal/domain/gateway/access/client/common/user_session_wrapper.go)
  - [x] `getIpStr()` -> [internal/domain/gateway/access/client/common/user_session_wrapper.go:GetIPStr()](../internal/domain/gateway/access/client/common/user_session_wrapper.go)
  - [x] `setUserSession(UserSession userSession)` -> [internal/domain/gateway/access/client/common/user_session_wrapper.go:SetUserSession(userSession *session.UserSession)](../internal/domain/gateway/access/client/common/user_session_wrapper.go)
  - [x] `hasUserSession()` -> [internal/domain/gateway/access/client/common/user_session_wrapper.go:HasUserSession()](../internal/domain/gateway/access/client/common/user_session_wrapper.go)

- **Policy.java** ([java/im/turms/gateway/access/client/common/authorization/policy/Policy.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/common/authorization/policy/Policy.java))
> [简述功能]

  - [x] `Policy(List<PolicyStatement> statements)` -> [internal/domain/gateway/access/client/common/authorization/policy.go:NewPolicy(statements []PolicyStatement)](../internal/domain/gateway/access/client/common/authorization/policy.go)

- **PolicyDeserializer.java** ([java/im/turms/gateway/access/client/common/authorization/policy/PolicyDeserializer.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/common/authorization/policy/PolicyDeserializer.java))
> [简述功能]

  - [x] `parse(Map<String, Object> map)` -> [internal/domain/gateway/access/client/common/authorization/policy.go:Parse(data map[string]interface{})](../internal/domain/gateway/access/client/common/authorization/policy.go)

- **PolicyStatement.java** ([java/im/turms/gateway/access/client/common/authorization/policy/PolicyStatement.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/common/authorization/policy/PolicyStatement.java))
> [简述功能]

  - [x] `PolicyStatement(PolicyStatementEffect effect, Set<PolicyStatementAction> actions, Set<PolicyStatementResource> resources)` -> [internal/domain/gateway/access/client/common/authorization/policy.go:NewPolicyStatement(effect PolicyStatementEffect, actions []PolicyStatementAction, resources []PolicyStatementResource)](../internal/domain/gateway/access/client/common/authorization/policy.go)

- **ServiceAvailabilityHandler.java** ([java/im/turms/gateway/access/client/common/channel/ServiceAvailabilityHandler.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/common/channel/ServiceAvailabilityHandler.java))
> [简述功能]

  - [x] `channelRegistered(ChannelHandlerContext ctx)` -> [internal/domain/gateway/access/client/common/service_availability.go:ChannelRegistered(isAvailable bool)](../internal/domain/gateway/access/client/common/service_availability.go)
  - [x] `exceptionCaught(ChannelHandlerContext ctx, Throwable cause)` -> [internal/domain/gateway/access/client/common/service_availability.go:ExceptionCaught(err error)](../internal/domain/gateway/access/client/common/service_availability.go)

- **NetConnection.java** ([java/im/turms/gateway/access/client/common/connection/NetConnection.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/common/connection/NetConnection.java))
> [简述功能]

  - [x] `getAddress()` -> [internal/domain/gateway/access/client/tcp/tcp_server.go:GetAddress()](../internal/domain/gateway/access/client/tcp/tcp_server.go)
  - [x] `send(ByteBuf buffer)` -> [internal/domain/gateway/access/client/tcp/tcp_server.go:Send(ctx context.Context, buffer []byte)](../internal/domain/gateway/access/client/tcp/tcp_server.go)
  - [x] `close(CloseReason closeReason)` -> [internal/domain/gateway/access/client/common/net_connection.go:CloseWithReason(reason CloseReason)](../internal/domain/gateway/access/client/common/net_connection.go)
  - [x] `close()` -> [internal/domain/gateway/access/client/common/net_connection.go:Close()](../internal/domain/gateway/access/client/common/net_connection.go)
  - [x] `switchToUdp()` -> [internal/domain/gateway/access/client/common/net_connection.go:SwitchToUdp()](../internal/domain/gateway/access/client/common/net_connection.go)
  - [x] `tryNotifyClientToRecover()` -> [internal/domain/gateway/access/client/common/net_connection.go:TryNotifyClientToRecover()](../internal/domain/gateway/access/client/common/net_connection.go)

- **ExtendedHAProxyMessageReader.java** ([java/im/turms/gateway/access/client/tcp/ExtendedHAProxyMessageReader.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/tcp/ExtendedHAProxyMessageReader.java))
> [简述功能]

  - [x] `channelRead(ChannelHandlerContext ctx, Object msg)` -> [internal/domain/gateway/access/client/tcp/haproxy.go:Read(conn net.Conn)](../internal/domain/gateway/access/client/tcp/haproxy.go)

- **HAProxyUtil.java** ([java/im/turms/gateway/access/client/tcp/HAProxyUtil.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/tcp/HAProxyUtil.java))
> [简述功能]

  - [x] `addProxyProtocolHandlers(ChannelPipeline pipeline, Consumer<InetSocketAddress> onRemoteAddressConfirmed)` -> [internal/domain/gateway/access/client/tcp/haproxy.go:AddProxyProtocolHandlers(callback func(net.Addr)](../internal/domain/gateway/access/client/tcp/haproxy.go)
  - [x] `addProxyProtocolDetectorHandler(ChannelPipeline pipeline, Consumer<InetSocketAddress> onRemoteAddressConfirmed)` -> [internal/domain/gateway/access/client/tcp/haproxy.go:AddProxyProtocolDetectorHandler(callback func(net.Addr)](../internal/domain/gateway/access/client/tcp/haproxy.go)

- **TcpConnection.java** ([java/im/turms/gateway/access/client/tcp/TcpConnection.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/tcp/TcpConnection.java))
> [简述功能]

  - [x] `getAddress()` -> [internal/domain/gateway/access/client/tcp/tcp_server.go:GetAddress()](../internal/domain/gateway/access/client/tcp/tcp_server.go)
  - [x] `send(ByteBuf buffer)` -> [internal/domain/gateway/access/client/tcp/tcp_server.go:Send(ctx context.Context, buffer []byte)](../internal/domain/gateway/access/client/tcp/tcp_server.go)
  - [x] `close(CloseReason closeReason)` -> [internal/domain/gateway/access/client/common/net_connection.go:CloseWithReason(reason CloseReason)](../internal/domain/gateway/access/client/common/net_connection.go)
  - [x] `close()` -> [internal/domain/gateway/access/client/common/net_connection.go:Close()](../internal/domain/gateway/access/client/common/net_connection.go)

- **TcpServerFactory.java** ([java/im/turms/gateway/access/client/tcp/TcpServerFactory.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/tcp/TcpServerFactory.java))
> [简述功能]

  - [x] `create(TcpProperties tcpProperties, BlocklistService blocklistService, ServerStatusManager serverStatusManager, SessionService sessionService, ConnectionListener connectionListener, int maxFrameLength)` -> [internal/domain/gateway/access/client/tcp/tcp_server.go:CreateWithArgs(tcpProperties any, blocklistService any, serverStatusManager any, sessionService any, connectionListener any, maxFrameLength int)](../internal/domain/gateway/access/client/tcp/tcp_server.go)

- **TcpUserSessionAssembler.java** ([java/im/turms/gateway/access/client/tcp/TcpUserSessionAssembler.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/tcp/TcpUserSessionAssembler.java))
> [简述功能]

  - [x] `getHost()` -> [internal/domain/gateway/access/client/tcp/tcp_server.go:GetHost()](../internal/domain/gateway/access/client/tcp/tcp_server.go)
  - [x] `getPort()` -> [internal/domain/gateway/access/client/tcp/tcp_server.go:GetPort()](../internal/domain/gateway/access/client/tcp/tcp_server.go)

- **UdpRequestDispatcher.java** ([java/im/turms/gateway/access/client/udp/UdpRequestDispatcher.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/udp/UdpRequestDispatcher.java))
> [简述功能]

  - [x] `sendSignal(InetSocketAddress address, UdpNotificationType signal)` -> [internal/domain/gateway/access/client/udp/udp_server.go:SendSignal(address net.Addr, signal UdpNotificationType)](../internal/domain/gateway/access/client/udp/udp_server.go)

- **UdpSignalResponseBufferPool.java** ([java/im/turms/gateway/access/client/udp/UdpSignalResponseBufferPool.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/udp/UdpSignalResponseBufferPool.java))
> [简述功能]

  - [x] `get(ResponseStatusCode code)` -> [internal/domain/common/cache/sharded_map.go:Get(key K)](../internal/domain/common/cache/sharded_map.go)
  - [x] `get(UdpNotificationType type)` -> [internal/domain/common/cache/sharded_map.go:Get(key K)](../internal/domain/common/cache/sharded_map.go)

- **UdpNotification.java** ([java/im/turms/gateway/access/client/udp/dto/UdpNotification.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/udp/dto/UdpNotification.java))
> [简述功能]

  - [x] `UdpNotification(InetSocketAddress recipientAddress, UdpNotificationType type)` -> [internal/domain/gateway/access/client/udp/udp_server.go:NewUdpNotification(recipientAddress net.Addr, notificationType UdpNotificationType)](../internal/domain/gateway/access/client/udp/udp_server.go)

- **UdpRequestType.java** ([java/im/turms/gateway/access/client/udp/dto/UdpRequestType.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/udp/dto/UdpRequestType.java))
> [简述功能]

  - [x] `parse(int number)` -> [internal/domain/gateway/access/client/udp/udp_server.go:ParseUdpRequestType(number int)](../internal/domain/gateway/access/client/udp/udp_server.go)
  - [x] `getNumber()` -> [internal/domain/gateway/access/client/udp/udp_server.go:GetNumber()](../internal/domain/gateway/access/client/udp/udp_server.go)

- **UdpSignalRequest.java** ([java/im/turms/gateway/access/client/udp/dto/UdpSignalRequest.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/udp/dto/UdpSignalRequest.java))
> [简述功能]

  - [x] `UdpSignalRequest(UdpRequestType type, long userId, DeviceType deviceType, int sessionId)` -> [internal/domain/gateway/access/client/udp/udp_server.go:NewUdpSignalRequest(reqType UdpRequestType, userID int64, deviceType protocol.DeviceType, sessionID int)](../internal/domain/gateway/access/client/udp/udp_server.go)

- **HttpForwardedHeaderHandler.java** ([java/im/turms/gateway/access/client/websocket/HttpForwardedHeaderHandler.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/websocket/HttpForwardedHeaderHandler.java))
> [简述功能]

  - [x] `apply(ConnectionInfo connectionInfo, HttpRequest request)` -> [internal/domain/gateway/access/client/ws/ws_server.go:Apply(connectionInfo any, request any)](../internal/domain/gateway/access/client/ws/ws_server.go)

- **WebSocketConnection.java** ([java/im/turms/gateway/access/client/websocket/WebSocketConnection.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/websocket/WebSocketConnection.java))
> [简述功能]

  - [x] `getAddress()` -> [internal/domain/gateway/access/client/tcp/tcp_server.go:GetAddress()](../internal/domain/gateway/access/client/tcp/tcp_server.go)
  - [x] `send(ByteBuf buffer)` -> [internal/domain/gateway/access/client/tcp/tcp_server.go:Send(ctx context.Context, buffer []byte)](../internal/domain/gateway/access/client/tcp/tcp_server.go)
  - [x] `close(CloseReason closeReason)` -> [internal/domain/gateway/access/client/common/net_connection.go:CloseWithReason(reason CloseReason)](../internal/domain/gateway/access/client/common/net_connection.go)
  - [x] `close()` -> [internal/domain/gateway/access/client/common/net_connection.go:Close()](../internal/domain/gateway/access/client/common/net_connection.go)

- **WebSocketServerFactory.java** ([java/im/turms/gateway/access/client/websocket/WebSocketServerFactory.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/websocket/WebSocketServerFactory.java))
> [简述功能]

  - [x] `create(WebSocketProperties webSocketProperties, BlocklistService blocklistService, ServerStatusManager serverStatusManager, SessionService sessionService, ConnectionListener connectionListener, int maxFramePayloadLength)` -> [internal/domain/gateway/access/client/ws/ws_server.go:Create(webSocketProperties any, blocklistService any, serverStatusManager any, sessionService *session.SessionService, connectionListener any, maxFramePayloadLength int)](../internal/domain/gateway/access/client/ws/ws_server.go)

- **NotificationService.java** ([java/im/turms/gateway/domain/notification/service/NotificationService.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/domain/notification/service/NotificationService.java))
> [简述功能]

  - [ ] `sendNotificationToLocalClients(TracingContext tracingContext, ByteBuf notificationData, Set<Long> recipientIds, Set<UserSessionId> excludedUserSessionIds, @Nullable DeviceType excludedDeviceType)`

- **StatisticsService.java** ([java/im/turms/gateway/domain/observation/service/StatisticsService.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/domain/observation/service/StatisticsService.java))
> [简述功能]

  - [ ] `countLocalOnlineUsers()`

- **ServiceRequestService.java** ([java/im/turms/gateway/domain/servicerequest/service/ServiceRequestService.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/domain/servicerequest/service/ServiceRequestService.java))
> [简述功能]

  - [x] `handleServiceRequest(UserSession session, ServiceRequest serviceRequest)` -> [internal/domain/gateway/servicerequest/service/servicerequest_service.go:HandleServiceRequest](../internal/domain/gateway/servicerequest/service/servicerequest_service.go)

- **SessionController.java** ([java/im/turms/gateway/domain/session/access/admin/controller/SessionController.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/domain/session/access/admin/controller/SessionController.java))
> [简述功能]

  - [ ] `deleteSessions(@QueryParam(required = false)`

- **SessionClientController.java** ([java/im/turms/gateway/domain/session/access/client/controller/SessionClientController.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/domain/session/access/client/controller/SessionClientController.java))
> [简述功能]

  - [x] `handleDeleteSessionRequest(UserSessionWrapper sessionWrapper)` -> [internal/domain/gateway/session/access/client/controller/session_client_controller.go:HandleDeleteSessionRequest](../internal/domain/gateway/session/access/client/controller/session_client_controller.go)
  - [x] `handleCreateSessionRequest(UserSessionWrapper sessionWrapper, CreateSessionRequest createSessionRequest)` -> [internal/domain/gateway/session/access/client/controller/session_client_controller.go:HandleCreateSessionRequest](../internal/domain/gateway/session/access/client/controller/session_client_controller.go)

- **UserLoginInfo.java** ([java/im/turms/gateway/domain/session/bo/UserLoginInfo.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/domain/session/bo/UserLoginInfo.java))
> [简述功能]

  - [ ] `UserLoginInfo(int version, Long userId, String password, DeviceType loggingInDeviceType, Map<String, String> deviceDetails, UserStatus userStatus, Location location, String ip)`

- **UserPermissionInfo.java** ([java/im/turms/gateway/domain/session/bo/UserPermissionInfo.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/domain/session/bo/UserPermissionInfo.java))
> [简述功能]

  - [ ] `UserPermissionInfo(ResponseStatusCode authenticationCode, Set<TurmsRequest.KindCase> permissions)`

- **HeartbeatManager.java** ([java/im/turms/gateway/domain/session/manager/HeartbeatManager.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/domain/session/manager/HeartbeatManager.java))
> [简述功能]

  - [ ] `setCloseIdleSessionAfterSeconds(int closeIdleSessionAfterSeconds)`
  - [ ] `setClientHeartbeatIntervalSeconds(int clientHeartbeatIntervalSeconds)`
  - [ ] `destroy()`
  - [ ] `estimatedSize()`
  - [ ] `next()`

- **UserSessionsManager.java** ([java/im/turms/gateway/domain/session/manager/UserSessionsManager.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/domain/session/manager/UserSessionsManager.java))
> [简述功能]

  - [ ] `addSessionIfAbsent(int version, Set<TurmsRequest.KindCase> permissions, DeviceType loggingInDeviceType, Map<String, String> deviceDetails, @Nullable Location location)`
  - [ ] `closeSession(@NotNull DeviceType deviceType, @NotNull CloseReason closeReason)`
  - [ ] `pushSessionNotification(DeviceType deviceType, String serverId)`
  - [x] `getSession(@NotNull DeviceType deviceType)` -> [internal/domain/gateway/session/sharded_map.go:GetSession(deviceType protocol.DeviceType)](../internal/domain/gateway/session/sharded_map.go)
  - [ ] `countSessions()`
  - [ ] `getLoggedInDeviceTypes()`

- **UserRepository.java** ([java/im/turms/gateway/domain/session/repository/UserRepository.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/domain/session/repository/UserRepository.java))
> [简述功能]

  - [ ] `findPassword(Long userId)`
  - [x] `isActiveAndNotDeleted(Long userId)` -> [internal/domain/user/repository/user_repository.go:IsActiveAndNotDeleted(ctx context.Context, userID int64)](../internal/domain/user/repository/user_repository.go)

- **HttpSessionIdentityAccessManager.java** ([java/im/turms/gateway/domain/session/service/HttpSessionIdentityAccessManager.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/domain/session/service/HttpSessionIdentityAccessManager.java))
> [简述功能]

  - [x] `verifyAndGrant(UserLoginInfo userLoginInfo)` -> [internal/domain/gateway/session/identity_access_manager.go:VerifyAndGrant(ctx context.Context, loginInfo *bo.UserLoginInfo)](../internal/domain/gateway/session/identity_access_manager.go)

- **JwtSessionIdentityAccessManager.java** ([java/im/turms/gateway/domain/session/service/JwtSessionIdentityAccessManager.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/domain/session/service/JwtSessionIdentityAccessManager.java))
> [简述功能]

  - [x] `verifyAndGrant(UserLoginInfo userLoginInfo)` -> [internal/domain/gateway/session/identity_access_manager.go:VerifyAndGrant(ctx context.Context, loginInfo *bo.UserLoginInfo)](../internal/domain/gateway/session/identity_access_manager.go)

- **LdapSessionIdentityAccessManager.java** ([java/im/turms/gateway/domain/session/service/LdapSessionIdentityAccessManager.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/domain/session/service/LdapSessionIdentityAccessManager.java))
> [简述功能]

  - [x] `verifyAndGrant(UserLoginInfo userLoginInfo)` -> [internal/domain/gateway/session/identity_access_manager.go:VerifyAndGrant(ctx context.Context, loginInfo *bo.UserLoginInfo)](../internal/domain/gateway/session/identity_access_manager.go)

- **NoopSessionIdentityAccessManager.java** ([java/im/turms/gateway/domain/session/service/NoopSessionIdentityAccessManager.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/domain/session/service/NoopSessionIdentityAccessManager.java))
> [简述功能]

  - [x] `verifyAndGrant(UserLoginInfo userLoginInfo)` -> [internal/domain/gateway/session/identity_access_manager.go:VerifyAndGrant(ctx context.Context, loginInfo *bo.UserLoginInfo)](../internal/domain/gateway/session/identity_access_manager.go)

- **PasswordSessionIdentityAccessManager.java** ([java/im/turms/gateway/domain/session/service/PasswordSessionIdentityAccessManager.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/domain/session/service/PasswordSessionIdentityAccessManager.java))
> [简述功能]

  - [x] `verifyAndGrant(UserLoginInfo userLoginInfo)` -> [internal/domain/gateway/session/identity_access_manager.go:VerifyAndGrant(ctx context.Context, loginInfo *bo.UserLoginInfo)](../internal/domain/gateway/session/identity_access_manager.go)
  - [x] `updateGlobalProperties(TurmsProperties properties)` -> [internal/domain/gateway/session/identity_access_manager.go:UpdateGlobalProperties(properties interface{})](../internal/domain/gateway/session/identity_access_manager.go)

- **SessionIdentityAccessManager.java** ([java/im/turms/gateway/domain/session/service/SessionIdentityAccessManager.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/domain/session/service/SessionIdentityAccessManager.java))
> [简述功能]

  - [x] `verifyAndGrant(int version, Long userId, @Nullable String password, DeviceType deviceType, @Nullable Map<String, String> deviceDetails, @Nullable UserStatus userStatus, @Nullable Location location, @Nullable String ip)` -> [internal/domain/gateway/session/identity_access_manager.go:VerifyAndGrant(ctx context.Context, loginInfo *bo.UserLoginInfo)](../internal/domain/gateway/session/identity_access_manager.go)

- **SessionService.java** ([java/im/turms/gateway/domain/session/service/SessionService.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/domain/session/service/SessionService.java))
> [简述功能]

  - [x] `destroy()` -> [internal/domain/gateway/session/service.go:Destroy(ctx context.Context)](../internal/domain/gateway/session/service.go)
  - [x] `handleHeartbeatUpdateRequest(UserSession session)` -> [internal/domain/gateway/session/service.go:HandleHeartbeatUpdateRequest(session *UserSession)](../internal/domain/gateway/session/service.go)
  - [x] `handleLoginRequest(int version, @NotNull ByteArrayWrapper ip, @NotNull Long userId, @Nullable String password, @NotNull DeviceType deviceType, @Nullable Map<String, String> deviceDetails, @Nullable UserStatus userStatus, @Nullable Location location, @Nullable String ipStr)` -> [internal/domain/gateway/session/service.go:HandleLoginRequest(ctx context.Context, ...)](../internal/domain/gateway/session/service.go)
  - [x] `closeLocalSessions(@NotNull List<byte[]> ips, @NotNull CloseReason closeReason)` -> [internal/domain/gateway/session/service.go:CloseLocalSessionsByIp(ctx context.Context, ips [][]byte, closeReason any)](../internal/domain/gateway/session/service.go)
  - [x] `closeLocalSessions(@NotNull byte[] ip, @NotNull CloseReason closeReason)` -> [internal/domain/gateway/session/service.go:CloseLocalSessionsByIp(ctx context.Context, ips [][]byte, closeReason any)](../internal/domain/gateway/session/service.go)
  - [x] `closeLocalSession(@NotNull Long userId, @NotNull @ValidDeviceType DeviceType deviceType, @NotNull SessionCloseStatus closeStatus)` -> [internal/domain/gateway/session/service.go:CloseLocalSession(ctx context.Context, userId int64, deviceTypes []protocol.DeviceType, closeReason any)](../internal/domain/gateway/session/service.go)
  - [x] `closeLocalSession(@NotNull Long userId, @NotNull @ValidDeviceType DeviceType deviceType, @NotNull CloseReason closeReason)` -> [internal/domain/gateway/session/service.go:CloseLocalSession(ctx context.Context, userId int64, deviceTypes []protocol.DeviceType, closeReason any)](../internal/domain/gateway/session/service.go)
  - [x] `closeLocalSession(@NotNull Long userId, @NotEmpty Set<@ValidDeviceType DeviceType> deviceTypes, @NotNull CloseReason closeReason)` -> [internal/domain/gateway/session/service.go:CloseLocalSession(ctx context.Context, userId int64, deviceTypes []protocol.DeviceType, closeReason any)](../internal/domain/gateway/session/service.go)
  - [x] `closeLocalSessions(@NotNull Set<Long> userIds, @NotNull CloseReason closeReason)` -> [internal/domain/gateway/session/service.go:CloseLocalSessionsByUserIds(ctx context.Context, userIds []int64, closeReason any)](../internal/domain/gateway/session/service.go)
  - [x] `authAndCloseLocalSession(@NotNull Long userId, @NotNull DeviceType deviceType, @NotNull CloseReason closeReason, int sessionId)` -> [internal/domain/gateway/session/service.go:AuthAndCloseLocalSession(ctx context.Context, userId int64, deviceType protocol.DeviceType, closeReason any, sessionId int)](../internal/domain/gateway/session/service.go)
  - [x] `closeAllLocalSessions(@NotNull CloseReason closeReason)` -> [internal/domain/gateway/session/service.go:CloseAllLocalSessions(ctx context.Context, closeReason any)](../internal/domain/gateway/session/service.go)
  - [x] `closeLocalSession(Long userId, SessionCloseStatus closeStatus)` -> [internal/domain/gateway/session/service.go:CloseLocalSession(ctx context.Context, userId int64, deviceTypes []protocol.DeviceType, closeReason any)](../internal/domain/gateway/session/service.go)
  - [x] `closeLocalSession(Long userId, CloseReason closeReason)` -> [internal/domain/gateway/session/service.go:CloseLocalSession(ctx context.Context, userId int64, deviceTypes []protocol.DeviceType, closeReason any)](../internal/domain/gateway/session/service.go)
  - [x] `getSessions(Set<Long> userIds)` -> [internal/domain/gateway/session/service.go:GetSessions(ctx context.Context, userIds []int64)](../internal/domain/gateway/session/service.go)
  - [x] `authAndUpdateHeartbeatTimestamp(long userId, @NotNull @ValidDeviceType DeviceType deviceType, int sessionId)` -> [internal/domain/gateway/session/service.go:AuthAndUpdateHeartbeatTimestamp(ctx context.Context, userId int64, deviceType protocol.DeviceType, sessionId int)](../internal/domain/gateway/session/service.go)
  - [x] `tryRegisterOnlineUser(int version, @NotNull Set<TurmsRequest.KindCase> permissions, @NotNull ByteArrayWrapper ip, @NotNull Long userId, @NotNull DeviceType deviceType, @Nullable Map<String, String> deviceDetails, @Nullable UserStatus userStatus, @Nullable Location location)` -> [internal/domain/gateway/session/service.go:TryRegisterOnlineUser(ctx context.Context, ...)](../internal/domain/gateway/session/service.go)
  - [x] `getUserSessionsManager(@NotNull Long userId)` -> [internal/domain/gateway/session/service.go:GetUserSessionsManager(ctx context.Context, userId int64)](../internal/domain/gateway/session/service.go)
  - [x] `getLocalUserSession(@NotNull Long userId, @NotNull DeviceType deviceType)` -> [internal/domain/gateway/session/service.go:GetLocalUserSession(ctx context.Context, userId int64, deviceType protocol.DeviceType)](../internal/domain/gateway/session/service.go)
  - [x] `getLocalUserSession(ByteArrayWrapper ip)` -> [internal/domain/gateway/session/service.go:GetLocalUserSession(ctx context.Context, userId int64, deviceType protocol.DeviceType)](../internal/domain/gateway/session/service.go)
  - [x] `countLocalOnlineUsers()` -> [internal/domain/gateway/session/service.go:CountOnlineUsers()](../internal/domain/gateway/session/service.go)
  - [x] `onSessionEstablished(@NotNull UserSessionsManager userSessionsManager, @NotNull @ValidDeviceType DeviceType deviceType)` -> [internal/domain/gateway/session/service.go:OnSessionEstablished(ctx context.Context, userSessionsManager any, deviceType protocol.DeviceType)](../internal/domain/gateway/session/service.go)
  - [x] `addOnSessionClosedListeners(Consumer<UserSession> onSessionClosed)` -> [internal/domain/gateway/session/service.go:AddOnSessionClosedListeners(ctx context.Context, onSessionClosed func(*UserSession))](../internal/domain/gateway/session/service.go)
  - [x] `invokeGoOnlineHandlers(@NotNull UserSessionsManager userSessionsManager, @NotNull UserSession userSession)` -> [internal/domain/gateway/session/service.go:InvokeGoOnlineHandlers(ctx context.Context, userSessionsManager any, userSession *UserSession)](../internal/domain/gateway/session/service.go)

- **UserService.java** ([java/im/turms/gateway/domain/session/service/UserService.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/domain/session/service/UserService.java))
> [简述功能]

  - [ ] `authenticate(@NotNull Long userId, @Nullable String rawPassword)`
  - [x] `isActiveAndNotDeleted(@NotNull Long userId)` -> [internal/domain/user/repository/user_repository.go:IsActiveAndNotDeleted(ctx context.Context, userID int64)](../internal/domain/user/repository/user_repository.go)

- **UserSimultaneousLoginService.java** ([java/im/turms/gateway/domain/session/service/UserSimultaneousLoginService.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/domain/session/service/UserSimultaneousLoginService.java))
> [简述功能]

  - [ ] `getConflictedDeviceTypes(@NotNull @ValidDeviceType DeviceType deviceType)`
  - [ ] `isForbiddenDeviceType(DeviceType deviceType)`
  - [ ] `shouldDisconnectLoggingInDeviceIfConflicts()`

- **ServiceAddressManager.java** ([java/im/turms/gateway/infra/address/ServiceAddressManager.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/address/ServiceAddressManager.java))
> [简述功能]

  - [ ] `getWsAddress()`
  - [ ] `getTcpAddress()`
  - [ ] `getUdpAddress()`

- **LdapClient.java** ([java/im/turms/gateway/infra/ldap/LdapClient.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/ldap/LdapClient.java))
> [简述功能]

  - [x] `isConnected()` -> [internal/domain/gateway/session/connection.go:IsConnected()](../internal/domain/gateway/session/connection.go)
  - [ ] `connect()`
  - [ ] `bind(boolean useFastBind, String dn, String password)`
  - [x] `search(String baseDn, Scope scope, DerefAliases derefAliases, int sizeLimit, int timeLimit, boolean typeOnly, List<String> attributes, String filter)` -> [internal/storage/elasticsearch/elasticsearch_client.go:Search(ctx context.Context)](../internal/storage/elasticsearch/elasticsearch_client.go)
  - [ ] `modify(String dn, List<ModifyOperationChange> changes)`

- **BerBuffer.java** ([java/im/turms/gateway/infra/ldap/asn1/BerBuffer.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/ldap/asn1/BerBuffer.java))
> [简述功能]

  - [ ] `skipTag()`
  - [ ] `skipTagAndLength()`
  - [ ] `skipTagAndLengthAndValue()`
  - [ ] `readTag()`
  - [ ] `peekAndCheckTag(int tag)`
  - [ ] `skipLength()`
  - [ ] `skipLengthAndValue()`
  - [ ] `writeLength(int length)`
  - [ ] `readLength()`
  - [ ] `tryReadLengthIfReadable()`
  - [ ] `beginSequence()`
  - [ ] `beginSequence(int tag)`
  - [ ] `endSequence()`
  - [ ] `writeBoolean(boolean value)`
  - [ ] `writeBoolean(int tag, boolean value)`
  - [ ] `readBoolean()`
  - [ ] `writeInteger(int value)`
  - [ ] `writeInteger(int tag, int value)`
  - [ ] `readInteger()`
  - [ ] `readIntWithTag(int tag)`
  - [ ] `writeOctetString(String value)`
  - [ ] `writeOctetString(byte[] value)`
  - [ ] `writeOctetString(int tag, byte[] value)`
  - [ ] `writeOctetString(byte[] value, int start, int length)`
  - [ ] `writeOctetString(int tag, byte[] value, int start, int length)`
  - [ ] `writeOctetString(int tag, String value)`
  - [ ] `writeOctetStrings(List<String> values)`
  - [ ] `readOctetString()`
  - [ ] `readOctetStringWithTag(int tag)`
  - [ ] `readOctetStringWithLength(int length)`
  - [ ] `writeEnumeration(int value)`
  - [ ] `readEnumeration()`
  - [ ] `getBytes()`
  - [ ] `skipBytes(int length)`
  - [x] `close()` -> [internal/domain/common/cache/ttl_cache.go:Close()](../internal/domain/common/cache/ttl_cache.go)
  - [ ] `refCnt()`
  - [ ] `retain()`
  - [ ] `retain(int increment)`
  - [ ] `touch()`
  - [ ] `touch(Object hint)`
  - [ ] `release()`
  - [ ] `release(int decrement)`
  - [ ] `isReadable(int length)`
  - [ ] `isReadable()`
  - [ ] `isReadableWithEnd(int end)`
  - [ ] `readerIndex()`

- **Attribute.java** ([java/im/turms/gateway/infra/ldap/element/common/Attribute.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/ldap/element/common/Attribute.java))
> [简述功能]

  - [x] `isEmpty()` -> [internal/domain/gateway/session/sharded_map.go:IsEmpty()](../internal/domain/gateway/session/sharded_map.go)
  - [ ] `decode(BerBuffer buffer)`

- **LdapMessage.java** ([java/im/turms/gateway/infra/ldap/element/common/LdapMessage.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/ldap/element/common/LdapMessage.java))
> [简述功能]

  - [ ] `estimateSize()`
  - [ ] `writeTo(BerBuffer buffer)`

- **LdapResult.java** ([java/im/turms/gateway/infra/ldap/element/common/LdapResult.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/ldap/element/common/LdapResult.java))
> [简述功能]

  - [ ] `isSuccess()`

- **Control.java** ([java/im/turms/gateway/infra/ldap/element/common/control/Control.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/ldap/element/common/control/Control.java))
> [简述功能]

  - [ ] `decode(BerBuffer buffer)`

- **BindRequest.java** ([java/im/turms/gateway/infra/ldap/element/operation/bind/BindRequest.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/ldap/element/operation/bind/BindRequest.java))
> [简述功能]

  - [ ] `estimateSize()`
  - [ ] `writeTo(BerBuffer buffer)`

- **BindResponse.java** ([java/im/turms/gateway/infra/ldap/element/operation/bind/BindResponse.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/ldap/element/operation/bind/BindResponse.java))
> [简述功能]

  - [ ] `decode(BerBuffer buffer)`

- **ModifyRequest.java** ([java/im/turms/gateway/infra/ldap/element/operation/modify/ModifyRequest.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/ldap/element/operation/modify/ModifyRequest.java))
> [简述功能]

  - [ ] `estimateSize()`
  - [ ] `writeTo(BerBuffer buffer)`

- **ModifyResponse.java** ([java/im/turms/gateway/infra/ldap/element/operation/modify/ModifyResponse.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/ldap/element/operation/modify/ModifyResponse.java))
> [简述功能]

  - [ ] `decode(BerBuffer buffer)`

- **Filter.java** ([java/im/turms/gateway/infra/ldap/element/operation/search/Filter.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/ldap/element/operation/search/Filter.java))
> [简述功能]

  - [ ] `write(BerBuffer buffer, String filter)`

- **SearchRequest.java** ([java/im/turms/gateway/infra/ldap/element/operation/search/SearchRequest.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/ldap/element/operation/search/SearchRequest.java))
> [简述功能]

  - [ ] `estimateSize()`
  - [ ] `writeTo(BerBuffer buffer)`

- **SearchResult.java** ([java/im/turms/gateway/infra/ldap/element/operation/search/SearchResult.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/ldap/element/operation/search/SearchResult.java))
> [简述功能]

  - [ ] `decode(BerBuffer buffer)`
  - [ ] `isComplete()`

- **ApiLoggingContext.java** ([java/im/turms/gateway/infra/logging/ApiLoggingContext.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/logging/ApiLoggingContext.java))
> [简述功能]

  - [ ] `shouldLogHeartbeatRequest()`
  - [x] `shouldLogRequest(TurmsRequest.KindCase requestType)` -> [internal/infra/logging/api_logging_context.go:ShouldLogRequest(requestType int)](../internal/infra/logging/api_logging_context.go)
  - [x] `shouldLogNotification(TurmsRequest.KindCase requestType)` -> [internal/infra/logging/api_logging_context.go:ShouldLogNotification(requestType int)](../internal/infra/logging/api_logging_context.go)

- **ClientApiLogging.java** ([java/im/turms/gateway/infra/logging/ClientApiLogging.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/logging/ClientApiLogging.java))
> [简述功能]

  - [x] `log(@Nullable Integer sessionId, @Nullable Long userId, @Nullable DeviceType deviceType, @Nullable Integer version, String ip, long requestId, TurmsRequest.KindCase requestType, int requestSize, long requestTime, TurmsNotification response, long processingTime)` -> [internal/infra/logging/client_api_logging.go:Log(request interface{}, serviceRequest interface{}, requestSize int64, requestTime int64, response interface{}, processingTime int64)](../internal/infra/logging/client_api_logging.go)
  - [x] `log(@Nullable Integer sessionId, @Nullable Long userId, @Nullable DeviceType deviceType, @Nullable Integer version, String ip, long requestId, TurmsRequest.KindCase requestType, int requestSize, long requestTime, int responseCode, long processingTime)` -> [internal/infra/logging/client_api_logging.go:Log(request interface{}, serviceRequest interface{}, requestSize int64, requestTime int64, response interface{}, processingTime int64)](../internal/infra/logging/client_api_logging.go)
  - [x] `log(@Nullable Integer sessionId, @Nullable Long userId, @Nullable DeviceType deviceType, @Nullable Integer version, String ip, long requestId, String requestType, int requestSize, long requestTime, int responseCode, @Nullable String responseDataType, int responseSize, long processingTime)` -> [internal/infra/logging/client_api_logging.go:Log(request interface{}, serviceRequest interface{}, requestSize int64, requestTime int64, response interface{}, processingTime int64)](../internal/infra/logging/client_api_logging.go)

- **NotificationLoggingManager.java** ([java/im/turms/gateway/infra/logging/NotificationLoggingManager.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/logging/NotificationLoggingManager.java))
> [简述功能]

  - [x] `log(SimpleTurmsNotification notification, int notificationBytes, int recipientCount, int onlineRecipientCount)` -> [internal/infra/logging/client_api_logging.go:Log(request interface{}, serviceRequest interface{}, requestSize int64, requestTime int64, response interface{}, processingTime int64)](../internal/infra/logging/client_api_logging.go)

- **SimpleTurmsNotification.java** ([java/im/turms/gateway/infra/proto/SimpleTurmsNotification.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/proto/SimpleTurmsNotification.java))
> [简述功能]

  - [ ] `SimpleTurmsNotification(long requesterId, Integer closeStatus, TurmsRequest.KindCase relayedRequestType)`

- **SimpleTurmsRequest.java** ([java/im/turms/gateway/infra/proto/SimpleTurmsRequest.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/proto/SimpleTurmsRequest.java))
> [简述功能]

  - [ ] `SimpleTurmsRequest(long requestId, TurmsRequest.KindCase type, CreateSessionRequest createSessionRequest)`
  - [x] `toString()` -> [internal/domain/gateway/session/connection.go:ToString()](../internal/domain/gateway/session/connection.go)

- **TurmsNotificationParser.java** ([java/im/turms/gateway/infra/proto/TurmsNotificationParser.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/proto/TurmsNotificationParser.java))
> [简述功能]

  - [ ] `parseSimpleNotification(CodedInputStream turmsRequestInputStream)`

- **TurmsRequestParser.java** ([java/im/turms/gateway/infra/proto/TurmsRequestParser.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/proto/TurmsRequestParser.java))
> [简述功能]

  - [ ] `parseSimpleRequest(CodedInputStream turmsRequestInputStream)`

- **MongoConfig.java** ([java/im/turms/gateway/storage/mongo/MongoConfig.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/storage/mongo/MongoConfig.java))
> [简述功能]

  - [x] `adminMongoClient(TurmsPropertiesManager propertiesManager)` -> [internal/storage/mongo/mongo_config.go:AdminMongoClient()](../internal/storage/mongo/mongo_config.go)
  - [x] `userMongoClient(TurmsPropertiesManager propertiesManager)` -> [internal/storage/mongo/mongo_config.go:UserMongoClient()](../internal/storage/mongo/mongo_config.go)
  - [ ] `mongoDataGenerator()`

### turms-service

> [简述功能]

#### Configurations

- **application-demo.yaml** ([resources/application-demo.yaml](../turms-orig/turms-service/src/main/resources/application-demo.yaml)): [简述功能]
- **application-dev.yaml** ([resources/application-dev.yaml](../turms-orig/turms-service/src/main/resources/application-dev.yaml)): [简述功能]
- **application-test.yaml** ([resources/application-test.yaml](../turms-orig/turms-service/src/main/resources/application-test.yaml)): [简述功能]
- **application.yaml** ([resources/application.yaml](../turms-orig/turms-service/src/main/resources/application.yaml)): [简述功能]

#### Java source tracking

- **TurmsServiceApplication.java** ([java/im/turms/service/TurmsServiceApplication.java](../turms-orig/turms-service/src/main/java/im/turms/service/TurmsServiceApplication.java))
> [简述功能]

  - [x] `main(String[] args)` -> [cmd/turms-gateway/main.go:main()](../cmd/turms-gateway/main.go)

- **ServiceRequestDispatcher.java** ([java/im/turms/service/access/servicerequest/dispatcher/ServiceRequestDispatcher.java](../turms-orig/turms-service/src/main/java/im/turms/service/access/servicerequest/dispatcher/ServiceRequestDispatcher.java))
> [简述功能]

  - [x] `dispatch(TracingContext context, ServiceRequest serviceRequest)` -> [internal/domain/common/infra/cluster/rpc/router.go:Dispatch(ctx context.Context, frame *codec.RpcFrame)](../internal/domain/common/infra/cluster/rpc/router.go)

- **ClientRequest.java** ([java/im/turms/service/access/servicerequest/dto/ClientRequest.java](../turms-orig/turms-service/src/main/java/im/turms/service/access/servicerequest/dto/ClientRequest.java))
> [简述功能]

  - [x] `toString()` -> [internal/domain/gateway/session/connection.go:ToString()](../internal/domain/gateway/session/connection.go)
  - [ ] `turmsRequest()`
  - [ ] `userId()`
  - [ ] `deviceType()`
  - [ ] `clientIp()`
  - [ ] `requestId()`
  - [ ] `equals(Object obj)`
  - [ ] `hashCode()`

- **RequestHandlerResult.java** ([java/im/turms/service/access/servicerequest/dto/RequestHandlerResult.java](../turms-orig/turms-service/src/main/java/im/turms/service/access/servicerequest/dto/RequestHandlerResult.java))
> [简述功能]

  - [ ] `RequestHandlerResult(ResponseStatusCode code, @Nullable String reason, @Nullable TurmsNotification.Data response, List<Notification> notifications)`
  - [x] `toString()` -> [internal/domain/gateway/session/connection.go:ToString()](../internal/domain/gateway/session/connection.go)
  - [ ] `of(@NotNull ResponseStatusCode code)`
  - [ ] `of(@NotNull ResponseStatusCode code, @Nullable String reason)`
  - [ ] `of(@NotNull TurmsNotification.Data response)`
  - [ ] `of(boolean forwardNotificationsToRequesterOtherOnlineSessions, @NotNull TurmsRequest notification)`
  - [ ] `of(boolean forwardNotificationsToRequesterOtherOnlineSessions, @NotNull Long recipientId, @NotNull TurmsRequest notification)`
  - [ ] `of(@NotNull Long recipientId, @NotNull TurmsRequest notification)`
  - [ ] `of(@NotEmpty Set<Long> recipientIds, @NotNull TurmsRequest dataForRecipient)`
  - [ ] `of(boolean forwardNotificationToRequesterOtherOnlineSessions, @NotEmpty Set<Long> recipientIds, @NotNull TurmsRequest notification)`
  - [ ] `of(TurmsNotification.Data response, @NotEmpty Set<Long> recipientIds, @NotNull TurmsRequest notificationForRecipients, @NotNull TurmsRequest notificationForRequesterOtherOnlineSessions)`
  - [ ] `of(@NotEmpty Set<Long> recipientIds, @NotNull TurmsRequest notificationForRecipients, @NotNull TurmsRequest notificationForRequesterOtherOnlineSessions)`
  - [ ] `of(TurmsNotification.Data response, boolean forwardNotificationToRequesterOtherOnlineSessions, @NotEmpty Set<Long> recipientIds, @NotNull TurmsRequest notification)`
  - [ ] `of(@NotNull ResponseStatusCode code, @NotNull Long recipientId, @NotNull TurmsRequest notification)`
  - [ ] `of(@NotNull ResponseStatusCode code, boolean forwardNotificationToRequesterOtherOnlineSessions, @NotNull Long recipientId, @NotNull TurmsRequest notification)`
  - [ ] `of(@NotNull List<Notification> notifications)`
  - [ ] `of(@NotNull Notification notification)`
  - [ ] `ofDataLong(@NotNull Long value)`
  - [ ] `ofDataLong(@NotNull Long value, @NotNull Long recipientId, @NotNull TurmsRequest notification)`
  - [ ] `ofDataLong(@NotNull Long value, boolean forwardNotificationToRequesterOtherOnlineSessions, @NotNull Long recipientId, @NotNull TurmsRequest notification)`
  - [ ] `ofDataLong(@NotNull Long value, boolean forwardDataForRecipientsToRequesterOtherOnlineSessions, @NotNull TurmsRequest notification)`
  - [ ] `ofDataLong(@NotNull Long value, boolean forwardNotificationsToRequesterOtherOnlineSessions, @NotEmpty Set<Long> recipients, TurmsRequest notification)`
  - [ ] `ofDataLongs(@NotNull Collection<Long> values)`
  - [ ] `Notification(boolean forwardToRequesterOtherOnlineSessions, Set<Long> recipients, TurmsRequest notification)`
  - [ ] `of(boolean forwardToRequesterOtherOnlineSessions, Set<Long> recipients, TurmsRequest notification)`
  - [ ] `of(boolean forwardToRequesterOtherOnlineSessions, Long recipient, TurmsRequest notification)`
  - [ ] `of(boolean forwardToRequesterOtherOnlineSessions, TurmsRequest notification)`
  - [x] `toString()` -> [internal/domain/gateway/session/connection.go:ToString()](../internal/domain/gateway/session/connection.go)

- **AdminController.java** ([java/im/turms/service/domain/admin/access/admin/controller/AdminController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/admin/access/admin/controller/AdminController.java))
> [简述功能]

  - [ ] `checkLoginNameAndPassword()`
  - [ ] `addAdmin(RequestContext requestContext, @RequestBody AddAdminDTO addAdminDTO)`
  - [ ] `queryAdmins(@QueryParam(required = false)`
  - [ ] `queryAdmins(@QueryParam(required = false)`
  - [ ] `updateAdmins(RequestContext requestContext, Set<Long> ids, @RequestBody UpdateAdminDTO updateAdminDTO)`
  - [ ] `deleteAdmins(RequestContext requestContext, Set<Long> ids)`

- **AdminPermissionController.java** ([java/im/turms/service/domain/admin/access/admin/controller/AdminPermissionController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/admin/access/admin/controller/AdminPermissionController.java))
> [简述功能]

  - [ ] `queryAdminPermissions()`

- **AdminRoleController.java** ([java/im/turms/service/domain/admin/access/admin/controller/AdminRoleController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/admin/access/admin/controller/AdminRoleController.java))
> [简述功能]

  - [ ] `addAdminRole(RequestContext requestContext, @RequestBody AddAdminRoleDTO addAdminRoleDTO)`
  - [ ] `queryAdminRoles(@QueryParam(required = false)`
  - [ ] `queryAdminRoles(@QueryParam(required = false)`
  - [ ] `updateAdminRole(RequestContext requestContext, Set<Long> ids, @RequestBody UpdateAdminRoleDTO updateAdminRoleDTO)`
  - [ ] `deleteAdminRoles(RequestContext requestContext, Set<Long> ids)`

- **AddAdminDTO.java** ([java/im/turms/service/domain/admin/access/admin/dto/request/AddAdminDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/admin/access/admin/dto/request/AddAdminDTO.java))
> [简述功能]

  - [ ] `AddAdminDTO(String loginName, @SensitiveProperty(SensitiveProperty.Access.ALLOW_DESERIALIZATION)`
  - [x] `toString()` -> [internal/domain/gateway/session/connection.go:ToString()](../internal/domain/gateway/session/connection.go)

- **AddAdminRoleDTO.java** ([java/im/turms/service/domain/admin/access/admin/dto/request/AddAdminRoleDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/admin/access/admin/dto/request/AddAdminRoleDTO.java))
> [简述功能]

  - [ ] `AddAdminRoleDTO(Long id, String name, Set<String> permissions, Integer rank)`

- **UpdateAdminDTO.java** ([java/im/turms/service/domain/admin/access/admin/dto/request/UpdateAdminDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/admin/access/admin/dto/request/UpdateAdminDTO.java))
> [简述功能]

  - [ ] `UpdateAdminDTO(@SensitiveProperty(SensitiveProperty.Access.ALLOW_DESERIALIZATION)`
  - [x] `toString()` -> [internal/domain/gateway/session/connection.go:ToString()](../internal/domain/gateway/session/connection.go)

- **UpdateAdminRoleDTO.java** ([java/im/turms/service/domain/admin/access/admin/dto/request/UpdateAdminRoleDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/admin/access/admin/dto/request/UpdateAdminRoleDTO.java))
> [简述功能]

  - [ ] `UpdateAdminRoleDTO(String name, Set<String> permissions, Integer rank)`

- **PermissionDTO.java** ([java/im/turms/service/domain/admin/access/admin/dto/response/PermissionDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/admin/access/admin/dto/response/PermissionDTO.java))
> [简述功能]

  - [ ] `PermissionDTO(String group, AdminPermission permission)`

- **AdminRepository.java** ([java/im/turms/service/domain/admin/repository/AdminRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/admin/repository/AdminRepository.java))
> [简述功能]

  - [ ] `updateAdmins(Set<Long> ids, @Nullable byte[] password, @Nullable String displayName, @Nullable Set<Long> roleIds)`
  - [ ] `countAdmins(@Nullable Set<Long> ids, @Nullable Set<Long> roleIds)`
  - [ ] `findAdmins(@Nullable Collection<Long> ids, @Nullable Collection<String> loginNames, @Nullable Collection<Long> roleIds, @Nullable Integer page, @Nullable Integer size)`

- **AdminRoleRepository.java** ([java/im/turms/service/domain/admin/repository/AdminRoleRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/admin/repository/AdminRoleRepository.java))
> [简述功能]

  - [ ] `updateAdminRoles(Set<Long> roleIds, String newName, @Nullable Set<AdminPermission> permissions, @Nullable Integer rank)`
  - [ ] `countAdminRoles(@Nullable Set<Long> ids, @Nullable Set<String> names, @Nullable Set<AdminPermission> includedPermissions, @Nullable Set<Integer> ranks)`
  - [ ] `findAdminRoles(@Nullable Set<Long> roleIds, @Nullable Set<String> names, @Nullable Set<AdminPermission> includedPermissions, @Nullable Set<Integer> ranks, @Nullable Integer page, @Nullable Integer size)`
  - [ ] `findAdminRolesByIdsAndRankGreaterThan(@NotNull Collection<Long> roleIds, @Nullable Integer rankGreaterThan)`
  - [ ] `findHighestRankByRoleIds(Set<Long> roleIds)`

- **AdminRoleService.java** ([java/im/turms/service/domain/admin/service/AdminRoleService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/admin/service/AdminRoleService.java))
> [简述功能]

  - [ ] `authAndAddAdminRole(@NotNull Long requesterId, @NotNull Long roleId, @NotNull @NoWhitespace @Size( min = MIN_ROLE_NAME_LIMIT, max = MAX_ROLE_NAME_LIMIT)`
  - [ ] `addAdminRole(@NotNull Long roleId, @NotNull @NoWhitespace @Size( min = MIN_ROLE_NAME_LIMIT, max = MAX_ROLE_NAME_LIMIT)`
  - [ ] `authAndDeleteAdminRoles(@NotNull Long requesterId, @NotEmpty Set<Long> roleIds)`
  - [ ] `deleteAdminRoles(@NotEmpty Set<Long> roleIds)`
  - [ ] `authAndUpdateAdminRoles(@NotNull Long requesterId, @NotEmpty Set<Long> roleIds, @Nullable @NoWhitespace @Size( min = MIN_ROLE_NAME_LIMIT, max = MAX_ROLE_NAME_LIMIT)`
  - [ ] `updateAdminRole(@NotEmpty Set<Long> roleIds, @Nullable @NoWhitespace @Size( min = MIN_ROLE_NAME_LIMIT, max = MAX_ROLE_NAME_LIMIT)`
  - [ ] `queryAdminRoles(@Nullable Set<Long> ids, @Nullable Set<String> names, @Nullable Set<AdminPermission> includedPermissions, @Nullable Set<Integer> ranks, @Nullable Integer page, @Nullable Integer size)`
  - [ ] `queryAndCacheRolesByRoleIdsAndRankGreaterThan(@NotNull Collection<Long> roleIds, @NotNull Integer rankGreaterThan)`
  - [ ] `countAdminRoles(@Nullable Set<Long> ids, @Nullable Set<String> names, @Nullable Set<AdminPermission> includedPermissions, @Nullable Set<Integer> ranks)`
  - [ ] `queryHighestRankByAdminId(@NotNull Long adminId)`
  - [ ] `queryHighestRankByRoleIds(@NotNull Set<Long> roleIds)`
  - [ ] `isAdminRankHigherThanRank(@NotNull Long adminId, @NotNull Integer rank)`
  - [ ] `queryPermissions(@NotNull Long adminId)`

- **AdminService.java** ([java/im/turms/service/domain/admin/service/AdminService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/admin/service/AdminService.java))
> [简述功能]

  - [ ] `queryRoleIdsByAdminIds(@NotEmpty Set<Long> adminIds)`
  - [ ] `authAndAddAdmin(@NotNull Long requesterId, @Nullable @NoWhitespace @Size( min = MIN_LOGIN_NAME_LIMIT, max = MAX_LOGIN_NAME_LIMIT)`
  - [ ] `addAdmin(@Nullable Long id, @Nullable @NoWhitespace @Size( min = MIN_LOGIN_NAME_LIMIT, max = MAX_LOGIN_NAME_LIMIT)`
  - [ ] `queryAdmins(@Nullable Collection<Long> ids, @Nullable Collection<String> loginNames, @Nullable Collection<Long> roleIds, @Nullable Integer page, @Nullable Integer size)`
  - [ ] `authAndDeleteAdmins(@NotNull Long requesterId, @NotEmpty Set<Long> adminIds)`
  - [ ] `authAndUpdateAdmins(@NotNull Long requesterId, @NotEmpty Set<Long> targetAdminIds, @Nullable @NoWhitespace @Size( min = MIN_PASSWORD_LIMIT, max = MAX_PASSWORD_LIMIT)`
  - [ ] `updateAdmins(@NotEmpty Set<Long> targetAdminIds, @Nullable @NoWhitespace @Size( min = MIN_PASSWORD_LIMIT, max = MAX_PASSWORD_LIMIT)`
  - [ ] `countAdmins(@Nullable Set<Long> ids, @Nullable Set<Long> roleIds)`
  - [ ] `errorRequesterNotExist()`

- **IpBlocklistController.java** ([java/im/turms/service/domain/blocklist/access/admin/controller/IpBlocklistController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/blocklist/access/admin/controller/IpBlocklistController.java))
> [简述功能]

  - [ ] `addBlockedIps(@RequestBody AddBlockedIpsDTO addBlockedIpsDTO)`
  - [ ] `queryBlockedIps(Set<String> ids)`
  - [ ] `queryBlockedIps(int page, @QueryParam(required = false)`
  - [ ] `deleteBlockedIps(@QueryParam(required = false)`

- **UserBlocklistController.java** ([java/im/turms/service/domain/blocklist/access/admin/controller/UserBlocklistController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/blocklist/access/admin/controller/UserBlocklistController.java))
> [简述功能]

  - [ ] `addBlockedUserIds(@RequestBody AddBlockedUserIdsDTO addBlockedUserIdsDTO)`
  - [x] `queryBlockedUsers(Set<Long> ids)` -> [internal/domain/group/service/group_blocklist_service.go:QueryBlockedUsers(ctx context.Context, groupID int64)](../internal/domain/group/service/group_blocklist_service.go)
  - [x] `queryBlockedUsers(int page, @QueryParam(required = false)` -> [internal/domain/group/service/group_blocklist_service.go:QueryBlockedUsers(ctx context.Context, groupID int64)](../internal/domain/group/service/group_blocklist_service.go)
  - [ ] `deleteBlockedUserIds(@QueryParam(required = false)`

- **AddBlockedIpsDTO.java** ([java/im/turms/service/domain/blocklist/access/admin/dto/request/AddBlockedIpsDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/blocklist/access/admin/dto/request/AddBlockedIpsDTO.java))
> [简述功能]

  - [ ] `AddBlockedIpsDTO(Set<String> ids, long blockDurationMillis)`

- **AddBlockedUserIdsDTO.java** ([java/im/turms/service/domain/blocklist/access/admin/dto/request/AddBlockedUserIdsDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/blocklist/access/admin/dto/request/AddBlockedUserIdsDTO.java))
> [简述功能]

  - [ ] `AddBlockedUserIdsDTO(Set<Long> ids, long blockDurationMillis)`

- **BlockedIpDTO.java** ([java/im/turms/service/domain/blocklist/access/admin/dto/response/BlockedIpDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/blocklist/access/admin/dto/response/BlockedIpDTO.java))
> [简述功能]

  - [ ] `BlockedIpDTO(String id, Date blockEndTime)`

- **BlockedUserDTO.java** ([java/im/turms/service/domain/blocklist/access/admin/dto/response/BlockedUserDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/blocklist/access/admin/dto/response/BlockedUserDTO.java))
> [简述功能]

  - [ ] `BlockedUserDTO(Long id, Date blockEndTime)`

- **BlockedClientSerializer.java** ([java/im/turms/service/domain/blocklist/codec/BlockedClientSerializer.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/blocklist/codec/BlockedClientSerializer.java))
> [简述功能]

  - [x] `serialize(BlockedClient value, JsonGenerator gen, SerializerProvider provider)` -> [internal/storage/elasticsearch/model/elasticsearch_model.go:Serialize()](../internal/storage/elasticsearch/model/elasticsearch_model.go)

- **MemberController.java** ([java/im/turms/service/domain/cluster/access/admin/controller/MemberController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/cluster/access/admin/controller/MemberController.java))
> [简述功能]

  - [ ] `queryMembers()`
  - [ ] `removeMembers(List<String> ids)`
  - [ ] `addMember(@RequestBody AddMemberDTO addMemberDTO)`
  - [ ] `updateMember(String id, @RequestBody UpdateMemberDTO updateMemberDTO)`
  - [ ] `queryLeader()`
  - [ ] `electNewLeader(@QueryParam(required = false)`

- **SettingController.java** ([java/im/turms/service/domain/cluster/access/admin/controller/SettingController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/cluster/access/admin/controller/SettingController.java))
> [简述功能]

  - [ ] `queryClusterSettings(boolean queryLocalSettings, boolean onlyMutable)`
  - [ ] `updateClusterSettings(boolean reset, boolean updateLocalSettings, @RequestBody(required = false)`
  - [ ] `queryClusterConfigMetadata(boolean queryLocalSettings, boolean onlyMutable, boolean withValue)`

- **AddMemberDTO.java** ([java/im/turms/service/domain/cluster/access/admin/dto/request/AddMemberDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/cluster/access/admin/dto/request/AddMemberDTO.java))
> [简述功能]

  - [ ] `AddMemberDTO(String nodeId, String zone, String name, NodeType nodeType, String version, boolean isSeed, boolean isLeaderEligible, Date registrationDate, int priority, String memberHost, int memberPort, String adminApiAddress, String wsAddress, String tcpAddress, String udpAddress, boolean isActive, boolean isHealthy)`

- **UpdateMemberDTO.java** ([java/im/turms/service/domain/cluster/access/admin/dto/request/UpdateMemberDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/cluster/access/admin/dto/request/UpdateMemberDTO.java))
> [简述功能]

  - [ ] `UpdateMemberDTO(String zone, String name, Boolean isSeed, Boolean isLeaderEligible, Boolean isActive, Integer priority)`

- **SettingsDTO.java** ([java/im/turms/service/domain/cluster/access/admin/dto/response/SettingsDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/cluster/access/admin/dto/response/SettingsDTO.java))
> [简述功能]

  - [ ] `SettingsDTO(int schemaVersion, Map<String, Object> settings)`

- **BaseController.java** ([java/im/turms/service/domain/common/access/admin/controller/BaseController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/common/access/admin/controller/BaseController.java))
> [简述功能]

  - [ ] `getPageSize(@Nullable Integer size)`
  - [ ] `queryBetweenDate(DateRange dateRange, DivideBy divideBy, Function3<DateRange, Boolean, Boolean, Mono<Long>> function, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages)`
  - [ ] `queryBetweenDate(DateRange dateRange, DivideBy divideBy, Function<DateRange, Mono<Long>> function)`
  - [ ] `checkAndQueryBetweenDate(DateRange dateRange, DivideBy divideBy, Function3<DateRange, Boolean, Boolean, Mono<Long>> function, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages)`
  - [ ] `checkAndQueryBetweenDate(DateRange dateRange, DivideBy divideBy, Function<DateRange, Mono<Long>> function)`

- **StatisticsRecordDTO.java** ([java/im/turms/service/domain/common/access/admin/dto/response/StatisticsRecordDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/common/access/admin/dto/response/StatisticsRecordDTO.java))
> [简述功能]

  - [ ] `StatisticsRecordDTO(Date date, Long total)`

- **ServicePermission.java** ([java/im/turms/service/domain/common/permission/ServicePermission.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/common/permission/ServicePermission.java))
> [简述功能]

  - [ ] `ServicePermission(ResponseStatusCode code, String reason)`
  - [x] `get(ResponseStatusCode code)` -> [internal/domain/common/cache/sharded_map.go:Get(key K)](../internal/domain/common/cache/sharded_map.go)
  - [x] `get(ResponseStatusCode code, String reason)` -> [internal/domain/common/cache/sharded_map.go:Get(key K)](../internal/domain/common/cache/sharded_map.go)

- **ExpirableEntityRepository.java** ([java/im/turms/service/domain/common/repository/ExpirableEntityRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/common/repository/ExpirableEntityRepository.java))
> [简述功能]

  - [ ] `isExpired(long creationDate)`
  - [ ] `getEntityExpirationDate()`
  - [x] `deleteExpiredData(String creationDateFieldName, Date expirationDate)` -> [internal/domain/user/repository/user_friend_request_repository.go:DeleteExpiredData(ctx context.Context, expirationDate time.Time)](../internal/domain/user/repository/user_friend_request_repository.go)
  - [x] `findMany(Filter filter)` -> [internal/domain/user/repository/user_repository.go:FindMany(ctx context.Context, filter bson.M)](../internal/domain/user/repository/user_repository.go)
  - [x] `findMany(Filter filter, QueryOptions options)` -> [internal/domain/user/repository/user_repository.go:FindMany(ctx context.Context, filter bson.M)](../internal/domain/user/repository/user_repository.go)

- **ExpirableEntityService.java** ([java/im/turms/service/domain/common/service/ExpirableEntityService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/common/service/ExpirableEntityService.java))
> [简述功能]

  - [ ] `getEntityExpirationDate()`

- **UserDefinedAttributesService.java** ([java/im/turms/service/domain/common/service/UserDefinedAttributesService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/common/service/UserDefinedAttributesService.java))
> [简述功能]

  - [ ] `updateGlobalProperties(UserDefinedAttributesProperties properties)`
  - [ ] `parseAttributesForUpsert(Map<String, Value> userDefinedAttributes)`

- **ExpirableRequestInspector.java** ([java/im/turms/service/domain/common/util/ExpirableRequestInspector.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/common/util/ExpirableRequestInspector.java))
> [简述功能]

  - [ ] `isProcessedByResponder(@Nullable RequestStatus status)`

- **DataValidator.java** ([java/im/turms/service/domain/common/validation/DataValidator.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/common/validation/DataValidator.java))
> [简述功能]

  - [x] `validRequestStatus(RequestStatus status)` -> [internal/infra/validator/validator.go:ValidRequestStatus(status interface{}, name string)](../internal/infra/validator/validator.go)
  - [ ] `validResponseAction(ResponseAction action)`
  - [ ] `validDeviceType(DeviceType deviceType)`
  - [ ] `validProfileAccess(ProfileAccessStrategy value)`
  - [ ] `validRelationshipKey(UserRelationship.Key key)`
  - [ ] `validRelationshipGroupKey(UserRelationshipGroup.Key key)`
  - [ ] `validGroupMemberKey(GroupMember.Key key)`
  - [ ] `validGroupMemberRole(GroupMemberRole role)`
  - [ ] `validGroupBlockedUserKey(GroupBlockedUser.Key key)`
  - [ ] `validNewGroupQuestion(NewGroupQuestion question)`
  - [ ] `validGroupQuestionIdAndAnswer(Map.Entry<Long, String> questionIdAndAnswer)`

- **CancelMeetingResult.java** ([java/im/turms/service/domain/conference/bo/CancelMeetingResult.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/conference/bo/CancelMeetingResult.java))
> [简述功能]

  - [ ] `CancelMeetingResult(boolean success, @Nullable Meeting meeting)`

- **UpdateMeetingInvitationResult.java** ([java/im/turms/service/domain/conference/bo/UpdateMeetingInvitationResult.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/conference/bo/UpdateMeetingInvitationResult.java))
> [简述功能]

  - [ ] `UpdateMeetingInvitationResult(boolean updated, @Nullable String accessToken, @Nullable Meeting meeting)`

- **UpdateMeetingResult.java** ([java/im/turms/service/domain/conference/bo/UpdateMeetingResult.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/conference/bo/UpdateMeetingResult.java))
> [简述功能]

  - [ ] `UpdateMeetingResult(boolean success, @Nullable Meeting meeting)`

- **ConferenceServiceController.java** ([java/im/turms/service/domain/conference/controller/ConferenceServiceController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/conference/controller/ConferenceServiceController.java))
> [简述功能]

  - [ ] `handleCreateMeetingRequest()`
  - [ ] `handleDeleteMeetingRequest()`
  - [ ] `handleUpdateMeetingRequest()`
  - [ ] `handleQueryMeetingsRequest()`
  - [ ] `handleUpdateMeetingInvitationRequest()`

- **MeetingRepository.java** ([java/im/turms/service/domain/conference/repository/MeetingRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/conference/repository/MeetingRepository.java))
> [简述功能]

  - [ ] `updateEndDate(Long meetingId, Date endDate)`
  - [ ] `updateCancelDateIfNotCanceled(Long meetingId, Date cancelDate)`
  - [ ] `updateMeeting(Long meetingId, @Nullable String name, @Nullable String intro, @Nullable String password)`
  - [ ] `find(@Nullable Collection<Long> ids, @Nullable Collection<Long> creatorIds, @Nullable Collection<Long> userIds, @Nullable Collection<Long> groupIds, @Nullable Date creationDateStart, @Nullable Date creationDateEnd, @Nullable Integer skip, @Nullable Integer limit)`
  - [ ] `find(@Nullable Collection<Long> ids, @NotNull Long creatorId, @NotNull Long userId, @Nullable Date creationDateStart, @Nullable Date creationDateEnd, @Nullable Integer skip, @Nullable Integer limit)`

- **ConferenceService.java** ([java/im/turms/service/domain/conference/service/ConferenceService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/conference/service/ConferenceService.java))
> [简述功能]

  - [ ] `onExtensionStarted(ConferenceServiceProvider extension)`
  - [ ] `authAndCancelMeeting(@NotNull Long requesterId, @NotNull Long meetingId)`
  - [ ] `queryMeetingParticipants(@Nullable Long userId, @Nullable Long groupId)`
  - [ ] `authAndUpdateMeeting(@NotNull Long requesterId, @NotNull Long meetingId, @Nullable String name, @Nullable String intro, @Nullable String password)`
  - [ ] `authAndUpdateMeetingInvitation(@NotNull Long requesterId, @NotNull Long meetingId, @Nullable String password, @NotNull ResponseAction responseAction)`
  - [ ] `authAndQueryMeetings(@NotNull Long requesterId, @Nullable Set<Long> ids, @Nullable Set<Long> creatorIds, @Nullable Set<Long> userIds, @Nullable Set<Long> groupIds, @Nullable Date creationDateStart, @Nullable Date creationDateEnd, @Nullable Integer skip, @Nullable Integer limit)`

- **ConversationController.java** ([java/im/turms/service/domain/conversation/access/admin/controller/ConversationController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/conversation/access/admin/controller/ConversationController.java))
> [简述功能]

  - [ ] `queryConversations(@QueryParam(required = false)`
  - [ ] `deleteConversations(@QueryParam(required = false)`
  - [ ] `updateConversations(@QueryParam(required = false)`

- **UpdateConversationDTO.java** ([java/im/turms/service/domain/conversation/access/admin/dto/request/UpdateConversationDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/conversation/access/admin/dto/request/UpdateConversationDTO.java))
> [简述功能]

  - [ ] `UpdateConversationDTO(Date readDate)`

- **ConversationsDTO.java** ([java/im/turms/service/domain/conversation/access/admin/dto/response/ConversationsDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/conversation/access/admin/dto/response/ConversationsDTO.java))
> [简述功能]

  - [ ] `ConversationsDTO(List<PrivateConversation> privateConversations, List<GroupConversation> groupConversations)`

- **ConversationServiceController.java** ([java/im/turms/service/domain/conversation/access/servicerequest/controller/ConversationServiceController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/conversation/access/servicerequest/controller/ConversationServiceController.java))
> [简述功能]

  - [ ] `handleQueryConversationsRequest()`
  - [ ] `handleUpdateTypingStatusRequest()`
  - [ ] `handleUpdateConversationRequest()`

- **ConversationSettingsServiceController.java** ([java/im/turms/service/domain/conversation/access/servicerequest/controller/ConversationSettingsServiceController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/conversation/access/servicerequest/controller/ConversationSettingsServiceController.java))
> [简述功能]

  - [ ] `handleUpdateConversationSettingsRequest()`
  - [ ] `handleDeleteConversationSettingsRequest()`
  - [ ] `handleQueryConversationSettingsRequest()`

- **ConversationSettingsRepository.java** ([java/im/turms/service/domain/conversation/repository/ConversationSettingsRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/conversation/repository/ConversationSettingsRepository.java))
> [简述功能]

  - [x] `upsertSettings(Long ownerId, Long targetId, Map<String, Object> settings)` -> [internal/domain/user/service/user_settings_service.go:UpsertSettings(ctx context.Context, userID int64, settings map[string]interface{})](../internal/domain/user/service/user_settings_service.go)
  - [x] `unsetSettings(Long ownerId, @Nullable Collection<Long> targetIds, @Nullable Collection<String> settingNames)` -> [internal/domain/user/service/user_settings_service.go:UnsetSettings(ctx context.Context, userID int64, keys []string)](../internal/domain/user/service/user_settings_service.go)
  - [ ] `findByIdAndSettingNames(Long ownerId, @Nullable Collection<String> settingNames, @Nullable Date lastUpdatedDateStart)`
  - [ ] `findByIdAndSettingNames(Collection<ConversationSettings.Key> keys, @Nullable Collection<String> settingNames, @Nullable Date lastUpdatedDateStart)`
  - [ ] `findSettingFields(Long ownerId, Long targetId, Collection<String> includedFields)`
  - [ ] `deleteByOwnerIds(Collection<Long> ownerIds, @Nullable ClientSession clientSession)`

- **GroupConversationRepository.java** ([java/im/turms/service/domain/conversation/repository/GroupConversationRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/conversation/repository/GroupConversationRepository.java))
> [简述功能]

  - [x] `upsert(Long groupId, Long memberId, Date readDate, boolean allowMoveReadDateForward)` -> [internal/domain/group/repository/group_version_repository.go:Upsert(ctx context.Context, groupID int64, timestamp time.Time)](../internal/domain/group/repository/group_version_repository.go)
  - [x] `upsert(Long groupId, Collection<Long> memberIds, Date readDate)` -> [internal/domain/group/repository/group_version_repository.go:Upsert(ctx context.Context, groupID int64, timestamp time.Time)](../internal/domain/group/repository/group_version_repository.go)
  - [ ] `deleteMemberConversations(Collection<Long> groupIds, Long memberId, ClientSession session)`

- **PrivateConversationRepository.java** ([java/im/turms/service/domain/conversation/repository/PrivateConversationRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/conversation/repository/PrivateConversationRepository.java))
> [简述功能]

  - [x] `upsert(Set<PrivateConversation.Key> keys, Date readDate, boolean allowMoveReadDateForward)` -> [internal/domain/group/repository/group_version_repository.go:Upsert(ctx context.Context, groupID int64, timestamp time.Time)](../internal/domain/group/repository/group_version_repository.go)
  - [ ] `deleteConversationsByOwnerIds(Set<Long> ownerIds, @Nullable ClientSession session)`
  - [ ] `findConversations(Collection<Long> ownerIds)`

- **ConversationService.java** ([java/im/turms/service/domain/conversation/service/ConversationService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/conversation/service/ConversationService.java))
> [简述功能]

  - [ ] `authAndUpsertGroupConversationReadDate(@NotNull Long groupId, @NotNull Long memberId, @Nullable @PastOrPresent Date readDate)`
  - [ ] `authAndUpsertPrivateConversationReadDate(@NotNull Long ownerId, @NotNull Long targetId, @Nullable @PastOrPresent Date readDate)`
  - [ ] `upsertGroupConversationReadDate(@NotNull Long groupId, @NotNull Long memberId, @Nullable @PastOrPresent Date readDate)`
  - [ ] `upsertGroupConversationsReadDate(@NotNull Set<GroupConversation.GroupConversionMemberKey> keys, @Nullable @PastOrPresent Date readDate)`
  - [ ] `upsertPrivateConversationReadDate(@NotNull Long ownerId, @NotNull Long targetId, @Nullable @PastOrPresent Date readDate)`
  - [ ] `upsertPrivateConversationsReadDate(@NotNull Set<PrivateConversation.Key> keys, @Nullable @PastOrPresent Date readDate)`
  - [x] `queryGroupConversations(@NotNull Collection<Long> groupIds)` -> [internal/domain/conversation/repository/group_conversation_repository.go:QueryGroupConversations(ctx context.Context, groupIDs []int64)](../internal/domain/conversation/repository/group_conversation_repository.go)
  - [ ] `queryPrivateConversationsByOwnerIds(@NotNull Set<Long> ownerIds)`
  - [x] `queryPrivateConversations(@NotNull Collection<Long> ownerIds, @NotNull Long targetId)` -> [internal/domain/conversation/repository/private_conversation_repository.go:QueryPrivateConversations(ctx context.Context, ownerIDs []int64)](../internal/domain/conversation/repository/private_conversation_repository.go)
  - [x] `queryPrivateConversations(@NotNull Set<PrivateConversation.Key> keys)` -> [internal/domain/conversation/repository/private_conversation_repository.go:QueryPrivateConversations(ctx context.Context, ownerIDs []int64)](../internal/domain/conversation/repository/private_conversation_repository.go)
  - [ ] `deletePrivateConversations(@NotNull Set<PrivateConversation.Key> keys)`
  - [ ] `deletePrivateConversations(@NotNull Set<Long> userIds, @Nullable ClientSession session)`
  - [ ] `deleteGroupConversations(@Nullable Set<Long> groupIds, @Nullable ClientSession session)`
  - [ ] `deleteGroupMemberConversations(@NotNull Collection<Long> userIds, @Nullable ClientSession session)`
  - [ ] `authAndUpdateTypingStatus(@NotNull Long requesterId, boolean isGroupMessage, @NotNull Long toId)`

- **ConversationSettingsService.java** ([java/im/turms/service/domain/conversation/service/ConversationSettingsService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/conversation/service/ConversationSettingsService.java))
> [简述功能]

  - [ ] `upsertPrivateConversationSettings(Long ownerId, Long userId, Map<String, Value> settings)`
  - [ ] `upsertGroupConversationSettings(Long ownerId, Long groupId, Map<String, Value> settings)`
  - [x] `deleteSettings(Collection<Long> ownerIds, @Nullable ClientSession clientSession)` -> [internal/domain/user/repository/user_settings_repository.go:DeleteSettings(ctx context.Context, filter interface{})](../internal/domain/user/repository/user_settings_repository.go)
  - [x] `unsetSettings(Long ownerId, @Nullable Set<Long> userIds, @Nullable Set<Long> groupIds, @Nullable Set<String> settingNames)` -> [internal/domain/user/service/user_settings_service.go:UnsetSettings(ctx context.Context, userID int64, keys []string)](../internal/domain/user/service/user_settings_service.go)
  - [x] `querySettings(Long ownerId, @Nullable Collection<Long> userIds, @Nullable Collection<Long> groupIds, @Nullable Set<String> settingNames, @Nullable Date lastUpdatedDateStart)` -> [internal/domain/user/service/user_settings_service.go:QuerySettings(ctx context.Context, filter bson.M)](../internal/domain/user/service/user_settings_service.go)

- **GroupBlocklistController.java** ([java/im/turms/service/domain/group/access/admin/controller/GroupBlocklistController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/access/admin/controller/GroupBlocklistController.java))
> [简述功能]

  - [ ] `addGroupBlockedUser(@RequestBody AddGroupBlockedUserDTO addGroupBlockedUserDTO)`
  - [ ] `queryGroupBlockedUsers(@QueryParam(required = false)`
  - [ ] `queryGroupBlockedUsers(@QueryParam(required = false)`
  - [ ] `updateGroupBlockedUsers(List<GroupBlockedUser.Key> keys, @RequestBody UpdateGroupBlockedUserDTO updateGroupBlockedUserDTO)`
  - [ ] `deleteGroupBlockedUsers(List<GroupBlockedUser.Key> keys)`

- **GroupController.java** ([java/im/turms/service/domain/group/access/admin/controller/GroupController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/access/admin/controller/GroupController.java))
> [简述功能]

  - [ ] `addGroup(@RequestBody AddGroupDTO addGroupDTO)`
  - [ ] `queryGroups(@QueryParam(required = false)`
  - [ ] `queryGroups(@QueryParam(required = false)`
  - [x] `countGroups(@QueryParam(required = false)` -> [internal/domain/user/repository/user_relationship_group_member_repository.go:CountGroups(ctx context.Context, ownerIDs []int64, relatedUserIDs []int64)](../internal/domain/user/repository/user_relationship_group_member_repository.go)
  - [ ] `updateGroups(Set<Long> ids, @RequestBody UpdateGroupDTO updateGroupDTO)`
  - [ ] `deleteGroups(@QueryParam(required = false)`

- **GroupInvitationController.java** ([java/im/turms/service/domain/group/access/admin/controller/GroupInvitationController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/access/admin/controller/GroupInvitationController.java))
> [简述功能]

  - [ ] `addGroupInvitation(@RequestBody AddGroupInvitationDTO addGroupInvitationDTO)`
  - [ ] `queryGroupInvitations(@QueryParam(required = false)`
  - [ ] `queryGroupInvitations(@QueryParam(required = false)`
  - [ ] `updateGroupInvitations(Set<Long> ids, @RequestBody UpdateGroupInvitationDTO updateGroupInvitationDTO)`
  - [ ] `deleteGroupInvitations(@QueryParam(required = false)`

- **GroupJoinRequestController.java** ([java/im/turms/service/domain/group/access/admin/controller/GroupJoinRequestController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/access/admin/controller/GroupJoinRequestController.java))
> [简述功能]

  - [ ] `addGroupJoinRequest(@RequestBody AddGroupJoinRequestDTO addGroupJoinRequestDTO)`
  - [ ] `queryGroupJoinRequests(@QueryParam(required = false)`
  - [ ] `queryGroupJoinRequests(@QueryParam(required = false)`
  - [ ] `updateGroupJoinRequests(Set<Long> ids, @RequestBody UpdateGroupJoinRequestDTO updateGroupJoinRequestDTO)`
  - [ ] `deleteGroupJoinRequests(@QueryParam(required = false)`

- **GroupMemberController.java** ([java/im/turms/service/domain/group/access/admin/controller/GroupMemberController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/access/admin/controller/GroupMemberController.java))
> [简述功能]

  - [x] `addGroupMember(@RequestBody AddGroupMemberDTO addGroupMemberDTO)` -> [internal/domain/group/repository/group_member_repository.go:AddGroupMember(ctx context.Context, member *po.GroupMember)](../internal/domain/group/repository/group_member_repository.go)
  - [ ] `queryGroupMembers(@QueryParam(required = false)`
  - [ ] `queryGroupMembers(@QueryParam(required = false)`
  - [x] `updateGroupMembers(List<GroupMember.Key> keys, @RequestBody UpdateGroupMemberDTO updateGroupMemberDTO)` -> [internal/domain/group/repository/group_member_repository.go:UpdateGroupMembers(ctx context.Context, keys []po.GroupMemberKey, name *string, role *protocol.GroupMemberRole, joinDate *time.Time, muteEndDate *time.Time)](../internal/domain/group/repository/group_member_repository.go)
  - [ ] `deleteGroupMembers(@QueryParam(required = false)`

- **GroupQuestionController.java** ([java/im/turms/service/domain/group/access/admin/controller/GroupQuestionController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/access/admin/controller/GroupQuestionController.java))
> [简述功能]

  - [ ] `queryGroupJoinQuestions(@QueryParam(required = false)`
  - [ ] `queryGroupJoinQuestions(@QueryParam(required = false)`
  - [ ] `addGroupJoinQuestion(@RequestBody AddGroupJoinQuestionDTO addGroupJoinQuestionDTO)`
  - [ ] `updateGroupJoinQuestions(Set<Long> ids, @RequestBody UpdateGroupJoinQuestionDTO updateGroupJoinQuestionDTO)`
  - [ ] `deleteGroupJoinQuestions(@QueryParam(required = false)`

- **GroupTypeController.java** ([java/im/turms/service/domain/group/access/admin/controller/GroupTypeController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/access/admin/controller/GroupTypeController.java))
> [简述功能]

  - [ ] `addGroupType(@RequestBody AddGroupTypeDTO addGroupTypeDTO)`
  - [ ] `queryGroupTypes(@QueryParam(required = false)`
  - [ ] `queryGroupTypes(int page, @QueryParam(required = false)`
  - [x] `updateGroupType(Set<Long> ids, @RequestBody UpdateGroupTypeDTO updateGroupTypeDTO)` -> [internal/domain/group/repository/group_type_repository.go:UpdateGroupType(ctx context.Context, typeID int64, update bson.M)](../internal/domain/group/repository/group_type_repository.go)
  - [ ] `deleteGroupType(Set<Long> ids)`

- **AddGroupBlockedUserDTO.java** ([java/im/turms/service/domain/group/access/admin/dto/request/AddGroupBlockedUserDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/access/admin/dto/request/AddGroupBlockedUserDTO.java))
> [简述功能]

  - [ ] `AddGroupBlockedUserDTO(Long groupId, Long userId, Date blockDate, Long requesterId)`

- **AddGroupDTO.java** ([java/im/turms/service/domain/group/access/admin/dto/request/AddGroupDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/access/admin/dto/request/AddGroupDTO.java))
> [简述功能]

  - [ ] `AddGroupDTO(Long typeId, Long creatorId, Long ownerId, String name, String intro, String announcement, Integer minimumScore, Date creationDate, Date deletionDate, Date muteEndDate, Boolean isActive)`

- **AddGroupInvitationDTO.java** ([java/im/turms/service/domain/group/access/admin/dto/request/AddGroupInvitationDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/access/admin/dto/request/AddGroupInvitationDTO.java))
> [简述功能]

  - [ ] `AddGroupInvitationDTO(Long id, String content, RequestStatus status, Date creationDate, Date responseDate, Long groupId, Long inviterId, Long inviteeId)`

- **AddGroupJoinQuestionDTO.java** ([java/im/turms/service/domain/group/access/admin/dto/request/AddGroupJoinQuestionDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/access/admin/dto/request/AddGroupJoinQuestionDTO.java))
> [简述功能]

  - [ ] `AddGroupJoinQuestionDTO(Long groupId, String question, LinkedHashSet<String> answers, Integer score)`

- **AddGroupJoinRequestDTO.java** ([java/im/turms/service/domain/group/access/admin/dto/request/AddGroupJoinRequestDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/access/admin/dto/request/AddGroupJoinRequestDTO.java))
> [简述功能]

  - [ ] `AddGroupJoinRequestDTO(Long id, String content, RequestStatus status, Date creationDate, Date responseDate, String responseReason, Long groupId, Long requesterId, Long responderId)`

- **AddGroupMemberDTO.java** ([java/im/turms/service/domain/group/access/admin/dto/request/AddGroupMemberDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/access/admin/dto/request/AddGroupMemberDTO.java))
> [简述功能]

  - [ ] `AddGroupMemberDTO(Long groupId, Long userId, String name, GroupMemberRole role, Date joinDate, Date muteEndDate)`

- **AddGroupTypeDTO.java** ([java/im/turms/service/domain/group/access/admin/dto/request/AddGroupTypeDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/access/admin/dto/request/AddGroupTypeDTO.java))
> [简述功能]

  - [ ] `AddGroupTypeDTO(String name, Integer groupSizeLimit, GroupInvitationStrategy invitationStrategy, GroupJoinStrategy joinStrategy, GroupUpdateStrategy groupInfoUpdateStrategy, GroupUpdateStrategy memberInfoUpdateStrategy, Boolean guestSpeakable, Boolean selfInfoUpdatable, Boolean enableReadReceipt, Boolean messageEditable)`

- **UpdateGroupBlockedUserDTO.java** ([java/im/turms/service/domain/group/access/admin/dto/request/UpdateGroupBlockedUserDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/access/admin/dto/request/UpdateGroupBlockedUserDTO.java))
> [简述功能]

  - [ ] `UpdateGroupBlockedUserDTO(Date blockDate, Long requesterId)`

- **UpdateGroupDTO.java** ([java/im/turms/service/domain/group/access/admin/dto/request/UpdateGroupDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/access/admin/dto/request/UpdateGroupDTO.java))
> [简述功能]

  - [ ] `UpdateGroupDTO(Long typeId, Long creatorId, Long ownerId, String name, String intro, String announcement, Integer minimumScore, Boolean isActive, Date creationDate, Date deletionDate, Date muteEndDate, Long successorId, Boolean quitAfterTransfer)`

- **UpdateGroupInvitationDTO.java** ([java/im/turms/service/domain/group/access/admin/dto/request/UpdateGroupInvitationDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/access/admin/dto/request/UpdateGroupInvitationDTO.java))
> [简述功能]

  - [ ] `UpdateGroupInvitationDTO(String content, RequestStatus status, Date creationDate, Date responseDate, Long groupId, Long inviterId, Long inviteeId)`

- **UpdateGroupJoinQuestionDTO.java** ([java/im/turms/service/domain/group/access/admin/dto/request/UpdateGroupJoinQuestionDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/access/admin/dto/request/UpdateGroupJoinQuestionDTO.java))
> [简述功能]

  - [ ] `UpdateGroupJoinQuestionDTO(Long groupId, String question, Set<String> answers, Integer score)`

- **UpdateGroupJoinRequestDTO.java** ([java/im/turms/service/domain/group/access/admin/dto/request/UpdateGroupJoinRequestDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/access/admin/dto/request/UpdateGroupJoinRequestDTO.java))
> [简述功能]

  - [ ] `UpdateGroupJoinRequestDTO(String content, RequestStatus status, Date creationDate, Date responseDate, Long groupId, Long requesterId, Long responderId)`

- **UpdateGroupMemberDTO.java** ([java/im/turms/service/domain/group/access/admin/dto/request/UpdateGroupMemberDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/access/admin/dto/request/UpdateGroupMemberDTO.java))
> [简述功能]

  - [ ] `UpdateGroupMemberDTO(String name, GroupMemberRole role, Date joinDate, Date muteEndDate)`

- **UpdateGroupTypeDTO.java** ([java/im/turms/service/domain/group/access/admin/dto/request/UpdateGroupTypeDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/access/admin/dto/request/UpdateGroupTypeDTO.java))
> [简述功能]

  - [ ] `UpdateGroupTypeDTO(String name, Integer groupSizeLimit, GroupInvitationStrategy invitationStrategy, GroupJoinStrategy joinStrategy, GroupUpdateStrategy groupInfoUpdateStrategy, GroupUpdateStrategy memberInfoUpdateStrategy, Boolean guestSpeakable, Boolean selfInfoUpdatable, Boolean enableReadReceipt, Boolean messageEditable)`

- **GroupStatisticsDTO.java** ([java/im/turms/service/domain/group/access/admin/dto/response/GroupStatisticsDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/access/admin/dto/response/GroupStatisticsDTO.java))
> [简述功能]

  - [ ] `GroupStatisticsDTO(Long deletedGroups, Long groupsThatSentMessages, Long createdGroups, List<StatisticsRecordDTO> deletedGroupsRecords, List<StatisticsRecordDTO> groupsThatSentMessagesRecords, List<StatisticsRecordDTO> createdGroupsRecords)`

- **GroupServiceController.java** ([java/im/turms/service/domain/group/access/servicerequest/controller/GroupServiceController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/access/servicerequest/controller/GroupServiceController.java))
> [简述功能]

  - [x] `handleCreateGroupRequest()` -> [internal/domain/group/controller/group_service_controller.go:HandleCreateGroupRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go)
  - [x] `handleDeleteGroupRequest()` -> [internal/domain/group/controller/group_service_controller.go:HandleDeleteGroupRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go)
  - [x] `handleQueryGroupsRequest()` -> [internal/domain/group/controller/group_service_controller.go:HandleQueryGroupsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go)
  - [x] `handleQueryJoinedGroupIdsRequest()` -> [internal/domain/group/controller/group_service_controller.go:HandleQueryJoinedGroupIdsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go)
  - [ ] `handleQueryJoinedGroupsRequest()`
  - [x] `handleUpdateGroupRequest()` -> [internal/domain/group/controller/group_service_controller.go:HandleUpdateGroupRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go)
  - [x] `handleCreateGroupBlockedUserRequest()` -> [internal/domain/group/controller/group_service_controller.go:HandleCreateGroupBlockedUserRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go)
  - [x] `handleDeleteGroupBlockedUserRequest()` -> [internal/domain/group/controller/group_service_controller.go:HandleDeleteGroupBlockedUserRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go)
  - [x] `handleQueryGroupBlockedUserIdsRequest()` -> [internal/domain/group/controller/group_service_controller.go:HandleQueryGroupBlockedUserIdsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go)
  - [ ] `handleQueryGroupBlockedUsersInfosRequest()`
  - [ ] `handleCheckGroupQuestionAnswerRequest()`
  - [ ] `handleCreateGroupInvitationRequestRequest()`
  - [x] `handleCreateGroupJoinRequestRequest()` -> [internal/domain/group/controller/group_service_controller.go:HandleCreateGroupJoinRequestRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go)
  - [ ] `handleCreateGroupQuestionsRequest()`
  - [x] `handleDeleteGroupInvitationRequest()` -> [internal/domain/group/controller/group_service_controller.go:HandleDeleteGroupInvitationRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go)
  - [x] `handleUpdateGroupInvitationRequest()` -> [internal/domain/group/controller/group_service_controller.go:HandleUpdateGroupInvitationRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go)
  - [x] `handleDeleteGroupJoinRequestRequest()` -> [internal/domain/group/controller/group_service_controller.go:HandleDeleteGroupJoinRequestRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go)
  - [x] `handleUpdateGroupJoinRequestRequest()` -> [internal/domain/group/controller/group_service_controller.go:HandleUpdateGroupJoinRequestRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go)
  - [x] `handleDeleteGroupJoinQuestionsRequest()` -> [internal/domain/group/controller/group_service_controller.go:HandleDeleteGroupJoinQuestionsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go)
  - [x] `handleQueryGroupInvitationsRequest()` -> [internal/domain/group/controller/group_service_controller.go:HandleQueryGroupInvitationsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go)
  - [x] `handleQueryGroupJoinRequestsRequest()` -> [internal/domain/group/controller/group_service_controller.go:HandleQueryGroupJoinRequestsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go)
  - [x] `handleQueryGroupJoinQuestionsRequest()` -> [internal/domain/group/controller/group_service_controller.go:HandleQueryGroupJoinQuestionsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go)
  - [x] `handleUpdateGroupJoinQuestionRequest()` -> [internal/domain/group/controller/group_service_controller.go:HandleUpdateGroupJoinQuestionRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go)
  - [x] `handleCreateGroupMembersRequest()` -> [internal/domain/group/controller/group_service_controller.go:HandleCreateGroupMembersRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go)
  - [x] `handleDeleteGroupMembersRequest()` -> [internal/domain/group/controller/group_service_controller.go:HandleDeleteGroupMembersRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go)
  - [x] `handleQueryGroupMembersRequest()` -> [internal/domain/group/controller/group_service_controller.go:HandleQueryGroupMembersRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go)
  - [x] `handleUpdateGroupMemberRequest()` -> [internal/domain/group/controller/group_service_controller.go:HandleUpdateGroupMemberRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go)

- **CheckGroupQuestionAnswerResult.java** ([java/im/turms/service/domain/group/bo/CheckGroupQuestionAnswerResult.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/bo/CheckGroupQuestionAnswerResult.java))
> [简述功能]

  - [ ] `CheckGroupQuestionAnswerResult(boolean joined, Long groupId, List<Long> questionIds, Integer score)`

- **GroupInvitationStrategy.java** ([java/im/turms/service/domain/group/bo/GroupInvitationStrategy.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/bo/GroupInvitationStrategy.java))
> [简述功能]

  - [x] `requiresApproval()` -> [internal/domain/group/constant/group_strategy.go:RequiresApproval()](../internal/domain/group/constant/group_strategy.go)

- **HandleHandleGroupInvitationResult.java** ([java/im/turms/service/domain/group/bo/HandleHandleGroupInvitationResult.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/bo/HandleHandleGroupInvitationResult.java))
> [简述功能]

  - [ ] `HandleHandleGroupInvitationResult(GroupInvitation groupInvitation, boolean requesterAddedAsNewMember)`

- **HandleHandleGroupJoinRequestResult.java** ([java/im/turms/service/domain/group/bo/HandleHandleGroupJoinRequestResult.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/bo/HandleHandleGroupJoinRequestResult.java))
> [简述功能]

  - [ ] `HandleHandleGroupJoinRequestResult(GroupJoinRequest groupJoinRequest, boolean requesterAddedAsNewMember)`

- **NewGroupQuestion.java** ([java/im/turms/service/domain/group/bo/NewGroupQuestion.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/bo/NewGroupQuestion.java))
> [简述功能]

  - [ ] `NewGroupQuestion(String question, LinkedHashSet<String> answers, Integer score)`

- **GroupBlocklistRepository.java** ([java/im/turms/service/domain/group/repository/GroupBlocklistRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/repository/GroupBlocklistRepository.java))
> [简述功能]

  - [ ] `updateBlockedUsers(Set<GroupBlockedUser.Key> keys, @Nullable Date blockDate, @Nullable Long requesterId)`
  - [x] `count(@Nullable Set<Long> groupIds, @Nullable Set<Long> userIds, @Nullable DateRange blockDateRange, @Nullable Set<Long> requesterIds)` -> [internal/domain/user/repository/user_repository.go:Count(ctx context.Context, filter bson.M)](../internal/domain/user/repository/user_repository.go)
  - [ ] `findBlockedUserIds(Long groupId)`
  - [ ] `findBlockedUsers(@Nullable Set<Long> groupIds, @Nullable Set<Long> userIds, @Nullable DateRange blockDateRange, @Nullable Set<Long> requesterIds, @Nullable Integer page, @Nullable Integer size)`

- **GroupInvitationRepository.java** ([java/im/turms/service/domain/group/repository/GroupInvitationRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/repository/GroupInvitationRepository.java))
> [简述功能]

  - [x] `getEntityExpireAfterSeconds()` -> [internal/domain/user/repository/user_friend_request_repository.go:GetEntityExpireAfterSeconds()](../internal/domain/user/repository/user_friend_request_repository.go)
  - [x] `updateStatusIfPending(Long invitationId, RequestStatus requestStatus, @Nullable String reason, @Nullable ClientSession session)` -> [internal/domain/group/repository/group_invitation_repository.go:UpdateStatusIfPending(ctx context.Context, invitationID int64, newStatus po.RequestStatus, reason *string, responseDate time.Time)](../internal/domain/group/repository/group_invitation_repository.go)
  - [x] `updateInvitations(Set<Long> invitationIds, @Nullable Long inviterId, @Nullable Long inviteeId, @Nullable String content, @Nullable RequestStatus status, @Nullable Date creationDate, @Nullable Date responseDate)` -> [internal/domain/group/repository/group_invitation_repository.go:UpdateInvitations(ctx context.Context, ids []int64, inviterID, inviteeID *int64, content *string, status *po.RequestStatus, creationDate, responseDate *time.Time)](../internal/domain/group/repository/group_invitation_repository.go)
  - [x] `count(@Nullable Set<Long> ids, @Nullable Set<Long> groupIds, @Nullable Set<Long> inviterIds, @Nullable Set<Long> inviteeIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange)` -> [internal/domain/user/repository/user_repository.go:Count(ctx context.Context, filter bson.M)](../internal/domain/user/repository/user_repository.go)
  - [x] `findGroupIdAndInviteeIdAndStatus(Long invitationId)` -> [internal/domain/group/repository/group_invitation_repository.go:FindGroupIdAndInviteeIdAndStatus(ctx context.Context, id int64)](../internal/domain/group/repository/group_invitation_repository.go)
  - [x] `findGroupIdAndInviterIdAndInviteeIdAndStatus(Long invitationId)` -> [internal/domain/group/repository/group_invitation_repository.go:FindGroupIdAndInviterIdAndInviteeIdAndStatus(ctx context.Context, id int64)](../internal/domain/group/repository/group_invitation_repository.go)
  - [x] `findInvitationsByInviteeId(Long inviteeId)` -> [internal/domain/group/repository/group_invitation_repository.go:FindInvitationsByInviteeID(ctx context.Context, inviteeID int64)](../internal/domain/group/repository/group_invitation_repository.go)
  - [ ] `findInvitationsByInviterId(Long inviterId)`
  - [x] `findInvitationsByGroupId(Long groupId)` -> [internal/domain/group/repository/group_invitation_repository.go:FindInvitationsByGroupID(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_invitation_repository.go)
  - [x] `findInviteeIdAndGroupIdAndCreationDateAndStatus(Long invitationId)` -> [internal/domain/group/repository/group_invitation_repository.go:FindInviteeIdAndGroupIdAndCreationDateAndStatus(ctx context.Context, id int64)](../internal/domain/group/repository/group_invitation_repository.go)
  - [x] `findInvitations(@Nullable Set<Long> ids, @Nullable Set<Long> groupIds, @Nullable Set<Long> inviterIds, @Nullable Set<Long> inviteeIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange, @Nullable Integer page, @Nullable Integer size)` -> [internal/domain/group/repository/group_invitation_repository.go:FindInvitations(ctx context.Context, groupID *int64, inviterID *int64, inviteeID *int64, status *po.RequestStatus, lastUpdatedDate *time.Time, page, size int)](../internal/domain/group/repository/group_invitation_repository.go)

- **GroupJoinRequestRepository.java** ([java/im/turms/service/domain/group/repository/GroupJoinRequestRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/repository/GroupJoinRequestRepository.java))
> [简述功能]

  - [x] `getEntityExpireAfterSeconds()` -> [internal/domain/user/repository/user_friend_request_repository.go:GetEntityExpireAfterSeconds()](../internal/domain/user/repository/user_friend_request_repository.go)
  - [x] `updateStatusIfPending(Long requestId, RequestStatus status, Long responderId, @Nullable String reason, @Nullable ClientSession session)` -> [internal/domain/group/repository/group_invitation_repository.go:UpdateStatusIfPending(ctx context.Context, invitationID int64, newStatus po.RequestStatus, reason *string, responseDate time.Time)](../internal/domain/group/repository/group_invitation_repository.go)
  - [ ] `updateRequests(Set<Long> requestIds, @Nullable Long requesterId, @Nullable Long responderId, @Nullable String content, @Nullable RequestStatus status, @Nullable Date creationDate, @Nullable Date responseDate)`
  - [ ] `countRequests(@Nullable Set<Long> ids, @Nullable Set<Long> groupIds, @Nullable Set<Long> requesterIds, @Nullable Set<Long> responderIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange)`
  - [ ] `findGroupId(Long requestId)`
  - [ ] `findRequesterIdAndStatusAndGroupId(Long requestId)`
  - [x] `findRequestsByGroupId(Long groupId)` -> [internal/domain/group/repository/group_join_request_repository.go:FindRequestsByGroupID(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_join_request_repository.go)
  - [x] `findRequestsByRequesterId(Long requesterId)` -> [internal/domain/group/repository/group_join_request_repository.go:FindRequestsByRequesterID(ctx context.Context, requesterID int64)](../internal/domain/group/repository/group_join_request_repository.go)
  - [x] `findRequests(@Nullable Set<Long> ids, @Nullable Set<Long> groupIds, @Nullable Set<Long> requesterIds, @Nullable Set<Long> responderIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange, @Nullable Integer page, @Nullable Integer size)` -> [internal/domain/group/repository/group_join_request_repository.go:FindRequests(ctx context.Context, groupID *int64, requesterID *int64, responderID *int64, status *po.RequestStatus, creationDate *time.Time, responseDate *time.Time, expirationDate *time.Time, page int, size int)](../internal/domain/group/repository/group_join_request_repository.go)

- **GroupMemberRepository.java** ([java/im/turms/service/domain/group/repository/GroupMemberRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/repository/GroupMemberRepository.java))
> [简述功能]

  - [x] `deleteAllGroupMembers(@Nullable Set<Long> groupIds, @Nullable ClientSession session)` -> [internal/domain/group/service/group_member_service.go:DeleteAllGroupMembers(ctx context.Context, groupIDs []int64, session mongo.SessionContext, updateVersion bool)](../internal/domain/group/service/group_member_service.go)
  - [x] `updateGroupMembers(Set<GroupMember.Key> keys, @Nullable String name, @Nullable GroupMemberRole role, @Nullable Date joinDate, @Nullable Date muteEndDate, @Nullable ClientSession session)` -> [internal/domain/group/repository/group_member_repository.go:UpdateGroupMembers(ctx context.Context, keys []po.GroupMemberKey, name *string, role *protocol.GroupMemberRole, joinDate *time.Time, muteEndDate *time.Time)](../internal/domain/group/repository/group_member_repository.go)
  - [x] `countMembers(Long groupId)` -> [internal/domain/group/repository/group_member_repository.go:CountMembers(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_member_repository.go)
  - [x] `countMembers(@Nullable Set<Long> groupIds, @Nullable Set<Long> userIds, @Nullable Set<GroupMemberRole> roles, @Nullable DateRange joinDateRange, @Nullable DateRange muteEndDateRange)` -> [internal/domain/group/repository/group_member_repository.go:CountMembers(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_member_repository.go)
  - [ ] `findGroupManagersAndOwnerId(Long groupId)`
  - [ ] `findGroupMembers(Long groupId)`
  - [ ] `findGroupMembers(Long groupId, Set<Long> memberIds)`
  - [ ] `findGroupsMembers(@Nullable Set<Long> groupIds, @Nullable Set<Long> userIds, @Nullable Set<@ValidGroupMemberRole GroupMemberRole> roles, @Nullable DateRange joinDateRange, @Nullable DateRange muteEndDateRange, @Nullable Integer page, @Nullable Integer size)`
  - [x] `findGroupMemberIds(Long groupId)` -> [internal/domain/group/repository/group_member_repository.go:FindGroupMemberIDs(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_member_repository.go)
  - [x] `findGroupMemberIds(Set<Long> groupIds)` -> [internal/domain/group/repository/group_member_repository.go:FindGroupMemberIDs(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_member_repository.go)
  - [ ] `findGroupMemberKeyAndRoleParis(Set<Long> userIds, Long groupId)`
  - [x] `findGroupMemberRole(Long userId, Long groupId)` -> [internal/domain/group/repository/group_member_repository.go:FindGroupMemberRole(ctx context.Context, groupID, userID int64)](../internal/domain/group/repository/group_member_repository.go)
  - [ ] `findMemberIdsByGroupId(Long groupId)`
  - [x] `findUserJoinedGroupIds(Long userId)` -> [internal/domain/group/repository/group_member_repository.go:FindUserJoinedGroupIDs(ctx context.Context, userID int64)](../internal/domain/group/repository/group_member_repository.go)
  - [ ] `findUsersJoinedGroupIds(@Nullable Set<Long> groupIds, @NotEmpty Set<Long> userIds, @Nullable Integer page, @Nullable Integer size)`
  - [x] `isMemberMuted(Long groupId, Long userId)` -> [internal/domain/group/repository/group_member_repository.go:IsMemberMuted(ctx context.Context, groupID, userID int64)](../internal/domain/group/repository/group_member_repository.go)

- **GroupQuestionRepository.java** ([java/im/turms/service/domain/group/repository/GroupQuestionRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/repository/GroupQuestionRepository.java))
> [简述功能]

  - [ ] `updateQuestion(Long questionId, @Nullable String question, @Nullable Set<String> answers, @Nullable Integer score)`
  - [ ] `updateQuestions(Set<Long> ids, @Nullable Long groupId, @Nullable String question, @Nullable Set<String> answers, @Nullable Integer score)`
  - [ ] `countQuestions(@Nullable Set<Long> ids, @Nullable Set<Long> groupIds)`
  - [ ] `checkQuestionAnswerAndGetScore(Long questionId, String answer, @Nullable Long groupId)`
  - [ ] `findGroupId(Long questionId)`
  - [ ] `findQuestions(@Nullable Set<Long> ids, @Nullable Set<Long> groupIds, @Nullable Integer page, @Nullable Integer size, boolean withAnswers)`

- **GroupRepository.java** ([java/im/turms/service/domain/group/repository/GroupRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/repository/GroupRepository.java))
> [简述功能]

  - [ ] `updateGroupsDeletionDate(@Nullable Collection<Long> groupIds, @Nullable ClientSession session)`
  - [ ] `updateGroups(Set<Long> groupIds, @Nullable Long typeId, @Nullable Long creatorId, @Nullable Long ownerId, @Nullable String name, @Nullable String intro, @Nullable String announcement, @Nullable Integer minimumScore, @Nullable Boolean isActive, @Nullable Date creationDate, @Nullable Date deletionDate, @Nullable Date muteEndDate, @Nullable Date lastUpdatedDate, @Nullable Map<String, Object> userDefinedAttributes, @Nullable ClientSession session)`
  - [ ] `countCreatedGroups(@Nullable DateRange dateRange)`
  - [ ] `countDeletedGroups(@Nullable DateRange dateRange)`
  - [x] `countGroups(@Nullable Set<Long> ids, @Nullable Set<Long> typeIds, @Nullable Set<Long> creatorIds, @Nullable Set<Long> ownerIds, @Nullable Boolean isActive, @Nullable DateRange creationDateRange, @Nullable DateRange deletionDateRange, @Nullable DateRange lastUpdatedDateRange, @Nullable DateRange muteEndDateRange)` -> [internal/domain/user/repository/user_relationship_group_member_repository.go:CountGroups(ctx context.Context, ownerIDs []int64, relatedUserIDs []int64)](../internal/domain/user/repository/user_relationship_group_member_repository.go)
  - [x] `countOwnedGroups(Long ownerId)` -> [internal/domain/group/repository/group_repository.go:CountOwnedGroups(ctx context.Context, ownerID int64)](../internal/domain/group/repository/group_repository.go)
  - [x] `countOwnedGroups(Long ownerId, Long groupTypeId)` -> [internal/domain/group/repository/group_repository.go:CountOwnedGroups(ctx context.Context, ownerID int64)](../internal/domain/group/repository/group_repository.go)
  - [x] `findGroups(@Nullable Set<Long> ids, @Nullable Set<Long> typeIds, @Nullable Set<Long> creatorIds, @Nullable Set<Long> ownerIds, @Nullable Boolean isActive, @Nullable DateRange creationDateRange, @Nullable DateRange deletionDateRange, @Nullable DateRange lastUpdatedDateRange, @Nullable DateRange muteEndDateRange, @Nullable Integer page, @Nullable Integer size)` -> [internal/domain/group/repository/group_repository.go:FindGroups(ctx context.Context, groupIDs []int64)](../internal/domain/group/repository/group_repository.go)
  - [ ] `findNotDeletedGroups(Collection<Long> ids, @Nullable Date lastUpdatedDate)`
  - [x] `findAllNames()` -> [internal/domain/user/repository/user_repository.go:FindAllNames(ctx context.Context)](../internal/domain/user/repository/user_repository.go)
  - [ ] `findTypeId(Long groupId)`
  - [ ] `findTypeIdAndGroupId(Collection<Long> groupIds)`
  - [ ] `findTypeIdIfActiveAndNotDeleted(Long groupId)`
  - [ ] `findMinimumScore(Long groupId)`
  - [ ] `findOwnerId(Long groupId)`
  - [ ] `isGroupMuted(Long groupId, Date muteEndDate)`
  - [ ] `isGroupActiveAndNotDeleted(Long groupId)`

- **GroupTypeRepository.java** ([java/im/turms/service/domain/group/repository/GroupTypeRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/repository/GroupTypeRepository.java))
> [简述功能]

  - [ ] `updateTypes(Set<Long> ids, @Nullable String name, @Nullable Integer groupSizeLimit, @Nullable GroupInvitationStrategy groupInvitationStrategy, @Nullable GroupJoinStrategy groupJoinStrategy, @Nullable GroupUpdateStrategy groupInfoUpdateStrategy, @Nullable GroupUpdateStrategy memberInfoUpdateStrategy, @Nullable Boolean guestSpeakable, @Nullable Boolean selfInfoUpdatable, @Nullable Boolean enableReadReceipt, @Nullable Boolean messageEditable)`

- **GroupVersionRepository.java** ([java/im/turms/service/domain/group/repository/GroupVersionRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/repository/GroupVersionRepository.java))
> [简述功能]

  - [ ] `updateVersions(String field)`
  - [ ] `updateVersions(@Nullable Set<Long> groupIds, String field)`
  - [x] `updateVersion(Long groupId, String field)` -> [internal/domain/group/repository/group_version_repository.go:UpdateVersion(ctx context.Context, groupID int64, field string)](../internal/domain/group/repository/group_version_repository.go)
  - [x] `updateVersion(Long groupId, boolean updateMembers, boolean updateBlocklist, boolean joinRequests, boolean joinQuestions)` -> [internal/domain/group/repository/group_version_repository.go:UpdateVersion(ctx context.Context, groupID int64, field string)](../internal/domain/group/repository/group_version_repository.go)
  - [ ] `findBlocklist(Long groupId)`
  - [x] `findInvitations(Long groupId)` -> [internal/domain/group/repository/group_invitation_repository.go:FindInvitations(ctx context.Context, groupID *int64, inviterID *int64, inviteeID *int64, status *po.RequestStatus, lastUpdatedDate *time.Time, page, size int)](../internal/domain/group/repository/group_invitation_repository.go)
  - [ ] `findJoinRequests(Long groupId)`
  - [ ] `findJoinQuestions(Long groupId)`
  - [ ] `findMembers(Long groupId)`

- **GroupBlocklistService.java** ([java/im/turms/service/domain/group/service/GroupBlocklistService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/service/GroupBlocklistService.java))
> [简述功能]

  - [x] `authAndBlockUser(@NotNull Long requesterId, @NotNull Long groupId, @NotNull Long userIdToBlock, @Nullable ClientSession session)` -> [internal/domain/group/service/group_blocklist_service.go:AuthAndBlockUser(ctx context.Context, requesterID int64, groupID int64, userID int64,)](../internal/domain/group/service/group_blocklist_service.go)
  - [x] `unblockUser(@NotNull Long requesterId, @NotNull Long groupId, @NotNull Long userIdToUnblock, @Nullable ClientSession session, boolean updateBlocklistVersion)` -> [internal/domain/group/service/group_blocklist_service.go:UnblockUser(ctx context.Context, groupID int64, userID int64)](../internal/domain/group/service/group_blocklist_service.go)
  - [ ] `findBlockedUserIds(@NotNull Long groupId, @NotNull Set<Long> userIds)`
  - [x] `isBlocked(@NotNull Long groupId, @NotNull Long userId)` -> [internal/domain/group/service/group_blocklist_service.go:IsBlocked(ctx context.Context, groupID int64, userID int64)](../internal/domain/group/service/group_blocklist_service.go)
  - [ ] `queryGroupBlockedUserIds(@NotNull Long groupId)`
  - [x] `queryBlockedUsers(@Nullable Set<Long> groupIds, @Nullable Set<Long> userIds, @Nullable DateRange blockDateRange, @Nullable Set<Long> requesterIds, @Nullable Integer page, @Nullable Integer size)` -> [internal/domain/group/service/group_blocklist_service.go:QueryBlockedUsers(ctx context.Context, groupID int64)](../internal/domain/group/service/group_blocklist_service.go)
  - [ ] `countBlockedUsers(@Nullable Set<Long> groupIds, @Nullable Set<Long> userIds, @Nullable DateRange blockDateRange, @Nullable Set<Long> requesterIds)`
  - [ ] `queryGroupBlockedUserIdsWithVersion(@NotNull Long groupId, @Nullable Date lastUpdatedDate)`
  - [ ] `queryGroupBlockedUserInfosWithVersion(@NotNull Long groupId, @Nullable Date lastUpdatedDate)`
  - [ ] `addBlockedUser(@NotNull Long groupId, @NotNull Long userId, @NotNull Long requesterId, @Nullable @PastOrPresent Date blockDate)`
  - [ ] `updateBlockedUsers(@NotEmpty Set<GroupBlockedUser.@ValidGroupBlockedUserKey Key> keys, @Nullable @PastOrPresent Date blockDate, @Nullable Long requesterId)`
  - [ ] `deleteBlockedUsers(@NotEmpty Set<GroupBlockedUser.@ValidGroupBlockedUserKey Key> keys)`

- **GroupInvitationService.java** ([java/im/turms/service/domain/group/service/GroupInvitationService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/service/GroupInvitationService.java))
> [简述功能]

  - [x] `authAndCreateGroupInvitation(@NotNull Long groupId, @NotNull Long inviterId, @NotNull Long inviteeId, @Nullable String content)` -> [internal/domain/group/service/group_invitation_service.go:AuthAndCreateGroupInvitation(ctx context.Context, requesterID int64, groupID int64, inviteeID int64, content string,)](../internal/domain/group/service/group_invitation_service.go)
  - [ ] `createGroupInvitation(@Nullable Long id, @NotNull Long groupId, @NotNull Long inviterId, @NotNull Long inviteeId, @Nullable String content, @Nullable @ValidRequestStatus RequestStatus status, @Nullable @PastOrPresent Date creationDate, @Nullable @PastOrPresent Date responseDate)`
  - [ ] `queryGroupIdAndInviterIdAndInviteeIdAndStatus(@NotNull Long invitationId)`
  - [ ] `queryGroupIdAndInviteeIdAndStatus(@NotNull Long invitationId)`
  - [x] `authAndRecallPendingGroupInvitation(@NotNull Long requesterId, @NotNull Long invitationId)` -> [internal/domain/group/service/group_invitation_service.go:AuthAndRecallPendingGroupInvitation(ctx context.Context, requesterID int64, invitationID int64,)](../internal/domain/group/service/group_invitation_service.go)
  - [ ] `queryGroupInvitationsByInviteeId(@NotNull Long inviteeId)`
  - [ ] `queryGroupInvitationsByInviterId(@NotNull Long inviterId)`
  - [ ] `queryGroupInvitationsByGroupId(@NotNull Long groupId)`
  - [x] `queryUserGroupInvitationsWithVersion(@NotNull Long userId, boolean areSentByUser, @Nullable Date lastUpdatedDate)` -> [internal/domain/group/service/group_invitation_service.go:QueryUserGroupInvitationsWithVersion(ctx context.Context, userID int64, areSentInvitations bool, lastUpdatedDate *time.Time)](../internal/domain/group/service/group_invitation_service.go)
  - [x] `authAndQueryGroupInvitationsWithVersion(@NotNull Long userId, @NotNull Long groupId, @Nullable Date lastUpdatedDate)` -> [internal/domain/group/service/group_invitation_service.go:AuthAndQueryGroupInvitationsWithVersion(ctx context.Context, requesterID int64, groupID int64, lastUpdatedDate *time.Time)](../internal/domain/group/service/group_invitation_service.go)
  - [ ] `queryInviteeIdAndGroupIdAndCreationDateAndStatusByInvitationId(@NotNull Long invitationId)`
  - [x] `queryInvitations(@Nullable Set<Long> ids, @Nullable Set<Long> groupIds, @Nullable Set<Long> inviterIds, @Nullable Set<Long> inviteeIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange, @Nullable Integer page, @Nullable Integer size)` -> [internal/domain/group/service/group_invitation_service.go:QueryInvitations(ctx context.Context, groupID *int64, inviterID *int64, inviteeID *int64, status *po.RequestStatus, lastUpdatedDate *time.Time, page int, size int)](../internal/domain/group/service/group_invitation_service.go)
  - [x] `countInvitations(@Nullable Set<Long> ids, @Nullable Set<Long> groupIds, @Nullable Set<Long> inviterIds, @Nullable Set<Long> inviteeIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange)` -> [internal/domain/group/repository/group_invitation_repository.go:CountInvitations(ctx context.Context, groupID, inviterID, inviteeID *int64, status *po.RequestStatus, lastUpdatedDate *time.Time)](../internal/domain/group/repository/group_invitation_repository.go)
  - [x] `deleteInvitations(@Nullable Set<Long> ids)` -> [internal/domain/group/repository/group_invitation_repository.go:DeleteInvitations(ctx context.Context, ids []int64)](../internal/domain/group/repository/group_invitation_repository.go)
  - [x] `authAndHandleInvitation(@NotNull Long requesterId, @NotNull Long invitationId, @NotNull @ValidResponseAction ResponseAction action, @Nullable String reason)` -> [internal/domain/group/service/group_invitation_service.go:AuthAndHandleInvitation(ctx context.Context, requesterID int64, invitationID int64, status po.RequestStatus, reason string,)](../internal/domain/group/service/group_invitation_service.go)
  - [ ] `updatePendingInvitationStatus(@NotNull Long groupId, @NotNull Long invitationId, @NotNull @ValidRequestStatus RequestStatus requestStatus, @Nullable String reason, @Nullable ClientSession session)`
  - [x] `updateInvitations(@NotEmpty Set<Long> invitationIds, @Nullable Long inviterId, @Nullable Long inviteeId, @Nullable String content, @Nullable @ValidRequestStatus RequestStatus status, @Nullable @PastOrPresent Date creationDate, @Nullable @PastOrPresent Date responseDate)` -> [internal/domain/group/repository/group_invitation_repository.go:UpdateInvitations(ctx context.Context, ids []int64, inviterID, inviteeID *int64, content *string, status *po.RequestStatus, creationDate, responseDate *time.Time)](../internal/domain/group/repository/group_invitation_repository.go)

- **GroupJoinRequestService.java** ([java/im/turms/service/domain/group/service/GroupJoinRequestService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/service/GroupJoinRequestService.java))
> [简述功能]

  - [ ] `authAndCreateGroupJoinRequest(@NotNull Long requesterId, @NotNull Long groupId, @Nullable String content)`
  - [ ] `authAndRecallPendingGroupJoinRequest(@NotNull Long requesterId, @NotNull Long requestId)`
  - [ ] `authAndQueryGroupJoinRequestsWithVersion(@NotNull Long requesterId, @Nullable Long groupId, @Nullable Date lastUpdatedDate)`
  - [ ] `queryGroupJoinRequestsByGroupId(@NotNull Long groupId)`
  - [ ] `queryGroupJoinRequestsByRequesterId(@NotNull Long requesterId)`
  - [ ] `queryGroupId(@NotNull Long requestId)`
  - [x] `queryJoinRequests(@Nullable Set<Long> ids, @Nullable Set<Long> groupIds, @Nullable Set<Long> requesterIds, @Nullable Set<Long> responderIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange, @Nullable Integer page, @Nullable Integer size)` -> [internal/domain/group/service/group_join_request_service.go:QueryJoinRequests(ctx context.Context, groupID *int64, requesterID *int64, responderID *int64, status *po.RequestStatus, creationDate *time.Time, page int, size int)](../internal/domain/group/service/group_join_request_service.go)
  - [ ] `countJoinRequests(@Nullable Set<Long> ids, @Nullable Set<Long> groupIds, @Nullable Set<Long> requesterIds, @Nullable Set<Long> responderIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange)`
  - [ ] `deleteJoinRequests(@Nullable Set<Long> ids)`
  - [x] `authAndHandleJoinRequest(@NotNull Long requesterId, @NotNull Long joinRequestId, @NotNull @ValidResponseAction ResponseAction action, @Nullable String responseReason)` -> [internal/domain/group/service/group_join_request_service.go:AuthAndHandleJoinRequest(ctx context.Context, responderID int64, requestID int64, status po.RequestStatus, reason string)](../internal/domain/group/service/group_join_request_service.go)
  - [ ] `updatePendingJoinRequestStatus(@NotNull Long groupId, @NotNull Long joinRequestId, @NotNull @ValidRequestStatus RequestStatus requestStatus, @NotNull Long responderId, @Nullable String responseReason, @Nullable ClientSession session)`
  - [ ] `updateJoinRequests(@NotEmpty Set<Long> requestIds, @Nullable Long requesterId, @Nullable Long responderId, @Nullable String content, @Nullable @ValidRequestStatus RequestStatus status, @Nullable @PastOrPresent Date creationDate, @Nullable @PastOrPresent Date responseDate)`
  - [ ] `createGroupJoinRequest(@Nullable Long id, @NotNull Long groupId, @NotNull Long requesterId, @NotNull Long responderId, @Nullable String content, @Nullable @ValidRequestStatus RequestStatus status, @Nullable @PastOrPresent Date creationDate, @Nullable @PastOrPresent Date responseDate, @Nullable String responseReason)`

- **GroupMemberService.java** ([java/im/turms/service/domain/group/service/GroupMemberService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/service/GroupMemberService.java))
> [简述功能]

  - [x] `addGroupMember(@NotNull Long groupId, @NotNull Long userId, @NotNull @ValidGroupMemberRole GroupMemberRole groupMemberRole, @Nullable String name, @Nullable @PastOrPresent Date joinDate, @Nullable Date muteEndDate, @Nullable ClientSession session)` -> [internal/domain/group/repository/group_member_repository.go:AddGroupMember(ctx context.Context, member *po.GroupMember)](../internal/domain/group/repository/group_member_repository.go)
  - [x] `addGroupMembers(@NotNull Long groupId, @NotNull Set<Long> userIds, @NotNull @ValidGroupMemberRole GroupMemberRole groupMemberRole, @Nullable String name, @Nullable @PastOrPresent Date joinDate, @Nullable Date muteEndDate, @Nullable ClientSession session)` -> [internal/domain/group/service/group_member_service.go:AddGroupMembers(ctx context.Context, groupID int64, userIDs []int64, role protocol.GroupMemberRole, name *string, joinTime *time.Time, muteEndDate *time.Time, session mongo.SessionContext,)](../internal/domain/group/service/group_member_service.go)
  - [x] `authAndAddGroupMembers(@NotNull Long requesterId, @NotNull Long groupId, @NotNull Set<Long> userIds, @Nullable @ValidGroupMemberRole GroupMemberRole groupMemberRole, @Nullable String name, @Nullable Date muteEndDate, @Nullable ClientSession session)` -> [internal/domain/group/service/group_member_service.go:AuthAndAddGroupMembers(ctx context.Context, requesterID int64, groupID int64, userIDs []int64, role protocol.GroupMemberRole, muteEndDate *time.Time,)](../internal/domain/group/service/group_member_service.go)
  - [x] `authAndDeleteGroupMembers(@NotNull Long requesterId, @NotNull Long groupId, @NotNull Set<Long> memberIdsToDelete, @Nullable Long successorId, @Nullable Boolean quitAfterTransfer)` -> [internal/domain/group/service/group_member_service.go:AuthAndDeleteGroupMembers(ctx context.Context, requesterID int64, groupID int64, userIDs []int64, successorID *int64, quitAfterTransfer bool,)](../internal/domain/group/service/group_member_service.go)
  - [x] `deleteGroupMember(@NotNull Long groupId, @NotNull Long memberId, @Nullable ClientSession session, boolean updateGroupMembersVersion)` -> [internal/domain/group/service/group_member_service.go:DeleteGroupMember(ctx context.Context, groupID, userID int64, session mongo.SessionContext, updateVersion bool,)](../internal/domain/group/service/group_member_service.go)
  - [ ] `deleteGroupMembers(@NotEmpty Collection<GroupMember.Key> keys, @Nullable ClientSession session, boolean updateGroupMembersVersion)`
  - [ ] `updateGroupMember(@NotNull Long groupId, @NotNull Long memberId, @Nullable String name, @Nullable @ValidGroupMemberRole GroupMemberRole role, @Nullable @PastOrPresent Date joinDate, @Nullable Date muteEndDate, @Nullable ClientSession session, boolean updateGroupMembersVersion)`
  - [x] `updateGroupMembers(@NotEmpty Set<GroupMember.Key> keys, @Nullable String name, @Nullable @ValidGroupMemberRole GroupMemberRole role, @Nullable @PastOrPresent Date joinDate, @Nullable Date muteEndDate, @Nullable ClientSession session, boolean updateGroupMembersVersion)` -> [internal/domain/group/repository/group_member_repository.go:UpdateGroupMembers(ctx context.Context, keys []po.GroupMemberKey, name *string, role *protocol.GroupMemberRole, joinDate *time.Time, muteEndDate *time.Time)](../internal/domain/group/repository/group_member_repository.go)
  - [x] `updateGroupMembers(@NotNull Long groupId, @NotEmpty Set<Long> memberIds, @Nullable String name, @Nullable @ValidGroupMemberRole GroupMemberRole role, @Nullable @PastOrPresent Date joinDate, @Nullable Date muteEndDate, @Nullable ClientSession session, boolean updateGroupMembersVersion)` -> [internal/domain/group/repository/group_member_repository.go:UpdateGroupMembers(ctx context.Context, keys []po.GroupMemberKey, name *string, role *protocol.GroupMemberRole, joinDate *time.Time, muteEndDate *time.Time)](../internal/domain/group/repository/group_member_repository.go)
  - [x] `isGroupMember(@NotNull Long groupId, @NotNull Long userId, boolean preferCache)` -> [internal/domain/group/repository/group_member_repository.go:IsGroupMember(ctx context.Context, groupID, userID int64)](../internal/domain/group/repository/group_member_repository.go)
  - [x] `isGroupMember(@NotEmpty Set<Long> groupIds, @NotNull Long userId)` -> [internal/domain/group/repository/group_member_repository.go:IsGroupMember(ctx context.Context, groupID, userID int64)](../internal/domain/group/repository/group_member_repository.go)
  - [ ] `findExistentMemberGroupIds(@NotEmpty Set<Long> groupIds, @NotNull Long userId)`
  - [ ] `isAllowedToInviteUser(@NotNull Long groupId, @NotNull Long inviterId)`
  - [ ] `isAllowedToBeInvited(@NotNull Long groupId, @NotNull Long inviteeId)`
  - [ ] `isAllowedToSendMessage(@NotNull Long groupId, @NotNull Long senderId)`
  - [x] `isMemberMuted(@NotNull Long groupId, @NotNull Long userId, boolean preferCache)` -> [internal/domain/group/repository/group_member_repository.go:IsMemberMuted(ctx context.Context, groupID, userID int64)](../internal/domain/group/repository/group_member_repository.go)
  - [ ] `queryGroupMemberKeyAndRolePairs(@NotNull Set<Long> userIds, @NotNull Long groupId)`
  - [x] `queryGroupMemberRole(@NotNull Long userId, @NotNull Long groupId, boolean preferCache)` -> [internal/domain/group/service/group_member_service.go:QueryGroupMemberRole(ctx context.Context, groupID, userID int64)](../internal/domain/group/service/group_member_service.go)
  - [x] `isOwner(@NotNull Long userId, @NotNull Long groupId, boolean preferCache)` -> [internal/domain/group/service/group_member_service.go:IsOwner(ctx context.Context, userID, groupID int64)](../internal/domain/group/service/group_member_service.go)
  - [x] `isOwnerOrManager(@NotNull Long userId, @NotNull Long groupId, boolean preferCache)` -> [internal/domain/group/service/group_member_service.go:IsOwnerOrManager(ctx context.Context, groupID, userID int64)](../internal/domain/group/service/group_member_service.go)
  - [ ] `isOwnerOrManagerOrMember(@NotNull Long userId, @NotNull Long groupId, boolean preferCache)`
  - [ ] `queryUserJoinedGroupIds(@NotNull Long userId)`
  - [ ] `queryUsersJoinedGroupIds(@Nullable Set<Long> groupIds, @NotEmpty Set<Long> userIds, @Nullable Integer page, @Nullable Integer size)`
  - [ ] `queryMemberIdsInUsersJoinedGroups(@NotEmpty Set<Long> userIds, boolean preferCache)`
  - [ ] `queryGroupMemberIds(@NotNull Long groupId, boolean preferCache)`
  - [ ] `queryGroupMemberIds(@NotEmpty Set<Long> groupIds, boolean preferCache)`
  - [ ] `queryGroupMembers(@Nullable Set<Long> groupIds, @Nullable Set<Long> userIds, @Nullable Set<@ValidGroupMemberRole GroupMemberRole> roles, @Nullable DateRange joinDateRange, @Nullable DateRange muteEndDateRange, @Nullable Integer page, @Nullable Integer size)`
  - [x] `countMembers(@Nullable Set<Long> groupIds, @Nullable Set<Long> userIds, @Nullable Set<@ValidGroupMemberRole GroupMemberRole> roles, @Nullable DateRange joinDateRange, @Nullable DateRange muteEndDateRange)` -> [internal/domain/group/repository/group_member_repository.go:CountMembers(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_member_repository.go)
  - [ ] `deleteGroupMembers(boolean updateGroupMembersVersion)`
  - [ ] `queryGroupMembers(@NotNull Long groupId, boolean preferCache)`
  - [ ] `queryGroupMembers(@NotNull Long groupId, @NotEmpty Set<Long> memberIds, boolean preferCache)`
  - [ ] `authAndQueryGroupMembers(@NotNull Long requesterId, @NotNull Long groupId, @NotEmpty Set<Long> memberIds, boolean withStatus)`
  - [ ] `authAndQueryGroupMembersWithVersion(@NotNull Long requesterId, @NotNull Long groupId, @Nullable Date lastUpdatedDate, boolean withStatus)`
  - [x] `authAndUpdateGroupMember(@NotNull Long requesterId, @NotNull Long groupId, @NotNull Long memberId, @Nullable String name, @Nullable @ValidGroupMemberRole GroupMemberRole role, @Nullable Date muteEndDate)` -> [internal/domain/group/service/group_member_service.go:AuthAndUpdateGroupMember(ctx context.Context, requesterID int64, groupID int64, memberID int64, name *string, role *protocol.GroupMemberRole, muteEndDate *time.Time,)](../internal/domain/group/service/group_member_service.go)
  - [x] `deleteAllGroupMembers(@Nullable Set<Long> groupIds, @Nullable ClientSession session, boolean updateMembersVersion)` -> [internal/domain/group/service/group_member_service.go:DeleteAllGroupMembers(ctx context.Context, groupIDs []int64, session mongo.SessionContext, updateVersion bool)](../internal/domain/group/service/group_member_service.go)
  - [ ] `queryGroupManagersAndOwnerId(@NotNull Long groupId)`

- **GroupQuestionService.java** ([java/im/turms/service/domain/group/service/GroupQuestionService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/service/GroupQuestionService.java))
> [简述功能]

  - [ ] `checkGroupQuestionAnswerAndGetScore(@NotNull Long questionId, @NotNull String answer, @Nullable Long groupId)`
  - [ ] `authAndCheckGroupQuestionAnswerAndJoin(@NotNull Long requesterId, @NotNull @ValidGroupQuestionIdAndAnswer Map<Long, String> questionIdToAnswer)`
  - [ ] `authAndCreateGroupJoinQuestions(@NotNull Long requesterId, @NotNull Long groupId, @NotNull List<NewGroupQuestion> questions)`
  - [ ] `createGroupJoinQuestions(@NotNull Long groupId, @NotNull List<NewGroupQuestion> questions)`
  - [ ] `queryGroupId(@NotNull Long questionId)`
  - [ ] `authAndDeleteGroupJoinQuestions(@NotNull Long requesterId, @NotNull Long groupId, @NotNull Set<Long> questionIds)`
  - [ ] `queryGroupJoinQuestions(@Nullable Set<Long> ids, @Nullable Set<Long> groupIds, @Nullable Integer page, @Nullable Integer size, boolean withAnswers)`
  - [ ] `countGroupJoinQuestions(@Nullable Set<Long> ids, @Nullable Set<Long> groupIds)`
  - [ ] `deleteGroupJoinQuestions(@Nullable Set<Long> ids)`
  - [ ] `authAndQueryGroupJoinQuestionsWithVersion(@NotNull Long requesterId, @NotNull Long groupId, boolean withAnswers, @Nullable Date lastUpdatedDate)`
  - [ ] `authAndUpdateGroupJoinQuestion(@NotNull Long requesterId, @NotNull Long questionId, @Nullable String question, @Nullable Set<String> answers, @Nullable @Min(0)`
  - [ ] `updateGroupJoinQuestions(@NotEmpty Set<Long> ids, @Nullable Long groupId, @Nullable String question, @Nullable Set<String> answers, @Nullable @Min(0)`

- **GroupService.java** ([java/im/turms/service/domain/group/service/GroupService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/service/GroupService.java))
> [简述功能]

  - [x] `createGroup(@NotNull Long creatorId, @NotNull Long ownerId, @Nullable String groupName, @Nullable String intro, @Nullable String announcement, @Nullable @Min(value = 0)` -> [internal/domain/group/service/group_service.go:CreateGroup(ctx context.Context, creatorID, groupID int64, name, intro *string, minimumScore *int32)](../internal/domain/group/service/group_service.go)
  - [x] `authAndDeleteGroup(boolean queryGroupMemberIds, @NotNull Long requesterId, @NotNull Long groupId)` -> [internal/domain/group/service/group_service.go:AuthAndDeleteGroup(ctx context.Context, requesterID int64, groupID int64)](../internal/domain/group/service/group_service.go)
  - [ ] `authAndCreateGroup(@NotNull Long creatorId, @NotNull Long ownerId, @Nullable String groupName, @Nullable String intro, @Nullable String announcement, @Nullable @Min(value = 0)`
  - [x] `deleteGroupsAndGroupMembers(@Nullable Set<Long> groupIds, @Nullable Boolean deleteLogically)` -> [internal/domain/group/service/group_service.go:DeleteGroupsAndGroupMembers(ctx context.Context, groupIDs []int64, session mongo.SessionContext)](../internal/domain/group/service/group_service.go)
  - [ ] `queryGroups(@Nullable Set<Long> ids, @Nullable Set<Long> typeIds, @Nullable Set<Long> creatorIds, @Nullable Set<Long> ownerIds, @Nullable Boolean isActive, @Nullable DateRange creationDateRange, @Nullable DateRange deletionDateRange, @Nullable DateRange lastUpdatedDateRange, @Nullable DateRange muteEndDateRange, @Nullable Set<Long> memberIds, @Nullable Integer page, @Nullable Integer size)`
  - [ ] `queryGroupTypeIfActiveAndNotDeleted(@NotNull Long groupId)`
  - [ ] `queryGroupTypeIfActiveAndNotDeleted(@NotNull Long groupId, boolean preferCache)`
  - [ ] `queryGroupTypeId(@NotNull Long groupId)`
  - [x] `queryGroupTypeIdIfActiveAndNotDeleted(@NotNull Long groupId)` -> [internal/domain/group/service/group_service.go:QueryGroupTypeIdIfActiveAndNotDeleted(ctx context.Context, groupID int64)](../internal/domain/group/service/group_service.go)
  - [ ] `queryGroupMinimumScore(@NotNull Long groupId)`
  - [x] `authAndTransferGroupOwnership(@NotNull Long requesterId, @NotNull Long groupId, @NotNull Long successorId, boolean quitAfterTransfer, @Nullable ClientSession session)` -> [internal/domain/group/service/group_service.go:AuthAndTransferGroupOwnership(ctx context.Context, requesterID, groupID, successorID int64, quitAfterTransfer bool, session mongo.SessionContext,)](../internal/domain/group/service/group_service.go)
  - [ ] `queryGroupOwnerId(@NotNull Long groupId)`
  - [ ] `checkAndTransferGroupOwnership(@NotEmpty Set<Long> groupIds, @NotNull Long successorId, boolean quitAfterTransfer)`
  - [ ] `checkAndTransferGroupOwnership(@Nullable Long auxiliaryCurrentOwnerId, @NotNull Long groupId, @NotNull Long successorId, boolean quitAfterTransfer, @Nullable ClientSession session)`
  - [ ] `updateGroupInformation(@NotNull Long groupId, @Nullable Long typeId, @Nullable Long creatorId, @Nullable Long ownerId, @Nullable String name, @Nullable String intro, @Nullable String announcement, @Nullable @Min(0)`
  - [ ] `updateGroupsInformation(@NotNull Set<Long> groupIds, @Nullable Long typeId, @Nullable Long creatorId, @Nullable Long ownerId, @Nullable String name, @Nullable String intro, @Nullable String announcement, @Nullable @Min(0)`
  - [ ] `authAndUpdateGroupInformation(@NotNull Long requesterId, @NotNull Long groupId, @Nullable Long typeId, @Nullable Long creatorId, @Nullable Long ownerId, @Nullable String name, @Nullable String intro, @Nullable String announcement, @Nullable @Min(0)`
  - [ ] `authAndQueryGroups(@Nullable Set<Long> groupIds, @Nullable String name, @Nullable Date lastUpdatedDate, @Nullable Integer skip, @Nullable Integer limit, @Nullable List<Integer> fieldsToHighlight)`
  - [ ] `queryJoinedGroups(@NotNull Long memberId)`
  - [ ] `queryJoinedGroupIdsWithVersion(@NotNull Long memberId, @Nullable Date lastUpdatedDate)`
  - [ ] `queryJoinedGroupsWithVersion(@NotNull Long memberId, @Nullable Date lastUpdatedDate)`
  - [ ] `isAllowedToCreateGroupAndHaveGroupType(@NotNull Long requesterId, @NotNull Long groupTypeId)`
  - [ ] `isAllowedToCreateGroup(@NotNull Long requesterId, @Nullable UserRole auxiliaryUserRole)`
  - [ ] `isAllowedCreateGroupWithGroupType(@NotNull Long requesterId, @NotNull Long groupTypeId, @Nullable UserRole auxiliaryUserRole)`
  - [ ] `isAllowedUpdateGroupToGroupType(@NotNull Long requesterId, @NotNull Long groupTypeId, @Nullable UserRole auxiliaryUserRole)`
  - [x] `countOwnedGroups(@NotNull Long ownerId)` -> [internal/domain/group/repository/group_repository.go:CountOwnedGroups(ctx context.Context, ownerID int64)](../internal/domain/group/repository/group_repository.go)
  - [x] `countOwnedGroups(@NotNull Long ownerId, @NotNull Long groupTypeId)` -> [internal/domain/group/repository/group_repository.go:CountOwnedGroups(ctx context.Context, ownerID int64)](../internal/domain/group/repository/group_repository.go)
  - [ ] `countCreatedGroups(@Nullable DateRange dateRange)`
  - [x] `countGroups(@Nullable Set<Long> ids, @Nullable Set<Long> typeIds, @Nullable Set<Long> creatorIds, @Nullable Set<Long> ownerIds, @Nullable Boolean isActive, @Nullable DateRange creationDateRange, @Nullable DateRange deletionDateRange, @Nullable DateRange lastUpdatedDateRange, @Nullable DateRange muteEndDateRange, @Nullable Set<Long> memberIds)` -> [internal/domain/user/repository/user_relationship_group_member_repository.go:CountGroups(ctx context.Context, ownerIDs []int64, relatedUserIDs []int64)](../internal/domain/user/repository/user_relationship_group_member_repository.go)
  - [ ] `countDeletedGroups(@Nullable DateRange dateRange)`
  - [x] `count()` -> [internal/domain/user/repository/user_repository.go:Count(ctx context.Context, filter bson.M)](../internal/domain/user/repository/user_repository.go)
  - [ ] `isGroupMuted(@NotNull Long groupId)`
  - [ ] `isGroupActiveAndNotDeleted(@NotNull Long groupId)`

- **GroupTypeService.java** ([java/im/turms/service/domain/group/service/GroupTypeService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/service/GroupTypeService.java))
> [简述功能]

  - [ ] `initGroupTypes()`
  - [ ] `queryGroupTypes(@Nullable Integer page, @Nullable Integer size)`
  - [ ] `addGroupType(@Nullable Long id, @NotNull @NoWhitespace String name, @NotNull @Min(1)`
  - [ ] `updateGroupTypes(@NotEmpty Set<Long> ids, @Nullable @NoWhitespace String name, @Nullable @Min(1)`
  - [ ] `deleteGroupTypes(@Nullable Set<Long> groupTypeIds)`
  - [ ] `queryGroupType(@NotNull Long groupTypeId)`
  - [ ] `queryGroupTypes(@NotNull Collection<Long> groupTypeIds)`
  - [ ] `groupTypeExists(@NotNull Long groupTypeId)`
  - [ ] `countGroupTypes()`

- **GroupVersionService.java** ([java/im/turms/service/domain/group/service/GroupVersionService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/service/GroupVersionService.java))
> [简述功能]

  - [ ] `queryMembersVersion(@NotNull Long groupId)`
  - [ ] `queryBlocklistVersion(@NotNull Long groupId)`
  - [x] `queryGroupJoinRequestsVersion(@NotNull Long groupId)` -> [internal/domain/user/service/user_version_service.go:QueryGroupJoinRequestsVersion(ctx context.Context, userID int64)](../internal/domain/user/service/user_version_service.go)
  - [ ] `queryGroupJoinQuestionsVersion(@NotNull Long groupId)`
  - [x] `queryGroupInvitationsVersion(@NotNull Long groupId)` -> [internal/domain/group/service/group_version_service.go:QueryGroupInvitationsVersion(ctx context.Context, groupID int64)](../internal/domain/group/service/group_version_service.go)
  - [x] `updateVersion(@NotNull Long groupId, boolean updateMembers, boolean updateBlocklist, boolean joinRequests, boolean joinQuestions)` -> [internal/domain/group/repository/group_version_repository.go:UpdateVersion(ctx context.Context, groupID int64, field string)](../internal/domain/group/repository/group_version_repository.go)
  - [x] `updateMembersVersion(@NotNull Long groupId)` -> [internal/domain/group/repository/group_version_repository.go:UpdateMembersVersion(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_version_repository.go)
  - [x] `updateMembersVersion(@Nullable Set<Long> groupIds)` -> [internal/domain/group/repository/group_version_repository.go:UpdateMembersVersion(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_version_repository.go)
  - [x] `updateMembersVersion()` -> [internal/domain/group/repository/group_version_repository.go:UpdateMembersVersion(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_version_repository.go)
  - [x] `updateBlocklistVersion(@NotNull Long groupId)` -> [internal/domain/group/repository/group_version_repository.go:UpdateBlocklistVersion(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_version_repository.go)
  - [x] `updateJoinRequestsVersion(@NotNull Long groupId)` -> [internal/domain/group/repository/group_version_repository.go:UpdateJoinRequestsVersion(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_version_repository.go)
  - [x] `updateJoinQuestionsVersion(@NotNull Long groupId)` -> [internal/domain/group/repository/group_version_repository.go:UpdateJoinQuestionsVersion(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_version_repository.go)
  - [ ] `updateGroupInvitationsVersion(@NotNull Long groupId)`
  - [x] `updateSpecificVersion(@NotNull Long groupId, @NotNull String field)` -> [internal/domain/user/repository/user_version_repository.go:UpdateSpecificVersion(ctx context.Context, userIDs []int64, field string, updateDate time.Time)](../internal/domain/user/repository/user_version_repository.go)
  - [x] `updateSpecificVersion(@NotNull String field)` -> [internal/domain/user/repository/user_version_repository.go:UpdateSpecificVersion(ctx context.Context, userIDs []int64, field string, updateDate time.Time)](../internal/domain/user/repository/user_version_repository.go)
  - [x] `updateSpecificVersion(@Nullable Set<Long> groupIds, @NotNull String field)` -> [internal/domain/user/repository/user_version_repository.go:UpdateSpecificVersion(ctx context.Context, userIDs []int64, field string, updateDate time.Time)](../internal/domain/user/repository/user_version_repository.go)
  - [x] `upsert(@NotNull Long groupId, @NotNull Date timestamp)` -> [internal/domain/group/repository/group_version_repository.go:Upsert(ctx context.Context, groupID int64, timestamp time.Time)](../internal/domain/group/repository/group_version_repository.go)
  - [x] `delete(@Nullable Set<Long> groupIds, @Nullable ClientSession session)` -> [internal/domain/common/cache/sharded_map.go:Delete(key K)](../internal/domain/common/cache/sharded_map.go)

- **MessageController.java** ([java/im/turms/service/domain/message/access/admin/controller/MessageController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/message/access/admin/controller/MessageController.java))
> [简述功能]

  - [ ] `createMessages(@QueryParam(defaultValue = "true")`
  - [x] `queryMessages(@QueryParam(required = false)` -> [internal/domain/message/repository/message_repository.go:QueryMessages(ctx context.Context, isGroupMessage *bool, senderIDs []int64, targetIDs []int64, deliveryDateAfter *time.Time, deliveryDateBefore *time.Time, size int64, ascending bool,)](../internal/domain/message/repository/message_repository.go)
  - [x] `queryMessages(@QueryParam(required = false)` -> [internal/domain/message/repository/message_repository.go:QueryMessages(ctx context.Context, isGroupMessage *bool, senderIDs []int64, targetIDs []int64, deliveryDateAfter *time.Time, deliveryDateBefore *time.Time, size int64, ascending bool,)](../internal/domain/message/repository/message_repository.go)
  - [ ] `countMessages(@QueryParam(required = false)`
  - [ ] `updateMessages(Set<Long> ids, @RequestBody UpdateMessageDTO updateMessageDTO)`
  - [ ] `deleteMessages(Set<Long> ids, @QueryParam(required = false)`

- **CreateMessageDTO.java** ([java/im/turms/service/domain/message/access/admin/dto/request/CreateMessageDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/message/access/admin/dto/request/CreateMessageDTO.java))
> [简述功能]

  - [ ] `CreateMessageDTO(Long id, Boolean isGroupMessage, Boolean isSystemMessage, String text, List<byte[]> records, Long senderId, String senderIp, DeviceType senderDeviceType, Long targetId, Integer burnAfter, Long referenceId, Long preMessageId)`

- **UpdateMessageDTO.java** ([java/im/turms/service/domain/message/access/admin/dto/request/UpdateMessageDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/message/access/admin/dto/request/UpdateMessageDTO.java))
> [简述功能]

  - [ ] `UpdateMessageDTO(Long senderId, String senderIp, DeviceType senderDeviceType, Boolean isSystemMessage, String text, List<byte[]> records, Integer burnAfter, Date recallDate)`

- **MessageStatisticsDTO.java** ([java/im/turms/service/domain/message/access/admin/dto/response/MessageStatisticsDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/message/access/admin/dto/response/MessageStatisticsDTO.java))
> [简述功能]

  - [ ] `MessageStatisticsDTO(Long sentMessagesOnAverage, Long acknowledgedMessages, Long acknowledgedMessagesOnAverage, Long sentMessages, List<StatisticsRecordDTO> sentMessagesOnAverageRecords, List<StatisticsRecordDTO> acknowledgedMessagesRecords, List<StatisticsRecordDTO> acknowledgedMessagesOnAverageRecords, List<StatisticsRecordDTO> sentMessagesRecords)`

- **MessageServiceController.java** ([java/im/turms/service/domain/message/access/servicerequest/controller/MessageServiceController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/message/access/servicerequest/controller/MessageServiceController.java))
> [简述功能]

  - [x] `handleCreateMessageRequest()` -> [internal/domain/message/controller/message_controller.go:HandleCreateMessageRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/message/controller/message_controller.go)
  - [ ] `handleQueryMessagesRequest()`
  - [ ] `handleUpdateMessageRequest()`
  - [ ] `handleCreateMessageReactionsRequest()`
  - [ ] `handleDeleteMessageReactionsRequest()`

- **MessageAndRecipientIds.java** ([java/im/turms/service/domain/message/bo/MessageAndRecipientIds.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/message/bo/MessageAndRecipientIds.java))
> [简述功能]

  - [ ] `MessageAndRecipientIds(Message message, Set<Long> recipientIds)`

- **Message.java** ([java/im/turms/service/domain/message/po/Message.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/message/po/Message.java))
> [简述功能]

  - [ ] `groupId()`

- **MessageRepository.java** ([java/im/turms/service/domain/message/repository/MessageRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/message/repository/MessageRepository.java))
> [简述功能]

  - [ ] `updateMessages(Set<Long> messageIds, @Nullable Boolean isSystemMessage, @Nullable Integer senderIp, @Nullable byte[] senderIpV6, @Nullable Date recallDate, @Nullable String text, @Nullable List<byte[]> records, @Nullable Integer burnAfter, @Nullable ClientSession session)`
  - [ ] `updateMessagesDeletionDate(@Nullable Set<Long> messageIds)`
  - [ ] `existsBySenderIdAndTargetId(Long senderId, Long targetId)`
  - [ ] `countMessages(@Nullable Set<Long> messageIds, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages, @Nullable Set<Long> senderIds, @Nullable Set<Long> targetIds, @Nullable DateRange deliveryDateRange, @Nullable DateRange deletionDateRange)`
  - [ ] `countUsersWhoSentMessage(@Nullable DateRange dateRange, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages)`
  - [ ] `countGroupsThatSentMessages(@Nullable DateRange dateRange)`
  - [ ] `countSentMessages(@Nullable DateRange dateRange, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages)`
  - [ ] `findDeliveryDate(Long messageId)`
  - [ ] `findExpiredMessageIds(Date expirationDate)`
  - [ ] `findMessageGroupId(Long messageId)`
  - [ ] `findMessageSenderIdAndTargetIdAndIsGroupMessage(Long messageId)`
  - [ ] `findMessages(@Nullable Collection<Long> messageIds, @Nullable Collection<byte[]> conversationIds, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages, @Nullable Set<Long> senderIds, @Nullable Set<Long> targetIds, @Nullable DateRange deliveryDateRange, @Nullable DateRange deletionDateRange, @Nullable DateRange recallDateRange, @Nullable Integer page, @Nullable Integer size, @Nullable Boolean ascending)`
  - [ ] `findIsGroupMessageAndTargetId(Long messageId, Long senderId)`
  - [ ] `findIsGroupMessageAndTargetIdAndDeliveryDate(Long messageId, Long senderId)`
  - [ ] `getGroupConversationId(long groupId)`
  - [ ] `getPrivateConversationId(long id1, long id2)`

- **MessageService.java** ([java/im/turms/service/domain/message/service/MessageService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/message/service/MessageService.java))
> [简述功能]

  - [ ] `isMessageRecipientOrSender(@NotNull Long messageId, @NotNull Long userId)`
  - [ ] `authAndQueryCompleteMessages(Long requesterId, @Nullable Collection<Long> messageIds, @NotNull Boolean areGroupMessages, @Nullable Boolean areSystemMessages, @Nullable Set<Long> fromIds, @Nullable DateRange deliveryDateRange, @Nullable Integer maxCount, boolean ascending, boolean withTotal)`
  - [ ] `queryMessage(@NotNull Long messageId)`
  - [x] `queryMessages(@Nullable Collection<Long> messageIds, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages, @Nullable Set<Long> senderIds, @Nullable Set<Long> targetIds, @Nullable DateRange deliveryDateRange, @Nullable DateRange deletionDateRange, @Nullable DateRange recallDateRange, @Nullable Integer page, @Nullable Integer size, @Nullable Boolean ascending)` -> [internal/domain/message/repository/message_repository.go:QueryMessages(ctx context.Context, isGroupMessage *bool, senderIDs []int64, targetIDs []int64, deliveryDateAfter *time.Time, deliveryDateBefore *time.Time, size int64, ascending bool,)](../internal/domain/message/repository/message_repository.go)
  - [ ] `saveMessage(@Nullable Long messageId, @NotNull Long senderId, @Nullable byte[] senderIp, @NotNull Long targetId, @NotNull Boolean isGroupMessage, @NotNull Boolean isSystemMessage, @Nullable String text, @Nullable List<byte[]> records, @Nullable @Min(0)`
  - [ ] `queryExpiredMessageIds(@NotNull Integer retentionPeriodHours)`
  - [ ] `deleteExpiredMessages(@NotNull Integer retentionPeriodHours)`
  - [ ] `deleteMessages(@Nullable Set<Long> messageIds, @Nullable Boolean deleteLogically)`
  - [ ] `updateMessages(@Nullable Long senderId, @Nullable DeviceType senderDeviceType, @NotEmpty Set<Long> messageIds, @Nullable Boolean isSystemMessage, @Nullable String text, @Nullable List<byte[]> records, @Nullable @Min(0)`
  - [ ] `hasPrivateMessage(Long senderId, Long targetId)`
  - [ ] `countMessages(@Nullable Set<Long> messageIds, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages, @Nullable Set<Long> senderIds, @Nullable Set<Long> targetIds, @Nullable DateRange deliveryDateRange, @Nullable DateRange deletionDateRange)`
  - [ ] `countUsersWhoSentMessage(@Nullable DateRange dateRange, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages)`
  - [ ] `countGroupsThatSentMessages(@Nullable DateRange dateRange)`
  - [ ] `countSentMessages(@Nullable DateRange dateRange, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages)`
  - [ ] `countSentMessagesOnAverage(@Nullable DateRange dateRange, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages)`
  - [ ] `authAndUpdateMessage(@NotNull Long senderId, @Nullable DeviceType senderDeviceType, @NotNull Long messageId, @Nullable String text, @Nullable List<byte[]> records, @Nullable @PastOrPresent Date recallDate)`
  - [ ] `queryMessageRecipients(@NotNull Long messageId)`
  - [x] `authAndSaveMessage(boolean queryRecipientIds, @Nullable Boolean persist, @Nullable Long messageId, @NotNull Long senderId, @Nullable byte[] senderIp, @NotNull Long targetId, @NotNull Boolean isGroupMessage, @NotNull Boolean isSystemMessage, @Nullable String text, @Nullable List<byte[]> records, @Nullable @Min(0)` -> [internal/domain/message/service/message_service.go:AuthAndSaveMessage(ctx context.Context, isGroupMessage bool, senderID int64, targetID int64, text string)](../internal/domain/message/service/message_service.go)
  - [ ] `saveMessage(boolean queryRecipientIds, @Nullable Boolean persist, @Nullable Long messageId, @NotNull Long senderId, @Nullable byte[] senderIp, @NotNull Long targetId, @NotNull Boolean isGroupMessage, @NotNull Boolean isSystemMessage, @Nullable String text, @Nullable List<byte[]> records, @Nullable @Min(0)`
  - [ ] `authAndCloneAndSaveMessage(boolean queryRecipientIds, @NotNull Long requesterId, @Nullable byte[] requesterIp, @NotNull Long referenceId, @NotNull Boolean isGroupMessage, @NotNull Boolean isSystemMessage, @NotNull Long targetId)`
  - [ ] `cloneAndSaveMessage(boolean queryRecipientIds, @NotNull Long senderId, @Nullable byte[] senderIp, @NotNull Long referenceId, @NotNull Boolean isGroupMessage, @NotNull Boolean isSystemMessage, @NotNull Long targetId)`
  - [x] `authAndSaveAndSendMessage(boolean send, @Nullable Boolean persist, @Nullable Long senderId, @Nullable DeviceType senderDeviceType, @Nullable byte[] senderIp, @Nullable Long messageId, @NotNull Boolean isGroupMessage, @NotNull Boolean isSystemMessage, @Nullable String text, @Nullable List<byte[]> records, @NotNull Long targetId, @Nullable @Min(0)` -> [internal/domain/message/service/message_service.go:AuthAndSaveAndSendMessage(ctx context.Context, isGroupMessage bool, senderID int64, targetID int64, text string)](../internal/domain/message/service/message_service.go)
  - [ ] `saveAndSendMessage(boolean send, @Nullable Boolean persist, @Nullable Long senderId, @Nullable DeviceType senderDeviceType, @Nullable byte[] senderIp, @Nullable Long messageId, @NotNull Boolean isGroupMessage, @NotNull Boolean isSystemMessage, @Nullable String text, @Nullable List<byte[]> records, @NotNull Long targetId, @Nullable @Min(0)`
  - [ ] `saveAndSendMessage(@Nullable Boolean persist, @Nullable Long senderId, @Nullable DeviceType senderDeviceType, @Nullable byte[] senderIp, @Nullable Long messageId, @NotNull Boolean isGroupMessage, @NotNull Boolean isSystemMessage, @Nullable String text, @Nullable List<byte[]> records, @NotNull Long targetId, @Nullable @Min(0)`
  - [ ] `deleteGroupMessageSequenceIds(Set<Long> groupIds)`
  - [ ] `deletePrivateMessageSequenceIds(Set<Long> userIds)`
  - [ ] `fetchGroupMessageSequenceId(Long groupId)`
  - [ ] `fetchPrivateMessageSequenceId(Long userId1, Long userId2)`

- **StatisticsService.java** ([java/im/turms/service/domain/observation/service/StatisticsService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/observation/service/StatisticsService.java))
> [简述功能]

  - [ ] `countOnlineUsersByNodes()`
  - [x] `countOnlineUsers()` -> [internal/domain/gateway/session/service.go:CountOnlineUsers()](../internal/domain/gateway/session/service.go)

- **StorageServiceController.java** ([java/im/turms/service/domain/storage/access/servicerequest/controller/StorageServiceController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/storage/access/servicerequest/controller/StorageServiceController.java))
> [简述功能]

  - [ ] `handleDeleteResourceRequest()`
  - [ ] `handleQueryResourceUploadInfoRequest()`
  - [ ] `handleQueryResourceDownloadInfoRequest()`
  - [ ] `handleUpdateMessageAttachmentInfoRequest()`
  - [ ] `handleQueryMessageAttachmentInfosRequest()`

- **StorageResourceInfo.java** ([java/im/turms/service/domain/storage/bo/StorageResourceInfo.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/storage/bo/StorageResourceInfo.java))
> [简述功能]

  - [ ] `StorageResourceInfo(@Nullable Long idNum, @Nullable String idStr, String name, String mediaType, Long uploaderId, Date creationDate)`

- **StorageService.java** ([java/im/turms/service/domain/storage/service/StorageService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/storage/service/StorageService.java))
> [简述功能]

  - [x] `deleteResource(Long requesterId, StorageResourceType resourceType, @Nullable Long resourceIdNum, @Nullable String resourceIdStr, List<Value> customAttributes)` -> [internal/domain/storage/provider/mock_storage_provider.go:DeleteResource(ctx context.Context, resourceType constants.StorageResourceType, keyStr string)](../internal/domain/storage/provider/mock_storage_provider.go)
  - [x] `queryResourceUploadInfo(Long requesterId, StorageResourceType resourceType, @Nullable Long resourceIdNum, @Nullable String resourceName, @Nullable String resourceMediaType, List<Value> customAttributes)` -> [internal/domain/storage/service/storage_service.go:QueryResourceUploadInfo(ctx context.Context, requesterID int64, resourceType constants.StorageResourceType, resourceName string, contentType string, maxSize int64,)](../internal/domain/storage/service/storage_service.go)
  - [x] `queryResourceDownloadInfo(Long requesterId, StorageResourceType resourceType, @Nullable Long resourceIdNum, @Nullable String resourceIdStr, List<Value> customAttributes)` -> [internal/domain/storage/service/storage_service.go:QueryResourceDownloadInfo(ctx context.Context, requesterID int64, resourceType constants.StorageResourceType, resourceIDStr string,)](../internal/domain/storage/service/storage_service.go)
  - [ ] `shareMessageAttachmentWithUser(Long requesterId, @Nullable Long messageAttachmentIdNum, @Nullable String messageAttachmentIdStr, Long userIdToShareWith)`
  - [ ] `shareMessageAttachmentWithGroup(Long requesterId, @Nullable Long messageAttachmentIdNum, @Nullable String messageAttachmentIdStr, Long groupIdToShareWith)`
  - [ ] `unshareMessageAttachmentWithUser(Long requesterId, @Nullable Long messageAttachmentIdNum, @Nullable String messageAttachmentIdStr, Long userIdToUnshareWith)`
  - [ ] `unshareMessageAttachmentWithGroup(Long requesterId, @Nullable Long messageAttachmentIdNum, @Nullable String messageAttachmentIdStr, Long groupIdToUnshareWith)`
  - [ ] `queryMessageAttachmentInfosUploadedByRequester(Long requesterId, @Nullable DateRange creationDateRange)`
  - [ ] `queryMessageAttachmentInfosInPrivateConversations(Long requesterId, @Nullable Set<Long> userIds, @Nullable DateRange creationDateRange, @Nullable Boolean areSharedByRequester)`
  - [ ] `queryMessageAttachmentInfosInGroupConversations(Long requesterId, @Nullable Set<Long> groupIds, @Nullable Set<Long> userIds, @Nullable DateRange creationDateRange)`

- **UserController.java** ([java/im/turms/service/domain/user/access/admin/controller/UserController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/access/admin/controller/UserController.java))
> [简述功能]

  - [x] `addUser(@RequestBody AddUserDTO addUserDTO)` -> [internal/domain/user/service/user_service.go:AddUser(ctx context.Context, id int64, password string, name string, intro string, profilePicture string, profileAccess int32, permissionGroupID int64, registrationDate time.Time, isActive bool)](../internal/domain/user/service/user_service.go)
  - [x] `queryUsers(@QueryParam(required = false)` -> [internal/domain/user/service/user_service.go:QueryUsers(ctx context.Context, userIDs []int64)](../internal/domain/user/service/user_service.go)
  - [x] `queryUsers(@QueryParam(required = false)` -> [internal/domain/user/service/user_service.go:QueryUsers(ctx context.Context, userIDs []int64)](../internal/domain/user/service/user_service.go)
  - [x] `countUsers(@QueryParam(required = false)` -> [internal/domain/user/repository/user_repository.go:CountUsers(ctx context.Context, startDate *time.Time, endDate *time.Time)](../internal/domain/user/repository/user_repository.go)
  - [x] `updateUser(Set<Long> ids, @RequestBody UpdateUserDTO updateUserDTO)` -> [internal/domain/user/service/user_service.go:UpdateUser(ctx context.Context, userID int64, update bson.M)](../internal/domain/user/service/user_service.go)
  - [x] `deleteUsers(Set<Long> ids, @QueryParam(required = false)` -> [internal/domain/user/service/user_service.go:DeleteUsers(ctx context.Context, userIDs []int64)](../internal/domain/user/service/user_service.go)

- **UserOnlineInfoController.java** ([java/im/turms/service/domain/user/access/admin/controller/UserOnlineInfoController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/access/admin/controller/UserOnlineInfoController.java))
> [简述功能]

  - [x] `countOnlineUsers(boolean countByNodes)` -> [internal/domain/gateway/session/service.go:CountOnlineUsers()](../internal/domain/gateway/session/service.go)
  - [x] `queryUserSessions(Set<Long> ids, boolean returnNonExistingUsers)` -> [internal/domain/user/service/onlineuser/session_service.go:QueryUserSessions(ctx context.Context, userIDs []int64)](../internal/domain/user/service/onlineuser/session_service.go)
  - [ ] `queryUserStatuses(Set<Long> ids, boolean returnNonExistingUsers)`
  - [ ] `queryUsersNearby(Long userId, @QueryParam(required = false)`
  - [ ] `queryUserLocations(Set<Long> ids, @QueryParam(required = false)`
  - [ ] `updateUserOnlineStatus(Set<Long> ids, @QueryParam(required = false)`

- **UserRoleController.java** ([java/im/turms/service/domain/user/access/admin/controller/UserRoleController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/access/admin/controller/UserRoleController.java))
> [简述功能]

  - [x] `addUserRole(@RequestBody AddUserRoleDTO addUserRoleDTO)` -> [internal/domain/user/service/user_role_service.go:AddUserRole(ctx context.Context, role *po.UserRole)](../internal/domain/user/service/user_role_service.go)
  - [x] `queryUserRoles(@QueryParam(required = false)` -> [internal/domain/user/service/user_role_service.go:QueryUserRoles(ctx context.Context, filter bson.M)](../internal/domain/user/service/user_role_service.go)
  - [ ] `queryUserRoleGroups(int page, @QueryParam(required = false)`
  - [ ] `updateUserRole(Set<Long> ids, @RequestBody UpdateUserRoleDTO updateUserRoleDTO)`
  - [ ] `deleteUserRole(Set<Long> ids)`

- **UserFriendRequestController.java** ([java/im/turms/service/domain/user/access/admin/controller/relationship/UserFriendRequestController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/access/admin/controller/relationship/UserFriendRequestController.java))
> [简述功能]

  - [x] `createFriendRequest(@RequestBody AddFriendRequestDTO addFriendRequestDTO)` -> [internal/domain/user/service/user_friend_request_service.go:CreateFriendRequest(ctx context.Context, requestID *int64, requesterID, recipientID int64, content string, status *po.RequestStatus, creationDate, responseDate *time.Time, reason *string)](../internal/domain/user/service/user_friend_request_service.go)
  - [x] `queryFriendRequests(@QueryParam(required = false)` -> [internal/domain/user/service/user_friend_request_service.go:QueryFriendRequests(ctx context.Context, ids, requesterIds, recipientIds []int64, statuses []po.RequestStatus, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd *time.Time, page, size *int)](../internal/domain/user/service/user_friend_request_service.go)
  - [x] `queryFriendRequests(@QueryParam(required = false)` -> [internal/domain/user/service/user_friend_request_service.go:QueryFriendRequests(ctx context.Context, ids, requesterIds, recipientIds []int64, statuses []po.RequestStatus, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd *time.Time, page, size *int)](../internal/domain/user/service/user_friend_request_service.go)
  - [x] `updateFriendRequests(Set<Long> ids, @RequestBody UpdateFriendRequestDTO updateFriendRequestDTO)` -> [internal/domain/user/repository/user_friend_request_repository.go:UpdateFriendRequests(ctx context.Context, requestIds []int64, requesterID, recipientID *int64, content *string, status *po.RequestStatus, reason *string, creationDate *time.Time)](../internal/domain/user/repository/user_friend_request_repository.go)
  - [x] `deleteFriendRequests(@QueryParam(required = false)` -> [internal/domain/user/service/user_friend_request_service.go:DeleteFriendRequests(ctx context.Context, ids []int64)](../internal/domain/user/service/user_friend_request_service.go)

- **UserRelationshipController.java** ([java/im/turms/service/domain/user/access/admin/controller/relationship/UserRelationshipController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/access/admin/controller/relationship/UserRelationshipController.java))
> [简述功能]

  - [ ] `addRelationship(@RequestBody AddRelationshipDTO addRelationshipDTO)`
  - [x] `queryRelationships(@QueryParam(required = false)` -> [internal/domain/user/service/user_relationship_service.go:QueryRelationships(ctx context.Context, ownerIDs []int64, relatedUserIDs []int64, groupIndexes []int32, isBlocked *bool, establishmentDateRange *turmsmongo.DateRange, page *int, size *int,)](../internal/domain/user/service/user_relationship_service.go)
  - [x] `queryRelationships(@QueryParam(required = false)` -> [internal/domain/user/service/user_relationship_service.go:QueryRelationships(ctx context.Context, ownerIDs []int64, relatedUserIDs []int64, groupIndexes []int32, isBlocked *bool, establishmentDateRange *turmsmongo.DateRange, page *int, size *int,)](../internal/domain/user/service/user_relationship_service.go)
  - [ ] `updateRelationships(List<UserRelationship.Key> keys, @RequestBody UpdateRelationshipDTO updateRelationshipDTO)`
  - [ ] `deleteRelationships(List<UserRelationship.Key> keys)`

- **UserRelationshipGroupController.java** ([java/im/turms/service/domain/user/access/admin/controller/relationship/UserRelationshipGroupController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/access/admin/controller/relationship/UserRelationshipGroupController.java))
> [简述功能]

  - [ ] `addRelationshipGroup(@RequestBody AddRelationshipGroupDTO addRelationshipGroupDTO)`
  - [x] `deleteRelationshipGroups(@QueryParam(required = false)` -> [internal/domain/user/repository/user_relationship_group_repository.go:DeleteRelationshipGroups(ctx context.Context, ownerID int64, groupIndexes []int32, session *mongo.Session)](../internal/domain/user/repository/user_relationship_group_repository.go)
  - [x] `updateRelationshipGroups(List<UserRelationshipGroup.Key> keys, @RequestBody UpdateRelationshipGroupDTO updateRelationshipGroupDTO)` -> [internal/domain/user/repository/user_relationship_group_repository.go:UpdateRelationshipGroups(ctx context.Context, keys []po.UserRelationshipGroupKey, newName string, session *mongo.Session)](../internal/domain/user/repository/user_relationship_group_repository.go)
  - [x] `queryRelationshipGroups(@QueryParam(required = false)` -> [internal/domain/user/service/user_relationship_group_service.go:QueryRelationshipGroups(ctx context.Context, ownerIDs []int64, groupIndexes []int32, page *int, size *int,)](../internal/domain/user/service/user_relationship_group_service.go)
  - [x] `queryRelationshipGroups(@QueryParam(required = false)` -> [internal/domain/user/service/user_relationship_group_service.go:QueryRelationshipGroups(ctx context.Context, ownerIDs []int64, groupIndexes []int32, page *int, size *int,)](../internal/domain/user/service/user_relationship_group_service.go)

- **AddFriendRequestDTO.java** ([java/im/turms/service/domain/user/access/admin/dto/request/AddFriendRequestDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/access/admin/dto/request/AddFriendRequestDTO.java))
> [简述功能]

  - [ ] `AddFriendRequestDTO(Long id, Long requesterId, Long recipientId, String content, RequestStatus status, String reason, Date creationDate, Date responseDate)`

- **AddRelationshipDTO.java** ([java/im/turms/service/domain/user/access/admin/dto/request/AddRelationshipDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/access/admin/dto/request/AddRelationshipDTO.java))
> [简述功能]

  - [ ] `AddRelationshipDTO(Long ownerId, Long relatedUserId, String name, Date blockDate, Date establishmentDate)`

- **AddRelationshipGroupDTO.java** ([java/im/turms/service/domain/user/access/admin/dto/request/AddRelationshipGroupDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/access/admin/dto/request/AddRelationshipGroupDTO.java))
> [简述功能]

  - [ ] `AddRelationshipGroupDTO(Long ownerId, Integer index, String name, Date creationDate)`

- **AddUserDTO.java** ([java/im/turms/service/domain/user/access/admin/dto/request/AddUserDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/access/admin/dto/request/AddUserDTO.java))
> [简述功能]

  - [ ] `AddUserDTO(Long id, @SensitiveProperty(SensitiveProperty.Access.ALLOW_DESERIALIZATION)`
  - [x] `toString()` -> [internal/domain/gateway/session/connection.go:ToString()](../internal/domain/gateway/session/connection.go)

- **AddUserRoleDTO.java** ([java/im/turms/service/domain/user/access/admin/dto/request/AddUserRoleDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/access/admin/dto/request/AddUserRoleDTO.java))
> [简述功能]

  - [ ] `AddUserRoleDTO(Long id, String name, Set<Long> creatableGroupTypeIds, Integer ownedGroupLimit, Integer ownedGroupLimitForEachGroupType, Map<Long, Integer> groupTypeIdToLimit)`

- **UpdateFriendRequestDTO.java** ([java/im/turms/service/domain/user/access/admin/dto/request/UpdateFriendRequestDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/access/admin/dto/request/UpdateFriendRequestDTO.java))
> [简述功能]

  - [ ] `UpdateFriendRequestDTO(Long requesterId, Long recipientId, String content, RequestStatus status, String reason, Date creationDate, Date responseDate)`

- **UpdateOnlineStatusDTO.java** ([java/im/turms/service/domain/user/access/admin/dto/request/UpdateOnlineStatusDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/access/admin/dto/request/UpdateOnlineStatusDTO.java))
> [简述功能]

  - [ ] `UpdateOnlineStatusDTO(UserStatus onlineStatus)`

- **UpdateRelationshipDTO.java** ([java/im/turms/service/domain/user/access/admin/dto/request/UpdateRelationshipDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/access/admin/dto/request/UpdateRelationshipDTO.java))
> [简述功能]

  - [ ] `UpdateRelationshipDTO(String name, Date blockDate, Date establishmentDate)`

- **UpdateRelationshipGroupDTO.java** ([java/im/turms/service/domain/user/access/admin/dto/request/UpdateRelationshipGroupDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/access/admin/dto/request/UpdateRelationshipGroupDTO.java))
> [简述功能]

  - [ ] `UpdateRelationshipGroupDTO(String name, Date creationDate)`

- **UpdateUserDTO.java** ([java/im/turms/service/domain/user/access/admin/dto/request/UpdateUserDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/access/admin/dto/request/UpdateUserDTO.java))
> [简述功能]

  - [ ] `UpdateUserDTO(@SensitiveProperty(SensitiveProperty.Access.ALLOW_DESERIALIZATION)`
  - [x] `toString()` -> [internal/domain/gateway/session/connection.go:ToString()](../internal/domain/gateway/session/connection.go)

- **UpdateUserRoleDTO.java** ([java/im/turms/service/domain/user/access/admin/dto/request/UpdateUserRoleDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/access/admin/dto/request/UpdateUserRoleDTO.java))
> [简述功能]

  - [ ] `UpdateUserRoleDTO(String name, Set<Long> creatableGroupTypeIds, Integer ownedGroupLimit, Integer ownedGroupLimitForEachGroupType, Map<Long, Integer> groupTypeIdToLimit)`

- **OnlineUserCountDTO.java** ([java/im/turms/service/domain/user/access/admin/dto/response/OnlineUserCountDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/access/admin/dto/response/OnlineUserCountDTO.java))
> [简述功能]

  - [ ] `OnlineUserCountDTO(Integer total, Map<String, Integer> nodeIdToUserCount)`

- **UserFriendRequestDTO.java** ([java/im/turms/service/domain/user/access/admin/dto/response/UserFriendRequestDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/access/admin/dto/response/UserFriendRequestDTO.java))
> [简述功能]

  - [ ] `UserFriendRequestDTO(Long id, String content, RequestStatus status, String reason, Date creationDate, Date responseDate, Long requesterId, Long recipientId, Date expirationDate)`

- **UserLocationDTO.java** ([java/im/turms/service/domain/user/access/admin/dto/response/UserLocationDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/access/admin/dto/response/UserLocationDTO.java))
> [简述功能]

  - [ ] `UserLocationDTO(Long userId, DeviceType deviceType, Double longitude, Double latitude)`

- **UserRelationshipDTO.java** ([java/im/turms/service/domain/user/access/admin/dto/response/UserRelationshipDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/access/admin/dto/response/UserRelationshipDTO.java))
> [简述功能]

  - [ ] `UserRelationshipDTO(Key key, String name, Date blockDate, Date establishmentDate, Set<Integer> groupIndexes)`
  - [ ] `fromDomain(UserRelationship relationship)`
  - [ ] `fromDomain(@NotNull UserRelationship relationship, @Nullable Set<Integer> groupIndexes)`
  - [ ] `Key(Long ownerId, Long relatedUserId)`

- **UserStatisticsDTO.java** ([java/im/turms/service/domain/user/access/admin/dto/response/UserStatisticsDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/access/admin/dto/response/UserStatisticsDTO.java))
> [简述功能]

  - [ ] `UserStatisticsDTO(Long deletedUsers, Long usersWhoSentMessages, Long loggedInUsers, Long maxOnlineUsers, Long registeredUsers, List<StatisticsRecordDTO> deletedUsersRecords, List<StatisticsRecordDTO> usersWhoSentMessagesRecords, List<StatisticsRecordDTO> loggedInUsersRecords, List<StatisticsRecordDTO> maxOnlineUsersRecords, List<StatisticsRecordDTO> registeredUsersRecords)`

- **UserStatusDTO.java** ([java/im/turms/service/domain/user/access/admin/dto/response/UserStatusDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/access/admin/dto/response/UserStatusDTO.java))
> [简述功能]

  - [ ] `UserStatusDTO(Long userId, UserStatus status, Map<DeviceType, String> deviceTypeToNodeId, Date loginDate, Location loginLocation)`

- **UserRelationshipServiceController.java** ([java/im/turms/service/domain/user/access/servicerequest/controller/UserRelationshipServiceController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/access/servicerequest/controller/UserRelationshipServiceController.java))
> [简述功能]

  - [x] `handleCreateFriendRequestRequest()` -> [internal/domain/user/controller/user_relationship_controller.go:HandleCreateFriendRequestRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_relationship_controller.go)
  - [x] `handleCreateRelationshipGroupRequest()` -> [internal/domain/user/controller/user_relationship_controller.go:HandleCreateRelationshipGroupRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_relationship_controller.go)
  - [x] `handleCreateRelationshipRequest()` -> [internal/domain/user/controller/user_relationship_controller.go:HandleCreateRelationshipRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_relationship_controller.go)
  - [x] `handleDeleteFriendRequestRequest()` -> [internal/domain/user/controller/user_relationship_controller.go:HandleDeleteFriendRequestRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_relationship_controller.go)
  - [x] `handleDeleteRelationshipGroupRequest()` -> [internal/domain/user/controller/user_relationship_controller.go:HandleDeleteRelationshipGroupRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_relationship_controller.go)
  - [x] `handleDeleteRelationshipRequest()` -> [internal/domain/user/controller/user_relationship_controller.go:HandleDeleteRelationshipRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_relationship_controller.go)
  - [x] `handleQueryFriendRequestsRequest()` -> [internal/domain/user/controller/user_relationship_controller.go:HandleQueryFriendRequestsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_relationship_controller.go)
  - [x] `handleQueryRelatedUserIdsRequest()` -> [internal/domain/user/controller/user_relationship_controller.go:HandleQueryRelatedUserIdsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_relationship_controller.go)
  - [x] `handleQueryRelationshipGroupsRequest()` -> [internal/domain/user/controller/user_relationship_controller.go:HandleQueryRelationshipGroupsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_relationship_controller.go)
  - [x] `handleQueryRelationshipsRequest()` -> [internal/domain/user/controller/user_relationship_controller.go:HandleQueryRelationshipsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_relationship_controller.go)
  - [x] `handleUpdateFriendRequestRequest()` -> [internal/domain/user/controller/user_relationship_controller.go:HandleUpdateFriendRequestRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_relationship_controller.go)
  - [x] `handleUpdateRelationshipGroupRequest()` -> [internal/domain/user/controller/user_relationship_controller.go:HandleUpdateRelationshipGroupRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_relationship_controller.go)
  - [x] `handleUpdateRelationshipRequest()` -> [internal/domain/user/controller/user_relationship_controller.go:HandleUpdateRelationshipRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_relationship_controller.go)

- **UserServiceController.java** ([java/im/turms/service/domain/user/access/servicerequest/controller/UserServiceController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/access/servicerequest/controller/UserServiceController.java))
> [简述功能]

  - [x] `handleQueryUserProfilesRequest()` -> [internal/domain/user/controller/user_service_controller.go:HandleQueryUserProfilesRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_service_controller.go)
  - [x] `handleQueryNearbyUsersRequest()` -> [internal/domain/user/controller/user_service_controller.go:HandleQueryNearbyUsersRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_service_controller.go)
  - [x] `handleQueryUserOnlineStatusesRequest()` -> [internal/domain/user/controller/user_service_controller.go:HandleQueryUserOnlineStatusesRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_service_controller.go)
  - [x] `handleUpdateUserLocationRequest()` -> [internal/domain/user/controller/user_service_controller.go:HandleUpdateUserLocationRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_service_controller.go)
  - [x] `handleUpdateUserOnlineStatusRequest()` -> [internal/domain/user/controller/user_service_controller.go:HandleUpdateUserOnlineStatusRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_service_controller.go)
  - [x] `handleUpdateUserRequest()` -> [internal/domain/user/controller/user_service_controller.go:HandleUpdateUserRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_service_controller.go)

- **UserSettingsServiceController.java** ([java/im/turms/service/domain/user/access/servicerequest/controller/UserSettingsServiceController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/access/servicerequest/controller/UserSettingsServiceController.java))
> [简述功能]

  - [ ] `handleDeleteUserSettingsRequest()`
  - [x] `handleUpdateUserSettingsRequest()` -> [internal/domain/user/controller/user_settings_controller.go:HandleUpdateUserSettingsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_settings_controller.go)
  - [x] `handleQueryUserSettingsRequest()` -> [internal/domain/user/controller/user_settings_controller.go:HandleQueryUserSettingsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_settings_controller.go)

- **HandleFriendRequestResult.java** ([java/im/turms/service/domain/user/bo/HandleFriendRequestResult.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/bo/HandleFriendRequestResult.java))
> [简述功能]

  - [ ] `HandleFriendRequestResult(UserFriendRequest friendRequest, @Nullable Integer newGroupIndexForFriendRequestRequester, @Nullable Integer newGroupIndexForFriendRequestRecipient)`

- **UpsertRelationshipResult.java** ([java/im/turms/service/domain/user/bo/UpsertRelationshipResult.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/bo/UpsertRelationshipResult.java))
> [简述功能]

  - [ ] `UpsertRelationshipResult(boolean createdNewRelationship, @Nullable Integer newRelationshipGroupIndex)`

- **UserFriendRequestRepository.java** ([java/im/turms/service/domain/user/repository/UserFriendRequestRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/repository/UserFriendRequestRepository.java))
> [简述功能]

  - [x] `getEntityExpireAfterSeconds()` -> [internal/domain/user/repository/user_friend_request_repository.go:GetEntityExpireAfterSeconds()](../internal/domain/user/repository/user_friend_request_repository.go)
  - [x] `updateFriendRequests(Set<Long> requestIds, @Nullable Long requesterId, @Nullable Long recipientId, @Nullable String content, @Nullable RequestStatus status, @Nullable String reason, @Nullable Date creationDate)` -> [internal/domain/user/repository/user_friend_request_repository.go:UpdateFriendRequests(ctx context.Context, requestIds []int64, requesterID, recipientID *int64, content *string, status *po.RequestStatus, reason *string, creationDate *time.Time)](../internal/domain/user/repository/user_friend_request_repository.go)
  - [x] `updateStatusIfPending(Long requestId, RequestStatus requestStatus, @Nullable String reason, @Nullable ClientSession session)` -> [internal/domain/group/repository/group_invitation_repository.go:UpdateStatusIfPending(ctx context.Context, invitationID int64, newStatus po.RequestStatus, reason *string, responseDate time.Time)](../internal/domain/group/repository/group_invitation_repository.go)
  - [x] `countFriendRequests(@Nullable Set<Long> ids, @Nullable Set<Long> requesterIds, @Nullable Set<Long> recipientIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange)` -> [internal/domain/user/repository/user_friend_request_repository.go:CountFriendRequests(ctx context.Context, ids, requesterIds, recipientIds []int64, statuses []po.RequestStatus, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd *time.Time)](../internal/domain/user/repository/user_friend_request_repository.go)
  - [x] `findFriendRequests(@Nullable Set<Long> ids, @Nullable Set<Long> requesterIds, @Nullable Set<Long> recipientIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange, @Nullable Integer page, @Nullable Integer size)` -> [internal/domain/user/repository/user_friend_request_repository.go:FindFriendRequests(ctx context.Context, ids, requesterIds, recipientIds []int64, statuses []po.RequestStatus, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd *time.Time, page, size *int)](../internal/domain/user/repository/user_friend_request_repository.go)
  - [x] `findFriendRequestsByRecipientId(Long recipientId)` -> [internal/domain/user/repository/user_friend_request_repository.go:FindFriendRequestsByRecipientId(ctx context.Context, recipientID int64)](../internal/domain/user/repository/user_friend_request_repository.go)
  - [x] `findFriendRequestsByRequesterId(Long requesterId)` -> [internal/domain/user/repository/user_friend_request_repository.go:FindFriendRequestsByRequesterId(ctx context.Context, requesterID int64)](../internal/domain/user/repository/user_friend_request_repository.go)
  - [x] `findRecipientId(Long requestId)` -> [internal/domain/user/repository/user_friend_request_repository.go:FindRecipientId(ctx context.Context, requestID int64)](../internal/domain/user/repository/user_friend_request_repository.go)
  - [x] `findRequesterIdAndRecipientIdAndStatus(Long requestId)` -> [internal/domain/user/repository/user_friend_request_repository.go:FindRequesterIdAndRecipientIdAndStatus(ctx context.Context, requestID int64)](../internal/domain/user/repository/user_friend_request_repository.go)
  - [x] `findRequesterIdAndRecipientIdAndCreationDateAndStatus(Long requestId)` -> [internal/domain/user/repository/user_friend_request_repository.go:FindRequesterIdAndRecipientIdAndCreationDateAndStatus(ctx context.Context, requestID int64)](../internal/domain/user/repository/user_friend_request_repository.go)
  - [x] `hasPendingFriendRequest(Long requesterId, Long recipientId)` -> [internal/domain/user/repository/user_friend_request_repository.go:HasPendingFriendRequest(ctx context.Context, requesterID, recipientID int64)](../internal/domain/user/repository/user_friend_request_repository.go)
  - [x] `hasPendingOrDeclinedOrIgnoredOrExpiredRequest(Long requesterId, Long recipientId)` -> [internal/domain/user/repository/user_friend_request_repository.go:HasPendingOrDeclinedOrIgnoredOrExpiredRequest(ctx context.Context, requesterID, recipientID int64)](../internal/domain/user/repository/user_friend_request_repository.go)

- **UserRelationshipGroupMemberRepository.java** ([java/im/turms/service/domain/user/repository/UserRelationshipGroupMemberRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/repository/UserRelationshipGroupMemberRepository.java))
> [简述功能]

  - [x] `deleteAllRelatedUserFromRelationshipGroup(Long ownerId, Integer groupIndex)` -> [internal/domain/user/repository/user_relationship_group_member_repository.go:DeleteAllRelatedUserFromRelationshipGroup(ctx context.Context, ownerID int64, session *mongo.Session)](../internal/domain/user/repository/user_relationship_group_member_repository.go)
  - [x] `deleteRelatedUserFromRelationshipGroup(Long ownerId, Long relatedUserId, Integer groupIndex, @Nullable ClientSession session)` -> [internal/domain/user/repository/user_relationship_group_member_repository.go:DeleteRelatedUserFromRelationshipGroup(ctx context.Context, ownerID int64, relatedUserID int64, groupIndexes []int32, session *mongo.Session)](../internal/domain/user/repository/user_relationship_group_member_repository.go)
  - [x] `deleteRelatedUsersFromAllRelationshipGroups(Long ownerId, Collection<Long> relatedUserIds, @Nullable ClientSession session)` -> [internal/domain/user/repository/user_relationship_group_member_repository.go:DeleteRelatedUsersFromAllRelationshipGroups(ctx context.Context, ownerID int64, relatedUserIDs []int64, session *mongo.Session)](../internal/domain/user/repository/user_relationship_group_member_repository.go)
  - [x] `countGroups(@Nullable Collection<Long> ownerIds, @Nullable Collection<Long> relatedUserIds)` -> [internal/domain/user/repository/user_relationship_group_member_repository.go:CountGroups(ctx context.Context, ownerIDs []int64, relatedUserIDs []int64)](../internal/domain/user/repository/user_relationship_group_member_repository.go)
  - [x] `countMembers(@Nullable Set<Long> ownerIds, @Nullable Set<Integer> groupIndexes)` -> [internal/domain/group/repository/group_member_repository.go:CountMembers(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_member_repository.go)
  - [x] `findGroupIndexes(Long ownerId, Long relatedUserId)` -> [internal/domain/user/repository/user_relationship_group_member_repository.go:FindGroupIndexes(ctx context.Context, ownerID int64, relatedUserID int64)](../internal/domain/user/repository/user_relationship_group_member_repository.go)
  - [x] `findRelationshipGroupMemberIds(Long ownerId, Integer groupIndex)` -> [internal/domain/user/repository/user_relationship_group_member_repository.go:FindRelationshipGroupMemberIds(ctx context.Context, ownerIDs []int64, groupIndexes []int32, page *int, size *int)](../internal/domain/user/repository/user_relationship_group_member_repository.go)
  - [x] `findRelationshipGroupMemberIds(@Nullable Set<Long> ownerIds, @Nullable Set<Integer> groupIndexes, @Nullable Integer page, @Nullable Integer size)` -> [internal/domain/user/repository/user_relationship_group_member_repository.go:FindRelationshipGroupMemberIds(ctx context.Context, ownerIDs []int64, groupIndexes []int32, page *int, size *int)](../internal/domain/user/repository/user_relationship_group_member_repository.go)
  - [x] `findRelationshipGroupMembers(Long ownerId, Integer groupIndex)` -> [internal/domain/user/repository/user_relationship_group_member_repository.go:FindRelationshipGroupMembers(ctx context.Context, ownerID int64, groupIndex int32)](../internal/domain/user/repository/user_relationship_group_member_repository.go)

- **UserRelationshipGroupRepository.java** ([java/im/turms/service/domain/user/repository/UserRelationshipGroupRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/repository/UserRelationshipGroupRepository.java))
> [简述功能]

  - [x] `deleteAllRelationshipGroups(Set<Long> ownerIds, @Nullable ClientSession session)` -> [internal/domain/user/repository/user_relationship_group_repository.go:DeleteAllRelationshipGroups(ctx context.Context, ownerIDs []int64, session *mongo.Session)](../internal/domain/user/repository/user_relationship_group_repository.go)
  - [x] `updateRelationshipGroupName(Long ownerId, Integer groupIndex, String newGroupName)` -> [internal/domain/user/repository/user_relationship_group_repository.go:UpdateRelationshipGroupName(ctx context.Context, ownerID int64, groupIndex int32, newName string, session *mongo.Session)](../internal/domain/user/repository/user_relationship_group_repository.go)
  - [x] `updateRelationshipGroups(Set<UserRelationshipGroup.Key> keys, @Nullable String name, @Nullable Date creationDate)` -> [internal/domain/user/repository/user_relationship_group_repository.go:UpdateRelationshipGroups(ctx context.Context, keys []po.UserRelationshipGroupKey, newName string, session *mongo.Session)](../internal/domain/user/repository/user_relationship_group_repository.go)
  - [x] `countRelationshipGroups(@Nullable Set<Long> ownerIds, @Nullable Set<Integer> indexes, @Nullable Set<String> names, @Nullable DateRange creationDateRange)` -> [internal/domain/user/repository/user_relationship_group_repository.go:CountRelationshipGroups(ctx context.Context, ownerIDs []int64, groupIndexes []int32)](../internal/domain/user/repository/user_relationship_group_repository.go)
  - [x] `findRelationshipGroups(@Nullable Set<Long> ownerIds, @Nullable Set<Integer> indexes, @Nullable Set<String> names, @Nullable DateRange creationDateRange, @Nullable Integer page, @Nullable Integer size)` -> [internal/domain/user/repository/user_relationship_group_repository.go:FindRelationshipGroups(ctx context.Context, ownerIDs []int64, groupIndexes []int32, page *int, size *int)](../internal/domain/user/repository/user_relationship_group_repository.go)
  - [x] `findRelationshipGroupsInfos(Long ownerId)` -> [internal/domain/user/repository/user_relationship_group_repository.go:FindRelationshipGroupsInfos(ctx context.Context, ownerID int64)](../internal/domain/user/repository/user_relationship_group_repository.go)

- **UserRelationshipRepository.java** ([java/im/turms/service/domain/user/repository/UserRelationshipRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/repository/UserRelationshipRepository.java))
> [简述功能]

  - [x] `deleteAllRelationships(Set<Long> userIds, @Nullable ClientSession session)` -> [internal/domain/user/repository/user_relationship_repository.go:DeleteAllRelationships(ctx context.Context, ownerIDs []int64, session *mongo.Session)](../internal/domain/user/repository/user_relationship_repository.go)
  - [x] `updateUserOneSidedRelationships(Set<UserRelationship.Key> keys, @Nullable String name, @Nullable Date blockDate, @Nullable Date establishmentDate)` -> [internal/domain/user/repository/user_relationship_repository.go:UpdateUserOneSidedRelationships(ctx context.Context, ownerID int64, relatedUserIDs []int64, blockDate *time.Time, establishmentDate *time.Time, name *string, session *mongo.Session,)](../internal/domain/user/repository/user_relationship_repository.go)
  - [x] `countRelationships(@Nullable Set<Long> ownerIds, @Nullable Set<Long> relatedUserIds, @Nullable Boolean isBlocked)` -> [internal/domain/user/repository/user_relationship_repository.go:CountRelationships(ctx context.Context, ownerIDs, relatedUserIDs []int64, groupIndexes []int32, isBlocked *bool, session *mongo.Session)](../internal/domain/user/repository/user_relationship_repository.go)
  - [x] `findRelatedUserIds(@Nullable Set<Long> ownerIds, @Nullable Boolean isBlocked)` -> [internal/domain/user/repository/user_relationship_repository.go:FindRelatedUserIDs(ctx context.Context, ownerIDs []int64, groupIndexes []int32, isBlocked *bool, page, size *int, session *mongo.Session)](../internal/domain/user/repository/user_relationship_repository.go)
  - [x] `findRelationships(@Nullable Set<Long> ownerIds, @Nullable Set<Long> relatedUserIds, @Nullable Boolean isBlocked, @Nullable DateRange establishmentDateRange, @Nullable Integer page, @Nullable Integer size)` -> [internal/domain/user/repository/user_relationship_repository.go:FindRelationships(ctx context.Context, ownerIDs []int64, relatedUserIDs []int64, groupIndexes []int32, isBlocked *bool, establishmentDateRange *turmsmongo.DateRange, page, size *int, session *mongo.Session,)](../internal/domain/user/repository/user_relationship_repository.go)
  - [x] `findRelationships(@Nullable Set<Long> ownerIds, @Nullable Set<Long> relatedUserIds, @Nullable Integer page, @Nullable Integer size)` -> [internal/domain/user/repository/user_relationship_repository.go:FindRelationships(ctx context.Context, ownerIDs []int64, relatedUserIDs []int64, groupIndexes []int32, isBlocked *bool, establishmentDateRange *turmsmongo.DateRange, page, size *int, session *mongo.Session,)](../internal/domain/user/repository/user_relationship_repository.go)
  - [x] `hasRelationshipAndNotBlocked(Long ownerId, Long relatedUserId)` -> [internal/domain/user/repository/user_relationship_repository.go:HasRelationshipAndNotBlocked(ctx context.Context, ownerID, relatedUserID int64, session *mongo.Session)](../internal/domain/user/repository/user_relationship_repository.go)
  - [x] `isBlocked(Long ownerId, Long relatedUserId)` -> [internal/domain/group/service/group_blocklist_service.go:IsBlocked(ctx context.Context, groupID int64, userID int64)](../internal/domain/group/service/group_blocklist_service.go)

- **UserRepository.java** ([java/im/turms/service/domain/user/repository/UserRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/repository/UserRepository.java))
> [简述功能]

  - [x] `updateUsers(Set<Long> userIds, @Nullable byte[] password, @Nullable String name, @Nullable String intro, @Nullable String profilePicture, @Nullable ProfileAccessStrategy profileAccessStrategy, @Nullable Long roleId, @Nullable Date registrationDate, @Nullable Boolean isActive, @Nullable Map<String, Object> userDefinedAttributes, @Nullable ClientSession session)` -> [internal/domain/user/repository/user_repository.go:UpdateUsers(ctx context.Context, userIDs []int64, update bson.M)](../internal/domain/user/repository/user_repository.go)
  - [x] `updateUsersDeletionDate(Set<Long> userIds)` -> [internal/domain/user/repository/user_repository.go:UpdateUsersDeletionDate(ctx context.Context, userIDs []int64)](../internal/domain/user/repository/user_repository.go)
  - [x] `checkIfUserExists(Long userId, boolean queryDeletedRecords)` -> [internal/domain/user/service/user_service.go:CheckIfUserExists(ctx context.Context, userID int64)](../internal/domain/user/service/user_service.go)
  - [x] `countRegisteredUsers(@Nullable DateRange dateRange, boolean queryDeletedRecords)` -> [internal/domain/user/repository/user_repository.go:CountRegisteredUsers(ctx context.Context, startDate *time.Time, endDate *time.Time, queryDeletedRecords bool)](../internal/domain/user/repository/user_repository.go)
  - [x] `countDeletedUsers(@Nullable DateRange dateRange)` -> [internal/domain/user/repository/user_repository.go:CountDeletedUsers(ctx context.Context, startDate *time.Time, endDate *time.Time)](../internal/domain/user/repository/user_repository.go)
  - [x] `countUsers(boolean queryDeletedRecords)` -> [internal/domain/user/repository/user_repository.go:CountUsers(ctx context.Context, startDate *time.Time, endDate *time.Time)](../internal/domain/user/repository/user_repository.go)
  - [x] `countUsers(@Nullable Set<Long> userIds, @Nullable DateRange registrationDateRange, @Nullable DateRange deletionDateRange, @Nullable Boolean isActive)` -> [internal/domain/user/repository/user_repository.go:CountUsers(ctx context.Context, startDate *time.Time, endDate *time.Time)](../internal/domain/user/repository/user_repository.go)
  - [x] `findName(Long userId)` -> [internal/domain/user/repository/user_repository.go:FindName(ctx context.Context, userID int64)](../internal/domain/user/repository/user_repository.go)
  - [x] `findAllNames()` -> [internal/domain/user/repository/user_repository.go:FindAllNames(ctx context.Context)](../internal/domain/user/repository/user_repository.go)
  - [x] `findProfileAccessIfNotDeleted(Long userId)` -> [internal/domain/user/repository/user_repository.go:FindProfileAccessIfNotDeleted(ctx context.Context, userID int64)](../internal/domain/user/repository/user_repository.go)
  - [x] `findUsers(@Nullable Collection<Long> userIds, @Nullable DateRange registrationDateRange, @Nullable DateRange deletionDateRange, @Nullable Boolean isActive, @Nullable Integer page, @Nullable Integer size, boolean queryDeletedRecords)` -> [internal/domain/user/repository/user_repository.go:FindUsers(ctx context.Context, userIDs []int64)](../internal/domain/user/repository/user_repository.go)
  - [x] `findNotDeletedUserProfiles(Collection<Long> userIds, Collection<String> includedUserDefinedAttributes, @Nullable Date lastUpdatedDate)` -> [internal/domain/user/repository/user_repository.go:FindNotDeletedUserProfiles(ctx context.Context, userIDs []int64)](../internal/domain/user/repository/user_repository.go)
  - [x] `findUsersProfile(Collection<Long> userIds, Collection<String> includedUserDefinedAttributes, boolean queryDeletedRecords)` -> [internal/domain/user/repository/user_repository.go:FindUsersProfile(ctx context.Context, userIDs []int64)](../internal/domain/user/repository/user_repository.go)
  - [x] `findUserRoleId(Long userId)` -> [internal/domain/user/repository/user_repository.go:FindUserRoleID(ctx context.Context, userID int64)](../internal/domain/user/repository/user_repository.go)
  - [x] `isActiveAndNotDeleted(Long userId)` -> [internal/domain/user/repository/user_repository.go:IsActiveAndNotDeleted(ctx context.Context, userID int64)](../internal/domain/user/repository/user_repository.go)

- **UserRoleRepository.java** ([java/im/turms/service/domain/user/repository/UserRoleRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/repository/UserRoleRepository.java))
> [简述功能]

  - [x] `updateUserRoles(Set<Long> groupIds, @Nullable String name, @Nullable Set<Long> creatableGroupTypeIds, @Nullable Integer ownedGroupLimit, @Nullable Integer ownedGroupLimitForEachGroupType, @Nullable Map<Long, Integer> groupTypeIdToLimit)` -> [internal/domain/user/repository/user_role_repository.go:UpdateUserRoles(ctx context.Context, roleIDs []int64, update interface{})](../internal/domain/user/repository/user_role_repository.go)

- **UserSettingsRepository.java** ([java/im/turms/service/domain/user/repository/UserSettingsRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/repository/UserSettingsRepository.java))
> [简述功能]

  - [x] `upsertSettings(Long userId, Map<String, Object> settings)` -> [internal/domain/user/repository/user_settings_repository.go:UpsertSettings(ctx context.Context, userID int64, settings map[string]interface{})](../internal/domain/user/repository/user_settings_repository.go)
  - [x] `unsetSettings(Long userId, @Nullable Collection<String> settingNames)` -> [internal/domain/user/repository/user_settings_repository.go:UnsetSettings(ctx context.Context, userID int64, settingsNames []string)](../internal/domain/user/repository/user_settings_repository.go)
  - [x] `findByIdAndSettingNames(Long userId, @Nullable Collection<String> settingNames, @Nullable Date lastUpdatedDateStart)` -> [internal/domain/user/repository/user_settings_repository.go:FindByIdAndSettingNames(ctx context.Context, userID int64, names []string)](../internal/domain/user/repository/user_settings_repository.go)

- **UserVersionRepository.java** ([java/im/turms/service/domain/user/repository/UserVersionRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/repository/UserVersionRepository.java))
> [简述功能]

  - [x] `updateSpecificVersion(Long userId, @Nullable ClientSession session, String... fields)` -> [internal/domain/user/repository/user_version_repository.go:UpdateSpecificVersion(ctx context.Context, userIDs []int64, field string, updateDate time.Time)](../internal/domain/user/repository/user_version_repository.go)
  - [x] `updateSpecificVersion(Long userId, @Nullable ClientSession session, String field)` -> [internal/domain/user/repository/user_version_repository.go:UpdateSpecificVersion(ctx context.Context, userIDs []int64, field string, updateDate time.Time)](../internal/domain/user/repository/user_version_repository.go)
  - [x] `updateSpecificVersion(Set<Long> userIds, @Nullable ClientSession session, String... fields)` -> [internal/domain/user/repository/user_version_repository.go:UpdateSpecificVersion(ctx context.Context, userIDs []int64, field string, updateDate time.Time)](../internal/domain/user/repository/user_version_repository.go)
  - [ ] `findGroupJoinRequests(Long userId)`
  - [ ] `findJoinedGroup(Long userId)`
  - [ ] `findReceivedGroupInvitations(Long userId)`
  - [x] `findRelationships(Long userId)` -> [internal/domain/user/repository/user_relationship_repository.go:FindRelationships(ctx context.Context, ownerIDs []int64, relatedUserIDs []int64, groupIndexes []int32, isBlocked *bool, establishmentDateRange *turmsmongo.DateRange, page, size *int, session *mongo.Session,)](../internal/domain/user/repository/user_relationship_repository.go)
  - [x] `findRelationshipGroups(Long userId)` -> [internal/domain/user/repository/user_relationship_group_repository.go:FindRelationshipGroups(ctx context.Context, ownerIDs []int64, groupIndexes []int32, page *int, size *int)](../internal/domain/user/repository/user_relationship_group_repository.go)
  - [ ] `findSentGroupInvitations(Long userId)`
  - [ ] `findSentFriendRequests(Long userId)`
  - [ ] `findReceivedFriendRequests(Long userId)`

- **UserFriendRequestService.java** ([java/im/turms/service/domain/user/service/UserFriendRequestService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/service/UserFriendRequestService.java))
> [简述功能]

  - [x] `removeAllExpiredFriendRequests(Date expirationDate)` -> [internal/domain/user/service/user_friend_request_service.go:RemoveAllExpiredFriendRequests(ctx context.Context, expirationDate time.Time)](../internal/domain/user/service/user_friend_request_service.go)
  - [x] `hasPendingFriendRequest(@NotNull Long requesterId, @NotNull Long recipientId)` -> [internal/domain/user/repository/user_friend_request_repository.go:HasPendingFriendRequest(ctx context.Context, requesterID, recipientID int64)](../internal/domain/user/repository/user_friend_request_repository.go)
  - [x] `createFriendRequest(@Nullable Long id, @NotNull Long requesterId, @NotNull Long recipientId, @NotNull String content, @Nullable @ValidRequestStatus RequestStatus status, @Nullable @PastOrPresent Date creationDate, @Nullable @PastOrPresent Date responseDate, @Nullable String reason)` -> [internal/domain/user/service/user_friend_request_service.go:CreateFriendRequest(ctx context.Context, requestID *int64, requesterID, recipientID int64, content string, status *po.RequestStatus, creationDate, responseDate *time.Time, reason *string)](../internal/domain/user/service/user_friend_request_service.go)
  - [x] `authAndCreateFriendRequest(@NotNull Long requesterId, @NotNull Long recipientId, @Nullable String content, @NotNull @PastOrPresent Date creationDate)` -> [internal/domain/user/service/user_friend_request_service.go:AuthAndCreateFriendRequest(ctx context.Context, requesterID, recipientID int64, content string, creationDate time.Time)](../internal/domain/user/service/user_friend_request_service.go)
  - [x] `authAndRecallFriendRequest(@NotNull Long requesterId, @NotNull Long requestId)` -> [internal/domain/user/service/user_friend_request_service.go:AuthAndRecallFriendRequest(ctx context.Context, requesterID, requestID int64)](../internal/domain/user/service/user_friend_request_service.go)
  - [x] `updatePendingFriendRequestStatus(@NotNull Long requestId, @NotNull @ValidRequestStatus RequestStatus requestStatus, @Nullable String reason, @Nullable ClientSession session)` -> [internal/domain/user/service/user_friend_request_service.go:UpdatePendingFriendRequestStatus(ctx context.Context, requestID int64, targetStatus po.RequestStatus, reason *string)](../internal/domain/user/service/user_friend_request_service.go)
  - [x] `updateFriendRequests(@NotEmpty Set<Long> requestIds, @Nullable Long requesterId, @Nullable Long recipientId, @Nullable String content, @Nullable @ValidRequestStatus RequestStatus status, @Nullable String reason, @Nullable @PastOrPresent Date creationDate, @Nullable @PastOrPresent Date responseDate)` -> [internal/domain/user/repository/user_friend_request_repository.go:UpdateFriendRequests(ctx context.Context, requestIds []int64, requesterID, recipientID *int64, content *string, status *po.RequestStatus, reason *string, creationDate *time.Time)](../internal/domain/user/repository/user_friend_request_repository.go)
  - [x] `queryRecipientId(@NotNull Long requestId)` -> [internal/domain/user/service/user_friend_request_service.go:QueryRecipientId(ctx context.Context, requestID int64)](../internal/domain/user/service/user_friend_request_service.go)
  - [ ] `queryRequesterIdAndRecipientIdAndStatus(@NotNull Long requestId)`
  - [ ] `queryRequesterIdAndRecipientIdAndCreationDateAndStatus(@NotNull Long requestId)`
  - [x] `authAndHandleFriendRequest(@NotNull Long friendRequestId, @NotNull Long requesterId, @NotNull @ValidResponseAction ResponseAction action, @Nullable String reason)` -> [internal/domain/user/service/user_friend_request_service.go:AuthAndHandleFriendRequest(ctx context.Context, friendRequestID int64, requesterID int64, action po.ResponseAction, reason *string)](../internal/domain/user/service/user_friend_request_service.go)
  - [x] `queryFriendRequestsWithVersion(@NotNull Long userId, boolean areSentByUser, @Nullable Date lastUpdatedDate)` -> [internal/domain/user/service/user_friend_request_service.go:QueryFriendRequestsWithVersion(ctx context.Context, userID int64, isRecipient bool, lastUpdatedDate *time.Time)](../internal/domain/user/service/user_friend_request_service.go)
  - [x] `queryFriendRequestsByRecipientId(@NotNull Long recipientId)` -> [internal/domain/user/service/user_friend_request_service.go:QueryFriendRequestsByRecipientId(ctx context.Context, recipientID int64)](../internal/domain/user/service/user_friend_request_service.go)
  - [x] `queryFriendRequestsByRequesterId(@NotNull Long requesterId)` -> [internal/domain/user/service/user_friend_request_service.go:QueryFriendRequestsByRequesterId(ctx context.Context, requesterID int64)](../internal/domain/user/service/user_friend_request_service.go)
  - [x] `deleteFriendRequests(@Nullable Set<Long> ids)` -> [internal/domain/user/service/user_friend_request_service.go:DeleteFriendRequests(ctx context.Context, ids []int64)](../internal/domain/user/service/user_friend_request_service.go)
  - [x] `queryFriendRequests(@Nullable Set<Long> ids, @Nullable Set<Long> requesterIds, @Nullable Set<Long> recipientIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange, @Nullable Integer page, @Nullable Integer size)` -> [internal/domain/user/service/user_friend_request_service.go:QueryFriendRequests(ctx context.Context, ids, requesterIds, recipientIds []int64, statuses []po.RequestStatus, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd *time.Time, page, size *int)](../internal/domain/user/service/user_friend_request_service.go)
  - [x] `countFriendRequests(@Nullable Set<Long> ids, @Nullable Set<Long> requesterIds, @Nullable Set<Long> recipientIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange)` -> [internal/domain/user/repository/user_friend_request_repository.go:CountFriendRequests(ctx context.Context, ids, requesterIds, recipientIds []int64, statuses []po.RequestStatus, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd *time.Time)](../internal/domain/user/repository/user_friend_request_repository.go)

- **UserRelationshipGroupService.java** ([java/im/turms/service/domain/user/service/UserRelationshipGroupService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/service/UserRelationshipGroupService.java))
> [简述功能]

  - [x] `createRelationshipGroup(@NotNull Long ownerId, @Nullable Integer groupIndex, @NotNull String groupName, @Nullable @PastOrPresent Date creationDate, @Nullable ClientSession session)` -> [internal/domain/user/service/user_relationship_group_service.go:CreateRelationshipGroup(ctx context.Context, ownerID int64, groupIndex *int32, groupName string, creationDate *time.Time, session *mongo.Session,)](../internal/domain/user/service/user_relationship_group_service.go)
  - [x] `queryRelationshipGroupsInfos(@NotNull Long ownerId)` -> [internal/domain/user/service/user_relationship_group_service.go:QueryRelationshipGroupsInfos(ctx context.Context, ownerID int64)](../internal/domain/user/service/user_relationship_group_service.go)
  - [x] `queryRelationshipGroupsInfosWithVersion(@NotNull Long ownerId, @Nullable Date lastUpdatedDate)` -> [internal/domain/user/service/user_relationship_group_service.go:QueryRelationshipGroupsInfosWithVersion(ctx context.Context, ownerID int64, lastUpdatedDate *time.Time,)](../internal/domain/user/service/user_relationship_group_service.go)
  - [x] `queryGroupIndexes(@NotNull Long ownerId, @NotNull Long relatedUserId)` -> [internal/domain/user/service/user_relationship_group_service.go:QueryGroupIndexes(ctx context.Context, ownerID int64, relatedUserID int64)](../internal/domain/user/service/user_relationship_group_service.go)
  - [x] `queryRelationshipGroupMemberIds(@NotNull Long ownerId, @NotNull Integer groupIndex)` -> [internal/domain/user/service/user_relationship_group_service.go:QueryRelationshipGroupMemberIds(ctx context.Context, ownerID int64, groupIndex int32,)](../internal/domain/user/service/user_relationship_group_service.go)
  - [x] `queryRelationshipGroupMemberIds(@Nullable Set<Long> ownerIds, @Nullable Set<Integer> groupIndexes, @Nullable Integer page, @Nullable Integer size)` -> [internal/domain/user/service/user_relationship_group_service.go:QueryRelationshipGroupMemberIds(ctx context.Context, ownerID int64, groupIndex int32,)](../internal/domain/user/service/user_relationship_group_service.go)
  - [x] `updateRelationshipGroupName(@NotNull Long ownerId, @NotNull Integer groupIndex, @NotNull String newGroupName)` -> [internal/domain/user/repository/user_relationship_group_repository.go:UpdateRelationshipGroupName(ctx context.Context, ownerID int64, groupIndex int32, newName string, session *mongo.Session)](../internal/domain/user/repository/user_relationship_group_repository.go)
  - [x] `upsertRelationshipGroupMember(@NotNull Long ownerId, @NotNull Long relatedUserId, @Nullable Integer newGroupIndex, @Nullable Integer deleteGroupIndex, @Nullable ClientSession session)` -> [internal/domain/user/service/user_relationship_group_service.go:UpsertRelationshipGroupMember(ctx context.Context, ownerID int64, relatedUserID int64, newGroupIndex *int32, deleteGroupIndex *int32, session *mongo.Session,)](../internal/domain/user/service/user_relationship_group_service.go)
  - [x] `updateRelationshipGroups(@NotEmpty Set<UserRelationshipGroup.@ValidUserRelationshipGroupKey Key> keys, @Nullable String name, @Nullable @PastOrPresent Date creationDate)` -> [internal/domain/user/repository/user_relationship_group_repository.go:UpdateRelationshipGroups(ctx context.Context, keys []po.UserRelationshipGroupKey, newName string, session *mongo.Session)](../internal/domain/user/repository/user_relationship_group_repository.go)
  - [ ] `addRelatedUserToRelationshipGroups(@NotNull Long ownerId, @NotNull Integer groupIndex, @NotNull Long relatedUserId, @Nullable ClientSession session)`
  - [x] `deleteRelationshipGroupAndMoveMembersToNewGroup(@NotNull Long ownerId, @NotNull Integer deleteGroupIndex, @NotNull Integer newGroupIndex)` -> [internal/domain/user/service/user_relationship_group_service.go:DeleteRelationshipGroupAndMoveMembersToNewGroup(ctx context.Context, ownerID int64, deleteGroupIndex int32, newGroupIndex int32,)](../internal/domain/user/service/user_relationship_group_service.go)
  - [x] `deleteAllRelationshipGroups(@NotEmpty Set<Long> ownerIds, @Nullable ClientSession session, boolean updateRelationshipGroupsVersion)` -> [internal/domain/user/repository/user_relationship_group_repository.go:DeleteAllRelationshipGroups(ctx context.Context, ownerIDs []int64, session *mongo.Session)](../internal/domain/user/repository/user_relationship_group_repository.go)
  - [x] `deleteRelatedUserFromRelationshipGroup(@NotNull Long ownerId, @NotNull Long relatedUserId, @NotNull Integer groupIndex, @Nullable ClientSession session, boolean updateRelationshipGroupsMembersVersion)` -> [internal/domain/user/repository/user_relationship_group_member_repository.go:DeleteRelatedUserFromRelationshipGroup(ctx context.Context, ownerID int64, relatedUserID int64, groupIndexes []int32, session *mongo.Session)](../internal/domain/user/repository/user_relationship_group_member_repository.go)
  - [ ] `deleteRelatedUserFromAllRelationshipGroups(@NotNull Long ownerId, @NotNull Long relatedUserId, @Nullable ClientSession session, boolean updateRelationshipGroupsMembersVersion)`
  - [x] `deleteRelatedUsersFromAllRelationshipGroups(@NotEmpty Set<UserRelationship.@ValidUserRelationshipKey Key> keys, @Nullable ClientSession session, boolean updateRelationshipGroupsMembersVersion)` -> [internal/domain/user/repository/user_relationship_group_member_repository.go:DeleteRelatedUsersFromAllRelationshipGroups(ctx context.Context, ownerID int64, relatedUserIDs []int64, session *mongo.Session)](../internal/domain/user/repository/user_relationship_group_member_repository.go)
  - [x] `moveRelatedUserToNewGroup(@NotNull Long ownerId, @NotNull Long relatedUserId, @NotNull Integer currentGroupIndex, @NotNull Integer targetGroupIndex, boolean suppressIfAlreadyExistsInTargetGroup, @Nullable ClientSession session)` -> [internal/domain/user/service/user_relationship_group_service.go:MoveRelatedUserToNewGroup(ctx context.Context, ownerID int64, relatedUserID int64, currentGroupIndex int32, targetGroupIndex int32, suppressIfAlreadyExists bool, session *mongo.Session,)](../internal/domain/user/service/user_relationship_group_service.go)
  - [x] `deleteRelationshipGroups()` -> [internal/domain/user/repository/user_relationship_group_repository.go:DeleteRelationshipGroups(ctx context.Context, ownerID int64, groupIndexes []int32, session *mongo.Session)](../internal/domain/user/repository/user_relationship_group_repository.go)
  - [x] `deleteRelationshipGroups(@NotEmpty Set<UserRelationshipGroup.@ValidUserRelationshipGroupKey Key> keys)` -> [internal/domain/user/repository/user_relationship_group_repository.go:DeleteRelationshipGroups(ctx context.Context, ownerID int64, groupIndexes []int32, session *mongo.Session)](../internal/domain/user/repository/user_relationship_group_repository.go)
  - [x] `queryRelationshipGroups(@Nullable Set<Long> ownerIds, @Nullable Set<Integer> indexes, @Nullable Set<String> names, @Nullable DateRange creationDateRange, @Nullable Integer page, @Nullable Integer size)` -> [internal/domain/user/service/user_relationship_group_service.go:QueryRelationshipGroups(ctx context.Context, ownerIDs []int64, groupIndexes []int32, page *int, size *int,)](../internal/domain/user/service/user_relationship_group_service.go)
  - [x] `countRelationshipGroups(@Nullable Set<Long> ownerIds, @Nullable Set<Integer> indexes, @Nullable Set<String> names, @Nullable DateRange creationDateRange)` -> [internal/domain/user/repository/user_relationship_group_repository.go:CountRelationshipGroups(ctx context.Context, ownerIDs []int64, groupIndexes []int32)](../internal/domain/user/repository/user_relationship_group_repository.go)
  - [x] `countRelationshipGroups(@Nullable Set<Long> ownerIds, @Nullable Set<Long> relatedUserIds)` -> [internal/domain/user/repository/user_relationship_group_repository.go:CountRelationshipGroups(ctx context.Context, ownerIDs []int64, groupIndexes []int32)](../internal/domain/user/repository/user_relationship_group_repository.go)
  - [x] `countRelationshipGroupMembers(@Nullable Set<Long> ownerIds, @Nullable Set<Integer> groupIndexes)` -> [internal/domain/user/service/user_relationship_group_service.go:CountRelationshipGroupMembers(ctx context.Context, ownerIDs []int64, groupIndexes []int32)](../internal/domain/user/service/user_relationship_group_service.go)

- **UserRelationshipService.java** ([java/im/turms/service/domain/user/service/UserRelationshipService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/service/UserRelationshipService.java))
> [简述功能]

  - [x] `deleteAllRelationships(@NotEmpty Set<Long> userIds, @Nullable ClientSession session, boolean updateRelationshipsVersion)` -> [internal/domain/user/repository/user_relationship_repository.go:DeleteAllRelationships(ctx context.Context, ownerIDs []int64, session *mongo.Session)](../internal/domain/user/repository/user_relationship_repository.go)
  - [x] `deleteOneSidedRelationships(@NotEmpty Set<UserRelationship.@ValidUserRelationshipKey Key> keys)` -> [internal/domain/user/repository/user_relationship_repository.go:DeleteOneSidedRelationships(ctx context.Context, ownerIDs, relatedUserIDs []int64, session *mongo.Session)](../internal/domain/user/repository/user_relationship_repository.go)
  - [x] `deleteOneSidedRelationship(@NotNull Long ownerId, @NotNull Long relatedUserId, @Nullable Integer groupIndex, @Nullable ClientSession session)` -> [internal/domain/user/service/user_relationship_service.go:DeleteOneSidedRelationship(ctx context.Context, ownerID int64, relatedUserID int64,)](../internal/domain/user/service/user_relationship_service.go)
  - [x] `tryDeleteTwoSidedRelationships(@NotNull Long requesterId, @NotNull Long relatedUserId, @Nullable Integer groupId)` -> [internal/domain/user/service/user_relationship_service.go:TryDeleteTwoSidedRelationships(ctx context.Context, user1ID int64, user2ID int64, session *mongo.Session,)](../internal/domain/user/service/user_relationship_service.go)
  - [x] `queryRelatedUserIdsWithVersion(@NotNull Long ownerId, @Nullable Set<Integer> groupIndexes, @Nullable Boolean isBlocked, @Nullable Date lastUpdatedDate)` -> [internal/domain/user/service/user_relationship_service.go:QueryRelatedUserIdsWithVersion(ctx context.Context, ownerID int64, groupIndexes []int32, isBlocked *bool, lastUpdatedDate *time.Time,)](../internal/domain/user/service/user_relationship_service.go)
  - [x] `queryRelationshipsWithVersion(@NotNull Long ownerId, @Nullable Set<Long> relatedUserIds, @Nullable Set<Integer> groupIndexes, @Nullable Boolean isBlocked, @Nullable Date lastUpdatedDate)` -> [internal/domain/user/service/user_relationship_service.go:QueryRelationshipsWithVersion(ctx context.Context, ownerID int64, relatedUserIDs []int64, groupIndexes []int32, isBlocked *bool, lastUpdatedDate *time.Time,)](../internal/domain/user/service/user_relationship_service.go)
  - [x] `queryRelatedUserIds(@Nullable Set<Long> ownerIds, @Nullable Boolean isBlocked)` -> [internal/domain/user/service/user_relationship_service.go:QueryRelatedUserIds(ctx context.Context, ownerIDs []int64, groupIndexes []int32, isBlocked *bool, page *int, size *int,)](../internal/domain/user/service/user_relationship_service.go)
  - [x] `queryRelatedUserIds(@Nullable Set<Long> ownerIds, @Nullable Set<Integer> groupIndexes, @Nullable Boolean isBlocked)` -> [internal/domain/user/service/user_relationship_service.go:QueryRelatedUserIds(ctx context.Context, ownerIDs []int64, groupIndexes []int32, isBlocked *bool, page *int, size *int,)](../internal/domain/user/service/user_relationship_service.go)
  - [x] `queryRelationships(@Nullable Set<Long> ownerIds, @Nullable Set<Long> relatedUserIds, @Nullable Set<Integer> groupIndexes, @Nullable Boolean isBlocked, @Nullable DateRange establishmentDateRange, @Nullable Integer page, @Nullable Integer size)` -> [internal/domain/user/service/user_relationship_service.go:QueryRelationships(ctx context.Context, ownerIDs []int64, relatedUserIDs []int64, groupIndexes []int32, isBlocked *bool, establishmentDateRange *turmsmongo.DateRange, page *int, size *int,)](../internal/domain/user/service/user_relationship_service.go)
  - [x] `queryMembersRelationships(@Nullable Set<Long> ownerIds, @Nullable Set<Integer> groupIndexes, @Nullable Integer page, @Nullable Integer size)` -> [internal/domain/user/repository/user_relationship_repository.go:QueryMembersRelationships(ctx context.Context, ownerID int64, groupIndexes []int32, page, size *int, session *mongo.Session)](../internal/domain/user/repository/user_relationship_repository.go)
  - [x] `countRelationships(@Nullable Set<Long> ownerIds, @Nullable Set<Long> relatedUserIds, @Nullable Set<Integer> groupIndexes, @Nullable Boolean isBlocked)` -> [internal/domain/user/repository/user_relationship_repository.go:CountRelationships(ctx context.Context, ownerIDs, relatedUserIDs []int64, groupIndexes []int32, isBlocked *bool, session *mongo.Session)](../internal/domain/user/repository/user_relationship_repository.go)
  - [x] `countRelationships(@Nullable Set<Long> ownerIds, @Nullable Set<Long> relatedUserIds, @Nullable Boolean isBlocked)` -> [internal/domain/user/repository/user_relationship_repository.go:CountRelationships(ctx context.Context, ownerIDs, relatedUserIDs []int64, groupIndexes []int32, isBlocked *bool, session *mongo.Session)](../internal/domain/user/repository/user_relationship_repository.go)
  - [x] `friendTwoUsers(@NotNull Long userOneId, @NotNull Long userTwoId, @Nullable ClientSession session)` -> [internal/domain/user/service/user_relationship_service.go:FriendTwoUsers(ctx context.Context, user1ID, user2ID int64)](../internal/domain/user/service/user_relationship_service.go)
  - [x] `upsertOneSidedRelationship(@NotNull Long ownerId, @NotNull Long relatedUserId, @Nullable String name, @Nullable @PastOrPresent Date blockDate, @Nullable Integer newGroupIndex, @Nullable Integer deleteGroupIndex, @Nullable @PastOrPresent Date establishmentDate, boolean upsert, @Nullable ClientSession session)` -> [internal/domain/user/service/user_relationship_service.go:UpsertOneSidedRelationship(ctx context.Context, ownerID int64, relatedUserID int64, blockDate *time.Time, groupIndex *int32, establishmentDate *time.Time, name *string, session *mongo.Session,)](../internal/domain/user/service/user_relationship_service.go)
  - [x] `isBlocked(@NotNull Long ownerId, @NotNull Long relatedUserId, boolean preferCache)` -> [internal/domain/group/service/group_blocklist_service.go:IsBlocked(ctx context.Context, groupID int64, userID int64)](../internal/domain/group/service/group_blocklist_service.go)
  - [x] `isNotBlocked(@NotNull Long ownerId, @NotNull Long relatedUserId, boolean preferCache)` -> [internal/domain/user/service/user_relationship_service.go:IsNotBlocked(ctx context.Context, ownerID, relatedUserID int64)](../internal/domain/user/service/user_relationship_service.go)
  - [x] `hasRelationshipAndNotBlocked(@NotNull Long ownerId, @NotNull Long relatedUserId)` -> [internal/domain/user/repository/user_relationship_repository.go:HasRelationshipAndNotBlocked(ctx context.Context, ownerID, relatedUserID int64, session *mongo.Session)](../internal/domain/user/repository/user_relationship_repository.go)
  - [x] `hasRelationshipAndNotBlocked(@NotNull Long ownerId, @NotNull Long relatedUserId, boolean preferCache)` -> [internal/domain/user/repository/user_relationship_repository.go:HasRelationshipAndNotBlocked(ctx context.Context, ownerID, relatedUserID int64, session *mongo.Session)](../internal/domain/user/repository/user_relationship_repository.go)
  - [x] `updateUserOneSidedRelationships(@NotEmpty Set<UserRelationship.@ValidUserRelationshipKey Key> keys, @Nullable String name, @Nullable @PastOrPresent Date blockDate, @Nullable @PastOrPresent Date establishmentDate)` -> [internal/domain/user/repository/user_relationship_repository.go:UpdateUserOneSidedRelationships(ctx context.Context, ownerID int64, relatedUserIDs []int64, blockDate *time.Time, establishmentDate *time.Time, name *string, session *mongo.Session,)](../internal/domain/user/repository/user_relationship_repository.go)
  - [x] `hasOneSidedRelationship(@NotNull Long ownerId, @NotNull Long relatedUserId)` -> [internal/domain/user/repository/user_relationship_repository.go:HasOneSidedRelationship(ctx context.Context, ownerID, relatedUserID int64, session *mongo.Session)](../internal/domain/user/repository/user_relationship_repository.go)

- **UserRoleService.java** ([java/im/turms/service/domain/user/service/UserRoleService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/service/UserRoleService.java))
> [简述功能]

  - [x] `queryUserRoles(@Nullable Integer page, @Nullable Integer size)` -> [internal/domain/user/service/user_role_service.go:QueryUserRoles(ctx context.Context, filter bson.M)](../internal/domain/user/service/user_role_service.go)
  - [x] `addUserRole(@Nullable Long groupId, @Nullable String name, @NotNull Set<Long> creatableGroupTypeIds, @NotNull Integer ownedGroupLimit, @NotNull Integer ownedGroupLimitForEachGroupType, @NotNull Map<Long, Integer> groupTypeIdToLimit)` -> [internal/domain/user/service/user_role_service.go:AddUserRole(ctx context.Context, role *po.UserRole)](../internal/domain/user/service/user_role_service.go)
  - [x] `updateUserRoles(@NotEmpty Set<Long> groupIds, @Nullable String name, @Nullable Set<Long> creatableGroupTypeIds, @Nullable Integer ownedGroupLimit, @Nullable Integer ownedGroupLimitForEachGroupType, @Nullable Map<Long, Integer> groupTypeIdToLimit)` -> [internal/domain/user/repository/user_role_repository.go:UpdateUserRoles(ctx context.Context, roleIDs []int64, update interface{})](../internal/domain/user/repository/user_role_repository.go)
  - [x] `deleteUserRoles(@Nullable Set<Long> groupIds)` -> [internal/domain/user/service/user_role_service.go:DeleteUserRoles(ctx context.Context, filter bson.M)](../internal/domain/user/service/user_role_service.go)
  - [x] `queryUserRoleById(@NotNull Long id)` -> [internal/domain/user/service/user_role_service.go:QueryUserRoleById(ctx context.Context, roleID int64)](../internal/domain/user/service/user_role_service.go)
  - [x] `queryStoredOrDefaultUserRoleByUserId(@NotNull Long userId)` -> [internal/domain/user/service/user_role_service.go:QueryStoredOrDefaultUserRoleByUserId(ctx context.Context, userID int64)](../internal/domain/user/service/user_role_service.go)
  - [x] `countUserRoles()` -> [internal/domain/user/service/user_role_service.go:CountUserRoles(ctx context.Context, filter bson.M)](../internal/domain/user/service/user_role_service.go)

- **UserService.java** ([java/im/turms/service/domain/user/service/UserService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/service/UserService.java))
> [简述功能]

  - [x] `isAllowedToSendMessageToTarget(@NotNull Boolean isGroupMessage, @NotNull Boolean isSystemMessage, @NotNull Long requesterId, @NotNull Long targetId)` -> [internal/domain/user/service/user_service.go:IsAllowedToSendMessageToTarget(ctx context.Context, isGroupMessage bool, isSystemMessage bool, requesterID int64, targetID int64)](../internal/domain/user/service/user_service.go)
  - [x] `createUser(@Nullable Long id, @Nullable String rawPassword, @Nullable String name, @Nullable String intro, @Nullable String profilePicture, @Nullable @ValidProfileAccess ProfileAccessStrategy profileAccessStrategy, @Nullable Long roleId, @Nullable @PastOrPresent Date registrationDate, @Nullable Boolean isActive)` -> [internal/domain/user/service/user_service.go:CreateUser(ctx context.Context, password string, name string, intro string, profilePicture string, profileAccess int32, permissionGroupID int64, isActive bool)](../internal/domain/user/service/user_service.go)
  - [x] `addUser(@Nullable Long id, @Nullable String rawPassword, @Nullable String name, @Nullable String intro, @Nullable String profilePicture, @Nullable @ValidProfileAccess ProfileAccessStrategy profileAccessStrategy, @Nullable Long roleId, @Nullable @PastOrPresent Date registrationDate, @Nullable Boolean isActive)` -> [internal/domain/user/service/user_service.go:AddUser(ctx context.Context, id int64, password string, name string, intro string, profilePicture string, profileAccess int32, permissionGroupID int64, registrationDate time.Time, isActive bool)](../internal/domain/user/service/user_service.go)
  - [x] `isAllowToQueryUserProfile(@NotNull Long requesterId, @NotNull Long targetUserId)` -> [internal/domain/user/service/user_service.go:IsAllowToQueryUserProfile(ctx context.Context, requesterID int64, targetID int64)](../internal/domain/user/service/user_service.go)
  - [x] `authAndQueryUsersProfile(@NotNull Long requesterId, @Nullable Set<Long> userIds, @Nullable String name, @Nullable Date lastUpdatedDate, @Nullable Integer skip, @Nullable Integer limit, @Nullable List<Integer> fieldsToHighlight)` -> [internal/domain/user/service/user_service.go:AuthAndQueryUsersProfile(ctx context.Context, requesterID int64, userIDs []int64, name string, lastUpdatedDate *time.Time, skip int, limit int)](../internal/domain/user/service/user_service.go)
  - [x] `queryUserName(@NotNull Long userId)` -> [internal/domain/user/service/user_service.go:QueryUserName(ctx context.Context, userID int64)](../internal/domain/user/service/user_service.go)
  - [x] `queryUsersProfile(@NotEmpty Collection<Long> userIds, boolean queryDeletedRecords)` -> [internal/domain/user/service/user_service.go:QueryUsersProfile(ctx context.Context, userIDs []int64)](../internal/domain/user/service/user_service.go)
  - [x] `queryUserRoleIdByUserId(@NotNull Long userId)` -> [internal/domain/user/service/user_service.go:QueryUserRoleIDByUserID(ctx context.Context, userID int64)](../internal/domain/user/service/user_service.go)
  - [x] `deleteUsers(@NotEmpty Set<Long> userIds, @Nullable Boolean deleteLogically)` -> [internal/domain/user/service/user_service.go:DeleteUsers(ctx context.Context, userIDs []int64)](../internal/domain/user/service/user_service.go)
  - [x] `checkIfUserExists(@NotNull Long userId, boolean queryDeletedRecords)` -> [internal/domain/user/service/user_service.go:CheckIfUserExists(ctx context.Context, userID int64)](../internal/domain/user/service/user_service.go)
  - [x] `updateUser(@NotNull Long userId, @Nullable String rawPassword, @Nullable String name, @Nullable String intro, @Nullable String profilePicture, @Nullable @ValidProfileAccess ProfileAccessStrategy profileAccessStrategy, @Nullable Long roleId, @Nullable Boolean isActive, @Nullable @PastOrPresent Date registrationDate, @Nullable Map<String, Value> userDefinedAttributes)` -> [internal/domain/user/service/user_service.go:UpdateUser(ctx context.Context, userID int64, update bson.M)](../internal/domain/user/service/user_service.go)
  - [x] `queryUsers(@Nullable Collection<Long> userIds, @Nullable DateRange registrationDateRange, @Nullable DateRange deletionDateRange, @Nullable Boolean isActive, @Nullable Integer page, @Nullable Integer size, boolean queryDeletedRecords)` -> [internal/domain/user/service/user_service.go:QueryUsers(ctx context.Context, userIDs []int64)](../internal/domain/user/service/user_service.go)
  - [x] `countRegisteredUsers(@Nullable DateRange dateRange, boolean queryDeletedRecords)` -> [internal/domain/user/repository/user_repository.go:CountRegisteredUsers(ctx context.Context, startDate *time.Time, endDate *time.Time, queryDeletedRecords bool)](../internal/domain/user/repository/user_repository.go)
  - [x] `countDeletedUsers(@Nullable DateRange dateRange)` -> [internal/domain/user/repository/user_repository.go:CountDeletedUsers(ctx context.Context, startDate *time.Time, endDate *time.Time)](../internal/domain/user/repository/user_repository.go)
  - [x] `countUsers(boolean queryDeletedRecords)` -> [internal/domain/user/repository/user_repository.go:CountUsers(ctx context.Context, startDate *time.Time, endDate *time.Time)](../internal/domain/user/repository/user_repository.go)
  - [x] `countUsers(@Nullable Set<Long> userIds, @Nullable DateRange registrationDateRange, @Nullable DateRange deletionDateRange, @Nullable Boolean isActive)` -> [internal/domain/user/repository/user_repository.go:CountUsers(ctx context.Context, startDate *time.Time, endDate *time.Time)](../internal/domain/user/repository/user_repository.go)
  - [x] `updateUsers(@NotEmpty Set<Long> userIds, @Nullable String rawPassword, @Nullable String name, @Nullable String intro, @Nullable String profilePicture, @Nullable @ValidProfileAccess ProfileAccessStrategy profileAccessStrategy, @Nullable Long roleId, @Nullable @PastOrPresent Date registrationDate, @Nullable Boolean isActive, @Nullable Map<String, Object> userDefinedAttributes)` -> [internal/domain/user/repository/user_repository.go:UpdateUsers(ctx context.Context, userIDs []int64, update bson.M)](../internal/domain/user/repository/user_repository.go)

- **UserSettingsService.java** ([java/im/turms/service/domain/user/service/UserSettingsService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/service/UserSettingsService.java))
> [简述功能]

  - [x] `upsertSettings(Long userId, Map<String, Value> settings)` -> [internal/domain/user/service/user_settings_service.go:UpsertSettings(ctx context.Context, userID int64, settings map[string]interface{})](../internal/domain/user/service/user_settings_service.go)
  - [x] `deleteSettings(Collection<Long> userIds, @Nullable ClientSession clientSession)` -> [internal/domain/user/repository/user_settings_repository.go:DeleteSettings(ctx context.Context, filter interface{})](../internal/domain/user/repository/user_settings_repository.go)
  - [x] `unsetSettings(Long userId, @Nullable Set<String> settingNames)` -> [internal/domain/user/service/user_settings_service.go:UnsetSettings(ctx context.Context, userID int64, keys []string)](../internal/domain/user/service/user_settings_service.go)
  - [x] `querySettings(Long userId, @Nullable Set<String> settingNames, @Nullable Date lastUpdatedDateStart)` -> [internal/domain/user/service/user_settings_service.go:QuerySettings(ctx context.Context, filter bson.M)](../internal/domain/user/service/user_settings_service.go)

- **UserVersionService.java** ([java/im/turms/service/domain/user/service/UserVersionService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/service/UserVersionService.java))
> [简述功能]

  - [x] `queryRelationshipsLastUpdatedDate(@NotNull Long userId)` -> [internal/domain/user/service/user_version_service.go:QueryRelationshipsLastUpdatedDate(ctx context.Context, userID int64)](../internal/domain/user/service/user_version_service.go)
  - [ ] `querySentGroupInvitationsLastUpdatedDate(@NotNull Long userId)`
  - [ ] `queryReceivedGroupInvitationsLastUpdatedDate(@NotNull Long userId)`
  - [x] `queryGroupJoinRequestsVersion(@NotNull Long userId)` -> [internal/domain/user/service/user_version_service.go:QueryGroupJoinRequestsVersion(ctx context.Context, userID int64)](../internal/domain/user/service/user_version_service.go)
  - [x] `queryRelationshipGroupsLastUpdatedDate(@NotNull Long userId)` -> [internal/domain/user/service/user_version_service.go:QueryRelationshipGroupsLastUpdatedDate(ctx context.Context, userID int64)](../internal/domain/user/service/user_version_service.go)
  - [x] `queryJoinedGroupVersion(@NotNull Long userId)` -> [internal/domain/user/service/user_version_service.go:QueryJoinedGroupVersion(ctx context.Context, userID int64)](../internal/domain/user/service/user_version_service.go)
  - [x] `querySentFriendRequestsVersion(@NotNull Long userId)` -> [internal/domain/user/service/user_version_service.go:QuerySentFriendRequestsVersion(ctx context.Context, userID int64)](../internal/domain/user/service/user_version_service.go)
  - [x] `queryReceivedFriendRequestsVersion(@NotNull Long userId)` -> [internal/domain/user/service/user_version_service.go:QueryReceivedFriendRequestsVersion(ctx context.Context, userID int64)](../internal/domain/user/service/user_version_service.go)
  - [x] `upsertEmptyUserVersion(@NotNull Long userId, @NotNull Date timestamp, @Nullable ClientSession session)` -> [internal/domain/user/repository/user_version_repository.go:UpsertEmptyUserVersion(ctx context.Context, userID int64)](../internal/domain/user/repository/user_version_repository.go)
  - [x] `updateRelationshipsVersion(@NotNull Long userId, @Nullable ClientSession session)` -> [internal/domain/user/service/user_version_service.go:UpdateRelationshipsVersion(ctx context.Context, userID int64)](../internal/domain/user/service/user_version_service.go)
  - [x] `updateRelationshipsVersion(@NotEmpty Set<Long> userIds, @Nullable ClientSession session)` -> [internal/domain/user/service/user_version_service.go:UpdateRelationshipsVersion(ctx context.Context, userID int64)](../internal/domain/user/service/user_version_service.go)
  - [x] `updateSentFriendRequestsVersion(@NotNull Long userId)` -> [internal/domain/user/service/user_version_service.go:UpdateSentFriendRequestsVersion(ctx context.Context, userID int64)](../internal/domain/user/service/user_version_service.go)
  - [x] `updateReceivedFriendRequestsVersion(@NotNull Long userId)` -> [internal/domain/user/service/user_version_service.go:UpdateReceivedFriendRequestsVersion(ctx context.Context, userID int64)](../internal/domain/user/service/user_version_service.go)
  - [x] `updateRelationshipGroupsVersion(@NotNull Long userId)` -> [internal/domain/user/service/user_version_service.go:UpdateRelationshipGroupsVersion(ctx context.Context, userID int64)](../internal/domain/user/service/user_version_service.go)
  - [x] `updateRelationshipGroupsVersion(@NotEmpty Set<Long> userIds)` -> [internal/domain/user/service/user_version_service.go:UpdateRelationshipGroupsVersion(ctx context.Context, userID int64)](../internal/domain/user/service/user_version_service.go)
  - [x] `updateRelationshipGroupsMembersVersion(@NotNull Long userId)` -> [internal/domain/user/service/user_version_service.go:UpdateRelationshipGroupsMembersVersion(ctx context.Context, userID int64)](../internal/domain/user/service/user_version_service.go)
  - [x] `updateRelationshipGroupsMembersVersion(@NotEmpty Set<Long> userIds)` -> [internal/domain/user/service/user_version_service.go:UpdateRelationshipGroupsMembersVersion(ctx context.Context, userID int64)](../internal/domain/user/service/user_version_service.go)
  - [x] `updateSentGroupInvitationsVersion(@NotNull Long userId)` -> [internal/domain/user/service/user_version_service.go:UpdateSentGroupInvitationsVersion(ctx context.Context, userID int64)](../internal/domain/user/service/user_version_service.go)
  - [x] `updateReceivedGroupInvitationsVersion(@NotNull Long userId)` -> [internal/domain/user/service/user_version_service.go:UpdateReceivedGroupInvitationsVersion(ctx context.Context, userID int64)](../internal/domain/user/service/user_version_service.go)
  - [x] `updateSentGroupJoinRequestsVersion(@NotNull Long userId)` -> [internal/domain/user/service/user_version_service.go:UpdateSentGroupJoinRequestsVersion(ctx context.Context, userID int64)](../internal/domain/user/service/user_version_service.go)
  - [x] `updateJoinedGroupsVersion(@NotNull Long userId)` -> [internal/domain/user/service/user_version_service.go:UpdateJoinedGroupsVersion(ctx context.Context, userID int64)](../internal/domain/user/service/user_version_service.go)
  - [x] `updateSpecificVersion(@NotNull Long userId, @Nullable ClientSession session, @NotEmpty String... fields)` -> [internal/domain/user/repository/user_version_repository.go:UpdateSpecificVersion(ctx context.Context, userIDs []int64, field string, updateDate time.Time)](../internal/domain/user/repository/user_version_repository.go)
  - [x] `updateSpecificVersion(@NotNull Long userId, @Nullable ClientSession session, @NotNull String field)` -> [internal/domain/user/repository/user_version_repository.go:UpdateSpecificVersion(ctx context.Context, userIDs []int64, field string, updateDate time.Time)](../internal/domain/user/repository/user_version_repository.go)
  - [x] `updateSpecificVersion(@NotEmpty Set<Long> userIds, @Nullable ClientSession session, @NotEmpty String... fields)` -> [internal/domain/user/repository/user_version_repository.go:UpdateSpecificVersion(ctx context.Context, userIDs []int64, field string, updateDate time.Time)](../internal/domain/user/repository/user_version_repository.go)
  - [x] `delete(@NotEmpty Set<Long> userIds, @Nullable ClientSession session)` -> [internal/domain/common/cache/sharded_map.go:Delete(key K)](../internal/domain/common/cache/sharded_map.go)

- **NearbyUserService.java** ([java/im/turms/service/domain/user/service/onlineuser/NearbyUserService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/service/onlineuser/NearbyUserService.java))
> [简述功能]

  - [x] `queryNearbyUsers(@NotNull Long userId, @NotNull DeviceType deviceType, @Nullable Float longitude, @Nullable Float latitude, @Nullable Short maxCount, @Nullable Integer maxDistance, boolean withCoordinates, boolean withDistance, boolean withUserInfo)` -> [internal/domain/user/service/onlineuser/nearby_user_service.go:QueryNearbyUsers(ctx context.Context, userID int64, deviceType protocol.DeviceType, longitude *float32, latitude *float32, maxCount *int, maxDistance *float64, withCoordinates bool, withDistance bool, withUserInfo bool)](../internal/domain/user/service/onlineuser/nearby_user_service.go)

- **SessionService.java** ([java/im/turms/service/domain/user/service/onlineuser/SessionService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/service/onlineuser/SessionService.java))
> [简述功能]

  - [x] `disconnect(@NotNull Long userId, @NotNull SessionCloseStatus closeStatus)` -> [internal/domain/user/service/onlineuser/session_service.go:Disconnect(ctx context.Context, userID int64, closeStatus int)](../internal/domain/user/service/onlineuser/session_service.go)
  - [x] `disconnect(@NotNull Long userId, @NotNull Set<@ValidDeviceType DeviceType> deviceTypes, @NotNull SessionCloseStatus closeStatus)` -> [internal/domain/user/service/onlineuser/session_service.go:Disconnect(ctx context.Context, userID int64, closeStatus int)](../internal/domain/user/service/onlineuser/session_service.go)
  - [x] `disconnect(@NotNull Long userId, @NotNull @ValidDeviceType DeviceType deviceType, @NotNull SessionCloseStatus closeStatus)` -> [internal/domain/user/service/onlineuser/session_service.go:Disconnect(ctx context.Context, userID int64, closeStatus int)](../internal/domain/user/service/onlineuser/session_service.go)
  - [x] `disconnect(@NotNull Set<Long> userIds, @NotNull SessionCloseStatus closeStatus)` -> [internal/domain/user/service/onlineuser/session_service.go:Disconnect(ctx context.Context, userID int64, closeStatus int)](../internal/domain/user/service/onlineuser/session_service.go)
  - [x] `disconnect(@NotNull Set<Long> userIds, @Nullable Set<@ValidDeviceType DeviceType> deviceTypes, @NotNull SessionCloseStatus closeStatus)` -> [internal/domain/user/service/onlineuser/session_service.go:Disconnect(ctx context.Context, userID int64, closeStatus int)](../internal/domain/user/service/onlineuser/session_service.go)
  - [x] `queryUserSessions(Set<Long> userIds)` -> [internal/domain/user/service/onlineuser/session_service.go:QueryUserSessions(ctx context.Context, userIDs []int64)](../internal/domain/user/service/onlineuser/session_service.go)

- **LocaleUtil.java** ([java/im/turms/service/infra/locale/LocaleUtil.java](../turms-orig/turms-service/src/main/java/im/turms/service/infra/locale/LocaleUtil.java))
> [简述功能]

  - [x] `isAvailableLanguage(String languageId)` -> [internal/infra/locale/locale_util.go:IsAvailableLanguage(languageID string)](../internal/infra/locale/locale_util.go)

- **ApiLoggingContext.java** ([java/im/turms/service/infra/logging/ApiLoggingContext.java](../turms-orig/turms-service/src/main/java/im/turms/service/infra/logging/ApiLoggingContext.java))
> [简述功能]

  - [x] `shouldLogRequest(TurmsRequest.KindCase requestType)` -> [internal/infra/logging/api_logging_context.go:ShouldLogRequest(requestType int)](../internal/infra/logging/api_logging_context.go)
  - [x] `shouldLogNotification(TurmsRequest.KindCase requestType)` -> [internal/infra/logging/api_logging_context.go:ShouldLogNotification(requestType int)](../internal/infra/logging/api_logging_context.go)

- **ClientApiLogging.java** ([java/im/turms/service/infra/logging/ClientApiLogging.java](../turms-orig/turms-service/src/main/java/im/turms/service/infra/logging/ClientApiLogging.java))
> [简述功能]

  - [x] `log(ClientRequest request, ServiceRequest serviceRequest, long requestSize, long requestTime, ServiceResponse response, long processingTime)` -> [internal/infra/logging/client_api_logging.go:Log(request interface{}, serviceRequest interface{}, requestSize int64, requestTime int64, response interface{}, processingTime int64)](../internal/infra/logging/client_api_logging.go)

- **AcceptMeetingInvitationResult.java** ([java/im/turms/service/infra/plugin/extension/model/AcceptMeetingInvitationResult.java](../turms-orig/turms-service/src/main/java/im/turms/service/infra/plugin/extension/model/AcceptMeetingInvitationResult.java))
> [简述功能]

  - [ ] `AcceptMeetingInvitationResult(String accessToken)`

- **CreateMeetingOptions.java** ([java/im/turms/service/infra/plugin/extension/model/CreateMeetingOptions.java](../turms-orig/turms-service/src/main/java/im/turms/service/infra/plugin/extension/model/CreateMeetingOptions.java))
> [简述功能]

  - [ ] `CreateMeetingOptions(@Nullable Integer maxParticipants, @Nullable Long idleTimeoutMillis // No plugins support this, so we hide it for now. // @Nullable Long maxDurationMillis)`

- **CreateMeetingResult.java** ([java/im/turms/service/infra/plugin/extension/model/CreateMeetingResult.java](../turms-orig/turms-service/src/main/java/im/turms/service/infra/plugin/extension/model/CreateMeetingResult.java))
> [简述功能]

  - [ ] `CreateMeetingResult(String accessToken)`

- **ProtoModelConvertor.java** ([java/im/turms/service/infra/proto/ProtoModelConvertor.java](../turms-orig/turms-service/src/main/java/im/turms/service/infra/proto/ProtoModelConvertor.java))
> [简述功能]

  - [x] `toList(Map<String, String> map)` -> [internal/infra/proto/proto_model_convertor.go:ToList(protoItems interface{})](../internal/infra/proto/proto_model_convertor.go)
  - [x] `value2proto(Value.Builder builder, Object object)` -> [internal/infra/proto/proto_model_convertor.go:Value2Proto(value interface{})](../internal/infra/proto/proto_model_convertor.go)

- **DefaultLanguageSettings.java** ([java/im/turms/service/storage/elasticsearch/DefaultLanguageSettings.java](../turms-orig/turms-service/src/main/java/im/turms/service/storage/elasticsearch/DefaultLanguageSettings.java))
> [简述功能]

  - [x] `getSetting(LanguageCode code)` -> [internal/storage/elasticsearch/default_language_settings.go:GetSetting()](../internal/storage/elasticsearch/default_language_settings.go)

- **ElasticsearchClient.java** ([java/im/turms/service/storage/elasticsearch/ElasticsearchClient.java](../turms-orig/turms-service/src/main/java/im/turms/service/storage/elasticsearch/ElasticsearchClient.java))
> [简述功能]

  - [x] `healthcheck()` -> [internal/storage/elasticsearch/elasticsearch_client.go:Healthcheck(ctx context.Context)](../internal/storage/elasticsearch/elasticsearch_client.go)
  - [x] `putIndex(String index, CreateIndexRequest request)` -> [internal/storage/elasticsearch/elasticsearch_client.go:PutIndex(ctx context.Context, request *model.CreateIndexRequest)](../internal/storage/elasticsearch/elasticsearch_client.go)
  - [x] `putDoc(String index, String id, Supplier<ByteBuf> payloadSupplier)` -> [internal/storage/elasticsearch/elasticsearch_client.go:PutDoc(ctx context.Context)](../internal/storage/elasticsearch/elasticsearch_client.go)
  - [x] `deleteDoc(String index, String id)` -> [internal/storage/elasticsearch/elasticsearch_client.go:DeleteDoc(ctx context.Context)](../internal/storage/elasticsearch/elasticsearch_client.go)
  - [x] `deleteByQuery(String index, DeleteByQueryRequest request)` -> [internal/storage/elasticsearch/elasticsearch_client.go:DeleteByQuery(ctx context.Context, request *model.DeleteByQueryRequest)](../internal/storage/elasticsearch/elasticsearch_client.go)
  - [x] `updateByQuery(String index, UpdateByQueryRequest request)` -> [internal/storage/elasticsearch/elasticsearch_client.go:UpdateByQuery(ctx context.Context)](../internal/storage/elasticsearch/elasticsearch_client.go)
  - [x] `search(String index, SearchRequest request, ObjectReader reader)` -> [internal/storage/elasticsearch/elasticsearch_client.go:Search(ctx context.Context)](../internal/storage/elasticsearch/elasticsearch_client.go)
  - [x] `bulk(BulkRequest request)` -> [internal/storage/elasticsearch/elasticsearch_client.go:Bulk(ctx context.Context, request *model.BulkRequest)](../internal/storage/elasticsearch/elasticsearch_client.go)
  - [x] `deletePit(String scrollId)` -> [internal/storage/elasticsearch/elasticsearch_client.go:DeletePit(ctx context.Context, request *model.ClosePointInTimeRequest)](../internal/storage/elasticsearch/elasticsearch_client.go)

- **ElasticsearchManager.java** ([java/im/turms/service/storage/elasticsearch/ElasticsearchManager.java](../turms-orig/turms-service/src/main/java/im/turms/service/storage/elasticsearch/ElasticsearchManager.java))
> [简述功能]

  - [x] `putUserDoc(Long userId, String name)` -> [internal/storage/elasticsearch/elasticsearch_manager.go:PutUserDoc(ctx context.Context)](../internal/storage/elasticsearch/elasticsearch_manager.go)
  - [x] `putUserDocs(Collection<Long> userIds, String name)` -> [internal/storage/elasticsearch/elasticsearch_manager.go:PutUserDocs(ctx context.Context)](../internal/storage/elasticsearch/elasticsearch_manager.go)
  - [x] `deleteUserDoc(Long userId)` -> [internal/storage/elasticsearch/elasticsearch_manager.go:DeleteUserDoc(ctx context.Context)](../internal/storage/elasticsearch/elasticsearch_manager.go)
  - [x] `deleteUserDocs(Collection<Long> userIds)` -> [internal/storage/elasticsearch/elasticsearch_manager.go:DeleteUserDocs(ctx context.Context)](../internal/storage/elasticsearch/elasticsearch_manager.go)
  - [x] `searchUserDocs(@Nullable Integer from, @Nullable Integer size, String name, @Nullable Collection<Long> ids, boolean highlight, @Nullable String scrollId, @Nullable String keepAlive)` -> [internal/storage/elasticsearch/elasticsearch_manager.go:SearchUserDocs(ctx context.Context)](../internal/storage/elasticsearch/elasticsearch_manager.go)
  - [x] `putGroupDoc(Long groupId, String name)` -> [internal/storage/elasticsearch/elasticsearch_manager.go:PutGroupDoc(ctx context.Context)](../internal/storage/elasticsearch/elasticsearch_manager.go)
  - [x] `putGroupDocs(Collection<Long> groupIds, String name)` -> [internal/storage/elasticsearch/elasticsearch_manager.go:PutGroupDocs(ctx context.Context)](../internal/storage/elasticsearch/elasticsearch_manager.go)
  - [x] `deleteGroupDocs(Collection<Long> groupIds)` -> [internal/storage/elasticsearch/elasticsearch_manager.go:DeleteGroupDocs(ctx context.Context)](../internal/storage/elasticsearch/elasticsearch_manager.go)
  - [x] `deleteAllGroupDocs()` -> [internal/storage/elasticsearch/elasticsearch_manager.go:DeleteAllGroupDocs(ctx context.Context)](../internal/storage/elasticsearch/elasticsearch_manager.go)
  - [x] `searchGroupDocs(@Nullable Integer from, @Nullable Integer size, String name, @Nullable Collection<Long> ids, boolean highlight, @Nullable String scrollId, @Nullable String keepAlive)` -> [internal/storage/elasticsearch/elasticsearch_manager.go:SearchGroupDocs(ctx context.Context)](../internal/storage/elasticsearch/elasticsearch_manager.go)
  - [x] `deletePitForUserDocs(String scrollId)` -> [internal/storage/elasticsearch/elasticsearch_manager.go:DeletePitForUserDocs(ctx context.Context)](../internal/storage/elasticsearch/elasticsearch_manager.go)

- **IndexTextFieldSetting.java** ([java/im/turms/service/storage/elasticsearch/IndexTextFieldSetting.java](../turms-orig/turms-service/src/main/java/im/turms/service/storage/elasticsearch/IndexTextFieldSetting.java))
> [简述功能]

  - [ ] `IndexTextFieldSetting(Map<String, Property> fieldToProperty, @Nullable IndexSettingsAnalysis analysis)`

- **BulkRequest.java** ([java/im/turms/service/storage/elasticsearch/model/BulkRequest.java](../turms-orig/turms-service/src/main/java/im/turms/service/storage/elasticsearch/model/BulkRequest.java))
> [简述功能]

  - [ ] `BulkRequest(List<Object> operations)`
  - [x] `serialize(BulkRequest value, JsonGenerator gen, SerializerProvider serializers)` -> [internal/storage/elasticsearch/model/elasticsearch_model.go:Serialize()](../internal/storage/elasticsearch/model/elasticsearch_model.go)

- **BulkResponse.java** ([java/im/turms/service/storage/elasticsearch/model/BulkResponse.java](../turms-orig/turms-service/src/main/java/im/turms/service/storage/elasticsearch/model/BulkResponse.java))
> [简述功能]

  - [ ] `BulkResponse(@JsonProperty("errors")`

- **BulkResponseItem.java** ([java/im/turms/service/storage/elasticsearch/model/BulkResponseItem.java](../turms-orig/turms-service/src/main/java/im/turms/service/storage/elasticsearch/model/BulkResponseItem.java))
> [简述功能]

  - [ ] `BulkResponseItem(@JsonProperty("_id")`

- **ClosePointInTimeRequest.java** ([java/im/turms/service/storage/elasticsearch/model/ClosePointInTimeRequest.java](../turms-orig/turms-service/src/main/java/im/turms/service/storage/elasticsearch/model/ClosePointInTimeRequest.java))
> [简述功能]

  - [ ] `ClosePointInTimeRequest(@JsonProperty("id")`

- **CreateIndexRequest.java** ([java/im/turms/service/storage/elasticsearch/model/CreateIndexRequest.java](../turms-orig/turms-service/src/main/java/im/turms/service/storage/elasticsearch/model/CreateIndexRequest.java))
> [简述功能]

  - [ ] `CreateIndexRequest(@JsonProperty("mappings")`

- **DeleteByQueryRequest.java** ([java/im/turms/service/storage/elasticsearch/model/DeleteByQueryRequest.java](../turms-orig/turms-service/src/main/java/im/turms/service/storage/elasticsearch/model/DeleteByQueryRequest.java))
> [简述功能]

  - [ ] `DeleteByQueryRequest(@JsonProperty("query")`

- **DeleteByQueryResponse.java** ([java/im/turms/service/storage/elasticsearch/model/DeleteByQueryResponse.java](../turms-orig/turms-service/src/main/java/im/turms/service/storage/elasticsearch/model/DeleteByQueryResponse.java))
> [简述功能]

  - [ ] `DeleteByQueryResponse(@JsonProperty("deleted")`

- **DeleteResponse.java** ([java/im/turms/service/storage/elasticsearch/model/DeleteResponse.java](../turms-orig/turms-service/src/main/java/im/turms/service/storage/elasticsearch/model/DeleteResponse.java))
> [简述功能]

  - [ ] `DeleteResponse(@JsonProperty("result")`

- **ErrorCause.java** ([java/im/turms/service/storage/elasticsearch/model/ErrorCause.java](../turms-orig/turms-service/src/main/java/im/turms/service/storage/elasticsearch/model/ErrorCause.java))
> [简述功能]

  - [ ] `ErrorCause(@JsonProperty("type")`

- **ErrorResponse.java** ([java/im/turms/service/storage/elasticsearch/model/ErrorResponse.java](../turms-orig/turms-service/src/main/java/im/turms/service/storage/elasticsearch/model/ErrorResponse.java))
> [简述功能]

  - [ ] `ErrorResponse(@JsonProperty("error")`

- **FieldCollapse.java** ([java/im/turms/service/storage/elasticsearch/model/FieldCollapse.java](../turms-orig/turms-service/src/main/java/im/turms/service/storage/elasticsearch/model/FieldCollapse.java))
> [简述功能]

  - [ ] `FieldCollapse(@JsonProperty("field")`

- **HealthResponse.java** ([java/im/turms/service/storage/elasticsearch/model/HealthResponse.java](../turms-orig/turms-service/src/main/java/im/turms/service/storage/elasticsearch/model/HealthResponse.java))
> [简述功能]

  - [ ] `HealthResponse(@JsonProperty("cluster_name")`

- **Highlight.java** ([java/im/turms/service/storage/elasticsearch/model/Highlight.java](../turms-orig/turms-service/src/main/java/im/turms/service/storage/elasticsearch/model/Highlight.java))
> [简述功能]

  - [ ] `Highlight(@JsonProperty("fields")`

- **IndexSettings.java** ([java/im/turms/service/storage/elasticsearch/model/IndexSettings.java](../turms-orig/turms-service/src/main/java/im/turms/service/storage/elasticsearch/model/IndexSettings.java))
> [简述功能]

  - [ ] `IndexSettings(@JsonProperty("index")`

- **IndexSettingsAnalysis.java** ([java/im/turms/service/storage/elasticsearch/model/IndexSettingsAnalysis.java](../turms-orig/turms-service/src/main/java/im/turms/service/storage/elasticsearch/model/IndexSettingsAnalysis.java))
> [简述功能]

  - [ ] `IndexSettingsAnalysis(@JsonProperty("analyzer")`
  - [x] `merge(IndexSettingsAnalysis analysis)` -> [internal/storage/elasticsearch/model/elasticsearch_model.go:Merge(other *IndexSettingsAnalysis)](../internal/storage/elasticsearch/model/elasticsearch_model.go)

- **PointInTimeReference.java** ([java/im/turms/service/storage/elasticsearch/model/PointInTimeReference.java](../turms-orig/turms-service/src/main/java/im/turms/service/storage/elasticsearch/model/PointInTimeReference.java))
> [简述功能]

  - [ ] `PointInTimeReference(String id, @Nullable String keepAlive)`

- **Property.java** ([java/im/turms/service/storage/elasticsearch/model/Property.java](../turms-orig/turms-service/src/main/java/im/turms/service/storage/elasticsearch/model/Property.java))
> [简述功能]

  - [ ] `Property(@JsonProperty("type")`

- **Script.java** ([java/im/turms/service/storage/elasticsearch/model/Script.java](../turms-orig/turms-service/src/main/java/im/turms/service/storage/elasticsearch/model/Script.java))
> [简述功能]

  - [ ] `Script(@JsonProperty("source")`

- **SearchRequest.java** ([java/im/turms/service/storage/elasticsearch/model/SearchRequest.java](../turms-orig/turms-service/src/main/java/im/turms/service/storage/elasticsearch/model/SearchRequest.java))
> [简述功能]

  - [ ] `SearchRequest(@JsonProperty("from")`

- **ShardFailure.java** ([java/im/turms/service/storage/elasticsearch/model/ShardFailure.java](../turms-orig/turms-service/src/main/java/im/turms/service/storage/elasticsearch/model/ShardFailure.java))
> [简述功能]

  - [ ] `ShardFailure(@JsonProperty("index")`

- **ShardStatistics.java** ([java/im/turms/service/storage/elasticsearch/model/ShardStatistics.java](../turms-orig/turms-service/src/main/java/im/turms/service/storage/elasticsearch/model/ShardStatistics.java))
> [简述功能]

  - [ ] `ShardStatistics(@JsonProperty("failed")`

- **TypeMapping.java** ([java/im/turms/service/storage/elasticsearch/model/TypeMapping.java](../turms-orig/turms-service/src/main/java/im/turms/service/storage/elasticsearch/model/TypeMapping.java))
> [简述功能]

  - [ ] `TypeMapping(@JsonProperty("dynamic")`

- **UpdateByQueryRequest.java** ([java/im/turms/service/storage/elasticsearch/model/UpdateByQueryRequest.java](../turms-orig/turms-service/src/main/java/im/turms/service/storage/elasticsearch/model/UpdateByQueryRequest.java))
> [简述功能]

  - [ ] `UpdateByQueryRequest(@JsonProperty("query")`

- **UpdateByQueryResponse.java** ([java/im/turms/service/storage/elasticsearch/model/UpdateByQueryResponse.java](../turms-orig/turms-service/src/main/java/im/turms/service/storage/elasticsearch/model/UpdateByQueryResponse.java))
> [简述功能]

  - [ ] `UpdateByQueryResponse(@JsonProperty("updated")`

- **MongoCollectionMigrator.java** ([java/im/turms/service/storage/mongo/MongoCollectionMigrator.java](../turms-orig/turms-service/src/main/java/im/turms/service/storage/mongo/MongoCollectionMigrator.java))
> [简述功能]

  - [x] `migrate(Set<String> existingCollectionNames)` -> [internal/storage/mongo/mongo_collection_migrator.go:Migrate()](../internal/storage/mongo/mongo_collection_migrator.go)

- **MongoConfig.java** ([java/im/turms/service/storage/mongo/MongoConfig.java](../turms-orig/turms-service/src/main/java/im/turms/service/storage/mongo/MongoConfig.java))
> [简述功能]

  - [x] `adminMongoClient(TurmsPropertiesManager propertiesManager)` -> [internal/storage/mongo/mongo_config.go:AdminMongoClient()](../internal/storage/mongo/mongo_config.go)
  - [x] `userMongoClient(TurmsPropertiesManager propertiesManager)` -> [internal/storage/mongo/mongo_config.go:UserMongoClient()](../internal/storage/mongo/mongo_config.go)
  - [x] `groupMongoClient(TurmsPropertiesManager propertiesManager)` -> [internal/storage/mongo/mongo_config.go:GroupMongoClient()](../internal/storage/mongo/mongo_config.go)
  - [x] `conversationMongoClient(TurmsPropertiesManager propertiesManager)` -> [internal/storage/mongo/mongo_config.go:ConversationMongoClient()](../internal/storage/mongo/mongo_config.go)
  - [x] `messageMongoClient(TurmsPropertiesManager propertiesManager)` -> [internal/storage/mongo/mongo_config.go:MessageMongoClient()](../internal/storage/mongo/mongo_config.go)
  - [x] `conferenceMongoClient(TurmsPropertiesManager propertiesManager)` -> [internal/storage/mongo/mongo_config.go:ConferenceMongoClient()](../internal/storage/mongo/mongo_config.go)

- **MongoFakeDataGenerator.java** ([java/im/turms/service/storage/mongo/MongoFakeDataGenerator.java](../turms-orig/turms-service/src/main/java/im/turms/service/storage/mongo/MongoFakeDataGenerator.java))
> [简述功能]

  - [x] `populateCollectionsWithFakeData()` -> [internal/storage/mongo/mongo_fake_data_generator.go:PopulateCollectionsWithFakeData()](../internal/storage/mongo/mongo_fake_data_generator.go)

- **RedisConfig.java** ([java/im/turms/service/storage/redis/RedisConfig.java](../turms-orig/turms-service/src/main/java/im/turms/service/storage/redis/RedisConfig.java))
> [简述功能]

  - [x] `newSequenceIdRedisClientManager(RedisProperties properties)` -> [internal/storage/redis/redis_config.go:NewSequenceIdRedisClientManager()](../internal/storage/redis/redis_config.go)
  - [x] `sequenceIdRedisClientManager()` -> [internal/storage/redis/redis_config.go:SequenceIdRedisClientManager()](../internal/storage/redis/redis_config.go)

