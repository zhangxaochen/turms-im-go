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

  - [x] `main(String[] args)` -> [main()](../cmd/turms-gateway/main.go#L8)

- **ClientRequestDispatcher.java** ([java/im/turms/gateway/access/client/common/ClientRequestDispatcher.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/common/ClientRequestDispatcher.java))
> [简述功能]

  - [x] `handleRequest(UserSessionWrapper sessionWrapper, ByteBuf serviceRequestBuffer)` -> [HandleRequest(ctx context.Context, sessionWrapper *UserSessionWrapper, serviceRequestBuffer []byte)](../internal/domain/gateway/access/client/common/client_request_dispatcher.go#L87)
  - [x] `handleRequest0(UserSessionWrapper sessionWrapper, ByteBuf serviceRequestBuffer)` -> [HandleRequest0(ctx context.Context, sessionWrapper *UserSessionWrapper, serviceRequestBuffer []byte)](../internal/domain/gateway/access/client/common/client_request_dispatcher.go#L98)
  - [x] `handleServiceRequest(UserSessionWrapper sessionWrapper, SimpleTurmsRequest request, ByteBuf serviceRequestBuffer, TracingContext tracingContext)` -> [HandleServiceRequest(ctx context.Context, sessionWrapper *UserSessionWrapper, request *protocol.TurmsRequest, serviceRequestBuffer []byte)](../internal/domain/gateway/access/client/common/client_request_dispatcher.go#L176)

- **IpRequestThrottler.java** ([java/im/turms/gateway/access/client/common/IpRequestThrottler.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/common/IpRequestThrottler.java))
> [简述功能]

  - [x] `tryAcquireToken(ByteArrayWrapper ip, long timestamp)` -> [TryAcquireToken(ip string)](../internal/domain/gateway/access/client/common/ip_request_throttler.go#L40)

- **NotificationFactory.java** ([java/im/turms/gateway/access/client/common/NotificationFactory.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/common/NotificationFactory.java))
> [简述功能]

  - [x] `init(TurmsPropertiesManager propertiesManager)` -> [NewNotificationFactory(props *config.GatewayProperties)](../internal/domain/gateway/access/client/common/notification_factory.go#L22)
  - [x] `create(ResponseStatusCode code, long requestId)` -> [Create(requestID *int64, code constant.ResponseStatusCode)](../internal/domain/gateway/access/client/common/notification_factory.go#L33)
  - [x] `create(ResponseStatusCode code, @Nullable String reason, long requestId)` -> [CreateWithReason(requestID *int64, code constant.ResponseStatusCode, reason string)](../internal/domain/gateway/access/client/common/notification_factory.go#L39)
  - [x] `create(ThrowableInfo info, long requestId)` -> [CreateFromError(err error, requestID *int64)](../internal/domain/gateway/access/client/common/notification_factory.go#L52)
  - [x] `createBuffer(CloseReason closeReason)` -> [CreateBuffer(requestID *int64, code constant.ResponseStatusCode, reason string)](../internal/domain/gateway/access/client/common/notification_factory.go#L75)
  - [x] `sessionClosed(long requestId)` -> [SessionClosed(requestID *int64)](../internal/domain/gateway/access/client/common/notification_factory.go#L82)

- **RequestHandlerResult.java** ([java/im/turms/gateway/access/client/common/RequestHandlerResult.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/common/RequestHandlerResult.java))
> [简述功能]

  - [x] `RequestHandlerResult(ResponseStatusCode code, String reason)` -> [NewRequestHandlerResult(code constant.ResponseStatusCode, reason string)](../internal/domain/gateway/access/client/common/request_handler_result.go#L13)

- **UserSession.java** ([java/im/turms/gateway/access/client/common/UserSession.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/common/UserSession.java))
> [简述功能]

  - [x] `setConnection(NetConnection connection, ByteArrayWrapper ip)` -> [SetConnection(connection Connection, ip string)](../internal/domain/gateway/session/connection.go#L92)
  - [x] `setLastHeartbeatRequestTimestampToNow()` -> [SetLastHeartbeatRequestTimestampToNow()](../internal/domain/gateway/session/connection.go#L41)
  - [x] `setLastRequestTimestampToNow()` -> [SetLastRequestTimestampToNow()](../internal/domain/gateway/session/connection.go#L52)
  - [x] `close(@NotNull CloseReason closeReason)` -> [Close()](../internal/domain/common/cache/ttl_cache.go#L72), [Close(closeReason any)](../internal/domain/gateway/session/connection.go#L125)
  - [x] `isOpen()` -> [IsOpen()](../internal/domain/gateway/session/connection.go#L63)
  - [x] `isConnected()` -> [IsConnected()](../internal/domain/gateway/session/connection.go#L97)
  - [x] `supportsSwitchingToUdp()` -> [SupportsSwitchingToUdp()](../internal/domain/gateway/session/connection.go#L102)
  - [x] `sendNotification(ByteBuf byteBuf)` -> [sendNotification(s *session.UserSession, requestID *int64, code int32, reason string)](../internal/domain/gateway/access/router/router.go#L135)
  - [x] `sendNotification(ByteBuf byteBuf, TracingContext tracingContext)` -> [sendNotification(s *session.UserSession, requestID *int64, code int32, reason string)](../internal/domain/gateway/access/router/router.go#L135)
  - [x] `acquireDeleteSessionRequestLoggingLock()` -> [AcquireDeleteSessionRequestLoggingLock()](../internal/domain/gateway/session/connection.go#L112)
  - [x] `hasPermission(TurmsRequest.KindCase requestType)` -> [HasPermission(requestType any)](../internal/domain/gateway/session/connection.go#L117)
  - [x] `toString()` -> [ToString()](../internal/domain/gateway/session/connection.go#L107), [ToString()](../internal/infra/proto/proto_parser.go#L38)

- **UserSessionWrapper.java** ([java/im/turms/gateway/access/client/common/UserSessionWrapper.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/common/UserSessionWrapper.java))
> [简述功能]

  - [x] `getIp()` -> [GetIP()](../internal/domain/gateway/access/client/common/user_session_wrapper.go#L15)
  - [x] `getIpStr()` -> [GetIPStr()](../internal/domain/gateway/access/client/common/user_session_wrapper.go#L20)
  - [x] `setUserSession(UserSession userSession)` -> [SetUserSession(userSession *session.UserSession)](../internal/domain/gateway/access/client/common/user_session_wrapper.go#L25)
  - [x] `hasUserSession()` -> [HasUserSession()](../internal/domain/gateway/access/client/common/user_session_wrapper.go#L30)

- **Policy.java** ([java/im/turms/gateway/access/client/common/authorization/policy/Policy.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/common/authorization/policy/Policy.java))
> [简述功能]

  - [x] `Policy(List<PolicyStatement> statements)` -> [NewPolicy(statements []PolicyStatement)](../internal/domain/gateway/access/client/common/authorization/policy.go#L155)

- **PolicyDeserializer.java** ([java/im/turms/gateway/access/client/common/authorization/policy/PolicyDeserializer.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/common/authorization/policy/PolicyDeserializer.java))
> [简述功能]

  - [x] `parse(Map<String, Object> map)` -> [Parse(data map[string]interface{})](../internal/domain/gateway/access/client/common/authorization/policy.go#L166)

- **PolicyStatement.java** ([java/im/turms/gateway/access/client/common/authorization/policy/PolicyStatement.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/common/authorization/policy/PolicyStatement.java))
> [简述功能]

  - [x] `PolicyStatement(PolicyStatementEffect effect, Set<PolicyStatementAction> actions, Set<PolicyStatementResource> resources)` -> [NewPolicyStatement(effect PolicyStatementEffect, actions []PolicyStatementAction, resources []PolicyStatementResource)](../internal/domain/gateway/access/client/common/authorization/policy.go#L145)

- **ServiceAvailabilityHandler.java** ([java/im/turms/gateway/access/client/common/channel/ServiceAvailabilityHandler.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/common/channel/ServiceAvailabilityHandler.java))
> [简述功能]

  - [x] `channelRegistered(ChannelHandlerContext ctx)` -> [ChannelRegistered(isAvailable bool)](../internal/domain/gateway/access/client/common/service_availability.go#L43)
  - [x] `exceptionCaught(ChannelHandlerContext ctx, Throwable cause)` -> [ExceptionCaught(err error)](../internal/domain/gateway/access/client/common/service_availability.go#L48)

- **NetConnection.java** ([java/im/turms/gateway/access/client/common/connection/NetConnection.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/common/connection/NetConnection.java))
> [简述功能]

  - [x] `getAddress()` -> [GetAddress()](../internal/domain/gateway/access/client/tcp/tcp_server.go#L31), [GetAddress()](../internal/domain/gateway/access/client/ws/websocket_connection.go#L30)
  - [x] `send(ByteBuf buffer)` -> [Send(ctx context.Context, buffer []byte)](../internal/domain/gateway/access/client/tcp/tcp_server.go#L36), [Send(ctx context.Context, buffer []byte)](../internal/domain/gateway/access/client/ws/websocket_connection.go#L35)
  - [x] `close(CloseReason closeReason)` -> [Close()](../internal/domain/common/cache/ttl_cache.go#L72), [CloseWithReason(reason CloseReason)](../internal/domain/gateway/access/client/common/net_connection.go#L58), [CloseWithReason(reason common.CloseReason)](../internal/domain/gateway/access/client/tcp/tcp_server.go#L43), [CloseWithReason(reason common.CloseReason)](../internal/domain/gateway/access/client/ws/websocket_connection.go#L44)
  - [x] `close()` -> [Close()](../internal/domain/common/cache/ttl_cache.go#L72), [Close()](../internal/domain/gateway/access/client/common/net_connection.go#L68), [Close()](../internal/domain/gateway/access/client/tcp/tcp_server.go#L82), [Close()](../internal/domain/gateway/access/client/ws/websocket_connection.go#L87)
  - [x] `switchToUdp()` -> [SwitchToUdp()](../internal/domain/gateway/access/client/common/net_connection.go#L78)
  - [x] `tryNotifyClientToRecover()` -> [TryNotifyClientToRecover()](../internal/domain/gateway/access/client/common/net_connection.go#L83)

- **ExtendedHAProxyMessageReader.java** ([java/im/turms/gateway/access/client/tcp/ExtendedHAProxyMessageReader.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/tcp/ExtendedHAProxyMessageReader.java))
> [简述功能]

  - [x] `channelRead(ChannelHandlerContext ctx, Object msg)` -> [Read(conn net.Conn)](../internal/domain/gateway/access/client/tcp/haproxy.go#L33)

- **HAProxyUtil.java** ([java/im/turms/gateway/access/client/tcp/HAProxyUtil.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/tcp/HAProxyUtil.java))
> [简述功能]

  - [x] `addProxyProtocolHandlers(ChannelPipeline pipeline, Consumer<InetSocketAddress> onRemoteAddressConfirmed)` -> [AddProxyProtocolHandlers(callback func(net.Addr)](../internal/domain/gateway/access/client/tcp/haproxy.go#L44)
  - [x] `addProxyProtocolDetectorHandler(ChannelPipeline pipeline, Consumer<InetSocketAddress> onRemoteAddressConfirmed)` -> [AddProxyProtocolDetectorHandler(callback func(net.Addr)](../internal/domain/gateway/access/client/tcp/haproxy.go#L49)

- **TcpConnection.java** ([java/im/turms/gateway/access/client/tcp/TcpConnection.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/tcp/TcpConnection.java))
> [简述功能]

  - [x] `getAddress()` -> [GetAddress()](../internal/domain/gateway/access/client/tcp/tcp_server.go#L31), [GetAddress()](../internal/domain/gateway/access/client/ws/websocket_connection.go#L30)
  - [x] `send(ByteBuf buffer)` -> [Send(ctx context.Context, buffer []byte)](../internal/domain/gateway/access/client/tcp/tcp_server.go#L36), [Send(ctx context.Context, buffer []byte)](../internal/domain/gateway/access/client/ws/websocket_connection.go#L35)
  - [x] `close(CloseReason closeReason)` -> [Close()](../internal/domain/common/cache/ttl_cache.go#L72), [CloseWithReason(reason CloseReason)](../internal/domain/gateway/access/client/common/net_connection.go#L58), [CloseWithReason(reason common.CloseReason)](../internal/domain/gateway/access/client/tcp/tcp_server.go#L43), [CloseWithReason(reason common.CloseReason)](../internal/domain/gateway/access/client/ws/websocket_connection.go#L44)
  - [x] `close()` -> [Close()](../internal/domain/common/cache/ttl_cache.go#L72), [Close()](../internal/domain/gateway/access/client/common/net_connection.go#L68), [Close()](../internal/domain/gateway/access/client/tcp/tcp_server.go#L82), [Close()](../internal/domain/gateway/access/client/ws/websocket_connection.go#L87)

- **TcpServerFactory.java** ([java/im/turms/gateway/access/client/tcp/TcpServerFactory.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/tcp/TcpServerFactory.java))
> [简述功能]

  - [ ] `create(TcpProperties tcpProperties, BlocklistService blocklistService, ServerStatusManager serverStatusManager, SessionService sessionService, ConnectionListener connectionListener, int maxFrameLength)`

- **TcpUserSessionAssembler.java** ([java/im/turms/gateway/access/client/tcp/TcpUserSessionAssembler.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/tcp/TcpUserSessionAssembler.java))
> [简述功能]

  - [x] `getHost()` -> [GetHost()](../internal/domain/gateway/access/client/tcp/tcp_server.go#L137)
  - [x] `getPort()` -> [GetPort()](../internal/domain/gateway/access/client/tcp/tcp_server.go#L145)

- **UdpRequestDispatcher.java** ([java/im/turms/gateway/access/client/udp/UdpRequestDispatcher.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/udp/UdpRequestDispatcher.java))
> [简述功能]

  - [x] `sendSignal(InetSocketAddress address, UdpNotificationType signal)` -> [SendSignal(address net.Addr, signal UdpNotificationType)](../internal/domain/gateway/access/client/udp/udp_server.go#L143)

- **UdpSignalResponseBufferPool.java** ([java/im/turms/gateway/access/client/udp/UdpSignalResponseBufferPool.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/udp/UdpSignalResponseBufferPool.java))
> [简述功能]

  - [x] `get(ResponseStatusCode code)` -> [Get(key K)](../internal/domain/common/cache/sharded_map.go#L53), [GetBufferFromStatusCode(code constant.ResponseStatusCode)](../internal/domain/gateway/access/client/udp/udp_server.go#L211)
  - [x] `get(UdpNotificationType type)` -> [Get(key K)](../internal/domain/common/cache/sharded_map.go#L53), [GetBufferFromNotificationType(t UdpNotificationType)](../internal/domain/gateway/access/client/udp/udp_server.go#L221)

- **UdpNotification.java** ([java/im/turms/gateway/access/client/udp/dto/UdpNotification.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/udp/dto/UdpNotification.java))
> [简述功能]

  - [x] `UdpNotification(InetSocketAddress recipientAddress, UdpNotificationType type)` -> [NewUdpNotification(recipientAddress net.Addr, notificationType UdpNotificationType)](../internal/domain/gateway/access/client/udp/udp_server.go#L37)

- **UdpRequestType.java** ([java/im/turms/gateway/access/client/udp/dto/UdpRequestType.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/udp/dto/UdpRequestType.java))
> [简述功能]

  - [x] `parse(int number)` -> [ParseUdpRequestType(number int)](../internal/domain/gateway/access/client/udp/udp_server.go#L63)
  - [x] `getNumber()` -> [GetNumber()](../internal/domain/gateway/access/client/udp/udp_server.go#L72)

- **UdpSignalRequest.java** ([java/im/turms/gateway/access/client/udp/dto/UdpSignalRequest.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/udp/dto/UdpSignalRequest.java))
> [简述功能]

  - [x] `UdpSignalRequest(UdpRequestType type, long userId, DeviceType deviceType, int sessionId)` -> [NewUdpSignalRequest(reqType UdpRequestType, userID int64, deviceType protocol.DeviceType, sessionID int)](../internal/domain/gateway/access/client/udp/udp_server.go#L53)

- **HttpForwardedHeaderHandler.java** ([java/im/turms/gateway/access/client/websocket/HttpForwardedHeaderHandler.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/websocket/HttpForwardedHeaderHandler.java))
> [简述功能]

  - [x] `apply(ConnectionInfo connectionInfo, HttpRequest request)` -> [Apply(connectionInfo any, request any)](../internal/domain/gateway/access/client/ws/ws_server.go#L12)

- **WebSocketConnection.java** ([java/im/turms/gateway/access/client/websocket/WebSocketConnection.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/websocket/WebSocketConnection.java))
> [简述功能]

  - [x] `getAddress()` -> [GetAddress()](../internal/domain/gateway/access/client/tcp/tcp_server.go#L31), [GetAddress()](../internal/domain/gateway/access/client/ws/websocket_connection.go#L30)
  - [x] `send(ByteBuf buffer)` -> [Send(ctx context.Context, buffer []byte)](../internal/domain/gateway/access/client/tcp/tcp_server.go#L36), [Send(ctx context.Context, buffer []byte)](../internal/domain/gateway/access/client/ws/websocket_connection.go#L35)
  - [x] `close(CloseReason closeReason)` -> [Close()](../internal/domain/common/cache/ttl_cache.go#L72), [CloseWithReason(reason CloseReason)](../internal/domain/gateway/access/client/common/net_connection.go#L58), [CloseWithReason(reason common.CloseReason)](../internal/domain/gateway/access/client/tcp/tcp_server.go#L43), [CloseWithReason(reason common.CloseReason)](../internal/domain/gateway/access/client/ws/websocket_connection.go#L44)
  - [x] `close()` -> [Close()](../internal/domain/common/cache/ttl_cache.go#L72), [Close()](../internal/domain/gateway/access/client/common/net_connection.go#L68), [Close()](../internal/domain/gateway/access/client/tcp/tcp_server.go#L82), [Close()](../internal/domain/gateway/access/client/ws/websocket_connection.go#L87)

- **WebSocketServerFactory.java** ([java/im/turms/gateway/access/client/websocket/WebSocketServerFactory.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/access/client/websocket/WebSocketServerFactory.java))
> [简述功能]

  - [ ] `create(WebSocketProperties webSocketProperties, BlocklistService blocklistService, ServerStatusManager serverStatusManager, SessionService sessionService, ConnectionListener connectionListener, int maxFramePayloadLength)`

- **NotificationService.java** ([java/im/turms/gateway/domain/notification/service/NotificationService.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/domain/notification/service/NotificationService.java))
> [简述功能]

  - [x] `sendNotificationToLocalClients(TracingContext tracingContext, ByteBuf notificationData, Set<Long> recipientIds, Set<UserSessionId> excludedUserSessionIds, @Nullable DeviceType excludedDeviceType)` -> [SendNotificationToLocalClients(ctx context.Context, notificationData []byte, recipientIds []int64, excludedUserSessionIds map[UserSessionID]struct{}, excludedDeviceType *protocol.DeviceType)](../internal/domain/gateway/notification/service/notification_service.go#L37)

- **StatisticsService.java** ([java/im/turms/gateway/domain/observation/service/StatisticsService.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/domain/observation/service/StatisticsService.java))
> [简述功能]

  - [x] `countLocalOnlineUsers()` -> [CountLocalOnlineUsers()](../internal/domain/gateway/observation/service/statistics_service.go#L19)

- **ServiceRequestService.java** ([java/im/turms/gateway/domain/servicerequest/service/ServiceRequestService.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/domain/servicerequest/service/ServiceRequestService.java))
> [简述功能]

  - [x] `handleServiceRequest(UserSession session, ServiceRequest serviceRequest)` -> [HandleServiceRequest(ctx context.Context, defaultSession *session.UserSession, serviceRequest any)](../internal/domain/gateway/servicerequest/service/servicerequest_service.go#L22)

- **SessionController.java** ([java/im/turms/gateway/domain/session/access/admin/controller/SessionController.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/domain/session/access/admin/controller/SessionController.java))
> [简述功能]

  - [ ] `deleteSessions(@QueryParam(required = false)`

- **SessionClientController.java** ([java/im/turms/gateway/domain/session/access/client/controller/SessionClientController.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/domain/session/access/client/controller/SessionClientController.java))
> [简述功能]

  - [x] `handleDeleteSessionRequest(UserSessionWrapper sessionWrapper)` -> [HandleDeleteSessionRequest(ctx context.Context, sessionWrapper *common.UserSessionWrapper)](../internal/domain/gateway/session/access/client/controller/session_client_controller.go#L25)
  - [x] `handleCreateSessionRequest(UserSessionWrapper sessionWrapper, CreateSessionRequest createSessionRequest)` -> [HandleCreateSessionRequest(ctx context.Context, sessionWrapper *common.UserSessionWrapper, req *protocol.CreateSessionRequest)](../internal/domain/gateway/session/access/client/controller/session_client_controller.go#L37)

- **UserLoginInfo.java** ([java/im/turms/gateway/domain/session/bo/UserLoginInfo.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/domain/session/bo/UserLoginInfo.java))
> [简述功能]

  - [x] `UserLoginInfo(int version, Long userId, String password, DeviceType loggingInDeviceType, Map<String, String> deviceDetails, UserStatus userStatus, Location location, String ip)` -> [NewUserLoginInfo(version int, userID int64, password *string, loggingInDeviceType protocol.DeviceType, deviceDetails map[string]string, userStatus *protocol.UserStatus /*, location *Location */, ip string)](../internal/domain/gateway/session/bo/user_login_info.go#L21)

- **UserPermissionInfo.java** ([java/im/turms/gateway/domain/session/bo/UserPermissionInfo.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/domain/session/bo/UserPermissionInfo.java))
> [简述功能]

  - [x] `UserPermissionInfo(ResponseStatusCode authenticationCode, Set<TurmsRequest.KindCase> permissions)` -> [NewUserPermissionInfo(authenticationCode constant.ResponseStatusCode, permissions map[any]bool)](../internal/domain/gateway/session/bo/user_permission_info.go#L15)

- **HeartbeatManager.java** ([java/im/turms/gateway/domain/session/manager/HeartbeatManager.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/domain/session/manager/HeartbeatManager.java))
> [简述功能]

  - [ ] `setCloseIdleSessionAfterSeconds(int closeIdleSessionAfterSeconds)`
  - [ ] `setClientHeartbeatIntervalSeconds(int clientHeartbeatIntervalSeconds)`
  - [x] `destroy()` -> [Destroy(ctx context.Context)](../internal/domain/gateway/session/service.go#L193)
  - [ ] `estimatedSize()`
  - [ ] `next()`

- **UserSessionsManager.java** ([java/im/turms/gateway/domain/session/manager/UserSessionsManager.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/domain/session/manager/UserSessionsManager.java))
> [简述功能]

  - [ ] `addSessionIfAbsent(int version, Set<TurmsRequest.KindCase> permissions, DeviceType loggingInDeviceType, Map<String, String> deviceDetails, @Nullable Location location)`
  - [ ] `closeSession(@NotNull DeviceType deviceType, @NotNull CloseReason closeReason)`
  - [ ] `pushSessionNotification(DeviceType deviceType, String serverId)`
  - [x] `getSession(@NotNull DeviceType deviceType)` -> [GetSession(deviceType protocol.DeviceType)](../internal/domain/gateway/session/sharded_map.go#L31)
  - [ ] `countSessions()`
  - [ ] `getLoggedInDeviceTypes()`

- **UserRepository.java** ([java/im/turms/gateway/domain/session/repository/UserRepository.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/domain/session/repository/UserRepository.java))
> [简述功能]

  - [x] `findPassword(Long userId)` -> [FindPassword(ctx context.Context, userID int64)](../internal/domain/user/repository/user_repository.go#L284)
  - [x] `isActiveAndNotDeleted(Long userId)` -> [IsActiveAndNotDeleted(ctx context.Context, userID int64)](../internal/domain/user/repository/user_repository.go#L271)

- **HttpSessionIdentityAccessManager.java** ([java/im/turms/gateway/domain/session/service/HttpSessionIdentityAccessManager.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/domain/session/service/HttpSessionIdentityAccessManager.java))
> [简述功能]

  - [x] `verifyAndGrant(UserLoginInfo userLoginInfo)` -> [VerifyAndGrant(ctx context.Context, loginInfo *bo.UserLoginInfo)](../internal/domain/gateway/session/identity_access_manager.go#L70)

- **JwtSessionIdentityAccessManager.java** ([java/im/turms/gateway/domain/session/service/JwtSessionIdentityAccessManager.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/domain/session/service/JwtSessionIdentityAccessManager.java))
> [简述功能]

  - [x] `verifyAndGrant(UserLoginInfo userLoginInfo)` -> [VerifyAndGrant(ctx context.Context, loginInfo *bo.UserLoginInfo)](../internal/domain/gateway/session/identity_access_manager.go#L70)

- **LdapSessionIdentityAccessManager.java** ([java/im/turms/gateway/domain/session/service/LdapSessionIdentityAccessManager.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/domain/session/service/LdapSessionIdentityAccessManager.java))
> [简述功能]

  - [x] `verifyAndGrant(UserLoginInfo userLoginInfo)` -> [VerifyAndGrant(ctx context.Context, loginInfo *bo.UserLoginInfo)](../internal/domain/gateway/session/identity_access_manager.go#L70)

- **NoopSessionIdentityAccessManager.java** ([java/im/turms/gateway/domain/session/service/NoopSessionIdentityAccessManager.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/domain/session/service/NoopSessionIdentityAccessManager.java))
> [简述功能]

  - [x] `verifyAndGrant(UserLoginInfo userLoginInfo)` -> [VerifyAndGrant(ctx context.Context, loginInfo *bo.UserLoginInfo)](../internal/domain/gateway/session/identity_access_manager.go#L70)

- **PasswordSessionIdentityAccessManager.java** ([java/im/turms/gateway/domain/session/service/PasswordSessionIdentityAccessManager.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/domain/session/service/PasswordSessionIdentityAccessManager.java))
> [简述功能]

  - [x] `verifyAndGrant(UserLoginInfo userLoginInfo)` -> [VerifyAndGrant(ctx context.Context, loginInfo *bo.UserLoginInfo)](../internal/domain/gateway/session/identity_access_manager.go#L70)
  - [x] `updateGlobalProperties(TurmsProperties properties)` -> [UpdateGlobalProperties(properties interface{})](../internal/domain/gateway/session/identity_access_manager.go#L33), [UpdateGlobalProperties(properties interface{})](../internal/domain/gateway/session/identity_access_manager.go#L102), [UpdateGlobalProperties(properties interface{})](../internal/domain/gateway/session/identity_access_manager.go#L111), [UpdateGlobalProperties(properties interface{})](../internal/domain/gateway/session/identity_access_manager.go#L119), [UpdateGlobalProperties(properties interface{})](../internal/domain/gateway/session/identity_access_manager.go#L127), [UpdateGlobalProperties(properties interface{})](../internal/domain/gateway/session/identity_access_manager.go#L135)

- **SessionIdentityAccessManager.java** ([java/im/turms/gateway/domain/session/service/SessionIdentityAccessManager.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/domain/session/service/SessionIdentityAccessManager.java))
> [简述功能]

  - [x] `verifyAndGrant(int version, Long userId, @Nullable String password, DeviceType deviceType, @Nullable Map<String, String> deviceDetails, @Nullable UserStatus userStatus, @Nullable Location location, @Nullable String ip)` -> [VerifyAndGrant(ctx context.Context, loginInfo *bo.UserLoginInfo)](../internal/domain/gateway/session/identity_access_manager.go#L108), [VerifyAndGrant(ctx context.Context, loginInfo *bo.UserLoginInfo)](../internal/domain/gateway/session/identity_access_manager.go#L116), [VerifyAndGrant(ctx context.Context, loginInfo *bo.UserLoginInfo)](../internal/domain/gateway/session/identity_access_manager.go#L124), [VerifyAndGrant(ctx context.Context, loginInfo *bo.UserLoginInfo)](../internal/domain/gateway/session/identity_access_manager.go#L132)

- **SessionService.java** ([java/im/turms/gateway/domain/session/service/SessionService.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/domain/session/service/SessionService.java))
> [简述功能]

  - [x] `destroy()` -> [Destroy(ctx context.Context)](../internal/domain/gateway/session/service.go#L193)
  - [x] `handleHeartbeatUpdateRequest(UserSession session)` -> [HandleHeartbeatUpdateRequest(session *UserSession)](../internal/domain/gateway/session/service.go#L198)
  - [x] `handleLoginRequest(int version, @NotNull ByteArrayWrapper ip, @NotNull Long userId, @Nullable String password, @NotNull DeviceType deviceType, @Nullable Map<String, String> deviceDetails, @Nullable UserStatus userStatus, @Nullable Location location, @Nullable String ipStr)` -> [HandleLoginRequest(ctx context.Context, version int, ip []byte, userId int64, password string, deviceType protocol.DeviceType, deviceDetails map[string]string, userStatus protocol.UserStatus, location any, ipStr string)](../internal/domain/gateway/session/service.go#L202)
  - [ ] `closeLocalSessions(@NotNull List<byte[]> ips, @NotNull CloseReason closeReason)`
  - [ ] `closeLocalSessions(@NotNull byte[] ip, @NotNull CloseReason closeReason)`
  - [x] `closeLocalSession(@NotNull Long userId, @NotNull @ValidDeviceType DeviceType deviceType, @NotNull SessionCloseStatus closeStatus)` -> [CloseLocalSession(ctx context.Context, userId int64, deviceTypes []protocol.DeviceType, closeReason any)](../internal/domain/gateway/session/service.go#L238)
  - [x] `closeLocalSession(@NotNull Long userId, @NotNull @ValidDeviceType DeviceType deviceType, @NotNull CloseReason closeReason)` -> [CloseLocalSession(ctx context.Context, userId int64, deviceTypes []protocol.DeviceType, closeReason any)](../internal/domain/gateway/session/service.go#L238)
  - [x] `closeLocalSession(@NotNull Long userId, @NotEmpty Set<@ValidDeviceType DeviceType> deviceTypes, @NotNull CloseReason closeReason)` -> [CloseLocalSession(ctx context.Context, userId int64, deviceTypes []protocol.DeviceType, closeReason any)](../internal/domain/gateway/session/service.go#L238)
  - [ ] `closeLocalSessions(@NotNull Set<Long> userIds, @NotNull CloseReason closeReason)`
  - [x] `authAndCloseLocalSession(@NotNull Long userId, @NotNull DeviceType deviceType, @NotNull CloseReason closeReason, int sessionId)` -> [AuthAndCloseLocalSession(ctx context.Context, userId int64, deviceType protocol.DeviceType, closeReason any, sessionId int)](../internal/domain/gateway/session/service.go#L287)
  - [x] `closeAllLocalSessions(@NotNull CloseReason closeReason)` -> [CloseAllLocalSessions(ctx context.Context, closeReason any)](../internal/domain/gateway/session/service.go#L323)
  - [x] `closeLocalSession(Long userId, SessionCloseStatus closeStatus)` -> [CloseLocalSession(ctx context.Context, userId int64, deviceTypes []protocol.DeviceType, closeReason any)](../internal/domain/gateway/session/service.go#L238)
  - [x] `closeLocalSession(Long userId, CloseReason closeReason)` -> [CloseLocalSession(ctx context.Context, userId int64, deviceTypes []protocol.DeviceType, closeReason any)](../internal/domain/gateway/session/service.go#L238)
  - [x] `getSessions(Set<Long> userIds)` -> [GetSessions(ctx context.Context, userIds []int64)](../internal/domain/gateway/session/service.go#L337)
  - [x] `authAndUpdateHeartbeatTimestamp(long userId, @NotNull @ValidDeviceType DeviceType deviceType, int sessionId)` -> [AuthAndUpdateHeartbeatTimestamp(ctx context.Context, userId int64, deviceType protocol.DeviceType, sessionId int)](../internal/domain/gateway/session/service.go#L380)
  - [x] `tryRegisterOnlineUser(int version, @NotNull Set<TurmsRequest.KindCase> permissions, @NotNull ByteArrayWrapper ip, @NotNull Long userId, @NotNull DeviceType deviceType, @Nullable Map<String, String> deviceDetails, @Nullable UserStatus userStatus, @Nullable Location location)` -> [TryRegisterOnlineUser(ctx context.Context, version int, permissions any, ip []byte, userId int64, deviceType protocol.DeviceType, deviceDetails map[string]string, userStatus protocol.UserStatus, location any)](../internal/domain/gateway/session/service.go#L390)
  - [x] `getUserSessionsManager(@NotNull Long userId)` -> [GetUserSessionsManager(ctx context.Context, userId int64)](../internal/domain/gateway/session/service.go#L531)
  - [x] `getLocalUserSession(@NotNull Long userId, @NotNull DeviceType deviceType)` -> [GetLocalUserSession(ctx context.Context, userId int64, deviceType protocol.DeviceType)](../internal/domain/gateway/session/service.go#L539)
  - [x] `getLocalUserSession(ByteArrayWrapper ip)` -> [GetLocalUserSession(ctx context.Context, userId int64, deviceType protocol.DeviceType)](../internal/domain/gateway/session/service.go#L539)
  - [x] `countLocalOnlineUsers()` -> [CountLocalOnlineUsers()](../internal/domain/gateway/observation/service/statistics_service.go#L19)
  - [x] `onSessionEstablished(@NotNull UserSessionsManager userSessionsManager, @NotNull @ValidDeviceType DeviceType deviceType)` -> [OnSessionEstablished(ctx context.Context, userSessionsManager any, deviceType protocol.DeviceType)](../internal/domain/gateway/session/service.go#L557)
  - [x] `addOnSessionClosedListeners(Consumer<UserSession> onSessionClosed)` -> [AddOnSessionClosedListeners(ctx context.Context, onSessionClosed func(*UserSession)](../internal/domain/gateway/session/service.go#L562)
  - [x] `invokeGoOnlineHandlers(@NotNull UserSessionsManager userSessionsManager, @NotNull UserSession userSession)` -> [InvokeGoOnlineHandlers(ctx context.Context, userSessionsManager any, userSession *UserSession)](../internal/domain/gateway/session/service.go#L568)

- **UserService.java** ([java/im/turms/gateway/domain/session/service/UserService.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/domain/session/service/UserService.java))
> [简述功能]

  - [x] `authenticate(@NotNull Long userId, @Nullable String rawPassword)` -> [Authenticate(ctx context.Context, userID int64, rawPassword string)](../internal/domain/gateway/session/user_service.go#L25)
  - [x] `isActiveAndNotDeleted(@NotNull Long userId)` -> [IsActiveAndNotDeleted(ctx context.Context, userID int64)](../internal/domain/gateway/session/user_service.go#L38), [IsActiveAndNotDeleted(ctx context.Context, userID int64)](../internal/domain/user/repository/user_repository.go#L271)

- **UserSimultaneousLoginService.java** ([java/im/turms/gateway/domain/session/service/UserSimultaneousLoginService.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/domain/session/service/UserSimultaneousLoginService.java))
> [简述功能]

  - [x] `getConflictedDeviceTypes(@NotNull @ValidDeviceType DeviceType deviceType)` -> [GetConflictedDeviceTypes(deviceType protocol.DeviceType)](../internal/domain/gateway/session/user_simultaneous_login_service.go#L141)
  - [x] `isForbiddenDeviceType(DeviceType deviceType)` -> [IsForbiddenDeviceType(deviceType protocol.DeviceType)](../internal/domain/gateway/session/user_simultaneous_login_service.go#L163)
  - [x] `shouldDisconnectLoggingInDeviceIfConflicts()` -> [ShouldDisconnectLoggingInDeviceIfConflicts()](../internal/domain/gateway/session/user_simultaneous_login_service.go#L169)

- **ServiceAddressManager.java** ([java/im/turms/gateway/infra/address/ServiceAddressManager.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/address/ServiceAddressManager.java))
> [简述功能]

  - [x] `getWsAddress()` -> [GetWsAddress()](../internal/infra/address/service_address_manager.go#L9)
  - [x] `getTcpAddress()` -> [GetTcpAddress()](../internal/infra/address/service_address_manager.go#L15)
  - [x] `getUdpAddress()` -> [GetUdpAddress()](../internal/infra/address/service_address_manager.go#L21)

- **LdapClient.java** ([java/im/turms/gateway/infra/ldap/LdapClient.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/ldap/LdapClient.java))
> [简述功能]

  - [x] `isConnected()` -> [IsConnected()](../internal/domain/gateway/session/connection.go#L97)
  - [x] `connect()` -> [Connect()](../internal/infra/ldap/ldap_client.go#L27)
  - [x] `bind(boolean useFastBind, String dn, String password)` -> [Bind(useFastBind bool, dn string, password string)](../internal/infra/ldap/ldap_client.go#L45)
  - [x] `search(String baseDn, Scope scope, DerefAliases derefAliases, int sizeLimit, int timeLimit, boolean typeOnly, List<String> attributes, String filter)` -> [Search(ctx context.Context)](../internal/storage/elasticsearch/elasticsearch_client.go#L48)
  - [x] `modify(String dn, List<ModifyOperationChange> changes)` -> [Modify(req *ldap.ModifyRequest)](../internal/infra/ldap/ldap_client.go#L56)

- **BerBuffer.java** ([java/im/turms/gateway/infra/ldap/asn1/BerBuffer.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/ldap/asn1/BerBuffer.java))
> [简述功能]

  - [x] `skipTag()` -> [SkipTag()](../internal/infra/ldap/asn1/ber_buffer.go#L9)
  - [x] `skipTagAndLength()` -> [SkipTagAndLength()](../internal/infra/ldap/asn1/ber_buffer.go#L13)
  - [x] `skipTagAndLengthAndValue()` -> [SkipTagAndLengthAndValue()](../internal/infra/ldap/asn1/ber_buffer.go#L17)
  - [x] `readTag()` -> [ReadTag()](../internal/infra/ldap/asn1/ber_buffer.go#L21)
  - [x] `peekAndCheckTag(int tag)` -> [PeekAndCheckTag(tag int)](../internal/infra/ldap/asn1/ber_buffer.go#L26)
  - [x] `skipLength()` -> [SkipLength()](../internal/infra/ldap/asn1/ber_buffer.go#L30)
  - [x] `skipLengthAndValue()` -> [SkipLengthAndValue()](../internal/infra/ldap/asn1/ber_buffer.go#L34)
  - [x] `writeLength(int length)` -> [WriteLength(length int)](../internal/infra/ldap/asn1/ber_buffer.go#L38)
  - [x] `readLength()` -> [ReadLength()](../internal/infra/ldap/asn1/ber_buffer.go#L42)
  - [x] `tryReadLengthIfReadable()` -> [TryReadLengthIfReadable()](../internal/infra/ldap/asn1/ber_buffer.go#L47)
  - [x] `beginSequence()` -> [BeginSequence()](../internal/infra/ldap/asn1/ber_buffer.go#L52)
  - [x] `beginSequence(int tag)` -> [BeginSequenceWithTag(tag int)](../internal/infra/ldap/asn1/ber_buffer.go#L56)
  - [x] `endSequence()` -> [EndSequence()](../internal/infra/ldap/asn1/ber_buffer.go#L60)
  - [x] `writeBoolean(boolean value)` -> [WriteBoolean(value bool)](../internal/infra/ldap/asn1/ber_buffer.go#L64)
  - [x] `writeBoolean(int tag, boolean value)` -> [WriteBooleanWithTag(tag int, value bool)](../internal/infra/ldap/asn1/ber_buffer.go#L68)
  - [x] `readBoolean()` -> [ReadBoolean()](../internal/infra/ldap/asn1/ber_buffer.go#L72)
  - [x] `writeInteger(int value)` -> [WriteInteger(value int)](../internal/infra/ldap/asn1/ber_buffer.go#L77)
  - [x] `writeInteger(int tag, int value)` -> [WriteIntegerWithTag(tag int, value int)](../internal/infra/ldap/asn1/ber_buffer.go#L81)
  - [x] `readInteger()` -> [ReadInteger()](../internal/infra/ldap/asn1/ber_buffer.go#L85)
  - [x] `readIntWithTag(int tag)` -> [ReadIntWithTag(tag int)](../internal/infra/ldap/asn1/ber_buffer.go#L90)
  - [x] `writeOctetString(String value)` -> [WriteOctetString(value string)](../internal/infra/ldap/asn1/ber_buffer.go#L95)
  - [x] `writeOctetString(byte[] value)` -> [WriteOctetStringBytes(value []byte)](../internal/infra/ldap/asn1/ber_buffer.go#L99)
  - [x] `writeOctetString(int tag, byte[] value)` -> [WriteOctetStringBytesWithTag(tag int, value []byte)](../internal/infra/ldap/asn1/ber_buffer.go#L103)
  - [x] `writeOctetString(byte[] value, int start, int length)` -> [WriteOctetStringBytesRange(value []byte, start int, length int)](../internal/infra/ldap/asn1/ber_buffer.go#L107)
  - [x] `writeOctetString(int tag, byte[] value, int start, int length)` -> [WriteOctetStringBytesRangeWithTag(tag int, value []byte, start int, length int)](../internal/infra/ldap/asn1/ber_buffer.go#L111)
  - [x] `writeOctetString(int tag, String value)` -> [WriteOctetStringWithTag(tag int, value string)](../internal/infra/ldap/asn1/ber_buffer.go#L115)
  - [x] `writeOctetStrings(List<String> values)` -> [WriteOctetStrings(values []string)](../internal/infra/ldap/asn1/ber_buffer.go#L119)
  - [x] `readOctetString()` -> [ReadOctetString()](../internal/infra/ldap/asn1/ber_buffer.go#L123)
  - [x] `readOctetStringWithTag(int tag)` -> [ReadOctetStringWithTag(tag int)](../internal/infra/ldap/asn1/ber_buffer.go#L128)
  - [x] `readOctetStringWithLength(int length)` -> [ReadOctetStringWithLength(length int)](../internal/infra/ldap/asn1/ber_buffer.go#L133)
  - [x] `writeEnumeration(int value)` -> [WriteEnumeration(value int)](../internal/infra/ldap/asn1/ber_buffer.go#L138)
  - [x] `readEnumeration()` -> [ReadEnumeration()](../internal/infra/ldap/asn1/ber_buffer.go#L142)
  - [x] `getBytes()` -> [GetBytes()](../internal/infra/ldap/asn1/ber_buffer.go#L147)
  - [x] `skipBytes(int length)` -> [SkipBytes(length int)](../internal/infra/ldap/asn1/ber_buffer.go#L152)
  - [x] `close()` -> [Close()](../internal/domain/common/cache/ttl_cache.go#L72), [Close()](../internal/domain/gateway/access/client/common/net_connection.go#L68), [Close()](../internal/domain/gateway/access/client/tcp/tcp_server.go#L82), [Close()](../internal/domain/gateway/access/client/ws/websocket_connection.go#L87)
  - [x] `refCnt()` -> [RefCnt()](../internal/infra/ldap/asn1/ber_buffer.go#L156)
  - [x] `retain()` -> [Retain()](../internal/infra/ldap/asn1/ber_buffer.go#L161)
  - [x] `retain(int increment)` -> [RetainIncrement(increment int)](../internal/infra/ldap/asn1/ber_buffer.go#L165)
  - [x] `touch()` -> [Touch()](../internal/infra/ldap/asn1/ber_buffer.go#L169)
  - [x] `touch(Object hint)` -> [TouchWithHint(hint any)](../internal/infra/ldap/asn1/ber_buffer.go#L173)
  - [x] `release()` -> [Release()](../internal/infra/ldap/asn1/ber_buffer.go#L177)
  - [x] `release(int decrement)` -> [ReleaseDecrement(decrement int)](../internal/infra/ldap/asn1/ber_buffer.go#L182)
  - [x] `isReadable(int length)` -> [IsReadableLen(length int)](../internal/infra/ldap/asn1/ber_buffer.go#L187)
  - [x] `isReadable()` -> [IsReadable()](../internal/infra/ldap/asn1/ber_buffer.go#L192)
  - [x] `isReadableWithEnd(int end)` -> [IsReadableWithEnd(end int)](../internal/infra/ldap/asn1/ber_buffer.go#L197)
  - [x] `readerIndex()` -> [ReaderIndex()](../internal/infra/ldap/asn1/ber_buffer.go#L202)

- **Attribute.java** ([java/im/turms/gateway/infra/ldap/element/common/Attribute.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/ldap/element/common/Attribute.java))
> [简述功能]

  - [x] `isEmpty()` -> [IsEmpty()](../internal/domain/gateway/session/sharded_map.go#L44)
  - [x] `decode(BerBuffer buffer)` -> [Decode(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L11), [Decode(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L44), [Decode(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L67), [Decode(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L90), [Decode(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L122)

- **LdapMessage.java** ([java/im/turms/gateway/infra/ldap/element/common/LdapMessage.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/ldap/element/common/LdapMessage.java))
> [简述功能]

  - [x] `estimateSize()` -> [EstimateSize()](../internal/infra/ldap/element/elements.go#L20), [EstimateSize()](../internal/infra/ldap/element/elements.go#L53), [EstimateSize()](../internal/infra/ldap/element/elements.go#L76), [EstimateSize()](../internal/infra/ldap/element/elements.go#L108)
  - [x] `writeTo(BerBuffer buffer)` -> [WriteTo(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L25), [WriteTo(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L58), [WriteTo(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L81), [WriteTo(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L113)

- **LdapResult.java** ([java/im/turms/gateway/infra/ldap/element/common/LdapResult.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/ldap/element/common/LdapResult.java))
> [简述功能]

  - [x] `isSuccess()` -> [IsSuccess()](../internal/infra/ldap/element/elements.go#L34)

- **Control.java** ([java/im/turms/gateway/infra/ldap/element/common/control/Control.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/ldap/element/common/control/Control.java))
> [简述功能]

  - [x] `decode(BerBuffer buffer)` -> [Decode(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L11), [Decode(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L44), [Decode(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L67), [Decode(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L90), [Decode(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L122)

- **BindRequest.java** ([java/im/turms/gateway/infra/ldap/element/operation/bind/BindRequest.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/ldap/element/operation/bind/BindRequest.java))
> [简述功能]

  - [x] `estimateSize()` -> [EstimateSize()](../internal/infra/ldap/element/elements.go#L20), [EstimateSize()](../internal/infra/ldap/element/elements.go#L53), [EstimateSize()](../internal/infra/ldap/element/elements.go#L76), [EstimateSize()](../internal/infra/ldap/element/elements.go#L108)
  - [x] `writeTo(BerBuffer buffer)` -> [WriteTo(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L25), [WriteTo(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L58), [WriteTo(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L81), [WriteTo(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L113)

- **BindResponse.java** ([java/im/turms/gateway/infra/ldap/element/operation/bind/BindResponse.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/ldap/element/operation/bind/BindResponse.java))
> [简述功能]

  - [x] `decode(BerBuffer buffer)` -> [Decode(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L11), [Decode(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L44), [Decode(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L67), [Decode(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L90), [Decode(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L122)

- **ModifyRequest.java** ([java/im/turms/gateway/infra/ldap/element/operation/modify/ModifyRequest.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/ldap/element/operation/modify/ModifyRequest.java))
> [简述功能]

  - [x] `estimateSize()` -> [EstimateSize()](../internal/infra/ldap/element/elements.go#L20), [EstimateSize()](../internal/infra/ldap/element/elements.go#L53), [EstimateSize()](../internal/infra/ldap/element/elements.go#L76), [EstimateSize()](../internal/infra/ldap/element/elements.go#L108)
  - [x] `writeTo(BerBuffer buffer)` -> [WriteTo(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L25), [WriteTo(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L58), [WriteTo(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L81), [WriteTo(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L113)

- **ModifyResponse.java** ([java/im/turms/gateway/infra/ldap/element/operation/modify/ModifyResponse.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/ldap/element/operation/modify/ModifyResponse.java))
> [简述功能]

  - [x] `decode(BerBuffer buffer)` -> [Decode(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L11), [Decode(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L44), [Decode(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L67), [Decode(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L90), [Decode(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L122)

- **Filter.java** ([java/im/turms/gateway/infra/ldap/element/operation/search/Filter.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/ldap/element/operation/search/Filter.java))
> [简述功能]

  - [x] `write(BerBuffer buffer, String filter)` -> [Write(buffer *asn1.BerBuffer, filter string)](../internal/infra/ldap/element/elements.go#L99)

- **SearchRequest.java** ([java/im/turms/gateway/infra/ldap/element/operation/search/SearchRequest.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/ldap/element/operation/search/SearchRequest.java))
> [简述功能]

  - [x] `estimateSize()` -> [EstimateSize()](../internal/infra/ldap/element/elements.go#L20), [EstimateSize()](../internal/infra/ldap/element/elements.go#L53), [EstimateSize()](../internal/infra/ldap/element/elements.go#L76), [EstimateSize()](../internal/infra/ldap/element/elements.go#L108)
  - [x] `writeTo(BerBuffer buffer)` -> [WriteTo(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L25), [WriteTo(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L58), [WriteTo(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L81), [WriteTo(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L113)

- **SearchResult.java** ([java/im/turms/gateway/infra/ldap/element/operation/search/SearchResult.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/ldap/element/operation/search/SearchResult.java))
> [简述功能]

  - [x] `decode(BerBuffer buffer)` -> [Decode(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L11), [Decode(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L44), [Decode(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L67), [Decode(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L90), [Decode(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go#L122)
  - [x] `isComplete()` -> [IsComplete()](../internal/infra/ldap/element/elements.go#L126)

- **ApiLoggingContext.java** ([java/im/turms/gateway/infra/logging/ApiLoggingContext.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/logging/ApiLoggingContext.java))
> [简述功能]

  - [x] `shouldLogHeartbeatRequest()` -> [ShouldLogHeartbeatRequest()](../internal/infra/logging/api_logging_context.go#L24)
  - [x] `shouldLogRequest(TurmsRequest.KindCase requestType)` -> [ShouldLogRequest(requestType int)](../internal/infra/logging/api_logging_context.go#L12)
  - [x] `shouldLogNotification(TurmsRequest.KindCase requestType)` -> [ShouldLogNotification(requestType int)](../internal/infra/logging/api_logging_context.go#L18)

- **ClientApiLogging.java** ([java/im/turms/gateway/infra/logging/ClientApiLogging.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/logging/ClientApiLogging.java))
> [简述功能]

  - [x] `log(@Nullable Integer sessionId, @Nullable Long userId, @Nullable DeviceType deviceType, @Nullable Integer version, String ip, long requestId, TurmsRequest.KindCase requestType, int requestSize, long requestTime, TurmsNotification response, long processingTime)` -> [Log(request interface{}, serviceRequest interface{}, requestSize int64, requestTime int64, response interface{}, processingTime int64)](../internal/infra/logging/client_api_logging.go#L12)
  - [x] `log(@Nullable Integer sessionId, @Nullable Long userId, @Nullable DeviceType deviceType, @Nullable Integer version, String ip, long requestId, TurmsRequest.KindCase requestType, int requestSize, long requestTime, int responseCode, long processingTime)` -> [Log(request interface{}, serviceRequest interface{}, requestSize int64, requestTime int64, response interface{}, processingTime int64)](../internal/infra/logging/client_api_logging.go#L12)
  - [x] `log(@Nullable Integer sessionId, @Nullable Long userId, @Nullable DeviceType deviceType, @Nullable Integer version, String ip, long requestId, String requestType, int requestSize, long requestTime, int responseCode, @Nullable String responseDataType, int responseSize, long processingTime)` -> [Log(request interface{}, serviceRequest interface{}, requestSize int64, requestTime int64, response interface{}, processingTime int64)](../internal/infra/logging/client_api_logging.go#L12)

- **NotificationLoggingManager.java** ([java/im/turms/gateway/infra/logging/NotificationLoggingManager.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/logging/NotificationLoggingManager.java))
> [简述功能]

  - [x] `log(SimpleTurmsNotification notification, int notificationBytes, int recipientCount, int onlineRecipientCount)` -> [Log(request interface{}, serviceRequest interface{}, requestSize int64, requestTime int64, response interface{}, processingTime int64)](../internal/infra/logging/client_api_logging.go#L12)

- **SimpleTurmsNotification.java** ([java/im/turms/gateway/infra/proto/SimpleTurmsNotification.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/proto/SimpleTurmsNotification.java))
> [简述功能]

  - [ ] `SimpleTurmsNotification(long requesterId, Integer closeStatus, TurmsRequest.KindCase relayedRequestType)`

- **SimpleTurmsRequest.java** ([java/im/turms/gateway/infra/proto/SimpleTurmsRequest.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/proto/SimpleTurmsRequest.java))
> [简述功能]

  - [ ] `SimpleTurmsRequest(long requestId, TurmsRequest.KindCase type, CreateSessionRequest createSessionRequest)`
  - [x] `toString()` -> [ToString()](../internal/domain/gateway/session/connection.go#L107), [ToString()](../internal/infra/proto/proto_parser.go#L38)

- **TurmsNotificationParser.java** ([java/im/turms/gateway/infra/proto/TurmsNotificationParser.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/proto/TurmsNotificationParser.java))
> [简述功能]

  - [x] `parseSimpleNotification(CodedInputStream turmsRequestInputStream)` -> [ParseSimpleNotification(turmsRequestInputStream []byte)](../internal/infra/proto/proto_parser.go#L47)

- **TurmsRequestParser.java** ([java/im/turms/gateway/infra/proto/TurmsRequestParser.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/infra/proto/TurmsRequestParser.java))
> [简述功能]

  - [x] `parseSimpleRequest(CodedInputStream turmsRequestInputStream)` -> [ParseSimpleRequest(turmsRequestInputStream []byte)](../internal/infra/proto/proto_parser.go#L57)

- **MongoConfig.java** ([java/im/turms/gateway/storage/mongo/MongoConfig.java](../turms-orig/turms-gateway/src/main/java/im/turms/gateway/storage/mongo/MongoConfig.java))
> [简述功能]

  - [x] `adminMongoClient(TurmsPropertiesManager propertiesManager)` -> [AdminMongoClient()](../internal/storage/mongo/mongo_config.go#L7)
  - [x] `userMongoClient(TurmsPropertiesManager propertiesManager)` -> [UserMongoClient()](../internal/storage/mongo/mongo_config.go#L12)
  - [x] `mongoDataGenerator()` -> [NewMongoDataGenerator()](../internal/infra/mongo/mongo_data_generator.go#L9)

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

  - [x] `main(String[] args)` -> [main()](../cmd/turms-gateway/main.go#L8)

- **ServiceRequestDispatcher.java** ([java/im/turms/service/access/servicerequest/dispatcher/ServiceRequestDispatcher.java](../turms-orig/turms-service/src/main/java/im/turms/service/access/servicerequest/dispatcher/ServiceRequestDispatcher.java))
> [简述功能]

  - [x] `dispatch(TracingContext context, ServiceRequest serviceRequest)` -> [Dispatch(ctx context.Context, frame *codec.RpcFrame)](../internal/domain/common/infra/cluster/rpc/router.go#L43)

- **ClientRequest.java** ([java/im/turms/service/access/servicerequest/dto/ClientRequest.java](../turms-orig/turms-service/src/main/java/im/turms/service/access/servicerequest/dto/ClientRequest.java))
> [简述功能]

  - [x] `toString()` -> [ToString()](../internal/domain/gateway/session/connection.go#L107), [ToString()](../internal/infra/proto/proto_parser.go#L38)
  - [x] `turmsRequest()` -> [TurmsRequest()](../internal/domain/common/dto/client_request.go#L16)
  - [x] `userId()` -> [UserId()](../internal/domain/common/dto/client_request.go#L21)
  - [x] `deviceType()` -> [DeviceType()](../internal/domain/common/dto/client_request.go#L26)
  - [x] `clientIp()` -> [ClientIp()](../internal/domain/common/dto/client_request.go#L31)
  - [x] `requestId()` -> [RequestId()](../internal/domain/common/dto/client_request.go#L36)
  - [x] `equals(Object obj)` -> [Equals(obj interface{})](../internal/domain/common/dto/client_request.go#L41)
  - [x] `hashCode()` -> [HashCode()](../internal/domain/common/dto/client_request.go#L46)

- **RequestHandlerResult.java** ([java/im/turms/service/access/servicerequest/dto/RequestHandlerResult.java](../turms-orig/turms-service/src/main/java/im/turms/service/access/servicerequest/dto/RequestHandlerResult.java))
> [简述功能]

  - [x] `RequestHandlerResult(ResponseStatusCode code, @Nullable String reason, @Nullable TurmsNotification.Data response, List<Notification> notifications)` -> [NewRequestHandlerResult(code constant.ResponseStatusCode, reason *string, response *protocol.TurmsNotification_Data, notifications []*Notification)](../internal/domain/common/dto/request_handler_result.go#L35)
  - [x] `toString()` -> [ToString()](../internal/domain/gateway/session/connection.go#L107), [ToString()](../internal/infra/proto/proto_parser.go#L38)
  - [x] `of(@NotNull ResponseStatusCode code)` -> [RequestHandlerResultOfCode(code constant.ResponseStatusCode)](../internal/domain/common/dto/request_handler_result.go#L47)
  - [x] `of(@NotNull ResponseStatusCode code, @Nullable String reason)` -> [RequestHandlerResultOfCodeReason(code constant.ResponseStatusCode, reason *string)](../internal/domain/common/dto/request_handler_result.go#L52)
  - [x] `of(@NotNull TurmsNotification.Data response)` -> [RequestHandlerResultOfResponse(response *protocol.TurmsNotification_Data)](../internal/domain/common/dto/request_handler_result.go#L60)
  - [x] `of(boolean forwardNotificationsToRequesterOtherOnlineSessions, @NotNull TurmsRequest notification)` -> [RequestHandlerResultOfForwardNotification(forward bool, notification *protocol.TurmsRequest)](../internal/domain/common/dto/request_handler_result.go#L65)
  - [x] `of(boolean forwardNotificationsToRequesterOtherOnlineSessions, @NotNull Long recipientId, @NotNull TurmsRequest notification)` -> [RequestHandlerResultOfForwardRecipientNotification(forward bool, recipientId int64, notification *protocol.TurmsRequest)](../internal/domain/common/dto/request_handler_result.go#L75)
  - [x] `of(@NotNull Long recipientId, @NotNull TurmsRequest notification)` -> [RequestHandlerResultOfRecipientNotification(recipientId int64, notification *protocol.TurmsRequest)](../internal/domain/common/dto/request_handler_result.go#L85)
  - [x] `of(@NotEmpty Set<Long> recipientIds, @NotNull TurmsRequest dataForRecipient)` -> [RequestHandlerResultOfRecipientsNotification(recipientIds []int64, dataForRecipient *protocol.TurmsRequest)](../internal/domain/common/dto/request_handler_result.go#L95)
  - [x] `of(boolean forwardNotificationToRequesterOtherOnlineSessions, @NotEmpty Set<Long> recipientIds, @NotNull TurmsRequest notification)` -> [RequestHandlerResultOfForwardRecipientsNotification(forward bool, recipientIds []int64, notification *protocol.TurmsRequest)](../internal/domain/common/dto/request_handler_result.go#L105)
  - [x] `of(TurmsNotification.Data response, @NotEmpty Set<Long> recipientIds, @NotNull TurmsRequest notificationForRecipients, @NotNull TurmsRequest notificationForRequesterOtherOnlineSessions)` -> [RequestHandlerResultOfResponseRecipientsNotifications(response *protocol.TurmsNotification_Data, recipientIds []int64, notificationForRecipients, notificationForRequesterOtherOnlineSessions *protocol.TurmsRequest)](../internal/domain/common/dto/request_handler_result.go#L115)
  - [x] `of(@NotEmpty Set<Long> recipientIds, @NotNull TurmsRequest notificationForRecipients, @NotNull TurmsRequest notificationForRequesterOtherOnlineSessions)` -> [RequestHandlerResultOfRecipientsNotifications(recipientIds []int64, notificationForRecipients, notificationForRequesterOtherOnlineSessions *protocol.TurmsRequest)](../internal/domain/common/dto/request_handler_result.go#L128)
  - [x] `of(TurmsNotification.Data response, boolean forwardNotificationToRequesterOtherOnlineSessions, @NotEmpty Set<Long> recipientIds, @NotNull TurmsRequest notification)` -> [RequestHandlerResultOfResponseForwardRecipientsNotification(response *protocol.TurmsNotification_Data, forward bool, recipientIds []int64, notification *protocol.TurmsRequest)](../internal/domain/common/dto/request_handler_result.go#L141)
  - [x] `of(@NotNull ResponseStatusCode code, @NotNull Long recipientId, @NotNull TurmsRequest notification)` -> [RequestHandlerResultOfCodeRecipientNotification(code constant.ResponseStatusCode, recipientId int64, notification *protocol.TurmsRequest)](../internal/domain/common/dto/request_handler_result.go#L151)
  - [x] `of(@NotNull ResponseStatusCode code, boolean forwardNotificationToRequesterOtherOnlineSessions, @NotNull Long recipientId, @NotNull TurmsRequest notification)` -> [RequestHandlerResultOfCodeForwardRecipientNotification(code constant.ResponseStatusCode, forward bool, recipientId int64, notification *protocol.TurmsRequest)](../internal/domain/common/dto/request_handler_result.go#L161)
  - [x] `of(@NotNull List<Notification> notifications)` -> [RequestHandlerResultOfNotifications(notifications []*Notification)](../internal/domain/common/dto/request_handler_result.go#L171)
  - [x] `of(@NotNull Notification notification)` -> [RequestHandlerResultOfNotification(notification *Notification)](../internal/domain/common/dto/request_handler_result.go#L176)
  - [x] `ofDataLong(@NotNull Long value)` -> [RequestHandlerResultOfDataLong(value int64)](../internal/domain/common/dto/request_handler_result.go#L181)
  - [x] `ofDataLong(@NotNull Long value, @NotNull Long recipientId, @NotNull TurmsRequest notification)` -> [RequestHandlerResultOfDataLongRecipientNotification(value int64, recipientId int64, notification *protocol.TurmsRequest)](../internal/domain/common/dto/request_handler_result.go#L193)
  - [x] `ofDataLong(@NotNull Long value, boolean forwardNotificationToRequesterOtherOnlineSessions, @NotNull Long recipientId, @NotNull TurmsRequest notification)` -> [RequestHandlerResultOfDataLongForwardRecipientNotification(value int64, forward bool, recipientId int64, notification *protocol.TurmsRequest)](../internal/domain/common/dto/request_handler_result.go#L205)
  - [x] `ofDataLong(@NotNull Long value, boolean forwardDataForRecipientsToRequesterOtherOnlineSessions, @NotNull TurmsRequest notification)` -> [RequestHandlerResultOfDataLongForwardNotification(value int64, forward bool, notification *protocol.TurmsRequest)](../internal/domain/common/dto/request_handler_result.go#L217)
  - [x] `ofDataLong(@NotNull Long value, boolean forwardNotificationsToRequesterOtherOnlineSessions, @NotEmpty Set<Long> recipients, TurmsRequest notification)` -> [RequestHandlerResultOfDataLongForwardRecipientsNotification(value int64, forward bool, recipients []int64, notification *protocol.TurmsRequest)](../internal/domain/common/dto/request_handler_result.go#L233)
  - [x] `ofDataLongs(@NotNull Collection<Long> values)` -> [RequestHandlerResultOfDataLongs(values []int64)](../internal/domain/common/dto/request_handler_result.go#L245)
  - [x] `Notification(boolean forwardToRequesterOtherOnlineSessions, Set<Long> recipients, TurmsRequest notification)` -> [NewNotification(forwardToRequesterOtherOnlineSessions bool, recipients []int64, notification *protocol.TurmsRequest)](../internal/domain/common/dto/request_handler_result.go#L17)
  - [x] `of(boolean forwardToRequesterOtherOnlineSessions, Set<Long> recipients, TurmsRequest notification)` -> [NotificationOfForwardRecipientsNotification(forward bool, recipients []int64, notification *protocol.TurmsRequest)](../internal/domain/common/dto/request_handler_result.go#L263)
  - [x] `of(boolean forwardToRequesterOtherOnlineSessions, Long recipient, TurmsRequest notification)` -> [NotificationOfForwardRecipientNotification(forward bool, recipient int64, notification *protocol.TurmsRequest)](../internal/domain/common/dto/request_handler_result.go#L268)
  - [x] `of(boolean forwardToRequesterOtherOnlineSessions, TurmsRequest notification)` -> [NotificationOfForwardNotification(forward bool, notification *protocol.TurmsRequest)](../internal/domain/common/dto/request_handler_result.go#L273)
  - [x] `toString()` -> [ToString()](../internal/domain/gateway/session/connection.go#L107), [ToString()](../internal/infra/proto/proto_parser.go#L38)

- **AdminController.java** ([java/im/turms/service/domain/admin/access/admin/controller/AdminController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/admin/access/admin/controller/AdminController.java))
> [简述功能]

  - [ ] `checkLoginNameAndPassword()`
  - [x] `addAdmin(RequestContext requestContext, @RequestBody AddAdminDTO addAdminDTO)` -> [AddAdmin(ctx context.Context, id *int64, loginName string, rawPassword string, displayName string, roleIds []int64)](../internal/domain/admin/service/admin_services.go#L182)
  - [x] `queryAdmins(@QueryParam(required = false)` -> [QueryAdmins(ctx context.Context, ids []int64, loginNames []string, roleIds []int64, page *int, size *int)](../internal/domain/admin/service/admin_services.go#L208)
  - [x] `queryAdmins(@QueryParam(required = false)` -> [QueryAdmins(ctx context.Context, ids []int64, loginNames []string, roleIds []int64, page *int, size *int)](../internal/domain/admin/service/admin_services.go#L208)
  - [x] `updateAdmins(RequestContext requestContext, Set<Long> ids, @RequestBody UpdateAdminDTO updateAdminDTO)` -> [UpdateAdmins(ctx context.Context, ids []int64, password []byte, displayName *string, roleIDs []int64)](../internal/domain/admin/repository/admin_repository.go#L37), [UpdateAdmins(ctx context.Context, targetAdminIds []int64, rawPassword *string, displayName *string, roleIds []int64)](../internal/domain/admin/service/admin_services.go#L220)
  - [x] `deleteAdmins(RequestContext requestContext, Set<Long> ids)` -> [DeleteAdmins(ctx context.Context, ids []int64)](../internal/domain/admin/repository/admin_repository.go#L123)

- **AdminPermissionController.java** ([java/im/turms/service/domain/admin/access/admin/controller/AdminPermissionController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/admin/access/admin/controller/AdminPermissionController.java))
> [简述功能]

  - [ ] `queryAdminPermissions()`

- **AdminRoleController.java** ([java/im/turms/service/domain/admin/access/admin/controller/AdminRoleController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/admin/access/admin/controller/AdminRoleController.java))
> [简述功能]

  - [ ] `addAdminRole(RequestContext requestContext, @RequestBody AddAdminRoleDTO addAdminRoleDTO)`
  - [x] `queryAdminRoles(@QueryParam(required = false)` -> [QueryAdminRoles(ctx context.Context, ids []int64, names []string, includedPermissions []permission.AdminPermission, ranks []int, page *int, size *int)](../internal/domain/admin/service/admin_services.go#L101)
  - [x] `queryAdminRoles(@QueryParam(required = false)` -> [QueryAdminRoles(ctx context.Context, ids []int64, names []string, includedPermissions []permission.AdminPermission, ranks []int, page *int, size *int)](../internal/domain/admin/service/admin_services.go#L101)
  - [ ] `updateAdminRole(RequestContext requestContext, Set<Long> ids, @RequestBody UpdateAdminRoleDTO updateAdminRoleDTO)`
  - [x] `deleteAdminRoles(RequestContext requestContext, Set<Long> ids)` -> [DeleteAdminRoles(ctx context.Context, ids []int64)](../internal/domain/admin/repository/admin_role_repository.go#L154)

- **AddAdminDTO.java** ([java/im/turms/service/domain/admin/access/admin/dto/request/AddAdminDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/admin/access/admin/dto/request/AddAdminDTO.java))
> [简述功能]

  - [ ] `AddAdminDTO(String loginName, @SensitiveProperty(SensitiveProperty.Access.ALLOW_DESERIALIZATION)`
  - [x] `toString()` -> [ToString()](../internal/domain/gateway/session/connection.go#L107), [ToString()](../internal/infra/proto/proto_parser.go#L38)

- **AddAdminRoleDTO.java** ([java/im/turms/service/domain/admin/access/admin/dto/request/AddAdminRoleDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/admin/access/admin/dto/request/AddAdminRoleDTO.java))
> [简述功能]

  - [ ] `AddAdminRoleDTO(Long id, String name, Set<String> permissions, Integer rank)`

- **UpdateAdminDTO.java** ([java/im/turms/service/domain/admin/access/admin/dto/request/UpdateAdminDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/admin/access/admin/dto/request/UpdateAdminDTO.java))
> [简述功能]

  - [ ] `UpdateAdminDTO(@SensitiveProperty(SensitiveProperty.Access.ALLOW_DESERIALIZATION)`
  - [x] `toString()` -> [ToString()](../internal/domain/gateway/session/connection.go#L107), [ToString()](../internal/infra/proto/proto_parser.go#L38)

- **UpdateAdminRoleDTO.java** ([java/im/turms/service/domain/admin/access/admin/dto/request/UpdateAdminRoleDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/admin/access/admin/dto/request/UpdateAdminRoleDTO.java))
> [简述功能]

  - [ ] `UpdateAdminRoleDTO(String name, Set<String> permissions, Integer rank)`

- **PermissionDTO.java** ([java/im/turms/service/domain/admin/access/admin/dto/response/PermissionDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/admin/access/admin/dto/response/PermissionDTO.java))
> [简述功能]

  - [ ] `PermissionDTO(String group, AdminPermission permission)`

- **AdminRepository.java** ([java/im/turms/service/domain/admin/repository/AdminRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/admin/repository/AdminRepository.java))
> [简述功能]

  - [x] `updateAdmins(Set<Long> ids, @Nullable byte[] password, @Nullable String displayName, @Nullable Set<Long> roleIds)` -> [UpdateAdmins(ctx context.Context, ids []int64, password []byte, displayName *string, roleIDs []int64)](../internal/domain/admin/repository/admin_repository.go#L37), [UpdateAdmins(ctx context.Context, targetAdminIds []int64, rawPassword *string, displayName *string, roleIds []int64)](../internal/domain/admin/service/admin_services.go#L220)
  - [x] `countAdmins(@Nullable Set<Long> ids, @Nullable Set<Long> roleIds)` -> [CountAdmins(ctx context.Context, ids []int64, roleIDs []int64)](../internal/domain/admin/repository/admin_repository.go#L79), [CountAdmins(ctx context.Context, ids []int64, roleIds []int64)](../internal/domain/admin/service/admin_services.go#L232)
  - [x] `findAdmins(@Nullable Collection<Long> ids, @Nullable Collection<String> loginNames, @Nullable Collection<Long> roleIds, @Nullable Integer page, @Nullable Integer size)` -> [FindAdmins(ctx context.Context, ids []int64, loginNames []string, roleIDs []int64, page *int, size *int)](../internal/domain/admin/repository/admin_repository.go#L90)

- **AdminRoleRepository.java** ([java/im/turms/service/domain/admin/repository/AdminRoleRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/admin/repository/AdminRoleRepository.java))
> [简述功能]

  - [x] `updateAdminRoles(Set<Long> roleIds, String newName, @Nullable Set<AdminPermission> permissions, @Nullable Integer rank)` -> [UpdateAdminRoles(ctx context.Context, roleIds []int64, newName *string, permissions []permission.AdminPermission, rank *int)](../internal/domain/admin/repository/admin_role_repository.go#L40)
  - [x] `countAdminRoles(@Nullable Set<Long> ids, @Nullable Set<String> names, @Nullable Set<AdminPermission> includedPermissions, @Nullable Set<Integer> ranks)` -> [CountAdminRoles(ctx context.Context, ids []int64, names []string, includedPermissions []permission.AdminPermission, ranks []int)](../internal/domain/admin/repository/admin_role_repository.go#L82), [CountAdminRoles(ctx context.Context, ids []int64, names []string, includedPermissions []permission.AdminPermission, ranks []int)](../internal/domain/admin/service/admin_services.go#L109)
  - [x] `findAdminRoles(@Nullable Set<Long> roleIds, @Nullable Set<String> names, @Nullable Set<AdminPermission> includedPermissions, @Nullable Set<Integer> ranks, @Nullable Integer page, @Nullable Integer size)` -> [FindAdminRoles(ctx context.Context, roleIds []int64, names []string, includedPermissions []permission.AdminPermission, ranks []int, page *int, size *int)](../internal/domain/admin/repository/admin_role_repository.go#L87)
  - [x] `findAdminRolesByIdsAndRankGreaterThan(@NotNull Collection<Long> roleIds, @Nullable Integer rankGreaterThan)` -> [FindAdminRolesByIdsAndRankGreaterThan(ctx context.Context, roleIds []int64, rankGreaterThan *int)](../internal/domain/admin/repository/admin_role_repository.go#L111)
  - [x] `findHighestRankByRoleIds(Set<Long> roleIds)` -> [FindHighestRankByRoleIds(ctx context.Context, roleIds []int64)](../internal/domain/admin/repository/admin_role_repository.go#L133)

- **AdminRoleService.java** ([java/im/turms/service/domain/admin/service/AdminRoleService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/admin/service/AdminRoleService.java))
> [简述功能]

  - [ ] `authAndAddAdminRole(@NotNull Long requesterId, @NotNull Long roleId, @NotNull @NoWhitespace @Size( min = MIN_ROLE_NAME_LIMIT, max = MAX_ROLE_NAME_LIMIT)`
  - [ ] `addAdminRole(@NotNull Long roleId, @NotNull @NoWhitespace @Size( min = MIN_ROLE_NAME_LIMIT, max = MAX_ROLE_NAME_LIMIT)`
  - [ ] `authAndDeleteAdminRoles(@NotNull Long requesterId, @NotEmpty Set<Long> roleIds)`
  - [x] `deleteAdminRoles(@NotEmpty Set<Long> roleIds)` -> [DeleteAdminRoles(ctx context.Context, ids []int64)](../internal/domain/admin/repository/admin_role_repository.go#L154)
  - [ ] `authAndUpdateAdminRoles(@NotNull Long requesterId, @NotEmpty Set<Long> roleIds, @Nullable @NoWhitespace @Size( min = MIN_ROLE_NAME_LIMIT, max = MAX_ROLE_NAME_LIMIT)`
  - [ ] `updateAdminRole(@NotEmpty Set<Long> roleIds, @Nullable @NoWhitespace @Size( min = MIN_ROLE_NAME_LIMIT, max = MAX_ROLE_NAME_LIMIT)`
  - [x] `queryAdminRoles(@Nullable Set<Long> ids, @Nullable Set<String> names, @Nullable Set<AdminPermission> includedPermissions, @Nullable Set<Integer> ranks, @Nullable Integer page, @Nullable Integer size)` -> [QueryAdminRoles(ctx context.Context, ids []int64, names []string, includedPermissions []permission.AdminPermission, ranks []int, page *int, size *int)](../internal/domain/admin/service/admin_services.go#L101)
  - [x] `queryAndCacheRolesByRoleIdsAndRankGreaterThan(@NotNull Collection<Long> roleIds, @NotNull Integer rankGreaterThan)` -> [QueryAndCacheRolesByRoleIdsAndRankGreaterThan(ctx context.Context, roleIds []int64, rankGreaterThan int)](../internal/domain/admin/service/admin_services.go#L105)
  - [x] `countAdminRoles(@Nullable Set<Long> ids, @Nullable Set<String> names, @Nullable Set<AdminPermission> includedPermissions, @Nullable Set<Integer> ranks)` -> [CountAdminRoles(ctx context.Context, ids []int64, names []string, includedPermissions []permission.AdminPermission, ranks []int)](../internal/domain/admin/repository/admin_role_repository.go#L82), [CountAdminRoles(ctx context.Context, ids []int64, names []string, includedPermissions []permission.AdminPermission, ranks []int)](../internal/domain/admin/service/admin_services.go#L109)
  - [x] `queryHighestRankByAdminId(@NotNull Long adminId)` -> [QueryHighestRankByAdminId(ctx context.Context, adminId int64)](../internal/domain/admin/service/admin_services.go#L113)
  - [x] `queryHighestRankByRoleIds(@NotNull Set<Long> roleIds)` -> [QueryHighestRankByRoleIds(ctx context.Context, roleIds []int64)](../internal/domain/admin/service/admin_services.go#L118)
  - [x] `isAdminRankHigherThanRank(@NotNull Long adminId, @NotNull Integer rank)` -> [IsAdminRankHigherThanRank(ctx context.Context, adminId int64, rank int)](../internal/domain/admin/service/admin_services.go#L122)
  - [x] `queryPermissions(@NotNull Long adminId)` -> [QueryPermissions(ctx context.Context, adminId int64)](../internal/domain/admin/service/admin_services.go#L133)

- **AdminService.java** ([java/im/turms/service/domain/admin/service/AdminService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/admin/service/AdminService.java))
> [简述功能]

  - [x] `queryRoleIdsByAdminIds(@NotEmpty Set<Long> adminIds)` -> [QueryRoleIdsByAdminIds(ctx context.Context, adminIds []int64)](../internal/domain/admin/service/admin_services.go#L165)
  - [x] `authAndAddAdmin(@NotNull Long requesterId, @Nullable @NoWhitespace @Size( min = MIN_LOGIN_NAME_LIMIT, max = MAX_LOGIN_NAME_LIMIT)` -> [AuthAndAddAdmin(ctx context.Context, requesterId int64, loginName string, rawPassword string, displayName string, roleIds []int64)](../internal/domain/admin/service/admin_services.go#L177)
  - [x] `addAdmin(@Nullable Long id, @Nullable @NoWhitespace @Size( min = MIN_LOGIN_NAME_LIMIT, max = MAX_LOGIN_NAME_LIMIT)` -> [AddAdmin(ctx context.Context, id *int64, loginName string, rawPassword string, displayName string, roleIds []int64)](../internal/domain/admin/service/admin_services.go#L182)
  - [x] `queryAdmins(@Nullable Collection<Long> ids, @Nullable Collection<String> loginNames, @Nullable Collection<Long> roleIds, @Nullable Integer page, @Nullable Integer size)` -> [QueryAdmins(ctx context.Context, ids []int64, loginNames []string, roleIds []int64, page *int, size *int)](../internal/domain/admin/service/admin_services.go#L208)
  - [x] `authAndDeleteAdmins(@NotNull Long requesterId, @NotEmpty Set<Long> adminIds)` -> [AuthAndDeleteAdmins(ctx context.Context, requesterId int64, adminIds []int64)](../internal/domain/admin/service/admin_services.go#L212)
  - [x] `authAndUpdateAdmins(@NotNull Long requesterId, @NotEmpty Set<Long> targetAdminIds, @Nullable @NoWhitespace @Size( min = MIN_PASSWORD_LIMIT, max = MAX_PASSWORD_LIMIT)` -> [AuthAndUpdateAdmins(ctx context.Context, requesterId int64, targetAdminIds []int64, rawPassword *string, displayName *string, roleIds []int64)](../internal/domain/admin/service/admin_services.go#L216)
  - [x] `updateAdmins(@NotEmpty Set<Long> targetAdminIds, @Nullable @NoWhitespace @Size( min = MIN_PASSWORD_LIMIT, max = MAX_PASSWORD_LIMIT)` -> [UpdateAdmins(ctx context.Context, ids []int64, password []byte, displayName *string, roleIDs []int64)](../internal/domain/admin/repository/admin_repository.go#L37), [UpdateAdmins(ctx context.Context, targetAdminIds []int64, rawPassword *string, displayName *string, roleIds []int64)](../internal/domain/admin/service/admin_services.go#L220)
  - [x] `countAdmins(@Nullable Set<Long> ids, @Nullable Set<Long> roleIds)` -> [CountAdmins(ctx context.Context, ids []int64, roleIDs []int64)](../internal/domain/admin/repository/admin_repository.go#L79), [CountAdmins(ctx context.Context, ids []int64, roleIds []int64)](../internal/domain/admin/service/admin_services.go#L232)
  - [x] `errorRequesterNotExist()` -> [ErrorRequesterNotExist()](../internal/domain/admin/service/admin_services.go#L236)

- **IpBlocklistController.java** ([java/im/turms/service/domain/blocklist/access/admin/controller/IpBlocklistController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/blocklist/access/admin/controller/IpBlocklistController.java))
> [简述功能]

  - [x] `addBlockedIps(@RequestBody AddBlockedIpsDTO addBlockedIpsDTO)` -> [AddBlockedIps()](../internal/domain/blocklist/access/admin/controller/blocklist_controllers.go#L9)
  - [x] `queryBlockedIps(Set<String> ids)` -> [QueryBlockedIpsByIds()](../internal/domain/blocklist/access/admin/controller/blocklist_controllers.go#L13)
  - [x] `queryBlockedIps(int page, @QueryParam(required = false)` -> [QueryBlockedIpsByPage()](../internal/domain/blocklist/access/admin/controller/blocklist_controllers.go#L17)
  - [x] `deleteBlockedIps(@QueryParam(required = false)` -> [DeleteBlockedIps()](../internal/domain/blocklist/access/admin/controller/blocklist_controllers.go#L21)

- **UserBlocklistController.java** ([java/im/turms/service/domain/blocklist/access/admin/controller/UserBlocklistController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/blocklist/access/admin/controller/UserBlocklistController.java))
> [简述功能]

  - [x] `addBlockedUserIds(@RequestBody AddBlockedUserIdsDTO addBlockedUserIdsDTO)` -> [AddBlockedUserIds()](../internal/domain/blocklist/access/admin/controller/blocklist_controllers.go#L30)
  - [x] `queryBlockedUsers(Set<Long> ids)` -> [IsBlocked(ctx context.Context, groupID int64, userID int64)](../internal/domain/group/service/group_blocklist_service.go#L126)
  - [x] `queryBlockedUsers(int page, @QueryParam(required = false)` -> [IsBlocked(ctx context.Context, groupID int64, userID int64)](../internal/domain/group/service/group_blocklist_service.go#L126)
  - [x] `deleteBlockedUserIds(@QueryParam(required = false)` -> [DeleteBlockedUserIds()](../internal/domain/blocklist/access/admin/controller/blocklist_controllers.go#L34)

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

  - [x] `serialize(BlockedClient value, JsonGenerator gen, SerializerProvider provider)` -> [Serialize()](../internal/storage/elasticsearch/model/elasticsearch_model.go#L8)

- **MemberController.java** ([java/im/turms/service/domain/cluster/access/admin/controller/MemberController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/cluster/access/admin/controller/MemberController.java))
> [简述功能]

  - [x] `queryMembers()` -> [QueryMembers()](../internal/domain/cluster/access/admin/controller/cluster_controllers.go#L9)
  - [x] `removeMembers(List<String> ids)` -> [RemoveMembers()](../internal/domain/cluster/access/admin/controller/cluster_controllers.go#L13)
  - [x] `addMember(@RequestBody AddMemberDTO addMemberDTO)` -> [AddMember()](../internal/domain/cluster/access/admin/controller/cluster_controllers.go#L17)
  - [x] `updateMember(String id, @RequestBody UpdateMemberDTO updateMemberDTO)` -> [UpdateMember()](../internal/domain/cluster/access/admin/controller/cluster_controllers.go#L21)
  - [x] `queryLeader()` -> [QueryLeader()](../internal/domain/cluster/access/admin/controller/cluster_controllers.go#L25)
  - [x] `electNewLeader(@QueryParam(required = false)` -> [ElectNewLeader()](../internal/domain/cluster/access/admin/controller/cluster_controllers.go#L29)

- **SettingController.java** ([java/im/turms/service/domain/cluster/access/admin/controller/SettingController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/cluster/access/admin/controller/SettingController.java))
> [简述功能]

  - [x] `queryClusterSettings(boolean queryLocalSettings, boolean onlyMutable)` -> [QueryClusterSettings()](../internal/domain/cluster/access/admin/controller/cluster_controllers.go#L38)
  - [x] `updateClusterSettings(boolean reset, boolean updateLocalSettings, @RequestBody(required = false)` -> [UpdateClusterSettings()](../internal/domain/cluster/access/admin/controller/cluster_controllers.go#L42)
  - [x] `queryClusterConfigMetadata(boolean queryLocalSettings, boolean onlyMutable, boolean withValue)` -> [QueryClusterConfigMetadata()](../internal/domain/cluster/access/admin/controller/cluster_controllers.go#L46)

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

  - [x] `getPageSize(@Nullable Integer size)` -> [GetPageSize()](../internal/domain/common/access/admin/controller/base_controller.go#L9)
  - [x] `queryBetweenDate(DateRange dateRange, DivideBy divideBy, Function3<DateRange, Boolean, Boolean, Mono<Long>> function, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages)` -> [QueryBetweenDate()](../internal/domain/common/access/admin/controller/base_controller.go#L13)
  - [x] `queryBetweenDate(DateRange dateRange, DivideBy divideBy, Function<DateRange, Mono<Long>> function)` -> [QueryBetweenDateFunc()](../internal/domain/common/access/admin/controller/base_controller.go#L17)
  - [x] `checkAndQueryBetweenDate(DateRange dateRange, DivideBy divideBy, Function3<DateRange, Boolean, Boolean, Mono<Long>> function, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages)` -> [CheckAndQueryBetweenDate()](../internal/domain/common/access/admin/controller/base_controller.go#L21)
  - [x] `checkAndQueryBetweenDate(DateRange dateRange, DivideBy divideBy, Function<DateRange, Mono<Long>> function)` -> [CheckAndQueryBetweenDateFunc()](../internal/domain/common/access/admin/controller/base_controller.go#L25)

- **StatisticsRecordDTO.java** ([java/im/turms/service/domain/common/access/admin/dto/response/StatisticsRecordDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/common/access/admin/dto/response/StatisticsRecordDTO.java))
> [简述功能]

  - [ ] `StatisticsRecordDTO(Date date, Long total)`

- **ServicePermission.java** ([java/im/turms/service/domain/common/permission/ServicePermission.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/common/permission/ServicePermission.java))
> [简述功能]

  - [x] `ServicePermission(ResponseStatusCode code, String reason)` -> [NewServicePermission()](../internal/domain/common/permission/service_permission.go#L9)
  - [x] `get(ResponseStatusCode code)` -> [Get(key K)](../internal/domain/common/cache/sharded_map.go#L53), [GetBufferFromStatusCode(code constant.ResponseStatusCode)](../internal/domain/gateway/access/client/udp/udp_server.go#L211)
  - [x] `get(ResponseStatusCode code, String reason)` -> [Get(key K)](../internal/domain/common/cache/sharded_map.go#L53)

- **ExpirableEntityRepository.java** ([java/im/turms/service/domain/common/repository/ExpirableEntityRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/common/repository/ExpirableEntityRepository.java))
> [简述功能]

  - [x] `isExpired(long creationDate)` -> [IsExpired()](../internal/domain/common/repository/expirable_entity_repository.go#L9)
  - [x] `getEntityExpirationDate()` -> [GetEntityExpirationDate()](../internal/domain/common/repository/expirable_entity_repository.go#L13), [GetEntityExpirationDate()](../internal/domain/common/service/common_services.go#L9)
  - [x] `deleteExpiredData(String creationDateFieldName, Date expirationDate)` -> [DeleteExpiredData(ctx context.Context, expirationDate time.Time)](../internal/domain/user/repository/user_friend_request_repository.go#L215)
  - [x] `findMany(Filter filter)` -> [FindMany(ctx context.Context, filter bson.M)](../internal/domain/user/repository/user_repository.go#L72)
  - [x] `findMany(Filter filter, QueryOptions options)` -> [FindMany(ctx context.Context, filter bson.M)](../internal/domain/user/repository/user_repository.go#L72)

- **ExpirableEntityService.java** ([java/im/turms/service/domain/common/service/ExpirableEntityService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/common/service/ExpirableEntityService.java))
> [简述功能]

  - [x] `getEntityExpirationDate()` -> [GetEntityExpirationDate()](../internal/domain/common/repository/expirable_entity_repository.go#L13), [GetEntityExpirationDate()](../internal/domain/common/service/common_services.go#L9)

- **UserDefinedAttributesService.java** ([java/im/turms/service/domain/common/service/UserDefinedAttributesService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/common/service/UserDefinedAttributesService.java))
> [简述功能]

  - [x] `updateGlobalProperties(UserDefinedAttributesProperties properties)` -> [UpdateGlobalProperties()](../internal/domain/common/service/common_services.go#L18)
  - [x] `parseAttributesForUpsert(Map<String, Value> userDefinedAttributes)` -> [ParseAttributesForUpsert()](../internal/domain/common/service/common_services.go#L22)

- **ExpirableRequestInspector.java** ([java/im/turms/service/domain/common/util/ExpirableRequestInspector.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/common/util/ExpirableRequestInspector.java))
> [简述功能]

  - [x] `isProcessedByResponder(@Nullable RequestStatus status)` -> [IsProcessedByResponder()](../internal/domain/common/util/expirable_request_inspector.go#L9)

- **DataValidator.java** ([java/im/turms/service/domain/common/validation/DataValidator.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/common/validation/DataValidator.java))
> [简述功能]

  - [x] `validRequestStatus(RequestStatus status)` -> [ValidRequestStatus(status interface{}, name string)](../internal/infra/validator/validator.go#L82)
  - [x] `validResponseAction(ResponseAction action)` -> [ValidResponseAction()](../internal/infra/validator/validator.go#L87)
  - [x] `validDeviceType(DeviceType deviceType)` -> [ValidDeviceType()](../internal/infra/validator/validator.go#L91)
  - [x] `validProfileAccess(ProfileAccessStrategy value)` -> [ValidProfileAccess()](../internal/infra/validator/validator.go#L95)
  - [x] `validRelationshipKey(UserRelationship.Key key)` -> [ValidRelationshipKey()](../internal/infra/validator/validator.go#L99)
  - [x] `validRelationshipGroupKey(UserRelationshipGroup.Key key)` -> [ValidRelationshipGroupKey()](../internal/infra/validator/validator.go#L103)
  - [x] `validGroupMemberKey(GroupMember.Key key)` -> [ValidGroupMemberKey()](../internal/infra/validator/validator.go#L107)
  - [x] `validGroupMemberRole(GroupMemberRole role)` -> [ValidGroupMemberRole()](../internal/infra/validator/validator.go#L111)
  - [x] `validGroupBlockedUserKey(GroupBlockedUser.Key key)` -> [ValidGroupBlockedUserKey()](../internal/infra/validator/validator.go#L115)
  - [x] `validNewGroupQuestion(NewGroupQuestion question)` -> [ValidNewGroupQuestion()](../internal/infra/validator/validator.go#L119)
  - [x] `validGroupQuestionIdAndAnswer(Map.Entry<Long, String> questionIdAndAnswer)` -> [ValidGroupQuestionIdAndAnswer()](../internal/infra/validator/validator.go#L123)

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

  - [x] `handleCreateMeetingRequest()` -> [HandleCreateMeetingRequest()](../internal/domain/conference/controller/conference_controller.go#L9)
  - [x] `handleDeleteMeetingRequest()` -> [HandleDeleteMeetingRequest()](../internal/domain/conference/controller/conference_controller.go#L13)
  - [x] `handleUpdateMeetingRequest()` -> [HandleUpdateMeetingRequest()](../internal/domain/conference/controller/conference_controller.go#L17)
  - [x] `handleQueryMeetingsRequest()` -> [HandleQueryMeetingsRequest()](../internal/domain/conference/controller/conference_controller.go#L21)
  - [x] `handleUpdateMeetingInvitationRequest()` -> [HandleUpdateMeetingInvitationRequest()](../internal/domain/conference/controller/conference_controller.go#L25)

- **MeetingRepository.java** ([java/im/turms/service/domain/conference/repository/MeetingRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/conference/repository/MeetingRepository.java))
> [简述功能]

  - [x] `updateEndDate(Long meetingId, Date endDate)` -> [UpdateEndDate()](../internal/domain/conference/repository/meeting_repository.go#L9)
  - [x] `updateCancelDateIfNotCanceled(Long meetingId, Date cancelDate)` -> [UpdateCancelDateIfNotCanceled()](../internal/domain/conference/repository/meeting_repository.go#L13)
  - [x] `updateMeeting(Long meetingId, @Nullable String name, @Nullable String intro, @Nullable String password)` -> [UpdateMeeting()](../internal/domain/conference/repository/meeting_repository.go#L17)
  - [x] `find(@Nullable Collection<Long> ids, @Nullable Collection<Long> creatorIds, @Nullable Collection<Long> userIds, @Nullable Collection<Long> groupIds, @Nullable Date creationDateStart, @Nullable Date creationDateEnd, @Nullable Integer skip, @Nullable Integer limit)` -> [Find()](../internal/domain/conference/repository/meeting_repository.go#L21)
  - [x] `find(@Nullable Collection<Long> ids, @NotNull Long creatorId, @NotNull Long userId, @Nullable Date creationDateStart, @Nullable Date creationDateEnd, @Nullable Integer skip, @Nullable Integer limit)` -> [FindByCreatorAndUser()](../internal/domain/conference/repository/meeting_repository.go#L25)

- **ConferenceService.java** ([java/im/turms/service/domain/conference/service/ConferenceService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/conference/service/ConferenceService.java))
> [简述功能]

  - [x] `onExtensionStarted(ConferenceServiceProvider extension)` -> [OnExtensionStarted()](../internal/domain/conference/service/conference_service.go#L9)
  - [x] `authAndCancelMeeting(@NotNull Long requesterId, @NotNull Long meetingId)` -> [AuthAndCancelMeeting()](../internal/domain/conference/service/conference_service.go#L13)
  - [x] `queryMeetingParticipants(@Nullable Long userId, @Nullable Long groupId)` -> [QueryMeetingParticipants()](../internal/domain/conference/service/conference_service.go#L17)
  - [x] `authAndUpdateMeeting(@NotNull Long requesterId, @NotNull Long meetingId, @Nullable String name, @Nullable String intro, @Nullable String password)` -> [AuthAndUpdateMeeting()](../internal/domain/conference/service/conference_service.go#L21)
  - [x] `authAndUpdateMeetingInvitation(@NotNull Long requesterId, @NotNull Long meetingId, @Nullable String password, @NotNull ResponseAction responseAction)` -> [AuthAndUpdateMeetingInvitation()](../internal/domain/conference/service/conference_service.go#L25)
  - [x] `authAndQueryMeetings(@NotNull Long requesterId, @Nullable Set<Long> ids, @Nullable Set<Long> creatorIds, @Nullable Set<Long> userIds, @Nullable Set<Long> groupIds, @Nullable Date creationDateStart, @Nullable Date creationDateEnd, @Nullable Integer skip, @Nullable Integer limit)` -> [AuthAndQueryMeetings()](../internal/domain/conference/service/conference_service.go#L29)

- **ConversationController.java** ([java/im/turms/service/domain/conversation/access/admin/controller/ConversationController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/conversation/access/admin/controller/ConversationController.java))
> [简述功能]

  - [x] `queryConversations(@QueryParam(required = false)` -> [QueryConversations()](../internal/domain/conversation/access/admin/controller/conversation_controller.go#L9)
  - [x] `deleteConversations(@QueryParam(required = false)` -> [DeleteConversations()](../internal/domain/conversation/access/admin/controller/conversation_controller.go#L13)
  - [x] `updateConversations(@QueryParam(required = false)` -> [UpdateConversations()](../internal/domain/conversation/access/admin/controller/conversation_controller.go#L17)

- **UpdateConversationDTO.java** ([java/im/turms/service/domain/conversation/access/admin/dto/request/UpdateConversationDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/conversation/access/admin/dto/request/UpdateConversationDTO.java))
> [简述功能]

  - [ ] `UpdateConversationDTO(Date readDate)`

- **ConversationsDTO.java** ([java/im/turms/service/domain/conversation/access/admin/dto/response/ConversationsDTO.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/conversation/access/admin/dto/response/ConversationsDTO.java))
> [简述功能]

  - [ ] `ConversationsDTO(List<PrivateConversation> privateConversations, List<GroupConversation> groupConversations)`

- **ConversationServiceController.java** ([java/im/turms/service/domain/conversation/access/servicerequest/controller/ConversationServiceController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/conversation/access/servicerequest/controller/ConversationServiceController.java))
> [简述功能]

  - [x] `handleQueryConversationsRequest()` -> [HandleQueryConversationsRequest()](../internal/domain/conversation/access/servicerequest/controller/conversation_service_controllers.go#L9)
  - [x] `handleUpdateTypingStatusRequest()` -> [HandleUpdateTypingStatusRequest()](../internal/domain/conversation/access/servicerequest/controller/conversation_service_controllers.go#L13)
  - [x] `handleUpdateConversationRequest()` -> [HandleUpdateConversationRequest()](../internal/domain/conversation/access/servicerequest/controller/conversation_service_controllers.go#L17)

- **ConversationSettingsServiceController.java** ([java/im/turms/service/domain/conversation/access/servicerequest/controller/ConversationSettingsServiceController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/conversation/access/servicerequest/controller/ConversationSettingsServiceController.java))
> [简述功能]

  - [x] `handleUpdateConversationSettingsRequest()` -> [HandleUpdateConversationSettingsRequest()](../internal/domain/conversation/access/servicerequest/controller/conversation_service_controllers.go#L26)
  - [x] `handleDeleteConversationSettingsRequest()` -> [HandleDeleteConversationSettingsRequest()](../internal/domain/conversation/access/servicerequest/controller/conversation_service_controllers.go#L30)
  - [x] `handleQueryConversationSettingsRequest()` -> [HandleQueryConversationSettingsRequest()](../internal/domain/conversation/access/servicerequest/controller/conversation_service_controllers.go#L34)

- **ConversationSettingsRepository.java** ([java/im/turms/service/domain/conversation/repository/ConversationSettingsRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/conversation/repository/ConversationSettingsRepository.java))
> [简述功能]

  - [x] `upsertSettings(Long ownerId, Long targetId, Map<String, Object> settings)` -> [UpsertSettings(ctx context.Context, userID int64, settings map[string]interface{})](../internal/domain/user/service/user_settings_service.go#L42)
  - [x] `unsetSettings(Long ownerId, @Nullable Collection<Long> targetIds, @Nullable Collection<String> settingNames)` -> [UnsetSettings(ctx context.Context, userID int64, keys []string)](../internal/domain/user/service/user_settings_service.go#L91)
  - [x] `findByIdAndSettingNames(Long ownerId, @Nullable Collection<String> settingNames, @Nullable Date lastUpdatedDateStart)` -> [FindByIdAndSettingNames()](../internal/domain/conversation/repository/conversation_settings_repository.go#L9)
  - [x] `findByIdAndSettingNames(Collection<ConversationSettings.Key> keys, @Nullable Collection<String> settingNames, @Nullable Date lastUpdatedDateStart)` -> [FindByIdAndSettingNamesWithKeys()](../internal/domain/conversation/repository/conversation_settings_repository.go#L13)
  - [x] `findSettingFields(Long ownerId, Long targetId, Collection<String> includedFields)` -> [FindSettingFields()](../internal/domain/conversation/repository/conversation_settings_repository.go#L17)
  - [x] `deleteByOwnerIds(Collection<Long> ownerIds, @Nullable ClientSession clientSession)` -> [DeleteByOwnerIds()](../internal/domain/conversation/repository/conversation_settings_repository.go#L21)

- **GroupConversationRepository.java** ([java/im/turms/service/domain/conversation/repository/GroupConversationRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/conversation/repository/GroupConversationRepository.java))
> [简述功能]

  - [x] `upsert(Long groupId, Long memberId, Date readDate, boolean allowMoveReadDateForward)` -> [Upsert(ctx context.Context, groupID int64, timestamp time.Time)](../internal/domain/group/repository/group_version_repository.go#L99)
  - [x] `upsert(Long groupId, Collection<Long> memberIds, Date readDate)` -> [Upsert(ctx context.Context, groupID int64, timestamp time.Time)](../internal/domain/group/repository/group_version_repository.go#L99)
  - [x] `deleteMemberConversations(Collection<Long> groupIds, Long memberId, ClientSession session)` -> [DeleteMemberConversations()](../internal/domain/conversation/repository/group_conversation_repository.go#L69)

- **PrivateConversationRepository.java** ([java/im/turms/service/domain/conversation/repository/PrivateConversationRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/conversation/repository/PrivateConversationRepository.java))
> [简述功能]

  - [x] `upsert(Set<PrivateConversation.Key> keys, Date readDate, boolean allowMoveReadDateForward)` -> [Upsert(ctx context.Context, groupID int64, timestamp time.Time)](../internal/domain/group/repository/group_version_repository.go#L99)
  - [x] `deleteConversationsByOwnerIds(Set<Long> ownerIds, @Nullable ClientSession session)` -> [DeleteConversationsByOwnerIds()](../internal/domain/conversation/repository/private_conversation_repository.go#L68)
  - [x] `findConversations(Collection<Long> ownerIds)` -> [FindConversations()](../internal/domain/conversation/repository/private_conversation_repository.go#L72)

- **ConversationService.java** ([java/im/turms/service/domain/conversation/service/ConversationService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/conversation/service/ConversationService.java))
> [简述功能]

  - [ ] `authAndUpsertGroupConversationReadDate(@NotNull Long groupId, @NotNull Long memberId, @Nullable @PastOrPresent Date readDate)`
  - [ ] `authAndUpsertPrivateConversationReadDate(@NotNull Long ownerId, @NotNull Long targetId, @Nullable @PastOrPresent Date readDate)`
  - [x] `upsertGroupConversationReadDate(@NotNull Long groupId, @NotNull Long memberId, @Nullable @PastOrPresent Date readDate)` -> [UpsertGroupConversationReadDate()](../internal/domain/conversation/service/conversation_service.go#L58)
  - [x] `upsertGroupConversationsReadDate(@NotNull Set<GroupConversation.GroupConversionMemberKey> keys, @Nullable @PastOrPresent Date readDate)` -> [UpsertGroupConversationsReadDate()](../internal/domain/conversation/service/conversation_service.go#L62)
  - [x] `upsertPrivateConversationReadDate(@NotNull Long ownerId, @NotNull Long targetId, @Nullable @PastOrPresent Date readDate)` -> [UpsertPrivateConversationReadDate()](../internal/domain/conversation/service/conversation_service.go#L66)
  - [x] `upsertPrivateConversationsReadDate(@NotNull Set<PrivateConversation.Key> keys, @Nullable @PastOrPresent Date readDate)` -> [UpsertPrivateConversationsReadDate()](../internal/domain/conversation/service/conversation_service.go#L70)
  - [x] `queryGroupConversations(@NotNull Collection<Long> groupIds)` -> [QueryGroupConversations(ctx context.Context, groupIDs []int64)](../internal/domain/conversation/repository/group_conversation_repository.go#L47)
  - [ ] `queryPrivateConversationsByOwnerIds(@NotNull Set<Long> ownerIds)`
  - [x] `queryPrivateConversations(@NotNull Collection<Long> ownerIds, @NotNull Long targetId)` -> [QueryPrivateConversations(ctx context.Context, ownerIDs []int64)](../internal/domain/conversation/repository/private_conversation_repository.go#L44)
  - [x] `queryPrivateConversations(@NotNull Set<PrivateConversation.Key> keys)` -> [QueryPrivateConversations(ctx context.Context, ownerIDs []int64)](../internal/domain/conversation/repository/private_conversation_repository.go#L44)
  - [x] `deletePrivateConversations(@NotNull Set<PrivateConversation.Key> keys)` -> [DeletePrivateConversationsByKeys()](../internal/domain/conversation/service/conversation_service.go#L74)
  - [x] `deletePrivateConversations(@NotNull Set<Long> userIds, @Nullable ClientSession session)` -> [DeletePrivateConversationsByUserIds()](../internal/domain/conversation/service/conversation_service.go#L78)
  - [x] `deleteGroupConversations(@Nullable Set<Long> groupIds, @Nullable ClientSession session)` -> [DeleteGroupConversations()](../internal/domain/conversation/service/conversation_service.go#L82)
  - [x] `deleteGroupMemberConversations(@NotNull Collection<Long> userIds, @Nullable ClientSession session)` -> [DeleteGroupMemberConversations()](../internal/domain/conversation/service/conversation_service.go#L86)
  - [x] `authAndUpdateTypingStatus(@NotNull Long requesterId, boolean isGroupMessage, @NotNull Long toId)` -> [AuthAndUpdateTypingStatus()](../internal/domain/conversation/service/conversation_service.go#L90)

- **ConversationSettingsService.java** ([java/im/turms/service/domain/conversation/service/ConversationSettingsService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/conversation/service/ConversationSettingsService.java))
> [简述功能]

  - [x] `upsertPrivateConversationSettings(Long ownerId, Long userId, Map<String, Value> settings)` -> [UpsertPrivateConversationSettings()](../internal/domain/conversation/service/conversation_settings_service.go#L9)
  - [x] `upsertGroupConversationSettings(Long ownerId, Long groupId, Map<String, Value> settings)` -> [UpsertGroupConversationSettings()](../internal/domain/conversation/service/conversation_settings_service.go#L13)
  - [x] `deleteSettings(Collection<Long> ownerIds, @Nullable ClientSession clientSession)` -> [DeleteSettings(ctx context.Context, filter interface{})](../internal/domain/user/repository/user_settings_repository.go#L46)
  - [x] `unsetSettings(Long ownerId, @Nullable Set<Long> userIds, @Nullable Set<Long> groupIds, @Nullable Set<String> settingNames)` -> [UnsetSettings(ctx context.Context, userID int64, keys []string)](../internal/domain/user/service/user_settings_service.go#L91)
  - [x] `querySettings(Long ownerId, @Nullable Collection<Long> userIds, @Nullable Collection<Long> groupIds, @Nullable Set<String> settingNames, @Nullable Date lastUpdatedDateStart)` -> [QuerySettings(ctx context.Context, filter bson.M)](../internal/domain/user/service/user_settings_service.go#L100)

- **GroupBlocklistController.java** ([java/im/turms/service/domain/group/access/admin/controller/GroupBlocklistController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/access/admin/controller/GroupBlocklistController.java))
> [简述功能]

  - [x] `addGroupBlockedUser(@RequestBody AddGroupBlockedUserDTO addGroupBlockedUserDTO)` -> [AddGroupBlockedUser(addGroupBlockedUserDTO any)](../internal/domain/group/access/admin/controller/group_controllers.go#L8)
  - [x] `queryGroupBlockedUsers(@QueryParam(required = false)` -> [QueryGroupBlockedUsers(page, size *int)](../internal/domain/group/access/admin/controller/group_controllers.go#L11)
  - [x] `queryGroupBlockedUsers(@QueryParam(required = false)` -> [QueryGroupBlockedUsers(page, size *int)](../internal/domain/group/access/admin/controller/group_controllers.go#L11)
  - [x] `updateGroupBlockedUsers(List<GroupBlockedUser.Key> keys, @RequestBody UpdateGroupBlockedUserDTO updateGroupBlockedUserDTO)` -> [UpdateGroupBlockedUsers(keys []any, updateGroupBlockedUserDTO any)](../internal/domain/group/access/admin/controller/group_controllers.go#L17)
  - [x] `deleteGroupBlockedUsers(List<GroupBlockedUser.Key> keys)` -> [DeleteGroupBlockedUsers(keys []any)](../internal/domain/group/access/admin/controller/group_controllers.go#L20)

- **GroupController.java** ([java/im/turms/service/domain/group/access/admin/controller/GroupController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/access/admin/controller/GroupController.java))
> [简述功能]

  - [x] `addGroup(@RequestBody AddGroupDTO addGroupDTO)` -> [AddGroup(addGroupDTO any)](../internal/domain/group/access/admin/controller/group_controllers.go#L27)
  - [x] `queryGroups(@QueryParam(required = false)` -> [QueryGroups(page, size *int)](../internal/domain/group/access/admin/controller/group_controllers.go#L28), [QueryGroups(ctx context.Context, groupIDs []int64, name *string, lastUpdatedDate *time.Time, skip *int32, limit *int32)](../internal/domain/group/repository/group_repository.go#L57)
  - [x] `queryGroups(@QueryParam(required = false)` -> [QueryGroups(page, size *int)](../internal/domain/group/access/admin/controller/group_controllers.go#L28), [QueryGroups(ctx context.Context, groupIDs []int64, name *string, lastUpdatedDate *time.Time, skip *int32, limit *int32)](../internal/domain/group/repository/group_repository.go#L57)
  - [x] `countGroups(@QueryParam(required = false)` -> [CountGroups(ctx context.Context, ownerIDs []int64, relatedUserIDs []int64)](../internal/domain/user/repository/user_relationship_group_member_repository.go#L195)
  - [x] `updateGroups(Set<Long> ids, @RequestBody UpdateGroupDTO updateGroupDTO)` -> [UpdateGroups(ids []int64, updateGroupDTO any)](../internal/domain/group/access/admin/controller/group_controllers.go#L32)
  - [x] `deleteGroups(@QueryParam(required = false)` -> [DeleteGroups(ids []int64, deleteLogical *bool)](../internal/domain/group/access/admin/controller/group_controllers.go#L35)

- **GroupInvitationController.java** ([java/im/turms/service/domain/group/access/admin/controller/GroupInvitationController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/access/admin/controller/GroupInvitationController.java))
> [简述功能]

  - [x] `addGroupInvitation(@RequestBody AddGroupInvitationDTO addGroupInvitationDTO)` -> [AddGroupInvitation(addGroupInvitationDTO any)](../internal/domain/group/access/admin/controller/group_controllers.go#L42)
  - [x] `queryGroupInvitations(@QueryParam(required = false)` -> [QueryGroupInvitations(page, size *int)](../internal/domain/group/access/admin/controller/group_controllers.go#L45)
  - [x] `queryGroupInvitations(@QueryParam(required = false)` -> [QueryGroupInvitations(page, size *int)](../internal/domain/group/access/admin/controller/group_controllers.go#L45)
  - [x] `updateGroupInvitations(Set<Long> ids, @RequestBody UpdateGroupInvitationDTO updateGroupInvitationDTO)` -> [UpdateGroupInvitations(ids []int64, updateGroupInvitationDTO any)](../internal/domain/group/access/admin/controller/group_controllers.go#L51)
  - [x] `deleteGroupInvitations(@QueryParam(required = false)` -> [DeleteGroupInvitations(ids []int64)](../internal/domain/group/access/admin/controller/group_controllers.go#L54)

- **GroupJoinRequestController.java** ([java/im/turms/service/domain/group/access/admin/controller/GroupJoinRequestController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/access/admin/controller/GroupJoinRequestController.java))
> [简述功能]

  - [x] `addGroupJoinRequest(@RequestBody AddGroupJoinRequestDTO addGroupJoinRequestDTO)` -> [AddGroupJoinRequest(addGroupJoinRequestDTO any)](../internal/domain/group/access/admin/controller/group_controllers.go#L61)
  - [x] `queryGroupJoinRequests(@QueryParam(required = false)` -> [QueryGroupJoinRequests(page, size *int)](../internal/domain/group/access/admin/controller/group_controllers.go#L64)
  - [x] `queryGroupJoinRequests(@QueryParam(required = false)` -> [QueryGroupJoinRequests(page, size *int)](../internal/domain/group/access/admin/controller/group_controllers.go#L64)
  - [x] `updateGroupJoinRequests(Set<Long> ids, @RequestBody UpdateGroupJoinRequestDTO updateGroupJoinRequestDTO)` -> [UpdateGroupJoinRequests(ids []int64, updateGroupJoinRequestDTO any)](../internal/domain/group/access/admin/controller/group_controllers.go#L70)
  - [x] `deleteGroupJoinRequests(@QueryParam(required = false)` -> [DeleteGroupJoinRequests(ids []int64)](../internal/domain/group/access/admin/controller/group_controllers.go#L73)

- **GroupMemberController.java** ([java/im/turms/service/domain/group/access/admin/controller/GroupMemberController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/access/admin/controller/GroupMemberController.java))
> [简述功能]

  - [x] `addGroupMember(@RequestBody AddGroupMemberDTO addGroupMemberDTO)` -> [AddGroupMember(ctx context.Context, member *po.GroupMember)](../internal/domain/group/repository/group_member_repository.go#L34)
  - [x] `queryGroupMembers(@QueryParam(required = false)` -> [QueryGroupMembers(page, size *int)](../internal/domain/group/access/admin/controller/group_controllers.go#L80)
  - [x] `queryGroupMembers(@QueryParam(required = false)` -> [QueryGroupMembers(page, size *int)](../internal/domain/group/access/admin/controller/group_controllers.go#L80)
  - [x] `updateGroupMembers(List<GroupMember.Key> keys, @RequestBody UpdateGroupMemberDTO updateGroupMemberDTO)` -> [UpdateGroupMembers(ctx context.Context, keys []po.GroupMemberKey, name *string, role *protocol.GroupMemberRole, joinDate *time.Time, muteEndDate *time.Time)](../internal/domain/group/repository/group_member_repository.go#L163)
  - [x] `deleteGroupMembers(@QueryParam(required = false)` -> [DeleteGroupMembers(keys []any, successorId *int64, quitAfterTransfer *bool)](../internal/domain/group/access/admin/controller/group_controllers.go#L86)

- **GroupQuestionController.java** ([java/im/turms/service/domain/group/access/admin/controller/GroupQuestionController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/access/admin/controller/GroupQuestionController.java))
> [简述功能]

  - [x] `queryGroupJoinQuestions(@QueryParam(required = false)` -> [QueryGroupJoinQuestions(page, size *int)](../internal/domain/group/access/admin/controller/group_controllers.go#L93)
  - [x] `queryGroupJoinQuestions(@QueryParam(required = false)` -> [QueryGroupJoinQuestions(page, size *int)](../internal/domain/group/access/admin/controller/group_controllers.go#L93)
  - [x] `addGroupJoinQuestion(@RequestBody AddGroupJoinQuestionDTO addGroupJoinQuestionDTO)` -> [AddGroupJoinQuestion(addGroupJoinQuestionDTO any)](../internal/domain/group/access/admin/controller/group_controllers.go#L99)
  - [x] `updateGroupJoinQuestions(Set<Long> ids, @RequestBody UpdateGroupJoinQuestionDTO updateGroupJoinQuestionDTO)` -> [UpdateGroupJoinQuestions(ids []int64, updateGroupJoinQuestionDTO any)](../internal/domain/group/access/admin/controller/group_controllers.go#L102)
  - [x] `deleteGroupJoinQuestions(@QueryParam(required = false)` -> [DeleteGroupJoinQuestions(ids []int64)](../internal/domain/group/access/admin/controller/group_controllers.go#L105)

- **GroupTypeController.java** ([java/im/turms/service/domain/group/access/admin/controller/GroupTypeController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/access/admin/controller/GroupTypeController.java))
> [简述功能]

  - [x] `addGroupType(@RequestBody AddGroupTypeDTO addGroupTypeDTO)` -> [AddGroupType(addGroupTypeDTO any)](../internal/domain/group/access/admin/controller/group_controllers.go#L112)
  - [x] `queryGroupTypes(@QueryParam(required = false)` -> [QueryGroupTypes(page, size *int)](../internal/domain/group/access/admin/controller/group_controllers.go#L113)
  - [x] `queryGroupTypes(int page, @QueryParam(required = false)` -> [QueryGroupTypes(page, size *int)](../internal/domain/group/access/admin/controller/group_controllers.go#L113)
  - [x] `updateGroupType(Set<Long> ids, @RequestBody UpdateGroupTypeDTO updateGroupTypeDTO)` -> [UpdateGroupType(ctx context.Context, typeID int64, update bson.M)](../internal/domain/group/repository/group_type_repository.go#L49)
  - [x] `deleteGroupType(Set<Long> ids)` -> [DeleteGroupType(ids []int64)](../internal/domain/group/access/admin/controller/group_controllers.go#L117)

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

  - [x] `handleCreateGroupRequest()` -> [HandleCreateGroupRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go#L95)
  - [x] `handleDeleteGroupRequest()` -> [HandleDeleteGroupRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go#L113)
  - [x] `handleQueryGroupsRequest()` -> [HandleQueryGroupsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go#L123)
  - [x] `handleQueryJoinedGroupIdsRequest()` -> [HandleQueryJoinedGroupIdsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go#L190)
  - [x] `handleQueryJoinedGroupsRequest()` -> [HandleQueryJoinedGroupsRequest()](../internal/domain/group/access/servicerequest/controller/group_service_controller.go#L6)
  - [x] `handleUpdateGroupRequest()` -> [HandleUpdateGroupRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go#L314)
  - [x] `handleCreateGroupBlockedUserRequest()` -> [HandleCreateGroupBlockedUserRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go#L490)
  - [x] `handleDeleteGroupBlockedUserRequest()` -> [HandleDeleteGroupBlockedUserRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go#L500)
  - [x] `handleQueryGroupBlockedUserIdsRequest()` -> [HandleQueryGroupBlockedUserIdsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go#L510)
  - [x] `handleQueryGroupBlockedUsersInfosRequest()` -> [HandleQueryGroupBlockedUsersInfosRequest()](../internal/domain/group/access/servicerequest/controller/group_service_controller.go#L7)
  - [x] `handleCheckGroupQuestionAnswerRequest()` -> [HandleCheckGroupQuestionAnswerRequest()](../internal/domain/group/access/servicerequest/controller/group_service_controller.go#L8)
  - [x] `handleCreateGroupInvitationRequestRequest()` -> [HandleCreateGroupInvitationRequestRequest()](../internal/domain/group/access/servicerequest/controller/group_service_controller.go#L9)
  - [x] `handleCreateGroupJoinRequestRequest()` -> [HandleCreateGroupJoinRequestRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go#L735)
  - [x] `handleCreateGroupQuestionsRequest()` -> [HandleCreateGroupQuestionsRequest()](../internal/domain/group/access/servicerequest/controller/group_service_controller.go#L10)
  - [x] `handleDeleteGroupInvitationRequest()` -> [HandleDeleteGroupInvitationRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go#L631)
  - [x] `handleUpdateGroupInvitationRequest()` -> [HandleUpdateGroupInvitationRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go#L722)
  - [x] `handleDeleteGroupJoinRequestRequest()` -> [HandleDeleteGroupJoinRequestRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go#L753)
  - [x] `handleUpdateGroupJoinRequestRequest()` -> [HandleUpdateGroupJoinRequestRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go#L841)
  - [x] `handleDeleteGroupJoinQuestionsRequest()` -> [HandleDeleteGroupJoinQuestionsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go#L877)
  - [x] `handleQueryGroupInvitationsRequest()` -> [HandleQueryGroupInvitationsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go#L641)
  - [x] `handleQueryGroupJoinRequestsRequest()` -> [HandleQueryGroupJoinRequestsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go#L763)
  - [x] `handleQueryGroupJoinQuestionsRequest()` -> [HandleQueryGroupJoinQuestionsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go#L889)
  - [x] `handleUpdateGroupJoinQuestionRequest()` -> [HandleUpdateGroupJoinQuestionRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go#L937)
  - [x] `handleCreateGroupMembersRequest()` -> [HandleCreateGroupMembersRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go#L340)
  - [x] `handleDeleteGroupMembersRequest()` -> [HandleDeleteGroupMembersRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go#L376)
  - [x] `handleQueryGroupMembersRequest()` -> [HandleQueryGroupMembersRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go#L394), [HandleQueryGroupMembersRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go#L394)
  - [x] `handleUpdateGroupMemberRequest()` -> [HandleUpdateGroupMemberRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/group/controller/group_service_controller.go#L460)

- **CheckGroupQuestionAnswerResult.java** ([java/im/turms/service/domain/group/bo/CheckGroupQuestionAnswerResult.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/bo/CheckGroupQuestionAnswerResult.java))
> [简述功能]

  - [ ] `CheckGroupQuestionAnswerResult(boolean joined, Long groupId, List<Long> questionIds, Integer score)`

- **GroupInvitationStrategy.java** ([java/im/turms/service/domain/group/bo/GroupInvitationStrategy.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/bo/GroupInvitationStrategy.java))
> [简述功能]

  - [x] `requiresApproval()` -> [RequiresApproval()](../internal/domain/group/constant/group_strategy.go#L19)

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

  - [x] `updateBlockedUsers(Set<GroupBlockedUser.Key> keys, @Nullable Date blockDate, @Nullable Long requesterId)` -> [UpdateBlockedUsers(ctx context.Context, keys []po.GroupBlockedUserKey, blockDate *time.Time, requesterId *int64)](../internal/domain/group/repository/group_blocked_user_repository.go#L108)
  - [x] `count(@Nullable Set<Long> groupIds, @Nullable Set<Long> userIds, @Nullable DateRange blockDateRange, @Nullable Set<Long> requesterIds)` -> [Count(ctx context.Context, filter bson.M)](../internal/domain/user/repository/user_repository.go#L104)
  - [x] `findBlockedUserIds(Long groupId)` -> [FindBlockedUserIds(ctx context.Context, groupId int64)](../internal/domain/group/repository/group_blocked_user_repository.go#L128)
  - [x] `findBlockedUsers(@Nullable Set<Long> groupIds, @Nullable Set<Long> userIds, @Nullable DateRange blockDateRange, @Nullable Set<Long> requesterIds, @Nullable Integer page, @Nullable Integer size)` -> [FindBlockedUsers(ctx context.Context, groupIds []int64, userIds []int64, blockDateRange *turmsmongo.DateRange, requesterIds []int64, page *int, size *int,)](../internal/domain/group/repository/group_blocked_user_repository.go#L152)

- **GroupInvitationRepository.java** ([java/im/turms/service/domain/group/repository/GroupInvitationRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/repository/GroupInvitationRepository.java))
> [简述功能]

  - [x] `getEntityExpireAfterSeconds()` -> [GetEntityExpireAfterSeconds()](../internal/domain/user/repository/user_friend_request_repository.go#L135)
  - [x] `updateStatusIfPending(Long invitationId, RequestStatus requestStatus, @Nullable String reason, @Nullable ClientSession session)` -> [UpdateStatusIfPending(ctx context.Context, invitationID int64, newStatus po.RequestStatus, reason *string, responseDate time.Time)](../internal/domain/group/repository/group_invitation_repository.go#L64)
  - [x] `updateInvitations(Set<Long> invitationIds, @Nullable Long inviterId, @Nullable Long inviteeId, @Nullable String content, @Nullable RequestStatus status, @Nullable Date creationDate, @Nullable Date responseDate)` -> [UpdateInvitations(ctx context.Context, ids []int64, inviterID, inviteeID *int64, content *string, status *po.RequestStatus, creationDate, responseDate *time.Time)](../internal/domain/group/repository/group_invitation_repository.go#L196)
  - [x] `count(@Nullable Set<Long> ids, @Nullable Set<Long> groupIds, @Nullable Set<Long> inviterIds, @Nullable Set<Long> inviteeIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange)` -> [Count(ctx context.Context, filter bson.M)](../internal/domain/user/repository/user_repository.go#L104)
  - [x] `findGroupIdAndInviteeIdAndStatus(Long invitationId)` -> [FindGroupIdAndInviteeIdAndStatus(ctx context.Context, id int64)](../internal/domain/group/repository/group_invitation_repository.go#L170)
  - [x] `findGroupIdAndInviterIdAndInviteeIdAndStatus(Long invitationId)` -> [FindGroupIdAndInviterIdAndInviteeIdAndStatus(ctx context.Context, id int64)](../internal/domain/group/repository/group_invitation_repository.go#L157)
  - [x] `findInvitationsByInviteeId(Long inviteeId)` -> [FindInvitationsByInviteeID(ctx context.Context, inviteeID int64)](../internal/domain/group/repository/group_invitation_repository.go#L88)
  - [x] `findInvitationsByInviterId(Long inviterId)` -> [FindInvitationsByInviterId(ctx context.Context, inviterId int64)](../internal/domain/group/repository/group_invitation_repository.go#L264)
  - [x] `findInvitationsByGroupId(Long groupId)` -> [FindInvitationsByGroupID(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_invitation_repository.go#L105)
  - [x] `findInviteeIdAndGroupIdAndCreationDateAndStatus(Long invitationId)` -> [FindInviteeIdAndGroupIdAndCreationDateAndStatus(ctx context.Context, id int64)](../internal/domain/group/repository/group_invitation_repository.go#L182)
  - [x] `findInvitations(@Nullable Set<Long> ids, @Nullable Set<Long> groupIds, @Nullable Set<Long> inviterIds, @Nullable Set<Long> inviteeIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange, @Nullable Integer page, @Nullable Integer size)` -> [FindInvitations(ctx context.Context, groupID *int64, inviterID *int64, inviteeID *int64, status *po.RequestStatus, lastUpdatedDate *time.Time, page, size int)](../internal/domain/group/repository/group_invitation_repository.go#L135)

- **GroupJoinRequestRepository.java** ([java/im/turms/service/domain/group/repository/GroupJoinRequestRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/repository/GroupJoinRequestRepository.java))
> [简述功能]

  - [x] `getEntityExpireAfterSeconds()` -> [GetEntityExpireAfterSeconds()](../internal/domain/user/repository/user_friend_request_repository.go#L135)
  - [x] `updateStatusIfPending(Long requestId, RequestStatus status, Long responderId, @Nullable String reason, @Nullable ClientSession session)` -> [UpdateStatusIfPending(ctx context.Context, invitationID int64, newStatus po.RequestStatus, reason *string, responseDate time.Time)](../internal/domain/group/repository/group_invitation_repository.go#L64)
  - [x] `updateRequests(Set<Long> requestIds, @Nullable Long requesterId, @Nullable Long responderId, @Nullable String content, @Nullable RequestStatus status, @Nullable Date creationDate, @Nullable Date responseDate)` -> [UpdateRequests(ctx context.Context, requestIds []int64, requesterId, responderId *int64, content *string, status *po.RequestStatus, creationDate, responseDate *time.Time)](../internal/domain/group/repository/group_join_request_repository.go#L189)
  - [x] `countRequests(@Nullable Set<Long> ids, @Nullable Set<Long> groupIds, @Nullable Set<Long> requesterIds, @Nullable Set<Long> responderIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange)` -> [CountRequests(ctx context.Context, ids, groupIds, requesterIds, responderIds []int64, statuses []po.RequestStatus, creationDateRange, responseDateRange, expirationDateRange *turmsmongo.DateRange)](../internal/domain/group/repository/group_join_request_repository.go#L222)
  - [x] `findGroupId(Long requestId)` -> [FindGroupId(ctx context.Context, requestId int64)](../internal/domain/group/repository/group_join_request_repository.go#L258)
  - [x] `findRequesterIdAndStatusAndGroupId(Long requestId)` -> [FindRequesterIdAndStatusAndGroupId(ctx context.Context, requestId int64)](../internal/domain/group/repository/group_join_request_repository.go#L274)
  - [x] `findRequestsByGroupId(Long groupId)` -> [FindRequestsByGroupID(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_join_request_repository.go#L93)
  - [x] `findRequestsByRequesterId(Long requesterId)` -> [FindRequestsByRequesterID(ctx context.Context, requesterID int64)](../internal/domain/group/repository/group_join_request_repository.go#L110)
  - [x] `findRequests(@Nullable Set<Long> ids, @Nullable Set<Long> groupIds, @Nullable Set<Long> requesterIds, @Nullable Set<Long> responderIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange, @Nullable Integer page, @Nullable Integer size)` -> [FindRequests(ctx context.Context, groupID *int64, requesterID *int64, responderID *int64, status *po.RequestStatus, creationDate *time.Time, responseDate *time.Time, expirationDate *time.Time, page int, size int)](../internal/domain/group/repository/group_join_request_repository.go#L136)

- **GroupMemberRepository.java** ([java/im/turms/service/domain/group/repository/GroupMemberRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/repository/GroupMemberRepository.java))
> [简述功能]

  - [x] `deleteAllGroupMembers(@Nullable Set<Long> groupIds, @Nullable ClientSession session)` -> [DeleteAllGroupMembers(ctx context.Context, groupIDs []int64, session mongo.SessionContext, updateVersion bool)](../internal/domain/group/service/group_member_service.go#L348)
  - [x] `updateGroupMembers(Set<GroupMember.Key> keys, @Nullable String name, @Nullable GroupMemberRole role, @Nullable Date joinDate, @Nullable Date muteEndDate, @Nullable ClientSession session)` -> [UpdateGroupMembers(ctx context.Context, keys []po.GroupMemberKey, name *string, role *protocol.GroupMemberRole, joinDate *time.Time, muteEndDate *time.Time)](../internal/domain/group/repository/group_member_repository.go#L163)
  - [x] `countMembers(Long groupId)` -> [CountMembers(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_member_repository.go#L208)
  - [x] `countMembers(@Nullable Set<Long> groupIds, @Nullable Set<Long> userIds, @Nullable Set<GroupMemberRole> roles, @Nullable DateRange joinDateRange, @Nullable DateRange muteEndDateRange)` -> [CountMembers(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_member_repository.go#L208)
  - [ ] `findGroupManagersAndOwnerId(Long groupId)`
  - [x] `findGroupMembers(Long groupId)` -> [FindGroupMembers(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_member_repository.go#L277)
  - [x] `findGroupMembers(Long groupId, Set<Long> memberIds)` -> [FindGroupMembersWithIds(ctx context.Context, groupID int64, memberIDs []int64)](../internal/domain/group/repository/group_member_repository.go#L297)
  - [x] `findGroupsMembers(@Nullable Set<Long> groupIds, @Nullable Set<Long> userIds, @Nullable Set<@ValidGroupMemberRole GroupMemberRole> roles, @Nullable DateRange joinDateRange, @Nullable DateRange muteEndDateRange, @Nullable Integer page, @Nullable Integer size)` -> [FindGroupsMembers(ctx context.Context, groupIds, userIds []int64, roles []int, joinDateRange, muteEndDateRange any, page, size *int)](../internal/domain/group/repository/group_member_repository.go#L317)
  - [x] `findGroupMemberIds(Long groupId)` -> [FindGroupMemberIDs(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_member_repository.go#L73)
  - [x] `findGroupMemberIds(Set<Long> groupIds)` -> [FindGroupMemberIDs(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_member_repository.go#L73)
  - [ ] `findGroupMemberKeyAndRoleParis(Set<Long> userIds, Long groupId)`
  - [x] `findGroupMemberRole(Long userId, Long groupId)` -> [FindGroupMemberRole(ctx context.Context, groupID, userID int64)](../internal/domain/group/repository/group_member_repository.go#L54)
  - [ ] `findMemberIdsByGroupId(Long groupId)`
  - [x] `findUserJoinedGroupIds(Long userId)` -> [FindUserJoinedGroupIDs(ctx context.Context, userID int64)](../internal/domain/group/repository/group_member_repository.go#L115)
  - [ ] `findUsersJoinedGroupIds(@Nullable Set<Long> groupIds, @NotEmpty Set<Long> userIds, @Nullable Integer page, @Nullable Integer size)`
  - [x] `isMemberMuted(Long groupId, Long userId)` -> [IsMemberMuted(ctx context.Context, groupID, userID int64)](../internal/domain/group/repository/group_member_repository.go#L100)

- **GroupQuestionRepository.java** ([java/im/turms/service/domain/group/repository/GroupQuestionRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/repository/GroupQuestionRepository.java))
> [简述功能]

  - [ ] `updateQuestion(Long questionId, @Nullable String question, @Nullable Set<String> answers, @Nullable Integer score)`
  - [x] `updateQuestions(Set<Long> ids, @Nullable Long groupId, @Nullable String question, @Nullable Set<String> answers, @Nullable Integer score)` -> [UpdateQuestions(ctx context.Context, ids []int64, groupID *int64, question *string, answers []string, score *int)](../internal/domain/group/repository/group_join_question_repository.go#L136), [UpdateQuestions(ctx context.Context, ids []int64, groupID *int64, question *string, answers []string, score *int)](../internal/domain/group/service/group_question_service.go#L295)
  - [x] `countQuestions(@Nullable Set<Long> ids, @Nullable Set<Long> groupIds)` -> [CountQuestions(ctx context.Context, ids []int64, groupIds []int64)](../internal/domain/group/repository/group_join_question_repository.go#L124), [CountQuestions(ctx context.Context, ids []int64, groupIds []int64)](../internal/domain/group/service/group_question_service.go#L278)
  - [x] `checkQuestionAnswerAndGetScore(Long questionId, String answer, @Nullable Long groupId)` -> [CheckQuestionAnswerAndGetScore(ctx context.Context, questionID int64, answer string, groupID *int64)](../internal/domain/group/repository/group_join_question_repository.go#L165), [CheckQuestionAnswerAndGetScore(ctx context.Context, questionId int64, answer string, groupID *int64)](../internal/domain/group/service/group_question_service.go#L261)
  - [x] `findGroupId(Long questionId)` -> [FindGroupId(ctx context.Context, requestId int64)](../internal/domain/group/repository/group_join_request_repository.go#L258)
  - [x] `findQuestions(@Nullable Set<Long> ids, @Nullable Set<Long> groupIds, @Nullable Integer page, @Nullable Integer size, boolean withAnswers)` -> [FindQuestions(ctx context.Context, ids []int64, groupIds []int64, page *int, size *int)](../internal/domain/group/repository/group_join_question_repository.go#L94), [FindQuestions(ctx context.Context, ids []int64, groupIds []int64, page *int, size *int, withAnswers bool)](../internal/domain/group/service/group_question_service.go#L282)

- **GroupRepository.java** ([java/im/turms/service/domain/group/repository/GroupRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/repository/GroupRepository.java))
> [简述功能]

  - [x] `updateGroupsDeletionDate(@Nullable Collection<Long> groupIds, @Nullable ClientSession session)` -> [UpdateGroupsDeletionDate(ctx context.Context, groupIDs []int64, deletionDate *time.Time, session mongo.SessionContext)](../internal/domain/group/repository/group_repository.go#L145)
  - [x] `updateGroups(Set<Long> groupIds, @Nullable Long typeId, @Nullable Long creatorId, @Nullable Long ownerId, @Nullable String name, @Nullable String intro, @Nullable String announcement, @Nullable Integer minimumScore, @Nullable Boolean isActive, @Nullable Date creationDate, @Nullable Date deletionDate, @Nullable Date muteEndDate, @Nullable Date lastUpdatedDate, @Nullable Map<String, Object> userDefinedAttributes, @Nullable ClientSession session)` -> [UpdateGroups(ids []int64, updateGroupDTO any)](../internal/domain/group/access/admin/controller/group_controllers.go#L32)
  - [x] `countCreatedGroups(@Nullable DateRange dateRange)` -> [CountCreatedGroups(ctx context.Context, dateRange *turmsmongo.DateRange)](../internal/domain/group/repository/group_repository.go#L164)
  - [x] `countDeletedGroups(@Nullable DateRange dateRange)` -> [CountDeletedGroups(ctx context.Context, dateRange *turmsmongo.DateRange)](../internal/domain/group/repository/group_repository.go#L182)
  - [x] `countGroups(@Nullable Set<Long> ids, @Nullable Set<Long> typeIds, @Nullable Set<Long> creatorIds, @Nullable Set<Long> ownerIds, @Nullable Boolean isActive, @Nullable DateRange creationDateRange, @Nullable DateRange deletionDateRange, @Nullable DateRange lastUpdatedDateRange, @Nullable DateRange muteEndDateRange)` -> [CountGroups(ctx context.Context, ownerIDs []int64, relatedUserIDs []int64)](../internal/domain/user/repository/user_relationship_group_member_repository.go#L195)
  - [x] `countOwnedGroups(Long ownerId)` -> [CountOwnedGroups(ctx context.Context, ownerID int64)](../internal/domain/group/repository/group_repository.go#L132)
  - [x] `countOwnedGroups(Long ownerId, Long groupTypeId)` -> [CountOwnedGroups(ctx context.Context, ownerID int64)](../internal/domain/group/repository/group_repository.go#L132)
  - [x] `findGroups(@Nullable Set<Long> ids, @Nullable Set<Long> typeIds, @Nullable Set<Long> creatorIds, @Nullable Set<Long> ownerIds, @Nullable Boolean isActive, @Nullable DateRange creationDateRange, @Nullable DateRange deletionDateRange, @Nullable DateRange lastUpdatedDateRange, @Nullable DateRange muteEndDateRange, @Nullable Integer page, @Nullable Integer size)` -> [FindGroups(ctx context.Context, groupIDs []int64)](../internal/domain/group/repository/group_repository.go#L37)
  - [x] `findNotDeletedGroups(Collection<Long> ids, @Nullable Date lastUpdatedDate)` -> [FindNotDeletedGroups(ctx context.Context, groupIDs []int64, lastUpdatedDate *time.Time)](../internal/domain/group/repository/group_repository.go#L200)
  - [x] `findAllNames()` -> [FindAllNames(ctx context.Context)](../internal/domain/user/repository/user_repository.go#L212)
  - [x] `findTypeId(Long groupId)` -> [FindTypeID(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_repository.go#L223)
  - [x] `findTypeIdAndGroupId(Collection<Long> groupIds)` -> [FindTypeIDAndGroupID(ctx context.Context, groupIDs []int64)](../internal/domain/group/repository/group_repository.go#L299)
  - [x] `findTypeIdIfActiveAndNotDeleted(Long groupId)` -> [FindTypeIDIfActiveAndNotDeleted(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_repository.go#L240)
  - [x] `findMinimumScore(Long groupId)` -> [FindMinimumScore(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_repository.go#L261)
  - [ ] `findOwnerId(Long groupId)`
  - [x] `isGroupMuted(Long groupId, Date muteEndDate)` -> [IsGroupMuted(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_repository.go#L278), [IsGroupMuted(ctx context.Context, groupID int64)](../internal/domain/group/service/group_service.go#L336)
  - [x] `isGroupActiveAndNotDeleted(Long groupId)` -> [IsGroupActiveAndNotDeleted(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_repository.go#L288), [IsGroupActiveAndNotDeleted(ctx context.Context, groupID int64)](../internal/domain/group/service/group_service.go#L351)

- **GroupTypeRepository.java** ([java/im/turms/service/domain/group/repository/GroupTypeRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/repository/GroupTypeRepository.java))
> [简述功能]

  - [x] `updateTypes(Set<Long> ids, @Nullable String name, @Nullable Integer groupSizeLimit, @Nullable GroupInvitationStrategy groupInvitationStrategy, @Nullable GroupJoinStrategy groupJoinStrategy, @Nullable GroupUpdateStrategy groupInfoUpdateStrategy, @Nullable GroupUpdateStrategy memberInfoUpdateStrategy, @Nullable Boolean guestSpeakable, @Nullable Boolean selfInfoUpdatable, @Nullable Boolean enableReadReceipt, @Nullable Boolean messageEditable)` -> [UpdateTypes(ctx context.Context, ids []int64, name *string, groupSizeLimit *int32, invitationStrategy *constant.GroupInvitationStrategy, joinStrategy *constant.GroupJoinStrategy, groupInfoUpdateStrategy *constant.GroupUpdateStrategy, memberInfoUpdateStrategy *constant.GroupUpdateStrategy, guestSpeakable *bool, selfInfoUpdatable *bool, enableReadReceipt *bool, messageEditable *bool)](../internal/domain/group/repository/group_type_repository.go#L56)

- **GroupVersionRepository.java** ([java/im/turms/service/domain/group/repository/GroupVersionRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/repository/GroupVersionRepository.java))
> [简述功能]

  - [x] `updateVersions(String field)` -> [UpdateVersions(ctx context.Context, groupIDs []int64, field string)](../internal/domain/group/repository/group_version_repository.go#L140)
  - [x] `updateVersions(@Nullable Set<Long> groupIds, String field)` -> [UpdateVersions(ctx context.Context, groupIDs []int64, field string)](../internal/domain/group/repository/group_version_repository.go#L140)
  - [x] `updateVersion(Long groupId, String field)` -> [UpdateVersion(ctx context.Context, groupID int64, field string)](../internal/domain/group/repository/group_version_repository.go#L48)
  - [x] `updateVersion(Long groupId, boolean updateMembers, boolean updateBlocklist, boolean joinRequests, boolean joinQuestions)` -> [UpdateVersion(ctx context.Context, groupID int64, field string)](../internal/domain/group/repository/group_version_repository.go#L48)
  - [x] `findBlocklist(Long groupId)` -> [FindBlocklist(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_version_repository.go#L151)
  - [x] `findInvitations(Long groupId)` -> [FindInvitations(ctx context.Context, groupID *int64, inviterID *int64, inviteeID *int64, status *po.RequestStatus, lastUpdatedDate *time.Time, page, size int)](../internal/domain/group/repository/group_invitation_repository.go#L135)
  - [x] `findJoinRequests(Long groupId)` -> [FindJoinRequests(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_version_repository.go#L156)
  - [x] `findJoinQuestions(Long groupId)` -> [FindJoinQuestions(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_version_repository.go#L161)
  - [x] `findMembers(Long groupId)` -> [FindMembers(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_version_repository.go#L166)

- **GroupBlocklistService.java** ([java/im/turms/service/domain/group/service/GroupBlocklistService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/service/GroupBlocklistService.java))
> [简述功能]

  - [x] `authAndBlockUser(@NotNull Long requesterId, @NotNull Long groupId, @NotNull Long userIdToBlock, @Nullable ClientSession session)` -> [AuthAndBlockUser(ctx context.Context, requesterID int64, groupID int64, userID int64,)](../internal/domain/group/service/group_blocklist_service.go#L50)
  - [x] `unblockUser(@NotNull Long requesterId, @NotNull Long groupId, @NotNull Long userIdToUnblock, @Nullable ClientSession session, boolean updateBlocklistVersion)` -> [UnblockUser(ctx context.Context, groupID int64, userID int64)](../internal/domain/group/service/group_blocklist_service.go#L116)
  - [x] `findBlockedUserIds(@NotNull Long groupId, @NotNull Set<Long> userIds)` -> [FindBlockedUserIds(ctx context.Context, groupId int64)](../internal/domain/group/repository/group_blocked_user_repository.go#L128)
  - [x] `isBlocked(@NotNull Long groupId, @NotNull Long userId)` -> [IsBlocked(ctx context.Context, groupID int64, userID int64)](../internal/domain/group/service/group_blocklist_service.go#L126)
  - [x] `queryGroupBlockedUserIds(@NotNull Long groupId)` -> [QueryGroupBlockedUserIds(ctx context.Context, groupID int64)](../internal/domain/group/service/group_blocklist_service.go#L220)
  - [x] `queryBlockedUsers(@Nullable Set<Long> groupIds, @Nullable Set<Long> userIds, @Nullable DateRange blockDateRange, @Nullable Set<Long> requesterIds, @Nullable Integer page, @Nullable Integer size)` -> [IsBlocked(ctx context.Context, groupID int64, userID int64)](../internal/domain/group/service/group_blocklist_service.go#L126)
  - [x] `countBlockedUsers(@Nullable Set<Long> groupIds, @Nullable Set<Long> userIds, @Nullable DateRange blockDateRange, @Nullable Set<Long> requesterIds)` -> [CountBlockedUsers(ctx context.Context, groupIds []int64, userIds []int64, blockDateRange *turmsmongo.DateRange, requesterIds []int64,)](../internal/domain/group/repository/group_blocked_user_repository.go#L207), [CountBlockedUsers(ctx context.Context, groupIds []int64, userIds []int64, blockDateRange *turmsmongo.DateRange, requesterIds []int64,)](../internal/domain/group/service/group_blocklist_service.go#L224)
  - [x] `queryGroupBlockedUserIdsWithVersion(@NotNull Long groupId, @Nullable Date lastUpdatedDate)` -> [QueryGroupBlockedUserIdsWithVersion(ctx context.Context, groupID int64, lastUpdatedDate *time.Time,)](../internal/domain/group/service/group_blocklist_service.go#L234)
  - [x] `queryGroupBlockedUserInfosWithVersion(@NotNull Long groupId, @Nullable Date lastUpdatedDate)` -> [QueryGroupBlockedUserInfosWithVersion(ctx context.Context, groupID int64, lastUpdatedDate *time.Time,)](../internal/domain/group/service/group_blocklist_service.go#L258)
  - [x] `addBlockedUser(@NotNull Long groupId, @NotNull Long userId, @NotNull Long requesterId, @Nullable @PastOrPresent Date blockDate)` -> [AddBlockedUser(ctx context.Context, groupID int64, userID int64, requesterID int64, blockDate *time.Time,)](../internal/domain/group/service/group_blocklist_service.go#L282)
  - [x] `updateBlockedUsers(@NotEmpty Set<GroupBlockedUser.@ValidGroupBlockedUserKey Key> keys, @Nullable @PastOrPresent Date blockDate, @Nullable Long requesterId)` -> [UpdateBlockedUsers(ctx context.Context, keys []po.GroupBlockedUserKey, blockDate *time.Time, requesterId *int64)](../internal/domain/group/repository/group_blocked_user_repository.go#L108)
  - [x] `deleteBlockedUsers(@NotEmpty Set<GroupBlockedUser.@ValidGroupBlockedUserKey Key> keys)` -> [DeleteBlockedUsers(ctx context.Context, keys []po.GroupBlockedUserKey)](../internal/domain/group/repository/group_blocked_user_repository.go#L241), [DeleteBlockedUsers(ctx context.Context, keys []po.GroupBlockedUserKey)](../internal/domain/group/service/group_blocklist_service.go#L308)

- **GroupInvitationService.java** ([java/im/turms/service/domain/group/service/GroupInvitationService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/service/GroupInvitationService.java))
> [简述功能]

  - [x] `authAndCreateGroupInvitation(@NotNull Long groupId, @NotNull Long inviterId, @NotNull Long inviteeId, @Nullable String content)` -> [AuthAndCreateGroupInvitation(ctx context.Context, requesterID int64, groupID int64, inviteeID int64, content string,)](../internal/domain/group/service/group_invitation_service.go#L48)
  - [ ] `createGroupInvitation(@Nullable Long id, @NotNull Long groupId, @NotNull Long inviterId, @NotNull Long inviteeId, @Nullable String content, @Nullable @ValidRequestStatus RequestStatus status, @Nullable @PastOrPresent Date creationDate, @Nullable @PastOrPresent Date responseDate)`
  - [x] `queryGroupIdAndInviterIdAndInviteeIdAndStatus(@NotNull Long invitationId)` -> [QueryGroupIdAndInviterIdAndInviteeIdAndStatus(ctx context.Context, invitationID int64)](../internal/domain/group/service/group_invitation_service.go#L305)
  - [x] `queryGroupIdAndInviteeIdAndStatus(@NotNull Long invitationId)` -> [QueryGroupIdAndInviteeIdAndStatus(ctx context.Context, invitationID int64)](../internal/domain/group/service/group_invitation_service.go#L309)
  - [x] `authAndRecallPendingGroupInvitation(@NotNull Long requesterId, @NotNull Long invitationId)` -> [AuthAndRecallPendingGroupInvitation(ctx context.Context, requesterID int64, invitationID int64,)](../internal/domain/group/service/group_invitation_service.go#L144)
  - [x] `queryGroupInvitationsByInviteeId(@NotNull Long inviteeId)` -> [QueryGroupInvitationsByInviteeId(ctx context.Context, inviteeID int64)](../internal/domain/group/service/group_invitation_service.go#L313)
  - [x] `queryGroupInvitationsByInviterId(@NotNull Long inviterId)` -> [QueryGroupInvitationsByInviterId(ctx context.Context, inviterID int64)](../internal/domain/group/service/group_invitation_service.go#L318)
  - [x] `queryGroupInvitationsByGroupId(@NotNull Long groupId)` -> [QueryGroupInvitationsByGroupId(ctx context.Context, groupID int64)](../internal/domain/group/service/group_invitation_service.go#L331)
  - [x] `queryUserGroupInvitationsWithVersion(@NotNull Long userId, boolean areSentByUser, @Nullable Date lastUpdatedDate)` -> [QueryUserGroupInvitationsWithVersion(ctx context.Context, userID int64, areSentInvitations bool, lastUpdatedDate *time.Time)](../internal/domain/group/service/group_invitation_service.go#L238)
  - [x] `authAndQueryGroupInvitationsWithVersion(@NotNull Long userId, @NotNull Long groupId, @Nullable Date lastUpdatedDate)` -> [AuthAndQueryGroupInvitationsWithVersion(ctx context.Context, requesterID int64, groupID int64, lastUpdatedDate *time.Time)](../internal/domain/group/service/group_invitation_service.go#L268)
  - [x] `queryInviteeIdAndGroupIdAndCreationDateAndStatusByInvitationId(@NotNull Long invitationId)` -> [QueryInviteeIdAndGroupIdAndCreationDateAndStatusByInvitationId(ctx context.Context, invitationID int64)](../internal/domain/group/service/group_invitation_service.go#L335)
  - [x] `queryInvitations(@Nullable Set<Long> ids, @Nullable Set<Long> groupIds, @Nullable Set<Long> inviterIds, @Nullable Set<Long> inviteeIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange, @Nullable Integer page, @Nullable Integer size)` -> [QueryInvitations(ctx context.Context, groupID *int64, inviterID *int64, inviteeID *int64, status *po.RequestStatus, lastUpdatedDate *time.Time, page int, size int)](../internal/domain/group/service/group_invitation_service.go#L233)
  - [x] `countInvitations(@Nullable Set<Long> ids, @Nullable Set<Long> groupIds, @Nullable Set<Long> inviterIds, @Nullable Set<Long> inviteeIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange)` -> [CountInvitations(ctx context.Context, groupID, inviterID, inviteeID *int64, status *po.RequestStatus, lastUpdatedDate *time.Time)](../internal/domain/group/repository/group_invitation_repository.go#L239)
  - [x] `deleteInvitations(@Nullable Set<Long> ids)` -> [DeleteInvitations(ctx context.Context, ids []int64)](../internal/domain/group/repository/group_invitation_repository.go#L230)
  - [x] `authAndHandleInvitation(@NotNull Long requesterId, @NotNull Long invitationId, @NotNull @ValidResponseAction ResponseAction action, @Nullable String reason)` -> [AuthAndHandleInvitation(ctx context.Context, requesterID int64, invitationID int64, status po.RequestStatus, reason string,)](../internal/domain/group/service/group_invitation_service.go#L184)
  - [x] `updatePendingInvitationStatus(@NotNull Long groupId, @NotNull Long invitationId, @NotNull @ValidRequestStatus RequestStatus requestStatus, @Nullable String reason, @Nullable ClientSession session)` -> [UpdatePendingInvitationStatus(ctx context.Context, invitationID int64, requestStatus po.RequestStatus, reason *string)](../internal/domain/group/service/group_invitation_service.go#L339)
  - [x] `updateInvitations(@NotEmpty Set<Long> invitationIds, @Nullable Long inviterId, @Nullable Long inviteeId, @Nullable String content, @Nullable @ValidRequestStatus RequestStatus status, @Nullable @PastOrPresent Date creationDate, @Nullable @PastOrPresent Date responseDate)` -> [UpdateInvitations(ctx context.Context, ids []int64, inviterID, inviteeID *int64, content *string, status *po.RequestStatus, creationDate, responseDate *time.Time)](../internal/domain/group/repository/group_invitation_repository.go#L196)

- **GroupJoinRequestService.java** ([java/im/turms/service/domain/group/service/GroupJoinRequestService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/service/GroupJoinRequestService.java))
> [简述功能]

  - [ ] `authAndCreateGroupJoinRequest(@NotNull Long requesterId, @NotNull Long groupId, @Nullable String content)`
  - [ ] `authAndRecallPendingGroupJoinRequest(@NotNull Long requesterId, @NotNull Long requestId)`
  - [x] `authAndQueryGroupJoinRequestsWithVersion(@NotNull Long requesterId, @Nullable Long groupId, @Nullable Date lastUpdatedDate)` -> [AuthAndQueryGroupJoinRequestsWithVersion(ctx context.Context, requesterID int64, groupID int64, lastUpdatedDate *time.Time)](../internal/domain/group/service/group_join_request_service.go#L185)
  - [x] `queryGroupJoinRequestsByGroupId(@NotNull Long groupId)` -> [QueryGroupJoinRequestsByGroupId(ctx context.Context, groupID int64)](../internal/domain/group/service/group_join_request_service.go#L234)
  - [x] `queryGroupJoinRequestsByRequesterId(@NotNull Long requesterId)` -> [QueryGroupJoinRequestsByRequesterId(ctx context.Context, requesterID int64)](../internal/domain/group/service/group_join_request_service.go#L238)
  - [x] `queryGroupId(@NotNull Long requestId)` -> [QueryGroupId(ctx context.Context, requestID int64)](../internal/domain/group/service/group_join_request_service.go#L242)
  - [x] `queryJoinRequests(@Nullable Set<Long> ids, @Nullable Set<Long> groupIds, @Nullable Set<Long> requesterIds, @Nullable Set<Long> responderIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange, @Nullable Integer page, @Nullable Integer size)` -> [QueryJoinRequests(ctx context.Context, groupID *int64, requesterID *int64, responderID *int64, status *po.RequestStatus, creationDate *time.Time, page int, size int)](../internal/domain/group/service/group_join_request_service.go#L180)
  - [x] `countJoinRequests(@Nullable Set<Long> ids, @Nullable Set<Long> groupIds, @Nullable Set<Long> requesterIds, @Nullable Set<Long> responderIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange)` -> [CountJoinRequests(ctx context.Context, ids, groupIds, requesterIds, responderIds []int64, statuses []po.RequestStatus, creationDateRange, responseDateRange, expirationDateRange *turmsmongo.DateRange)](../internal/domain/group/service/group_join_request_service.go#L246)
  - [x] `deleteJoinRequests(@Nullable Set<Long> ids)` -> [DeleteJoinRequests(ctx context.Context, ids []int64)](../internal/domain/group/service/group_join_request_service.go#L250)
  - [x] `authAndHandleJoinRequest(@NotNull Long requesterId, @NotNull Long joinRequestId, @NotNull @ValidResponseAction ResponseAction action, @Nullable String responseReason)` -> [AuthAndHandleJoinRequest(ctx context.Context, responderID int64, requestID int64, status po.RequestStatus, reason string)](../internal/domain/group/service/group_join_request_service.go#L138)
  - [x] `updatePendingJoinRequestStatus(@NotNull Long groupId, @NotNull Long joinRequestId, @NotNull @ValidRequestStatus RequestStatus requestStatus, @NotNull Long responderId, @Nullable String responseReason, @Nullable ClientSession session)` -> [UpdatePendingJoinRequestStatus(ctx context.Context, requestID int64, responderID int64, requestStatus po.RequestStatus, responseReason *string)](../internal/domain/group/service/group_join_request_service.go#L254)
  - [x] `updateJoinRequests(@NotEmpty Set<Long> requestIds, @Nullable Long requesterId, @Nullable Long responderId, @Nullable String content, @Nullable @ValidRequestStatus RequestStatus status, @Nullable @PastOrPresent Date creationDate, @Nullable @PastOrPresent Date responseDate)` -> [UpdateJoinRequests(ctx context.Context, requestIds []int64, requesterId, responderId *int64, content *string, status *po.RequestStatus, creationDate, responseDate *time.Time)](../internal/domain/group/service/group_join_request_service.go#L258)
  - [ ] `createGroupJoinRequest(@Nullable Long id, @NotNull Long groupId, @NotNull Long requesterId, @NotNull Long responderId, @Nullable String content, @Nullable @ValidRequestStatus RequestStatus status, @Nullable @PastOrPresent Date creationDate, @Nullable @PastOrPresent Date responseDate, @Nullable String responseReason)`

- **GroupMemberService.java** ([java/im/turms/service/domain/group/service/GroupMemberService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/service/GroupMemberService.java))
> [简述功能]

  - [x] `addGroupMember(@NotNull Long groupId, @NotNull Long userId, @NotNull @ValidGroupMemberRole GroupMemberRole groupMemberRole, @Nullable String name, @Nullable @PastOrPresent Date joinDate, @Nullable Date muteEndDate, @Nullable ClientSession session)` -> [AddGroupMember(ctx context.Context, member *po.GroupMember)](../internal/domain/group/repository/group_member_repository.go#L34)
  - [x] `addGroupMembers(@NotNull Long groupId, @NotNull Set<Long> userIds, @NotNull @ValidGroupMemberRole GroupMemberRole groupMemberRole, @Nullable String name, @Nullable @PastOrPresent Date joinDate, @Nullable Date muteEndDate, @Nullable ClientSession session)` -> [AddGroupMembers(ctx context.Context, groupID int64, userIDs []int64, role protocol.GroupMemberRole, name *string, joinTime *time.Time, muteEndDate *time.Time, session mongo.SessionContext,)](../internal/domain/group/service/group_member_service.go#L365)
  - [x] `authAndAddGroupMembers(@NotNull Long requesterId, @NotNull Long groupId, @NotNull Set<Long> userIds, @Nullable @ValidGroupMemberRole GroupMemberRole groupMemberRole, @Nullable String name, @Nullable Date muteEndDate, @Nullable ClientSession session)` -> [AuthAndAddGroupMembers(ctx context.Context, requesterID int64, groupID int64, userIDs []int64, role protocol.GroupMemberRole, muteEndDate *time.Time,)](../internal/domain/group/service/group_member_service.go#L405)
  - [x] `authAndDeleteGroupMembers(@NotNull Long requesterId, @NotNull Long groupId, @NotNull Set<Long> memberIdsToDelete, @Nullable Long successorId, @Nullable Boolean quitAfterTransfer)` -> [AuthAndDeleteGroupMembers(ctx context.Context, requesterID int64, groupID int64, userIDs []int64, successorID *int64, quitAfterTransfer bool,)](../internal/domain/group/service/group_member_service.go#L475)
  - [x] `deleteGroupMember(@NotNull Long groupId, @NotNull Long memberId, @Nullable ClientSession session, boolean updateGroupMembersVersion)` -> [DeleteGroupMember(ctx context.Context, groupID, userID int64, session mongo.SessionContext, updateVersion bool,)](../internal/domain/group/service/group_member_service.go#L325)
  - [x] `deleteGroupMembers(@NotEmpty Collection<GroupMember.Key> keys, @Nullable ClientSession session, boolean updateGroupMembersVersion)` -> [DeleteGroupMembers(keys []any, successorId *int64, quitAfterTransfer *bool)](../internal/domain/group/access/admin/controller/group_controllers.go#L86)
  - [ ] `updateGroupMember(@NotNull Long groupId, @NotNull Long memberId, @Nullable String name, @Nullable @ValidGroupMemberRole GroupMemberRole role, @Nullable @PastOrPresent Date joinDate, @Nullable Date muteEndDate, @Nullable ClientSession session, boolean updateGroupMembersVersion)`
  - [x] `updateGroupMembers(@NotEmpty Set<GroupMember.Key> keys, @Nullable String name, @Nullable @ValidGroupMemberRole GroupMemberRole role, @Nullable @PastOrPresent Date joinDate, @Nullable Date muteEndDate, @Nullable ClientSession session, boolean updateGroupMembersVersion)` -> [UpdateGroupMembers(ctx context.Context, keys []po.GroupMemberKey, name *string, role *protocol.GroupMemberRole, joinDate *time.Time, muteEndDate *time.Time)](../internal/domain/group/repository/group_member_repository.go#L163)
  - [x] `updateGroupMembers(@NotNull Long groupId, @NotEmpty Set<Long> memberIds, @Nullable String name, @Nullable @ValidGroupMemberRole GroupMemberRole role, @Nullable @PastOrPresent Date joinDate, @Nullable Date muteEndDate, @Nullable ClientSession session, boolean updateGroupMembersVersion)` -> [UpdateGroupMembers(ctx context.Context, keys []po.GroupMemberKey, name *string, role *protocol.GroupMemberRole, joinDate *time.Time, muteEndDate *time.Time)](../internal/domain/group/repository/group_member_repository.go#L163)
  - [x] `isGroupMember(@NotNull Long groupId, @NotNull Long userId, boolean preferCache)` -> [IsGroupMember(ctx context.Context, groupID, userID int64)](../internal/domain/group/repository/group_member_repository.go#L239)
  - [x] `isGroupMember(@NotEmpty Set<Long> groupIds, @NotNull Long userId)` -> [IsGroupMember(ctx context.Context, groupID, userID int64)](../internal/domain/group/repository/group_member_repository.go#L239)
  - [x] `findExistentMemberGroupIds(@NotEmpty Set<Long> groupIds, @NotNull Long userId)` -> [FindExistentMemberGroupIds(ctx context.Context, groupIDs []int64, userID int64)](../internal/domain/group/repository/group_member_repository.go#L356), [FindExistentMemberGroupIds(ctx context.Context, groupIDs []int64, userID int64)](../internal/domain/group/service/group_member_service.go#L147)
  - [x] `isAllowedToInviteUser(@NotNull Long groupId, @NotNull Long inviterId)` -> [IsAllowedToInviteUser(ctx context.Context, groupID, inviterID int64)](../internal/domain/group/service/group_member_service.go#L153)
  - [x] `isAllowedToBeInvited(@NotNull Long groupId, @NotNull Long inviteeId)` -> [IsAllowedToBeInvited(ctx context.Context, groupID, inviteeID int64)](../internal/domain/group/service/group_member_service.go#L198)
  - [x] `isAllowedToSendMessage(@NotNull Long groupId, @NotNull Long senderId)` -> [IsAllowedToSendMessage(ctx context.Context, groupID, senderID int64)](../internal/domain/group/service/group_member_service.go#L220)
  - [x] `isMemberMuted(@NotNull Long groupId, @NotNull Long userId, boolean preferCache)` -> [IsMemberMuted(ctx context.Context, groupID, userID int64)](../internal/domain/group/repository/group_member_repository.go#L100)
  - [x] `queryGroupMemberKeyAndRolePairs(@NotNull Set<Long> userIds, @NotNull Long groupId)` -> [QueryGroupMemberKeyAndRolePairs(ctx context.Context, userIDs []int64, groupID int64)](../internal/domain/group/service/group_member_service.go#L279)
  - [x] `queryGroupMemberRole(@NotNull Long userId, @NotNull Long groupId, boolean preferCache)` -> [QueryGroupMemberRole(ctx context.Context, groupID, userID int64)](../internal/domain/group/service/group_member_service.go#L137)
  - [x] `isOwner(@NotNull Long userId, @NotNull Long groupId, boolean preferCache)` -> [IsOwner(ctx context.Context, userID, groupID int64)](../internal/domain/group/service/group_member_service.go#L109)
  - [x] `isOwnerOrManager(@NotNull Long userId, @NotNull Long groupId, boolean preferCache)` -> [IsOwnerOrManager(ctx context.Context, groupID, userID int64)](../internal/domain/group/service/group_member_service.go#L121)
  - [x] `isOwnerOrManagerOrMember(@NotNull Long userId, @NotNull Long groupId, boolean preferCache)` -> [IsOwnerOrManagerOrMember(ctx context.Context, userID, groupID int64)](../internal/domain/group/service/group_member_service.go#L285)
  - [x] `queryUserJoinedGroupIds(@NotNull Long userId)` -> [QueryUserJoinedGroupIds(ctx context.Context, userID int64)](../internal/domain/group/service/group_member_service.go#L616)
  - [x] `queryUsersJoinedGroupIds(@Nullable Set<Long> groupIds, @NotEmpty Set<Long> userIds, @Nullable Integer page, @Nullable Integer size)` -> [FindUsersJoinedGroupIds(ctx context.Context, groupIDs []int64, userIDs []int64, page, size *int)](../internal/domain/group/repository/group_member_repository.go#L324), [QueryUsersJoinedGroupIds(ctx context.Context, groupIDs []int64, userIDs []int64, page, size *int)](../internal/domain/group/service/group_member_service.go#L622)
  - [x] `queryMemberIdsInUsersJoinedGroups(@NotEmpty Set<Long> userIds, boolean preferCache)` -> [QueryMemberIdsInUsersJoinedGroups(ctx context.Context, userIDs []int64)](../internal/domain/group/service/group_member_service.go#L628)
  - [ ] `queryGroupMemberIds(@NotNull Long groupId, boolean preferCache)`
  - [x] `queryGroupMemberIds(@NotEmpty Set<Long> groupIds, boolean preferCache)` -> [FindMemberIdsByGroupIds(ctx context.Context, groupIDs []int64)](../internal/domain/group/repository/group_member_repository.go#L383), [QueryGroupMemberIds(ctx context.Context, groupIDs []int64)](../internal/domain/group/service/group_member_service.go#L641)
  - [x] `queryGroupMembers(@Nullable Set<Long> groupIds, @Nullable Set<Long> userIds, @Nullable Set<@ValidGroupMemberRole GroupMemberRole> roles, @Nullable DateRange joinDateRange, @Nullable DateRange muteEndDateRange, @Nullable Integer page, @Nullable Integer size)` -> [QueryGroupMembers(page, size *int)](../internal/domain/group/access/admin/controller/group_controllers.go#L80)
  - [x] `countMembers(@Nullable Set<Long> groupIds, @Nullable Set<Long> userIds, @Nullable Set<@ValidGroupMemberRole GroupMemberRole> roles, @Nullable DateRange joinDateRange, @Nullable DateRange muteEndDateRange)` -> [CountMembers(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_member_repository.go#L208)
  - [x] `deleteGroupMembers(boolean updateGroupMembersVersion)` -> [DeleteGroupMembers(keys []any, successorId *int64, quitAfterTransfer *bool)](../internal/domain/group/access/admin/controller/group_controllers.go#L86)
  - [x] `queryGroupMembers(@NotNull Long groupId, boolean preferCache)` -> [QueryGroupMembers(page, size *int)](../internal/domain/group/access/admin/controller/group_controllers.go#L80)
  - [x] `queryGroupMembers(@NotNull Long groupId, @NotEmpty Set<Long> memberIds, boolean preferCache)` -> [QueryGroupMembers(page, size *int)](../internal/domain/group/access/admin/controller/group_controllers.go#L80)
  - [x] `authAndQueryGroupMembers(@NotNull Long requesterId, @NotNull Long groupId, @NotEmpty Set<Long> memberIds, boolean withStatus)` -> [AuthAndQueryGroupMembers(ctx context.Context, requesterID int64, groupID int64, memberIDs []int64, withStatus bool)](../internal/domain/group/service/group_member_service.go#L647)
  - [ ] `authAndQueryGroupMembersWithVersion(@NotNull Long requesterId, @NotNull Long groupId, @Nullable Date lastUpdatedDate, boolean withStatus)`
  - [x] `authAndUpdateGroupMember(@NotNull Long requesterId, @NotNull Long groupId, @NotNull Long memberId, @Nullable String name, @Nullable @ValidGroupMemberRole GroupMemberRole role, @Nullable Date muteEndDate)` -> [AuthAndUpdateGroupMember(ctx context.Context, requesterID int64, groupID int64, memberID int64, name *string, role *protocol.GroupMemberRole, muteEndDate *time.Time,)](../internal/domain/group/service/group_member_service.go#L563)
  - [x] `deleteAllGroupMembers(@Nullable Set<Long> groupIds, @Nullable ClientSession session, boolean updateMembersVersion)` -> [DeleteAllGroupMembers(ctx context.Context, groupIDs []int64, session mongo.SessionContext, updateVersion bool)](../internal/domain/group/service/group_member_service.go#L348)
  - [x] `queryGroupManagersAndOwnerId(@NotNull Long groupId)` -> [FindGroupManagersAndOwnerId(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_member_repository.go#L252), [QueryGroupManagersAndOwnerId(ctx context.Context, groupID int64)](../internal/domain/group/service/group_member_service.go#L675)

- **GroupQuestionService.java** ([java/im/turms/service/domain/group/service/GroupQuestionService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/service/GroupQuestionService.java))
> [简述功能]

  - [ ] `checkGroupQuestionAnswerAndGetScore(@NotNull Long questionId, @NotNull String answer, @Nullable Long groupId)`
  - [ ] `authAndCheckGroupQuestionAnswerAndJoin(@NotNull Long requesterId, @NotNull @ValidGroupQuestionIdAndAnswer Map<Long, String> questionIdToAnswer)`
  - [ ] `authAndCreateGroupJoinQuestions(@NotNull Long requesterId, @NotNull Long groupId, @NotNull List<NewGroupQuestion> questions)`
  - [ ] `createGroupJoinQuestions(@NotNull Long groupId, @NotNull List<NewGroupQuestion> questions)`
  - [x] `queryGroupId(@NotNull Long questionId)` -> [QueryGroupId(ctx context.Context, requestID int64)](../internal/domain/group/service/group_join_request_service.go#L242)
  - [ ] `authAndDeleteGroupJoinQuestions(@NotNull Long requesterId, @NotNull Long groupId, @NotNull Set<Long> questionIds)`
  - [x] `queryGroupJoinQuestions(@Nullable Set<Long> ids, @Nullable Set<Long> groupIds, @Nullable Integer page, @Nullable Integer size, boolean withAnswers)` -> [QueryGroupJoinQuestions(page, size *int)](../internal/domain/group/access/admin/controller/group_controllers.go#L93)
  - [ ] `countGroupJoinQuestions(@Nullable Set<Long> ids, @Nullable Set<Long> groupIds)`
  - [x] `deleteGroupJoinQuestions(@Nullable Set<Long> ids)` -> [DeleteGroupJoinQuestions(ids []int64)](../internal/domain/group/access/admin/controller/group_controllers.go#L105)
  - [x] `authAndQueryGroupJoinQuestionsWithVersion(@NotNull Long requesterId, @NotNull Long groupId, boolean withAnswers, @Nullable Date lastUpdatedDate)` -> [AuthAndQueryGroupJoinQuestionsWithVersion(ctx context.Context, requesterID int64, groupID int64, withAnswers bool, lastUpdatedDate *time.Time)](../internal/domain/group/service/group_question_service.go#L212)
  - [ ] `authAndUpdateGroupJoinQuestion(@NotNull Long requesterId, @NotNull Long questionId, @Nullable String question, @Nullable Set<String> answers, @Nullable @Min(0)`
  - [x] `updateGroupJoinQuestions(@NotEmpty Set<Long> ids, @Nullable Long groupId, @Nullable String question, @Nullable Set<String> answers, @Nullable @Min(0)` -> [UpdateGroupJoinQuestions(ids []int64, updateGroupJoinQuestionDTO any)](../internal/domain/group/access/admin/controller/group_controllers.go#L102)

- **GroupService.java** ([java/im/turms/service/domain/group/service/GroupService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/service/GroupService.java))
> [简述功能]

  - [x] `createGroup(@NotNull Long creatorId, @NotNull Long ownerId, @Nullable String groupName, @Nullable String intro, @Nullable String announcement, @Nullable @Min(value = 0)` -> [CreateGroup(ctx context.Context, creatorID, groupID int64, name, intro *string, minimumScore *int32)](../internal/domain/group/service/group_service.go#L45)
  - [x] `authAndDeleteGroup(boolean queryGroupMemberIds, @NotNull Long requesterId, @NotNull Long groupId)` -> [AuthAndDeleteGroup(ctx context.Context, requesterID int64, groupID int64)](../internal/domain/group/service/group_service.go#L184)
  - [ ] `authAndCreateGroup(@NotNull Long creatorId, @NotNull Long ownerId, @Nullable String groupName, @Nullable String intro, @Nullable String announcement, @Nullable @Min(value = 0)`
  - [x] `deleteGroupsAndGroupMembers(@Nullable Set<Long> groupIds, @Nullable Boolean deleteLogically)` -> [DeleteGroupsAndGroupMembers(ctx context.Context, groupIDs []int64, session mongo.SessionContext)](../internal/domain/group/service/group_service.go#L202)
  - [x] `queryGroups(@Nullable Set<Long> ids, @Nullable Set<Long> typeIds, @Nullable Set<Long> creatorIds, @Nullable Set<Long> ownerIds, @Nullable Boolean isActive, @Nullable DateRange creationDateRange, @Nullable DateRange deletionDateRange, @Nullable DateRange lastUpdatedDateRange, @Nullable DateRange muteEndDateRange, @Nullable Set<Long> memberIds, @Nullable Integer page, @Nullable Integer size)` -> [QueryGroups(page, size *int)](../internal/domain/group/access/admin/controller/group_controllers.go#L28), [QueryGroups(ctx context.Context, groupIDs []int64, name *string, lastUpdatedDate *time.Time, skip *int32, limit *int32)](../internal/domain/group/repository/group_repository.go#L57)
  - [ ] `queryGroupTypeIfActiveAndNotDeleted(@NotNull Long groupId)`
  - [ ] `queryGroupTypeIfActiveAndNotDeleted(@NotNull Long groupId, boolean preferCache)`
  - [ ] `queryGroupTypeId(@NotNull Long groupId)`
  - [x] `queryGroupTypeIdIfActiveAndNotDeleted(@NotNull Long groupId)` -> [QueryGroupTypeIdIfActiveAndNotDeleted(ctx context.Context, groupID int64)](../internal/domain/group/service/group_service.go#L103)
  - [ ] `queryGroupMinimumScore(@NotNull Long groupId)`
  - [x] `authAndTransferGroupOwnership(@NotNull Long requesterId, @NotNull Long groupId, @NotNull Long successorId, boolean quitAfterTransfer, @Nullable ClientSession session)` -> [AuthAndTransferGroupOwnership(ctx context.Context, requesterID, groupID, successorID int64, quitAfterTransfer bool, session mongo.SessionContext,)](../internal/domain/group/service/group_service.go#L131)
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
  - [x] `countOwnedGroups(@NotNull Long ownerId)` -> [CountOwnedGroups(ctx context.Context, ownerID int64)](../internal/domain/group/repository/group_repository.go#L132)
  - [x] `countOwnedGroups(@NotNull Long ownerId, @NotNull Long groupTypeId)` -> [CountOwnedGroups(ctx context.Context, ownerID int64)](../internal/domain/group/repository/group_repository.go#L132)
  - [x] `countCreatedGroups(@Nullable DateRange dateRange)` -> [CountCreatedGroups(ctx context.Context, dateRange *turmsmongo.DateRange)](../internal/domain/group/repository/group_repository.go#L164)
  - [x] `countGroups(@Nullable Set<Long> ids, @Nullable Set<Long> typeIds, @Nullable Set<Long> creatorIds, @Nullable Set<Long> ownerIds, @Nullable Boolean isActive, @Nullable DateRange creationDateRange, @Nullable DateRange deletionDateRange, @Nullable DateRange lastUpdatedDateRange, @Nullable DateRange muteEndDateRange, @Nullable Set<Long> memberIds)` -> [CountGroups(ctx context.Context, ownerIDs []int64, relatedUserIDs []int64)](../internal/domain/user/repository/user_relationship_group_member_repository.go#L195)
  - [x] `countDeletedGroups(@Nullable DateRange dateRange)` -> [CountDeletedGroups(ctx context.Context, dateRange *turmsmongo.DateRange)](../internal/domain/group/repository/group_repository.go#L182)
  - [x] `count()` -> [Count(ctx context.Context, filter bson.M)](../internal/domain/user/repository/user_repository.go#L104)
  - [x] `isGroupMuted(@NotNull Long groupId)` -> [IsGroupMuted(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_repository.go#L278), [IsGroupMuted(ctx context.Context, groupID int64)](../internal/domain/group/service/group_service.go#L336)
  - [x] `isGroupActiveAndNotDeleted(@NotNull Long groupId)` -> [IsGroupActiveAndNotDeleted(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_repository.go#L288), [IsGroupActiveAndNotDeleted(ctx context.Context, groupID int64)](../internal/domain/group/service/group_service.go#L351)

- **GroupTypeService.java** ([java/im/turms/service/domain/group/service/GroupTypeService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/service/GroupTypeService.java))
> [简述功能]

  - [ ] `initGroupTypes()`
  - [x] `queryGroupTypes(@Nullable Integer page, @Nullable Integer size)` -> [QueryGroupTypes(page, size *int)](../internal/domain/group/access/admin/controller/group_controllers.go#L113)
  - [x] `addGroupType(@Nullable Long id, @NotNull @NoWhitespace String name, @NotNull @Min(1)` -> [AddGroupType(addGroupTypeDTO any)](../internal/domain/group/access/admin/controller/group_controllers.go#L112)
  - [ ] `updateGroupTypes(@NotEmpty Set<Long> ids, @Nullable @NoWhitespace String name, @Nullable @Min(1)`
  - [ ] `deleteGroupTypes(@Nullable Set<Long> groupTypeIds)`
  - [ ] `queryGroupType(@NotNull Long groupTypeId)`
  - [x] `queryGroupTypes(@NotNull Collection<Long> groupTypeIds)` -> [QueryGroupTypes(page, size *int)](../internal/domain/group/access/admin/controller/group_controllers.go#L113)
  - [ ] `groupTypeExists(@NotNull Long groupTypeId)`
  - [ ] `countGroupTypes()`

- **GroupVersionService.java** ([java/im/turms/service/domain/group/service/GroupVersionService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/group/service/GroupVersionService.java))
> [简述功能]

  - [x] `queryMembersVersion(@NotNull Long groupId)` -> [QueryGroupMembersVersion(ctx context.Context, groupID int64)](../internal/domain/group/service/group_version_service.go#L76)
  - [x] `queryBlocklistVersion(@NotNull Long groupId)` -> [QueryGroupBlocklistVersion(ctx context.Context, groupID int64)](../internal/domain/group/service/group_version_service.go#L85)
  - [x] `queryGroupJoinRequestsVersion(@NotNull Long groupId)` -> [QueryGroupJoinRequestsVersion(ctx context.Context, groupID int64)](../internal/domain/group/service/group_version_service.go#L58), [QueryGroupJoinRequestsVersion(ctx context.Context, userID int64)](../internal/domain/user/service/user_version_service.go#L53)
  - [x] `queryGroupJoinQuestionsVersion(@NotNull Long groupId)` -> [QueryGroupJoinQuestionsVersion(ctx context.Context, groupID int64)](../internal/domain/group/service/group_version_service.go#L67)
  - [x] `queryGroupInvitationsVersion(@NotNull Long groupId)` -> [QueryGroupInvitationsVersion(ctx context.Context, groupID int64)](../internal/domain/group/service/group_version_service.go#L49)
  - [x] `updateVersion(@NotNull Long groupId, boolean updateMembers, boolean updateBlocklist, boolean joinRequests, boolean joinQuestions)` -> [UpdateVersion(ctx context.Context, groupID int64, field string)](../internal/domain/group/repository/group_version_repository.go#L48)
  - [x] `updateMembersVersion(@NotNull Long groupId)` -> [UpdateMembersVersion(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_version_repository.go#L61)
  - [x] `updateMembersVersion(@Nullable Set<Long> groupIds)` -> [UpdateMembersVersion(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_version_repository.go#L61)
  - [x] `updateMembersVersion()` -> [UpdateMembersVersion(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_version_repository.go#L61)
  - [x] `updateBlocklistVersion(@NotNull Long groupId)` -> [UpdateBlocklistVersion(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_version_repository.go#L73)
  - [x] `updateJoinRequestsVersion(@NotNull Long groupId)` -> [UpdateJoinRequestsVersion(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_version_repository.go#L79)
  - [x] `updateJoinQuestionsVersion(@NotNull Long groupId)` -> [UpdateJoinQuestionsVersion(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_version_repository.go#L85)
  - [ ] `updateGroupInvitationsVersion(@NotNull Long groupId)`
  - [x] `updateSpecificVersion(@NotNull Long groupId, @NotNull String field)` -> [UpdateSpecificVersion(ctx context.Context, userIDs []int64, field string, updateDate time.Time)](../internal/domain/user/repository/user_version_repository.go#L104)
  - [x] `updateSpecificVersion(@NotNull String field)` -> [UpdateSpecificVersion(ctx context.Context, userIDs []int64, field string, updateDate time.Time)](../internal/domain/user/repository/user_version_repository.go#L104)
  - [x] `updateSpecificVersion(@Nullable Set<Long> groupIds, @NotNull String field)` -> [UpdateSpecificVersion(ctx context.Context, userIDs []int64, field string, updateDate time.Time)](../internal/domain/user/repository/user_version_repository.go#L104)
  - [x] `upsert(@NotNull Long groupId, @NotNull Date timestamp)` -> [Upsert(ctx context.Context, groupID int64, timestamp time.Time)](../internal/domain/group/repository/group_version_repository.go#L99)
  - [x] `delete(@Nullable Set<Long> groupIds, @Nullable ClientSession session)` -> [Delete(key K)](../internal/domain/common/cache/sharded_map.go#L63)

- **MessageController.java** ([java/im/turms/service/domain/message/access/admin/controller/MessageController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/message/access/admin/controller/MessageController.java))
> [简述功能]

  - [x] `createMessages(@QueryParam(defaultValue = "true")` -> [CreateMessages()](../internal/domain/message/access/admin/controller/message_controller.go#L9)
  - [x] `queryMessages(@QueryParam(required = false)` -> [QueryMessages()](../internal/domain/message/access/admin/controller/message_controller.go#L14), [QueryMessages(ctx context.Context, isGroupMessage *bool, senderIDs []int64, targetIDs []int64, deliveryDateAfter *time.Time, deliveryDateBefore *time.Time, size int64, ascending bool,)](../internal/domain/message/repository/message_repository.go#L71)
  - [x] `queryMessages(@QueryParam(required = false)` -> [QueryMessages()](../internal/domain/message/access/admin/controller/message_controller.go#L14), [QueryMessages(ctx context.Context, isGroupMessage *bool, senderIDs []int64, targetIDs []int64, deliveryDateAfter *time.Time, deliveryDateBefore *time.Time, size int64, ascending bool,)](../internal/domain/message/repository/message_repository.go#L71)
  - [x] `countMessages(@QueryParam(required = false)` -> [CountMessages()](../internal/domain/message/access/admin/controller/message_controller.go#L19)
  - [x] `updateMessages(Set<Long> ids, @RequestBody UpdateMessageDTO updateMessageDTO)` -> [UpdateMessages()](../internal/domain/message/access/admin/controller/message_controller.go#L24)
  - [x] `deleteMessages(Set<Long> ids, @QueryParam(required = false)` -> [DeleteMessages()](../internal/domain/message/access/admin/controller/message_controller.go#L29)

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

  - [x] `handleCreateMessageRequest()` -> [HandleCreateMessageRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/message/controller/message_controller.go#L28)
  - [x] `handleQueryMessagesRequest()` -> [HandleQueryMessagesRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/message/controller/message_controller.go#L116)
  - [x] `handleUpdateMessageRequest()` -> [HandleUpdateMessageRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/message/controller/message_controller.go#L250)
  - [x] `handleCreateMessageReactionsRequest()` -> [HandleCreateMessageReactionsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/message/controller/message_controller.go#L272)
  - [x] `handleDeleteMessageReactionsRequest()` -> [HandleDeleteMessageReactionsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/message/controller/message_controller.go#L278)

- **MessageAndRecipientIds.java** ([java/im/turms/service/domain/message/bo/MessageAndRecipientIds.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/message/bo/MessageAndRecipientIds.java))
> [简述功能]

  - [ ] `MessageAndRecipientIds(Message message, Set<Long> recipientIds)`

- **Message.java** ([java/im/turms/service/domain/message/po/Message.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/message/po/Message.java))
> [简述功能]

  - [ ] `groupId()`

- **MessageRepository.java** ([java/im/turms/service/domain/message/repository/MessageRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/message/repository/MessageRepository.java))
> [简述功能]

  - [x] `updateMessages(Set<Long> messageIds, @Nullable Boolean isSystemMessage, @Nullable Integer senderIp, @Nullable byte[] senderIpV6, @Nullable Date recallDate, @Nullable String text, @Nullable List<byte[]> records, @Nullable Integer burnAfter, @Nullable ClientSession session)` -> [UpdateMessages(ctx context.Context, messageIDs []int64, isSystemMessage *bool, senderIP *int32, senderIPv6 []byte, recallDate *time.Time, text *string, records [][]byte, burnAfter *int32,)](../internal/domain/message/repository/message_repository.go#L304), [UpdateMessages(ctx context.Context, senderID *int64, senderDeviceType *int32, messageIDs []int64, isSystemMessage *bool, text *string, records [][]byte, burnAfter *int32,)](../internal/domain/message/service/message_service.go#L417)
  - [x] `updateMessagesDeletionDate(@Nullable Set<Long> messageIds)` -> [UpdateMessagesDeletionDate(ctx context.Context, messageIDs []int64, deletionDate *time.Time)](../internal/domain/message/repository/message_repository.go#L357)
  - [x] `existsBySenderIdAndTargetId(Long senderId, Long targetId)` -> [ExistsBySenderIDAndTargetID(ctx context.Context, senderID int64, targetID int64)](../internal/domain/message/repository/message_repository.go#L374)
  - [x] `countMessages(@Nullable Set<Long> messageIds, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages, @Nullable Set<Long> senderIds, @Nullable Set<Long> targetIds, @Nullable DateRange deliveryDateRange, @Nullable DateRange deletionDateRange)` -> [CountMessages(ctx context.Context, isGroupMessage *bool, senderIDs []int64, targetIDs []int64, deliveryDateAfter *time.Time, deliveryDateBefore *time.Time,)](../internal/domain/message/repository/message_repository.go#L128), [CountMessages(ctx context.Context, isGroupMessage *bool, senderIDs []int64, targetIDs []int64, deliveryDateAfter *time.Time, deliveryDateBefore *time.Time,)](../internal/domain/message/service/message_service.go#L249)
  - [x] `countUsersWhoSentMessage(@Nullable DateRange dateRange, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages)` -> [CountUsersWhoSentMessage(ctx context.Context, deliveryDateAfter *time.Time, deliveryDateBefore *time.Time, areGroupMessages *bool, areSystemMessages *bool)](../internal/domain/message/repository/message_repository.go#L507), [CountUsersWhoSentMessage(ctx context.Context, deliveryDateAfter *time.Time, deliveryDateBefore *time.Time, areGroupMessages *bool, areSystemMessages *bool)](../internal/domain/message/service/message_service.go#L474)
  - [x] `countGroupsThatSentMessages(@Nullable DateRange dateRange)` -> [CountGroupsThatSentMessages(ctx context.Context, deliveryDateAfter *time.Time, deliveryDateBefore *time.Time)](../internal/domain/message/repository/message_repository.go#L534), [CountGroupsThatSentMessages(ctx context.Context, deliveryDateAfter *time.Time, deliveryDateBefore *time.Time)](../internal/domain/message/service/message_service.go#L479)
  - [x] `countSentMessages(@Nullable DateRange dateRange, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages)` -> [CountSentMessages(ctx context.Context, deliveryDateAfter *time.Time, deliveryDateBefore *time.Time, areGroupMessages *bool, areSystemMessages *bool)](../internal/domain/message/repository/message_repository.go#L555), [CountSentMessages(ctx context.Context, deliveryDateAfter *time.Time, deliveryDateBefore *time.Time, areGroupMessages *bool, areSystemMessages *bool)](../internal/domain/message/service/message_service.go#L484)
  - [x] `findDeliveryDate(Long messageId)` -> [FindDeliveryDate(ctx context.Context, messageID int64)](../internal/domain/message/repository/message_repository.go#L387)
  - [x] `findExpiredMessageIds(Date expirationDate)` -> [FindExpiredMessageIds(ctx context.Context, expirationDate time.Time)](../internal/domain/message/repository/message_repository.go#L400)
  - [x] `findMessageGroupId(Long messageId)` -> [FindMessageGroupId(ctx context.Context, messageID int64)](../internal/domain/message/repository/message_repository.go#L421)
  - [x] `findMessageSenderIdAndTargetIdAndIsGroupMessage(Long messageId)` -> [FindMessageSenderIDAndTargetIDAndIsGroupMessage(ctx context.Context, messageID int64)](../internal/domain/message/repository/message_repository.go#L434)
  - [x] `findMessages(@Nullable Collection<Long> messageIds, @Nullable Collection<byte[]> conversationIds, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages, @Nullable Set<Long> senderIds, @Nullable Set<Long> targetIds, @Nullable DateRange deliveryDateRange, @Nullable DateRange deletionDateRange, @Nullable DateRange recallDateRange, @Nullable Integer page, @Nullable Integer size, @Nullable Boolean ascending)` -> [FindMessages(ctx context.Context, messageIDs []int64, conversationIDs [][]byte, areGroupMessages *bool, areSystemMessages *bool, senderIDs []int64, targetIDs []int64, deliveryDateRange *turmsmongo.DateRange, deletionDateRange *turmsmongo.DateRange, recallDateRange *turmsmongo.DateRange, page *int, size *int, ascending *bool,)](../internal/domain/message/repository/message_repository.go#L167)
  - [x] `findIsGroupMessageAndTargetId(Long messageId, Long senderId)` -> [FindIsGroupMessageAndTargetID(ctx context.Context, messageID int64, senderID int64)](../internal/domain/message/repository/message_repository.go#L445)
  - [x] `findIsGroupMessageAndTargetIdAndDeliveryDate(Long messageId, Long senderId)` -> [FindIsGroupMessageAndTargetIDAndDeliveryDate(ctx context.Context, messageID int64, senderID int64)](../internal/domain/message/repository/message_repository.go#L456)
  - [x] `getGroupConversationId(long groupId)` -> [GetGroupConversationID(groupID int64)](../internal/domain/message/repository/message_repository.go#L467)
  - [x] `getPrivateConversationId(long id1, long id2)` -> [GetPrivateConversationID(id1 int64, id2 int64)](../internal/domain/message/repository/message_repository.go#L477)

- **MessageService.java** ([java/im/turms/service/domain/message/service/MessageService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/message/service/MessageService.java))
> [简述功能]

  - [x] `isMessageRecipientOrSender(@NotNull Long messageId, @NotNull Long userId)` -> [IsMessageRecipientOrSender(ctx context.Context, messageID int64, userID int64)](../internal/domain/message/service/message_service.go#L309)
  - [x] `authAndQueryCompleteMessages(Long requesterId, @Nullable Collection<Long> messageIds, @NotNull Boolean areGroupMessages, @Nullable Boolean areSystemMessages, @Nullable Set<Long> fromIds, @Nullable DateRange deliveryDateRange, @Nullable Integer maxCount, boolean ascending, boolean withTotal)` -> [AuthAndQueryCompleteMessages(ctx context.Context, requesterID int64, isGroupMessage *bool, areSystemMessages *bool, senderIDs []int64, deliveryDateAfter *time.Time, deliveryDateBefore *time.Time, size int64, ascending bool,)](../internal/domain/message/service/message_service.go#L508)
  - [x] `queryMessage(@NotNull Long messageId)` -> [QueryMessage(ctx context.Context, messageID int64)](../internal/domain/message/service/message_service.go#L327)
  - [x] `queryMessages(@Nullable Collection<Long> messageIds, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages, @Nullable Set<Long> senderIds, @Nullable Set<Long> targetIds, @Nullable DateRange deliveryDateRange, @Nullable DateRange deletionDateRange, @Nullable DateRange recallDateRange, @Nullable Integer page, @Nullable Integer size, @Nullable Boolean ascending)` -> [QueryMessages(ctx context.Context, isGroupMessage *bool, senderIDs []int64, targetIDs []int64, deliveryDateAfter *time.Time, deliveryDateBefore *time.Time, size int64, ascending bool,)](../internal/domain/message/repository/message_repository.go#L71)
  - [x] `saveMessage(@Nullable Long messageId, @NotNull Long senderId, @Nullable byte[] senderIp, @NotNull Long targetId, @NotNull Boolean isGroupMessage, @NotNull Boolean isSystemMessage, @Nullable String text, @Nullable List<byte[]> records, @Nullable @Min(0)` -> [SaveMessage(ctx context.Context, isGroupMessage bool, senderID int64, targetID int64, text string, records [][]byte, burnAfter *int32, deliveryDate *time.Time, preMessageID *int64,)](../internal/domain/message/service/message_service.go#L332)
  - [x] `queryExpiredMessageIds(@NotNull Integer retentionPeriodHours)` -> [QueryExpiredMessageIds(ctx context.Context, retentionPeriodHours int)](../internal/domain/message/service/message_service.go#L389)
  - [x] `deleteExpiredMessages(@NotNull Integer retentionPeriodHours)` -> [DeleteExpiredMessages(ctx context.Context, retentionPeriodHours int)](../internal/domain/message/service/message_service.go#L395)
  - [x] `deleteMessages(@Nullable Set<Long> messageIds, @Nullable Boolean deleteLogically)` -> [DeleteMessages(ctx context.Context, messageIDs []int64)](../internal/domain/message/repository/message_repository.go#L497), [DeleteMessages(ctx context.Context, messageIDs []int64, deleteLogically *bool)](../internal/domain/message/service/message_service.go#L405)
  - [x] `updateMessages(@Nullable Long senderId, @Nullable DeviceType senderDeviceType, @NotEmpty Set<Long> messageIds, @Nullable Boolean isSystemMessage, @Nullable String text, @Nullable List<byte[]> records, @Nullable @Min(0)` -> [UpdateMessages(ctx context.Context, messageIDs []int64, isSystemMessage *bool, senderIP *int32, senderIPv6 []byte, recallDate *time.Time, text *string, records [][]byte, burnAfter *int32,)](../internal/domain/message/repository/message_repository.go#L304), [UpdateMessages(ctx context.Context, senderID *int64, senderDeviceType *int32, messageIDs []int64, isSystemMessage *bool, text *string, records [][]byte, burnAfter *int32,)](../internal/domain/message/service/message_service.go#L417)
  - [x] `hasPrivateMessage(Long senderId, Long targetId)` -> [HasPrivateMessage(ctx context.Context, senderID int64, targetID int64)](../internal/domain/message/service/message_service.go#L431)
  - [x] `countMessages(@Nullable Set<Long> messageIds, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages, @Nullable Set<Long> senderIds, @Nullable Set<Long> targetIds, @Nullable DateRange deliveryDateRange, @Nullable DateRange deletionDateRange)` -> [CountMessages(ctx context.Context, isGroupMessage *bool, senderIDs []int64, targetIDs []int64, deliveryDateAfter *time.Time, deliveryDateBefore *time.Time,)](../internal/domain/message/repository/message_repository.go#L128), [CountMessages(ctx context.Context, isGroupMessage *bool, senderIDs []int64, targetIDs []int64, deliveryDateAfter *time.Time, deliveryDateBefore *time.Time,)](../internal/domain/message/service/message_service.go#L249)
  - [x] `countUsersWhoSentMessage(@Nullable DateRange dateRange, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages)` -> [CountUsersWhoSentMessage(ctx context.Context, deliveryDateAfter *time.Time, deliveryDateBefore *time.Time, areGroupMessages *bool, areSystemMessages *bool)](../internal/domain/message/repository/message_repository.go#L507), [CountUsersWhoSentMessage(ctx context.Context, deliveryDateAfter *time.Time, deliveryDateBefore *time.Time, areGroupMessages *bool, areSystemMessages *bool)](../internal/domain/message/service/message_service.go#L474)
  - [x] `countGroupsThatSentMessages(@Nullable DateRange dateRange)` -> [CountGroupsThatSentMessages(ctx context.Context, deliveryDateAfter *time.Time, deliveryDateBefore *time.Time)](../internal/domain/message/repository/message_repository.go#L534), [CountGroupsThatSentMessages(ctx context.Context, deliveryDateAfter *time.Time, deliveryDateBefore *time.Time)](../internal/domain/message/service/message_service.go#L479)
  - [x] `countSentMessages(@Nullable DateRange dateRange, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages)` -> [CountSentMessages(ctx context.Context, deliveryDateAfter *time.Time, deliveryDateBefore *time.Time, areGroupMessages *bool, areSystemMessages *bool)](../internal/domain/message/repository/message_repository.go#L555), [CountSentMessages(ctx context.Context, deliveryDateAfter *time.Time, deliveryDateBefore *time.Time, areGroupMessages *bool, areSystemMessages *bool)](../internal/domain/message/service/message_service.go#L484)
  - [x] `countSentMessagesOnAverage(@Nullable DateRange dateRange, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages)` -> [CountSentMessagesOnAverage(ctx context.Context, deliveryDateAfter *time.Time, deliveryDateBefore *time.Time, areGroupMessages *bool, areSystemMessages *bool)](../internal/domain/message/service/message_service.go#L489)
  - [x] `authAndUpdateMessage(@NotNull Long senderId, @Nullable DeviceType senderDeviceType, @NotNull Long messageId, @Nullable String text, @Nullable List<byte[]> records, @Nullable @PastOrPresent Date recallDate)` -> [AuthAndUpdateMessage(ctx context.Context, senderID int64, senderDeviceType *int32, messageID int64, text *string, records [][]byte, recallDate *time.Time,)](../internal/domain/message/service/message_service.go#L436)
  - [x] `queryMessageRecipients(@NotNull Long messageId)` -> [QueryMessageRecipients(ctx context.Context, messageID int64)](../internal/domain/message/service/message_service.go#L524)
  - [x] `authAndSaveMessage(boolean queryRecipientIds, @Nullable Boolean persist, @Nullable Long messageId, @NotNull Long senderId, @Nullable byte[] senderIp, @NotNull Long targetId, @NotNull Boolean isGroupMessage, @NotNull Boolean isSystemMessage, @Nullable String text, @Nullable List<byte[]> records, @Nullable @Min(0)` -> [AuthAndSaveMessage(ctx context.Context, isGroupMessage bool, senderID int64, targetID int64, text string, records [][]byte, burnAfter *int32, deliveryDate *time.Time, preMessageID *int64,)](../internal/domain/message/service/message_service.go#L73)
  - [x] `saveMessage(boolean queryRecipientIds, @Nullable Boolean persist, @Nullable Long messageId, @NotNull Long senderId, @Nullable byte[] senderIp, @NotNull Long targetId, @NotNull Boolean isGroupMessage, @NotNull Boolean isSystemMessage, @Nullable String text, @Nullable List<byte[]> records, @Nullable @Min(0)` -> [SaveMessage(ctx context.Context, isGroupMessage bool, senderID int64, targetID int64, text string, records [][]byte, burnAfter *int32, deliveryDate *time.Time, preMessageID *int64,)](../internal/domain/message/service/message_service.go#L332)
  - [x] `authAndCloneAndSaveMessage(boolean queryRecipientIds, @NotNull Long requesterId, @Nullable byte[] requesterIp, @NotNull Long referenceId, @NotNull Boolean isGroupMessage, @NotNull Boolean isSystemMessage, @NotNull Long targetId)` -> [AuthAndCloneAndSaveMessage(ctx context.Context, requesterID int64, referenceID int64, isGroupMessage bool, isSystemMessage bool, targetID int64,)](../internal/domain/message/service/message_service.go#L610)
  - [x] `cloneAndSaveMessage(boolean queryRecipientIds, @NotNull Long senderId, @Nullable byte[] senderIp, @NotNull Long referenceId, @NotNull Boolean isGroupMessage, @NotNull Boolean isSystemMessage, @NotNull Long targetId)` -> [CloneAndSaveMessage(ctx context.Context, senderID int64, referenceID int64, isGroupMessage bool, isSystemMessage bool, targetID int64,)](../internal/domain/message/service/message_service.go#L583)
  - [x] `authAndSaveAndSendMessage(boolean send, @Nullable Boolean persist, @Nullable Long senderId, @Nullable DeviceType senderDeviceType, @Nullable byte[] senderIp, @Nullable Long messageId, @NotNull Boolean isGroupMessage, @NotNull Boolean isSystemMessage, @Nullable String text, @Nullable List<byte[]> records, @NotNull Long targetId, @Nullable @Min(0)` -> [AuthAndSaveAndSendMessage(ctx context.Context, isGroupMessage bool, senderID int64, targetID int64, text string, records [][]byte, burnAfter *int32, deliveryDate *time.Time, preMessageID *int64,)](../internal/domain/message/service/message_service.go#L141)
  - [x] `saveAndSendMessage(boolean send, @Nullable Boolean persist, @Nullable Long senderId, @Nullable DeviceType senderDeviceType, @Nullable byte[] senderIp, @Nullable Long messageId, @NotNull Boolean isGroupMessage, @NotNull Boolean isSystemMessage, @Nullable String text, @Nullable List<byte[]> records, @NotNull Long targetId, @Nullable @Min(0)` -> [SaveAndSendMessage(ctx context.Context, send bool, persist bool, senderID int64, isGroupMessage bool, isSystemMessage bool, // Note: not persisted separately yet text string, records [][]byte, targetID int64, burnAfter *int32, deliveryDate *time.Time, preMessageID *int64,)](../internal/domain/message/service/message_service.go#L539)
  - [x] `saveAndSendMessage(@Nullable Boolean persist, @Nullable Long senderId, @Nullable DeviceType senderDeviceType, @Nullable byte[] senderIp, @Nullable Long messageId, @NotNull Boolean isGroupMessage, @NotNull Boolean isSystemMessage, @Nullable String text, @Nullable List<byte[]> records, @NotNull Long targetId, @Nullable @Min(0)` -> [SaveAndSendMessage(ctx context.Context, send bool, persist bool, senderID int64, isGroupMessage bool, isSystemMessage bool, // Note: not persisted separately yet text string, records [][]byte, targetID int64, burnAfter *int32, deliveryDate *time.Time, preMessageID *int64,)](../internal/domain/message/service/message_service.go#L539)
  - [x] `deleteGroupMessageSequenceIds(Set<Long> groupIds)` -> [DeleteGroupMessageSequenceIDs(ctx context.Context, groupIDs []int64)](../internal/domain/message/service/message_service.go#L631)
  - [x] `deletePrivateMessageSequenceIds(Set<Long> userIds)` -> [DeletePrivateMessageSequenceIDs(ctx context.Context, userIDs []int64)](../internal/domain/message/service/message_service.go#L638)
  - [x] `fetchGroupMessageSequenceId(Long groupId)` -> [FetchGroupMessageSequenceID(ctx context.Context, groupID int64)](../internal/domain/message/service/message_service.go#L643)
  - [x] `fetchPrivateMessageSequenceId(Long userId1, Long userId2)` -> [FetchPrivateMessageSequenceID(ctx context.Context, userID1 int64, userID2 int64)](../internal/domain/message/service/message_service.go#L648)

- **StatisticsService.java** ([java/im/turms/service/domain/observation/service/StatisticsService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/observation/service/StatisticsService.java))
> [简述功能]

  - [ ] `countOnlineUsersByNodes()`
  - [x] `countOnlineUsers()` -> [CountOnlineUsers()](../internal/domain/gateway/session/service.go#L189), [CountOnlineUsers()](../internal/domain/gateway/session/sharded_map.go#L144)

- **StorageServiceController.java** ([java/im/turms/service/domain/storage/access/servicerequest/controller/StorageServiceController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/storage/access/servicerequest/controller/StorageServiceController.java))
> [简述功能]

  - [x] `handleDeleteResourceRequest()` -> [HandleDeleteResourceRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/storage/controller/storage_controller.go#L28)
  - [x] `handleQueryResourceUploadInfoRequest()` -> [HandleQueryResourceUploadInfoRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/storage/controller/storage_controller.go#L52)
  - [x] `handleQueryResourceDownloadInfoRequest()` -> [HandleQueryResourceDownloadInfoRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/storage/controller/storage_controller.go#L84)
  - [x] `handleUpdateMessageAttachmentInfoRequest()` -> [HandleUpdateMessageAttachmentInfoRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/storage/controller/storage_controller.go#L114)
  - [x] `handleQueryMessageAttachmentInfosRequest()` -> [HandleQueryMessageAttachmentInfosRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/storage/controller/storage_controller.go#L140)

- **StorageResourceInfo.java** ([java/im/turms/service/domain/storage/bo/StorageResourceInfo.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/storage/bo/StorageResourceInfo.java))
> [简述功能]

  - [ ] `StorageResourceInfo(@Nullable Long idNum, @Nullable String idStr, String name, String mediaType, Long uploaderId, Date creationDate)`

- **StorageService.java** ([java/im/turms/service/domain/storage/service/StorageService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/storage/service/StorageService.java))
> [简述功能]

  - [x] `deleteResource(Long requesterId, StorageResourceType resourceType, @Nullable Long resourceIdNum, @Nullable String resourceIdStr, List<Value> customAttributes)` -> [DeleteResource(ctx context.Context, resourceType constants.StorageResourceType, keyStr string)](../internal/domain/storage/provider/mock_storage_provider.go#L27)
  - [x] `queryResourceUploadInfo(Long requesterId, StorageResourceType resourceType, @Nullable Long resourceIdNum, @Nullable String resourceName, @Nullable String resourceMediaType, List<Value> customAttributes)` -> [QueryResourceUploadInfo(ctx context.Context, requesterID int64, resourceType constants.StorageResourceType, resourceName string, contentType string, maxSize int64,)](../internal/domain/storage/service/storage_service.go#L38)
  - [x] `queryResourceDownloadInfo(Long requesterId, StorageResourceType resourceType, @Nullable Long resourceIdNum, @Nullable String resourceIdStr, List<Value> customAttributes)` -> [QueryResourceDownloadInfo(ctx context.Context, requesterID int64, resourceType constants.StorageResourceType, resourceIDStr string,)](../internal/domain/storage/service/storage_service.go#L54)
  - [x] `shareMessageAttachmentWithUser(Long requesterId, @Nullable Long messageAttachmentIdNum, @Nullable String messageAttachmentIdStr, Long userIdToShareWith)` -> [ShareMessageAttachmentWithUser(ctx context.Context, requesterID int64, messageAttachmentIDNum *int64, messageAttachmentIDStr *string, userIDToShareWith int64)](../internal/domain/storage/service/storage_service.go#L68)
  - [x] `shareMessageAttachmentWithGroup(Long requesterId, @Nullable Long messageAttachmentIdNum, @Nullable String messageAttachmentIdStr, Long groupIdToShareWith)` -> [ShareMessageAttachmentWithGroup(ctx context.Context, requesterID int64, messageAttachmentIDNum *int64, messageAttachmentIDStr *string, groupIDToShareWith int64)](../internal/domain/storage/service/storage_service.go#L73)
  - [x] `unshareMessageAttachmentWithUser(Long requesterId, @Nullable Long messageAttachmentIdNum, @Nullable String messageAttachmentIdStr, Long userIdToUnshareWith)` -> [UnshareMessageAttachmentWithUser(ctx context.Context, requesterID int64, messageAttachmentIDNum *int64, messageAttachmentIDStr *string, userIDToUnshareWith int64)](../internal/domain/storage/service/storage_service.go#L78)
  - [x] `unshareMessageAttachmentWithGroup(Long requesterId, @Nullable Long messageAttachmentIdNum, @Nullable String messageAttachmentIdStr, Long groupIdToUnshareWith)` -> [UnshareMessageAttachmentWithGroup(ctx context.Context, requesterID int64, messageAttachmentIDNum *int64, messageAttachmentIDStr *string, groupIDToUnshareWith int64)](../internal/domain/storage/service/storage_service.go#L83)
  - [x] `queryMessageAttachmentInfosUploadedByRequester(Long requesterId, @Nullable DateRange creationDateRange)` -> [QueryMessageAttachmentInfosUploadedByRequester(ctx context.Context, requesterID int64, creationDateStart *time.Time, creationDateEnd *time.Time)](../internal/domain/storage/service/storage_service.go#L88)
  - [x] `queryMessageAttachmentInfosInPrivateConversations(Long requesterId, @Nullable Set<Long> userIds, @Nullable DateRange creationDateRange, @Nullable Boolean areSharedByRequester)` -> [QueryMessageAttachmentInfosInPrivateConversations(ctx context.Context, requesterID int64, userIDs []int64, creationDateStart *time.Time, creationDateEnd *time.Time, areSharedByRequester *bool)](../internal/domain/storage/service/storage_service.go#L93)
  - [x] `queryMessageAttachmentInfosInGroupConversations(Long requesterId, @Nullable Set<Long> groupIds, @Nullable Set<Long> userIds, @Nullable DateRange creationDateRange)` -> [QueryMessageAttachmentInfosInGroupConversations(ctx context.Context, requesterID int64, groupIDs []int64, userIDs []int64, creationDateStart *time.Time, creationDateEnd *time.Time)](../internal/domain/storage/service/storage_service.go#L98)

- **UserController.java** ([java/im/turms/service/domain/user/access/admin/controller/UserController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/access/admin/controller/UserController.java))
> [简述功能]

  - [x] `addUser(@RequestBody AddUserDTO addUserDTO)` -> [AddUser()](../internal/domain/user/access/admin/controller/user_controllers.go#L9), [AddUser(ctx context.Context, id int64, password string, name string, intro string, profilePicture string, profileAccess int32, permissionGroupID int64, registrationDate time.Time, isActive bool)](../internal/domain/user/service/user_service.go#L71)
  - [x] `queryUsers(@QueryParam(required = false)` -> [QueryUsers()](../internal/domain/user/access/admin/controller/user_controllers.go#L14), [QueryUsers(ctx context.Context, userIDs []int64)](../internal/domain/user/service/user_service.go#L149)
  - [x] `queryUsers(@QueryParam(required = false)` -> [QueryUsers()](../internal/domain/user/access/admin/controller/user_controllers.go#L14), [QueryUsers(ctx context.Context, userIDs []int64)](../internal/domain/user/service/user_service.go#L149)
  - [x] `countUsers(@QueryParam(required = false)` -> [CountUsers()](../internal/domain/user/access/admin/controller/user_controllers.go#L19), [CountUsers(ctx context.Context, startDate *time.Time, endDate *time.Time)](../internal/domain/user/repository/user_repository.go#L183)
  - [x] `updateUser(Set<Long> ids, @RequestBody UpdateUserDTO updateUserDTO)` -> [UpdateUser()](../internal/domain/user/access/admin/controller/user_controllers.go#L24), [UpdateUser(ctx context.Context, userID int64, update bson.M)](../internal/domain/user/service/user_service.go#L103)
  - [x] `deleteUsers(Set<Long> ids, @QueryParam(required = false)` -> [DeleteUsers()](../internal/domain/user/access/admin/controller/user_controllers.go#L29), [DeleteUsers(ctx context.Context, userIDs []int64)](../internal/domain/user/service/user_service.go#L117)

- **UserOnlineInfoController.java** ([java/im/turms/service/domain/user/access/admin/controller/UserOnlineInfoController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/access/admin/controller/UserOnlineInfoController.java))
> [简述功能]

  - [x] `countOnlineUsers(boolean countByNodes)` -> [CountOnlineUsers()](../internal/domain/user/access/admin/controller/user_controllers.go#L39)
  - [x] `queryUserSessions(Set<Long> ids, boolean returnNonExistingUsers)` -> [QueryUserSessions()](../internal/domain/user/access/admin/controller/user_controllers.go#L44), [QueryUserSessions(ctx context.Context, userIDs []int64)](../internal/domain/user/service/onlineuser/session_service.go#L66)
  - [x] `queryUserStatuses(Set<Long> ids, boolean returnNonExistingUsers)` -> [QueryUserStatuses()](../internal/domain/user/access/admin/controller/user_controllers.go#L49)
  - [x] `queryUsersNearby(Long userId, @QueryParam(required = false)` -> [QueryUsersNearby()](../internal/domain/user/access/admin/controller/user_controllers.go#L54)
  - [x] `queryUserLocations(Set<Long> ids, @QueryParam(required = false)` -> [QueryUserLocations()](../internal/domain/user/access/admin/controller/user_controllers.go#L59)
  - [x] `updateUserOnlineStatus(Set<Long> ids, @QueryParam(required = false)` -> [UpdateUserOnlineStatus()](../internal/domain/user/access/admin/controller/user_controllers.go#L64)

- **UserRoleController.java** ([java/im/turms/service/domain/user/access/admin/controller/UserRoleController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/access/admin/controller/UserRoleController.java))
> [简述功能]

  - [x] `addUserRole(@RequestBody AddUserRoleDTO addUserRoleDTO)` -> [AddUserRole()](../internal/domain/user/access/admin/controller/user_controllers.go#L74), [AddUserRole(ctx context.Context, role *po.UserRole)](../internal/domain/user/service/user_role_service.go#L29)
  - [x] `queryUserRoles(@QueryParam(required = false)` -> [QueryUserRoles()](../internal/domain/user/access/admin/controller/user_controllers.go#L79), [QueryUserRoles(ctx context.Context, filter bson.M)](../internal/domain/user/service/user_role_service.go#L23)
  - [x] `queryUserRoleGroups(int page, @QueryParam(required = false)` -> [QueryUserRoleGroups()](../internal/domain/user/access/admin/controller/user_controllers.go#L84)
  - [x] `updateUserRole(Set<Long> ids, @RequestBody UpdateUserRoleDTO updateUserRoleDTO)` -> [UpdateUserRole()](../internal/domain/user/access/admin/controller/user_controllers.go#L89)
  - [x] `deleteUserRole(Set<Long> ids)` -> [DeleteUserRole()](../internal/domain/user/access/admin/controller/user_controllers.go#L94)

- **UserFriendRequestController.java** ([java/im/turms/service/domain/user/access/admin/controller/relationship/UserFriendRequestController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/access/admin/controller/relationship/UserFriendRequestController.java))
> [简述功能]

  - [x] `createFriendRequest(@RequestBody AddFriendRequestDTO addFriendRequestDTO)` -> [CreateFriendRequest()](../internal/domain/user/access/admin/controller/user_controllers.go#L104), [CreateFriendRequest(ctx context.Context, requestID *int64, requesterID, recipientID int64, content string, status *po.RequestStatus, creationDate, responseDate *time.Time, reason *string)](../internal/domain/user/service/user_friend_request_service.go#L80)
  - [x] `queryFriendRequests(@QueryParam(required = false)` -> [QueryFriendRequests()](../internal/domain/user/access/admin/controller/user_controllers.go#L109), [QueryFriendRequests(ctx context.Context, ids, requesterIds, recipientIds []int64, statuses []po.RequestStatus, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd *time.Time, page, size *int)](../internal/domain/user/service/user_friend_request_service.go#L405)
  - [x] `queryFriendRequests(@QueryParam(required = false)` -> [QueryFriendRequests()](../internal/domain/user/access/admin/controller/user_controllers.go#L109), [QueryFriendRequests(ctx context.Context, ids, requesterIds, recipientIds []int64, statuses []po.RequestStatus, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd *time.Time, page, size *int)](../internal/domain/user/service/user_friend_request_service.go#L405)
  - [x] `updateFriendRequests(Set<Long> ids, @RequestBody UpdateFriendRequestDTO updateFriendRequestDTO)` -> [UpdateFriendRequests()](../internal/domain/user/access/admin/controller/user_controllers.go#L114), [UpdateFriendRequests(ctx context.Context, requestIds []int64, requesterID, recipientID *int64, content *string, status *po.RequestStatus, reason *string, creationDate *time.Time)](../internal/domain/user/repository/user_friend_request_repository.go#L106)
  - [x] `deleteFriendRequests(@QueryParam(required = false)` -> [DeleteFriendRequests()](../internal/domain/user/access/admin/controller/user_controllers.go#L119), [DeleteFriendRequests(ctx context.Context, ids []int64)](../internal/domain/user/service/user_friend_request_service.go#L399)

- **UserRelationshipController.java** ([java/im/turms/service/domain/user/access/admin/controller/relationship/UserRelationshipController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/access/admin/controller/relationship/UserRelationshipController.java))
> [简述功能]

  - [x] `addRelationship(@RequestBody AddRelationshipDTO addRelationshipDTO)` -> [AddRelationship()](../internal/domain/user/access/admin/controller/user_controllers.go#L129)
  - [x] `queryRelationships(@QueryParam(required = false)` -> [QueryRelationships()](../internal/domain/user/access/admin/controller/user_controllers.go#L134), [QueryRelationships(ctx context.Context, ownerIDs []int64, relatedUserIDs []int64, groupIndexes []int32, isBlocked *bool, establishmentDateRange *turmsmongo.DateRange, page *int, size *int,)](../internal/domain/user/service/user_relationship_service.go#L345)
  - [x] `queryRelationships(@QueryParam(required = false)` -> [QueryRelationships()](../internal/domain/user/access/admin/controller/user_controllers.go#L134), [QueryRelationships(ctx context.Context, ownerIDs []int64, relatedUserIDs []int64, groupIndexes []int32, isBlocked *bool, establishmentDateRange *turmsmongo.DateRange, page *int, size *int,)](../internal/domain/user/service/user_relationship_service.go#L345)
  - [x] `updateRelationships(List<UserRelationship.Key> keys, @RequestBody UpdateRelationshipDTO updateRelationshipDTO)` -> [UpdateRelationships()](../internal/domain/user/access/admin/controller/user_controllers.go#L139)
  - [x] `deleteRelationships(List<UserRelationship.Key> keys)` -> [DeleteRelationships()](../internal/domain/user/access/admin/controller/user_controllers.go#L144)

- **UserRelationshipGroupController.java** ([java/im/turms/service/domain/user/access/admin/controller/relationship/UserRelationshipGroupController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/access/admin/controller/relationship/UserRelationshipGroupController.java))
> [简述功能]

  - [x] `addRelationshipGroup(@RequestBody AddRelationshipGroupDTO addRelationshipGroupDTO)` -> [AddRelationshipGroup()](../internal/domain/user/access/admin/controller/user_controllers.go#L154)
  - [x] `deleteRelationshipGroups(@QueryParam(required = false)` -> [DeleteRelationshipGroups()](../internal/domain/user/access/admin/controller/user_controllers.go#L159), [DeleteRelationshipGroups(ctx context.Context, ownerID int64, groupIndexes []int32, session *mongo.Session)](../internal/domain/user/repository/user_relationship_group_repository.go#L104)
  - [x] `updateRelationshipGroups(List<UserRelationshipGroup.Key> keys, @RequestBody UpdateRelationshipGroupDTO updateRelationshipGroupDTO)` -> [UpdateRelationshipGroups()](../internal/domain/user/access/admin/controller/user_controllers.go#L164), [UpdateRelationshipGroups(ctx context.Context, keys []po.UserRelationshipGroupKey, newName string, session *mongo.Session)](../internal/domain/user/repository/user_relationship_group_repository.go#L130)
  - [x] `queryRelationshipGroups(@QueryParam(required = false)` -> [QueryRelationshipGroups()](../internal/domain/user/access/admin/controller/user_controllers.go#L169), [QueryRelationshipGroups(ctx context.Context, ownerIDs []int64, groupIndexes []int32, page *int, size *int,)](../internal/domain/user/service/user_relationship_group_service.go#L562)
  - [x] `queryRelationshipGroups(@QueryParam(required = false)` -> [QueryRelationshipGroups()](../internal/domain/user/access/admin/controller/user_controllers.go#L169), [QueryRelationshipGroups(ctx context.Context, ownerIDs []int64, groupIndexes []int32, page *int, size *int,)](../internal/domain/user/service/user_relationship_group_service.go#L562)

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
  - [x] `toString()` -> [ToString()](../internal/domain/gateway/session/connection.go#L107), [ToString()](../internal/infra/proto/proto_parser.go#L38)

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
  - [x] `toString()` -> [ToString()](../internal/domain/gateway/session/connection.go#L107), [ToString()](../internal/infra/proto/proto_parser.go#L38)

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

  - [x] `handleCreateFriendRequestRequest()` -> [HandleCreateFriendRequestRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_relationship_controller.go#L52)
  - [x] `handleCreateRelationshipGroupRequest()` -> [HandleCreateRelationshipGroupRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_relationship_controller.go#L62)
  - [x] `handleCreateRelationshipRequest()` -> [HandleCreateRelationshipRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_relationship_controller.go#L82)
  - [x] `handleDeleteFriendRequestRequest()` -> [HandleDeleteFriendRequestRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_relationship_controller.go#L100)
  - [x] `handleDeleteRelationshipGroupRequest()` -> [HandleDeleteRelationshipGroupRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_relationship_controller.go#L110)
  - [x] `handleDeleteRelationshipRequest()` -> [HandleDeleteRelationshipRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_relationship_controller.go#L127)
  - [x] `handleQueryFriendRequestsRequest()` -> [HandleQueryFriendRequestsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_relationship_controller.go#L150)
  - [x] `handleQueryRelatedUserIdsRequest()` -> [HandleQueryRelatedUserIdsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_relationship_controller.go#L183)
  - [x] `handleQueryRelationshipGroupsRequest()` -> [HandleQueryRelationshipGroupsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_relationship_controller.go#L218)
  - [x] `handleQueryRelationshipsRequest()` -> [HandleQueryRelationshipsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_relationship_controller.go#L257)
  - [x] `handleUpdateFriendRequestRequest()` -> [HandleUpdateFriendRequestRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_relationship_controller.go#L297)
  - [x] `handleUpdateRelationshipGroupRequest()` -> [HandleUpdateRelationshipGroupRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_relationship_controller.go#L310)
  - [x] `handleUpdateRelationshipRequest()` -> [HandleUpdateRelationshipRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_relationship_controller.go#L322)

- **UserServiceController.java** ([java/im/turms/service/domain/user/access/servicerequest/controller/UserServiceController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/access/servicerequest/controller/UserServiceController.java))
> [简述功能]

  - [x] `handleQueryUserProfilesRequest()` -> [HandleQueryUserProfilesRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_service_controller.go#L60)
  - [x] `handleQueryNearbyUsersRequest()` -> [HandleQueryNearbyUsersRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_service_controller.go#L113)
  - [x] `handleQueryUserOnlineStatusesRequest()` -> [HandleQueryUserOnlineStatusesRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_service_controller.go#L180)
  - [x] `handleUpdateUserLocationRequest()` -> [HandleUpdateUserLocationRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_service_controller.go#L217)
  - [x] `handleUpdateUserOnlineStatusRequest()` -> [HandleUpdateUserOnlineStatusRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_service_controller.go#L228)
  - [x] `handleUpdateUserRequest()` -> [HandleUpdateUserRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_service_controller.go#L265)

- **UserSettingsServiceController.java** ([java/im/turms/service/domain/user/access/servicerequest/controller/UserSettingsServiceController.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/access/servicerequest/controller/UserSettingsServiceController.java))
> [简述功能]

  - [ ] `handleDeleteUserSettingsRequest()`
  - [x] `handleUpdateUserSettingsRequest()` -> [HandleUpdateUserSettingsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_settings_controller.go#L32)
  - [x] `handleQueryUserSettingsRequest()` -> [HandleQueryUserSettingsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest)](../internal/domain/user/controller/user_settings_controller.go#L39)

- **HandleFriendRequestResult.java** ([java/im/turms/service/domain/user/bo/HandleFriendRequestResult.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/bo/HandleFriendRequestResult.java))
> [简述功能]

  - [ ] `HandleFriendRequestResult(UserFriendRequest friendRequest, @Nullable Integer newGroupIndexForFriendRequestRequester, @Nullable Integer newGroupIndexForFriendRequestRecipient)`

- **UpsertRelationshipResult.java** ([java/im/turms/service/domain/user/bo/UpsertRelationshipResult.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/bo/UpsertRelationshipResult.java))
> [简述功能]

  - [ ] `UpsertRelationshipResult(boolean createdNewRelationship, @Nullable Integer newRelationshipGroupIndex)`

- **UserFriendRequestRepository.java** ([java/im/turms/service/domain/user/repository/UserFriendRequestRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/repository/UserFriendRequestRepository.java))
> [简述功能]

  - [x] `getEntityExpireAfterSeconds()` -> [GetEntityExpireAfterSeconds()](../internal/domain/user/repository/user_friend_request_repository.go#L135)
  - [x] `updateFriendRequests(Set<Long> requestIds, @Nullable Long requesterId, @Nullable Long recipientId, @Nullable String content, @Nullable RequestStatus status, @Nullable String reason, @Nullable Date creationDate)` -> [UpdateFriendRequests(ctx context.Context, requestIds []int64, requesterID, recipientID *int64, content *string, status *po.RequestStatus, reason *string, creationDate *time.Time)](../internal/domain/user/repository/user_friend_request_repository.go#L106)
  - [x] `updateStatusIfPending(Long requestId, RequestStatus requestStatus, @Nullable String reason, @Nullable ClientSession session)` -> [UpdateStatusIfPending(ctx context.Context, invitationID int64, newStatus po.RequestStatus, reason *string, responseDate time.Time)](../internal/domain/group/repository/group_invitation_repository.go#L64)
  - [x] `countFriendRequests(@Nullable Set<Long> ids, @Nullable Set<Long> requesterIds, @Nullable Set<Long> recipientIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange)` -> [CountFriendRequests(ctx context.Context, ids, requesterIds, recipientIds []int64, statuses []po.RequestStatus, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd *time.Time)](../internal/domain/user/repository/user_friend_request_repository.go#L292)
  - [x] `findFriendRequests(@Nullable Set<Long> ids, @Nullable Set<Long> requesterIds, @Nullable Set<Long> recipientIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange, @Nullable Integer page, @Nullable Integer size)` -> [FindFriendRequests(ctx context.Context, ids, requesterIds, recipientIds []int64, statuses []po.RequestStatus, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd *time.Time, page, size *int)](../internal/domain/user/repository/user_friend_request_repository.go#L268)
  - [x] `findFriendRequestsByRecipientId(Long recipientId)` -> [FindFriendRequestsByRecipientId(ctx context.Context, recipientID int64)](../internal/domain/user/repository/user_friend_request_repository.go#L140)
  - [x] `findFriendRequestsByRequesterId(Long requesterId)` -> [FindFriendRequestsByRequesterId(ctx context.Context, requesterID int64)](../internal/domain/user/repository/user_friend_request_repository.go#L157)
  - [x] `findRecipientId(Long requestId)` -> [FindRecipientId(ctx context.Context, requestID int64)](../internal/domain/user/repository/user_friend_request_repository.go#L174)
  - [x] `findRequesterIdAndRecipientIdAndStatus(Long requestId)` -> [FindRequesterIdAndRecipientIdAndStatus(ctx context.Context, requestID int64)](../internal/domain/user/repository/user_friend_request_repository.go#L185)
  - [x] `findRequesterIdAndRecipientIdAndCreationDateAndStatus(Long requestId)` -> [FindRequesterIdAndRecipientIdAndCreationDateAndStatus(ctx context.Context, requestID int64)](../internal/domain/user/repository/user_friend_request_repository.go#L200)
  - [x] `hasPendingFriendRequest(Long requesterId, Long recipientId)` -> [HasPendingFriendRequest(ctx context.Context, requesterID, recipientID int64)](../internal/domain/user/repository/user_friend_request_repository.go#L50)
  - [x] `hasPendingOrDeclinedOrIgnoredOrExpiredRequest(Long requesterId, Long recipientId)` -> [HasPendingOrDeclinedOrIgnoredOrExpiredRequest(ctx context.Context, requesterID, recipientID int64)](../internal/domain/user/repository/user_friend_request_repository.go#L64)

- **UserRelationshipGroupMemberRepository.java** ([java/im/turms/service/domain/user/repository/UserRelationshipGroupMemberRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/repository/UserRelationshipGroupMemberRepository.java))
> [简述功能]

  - [x] `deleteAllRelatedUserFromRelationshipGroup(Long ownerId, Integer groupIndex)` -> [DeleteAllRelatedUserFromRelationshipGroup(ctx context.Context, ownerID int64, session *mongo.Session)](../internal/domain/user/repository/user_relationship_group_member_repository.go#L117)
  - [x] `deleteRelatedUserFromRelationshipGroup(Long ownerId, Long relatedUserId, Integer groupIndex, @Nullable ClientSession session)` -> [DeleteRelatedUserFromRelationshipGroup(ctx context.Context, ownerID int64, relatedUserID int64, groupIndexes []int32, session *mongo.Session)](../internal/domain/user/repository/user_relationship_group_member_repository.go#L139)
  - [x] `deleteRelatedUsersFromAllRelationshipGroups(Long ownerId, Collection<Long> relatedUserIds, @Nullable ClientSession session)` -> [DeleteRelatedUsersFromAllRelationshipGroups(ctx context.Context, ownerID int64, relatedUserIDs []int64, session *mongo.Session)](../internal/domain/user/repository/user_relationship_group_member_repository.go#L167)
  - [x] `countGroups(@Nullable Collection<Long> ownerIds, @Nullable Collection<Long> relatedUserIds)` -> [CountGroups(ctx context.Context, ownerIDs []int64, relatedUserIDs []int64)](../internal/domain/user/repository/user_relationship_group_member_repository.go#L195)
  - [x] `countMembers(@Nullable Set<Long> ownerIds, @Nullable Set<Integer> groupIndexes)` -> [CountMembers(ctx context.Context, groupID int64)](../internal/domain/group/repository/group_member_repository.go#L208)
  - [x] `findGroupIndexes(Long ownerId, Long relatedUserId)` -> [FindGroupIndexes(ctx context.Context, ownerID int64, relatedUserID int64)](../internal/domain/user/repository/user_relationship_group_member_repository.go#L218)
  - [x] `findRelationshipGroupMemberIds(Long ownerId, Integer groupIndex)` -> [FindRelationshipGroupMemberIds(ctx context.Context, ownerIDs []int64, groupIndexes []int32, page *int, size *int)](../internal/domain/user/repository/user_relationship_group_member_repository.go#L242)
  - [x] `findRelationshipGroupMemberIds(@Nullable Set<Long> ownerIds, @Nullable Set<Integer> groupIndexes, @Nullable Integer page, @Nullable Integer size)` -> [FindRelationshipGroupMemberIds(ctx context.Context, ownerIDs []int64, groupIndexes []int32, page *int, size *int)](../internal/domain/user/repository/user_relationship_group_member_repository.go#L242)
  - [x] `findRelationshipGroupMembers(Long ownerId, Integer groupIndex)` -> [FindRelationshipGroupMembers(ctx context.Context, ownerID int64, groupIndex int32)](../internal/domain/user/repository/user_relationship_group_member_repository.go#L276)

- **UserRelationshipGroupRepository.java** ([java/im/turms/service/domain/user/repository/UserRelationshipGroupRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/repository/UserRelationshipGroupRepository.java))
> [简述功能]

  - [x] `deleteAllRelationshipGroups(Set<Long> ownerIds, @Nullable ClientSession session)` -> [DeleteAllRelationshipGroups(ctx context.Context, ownerIDs []int64, session *mongo.Session)](../internal/domain/user/repository/user_relationship_group_repository.go#L78)
  - [x] `updateRelationshipGroupName(Long ownerId, Integer groupIndex, String newGroupName)` -> [UpdateRelationshipGroupName(ctx context.Context, ownerID int64, groupIndex int32, newName string, session *mongo.Session)](../internal/domain/user/repository/user_relationship_group_repository.go#L163)
  - [x] `updateRelationshipGroups(Set<UserRelationshipGroup.Key> keys, @Nullable String name, @Nullable Date creationDate)` -> [UpdateRelationshipGroups(ctx context.Context, keys []po.UserRelationshipGroupKey, newName string, session *mongo.Session)](../internal/domain/user/repository/user_relationship_group_repository.go#L130)
  - [x] `countRelationshipGroups(@Nullable Set<Long> ownerIds, @Nullable Set<Integer> indexes, @Nullable Set<String> names, @Nullable DateRange creationDateRange)` -> [CountRelationshipGroups(ctx context.Context, ownerIDs []int64, groupIndexes []int32)](../internal/domain/user/repository/user_relationship_group_repository.go#L189)
  - [x] `findRelationshipGroups(@Nullable Set<Long> ownerIds, @Nullable Set<Integer> indexes, @Nullable Set<String> names, @Nullable DateRange creationDateRange, @Nullable Integer page, @Nullable Integer size)` -> [FindRelationshipGroups(ctx context.Context, ownerIDs []int64, groupIndexes []int32, page *int, size *int)](../internal/domain/user/repository/user_relationship_group_repository.go#L202)
  - [x] `findRelationshipGroupsInfos(Long ownerId)` -> [FindRelationshipGroupsInfos(ctx context.Context, ownerID int64)](../internal/domain/user/repository/user_relationship_group_repository.go#L231)

- **UserRelationshipRepository.java** ([java/im/turms/service/domain/user/repository/UserRelationshipRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/repository/UserRelationshipRepository.java))
> [简述功能]

  - [x] `deleteAllRelationships(Set<Long> userIds, @Nullable ClientSession session)` -> [DeleteAllRelationships(ctx context.Context, ownerIDs []int64, session *mongo.Session)](../internal/domain/user/repository/user_relationship_repository.go#L222)
  - [x] `updateUserOneSidedRelationships(Set<UserRelationship.Key> keys, @Nullable String name, @Nullable Date blockDate, @Nullable Date establishmentDate)` -> [UpdateUserOneSidedRelationships(ctx context.Context, ownerID int64, relatedUserIDs []int64, blockDate *time.Time, establishmentDate *time.Time, name *string, session *mongo.Session,)](../internal/domain/user/repository/user_relationship_repository.go#L245)
  - [x] `countRelationships(@Nullable Set<Long> ownerIds, @Nullable Set<Long> relatedUserIds, @Nullable Boolean isBlocked)` -> [CountRelationships(ctx context.Context, ownerIDs, relatedUserIDs []int64, groupIndexes []int32, isBlocked *bool, session *mongo.Session)](../internal/domain/user/repository/user_relationship_repository.go#L308), [CountRelationships(ctx context.Context, ownerIDs []int64, relatedUserIDs []int64, groupIndexes []int32, isBlocked *bool,)](../internal/domain/user/service/user_relationship_service.go#L412)
  - [x] `findRelatedUserIds(@Nullable Set<Long> ownerIds, @Nullable Boolean isBlocked)` -> [FindRelatedUserIDs(ctx context.Context, ownerIDs []int64, groupIndexes []int32, isBlocked *bool, page, size *int, session *mongo.Session)](../internal/domain/user/repository/user_relationship_repository.go#L102)
  - [x] `findRelationships(@Nullable Set<Long> ownerIds, @Nullable Set<Long> relatedUserIds, @Nullable Boolean isBlocked, @Nullable DateRange establishmentDateRange, @Nullable Integer page, @Nullable Integer size)` -> [FindRelationships(ctx context.Context, ownerIDs []int64, relatedUserIDs []int64, groupIndexes []int32, isBlocked *bool, establishmentDateRange *turmsmongo.DateRange, page, size *int, session *mongo.Session,)](../internal/domain/user/repository/user_relationship_repository.go#L181)
  - [x] `findRelationships(@Nullable Set<Long> ownerIds, @Nullable Set<Long> relatedUserIds, @Nullable Integer page, @Nullable Integer size)` -> [FindRelationships(ctx context.Context, ownerIDs []int64, relatedUserIDs []int64, groupIndexes []int32, isBlocked *bool, establishmentDateRange *turmsmongo.DateRange, page, size *int, session *mongo.Session,)](../internal/domain/user/repository/user_relationship_repository.go#L181)
  - [x] `hasRelationshipAndNotBlocked(Long ownerId, Long relatedUserId)` -> [HasRelationshipAndNotBlocked(ctx context.Context, ownerID, relatedUserID int64, session *mongo.Session)](../internal/domain/user/repository/user_relationship_repository.go#L48)
  - [x] `isBlocked(Long ownerId, Long relatedUserId)` -> [IsBlocked(ctx context.Context, groupID int64, userID int64)](../internal/domain/group/service/group_blocklist_service.go#L126)

- **UserRepository.java** ([java/im/turms/service/domain/user/repository/UserRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/repository/UserRepository.java))
> [简述功能]

  - [x] `updateUsers(Set<Long> userIds, @Nullable byte[] password, @Nullable String name, @Nullable String intro, @Nullable String profilePicture, @Nullable ProfileAccessStrategy profileAccessStrategy, @Nullable Long roleId, @Nullable Date registrationDate, @Nullable Boolean isActive, @Nullable Map<String, Object> userDefinedAttributes, @Nullable ClientSession session)` -> [UpdateUsers(ctx context.Context, userIDs []int64, update bson.M)](../internal/domain/user/repository/user_repository.go#L128)
  - [x] `updateUsersDeletionDate(Set<Long> userIds)` -> [UpdateUsersDeletionDate(ctx context.Context, userIDs []int64)](../internal/domain/user/repository/user_repository.go#L137)
  - [x] `checkIfUserExists(Long userId, boolean queryDeletedRecords)` -> [CheckIfUserExists(ctx context.Context, userID int64)](../internal/domain/user/service/user_service.go#L111)
  - [x] `countRegisteredUsers(@Nullable DateRange dateRange, boolean queryDeletedRecords)` -> [CountRegisteredUsers(ctx context.Context, startDate *time.Time, endDate *time.Time, queryDeletedRecords bool)](../internal/domain/user/repository/user_repository.go#L146)
  - [x] `countDeletedUsers(@Nullable DateRange dateRange)` -> [CountDeletedUsers(ctx context.Context, startDate *time.Time, endDate *time.Time)](../internal/domain/user/repository/user_repository.go#L165)
  - [x] `countUsers(boolean queryDeletedRecords)` -> [CountUsers(ctx context.Context, startDate *time.Time, endDate *time.Time)](../internal/domain/user/repository/user_repository.go#L183)
  - [x] `countUsers(@Nullable Set<Long> userIds, @Nullable DateRange registrationDateRange, @Nullable DateRange deletionDateRange, @Nullable Boolean isActive)` -> [CountUsers(ctx context.Context, startDate *time.Time, endDate *time.Time)](../internal/domain/user/repository/user_repository.go#L183)
  - [x] `findName(Long userId)` -> [FindName(ctx context.Context, userID int64)](../internal/domain/user/repository/user_repository.go#L203)
  - [x] `findAllNames()` -> [FindAllNames(ctx context.Context)](../internal/domain/user/repository/user_repository.go#L212)
  - [x] `findProfileAccessIfNotDeleted(Long userId)` -> [FindProfileAccessIfNotDeleted(ctx context.Context, userID int64)](../internal/domain/user/repository/user_repository.go#L225)
  - [x] `findUsers(@Nullable Collection<Long> userIds, @Nullable DateRange registrationDateRange, @Nullable DateRange deletionDateRange, @Nullable Boolean isActive, @Nullable Integer page, @Nullable Integer size, boolean queryDeletedRecords)` -> [FindUsers(ctx context.Context, userIDs []int64)](../internal/domain/user/repository/user_repository.go#L234)
  - [x] `findNotDeletedUserProfiles(Collection<Long> userIds, Collection<String> includedUserDefinedAttributes, @Nullable Date lastUpdatedDate)` -> [FindNotDeletedUserProfiles(ctx context.Context, userIDs []int64)](../internal/domain/user/repository/user_repository.go#L243)
  - [x] `findUsersProfile(Collection<Long> userIds, Collection<String> includedUserDefinedAttributes, boolean queryDeletedRecords)` -> [FindUsersProfile(ctx context.Context, userIDs []int64)](../internal/domain/user/repository/user_repository.go#L252)
  - [x] `findUserRoleId(Long userId)` -> [FindUserRoleID(ctx context.Context, userID int64)](../internal/domain/user/repository/user_repository.go#L261)
  - [x] `isActiveAndNotDeleted(Long userId)` -> [IsActiveAndNotDeleted(ctx context.Context, userID int64)](../internal/domain/user/repository/user_repository.go#L271)

- **UserRoleRepository.java** ([java/im/turms/service/domain/user/repository/UserRoleRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/repository/UserRoleRepository.java))
> [简述功能]

  - [x] `updateUserRoles(Set<Long> groupIds, @Nullable String name, @Nullable Set<Long> creatableGroupTypeIds, @Nullable Integer ownedGroupLimit, @Nullable Integer ownedGroupLimitForEachGroupType, @Nullable Map<Long, Integer> groupTypeIdToLimit)` -> [UpdateUserRoles(ctx context.Context, roleIDs []int64, update interface{})](../internal/domain/user/repository/user_role_repository.go#L79)

- **UserSettingsRepository.java** ([java/im/turms/service/domain/user/repository/UserSettingsRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/repository/UserSettingsRepository.java))
> [简述功能]

  - [x] `upsertSettings(Long userId, Map<String, Object> settings)` -> [UpsertSettings(ctx context.Context, userID int64, settings map[string]interface{})](../internal/domain/user/repository/user_settings_repository.go#L31)
  - [x] `unsetSettings(Long userId, @Nullable Collection<String> settingNames)` -> [UnsetSettings(ctx context.Context, userID int64, settingsNames []string)](../internal/domain/user/repository/user_settings_repository.go#L67)
  - [x] `findByIdAndSettingNames(Long userId, @Nullable Collection<String> settingNames, @Nullable Date lastUpdatedDateStart)` -> [FindByIdAndSettingNames(ctx context.Context, userID int64, names []string)](../internal/domain/user/repository/user_settings_repository.go#L81)

- **UserVersionRepository.java** ([java/im/turms/service/domain/user/repository/UserVersionRepository.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/repository/UserVersionRepository.java))
> [简述功能]

  - [x] `updateSpecificVersion(Long userId, @Nullable ClientSession session, String... fields)` -> [UpdateSpecificVersion(ctx context.Context, userIDs []int64, field string, updateDate time.Time)](../internal/domain/user/repository/user_version_repository.go#L104)
  - [x] `updateSpecificVersion(Long userId, @Nullable ClientSession session, String field)` -> [UpdateSpecificVersion(ctx context.Context, userIDs []int64, field string, updateDate time.Time)](../internal/domain/user/repository/user_version_repository.go#L104)
  - [x] `updateSpecificVersion(Set<Long> userIds, @Nullable ClientSession session, String... fields)` -> [UpdateSpecificVersion(ctx context.Context, userIDs []int64, field string, updateDate time.Time)](../internal/domain/user/repository/user_version_repository.go#L104)
  - [x] `findGroupJoinRequests(Long userId)` -> [FindGroupJoinRequestsVersion(ctx context.Context, userID int64)](../internal/domain/user/repository/user_version_repository.go#L135)
  - [x] `findJoinedGroup(Long userId)` -> [FindJoinedGroupVersion(ctx context.Context, userID int64)](../internal/domain/user/repository/user_version_repository.go#L144)
  - [x] `findReceivedGroupInvitations(Long userId)` -> [FindReceivedGroupInvitationsVersion(ctx context.Context, userID int64)](../internal/domain/user/repository/user_version_repository.go#L153)
  - [x] `findRelationships(Long userId)` -> [FindRelationships(ctx context.Context, ownerIDs []int64, relatedUserIDs []int64, groupIndexes []int32, isBlocked *bool, establishmentDateRange *turmsmongo.DateRange, page, size *int, session *mongo.Session,)](../internal/domain/user/repository/user_relationship_repository.go#L181)
  - [x] `findRelationshipGroups(Long userId)` -> [FindRelationshipGroups(ctx context.Context, ownerIDs []int64, groupIndexes []int32, page *int, size *int)](../internal/domain/user/repository/user_relationship_group_repository.go#L202)
  - [x] `findSentGroupInvitations(Long userId)` -> [FindSentGroupInvitationsVersion(ctx context.Context, userID int64)](../internal/domain/user/repository/user_version_repository.go#L178)
  - [x] `findSentFriendRequests(Long userId)` -> [FindSentFriendRequestsVersion(ctx context.Context, userID int64)](../internal/domain/user/repository/user_version_repository.go#L187)
  - [x] `findReceivedFriendRequests(Long userId)` -> [FindReceivedFriendRequestsVersion(ctx context.Context, userID int64)](../internal/domain/user/repository/user_version_repository.go#L196)

- **UserFriendRequestService.java** ([java/im/turms/service/domain/user/service/UserFriendRequestService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/service/UserFriendRequestService.java))
> [简述功能]

  - [x] `removeAllExpiredFriendRequests(Date expirationDate)` -> [RemoveAllExpiredFriendRequests(ctx context.Context, expirationDate time.Time)](../internal/domain/user/service/user_friend_request_service.go#L70)
  - [x] `hasPendingFriendRequest(@NotNull Long requesterId, @NotNull Long recipientId)` -> [HasPendingFriendRequest(ctx context.Context, requesterID, recipientID int64)](../internal/domain/user/repository/user_friend_request_repository.go#L50)
  - [x] `createFriendRequest(@Nullable Long id, @NotNull Long requesterId, @NotNull Long recipientId, @NotNull String content, @Nullable @ValidRequestStatus RequestStatus status, @Nullable @PastOrPresent Date creationDate, @Nullable @PastOrPresent Date responseDate, @Nullable String reason)` -> [CreateFriendRequest(ctx context.Context, requestID *int64, requesterID, recipientID int64, content string, status *po.RequestStatus, creationDate, responseDate *time.Time, reason *string)](../internal/domain/user/service/user_friend_request_service.go#L80)
  - [x] `authAndCreateFriendRequest(@NotNull Long requesterId, @NotNull Long recipientId, @Nullable String content, @NotNull @PastOrPresent Date creationDate)` -> [AuthAndCreateFriendRequest(ctx context.Context, requesterID, recipientID int64, content string, creationDate time.Time)](../internal/domain/user/service/user_friend_request_service.go#L151)
  - [x] `authAndRecallFriendRequest(@NotNull Long requesterId, @NotNull Long requestId)` -> [AuthAndRecallFriendRequest(ctx context.Context, requesterID, requestID int64)](../internal/domain/user/service/user_friend_request_service.go#L209)
  - [x] `updatePendingFriendRequestStatus(@NotNull Long requestId, @NotNull @ValidRequestStatus RequestStatus requestStatus, @Nullable String reason, @Nullable ClientSession session)` -> [UpdatePendingFriendRequestStatus(ctx context.Context, requestID int64, targetStatus po.RequestStatus, reason *string)](../internal/domain/user/service/user_friend_request_service.go#L292)
  - [x] `updateFriendRequests(@NotEmpty Set<Long> requestIds, @Nullable Long requesterId, @Nullable Long recipientId, @Nullable String content, @Nullable @ValidRequestStatus RequestStatus status, @Nullable String reason, @Nullable @PastOrPresent Date creationDate, @Nullable @PastOrPresent Date responseDate)` -> [UpdateFriendRequests(ctx context.Context, requestIds []int64, requesterID, recipientID *int64, content *string, status *po.RequestStatus, reason *string, creationDate *time.Time)](../internal/domain/user/repository/user_friend_request_repository.go#L106)
  - [x] `queryRecipientId(@NotNull Long requestId)` -> [QueryRecipientId(ctx context.Context, requestID int64)](../internal/domain/user/service/user_friend_request_service.go#L277)
  - [x] `queryRequesterIdAndRecipientIdAndStatus(@NotNull Long requestId)` -> [QueryRequesterIdAndRecipientIdAndStatus(ctx context.Context, requestID int64)](../internal/domain/user/service/user_friend_request_service.go#L282)
  - [x] `queryRequesterIdAndRecipientIdAndCreationDateAndStatus(@NotNull Long requestId)` -> [QueryRequesterIdAndRecipientIdAndCreationDateAndStatus(ctx context.Context, requestID int64)](../internal/domain/user/service/user_friend_request_service.go#L287)
  - [x] `authAndHandleFriendRequest(@NotNull Long friendRequestId, @NotNull Long requesterId, @NotNull @ValidResponseAction ResponseAction action, @Nullable String reason)` -> [AuthAndHandleFriendRequest(ctx context.Context, friendRequestID int64, requesterID int64, action po.ResponseAction, reason *string)](../internal/domain/user/service/user_friend_request_service.go#L320)
  - [x] `queryFriendRequestsWithVersion(@NotNull Long userId, boolean areSentByUser, @Nullable Date lastUpdatedDate)` -> [QueryFriendRequestsWithVersion(ctx context.Context, userID int64, isRecipient bool, lastUpdatedDate *time.Time)](../internal/domain/user/service/user_friend_request_service.go#L390)
  - [x] `queryFriendRequestsByRecipientId(@NotNull Long recipientId)` -> [QueryFriendRequestsByRecipientId(ctx context.Context, recipientID int64)](../internal/domain/user/service/user_friend_request_service.go#L380)
  - [x] `queryFriendRequestsByRequesterId(@NotNull Long requesterId)` -> [QueryFriendRequestsByRequesterId(ctx context.Context, requesterID int64)](../internal/domain/user/service/user_friend_request_service.go#L385)
  - [x] `deleteFriendRequests(@Nullable Set<Long> ids)` -> [DeleteFriendRequests(ctx context.Context, ids []int64)](../internal/domain/user/service/user_friend_request_service.go#L399)
  - [x] `queryFriendRequests(@Nullable Set<Long> ids, @Nullable Set<Long> requesterIds, @Nullable Set<Long> recipientIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange, @Nullable Integer page, @Nullable Integer size)` -> [QueryFriendRequests(ctx context.Context, ids, requesterIds, recipientIds []int64, statuses []po.RequestStatus, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd *time.Time, page, size *int)](../internal/domain/user/service/user_friend_request_service.go#L405)
  - [x] `countFriendRequests(@Nullable Set<Long> ids, @Nullable Set<Long> requesterIds, @Nullable Set<Long> recipientIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange)` -> [CountFriendRequests(ctx context.Context, ids, requesterIds, recipientIds []int64, statuses []po.RequestStatus, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd *time.Time)](../internal/domain/user/repository/user_friend_request_repository.go#L292)

- **UserRelationshipGroupService.java** ([java/im/turms/service/domain/user/service/UserRelationshipGroupService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/service/UserRelationshipGroupService.java))
> [简述功能]

  - [x] `createRelationshipGroup(@NotNull Long ownerId, @Nullable Integer groupIndex, @NotNull String groupName, @Nullable @PastOrPresent Date creationDate, @Nullable ClientSession session)` -> [CreateRelationshipGroup(ctx context.Context, ownerID int64, groupIndex *int32, groupName string, creationDate *time.Time, session *mongo.Session,)](../internal/domain/user/service/user_relationship_group_service.go#L58)
  - [x] `queryRelationshipGroupsInfos(@NotNull Long ownerId)` -> [QueryRelationshipGroupsInfos(ctx context.Context, ownerID int64)](../internal/domain/user/service/user_relationship_group_service.go#L112)
  - [x] `queryRelationshipGroupsInfosWithVersion(@NotNull Long ownerId, @Nullable Date lastUpdatedDate)` -> [QueryRelationshipGroupsInfosWithVersion(ctx context.Context, ownerID int64, lastUpdatedDate *time.Time,)](../internal/domain/user/service/user_relationship_group_service.go#L120)
  - [x] `queryGroupIndexes(@NotNull Long ownerId, @NotNull Long relatedUserId)` -> [QueryGroupIndexes(ctx context.Context, ownerID int64, relatedUserID int64)](../internal/domain/user/service/user_relationship_group_service.go#L140)
  - [x] `queryRelationshipGroupMemberIds(@NotNull Long ownerId, @NotNull Integer groupIndex)` -> [QueryRelationshipGroupMemberIds(ctx context.Context, ownerID int64, groupIndex int32,)](../internal/domain/user/service/user_relationship_group_service.go#L152)
  - [x] `queryRelationshipGroupMemberIds(@Nullable Set<Long> ownerIds, @Nullable Set<Integer> groupIndexes, @Nullable Integer page, @Nullable Integer size)` -> [QueryRelationshipGroupMemberIds(ctx context.Context, ownerID int64, groupIndex int32,)](../internal/domain/user/service/user_relationship_group_service.go#L152)
  - [x] `updateRelationshipGroupName(@NotNull Long ownerId, @NotNull Integer groupIndex, @NotNull String newGroupName)` -> [UpdateRelationshipGroupName(ctx context.Context, ownerID int64, groupIndex int32, newName string, session *mongo.Session)](../internal/domain/user/repository/user_relationship_group_repository.go#L163)
  - [x] `upsertRelationshipGroupMember(@NotNull Long ownerId, @NotNull Long relatedUserId, @Nullable Integer newGroupIndex, @Nullable Integer deleteGroupIndex, @Nullable ClientSession session)` -> [UpsertRelationshipGroupMember(ctx context.Context, ownerID int64, relatedUserID int64, newGroupIndex *int32, deleteGroupIndex *int32, session *mongo.Session,)](../internal/domain/user/service/user_relationship_group_service.go#L196)
  - [x] `updateRelationshipGroups(@NotEmpty Set<UserRelationshipGroup.@ValidUserRelationshipGroupKey Key> keys, @Nullable String name, @Nullable @PastOrPresent Date creationDate)` -> [UpdateRelationshipGroups(ctx context.Context, keys []po.UserRelationshipGroupKey, newName string, session *mongo.Session)](../internal/domain/user/repository/user_relationship_group_repository.go#L130)
  - [ ] `addRelatedUserToRelationshipGroups(@NotNull Long ownerId, @NotNull Integer groupIndex, @NotNull Long relatedUserId, @Nullable ClientSession session)`
  - [x] `deleteRelationshipGroupAndMoveMembersToNewGroup(@NotNull Long ownerId, @NotNull Integer deleteGroupIndex, @NotNull Integer newGroupIndex)` -> [DeleteRelationshipGroupAndMoveMembersToNewGroup(ctx context.Context, ownerID int64, deleteGroupIndex int32, newGroupIndex int32,)](../internal/domain/user/service/user_relationship_group_service.go#L339)
  - [x] `deleteAllRelationshipGroups(@NotEmpty Set<Long> ownerIds, @Nullable ClientSession session, boolean updateRelationshipGroupsVersion)` -> [DeleteAllRelationshipGroups(ctx context.Context, ownerIDs []int64, session *mongo.Session)](../internal/domain/user/repository/user_relationship_group_repository.go#L78)
  - [x] `deleteRelatedUserFromRelationshipGroup(@NotNull Long ownerId, @NotNull Long relatedUserId, @NotNull Integer groupIndex, @Nullable ClientSession session, boolean updateRelationshipGroupsMembersVersion)` -> [DeleteRelatedUserFromRelationshipGroup(ctx context.Context, ownerID int64, relatedUserID int64, groupIndexes []int32, session *mongo.Session)](../internal/domain/user/repository/user_relationship_group_member_repository.go#L139)
  - [x] `deleteRelatedUserFromAllRelationshipGroups(@NotNull Long ownerId, @NotNull Long relatedUserId, @Nullable ClientSession session, boolean updateRelationshipGroupsMembersVersion)` -> [DeleteRelatedUserFromAllRelationshipGroups(ctx context.Context, ownerID int64, relatedUserID int64, session *mongo.Session, updateVersion bool,)](../internal/domain/user/service/user_relationship_group_service.go#L450)
  - [x] `deleteRelatedUsersFromAllRelationshipGroups(@NotEmpty Set<UserRelationship.@ValidUserRelationshipKey Key> keys, @Nullable ClientSession session, boolean updateRelationshipGroupsMembersVersion)` -> [DeleteRelatedUsersFromAllRelationshipGroups(ctx context.Context, ownerID int64, relatedUserIDs []int64, session *mongo.Session)](../internal/domain/user/repository/user_relationship_group_member_repository.go#L167)
  - [x] `moveRelatedUserToNewGroup(@NotNull Long ownerId, @NotNull Long relatedUserId, @NotNull Integer currentGroupIndex, @NotNull Integer targetGroupIndex, boolean suppressIfAlreadyExistsInTargetGroup, @Nullable ClientSession session)` -> [MoveRelatedUserToNewGroup(ctx context.Context, ownerID int64, relatedUserID int64, currentGroupIndex int32, targetGroupIndex int32, suppressIfAlreadyExists bool, session *mongo.Session,)](../internal/domain/user/service/user_relationship_group_service.go#L506)
  - [x] `deleteRelationshipGroups()` -> [DeleteRelationshipGroups(ctx context.Context, ownerID int64, groupIndexes []int32, session *mongo.Session)](../internal/domain/user/repository/user_relationship_group_repository.go#L104)
  - [x] `deleteRelationshipGroups(@NotEmpty Set<UserRelationshipGroup.@ValidUserRelationshipGroupKey Key> keys)` -> [DeleteRelationshipGroups(ctx context.Context, ownerID int64, groupIndexes []int32, session *mongo.Session)](../internal/domain/user/repository/user_relationship_group_repository.go#L104)
  - [x] `queryRelationshipGroups(@Nullable Set<Long> ownerIds, @Nullable Set<Integer> indexes, @Nullable Set<String> names, @Nullable DateRange creationDateRange, @Nullable Integer page, @Nullable Integer size)` -> [QueryRelationshipGroups(ctx context.Context, ownerIDs []int64, groupIndexes []int32, page *int, size *int,)](../internal/domain/user/service/user_relationship_group_service.go#L562)
  - [x] `countRelationshipGroups(@Nullable Set<Long> ownerIds, @Nullable Set<Integer> indexes, @Nullable Set<String> names, @Nullable DateRange creationDateRange)` -> [CountRelationshipGroups(ctx context.Context, ownerIDs []int64, groupIndexes []int32)](../internal/domain/user/repository/user_relationship_group_repository.go#L189)
  - [x] `countRelationshipGroups(@Nullable Set<Long> ownerIds, @Nullable Set<Long> relatedUserIds)` -> [CountRelationshipGroups(ctx context.Context, ownerIDs []int64, groupIndexes []int32)](../internal/domain/user/repository/user_relationship_group_repository.go#L189)
  - [x] `countRelationshipGroupMembers(@Nullable Set<Long> ownerIds, @Nullable Set<Integer> groupIndexes)` -> [CountRelationshipGroupMembers(ctx context.Context, ownerIDs []int64, groupIndexes []int32)](../internal/domain/user/service/user_relationship_group_service.go#L556)

- **UserRelationshipService.java** ([java/im/turms/service/domain/user/service/UserRelationshipService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/service/UserRelationshipService.java))
> [简述功能]

  - [x] `deleteAllRelationships(@NotEmpty Set<Long> userIds, @Nullable ClientSession session, boolean updateRelationshipsVersion)` -> [DeleteAllRelationships(ctx context.Context, ownerIDs []int64, session *mongo.Session)](../internal/domain/user/repository/user_relationship_repository.go#L222), [DeleteAllRelationships(ctx context.Context, userIDs []int64, session *mongo.Session,)](../internal/domain/user/service/user_relationship_service.go#L216)
  - [x] `deleteOneSidedRelationships(@NotEmpty Set<UserRelationship.@ValidUserRelationshipKey Key> keys)` -> [DeleteOneSidedRelationships(ctx context.Context, ownerIDs, relatedUserIDs []int64, session *mongo.Session)](../internal/domain/user/repository/user_relationship_repository.go#L230), [DeleteOneSidedRelationships(ctx context.Context, ownerID int64, relatedUserIDs []int64, session *mongo.Session,)](../internal/domain/user/service/user_relationship_service.go#L241)
  - [x] `deleteOneSidedRelationship(@NotNull Long ownerId, @NotNull Long relatedUserId, @Nullable Integer groupIndex, @Nullable ClientSession session)` -> [DeleteOneSidedRelationship(ctx context.Context, ownerID int64, relatedUserID int64,)](../internal/domain/user/service/user_relationship_service.go#L274)
  - [x] `tryDeleteTwoSidedRelationships(@NotNull Long requesterId, @NotNull Long relatedUserId, @Nullable Integer groupId)` -> [TryDeleteTwoSidedRelationships(ctx context.Context, user1ID int64, user2ID int64, session *mongo.Session,)](../internal/domain/user/service/user_relationship_service.go#L192)
  - [x] `queryRelatedUserIdsWithVersion(@NotNull Long ownerId, @Nullable Set<Integer> groupIndexes, @Nullable Boolean isBlocked, @Nullable Date lastUpdatedDate)` -> [QueryRelatedUserIdsWithVersion(ctx context.Context, ownerID int64, groupIndexes []int32, isBlocked *bool, lastUpdatedDate *time.Time,)](../internal/domain/user/service/user_relationship_service.go#L392)
  - [x] `queryRelationshipsWithVersion(@NotNull Long ownerId, @Nullable Set<Long> relatedUserIds, @Nullable Set<Integer> groupIndexes, @Nullable Boolean isBlocked, @Nullable Date lastUpdatedDate)` -> [QueryRelationshipsWithVersion(ctx context.Context, ownerID int64, relatedUserIDs []int64, groupIndexes []int32, isBlocked *bool, lastUpdatedDate *time.Time,)](../internal/domain/user/service/user_relationship_service.go#L372)
  - [x] `queryRelatedUserIds(@Nullable Set<Long> ownerIds, @Nullable Boolean isBlocked)` -> [QueryRelatedUserIds(ctx context.Context, ownerIDs []int64, groupIndexes []int32, isBlocked *bool, page *int, size *int,)](../internal/domain/user/service/user_relationship_service.go#L360)
  - [x] `queryRelatedUserIds(@Nullable Set<Long> ownerIds, @Nullable Set<Integer> groupIndexes, @Nullable Boolean isBlocked)` -> [QueryRelatedUserIds(ctx context.Context, ownerIDs []int64, groupIndexes []int32, isBlocked *bool, page *int, size *int,)](../internal/domain/user/service/user_relationship_service.go#L360)
  - [x] `queryRelationships(@Nullable Set<Long> ownerIds, @Nullable Set<Long> relatedUserIds, @Nullable Set<Integer> groupIndexes, @Nullable Boolean isBlocked, @Nullable DateRange establishmentDateRange, @Nullable Integer page, @Nullable Integer size)` -> [QueryRelationships(ctx context.Context, ownerIDs []int64, relatedUserIDs []int64, groupIndexes []int32, isBlocked *bool, establishmentDateRange *turmsmongo.DateRange, page *int, size *int,)](../internal/domain/user/service/user_relationship_service.go#L345)
  - [x] `queryMembersRelationships(@Nullable Set<Long> ownerIds, @Nullable Set<Integer> groupIndexes, @Nullable Integer page, @Nullable Integer size)` -> [QueryMembersRelationships(ctx context.Context, ownerID int64, groupIndexes []int32, page, size *int, session *mongo.Session)](../internal/domain/user/repository/user_relationship_repository.go#L316)
  - [x] `countRelationships(@Nullable Set<Long> ownerIds, @Nullable Set<Long> relatedUserIds, @Nullable Set<Integer> groupIndexes, @Nullable Boolean isBlocked)` -> [CountRelationships(ctx context.Context, ownerIDs, relatedUserIDs []int64, groupIndexes []int32, isBlocked *bool, session *mongo.Session)](../internal/domain/user/repository/user_relationship_repository.go#L308), [CountRelationships(ctx context.Context, ownerIDs []int64, relatedUserIDs []int64, groupIndexes []int32, isBlocked *bool,)](../internal/domain/user/service/user_relationship_service.go#L412)
  - [x] `countRelationships(@Nullable Set<Long> ownerIds, @Nullable Set<Long> relatedUserIds, @Nullable Boolean isBlocked)` -> [CountRelationships(ctx context.Context, ownerIDs, relatedUserIDs []int64, groupIndexes []int32, isBlocked *bool, session *mongo.Session)](../internal/domain/user/repository/user_relationship_repository.go#L308), [CountRelationships(ctx context.Context, ownerIDs []int64, relatedUserIDs []int64, groupIndexes []int32, isBlocked *bool,)](../internal/domain/user/service/user_relationship_service.go#L412)
  - [x] `friendTwoUsers(@NotNull Long userOneId, @NotNull Long userTwoId, @Nullable ClientSession session)` -> [FriendTwoUsers(ctx context.Context, user1ID, user2ID int64)](../internal/domain/user/service/user_relationship_service.go#L308)
  - [x] `upsertOneSidedRelationship(@NotNull Long ownerId, @NotNull Long relatedUserId, @Nullable String name, @Nullable @PastOrPresent Date blockDate, @Nullable Integer newGroupIndex, @Nullable Integer deleteGroupIndex, @Nullable @PastOrPresent Date establishmentDate, boolean upsert, @Nullable ClientSession session)` -> [UpsertOneSidedRelationship(ctx context.Context, ownerID int64, relatedUserID int64, blockDate *time.Time, groupIndex *int32, establishmentDate *time.Time, name *string, session *mongo.Session,)](../internal/domain/user/service/user_relationship_service.go#L72), [upsertOneSidedRelationship(ctx context.Context, ownerID int64, relatedUserID int64, blockDate *time.Time, groupIndex *int32, establishmentDate *time.Time, name *string, session *mongo.Session,)](../internal/domain/user/service/user_relationship_service.go#L97)
  - [x] `isBlocked(@NotNull Long ownerId, @NotNull Long relatedUserId, boolean preferCache)` -> [IsBlocked(ctx context.Context, groupID int64, userID int64)](../internal/domain/group/service/group_blocklist_service.go#L126), [IsBlocked(ctx context.Context, ownerID, relatedUserID int64)](../internal/domain/user/service/user_relationship_service.go#L283)
  - [x] `isNotBlocked(@NotNull Long ownerId, @NotNull Long relatedUserId, boolean preferCache)` -> [IsNotBlocked(ctx context.Context, ownerID, relatedUserID int64)](../internal/domain/user/service/user_relationship_service.go#L299)
  - [x] `hasRelationshipAndNotBlocked(@NotNull Long ownerId, @NotNull Long relatedUserId)` -> [HasRelationshipAndNotBlocked(ctx context.Context, ownerID, relatedUserID int64, session *mongo.Session)](../internal/domain/user/repository/user_relationship_repository.go#L48), [HasRelationshipAndNotBlocked(ctx context.Context, ownerID, relatedUserID int64)](../internal/domain/user/service/user_relationship_service.go#L430)
  - [x] `hasRelationshipAndNotBlocked(@NotNull Long ownerId, @NotNull Long relatedUserId, boolean preferCache)` -> [HasRelationshipAndNotBlocked(ctx context.Context, ownerID, relatedUserID int64, session *mongo.Session)](../internal/domain/user/repository/user_relationship_repository.go#L48), [HasRelationshipAndNotBlocked(ctx context.Context, ownerID, relatedUserID int64)](../internal/domain/user/service/user_relationship_service.go#L430)
  - [x] `updateUserOneSidedRelationships(@NotEmpty Set<UserRelationship.@ValidUserRelationshipKey Key> keys, @Nullable String name, @Nullable @PastOrPresent Date blockDate, @Nullable @PastOrPresent Date establishmentDate)` -> [UpdateUserOneSidedRelationships(ctx context.Context, ownerID int64, relatedUserIDs []int64, blockDate *time.Time, establishmentDate *time.Time, name *string, session *mongo.Session,)](../internal/domain/user/repository/user_relationship_repository.go#L245), [UpdateUserOneSidedRelationships(ctx context.Context, userID int64, relatedUserIDs []int64, blockDate *time.Time, groupIndex *int32, deleteGroupIndex *int32, name *string, lastUpdatedDate *time.Time,)](../internal/domain/user/service/user_relationship_service.go#L141)
  - [x] `hasOneSidedRelationship(@NotNull Long ownerId, @NotNull Long relatedUserId)` -> [HasOneSidedRelationship(ctx context.Context, ownerID, relatedUserID int64, session *mongo.Session)](../internal/domain/user/repository/user_relationship_repository.go#L341)

- **UserRoleService.java** ([java/im/turms/service/domain/user/service/UserRoleService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/service/UserRoleService.java))
> [简述功能]

  - [x] `queryUserRoles(@Nullable Integer page, @Nullable Integer size)` -> [QueryUserRoles(ctx context.Context, filter bson.M)](../internal/domain/user/service/user_role_service.go#L23)
  - [x] `addUserRole(@Nullable Long groupId, @Nullable String name, @NotNull Set<Long> creatableGroupTypeIds, @NotNull Integer ownedGroupLimit, @NotNull Integer ownedGroupLimitForEachGroupType, @NotNull Map<Long, Integer> groupTypeIdToLimit)` -> [AddUserRole(ctx context.Context, role *po.UserRole)](../internal/domain/user/service/user_role_service.go#L29)
  - [x] `updateUserRoles(@NotEmpty Set<Long> groupIds, @Nullable String name, @Nullable Set<Long> creatableGroupTypeIds, @Nullable Integer ownedGroupLimit, @Nullable Integer ownedGroupLimitForEachGroupType, @Nullable Map<Long, Integer> groupTypeIdToLimit)` -> [UpdateUserRoles(ctx context.Context, roleIDs []int64, update interface{})](../internal/domain/user/repository/user_role_repository.go#L79)
  - [x] `deleteUserRoles(@Nullable Set<Long> groupIds)` -> [DeleteUserRoles(ctx context.Context, filter bson.M)](../internal/domain/user/service/user_role_service.go#L43)
  - [x] `queryUserRoleById(@NotNull Long id)` -> [QueryUserRoleById(ctx context.Context, roleID int64)](../internal/domain/user/service/user_role_service.go#L48)
  - [x] `queryStoredOrDefaultUserRoleByUserId(@NotNull Long userId)` -> [QueryStoredOrDefaultUserRoleByUserId(ctx context.Context, userID int64)](../internal/domain/user/service/user_role_service.go#L53)
  - [x] `countUserRoles()` -> [CountUserRoles(ctx context.Context, filter bson.M)](../internal/domain/user/service/user_role_service.go#L61)

- **UserService.java** ([java/im/turms/service/domain/user/service/UserService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/service/UserService.java))
> [简述功能]

  - [x] `isAllowedToSendMessageToTarget(@NotNull Boolean isGroupMessage, @NotNull Boolean isSystemMessage, @NotNull Long requesterId, @NotNull Long targetId)` -> [IsAllowedToSendMessageToTarget(ctx context.Context, isGroupMessage bool, isSystemMessage bool, requesterID int64, targetID int64)](../internal/domain/user/service/user_service.go#L166)
  - [x] `createUser(@Nullable Long id, @Nullable String rawPassword, @Nullable String name, @Nullable String intro, @Nullable String profilePicture, @Nullable @ValidProfileAccess ProfileAccessStrategy profileAccessStrategy, @Nullable Long roleId, @Nullable @PastOrPresent Date registrationDate, @Nullable Boolean isActive)` -> [CreateUser(ctx context.Context, password string, name string, intro string, profilePicture string, profileAccess int32, permissionGroupID int64, isActive bool)](../internal/domain/user/service/user_service.go#L47)
  - [x] `addUser(@Nullable Long id, @Nullable String rawPassword, @Nullable String name, @Nullable String intro, @Nullable String profilePicture, @Nullable @ValidProfileAccess ProfileAccessStrategy profileAccessStrategy, @Nullable Long roleId, @Nullable @PastOrPresent Date registrationDate, @Nullable Boolean isActive)` -> [AddUser(ctx context.Context, id int64, password string, name string, intro string, profilePicture string, profileAccess int32, permissionGroupID int64, registrationDate time.Time, isActive bool)](../internal/domain/user/service/user_service.go#L71)
  - [x] `isAllowToQueryUserProfile(@NotNull Long requesterId, @NotNull Long targetUserId)` -> [IsAllowToQueryUserProfile(ctx context.Context, requesterID int64, targetID int64)](../internal/domain/user/service/user_service.go#L175)
  - [x] `authAndQueryUsersProfile(@NotNull Long requesterId, @Nullable Set<Long> userIds, @Nullable String name, @Nullable Date lastUpdatedDate, @Nullable Integer skip, @Nullable Integer limit, @Nullable List<Integer> fieldsToHighlight)` -> [AuthAndQueryUsersProfile(ctx context.Context, requesterID int64, userIDs []int64, name string, lastUpdatedDate *time.Time, skip int, limit int)](../internal/domain/user/service/user_service.go#L181)
  - [x] `queryUserName(@NotNull Long userId)` -> [QueryUserName(ctx context.Context, userID int64)](../internal/domain/user/service/user_service.go#L136)
  - [x] `queryUsersProfile(@NotEmpty Collection<Long> userIds, boolean queryDeletedRecords)` -> [QueryUsersProfile(ctx context.Context, userIDs []int64)](../internal/domain/user/service/user_service.go#L127)
  - [x] `queryUserRoleIdByUserId(@NotNull Long userId)` -> [QueryUserRoleIDByUserID(ctx context.Context, userID int64)](../internal/domain/user/service/user_service.go#L194)
  - [x] `deleteUsers(@NotEmpty Set<Long> userIds, @Nullable Boolean deleteLogically)` -> [DeleteUsers(ctx context.Context, userIDs []int64)](../internal/domain/user/service/user_service.go#L117)
  - [x] `checkIfUserExists(@NotNull Long userId, boolean queryDeletedRecords)` -> [CheckIfUserExists(ctx context.Context, userID int64)](../internal/domain/user/service/user_service.go#L111)
  - [x] `updateUser(@NotNull Long userId, @Nullable String rawPassword, @Nullable String name, @Nullable String intro, @Nullable String profilePicture, @Nullable @ValidProfileAccess ProfileAccessStrategy profileAccessStrategy, @Nullable Long roleId, @Nullable Boolean isActive, @Nullable @PastOrPresent Date registrationDate, @Nullable Map<String, Value> userDefinedAttributes)` -> [UpdateUser(ctx context.Context, userID int64, update bson.M)](../internal/domain/user/service/user_service.go#L103)
  - [x] `queryUsers(@Nullable Collection<Long> userIds, @Nullable DateRange registrationDateRange, @Nullable DateRange deletionDateRange, @Nullable Boolean isActive, @Nullable Integer page, @Nullable Integer size, boolean queryDeletedRecords)` -> [QueryUsers(ctx context.Context, userIDs []int64)](../internal/domain/user/service/user_service.go#L149)
  - [x] `countRegisteredUsers(@Nullable DateRange dateRange, boolean queryDeletedRecords)` -> [CountRegisteredUsers(ctx context.Context, startDate *time.Time, endDate *time.Time, queryDeletedRecords bool)](../internal/domain/user/repository/user_repository.go#L146)
  - [x] `countDeletedUsers(@Nullable DateRange dateRange)` -> [CountDeletedUsers(ctx context.Context, startDate *time.Time, endDate *time.Time)](../internal/domain/user/repository/user_repository.go#L165)
  - [x] `countUsers(boolean queryDeletedRecords)` -> [CountUsers(ctx context.Context, startDate *time.Time, endDate *time.Time)](../internal/domain/user/repository/user_repository.go#L183)
  - [x] `countUsers(@Nullable Set<Long> userIds, @Nullable DateRange registrationDateRange, @Nullable DateRange deletionDateRange, @Nullable Boolean isActive)` -> [CountUsers(ctx context.Context, startDate *time.Time, endDate *time.Time)](../internal/domain/user/repository/user_repository.go#L183)
  - [x] `updateUsers(@NotEmpty Set<Long> userIds, @Nullable String rawPassword, @Nullable String name, @Nullable String intro, @Nullable String profilePicture, @Nullable @ValidProfileAccess ProfileAccessStrategy profileAccessStrategy, @Nullable Long roleId, @Nullable @PastOrPresent Date registrationDate, @Nullable Boolean isActive, @Nullable Map<String, Object> userDefinedAttributes)` -> [UpdateUsers(ctx context.Context, userIDs []int64, update bson.M)](../internal/domain/user/repository/user_repository.go#L128)

- **UserSettingsService.java** ([java/im/turms/service/domain/user/service/UserSettingsService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/service/UserSettingsService.java))
> [简述功能]

  - [x] `upsertSettings(Long userId, Map<String, Value> settings)` -> [UpsertSettings(ctx context.Context, userID int64, settings map[string]interface{})](../internal/domain/user/service/user_settings_service.go#L42)
  - [x] `deleteSettings(Collection<Long> userIds, @Nullable ClientSession clientSession)` -> [DeleteSettings(ctx context.Context, filter interface{})](../internal/domain/user/repository/user_settings_repository.go#L46)
  - [x] `unsetSettings(Long userId, @Nullable Set<String> settingNames)` -> [UnsetSettings(ctx context.Context, userID int64, keys []string)](../internal/domain/user/service/user_settings_service.go#L91)
  - [x] `querySettings(Long userId, @Nullable Set<String> settingNames, @Nullable Date lastUpdatedDateStart)` -> [QuerySettings(ctx context.Context, filter bson.M)](../internal/domain/user/service/user_settings_service.go#L100)

- **UserVersionService.java** ([java/im/turms/service/domain/user/service/UserVersionService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/service/UserVersionService.java))
> [简述功能]

  - [x] `queryRelationshipsLastUpdatedDate(@NotNull Long userId)` -> [QueryRelationshipsLastUpdatedDate(ctx context.Context, userID int64)](../internal/domain/user/service/user_version_service.go#L25)
  - [x] `querySentGroupInvitationsLastUpdatedDate(@NotNull Long userId)` -> [QuerySentGroupInvitationsVersion(ctx context.Context, userID int64)](../internal/domain/user/service/user_version_service.go#L34)
  - [x] `queryReceivedGroupInvitationsLastUpdatedDate(@NotNull Long userId)` -> [QueryReceivedGroupInvitationsVersion(ctx context.Context, userID int64)](../internal/domain/user/service/user_version_service.go#L43)
  - [x] `queryGroupJoinRequestsVersion(@NotNull Long userId)` -> [QueryGroupJoinRequestsVersion(ctx context.Context, userID int64)](../internal/domain/user/service/user_version_service.go#L53)
  - [x] `queryRelationshipGroupsLastUpdatedDate(@NotNull Long userId)` -> [QueryRelationshipGroupsLastUpdatedDate(ctx context.Context, userID int64)](../internal/domain/user/service/user_version_service.go#L62)
  - [x] `queryJoinedGroupVersion(@NotNull Long userId)` -> [QueryJoinedGroupVersion(ctx context.Context, userID int64)](../internal/domain/user/service/user_version_service.go#L71)
  - [x] `querySentFriendRequestsVersion(@NotNull Long userId)` -> [QuerySentFriendRequestsVersion(ctx context.Context, userID int64)](../internal/domain/user/service/user_version_service.go#L80)
  - [x] `queryReceivedFriendRequestsVersion(@NotNull Long userId)` -> [QueryReceivedFriendRequestsVersion(ctx context.Context, userID int64)](../internal/domain/user/service/user_version_service.go#L89)
  - [x] `upsertEmptyUserVersion(@NotNull Long userId, @NotNull Date timestamp, @Nullable ClientSession session)` -> [UpsertEmptyUserVersion(ctx context.Context, userID int64)](../internal/domain/user/repository/user_version_repository.go#L42)
  - [x] `updateRelationshipsVersion(@NotNull Long userId, @Nullable ClientSession session)` -> [UpdateRelationshipsVersion(ctx context.Context, userID int64)](../internal/domain/user/service/user_version_service.go#L99)
  - [x] `updateRelationshipsVersion(@NotEmpty Set<Long> userIds, @Nullable ClientSession session)` -> [UpdateRelationshipsVersion(ctx context.Context, userID int64)](../internal/domain/user/service/user_version_service.go#L99)
  - [x] `updateSentFriendRequestsVersion(@NotNull Long userId)` -> [UpdateSentFriendRequestsVersion(ctx context.Context, userID int64)](../internal/domain/user/service/user_version_service.go#L106)
  - [x] `updateReceivedFriendRequestsVersion(@NotNull Long userId)` -> [UpdateReceivedFriendRequestsVersion(ctx context.Context, userID int64)](../internal/domain/user/service/user_version_service.go#L113)
  - [x] `updateRelationshipGroupsVersion(@NotNull Long userId)` -> [UpdateRelationshipGroupsVersion(ctx context.Context, userID int64)](../internal/domain/user/service/user_version_service.go#L121)
  - [x] `updateRelationshipGroupsVersion(@NotEmpty Set<Long> userIds)` -> [UpdateRelationshipGroupsVersion(ctx context.Context, userID int64)](../internal/domain/user/service/user_version_service.go#L121)
  - [x] `updateRelationshipGroupsMembersVersion(@NotNull Long userId)` -> [UpdateRelationshipGroupsMembersVersion(ctx context.Context, userID int64)](../internal/domain/user/service/user_version_service.go#L129)
  - [x] `updateRelationshipGroupsMembersVersion(@NotEmpty Set<Long> userIds)` -> [UpdateRelationshipGroupsMembersVersion(ctx context.Context, userID int64)](../internal/domain/user/service/user_version_service.go#L129)
  - [x] `updateSentGroupInvitationsVersion(@NotNull Long userId)` -> [UpdateSentGroupInvitationsVersion(ctx context.Context, userID int64)](../internal/domain/user/service/user_version_service.go#L136)
  - [x] `updateReceivedGroupInvitationsVersion(@NotNull Long userId)` -> [UpdateReceivedGroupInvitationsVersion(ctx context.Context, userID int64)](../internal/domain/user/service/user_version_service.go#L143)
  - [x] `updateSentGroupJoinRequestsVersion(@NotNull Long userId)` -> [UpdateSentGroupJoinRequestsVersion(ctx context.Context, userID int64)](../internal/domain/user/service/user_version_service.go#L150)
  - [x] `updateJoinedGroupsVersion(@NotNull Long userId)` -> [UpdateJoinedGroupsVersion(ctx context.Context, userID int64)](../internal/domain/user/service/user_version_service.go#L157)
  - [x] `updateSpecificVersion(@NotNull Long userId, @Nullable ClientSession session, @NotEmpty String... fields)` -> [UpdateSpecificVersion(ctx context.Context, userIDs []int64, field string, updateDate time.Time)](../internal/domain/user/repository/user_version_repository.go#L104)
  - [x] `updateSpecificVersion(@NotNull Long userId, @Nullable ClientSession session, @NotNull String field)` -> [UpdateSpecificVersion(ctx context.Context, userIDs []int64, field string, updateDate time.Time)](../internal/domain/user/repository/user_version_repository.go#L104)
  - [x] `updateSpecificVersion(@NotEmpty Set<Long> userIds, @Nullable ClientSession session, @NotEmpty String... fields)` -> [UpdateSpecificVersion(ctx context.Context, userIDs []int64, field string, updateDate time.Time)](../internal/domain/user/repository/user_version_repository.go#L104)
  - [x] `delete(@NotEmpty Set<Long> userIds, @Nullable ClientSession session)` -> [Delete(key K)](../internal/domain/common/cache/sharded_map.go#L63)

- **NearbyUserService.java** ([java/im/turms/service/domain/user/service/onlineuser/NearbyUserService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/service/onlineuser/NearbyUserService.java))
> [简述功能]

  - [x] `queryNearbyUsers(@NotNull Long userId, @NotNull DeviceType deviceType, @Nullable Float longitude, @Nullable Float latitude, @Nullable Short maxCount, @Nullable Integer maxDistance, boolean withCoordinates, boolean withDistance, boolean withUserInfo)` -> [QueryNearbyUsers(ctx context.Context, userID int64, deviceType protocol.DeviceType, longitude *float32, latitude *float32, maxCount *int, maxDistance *float64, withCoordinates bool, withDistance bool, withUserInfo bool)](../internal/domain/user/service/onlineuser/nearby_user_service.go#L43)

- **SessionService.java** ([java/im/turms/service/domain/user/service/onlineuser/SessionService.java](../turms-orig/turms-service/src/main/java/im/turms/service/domain/user/service/onlineuser/SessionService.java))
> [简述功能]

  - [x] `disconnect(@NotNull Long userId, @NotNull SessionCloseStatus closeStatus)` -> [Disconnect(ctx context.Context, userID int64, closeStatus int)](../internal/domain/user/service/onlineuser/session_service.go#L43)
  - [x] `disconnect(@NotNull Long userId, @NotNull Set<@ValidDeviceType DeviceType> deviceTypes, @NotNull SessionCloseStatus closeStatus)` -> [Disconnect(ctx context.Context, userID int64, closeStatus int)](../internal/domain/user/service/onlineuser/session_service.go#L43)
  - [x] `disconnect(@NotNull Long userId, @NotNull @ValidDeviceType DeviceType deviceType, @NotNull SessionCloseStatus closeStatus)` -> [Disconnect(ctx context.Context, userID int64, closeStatus int)](../internal/domain/user/service/onlineuser/session_service.go#L43)
  - [x] `disconnect(@NotNull Set<Long> userIds, @NotNull SessionCloseStatus closeStatus)` -> [Disconnect(ctx context.Context, userID int64, closeStatus int)](../internal/domain/user/service/onlineuser/session_service.go#L43)
  - [x] `disconnect(@NotNull Set<Long> userIds, @Nullable Set<@ValidDeviceType DeviceType> deviceTypes, @NotNull SessionCloseStatus closeStatus)` -> [Disconnect(ctx context.Context, userID int64, closeStatus int)](../internal/domain/user/service/onlineuser/session_service.go#L43)
  - [x] `queryUserSessions(Set<Long> userIds)` -> [QueryUserSessions(ctx context.Context, userIDs []int64)](../internal/domain/user/service/onlineuser/session_service.go#L66)

- **LocaleUtil.java** ([java/im/turms/service/infra/locale/LocaleUtil.java](../turms-orig/turms-service/src/main/java/im/turms/service/infra/locale/LocaleUtil.java))
> [简述功能]

  - [x] `isAvailableLanguage(String languageId)` -> [IsAvailableLanguage(languageID string)](../internal/infra/locale/locale_util.go#L16)

- **ApiLoggingContext.java** ([java/im/turms/service/infra/logging/ApiLoggingContext.java](../turms-orig/turms-service/src/main/java/im/turms/service/infra/logging/ApiLoggingContext.java))
> [简述功能]

  - [x] `shouldLogRequest(TurmsRequest.KindCase requestType)` -> [ShouldLogRequest(requestType int)](../internal/infra/logging/api_logging_context.go#L12)
  - [x] `shouldLogNotification(TurmsRequest.KindCase requestType)` -> [ShouldLogNotification(requestType int)](../internal/infra/logging/api_logging_context.go#L18)

- **ClientApiLogging.java** ([java/im/turms/service/infra/logging/ClientApiLogging.java](../turms-orig/turms-service/src/main/java/im/turms/service/infra/logging/ClientApiLogging.java))
> [简述功能]

  - [x] `log(ClientRequest request, ServiceRequest serviceRequest, long requestSize, long requestTime, ServiceResponse response, long processingTime)` -> [Log(request interface{}, serviceRequest interface{}, requestSize int64, requestTime int64, response interface{}, processingTime int64)](../internal/infra/logging/client_api_logging.go#L12)

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

  - [x] `toList(Map<String, String> map)` -> [ToList(protoItems interface{})](../internal/infra/proto/proto_model_convertor.go#L4)
  - [x] `value2proto(Value.Builder builder, Object object)` -> [Value2Proto(value interface{})](../internal/infra/proto/proto_model_convertor.go#L10)

- **DefaultLanguageSettings.java** ([java/im/turms/service/storage/elasticsearch/DefaultLanguageSettings.java](../turms-orig/turms-service/src/main/java/im/turms/service/storage/elasticsearch/DefaultLanguageSettings.java))
> [简述功能]

  - [x] `getSetting(LanguageCode code)` -> [GetSetting()](../internal/storage/elasticsearch/default_language_settings.go#L8)

- **ElasticsearchClient.java** ([java/im/turms/service/storage/elasticsearch/ElasticsearchClient.java](../turms-orig/turms-service/src/main/java/im/turms/service/storage/elasticsearch/ElasticsearchClient.java))
> [简述功能]

  - [x] `healthcheck()` -> [Healthcheck(ctx context.Context)](../internal/storage/elasticsearch/elasticsearch_client.go#L17)
  - [x] `putIndex(String index, CreateIndexRequest request)` -> [PutIndex(ctx context.Context, request *model.CreateIndexRequest)](../internal/storage/elasticsearch/elasticsearch_client.go#L22)
  - [x] `putDoc(String index, String id, Supplier<ByteBuf> payloadSupplier)` -> [PutDoc(ctx context.Context)](../internal/storage/elasticsearch/elasticsearch_client.go#L27)
  - [x] `deleteDoc(String index, String id)` -> [DeleteDoc(ctx context.Context)](../internal/storage/elasticsearch/elasticsearch_client.go#L32)
  - [x] `deleteByQuery(String index, DeleteByQueryRequest request)` -> [DeleteByQuery(ctx context.Context, request *model.DeleteByQueryRequest)](../internal/storage/elasticsearch/elasticsearch_client.go#L37)
  - [x] `updateByQuery(String index, UpdateByQueryRequest request)` -> [UpdateByQuery(ctx context.Context)](../internal/storage/elasticsearch/elasticsearch_client.go#L42)
  - [x] `search(String index, SearchRequest request, ObjectReader reader)` -> [Search(ctx context.Context)](../internal/storage/elasticsearch/elasticsearch_client.go#L48)
  - [x] `bulk(BulkRequest request)` -> [Bulk(ctx context.Context, request *model.BulkRequest)](../internal/storage/elasticsearch/elasticsearch_client.go#L53)
  - [x] `deletePit(String scrollId)` -> [DeletePit(ctx context.Context, request *model.ClosePointInTimeRequest)](../internal/storage/elasticsearch/elasticsearch_client.go#L58)

- **ElasticsearchManager.java** ([java/im/turms/service/storage/elasticsearch/ElasticsearchManager.java](../turms-orig/turms-service/src/main/java/im/turms/service/storage/elasticsearch/ElasticsearchManager.java))
> [简述功能]

  - [x] `putUserDoc(Long userId, String name)` -> [PutUserDoc(ctx context.Context)](../internal/storage/elasticsearch/elasticsearch_manager.go#L18)
  - [x] `putUserDocs(Collection<Long> userIds, String name)` -> [PutUserDocs(ctx context.Context)](../internal/storage/elasticsearch/elasticsearch_manager.go#L23)
  - [x] `deleteUserDoc(Long userId)` -> [DeleteUserDoc(ctx context.Context)](../internal/storage/elasticsearch/elasticsearch_manager.go#L28)
  - [x] `deleteUserDocs(Collection<Long> userIds)` -> [DeleteUserDocs(ctx context.Context)](../internal/storage/elasticsearch/elasticsearch_manager.go#L33)
  - [x] `searchUserDocs(@Nullable Integer from, @Nullable Integer size, String name, @Nullable Collection<Long> ids, boolean highlight, @Nullable String scrollId, @Nullable String keepAlive)` -> [SearchUserDocs(ctx context.Context)](../internal/storage/elasticsearch/elasticsearch_manager.go#L38)
  - [x] `putGroupDoc(Long groupId, String name)` -> [PutGroupDoc(ctx context.Context)](../internal/storage/elasticsearch/elasticsearch_manager.go#L43)
  - [x] `putGroupDocs(Collection<Long> groupIds, String name)` -> [PutGroupDocs(ctx context.Context)](../internal/storage/elasticsearch/elasticsearch_manager.go#L48)
  - [x] `deleteGroupDocs(Collection<Long> groupIds)` -> [DeleteGroupDocs(ctx context.Context)](../internal/storage/elasticsearch/elasticsearch_manager.go#L53)
  - [x] `deleteAllGroupDocs()` -> [DeleteAllGroupDocs(ctx context.Context)](../internal/storage/elasticsearch/elasticsearch_manager.go#L58)
  - [x] `searchGroupDocs(@Nullable Integer from, @Nullable Integer size, String name, @Nullable Collection<Long> ids, boolean highlight, @Nullable String scrollId, @Nullable String keepAlive)` -> [SearchGroupDocs(ctx context.Context)](../internal/storage/elasticsearch/elasticsearch_manager.go#L63)
  - [x] `deletePitForUserDocs(String scrollId)` -> [DeletePitForUserDocs(ctx context.Context)](../internal/storage/elasticsearch/elasticsearch_manager.go#L68)

- **IndexTextFieldSetting.java** ([java/im/turms/service/storage/elasticsearch/IndexTextFieldSetting.java](../turms-orig/turms-service/src/main/java/im/turms/service/storage/elasticsearch/IndexTextFieldSetting.java))
> [简述功能]

  - [ ] `IndexTextFieldSetting(Map<String, Property> fieldToProperty, @Nullable IndexSettingsAnalysis analysis)`

- **BulkRequest.java** ([java/im/turms/service/storage/elasticsearch/model/BulkRequest.java](../turms-orig/turms-service/src/main/java/im/turms/service/storage/elasticsearch/model/BulkRequest.java))
> [简述功能]

  - [ ] `BulkRequest(List<Object> operations)`
  - [x] `serialize(BulkRequest value, JsonGenerator gen, SerializerProvider serializers)` -> [Serialize()](../internal/storage/elasticsearch/model/elasticsearch_model.go#L8)

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
  - [x] `merge(IndexSettingsAnalysis analysis)` -> [Merge(other *IndexSettingsAnalysis)](../internal/storage/elasticsearch/model/elasticsearch_model.go#L55)

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

  - [x] `migrate(Set<String> existingCollectionNames)` -> [Migrate()](../internal/storage/mongo/mongo_collection_migrator.go#L7)

- **MongoConfig.java** ([java/im/turms/service/storage/mongo/MongoConfig.java](../turms-orig/turms-service/src/main/java/im/turms/service/storage/mongo/MongoConfig.java))
> [简述功能]

  - [x] `adminMongoClient(TurmsPropertiesManager propertiesManager)` -> [AdminMongoClient()](../internal/storage/mongo/mongo_config.go#L7)
  - [x] `userMongoClient(TurmsPropertiesManager propertiesManager)` -> [UserMongoClient()](../internal/storage/mongo/mongo_config.go#L12)
  - [x] `groupMongoClient(TurmsPropertiesManager propertiesManager)` -> [GroupMongoClient()](../internal/storage/mongo/mongo_config.go#L17)
  - [x] `conversationMongoClient(TurmsPropertiesManager propertiesManager)` -> [ConversationMongoClient()](../internal/storage/mongo/mongo_config.go#L22)
  - [x] `messageMongoClient(TurmsPropertiesManager propertiesManager)` -> [MessageMongoClient()](../internal/storage/mongo/mongo_config.go#L27)
  - [x] `conferenceMongoClient(TurmsPropertiesManager propertiesManager)` -> [ConferenceMongoClient()](../internal/storage/mongo/mongo_config.go#L32)

- **MongoFakeDataGenerator.java** ([java/im/turms/service/storage/mongo/MongoFakeDataGenerator.java](../turms-orig/turms-service/src/main/java/im/turms/service/storage/mongo/MongoFakeDataGenerator.java))
> [简述功能]

  - [x] `populateCollectionsWithFakeData()` -> [PopulateCollectionsWithFakeData()](../internal/storage/mongo/mongo_fake_data_generator.go#L7)

- **RedisConfig.java** ([java/im/turms/service/storage/redis/RedisConfig.java](../turms-orig/turms-service/src/main/java/im/turms/service/storage/redis/RedisConfig.java))
> [简述功能]

  - [x] `newSequenceIdRedisClientManager(RedisProperties properties)` -> [NewSequenceIdRedisClientManager()](../internal/storage/redis/redis_config.go#L7)
  - [x] `sequenceIdRedisClientManager()` -> [SequenceIdRedisClientManager()](../internal/storage/redis/redis_config.go#L12)

