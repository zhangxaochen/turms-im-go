
# TurmsGatewayApplication.java
*Checked methods: main(String[] args)*

## main(String[] args)

- [ ] **Missing static initializer logic**: The Java version has a static initializer block that sets `TimeZone.setDefault(TimeZoneConst.ZONE)`, sets `io.netty.maxDirectMemory` to `"0"`, sets `spring.main.banner-mode` to `"off"`, and sets `spring.main.web-application-type` to `"none"`. None of this is present in the Go version.

- [ ] **Missing environment validation**: The Java version calls `validateEnv()` which loads utility classes (`CollectionUtil`, `ClassUtil`, `StringUtil`) to trigger JVM compatibility checks. The Go version has no equivalent environment validation.

- [ ] **Missing application bootstrap/initialization**: The Java version calls `SpringApplication.run(applicationClass, args)` to start the full Spring Boot application with component scanning across the `GATEWAY` and `SERVER_COMMON` packages. The Go version only prints a log message and has a TODO comment — no actual server initialization, dependency injection, or server startup occurs.

- [ ] **Missing error handling with graceful logger shutdown**: The Java version has a `catch` block that checks if `LoggerFactory` is initialized, attempts to close it with a 50-second timeout, falls back to `printStackTrace()`, and calls `System.exit(1)` to ensure the process terminates. The Go version has no error handling or graceful shutdown logic.

- [ ] **Missing `@Application(nodeType = NodeType.GATEWAY)` configuration**: The Java version declares the node type as `GATEWAY` via the `@Application` annotation, which is used at runtime for cluster node identification. The Go version has no equivalent node type configuration.

# ClientRequestDispatcher.java
*Checked methods: handleRequest(UserSessionWrapper sessionWrapper, ByteBuf serviceRequestBuffer), handleRequest0(UserSessionWrapper sessionWrapper, ByteBuf serviceRequestBuffer), handleServiceRequest(UserSessionWrapper sessionWrapper, SimpleTurmsRequest request, ByteBuf serviceRequestBuffer, TracingContext tracingContext)*

Now I have a thorough understanding of both implementations. Let me do the detailed comparison.

## HandleRequest

- [ ] **Missing pending request counting**: The Java `handleRequest` increments `pendingRequestCount` and wraps `handleRequest0` with `doFinally` to decrement it (shutdown hook coordination). The Go `HandleRequest` has no equivalent pending request tracking or `onPendingRequestHandled()` call. It only has a `defer` for panic recovery, which is not present in the Java version.

- [ ] **Error handling mismatch**: In Java, `handleRequest` catches synchronous exceptions from `handleRequest0` and returns `Mono.error(e)` while still calling `onPendingRequestHandled()`. In Go, `HandleRequest` returns whatever `HandleRequest0` returns directly. The panic recovery via `defer recover()` is Go-specific and doesn't correspond to the Java exception handling, which catches only `Exception` (not `Error`/throwable).

## HandleRequest0

- [ ] **Missing server availability check in heartbeat with server unavailable response encoding**: In the Java heartbeat-unavailable path (line 161-164), the server returns `ClientMessageEncoder.encodeResponse(System.currentTimeMillis(), HEARTBEAT_FAILURE_REQUEST_ID, ResponseStatusCode.SERVER_UNAVAILABLE, serviceAvailability.reason())` which is a proper TurmsNotification with timestamp, requestId, code, and reason. The Go version (line 102-103) creates a notification via `d.NotificationFactory.CreateWithReason(&HeartbeatFailureRequestId, ...)` and then marshals it. This looks potentially correct in spirit, but the Java version explicitly includes `System.currentTimeMillis()` as the timestamp. Need to verify if `CreateWithReason` also sets timestamp.

- [ ] **Missing UNRECOGNIZED_REQUEST fallback and requestType tracking on parse failure**: In Java, when parsing fails, `tempRequest` is set to `UNRECOGNIZED_REQUEST` (with `requestId=-1` and `type=KIND_NOT_SET`), and these values are used for subsequent metrics/logging. In Go, on parse failure, `requestType` remains `nil` and `requestID` may be 0 (if `req.RequestId` is nil), and no equivalent logging/metrics are performed for the error case. The Java code logs the error and records metrics even for corrupted/invalid requests, while the Go code only creates a notification and returns early without logging.

- [ ] **Missing error logging for server errors on parse failure path**: Java code at lines 218-226 uses `.onErrorResume()` to log server errors with `LOGGER.error(...)` using the tracing context. The Go code does not log any server errors from `HandleServiceRequest` — it just converts errors to notifications via `CreateFromError`.

- [ ] **Missing metrics recording**: Java code at lines 213-216 uses `.name(TURMS_CLIENT_REQUEST).tag(TURMS_CLIENT_REQUEST_TAG_TYPE, requestType.name()).metrics()` to record metrics for every request. The Go code has no equivalent metrics recording.

- [ ] **Missing TracingContext propagation**: Java code at lines 191-193 creates a `TracingContext` based on `supportsTracing(requestType)`, updates it at line 304, clears it in the `finally` block at line 320, and propagates it via `.contextWrite()` at lines 260-264. The Go code has no tracing context at all.

- [ ] **Missing permission check**: Java code at line 198 checks `session.hasPermission(requestType)` and sets `notificationMono = UNAUTHORIZED_REQUEST_ERROR_MONO` if permission is denied. The Go code at line 137-139 has this commented out (`// if !sessionWrapper.UserSession.HasPermission(requestType) { ... }`), meaning unauthorized requests will proceed to be handled instead of being rejected.

- [ ] **DeleteSessionRequest logging lock not properly implemented**: Java code at line 202 calls `session.acquireDeleteSessionRequestLoggingLock()` which is an atomic compare-and-swap that returns `false` if already locked (preventing duplicate logging of delete session). The Go code at line 140-142 just sets `canLogRequest = true` unconditionally, with a `// MOCK` comment. This means duplicate logging is not prevented.

- [ ] **Missing `version` and `sessionId` extraction in logging**: Java code at lines 235-244 extracts `version` and `sessionId` from `userSession` when logging. The Go code at lines 154-164 does not extract `version` or `sessionId` from the user session — `version` and `sessionId` remain `nil`.

- [ ] **HandleServiceRequest called even when parse fails**: In Java, when parsing fails, `notificationMono` is already set to an error mono, and `handleServiceRequest` is NOT called (line 206 check: `if notificationMono == null`). In Go, when parsing fails, the code goes to the `else` branch... wait, actually in Go the parse failure sets `notification` and then the `else` block (lines 129-168) is skipped. However, looking more carefully: the Go code's structure is `if err != nil { ... } else { ... }`, so the service request is NOT called on parse failure. This is actually correct behavior. But the logging/metrics for the error case is still missing (see above).

- [ ] **Response encoding difference**: Java code at line 258 uses `ProtoEncoder.getDirectByteBuffer(notification)` to encode the notification as a direct ByteBuf. Go code at line 171 uses `proto.Marshal(notification)`. The encoding approach is different but functionally equivalent for protobuf.

## HandleServiceRequest

- [ ] **Missing `serviceRequestBuffer.retain()` for default (generic) case**: In Java at line 313, when the request falls through to the `default` case, `serviceRequestBuffer.retain()` is called before passing to the inner `handleServiceRequest` method, because the outer `finally` block (line 319) releases the buffer. This retain+release pattern ensures the buffer lives long enough for the async service request. In Go at line 210, `d.handleGenericServiceRequest(sessionWrapper, request, serviceRequestBuffer)` is called but the buffer is a `[]byte` (not reference-counted), so `retain`/`release` isn't needed. This is actually correct for Go since `[]byte` is garbage-collected. **Not a bug** in this context.

- [ ] **DeleteSessionRequest returns different result**: In Java at lines 310-311, `DELETE_SESSION_REQUEST` calls `sessionController.handleDeleteSessionRequest(sessionWrapper)` which returns a `Mono<TurmsNotification>` directly (no mapping through `getNotificationFromHandlerResult`). In Go at lines 203-208, `HandleDeleteSessionRequest` returns a `(*RequestHandlerResult, error)`, and then `getNotificationFromHandlerResult` is called on it. The Java `handleDeleteSessionRequest` returns the notification directly, not a `RequestHandlerResult`. This means the Go code wraps the result through `getNotificationFromHandlerResult` which may produce a different notification structure than what the Java controller returns directly.

- [ ] **Missing `tracingContext.updateThreadContext()` and `tracingContext.clearThreadContext()`**: Java code at lines 304 and 320 updates and clears the tracing context. Go has no tracing context handling.

Here is the consolidated bug list:

## HandleRequest

- [ ] **Missing pending request counting**: The Java version tracks `pendingRequestCount` (increment on entry, decrement on completion via `doFinally`) with a shutdown hook that waits for all pending requests. The Go version has no equivalent pending request tracking.

## HandleRequest0

- [ ] **Missing metrics recording**: Java records metrics via `.name(TURMS_CLIENT_REQUEST).tag(TURMS_CLIENT_REQUEST_TAG_TYPE, requestType.name()).metrics()` for all requests. Go has no equivalent.

- [ ] **Missing server error logging**: Java uses `.onErrorResume()` to log server errors via `LOGGER.error("Failed to handle the service request: {}", request, throwable)` with tracing context. Go does not log server errors.

- [ ] **Missing TracingContext**: Java creates and propagates a `TracingContext` based on request type, updates thread context, and clears it in finally. Go has no tracing at all.

- [ ] **Missing permission check (commented out)**: Java checks `session.hasPermission(requestType)` and rejects unauthorized requests with `UNAUTHORIZED_REQUEST`. Go has this check commented out at lines 137-139.

- [ ] **DeleteSessionRequest logging lock not implemented**: Java calls `session.acquireDeleteSessionRequestLoggingLock()` (an atomic CAS) to prevent duplicate logging. Go unconditionally sets `canLogRequest = true`.

- [ ] **Missing `version` and `sessionId` in logging**: Java extracts these fields from `userSession` for logging. Go leaves `version` and `sessionId` as `nil`.

- [ ] **No logging/metrics for corrupted request path**: When parsing fails, Java still records metrics (with `requestType=KIND_NOT_SET`) and potentially logs. Go skips all logging/metrics for the error branch.

## HandleServiceRequest

- [ ] **DeleteSessionRequest result handling differs**: Java's `handleDeleteSessionRequest(sessionWrapper)` returns a `Mono<TurmsNotification>` directly. Go calls `getNotificationFromHandlerResult()` on the result, wrapping it differently than Java which returns the notification as-is from the controller.

- [ ] **Missing tracing context update and clear**: Java calls `tracingContext.updateThreadContext()` before the switch and `tracingContext.clearThreadContext()` in the finally block. Go has no tracing.

# IpRequestThrottler.java
*Checked methods: tryAcquireToken(ByteArrayWrapper ip, long timestamp)*

## TryAcquireToken

- [x] **Missing `timestamp` parameter**: The Java method signature is `tryAcquireToken(ByteArrayWrapper ip, long timestamp)` and passes `timestamp` to `bucket.tryAcquire(timestamp)` for refill calculation. The Go version `TryAcquireToken(ip string)` drops the `timestamp` parameter entirely and uses `rate.Limiter.Allow()` which uses `time.Now()` internally. This means the Go version cannot honor caller-provided timestamps for refill computation, changing behavior when timestamps are externally controlled or batched.

- [x] **Different rate limiting algorithm**: Java uses a custom `TokenBucket` with explicit token counting, CAS-based refill logic, configurable `capacity`, `tokensPerPeriod`, and `refillIntervalNanos`. Go uses `golang.org/x/time/rate.Limiter` which is a token bucket with a continuous refill model (not discrete period-based). This produces different throttling behavior: Java refills tokens in discrete batches per period (e.g., 10 tokens every 1 second), while Go's `rate.Limiter` adds tokens continuously at a steady rate.

- [x] **Shared context vs. static configuration**: In Java, all `TokenBucket` instances share a single `TokenBucketContext` that can be dynamically updated at runtime (via `propertiesManager.addGlobalPropertiesChangeListener`). The Go version stores `Limit` and `Burst` as struct fields at construction time, and new limiters created after an update would use the updated values, but **existing per-IP limiters are never updated** when `Limit`/`Burst` fields change. Java's shared context means all buckets immediately reflect new rate limiting settings.

- [x] **Unlimited condition is wrong**: Go returns `true` (unlimited) when `t.Burst <= 0 || t.Limit == 0`. In Java, whether the bucket is "unlimited" depends on the `TokenBucketContext` configuration — specifically, if `refillIntervalNanos <= 0`, the bucket returns `false` when empty (no refill). When `capacity` and `tokensPerPeriod` allow unlimited access, it's the initial token count that matters. The Go condition `t.Limit == 0` treating zero limit as "unlimited" is backwards — `rate.Limit == 0` means zero rate (no tokens), not infinite rate. And `Burst <= 0` as unlimited is also incorrect since `rate.NewLimiter` with burst 0 would allow zero requests.

- [x] **Cleanup logic is fundamentally different**: Java's cleanup iterates entries and only removes those that are both idle for 30+ minutes **and** have tokens >= initial tokens (i.e., fully replenished). Go's cleanup does a full map reset every 10 minutes, wiping **all** entries including actively-used ones. This means active IPs lose their rate limiting state and start fresh every 10 minutes, which is a behavioral difference.

- [x] **Missing session-closed listener**: Java registers a listener on `SessionService` that removes an IP's token bucket when a session closes (if tokens are replenished). The Go version has no equivalent mechanism for cleaning up on session close.

# NotificationFactory.java
*Checked methods: init(TurmsPropertiesManager propertiesManager), create(ResponseStatusCode code, long requestId), create(ResponseStatusCode code, @Nullable String reason, long requestId), create(ThrowableInfo info, long requestId), createBuffer(CloseReason closeReason), sessionClosed(long requestId)*

Now I have a complete picture. Let me compile the bugs.

## init (NewNotificationFactory)

- [ ] **Missing dynamic config updates**: The Java `init` method registers a `notifyAndAddLocalPropertiesChangeListener` that dynamically updates `returnReasonForServerError` whenever properties change at runtime. The Go `NewNotificationFactory` captures a static `*config.GatewayProperties` snapshot at construction time and never updates it. If the `ReturnReasonForServerError` property changes at runtime, the Go version will not reflect the change.

## Create (create(ResponseStatusCode code, long requestId))

- [ ] **Missing default reason from status code**: The Java `create(code, requestId)` calls `trySetReason(builder, code, code.getReason())`, which passes the status code's built-in default reason (e.g., "ok" for OK, "The client request is invalid" for INVALID_REQUEST). The Go `Create(requestID, code)` delegates to `CreateWithReason(requestID, code, "")`, passing an empty string instead of the status code's default reason. Since `trySetReason` returns early when `reason == ""`, **no reason is ever set** for any status code in this method. In Java, non-server-error codes would always have their default reason included.

## CreateWithReason (create(ResponseStatusCode code, @Nullable String reason, long requestId))

- [x] **Empty-string vs nil/null semantics mismatch**: The Java version uses `@Nullable String reason` where `null` triggers the fallback to `code.getReason()`: `reason == null ? code.getReason() : reason`. The Go version uses `reason string` (empty string `""` as zero value). When a caller passes an empty reason, the Go code treats it like Java's `null` (no reason set). But the Java version would still set `code.getReason()` as the reason even when the explicit reason is non-null but empty. More critically, when the Java `reason` parameter is null, it falls back to `code.getReason()` (the default reason for that status code). The Go version has no such fallback — it just uses the empty string directly.

## CreateFromError (create(ThrowableInfo info, long requestId))

- [ ] **Wrong default error code for non-TurmsError errors**: In Java, `create(ThrowableInfo info, long requestId)` always extracts `info.code()` from the `ThrowableInfo` record, which already contains the correct `ResponseStatusCode` (resolved by `ThrowableInfo.get(Throwable)`). The Go version defaults to `ResponseStatusCode_SERVER_INTERNAL_ERROR` for non-`TurmsError` errors and falls back to the generic error message. While this is architecturally different (Go uses `error` interface vs Java's `ThrowableInfo` record), it means any custom error types with specific status codes (like the Java equivalents of `RECORD_CONTAINS_DUPLICATE_KEY`, `RESOURCE_NOT_FOUND`, etc.) will all map to `SERVER_INTERNAL_ERROR` instead of their proper codes.
- [x] **Missing fallback to code's default reason**: The Java version passes `info.reason()` to `trySetReason`, which may be `null` — in which case `trySetReason` returns without setting a reason. However, the Java version could have a non-null reason from `ThrowableInfo`. The Go version sets `reason = err.Error()` for non-TurmsError errors, which is a reasonable but different behavior (Java would use the throwable's message via `ThrowableInfo.get()`). For TurmsError cases, `te.Message` is used which maps to `info.reason()`, which is correct.

## CreateBuffer (createBuffer(CloseReason closeReason))

- [ ] **Completely different method signature and missing CloseReason integration**: The Java `createBuffer(CloseReason closeReason)` takes a single `CloseReason` parameter and extracts `closeReason.businessStatusCode()`, `closeReason.closeStatus()`, and `closeReason.reason()` from it. It calls `ClientMessageEncoder.encodeCloseNotification(timestamp, closeStatus, code, getReason(code, closeReason.reason()))` — which is a specialized encoding that includes the close status. The Go version takes `(requestID *int64, code ResponseStatusCode, reason string)` and simply marshals a standard notification, losing the close status entirely. This means the serialized output is structurally different from the Java version.
- [ ] **Missing close status in the encoded output**: The Java `encodeCloseNotification` includes a `SessionCloseStatus` in the encoded notification (likely as a `closeStatus` field or via the data section). The Go version produces a plain `TurmsNotification` via `CreateWithReason` + `proto.Marshal`, which does not include any close status information.
- [ ] **Missing getReason logic for server errors**: The Java `createBuffer` uses a private `getReason(code, closeReason.reason())` method that applies the same `returnReasonForServerError` filter for server errors (returning `null` for server errors if the config is disabled). The Go `CreateBuffer` delegates to `CreateWithReason`, which does call `trySetReason` — so this part is actually handled. However, the input parameters are fundamentally different (no `CloseReason` object).

## SessionClosed (sessionClosed(long requestId))

- [ ] **No bug in logic — matches Java behavior**: Sets timestamp, requestId, and `SERVER_INTERNAL_ERROR` code. Does not set a reason, which matches the Java version (it does not call `trySetReason`). This method is correct.

# UserSession.java
*Checked methods: setConnection(NetConnection connection, ByteArrayWrapper ip), setLastHeartbeatRequestTimestampToNow(), setLastRequestTimestampToNow(), close(@NotNull CloseReason closeReason), isOpen(), isConnected(), supportsSwitchingToUdp(), sendNotification(ByteBuf byteBuf), sendNotification(ByteBuf byteBuf, TracingContext tracingContext), acquireDeleteSessionRequestLoggingLock(), hasPermission(TurmsRequest.KindCase requestType), toString()*

Now let me carefully compare each method.

## setConnection(NetConnection connection, ByteArrayWrapper ip)
- [x] **Missing IP assignment**: The Java version assigns both `this.connection = connection` and `this.ip = ip`, but the Go version at `connection.go:92-94` assigns `s.Conn = connection` but never assigns the `ip` parameter to `s.IP`. The `ip string` parameter is received but discarded.

## setLastHeartbeatRequestTimestampToNow()
- [x] **Missing nanosecond timestamp tracking**: The Java version updates both `lastHeartbeatRequestTimestampMillis` (via `System.currentTimeMillis()`) and `lastHeartbeatRequestTimestampNanos` (via `System.nanoTime()`). The Go version at `connection.go:41-43` only stores a millisecond timestamp in `lastHeartbeat`. The nanosecond timestamp (`lastHeartbeatRequestTimestampNanos`) is not tracked at all.

## setLastRequestTimestampToNow()
- [x] **Missing nanosecond timestamp tracking**: Same as above. The Java version updates both `lastRequestTimestampMillis` and `lastRequestTimestampNanos`. The Go version at `connection.go:52-54` only stores a millisecond timestamp in `lastRequest`. The nanosecond timestamp (`lastRequestTimestampNanos`) is not tracked at all.

## close(@NotNull CloseReason closeReason)
- [x] **Missing `isSessionOpen` state tracking**: The Java version maintains a separate `isSessionOpen` volatile boolean that is set to `false` on close, and the method returns `true` only if the session was previously open (i.e., the first close succeeds). The Go version at `connection.go:125-129` has no `isSessionOpen` equivalent — it only checks `s.Conn != nil` and does not track whether the session has already been closed.
- [x] **Missing return value**: The Java version returns `boolean` indicating whether the session was actually closed (was open). The Go version returns nothing (`void`).
- [x] **Missing close-reason propagation**: The Java version passes `closeReason` to `connection.close(closeReason)`. The Go version ignores the `closeReason` parameter entirely and calls `s.Conn.Close()` with no arguments.
- [x] **Missing warning log when connection is null**: The Java version logs a warning `"The connection is missing for the user session: {}"` when `isSessionOpen` is true but `connection == null`. The Go version silently does nothing when `Conn` is nil.

## isOpen()
- [x] **Wrong semantics**: The Java version at line 175-177 returns the `isSessionOpen` volatile boolean, which tracks whether the session is open (independently of whether a connection exists — a session can be open with UDP heartbeats even without a connection). The Go version at `connection.go:63-65` returns `s.Conn != nil`, which checks for the presence of a connection, not session openness. These are semantically different: after `close()` is called, Java returns `false` (session closed) while Go would still return `true` if the connection object hasn't been nulled out.

## isConnected()
- [x] **Missing `connection.isConnected()` check**: The Java version at line 179-181 returns `connection != null && connection.isConnected()` — it checks both that the connection exists AND that it is connected. The Go version at `connection.go:97-99` only checks `s.Conn != nil`, without calling any `IsActive()` or equivalent method on the connection.

## supportsSwitchingToUdp()
- No bugs. The Go version at `connection.go:102-104` correctly checks `s.DeviceType != protocol.DeviceType_BROWSER`, matching the Java logic `deviceType != DeviceType.BROWSER`.

## sendNotification(ByteBuf byteBuf)
- [x] **Method is completely missing**: The Java version has `sendNotification(ByteBuf byteBuf)` that calls `notificationConsumer.apply(byteBuf, TracingContext.NOOP)`. There is no corresponding method on `UserSession` in the Go code. The `sendNotification` in `router.go` is a completely different method on the `Router` struct that creates a new notification from scratch, rather than forwarding a pre-built ByteBuf via a consumer function.

## sendNotification(ByteBuf byteBuf, TracingContext tracingContext)
- [x] **Method is completely missing on UserSession**: The Java version stores a `BiFunction<ByteBuf, TracingContext, Mono<Void>> notificationConsumer` field and uses it in `sendNotification`. The Go version has no `notificationConsumer` field and no `SendNotification` method on `UserSession`. The `sendNotification` in `router.go:135-140` is a `Router` method that creates notifications via a factory — it is architecturally different from the Java version which forwards pre-built ByteBuf notifications from turms-service servers.

## acquireDeleteSessionRequestLoggingLock()
- No bugs. The Go version at `connection.go:112-114` correctly uses `atomic.CompareAndSwapUint32(&s.isDeleteSessionLockAcquired, 0, 1)`, matching the Java `AtomicIntegerFieldUpdater` compareAndSet logic.

## hasPermission(TurmsRequest.KindCase requestType)
- [ ] **Changed behavior on nil map**: The Java version uses a `Set<TurmsRequest.KindCase>` and calls `permissions.contains(requestType)`, which would throw `NullPointerException` if `permissions` were null (it's never null due to constructor). The Go version at `connection.go:117-122` adds a nil-check for `s.Permissions` and returns `false`. While this is arguably more defensive, the behavior differs — Java would crash on null permissions, Go silently returns false. More importantly, the Go `map[any]bool` type is less type-safe than Java's `Set<TurmsRequest.KindCase>`, but this is an idiomatic Go trade-off rather than a bug.

## toString()
- [x] **Stub implementation — missing all fields**: The Java version at lines 206-225 returns a string with `id`, `version`, `userId`, `deviceType`, `loginDate`, `loginLocation`, `isSessionOpen`, and `connection`. The Go version at `connection.go:107-109` returns the hardcoded literal `"UserSession{}"`, omitting all field values entirely.

# ServiceAvailabilityHandler.java
*Checked methods: channelRegistered(ChannelHandlerContext ctx), exceptionCaught(ChannelHandlerContext ctx, Throwable cause)*

Now I have all the information needed for a thorough comparison. Let me analyze the bugs.

## channelRegistered

- [x] **Missing `ctx.fireChannelRegistered()` equivalent when the IP is blocked**: In Java, when the service is available but the IP is blocked, `ctx.close()` is called and the method returns **without** calling `ctx.fireChannelRegistered()`. In Go, `HandleConnection` returns `false` (which maps to closing the connection), but the caller semantics may differ. However, this appears functionally equivalent since returning `false` signals rejection.

- [ ] **IP comparison uses string representation instead of raw bytes**: Java's `blocklistService.isIpBlocked()` takes `byte[]` (via `ByteArrayWrapper` wrapping `InetAddress.getAddress()`), comparing raw IP bytes. Go's `IsIpBlocked(ip string)` takes a string (`tcpAddr.IP.String()`). This is a behavioral difference — a string like `"192.168.001.001"` vs `"192.168.1.1"` or IPv6 representations could cause mismatched blocklist lookups compared to the Java byte-array comparison.

## exceptionCaught

- [ ] **Missing nil-check for sessions before iterating**: Java explicitly checks `if (sessions != null)` before iterating. In Go, `GetLocalUserSession` returns `[]*session.UserSession` (a slice), and if it returns `nil`, the `for range` loop on a nil slice is a no-op in Go, so this is actually safe. However, the Java code checks for `null` because `getLocalUserSession` can return `null`. In Go, the interface contract should ensure the same behavior, so this is not a real bug.

- [ ] **Unconditional propagation of the exception**: Java calls `ctx.fireExceptionCaught(cause)` at the end of the method regardless of the exception type. Go's `HandleException` returns the error unconditionally (`return cause`), which is the semantic equivalent. This is correct.

- [ ] **Unsafe type assertion without checking**: On line 55, `tcpAddr := addr.(*net.TCPAddr)` is an unchecked type assertion that will panic if `addr` is not a `*net.TCPAddr`. In Java, an explicit `(InetSocketAddress)` cast is used which would throw a `ClassCastException`. The Go code replicates the fail-fast behavior via panic, but the comment says "Replicate Java's implicit NullPointerException / ClassCastException behavior fail-fast" — so this is intentional and matches Java's behavior.

- [ ] **`conn.Close()` is called but the error is still propagated**: In Java, when `OutOfDirectMemoryError` occurs, `ctx.close()` is called, and then `ctx.fireExceptionCaught(cause)` is still called. In Go, `conn.Close()` is called, and then `return cause` propagates the error. This matches the Java behavior.

After careful analysis, the only substantive bug is the IP comparison difference:

## channelRegistered (HandleConnection)

- [ ] **IP address comparison uses string representation (`tcpAddr.IP.String()`) instead of raw bytes**: Java compares raw IP bytes via `ByteArrayWrapper(address.getAddress().getAddress())`, while Go compares using the string form `tcpAddr.IP.String()`. This can cause behavioral differences in blocklist matching — e.g., IPv4-mapped IPv6 addresses like `::ffff:192.168.1.1` may produce a different string than `192.168.1.1`, and different normalizations of the same IP may not match, whereas the Java byte-array comparison would match consistently. The Go `IsIpBlocked` interface takes `string` rather than `[]byte`, diverging from the Java `isIpBlocked(ByteArrayWrapper)` contract.

## exceptionCaught (HandleException)

- [ ] **No functional bugs found** — the Go implementation correctly mirrors the Java logic for this method.

# NetConnection.java
*Checked methods: getAddress(), send(ByteBuf buffer), close(CloseReason closeReason), close(), switchToUdp(), tryNotifyClientToRecover()*

Now I have all the information needed. Let me analyze each method carefully.

## getAddress()

**Java:** Returns `(InetSocketAddress) connection.address()` — returns the remote address as `InetSocketAddress`.
**Go:** Returns `c.conn.RemoteAddr()` — returns the remote address as `net.Addr`.

This is a faithful port. No bug.

## send(ByteBuf buffer)

**Java:** `connection.sendObject(buffer).then()` — sends the buffer reactively via Netty's ChannelOperations. No timeout is set at this level.
**Go:** Sets a `5 * time.Second` write deadline, then calls `c.conn.Write(buffer)`. Returns the error directly.

The Go version adds a hardcoded 5-second write deadline that does not exist in the Java version. The Java version relies on reactive Netty send without an explicit per-write timeout. This is a behavioral difference, but could be considered a reasonable adaptation. However, the key issue is that the Java `send()` is a simple reactive send with no timeout, while Go adds an arbitrary 5s deadline. This is a minor behavioral difference but not a "missing core logic" bug — it's a Go-appropriate addition.

No critical bug.

## close(CloseReason closeReason)

**Java (TcpConnection):**
1. Checks `!isConnected() || connection.isDisposed()` → returns early
2. Calls `super.close(closeReason)` — sets `isConnected=false`, `isConnectionRecovering=false`, `isSwitchingToUdp` based on SWITCH status
3. Creates a close notification buffer via `NotificationFactory.createBuffer(closeReason)`
4. Sends it with retry (2 retries, 3s backoff)
5. If `closeTimeout.isZero()` → calls `close()` in `doFinally`
6. If `closeTimeout` is non-negative → waits for `onTerminate` with timeout, then calls `close()` in `doFinally`
7. If `closeTimeout` is negative → mono is never subscribed with a `doFinally`, meaning `close()` is never called (no cleanup path)

**Go (TcpConnection):**
1. Checks `!c.IsConnected()` → returns early (missing `connection.isDisposed()` check)
2. Calls `c.BaseNetConnection.CloseWithReason(reason)`
3. If reason.Status != `UNKNOWN_ERROR` → goroutine: retries send 3 times with 3s sleep, then sleeps closeTimeout, then closes conn
4. If reason.Status == `UNKNOWN_ERROR` → immediately closes conn

**Bugs identified:**

- The Go version checks `reason.Status != constant.SessionCloseStatus_UNKNOWN_ERROR` to decide between graceful and immediate close. The Java version does **not** branch on the close reason type for the close logic — it always attempts to send the notification (unless already disconnected/disposed). The Go code skips the notification-send-retry entirely for `UNKNOWN_ERROR` status and just closes the connection directly.
- The Go version does not check `connection.isDisposed()` (or equivalent) in the guard condition.
- The Go version sends `[]byte{byte(reason.Status)}` as the notification, while Java uses `NotificationFactory.createBuffer(closeReason)`. This is a different notification format.
- When `closeTimeout` is negative in Java, the `close()` (dispose) is never called — the connection is just abandoned to be cleaned up by other means. In Go, when reason != UNKNOWN_ERROR, it always eventually closes via `c.conn.Close()`. When reason == UNKNOWN_ERROR, it also closes. The Java negative timeout behavior is not preserved.
- The Java version has a `doFinally` that calls `close()` (dispose) after the notification send completes (for zero and positive timeout). The Go version calls `c.conn.Close()` after sending the notification, but does NOT call `c.BaseNetConnection.Close()` (the no-arg close that resets flags) after the notification — it only closes the underlying socket.
- The `close()` method in Java's TcpConnection does NOT call `super.close()` (no-arg), it just disposes the connection. The Go `Close()` has a guard `if !c.IsConnected() { return nil }` and does NOT call `c.BaseNetConnection.Close()` either. This is actually consistent.

## close()

**Java (TcpConnection):** Just calls `connection.dispose()` with error handling. Does NOT call `super.close()` (which would reset the flags).
**Go (TcpConnection):** Checks `!c.IsConnected()` → returns early, then calls `c.conn.Close()`. Does NOT call `c.BaseNetConnection.Close()`.

The Go version adds an `IsConnected()` guard that the Java version does not have. Java's `close()` simply disposes unconditionally. This means in Java, `close()` can be called multiple times safely (dispose is idempotent), while in Go, after the first `Close()` call, `isConnected` may still be true if `BaseNetConnection.Close()` was never called, leading to inconsistent state. Actually wait — `close(CloseReason)` calls `super.close()` which sets `isConnected=false`, and then the Java `close()` override (no-arg) just disposes. In Go, `CloseWithReason` sets `isConnected=false`, and then `Close()` checks `!c.IsConnected()` which would return early. This means in Go, the final socket `conn.Close()` in the `CloseWithReason` goroutine would actually work because it's called directly, not via `Close()`. But the standalone `Close()` method has an extra guard.

The key bug: In Java, `TcpConnection.close()` (no-arg) does NOT call `super.close()`. It just disposes the connection. In Go, `TcpConnection.Close()` also does NOT call `BaseNetConnection.Close()`. However, the Go version has an `if !c.IsConnected() { return nil }` guard that Java doesn't have. This means if `Close()` is called independently (not via `CloseWithReason`), Java will still dispose the connection, but Go will only close if `isConnected` is true. This is a behavioral difference.

## switchToUdp()

**Java:** `close(CloseReason.get(SessionCloseStatus.SWITCH))` — calls the `close(CloseReason)` on the base class (NetConnection), not on TcpConnection's override.
**Go:** `b.CloseWithReason(NewCloseReason(constant.SessionCloseStatus_SWITCH))` — calls the base's `CloseWithReason`.

Wait, in Java, `switchToUdp()` is defined in `NetConnection` and calls `close(CloseReason)`. Since `close(CloseReason)` is overridden in `TcpConnection`, Java's `switchToUdp()` would actually call `TcpConnection.close(CloseReason)`, which does the full notification + dispose flow.

In Go, `SwitchToUdp()` is defined on `BaseNetConnection` and calls `b.CloseWithReason(...)`, which only sets flags. It does NOT call the `TcpConnection.CloseWithReason(...)` override.

**This is a significant bug.** In Java, `switchToUdp()` triggers the full TCP close flow (send notification, retry, dispose connection). In Go, `SwitchToUdp()` only updates the base flags and never actually closes the TCP connection or sends the SWITCH notification.

## tryNotifyClientToRecover()

**Java:** Checks `!isConnected && !isConnectionRecovering && udpAddress != null`, then calls `UdpRequestDispatcher.instance.sendSignal(udpAddress, UdpNotificationType.OPEN_CONNECTION)` and sets `isConnectionRecovering = true`.
**Go:** Same logic: checks `!b.isConnected && !b.isConnectionRecovering && b.udpAddress != nil`, then calls `b.udpSignalDispatcher(b.udpAddress)` and sets `b.isConnectionRecovering = true`.

The Go version uses an injectable callback pattern instead of a singleton, but the logic is equivalent. The callback is set up in `CreateConnection` to call `udp.Instance.SendSignal(addr, udp.OpenConnectionNotification)`, which mirrors `UdpRequestDispatcher.instance.sendSignal(udpAddress, UdpNotificationType.OPEN_CONNECTION)`. No bug.

---

## switchToUdp()

- [ ] **Critical Bug**: In Java, `switchToUdp()` calls `close(CloseReason)`, which resolves to `TcpConnection.close(CloseReason)` (the override), triggering the full close flow: send notification buffer, retry on failure, wait for timeout, then dispose connection. In Go, `SwitchToUdp()` is defined on `BaseNetConnection` and calls `b.CloseWithReason(...)` directly on the base struct, which only sets the `isConnected`/`isSwitchingToUdp`/`isConnectionRecovering` flags. It never calls the `TcpConnection.CloseWithReason(...)` override, so the TCP connection is never actually closed, no SWITCH notification is sent to the client, and the underlying socket is leaked.

## close(CloseReason closeReason)

- [ ] **Wrong condition for notification path**: The Go code branches on `reason.Status != constant.SessionCloseStatus_UNKNOWN_ERROR` to decide between graceful close (with notification retry) and immediate close. The Java version does NOT branch on the close reason — it always attempts to send the close notification (with retry) for any `CloseReason`, as long as `isConnected()` is true and the connection is not disposed. The `UNKNOWN_ERROR` check in Go is not present in the Java source.
- [ ] **Missing `connection.isDisposed()` guard**: Java checks `!isConnected() || connection.isDisposed()` before proceeding. Go only checks `!c.IsConnected()`, missing the disposed-connection check.
- [ ] **Incorrect notification payload**: Go sends `[]byte{byte(reason.Status)}` (a single byte of the status code). Java uses `NotificationFactory.createBuffer(closeReason)` which produces a properly formatted Turms protocol notification buffer. The Go version sends a raw status byte that the client cannot interpret as a valid Turms notification.
- [ ] **Missing `closeTimeout` negative/zero handling**: In Java, when `closeTimeout.isZero()`, the connection is closed immediately after the notification completes. When `closeTimeout` is positive, it waits for `onTerminate` with a timeout, then closes. When `closeTimeout` is negative, neither `doFinally` nor the termination wait is attached, so the connection is not explicitly disposed. Go always sleeps `closeTimeout` then closes the socket, not distinguishing between zero, positive, and negative timeout values.
- [ ] **Missing final `close()` call after notification**: In Java, after the notification + retry sequence completes (for zero or positive timeout), `close()` is called via `doFinally` to dispose the connection. In Go, the goroutine calls `c.conn.Close()` directly but does NOT call `c.Close()` (the no-arg version), so the Go `Close()` method's `IsConnected()` guard is never reached, meaning the `isConnected`/`isSwitchingToUdp`/`isConnectionRecovering` flags are not reset by the final cleanup.
- [ ] **Retry does not filter disconnected-client errors**: Java's `RETRY_SEND_CLOSE_NOTIFICATION` filters out retries for disconnected-client errors (`!ThrowableUtil.isDisconnectedClientError(throwable)`). Go's retry loop unconditionally retries on any error, wasting time retrying when the client is already gone.

## close()

- [ ] **Extra `IsConnected()` guard not present in Java**: Java's `TcpConnection.close()` unconditionally calls `connection.dispose()` without checking `isConnected`. Go's `Close()` returns early if `!c.IsConnected()`. This changes behavior: in Java, `close()` can be called as cleanup even after `close(CloseReason)` has set `isConnected=false`. In Go, `Close()` would be a no-op in that case.
- [ ] **Does not call `BaseNetConnection.Close()`**: Java's `TcpConnection.close()` also does not call `super.close()`, so this is actually consistent behavior. However, the flag state (`isConnected`, `isSwitchingToUdp`, `isConnectionRecovering`) is never reset by a standalone `close()` call in either language — this is by design in Java but worth noting in Go since the `IsConnected()` guard makes the standalone `Close()` unreachable after `CloseWithReason` has already set `isConnected=false`.

## send(ByteBuf buffer)

- [ ] **Added hardcoded 5-second write deadline not present in Java**: The Java `send()` method uses reactive Netty's `sendObject` without any per-write timeout. The Go version adds `c.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))` which introduces a timeout behavior not present in the original Java code.

# ExtendedHAProxyMessageReader.java
*Checked methods: channelRead(ChannelHandlerContext ctx, Object msg)*

## channelRead

- [ ] **Missing source address/port extraction from PROXY protocol header**: The Java version explicitly extracts `sourceAddress` and `sourcePort` from the `HAProxyMessage` and creates an unresolved `InetSocketAddress` from them. The Go version simply calls `conn.RemoteAddr()` on the raw connection, which returns the **direct TCP remote address**, NOT the address from the PROXY protocol header. It should instead check if the connection is a `*proxyproto.Conn`, extract the header via `conn.(*proxyproto.Conn).ProxyHeader()`, and use the source address from that header (falling back to `conn.RemoteAddr()` only when the header is absent or invalid).

- [ ] **Missing null/validity check on source address**: The Java code checks `if (sourceAddress != null && sourcePort > 0)` before using the proxy-provided address, falling back to `ctx.channel().remoteAddress()` when the source is null or port is invalid. The Go version has no equivalent guard — it unconditionally uses whatever `conn.RemoteAddr()` returns, never falling back to an alternative.

- [ ] **Missing handler self-removal after processing**: The Java code calls `ctx.channel().pipeline().remove(this)` after processing the HAProxy message, ensuring this handler runs only once. The Go `Read` method has no equivalent — calling it multiple times will invoke the callback repeatedly.

- [ ] **Missing `ctx.read()` continuation call**: After processing the proxy message and invoking the callback, the Java code calls `ctx.read()` to resume reading on the channel. The Go version has no equivalent mechanism to signal that reading should continue.

- [ ] **Missing pass-through for non-PROXY messages**: The Java code has an `else` branch that calls `super.channelRead(ctx, msg)` for messages that are not `HAProxyMessage` instances, ensuring they flow through the pipeline normally. The Go `Read` method always invokes the callback regardless, with no message-type discrimination.

# HAProxyUtil.java
*Checked methods: addProxyProtocolHandlers(ChannelPipeline pipeline, Consumer<InetSocketAddress> onRemoteAddressConfirmed), addProxyProtocolDetectorHandler(ChannelPipeline pipeline, Consumer<InetSocketAddress> onRemoteAddressConfirmed)*

Now I have a complete understanding of both implementations. Here is my review:

## AddProxyProtocolHandlers

- [ ] **Method body is completely empty** — The Java implementation adds a `HAProxyMessageDecoder` handler first, then an `ExtendedHAProxyMessageReader` handler after it in the channel pipeline. The Go version is a no-op stub with the comment "Replaced by WrapWithProxyProtocol interceptor above." While `WrapWithProxyProtocol` wraps the listener at the transport level, the critical `onRemoteAddressConfirmed` callback is **never invoked** from `AddProxyProtocolHandlers` or from `WrapWithProxyProtocol`. In Java, `ExtendedHAProxyMessageReader.channelRead()` calls `onRemoteAddressConfirmed` with the proxy-provided source address (or falls back to the channel's remote address). The Go code has an `ExtendedHAProxyMessageReader.Read()` method but it is **never called** anywhere — it is dead code. The callback parameter accepted by `AddProxyProtocolHandlers` is silently discarded.

- [ ] **Missing proxy header source address extraction** — In the Java `ExtendedHAProxyMessageReader.channelRead()`, the `HAProxyMessage` is parsed to extract `sourceAddress` and `sourcePort` via `proxyMessage.sourceAddress()` / `proxyMessage.sourcePort()`, and these are passed to `onRemoteAddressConfirmed` as an unresolved `InetSocketAddress`. In Go, `WrapWithProxyProtocol` uses `proxyproto.Listener` which makes `conn.RemoteAddr()` return the proxied address transparently. However, the callback (`onRemoteAddressConfirmed`) is never wired up to be called with this address, so the remote address confirmation logic is entirely absent.

## AddProxyProtocolDetectorHandler

- [ ] **Method body is completely empty** — The Java implementation adds an `ExtendedHAProxyMessageDetector` handler to the pipeline. This detector performs protocol detection: if the incoming data is **not** a PROXY protocol header, it removes itself and calls `onRemoteAddressConfirmed` with the direct channel remote address. If it **is** a PROXY protocol header, it replaces itself with the full `HAProxyMessageDecoder` + `ExtendedHAProxyMessageReader` pipeline. The Go version is a no-op stub. The callback parameter is silently discarded, meaning the `onRemoteAddressConfirmed` notification for non-proxy connections is never fired.

- [ ] **Missing fallback path for non-PROXY connections** — The Java `ExtendedHAProxyMessageDetector.decode()` has a critical branch: when `detectionResult` is `invalid` (i.e., the client is not sending PROXY protocol), it still calls `onRemoteAddressConfirmed` with `ctx.channel().remoteAddress()`. The Go code has no equivalent of this detection + fallback. Since `WrapWithProxyProtocol` unconditionally enables PROXY protocol parsing on the listener, non-PROXY connections would fail the handshake rather than being gracefully handled with a fallback callback.

# TcpConnection.java
*Checked methods: getAddress(), send(ByteBuf buffer), close(CloseReason closeReason), close()*

Now I have all the information needed for a thorough comparison. Let me analyze each method carefully.

## getAddress()

- [ ] **Return type mismatch**: Java returns `InetSocketAddress` (casting from `connection.address()`), while Go returns `net.Addr` (the generic interface). The Java version explicitly returns `InetSocketAddress`, providing access to IP address and port. Go's `net.Addr` is a generic interface — callers needing `*net.TCPAddr` (equivalent to `InetSocketAddress`) would need a type assertion.

## send(ByteBuf buffer)

- [ ] **Hardcoded write deadline not present in Java**: The Go version sets a `5 * time.Second` write deadline via `c.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))` before writing. The Java version uses `connection.sendObject(buffer).then()` with no timeout, relying on Netty's channel pipeline and backpressure. The 5-second hardcoded timeout is an invented constraint not present in the Java source.

## close(CloseReason closeReason)

- [ ] **Missing `connection.isDisposed()` guard**: Java checks `!isConnected() || connection.isDisposed()` before proceeding. Go only checks `!c.IsConnected()`, missing the dual guard for a disposed connection.
- [ ] **`super.close(closeReason)` called before the `isConnected` check in Go**: In Java, `super.close(closeReason)` is called *after* the early return guard (`if (!isConnected() || connection.isDisposed()) return`), meaning the state fields are only updated once. In Go, `c.BaseNetConnection.CloseWithReason(reason)` is called inside `CloseWithReason` on the `TcpConnection` after the `!c.IsConnected()` check — but this means it correctly mirrors the order. However, there is a subtle bug: the Go `BaseNetConnection.CloseWithReason` acquires the mutex and sets `isConnected = false`. But `TcpConnection.CloseWithReason` first checks `!c.IsConnected()` (which acquires a read lock), then calls `c.BaseNetConnection.CloseWithReason(reason)` (which acquires a write lock). Between the read lock release and write lock acquisition, another goroutine could interleave. The Java version uses `volatile` fields with no lock, accepting non-thread-safety (documented: "It is acceptable that the method isn't thread-safe"). The Go version adds mutex synchronization but doesn't hold the lock across the check-and-act sequence, making it partially but inconsistently thread-safe — neither matching Java's deliberately non-thread-safe semantics nor being fully thread-safe.
- [ ] **Missing retry filter for disconnected client errors**: Java uses `RETRY_SEND_CLOSE_NOTIFICATION` which is a `Retry.backoff(2, Duration.ofSeconds(3)).filter(throwable -> !ThrowableUtil.isDisconnectedClientError(throwable))`. This means retries are skipped if the error is a disconnected client error. The Go version retries unconditionally 3 times with `time.Sleep(3 * time.Second)` and no filtering — it will retry even when the client has disconnected, wasting resources and producing misleading log messages.
- [ ] **Missing error logging filter for disconnected clients**: Java logs "Failed to send the close notification" only if `!ThrowableUtil.isDisconnectedClientError(throwable)` — i.e., it suppresses logging for expected disconnection errors. The Go version logs all errors indiscriminately with `log.Printf("Failed to send close notification attempt %d: %v", ...)`.
- [ ] **Missing error logging in the subscribe handler**: Java has a separate `.subscribe()` error handler that logs "Failed to send the close notification after (2) attempts" (with the max attempts count from `RETRY_SEND_CLOSE_NOTIFICATION`), also filtered by `isDisconnectedClientError`. The Go version has no equivalent final failure log after all retries are exhausted.
- [ ] **Notification payload differs**: Java uses `NotificationFactory.createBuffer(closeReason)` to create a properly formatted notification buffer from the `CloseReason`. Go uses `[]byte{byte(reason.Status)}`, which is a raw single byte — this is a simplification that may not match the wire protocol expected by the client.
- [ ] **`closeTimeout == 0` branch not handled correctly**: Java has three branches: (1) `closeTimeout.isZero()` → send notification then immediately close, (2) `!closeTimeout.isNegative()` → send notification, wait for `connection.onTerminate()`, apply timeout, then close, (3) negative → mono is never subscribed/assigned, meaning no close happens via this path. Go has two branches: `closeTimeout > 0` → sleep then close, else fall through to immediate close. The Java `closeTimeout == 0` case still sends the notification and calls `close()` in `doFinally`, which the Go version handles by falling through to `c.conn.Close()`. But Go's `closeTimeout == 0` path does NOT send a notification first — it only does a `c.conn.Close()`. The notification send is gated by `reason.Status != constant.SessionCloseStatus_UNKNOWN_ERROR`, not by closeTimeout. The branching logic between Java's three timeout paths is not faithfully reproduced.
- [ ] **Missing `connection.onTerminate()` wait**: When `closeTimeout > 0`, Java waits for `connection.onTerminate()` (i.e., waits for the peer to close their side after receiving the notification), with a timeout fallback. Go simply does `time.Sleep(c.closeTimeout)` — it does not wait for the peer to actually terminate the connection, just sleeps a fixed duration.

## close()

- [ ] **Extra `!c.IsConnected()` guard not in Java**: Go's `Close()` starts with `if !c.IsConnected() { return nil }`. Java's `close()` has no such guard — it unconditionally calls `connection.dispose()`. This means if `close()` is called multiple times in Java, it will attempt `dispose()` each time (idempotent on Netty's `ChannelOperations`). In Go, subsequent calls are silently no-ops. The Java version does NOT call `super.close()` (it only calls `connection.dispose()`), meaning Java's `close()` does NOT set `isConnected = false`, `isConnectionRecovering = false`, `isSwitchingToUdp = false`. The Go version's `c.conn.Close()` also doesn't update those fields directly, but the `!c.IsConnected()` guard adds behavior not present in the original.
- [ ] **Go calls `c.conn.Close()` which updates the base connection state via `CloseWithReason` but Java's `close()` does NOT update parent class fields**: In Java, `TcpConnection.close()` only calls `connection.dispose()` — it does NOT call `super.close()`, meaning `isConnected`, `isSwitchingToUdp`, `isConnectionRecovering` remain unchanged. The Go code comment acknowledges this ("Java's TcpConnection.close() does NOT call super.close()"), and the Go version similarly doesn't call `c.BaseNetConnection.Close()`. However, the Go version adds the `if !c.IsConnected() { return nil }` guard that Java does not have, which changes behavior on repeated invocations.

# TcpServerFactory.java
*Checked methods: create(TcpProperties tcpProperties, BlocklistService blocklistService, ServerStatusManager serverStatusManager, SessionService sessionService, ConnectionListener connectionListener, int maxFrameLength)*

Now I have a thorough understanding of both implementations. Let me analyze the differences systematically.

## `create` (Java: static method `create(...)`, Go: `func (f *TcpServerFactory) Create(...)`)

- [ ] **Missing `ServiceAvailabilityHandler`**: The Java code creates a `ServiceAvailabilityHandler` with `blocklistService`, `serverStatusManager`, and `sessionService`, and adds it as the first handler in the pipeline (`pipeline.addFirst("serviceAvailabilityHandler", serviceAvailabilityHandler)`). The Go `Create` method does not accept `blocklistService`, `serverStatusManager`, or `sessionService` parameters, and does not create or wire any equivalent handler.

- [ ] **Missing `BlocklistService` IP blocking check**: The Java code checks `blocklistService.isIpBlocked(address.getAddress().getAddress())` in all three proxy protocol branches (REQUIRED, OPTIONAL, and implicitly via `ServiceAvailabilityHandler`). When the IP is blocked, it emits an empty signal (`remoteAddressSink.tryEmitEmpty()`), which prevents the connection from being processed. The Go code does not perform any IP blocklist check.

- [ ] **Missing `SessionService` and `ServerStatusManager` parameters**: The Java `create` method accepts `BlocklistService blocklistService`, `ServerStatusManager serverStatusManager`, and `SessionService sessionService`. The Go `Create` method replaces all of these with a simple `callback func(net.Conn)`, discarding all three service dependencies.

- [ ] **Missing `maxFrameLength` parameter and varint frame codec pipeline**: The Java method accepts `int maxFrameLength` and configures a full Netty channel pipeline with: a `varintLengthBasedFrameDecoder` (extended varint), a `varintLengthFieldPrepender`, and a `protobufFrameEncoder`. The Go method does not accept `maxFrameLength` and sets up no codec/frame decoding pipeline — raw connections are passed directly to the callback.

- [ ] **Missing `remoteAddressSink` / remote address resolution logic**: The Java code uses a `Sinks.One<InetSocketAddress>` to asynchronously resolve the remote address, with three branches: (1) REQUIRED proxy protocol — uses `HAProxyUtil.addProxyProtocolHandlers` with IP blocklist check, (2) OPTIONAL proxy protocol — uses `HAProxyUtil.addProxyProtocolDetectorHandler` with IP blocklist check, (3) no proxy — directly uses `channel.remoteAddress()`. The Go code collapses this into a single `bool proxy` flag that wraps with `go-proxyproto`, losing the REQUIRED vs OPTIONAL distinction and the blocklist check.

- [ ] **Missing `ConnectionListener.onAdded(...)` integration**: The Java code's `.handle()` block calls `connectionListener.onAdded(connection, remoteAddress, in.receive(), out, connection.onDispose())` after resolving the remote address. The Go code replaces this with a simple `callback func(net.Conn)`, losing the connection lifecycle integration, inbound stream subscription, and on-dispose handling.

- [ ] **Missing `setAutoRead(true)` equivalent**: The Java code explicitly calls `connection.channel().config().setAutoRead(true)` inside the `.handle()` block to trigger the read event, with detailed comments explaining why. The Go code has no equivalent — it relies on `net.Listener.Accept()` which has different read semantics.

- [ ] **Missing TCP socket options**: The Java code sets `CONNECT_TIMEOUT_MILLIS`, `SO_REUSEADDR` (both server and child option), `SO_BACKLOG`, `SO_LINGER=0` (child), and `TCP_NODELAY=true` (child). The Go code uses `net.Listen("tcp", addr)` which uses OS defaults and does not set any of these options (notably missing `TCP_NODELAY` and `SO_LINGER=0` on child connections).

- [ ] **Missing `wiretap` configuration**: The Java code applies `.wiretap(tcpProperties.isWiretap())`. The Go code has no equivalent.

- [ ] **Missing metrics recording**: The Java code enables metrics via `.metrics(true, () -> new TurmsMicrometerChannelMetricsRecorder(MetricNameConst.TURMS_GATEWAY_SERVER_TCP))`. The Go code has no metrics instrumentation.

- [ ] **Missing dedicated event loop / thread pool**: The Java code uses `.runOn(LoopResourcesFactory.createForServer(ThreadNameConst.GATEWAY_TCP_PREFIX))` to create a dedicated event loop with a named thread prefix. The Go code spawns goroutines directly from `go callback(conn)`, with no dedicated worker pool or thread naming.

- [ ] **Missing SSL/TLS configuration**: The Java code checks `tcpProperties.getSsl().isEnabled()` and conditionally calls `server.secure(...)` with `SslUtil.configureSslContextSpec(...)`. The Go code has no SSL/TLS support.

- [ ] **Missing `RemoteAddressSourceProxyProtocolMode` three-way branch**: The Java code distinguishes between `REQUIRED`, `OPTIONAL`, and neither (default). The Go code uses a boolean `proxy` parameter, which collapses REQUIRED and OPTIONAL into one case, losing the behavioral distinction (REQUIRED always expects a proxy header, OPTIONAL detects whether one is present).

- [ ] **Missing error wrapping with `BindException`**: The Java code catches bind failures and wraps them in a `BindException` with a descriptive message `"Failed to bind the TCP server on: " + host + ":" + port`. The Go code returns the raw `net.Listen` error without wrapping or contextual message.

- [ ] **Proxy protocol handler is added unconditionally to the entire listener**: In Java, proxy protocol handlers are added per-connection in the channel pipeline, with distinct behavior for REQUIRED vs OPTIONAL. In Go, `WrapWithProxyProtocol` wraps the entire listener, meaning all connections go through proxy protocol parsing. The Java default path (no proxy) directly uses `channel.remoteAddress()` without any proxy parsing overhead.

- [ ] **Missing `doOnChannelInit` pipeline setup ordering**: The Java code carefully orders handlers: `serviceAvailabilityHandler` first, then `varintLengthBasedFrameDecoder` before `ReactiveBridge`, then outbound handlers (`varintLengthFieldPrepender`, `protobufFrameEncoder`). This ordering is critical for correct protocol behavior. The Go code has no pipeline concept at all.

# TcpUserSessionAssembler.java
*Checked methods: getHost(), getPort()*

## getHost()

- [ ] **Missing disabled check semantics**: The Java version throws `FeatureDisabledException` (a specific exception type), while the Go version returns a generic `fmt.Errorf("TCP server is disabled")`. This loses the semantic distinction of a feature-disabled error vs. a generic error, which callers may rely on for conditional handling.

- [ ] **Empty string fallback instead of null**: When the server is disabled, the Java version sets `host = null` and would return `null` from `getHost()` if the disabled check were bypassed. The Go version initializes `Host` to `""` (empty string) and returns `""` in the disabled case. While Go doesn't have null strings, this is a behavioral difference that downstream consumers should be aware of.

## getPort()

- [ ] **Missing disabled check semantics**: Same as `getHost()` — the Java version throws `FeatureDisabledException` while Go returns a generic `fmt.Errorf`. The exception type specificity is lost.

- [ ] **Return value on disabled path**: The Go version returns `-1` along with an error when the server is disabled. The Java version also sets `port = -1` in the disabled path, so this is consistent. However, in Go, the caller receives both the value `-1` **and** an error, whereas Java throws and the value is never returned. A Go caller could ignore the error and use `-1`, which is a behavioral difference.

# UdpRequestDispatcher.java
*Checked methods: sendSignal(InetSocketAddress address, UdpNotificationType signal)*

## SendSignal

- [ ] **Missing `tryEmitNext` return value handling**: The Java version uses `notificationSink.tryEmitNext(...)` which returns a `Sinks.Emission` result that can indicate failures (e.g., `FAIL_TERMINATED`, `FAIL_OVERFLOW`). While the Java code ignores the return value here, the Go code's approach of spawning a goroutine on buffer-full is actually a *stricter* behavioral guarantee than the Java version. However, the Go fallback goroutine will **block indefinitely** if the channel is closed or never drained (unlike Java's `tryEmitNext` which would fail fast). This is a subtle behavioral difference — the Go version can leak goroutines on shutdown.

- [ ] **`notificationSink` nil check is insufficient**: When UDP is disabled, the Go constructor `NewUdpRequestDispatcher` returns `&UdpRequestDispatcher{}` without initializing `notificationSink`. Since `notificationSink` is a channel (a reference type), it will be `nil`, so the nil check `d.notificationSink != nil` works correctly. However, the `sessionService` field will also be `nil`, and `Instance` is **not set** when disabled (line 112 is inside the `enabled` branch), while Java always sets `instance = this` (line 77). This means Go code calling `udp.Instance.SendSignal(...)` when disabled will panic on nil pointer. Java's `instance` is always set regardless of enabled state.

# UdpSignalResponseBufferPool.java
*Checked methods: get(ResponseStatusCode code), get(UdpNotificationType type)*

Now I have all the information needed for the comparison. Let me analyze both methods.

## GetBufferFromStatusCode (mapped from `get(ResponseStatusCode code)`)

- [ ] **Missing caching/lazy initialization logic**: The Java version uses a `FastEnumMap` cache (`CODE_POOL`) with double-checked locking to lazily cache and reuse the same `ByteBuf` for each status code. The Go version creates a new `[]byte` slice on every call, defeating the purpose of a "buffer pool".
- [ ] **Missing `OK` code returns empty buffer with shared reference semantics**: The Java version returns `Unpooled.EMPTY_BUFFER` (a shared singleton empty buffer) for `OK`. The Go version returns `[]byte{}` which is functionally equivalent for OK, so this is acceptable. However, for non-OK codes, the Go version writes the raw enum int value (`uint16(code)`) instead of the `businessCode`. Looking at the Java, `code.getBusinessCode()` returns the `businessCode` field. In the Go code, the `ResponseStatusCode` constants are defined with their business code values directly (e.g., `ResponseStatusCode_OK = 1000`), and the conversion `uint16(code)` uses those values. So the business code mapping is correct. This is NOT a bug for the value itself.

## GetBufferFromNotificationType (mapped from `get(UdpNotificationType type)`)

- [ ] **Missing caching/pool semantics**: The Java version pre-allocates all notification buffers once in a static initializer and stores them in `NOTIFICATION_POOL` (a `FastEnumMap`). The Go version does pre-allocate in `init()` via `udpNotificationBuffers`, which is equivalent. The lookup returns the cached slice. This is functionally equivalent for the existing single enum value.
- [ ] **Fallback path differs from Java**: The Go version has a fallback `return []byte{byte(t) + 1}` when the index is out of bounds, which would create a new slice each time. The Java version simply does `NOTIFICATION_POOL.get(type)` which returns `null` for unknown types (since all enum values are pre-populated). The Go fallback silently produces a value instead of returning nil/empty, changing behavior for unexpected types.

Wait, let me re-examine the `get(ResponseStatusCode code)` more carefully. The Java version:
1. For `OK`: returns `Unpooled.EMPTY_BUFFER`
2. For non-OK: `writeShort(code.getBusinessCode())` — writes a 2-byte big-endian short of the business code

The Go version:
1. For `OK`: returns `[]byte{}` — correct
2. For non-OK: `uint16(code)` then `byte(val >> 8), byte(val)` — writes a 2-byte big-endian uint16 of the code value

Since the Go `ResponseStatusCode` constants store the business code directly as their value, `uint16(code)` is equivalent to `code.getBusinessCode()`. The byte encoding is also big-endian (`writeShort` in Java is big-endian). So the actual data is correct.

However, the most significant bug is:

## GetBufferFromStatusCode

- [ ] **Missing caching**: Every call allocates a new `[]byte` slice. The Java version caches buffers in `CODE_POOL` with double-checked locking so the same buffer is reused across calls. This is the core purpose of the "BufferPool" class. The Go implementation completely misses this caching behavior, which is the fundamental design of the original class.

## GetBufferFromNotificationType

- [ ] **Missing out-of-bounds protection consistent with Java**: Java's `NOTIFICATION_POOL.get(type)` returns `null` for non-existent types (which shouldn't happen since all enum values are pre-populated). The Go fallback creates a new buffer with `byte(t) + 1` instead of returning `nil`, which masks programming errors rather than surfacing them.

Let me format the final answer:

## GetBufferFromStatusCode

- [ ] **Missing lazy caching with double-checked locking**: The Java version caches the generated `ByteBuf` in `CODE_POOL` using double-checked locking (`synchronized`), so the same buffer instance is reused on every call. The Go version creates a brand-new `[]byte` slice on every invocation, completely defeating the "buffer pool" purpose of the class.
- [ ] **Missing synchronization for concurrent access**: The Java version uses `synchronized (CODE_POOL)` to ensure thread-safe lazy initialization. The Go version has no synchronization or concurrency protection for the cached buffers.

## GetBufferFromNotificationType

- [ ] **Fallback behavior differs**: The Java version returns `null` for an unrecognized `UdpNotificationType` (via `NOTIFICATION_POOL.get(type)`), while the Go version has a fallback that creates a new buffer `[]byte{byte(t) + 1}` for out-of-bounds indices. This silently produces data for unknown notification types instead of signaling an error or returning an empty/nil value, diverging from the Java contract.

# HttpForwardedHeaderHandler.java
*Checked methods: apply(ConnectionInfo connectionInfo, HttpRequest request)*

## Apply

- [ ] The method body is entirely unimplemented — it returns `nil` with a comment "Pending implementation" instead of containing any of the Java logic for parsing `Forwarded` or `X-Forwarded-*` headers.
- [ ] Missing the `isForwardedIpRequired` field on the struct. The Java constructor accepts `boolean isForwardedIpRequired` which controls whether an `IllegalArgumentException` is thrown when no forwarded IP is found. The Go struct has no fields at all.
- [ ] Missing the `Forwarded` header parsing path (`parseForwardedInfo`): splitting on `,`, regex-matching `for=`, `proto=`, `host=` directives, updating remote address, scheme, and host address accordingly.
- [ ] Missing the `X-Forwarded-For` / `X-Forwarded-Proto` / `X-Forwarded-Host` / `X-Forwarded-Port` header parsing path (`parseXForwardedInfo`): splitting on `,`, trimming, updating remote address, scheme, host address, and host port.
- [ ] Missing the `isForwardedIpRequired` validation: in both `parseForwardedInfo` (when no `for=` directive) and `parseXForwardedInfo` (when no `X-Forwarded-For` header), the Java code throws `IllegalArgumentException` if `isForwardedIpRequired` is true. The Go code has none of this.
- [ ] Missing the `X-Forwarded-Port` handling: the Java code parses the port, validates it's a valid integer (using `IntUtil.tryParse`), and calls `withHostAddress` with the parsed port and host string. The Go code has no equivalent.
- [ ] Method signature uses `any` types instead of concrete types for `connectionInfo` and `request` parameters and return value, losing type safety and making it impossible to actually implement the logic correctly without type assertions.

# WebSocketConnection.java
*Checked methods: getAddress(), send(ByteBuf buffer), close(CloseReason closeReason), close()*

Now I have a thorough understanding of both the Java and Go code. Here is my analysis:

## getAddress()

- [ ] **Return type mismatch**: Java's `getAddress()` returns `InetSocketAddress` (which provides both IP address and port), while Go's `GetAddress()` returns `net.Addr`. For TCP connections, `RemoteAddr()` returns a `*net.TCPAddr`, which is functionally similar but the interface type `net.Addr` is less specific than the Java version. This is minor but worth noting the return type is less constrained.

## send(ByteBuf buffer)

- [ ] **Missing framing/protocol wrapping for WebSocket**: The Java `WebSocketConnection.send()` wraps the buffer in a `BinaryWebSocketFrame` before sending. The Java `TcpConnection.send()` sends raw bytes. The Go `TcpConnection.Send()` sends raw bytes, which is correct for TCP. However, the Go code is being compared against the `WebSocketConnection` Java file, so if this is meant to be a WebSocket-equivalent, it's missing the WebSocket frame wrapping entirely.
- [ ] **Hardcoded write deadline**: The Go code sets a `5 * time.Second` write deadline (`c.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))`) that does not exist in either Java implementation. Neither `TcpConnection.send()` nor `WebSocketConnection.send()` in Java imposes any timeout on the send operation.
- [ ] **Non-error-handling difference**: Java returns `Mono<Void>` (reactive, non-blocking), Go returns `error` (synchronous, blocking). This is an architectural difference rather than a logic bug.

## close(CloseReason closeReason)

- [ ] **Missing `connection.isDisposed()` guard**: Java checks `!isConnected() || connection.isDisposed()` before proceeding. The Go code only checks `!c.IsConnected()` — it does not check whether the underlying connection is already closed/disposed (equivalent to checking `net.Conn` state), which could cause a panic or error when writing to an already-closed connection.
- [ ] **Missing notification content — sends raw byte instead of `NotificationFactory.createBuffer()`**: Java calls `NotificationFactory.createBuffer(closeReason)` which encodes a full `TurmsNotification` protobuf message with timestamp, close status, business status code, and optional reason string. The Go code sends only `[]byte{byte(reason.Status)}` — a single byte representing the status. This loses all the rich notification data (timestamp, business code, reason string).
- [ ] **Missing retry-with-backoff filter for disconnected-client errors**: Java uses `RETRY_SEND_CLOSE_NOTIFICATION` which is `Retry.backoff(2, Duration.ofSeconds(3)).filter(throwable -> !ThrowableUtil.isDisconnectedClientError(throwable))` — it only retries on non-disconnected-client errors. Go retries unconditionally with `time.Sleep(3 * time.Second)` between attempts, retrying even when the client has disconnected, wasting resources.
- [ ] **Missing `closeTimeout.isZero()` path (immediate close)**: In Java, when `closeTimeout` is zero, `close()` is called immediately in `doFinally` after the send attempt completes. In Go, when `closeTimeout` is zero, it skips the `if c.closeTimeout > 0` block and still calls `c.conn.Close()` — this path is actually correct in outcome but the flow differs: Java sends the notification and immediately closes in `doFinally`, while Go sends notification then immediately closes (no wait). This is functionally similar.
- [ ] **Missing `closeTimeout.isNegative()` (disabled close) path**: Java has three branches: `isZero()` → immediate close after send, `!isNegative()` → wait for close status/timeout then close, and implicitly negative → never close (the `mono` never calls `close()`). Go's code only has two paths: `> 0` (wait) or immediate close. If `closeTimeout` is negative in Go, it falls through to immediate `c.conn.Close()`, whereas Java would never close the connection. This is a behavioral difference.
- [ ] **Missing WebSocket close status wait (`receiveCloseStatus`)**: In the `WebSocketConnection` Java version, when `closeTimeout` is positive, it waits for `receiveCloseStatus()` from the client before closing. The Go TCP version has no equivalent — it just sleeps for `closeTimeout` duration. (This is somewhat expected for TCP vs WebSocket, but if comparing against the WebSocket Java version, it's missing.)
- [ ] **Missing `onTerminate()` wait for TCP version**: In the `TcpConnection` Java version, when `closeTimeout` is positive, it waits for `connection.onTerminate()` (i.e., waits for the connection to actually terminate) with a timeout. Go just does a blind `time.Sleep(c.closeTimeout)` instead of waiting for the connection to actually finish.
- [ ] **Missing `isSwitchingToUdp` flag is not set before sending notification**: In Go, `c.BaseNetConnection.CloseWithReason(reason)` is called which sets `isSwitchingToUdp` based on `SessionCloseStatus_SWITCH`, then the notification is sent in a goroutine. However, the check `if reason.Status != constant.SessionCloseStatus_UNKNOWN_ERROR` means for `UNKNOWN_ERROR` status, the connection is closed immediately without sending a notification. In Java, a notification is **always** sent (for any close reason), and then the close happens. The Go code skips notification entirely for `UNKNOWN_ERROR`, which differs from Java behavior.

## close()

- [ ] **Missing `isConnected` check in Java WebSocketConnection**: Java's `WebSocketConnection.close()` does **not** check `isConnected()` — it always sends a close frame. Go's `TcpConnection.Close()` checks `if !c.IsConnected() { return nil }` and returns early. This means if `close()` is called on an already-marked-disconnected connection in Go, it does nothing, whereas in Java WebSocket it still tries to send the close frame.
- [ ] **Missing `isConnected` check in Java TcpConnection**: Similarly, Java's `TcpConnection.close()` does **not** check `isConnected()` — it directly calls `connection.dispose()`. Go adds an `isConnected` guard not present in the Java original.
- [ ] **Go's `BaseNetConnection.Close()` is never called from `TcpConnection.Close()`**: The Go `TcpConnection.Close()` only closes the `net.Conn` but does **not** call `c.BaseNetConnection.Close()`, so `isConnected`, `isSwitchingToUdp`, and `isConnectionRecovering` flags in the base struct are never updated. The comment says "Java's TcpConnection.close() does NOT call super.close()" — but this is incorrect. Looking at Java `TcpConnection.close()`, it calls `connection.dispose()` but indeed does NOT call `super.close()`. However, in Java, `close(CloseReason)` always calls `super.close(closeReason)` first, setting `isConnected = false`, and then `close()` is called later from the `doFinally` callback. So by the time `close()` runs, `isConnected` is already false. In Go, the same pattern holds (closeWithReason sets isConnected=false), but the `Close()` method has its own guard that would prevent re-closing, which is actually safe. However, if `Close()` is ever called directly without first calling `CloseWithReason()`, the base struct flags will never be updated — unlike Java where `close()` from `NetConnection` base class sets `isConnected = false`.
- [ ] **Error handling logs on all errors vs filtering disconnected-client errors**: Java filters out "disconnected client" errors before logging (using `ThrowableUtil.isDisconnectedClientError`). Go logs all errors unconditionally, which will produce noisy log output when clients disconnect normally.

# WebSocketServerFactory.java
*Checked methods: create(WebSocketProperties webSocketProperties, BlocklistService blocklistService, ServerStatusManager serverStatusManager, SessionService sessionService, ConnectionListener connectionListener, int maxFramePayloadLength)*

Now I have a complete picture of both implementations. Let me analyze the differences carefully.

## Create

- [ ] **Missing parameter: `webSocketProperties`** — The Go `Create` method only takes `addr`, `handler`, and `sessionService`. The Java `create` takes `WebSocketProperties` which provides: host, port, connectTimeoutMillis, backlog, SSL config, and remote address source properties. None of these configuration-driven options are passed or used in the Go version.
- [ ] **Missing parameter: `blocklistService`** — The Java `create` accepts a `BlocklistService` which is used both to create `ServiceAvailabilityHandler` and to check if a client IP is blocked during handshake (`handleHttpRequest`). The Go version has no blocklist/IP blocking check anywhere.
- [ ] **Missing parameter: `serverStatusManager`** — The Java `create` uses `ServerStatusManager` to create `ServiceAvailabilityHandler` which is added as a Netty pipeline handler for service availability checks. The Go version has no service availability handler.
- [ ] **Missing parameter: `connectionListener`** — The Java `create` passes a `ConnectionListener` to handle the upgraded WebSocket connection (via `connectionListener.onAdded`). The Go version uses a simpler `session.MessageHandler` callback instead, losing the `onClose` lifecycle and the `inbound`/`outbound` Flux-based streaming model.
- [ ] **Missing parameter: `maxFramePayloadLength`** — The Java `create` passes this to `WebsocketServerSpec` and uses it in `aggregateFrames()` during frame processing. The Go upgrader uses hardcoded `ReadBufferSize: 1024` and `WriteBufferSize: 1024` instead of the configurable `maxFramePayloadLength`.
- [ ] **Missing: `ServiceAvailabilityHandler` pipeline initialization** — Java adds `serviceAvailabilityHandler` as the first handler in the Netty channel pipeline via `doOnChannelInit`. Go has no equivalent — no service availability check occurs during connection setup.
- [ ] **Missing: Proxy protocol mode configuration** — Java reads `RemoteAddressSourceProxyProtocolMode` (REQUIRED/OPTIONAL/DISABLED) and maps it to `ProxyProtocolSupportType`. Go hardcodes `&proxyproto.Listener{Listener: ln}` (always enabled) with no configuration option.
- [ ] **Missing: Forwarded header handler** — Java reads `RemoteAddressSourceHttpHeaderMode` (REQUIRED/OPTIONAL) and conditionally applies `HttpForwardedHeaderHandler` via `server.forwarded(...)`. The Go `HttpForwardedHeaderHandler` struct exists but its `Apply` method is a stub returning `nil`, and it is never wired into the server.
- [ ] **Missing: SSL/TLS support** — Java conditionally configures SSL via `SslUtil.configureSslContextSpec(...)` when `ssl.isEnabled()` is true. Go has no TLS/SSL configuration.
- [ ] **Missing: Socket options** — Java sets `CONNECT_TIMEOUT_MILLIS`, `SO_REUSEADDR` (option + childOption), `SO_BACKLOG`, `SO_LINGER=0`, and `TCP_NODELAY=true`. Go uses default `net.Listen("tcp", ...)` with no socket option tuning.
- [ ] **Missing: Custom event loop threads** — Java uses `LoopResourcesFactory.createForServer(ThreadNameConst.GATEWAY_WS_PREFIX)` to configure a named event loop for the server. Go uses the default Go scheduler with no custom thread/event loop naming.
- [ ] **Missing: Metrics recording** — Java enables metrics via `.metrics(true, ...)` with a `TurmsMicrometerChannelMetricsRecorder`. Go has no metrics instrumentation.
- [ ] **Missing: CORS preflight handling** — Java's `handleHttpRequest` checks `isPreFlightRequest(request)` and responds with CORS headers (`Access-Control-Allow-Origin: *`, etc.) before returning `Mono.never()`. The Go `handleHTTPFunc` has no CORS handling.
- [ ] **Missing: Handshake request validation** — Java validates the HTTP method is GET, checks `Upgrade: websocket` header, `Connection: upgrade` header, and `Sec-WebSocket-Key` header presence. Go's `upgrader.Upgrade()` handles some basic WebSocket validation internally but does not perform the same explicit error responses with custom status messages.
- [ ] **Missing: IP blocklist check before upgrade** — Java checks `blocklistService.isIpBlocked(...)` after handshake validation and silently drops blocked IPs by returning `Mono.empty()`. Go has no IP blocklist check.
- [ ] **Missing: Frame aggregation and filtering** — Java uses `in.aggregateFrames(maxFramePayloadLength)` and filters for `BinaryWebSocketFrame` only, discarding other frame types via `flatMap`. Go reads all messages with `conn.ReadMessage()` and only processes `BinaryMessage`, but does not handle frame aggregation for fragmented frames.
- [ ] **Missing: Close status handling** — Java captures `in.receiveCloseStatus()` and passes the `onClose` mono to `connectionListener.onAdded()`. Go uses a `CloseChan` on the `UserSession` but does not capture or propagate the WebSocket close status code.
- [ ] **Missing: Fallback remote address resolution** — Java has a fallback: `remoteAddress == null ? (InetSocketAddress) connection.channel().remoteAddress() : remoteAddress` for cases where the request's remote address is null after proxy protocol parsing. Go relies solely on `r.RemoteAddr` from the HTTP request.
- [ ] **Different: Bind error handling** — Java catches bind exceptions and wraps them in a custom `BindException` with a descriptive message including host and port. Go returns the raw `net.Listen` error without wrapping.
- [ ] **Missing: `Connection` object passthrough** — Java passes the reactor-netty `Connection` object to `connectionListener.onAdded()`, giving the listener full access to the underlying Netty channel. Go wraps the connection in `WSConnection` with limited methods, losing access to the raw connection/channel capabilities.

# NotificationService.java
*Checked methods: sendNotificationToLocalClients(TracingContext tracingContext, ByteBuf notificationData, Set<Long> recipientIds, Set<UserSessionId> excludedUserSessionIds, @Nullable DeviceType excludedDeviceType)*

Now I have a thorough understanding of both implementations. Let me do the detailed comparison.

## SendNotificationToLocalClients

- [ ] **Bug: Nil validation for `excludedUserSessionIds` is missing.** The Java version validates that `excludedUserSessionIds` is not null (`Validator.notNull(excludedUserSessionIds, "excludedUserSessionIds")`), returning an error if it is. The Go version does not check for a nil `excludedUserSessionIds` map, which would cause a panic on `len(excludedUserSessionIds)` if nil is passed (though in practice Go's `len(nil map)` returns 0, so this is minor — but semantically the Java version explicitly rejects null).

- [ ] **Bug: Nil validation for `notificationData` is incorrect.** The Java version checks `Validator.notNull(notificationData, "notificationData")` — i.e., it rejects null references, not empty content. The Go version checks `len(notificationData) == 0`, which is a different condition (rejects empty/nil byte slices). The Java version accepts a non-null ByteBuf with zero readable bytes.

- [ ] **Bug: Uses `GetAllUserSessions` instead of `GetUserSessionsManager`.** The Java version calls `sessionService.getUserSessionsManager(recipientId)` and then iterates over `userSessionsManager.getDeviceTypeToSession().values()`. The Go version calls `s.sessionService.GetAllUserSessions(recipientID)`. While functionally similar (both iterate over all sessions for a user), the Java version first checks if the manager is null (user offline) before iterating. The Go `GetAllUserSessions` returns `[]*UserSession` directly, conflating the null-check and the iteration — this is acceptable *only if* `GetAllUserSessions` returns an empty slice (not nil) for users with no sessions. However, the semantic difference is that the Java code uses the session manager's concurrency model, while the Go version bypasses it.

- [ ] **Bug: `TryNotifyClientToRecover` is called only on success, but in Java it's called unconditionally for every sent session.** In the Java code (line 160-161), `userSession.getConnection().tryNotifyClientToRecover()` is called unconditionally *after* adding the send mono to the list — meaning it fires regardless of whether the send succeeds or fails. In the Go code (line 72), `TryNotifyClientToRecover` is only called in the `else` (success) branch. This is a behavioral difference: in Java, the client recovery notification is attempted even if the notification send might fail asynchronously.

- [ ] **Bug: No plugin extension point invocation.** The Java version invokes `NotificationHandler` plugin extensions via `invokeExtensionPointHandlers` after all notifications are sent. The Go version omits this entirely (noted as "omitted as per stubbing strategy" in the comment). This is a missing core feature, not just logging.

- [ ] **Bug: No notification logging.** The Java version logs notification details via `notificationLoggingManager.log(...)` when notification logging is enabled, and also logs errors via `LOGGER.error(...)`. The Go version omits all logging. The `isNotificationLoggingEnabled` field is commented out, and there's no error logging at all (the error logging in the `if userSession.IsOpen()` block is just a comment `// log error`).

- [ ] **Bug: No error aggregation/propagation.** The Java version uses `Mono.whenDelayError(monos)` to wait for all sends to complete (collecting errors), and uses `.onErrorComplete(t -> true)` to suppress errors while still returning results. The Go version sends synchronously and if a send fails for one session, it still continues to others (which is correct), but it never aggregates or returns the error — it always returns `nil` as the error. In Java, errors from failed sends that occur while the session is still open are propagated as a combined error.

- [ ] **Bug: No tracing context propagation.** The Java version receives a `TracingContext` and passes it to `userSession.sendNotification(notificationData, tracingContext)`, and uses `TracingCloseableContext` in the logging callbacks. The Go version accepts a `context.Context` parameter but never passes it to any child operation or uses it for tracing.

- [ ] **Bug: No reference counting for notification data buffer.** The Java version uses Netty's `ByteBuf` reference counting: it calls `notificationData.retain()` before each send (line 144) and `notificationData.release()` in `doFinally` (line 209). This is critical for shared buffer management. The Go version uses `[]byte` (value type), so reference counting is not applicable — but this means every session gets the same byte slice without copying. If the downstream `WriteMessage` is asynchronous and the `notificationData` slice could be modified by the caller, this is a potential data race. This may or may not be a bug depending on Go's `WriteMessage` implementation (if it copies the bytes internally, it's fine).

- [ ] **Bug: Sending is synchronous instead of asynchronous/concurrent.** The Java version collects `Mono<Void>` for each send and executes them concurrently with `Mono.whenDelayError(monos)`. The Go version sends to each session sequentially in a for loop. For high-throughput scenarios with many recipients, this is a significant behavioral/performance difference — the Java version sends notifications to all sessions in parallel, while the Go version sends them one at a time.

# ServiceRequestService.java
*Checked methods: handleServiceRequest(UserSession session, ServiceRequest serviceRequest)*

## HandleServiceRequest

- [ ] **Missing buffer retain/release lifecycle**: The Java code calls `serviceRequest.getTurmsRequestBuffer().retain()` before the async operation and `release()` in a `finally` block. The Go code has these as TODO comments and does not implement any buffer reference counting.
- [ ] **Missing RPC call to forward the request**: The Java code creates a `HandleServiceRequest` wrapper and calls `node.getRpcService().requestResponse(request)` to forward the request to the cluster. The Go code has this as a TODO and returns an empty stub notification instead.
- [ ] **Missing `defaultIfEmpty(REQUEST_RESPONSE_NO_CONTENT)` fallback**: The Java code calls `.defaultIfEmpty(REQUEST_RESPONSE_NO_CONTENT)` so that if the RPC returns an empty response, a `NO_CONTENT` response is used. The Go code does not implement this fallback.
- [ ] **Missing `getNotificationFromResponse` mapping**: The Java code maps the `ServiceResponse` into a `TurmsNotification` via `getNotificationFromResponse`, which sets `timestamp` (current millis), `code` (business code), `requestId`, `reason`, and `data`. The Go code's `getNotificationFromResponse` is a TODO stub returning `nil`.
- [ ] **Missing error handling via try/catch/finally equivalent**: The Java code wraps the logic in try/catch/finally, returning `Mono.error(e)` on exception while always releasing the buffer in `finally`. The Go code does not recover from panics or handle errors from the RPC call.
- [ ] **Returns a zero-value `TurmsNotification` instead of a properly constructed response**: The Go code returns `&notification` where `notification` is a zero-valued struct, which means all fields (code, timestamp, requestId, etc.) are zero/nil. The Java code always returns a notification with at minimum a timestamp, code, and requestId set.

## getNotificationFromResponse

- [ ] **Entire method is a stub returning `nil`**: The Java version validates that `response.code()` is non-null (throwing `IllegalArgumentException` otherwise), builds a `TurmsNotification` with `reason`, `data`, `timestamp`, `code` (business code), and `requestId`. The Go version does nothing and returns `nil`.

# SessionController.java
*Checked methods: deleteSessions(@QueryParam(required = false)*

## deleteSessions

- [ ] **Missing logic: When both `ids` and `ips` are empty, the Java code calls `sessionService.closeAllLocalSessions(closeReason)` to close ALL local sessions. The Go code does nothing in this case — it simply returns `0` with no error, silently dropping the "close all" behavior.**
- [ ] **Incorrect return value: The Go code returns `len(ids)` and `len(ips)` as the count of closed sessions, but the Java code returns the actual count from `sessionService.closeLocalSessions()` (i.e., the number of sessions actually closed). The Go code assumes the number of input IDs/IPs equals the number of closed sessions, which is incorrect — not every ID or IP necessarily has an active session.**
- [ ] **Incorrect IP conversion: The Go code converts IP strings to `[]byte` by simply casting the string to bytes (`[]byte(ip)`), which produces the UTF-8 byte representation of the string (e.g., `"127.0.0.1"` → `[]byte{49, 50, 55, ...}`). The Java code uses `InetAddressUtil::ipStringToBytes`, which parses the IP string into its actual 4-byte (IPv4) or 16-byte (IPv6) binary representation (e.g., `"127.0.0.1"` → `[]byte{127, 0, 0, 1}`). This is a functional bug — the Go service will never match any sessions by IP.**
- [ ] **Missing CloseReason: The Go code passes `nil` as the close reason instead of constructing a `CloseReason` equivalent to `CloseReason.get(SessionCloseStatus.DISCONNECTED_BY_ADMIN)` (noted as a TODO but still a bug).**
- [ ] **Incorrect branch logic when both `ids` and `ips` are non-empty: In Java, when both are provided, it uses `Mono.zip(..., Integer::sum)` to run both close operations concurrently and sum their **actual** return values. The Go code runs them sequentially (not necessarily wrong for Go), but the summed count is wrong because it uses input lengths instead of actual closed session counts (see second bug above).**

# SessionClientController.java
*Checked methods: handleDeleteSessionRequest(UserSessionWrapper sessionWrapper), handleCreateSessionRequest(UserSessionWrapper sessionWrapper, CreateSessionRequest createSessionRequest)*

Now I have all the information needed for a thorough comparison. Here are the findings:

## HandleDeleteSessionRequest

- [ ] **Missing error handling/logging for session close failure**: The Java version subscribes to the `closeLocalSession` Mono and logs an error with the user ID if the close operation fails (`t -> LOGGER.error("Caught an error while closing the session with the user ID: " + userId, t)`). The Go version calls `CloseLocalSession` but ignores the returned error entirely — the error is not logged or handled.

## HandleCreateSessionRequest

- [ ] **Location data passed incorrectly — type mismatch causes silent failure**: The Java version constructs a `Location` object with `longitude`, `latitude`, optional `timestamp` (as `Date`), and `details` map, then passes it to `handleLoginRequest`. The Go version extracts only `req.Location.Details` (a `map[string]string`) and passes it as the `location` parameter. However, `HandleLoginRequest` expects `location any` and internally does a type assertion to `*protocol.UserLocation`. Since a `map[string]string` is passed instead of `*protocol.UserLocation`, the type assertion `loc, ok := location.(*protocol.UserLocation)` will always fail, causing location data (longitude, latitude) to never be stored in the session or persisted to Redis via `UpsertUserLocation`. The fix should pass `req.Location` (the `*protocol.UserLocation` pointer) directly.

- [ ] **Location timestamp and details fields lost in processing**: Related to the above, even if the location were passed correctly as `*protocol.UserLocation`, the downstream code (`addOnlineDeviceIfAbsent` and `TryRegisterOnlineUser`) only uses `Longitude` and `Latitude` from the `UserLocation`. The Java version passes a rich `Location` BO containing `timestamp` (as `Date`) and `details` (as `Map<String,String>`) which may be used elsewhere in the Java codebase. In the Go version, the `timestamp` from `UserLocation` and the `details` map are never extracted or passed — only `Details` is incorrectly extracted as the entire location value. The `sessionbo.UserLocation` struct only has `Longitude` and `Latitude` fields, missing `Timestamp` and `Details`.

- [ ] **DeviceType UNRECOGNIZED check is a hardcoded magic number**: The Java version checks `deviceType == DeviceType.UNRECOGNIZED`, which is a protobuf-generated sentinel value for unknown enum values. The Go version uses a hardcoded `protocol.DeviceType(5)` with a comment "Assuming 5 is UNKNOWN". In Go protobuf, `DeviceType_UNKNOWN` is explicitly value 5, so `protocol.DeviceType(5)` is equivalent to `protocol.DeviceType_UNKNOWN`, making the check `deviceType == protocol.DeviceType_UNKNOWN` which then sets `deviceType = protocol.DeviceType_UNKNOWN` — a no-op. The intent should be to check for out-of-range/unrecognized values (e.g., any value not in 0–5), but Go protobuf doesn't generate an `UNRECOGNIZED` constant the same way Java does. This check is effectively dead code.

- [ ] **Session establishment timeout logic is hardcoded to `false` (dead code)**: The Java version checks `sessionEstablishTimeout == null || sessionEstablishTimeout.cancel()` to determine if the session establishment timed out during the login process. If the timeout already fired (cancel returns false), it closes the session with `LOGIN_TIMEOUT` and returns an error result. The Go version hardcodes `isTimeout := false` (line 88), making this branch unreachable dead code. If a timeout mechanism exists in the Go connection layer, it is not wired into this check.

- [ ] **Connection alive check is hardcoded to `true` (dead code)**: The Java version checks `sessionWrapper.getConnection().isConnected()` to verify the TCP/WebSocket connection is still open before committing the session. If the connection dropped during login, it cleans up with `closeLocalSession` and returns empty. The Go version hardcodes `isConnectionAlive := true` (line 95), making the connection-drop cleanup path (lines 111–112) unreachable dead code.

- [ ] **Error handling for `InvokeGoOnlineHandlers` is missing**: The Java version subscribes to `invokeGoOnlineHandlers` and logs errors: `.subscribe(null, t -> LOGGER.error(ERROR_INVOKE_GO_ONLINE, t))`. The Go version calls `InvokeGoOnlineHandlers` synchronously and ignores any potential error or panic. If the Go `InvokeGoOnlineHandlers` can fail, the error is silently swallowed.

- [ ] **Error from `OnSessionEstablished` is silently ignored**: The Go version calls `c.sessionService.OnSessionEstablished(ctx, userSessionsManager, session.DeviceType)` but does not check or handle any error return value. The Java version uses a fire-and-forget pattern for similar async calls, but the Go version should at least log errors for observability.

- [ ] **`GetUserSessionsManager` may return nil without check**: The Go code calls `c.sessionService.GetUserSessionsManager(ctx, userID)` and immediately passes the result to `OnSessionEstablished` and `InvokeGoOnlineHandlers`. If the manager is nil (e.g., due to a race condition where the session was concurrently removed), this would cause a nil pointer panic in downstream methods.

# UserPermissionInfo.java
*Checked methods: UserPermissionInfo(...)*

## UserPermissionInfo (constructor / static fields)

- [ ] **`GrantedWithAllPermissions` uses wrong constructor**: In Java, `GRANTED_WITH_ALL_PERMISSIONS` uses the two-arg constructor `new UserPermissionInfo(ResponseStatusCode.OK, TurmsRequestTypePool.ALL)`, passing the full set of all permissions. In Go, `GrantedWithAllPermissions` calls `NewUserPermissionInfoCodeOnly` which creates an **empty** permissions map. This means the Go version has zero permissions instead of all permissions.

- [ ] **Missing `TurmsRequestTypePool.ALL` equivalent**: The Go code does not import or reference any equivalent of `TurmsRequestTypePool.ALL`. The `Permissions` field is typed as `map[any]bool` but there is no constant or variable providing the "all permissions" set to pass to `NewUserPermissionInfo` for the `GrantedWithAllPermissions` static field.

- [ ] **`Permissions` field typed as `map[any]bool` instead of a typed set**: Java uses `Set<TurmsRequest.KindCase>` for type safety. The Go version uses `map[any]bool`, which loses type safety — any key type can be inserted, not just `TurmsRequest.KindCase` values. This is a design divergence that could lead to runtime bugs where incorrect keys are stored.

# UserRepository.java
*Checked methods: findPassword(Long userId), isActiveAndNotDeleted(Long userId)*

Now I have all the information needed to compare. Let me analyze both methods carefully.

## FindPassword

**Java** (`findPassword(Long userId)`):
- Creates a filter: `eq(ID, userId)`
- Creates query options: `.include(User.Fields.PASSWORD)` — this is a **projection** that only fetches the `password` field
- Returns `Mono<User>` — the returned User object will **only** have the password field populated; all other fields will be null/default

**Go** (`FindPassword(ctx, userID)`):
- Calls `r.FindByID(ctx, userID)` which fetches the **entire user document** (no projection)
- Returns `*string` (the password pointer) instead of `*User`

The difference in return type (`*User` vs `*string`) is an acceptable API design choice — the Go version extracts just the password string rather than returning the full user object. The caller only needs the password, so this is fine.

However, there is a functional difference: the Java version uses a projection to only fetch the `password` field from MongoDB (bandwidth/latency optimization), while the Go version fetches all fields. This is a performance concern but not a correctness bug.

There's also a subtle logic issue: the Go code returns `nil` when `user.Password == ""`. But in Java, the password field could legitimately be `null` (no password set) vs an empty string. The Go code treats an empty string as "no password" and returns nil, but the Java version would still return the User object with a null password field. The Go code conflates "password is empty string" with "no password found" — both return nil. This changes behavior if a password could ever be an empty string in the database.

## IsActiveAndNotDeleted

**Java** (`isActiveAndNotDeleted(Long userId)`):
- Creates a filter with 3 conditions: `eq(ID, userId) AND eq(IS_ACTIVE, true) AND eq(DELETION_DATE, null)`
- Uses `mongoClient.exists()` — this is a MongoDB **existence check** that runs entirely on the server. It checks if a document matching all 3 criteria exists.

**Go** (`IsActiveAndNotDeleted(ctx, userID)`):
- Calls `r.FindByID(ctx, userID)` which fetches the **entire user document** with only the ID filter
- Then checks `user.IsActive && user.DeletionDate == nil` in Go code

The Go version is functionally equivalent in terms of the logical conditions checked — it verifies the same three things (user exists, is active, and has no deletion date). However:
1. It fetches the entire document instead of doing a server-side existence check (performance concern, not a correctness bug)
2. The logic is correct: if user is nil (not found), returns false. If user exists, returns `IsActive && DeletionDate == nil`.

This is correct behavior — no logical bug.

---

## FindPassword

- [ ] **Missing projection**: The Java version uses `.include(User.Fields.PASSWORD)` to create a MongoDB projection that only fetches the `password` field from the database. The Go version calls `FindByID` which fetches all fields of the user document. This is a performance regression — the Go code transfers unnecessary data from MongoDB on every password lookup.

- [ ] **Empty string vs nil conflation**: The Go code returns `nil` when `user.Password == ""` (line 289-291). In the Java version, the User object is returned regardless of the password value — it's up to the caller to inspect the password field. A user record with a `null` password (no password set) and a user record that doesn't exist are distinguishable in Java (empty Mono vs Mono with User having null password). In Go, both `user not found` and `user found with empty password` return `nil, nil`, losing the distinction. If the caller needs to differentiate between "user doesn't exist" and "user exists but has no password", this is a behavioral difference.

# HttpSessionIdentityAccessManager.java
*Checked methods: verifyAndGrant(UserLoginInfo userLoginInfo)*

Now I have the full picture. Let me analyze the `verifyAndGrant` method across all three relevant classes: the orchestrator `SessionIdentityAccessManager`, the `PasswordSessionIdentityAccessManager`, and the `HttpSessionIdentityAccessManager`.

## HttpSessionIdentityAccessManager.VerifyAndGrant

- [ ] **Entire method is a stub returning `(nil, nil)` instead of implementing the HTTP authentication flow.** The Java version: (1) sends an HTTP request with the serialized `UserLoginInfo` as JSON body, (2) validates the response status code against `httpAuthenticationExpectedStatusCodes`, (3) validates expected response headers, (4) parses the response body as a JSON map, (5) calls `PolicyDeserializer.parse(map)` to get a `Policy`, (6) validates the response body fields against `httpAuthenticationExpectedBodyFields` using loose comparison, and (7) returns `UserPermissionInfo` with the allowed request types from `policyManager.findAllowedRequestTypes(policy)`. The Go version returns `(nil, nil)` with none of this logic.

- [ ] **Missing all struct fields.** The Java `HttpSessionIdentityAccessManager` holds fields: `httpIdentityAccessManagementClient`, `httpIdentityAccessManagementHttpMethod`, `httpAuthenticationExpectedStatusCodes`, `httpAuthenticationExpectedHeaders`, `httpAuthenticationExpectedBodyFields`, and `policyManager`. The Go struct `HttpSessionIdentityAccessManager` is empty (`struct{}`).

## PasswordSessionIdentityAccessManager.VerifyAndGrant

- [ ] **Password comparison is a plain string equality check instead of using a proper password encoder.** The Java version calls `userService.authenticate(userId, password)` which uses Spring Security's `PasswordEncoder` (typically bcrypt). The Go version does a direct `user.Password != *loginInfo.Password` comparison, which will fail for any bcrypt-hashed password stored in the database.

- [ ] **User lookup uses `FindUser` and checks `IsActive` locally instead of calling `isActiveAndNotDeleted`.** The Java version calls `userService.isActiveAndNotDeleted(userId)` which is a single service call that checks both active and not-deleted status. The Go version calls `FindUser` (fetching the full user record) and then only checks `IsActive`, missing the "not deleted" check. A deleted user could still authenticate.

- [ ] **Missing nil password handling difference.** The Go version returns `LOGIN_AUTHENTICATION_FAILED` when password is nil/empty, but the Java version passes the password directly to `userService.authenticate()` and lets it handle nil. If `password` is nil but the user account has no password requirement (edge case in Java's design), the Go version would incorrectly reject while Java might not.

## SessionIdentityAccessManager.VerifyAndGrant

- [ ] **Admin user ID check uses `== 0` instead of comparing against `AdminConst.ADMIN_REQUESTER_ID`.** The Java version explicitly compares `userId.equals(AdminConst.ADMIN_REQUESTER_ID)`. The Go version hardcodes `loginInfo.UserID == 0` with a comment, but this should use the Go equivalent of the admin constant for correctness and clarity.

- [ ] **Plugin-based authentication is entirely stubbed out with a TODO.** The Java version checks `pluginManager.hasRunningExtensions(UserAuthenticator.class)` and, if true, invokes `authenticate` extension points sequentially, falling back to the default handler via `switchIfEmpty`. The Go version has this entirely commented out with a TODO. This means plugins cannot intercept authentication.

- [ ] **`GRANTED_WITH_ALL_PERMISSIONS` uses `map[any]bool{}` instead of the proper "all permissions" sentinel.** The Java version returns `GRANTED_WITH_ALL_PERMISSIONS_MONO` which is a predefined constant. The Go version constructs `map[any]bool{}` (an empty map), which may not be semantically equivalent to "all permissions granted." The Java constant likely represents a special sentinel value that means "all request types are allowed," while an empty Go map could be interpreted as "no permissions."

- [ ] **Go returns `(nil, nil)` instead of `LOGIN_AUTHENTICATION_FAILED` when support is nil.** At the end of `VerifyAndGrant` on the orchestrator, if `m.support` is nil, Go returns `(nil, nil)`, while Java would never reach that path (the support is always initialized in the constructor via a switch). This is a potential nil-pointer crash for any caller that doesn't check for nil `UserPermissionInfo`.

# JwtSessionIdentityAccessManager.java
*Checked methods: verifyAndGrant(UserLoginInfo userLoginInfo)*

## `verifyAndGrant` (on `JwtSessionIdentityAccessManager`)

- [ ] **Method body is completely unimplemented**: The Go version returns `nil, nil` without any logic. The Java version performs: (1) JWT blank check, (2) JWT decode + signature verification, (3) subject claim validation, (4) subject-to-userId match, (5) `expiresAt`/`notBefore` time validation, (6) custom claims loose-comparison check against expected claims, (7) policy deserialization from custom claims, and (8) returning `UserPermissionInfo` with allowed request types from `policyManager.findAllowedRequestTypes(policy)`. None of this logic exists in the Go code.

- [ ] **Missing fields on `JwtSessionIdentityAccessManager` struct**: The Java class has three fields: `jwtManager`, `policyManager`, and `jwtAuthenticationExpectedCustomPayloadClaims`. The Go struct has none of these fields.

- [ ] **Missing constructor/initialization**: The Java constructor initializes `jwtManager` (with all algorithm properties), `policyManager`, and `jwtAuthenticationExpectedCustomPayloadClaims`. No equivalent constructor or factory function exists in Go.

- [ ] **Missing `NewUserPermissionInfo` with `allowedRequestTypes`**: The Java code returns `UserPermissionInfo` populated with `policyManager.findAllowedRequestTypes(policy)`. The Go code returns `nil, nil`, which will likely cause a nil pointer dereference at the call site since callers expect a valid `*bo.UserPermissionInfo`.

# LdapSessionIdentityAccessManager.java
*Checked methods: verifyAndGrant(UserLoginInfo userLoginInfo)*

## LdapSessionIdentityAccessManager.VerifyAndGrant

- [ ] **Missing password blank check**: The Java version checks `StringUtil.isBlank(password)` at the start and returns `LOGIN_AUTHENTICATION_FAILED_MONO` if the password is blank. The Go version returns `nil, nil` immediately without any logic, skipping this check entirely.
- [ ] **Missing user search filter replacement**: The Java version replaces the `SEARCH_FILTER_PLACEHOLDER_USER_ID` placeholder in `userSearchFilter` with the user ID string to build the LDAP search filter. The Go version has no such logic.
- [ ] **Missing admin LDAP search**: The Java version performs an LDAP search using `adminLdapClient.search(baseDn, Scope.WHOLE_SUBTREE, DerefAliases.ALWAYS, 2, 0, false, SearchRequest.NO_ATTRIBUTES, filter)` to find the user's LDAP entry. The Go version has no LDAP client or search logic at all.
- [ ] **Missing entry count checks**: The Java version checks: (1) if 0 entries found, returns `LOGGING_IN_USER_NOT_ACTIVE_MONO`; (2) if more than 1 entry found, returns an error `SERVER_INTERNAL_ERROR` with a descriptive message. The Go version has none of this.
- [ ] **Missing DN extraction and user authentication**: The Java version extracts `objectName` (DN) from the first search result entry, then calls `authenticateUser(objectName, password)` which binds to the LDAP server with the user's DN and password via a task scheduler. The Go version has no authentication logic.
- [ ] **Missing clientBindTaskScheduler**: The Java version uses a `TaskScheduler` to serialize LDAP bind requests (to comply with RFC 4511 section 4.2.1). The Go `LdapSessionIdentityAccessManager` struct has no fields at all — no LDAP clients, no base DN, no search filter, no scheduler.
- [ ] **Entire method is a stub**: The Go `LdapSessionIdentityAccessManager.VerifyAndGrant` simply returns `nil, nil` without implementing any of the Java logic. It should implement the full LDAP verification flow: blank password check → admin search → entry count validation → user bind authentication → return appropriate `UserPermissionInfo`.

## LdapSessionIdentityAccessManager (structural)

- [ ] **Missing struct fields**: The Go `LdapSessionIdentityAccessManager` struct is empty, but the Java version holds `baseDn`, `userSearchFilter`, `adminLdapClient`, `userLdapClient`, and `clientBindTaskScheduler`. None of these are present in the Go struct.
- [ ] **Missing constructor/initialization logic**: The Java constructor validates the search filter contains the placeholder, creates admin and user LDAP clients, performs startup health checks (admin bind + search, user bind) with a 1-minute timeout, and initializes the task scheduler. The Go version has no constructor or any of this initialization.

# NoopSessionIdentityAccessManager.java
*Checked methods: verifyAndGrant(UserLoginInfo userLoginInfo)*

Now I have a clear picture. Let me verify the key details:

**Java `NoopSessionIdentityAccessManager.verifyAndGrant`:** Returns `GRANTED_WITH_ALL_PERMISSIONS_MONO`, which wraps `new UserPermissionInfo(ResponseStatusCode.OK, TurmsRequestTypePool.ALL)` — that is, a `UserPermissionInfo` with status `OK` and a permissions set containing ALL request types.

**Go `NoopSessionIdentityAccessManager.VerifyAndGrant`:** Returns `bo.NewUserPermissionInfo(constant.ResponseStatusCode_OK, nil)` — that is, a `UserPermissionInfo` with status `OK` and `nil` permissions.

The Go `GrantedWithAllPermissions` uses `NewUserPermissionInfoCodeOnly` which creates an empty permissions map (`make(map[any]bool)`), NOT populated with all request types. And the Noop variant passes `nil` rather than even an empty map.

## NoopSessionIdentityAccessManager.VerifyAndGrant

- [ ] **Missing permissions set**: The Java version returns `GRANTED_WITH_ALL_PERMISSIONS` which contains `TurmsRequestTypePool.ALL` (the full set of all request types) as the permissions. The Go version passes `nil` for permissions. Even if `nil` is treated as "all permissions" downstream, this is inconsistent with the Java contract where a populated set is explicitly provided. The Go code should pass a set populated with all TurmsRequest types, or at minimum use the existing `bo.GrantedWithAllPermissions` constant (though that itself is also wrong — see next point).

- [ ] **`GrantedWithAllPermissions` global constant is incorrect**: In `user_permission_info.go:32`, `GrantedWithAllPermissions` is defined as `NewUserPermissionInfoCodeOnly(constant.ResponseStatusCode_OK)` which creates a permissions set via `make(map[any]bool)` — an **empty** map, not one populated with all request types like Java's `TurmsRequestTypePool.ALL`. The NoopSessionIdentityAccessManager should ideally use this constant but the constant itself needs to be fixed to include all permissions.

# PasswordSessionIdentityAccessManager.java
*Checked methods: verifyAndGrant(UserLoginInfo userLoginInfo), updateGlobalProperties(TurmsProperties properties)*

Now I have a complete picture. Let me compile the findings.

## verifyAndGrant (PasswordSessionIdentityAccessManager)

- [x] **Wrong error code when user is not active**: When `!user.IsActive`, the Go code returns `ResponseStatusCode_LOGIN_AUTHENTICATION_FAILED` (2000), but the Java code returns `ResponseStatusCode_LOGGING_IN_USER_NOT_ACTIVE` (2002). In Java, `isActiveAndNotDeleted` returns false, which maps to `LOGGING_IN_USER_NOT_ACTIVE_MONO` — a different status code.

- [x] **Missing "deleted" check**: The Java code calls `isActiveAndNotDeleted(userId)` which checks both conditions. The Go code only checks `!user.IsActive` but does not explicitly check if the user is deleted. If the `FindUser` query returns deleted users, the logic would differ.

- [x] **Password comparison is plain string equality instead of using PasswordManager**: The Java code uses `passwordManager.matchesUserPassword(rawPassword, user.getPassword())` which handles encoding (e.g., bcrypt). The Go code does a direct `user.Password != *loginInfo.Password` string comparison, which will fail for any hashed password storage.

- [x] **Granted response returns nil permissions instead of all permissions**: On success, the Go code returns `NewUserPermissionInfo(constant.ResponseStatusCode_OK, nil)` (nil permissions map), but the Java code returns `GRANTED_WITH_ALL_PERMISSIONS` which has `TurmsRequestTypePool.ALL` as the permissions set — meaning all request types are permitted. Returning nil means no permissions are granted.

- [x] **Finds full user record instead of separate targeted queries**: The Java code makes two separate, optimized calls: `isActiveAndNotDeleted(userId)` and `authenticate(userId, password)`. The Go code calls `FindUser(ctx, loginInfo.UserID)` to load the entire user entity. While functionally similar for this simple case, it's a behavioral difference that loads more data than needed and combines the two-step logic into one.

## updateGlobalProperties (PasswordSessionIdentityAccessManager)

- [ ] **Method body is completely empty**: The Go `PasswordSessionIdentityAccessManager.UpdateGlobalProperties` is a no-op (empty body), while the Java version contains critical logic: it reads `enableIdentityAccessManagement` from properties, checks if the user service is enabled, logs an error and returns `false` if trying to enable a previously disabled service, or returns `enableIdentityAccessManagement` otherwise.

## updateGlobalProperties (SessionIdentityAccessManager orchestrator)

- [ ] **Does not update `enableIdentityAccessManagement` from the support's return value**: The Java orchestrator's `updateGlobalProperties` sets `enableIdentityAccessManagement` from the return value of `sessionIdentityAccessManagementSupport.updateGlobalProperties(properties)`. The Go version calls `m.support.UpdateGlobalProperties(properties)` but the interface returns void — it has no way to propagate the enabled/disabled state back to the orchestrator's `enableIdentityAccessManagement` field. The `UpdateGlobalProperties` on `SessionIdentityAccessManager` also has a comment `// m.enableIdentityAccessManagement = ...` suggesting this assignment was never implemented.

# SessionIdentityAccessManager.java
*Checked methods: verifyAndGrant(int version, Long userId, @Nullable String password, DeviceType deviceType, @Nullable Map<String, String> deviceDetails, @Nullable UserStatus userStatus, @Nullable Location location, @Nullable String ip)*

Now I have all the information needed for a thorough comparison. Let me analyze the Go code against the Java code for the `verifyAndGrant` method.

## verifyAndGrant

- [ ] **Admin user ID check uses `0` literal instead of constant, and uses `==` instead of `.equals()`**: In Java, `userId.equals(AdminConst.ADMIN_REQUESTER_ID)` compares `Long` objects. Since `AdminConst.ADMIN_REQUESTER_ID` is `0L`, the Go code's `loginInfo.UserID == 0` is functionally equivalent. However, the Go code passes the *entire* `loginInfo` object as the method parameter rather than the individual parameters the Java method accepts. The Java `verifyAndGrant` takes 8 individual parameters (`int version`, `Long userId`, `String password`, `DeviceType deviceType`, `Map<String, String> deviceDetails`, `UserStatus userStatus`, `Location location`, `String ip`), while the Go version takes `(ctx context.Context, loginInfo *bo.UserLoginInfo)`. This is an acceptable API adaptation for Go, but the caller must ensure all fields are populated on `loginInfo` before calling.

- [ ] **Missing `GRANTED_WITH_ALL_PERMISSIONS` returns full permission set, but Go returns empty permissions map**: In Java, when identity access management is disabled, it returns `GRANTED_WITH_ALL_PERMISSIONS_MONO` which maps to `new UserPermissionInfo(ResponseStatusCode.OK, TurmsRequestTypePool.ALL)` — meaning the permissions set contains ALL request types. In Go, line 48 returns `bo.NewUserPermissionInfo(constant.ResponseStatusCode_OK, map[any]bool{})` which creates an **empty** permissions map. The Go predefined `GrantedWithAllPermissions` variable (line 32 of user_permission_info.go) also uses `NewUserPermissionInfoCodeOnly` which creates an empty map. This is a behavioral difference: Java grants all permissions, Go grants none.

- [ ] **Missing plugin authenticator fallback logic**: The Java code checks `pluginManager.hasRunningExtensions(UserAuthenticator.class)` and if true, invokes plugin authenticators sequentially. If a plugin authenticator returns `true`, it returns `GRANTED_WITH_ALL_PERMISSIONS`; if `false`, `LOGIN_AUTHENTICATION_FAILED`. The plugin result is used via `authenticate.switchIfEmpty(defaultVerifyAndGrantHandler)`, meaning if no plugin produces a result, it falls back to the default handler. The Go code has this as a TODO comment and skips it entirely. While this is noted as TODO, the Go code currently always falls through to `m.support.VerifyAndGrant()`, which means plugin-based authentication is completely non-functional.

- [ ] **`GRANTED_WITH_ALL_PERMISSIONS` permissions set is empty instead of containing all request types**: The Java `GRANTED_WITH_ALL_PERMISSIONS` constant is initialized with `TurmsRequestTypePool.ALL` (the complete set of all TurmsRequest kind cases). The Go `GrantedWithAllPermissions` is initialized with `NewUserPermissionInfoCodeOnly(constant.ResponseStatusCode_OK)` which uses `make(map[any]bool)` — an empty set. This means even when authentication succeeds with "all permissions", no actual permissions are granted.

- [ ] **`LOGIN_AUTHENTICATION_FAILED` return uses `NewUserPermissionInfo` with `nil` permissions instead of `NewUserPermissionInfoCodeOnly`**: In Java, `LOGIN_AUTHENTICATION_FAILED` is `new UserPermissionInfo(ResponseStatusCode.LOGIN_AUTHENTICATION_FAILED)` which uses the secondary constructor that sets `Collections.emptySet()`. In Go, the admin check on line 44 uses `bo.NewUserPermissionInfo(constant.ResponseStatusCode_LOGIN_AUTHENTICATION_FAILED, nil)` which passes `nil` for permissions. While functionally a nil map and empty map behave differently in Go (nil map is read-only, empty map can be written to), neither should be written to after return, so this is a minor issue. However, the predefined `LoginAuthenticationFailed` variable exists and should be used for consistency.

# SessionService.java
*Checked methods: destroy(), handleHeartbeatUpdateRequest(UserSession session), handleLoginRequest(int version, @NotNull ByteArrayWrapper ip, @NotNull Long userId, @Nullable String password, @NotNull DeviceType deviceType, @Nullable Map<String, String> deviceDetails, @Nullable UserStatus userStatus, @Nullable Location location, @Nullable String ipStr), closeLocalSessions(@NotNull List<byte[]> ips, @NotNull CloseReason closeReason), closeLocalSessions(@NotNull byte[] ip, @NotNull CloseReason closeReason), closeLocalSession(@NotNull Long userId, @NotNull @ValidDeviceType DeviceType deviceType, @NotNull SessionCloseStatus closeStatus), closeLocalSession(@NotNull Long userId, @NotNull @ValidDeviceType DeviceType deviceType, @NotNull CloseReason closeReason), closeLocalSession(@NotNull Long userId, @NotEmpty Set<@ValidDeviceType DeviceType> deviceTypes, @NotNull CloseReason closeReason), closeLocalSessions(@NotNull Set<Long> userIds, @NotNull CloseReason closeReason), authAndCloseLocalSession(@NotNull Long userId, @NotNull DeviceType deviceType, @NotNull CloseReason closeReason, int sessionId), closeAllLocalSessions(@NotNull CloseReason closeReason), closeLocalSession(Long userId, SessionCloseStatus closeStatus), closeLocalSession(Long userId, CloseReason closeReason), getSessions(Set<Long> userIds), authAndUpdateHeartbeatTimestamp(long userId, @NotNull @ValidDeviceType DeviceType deviceType, int sessionId), tryRegisterOnlineUser(int version, @NotNull Set<TurmsRequest.KindCase> permissions, @NotNull ByteArrayWrapper ip, @NotNull Long userId, @NotNull DeviceType deviceType, @Nullable Map<String, String> deviceDetails, @Nullable UserStatus userStatus, @Nullable Location location), getUserSessionsManager(@NotNull Long userId), getLocalUserSession(@NotNull Long userId, @NotNull DeviceType deviceType), getLocalUserSession(ByteArrayWrapper ip), countLocalOnlineUsers(), onSessionEstablished(@NotNull UserSessionsManager userSessionsManager, @NotNull @ValidDeviceType DeviceType deviceType), addOnSessionClosedListeners(Consumer<UserSession> onSessionClosed), invokeGoOnlineHandlers(@NotNull UserSessionsManager userSessionsManager, @NotNull UserSession userSession)*

Now I have a thorough understanding of both files. Let me compile the bug report.

## destroy

- [ ] The Go `Destroy` method does not destroy the heartbeat manager before closing sessions. The Java `destroy()` first calls `heartbeatManager.destroy()` and only then proceeds to `closeAllLocalSessions(closeReason)`. The Go code skips the heartbeat manager shutdown entirely.
- [ ] The Go `Destroy` does not construct a proper `CloseReason` from `SessionCloseStatus.SERVER_CLOSED`. It passes `nil` as `closeReason` to `CloseAllLocalSessions`, whereas Java creates `CloseReason.get(SessionCloseStatus.SERVER_CLOSED)`.

## handleHeartbeatUpdateRequest

- [ ] No bugs found. The Go implementation matches the Java logic.

## handleLoginRequest

- [ ] The Go code does not check the `authenticationCode` (equivalent of `statusCode`) from `permissionInfo` before proceeding. In Java, after `verifyAndGrant`, it checks `statusCode == ResponseStatusCode.OK` and only then calls `tryRegisterOnlineUser`; otherwise it returns `Mono.error(ResponseException.get(statusCode))`. The Go code unconditionally calls `TryRegisterOnlineUser` regardless of the authentication code returned by `VerifyAndGrant`.
- [ ] The Go code passes `ip` (raw bytes) and `ipStr` separately, but the Java code passes `ip` as `ByteArrayWrapper` into `tryRegisterOnlineUser`. In Go, `ip` is only used in `addOnlineDeviceIfAbsent` to construct `net.IP(ip)`, while the Java also uses `ip` for the `ipToSessions` map key. This is not a logic bug per se but the `ip` parameter is not passed to the `ipToSessions` registration in the same way (it goes through `RegisterSession` which uses `session.IP`).
- [ ] The Go `handleLoginRequest` passes `location` as `any` type to `TryRegisterOnlineUser`, but in Java `location` is passed as `@Nullable Location location`. The `location` parameter is not included in the `UserLoginInfo` in Go (it's commented out), while in Java `location` is part of `verifyAndGrant`.

## closeLocalSessions(List<byte[]> ips, CloseReason closeReason)

- [ ] The Go `CloseLocalSessionsByIp` does not return the count of closed sessions. The Java method returns `Mono<Integer>` with the total count of closed sessions. The Go method returns `error` only.
- [ ] The Go method iterates IPs and for each IP's sessions calls `UnregisterSession` with the session's connection, but the Java version calls `closeLocalSession(userId, DeviceTypeUtil.ALL_AVAILABLE_DEVICE_TYPES_SET, closeReason)` which closes ALL device types for that user (not just the one from the IP mapping). The Go version only closes the specific device type session found via the IP, which is a behavioral difference.
- [ ] The Go method does not validate that `ips` or `closeReason` are non-nil, whereas Java validates both parameters.
- [ ] The Go method does not handle the empty `ips` case as an early return (though it naturally does nothing).

## closeLocalSessions(byte[] ip, CloseReason closeReason)

- [ ] The Go implementation uses `CloseLocalSessionsByIp` for both single-IP and multi-IP variants. The Java single-IP version looks up sessions via `ipToSessions.get(new ByteArrayWrapper(ip))`, iterates them, and calls `closeLocalSession(userId, ALL_AVAILABLE_DEVICE_TYPES_SET, closeReason)` for each userId. The Go version calls `UnregisterSession` which only closes the specific device type session, not all device types for that user.
- [ ] The Go version does not return the count of closed sessions. Java returns `Mono<Integer>` with the count.
- [ ] The Go version does not aggregate errors with `Mono.whenDelayError` semantics — errors from closing one session may prevent closing others.

## closeLocalSession(Long userId, DeviceType deviceType, SessionCloseStatus closeStatus)

- [ ] The Go code does not have a dedicated method for this signature. It uses `CloseLocalSession(ctx, userId, []DeviceType{deviceType}, closeReason)` via the RPC handler registration, but there is no separate method matching this exact Java signature that converts `SessionCloseStatus` to `CloseReason`.

## closeLocalSession(Long userId, DeviceType deviceType, CloseReason closeReason)

- [ ] Same as above — no dedicated method for this signature. The Go code folds this into the multi-device-type `CloseLocalSession`.

## closeLocalSession(Long userId, Set<DeviceType> deviceTypes, CloseReason closeReason)

- [ ] The Go `CloseLocalSession` does not call `userStatusService.removeStatusByUserIdAndDeviceTypes(userId, deviceTypes)` before closing sessions locally. The Java code calls this Redis status removal first ("Don't close the session first and then remove the session status in Redis because it will make trouble if a client logins again while the session status in Redis hasn't been removed"). This is a critical missing step.
- [ ] The Go method does not call `sessionLocationService.removeUserLocation(userId, deviceType)` for each session being closed. The Java code does this when `sessionLocationService.isLocationEnabled()`.
- [ ] The Go method does not call `notifyOnSessionClosedListeners` inside the session-closing loop in the same way. The Java code calls `notifyOnSessionClosedListeners(session)` only when `wasSessionOpen` is true (i.e., the session was actually open before closing). The Go code invokes listeners unconditionally after closing each session's connection.
- [ ] The Go method does not return the count of closed sessions. Java returns `Mono<Integer>` with the count.
- [ ] The `removeSessionsManagerIfEmpty` logic in Java only removes the manager from the map when `manager.countSessions() == 0`, but always invokes the plugin goOffline handler. In Go, `RemoveIfEmpty` returns non-nil only when the manager is empty and was removed, so `InvokeGoOfflineHandlers` is only called when sessions are empty — matching the Java behavior for removal but the Java always invokes the plugin regardless of emptiness.

## closeLocalSessions(Set<Long> userIds, CloseReason closeReason)

- [ ] The Go `CloseLocalSessionsByUserIds` calls `UnregisterSession` per device type, while the Java version calls `closeLocalSession(userId, closeReason)` which uses `ALL_AVAILABLE_DEVICE_TYPES_SET`. The Go version gets sessions then unregisters each individually, which is functionally similar but uses a different code path.
- [ ] The Go method does not return the count of closed sessions. Java returns `Mono<Integer>` with the aggregated count.
- [ ] The Go method does not validate that `userIds` and `closeReason` are non-nil. Java validates both.
- [ ] The Go method does not handle the empty `userIds` early return.

## authAndCloseLocalSession

- [ ] The Go method does not validate parameters (userId, deviceType, closeReason) are non-nil. Java validates all three.
- [ ] The Go method does not return the count of closed sessions. Java returns `Mono<Integer>` with the count.
- [ ] The Go method does not call `userStatusService.removeStatusByUserIdAndDeviceTypes` before closing (same issue as `closeLocalSession`).
- [ ] The Go method does not call `sessionLocationService.removeUserLocation`.

## closeAllLocalSessions

- [ ] The Go method does not validate `closeReason` is non-nil. Java validates it.
- [ ] The Go method does not return the count of closed sessions. Java returns `Mono<Integer>` with the aggregated count.
- [ ] The Go method gets `toClose` device types from `manager.GetAllSessions()` rather than `manager.getLoggedInDeviceTypes()` as in Java. These may differ if there are sessions that are registered but not yet "logged in."

## closeLocalSession(Long userId, SessionCloseStatus closeStatus)

- [ ] The Go code does not have a dedicated method for this. In Java, this calls `closeLocalSession(userId, ALL_AVAILABLE_DEVICE_TYPES_SET, CloseReason.get(closeStatus))`. The Go RPC handler partially covers this.

## closeLocalSession(Long userId, CloseReason closeReason)

- [ ] The Go code does not have a dedicated method for this. In Java, this calls `closeLocalSession(userId, ALL_AVAILABLE_DEVICE_TYPES_SET, closeReason)`.

## getSessions

- [ ] The Go method returns `nil` for users not found in the map, while the Java code returns `new UserSessionsInfo(userId, UserStatus.OFFLINE, Collections.emptyList())` for offline users. This means Go omits offline users entirely from the result, whereas Java includes them with OFFLINE status.
- [ ] The Go `UserSessionsInfo` does not include the user's status (`UserStatus`) field, while Java's `UserSessionsInfo(userId, manager.getUserStatus(), sessionInfos)` includes it.
- [ ] The Go `UserSessionInfo` is missing several fields compared to Java: `DeviceDetails`, `LastHeartbeatRequestTimestampMillis`, `LastRequestTimestampMillis`, `IsSessionOpen`, `IP` bytes. Java includes all these in each `UserSessionInfo`.

## authAndUpdateHeartbeatTimestamp

- [ ] The Go method does not check `!session.Conn.IsActive()` to verify the connection is NOT in a recovering state. Java checks `!session.getConnection().isConnectionRecovering()` before updating the heartbeat. The Go code only checks `session.ID == sessionId`.
- [ ] The Go method uses `s.GetUserSession` which goes through the sharded map, while Java uses `getUserSessionsManager(userId)` and then `manager.getSession(deviceType)`. This is functionally equivalent but the Java also validates the device type.

## tryRegisterOnlineUser

- [ ] The Go method does not validate input parameters (ip not null, deviceType not null, valid device type, userStatus not UNRECOGNIZED, userStatus not OFFLINE, location range). Java performs all these validations upfront.
- [ ] The Go method uses `sessionsStatus.OnlineDeviceTypeToSessionInfo` but the Java uses `sessionsStatus.getDeviceTypeToSessionInfo()` — the field naming suggests they map the same data but the Go code checks `info.IsActive` and `info.NodeID != s.nodeID` directly. In the Java version, it checks `sessionInfo.isActive()` and `!node.getLocalMemberId().equals(sessionInfo.getNodeId())`. The logic is the same.
- [ ] In the `isClosedSessionOnLocal` branch, the Go code calls `s.userStatusService.UpdateStatus` instead of `userStatusService.updateOnlineUserStatusIfPresent`. The Java method is `updateOnlineUserStatusIfPresent` which has "if present" semantics. The Go code calls `UpdateStatus` which may have different semantics.
- [ ] In the `isClosedSessionOnLocal` branch, when updating user status, the Go code checks `userStatus != 0` while Java checks `userStatus == null || existingUserStatus == userStatus` and skips the update if the userStatus is null or already matches. The Go code uses the zero-value of the UserStatus enum as the "null" check which may not be equivalent.
- [ ] In the `isClosedSessionOnLocal` branch, the Go code does not handle errors from status update or location upsert. Java wraps these in `onErrorComplete` handlers that log but don't propagate errors.
- [ ] The `closeSessionsWithConflictedDeviceTypes` method in Go iterates `sessionsStatus.OnlineDeviceTypeToSessionInfo` and checks `s.userSimultaneousLoginService.IsConflicted(deviceType, existingDeviceType)`. The Java version calls `userSimultaneousLoginService.getConflictedDeviceTypes(deviceType)` to get ALL conflicted types, then groups them by nodeId, then sends ONE RPC request per node with all conflicting device types. The Go version sends one RPC per conflicting device type per node, which is less efficient but also means the RPC `SetUserOfflineRequest` only contains a single device type rather than a set.
- [ ] The `closeSessionsWithConflictedDeviceTypes` in Go does not handle `ConnectionNotFound` errors with the isKnownMember fallback logic that Java has. In Java, if a `ConnectionNotFound` error occurs, it checks `node.getDiscoveryService().isKnownMember(nodeId)` — if the member is known, it returns the error; if unknown (dead node), it returns `true` to allow login. The Go code silently ignores errors.
- [ ] The `closeSessionsWithConflictedDeviceTypes` in Go returns `(true, nil)` on success but does not aggregate multiple RPC results. Java uses `PublisherUtil.areAllTrue(requests)` to ensure all RPCs succeed. The Go version returns `true` even if some RPCs fail silently.
- [ ] The `addOnlineDeviceIfAbsent` in Go does not pass `deviceDetails` (filtered by `deviceDetailsItemPropertiesList`) to `userStatusService.AddOnlineDevice`. Java filters `deviceDetails` based on `deviceDetailsItemPropertiesList` and passes the filtered `details` map to `addOnlineDeviceIfAbsent`. The Go code does not filter device details at all.
- [ ] The `addOnlineDeviceIfAbsent` in Go does not pass `closeIdleSessionAfterSeconds`, `expectedNodeId`, or `expectedDeviceTimestamp` to `userStatusService.AddOnlineDevice`. The Java code passes all of these to `userStatusService.addOnlineDeviceIfAbsent(userId, deviceType, details, userStatus, closeIdleSessionAfterSeconds, expectedNodeId, expectedDeviceTimestamp)`. The Go `AddOnlineDevice` only takes `(ctx, userId, deviceType, userStatus, nodeID, &now)`.
- [ ] The `addOnlineDeviceIfAbsent` in Go does not handle the case where `session == nil` from session creation (Java has a retry path: if `addSessionIfAbsent` returns null, it closes the existing session, cleans ipToSessions, and retries). The Go code does not have this retry/fallback logic.
- [ ] The `addOnlineDeviceIfAbsent` in Go does not set `Version` or `Permissions` on the created `UserSession`. Java calls `manager.addSessionIfAbsent(version, permissions, deviceType, deviceDetails, location)` which sets these fields on the session.
- [ ] The `addOnlineDeviceIfAbsent` in Go does not set `DeviceDetails` on the session. Java passes `deviceDetails` to `addSessionIfAbsent`.
- [ ] The `addOnlineDeviceIfAbsent` in Go always creates a new `UserSessionsManager` via `GetOrAdd`, but Java computes `finalUserStatus` (defaulting to `AVAILABLE` if null) and passes it to the constructor. The Go `NewUserSessionsManager` does not accept or set user status.

## getUserSessionsManager

- [ ] The Go method returns `any` instead of `*UserSessionsManager`. While this is a type signature issue rather than a logic bug, it loses type safety.
- [ ] The Go method does not validate that `userId` is non-null. Java calls `Validator.notNull(userId, "userId")`.

## getLocalUserSession(Long userId, DeviceType deviceType)

- [ ] The Go method does not validate that `userId` and `deviceType` are non-null. Java validates both.

## getLocalUserSession(ByteArrayWrapper ip)

- [ ] The Go `GetLocalUserSessionsByIp` returns `[]*UserSession` (a slice), while Java returns `Queue<UserSession>` from `ipToSessions.get(ip)`. The Java version returns the actual queue reference (which can be mutated), while Go returns a new slice. This is a behavioral difference — the Java version returns a live reference to the queue.

## countLocalOnlineUsers

- [ ] No bugs found. The Go implementation `CountOnlineUsers()` correctly returns the size of the sharded map.

## onSessionEstablished

- [ ] The Go method is completely stubbed out with TODOs. It does not increment the `loggedInUsersCounter` metric (Java: `loggedInUsersCounter.increment()`).
- [ ] The Go method does not push session notifications to clients when `notifyClientsOfSessionInfoAfterConnected` is true. Java calls `userSessionsManager.pushSessionNotification(deviceType, serverId)`.

## addOnSessionClosedListeners

- [ ] The Go method acquires a mutex lock (`s.mu.Lock()`) before appending, while Java simply calls `onSessionClosedListeners.add(onSessionClosed)` without synchronization. This is actually safer in Go but differs from the Java pattern. However, the Java `LinkedList` is not thread-safe either, so the Go version is arguably more correct. Not a bug per se.

## invokeGoOnlineHandlers

- [ ] The Go method is completely stubbed out with a TODO comment. It does not invoke the plugin manager's extension points. Java calls `pluginManager.invokeExtensionPointsSimultaneously(UserOnlineStatusChangeHandler.class, GO_ONLINE_METHOD, handler -> handler.goOnline(userSessionsManager, userSession))`.

## Additional cross-cutting issues

- [ ] The Go `addOnSessionClosedListeners` acquires a write lock on `s.mu` while `notifyOnSessionClosedListeners` in Java iterates without synchronization (using a `LinkedList`). In Go, the `CloseLocalSession` method reads `s.onSessionClosedListeners` under `s.mu.RLock()`, which is correct for concurrent safety. However, the Go code does NOT wrap listener invocations in try-catch like Java does (`try { onSessionClosedListener.accept(session) } catch (Exception e) { LOGGER.error(...) }`). A panicking listener in Go would crash the entire close operation.

# UserService.java
*Checked methods: authenticate(@NotNull Long userId, @Nullable String rawPassword), isActiveAndNotDeleted(@NotNull Long userId)*

Now I have all the information needed. Let me analyze the differences.

## Authenticate

-[x] **Password comparison uses plaintext string equality instead of proper password hashing**: The Java version delegates to `PasswordManager.matchesUserPassword()` which supports BCRYPT, SALTED_SHA256, and NOOP encoding algorithms. The Go version at `user_service.go:34` does a raw string comparison `*passwordHash == rawPassword`, comparing the stored hash directly against the raw password. This means authentication will fail for any user whose password is stored as a BCrypt or Salted-SHA256 hash (the common case), and will only work for NOOP (plaintext) stored passwords.

-[x] **Handling of nil/null password from DB differs**: In Java, `PasswordManager.matchesPassword` returns `true` when both passwords are null, and `false` when the encoded password is non-null but raw password is null. In Go, if `passwordHash` is nil (user not found), it correctly returns `false`, but if the user exists and `FindPassword` returns a nil pointer (empty password field at `user_repository.go:289-291`), the Go code returns `false, nil`. In Java, if the user exists but `user.getPassword()` returns null, then `matchesUserPassword(rawPassword, null)` returns `rawPassword == null` — meaning a user with no stored password and an empty raw password would authenticate successfully in Java but not in Go.

## IsActiveAndNotDeleted

-[x] **Performance: Go fetches full user document instead of using an efficient existence query**: The Java version in `UserRepository.isActiveAndNotDeleted` builds a MongoDB filter with three conditions (`_id`, `isActive=true`, `deletionDate=null`) and uses `mongoClient.exists()` — a highly efficient existence check. The Go version at `user_repository.go:271-277` calls `FindByID` which fetches and decodes the entire user document, then checks `user.IsActive && user.DeletionDate == nil` in Go code. This is functionally equivalent but fetches significantly more data from MongoDB than needed.

-[x] **Missing filter on `isActive` and `deletionDate` at the database query level**: The Java version sends `isActive=true` and `deletionDate=null` as query predicates to MongoDB, so the database only returns a match if all conditions hold. The Go version queries by ID only and checks the fields in application code. While the end result is the same for existing users, the Go approach transfers unnecessary data and shifts filtering responsibility from the database to the application.

# ServiceAddressManager.java
*Checked methods: getWsAddress(), getTcpAddress(), getUdpAddress()*

Now I need to understand the full picture. The Java `getWsAddress()`, `getTcpAddress()`, and `getUdpAddress()` are simple getters, but the core logic is in `updateCustomAddresses()` which populates the fields `wsAddress`, `tcpAddress`, and `udpAddress`. The Go code has none of that logic. Let me verify the Go struct has no fields and no `updateCustomAddresses` equivalent.

## GetWsAddress

-[x] Missing struct fields: The Java class stores `wsAddress` as a `@Nullable String` field updated by `updateCustomAddresses()`. The Go struct `ServiceAddressManager` has **no fields at all**, so it cannot store or return a WebSocket address.
-[x] Stub-only implementation: The method returns a hardcoded empty string `""` instead of returning a stored field value. In Java, `getWsAddress()` returns the `wsAddress` field which is populated with a full `"ws://host:port"` or `"wss://host:port"` URL.
-[x] Missing the `updateCustomAddresses` logic that populates `wsAddress`: The Java code builds the address by calling `queryHost(advertiseStrategy, webSocketProperties.getHost(), advertiseHost)` and formats it with `"ws://"` or `"wss://"` prefix (depending on `adminHttpProperties.getSsl().isEnabled()`), optionally appending `":" + port` based on `attachPortToHost`. None of this logic exists in Go.

## GetTcpAddress

-[x] Missing struct fields: Same as above — no `tcpAddress` field on the Go struct.
-[x] Stub-only implementation: Returns hardcoded `""` instead of a stored address value.
-[x] Missing the `updateCustomAddresses` logic that populates `tcpAddress`: The Java code builds the address by calling `queryHost(advertiseStrategy, tcpProperties.getHost(), advertiseHost)` and formats it as `host + (attachPortToHost ? ":" + port : "")`. None of this logic exists in Go.

## GetUdpAddress

-[x] Missing struct fields: Same as above — no `udpAddress` field on the Go struct.
-[x] Stub-only implementation: Returns hardcoded `""` instead of a stored address value.
-[x] Missing the `updateCustomAddresses` logic that populates `udpAddress`: The Java code builds the address by calling `queryHost(advertiseStrategy, udpProperties.getHost(), advertiseHost)` and formats it as `host + (attachPortToHost ? ":" + port : "")`. None of this logic exists in Go.

# LdapClient.java
*Checked methods: isConnected(), connect(), bind(boolean useFastBind, String dn, String password), search(String baseDn, Scope scope, DerefAliases derefAliases, int sizeLimit, int timeLimit, boolean typeOnly, List<String> attributes, String filter), modify(String dn, List<ModifyOperationChange> changes)*

## isConnected()

- [ ] **Missing method**: The Go `LdapClient` does not have an `IsConnected()` method at all. The Java version checks `connection != null && !connection.isDisposed()`. The Go client has no equivalent — callers cannot query whether the LDAP connection is alive.

## connect()

- [ ] **Missing connection-sharing/caching semantics**: The Java `connect()` uses an atomic `CONNECTION_MONO_UPDATER` to ensure concurrent callers share the same connection attempt (CAS pattern — only one `connect` mono is created, all callers subscribe to it). The Go `Connect()` has no such protection — every call creates a new `ldap.Conn`, overwriting `c.Conn` and leaking the previous connection if any.
- [ ] **No error handling on connection close/replacement**: The Java version disposes the old connection on errors and stores the connection atomically. The Go version unconditionally sets `c.Conn = l` without closing any previously held connection, which leaks the old connection.
- [ ] **Missing LDAP message encoder/decoder pipeline setup**: The Java version explicitly adds `LdapMessageDecoder` and `LdapMessageEncoder` handlers to the pipeline and subscribes to `receiveObject()` with an error handler that disposes the connection. The Go version delegates this entirely to `go-ldap`, which is acceptable as an abstraction choice, but loses the custom error handling that closes the connection on decode errors.

## bind(boolean useFastBind, String dn, String password)

- [ ] **Fast bind control completely ignored**: The Java version passes `REQUEST_CONTROLS_FAST_BIND` (a control with OID `ControlOidConst.FAST_BIND`) when `useFastBind` is `true`. The Go version ignores the `useFastBind` parameter entirely — it performs a standard `SimpleBind` regardless. The comment says "We use Simple Bind regardless" but this changes behavior: fast bind skips the final bind result response for better performance.
- [ ] **Missing result code handling**: The Java version checks `response.isSuccess()`, returns `true` on success, `false` on `INVALID_CREDENTIALS`, and throws `LdapException` for all other result codes with diagnostic message. The Go version simply returns the raw error from `c.Conn.Bind()` — it cannot distinguish between invalid credentials (should return false/nil) and a protocol error (should return an error). The Java version returns `Mono<Boolean>` for exactly this reason.
- [ ] **Wrong return type**: The Java `bind()` returns `Mono<Boolean>` (true = success, false = invalid credentials, error = other failure). The Go version returns `error` only, losing the ability to signal "invalid credentials" as a non-error condition.

## search(String baseDn, Scope scope, DerefAliases derefAliases, int sizeLimit, int timeLimit, boolean typeOnly, List<String> attributes, String filter)

- [ ] **Method completely missing**: The Go `LdapClient` has no `Search` method. The `search` annotation appears on the `ElasticsearchClient.Search()` in `elasticsearch_client.go` (line 46), which is the wrong class — LDAP search and Elasticsearch search are completely different operations. The LDAP client has no search capability at all.

## modify(String dn, List<ModifyOperationChange> changes)

- [ ] **Missing validation for empty changes**: The Java version returns `Mono.empty()` when `changes.isEmpty()`. The Go version accepts a pre-built `*ldap.ModifyRequest` and passes it through without any empty-check.
- [ ] **Missing validation for ADD with empty attribute values**: The Java version iterates over changes and throws `LdapException(INVALID_ATTRIBUTE_SYNTAX, ...)` if any ADD operation has an attribute with no values. The Go version delegates entirely to the `go-ldap` library's `Modify()` with no equivalent validation.
- [ ] **Wrong signature — accepts pre-built request instead of primitive parameters**: The Java version takes `(String dn, List<ModifyOperationChange> changes)` and constructs the `ModifyRequest` internally. The Go version takes a pre-built `*ldap.ModifyRequest`, shifting the construction and validation burden to the caller and bypassing the validation logic that exists in the Java version.

# BerBuffer.java
*Checked methods: skipTag(), skipTagAndLength(), skipTagAndLengthAndValue(), readTag(), peekAndCheckTag(int tag), skipLength(), skipLengthAndValue(), writeLength(int length), readLength(), tryReadLengthIfReadable(), beginSequence(), beginSequence(int tag), endSequence(), writeBoolean(boolean value), writeBoolean(int tag, boolean value), readBoolean(), writeInteger(int value), writeInteger(int tag, int value), readInteger(), readIntWithTag(int tag), writeOctetString(String value), writeOctetString(byte[] value), writeOctetString(int tag, byte[] value), writeOctetString(byte[] value, int start, int length), writeOctetString(int tag, byte[] value, int start, int length), writeOctetString(int tag, String value), writeOctetStrings(List<String> values), readOctetString(), readOctetStringWithTag(int tag), readOctetStringWithLength(int length), writeEnumeration(int value), readEnumeration(), getBytes(), skipBytes(int length), close(), refCnt(), retain(), retain(int increment), touch(), touch(Object hint), release(), release(int decrement), isReadable(int length), isReadable(), isReadableWithEnd(int end), readerIndex()*

The Go `ber_buffer.go` file is entirely stubbed out — every method body is empty, returns zero values, or has no struct fields. This is not a faithful port; it's a skeleton with no implementation at all.

## skipTag()
- [ ] Method body is empty — missing `buffer.skipBytes(1)` equivalent logic

## skipTagAndLength()
- [ ] Method body is empty — missing `buffer.skipBytes(1)` followed by `readLength()` equivalent logic

## skipTagAndLengthAndValue()
- [ ] Method body is empty — missing `buffer.skipBytes(1)`, `readLength()`, and `buffer.skipBytes(length)` equivalent logic

## readTag()
- [ ] Method returns hardcoded `0` instead of reading and returning a byte from the buffer

## peekAndCheckTag(int tag)
- [ ] Method has no return value — Java returns `boolean`. Missing the logic to check `buffer.isReadable() && buffer.getByte(buffer.readerIndex()) == tag`

## skipLength()
- [ ] Method body is empty — missing delegation to `readLength()` equivalent

## skipLengthAndValue()
- [ ] Method body is empty — missing `buffer.skipBytes(readLength())` equivalent logic

## writeLength(int length)
- [ ] Method body is empty — missing all the BER length encoding logic (short form vs long form with 1–4 length bytes)

## readLength()
- [ ] Method returns hardcoded `0` instead of implementing the BER length decoding logic with indefinite-length rejection, bounds checking, and multi-byte length parsing

## tryReadLengthIfReadable()
- [ ] Method returns hardcoded `0` instead of `-1` when not readable — Java returns `-1` for the not-readable case. Missing the full BER length decoding logic with early `-1` return

## beginSequence()
- [ ] Method body is empty — missing delegation to `BeginSequenceWithTag` with the sequence+constructed tag

## beginSequence(int tag)
- [ ] Method body is empty — missing sequence-length writer index tracking, buffer write of tag, and writer index advancement by 3 bytes. The struct also lacks the `sequenceLengthWriterIndexes []int` and `currentSequenceLengthIndex int` fields

## endSequence()
- [ ] Method body is empty — missing the sequence finalization logic: calculating value length, writing 0x82 prefix + 2-byte length, and writer index manipulation

## writeBoolean(boolean value)
- [ ] Method body is empty — missing write of `TAG_BOOLEAN`, length byte `1`, and value byte (`0xFF` or `0`)

## writeBoolean(int tag, boolean value)
- [ ] Method body is empty — missing write of tag byte, length byte `1`, and value byte (`0xFF` or `0`)

## readBoolean()
- [ ] Method returns hardcoded `false` instead of reading tag, validating it, reading length, and returning `buffer.readByte() != 0`

## writeInteger(int value)
- [ ] Method body is empty — missing delegation to `WriteIntegerWithTag` with `TAG_INTEGER`

## writeInteger(int tag, int value)
- [ ] Method body is empty — missing the full signed/unsigned integer BER encoding logic with variable-length byte handling and sign-bit masking

## readInteger()
- [ ] Method returns hardcoded `0` instead of delegating to `ReadIntWithTag` with `TAG_INTEGER`

## readIntWithTag(int tag)
- [ ] Method returns hardcoded `0` — missing tag validation, length reading, and signed integer decoding logic

## writeOctetString(String value)
- [ ] Method body is empty — missing delegation to `WriteOctetStringWithTag` with `TAG_OCTET_STRING`

## writeOctetString(byte[] value)
- [ ] Method body is empty — missing write of `TAG_OCTET_STRING`, length, and byte data

## writeOctetString(int tag, byte[] value)
- [ ] Method body is empty — missing write of tag, length, and byte data

## writeOctetString(byte[] value, int start, int length)
- [ ] Method body is empty — missing write of `TAG_OCTET_STRING`, length, and byte slice data

## writeOctetString(int tag, byte[] value, int start, int length)
- [ ] Method body is empty — missing write of tag, length, and byte slice data

## writeOctetString(int tag, String value)
- [ ] Method body is empty — missing the deferred-length-write logic: advancing writer index, writing UTF-8 bytes, then backfilling the 0x82 + 2-byte length prefix

## writeOctetStrings(List<String> values)
- [ ] Method body is empty — missing iteration over values calling `writeOctetString(TAG_OCTET_STRING, value)`

## readOctetString()
- [ ] Method returns hardcoded `""` instead of delegating to `ReadOctetStringWithTag` with `TAG_OCTET_STRING`

## readOctetStringWithTag(int tag)
- [ ] Method returns hardcoded `""` — missing tag validation, length reading, readability check, and UTF-8 string decoding

## readOctetStringWithLength(int length)
- [ ] Method returns hardcoded `""` — missing length-0 check and UTF-8 byte reading

## writeEnumeration(int value)
- [ ] Method body is empty — missing delegation to `WriteIntegerWithTag` with `TAG_ENUMERATED`

## readEnumeration()
- [ ] Method returns hardcoded `0` instead of delegating to `ReadIntWithTag` with `TAG_ENUMERATED`

## getBytes()
- [ ] Method returns `nil` instead of returning a copy of the buffer's readable bytes

## skipBytes(int length)
- [ ] Method body is empty — missing buffer skip-forward logic

## close()
- [ ] Method is missing entirely — no `Close()` or `close()` method exists in the Go code

## refCnt()
- [ ] Method returns hardcoded `0` instead of delegating to the underlying buffer's reference count

## retain()
- [ ] Method body is empty — missing reference count increment logic

## retain(int increment)
- [ ] Method body is empty — missing reference count increment-by-argument logic

## touch()
- [ ] Method body is empty — missing resource leak detection hint logic

## touch(Object hint)
- [ ] Method body is empty — missing resource leak detection hint logic with hint argument

## release()
- [ ] Method returns hardcoded `false` instead of delegating to the underlying buffer's release

## release(int decrement)
- [ ] Method returns hardcoded `false` instead of delegating to the underlying buffer's release with decrement

## isReadable(int length)
- [ ] Method returns hardcoded `false` instead of checking if the buffer has `length` readable bytes

## isReadable()
- [ ] Method returns hardcoded `false` instead of checking if the buffer has any readable bytes

## isReadableWithEnd(int end)
- [ ] Method returns hardcoded `false` instead of checking `buffer.readerIndex() < end`

## readerIndex()
- [ ] Method returns hardcoded `0` instead of returning the buffer's current reader index

## BerBuffer struct
- [ ] Struct has no fields — missing `buffer` (underlying byte buffer), `sequenceLengthWriterIndexes []int`, and `currentSequenceLengthIndex int` fields required by nearly all methods

# Attribute.java
*Checked methods: isEmpty(), decode(BerBuffer buffer)*

Now I have the full picture. Let me analyze the two specific methods requested:

**1. `Attribute.isEmpty()` (Java) vs Go port**

The Java `Attribute` class has an `isEmpty()` method that checks `values.isEmpty()`. The Go `Attribute` struct has **no fields** (`type`, `values` are missing) and **no `IsEmpty()` method** at all.

**2. `Attribute.decode(BerBuffer buffer)` (Java) vs Go port**

The Java `decode` method is a static factory that:
1. Skips tag and length
2. Reads an octet string as `type`
3. Reads a tag, checks if it's a SET+CONSTRUCTED tag
4. If not, skips length+value and returns `Attribute(type, emptyList)`
5. If yes, reads length, computes end index, loops reading octet strings into a values list
6. Returns `new Attribute(type, values)`

The Go `Decode` method body is **completely empty** — it does nothing.

**3. `isEmpty()` on `UserSessionsManager` (Java via map) vs Go**

The Java code uses `ConcurrentEnumMap.isEmpty()` (which checks `size == 0`). The Go code checks `len(m.Sessions) == 0` with an RLock. This is functionally equivalent — the Go implementation is correct for this method.

---

## Attribute.isEmpty

- [ ] **Missing method entirely**: The Go `Attribute` struct has no `IsEmpty()` method. The Java class has `isEmpty()` that returns `values.isEmpty()`.
- [ ] **Missing struct fields**: The Go `Attribute` struct is empty — it lacks the `Type string` and `Values []string` fields that the Java class has (`private final String type` and `private final List<String> values`). Without these fields, neither `IsEmpty()` nor `Decode()` can be meaningfully implemented.

## Attribute.Decode

- [ ] **Empty method body**: The Go `Decode` method does nothing. The Java `decode` is a static factory method that parses a BER-encoded LDAP attribute from the buffer, including: skipping tag/length, reading the attribute type as an octet string, reading and validating a SET+CONSTRUCTED tag, and looping to read all value octet strings. None of this logic is present.
- [ ] **Wrong method signature**: The Java `decode` is a `static` factory method returning a new `Attribute` instance (`public static Attribute decode(BerBuffer buffer)`). The Go version is a method on `*Attribute` with no return value (`func (a *Attribute) Decode(buffer *asn1.BerBuffer)`). It should likely be a function that returns `*Attribute`.
- [ ] **Missing tag validation logic**: The Java code checks `if (tag != (Asn1IdConst.TAG_SET | Asn1IdConst.FORM_CONSTRUCTED))` and handles the fallback case by calling `buffer.skipLengthAndValue()` and returning an attribute with an empty values list. This entire branch is absent in Go.
- [ ] **Missing values loop**: The Java code loops with `do { values.add(buffer.readOctetString()); } while (buffer.isReadableWithEnd(end))` to collect all attribute values. This is completely missing from the Go implementation.
- [ ] **No return value / field population**: Even if restructured as a method on `*Attribute`, the Go version never assigns to `a.Type` or `a.Values` (which don't exist on the struct), so any decoded data is silently discarded.

# LdapMessage.java
*Checked methods: estimateSize(), writeTo(BerBuffer buffer)*

## estimateSize()
- [ ] **Matching return value but structurally incomplete**: Both Java and Go return `0`, so the logic is technically identical. However, this is correct — no bug here.

## WriteTo()
- [ ] **Method body is completely empty**: The Java `writeTo` contains the full LDAP message serialization logic: `buffer.beginSequence()`, `writeInteger(messageId)`, calling `protocolOperation.writeTo(buffer)`, iterating and writing controls, and calling `buffer.endSequence()`. The Go `WriteTo` method body is entirely empty — it does nothing.
- [ ] **Missing struct fields**: The `LdapMessage` struct has no fields. The Java class has `messageId` (int), `protocolOperation` (generic T extending ProtocolOperation), and `controls` (List<Control>). None of these are present in the Go struct.
- [ ] **Missing outer sequence write**: Java calls `buffer.beginSequence()` at the start and `buffer.endSequence()` at the end. Go has neither.
- [ ] **Missing messageId write**: Java calls `buffer.writeInteger(messageId)`. Go omits this entirely.
- [ ] **Missing protocolOperation.writeTo() delegation**: Java delegates to `protocolOperation.writeTo(buffer)`. Go does not call any operation's WriteTo.
- [ ] **Missing controls serialization**: Java iterates over controls, writes each control's OID as an octet string, writes criticality as a boolean if true, all wrapped in a `beginSequence(LdapTagConst.CONTROLS)` / `endSequence()` block. Go has none of this logic.

# LdapResult.java
*Checked methods: isSuccess()*

## IsSuccess

- [ ] **Missing fields**: The `LdapResult` struct has no fields at all. The Java version has `resultCode`, `matchedDn`, `diagnosticMessage`, and `referrals`. The `IsSuccess()` method needs `resultCode` to function.
- [ ] **Hardcoded return value**: `IsSuccess()` always returns `false`. The Java version returns `resultCode == ResultCodeConst.SUCCESS`. Since there is no `resultCode` field, the comparison logic is completely missing.
- [ ] **Missing `ResultCodeConst` reference**: The Go implementation does not reference or define the `ResultCodeConst.SUCCESS` constant needed for the success comparison.

# Control.java
*Checked methods: decode(BerBuffer buffer)*

## decode

- [ ] The Go `Control.Decode` method is completely empty (no body), while the Java `decode` is a static method with full parsing logic: it peeks for the `CONTROLS` tag, reads tag/length, iterates over control elements reading OID, criticality, and optional value, and returns a list of `Control` objects.
- [ ] The Go `Control` struct has no fields (`oid`, `criticality`) — the Java class has `oid` (String) and `criticality` (boolean) fields.
- [ ] The Java `decode` is a static factory method returning `List<Control>`, but the Go version is an instance method `func (c *Control) Decode(buffer *asn1.BerBuffer)` with no return value — the method signature does not match the Java design.
- [ ] All core logic is missing: reading the outer tag and length, looping over control entries, reading OID via `readOctetString()`, conditionally reading `criticality` via `readBoolean()`, skipping optional octet string values, and collecting `Control` instances into a slice.
- [ ] The guard clause `if (!buffer.isReadable() || !buffer.peekAndCheckTag(LdapTagConst.CONTROLS))` returning an empty list is missing.

# BindRequest.java
*Checked methods: estimateSize(), writeTo(BerBuffer buffer)*

## EstimateSize
- [ ] Returns a hardcoded `0` instead of `dn.length() + password.length() + 16`. The Go `BindRequest` struct has no `dn` or `password` fields, so the calculation cannot be performed at all.

## WriteTo
- [ ] Method body is completely empty. The Java version writes: `buffer.beginSequence(LdapTagConst.BIND_REQUEST).writeInteger(LDAP_VERSION3).writeOctetString(dn).writeOctetString(Asn1IdConst.TAG_CLASS_CONTEXT, password).endSequence()`. None of this logic is present.
- [ ] The Go `BindRequest` struct lacks the `dn` and `password` fields that are required by both `estimateSize()` and `writeTo()`.

# BindResponse.java
*Checked methods: decode(BerBuffer buffer)*

Now I have a complete picture of both the Java and Go implementations. Let me compare them.

## BindResponse.Decode

- [ ] **Missing struct fields**: The Go `BindResponse` struct at line 63-64 is empty (`struct{}`), but the Java `BindResponse` extends `LdapResult` which holds `resultCode`, `matchedDn`, `diagnosticMessage`, and `referrals` fields. The Go struct lacks all these fields entirely.

- [ ] **Method body is empty / all core logic missing**: The Go `Decode` method at line 67-68 is a no-op (`func (r *BindResponse) Decode(buffer *asn1.BerBuffer) {}`). The Java `decode` method performs: (1) `buffer.skipTagAndLength()`, (2) calls `LdapResult.decodeResult(buffer)` which reads `resultCode` (enumeration), `matchedDn` (octet string), `diagnosticMessage` (octet string), and optionally parses referrals with tag checking, length validation, and loop-based referral reading, then (3) returns a new `BindResponse` with the decoded fields. None of this logic exists in the Go version.

- [ ] **Does not return the decoded result**: The Java method returns a new `BindResponse` instance with the decoded data. The Go method has no return value and does not populate the receiver struct with any decoded data.

# ModifyRequest.java
*Checked methods: estimateSize(), writeTo(BerBuffer buffer)*

## EstimateSize

- [ ] Go method returns hardcoded `0` instead of `dn.length() + changes.size() * 32`. Missing `dn` and `changes` fields on the struct, and missing the size estimation logic.

## WriteTo

- [ ] Go method is completely empty (no-op). Missing all logic: the outer `beginSequence` with `MODIFY_REQUEST` tag, writing the DN as octet string, iterating over changes, writing enumeration for change type, writing attribute type and values with proper sequence nesting (5 total `beginSequence`/`endSequence` pairs), and all corresponding `endSequence` calls.
- [ ] Missing `dn` field on the `ModifyRequest` struct.
- [ ] Missing `changes` field (list of `ModifyOperationChange`) on the `ModifyRequest` struct.
- [ ] Missing `ModifyOperationChange` type/struct with `getType()` and `getAttribute()` accessors.
- [ ] Missing `ModifyOperationType` enum with `getValue()`.
- [ ] Missing `Attribute` struct fields (`Type` string and `Values` slice) — currently an empty struct.

# ModifyResponse.java
*Checked methods: decode(BerBuffer buffer)*

## ModifyResponse.Decode

- [ ] **Entire decode logic is missing**: The Go method body is empty (`func (r *ModifyResponse) Decode(buffer *asn1.BerBuffer) {}`). The Java version calls `buffer.skipTagAndLength()`, then `LdapResult.decodeResult(buffer)`, and returns a new `ModifyResponse` populated with the result's fields. None of this logic is present in Go.
- [ ] **ModifyResponse struct has no fields**: The Java `ModifyResponse` extends `LdapResult`, which holds `resultCode`, `matchedDn`, `diagnosticMessage`, and `referrals`. The Go struct `ModifyResponse` is completely empty, so even if decode logic were added, there would be nowhere to store the decoded values.
- [ ] **LdapResult struct has no fields and decodeResult is not implemented**: The Go `LdapResult` struct is empty and has no `decodeResult` equivalent. The Java `decode` depends on `LdapResult.decodeResult(buffer)` to parse the LDAP result fields, which is not ported at all.
- [ ] **Decode does not return a value**: The Java `decode` method returns a new `ModifyResponse` instance. The Go method has a `void`-like signature with no return value, making it impossible to use the decoded result.

# Filter.java
*Checked methods: write(BerBuffer buffer, String filter)*

## Filter.Write

- [ ] The Go method `Write` has an entirely empty body (`func (f *Filter) Write(buffer *asn1.BerBuffer, filter string) {}`). The Java version contains the complete LDAP filter parsing and BER encoding logic including: converting the filter string to bytes, calling `writeFilter`, which in turn handles parenthesis balancing, escape sequences, filter types (AND/OR/NOT/equality/substring/greater/less/approximate/extensible-match/present), and delegates to `writeFilterSet`, `writeFilterInSet`, `writeSubstringFilter`, `writeExtensibleMatchFilter`, and `unescapeFilterValue`. None of this logic is ported.

# SearchRequest.java
*Checked methods: estimateSize(), writeTo(BerBuffer buffer)*

## EstimateSize

- [x] Returns `0` instead of `128`. The Java version returns `128` as a fixed estimate.

## WriteTo

- [x] Method body is completely empty (no-op). The Java version writes a full LDAP search request: begins a sequence with `SEARCH_REQUEST` tag, writes `baseDn` as octet string, `scope` as enumeration, `derefAliases` as enumeration, `sizeLimit` as integer, `timeLimit` as integer, `typesOnly` as boolean, writes the filter via `Filter.write()`, begins an inner sequence for attributes writing them as octet strings, and ends both sequences. None of this logic is present.
- [x] `SearchRequest` struct has no fields at all. The Java version has `baseDn`, `scope`, `derefAliases`, `sizeLimit`, `timeLimit`, `typesOnly`, `attributes`, and `filter` fields — all unported.

# SearchResult.java
*Checked methods: decode(BerBuffer buffer), isComplete()*

## Decode

- [x] **Method body is entirely empty** — The Go `Decode` method has an empty body `{}` and discards the `buffer` parameter, whereas the Java version contains the full LDAP search result decoding logic: reading a tag, skipping length, switching on `SEARCH_RESULT_ENTRY` vs `SEARCH_RESULT_DONE`, parsing attributes, object names, controls, and constructing a new `SearchResult` with the parsed data.
- [x] **Missing return value / result propagation** — The Java `decode` returns a `SearchResult` (either a partial result with entries or a completed result with `LdapResult` data). The Go version returns nothing (`void` equivalent) and does not produce or return any decoded result.
- [x] **Missing `SEARCH_RESULT_ENTRY` branch** — The Java code parses an entry's object name, its attributes (looping with `isReadableWithEnd`), decodes controls, creates a `SearchResultEntry`, wraps it in a list, and returns a `SearchResult` with an overridden `isComplete()` returning `false`. None of this logic exists in Go.
- [x] **Missing `SEARCH_RESULT_DONE` branch** — The Java code calls `LdapResult.decodeResult(buffer)`, constructs a new `SearchResult` with the decoded result code, matched DN, diagnostic message, referrals, and the previously accumulated entries. None of this logic exists in Go.
- [x] **Missing error handling for unexpected tags** — The Java code throws an `LdapException` with `ResultCode.PROTOCOL_ERROR` for unexpected tags and for receiving `SEARCH_RESULT_DONE` when entries haven't been accumulated yet. The Go version has no error handling at all.
- [x] **Missing fields on `SearchResult` struct** — The Go struct has no fields at all. The Java class has fields inherited from `LdapResult` (`resultCode`, `matchedDn`, `diagnosticMessage`, `referrals`) and its own `entries` field. None of these are declared in Go.
- [x] **Missing `SearchResultEntry` type** — The Java code references `SearchResultEntry` (with `objectName`, `attributes`, `controls` fields). This type does not exist in the Go file.

## IsComplete

- [x] **Always returns `true`, but Java behavior is context-dependent** — In Java, the base `SearchResult` class does not override `isComplete()`, so it inherits the default (which returns `true` from `LdapResult`). However, the `decode` method returns an **anonymous subclass** when parsing a `SEARCH_RESULT_ENTRY` where `isComplete()` is overridden to return `false`. The Go version hardcodes `true` with no mechanism to represent an incomplete/partial search result. This means a caller cannot distinguish between an intermediate entry result and a final done result.

# ApiLoggingContext.java
*Checked methods: shouldLogHeartbeatRequest(), shouldLogRequest(TurmsRequest.KindCase requestType), shouldLogNotification(TurmsRequest.KindCase requestType)*

## ShouldLogHeartbeatRequest

- [x] **Stub returns `true` unconditionally; Java uses probabilistic sampling.** The Java version calls `shouldLog(heartbeatSampleRate)` which applies a sample-rate check: returns `false` if rate ≤ 0, returns `true` if rate ≥ 1.0, otherwise returns `ThreadLocalRandom.current().nextFloat() < sampleRate`. The Go stub always returns `true`, ignoring the heartbeat sample rate entirely.
- [x] **Missing `heartbeatSampleRate` field.** The Go struct has no fields at all. The Java class stores `heartbeatSampleRate` from configuration and uses it in `shouldLogHeartbeatRequest()`.

## ShouldLogRequest

- [x] **Stub returns `true` unconditionally; Java uses a filtered map lookup with sampling.** The Java version looks up `requestType` in `typeToSupportedLoggingRequestProperties`. If the type is not in the map (excluded or never included), it returns `false`. If found, it applies the per-request sample rate. The Go stub always returns `true`.
- [x] **Missing `typeToSupportedLoggingRequestProperties` field.** The Go struct has no fields. The Java class builds this map from included/excluded categories and request types during construction.
- [x] **Missing `getSupportedLoggingRequestProperties` construction logic.** The Java constructor builds the map via a 4-step process: (1) add included categories, (2) add included individual requests, (3) remove excluded categories, (4) remove excluded individual request types. None of this logic exists in Go.

## ShouldLogNotification

- [x] **Stub returns `true` unconditionally; Java uses a separate filtered map lookup with sampling.** The Java version looks up `requestType` in `typeToSupportedLoggingNotificationsProperties` (a separate map from notification-specific config). If not found, returns `false`; if found, applies the per-notification sample rate. The Go stub always returns `true`.
- [x] **Missing `typeToSupportedLoggingNotificationsProperties` field.** The Go struct has no fields. The Java class maintains a separate map for notification logging, built from notification-specific included/excluded categories and types.
- [x] **Missing the entire `shouldLog(requestType, propertiesMap)` and `shouldLog(sampleRate)` base logic.** The underlying probabilistic sampling logic (sample rate ≤ 0 → false, ≥ 1.0 → true, else random float comparison) is not implemented anywhere in the Go code.

# ClientApiLogging.java
*Checked methods: log(@Nullable Integer sessionId, @Nullable Long userId, @Nullable DeviceType deviceType, @Nullable Integer version, String ip, long requestId, TurmsRequest.KindCase requestType, int requestSize, long requestTime, TurmsNotification response, long processingTime), log(@Nullable Integer sessionId, @Nullable Long userId, @Nullable DeviceType deviceType, @Nullable Integer version, String ip, long requestId, TurmsRequest.KindCase requestType, int requestSize, long requestTime, int responseCode, long processingTime), log(@Nullable Integer sessionId, @Nullable Long userId, @Nullable DeviceType deviceType, @Nullable Integer version, String ip, long requestId, String requestType, int requestSize, long requestTime, int responseCode, @Nullable String responseDataType, int responseSize, long processingTime)*

## Log (TurmsNotification response variant)

- [x] **Stub implementation with no real logic**: The Go code is a single stub function `Log(request interface{}, serviceRequest interface{}, ...)` that only calls `log.Printf("mock client api log: %v", request)`. It does not implement any of the actual logging logic from the Java `log(..., TurmsNotification response, ...)` method, which joins 14 fields (sessionId, userId, deviceType, version, ip, requestId, requestType, requestSize, requestTime, response.getCode(), response data kind case, response serialized size, processingTime) with delimiters and logs via `CLIENT_API_LOGGER.info()`.

## Log (int responseCode variant)

- [x] **Stub implementation with no real logic**: The Go code does not implement the Java `log(..., int responseCode, long processingTime)` method, which joins 14 fields including `null` for response data type, `'0'` for response serialized size, and the numeric `responseCode` and `processingTime`.

## Log (String requestType variant)

- [x] **Stub implementation with no real logic**: The Go code does not implement the Java `log(..., String requestType, ..., String responseDataType, int responseSize, ...)` method, which joins 14 fields including the string `requestType`, numeric `responseCode`, `responseDataType`, `responseSize`, and `processingTime`.

- [x] **No separate method signatures**: All three Java overloads are collapsed into a single generic `Log(... interface{})` stub with no field extraction, no delimiter-based formatting, and no structured field output — entirely different behavior from the Java versions which produce structured, delimited log lines with specific field ordering.

# NotificationLoggingManager.java
*Checked methods: log(SimpleTurmsNotification notification, int notificationBytes, int recipientCount, int onlineRecipientCount)*

## Log

- [x] The Go code is a stub implementation with no actual logic. It receives generic `interface{}` parameters and simply calls `log.Printf`, completely discarding the structured logging format from Java. None of the Java method's core logic is implemented:
  - Missing: extraction of `closeStatus` from the notification
  - Missing: formatting of `requesterId` from the notification
  - Missing: formatting of `recipientCount` and `onlineRecipientCount`
  - Missing: formatting of `notificationBytes`
  - Missing: extraction and logging of `relayedRequestType` name
  - Missing: field-delimited output format (joined by `LOG_FIELD_DELIMITER`)
  - Missing: writing to `NOTIFICATION_LOGGER` (uses Go's default `log.Printf` instead)
- [x] The method signature doesn't match the Java `log(SimpleTurmsNotification notification, int notificationBytes, int recipientCount, int onlineRecipientCount)` — instead of typed parameters for the notification fields, it uses generic `interface{}` parameters, making it impossible to extract the specific fields (closeStatus, requesterId, relayedRequestType) that the Java code logs.

# SimpleTurmsNotification.java
*Checked methods: SimpleTurmsNotification(long requesterId, Integer closeStatus, TurmsRequest.KindCase relayedRequestType)*

The Java `SimpleTurmsNotification` is a simple record (data class) with three fields. The Go struct and constructor map to it correctly. The core parsing logic lives in `TurmsNotificationParser.parseSimpleNotification`, which in Go is a stub returning `nil, nil`.

The question asks me to focus specifically on the `SimpleTurmsNotification` constructor (the record with parameters `long requesterId, Integer closeStatus, TurmsRequest.KindCase relayedRequestType`).

## SimpleTurmsNotification Constructor (NewSimpleTurmsNotification)

- [x] **`RelayedRequestType` field typed as `any` instead of a proper enum type**: In Java, `relayedRequestType` is `TurmsRequest.KindCase`, a protobuf enum type. In Go, it's declared as `any` (empty interface), losing type safety and making it impossible to distinguish between valid `KindCase` values and arbitrary types. The Go code should use a specific enum type (e.g., a generated protobuf enum or a dedicated type representing `TurmsRequest.KindCase`).

- [x] **`CloseStatus` field typed as `*int32` instead of `*int`**: In Java, `closeStatus` is `Integer` (a boxed 32-bit int). Go's `*int32` is technically correct for size, but the Java tests assign an `int` (e.g., `closeStatus = 1` in tests, and `readInt32()` in the parser). The `*int32` choice is acceptable for bit-width matching, so this is borderline. However, since the protobuf definition uses `int32` and Java uses `Integer`, `*int32` is reasonable — not a real bug.

- [x] **No validation in the Go constructor**: The Java `TurmsNotificationParser.parseSimpleNotification` enforces that `requesterId` must not be `UNSET_VALUE` (i.e., must have been set), `type` must not be `null` or `KIND_NOT_SET`, and duplicate fields throw errors. The Go `NewSimpleTurmsNotification` performs no validation at all. While the constructor itself in Java also performs no validation (it's a simple record), the parsing function that *creates* instances does. The Go `ParseSimpleNotification` is a stub (`return nil, nil`), which means **all core parsing logic is missing** — the tag-based field extraction, duplicate detection, the `requesterId`/`type` presence validation, and the relayed request sub-message parsing.

Since the question specifically scopes to the `SimpleTurmsNotification` constructor itself (the record), not the parser:

## SimpleTurmsNotification

- [x] **`RelayedRequestType` uses `any` instead of a typed enum equivalent to `TurmsRequest.KindCase`**: The Java record field is `TurmsRequest.KindCase`, a strongly-typed enum. The Go struct uses `any`, which provides no type safety and no way to represent the `KIND_NOT_SET` sentinel value that the Java code checks against. This means downstream consumers cannot reliably switch on the request type or validate it, and any value can be passed without compilation errors.

## ParseSimpleNotification (TurmsNotificationParser)

- [x] **Entire method body is a stub (`return nil, nil`)**: The Java implementation contains ~50 lines of protobuf wire-format parsing logic including: reading tags in a loop (`REQUESTER_ID_TAG=80`, `CLOSE_STATUS_TAG=88`, `RELAYED_REQUEST_TAG=98`), extracting `requesterId` via `readInt64()`, extracting `closeStatus` via `readInt32()`, parsing the relayed request sub-message to determine the `KindCase`, duplicate field detection with error throwing, `requesterId` presence validation, `type` null/KIND_NOT_SET validation, and proper error wrapping for `IOException`. The Go stub returns `nil, nil` unconditionally, meaning every call silently produces no notification and no error.

# SimpleTurmsRequest.java
*Checked methods: SimpleTurmsRequest(long requestId, TurmsRequest.KindCase type, CreateSessionRequest createSessionRequest), toString()*

## NewSimpleTurmsRequest (constructor)

- [x] **Missing `Type` field type safety**: Java uses `TurmsRequest.KindCase` (a typed enum), while Go uses `any` for the `Type` field. This is a type safety regression — any arbitrary value can be passed as `reqType`, losing the compile-time guarantee that only valid `KindCase` values are accepted.

## ToString

- [x] **Entire method body is a stub returning an empty string**: The Java `toString()` formats the record as `"SimpleTurmsRequest[requestId=<id>, type=<type>, createSessionRequest=<formatted>]"`, including all three fields with `ProtoFormatter.toLogString` for the `createSessionRequest`. The Go `ToString()` returns `""` (empty string), providing zero diagnostic value.

## SetConnection (related field mapping)

- [x] **`ip` parameter is ignored**: The Java `setConnection(NetConnection connection, ByteArrayWrapper ip)` stores the IP address, but the Go `SetConnection(connection Connection, ip string)` accepts `ip` as a parameter yet never assigns it to `s.IP`. The IP field on `UserSession` is never populated through this method, meaning session IP tracking is broken.

# TurmsNotificationParser.java
*Checked methods: parseSimpleNotification(CodedInputStream turmsRequestInputStream)*

## ParseSimpleNotification

- [x] **Entire method is a stub** — The Go implementation of `ParseSimpleNotification` returns `nil, nil` with a `// Stub implementation` comment. None of the Java logic has been ported. The following core logic is completely missing:
  - Parsing the protobuf wire format from the input byte slice (reading tags in a loop)
  - Extracting `requesterId` from tag `80` (field 10, wire type 0 / int64)
  - Extracting `closeStatus` from tag `88` (field 11, wire type 0 / int32)
  - Extracting the relayed request type from tag `98` (field 12, wire type 2 / length-delimited), which involves reading the length varint, then reading the inner tag to get the `KindCase` field number
  - Duplicate field detection (throwing/returning error if `requesterId` or `closeStatus` appear more than once)
  - Validation that `requesterId` is present (error if missing)
  - Validation that the request `type` is not null/nil and not `KIND_NOT_SET` (error if unknown or unset)
  - Constructing and returning a `SimpleTurmsNotification` with the parsed values
- [x] **Return type mismatch** — The Java method returns a `SimpleTurmsNotification` with `(requesterId, closeStatus, type)` where `type` is a `TurmsRequest.KindCase` enum. The Go struct stores `RelayedRequestType` as `any` rather than a typed enum equivalent, losing type safety and making it impossible to compare against a `KIND_NOT_SET` sentinel value as the Java code does.

# TurmsRequestParser.java
*Checked methods: parseSimpleRequest(CodedInputStream turmsRequestInputStream)*

## ParseSimpleRequest

- [x] **Entire method is a stub**: The Go `ParseSimpleRequest` method returns `nil, nil` immediately without any logic. The Java version implements a full protobuf wire-format parser that: reads tags in a loop, extracts `requestId` (field 1, varint), detects duplicate request IDs, identifies the request type from the field number of the first non-requestId tag, validates that a request ID was found, validates that the request type is known, and conditionally parses a `CreateSessionRequest` sub-message. None of this logic exists in the Go code.

- [x] **Missing duplicate request ID detection**: The Java code throws `ILLEGAL_ARGUMENT` with "Duplicate request ID" if the `requestId` field appears more than once. The Go stub has no such check.

- [x] **Missing request ID validation**: The Java code validates that a request ID was actually present (not left at `UNDEFINED_REQUEST_ID` / `Long.MIN_VALUE`) and throws "No request ID" if missing. The Go stub has no such validation.

- [x] **Missing request type validation**: The Java code checks that `type` is not null and not `KIND_NOT_SET`, throwing "Unknown request type" if so. The Go stub has no such validation.

- [x] **Missing CreateSessionRequest parsing**: When the request type is `CREATE_SESSION_REQUEST`, the Java code reads and parses the full `CreateSessionRequest` sub-message from the stream. The Go stub does not parse any sub-message.

- [x] **Missing error wrapping for IO/parse failures**: The Java code catches `IOException` and wraps it in a `ResponseException` with `ILLEGAL_ARGUMENT` status. The Go stub has no error handling for malformed input.

- [x] **Method signature uses `[]byte` instead of a stream**: The Java version takes `CodedInputStream` (a streaming protobuf decoder), while the Go version takes a raw `[]byte`. While this is a reasonable adaptation, the implementation still needs to actually parse the protobuf wire format from those bytes — which it does not do at all.

# MongoConfig.java
*Checked methods: adminMongoClient(TurmsPropertiesManager propertiesManager), userMongoClient(TurmsPropertiesManager propertiesManager), mongoDataGenerator()*

Now I have all the context needed to compare the three methods. Here is the analysis:

## adminMongoClient

- [x] **Method body is entirely unimplemented** — The Go method `AdminMongoClient()` simply returns `nil` with no logic. The Java version: (1) retrieves `MongoProperties` from `propertiesManager.getLocalProperties().getGateway().getMongo().getAdmin()`, (2) calls `getMongoClient(properties, "admin", Collections.emptySet())` to create a `TurmsMongoClient`, and (3) calls `mongoClient.registerEntitiesByClasses(Admin.class, AdminRole.class)` to register entity classes. None of this logic is present in Go.

## userMongoClient

- [x] **Method body is entirely unimplemented** — The Go method `UserMongoClient()` simply returns `nil` with no logic. The Java version: (1) retrieves `IdentityAccessManagementProperties` from the local properties, (2) checks if authentication is enabled AND type is `PASSWORD` — returns `null` if not, (3) retrieves `MongoProperties` from `getGateway().getMongo().getUser()`, (4) calls `getMongoClient(properties, "user", Set.of(ClusterType.SHARDED, ClusterType.LOAD_BALANCED))` with a non-empty cluster type set, and (5) calls `mongoClient.registerEntitiesByClasses(User.class)`. None of this logic is present in Go.
- [x] **Missing conditional null-return logic** — The Java version conditionally returns `null` when `IdentityAccessManagement` is disabled or not of type `PASSWORD`. The Go stub always returns `nil` with no conditional evaluation.

## mongoDataGenerator

- [x] **Missing the `IMongoCollectionInitializer` interface/contract** — The Java version returns an anonymous implementation of `IMongoCollectionInitializer` registered as a Spring bean with name `IMongoCollectionInitializer.BEAN_NAME`. This serves as a placeholder initializer. The Go version returns a bare `*MongoDataGenerator` struct with no interface implementation and no indication of the collection initializer contract.

# TurmsServiceApplication.java
*Checked methods: main(String[] args)*

## main(String[] args)

- [ ] **Missing bootstrap/application initialization**: The Java `main` calls `bootstrap(TurmsServiceApplication.class, args)` which triggers the full Spring Boot application context initialization (dependency injection, component scanning, configuration loading, cluster node startup). The Go version only prints a log message and has no initialization logic.
- [ ] **Missing node type declaration**: The Java class uses `@Application(nodeType = NodeType.SERVICE)` to declare itself as a SERVICE node type. The Go version does not declare any node type.
- [ ] **Missing component scanning**: The Java class uses `@SpringBootApplication(scanBasePackages = {PackageConst.SERVICE, PackageConst.SERVER_COMMON})` to scan for service and common server components. The Go version has no equivalent component/service registration.
- [ ] **Missing class reference passing**: The Java `main` passes `TurmsServiceApplication.class` as a configuration source to the bootstrap. The Go version has no equivalent application configuration passing mechanism.
- [ ] **Command-line arguments not used**: The Java `main` passes `args` to `bootstrap()` for configuration overrides and Spring profile activation. The Go `main()` does not accept or process any command-line arguments.

# ServiceRequestDispatcher.java
*Checked methods: dispatch(TracingContext context, ServiceRequest serviceRequest)*

Now I have a clear picture. The Java `ServiceRequestDispatcher.dispatch` is a method in `turms-service` that handles client requests forwarded from `turms-gateway` via RPC. The Go `Router.Dispatch` is an RPC frame router that simply looks up a handler by codec ID and calls it. These are fundamentally different in scope.

The Java `dispatch` method does all of the following:
1. Pending request counting (increment/decrement with shutdown hook)
2. Buffer touch tracking
3. Buffer release in finally block
4. Dispatching to `dispatch0` which:
   - Validates userId and deviceType
   - Checks service availability
   - Decodes protobuf request (with blocklist on failure)
   - Transforms request via plugin extensions
   - Looks up handler by request type
   - Invokes handler with plugin extension points
   - Metrics recording
   - Notifies related users of actions
   - Error handling with logging
   - API logging

The Go `Router.Dispatch` does none of this — it's just a generic handler lookup by codec ID.

## Router.Dispatch

- [ ] **Missing: All core dispatch logic from Java's `dispatch` method.** The Java `dispatch(TracingContext, ServiceRequest)` is a complex method that: (1) tracks pending request count with shutdown coordination, (2) touches/releases the request buffer for resource tracking, (3) delegates to `dispatch0` which validates userId/deviceType, checks service availability, decodes protobuf, handles blocklist on corrupted requests, transforms requests via plugin extensions, routes to a handler by request type, records metrics, notifies related users, handles errors, and performs API logging. The Go `Dispatch` method only performs a simple handler lookup by codec ID and invokes it — entirely missing the request dispatch semantics of the Java version.
- [ ] **Missing: Pending request counting.** Java tracks `pendingRequestCount` with `AtomicInteger`, incrementing on entry and decrementing via `doFinally`/`onPendingRequestHandled()`, with a shutdown hook that waits for all pending requests. Go has no equivalent.
- [ ] **Missing: Service availability check.** Java validates `serverStatusManager.getServiceAvailability()` before processing. Go does not check service availability.
- [ ] **Missing: User ID and device type validation.** Java returns `SERVER_INTERNAL_ERROR` if `userId` or `deviceType` is null. Go does not validate these.
- [ ] **Missing: Protobuf decoding with blocklist on failure.** Java decodes the request buffer into `TurmsRequest`, and on `IOException` calls `blocklistService.tryBlockIpForCorruptedRequest` and `tryBlockUserIdForCorruptedRequest`. Go does not decode or handle blocklisting.
- [ ] **Missing: Plugin extension point integration.** Java checks for running `ClientRequestTransformer` extensions and invokes them to transform the request, then invokes `ClientRequestHandler` extensions. Go has no plugin system integration.
- [ ] **Missing: Request type-based handler routing.** Java looks up handlers via `requestTypeToHandler.get(requestType)` using `TurmsRequest.KindCase`. Go routes by numeric `CodecID` instead, which is a different dispatch mechanism entirely.
- [ ] **Missing: Notification of related users.** Java calls `notifyRelatedUsersOfAction` to forward notifications to recipients and other sessions of the requester, with plugin extension points for before/after notify. Go has no notification logic.
- [ ] **Missing: Error handling with status code mapping.** Java maps all `Throwable` to `ServiceResponse` via `ThrowableInfo.get(t)` with server error logging. Go simply returns the error from the handler.
- [ ] **Missing: API request logging.** Java logs requests via `ClientApiLogging.log(...)` when the response is a server error or when `apiLoggingContext.shouldLogRequest(requestType)` is true. Go has no request logging.
- [ ] **Missing: Metrics recording.** Java uses `.name(TURMS_CLIENT_REQUEST).tag(TURMS_CLIENT_REQUEST_TAG_TYPE, requestType.name()).metrics()` for request metrics. Go has no metrics.
- [ ] **Different: The Go method signature accepts `*codec.RpcFrame` instead of `ServiceRequest`.** The Java method operates on a `ServiceRequest` containing userId, deviceType, IP, request buffer, etc. The Go method operates on a generic RPC frame with just a codec ID and payload bytes, making it impossible to implement the same logic without the caller providing the equivalent context.

# ClientRequest.java
*Checked methods: toString(), turmsRequest(), userId(), deviceType(), clientIp(), requestId(), equals(Object obj), hashCode()*

## toString()

- [ ] **Missing `String()` method on `ClientRequest`**: The Java `toString()` returns a formatted string with userId, deviceType, clientIp (as array), requestId, and turmsRequest (via ProtoFormatter). The Go `ClientRequest` has no `String()` / `ToString()` method at all.

## turmsRequest()

- [ ] **Missing lazy-build logic**: The Java version lazily builds `turmsRequest` from `turmsRequestBuilder` if `turmsRequest` is null. The Go version has no `turmsRequestBuilder` field and no lazy-build logic — it simply returns the stored `turmsRequest` pointer, which means if `turmsRequest` is nil it will always return nil rather than building from a builder.

## clientIp()

- [ ] **Type mismatch — `*string` vs `[]byte` / `net.IP`**: In Java, `clientIp` is `byte[]`. In Go, `ClientRequest` stores it as `*string`. This is a semantic type change — the Java version stores raw bytes (an IP address in byte form), while the Go version stores a string pointer, losing the raw byte representation semantics.

## equals(Object obj)

- [ ] **Stub implementation — always returns `false`**: The Go `Equals()` method is a stub that unconditionally returns `false`. The Java version performs a proper equality check comparing userId, deviceType, clientIp (via `Arrays.equals`), requestId, turmsRequestBuilder, and turmsRequest.

## hashCode()

- [ ] **Stub implementation — always returns `0`**: The Go `HashCode()` method is a stub that unconditionally returns `0`. The Java version computes a proper hash using `Objects.hash()` for most fields combined with `Arrays.hashCode(clientIp)`, following the standard `31 * hash + arrayHash` pattern.

# RequestHandlerResult.java
*Checked methods: RequestHandlerResult(ResponseStatusCode code, @Nullable String reason, @Nullable TurmsNotification.Data response, List<Notification> notifications), toString(), of(@NotNull ResponseStatusCode code), of(@NotNull ResponseStatusCode code, @Nullable String reason), of(@NotNull TurmsNotification.Data response), of(boolean forwardNotificationsToRequesterOtherOnlineSessions, @NotNull TurmsRequest notification), of(boolean forwardNotificationsToRequesterOtherOnlineSessions, @NotNull Long recipientId, @NotNull TurmsRequest notification), of(@NotNull Long recipientId, @NotNull TurmsRequest notification), of(@NotEmpty Set<Long> recipientIds, @NotNull TurmsRequest dataForRecipient), of(boolean forwardNotificationToRequesterOtherOnlineSessions, @NotEmpty Set<Long> recipientIds, @NotNull TurmsRequest notification), of(TurmsNotification.Data response, @NotEmpty Set<Long> recipientIds, @NotNull TurmsRequest notificationForRecipients, @NotNull TurmsRequest notificationForRequesterOtherOnlineSessions), of(@NotEmpty Set<Long> recipientIds, @NotNull TurmsRequest notificationForRecipients, @NotNull TurmsRequest notificationForRequesterOtherOnlineSessions), of(TurmsNotification.Data response, boolean forwardNotificationToRequesterOtherOnlineSessions, @NotEmpty Set<Long> recipientIds, @NotNull TurmsRequest notification), of(@NotNull ResponseStatusCode code, @NotNull Long recipientId, @NotNull TurmsRequest notification), of(@NotNull ResponseStatusCode code, boolean forwardNotificationToRequesterOtherOnlineSessions, @NotNull Long recipientId, @NotNull TurmsRequest notification), of(@NotNull List<Notification> notifications), of(@NotNull Notification notification), ofDataLong(@NotNull Long value), ofDataLong(@NotNull Long value, @NotNull Long recipientId, @NotNull TurmsRequest notification), ofDataLong(@NotNull Long value, boolean forwardNotificationToRequesterOtherOnlineSessions, @NotNull Long recipientId, @NotNull TurmsRequest notification), ofDataLong(@NotNull Long value, boolean forwardDataForRecipientsToRequesterOtherOnlineSessions, @NotNull TurmsRequest notification), ofDataLong(@NotNull Long value, boolean forwardNotificationsToRequesterOtherOnlineSessions, @NotEmpty Set<Long> recipients, TurmsRequest notification), ofDataLongs(@NotNull Collection<Long> values), Notification(boolean forwardToRequesterOtherOnlineSessions, Set<Long> recipients, TurmsRequest notification), of(boolean forwardToRequesterOtherOnlineSessions, Set<Long> recipients, TurmsRequest notification), of(boolean forwardToRequesterOtherOnlineSessions, Long recipient, TurmsRequest notification), of(boolean forwardToRequesterOtherOnlineSessions, TurmsRequest notification), toString()*

Confirmed: no `String()` methods exist for either `RequestHandlerResult` or `Notification`.

## RequestHandlerResult.toString()
- [ ] The Go `RequestHandlerResult` struct has no `String()` method. The Java version returns `"RequestHandlerResult[code=" + code + ", reason='" + reason + "', response=" + ProtoFormatter.toLogString(response) + ", notifications=" + notifications + "]"`. This is a missing method port.

## Notification.toString()
- [ ] The Go `Notification` struct has no `String()` method. The Java version returns `"Notification{forwardToRequesterOtherOnlineSessions=" + forwardToRequesterOtherOnlineSessions + ", recipients=" + recipients + ", notification=" + ProtoFormatter.toLogString(notification) + "}"`. This is a missing method port.

# AdminController.java
*Checked methods: checkLoginNameAndPassword(), addAdmin(RequestContext requestContext, @RequestBody AddAdminDTO addAdminDTO), queryAdmins(@QueryParam(required = false), queryAdmins(@QueryParam(required = false), updateAdmins(RequestContext requestContext, Set<Long> ids, @RequestBody UpdateAdminDTO updateAdminDTO), deleteAdmins(RequestContext requestContext, Set<Long> ids)*

## CheckLoginNameAndPassword
- [ ] Method returns `nil` instead of an OK/success result. The Java version returns `HttpHandlerResult.OK` (a successful response), but the Go version returns `nil` with a comment "Requires further RequestContext integration".
- [ ] Method signature takes `(loginName string, password string)` parameters that don't exist in the Java version. The Java `checkLoginNameAndPassword()` takes no parameters — it's a HEAD endpoint that simply returns OK to confirm the admin API is reachable.

## AddAdmin
- [ ] Method body is entirely unimplemented — returns `nil` instead of calling `adminService.AuthAndAddAdmin`. The Java version calls `adminService.authAndAddAdmin(requesterId, loginName, password, roleIds, displayName, new Date(), false)` and wraps the result.
- [ ] Missing `displayName` parameter extraction from `AddAdminDTO`. The Java version passes `addAdminDTO.displayName()` to the service.
- [ ] Missing `roleIds` parameter extraction from `AddAdminDTO`. The Java version passes `addAdminDTO.roleIds()` to the service.
- [ ] Missing `new Date()` (current timestamp) and `false` arguments that the Java version passes to `authAndAddAdmin`.

## QueryAdmins
- [ ] Both `QueryAdminsWithQuery` and `QueryAdmins` methods return `nil` — no call to `adminService.queryAdmins()`. The Java version calls `adminService.queryAdmins(ids, loginNames, roleIds, 0, size)`.
- [ ] Missing `withPassword` parameter and the password-stripping logic. The Java version conditionally sets `admin.password` to `null` based on the `withPassword` boolean flag using `admin.toBuilder().password(null).build()`.
- [ ] Missing `getPageSize(size)` call. The Java version applies `size = getPageSize(size)` before passing it to the service.
- [ ] Missing `loginNames` and `roleIds` parameters from `QueryAdminsWithQuery` signature — the method accepts them but the Java query also supports these filters; however the Go signature does include them (ids, loginNames, roleIds) but the method body is `nil`.

## QueryAdmins (page variant)
- [ ] The Java version has a separate `queryAdmins` page variant (GET `/page`) that calls both `adminService.countAdmins(ids, roleIds)` and `adminService.queryAdmins(ids, loginNames, roleIds, page, size)`, returning a paginated result. The Go code has no corresponding page-based query method.

## UpdateAdmins
- [ ] Method body is entirely unimplemented — returns `nil` instead of calling `adminService.authAndUpdateAdmins`. The Java version calls `adminService.authAndUpdateAdmins(requesterId, ids, password, displayName, roleIds)`.

## DeleteAdmins
- [ ] Method body is entirely unimplemented — returns `nil` instead of calling `adminService.authAndDeleteAdmins`. The Java version calls `adminService.authAndDeleteAdmins(requesterId, ids)`.

# AdminPermissionController.java
*Checked methods: queryAdminPermissions()*

Now I have all the information needed to compare the two implementations. Let me analyze the differences.

## queryAdminPermissions

- [ ] **Method body is entirely unimplemented (returns `nil`)**: The Java version returns `ResponseDTO<Collection<PermissionDTO>>` containing a static list of all permissions (`ALL_PERMISSIONS`) mapped as `PermissionDTO(permission.getGroup(), permission)`. The Go version simply returns `nil` with a comment placeholder.

- [ ] **Missing `Group` field in `PermissionDTO`**: The Java `PermissionDTO` is a record with two fields: `String group` and `AdminPermission permission`. The Go `PermissionDTO` only has a single `Name string` field — it is missing the `Group` field entirely. Even if the method were implemented, it cannot produce the same response structure.

- [ ] **Missing group information on `AdminPermission`**: The Java `AdminPermission` enum has a `group` field (via `@Getter`) that maps each permission constant to its group (e.g., `USER_CREATE` → group `"USER"`). The Go `AdminPermission` is a plain `string` type with no group metadata. There is no way to derive the group from a Go `AdminPermission` value, so the equivalent of `permission.getGroup()` cannot be called.

- [ ] **Static `ALL_PERMISSIONS` list not used**: The Java version pre-computes `ALL_PERMISSIONS` as a static list of `PermissionDTO` objects at class load time. The Go side has `AllAdminPermissions` (a `[]AdminPermission` slice) available but the method does not reference or transform it into `[]PermissionDTO`.

# AdminRoleController.java
*Checked methods: addAdminRole(RequestContext requestContext, @RequestBody AddAdminRoleDTO addAdminRoleDTO), queryAdminRoles(@QueryParam(required = false), queryAdminRoles(@QueryParam(required = false), updateAdminRole(RequestContext requestContext, Set<Long> ids, @RequestBody UpdateAdminRoleDTO updateAdminRoleDTO), deleteAdminRoles(RequestContext requestContext, Set<Long> ids)*

Now I have a complete picture. Let me analyze the differences systematically.

## AddAdminRole

- [ ] **All method bodies are stubs returning `nil`**: The `AddAdminRole` method (line 77-79) returns `nil` without calling `c.adminRoleService.AuthAndAddAdminRole`. The Java version calls `adminRoleService.authAndAddAdminRole(requestContext.getRequesterId(), addAdminRoleDTO.id(), addAdminRoleDTO.name(), permissions, addAdminRoleDTO.rank())` and wraps the result with `HttpHandlerResult.okIfTruthy`. The Go method accepts the correct parameters but does not invoke the service.

- [ ] **Missing permission transformation**: In Java, `addAdminRoleDTO.permissions()` returns `Set<String>` and is transformed via `AdminPermission.matchPermissions(addAdminRoleDTO.permissions())` before passing to the service. The Go DTO already uses `[]permission.AdminPermission` (typed), so this step may be intentionally skipped. However, the Java code does a null-check before transforming (`permissions() == null ? null : matchPermissions(...)`), and the Go code doesn't handle this equivalence at all since the method body is empty.

## QueryAdminRolesWithQuery

- [ ] **Stub method returning `nil`**: The `QueryAdminRolesWithQuery` method (line 82-84) returns `nil` without calling `c.adminRoleService.QueryAdminRoles`. The Java version calls `adminRoleService.queryAdminRoles(ids, names, includedPermissions, ranks, 0, size)` (with page=0 for the non-paginated variant) and wraps with `HttpHandlerResult.okIfTruthy`.

- [ ] **Missing `getPageSize` call**: The Java version calls `size = getPageSize(size)` before using the size parameter. The Go method does not apply any page size normalization.

- [ ] **Missing non-paginated query with page=0**: The Java `@GetMapping` (non-paged) variant passes `0` as the page parameter to the service. The Go `QueryAdminRolesWithQuery` accepts a `page *int` parameter but doesn't enforce `0` for the non-paginated endpoint variant.

- [ ] **Missing `QueryAdminRoles` paginated endpoint**: The Go file has `QueryAdminRoles` (line 87-89) and `QueryAdminRolesWithQuery` (line 82-84) as separate methods, but the Java code has two overloaded `queryAdminRoles` methods: one non-paginated (`@GetMapping`) and one paginated (`@GetMapping("page")`). The Go `QueryAdminRolesWithQuery` seems to serve the non-paginated query but accepts `page *int` which should not exist for the non-paginated variant. The paginated variant (Java's `@GetMapping("page")`) which calls `countAdminRoles` + `queryAdminRoles` + `HttpHandlerResult.page` is not properly implemented.

## UpdateAdminRole

- [ ] **Stub method returning `nil`**: The `UpdateAdminRole` method (line 92-94) returns `nil` without calling `c.adminRoleService.AuthAndUpdateAdminRoles`. The Java version calls `adminRoleService.authAndUpdateAdminRoles(requesterId, ids, name, permissions, rank)` and wraps with `HttpHandlerResult.updateResult`.

- [ ] **Missing permission transformation**: Same as `AddAdminRole` — the Java code transforms `Set<String>` permissions via `AdminPermission.matchPermissions()` with a null-check. The Go method doesn't do this (though the Go DTO already uses `[]permission.AdminPermission`, the method body is empty).

## DeleteAdminRoles

- [ ] **Stub method returning `nil`**: The `DeleteAdminRoles` method (line 97-99) returns `nil` without calling `c.adminRoleService.AuthAndDeleteAdminRoles`. The Java version calls `adminRoleService.authAndDeleteAdminRoles(requesterId, ids)` and wraps with `HttpHandlerResult.deleteResult`.

# AddAdminDTO.java
*Checked methods: AddAdminDTO(String loginName, @SensitiveProperty(SensitiveProperty.Access.ALLOW_DESERIALIZATION), toString()*

Now I have all the information needed for the comparison. Let me summarize the bugs found.

## AddAdminDTO (Constructor / Field Mapping)

- [ ] **Missing field `displayName`**: The Java record has a `displayName` field, but the Go struct maps it to `Name` with JSON tag `"name"`. The Java JSON tag would be `"displayName"` (matching the record component name), but the Go struct uses `"name"`. If the API uses Jackson deserialization by record component names, the Go JSON tag should be `"displayName"` to match.
- [ ] **Extra field `ID`**: The Go struct has an `ID *int64` field that does not exist in the Java `AddAdminDTO` record. This is not present in the original Java source.
- [ ] **Missing `toString()` / password masking**: The Java `toString()` masks the `password` field with `SecurityValueConst.SENSITIVE_VALUE`. The Go struct has no `String()` method (or equivalent), meaning password values could be logged in plaintext if the struct is printed/logged via default formatting.

## SetConnection

- [ ] **Missing `ip` assignment**: The Java `setConnection(NetConnection connection, ByteArrayWrapper ip)` sets both `this.connection = connection` and `this.ip = ip`. The Go `SetConnection` receives an `ip string` parameter but never assigns it to the `UserSession.IP` field.

## ToString (UserSession)

- [ ] **Stub implementation**: The Java `toString()` returns a detailed string with all fields (`id`, `version`, `userId`, `deviceType`, `loginDate`, `loginLocation`, `isSessionOpen`, `connection`). The Go `ToString()` returns only `"UserSession{}"`, losing all diagnostic information.

# AddAdminRoleDTO.java
*Checked methods: AddAdminRoleDTO(Long id, String name, Set<String> permissions, Integer rank)*

## AddAdminRoleDTO

- [ ] **`Rank` field type mismatch**: In Java, `rank` is `Integer` (nullable/optional). In Go, `Rank` is `int` (non-pointer, zero-value default). This means the Go version cannot distinguish between "rank not provided" and "rank explicitly set to 0". It should be `*int` to match the optional semantics of Java's `Integer` wrapper type.

- [ ] **`Name` field missing nullability**: In Java, `name` is `String` which is nullable. In Go, `Name` is `string` (non-pointer). For a DTO used in an "add" operation this may be acceptable since `string` empty-check serves a similar purpose, but it differs from the Java semantics where `null` and `""` are distinguishable. Minor concern for an add DTO.

- [ ] **`Permissions` field type difference**: In Java, `permissions` is `Set<String>` (a set of strings). In Go, `Permissions` is `[]permission.AdminPermission` (a typed slice). This is a design-level choice (using typed enums vs raw strings), which is intentional and acceptable as a refactoring improvement, but worth noting the behavioral difference — the Java version accepts arbitrary strings while the Go version restricts to known `AdminPermission` enum values.

# UpdateAdminDTO.java
*Checked methods: UpdateAdminDTO(@SensitiveProperty(SensitiveProperty.Access.ALLOW_DESERIALIZATION), toString()*

Now I have a thorough understanding of both the Java and Go code. Let me analyze the two specific methods/behaviors:

1. **`@SensitiveProperty(SensitiveProperty.Access.ALLOW_DESERIALIZATION)`** on the `password` field
2. **`toString()`** method

Here's the analysis:

## `@SensitiveProperty(SensitiveProperty.Access.ALLOW_DESERIALIZATION)`

In Java, the `@SensitiveProperty` annotation on the `password` field with `ALLOW_DESERIALIZATION` means:
- **Serialization (Java → JSON response)**: The password field is **ignored** (not serialized). The `hasIgnoreMarker` for serialization returns `true` when `value != ALLOW_SERIALIZATION`, and `ALLOW_DESERIALIZATION != ALLOW_SERIALIZATION`, so the password is **never serialized** in responses.
- **Deserialization (JSON request → Java)**: The password field is **allowed** (can be deserialized). The `hasIgnoreMarker` for deserialization returns `true` when `value != ALLOW_DESERIALIZATION`, and since the value IS `ALLOW_DESERIALIZATION`, the marker returns `false`, so the field **is deserialized** from requests.

In Go, the `UpdateAdminDTO.Password` field has tag `json:"password,omitempty"`, which means:
- **Serialization**: The password **will be serialized** in responses (when non-nil). This is a behavioral difference — the Java version never serializes the password in responses.
- **Deserialization**: The password is correctly deserialized from requests.

**Bug**: The Go struct lacks write-only semantics for the password field. In Java, `@SensitiveProperty(ALLOW_DESERIALIZATION)` prevents the password from ever appearing in JSON responses. In Go, `json:"password,omitempty"` will include the password in any JSON serialization when it's non-nil.

## `toString()`

The Java `toString()` explicitly masks the password with `SecurityValueConst.SENSITIVE_VALUE` (`"***"`), while including `displayName` and `roleIds` in the output. The Go struct has no `String()` method at all, meaning default formatting (e.g., `fmt.Sprintf("%v", dto)`) would print all fields including the raw password.

**Bug**: Missing `String()` / `fmt.Stringer` implementation that masks the password.

---

## UpdateAdminDTO

- [ ] **Missing write-only semantics for `password` field**: The Java `@SensitiveProperty(ALLOW_DESERIALIZATION)` prevents the password from being serialized to JSON responses (write-only). The Go field `json:"password,omitempty"` allows the password to be serialized in responses when non-nil, leaking the password in any JSON response. The Go code needs a mechanism to prevent password serialization (e.g., a custom `MarshalJSON` method or using a separate input-only DTO type).
- [ ] **Missing `String()` method with password masking**: The Java `toString()` masks the password with `"***"` (`SecurityValueConst.SENSITIVE_VALUE`). The Go `UpdateAdminDTO` has no `String()` method, so `fmt.Sprintf("%+v", dto)` or `%#v` will print the raw password in plaintext. This is a security risk in logging/debugging scenarios.

# UpdateAdminRoleDTO.java
*Checked methods: UpdateAdminRoleDTO(String name, Set<String> permissions, Integer rank)*

## UpdateAdminRoleDTO

- [ ] **`permissions` field type mismatch**: The Java version uses `Set<String>` for `permissions`, while the Go version uses `[]permission.AdminPermission` (a typed slice). The Java version treats permissions as raw strings (`Set<String>`), whereas the Go version deserializes them into a custom enum type `permission.AdminPermission`. This is a behavioral difference — in Java the raw string values are preserved, while in Go they are parsed/validated as enum values during deserialization. If any invalid permission string is sent that exists in neither mapping, the Go version would fail at deserialization time whereas the Java version would accept the raw string.

- [ ] **`name` and `rank` nullability mismatch**: In the Java record, all three fields (`name`, `permissions`, `rank`) are non-pointer nullable types — `String` can be `null`, `Set<String>` can be `null`, and `Integer` (boxed) can be `null`. This signals that all fields are optional for a partial update. In Go, `Name` is `*string` (correctly optional via pointer), `Rank` is `*int` (correctly optional via pointer), but `Permissions` is `[]permission.AdminPermission` — a non-pointer slice. In Go, a nil/omitted JSON field for a slice results in `nil`, and a non-nil empty array results in `[]T{}`. This distinction means the Go code **can** distinguish "not provided" (`nil`) from "explicitly set to empty" (`[]permission.AdminPermission{}`), which is functionally equivalent to Java's `null` vs. `Set.of()`. However, the type difference from `Set<String>` to `[]permission.AdminPermission` remains a logic divergence as noted above.

- [ ] **No deduplication of `permissions`**: The Java version uses `Set<String>`, which inherently deduplicates entries. The Go version uses a slice `[]permission.AdminPermission`, which allows duplicate values. If the caller sends duplicate permission strings, the Java version would silently deduplicate, while the Go version would preserve duplicates. This could lead to unexpected behavior downstream if the consumer expects unique permissions.

# PermissionDTO.java
*Checked methods: PermissionDTO(String group, AdminPermission permission)*

## PermissionDTO

- [ ] **Missing `group` field**: The Java record has a `String group` field, but the Go struct does not have a `Group` field at all.
- [ ] **Missing `permission` field**: The Java record has an `AdminPermission permission` field, but the Go struct does not have a `Permission` field.
- [ ] **Spurious `name` field**: The Go struct has a `Name string` field that does not exist in the Java record. The Java record only has `group` and `permission`.

# AdminRoleRepository.java
*Checked methods: updateAdminRoles(Set<Long> roleIds, String newName, @Nullable Set<AdminPermission> permissions, @Nullable Integer rank), countAdminRoles(@Nullable Set<Long> ids, @Nullable Set<String> names, @Nullable Set<AdminPermission> includedPermissions, @Nullable Set<Integer> ranks), findAdminRoles(@Nullable Set<Long> roleIds, @Nullable Set<String> names, @Nullable Set<AdminPermission> includedPermissions, @Nullable Set<Integer> ranks, @Nullable Integer page, @Nullable Integer size), findAdminRolesByIdsAndRankGreaterThan(@NotNull Collection<Long> roleIds, @Nullable Integer rankGreaterThan), findHighestRankByRoleIds(Set<Long> roleIds)*

Now let me verify how the Java `AdminPermission` enum is serialized to MongoDB. Since the Go version uses `AdminPermission` as a `string` type, and the Java version stores permissions as enum strings (via `encodeAsStrings`), they should be compatible. The key question is about the `updateAdminRoles` method's handling of permissions when converting to `$in`.

Let me also verify what happens in the `updateAdminRoles` with the `$in` filter on `includedPermissions` in `buildFilter` — this is for count/find, not update. In the update, permissions are set directly.

Now I have all the information to provide the complete comparison. Here are the findings:

## updateAdminRoles

- [ ] **`setIfNotNullForEnumStrings` semantics mismatch for permissions**: In Java, `setIfNotNullForEnumStrings` only sets the field if `permissions != null && !permissions.isEmpty()`. If `permissions` is an empty set, it does NOT add any update for permissions. In Go, when `permissions` is a non-nil empty slice (`[]permission.AdminPermission{}`), it adds `$unset` for permissions. This means the Go code will unset/clear the permissions field when passed an empty slice, while the Java code simply skips the update entirely for an empty collection.

## countAdminRoles

- [ ] **`inIfNotNullForEnumStrings` uses `$in` on an array field, but Java queries individual string elements**: In Java, `inIfNotNullForEnumStrings` converts enum values to their string names and uses `$in` — `{perm: {$in: ["USER_CREATE", "USER_DELETE"]}}`. In Go, `buildFilter` does the same: `{perm: {$in: [permission1, permission2]}}`. Since `AdminPermission` is a `string` type in Go, the BSON serialization should produce the same string values. However, there is a subtlety: the Go `$in` on `includedPermissions` will serialize `AdminPermission` string values, which should match the Java `encodeAsStrings` behavior. This appears correct as long as the Go BSON driver serializes `[]AdminPermission` as `[]string`. This is likely fine since `AdminPermission` is `type AdminPermission string`. **No bug here** — the permissions are stored as strings in MongoDB in both cases.

## findAdminRoles

- [ ] **Pagination skip/limit logic differs from Java's `paginateIfNotNull`**: In Java's `paginateIfNotNull`, if `size != null && page == null`, it defaults `page` to `0`, so `skip = size * 0 = 0`. In Go, if `size != nil && page == nil`, the `SetSkip` is not called at all (only `SetLimit` is set), which effectively means skip=0. The behavior is equivalent. **No bug here for this specific case.**

## findAdminRolesByIdsAndRankGreaterThan

- [ ] **Missing filter when `roleIds` is empty**: In Java, `.in(DomainFieldName.ID, roleIds)` is called unconditionally — even if `roleIds` is empty, it will create an `$in` with an empty array, which matches no documents. In Go, when `roleIds` is empty (`len(roleIds) == 0`), the `$in` filter for IDs is skipped entirely, meaning the query has no ID filter and could match ALL documents. This is a behavioral difference: Java would return an empty result, Go could return all roles with rank > X.

## findHighestRankByRoleIds

- [ ] **Different implementation strategy — Go uses sort+limit instead of projection+in-memory max**: Java fetches only the `rank` field (via `.include(AdminRole.Fields.RANK)`) for all matching roles, then iterates in memory to find the max. Go uses `FindOne` with `sort(rank: -1)` to get the highest rank directly from MongoDB. The Go approach is actually more efficient, but there is a subtle difference: the Java version initializes `highestRank = 0` and returns `0` if all roles have null/zero rank, while the Go version returns the actual highest rank from the database. If a role has a negative rank (unlikely but possible), Java would return `0` (since `rank > highestRank` would be false for negative values), while Go would correctly return the negative rank. **This is a minor behavioral edge case.**

- [ ] **Java uses projection to only fetch the `rank` field; Go fetches the entire document**: Java does `.include(AdminRole.Fields.RANK)` which is a projection that only returns the rank field. Go's `FindOne` fetches the full document. This is a performance difference but not a logic bug per se.

- [ ] **Java returns `Mono.empty()` when roles list is empty (after query); Go returns `nil` early when `roleIds` is empty**: Java queries MongoDB even with an empty `roleIds` set (which produces an `$in: []` matching nothing), then checks if the result is empty. Go short-circuits with `nil` before querying. The net effect is the same (nil/empty result), but via different paths. **No functional bug.**

## updateAdminRoles

- [ ] **Java applies filter unconditionally; Go skips `$in` when `roleIds` is empty**: In Java, `.in(DomainFieldName.ID, roleIds)` is called unconditionally regardless of whether `roleIds` is empty. With an empty set, this creates `{_id: {$in: []}}` which matches no documents, so the update becomes a no-op. In Go, when `roleIds` is empty, the `$in` filter is skipped entirely, resulting in an empty filter `{}`, which would match ALL documents. This is a critical bug — an update with empty roleIds would update every admin role in the database.

## buildFilter (used by countAdminRoles and findAdminRoles)

- [ ] **`$in` on array field `perm` for `includedPermissions` uses element match semantics that differ from Java**: In Java, `inIfNotNullForEnumStrings(AdminRole.Fields.PERMISSIONS, includedPermissions)` generates `{perm: {$in: ["STRING1", "STRING2"]}}`. MongoDB's `$in` on an array field checks if any element of the array matches any value in the `$in` list. In Go, `buildFilter` generates `{perm: {$in: includedPermissions}}` where `includedPermissions` is `[]permission.AdminPermission`. Since `AdminPermission` is `type AdminPermission string`, the Go BSON driver should serialize these as strings. The behavior should be equivalent. **No bug here.**

- [ ] **Java's `inIfNotNull` skips on null or empty collections; Go's `buildFilter` skips on empty slices but accepts nil slices**: Java uses `collection != null && !collection.isEmpty()`. Go uses `len(ids) > 0`. For nil slices in Go, `len(nil) == 0`, so the filter is skipped — matching Java's null check. For empty slices, `len([]int64{}) == 0`, also skipped — matching Java's `isEmpty()` check. **No bug here.**

Here is the consolidated bug report:

---

## updateAdminRoles

- [ ] **Critical: Empty `roleIds` slice causes update of ALL documents**. Java calls `.in(DomainFieldName.ID, roleIds)` unconditionally, which with an empty set produces `{_id: {$in: []}}` matching no documents. Go skips the `$in` filter when `len(roleIds) == 0`, producing an empty filter `{}` that matches ALL documents. This could update every admin role in the database.

- [ ] **Behavioral difference for empty `permissions` slice**: Java's `setIfNotNullForEnumStrings` skips the update when permissions is non-null but empty (`!value.isEmpty()` check at line 80 of Update.java). Go code explicitly `$unset`s permissions when `len(permissions) == 0` (lines 53-54). This means passing an empty permissions slice in Go will clear the permissions field, while in Java it leaves permissions unchanged.

## findAdminRolesByIdsAndRankGreaterThan

- [ ] **Missing `$in` filter when `roleIds` is empty**: Java applies `.in(DomainFieldName.ID, roleIds)` unconditionally. With an empty collection, this creates `{_id: {$in: []}}` which matches nothing. Go skips the filter when `len(roleIds) == 0`, leaving no ID constraint, so `findAdminRolesByIdsAndRankGreaterThan(ctx, []int64{}, &rank)` would return ALL admin roles with rank > X instead of returning an empty list.

## findHighestRankByRoleIds

- [ ] **Returns `0` for the "no matching roles" case vs `nil`**: Java initializes `highestRank = 0` and iterates all results. If roles exist but all have rank 0, it returns 0. If no roles are found, it returns `Mono.empty()` (nil). Go returns `nil` when no documents found (correct) and otherwise returns the actual max rank. However, Go also uses `FindOne` with sort instead of fetching all and computing max in-memory like Java does. While the Go approach is more efficient, there's an edge case: Java initializes `highestRank = 0`, so if all returned roles have a rank ≤ 0, Java would return 0 rather than the actual highest rank. Go would return the correct actual highest rank (e.g., -1). This is a behavioral difference but likely inconsequential since ranks are presumably always positive.

# AdminRoleService.java
*Checked methods: authAndAddAdminRole(@NotNull Long requesterId, @NotNull Long roleId, @NotNull @NoWhitespace @Size( min = MIN_ROLE_NAME_LIMIT, max = MAX_ROLE_NAME_LIMIT), addAdminRole(@NotNull Long roleId, @NotNull @NoWhitespace @Size( min = MIN_ROLE_NAME_LIMIT, max = MAX_ROLE_NAME_LIMIT), authAndDeleteAdminRoles(@NotNull Long requesterId, @NotEmpty Set<Long> roleIds), deleteAdminRoles(@NotEmpty Set<Long> roleIds), authAndUpdateAdminRoles(@NotNull Long requesterId, @NotEmpty Set<Long> roleIds, @Nullable @NoWhitespace @Size( min = MIN_ROLE_NAME_LIMIT, max = MAX_ROLE_NAME_LIMIT), updateAdminRole(@NotEmpty Set<Long> roleIds, @Nullable @NoWhitespace @Size( min = MIN_ROLE_NAME_LIMIT, max = MAX_ROLE_NAME_LIMIT), queryAdminRoles(@Nullable Set<Long> ids, @Nullable Set<String> names, @Nullable Set<AdminPermission> includedPermissions, @Nullable Set<Integer> ranks, @Nullable Integer page, @Nullable Integer size), queryAndCacheRolesByRoleIdsAndRankGreaterThan(@NotNull Collection<Long> roleIds, @NotNull Integer rankGreaterThan), countAdminRoles(@Nullable Set<Long> ids, @Nullable Set<String> names, @Nullable Set<AdminPermission> includedPermissions, @Nullable Set<Integer> ranks), queryHighestRankByAdminId(@NotNull Long adminId), queryHighestRankByRoleIds(@NotNull Set<Long> roleIds), isAdminRankHigherThanRank(@NotNull Long adminId, @NotNull Integer rank), queryPermissions(@NotNull Long adminId)*

Now I have both files. Let me systematically compare each method.

## authAndAddAdminRole

- [ ] Missing the `roleId == ADMIN_ROLE_ROOT_ID` check: Java rejects root role ID with an unauthorized error, Go does not.
- [ ] Missing the `IsAdminRankHigherThanRank` authorization check: Java verifies the requester's rank is higher than the target rank, and returns an error if the requester doesn't exist (via `switchIfEmpty(adminService.errorRequesterNotExist())`). Go skips all auth and goes straight to `AddAdminRole`.
- [ ] Missing the requester permissions containment check: Java queries the requester's permissions and verifies they contain all the requested permissions for the new role, returning an error with the missing permissions listed. Go omits this entirely.
- [ ] Missing validation: Java validates `roleId` not null, `name` not null / no whitespace / length 1-32, and `permissions` not empty. Go performs none of these validations.
- [ ] `roleId` parameter semantics differ: Go's `roleId` is `*int64` (optional, generates one if nil), while Java's is `@NotNull Long roleId` (required, never auto-generated). The auto-generation fallback in Go does not exist in the Java version.

## addAdminRole

- [ ] Missing the `roleId == ADMIN_ROLE_ROOT_ID` check: Java rejects root role ID. Go does not.
- [ ] Missing validation: Java validates `roleId` not null, `name` not null / no whitespace / length 1-32, `permissions` not empty, and `rank` not null. Go performs none of these validations.
- [ ] Missing the in-memory cache update (`idToRole.put`): After inserting, Java updates an in-memory cache `idToRole` with the new role. Go does not maintain or update any equivalent cache.

## authAndDeleteAdminRoles

- [ ] Missing the `checkIfAllowedToManageRoles` authorization: Java verifies the requester's rank is higher than all target roles before deleting. Go simply delegates to `DeleteAdminRoles` with no auth.
- [ ] Missing validation: Java validates `roleIds` not empty. Go does not validate input.

## deleteAdminRoles

- [ ] Missing the `ADMIN_ROLE_ROOT_ID` containment check: Java rejects deletion if roleIds contains the root role ID ("The root admin is reserved and cannot be deleted"). Go does not.
- [ ] Missing validation: Java validates `roleIds` not empty and does not contain root ID. Go does not validate.
- [ ] Missing the in-memory cache invalidation (`idToRole.remove`): After deleting, Java removes each deleted role from the `idToRole` cache. Go does not update any equivalent cache.

## authAndUpdateAdminRoles

- [ ] Missing the `checkIfAllowedToManageRoles` authorization: Java verifies the requester's rank is higher than all target roles before updating. Go delegates directly to `UpdateAdminRole` with no auth.
- [ ] Missing the requester permissions containment check for non-null `permissions`: Java queries the requester's permissions and verifies they contain all requested permissions. Go omits this entirely.
- [ ] Missing validation: Java validates `roleIds` not empty, `newName` no whitespace / length 1-32. Go does not validate.
- [ ] The `noOpResult` early-return when all update fields are null is bypassed: In Java, if `permissions == null`, it calls `updateAdminRole` directly (which checks `areAllFalsy`), but the auth flow goes through `checkIfAllowedToManageRoles` first. Go skips all auth and calls `UpdateAdminRole` directly, which also lacks the `areAllFalsy` early return (though the repo may handle it).

## updateAdminRole

- [ ] Missing the `ADMIN_ROLE_ROOT_ID` containment check: Java rejects update if roleIds contains the root role ID ("The root admin is reserved and cannot be updated"). Go does not.
- [ ] Missing validation: Java validates `roleIds` not empty and does not contain root ID, `newName` no whitespace / length 1-32. Go does not validate.
- [ ] Missing the `areAllFalsy` early return: Java returns `ACKNOWLEDGED_UPDATE_RESULT` (a no-op success) if `newName`, `permissions`, and `rank` are all null/falsy. Go always calls the repo update regardless.
- [ ] Missing the in-memory cache invalidation (`idToRole.remove`): After updating, Java removes each updated role from the `idToRole` cache. Go does not update any equivalent cache.

## queryAdminRoles

- [ ] Missing the root role inclusion logic: Java checks `isRootRoleQualified(ids, names, includedPermissions, ranks)` and prepends the root role to results if it qualifies. Go returns only the repository results, never including the root role.

## queryAndCacheRolesByRoleIdsAndRankGreaterThan

- [ ] Missing the `ADMIN_ROLE_ROOT_ID` special handling: Java checks if roleIds contains the root role ID and, if so, handles the root role separately (checking its rank against `rankGreaterThan`), removes it from the query set, and prepends it to results if qualified. Go passes all roleIds directly to the repo without this logic.
- [ ] Missing validation: Java validates `roleIds` not null and `rankGreaterThan` not null. Go does not validate.
- [ ] Missing the early return for empty `roleIds`: Java returns `Flux.empty()` if roleIds is empty. Go does not have this check.
- [ ] Missing the in-memory cache update (`idToRole.put`): Java caches each returned role in `idToRole`. Go does not update any equivalent cache.

## countAdminRoles

- [ ] Missing the `+1` for the builtin root role: Java adds 1 to the count (`number + 1`) because the root role is not stored in the database. Go returns the raw repository count without this adjustment.

## queryHighestRankByAdminId

- [ ] Stub implementation returning `nil, nil`: Java queries `adminService.queryRoleIdsByAdminId(adminId)` then calls `queryHighestRankByRoleIds`. Go returns `nil, nil` as a TODO stub, meaning all dependent methods (`IsAdminRankHigherThanRank`, `queryPermissions`, auth methods) will fail silently.

## queryHighestRankByRoleIds

- [ ] Missing the `ADMIN_ROLE_ROOT_ID` shortcut: Java checks if roleIds contains the root role ID and returns the root role's rank immediately. Go always queries the repository.
- [ ] Missing the empty set early return: Java returns `Mono.empty()` if roleIds is empty. Go does not have this check.

## isAdminRankHigherThanRank

- [ ] Missing validation: Java validates `rank` not null. Go does not validate.
- [ ] Behavior difference when admin has no roles: In Java, `queryHighestRankByAdminId` returns `Mono.empty()` if the admin has no roles, which causes `isAdminRankHigherThanRank` to return `Mono.empty()` (an empty signal, not `false`). In Go, when `highest == nil` (from the stub), it returns `false, nil`. This changes the behavior in `authAndAddAdminRole` where Java uses `switchIfEmpty` to detect a non-existent requester.

## queryPermissions

- [ ] Stub implementation returning `nil, nil`: Java queries `adminService.queryRoleIdsByAdminId(adminId)` then aggregates permissions from all roles. Go returns `nil, nil` as a TODO stub, meaning all permission checks in auth methods will silently pass or fail incorrectly.

# AdminService.java
*Checked methods: queryRoleIdsByAdminIds(@NotEmpty Set<Long> adminIds), authAndAddAdmin(@NotNull Long requesterId, @Nullable @NoWhitespace @Size( min = MIN_LOGIN_NAME_LIMIT, max = MAX_LOGIN_NAME_LIMIT), addAdmin(@Nullable Long id, @Nullable @NoWhitespace @Size( min = MIN_LOGIN_NAME_LIMIT, max = MAX_LOGIN_NAME_LIMIT), queryAdmins(@Nullable Collection<Long> ids, @Nullable Collection<String> loginNames, @Nullable Collection<Long> roleIds, @Nullable Integer page, @Nullable Integer size), authAndDeleteAdmins(@NotNull Long requesterId, @NotEmpty Set<Long> adminIds), authAndUpdateAdmins(@NotNull Long requesterId, @NotEmpty Set<Long> targetAdminIds, @Nullable @NoWhitespace @Size( min = MIN_PASSWORD_LIMIT, max = MAX_PASSWORD_LIMIT), updateAdmins(@NotEmpty Set<Long> targetAdminIds, @Nullable @NoWhitespace @Size( min = MIN_PASSWORD_LIMIT, max = MAX_PASSWORD_LIMIT), countAdmins(@Nullable Set<Long> ids, @Nullable Set<Long> roleIds), errorRequesterNotExist()*

## QueryRoleIdsByAdminIds

- [ ] **Missing in-memory cache lookup optimization**: The Java version first attempts to resolve all admin role IDs from the in-memory `idToAdmin` map. Only when an admin is not found in the cache does it fall back to querying the database. The Go version always queries the database directly, completely bypassing any cache-first logic.
- [ ] **Missing early-return on all-found path**: In Java, if all admins are found in the in-memory map, it returns immediately without any DB query. The Go version always performs a DB query regardless.

## AuthAndAddAdmin

- [ ] **Missing requester existence check**: Java calls `checkIfAllowedToAddRolesToAdmins` which zips `adminRoleService.queryHighestRankByAdminId(requesterId).switchIfEmpty(errorRequesterNotExist())` — if the requester doesn't exist, it returns an unauthorized error. The Go version has a comment `// auth omitted` and skips all authorization.
- [ ] **Missing root role ID validation**: Java validates that `roleIds` does not contain `ADMIN_ROLE_ROOT_ID` (`Validator.notContains(roleIds, ADMIN_ROLE_ROOT_ID, ...)`). The Go version does not check for this.
- [ ] **Missing rank-based permission check**: Java calls `checkIfAllowedToAddRolesToAdmins` which verifies the target role ranks are all lower than the requester's rank. The Go version skips this entirely.
- [ ] **Missing validation of roleIds existence**: Java queries `adminRoleService.queryAndCacheRolesByRoleIds(roleIds)` and verifies all requested role IDs actually exist, returning an error with the missing IDs if not. The Go version does not validate role existence.

## AddAdmin

- [ ] **Missing default random login name**: When `loginName` is empty/null, Java generates `RandomStringUtils.randomAlphabetic(16)` as a default. The Go version uses the empty string directly with no fallback.
- [ ] **Missing default random password**: When `rawPassword` is empty/null, Java generates `RandomStringUtils.randomAlphabetic(10)` and encodes that. The Go version would pass an empty string to bcrypt, which is different behavior.
- [ ] **Missing default displayName fallback to loginName**: In Java, if `displayName` is null/empty, it falls back to `loginName`. The Go version does not implement this fallback.
- [ ] **Missing upsert support**: Java accepts a `boolean upsert` parameter and conditionally calls `adminRepository.upsert(admin)` vs `adminRepository.insert(admin)`. The Go version always calls `s.repo.Insert` with no upsert option.
- [ ] **Missing `registrationDate` parameter**: Java accepts an optional `registrationDate` parameter and only defaults to `new Date()` if null. The Go version always uses `time.Now()`, ignoring any passed-in value. The Go function signature doesn't even accept a `registrationDate` parameter.

## QueryAdmins

- [ ] No bugs found — the Go version correctly delegates to the repository with the same parameters.

## AuthAndDeleteAdmins

- [ ] **Missing requester existence check**: Java uses `adminRoleService.queryHighestRankByAdminId(requesterId).switchIfEmpty(errorRequesterNotExist())` to verify the requester exists. The Go version does not check if the requester exists.
- [ ] **Missing root admin deletion guard**: Java validates that `adminIds` does not contain `ROOT_ADMIN_ID` with the message "The root admin is reserved and cannot be deleted". The Go version has no such check.
- [ ] **Missing rank-based permission check**: Java calls `checkIfAllowedToManageAdmins` which queries target admins' roles and verifies the requester has a higher rank than all target admins. The Go version skips this entirely and directly deletes.
- [ ] **Missing in-memory cache invalidation**: Java's `deleteAdmins` removes deleted admin IDs from the `idToAdmin` map. The Go version has no equivalent cache invalidation.

## AuthAndUpdateAdmins

- [ ] **Missing requester existence check**: Java calls `adminRoleService.queryHighestRankByAdminId(requesterId).switchIfEmpty(errorRequesterNotExist())`. The Go version has no requester existence check.
- [ ] **Missing early return when all params are null**: Java checks `if (Validator.areAllNull(rawPassword, displayName) && noRoleIds)` and returns `ACKNOWLEDGED_UPDATE_RESULT` immediately. The Go version does not check for this no-op case.
- [ ] **Missing self-role-update guard**: Java checks `if (targetAdminIds.contains(requesterId))` when `roleIds` is non-null and returns `MONO_ERROR_UPDATE_ONE_OWN_ROLE_ID` ("It is forbidden to update one's own role ID"). The Go version has no such check.
- [ ] **Missing rank-based permission check for managing admins**: When `roleIds` is null, Java calls `checkIfAllowedToManageAdmins` to verify the requester outranks all target admins. The Go version skips this.
- [ ] **Missing rank-based permission check for adding roles**: When `roleIds` is non-null, Java calls `checkIfAllowedToAddRolesToAdmins` followed by `checkIfAllowedToManageAdmins`. The Go version skips both checks.
- [ ] **Missing roleIds in the update call**: In the `noRoleIds` branch, Java calls `updateAdmins(targetAdminIds, rawPassword, displayName, null)` — passing null for roleIds. But in the `roleIds` branch, Java also calls `updateAdmins(targetAdminIds, rawPassword, displayName, null)` — notably also passing null for roleIds (the role check happens separately). The Go version passes `roleIds` through directly to `UpdateAdmins`, which means it updates role IDs in the database without having verified the requester's rank vs. the new roles' ranks.

## UpdateAdmins

- [ ] **Missing early return when all params are null**: Java checks `if (Validator.areAllNull(rawPassword, displayName) && CollectionUtil.isEmpty(roleIds))` and returns `ACKNOWLEDGED_UPDATE_RESULT`. The Go version does not handle this no-op case.
- [ ] **Missing in-memory cache invalidation**: Java iterates over `targetAdminIds` and removes each from `idToAdmin` in a `doOnNext` callback. The Go version has no equivalent cache invalidation.

## CountAdmins

- [ ] No bugs found — the Go version correctly delegates to the repository with the same parameters.

## ErrorRequesterNotExist

- [ ] **Error type differs**: Java returns `ResponseException` with `ResponseStatusCode.UNAUTHORIZED` and a message "The requester does not exist". The Go version returns a plain `errors.New("requester does not exist")` without any status code equivalent. While this is a minor adaptation, the Java version specifically uses an UNAUTHORIZED status code that would drive HTTP response behavior — the Go plain error loses this semantic information.

# IpBlocklistController.java
*Checked methods: addBlockedIps(@RequestBody AddBlockedIpsDTO addBlockedIpsDTO), queryBlockedIps(Set<String> ids), queryBlockedIps(int page, @QueryParam(required = false), deleteBlockedIps(@QueryParam(required = false)*

Now I have a complete picture. The Go code is entirely stub implementations with empty method bodies.

## AddBlockedIps

- [ ] Method body is completely empty — missing all core logic: does not call `blocklistService.BlockIpStrings()` with the DTO's `IDs` and `BlockDurationMillis` fields, and does not return an `HttpHandlerResult.okIfTruthy()`-style response.

## QueryBlockedIpsByIds

- [ ] Method body is completely empty — missing all core logic: does not call `blocklistService.GetBlockedIpStrings(ids)`, does not convert `BlockedClient` results to `BlockedIpDTO` (IP bytes to string + block end time), and does not return an `HttpHandlerResult.okIfTruthy()` response with the collection.

## QueryBlockedIpsByPage

- [ ] Method body is completely empty — missing all core logic: does not resolve page size via `getPageSize(size)`, does not call `blocklistService.CountBlockIps()` for the total count, does not call `blocklistService.GetBlockedIps(page, size)` to fetch the page, does not convert results to `BlockedIpDTO`, and does not return an `HttpHandlerResult.page()` paginated response.

## DeleteBlockedIps

- [ ] Method body is completely empty — missing all core logic: does not handle the `deleteAll` flag to call `blocklistService.UnblockAllIps()`, does not handle the `ids` set to call `blocklistService.UnblockIpStrings(ids)` when `deleteAll` is false and `ids` is non-empty, and does not return an `HttpHandlerResult.okIfTruthy()` response.

# UserBlocklistController.java
*Checked methods: addBlockedUserIds(@RequestBody AddBlockedUserIdsDTO addBlockedUserIdsDTO), queryBlockedUsers(Set<Long> ids), queryBlockedUsers(int page, @QueryParam(required = false), deleteBlockedUserIds(@QueryParam(required = false)*

## addBlockedUserIds(@RequestBody AddBlockedUserIdsDTO addBlockedUserIdsDTO)

- [ ] **Method body is completely empty.** The Go method `AddBlockedUserIds()` has no implementation at all — no parameters, no logic. The Java version takes an `AddBlockedUserIdsDTO`, calls `blocklistService.blockUserIds(addBlockedUserIdsDTO.ids(), addBlockedUserIdsDTO.blockDurationMillis())`, and returns `HttpHandlerResult.okIfTruthy(result)`. The Go stub accepts no arguments and performs no service call.

## queryBlockedUsers(Set\<Long\> ids)

- [ ] **Method is entirely missing.** There is no Go method corresponding to the Java non-paginated `queryBlockedUsers(Set<Long> ids)`. The Java version calls `blocklistService.getBlockedUsers(ids)`, maps results to `BlockedUserDTO` via `clients2users()`, and returns them. No equivalent exists in the Go controller.

## queryBlockedUsers(int page, @QueryParam(required = false) Integer size)

- [ ] **Method is entirely missing.** There is no Go method corresponding to the Java paginated `queryBlockedUsers(int page, Integer size)`. The Java version resolves the page size via `getPageSize(size)`, calls `blocklistService.countBlockUsers()` and `blocklistService.getBlockedUsers(page, size)`, then returns `HttpHandlerResult.page(blockUserCount, clients2users(blockedUsers))`. No equivalent exists in the Go controller.

## deleteBlockedUserIds(@QueryParam(required = false) Set\<Long\> ids, boolean deleteAll)

- [ ] **Method body is completely empty.** The Go method `DeleteBlockedUserIds()` has no implementation — no parameters, no logic. The Java version accepts optional `ids` and a `deleteAll` flag. When `deleteAll` is true, it calls `blocklistService.unblockAllUserIds()`. When `deleteAll` is false and `ids` is non-empty, it calls `blocklistService.unblockUserIds(ids)`. If neither condition is met, it returns an empty result. The Go stub performs none of this logic.

# BlockedClientSerializer.java
*Checked methods: serialize(BlockedClient value, JsonGenerator gen, SerializerProvider provider)*

## `serialize(BlockedClient value, JsonGenerator gen, SerializerProvider provider)`

- [ ] **The Go code is a stub with no implementation.** The `BulkRequest.Serialize()` method at line 8-10 of `elasticsearch_model.go` is annotated with `@MappedFrom serialize(BlockedClient value, JsonGenerator gen, SerializerProvider provider)` but the method body is simply `return nil`. None of the Java serializer logic is implemented.
- [ ] **Missing BlockedClient struct/type.** The Java code operates on a `BlockedClient` domain object, but the Go code has no `BlockedClient` type. The serialization is incorrectly placed on `BulkRequest` instead.
- [ ] **Missing ID field serialization with type discrimination.** The Java code checks whether the `id` is a `Long` (user ID) or a `ByteArrayWrapper` (IP address bytes), writing a number field for user IDs and converting IP bytes to a string via `InetAddressUtil.ipBytesToString()`. The Go code has none of this logic.
- [ ] **Missing `blockEndTime` field serialization.** The Java code writes a `blockEndTime` string field by converting `blockEndTimeMillis()` using `DateTimeUtil.toStr()`. The Go code does not serialize this field at all.

# MemberController.java
*Checked methods: queryMembers(), removeMembers(List<String> ids), addMember(@RequestBody AddMemberDTO addMemberDTO), updateMember(String id, @RequestBody UpdateMemberDTO updateMemberDTO), queryLeader(), electNewLeader(@QueryParam(required = false)*

## QueryMembers
- [ ] Method is an empty stub — no logic at all. Java version calls `discoveryService.getAllKnownMembers().values()` and wraps in `HttpHandlerResult.okIfTruthy()`. Go version has no `DiscoveryService` field, no method body, and no return value.
- [ ] Missing `DiscoveryService` field on `MemberController` struct. Java constructor injects `Node` and extracts `discoveryService` from it.

## RemoveMembers
- [ ] Method is an empty stub — no logic at all. Java version calls `discoveryService.unregisterMembers(CollectionUtil.newSet(ids))` and returns a `DeleteResultDTO`. Go version has no parameters, no body, and no return value.
- [ ] Missing `ids` parameter (`List<String> ids` in Java). The Go method signature is `RemoveMembers()` with no arguments.

## AddMember
- [ ] Method is an empty stub — no logic at all. Java version validates nodeType vs isLeaderEligible, constructs a full `Member` object with all fields (clusterId, nodeId, zone, name, nodeType, NodeVersion, isSeed, isLeaderEligible, registrationDate, priority, memberHost, memberPort, adminApiAddress, wsAddress, tcpAddress, udpAddress, false for isLeader, isActive, isHealthy), and registers it via `discoveryService.registerMember()`. Go version has no body.
- [ ] Missing `AddMemberDTO` parameter in method signature.
- [ ] Missing validation: Java checks `if (nodeType != NodeType.SERVICE && addMemberDTO.isLeaderEligible())` and returns an error `"Only turms-service servers can be the leader"`. Entirely absent in Go.
- [ ] Missing `clusterId` resolution: Java calls `discoveryService.getLocalMember().getClusterId()` to get the cluster ID for the new member.
- [ ] Missing `Member` construction with all 17+ fields from the DTO, including the hardcoded `false` for `isLeader` and the `NodeVersion.parse()` call.

## UpdateMember
- [ ] Method is an empty stub — no logic at all. Java version calls `discoveryService.updateMemberInfo(id, zone, name, isSeed, isLeaderEligible, isActive, priority)` and returns `RESPONSE_OK`. Go version has no body.
- [ ] Missing `id` parameter (`String id` in Java). The Go method signature is `UpdateMember()` with no arguments.
- [ ] Missing `UpdateMemberDTO` parameter in method signature.

## QueryLeader
- [ ] Method is an empty stub — no logic at all. Java version gets the leader from `discoveryService.getLeader()`, checks for null (throws `NO_CONTENT` if null), extracts `nodeId`, looks up the corresponding `Member` from `getAllKnownMembers()`, and returns it. Go version has no body.

## ElectNewLeader
- [ ] Method is an empty stub — no logic at all. Java version accepts an optional `id` query parameter, branches: if `id == null` calls `discoveryService.electNewLeaderByPriority()`, otherwise calls `discoveryService.electNewLeaderByNodeId(id)`, and wraps the result in `HttpHandlerResult.okIfTruthy()`. Go version has no body.
- [ ] Missing `id` parameter (`@QueryParam(required = false) String id` in Java).

# SettingController.java
*Checked methods: queryClusterSettings(boolean queryLocalSettings, boolean onlyMutable), updateClusterSettings(boolean reset, boolean updateLocalSettings, @RequestBody(required = false), queryClusterConfigMetadata(boolean queryLocalSettings, boolean onlyMutable, boolean withValue)*

There is only one file and no other implementation exists. The Go code consists entirely of empty stub methods with no logic.

## queryClusterSettings

- [ ] Method body is completely empty — no implementation of the logic to select local vs. global properties based on `queryLocalSettings`, no call to `convertPropertiesToValueMap(properties, onlyMutable)`, and no return of a `SettingsDTO` with `SCHEMA_VERSION` and the value map.

## updateClusterSettings

- [ ] Method body is completely empty — no implementation of the `updateLocalSettings` branch that calls `propertiesManager.updateLocalProperties(reset, turmsProperties)` and returns a sync response, or the else-branch that calls `propertiesManager.updateGlobalProperties(reset, turmsProperties)` and returns an async response.

## queryClusterConfigMetadata

- [ ] Method body is completely empty — no implementation of the logic to select `ONLY_MUTABLE_METADATA` vs. `METADATA` based on `onlyMutable`, no conditional merge with property values when `withValue` is true (using local vs. global properties based on `queryLocalSettings`), and no return of a `SettingsDTO` with `SCHEMA_VERSION` and the settings map.

# BaseController.java
*Checked methods: getPageSize(@Nullable Integer size), queryBetweenDate(DateRange dateRange, DivideBy divideBy, Function3<DateRange, Boolean, Boolean, Mono<Long>> function, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages), queryBetweenDate(DateRange dateRange, DivideBy divideBy, Function<DateRange, Mono<Long>> function), checkAndQueryBetweenDate(DateRange dateRange, DivideBy divideBy, Function3<DateRange, Boolean, Boolean, Mono<Long>> function, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages), checkAndQueryBetweenDate(DateRange dateRange, DivideBy divideBy, Function<DateRange, Mono<Long>> function)*

## GetPageSize
- [ ] Method has no parameters — missing `size *int` parameter (Java: `@Nullable Integer size`)
- [ ] Method has no return value — missing `int` return type (Java returns `int`)
- [ ] Missing core logic: null/zero/negative check returning `defaultAvailableRecordsPerRequest`, and `Math.min(size, maxAvailableRecordsPerRequest)` cap
- [ ] Missing field: `defaultAvailableRecordsPerRequest` on the struct
- [ ] Missing field: `maxAvailableRecordsPerRequest` on the struct

## QueryBetweenDate
- [ ] Method has no parameters — missing all parameters: `dateRange DateRange`, `divideBy DivideBy`, `function func(DateRange, bool, bool)`, `areGroupMessages *bool`, `areSystemMessages *bool`
- [ ] Method has no return value — missing return of `[]StatisticsRecordDTO` and `error` (Java: `Mono<List<StatisticsRecordDTO>>`)
- [ ] Missing core logic: calling `DateTimeUtil.divideDuration(dateRange.start, dateRange.end, divideBy)` to get date pairs
- [ ] Missing core logic: iterating over date pairs and calling the function with `DateRange.of(pair.left, pair.right)`, `areGroupMessages`, `areSystemMessages`
- [ ] Missing core logic: mapping each result to `StatisticsRecordDTO{Date: pair.left, Total: total}`
- [ ] Missing core logic: calling `mergeStaticsRecords` to collect and sort results by date

## QueryBetweenDateFunc
- [ ] Method has no parameters — missing all parameters: `dateRange DateRange`, `divideBy DivideBy`, `function func(DateRange)`
- [ ] Method has no return value — missing return of `[]StatisticsRecordDTO` and `error` (Java: `Mono<List<StatisticsRecordDTO>>`)
- [ ] Missing core logic: calling `DateTimeUtil.divideDuration(dateRange.start, dateRange.end, divideBy)` to get date pairs
- [ ] Missing core logic: iterating over date pairs, calling the function with `DateRange.of(pair.left, pair.right)`
- [ ] Missing core logic: mapping each result to `StatisticsRecordDTO{Date: pair.left, Total: total}`
- [ ] Missing core logic: calling `mergeStaticsRecords` to collect and sort results by date

## CheckAndQueryBetweenDate
- [ ] Method has no parameters — missing all parameters: `dateRange DateRange`, `divideBy DivideBy`, `function func(DateRange, bool, bool)`, `areGroupMessages *bool`, `areSystemMessages *bool`
- [ ] Method has no return value — missing return of `[]StatisticsRecordDTO` and `error` (Java: `Mono<List<StatisticsRecordDTO>>`)
- [ ] Missing core logic: calling `isDurationNotGreaterThanMax` with `maxHourDifferencePerCountRequest`, `maxDayDifferencePerCountRequest`, `maxMonthDifferencePerCountRequest`
- [ ] Missing core logic: if duration check passes, delegating to `QueryBetweenDate`; otherwise returning error `ADMIN_REQUESTS_TOO_FREQUENT`
- [ ] Missing fields on struct: `maxHourDifferencePerCountRequest`, `maxDayDifferencePerCountRequest`, `maxMonthDifferencePerCountRequest`

## CheckAndQueryBetweenDateFunc
- [ ] Method has no parameters — missing all parameters: `dateRange DateRange`, `divideBy DivideBy`, `function func(DateRange)`
- [ ] Method has no return value — missing return of `[]StatisticsRecordDTO` and `error` (Java: `Mono<List<StatisticsRecordDTO>>`)
- [ ] Missing core logic: calling `isDurationNotGreaterThanMax` with `maxHourDifferencePerCountRequest`, `maxDayDifferencePerCountRequest`, `maxMonthDifferencePerCountRequest`
- [ ] Missing core logic: if duration check passes, delegating to `QueryBetweenDateFunc`; otherwise returning error `ADMIN_REQUESTS_TOO_FREQUENT`

# ServicePermission.java
*Checked methods: ServicePermission(ResponseStatusCode code, String reason), get(ResponseStatusCode code), get(ResponseStatusCode code, String reason)*

## NewServicePermission (constructor)

- [ ] **Missing `Code` field**: The Go struct `ServicePermission` has no fields at all. The Java record has two fields: `ResponseStatusCode code` and `String reason`. The Go struct should have equivalent fields (e.g., `Code` and `Reason`).
- [ ] **Constructor accepts no parameters**: `NewServicePermission()` takes no arguments, while the Java constructor `ServicePermission(ResponseStatusCode code, String reason)` takes two parameters and assigns them to the record fields. The Go constructor should accept corresponding parameters and assign them to the struct fields.
- [ ] **Missing static constant `OK`**: The Java code defines `public static final ServicePermission OK = new ServicePermission(ResponseStatusCode.OK, null)`. There is no equivalent in the Go code.

## Get (single-parameter overload)

- [ ] **Method completely missing**: The Java static method `get(ResponseStatusCode code)` returns `new ServicePermission(code, null)`. There is no corresponding Go function. The `@MappedFrom` annotation on `ShardedMap.Get` in `sharded_map.go` is incorrectly attributed — that is a generic map `Get` method, not the `ServicePermission.get()` factory method.

## Get (two-parameter overload)

- [ ] **Method completely missing**: The Java static method `get(ResponseStatusCode code, String reason)` returns `new ServicePermission(code, reason)`. There is no corresponding Go function.

# ExpirableEntityRepository.java
*Checked methods: isExpired(long creationDate), getEntityExpirationDate(), deleteExpiredData(String creationDateFieldName, Date expirationDate), findMany(Filter filter), findMany(Filter filter, QueryOptions options)*

Now I have all the information needed to produce the review.

## IsExpired

- [ ] **Method has no implementation.** The Go stub `func (r *ExpirableEntityRepository) IsExpired() {}` is completely empty — it takes no `creationDate` parameter, has no return value, and contains no logic. The Java version accepts `long creationDate`, calls `getEntityExpireAfterSeconds()`, and returns `expireAfterSeconds > 0 && creationDate < System.currentTimeMillis() - expireAfterSeconds * 1000L`.

## GetEntityExpirationDate

- [ ] **Method has no implementation.** The Go stubs in both `ExpirableEntityRepository.GetEntityExpirationDate()` and `ExpirableEntityService.GetEntityExpirationDate()` are empty — no parameters, no return value, no logic. The Java version calls `getEntityExpireAfterSeconds()`, returns `nil` if `expireAfterSeconds <= 0`, otherwise returns `new Date(System.currentTimeMillis() - expireAfterSeconds * 1000L)`.

## DeleteExpiredData

- [ ] **Hardcodes `"cd"` as the creation date field name instead of using a parameter.** The Java method accepts `String creationDateFieldName` as a parameter and passes it to the filter builder, making it flexible for any entity. The Go implementation hardcodes `"cd"`, which couples it to a specific BSON field name and will break for entities that use a different field name for their creation date.

## FindMany (Filter filter)

- [ ] **Drops the `QueryOptions` parameter.** The Java `findMany(Filter filter)` is a simple delegation to `mongoClient.findMany(entityClass, filter)` with no options. The Go `UserRepository.FindMany` also takes only a filter, which matches this signature. However, the Java class also defines a second overload `findMany(Filter filter, QueryOptions options)` that passes query options (projection, sort, limit, etc.). The Go `FindMany` merges both Java `@MappedFrom` annotations into a single method that only supports the no-options variant, losing the ability to pass `QueryOptions` (projection, pagination, sorting).

## FindMany (Filter filter, QueryOptions options)

- [ ] **Method is entirely missing.** The Java `ExpirableEntityRepository` defines two overloads: `findMany(Filter filter)` and `findMany(Filter filter, QueryOptions options)`. The Go code only has one `FindMany` that accepts a filter only. There is no second method accepting query options, so callers cannot pass projection, sort, limit, or skip options — all of which are supported in the Java version.

# ExpirableEntityService.java
*Checked methods: getEntityExpirationDate()*

Now I have a complete picture. Let me compare the Java and Go implementations of `getEntityExpirationDate()`.

**Java implementation (in `ExpirableEntityRepository`):**
```java
@Nullable
public Date getEntityExpirationDate() {
    int expireAfterSeconds = getEntityExpireAfterSeconds();
    if (expireAfterSeconds <= 0) {
        return null;
    }
    return new Date(System.currentTimeMillis() - expireAfterSeconds * 1000L);
}
```

And in `ExpirableEntityService`:
```java
@Nullable
public Date getEntityExpirationDate() {
    return expirableEntityRepository.getEntityExpirationDate();
}
```

**Go implementation (in service):**
```go
func (s *ExpirableEntityService) GetEntityExpirationDate() {
}
```

**Go implementation (in repository):**
```go
func (r *ExpirableEntityRepository) GetEntityExpirationDate() {
}
```

## GetEntityExpirationDate

- [ ] **Service method is a no-op stub**: `ExpirableEntityService.GetEntityExpirationDate()` has an empty body with no parameters and no return value. The Java version accepts no parameters and returns a nullable `Date` (the expiration date). The Go version should return `*time.Time` (or equivalent) and delegate to the repository.

- [ ] **Repository method is a no-op stub**: `ExpirableEntityRepository.GetEntityExpirationDate()` has an empty body. The Java version computes `System.currentTimeMillis() - expireAfterSeconds * 1000L` and returns it as a `Date`, or returns `null` if `expireAfterSeconds <= 0`. The Go version does nothing and returns nothing.

- [ ] **Missing `getEntityExpireAfterSeconds()` abstract method**: The repository does not define or call the abstract `getEntityExpireAfterSeconds()` method that `getEntityExpirationDate()` depends on to determine the expiration time window.

- [ ] **Missing expiration calculation logic**: The core logic of computing `currentTime - expireAfterSeconds * 1000` and returning it (or returning nil when `expireAfterSeconds <= 0`) is entirely absent.

- [ ] **Service does not hold a reference to the repository**: The Go `ExpirableEntityService` struct has no fields. The Java version holds a `private final ExpirableEntityRepository<T, ?> expirableEntityRepository` field and delegates `getEntityExpirationDate()` to it.

# UserDefinedAttributesService.java
*Checked methods: updateGlobalProperties(UserDefinedAttributesProperties properties), parseAttributesForUpsert(Map<String, Value> userDefinedAttributes)*

## UpdateGlobalProperties
- [ ] Method has an empty body with no parameters — the Java version takes a `UserDefinedAttributesProperties properties` parameter, iterates over `properties.getAllowedAttributes()`, builds a map of source names to attribute properties (handling duplicate detection and defaulting `storedName` to `sourceName` when empty), collects immutable attributes, and assigns to `sourceNameToAttributeProperties`, `knownAttributes`, `immutableAttributes`, and `ignoreUnknownAttributesOnUpsert` fields. None of this logic is present.
- [ ] The `UserDefinedAttributesService` struct has no fields — the Java class declares `knownAttributes`, `sourceNameToAttributeProperties`, `immutableAttributes`, and `ignoreUnknownAttributesOnUpsert` as instance fields. None are present in the Go struct.

## ParseAttributesForUpsert
- [ ] Method has an empty body with no parameters — the Java version takes `Map<String, Value> userDefinedAttributes`, validates non-null, handles the empty case, checks for immutable attributes in the input keys, calls `findUserDefinedAttributes` to check existing values for immutability conflicts, and calls `parseAttributes`. None of this logic is present.
- [ ] Missing the `parseAttributes` helper method entirely — the Java class has a `parseAttributes(boolean ignoreUnknownAttributes, Map<String, Value> inputAttributes)` method that handles unknown attribute detection (with different error behaviors for known/unknown) and value parsing. This method is not implemented in Go.
- [ ] Missing the abstract `findUserDefinedAttributes` method — the Java class declares `protected abstract Mono<List<String>> findUserDefinedAttributes(Collection<String> immutableAttributesForUpsert)` which is called to check existing immutable attribute values. No equivalent exists in Go.

# ExpirableRequestInspector.java
*Checked methods: isProcessedByResponder(@Nullable RequestStatus status)*

## IsProcessedByResponder

- [ ] **Missing parameter**: The Go method `IsProcessedByResponder()` has no parameters, but the Java version takes a `RequestStatus status` parameter. It should accept a status argument.
- [ ] **Missing return value**: The Go method returns nothing (`void`-like), but the Java version returns `bool`. It should return `bool`.
- [ ] **Missing core logic**: The method body is entirely empty. The Java version checks if `status == ACCEPTED || status == DECLINED || status == IGNORED`. None of this comparison logic is present in the Go implementation.

# DataValidator.java
*Checked methods: validRequestStatus(RequestStatus status), validResponseAction(ResponseAction action), validDeviceType(DeviceType deviceType), validProfileAccess(ProfileAccessStrategy value), validRelationshipKey(UserRelationship.Key key), validRelationshipGroupKey(UserRelationshipGroup.Key key), validGroupMemberKey(GroupMember.Key key), validGroupMemberRole(GroupMemberRole role), validGroupBlockedUserKey(GroupBlockedUser.Key key), validNewGroupQuestion(NewGroupQuestion question), validGroupQuestionIdAndAnswer(Map.Entry<Long, String> questionIdAndAnswer)*

## ValidRequestStatus
- [ ] Method is a stub returning `nil` unconditionally. The Java version checks if `status == RequestStatus.UNRECOGNIZED` and throws `ILLEGAL_ARGUMENT` if so. The Go version accepts an `interface{}` parameter with a `name` string but performs no validation logic at all.

## ValidResponseAction
- [ ] Method is a stub with no parameters and no body. The Java version takes a `ResponseAction action` parameter, checks if it equals `ResponseAction.UNRECOGNIZED`, and throws `ILLEGAL_ARGUMENT` if so.

## ValidDeviceType
- [ ] Method is a stub with no parameters and no body. The Java version takes a `DeviceType deviceType` parameter, checks if it equals `DeviceType.UNRECOGNIZED`, and throws `ILLEGAL_ARGUMENT` if so.

## ValidProfileAccess
- [ ] Method is a stub with no parameters and no body. The Java version takes a `ProfileAccessStrategy value` parameter, checks if it equals `ProfileAccessStrategy.UNRECOGNIZED`, and throws `ILLEGAL_ARGUMENT` if so.

## ValidRelationshipKey
- [ ] Method is a stub with no parameters and no body. The Java version takes a `UserRelationship.Key key` parameter, checks if `key == null || key.getOwnerId() == null || key.getRelatedUserId() == null`, and throws `ILLEGAL_ARGUMENT` if so.

## ValidRelationshipGroupKey
- [ ] Method is a stub with no parameters and no body. The Java version takes a `UserRelationshipGroup.Key key` parameter, checks if `key == null || key.getOwnerId() == null || key.getGroupIndex() == null`, and throws `ILLEGAL_ARGUMENT` if so.

## ValidGroupMemberKey
- [ ] Method is a stub with no parameters and no body. The Java version takes a `GroupMember.Key key` parameter, checks if `key == null || key.getGroupId() == null || key.getUserId() == null`, and throws `ILLEGAL_ARGUMENT` if so.

## ValidGroupMemberRole
- [ ] Method is a stub with no parameters and no body. The Java version takes a `GroupMemberRole role` parameter, checks if it equals `GroupMemberRole.UNRECOGNIZED`, and throws `ILLEGAL_ARGUMENT` if so.

## ValidGroupBlockedUserKey
- [ ] Method is a stub with no parameters and no body. The Java version takes a `GroupBlockedUser.Key key` parameter, checks if `key == null || key.getGroupId() == null || key.getUserId() == null`, and throws `ILLEGAL_ARGUMENT` if so.

## ValidNewGroupQuestion
- [ ] Method is a stub with no parameters and no body. The Java version takes a `NewGroupQuestion question` parameter and performs two validations: (1) checks that `question.answers()` is not empty, and (2) checks that `question.score()` is not null and >= 0, throwing `ILLEGAL_ARGUMENT` on failure.

## ValidGroupQuestionIdAndAnswer
- [ ] Method is a stub with no parameters and no body. The Java version takes a `Map.Entry<Long, String> questionIdAndAnswer` parameter, checks if it is null, or if its key or value is null, and throws `ILLEGAL_ARGUMENT` if so.

# CancelMeetingResult.java
*Checked methods: CancelMeetingResult(boolean success, @Nullable Meeting meeting)*

## CancelMeetingResult

- [ ] Missing the `FAILED` static constant. The Java code defines `public static final CancelMeetingResult FAILED = new CancelMeetingResult(false, null);` but the Go version has no equivalent (e.g., a package-level `var CancelMeetingResultFAILED = CancelMeetingResult{Success: false, Meeting: nil}`).
- [ ] `Meeting` field is typed as `interface{}` instead of the proper `*Meeting` type (or its Go equivalent). The Java code uses `@Nullable Meeting meeting` referring to the `im.turms.service.domain.conference.po.Meeting` entity. Using `interface{}` loses type safety and deviates from the original design.

# UpdateMeetingResult.java
*Checked methods: UpdateMeetingResult(boolean success, @Nullable Meeting meeting)*

## UpdateMeetingResult

- [ ] **Missing `FAILED` static constant**: The Java version defines a `public static final UpdateMeetingResult FAILED = new UpdateMeetingResult(false, null);` constant. The Go version has no equivalent `FAILED` singleton/constant (e.g., `var FailedUpdateMeetingResult = &UpdateMeetingResult{Success: false, Meeting: nil}`).
- [ ] **`Meeting` field typed as `interface{}` instead of the concrete `Meeting` type**: The Java version uses the concrete `Meeting` PO type for the `meeting` field. The Go version uses `interface{}` with a comment "Replace with actual Meeting type", meaning the field lacks proper typing and type safety compared to the original.

# ConferenceServiceController.java
*Checked methods: handleCreateMeetingRequest(), handleDeleteMeetingRequest(), handleUpdateMeetingRequest(), handleQueryMeetingsRequest(), handleUpdateMeetingInvitationRequest()*

The Go file contains only empty stub methods with no implementation whatsoever.

## HandleCreateMeetingRequest
- [ ] Method body is completely empty. Missing all logic: extracting `CreateMeetingRequest` fields (userId, groupId, name, intro, password, startDate) with null-checks via `hasX()`, calling `conferenceService.authAndCreateMeeting()`, and mapping the result to `RequestHandlerResult.ofDataLong(meeting.getId())`.

## HandleDeleteMeetingRequest
- [ ] Method body is completely empty. Missing all logic: extracting `DeleteMeetingRequest`, calling `conferenceService.authAndCancelMeeting()`, and the full notification branching logic — (1) `notifyMeetingParticipantsOfMeetingCanceled`: querying participants, removing requester, returning `RequestHandlerResult.of(true, participantIds, turmsRequest)`; (2) `notifyRequesterOtherOnlineSessionsOfMeetingCanceled`: returning `RequestHandlerResult.of(true, turmsRequest)` on success; (3) default: returning `RequestHandlerResult.OK`.

## HandleUpdateMeetingRequest
- [ ] Method body is completely empty. Missing all logic: extracting `UpdateMeetingRequest` fields (id, name, intro, password) with `hasX()` null-checks, calling `conferenceService.authAndUpdateMeeting()`, and the full notification branching logic — (1) `notifyMeetingParticipantsOfMeetingUpdated`: checking `result.success() && (request.hasName() || request.hasIntro())`, querying participants, removing requester, and conditionally building a notification with `clearPassword()` when password is present; (2) `notifyRequesterOtherOnlineSessionsOfMeetingUpdated`: similar success+fields check; (3) default: `RequestHandlerResult.OK`.

## HandleQueryMeetingsRequest
- [ ] Method body is completely empty. Missing all logic: extracting all `QueryMeetingsRequest` fields (ids, creatorIds, userIds, groupIds as sets from count>0 checks; creationDateStart/End as dates from hasX(); skip/limit as nullable ints), calling `conferenceService.authAndQueryMeetings()` with all parameters, converting each meeting to proto via `meeting2proto()`, collecting into a list, wrapping in `TurmsNotification.Data` with `MeetingsBuilder`, and returning as `RequestHandlerResult.of()`.

## HandleUpdateMeetingInvitationRequest
- [ ] Method body is completely empty. Missing all logic: extracting `UpdateMeetingInvitationRequest` fields (meetingId, password with `hasPassword()` check, responseAction), calling `conferenceService.authAndUpdateMeetingInvitation()`, and the full notification branching logic — (1) `notifyMeetingParticipantsOfMeetingInvitationUpdated`: checking `result.updated()`, querying participants, removing requester, conditionally returning response data with `accessToken` when `responseAction == ACCEPT`, and building a notification with `clearPassword()` when password is present and participants are non-empty; (2) `notifyRequesterOtherOnlineSessionsOfMeetingInvitationUpdated`: returning response with `accessToken` on ACCEPT, empty set for participants; (3) default: returning response data with `accessToken` on ACCEPT only.

# MeetingRepository.java
*Checked methods: updateEndDate(Long meetingId, Date endDate), updateCancelDateIfNotCanceled(Long meetingId, Date cancelDate), updateMeeting(Long meetingId, @Nullable String name, @Nullable String intro, @Nullable String password), find(@Nullable Collection<Long> ids, @Nullable Collection<Long> creatorIds, @Nullable Collection<Long> userIds, @Nullable Collection<Long> groupIds, @Nullable Date creationDateStart, @Nullable Date creationDateEnd, @Nullable Integer skip, @Nullable Integer limit), find(@Nullable Collection<Long> ids, @NotNull Long creatorId, @NotNull Long userId, @Nullable Date creationDateStart, @Nullable Date creationDateEnd, @Nullable Integer skip, @Nullable Integer limit)*

All five Go methods are empty stubs with no parameters, no return values, and no logic. Every single method is missing its entire implementation.

## UpdateEndDate
- [ ] Method has no parameters (missing `meetingId int64, endDate time.Time` or equivalent)
- [ ] Method has no return value (should return an update result or error)
- [ ] Missing filter: `.eq(ID, meetingId)`
- [ ] Missing update: `.set(END_DATE, endDate)`
- [ ] Missing MongoDB updateOne call

## UpdateCancelDateIfNotCanceled
- [ ] Method has no parameters (missing `meetingId int64, cancelDate time.Time` or equivalent)
- [ ] Method has no return value (should return an update result or error)
- [ ] Missing filter: `.eq(ID, meetingId).eq(CANCEL_DATE, null)` — the critical "if not canceled" null check on `CANCEL_DATE` is missing
- [ ] Missing update: `.set(CANCEL_DATE, cancelDate)`
- [ ] Missing MongoDB updateOne call

## UpdateMeeting
- [ ] Method has no parameters (missing `meetingId int64, name *string, intro *string, password *string` or equivalent)
- [ ] Method has no return value (should return an update result or error)
- [ ] Missing filter: `.eq(ID, meetingId)`
- [ ] Missing update: `.setIfNotNull(NAME, name).setIfNotNull(INTRO, intro).setIfNotNull(PASSWORD, password)`
- [ ] Missing MongoDB updateOne call

## Find (multi-criteria overload)
- [ ] Method has no parameters (missing `ids, creatorIds, userIds, groupIds []int64, creationDateStart, creationDateEnd *time.Time, skip, limit *int` or equivalent)
- [ ] Method has no return value (should return a slice of Meeting entities or error)
- [ ] Missing filter: `.inIfNotNull(ID, ids).inIfNotNull(CREATOR_ID, creatorIds).inIfNotNull(USER_ID, userIds).inIfNotNull(GROUP_ID, groupIds).addBetweenIfNotNull(CREATION_DATE, creationDateStart, creationDateEnd)`
- [ ] Missing query options: `.skipIfNotNull(skip).limitIfNotNull(limit)`
- [ ] Missing MongoDB findMany call

## FindByCreatorAndUser
- [ ] Method has no parameters (missing `ids []int64, creatorId int64, userId int64, creationDateStart, creationDateEnd *time.Time, skip, limit *int` or equivalent)
- [ ] Method has no return value (should return a slice of Meeting entities or error)
- [ ] Missing filter with OR logic: `.inIfNotNull(ID, ids).or(eq(CREATOR_ID, creatorId).eq(USER_ID, userId)).addBetweenIfNotNull(CREATION_DATE, creationDateStart, creationDateEnd)` — the OR sub-filter combining creatorId and userId is missing
- [ ] Missing query options: `.skipIfNotNull(skip).limitIfNotNull(limit)`
- [ ] Missing MongoDB findMany call

# ConferenceService.java
*Checked methods: onExtensionStarted(ConferenceServiceProvider extension), authAndCancelMeeting(@NotNull Long requesterId, @NotNull Long meetingId), queryMeetingParticipants(@Nullable Long userId, @Nullable Long groupId), authAndUpdateMeeting(@NotNull Long requesterId, @NotNull Long meetingId, @Nullable String name, @Nullable String intro, @Nullable String password), authAndUpdateMeetingInvitation(@NotNull Long requesterId, @NotNull Long meetingId, @Nullable String password, @NotNull ResponseAction responseAction), authAndQueryMeetings(@NotNull Long requesterId, @Nullable Set<Long> ids, @Nullable Set<Long> creatorIds, @Nullable Set<Long> userIds, @Nullable Set<Long> groupIds, @Nullable Date creationDateStart, @Nullable Date creationDateEnd, @Nullable Integer skip, @Nullable Integer limit)*

## OnExtensionStarted

- [ ] The Go method `OnExtensionStarted` is completely empty with no parameters, while the Java version receives a `ConferenceServiceProvider extension` parameter and calls `extension.addMeetingEndedEventListener(ConferenceService.this::handleMeetingEndedEvent)` to register a meeting-ended event listener with error logging.

## AuthAndCancelMeeting

- [ ] The Go method `AuthAndCancelMeeting` is completely empty with no parameters, while the Java version takes `requesterId` and `meetingId` parameters and implements full logic: validates non-null parameters, checks `allowCancel` flag, checks for a conference service provider, looks up the meeting by ID, verifies the requester is the creator (otherwise checks if they're allowed to view meeting info for appropriate error responses), updates the cancel date in a transaction, invokes the plugin extension point to cancel the meeting, and returns a `CancelMeetingResult`.

## QueryMeetingParticipants

- [ ] The Go method `QueryMeetingParticipants` is completely empty with no parameters, while the Java version takes `userId` and `groupId` parameters and implements logic: if `userId` is non-null returns it as a singleton set, else if `groupId` is non-null queries group member IDs via `groupMemberService.queryGroupMemberIds`, otherwise returns an empty set.

## AuthAndUpdateMeeting

- [ ] The Go method `AuthAndUpdateMeeting` is completely empty with no parameters, while the Java version takes `requesterId`, `meetingId`, `name`, `intro`, and `password` parameters and implements full logic: validates parameters (not-null for IDs, length ranges for name/intro, validates password if non-null), returns `FAILED` if all update fields are null, looks up the meeting, checks if password update is requested (requiring creator privileges), checks `isAllowedToViewMeetingInfo` for non-password updates, and calls `meetingRepository.updateMeeting` returning an `UpdateMeetingResult`.

## AuthAndUpdateMeetingInvitation

- [ ] The Go method `AuthAndUpdateMeetingInvitation` is completely empty with no parameters, while the Java version takes `requesterId`, `meetingId`, `password`, and `responseAction` parameters and implements full logic: validates non-null parameters, checks for conference service provider, returns early for `IGNORE`/`UNRECOGNIZED` actions, looks up the meeting, verifies requester authorization (checks if userId matches requesterId or if requester is a group member), validates password matching, handles DECLINE action, checks meeting status (expired/canceled/pending/ended) for ACCEPT action, and invokes the plugin extension point for accepting the invitation.

## AuthAndQueryMeetings

- [ ] The Go method `AuthAndQueryMeetings` is completely empty with no parameters, while the Java version takes `requesterId`, `ids`, `creatorIds`, `userIds`, `groupIds`, `creationDateStart`, `creationDateEnd`, `skip`, and `limit` parameters and implements complex authorization-based query logic: validates `requesterId` is non-null, then branches based on userId/groupId/creatorId counts — when `userIds` is present and contains only the requester it allows the query, when `userIds` and `groupIds` are both non-empty returns empty, when no `userIds`/`groupIds` with creatorIds it scopes queries to the requester, when `groupIds` are present it queries the user's joined groups and intersects with requested groupIds, and delegates to `meetingRepository.find` with appropriate filters.

# ConversationController.java
*Checked methods: queryConversations(@QueryParam(required = false), deleteConversations(@QueryParam(required = false), updateConversations(@QueryParam(required = false)*

## QueryConversations
- [ ] Method body is completely empty (stub only). Missing all core logic: parameters (privateConversationKeys, ownerIds, groupIds), querying private conversations via conversationService.queryPrivateConversations and queryPrivateConversationsByOwnerIds, querying group conversations via conversationService.queryGroupConversations, combining results into ConversationsDTO, and returning HTTP response.
- [ ] Missing struct fields: `conversationService` dependency and `propertiesManager` (via BaseController) are not present on `ConversationController`.
- [ ] Missing parameter parsing for `privateConversationKeys` (List<PrivateConversation.Key>), `ownerIds` (Set<Long>), and `groupIds` (Set<Long>).

## DeleteConversations
- [ ] Method body is completely empty (stub only). Missing all core logic: deleting private conversations by keys via conversationService.deletePrivateConversations, deleting private conversations by ownerIds via conversationService.deletePrivateConversations(ownerIds, null), deleting group conversations via conversationService.deleteGroupConversations, merging DeleteResults via OperationResultConvertor.merge, and returning HttpHandlerResult.deleteResult.
- [ ] Missing parameter parsing for `privateConversationKeys`, `ownerIds`, and `groupIds`.

## UpdateConversations
- [ ] Method body is completely empty (stub only). Missing all core logic: upserting private conversations read date via conversationService.upsertPrivateConversationsReadDate, upserting group conversations read date via conversationService.upsertGroupConversationsReadDate, parallel execution with `whenDelayError` equivalent, and returning RESPONSE_OK.
- [ ] Missing parameter parsing for `privateConversationKeys`, `groupConversationMemberKeys`, and request body `UpdateConversationDTO` with `readDate` field.

# ConversationsDTO.java
*Checked methods: ConversationsDTO(List<PrivateConversation> privateConversations, List<GroupConversation> groupConversations)*

## ConversationsDTO

- [ ] **Field types are `[]interface{}` instead of properly typed slices.** The Java code uses `List<PrivateConversation>` and `List<GroupConversation>` with full type safety. The Go code uses `[]interface{}` for both `PrivateConversations` and `GroupConversations`, losing all type information. These should be `[]po.PrivateConversation` and `[]po.GroupConversation` respectively, matching the Go domain model types that exist in the codebase.

- [ ] **Missing import for the `po` package.** Since the fields should reference `po.PrivateConversation` and `po.GroupConversation`, the file needs to import the `po` package (currently it only imports `"time"`).

# ConversationServiceController.java
*Checked methods: handleQueryConversationsRequest(), handleUpdateTypingStatusRequest(), handleUpdateConversationRequest()*

## HandleQueryConversationsRequest

- [ ] **All core logic is missing**: The method body is completely empty (`func (c *ConversationServiceController) HandleQueryConversationsRequest() {}`). The Java version implements: parsing `QueryConversationsRequest` to get `targetIds` (user IDs) and `groupIds`, returning `NO_CONTENT` when both are empty, calling `conversationService.queryGroupConversations(groupIds)` for group conversations or `conversationService.queryPrivateConversations(targetIds, userId)` for private conversations, converting results to proto `Conversations` with either `groupConversations` or `privateConversations`, and wrapping in `RequestHandlerResult`.

## HandleUpdateTypingStatusRequest

- [ ] **All core logic is missing**: The method body is completely empty. The Java version implements: parsing `UpdateTypingStatusRequest` to get `isGroupMessage` and `toId`, calling `conversationService.authAndUpdateTypingStatus(userId, isGroupMessage, toId)`, and mapping the returned `recipientIds` into `RequestHandlerResult.of(recipientIds, turmsRequest)`.

## HandleUpdateConversationRequest

- [ ] **All core logic is missing**: The method body is completely empty. The Java version implements: validating that at least one of `userId` or `groupId` is present (returning `ILLEGAL_ARGUMENT` error otherwise), extracting `readDate`, branching on `hasUserId` to call either `conversationService.authAndUpsertPrivateConversationReadDate(requesterId, targetId, readDate)` or `conversationService.authAndUpsertGroupConversationReadDate(targetId, requesterId, readDate)`, and after completion, conditionally returning `RequestHandlerResult` with notification recipients based on four configurable boolean flags (`notifyContactOfPrivateConversationReadDateUpdated`, `notifyRequesterOtherOnlineSessionsOfPrivateConversationReadDateUpdated`, `notifyOtherGroupMembersOfGroupConversationReadDateUpdated`, `notifyRequesterOtherOnlineSessionsOfGroupConversationReadDateUpdated`).
- [ ] **Configuration properties are missing**: The Go `ConversationServiceController` struct has no fields for the notification properties or dependencies (`conversationService`, `groupMemberService`) that the Java version holds and uses in `handleUpdateConversationRequest()`.

# ConversationSettingsServiceController.java
*Checked methods: handleUpdateConversationSettingsRequest(), handleDeleteConversationSettingsRequest(), handleQueryConversationSettingsRequest()*

## HandleUpdateConversationSettingsRequest
- [ ] Method body is completely empty — missing all logic including: request parsing, validation that userId and groupId must not both be null, branching on hasUserId to call either `upsertPrivateConversationSettings` or `upsertGroupConversationSettings`, passing the settings map, mapping the result to `RequestHandlerResult` with notification flags for private/group conversation setting updated, and returning `RequestHandlerResult.OK` when not updated.

## HandleDeleteConversationSettingsRequest
- [ ] Method body is completely empty — missing all logic including: request parsing, extracting userIds/groupIds/names from the delete request, calling `conversationSettingsService.unsetSettings` with the extracted sets (or null), mapping the result to `RequestHandlerResult` with notification logic that checks `hasUserId && notifyRequesterOtherOnlineSessionsOfPrivateConversationSettingDeleted || hasGroupId && notifyRequesterOtherOnlineSessionsOfGroupConversationSettingDeleted`, and passing the turmsRequest.

## HandleQueryConversationSettingsRequest
- [ ] Method body is completely empty — missing all logic including: request parsing, extracting userIds/groupIds/names/lastUpdatedDateStart from the query request, calling `conversationSettingsService.querySettings` with the extracted parameters (sets or null, date or null), collecting results into a list, building a `ConversationSettingsList` proto response by iterating settings and converting via `ProtoModelConvertor.conversationSettings2proto`, and returning a `RequestHandlerResult` with the built `TurmsNotificationData`.

# ConversationSettingsRepository.java
*Checked methods: upsertSettings(Long ownerId, Long targetId, Map<String, Object> settings), unsetSettings(Long ownerId, @Nullable Collection<Long> targetIds, @Nullable Collection<String> settingNames), findByIdAndSettingNames(Long ownerId, @Nullable Collection<String> settingNames, @Nullable Date lastUpdatedDateStart), findByIdAndSettingNames(Collection<ConversationSettings.Key> keys, @Nullable Collection<String> settingNames, @Nullable Date lastUpdatedDateStart), findSettingFields(Long ownerId, Long targetId, Collection<String> includedFields), deleteByOwnerIds(Collection<Long> ownerIds, @Nullable ClientSession clientSession)*

Now I have a complete picture. Let me analyze each method carefully.

## upsertSettings(Long ownerId, Long targetId, Map<String, Object> settings)

**Java logic:**
1. Filter by composite `_id` = `ConversationSettings.Key(ownerId, targetId)` 
2. Set `lastUpdatedDate` to current date
3. Set each setting key as `"s." + key`
4. Upsert

**Go code (user_settings_service.go):**
The `UpsertSettings` in the service delegates to `userSettingsRepository.UpsertSettings`, which is the **User** settings repo, not the Conversation settings repo. The ConversationSettingsRepository has no `UpsertSettings` method at all — it's a stub with empty methods. But the user asked about the mapped method. Looking at the `userSettingsRepository.UpsertSettings`:
- It does NOT set `lastUpdatedDate` (the Java version always sets `lud` to `new Date()`)
- The projection key uses `"s."+k` correctly
- Uses upsert correctly

## unsetSettings(Long ownerId, Collection<Long> targetIds, Collection<String> settingNames)

**Java logic:**
1. If targetIds null/empty: filter by `ID_OWNER_ID` (ownerId)
2. Else: build composite keys from ownerId+targetId, filter by `_id` IN keys
3. Always set `lastUpdatedDate` to current date
4. If settingNames null/empty: unset entire `"s"` field
5. Else: unset each `"s." + settingName`

**Go code (user_settings_service.go):**
The `UnsetSettings` takes `(ctx, userID int64, keys []string)` — only userID and setting names, no targetIds parameter. It delegates to `userSettingsRepository.UnsetSettings` which:
- Does NOT handle the targetIds parameter at all (the Java version has complex targetIds logic)
- Does NOT set `lastUpdatedDate` (Java always sets `lud`)
- Does NOT handle the case where settingNames is null/empty to unset the entire `"s"` field (Go returns nil if keys empty, while Java unsets entire settings)
- The signature is completely different — missing targetIds

## findByIdAndSettingNames(Long ownerId, Collection<String> settingNames, Date lastUpdatedDateStart)

**Java logic:**
1. Filter by `ID_OWNER_ID` = ownerId
2. Add `gte` filter on `LAST_UPDATED_DATE` if lastUpdatedDateStart not null
3. If settingNames not null/empty: project `lud` + each settingName field
4. Return multiple results

**Go code:** The ConversationSettingsRepository.`FindByIdAndSettingNames` is an empty stub with no parameters, no logic at all.

The `userSettingsRepository.FindByIdAndSettingNames` (for user settings, not conversation settings):
- Does NOT filter by `lastUpdatedDateStart` (Java applies `gte` on `lud`)
- Projection: when names is non-empty, projects `"s."+name` but does NOT project `lud` (Java always projects `lud` when settingNames is non-empty)
- Returns single result (`FindOne`) vs Java returns `Flux` (multiple results)

## findByIdAndSettingNames(Collection<ConversationSettings.Key> keys, Collection<String> settingNames, Date lastUpdatedDateStart)

**Java logic:**
1. Filter by `_id` IN keys
2. Add `gte` filter on `LAST_UPDATED_DATE` if lastUpdatedDateStart not null
3. If settingNames not null/empty: project `lud` + each settingName field
4. Return multiple results

**Go code:** `FindByIdAndSettingNamesWithKeys` is an empty stub with no parameters, no logic.

## findSettingFields(Long ownerId, Long targetId, Collection<String> includedFields)

**Java logic:**
1. Filter by `_id` = `Key(ownerId, targetId)`
2. Call `findObjectFields` with settings field `"s"` and includedFields

**Go code:** `FindSettingFields` is an empty stub with no parameters, no logic.

## deleteByOwnerIds(Collection<Long> ownerIds, ClientSession clientSession)

**Java logic:**
1. Filter by `ID_OWNER_ID` IN ownerIds
2. Delete many with optional clientSession

**Go code:** `DeleteByOwnerIds` is an empty stub with no parameters, no logic.

---

Here is the complete bug report:

## UpsertSettings

- [ ] **Missing `lastUpdatedDate` (`lud`) field assignment**: The Java version always sets `ConversationSettings.Fields.LAST_UPDATED_DATE` (`lud`) to `new Date()` in every upsert. The Go `userSettingsRepository.UpsertSettings` does not include `lud` in the `$set` update document.
- [ ] **Mapped to wrong domain**: The `@MappedFrom upsertSettings(Long ownerId, Long targetId, Map<String, Object> settings)` annotation on the Go `userSettingsService.UpsertSettings` maps it to the **user** settings domain instead of the **conversation** settings domain. The Java method takes both `ownerId` and `targetId` (composite key), but the Go version only takes `userID`, losing the `targetId` parameter entirely.
- [ ] **ConversationSettingsRepository is missing the `UpsertSettings` method entirely**: The Java ConversationSettingsRepository has `upsertSettings(Long ownerId, Long targetId, Map<String, Object> settings)` but the Go ConversationSettingsRepository has no corresponding method.

## UnsetSettings

- [ ] **Missing `targetIds` parameter**: The Java version has `unsetSettings(Long ownerId, @Nullable Collection<Long> targetIds, @Nullable Collection<String> settingNames)`. The Go version takes only `(ctx, userID int64, keys []string)`, completely dropping the `targetIds` parameter and its associated filter logic (building composite keys from ownerId+targetId, filtering by `_id` IN keys vs by `_id.oid`).
- [ ] **Missing `lastUpdatedDate` (`lud`) field assignment**: The Java version always sets `ConversationSettings.Fields.LAST_UPDATED_DATE` (`lud`) to `new Date()` during unset operations. The Go version does not set `lud`.
- [ ] **Incorrect handling of null/empty `settingNames`**: When `settingNames` is null or empty, the Java version unsets the entire `"s"` (settings) field via `update.unset(ConversationSettings.Fields.SETTINGS)`. The Go version returns `nil` immediately when `keys` is empty, performing no database operation at all.
- [ ] **ConversationSettingsRepository is missing the `UnsetSettings` method entirely**: The Java ConversationSettingsRepository has this method but the Go ConversationSettingsRepository has no corresponding method.

## FindByIdAndSettingNames (by ownerId)

- [ ] **Empty stub with no implementation**: The Go `ConversationSettingsRepository.FindByIdAndSettingNames()` has no parameters, no body, and no return values. The Java version filters by `_id.oid = ownerId`, optionally applies `gte` on `lud`, and optionally projects specific fields.
- [ ] **Missing `lastUpdatedDateStart` filter**: The Java version applies `.gteIfNotNull(LAST_UPDATED_DATE, lastUpdatedDateStart)` to filter records. The Go version (in the user settings repo it was ported to) does not apply this filter.
- [ ] **Missing `lud` in projection**: When `settingNames` is non-empty, the Java version always includes `LAST_UPDATED_DATE` (`lud`) in the projection. The Go version (user settings repo) does not include `lud` in the projection.
- [ ] **Returns single result instead of multiple**: The Java version returns `Flux<ConversationSettings>` (multiple results via `findMany`). The Go version in the user settings repo uses `FindOne`, returning only a single document.

## FindByIdAndSettingNamesWithKeys (by collection of keys)

- [ ] **Empty stub with no implementation**: The Go `ConversationSettingsRepository.FindByIdAndSettingNamesWithKeys()` has no parameters, no body, and no return values. The Java version filters by `_id IN keys`, optionally applies `gte` on `lud`, and optionally projects specific fields.
- [ ] **Missing `lastUpdatedDateStart` filter**: Same as the ownerId overload — the `gte` filter on `lud` is not implemented.
- [ ] **Missing `lud` in projection**: Same as the ownerId overload — `lud` is not included in projection when settingNames is provided.

## FindSettingFields

- [ ] **Empty stub with no implementation**: The Go `ConversationSettingsRepository.FindSettingFields()` has no parameters, no body, and no return values. The Java version filters by `_id = Key(ownerId, targetId)` and calls `findObjectFields` to retrieve specific setting fields.

## DeleteByOwnerIds

- [ ] **Empty stub with no implementation**: The Go `ConversationSettingsRepository.DeleteByOwnerIds()` has no parameters, no body, and no return values. The Java version filters by `_id.oid IN ownerIds` and deletes matching documents, with optional `clientSession` support.

# GroupConversationRepository.java
*Checked methods: upsert(Long groupId, Long memberId, Date readDate, boolean allowMoveReadDateForward), upsert(Long groupId, Collection<Long> memberIds, Date readDate), deleteMemberConversations(Collection<Long> groupIds, Long memberId, ClientSession session)*

## `upsert(Long groupId, Long memberId, Date readDate, boolean allowMoveReadDateForward)` (Go: `UpsertReadDate`)

- [ ] **Missing `allowMoveReadDateForward` parameter and conditional filter logic.** The Java version accepts a `boolean allowMoveReadDateForward` parameter. When it is `false`, the filter includes `.ltOrNull(fieldKey, readDate)` to ensure the read date is only moved forward (i.e., the update only applies if the existing value is `null` or less than `readDate`). When it is `true`, no such condition is added. The Go version (`UpsertReadDate`) has no `allowMoveReadDateForward` parameter at all — it always unconditionally sets the read date, which is equivalent to only the `allowMoveReadDateForward = true` behavior. The `false` branch is completely missing.

## `upsert(Long groupId, Collection<Long> memberIds, Date readDate)` (Go: not implemented)

- [ ] **Method is completely missing.** The Java version has a bulk upsert that sets the read date for multiple members in a single update operation. There is no equivalent Go method. The `@MappedFrom` comment in `group_version_repository.go` references this method signature but the `Upsert` method there is for `GroupVersion`, not `GroupConversation` — it updates version fields (`mbr`, `bl`, `jr`, `jq`, `invt`, `info`), which is unrelated.

## `deleteMemberConversations(Collection<Long> groupIds, Long memberId, ClientSession session)` (Go: `DeleteMemberConversations`)

- [ ] **Method is a stub with no implementation.** The Go version is declared as `func (r *GroupConversationRepository) DeleteMemberConversations()` with an empty body. It is missing all parameters (`groupIDs []int64`, `memberID int64`, and a session/transaction context), all filter logic (`_id $in groupIds`), all update logic (`$unset` on the member's read date field), and all MongoDB execution logic.

# PrivateConversationRepository.java
*Checked methods: upsert(Set<PrivateConversation.Key> keys, Date readDate, boolean allowMoveReadDateForward), deleteConversationsByOwnerIds(Set<Long> ownerIds, @Nullable ClientSession session), findConversations(Collection<Long> ownerIds)*

Now I have all the information needed to perform a thorough comparison.

## UpsertReadDate (Java: `upsert(Set<PrivateConversation.Key> keys, Date readDate, boolean allowMoveReadDateForward)`)

- [ ] **Missing batch operation**: Java takes a `Set<PrivateConversation.Key>` and uses `mongoClient.upsert()` which operates on multiple keys at once. The Go version only takes a single `ownerID`/`targetID` pair and calls `UpdateOne`, not `UpdateMany`. This changes the method from a batch upsert to a single-record upsert.
- [ ] **Missing `allowMoveReadDateForward` parameter and conditional filter logic**: The Java version has critical conditional logic: when `allowMoveReadDateForward` is `false`, it adds a filter `ltOrNull(PrivateConversation.Fields.READ_DATE, readDate)` to ensure the read date is only updated if the existing date is earlier (or null). The Go version unconditionally sets the read date, meaning it can move the read date forward even when it shouldn't. This is a behavioral difference that can cause data corruption.
- [ ] **Filter uses wrong field reference**: Java filters on `DomainFieldName.ID` (the `_id` field) for the `in` query, matching against the full compound key objects. The Go version filters on `bson.M{"_id": filter}` where `filter` is a single key struct — correct for single-record but wrong for batch semantics.

## DeleteConversationsByOwnerIds (Java: `deleteConversationsByOwnerIds(Set<Long> ownerIds, @Nullable ClientSession session)`)

- [ ] **Method is a no-op stub**: The Go method has an empty body with no parameters, no filter, and no database call. The Java version accepts `Set<Long> ownerIds` and an optional `ClientSession`, builds a filter `in(PrivateConversation.Fields.ID_OWNER_ID, ownerIds)`, and calls `mongoClient.deleteMany()`. The Go version does nothing.

## FindConversations (Java: `findConversations(Collection<Long> ownerIds)`)

- [ ] **Method is a no-op stub**: The Go method has an empty body with no parameters, no filter, and no database call. The Java version accepts `Collection<Long> ownerIds`, builds a filter `in(PrivateConversation.Fields.ID_OWNER_ID, ownerIds)`, and calls `mongoClient.findMany()`. The Go version does nothing.

Note: `QueryPrivateConversations` in the Go file appears to be a partial implementation of `FindConversations`/`findConversations`, but it has different naming and the comment maps it to different Java methods (`queryPrivateConversations`), so it does not serve as a valid replacement for `findConversations(Collection<Long> ownerIds)`.

# ConversationService.java
*Checked methods: authAndUpsertGroupConversationReadDate(@NotNull Long groupId, @NotNull Long memberId, @Nullable @PastOrPresent Date readDate), authAndUpsertPrivateConversationReadDate(@NotNull Long ownerId, @NotNull Long targetId, @Nullable @PastOrPresent Date readDate), upsertGroupConversationReadDate(@NotNull Long groupId, @NotNull Long memberId, @Nullable @PastOrPresent Date readDate), upsertGroupConversationsReadDate(@NotNull Set<GroupConversation.GroupConversionMemberKey> keys, @Nullable @PastOrPresent Date readDate), upsertPrivateConversationReadDate(@NotNull Long ownerId, @NotNull Long targetId, @Nullable @PastOrPresent Date readDate), upsertPrivateConversationsReadDate(@NotNull Set<PrivateConversation.Key> keys, @Nullable @PastOrPresent Date readDate), queryGroupConversations(@NotNull Collection<Long> groupIds), queryPrivateConversationsByOwnerIds(@NotNull Set<Long> ownerIds), queryPrivateConversations(@NotNull Collection<Long> ownerIds, @NotNull Long targetId), queryPrivateConversations(@NotNull Set<PrivateConversation.Key> keys), deletePrivateConversations(@NotNull Set<PrivateConversation.Key> keys), deletePrivateConversations(@NotNull Set<Long> userIds, @Nullable ClientSession session), deleteGroupConversations(@Nullable Set<Long> groupIds, @Nullable ClientSession session), deleteGroupMemberConversations(@NotNull Collection<Long> userIds, @Nullable ClientSession session), authAndUpdateTypingStatus(@NotNull Long requesterId, boolean isGroupMessage, @NotNull Long toId)*

Now I have all the context needed for a thorough comparison. Let me systematically analyze each method.

## authAndUpsertGroupConversationReadDate

**Java logic:**
1. Check `isReadReceiptEnabled` → error if disabled
2. Query group type via `groupService.queryGroupTypeIfActiveAndNotDeleted(groupId, true)` → error if not found
3. Check `groupMemberService.isGroupMember(groupId, memberId, true)` → error if not member
4. Check `groupType.getEnableReadReceipt()` → error if disabled by group
5. Call `upsertGroupConversationReadDate` with `useServerTime ? new Date() : readDate`

**Go logic:**
1. No `isReadReceiptEnabled` check
2. No group type query / active+not-deleted check
3. No group membership check
4. No group-type read-receipt check
5. No `useServerTime` conditional — always uses the passed `readDate` directly
6. No `allowMoveReadDateForward` logic
7. Directly calls `groupConvRepo.UpsertReadDate`

## authAndUpsertPrivateConversationReadDate

**Java logic:**
1. Validate `ownerId` and `targetId` not null
2. Check `isReadReceiptEnabled` → error if disabled
3. Check `messageService.hasPrivateMessage(targetId, ownerId)` → only proceed if message exists
4. Call `upsertPrivateConversationReadDate` with `useServerTime ? new Date() : readDate`

**Go logic:**
1. No null validation (Go handles this differently, acceptable)
2. No `isReadReceiptEnabled` check
3. No `hasPrivateMessage` check — skips the authorization that requires an existing private message
4. No `useServerTime` conditional
5. No `allowMoveReadDateForward` logic
6. Directly calls `privateConvRepo.UpsertReadDate`

## upsertGroupConversationReadDate

**Java logic:**
1. Validate `groupId`, `memberId` not null, `readDate` past-or-present
2. Default `readDate` to `new Date()` if null
3. Call `groupConversationRepository.upsert(groupId, memberId, finalReadDate, allowMoveReadDateForward)`
4. On `DuplicateKeyException`: if `readDate == null` → empty, else → error "MOVING_READ_DATE_FORWARD_IS_DISABLED"

**Go logic:**
- Method body is completely empty (`func (s *ConversationService) UpsertGroupConversationReadDate() {}`) — **not implemented at all**

## upsertGroupConversationsReadDate

**Java logic:**
1. Validate `keys` not null, `readDate` past-or-present
2. Return empty if keys empty
3. Default `readDate` to `new Date()` if null
4. Group keys by `groupId` → map of `groupId` to `List<memberId>`
5. For each group, call `groupConversationRepository.upsert(groupId, memberIds, readDate)`
6. Use `Mono.whenDelayError` to run all in parallel with delay-error semantics

**Go logic:**
- Method body is completely empty — **not implemented at all**

## upsertPrivateConversationReadDate

**Java logic:**
1. Validate `ownerId`, `targetId` not null
2. Delegate to `upsertPrivateConversationsReadDate(Set.of(new PrivateConversation.Key(ownerId, targetId)), readDate)`

**Go logic:**
- Method body is completely empty — **not implemented at all**

## upsertPrivateConversationsReadDate

**Java logic:**
1. Validate `keys` not null, `readDate` past-or-present
2. Return empty if keys empty
3. Default `readDate` to `new Date()` if null
4. Call `privateConversationRepository.upsert(keys, finalReadDate, allowMoveReadDateForward)`
5. On `DuplicateKeyException`: if `readDate == null` → empty, else → error "MOVING_READ_DATE_FORWARD_IS_DISABLED"

**Go logic:**
- Method body is completely empty — **not implemented at all**

## queryGroupConversations

**Java logic:**
1. Validate `groupIds` not null
2. Return empty if `groupIds.isEmpty()`
3. Call `groupConversationRepository.findByIds(groupIds)`

**Go logic:**
1. No null validation (acceptable in Go)
2. Empty check present ✓
3. Calls `groupConvRepo.QueryGroupConversations(ctx, groupIDs)` ✓

This method appears correctly implemented.

## queryPrivateConversationsByOwnerIds

**Java logic:**
1. Validate `ownerIds` not null
2. Return empty if `ownerIds.isEmpty()`
3. Call `privateConversationRepository.findConversations(ownerIds)`

**Go logic:**
1. No null validation (acceptable in Go)
2. No empty check (the repo does it, but Java does it at service level)
3. Calls `s.privateConvRepo.QueryPrivateConversations(ctx, ownerIDs)` — but the Java version calls `findConversations(ownerIds)` which is a different method from `findByIds`. The Go `QueryPrivateConversations` queries by `"_id.oid": {"$in": ownerIDs}` which matches the Java `findConversations` behavior (finding all conversations where the owner matches). ✓

This method is essentially correct, though the Go repository stubs out `FindConversations()` as empty.

## queryPrivateConversations(ownerIds, targetId)

**Java logic:**
1. Validate `ownerIds` and `targetId` not null
2. Return empty if `ownerIds.isEmpty()`
3. Build `Set<PrivateConversation.Key>` from each `ownerId` paired with `targetId`
4. Delegate to `queryPrivateConversations(keys)` which calls `privateConversationRepository.findByIds(keys)`

**Go logic:**
- Not present as a separate method. The Go `QueryPrivateConversations` only takes `ownerIDs` and queries by `"_id.oid"`, which is equivalent to `queryPrivateConversationsByOwnerIds`, not this overloaded version that filters by both ownerId AND targetId.

## queryPrivateConversations(keys)

**Java logic:**
1. Validate `keys` not null
2. Return empty if keys empty
3. Call `privateConversationRepository.findByIds(keys)`

**Go logic:**
- Not present as a separate method. There is no method that takes conversation keys to query by specific key pairs.

## deletePrivateConversations(keys)

**Java logic:**
1. Validate `keys` not null
2. Return `ACKNOWLEDGED_DELETE_RESULT` if keys empty
3. Call `privateConversationRepository.deleteByIds(keys)`

**Go logic:**
- Method body is completely empty — **not implemented at all**

## deletePrivateConversations(userIds, session)

**Java logic:**
1. Validate `userIds` not null
2. Return `ACKNOWLEDGED_DELETE_RESULT` if empty
3. Call `privateConversationRepository.deleteConversationsByOwnerIds(userIds, session)`

**Go logic:**
- Method body is completely empty — **not implemented at all**
- The repository method `DeleteConversationsByOwnerIds()` is also empty

## deleteGroupConversations(groupIds, session)

**Java logic:**
1. No validation (accepts nullable)
2. Call `groupConversationRepository.deleteByIds(groupIds, session)`

**Go logic:**
- Method body is completely empty — **not implemented at all**

## deleteGroupMemberConversations(userIds, session)

**Java logic:**
1. Validate `userIds` not null
2. For each userId: query joined group IDs, then call `groupConversationRepository.deleteMemberConversations(groupIds, userId, session)`
3. Chains operations sequentially with `Mono.then()`

**Go logic:**
- Method body is completely empty — **not implemented at all**
- The repository method `DeleteMemberConversations()` is also empty

## authAndUpdateTypingStatus

**Java logic:**
1. Validate `requesterId` and `toId` not null
2. Check `isTypingStatusEnabled` → error if disabled
3. If group message: check group membership → return group member IDs
4. If private message: check relationship not blocked → return Set.of(toId)

**Go logic:**
- Method body is completely empty — **not implemented at all**

---

Additionally, there are bugs in the `UpsertReadDate` repository methods:

### GroupConversationRepository.UpsertReadDate
- **Java** passes `allowMoveReadDateForward` to control whether moving the read date forward is allowed. The Go version has no such parameter and always sets the date unconditionally. The Java version uses a conditional update that only sets the date if the new date is after the existing one (when `allowMoveReadDateForward` is enabled), and handles `DuplicateKeyException`.
- The Go version does `$set` unconditionally, meaning it always overwrites even if the new date is older.

### PrivateConversationRepository.UpsertReadDate
- Same issue: no `allowMoveReadDateForward` logic. Always does an unconditional `$set`.

## authAndUpsertGroupConversationReadDate
- [ ] Missing `isReadReceiptEnabled` check — Java returns error `UPDATING_READ_DATE_IS_DISABLED` when disabled
- [ ] Missing group active/not-deleted validation via `groupService.queryGroupTypeIfActiveAndNotDeleted(groupId, true)` — Java returns error `UPDATING_READ_DATE_OF_NONEXISTENT_GROUP_CONVERSATION` if group not found
- [ ] Missing group membership check via `groupMemberService.isGroupMember(groupId, memberId, true)` — Java returns error `NOT_GROUP_MEMBER_TO_UPDATE_READ_DATE_OF_GROUP_CONVERSATION`
- [ ] Missing group-type read-receipt check via `groupType.getEnableReadReceipt()` — Java returns error `UPDATING_READ_DATE_IS_DISABLED_BY_GROUP`
- [ ] Missing `useServerTime` conditional — Java uses `new Date()` (server time) when `useServerTime` is true, otherwise uses client-provided `readDate`
- [ ] Missing `allowMoveReadDateForward` logic in the underlying repo `UpsertReadDate` — the Go repo does an unconditional `$set`, while Java only updates if the new date is forward, and handles `DuplicateKeyException`

## authAndUpsertPrivateConversationReadDate
- [ ] Missing `isReadReceiptEnabled` check — Java returns error `UPDATING_READ_DATE_IS_DISABLED` when disabled
- [ ] Missing `hasPrivateMessage(targetId, ownerId)` authorization check — Java only proceeds if a private message exists between the users, otherwise returns `Mono.empty()`
- [ ] Missing `useServerTime` conditional — Java uses `new Date()` when `useServerTime` is true
- [ ] Missing `allowMoveReadDateForward` logic in the underlying repo `UpsertReadDate`

## upsertGroupConversationReadDate
- [ ] Method body is completely empty — not implemented. Java validates inputs, defaults `readDate` to `new Date()` when null, calls `groupConversationRepository.upsert()`, and handles `DuplicateKeyException` with `allowMoveReadDateForward` logic

## upsertGroupConversationsReadDate
- [ ] Method body is completely empty — not implemented. Java validates inputs, groups keys by `groupId`, and calls `groupConversationRepository.upsert(groupId, memberIds, readDate)` for each group with `Mono.whenDelayError`

## upsertPrivateConversationReadDate
- [ ] Method body is completely empty — not implemented. Java delegates to `upsertPrivateConversationsReadDate` with a single-key set

## upsertPrivateConversationsReadDate
- [ ] Method body is completely empty — not implemented. Java validates inputs, defaults `readDate` to `new Date()` when null, calls `privateConversationRepository.upsert(keys, finalReadDate, allowMoveReadDateForward)`, and handles `DuplicateKeyException`

## queryPrivateConversations(ownerIds, targetId)
- [ ] Method is missing entirely — Java builds a set of `PrivateConversation.Key` from each ownerId paired with targetId and delegates to `queryPrivateConversations(keys)`. The Go `QueryPrivateConversations` only filters by `ownerID`, not by `(ownerID, targetID)` pairs

## queryPrivateConversations(keys)
- [ ] Method is missing entirely — Java accepts `Set<PrivateConversation.Key>` and calls `privateConversationRepository.findByIds(keys)` to query by specific conversation key pairs. No equivalent exists in Go

## deletePrivateConversations(keys)
- [ ] Method body is completely empty — not implemented. Java validates keys, returns acknowledged result if empty, and calls `privateConversationRepository.deleteByIds(keys)`

## deletePrivateConversations(userIds, session)
- [ ] Method body is completely empty — not implemented. Java validates userIds, returns acknowledged result if empty, and calls `privateConversationRepository.deleteConversationsByOwnerIds(userIds, session)`

## deleteGroupConversations(groupIds, session)
- [ ] Method body is completely empty — not implemented. Java accepts nullable groupIds and nullable session, and calls `groupConversationRepository.deleteByIds(groupIds, session)`

## deleteGroupMemberConversations(userIds, session)
- [ ] Method body is completely empty — not implemented. Java iterates each userId, queries joined group IDs via `groupMemberService.queryUserJoinedGroupIds(userId)`, then calls `groupConversationRepository.deleteMemberConversations(groupIds, userId, session)`

## authAndUpdateTypingStatus
- [ ] Method body is completely empty — not implemented. Java validates inputs, checks `isTypingStatusEnabled`, then for group messages checks group membership and returns group member IDs, or for private messages checks relationship not-blocked status and returns the target ID

## GroupConversationRepository.UpsertReadDate
- [ ] Missing `allowMoveReadDateForward` parameter — Go always does an unconditional `$set` on the member's read date. Java uses this flag to conditionally only update if the new date is after the existing date, and relies on a MongoDB unique index / conditional update to prevent forward-date movement
- [ ] Missing `DuplicateKeyException` handling for the case when moving read date forward is disabled

## PrivateConversationRepository.UpsertReadDate
- [ ] Missing `allowMoveReadDateForward` parameter — Go always does an unconditional `$set` on read date. Java conditionally prevents moving the read date forward based on this flag
- [ ] Missing `DuplicateKeyException` handling for the case when moving read date forward is disabled

## GroupConversationRepository.DeleteMemberConversations
- [ ] Method body is completely empty — not implemented. Called by `deleteGroupMemberConversations`

## PrivateConversationRepository.DeleteConversationsByOwnerIds
- [ ] Method body is completely empty — not implemented. Called by `deletePrivateConversations(userIds, session)`

## PrivateConversationRepository.FindConversations
- [ ] Method body is completely empty — not implemented. Called by `queryPrivateConversationsByOwnerIds`

# ConversationSettingsService.java
*Checked methods: upsertPrivateConversationSettings(Long ownerId, Long userId, Map<String, Value> settings), upsertGroupConversationSettings(Long ownerId, Long groupId, Map<String, Value> settings), deleteSettings(Collection<Long> ownerIds, @Nullable ClientSession clientSession), unsetSettings(Long ownerId, @Nullable Set<Long> userIds, @Nullable Set<Long> groupIds, @Nullable Set<String> settingNames), querySettings(Long ownerId, @Nullable Collection<Long> userIds, @Nullable Collection<Long> groupIds, @Nullable Set<String> settingNames, @Nullable Date lastUpdatedDateStart)*

Now I have a complete picture. Let me compile the findings.

## UpsertPrivateConversationSettings

- [ ] **Method is a stub with no implementation** — The Go method `UpsertPrivateConversationSettings()` has no parameters, no return values, and an empty body. The Java version accepts `(Long ownerId, Long userId, Map<String, Value> settings)` and returns `Mono<Boolean>`. All core logic is missing: null validation, empty-settings early return, immutable-settings detection and conflict checking, relationship verification via `userRelationshipService.hasOneSidedRelationship`, `parseSettings` call, repository `upsertSettings` call, and the `modifiedCount > 0 || upsertedId != null` result mapping.
- [ ] **Missing `getTargetIdFromGroupId` utility** — The Java code uses `getTargetIdFromGroupId(groupId)` which negates group IDs (`-groupId`). No equivalent utility exists in the Go codebase.
- [ ] **Missing `getTargetIds` helper** — The Java code has a `getTargetIds(userIds, groupIds)` method that combines user IDs and negated group IDs into a single `List<Long>`. No equivalent exists in Go.

## UpsertGroupConversationSettings

- [ ] **Method is a stub with no implementation** — The Go method `UpsertGroupConversationSettings()` has no parameters, no return values, and an empty body. The Java version accepts `(Long ownerId, Long groupId, Map<String, Value> settings)` and returns `Mono<Boolean>`. All core logic is missing: null validation, empty-settings early return, immutable-settings detection and conflict checking, group-membership verification via `groupMemberService.isGroupMember(groupId, ownerId, false)`, `parseSettings` call, repository `upsertSettings` call with negated group ID, and the `modifiedCount > 0 || upsertedId != null` result mapping.

## DeleteSettings

- [ ] **Method does not exist on ConversationSettingsService** — The Java `deleteSettings(Collection<Long> ownerIds, @Nullable ClientSession clientSession)` method on `ConversationSettingsService` calls `conversationSettingsRepository.deleteByOwnerIds(ownerIds, clientSession)`. The Go `ConversationSettingsService` has no `DeleteSettings` method at all. The repository's `DeleteByOwnerIds()` is also a stub with no parameters, no body, and no return values. The Java logic of deleting by owner IDs and checking `deletedCount > 0` is entirely absent.

## UnsetSettings

- [ ] **Method does not exist on ConversationSettingsService** — The Java `unsetSettings(Long ownerId, @Nullable Set<Long> userIds, @Nullable Set<Long> groupIds, @Nullable Set<String> settingNames)` has no equivalent on the Go `ConversationSettingsService`. The Go `UserSettingsService.UnsetSettings` is a completely different method (for user settings, not conversation settings) and does not implement the conversation-specific logic. All core logic is missing: the `settingNames == null → use deletableSettings` fallback, the `ignoreUnknownSettingsOnDelete` branch that only checks non-deletable settings, the `!ignoreUnknownSettingsOnDelete` branch that also checks for unknown settings and produces combined error messages, the `getTargetIds(userIds, groupIds)` call to combine and negate IDs, and the `modifiedCount > 0` result mapping.

## QuerySettings

- [ ] **Method does not exist on ConversationSettingsService** — The Java `querySettings(Long ownerId, @Nullable Collection<Long> userIds, @Nullable Collection<Long> groupIds, @Nullable Set<String> settingNames, @Nullable Date lastUpdatedDateStart)` has no equivalent on the Go `ConversationSettingsService`. The Go `UserSettingsService.QuerySettings` is for user settings (takes a `bson.M` filter), not conversation settings. All core logic is missing: the three-way branching on `userIds`/`groupIds` to build composite keys, the `getTargetIdFromGroupId(groupId)` negation of group IDs, the construction of `ConversationSettings.Key` objects, the `findByIdAndSettingNames` calls with different key signatures (single ownerId vs collection of keys), and the `lastUpdatedDateStart` filtering parameter.

# GroupBlocklistController.java
*Checked methods: addGroupBlockedUser(@RequestBody AddGroupBlockedUserDTO addGroupBlockedUserDTO), queryGroupBlockedUsers(@QueryParam(required = false), queryGroupBlockedUsers(@QueryParam(required = false), updateGroupBlockedUsers(List<GroupBlockedUser.Key> keys, @RequestBody UpdateGroupBlockedUserDTO updateGroupBlockedUserDTO), deleteGroupBlockedUsers(List<GroupBlockedUser.Key> keys)*

Now I have all the information needed for a thorough comparison. Here are the findings:

## AddGroupBlockedUser

- [ ] **All logic is missing** — The Go method body is `return nil` with no implementation. The Java version calls `groupBlocklistService.addBlockedUser(addGroupBlockedUserDTO.groupId(), addGroupBlockedUserDTO.userId(), addGroupBlockedUserDTO.requesterId(), addGroupBlockedUserDTO.blockDate())` and wraps the result with `HttpHandlerResult.okIfTruthy()`. The Go stub doesn't extract any fields (groupId, userId, requesterId, blockDate) from the DTO and doesn't call the service layer at all.

## QueryGroupBlockedUsers

- [ ] **Non-paginated query method is missing** — Java has two overloads of `queryGroupBlockedUsers`: one without `page` (non-paginated, returns `Collection<GroupBlockedUser>`) and one with `page` (paginated, returns `PaginationDTO<GroupBlockedUser>`). The Go code has `QueryGroupBlockedUsers(page, size *int)` which takes page/size parameters (looks like the paginated version) and `QueryGroupBlockedUsersWithQuery(...)` which also takes page/size. There is no non-paginated query method (i.e., one that corresponds to the Java `@GetMapping` without `"page"` that returns all results with just a `size` limit and `offset=0`).
- [ ] **Query parameters missing from `QueryGroupBlockedUsers`** — The Java non-paginated version accepts `groupIds`, `userIds`, `blockDateStart`, `blockDateEnd`, `requesterIds`, and `size`. The Go `QueryGroupBlockedUsers` only takes `page` and `size`, missing all filter parameters.
- [ ] **`QueryGroupBlockedUsersWithQuery` is missing `blockDateStart`/`blockDateEnd` parameters for date range filtering** — Looking more closely, these ARE present as `blockDateStart, blockDateEnd *int64`. However, this method conflates the paginated and non-paginated Java versions into one method.
- [ ] **Paginated version missing count query** — The Java paginated version (`@GetMapping("page")`) performs both a `countBlockedUsers` call and a `queryBlockedUsers` call, then returns `HttpHandlerResult.page(count, userFlux)`. The Go code has no implementation at all, so this count-then-query pattern is absent.
- [ ] **`getPageSize(size)` logic missing** — Both Java overloads call `getPageSize(size)` to normalize/validate the page size. The Go stubs don't implement this.
- [ ] **`DateRange.of(blockDateStart, blockDateEnd)` construction missing** — Java wraps the two date parameters into a `DateRange` object before passing to the service. The Go code passes raw `*int64` values with no such wrapping.

## UpdateGroupBlockedUsers

- [ ] **All logic is missing** — The Go method body is `return nil`. The Java version calls `groupBlocklistService.updateBlockedUsers(CollectionUtil.newSet(keys), updateGroupBlockedUserDTO.blockDate(), updateGroupBlockedUserDTO.requesterId())` and maps the result with `UpdateResultDTO::get`, then wraps with `HttpHandlerResult.okIfTruthy()`. None of this is implemented.
- [ ] **Keys parameter type is `[]any` instead of a typed key struct** — Java uses `List<GroupBlockedUser.Key>` where `Key` contains `groupId` and `userId`. The Go version uses `[]any`, losing type safety and field structure.
- [ ] **Missing `CollectionUtil.newSet(keys)` conversion** — Java converts the list of keys to a Set. The Go code has no such conversion.

## DeleteGroupBlockedUsers

- [ ] **All logic is missing** — The Go method body is `return nil`. The Java version calls `groupBlocklistService.deleteBlockedUsers(CollectionUtil.newSet(keys))`, maps with `DeleteResultDTO::get`, and wraps with `HttpHandlerResult.okIfTruthy()`. None of this is implemented.
- [ ] **Keys parameter type is `[]any` instead of a typed key struct** — Java uses `List<GroupBlockedUser.Key>` (composite key with `groupId` and `userId`). The Go version uses `[]any`, losing the typed key structure.
- [ ] **Missing `CollectionUtil.newSet(keys)` conversion** — Java converts the list of keys to a Set before passing to the service. The Go code has no such conversion.

# GroupController.java
*Checked methods: addGroup(@RequestBody AddGroupDTO addGroupDTO), queryGroups(@QueryParam(required = false), queryGroups(@QueryParam(required = false), countGroups(@QueryParam(required = false), updateGroups(Set<Long> ids, @RequestBody UpdateGroupDTO updateGroupDTO), deleteGroups(@QueryParam(required = false)*

I now have a comprehensive understanding of both codebases. Let me compile the findings.

## addGroup

- [ ] **All method logic is missing** — The Go `AddGroup` method (`group_controllers.go:27`) is a stub that accepts `any` and returns `nil`. The Java `addGroup` calls `groupService.authAndCreateGroup(...)` with 12 parameters (creatorId, ownerId with fallback-to-creatorId logic, name, intro, announcement, minimumScore, typeId, creationDate, deletionDate, muteEndDate, isActive). None of this is implemented.
- [ ] **Missing `ownerId` fallback-to-`creatorId` logic** — Java: `ownerId == null ? addGroupDTO.creatorId() : ownerId`. The Go stub doesn't have this.
- [ ] **Go `AuthAndCreateGroup` service method is incomplete** — Java `authAndCreateGroup` takes 12 parameters including `groupTypeId`, `creationDate`, `deletionDate`, `muteEndDate`, `isActive`. The Go `AuthAndCreateGroup` (`group_service.go:115`) only takes `creatorID, ownerID, name, intro, minimumScore` — missing `announcement`, `groupTypeId`, `creationDate`, `deletionDate`, `muteEndDate`, `isActive`.
- [ ] **Missing validation** — Java validates `notNull(ownerId)`, `min(minimumScore, 0)`, `pastOrPresent(creationDate)`, `pastOrPresent(deletionDate)`, `before(creationDate, deletionDate)`. None of these validations exist in Go.
- [ ] **Missing default groupTypeId** — Java defaults `groupTypeId` to `DEFAULT_GROUP_TYPE_ID` when null. Go `CreateGroup` never sets `TypeID` on the group struct.
- [ ] **Missing field assignments in `CreateGroup`** — Java `createGroup` sets `announcement`, `groupTypeId`, `deletionDate`, `muteEndDate`, `isActive`, `creationDate` (if provided, else now). Go `CreateGroup` (`group_service.go:51`) only sets `ID`, `CreatorID`, `OwnerID`, `Name`, `Intro`, `MinimumScore`, `CreationDate` — missing `Announcement`, `TypeID`, `DeletionDate`, `MuteEndDate`, `IsActive`.

## queryGroups (non-page)

- [ ] **All method logic is missing** — The Go `QueryGroups` method (`group_controllers.go:28`) is a stub returning `nil`.
- [ ] **Missing query parameters** — Java accepts 15 parameters (ids, typeIds, creatorIds, ownerIds, isActive, 4 date ranges with start/end, memberIds, size). Go `QueryGroupsWithQuery` (`group_controllers.go:29`) accepts `lastUpdatedDateStart/End` but is missing `lastUpdatedDateStart` and `lastUpdatedDateEnd` is misspelled — the method signature shows `deletionDateEnd, muteEndDateStart, muteEndDateEnd` but omits `lastUpdatedDateStart`, `lastUpdatedDateEnd`.
- [ ] **`memberIds` filtering logic missing** — Java `queryGroups` in `GroupService` intersects results with group IDs from member IDs via `queryGroupIdsFromGroupIdsAndMemberIds`. No equivalent exists in Go.
- [ ] **Repository `QueryGroups` is incomplete** — The Go `group_repository.go:57` `QueryGroups` only supports filtering by `groupIDs`, `name`, `lastUpdatedDate` — missing `typeIds`, `creatorIds`, `ownerIds`, `isActive`, all 4 date ranges. The Java repository `findGroups` supports all these filters.
- [ ] **`getPageSize` pagination defaulting missing** — Java calls `size = getPageSize(size)` to apply a configured default. No equivalent exists in Go.
- [ ] **Missing the `0` page offset** — Java passes `page=0` explicitly for the non-page query endpoint. No equivalent in Go.

## queryGroups (page)

- [ ] **All method logic is missing** — No separate page-based query method exists in Go. The `QueryGroups(page, size *int)` stub (`group_controllers.go:28`) takes no filter parameters at all.
- [ ] **Missing `countGroups` call for pagination** — Java calls `groupService.countGroups(...)` for total count and then `groupService.queryGroups(...)` for the page data, returning via `HttpHandlerResult.page(count, groupsFlux)`. No equivalent in Go.
- [ ] **Missing `countGroups` with full parameter set in service** — Java `countGroups` accepts (ids, typeIds, creatorIds, ownerIds, isActive, 4 date ranges, memberIds). Go `CountGroups` (`group_service.go:584`) only accepts a single `dateRange` parameter.

## countGroups

- [ ] **All method logic is missing** — No `CountGroups` method exists on the Go controller at all. The Java `countGroups` is a separate `@GetMapping("count")` endpoint returning `GroupStatisticsDTO`.
- [ ] **Missing `DivideBy` branch** — Java handles two modes: `NOOP` (returns raw counts for deletedGroups, groupsThatSentMessages, createdGroups) and a `DivideBy` mode (returns `StatisticsRecordDTO` lists with date-divided counts). No equivalent exists in Go.
- [ ] **Missing conditional count logic** — Java has complex conditional logic: only counts deleted groups if `deletedStartDate/End` provided, only counts groups that sent messages if `sentMessageStartDate/End` provided, always counts created groups (unless other counts present and no created date filter). No equivalent in Go.
- [ ] **Missing `messageService.countGroupsThatSentMessages` call** — Java calls `messageService.countGroupsThatSentMessages(DateRange)`. No equivalent integration in Go.
- [ ] **Missing `checkAndQueryBetweenDate` call** — Java uses `BaseController.checkAndQueryBetweenDate` for the `DivideBy` mode to split date ranges. No equivalent in Go.
- [ ] **Missing `GroupStatisticsDTO` population** — The Go DTO struct exists but is never populated or returned.

## updateGroups

- [ ] **All method logic is missing** — The Go `UpdateGroups` method (`group_controllers.go:32`) is a stub returning `nil`.
- [ ] **Missing `successorId` branching logic** — Java checks `successorId == null`: if null, calls `updateGroupsInformation`; if non-null, calls `checkAndTransferGroupOwnership`. No equivalent in Go.
- [ ] **Missing `creationDate`, `deletionDate`, `muteEndDate`, `userDefinedAttributes`, `session` parameters** — Java `updateGroupsInformation` takes 15 parameters including these. Go `UpdateGroupsInformation` (`group_service.go:441`) only takes 9 parameters — missing `creationDate`, `deletionDate`, `muteEndDate`, `userDefinedAttributes`, `session`.
- [ ] **Missing field assignments for `creationDate`, `deletionDate`, `muteEndDate`** — Java's repository `updateGroups` sets `creationDate` (cd), `deletionDate` (dd), and `muteEndDate` (med) fields. Go's `UpdateGroupsInformation` doesn't set these fields in the update bson.M.
- [ ] **Missing validation** — Java validates `notEmpty(groupIds)`, `min(minimumScore, 0)`, `pastOrPresent(creationDate)`, `pastOrPresent(deletionDate)`, `before(creationDate, deletionDate)`. None exist in Go.
- [ ] **Missing Elasticsearch integration** — Java handles Elasticsearch sync when updating group names. No equivalent in Go.
- [ ] **Missing `UpdateResult` return** — Java returns `Mono<UpdateResult>` (MongoDB update result with matched/modified counts). Go's `UpdateGroupsInformation` returns only `error`, losing the update statistics.
- [ ] **`checkAndTransferGroupOwnership` returns different type** — Java returns `Mono<UpdateResult>` from `checkAndTransferGroupOwnership(Set<Long> groupIds, ...)`. Go's `CheckAndTransferGroupOwnership` returns `error` only, losing the matched/modified count aggregation that Java builds by collecting signals.

## deleteGroups

- [ ] **All method logic is missing** — The Go `DeleteGroups` method (`group_controllers.go:35`) is a stub returning `nil`.
- [ ] **Missing `deleteLogically` parameter handling** — Java defaults `deleteLogically` to a configured `deleteGroupLogicallyByDefault` when null, then branches: logical deletion calls `updateGroupsDeletionDate`, physical deletion calls `deleteByIds`. Go's `DeleteGroupsAndGroupMembers` (`group_service.go:291`) ignores the `deleteLogically` parameter entirely — the signature takes `session mongo.SessionContext` instead of `deleteLogically *bool`, and always does hard delete.
- [ ] **Missing transaction wrapping** — Java wraps the entire delete operation in `groupRepository.inTransaction(...)` with retry. Go has no transaction support.
- [ ] **Missing `DeleteResultDTO` return** — Java maps the MongoDB `DeleteResult` to `DeleteResultDTO` and returns it. Go returns `error` only.
- [ ] **Missing cascading deletes for message sequence IDs and conversations** — Java cascades to delete message sequence IDs (`messageService.deleteGroupMessageSequenceIds`) and group conversations (`conversationService.deleteGroupConversations`). Go only cascades to group members and group versions.
- [ ] **Missing counter increment** — Java increments `deletedGroupsCounter`. No equivalent in Go.
- [ ] **Missing Elasticsearch cleanup** — Java deletes group docs from Elasticsearch. No equivalent in Go.
- [ ] **Missing `deleteAllGroupDocs` for null ids** — Java handles `groupIds == null` to delete all groups and all Elasticsearch docs. Go doesn't support null/empty IDs as "delete all".

# GroupInvitationController.java
*Checked methods: addGroupInvitation(@RequestBody AddGroupInvitationDTO addGroupInvitationDTO), queryGroupInvitations(@QueryParam(required = false), queryGroupInvitations(@QueryParam(required = false), updateGroupInvitations(Set<Long> ids, @RequestBody UpdateGroupInvitationDTO updateGroupInvitationDTO), deleteGroupInvitations(@QueryParam(required = false)*

## AddGroupInvitation
- [ ] Method body is a stub returning `nil`. It does not call `groupInvitationService.createGroupInvitation()` with the DTO fields (`id`, `groupId`, `inviterId`, `inviteeId`, `content`, `status`, `creationDate`, `responseDate`), nor does it map the result into a `GroupInvitationDTO` with `getEntityExpirationDate()`, nor return it via `okIfTruthy`.

## QueryGroupInvitations
- [ ] Method body is a stub returning `nil`. It does not call `groupInvitationService.queryInvitations()` with the filter parameters (`ids`, `groupIds`, `inviterIds`, `inviteeIds`, `statuses`, date ranges for `creationDate`, `responseDate`, `expirationDate`, `page=0`, `size`), nor map results to `GroupInvitationDTO` with `getEntityExpirationDate()`, nor return via `okIfTruthy`.

## QueryGroupInvitationsWithQuery
- [ ] Method body is a stub returning `nil`. It does not call `groupInvitationService.countInvitations()` for the total count, nor call `groupInvitationService.queryInvitations()` with all filter parameters and pagination (`page`, `size`), nor map results to `GroupInvitationDTO` with `getEntityExpirationDate()`, nor return via `HttpHandlerResult.page(count, invitationFlux)`.

## UpdateGroupInvitations
- [ ] Method body is a stub returning `nil`. It does not call `groupInvitationService.updateInvitations()` with `ids` and the DTO fields (`inviterId`, `inviteeId`, `content`, `status`, `creationDate`, `responseDate`), nor map the result to `UpdateResultDTO`, nor return via `okIfTruthy`.

## DeleteGroupInvitations
- [ ] Method body is a stub returning `nil`. It does not call `groupInvitationService.deleteInvitations(ids)`, nor map the result to `DeleteResultDTO`, nor return via `okIfTruthy`.

# GroupJoinRequestController.java
*Checked methods: addGroupJoinRequest(@RequestBody AddGroupJoinRequestDTO addGroupJoinRequestDTO), queryGroupJoinRequests(@QueryParam(required = false), queryGroupJoinRequests(@QueryParam(required = false), updateGroupJoinRequests(Set<Long> ids, @RequestBody UpdateGroupJoinRequestDTO updateGroupJoinRequestDTO), deleteGroupJoinRequests(@QueryParam(required = false)*

Now I have all the information needed for the comparison. The Go code consists entirely of stub methods that return `nil` with no implementation. Let me provide the detailed analysis.

## AddGroupJoinRequest

- [ ] **All core logic is missing**: The Go method is a stub that returns `nil` with no implementation. The Java version calls `groupJoinRequestService.createGroupJoinRequest()` with 9 parameters (id, groupId, requesterId, responderId, content, status, creationDate, responseDate, responseReason) and wraps the result in a `GroupJoinRequestDTO` that includes the entity expiration date. None of this logic exists in Go.

## QueryGroupJoinRequests (non-paginated)

- [ ] **All core logic is missing**: The Go method is a stub that returns `nil`. The Java version calls `groupJoinRequestService.queryJoinRequests()` with 10 parameters (ids, groupIds, requesterIds, responderIds, statuses, 3 DateRanges for creation/response/expiration dates, page=0, size) and maps each result to `GroupJoinRequestDTO` with expiration date. None of this logic exists in Go.

- [ ] **Missing parameter handling**: Java passes `expirationDateStart` and `expirationDateEnd` as query parameters with a `DateRange` filter. The Go signature for `QueryGroupJoinRequestsWithQuery` does include these parameters, but they are unused since the method body is empty.

- [ ] **Missing page size normalization**: Java calls `size = getPageSize(size)` to apply default page size from configuration. This is absent in Go.

## QueryGroupJoinRequests (paginated, "page" endpoint)

- [ ] **All core logic is missing**: The Java version has a separate `@GetMapping("page")` overload that first calls `countJoinRequests()` to get a total count, then calls `queryJoinRequests()` with pagination, and wraps the result via `HttpHandlerResult.page()`. The Go version has no separate paginated endpoint — there is only `QueryGroupJoinRequests` and `QueryGroupJoinRequestsWithQuery`, neither of which implements count-then-query pagination logic.

## UpdateGroupJoinRequests

- [ ] **All core logic is missing**: The Go method is a stub that returns `nil`. The Java version calls `groupJoinRequestService.updateJoinRequests()` with 7 parameters from the DTO (requesterId, responderId, content, status, creationDate, responseDate) plus the `ids` set. None of this logic exists in Go.

## DeleteGroupJoinRequests

- [ ] **All core logic is missing**: The Go method is a stub that returns `nil`. The Java version calls `groupJoinRequestService.deleteJoinRequests(ids)` and wraps the result in `DeleteResultDTO`. None of this logic exists in Go.

# GroupMemberController.java
*Checked methods: addGroupMember(@RequestBody AddGroupMemberDTO addGroupMemberDTO), queryGroupMembers(@QueryParam(required = false), queryGroupMembers(@QueryParam(required = false), updateGroupMembers(List<GroupMember.Key> keys, @RequestBody UpdateGroupMemberDTO updateGroupMemberDTO), deleteGroupMembers(@QueryParam(required = false)*

Now I have all the information needed. Let me analyze the Go controller and repository against the Java controller and service/repository.

## addGroupMember

**Java controller** calls `groupMemberService.addGroupMember(groupId, userId, role, name, joinDate, muteEndDate, null)` which:
1. Validates groupId, userId, role are not null
2. Validates role is valid (not OWNER)
3. Validates joinDate is past or present
4. If joinDate is null, sets it to `new Date()`
5. Creates GroupMember with groupId, userId, name, role, joinDate, muteEndDate
6. Inserts via repository (not upsert)
7. Updates group members version
8. Caches the member

**Go controller** (`GroupMemberController`) - has **no `AddGroupMember` method at all**. The controller file only has `QueryGroupMembers`, `QueryGroupMembersWithQuery`, and `DeleteGroupMembers` - the `addGroupMember` endpoint is completely missing.

## queryGroupMembers (non-paginated)

**Java controller** calls `groupMemberService.queryGroupMembers(groupIds, userIds, roles, DateRange.of(joinDateStart, joinDateEnd), DateRange.of(muteEndDateStart, muteEndDateEnd), 0, size)`.

**Go controller** has `QueryGroupMembers` (page/size only, no filters) and `QueryGroupMembersWithQuery` (with filters). The `QueryGroupMembersWithQuery` takes `groupIds, userIds, roles, joinDateStart, joinDateEnd, muteEndDateStart, muteEndDateEnd, page, size` - but both methods are stubs returning `nil`.

The repository's `FindGroupsMembers` is also a stub returning `nil, nil`.

## queryGroupMembers (paginated - "page" endpoint)

**Java controller** calls `countMembers` + `queryGroupMembers` and returns a `PaginationDTO`. This endpoint is completely missing from the Go controller.

## updateGroupMembers

**Java controller** calls `groupMemberService.updateGroupMembers(CollectionUtil.newSet(keys), name, role, joinDate, muteEndDate, null, true)`.

**Go controller** has **no `UpdateGroupMembers` method**. The `UpdateGroupMembers` is missing entirely from the controller.

**Go repository** `UpdateGroupMembers` has the muteEndDate unset logic matching Java, but the service layer `UpdateGroupMemberRole` is the only caller and it doesn't handle all the fields from the admin controller DTO.

## deleteGroupMembers

**Java controller**: If keys are empty, calls `deleteGroupMembers(true)` (deletes ALL group members). If keys are present, calls `deleteGroupMembers(CollectionUtil.newSet(keys), null, true)`.

**Go controller** `DeleteGroupMembers(keys []any, successorId *int64, quitAfterTransfer *bool)` - this is a stub returning `nil`. The parameters include `successorId` and `quitAfterTransfer` which are from a different Java method (`authAndDeleteGroupMembers`), not from the admin controller's `deleteGroupMembers`. But critically, it doesn't implement the "delete all if keys empty" vs "delete by keys" logic.

Here's the full bug report:

## addGroupMember
- [ ] **Missing method**: The Go `GroupMemberController` has no `AddGroupMember` method at all. The Java controller has a `@PostMapping addGroupMember` that accepts `AddGroupMemberDTO` with `groupId`, `userId`, `name`, `role`, `joinDate`, `muteEndDate` and calls `groupMemberService.addGroupMember(...)`. This endpoint is completely absent from the Go controller.
- [ ] **Missing `name` field assignment in service**: The Go `GroupMemberService.AddGroupMember` does not accept or set the `name` parameter. The Java `addGroupMember` service method accepts `name` and passes it into the `GroupMember` constructor (`new GroupMember(groupId, userId, name, groupMemberRole, joinDate, muteEndDate)`). The Go service ignores the `name` field entirely.
- [ ] **Missing `joinDate` passthrough from DTO**: The Java controller passes `addGroupMemberDTO.joinDate()` directly to the service, where it's used if non-null (otherwise defaults to `new Date()`). The Go service always sets `JoinDate = &now` and has no parameter to accept a custom joinDate from the DTO.
- [ ] **Missing group version update**: The Java `addGroupMember` service calls `groupVersionService.updateMembersVersion(groupId)` after inserting. The Go `GroupMemberService.AddGroupMember` does not update the group members version after adding a member.
- [ ] **Uses upsert instead of insert**: The Go repository `AddGroupMember` uses `UpdateOne` with `SetUpsert(true)`, whereas the Java code uses `groupMemberRepository.insert(groupMember, session)` which is a plain insert. Using upsert means an existing member's data would be silently overwritten instead of potentially failing.
- [ ] **Missing cache update**: The Java service calls `cacheMember()` after insert. The Go `AddGroupMember` service method does not update the member cache.

## queryGroupMembers
- [ ] **Non-paginated query method is a stub**: The Go `GroupMemberController.QueryGroupMembers` and `QueryGroupMembersWithQuery` both return `nil` without any implementation. In Java, this calls `groupMemberService.queryGroupMembers(...)` which queries MongoDB with filters for groupIds, userIds, roles, date ranges, and pagination.
- [ ] **Repository `FindGroupsMembers` is a stub**: The Go repository `FindGroupsMembers` returns `nil, nil` with a TODO comment. This is the repository method that should implement the actual MongoDB query with filters matching the Java `findGroupsMembers`.
- [ ] **Missing paginated "page" endpoint**: The Java controller has a separate `@GetMapping("page") queryGroupMembers` method that returns `PaginationDTO<GroupMember>` with a `count` + `query` pattern. This endpoint is completely missing from the Go controller.
- [ ] **Missing `CountMembers` with filters**: The Go repository `CountMembers` only takes a single `groupID int64` parameter. The Java `countMembers` used by the paginated query takes `groupIds, userIds, roles, joinDateRange, muteEndDateRange` for filtering. The Go version cannot support the same query filtering for pagination.

## updateGroupMembers
- [ ] **Missing method in controller**: The Go `GroupMemberController` has no `UpdateGroupMembers` method. The Java controller has a `@PutMapping updateGroupMembers` that accepts `List<GroupMember.Key> keys` and `UpdateGroupMemberDTO` and calls the service. This endpoint is completely absent from the Go controller.
- [ ] **Missing service-level admin update method**: The Go `GroupMemberService` has no method matching the Java `updateGroupMembers(Set<GroupMember.Key> keys, name, role, joinDate, muteEndDate, session, updateGroupMembersVersion)` that the admin controller calls. The Go `UpdateGroupMemberRole` only handles role updates for a single member.
- [ ] **Missing group version update in repository-layer update**: The Go repository `UpdateGroupMembers` does not trigger a group version update. In Java, the service-layer `updateGroupMembers` calls `groupVersionService.updateMembersVersion(...)` after the repository update.
- [ ] **Missing cache update**: The Java service calls `updateMembersCache(keys, name, role, joinDate, muteEndDate)` after a successful update. The Go code has no equivalent cache invalidation/update for the update path (the repository method doesn't touch the cache, and there's no service-level method for this).

## deleteGroupMembers
- [ ] **Controller method is a stub returning nil**: The Go `GroupMemberController.DeleteGroupMembers` returns `nil` with no implementation. It should implement the same logic as Java: if keys are empty, delete all group members; otherwise delete by the specified keys.
- [ ] **Missing "delete all" path**: The Java controller calls `groupMemberService.deleteGroupMembers(true)` when keys are empty, which deletes ALL group members in the database. The Go controller has no logic to handle this case.
- [ ] **Incorrect parameters**: The Go controller `DeleteGroupMembers` takes `successorId *int64, quitAfterTransfer *bool` parameters which belong to `authAndDeleteGroupMembers` (a client-facing method), not the admin `deleteGroupMembers`. The admin endpoint should only take `keys` as a query parameter.
- [ ] **Missing `DeleteAll` repository method**: The Go `GroupMemberRepository` has no `DeleteAll` method. Java calls `groupMemberRepository.deleteAll()` for the "delete all" path. The Go repo only has `DeleteByIds` and `DeleteByGroupIDs`.

# GroupQuestionController.java
*Checked methods: queryGroupJoinQuestions(@QueryParam(required = false), queryGroupJoinQuestions(@QueryParam(required = false), addGroupJoinQuestion(@RequestBody AddGroupJoinQuestionDTO addGroupJoinQuestionDTO), updateGroupJoinQuestions(Set<Long> ids, @RequestBody UpdateGroupJoinQuestionDTO updateGroupJoinQuestionDTO), deleteGroupJoinQuestions(@QueryParam(required = false)*

Now I have all the information needed to provide a complete review. The Go controller file contains only stub methods (all return `nil`), so every method has bugs relative to the Java implementation. Let me enumerate them precisely.

## QueryGroupJoinQuestions (non-paginated)

- [ ] **All logic is missing**: The method returns `nil` instead of calling `FindQuestions(ids, groupIds, 0, size, true)` (with `page=0`). In Java, it queries questions with the provided `ids`, `groupIds`, `size` parameters (after `getPageSize`), and returns via `HttpHandlerResult.okIfTruthy`. The Go version accepts `page` and `size` parameters but doesn't match the Java signature: Java takes `ids`, `groupIds`, and `size` (no `page`), while Go splits into `QueryGroupJoinQuestions(page, size)` and `QueryGroupJoinQuestionsWithQuery(ids, groupIds, score, page, size)`.
- [ ] **Method signature mismatch - missing query parameters**: The Java non-paginated `queryGroupJoinQuestions` accepts `ids` and `groupIds` as query parameters. The Go version `QueryGroupJoinQuestions(page, size *int)` has no `ids` or `groupIds` parameters.
- [ ] **Missing `getPageSize` logic**: Java normalizes the `size` parameter via `getPageSize` (defaults to `defaultAvailableRecordsPerRequest` if null/<=0, capped at `maxAvailableRecordsPerRequest`). Go does not implement this size normalization.
- [ ] **Hardcoded `page=0` not passed correctly**: Java calls `queryGroupJoinQuestions(ids, groupIds, 0, size, true)` — passing `0` as the page. The Go non-paginated method receives `page` as a parameter instead of hardcoding 0.
- [ ] **Missing `withAnswers=true`**: Java passes `true` for the `withAnswers` parameter. The Go `FindQuestions` service method accepts `withAnswers`, but the controller stub doesn't call it.

## QueryGroupJoinQuestionsWithQuery (paginated)

- [ ] **All logic is missing**: Returns `nil` instead of calling `CountQuestions` + `FindQuestions` to build a paginated response. In Java, the paginated endpoint calls both `countGroupJoinQuestions(ids, groupIds)` and `queryGroupJoinQuestions(ids, groupIds, page, size, true)` to return a `PaginationDTO`.
- [ ] **Spurious `score` parameter**: The Go method signature has `score *int` which does not exist in the Java paginated method (`queryGroupJoinQuestions` page variant accepts `ids`, `groupIds`, `page`, `size`). The `score` parameter is not used anywhere in the Java controller or service for querying questions.
- [ ] **Missing count query**: The Java paginated version calls `countGroupJoinQuestions(ids, groupIds)` before fetching data to provide total count for pagination. The Go stub doesn't call `CountQuestions`.

## AddGroupJoinQuestion

- [ ] **All logic is missing**: Returns `nil` instead of calling a create method. In Java, it calls `groupQuestionService.createGroupJoinQuestions(groupId, List.of(new NewGroupQuestion(question, answers, score)))` and maps the result via `.map(List::getFirst)`.
- [ ] **DTO typed as `any`**: The Go method parameter is `addGroupJoinQuestionDTO any` instead of `*dto.AddGroupJoinQuestionDTO`, losing type safety and making it impossible to access `GroupId`, `Question`, `Answers`, `Score` fields.
- [ ] **Missing `createGroupJoinQuestions` service call**: Java calls `createGroupJoinQuestions(groupId, questions)` which is a batch create method returning `List<GroupJoinQuestion>`. The Go `CreateJoinQuestion` service method only creates a single question. There is no `createGroupJoinQuestions` (plural) equivalent in the Go service.
- [ ] **Missing `.map(List::getFirst)` equivalent**: Java creates a list with one question and returns only the first element. Go stub doesn't implement this.
- [ ] **Missing `NewGroupQuestion` wrapper**: Java wraps the DTO fields in a `NewGroupQuestion` record before passing to the service. Go doesn't have an equivalent intermediate type.

## UpdateGroupJoinQuestions

- [ ] **All logic is missing**: Returns `nil` instead of calling `UpdateQuestions`. In Java, it calls `groupQuestionService.updateGroupJoinQuestions(ids, groupId, question, answers, score)` and maps the result via `.map(UpdateResultDTO::get)`.
- [ ] **DTO typed as `any`**: The Go method parameter is `updateGroupJoinQuestionDTO any` instead of `*dto.UpdateGroupJoinQuestionDTO`, losing type safety and making it impossible to access `GroupId`, `Question`, `Answers`, `Score` fields.
- [ ] **Missing `UpdateResultDTO` transformation**: Java wraps the `UpdateResult` from the service into an `UpdateResultDTO` via `.map(UpdateResultDTO::get)`. The Go stub doesn't implement this response wrapping.

## DeleteGroupJoinQuestions

- [ ] **All logic is missing**: Returns `nil` instead of calling a batch delete. In Java, it calls `groupQuestionService.deleteGroupJoinQuestions(ids)` and maps the result via `.map(DeleteResultDTO::get)`.
- [ ] **Missing batch delete repository method**: Java calls `deleteGroupJoinQuestions(ids)` which delegates to `groupQuestionRepository.deleteByIds(ids)` — a batch delete by a set of IDs. The Go repository only has `Delete(ctx, questionID int64)` for single deletion. There is no `DeleteByIds` or batch delete method in the Go repository.
- [ ] **Missing `DeleteResultDTO` transformation**: Java wraps the `DeleteResult` into a `DeleteResultDTO` via `.map(DeleteResultDTO::get)`. The Go stub doesn't implement this response wrapping.

# GroupTypeController.java
*Checked methods: addGroupType(@RequestBody AddGroupTypeDTO addGroupTypeDTO), queryGroupTypes(@QueryParam(required = false), queryGroupTypes(int page, @QueryParam(required = false), updateGroupType(Set<Long> ids, @RequestBody UpdateGroupTypeDTO updateGroupTypeDTO), deleteGroupType(Set<Long> ids)*

Now I have all the information needed. Let me compile the review.

## addGroupType

- [ ] **Controller method body is completely unimplemented**: The Go `GroupTypeController.AddGroupType` just returns `nil` with parameter type `any` instead of using the `AddGroupTypeDTO`. In Java, the controller extracts all 10 fields from `AddGroupTypeDTO` and passes them to `groupTypeService.addGroupType(null, name, groupSizeLimit, invitationStrategy, joinStrategy, groupInfoUpdateStrategy, memberInfoUpdateStrategy, guestSpeakable, selfInfoUpdatable, enableReadReceipt, messageEditable)`.
- [ ] **Service `AddGroupType` is missing ID generation**: In Java, when `id == null`, the service generates one via `node.nextLargeGapId(ServiceType.GROUP_TYPE)`. The Go `AddGroupType` service method does not generate an ID — it just inserts whatever `po.GroupType` is passed, leaving ID generation entirely to the caller or MongoDB.
- [ ] **Service `AddGroupType` is missing all input validation**: The Java service validates all 10 parameters (`Validator.notNull`, `Validator.noWhitespace`, `Validator.min`). The Go service performs no validation at all.
- [ ] **Service `AddGroupType` does not return the created entity**: Java returns `Mono<GroupType>` (the created object). Go returns `error` only, losing the created entity data that the controller needs to return in the response.
- [ ] **`EnsureDefaultGroupType` has hardcoded defaults missing most fields**: The Go `EnsureDefaultGroupType` only sets `ID`, `Name`, and `GroupSizeLimit` on the default type. The Java version sets all 10 fields including `InvitationStrategy` (`OWNER_MANAGER_MEMBER_REQUIRING_APPROVAL`), `JoinStrategy` (`INVITATION`), `GroupInfoUpdateStrategy` (`OWNER_MANAGER`), `MemberInfoUpdateStrategy` (`OWNER_MANAGER`), `GuestSpeakable` (`false`), `SelfInfoUpdatable` (`true`), `EnableReadReceipt` (`true`), `MessageEditable` (`true`).

## queryGroupTypes (non-paginated, `@GetMapping`)

- [ ] **Missing method entirely**: The Java controller has TWO overloaded `queryGroupTypes` methods — one non-paginated (`@GetMapping`, takes only `size`) and one paginated (`@GetMapping("page")`, takes `page` + `size`). In Go, the controller has `QueryGroupTypes(page, size *int)` and `QueryGroupTypesWithQuery(page *int, pageable any)`, but the naming and mapping are incorrect. The non-paginated Java version calls `groupTypeService.queryGroupTypes(0, size)` (page=0), while the paginated version calls `groupTypeService.queryGroupTypes(page, size)`. The Go controller does not distinguish these two endpoints properly.
- [ ] **Repository `FindGroupTypes` is missing pagination (skip/limit)**: The Java `findAll(page, size)` applies MongoDB skip/limit for pagination. The Go `FindGroupTypes` has a comment `// Add skip and limit if page/size exist... (Skipped detailed impl for brevity)` but **never actually applies skip or limit** — it always returns all documents.

## queryGroupTypes (paginated, `@GetMapping("page")`)

- [ ] **Controller method `QueryGroupTypesWithQuery` is incorrectly mapped**: The Java paginated endpoint calls `groupTypeService.countGroupTypes()` to get the total count and wraps results in `PaginationDTO`. The Go `QueryGroupTypesWithQuery` takes a `pageable any` parameter and returns `nil` — the pagination metadata (count + records) is completely unimplemented.
- [ ] **Service `QueryGroupTypes` does not return count information**: The Java service for the paginated endpoint returns both the total count and the page of results via `HttpHandlerResult.page(count, groupTypesFlux)`. The Go implementation only returns the records without any count, making pagination impossible.

## updateGroupType

- [ ] **Controller method `UpdateGroupType` is missing entirely**: The Java controller has `updateGroupType(Set<Long> ids, @RequestBody UpdateGroupTypeDTO updateGroupTypeDTO)` but the Go controller has **no `UpdateGroupType` method at all** on `GroupTypeController`.
- [ ] **Service `UpdateGroupTypes` is missing input validation**: Java validates `ids` is not empty, validates `name` has no whitespace, and validates `groupSizeLimit >= 1`. Go performs no validation.
- [ ] **Service `UpdateGroupTypes` is missing the "all-null fields" early return**: Java returns `ACKNOWLEDGED_UPDATE_RESULT` (a successful no-op `UpdateResult`) when all update fields are null, so `updateTypes` is never called. Go's `UpdateTypes` repository method does check `len(updateOps) == 0` and returns `nil`, but the service layer doesn't replicate the early return with an `UpdateResult`-like response — it just returns `nil` as error, losing the information about matched/modified count.
- [ ] **Service `UpdateGroupTypes` does not return `UpdateResult`**: Java returns `Mono<UpdateResult>` with `matchedCount`/`modifiedCount` that gets mapped to `UpdateResultDTO`. Go returns only `error`, discarding the MongoDB update result (matched count, modified count).

## deleteGroupType

- [ ] **Service `DeleteGroupTypes` is missing protection for default group type**: Java checks `if (groupTypeIds != null && groupTypeIds.contains(DEFAULT_GROUP_TYPE_ID))` and throws an `ILLEGAL_ARGUMENT` error preventing deletion of the default group type. The Go `DeleteGroupTypes` has **no such check** — it will blindly delete the default group type.
- [ ] **Service `DeleteGroupTypes` does not return `DeleteResult`**: Java returns `Mono<DeleteResult>` with `deletedCount` that gets mapped to `DeleteResultDTO`. Go returns only `error`, discarding the MongoDB delete result (deleted count).
- [ ] **Repository `DeleteTypes` behavior differs when ids is empty**: In Java, `deleteByIds` with null/empty ids likely deletes nothing or all. In Go, when `ids` is empty, the filter is `bson.M{}` which matches **ALL documents**, deleting every group type. The Java service passes `Set<Long> groupTypeIds` which is `@Nullable`, and when null, the Java repository's `deleteByIds` likely handles null differently than Go's empty-slice-means-no-filter logic.

# GroupStatisticsDTO.java
*Checked methods: GroupStatisticsDTO(Long deletedGroups, Long groupsThatSentMessages, Long createdGroups, List<StatisticsRecordDTO> deletedGroupsRecords, List<StatisticsRecordDTO> groupsThatSentMessagesRecords, List<StatisticsRecordDTO> createdGroupsRecords)*

Now I have all the information needed. Let me compare the two implementations.

**Java `GroupStatisticsDTO`:**
- `Long deletedGroups` — nullable Long
- `Long groupsThatSentMessages` — nullable Long
- `Long createdGroups` — nullable Long
- `List<StatisticsRecordDTO> deletedGroupsRecords` — typed list of `StatisticsRecordDTO`
- `List<StatisticsRecordDTO> groupsThatSentMessagesRecords` — typed list of `StatisticsRecordDTO`
- `List<StatisticsRecordDTO> createdGroupsRecords` — typed list of `StatisticsRecordDTO`

**Go `GroupStatisticsDTO`:**
- `DeletedGroups *int64` — pointer to int64 (nullable) ✓
- `GroupsThatSentMessages *int64` — pointer to int64 (nullable) ✓
- `CreatedGroups *int64` — pointer to int64 (nullable) ✓
- `DeletedGroupsRecords []any` — **uses `[]any` instead of `[]StatisticsRecordDTO`**
- `GroupsThatSentMessagesRecords []any` — **uses `[]any` instead of `[]StatisticsRecordDTO`**
- `CreatedGroupsRecords []any` — **uses `[]any` instead of `[]StatisticsRecordDTO`**

The Go codebase already has a proper `StatisticsRecordDTO` struct defined in `common_dtos.go` with the correct fields (`Date time.Time` and `Total int64`), but the `GroupStatisticsDTO` uses `[]any` instead of `[]StatisticsRecordDTO` for the three records fields.

## GroupStatisticsDTO

- [ ] `DeletedGroupsRecords`, `GroupsThatSentMessagesRecords`, and `CreatedGroupsRecords` fields are typed as `[]any` instead of `[]StatisticsRecordDTO`. The Go codebase already defines the correct `StatisticsRecordDTO` struct in `internal/domain/common/access/admin/dto/common_dtos.go`, so these fields should use `[]StatisticsRecordDTO` to match the Java type `List<StatisticsRecordDTO>`. Using `[]any` loses type safety and does not accurately reflect the Java source's typed list.

# GroupServiceController.java
*Checked methods: handleCreateGroupRequest(), handleDeleteGroupRequest(), handleQueryGroupsRequest(), handleQueryJoinedGroupIdsRequest(), handleQueryJoinedGroupsRequest(), handleUpdateGroupRequest(), handleCreateGroupBlockedUserRequest(), handleDeleteGroupBlockedUserRequest(), handleQueryGroupBlockedUserIdsRequest(), handleQueryGroupBlockedUsersInfosRequest(), handleCheckGroupQuestionAnswerRequest(), handleCreateGroupInvitationRequestRequest(), handleCreateGroupJoinRequestRequest(), handleCreateGroupQuestionsRequest(), handleDeleteGroupInvitationRequest(), handleUpdateGroupInvitationRequest(), handleDeleteGroupJoinRequestRequest(), handleUpdateGroupJoinRequestRequest(), handleDeleteGroupJoinQuestionsRequest(), handleQueryGroupInvitationsRequest(), handleQueryGroupJoinRequestsRequest(), handleQueryGroupJoinQuestionsRequest(), handleUpdateGroupJoinQuestionRequest(), handleCreateGroupMembersRequest(), handleDeleteGroupMembersRequest(), handleQueryGroupMembersRequest(), handleUpdateGroupMemberRequest()*

Now I have read both the Java source and both Go files completely. Let me compare each method systematically.

## handleCreateGroupRequest

**Java:** Calls `groupService.authAndCreateGroup(creatorIdAndOwnerId, creatorIdAndOwnerId, request.getName(), intro, announcement, minScore, typeId, muteEndDate, null, null, null)` - passes announcement, typeId, muteEndDate.

**Go:** Calls `c.groupService.CreateGroup(ctx, s.UserID, 0, &createReq.Name, createReq.Intro, createReq.MinScore)` - missing announcement, typeId, muteEndDate. Also passes `0` as typeId instead of from request.

## handleDeleteGroupRequest

**Java:** Calls `groupService.authAndDeleteGroup(notifyGroupMembersOfGroupDeleted, clientRequest.userId(), request.getGroupId())` and handles notification logic based on `notifyGroupMembersOfGroupDeleted`.

**Go:** Calls `c.groupService.DeleteGroup(ctx, s.UserID, deleteReq.GetGroupId())` - missing notification logic entirely.

## handleQueryGroupsRequest

**Java:** Passes `groupIds` (as Set), `name`, `lastUpdatedDate`, `skip`, `limit`, `fieldsToHighlight` to `groupService.authAndQueryGroups`. Returns data without NO_CONTENT check - always returns `GroupsWithVersion` even if empty.

**Go:** Passes `groupIds`, nil for name, `lastUpdatedDate`, nil for skip/limit/fieldsToHighlight. Missing `name`, `skip`, `limit`, `fieldsToHighlight` from request. Returns NO_CONTENT (204) when groups are empty, which differs from Java (Java always returns `GroupsWithVersion`).

## handleQueryJoinedGroupIdsRequest

**Java:** Calls `groupService.queryJoinedGroupIdsWithVersion(clientRequest.userId(), lastUpdatedDate)` - delegates to groupService, which handles version comparison internally.

**Go:** Calls `c.groupMemberService.QueryUserJoinedGroupIds(ctx, s.UserID)` - ignores `lastUpdatedDate` entirely, and calls a different service method. The version handling logic (commented out/empty) is not implemented.

## handleQueryJoinedGroupsRequest

**Java:** Calls `groupService.queryJoinedGroupsWithVersion(clientRequest.userId(), lastUpdatedDate)` - single call to groupService which handles everything.

**Go:** Makes two separate calls: first `c.groupMemberService.QueryUserJoinedGroupIds`, then `c.groupService.AuthAndQueryGroups` with those IDs. This is architecturally different from Java. Also returns NO_CONTENT when groups is empty, different from Java.

## handleUpdateGroupRequest

**Java:** Has two branches: if `successorId == null`, calls `groupService.authAndUpdateGroupInformation` with typeId, name, intro, announcement, minimumScore, muteEndDate, userDefinedAttributes. If `successorId != null`, calls `groupService.authAndTransferGroupOwnership`. Also has notification logic.

**Go:** Calls `c.groupService.AuthAndUpdateGroup` passing successorId along with other fields - does NOT branch between update vs transfer ownership. Missing muteEndDate, announcement, userDefinedAttributes handling. Missing notification logic.

## handleCreateGroupBlockedUserRequest

**Java:** Calls `groupBlocklistService.authAndBlockUser(clientRequest.userId(), groupId, userIdToBlock, null)` with complex notification logic.

**Go:** Calls `c.groupBlocklistService.BlockUser(ctx, createReq.GetGroupId(), createReq.GetUserId(), s.UserID)` - parameter order differs (groupId, userId, requesterId vs requesterId, groupId, userId). Missing notification logic.

## handleDeleteGroupBlockedUserRequest

**Java:** Calls `groupBlocklistService.unblockUser(clientRequest.userId(), groupId, userIdToUnblock, null, true)` and handles `wasBlocked` return + notification logic.

**Go:** Calls `c.groupBlocklistService.UnblockUser(ctx, deleteReq.GetGroupId(), deleteReq.GetUserId())` - different method signature (missing requesterId, auth check, and notification logic).

## handleQueryGroupBlockedUserIdsRequest

**Java:** Calls `groupBlocklistService.queryGroupBlockedUserIdsWithVersion(request.getGroupId(), lastUpdatedDate)` - no auth check (no userId parameter).

**Go:** Calls `c.groupBlocklistService.AuthAndQueryGroupBlockedUserIds(ctx, s.UserID, queryReq.GetGroupId(), lastUpdatedDate)` - adds auth check not present in Java.

## handleQueryGroupBlockedUsersInfosRequest

**Java:** Calls `groupBlocklistService.queryGroupBlockedUserInfosWithVersion(request.getGroupId(), lastUpdatedDate)` - no auth check. Returns the proto directly.

**Go:** Calls `c.groupBlocklistService.AuthAndQueryGroupBlockedUserInfos(ctx, s.UserID, queryReq.GetGroupId(), lastUpdatedDate)` - adds auth check. Manually constructs `UserInfo` protos with only IDs instead of using the full proto from the service.

## handleCheckGroupQuestionAnswerRequest

**Java:** Calls `groupQuestionService.authAndCheckGroupQuestionAnswerAndJoin` and builds a detailed `GroupJoinQuestionsAnswerResult` with `joined`, `questionIds`, `score`. Then has complex notification logic when joined=true (sends CreateGroupMembers notification to relevant users).

**Go:** Calls `c.groupQuestionService.CheckGroupJoinQuestionsAnswersAndJoin` and returns the result directly. Missing the detailed result construction (joined, questionIds, score fields). Missing all notification logic.

## handleCreateGroupInvitationRequestRequest

**Java:** Calls `groupInvitationService.authAndCreateGroupInvitation` with complex notification logic for group members/owner/managers/invitee.

**Go:** Calls `c.groupInvitationService.CreateInvitation` - missing `auth` prefix, missing notification logic.

## handleCreateGroupJoinRequestRequest

**Java:** Calls `groupJoinRequestService.authAndCreateGroupJoinRequest` with complex notification logic for group members/owner/managers.

**Go:** Calls `c.groupJoinRequestService.CreateJoinRequest` - missing `auth` prefix, missing notification logic.

## handleCreateGroupQuestionsRequest

**Java:** Creates a list of `NewGroupQuestion` objects, calls `groupQuestionService.authAndCreateGroupJoinQuestions(clientRequest.userId(), request.getGroupId(), questions)` - single batch call. Returns `RequestHandlerResult.ofDataLongs(questionIds)`.

**Go:** Loops through questions individually calling `c.groupQuestionService.CreateJoinQuestion` one at a time. Missing `auth` check. Not atomic - if one fails partway, some questions are already created.

## handleDeleteGroupInvitationRequest

**Java:** Calls `groupInvitationService.authAndRecallPendingGroupInvitation` with complex notification logic.

**Go:** Calls `c.groupInvitationService.RecallPendingInvitation` - missing `auth` prefix. Missing notification logic. The method name is different (RecallPendingInvitation vs authAndRecallPendingGroupInvitation).

## handleUpdateGroupInvitationRequest

**Java:** Calls `groupInvitationService.authAndHandleInvitation(clientRequest.userId(), request.getInvitationId(), request.getResponseAction(), request.getReason())` - passes ResponseAction enum and reason string. Complex notification logic.

**Go:** Converts ResponseAction to boolean accept, calls `c.groupInvitationService.ReplyToInvitation` - loses the reason parameter. Missing notification logic.

## handleDeleteGroupJoinRequestRequest

**Java:** Calls `groupJoinRequestService.authAndRecallPendingGroupJoinRequest` with complex notification logic.

**Go:** Calls `c.groupJoinRequestService.RecallPendingJoinRequest` - missing `auth`. Missing notification logic.

## handleUpdateGroupJoinRequestRequest

**Java:** Calls `groupJoinRequestService.authAndHandleJoinRequest` with ResponseAction and reason. Complex notification logic.

**Go:** Converts ResponseAction to boolean, calls `c.groupJoinRequestService.ReplyToJoinRequest` - loses reason parameter. Missing notification logic.

## handleDeleteGroupJoinQuestionsRequest

**Java:** Calls `groupQuestionService.authAndDeleteGroupJoinQuestions(clientRequest.userId(), request.getGroupId(), CollectionUtil.newSet(request.getQuestionIdsList()))` - single batch call with auth and groupId.

**Go:** Loops through question IDs individually calling `c.groupQuestionService.DeleteJoinQuestion` one at a time. Missing auth. Missing groupId parameter. Not atomic.

## handleQueryGroupInvitationsRequest

**Java:** Two branches: if groupId==null, calls `queryUserGroupInvitationsWithVersion`; otherwise calls `authAndQueryGroupInvitationsWithVersion`. Both return proto directly.

**Go:** Same branching logic but calls different method names. Manually reconstructs invitation protos from POJOs instead of using proto directly. The Java version delegates proto conversion to the service; Go does it in the controller.

## handleQueryGroupJoinRequestsRequest

**Java:** Single call to `groupJoinRequestService.authAndQueryGroupJoinRequestsWithVersion` with groupId (nullable) and lastUpdatedDate.

**Go:** Two branches: if groupId!=nil calls `AuthAndQueryGroupJoinRequestsWithVersion`, else calls `QueryUserGroupJoinRequestsWithVersion`. Java has a single method that handles both cases internally.

## handleQueryGroupJoinQuestionsRequest

**Java:** Calls `groupQuestionService.authAndQueryGroupJoinQuestionsWithVersion(clientRequest.userId(), request.getGroupId(), request.getWithAnswers(), lastUpdatedDate)` and returns the proto directly.

**Go:** Calls similar method but manually constructs proto `GroupJoinQuestion` objects from the POJO. The `LastUpdatedDate` in the version is set as `*int64` from `questionsWithVersion.LastUpdatedDate` but this may not match the proto format.

## handleUpdateGroupJoinQuestionRequest

**Java:** Calls `groupQuestionService.authAndUpdateGroupJoinQuestion(clientRequest.userId(), request.getQuestionId(), question, answers, score)`.

**Go:** Calls `c.groupQuestionService.UpdateJoinQuestion(ctx, updateReq.GetQuestionId(), 0, updateReq.Question, updateReq.Answers, score)` - passing `0` for groupId instead of the user ID (missing auth). Missing auth check.

## handleCreateGroupMembersRequest

**Java:** Passes `name` from request, `muteEndDate`, `userIds`, `groupId`, `role` (nullable), `name`, `muteEndDate`, and null. Complex notification logic.

**Go:** Missing `name` parameter from request. Calls `c.groupMemberService.AuthAndAddGroupMembers` without `name`. Missing notification logic. Returns LongsWithVersion (IDs) while Java returns RequestHandlerResult with notification recipients.

## handleDeleteGroupMembersRequest

**Java:** Extracts successorId, quitAfterTransfer, memberIdsToDelete. Calls `groupMemberService.authAndDeleteGroupMembers`. Handles empty deletedUserIds case and notification logic.

**Go:** Calls `c.groupMemberService.AuthAndDeleteGroupMembers` with similar parameters. Missing the empty result check and notification logic.

## handleQueryGroupMembersRequest

**Java:** Has two branches based on `request.getMemberIdsCount() > 0`: calls `authAndQueryGroupMembers` (for specific members) or `authAndQueryGroupMembersWithVersion` (for all with version). Also handles `withStatus` parameter.

**Go:** Single call to `c.groupMemberService.AuthAndQueryGroupMembersWithVersion` - missing the branch for specific memberIds query. Missing `withStatus` parameter.

## handleUpdateGroupMemberRequest

**Java:** Calls `groupMemberService.authAndUpdateGroupMember` with name, role, muteEndDate. Then has notification logic.

**Go:** Calls `c.groupMemberService.AuthAndUpdateGroupMember` with similar params. Missing notification logic.

Here is the consolidated bug report:

## handleCreateGroupRequest
- [ ] Missing `announcement` field from request (Java passes `request.getAnnouncement()`)
- [ ] Missing `typeId` field from request (Java passes `request.getTypeId()`); instead hardcoded `0`
- [ ] Missing `muteEndDate` field from request (Java passes `new Date(request.getMuteEndDate())`)

## handleDeleteGroupRequest
- [ ] Missing notification logic: Java conditionally notifies group members or requester sessions based on `notifyGroupMembersOfGroupDeleted`/`notifyRequesterOtherOnlineSessionsOfGroupDeleted` flags

## handleQueryGroupsRequest
- [ ] Missing `name` query parameter (Java passes `request.getName()`)
- [ ] Missing `skip` parameter (Java passes `request.getSkip()`)
- [ ] Missing `limit` parameter (Java passes `request.getLimit()`)
- [ ] Missing `fieldsToHighlight` parameter (Java passes `request.getFieldsToHighlightList()`)
- [ ] Returns NO_CONTENT (code 204) when groups is empty, but Java always returns `GroupsWithVersion` (even with empty list) — behavioral difference

## handleQueryJoinedGroupIdsRequest
- [ ] Ignores `lastUpdatedDate` parameter — calls `QueryUserJoinedGroupIds` without it, while Java passes it to `groupService.queryJoinedGroupIdsWithVersion` for version-based filtering
- [ ] Calls a different service (`groupMemberService.QueryUserJoinedGroupIds`) than Java (`groupService.queryJoinedGroupIdsWithVersion`)

## handleQueryJoinedGroupsRequest
- [ ] Uses two-step query (fetch IDs then fetch groups) instead of Java's single call to `groupService.queryJoinedGroupsWithVersion` — different architecture that may produce inconsistent results
- [ ] Returns NO_CONTENT when groups is empty, while Java always returns `GroupsWithVersion`

## handleUpdateGroupRequest
- [ ] Missing branching logic: Java has two distinct paths — `authAndUpdateGroupInformation` when no successorId, vs `authAndTransferGroupOwnership` when successorId is present. Go passes successorId to a single method.
- [ ] Missing `muteEndDate` parameter (Java passes `new Date(request.getMuteEndDate())`)
- [ ] Missing `announcement` field from request
- [ ] Missing `userDefinedAttributes` parameter (Java passes `request.getUserDefinedAttributesMap()`)
- [ ] Missing notification logic for group members/requester sessions

## handleCreateGroupBlockedUserRequest
- [ ] Missing notification logic for blocked user, group members, and requester sessions

## handleDeleteGroupBlockedUserRequest
- [ ] Missing auth parameter — Java calls `unblockUser(clientRequest.userId(), groupId, userIdToUnblock, null, true)` with requesterId, Go calls `UnblockUser(groupId, userIdToUnblock)` without requesterId
- [ ] Missing `wasBlocked` return handling — Java returns `RequestHandlerResult.OK` if not blocked, Go always returns success
- [ ] Missing notification logic

## handleQueryGroupBlockedUserIdsRequest
- [ ] Adds auth check not present in Java — Java calls `queryGroupBlockedUserIdsWithVersion(groupId, lastUpdatedDate)` without userId, Go calls `AuthAndQueryGroupBlockedUserIds` with userId

## handleQueryGroupBlockedUsersInfosRequest
- [ ] Adds auth check not present in Java — Java calls `queryGroupBlockedUserInfosWithVersion` without userId
- [ ] Manually constructs `UserInfo` protos with only IDs, while Java delegates proto conversion to the service which returns complete user info

## handleCheckGroupQuestionAnswerRequest
- [ ] Missing detailed result construction — Java builds `GroupJoinQuestionsAnswerResult` with `joined`, `questionIds`, `score` fields; Go returns a result object directly
- [ ] Missing notification logic when `joined=true` — Java sends `CreateGroupMembers` notification to relevant group members

## handleCreateGroupInvitationRequestRequest
- [ ] Missing notification logic — Java has complex notification to group members/owner/managers/invitee based on config flags

## handleCreateGroupJoinRequestRequest
- [ ] Missing notification logic — Java has complex notification to group members/owner/managers based on config flags

## handleCreateGroupQuestionsRequest
- [ ] Creates questions one-by-one in a loop instead of batch — not atomic; partial creation possible on error
- [ ] Missing auth check — Java passes `clientRequest.userId()` for authorization, Go does not

## handleDeleteGroupInvitationRequest
- [ ] Missing notification logic — Java has complex notification to group members/invitee/owner/managers

## handleUpdateGroupInvitationRequest
- [ ] Missing `reason` parameter — Java passes `request.getReason()` to service
- [ ] Converts `ResponseAction` to boolean `accept` — loses information vs Java's enum-based handling
- [ ] Missing notification logic — Java has complex notification including member-added notifications when requester joins

## handleDeleteGroupJoinRequestRequest
- [ ] Missing notification logic — Java has complex notification to group members/owner/managers

## handleUpdateGroupJoinRequestRequest
- [ ] Missing `reason` parameter — Java passes `request.getReason()` to service
- [ ] Converts `ResponseAction` to boolean — loses information vs Java's enum-based handling
- [ ] Missing notification logic — Java has complex notification including member-added notifications when requester joins

## handleDeleteGroupJoinQuestionsRequest
- [ ] Deletes questions one-by-one in a loop instead of batch — not atomic
- [ ] Missing auth check — Java passes `clientRequest.userId()` for authorization
- [ ] Missing `groupId` parameter — Java passes `request.getGroupId()` for auth validation

## handleQueryGroupInvitationsRequest
- [ ] Manually reconstructs invitation protos from POJOs instead of using service-returned proto, risking field mismatches

## handleQueryGroupJoinRequestsRequest
- [ ] Two-branch logic differs from Java — Java has a single `authAndQueryGroupJoinRequestsWithVersion` method that handles both nullable groupId and non-null cases internally

## handleUpdateGroupJoinQuestionRequest
- [ ] Missing auth check — Java passes `clientRequest.userId()`, Go passes `0` as groupId (likely placeholder for missing userId)
- [ ] `answers` passed directly as `updateReq.Answers` without converting empty to nil — Java converts empty list to null

## handleCreateGroupMembersRequest
- [ ] Missing `name` parameter from request — Java passes `request.getName()`
- [ ] Returns IDs as `LongsWithVersion` in data, but Java's notification-based return includes member notification recipients, not just the member IDs
- [ ] Missing notification logic

## handleDeleteGroupMembersRequest
- [ ] Missing empty result check — Java returns `RequestHandlerResult.OK` when `deletedUserIds.isEmpty()` (no notification)
- [ ] Missing notification logic

## handleQueryGroupMembersRequest
- [ ] Missing branch for specific memberIds query — Java has two paths: `authAndQueryGroupMembers` (when memberIds present) vs `authAndQueryGroupMembersWithVersion` (when querying all with version). Go only has the version path.
- [ ] Missing `withStatus` parameter — Java passes `request.getWithStatus()`

## handleUpdateGroupMemberRequest
- [ ] Missing notification logic — Java notifies group members/updated member based on config flags

# NewGroupQuestion.java
*Checked methods: NewGroupQuestion(String question, LinkedHashSet<String> answers, Integer score)*

## NewGroupQuestion

- [ ] **`validNewGroupQuestion` is an empty stub** — The Go version at `validator.go:119-120` is a no-op function with no parameters and no body logic. The Java version (`DataValidator.validNewGroupQuestion`) takes a `NewGroupQuestion` parameter and validates: (1) `answers` must not be empty (throws `EMPTY_GROUP_QUESTION_ANSWERS`), and (2) `score` must not be null and must be >= 0 (throws `ILLEGAL_GROUP_QUESTION_SCORE`). The Go implementation performs zero validation.
- [ ] **Missing all field-level validation from the service layer** — The Java `GroupQuestionService.createGroupJoinQuestions` method performs additional validation for each `NewGroupQuestion`: `notNull(question)`, `maxLength(question, questionContentLimit)`, `notNull(answers)`, `inSizeRange(answers, 1, maxAnswerCount)`, `maxLength(answers, answerContentLimit)`, `notNull(score)`, `min(score, 0)`. There is no equivalent validation logic implemented in the Go code for these fields.
- [ ] **`LinkedHashSet` semantics not preserved** — The Java `NewGroupQuestion` uses `LinkedHashSet<String>` for `answers`, which guarantees insertion order and uniqueness (no duplicate answers). The Go struct uses `[]string` (a plain slice), which allows duplicate entries and does not enforce uniqueness. This is a behavioral difference that could lead to data integrity issues.

# GroupInvitationRepository.java
*Checked methods: getEntityExpireAfterSeconds(), updateStatusIfPending(Long invitationId, RequestStatus requestStatus, @Nullable String reason, @Nullable ClientSession session), updateInvitations(Set<Long> invitationIds, @Nullable Long inviterId, @Nullable Long inviteeId, @Nullable String content, @Nullable RequestStatus status, @Nullable Date creationDate, @Nullable Date responseDate), count(@Nullable Set<Long> ids, @Nullable Set<Long> groupIds, @Nullable Set<Long> inviterIds, @Nullable Set<Long> inviteeIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange), findGroupIdAndInviteeIdAndStatus(Long invitationId), findGroupIdAndInviterIdAndInviteeIdAndStatus(Long invitationId), findInvitationsByInviteeId(Long inviteeId), findInvitationsByInviterId(Long inviterId), findInvitationsByGroupId(Long groupId), findInviteeIdAndGroupIdAndCreationDateAndStatus(Long invitationId), findInvitations(@Nullable Set<Long> ids, @Nullable Set<Long> groupIds, @Nullable Set<Long> inviterIds, @Nullable Set<Long> inviteeIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange, @Nullable Integer page, @Nullable Integer size)*

Now I have all the context needed. Let me analyze each method carefully.

## getEntityExpireAfterSeconds()

**Java:** Returns `expireAfterSeconds` which is set from `properties.getService().getGroup().getInvitation().getExpireAfterSeconds()`.

**Go:** The `group_invitation_repository.go` does **NOT** have a `GetEntityExpireAfterSeconds()` method at all. There is no equivalent of the expiration property management.

## updateStatusIfPending()

**Java:** 
1. Filter includes `isNotExpired(UserFriendRequest.Fields.CREATION_DATE, getEntityExpirationDate())` — checks that the invitation hasn't expired.
2. Sets `RESPONSE_DATE` to `new Date()` (always, via `setIfNotNull` which will never be null since `new Date()` is always non-null).

**Go:**
1. Filter does NOT check for expiration (`isNotExpired` equivalent is missing).
2. `responseDate` is passed as a parameter instead of being set to current time inside the method.

## updateInvitations()

**Java:** Uses `updateResponseDateBasedOnStatus` which:
- If `status` is ACCEPTED/DECLINED/IGNORED: sets `responseDate` (defaults to `new Date()` if null).
- If `status` is not processed by responder (e.g., PENDING): **unsets** the `responseDate` field.

**Go:** Simply sets `responseDate` if non-nil, without any status-based logic. Missing the `updateResponseDateBasedOnStatus` logic entirely.

## count()

**Java:**
1. Takes `Set<Long>` for ids, groupIds, inviterIds, inviteeIds, statuses — allows filtering by multiple values per field using `inIfNotNull`.
2. Takes `DateRange` for creationDateRange, responseDateRange, expirationDateRange.
3. Calls `getCreationDateRange(creationDateRange, expirationDateRange)` which merges expiration date range into creation date range.
4. Calls `isExpiredOrNot(statuses, ...)` which conditionally adds expiration filtering.

**Go:** `CountInvitations` takes single `*int64` values (not sets), single `*po.RequestStatus` (not a set), and a `lastUpdatedDate` (not date ranges). Completely different signature and logic — no `getCreationDateRange`, no `isExpiredOrNot`, no multi-value filtering.

## findGroupIdAndInviteeIdAndStatus()

**Java:** Includes `CREATION_DATE` in projection for `findExpirableDoc` (which calls `transformExpiredRequest` to check expiration status and change PENDING to EXPIRED if expired).

**Go:** Does NOT include `cd` (creation date) in the projection, so it cannot perform expiration checking. Also does not call any equivalent of `transformExpiredRequest`.

## findGroupIdAndInviterIdAndInviteeIdAndStatus()

**Java:** Includes `CREATION_DATE` in projection for `findExpirableDoc`.

**Go:** Does NOT include `cd` (creation date) in projection. Missing expiration transform logic.

## findInvitationsByInviteeId()

**Java:** Calls `findExpirableDocs(filter)` which applies `transformExpiredRequest` to transform expired PENDING records to EXPIRED status.

**Go:** Does NOT apply any expiration transform. Also adds sorting by `{"cd": -1}` which is not in the Java version.

## findInvitationsByInviterId()

**Java:** Returns `findExpirableDocs(filter)`.

**Go:** Stub implementation `return nil, nil` — completely unimplemented.

## findInvitationsByGroupId()

**Java:** Calls `findExpirableDocs(filter)` which applies expiration transform.

**Go:** Does NOT apply any expiration transform. Also adds sorting by `{"cd": -1}` not in Java.

## findInviteeIdAndGroupIdAndCreationDateAndStatus()

**Java:** Calls `mongoClient.findOne(entityClass, filter, options)` directly — does NOT call `findExpirableDoc`, so no expiration transform. The projection includes INVITEE_ID, GROUP_ID, CREATION_DATE, STATUS.

**Go:** Equivalent — calls `FindOne` with the same projection fields. This one appears correct.

## findInvitations()

**Java:**
1. Takes sets for ids, groupIds, inviterIds, inviteeIds, statuses — multi-value filtering.
2. Takes DateRange objects for date filtering.
3. Uses `getCreationDateRange` and `isExpiredOrNot`.
4. Uses `paginateIfNotNull` for pagination.
5. Calls `findExpirableDocs(filter, options)` for expiration transform.

**Go:** Takes single values instead of sets. Uses a `lastUpdatedDate` instead of date ranges. No `getCreationDateRange`, no `isExpiredOrNot`, no expiration transform.

---

Here is the complete bug report:

## getEntityExpireAfterSeconds
- [ ] Missing entirely. The Go repository has no `GetEntityExpireAfterSeconds()` method, no `expireAfterSeconds` field, and no properties-based initialization. The entire expiration configuration mechanism is absent.

## updateStatusIfPending
- [ ] Missing expiration check in filter. The Java version filters with `isNotExpired(CREATION_DATE, getEntityExpirationDate())` to ensure only non-expired pending invitations are updated. The Go version omits this filter condition entirely.
- [ ] `responseDate` is set to the caller-provided value instead of being unconditionally set to the current time (`new Date()` in Java). In Java, `setIfNotNull(GroupInvitation.Fields.RESPONSE_DATE, new Date())` always sets the response date to now because `new Date()` is never null.

## updateInvitations
- [ ] Missing `updateResponseDateBasedOnStatus` logic. Java conditionally handles the response date based on status: if status is ACCEPTED/DECLINED/IGNORED it sets responseDate (defaulting to `new Date()` if null); if status is PENDING or other non-responder-processed statuses, it **unsets** the responseDate field. The Go version simply sets responseDate if non-nil with no status-based conditional logic.

## count
- [ ] Takes single `*int64` filter parameters instead of `[]int64` (sets) for groupID, inviterID, inviteeID. The Java version uses `inIfNotNull` with `Set<Long>` to filter by multiple values per field. The Go version can only filter by a single value per field.
- [ ] Takes a single `*po.RequestStatus` instead of `[]po.RequestStatus` (set), losing the ability to filter by multiple statuses.
- [ ] Takes a single `*time.Time` (`lastUpdatedDate`) instead of separate `creationDateRange`, `responseDateRange`, and `expirationDateRange` `DateRange` parameters. The Java version supports filtering by creation date range, response date range, and expiration date range independently.
- [ ] Missing `getCreationDateRange` logic that merges the expiration date range into the creation date range.
- [ ] Missing `isExpiredOrNot` logic that conditionally adds expiration filtering based on whether the statuses set contains EXPIRED and/or PENDING.
- [ ] Missing `ids` (`_id`) filtering parameter entirely.

## findGroupIdAndInviteeIdAndStatus
- [ ] Missing `cd` (creation date) in the projection. The Java version explicitly includes `CREATION_DATE` in the projection because it is "Required by findExpirableDoc" to check expiration.
- [ ] Missing `transformExpiredRequest` logic (from `findExpirableDoc`). The Java version checks if the invitation's PENDING status should be transformed to EXPIRED based on the creation date and `expireAfterSeconds`. The Go version returns the raw database status without this transformation.

## findGroupIdAndInviterIdAndInviteeIdAndStatus
- [ ] Missing `cd` (creation date) in the projection. The Java version includes `CREATION_DATE` because it is "Required by findExpirableDoc".
- [ ] Missing `transformExpiredRequest` logic from `findExpirableDoc`. The Java version transforms PENDING to EXPIRED when the record is expired; the Go version returns raw data.

## findInvitationsByInviteeId
- [ ] Missing expiration transform logic. The Java version calls `findExpirableDocs(filter)` which applies `transformExpiredRequest` to convert expired PENDING records to EXPIRED status. The Go version returns raw database records.
- [ ] Adds sorting `{"cd": -1}` that is not present in the Java version.

## findInvitationsByInviterId
- [ ] Completely unimplemented — returns `nil, nil` stub. The Java version filters by `INVITER_ID` and calls `findExpirableDocs(filter)`.

## findInvitationsByGroupId
- [ ] Missing expiration transform logic. The Java version calls `findExpirableDocs(filter)` which applies `transformExpiredRequest`. The Go version returns raw database records.
- [ ] Adds sorting `{"cd": -1}` that is not present in the Java version.

## findInvitations
- [ ] Takes single `*int64` filter parameters instead of `[]int64` (sets) for groupID, inviterID, inviteeID. The Java version uses `inIfNotNull` with `Set<Long>` to filter by multiple values. The Go version only supports single-value filtering.
- [ ] Takes a single `*po.RequestStatus` instead of `[]po.RequestStatus` (set), losing multi-status filtering capability.
- [ ] Takes a single `*time.Time` (`lastUpdatedDate`) instead of separate `creationDateRange`, `responseDateRange`, and `expirationDateRange` `DateRange` parameters.
- [ ] Missing `getCreationDateRange` logic that merges expiration date range into creation date range.
- [ ] Missing `isExpiredOrNot` logic for conditional expiration filtering.
- [ ] Missing `ids` (`_id`) filtering parameter entirely.
- [ ] Missing expiration transform logic from `findExpirableDocs`. The Java version transforms PENDING to EXPIRED for expired records.
- [ ] Always applies pagination (`SetSkip`/`SetLimit`) even when page/size are zero, unlike Java's `paginateIfNotNull` which only applies when page and size are non-null.

# GroupJoinRequestRepository.java
*Checked methods: getEntityExpireAfterSeconds(), updateStatusIfPending(Long requestId, RequestStatus status, Long responderId, @Nullable String reason, @Nullable ClientSession session), updateRequests(Set<Long> requestIds, @Nullable Long requesterId, @Nullable Long responderId, @Nullable String content, @Nullable RequestStatus status, @Nullable Date creationDate, @Nullable Date responseDate), countRequests(@Nullable Set<Long> ids, @Nullable Set<Long> groupIds, @Nullable Set<Long> requesterIds, @Nullable Set<Long> responderIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange), findGroupId(Long requestId), findRequesterIdAndStatusAndGroupId(Long requestId), findRequestsByGroupId(Long groupId), findRequestsByRequesterId(Long requesterId), findRequests(@Nullable Set<Long> ids, @Nullable Set<Long> groupIds, @Nullable Set<Long> requesterIds, @Nullable Set<Long> responderIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange, @Nullable Integer page, @Nullable Integer size)*

Now I have a comprehensive understanding of both implementations. Let me compile the review.

## getEntityExpireAfterSeconds()

- [ ] **Go implementation is missing entirely.** The Java version reads `expireAfterSeconds` from `TurmsProperties` via a properties change listener and returns it in `getEntityExpireAfterSeconds()`. The Go repository has no `expireAfterSeconds` field, no property-reading logic, and no equivalent method. This breaks all expiration-dependent behavior (expiration date calculation, `isExpiredOrNot` filtering, `findExpirableDocs` transformation, etc.).

## updateStatusIfPending()

- [ ] **Go version takes `responseDate` as a parameter instead of generating it internally.** The Java version sets `RESPONSE_DATE` to `new Date()` internally (line 84), while the Go version requires `responseDate time.Time` as a caller-provided parameter (line 68). If the caller passes a zero-value or incorrect time, the behavior diverges from Java.

## updateRequests()

- [ ] **Missing `updateResponseDateBasedOnStatus` logic.** The Java version calls `updateResponseDateBasedOnStatus` (lines 105-108), which conditionally sets or *unsets* the response date based on whether the status is "processed by responder." The Go version simply sets `responseDate` if non-nil (line 212) and never unsets it. This means if a status change requires the response date to be cleared (e.g., reverting to a non-responder-processed status), the Go code will not clear it, producing incorrect data.

## countRequests()

- [ ] **Missing `isExpiredOrNot` filter logic.** The Java version calls `.isExpiredOrNot(statuses, GroupJoinRequest.Fields.CREATION_DATE, expirationDate)` (line 131), which adds a creation-date-based expiration filter: if `EXPIRED` is in statuses (without `PENDING`), it adds `cd < expirationDate`; if `PENDING` is in statuses (without `EXPIRED`), it adds `cd >= expirationDate OR cd IS NULL`; if both or neither, no additional filter. The Go `CountRequests` has no equivalent logic — it only applies the `expirationDateRange` as a direct filter on `ed` (line 251), which is a completely different field and semantics.
- [ ] **Missing `getCreationDateRange` merging logic.** The Java version merges `creationDateRange` with `expirationDateRange` via `getCreationDateRange(creationDateRange, expirationDateRange)` (line 129), which adjusts the expiration range by `expireAfterSeconds` and intersects it with the creation date range. The Go version treats `creationDateRange` and `expirationDateRange` as independent filters on separate fields (`cd` and `ed`), missing the cross-field range calculation.
- [ ] **Missing `getEntityExpirationDate()` usage.** The Java version calls `getEntityExpirationDate()` (line 121) and passes the result to `isExpiredOrNot`. The Go version has no equivalent — there is no computed expiration date used for filtering.

## findGroupId()

- [ ] No bugs found. The Go implementation correctly projects `gid` and returns it.

## findRequesterIdAndStatusAndGroupId()

- [ ] **Missing `creationDate` in projection and missing expiration transformation.** The Java version includes `CREATION_DATE` in the projection (line 152) with the comment "Required by findExpirableDoc," and then passes the result through `findExpirableDoc` which calls `transformExpiredRequest`. This means if the request status is `PENDING` and the creation date indicates it has expired, the status is overridden to `EXPIRED`. The Go version (line 276) does not include `cd` in the projection and does not perform any expiration transformation, so expired pending requests will incorrectly retain `PENDING` status.

## findRequestsByGroupId()

- [ ] **Missing expiration transformation (`findExpirableDocs` behavior).** The Java version calls `findExpirableDocs(filter)` (line 164), which transforms returned documents — any request with status `PENDING` that has expired has its status changed to `EXPIRED`. The Go version (line 93) returns raw documents with no such transformation, so expired pending requests will be returned with the incorrect `PENDING` status.

## findRequestsByRequesterId()

- [ ] **Missing expiration transformation (`findExpirableDocs` behavior).** Same issue as `findRequestsByGroupId`. The Java version calls `findExpirableDocs(filter)` (line 170) to transform expired pending requests to `EXPIRED` status. The Go version (line 110) returns raw documents without this transformation.

## findRequests()

- [ ] **Missing `isExpiredOrNot` filter logic.** Same as `countRequests`. The Java version calls `.isExpiredOrNot(statuses, ..., expirationDate)` (line 194), which adds creation-date-based expiration filtering. The Go version has no equivalent.
- [ ] **Missing `getCreationDateRange` merging logic.** Same as `countRequests`. The Java version merges creation and expiration date ranges (line 192). The Go version treats them as independent filters on different fields.
- [ ] **Missing `getEntityExpirationDate()` usage.** Same as `countRequests`. The Java version computes and uses an expiration date (line 184).
- [ ] **Missing expiration transformation (`findExpirableDocs` behavior).** The Java version calls `findExpirableDocs(filter, options)` (line 197) which transforms expired pending requests to `EXPIRED` status. The Go version returns raw documents without this transformation.
- [ ] **Method signature mismatch: Java takes `Set<Long>` for IDs/groups/requesters/responders, Go takes `*int64` single values.** The Java `findRequests` accepts `Set<Long>` for `ids`, `groupIds`, `requesterIds`, `responderIds` (allowing multi-value filters with `$in`), while the Go `FindRequests` accepts `*int64` for `groupID`, `requesterID`, `responderID` (single value equality match). This prevents querying by multiple groups/requesters/responders at once.
- [ ] **Method signature mismatch: Java takes `Set<RequestStatus>`, Go takes `*RequestStatus`.** The Java version accepts a set of statuses (enabling `$in` queries for multiple statuses), while the Go version accepts a single `*po.RequestStatus` (equality match). This prevents filtering by multiple statuses simultaneously.
- [ ] **Date range semantics completely different.** The Java version takes `DateRange` objects (with start/end) for `creationDateRange`, `responseDateRange`, `expirationDateRange`, which support `$gte`/`$lte` range queries. The Go version takes single `*time.Time` pointers and uses `$gte` for creation/response and `$lt` for expiration, which is not equivalent to range queries.
- [ ] **Always applies pagination even when page/size are zero.** The Java version uses `paginateIfNotNull(page, size)` (line 196), which skips pagination when page or size is null. The Go version always sets `skip = page * size` and `limit = size` (lines 169-170), meaning when `page=0` and `size=0`, it sets `limit=0` which returns no results.

# GroupMemberRepository.java
*Checked methods: deleteAllGroupMembers(@Nullable Set<Long> groupIds, @Nullable ClientSession session), updateGroupMembers(Set<GroupMember.Key> keys, @Nullable String name, @Nullable GroupMemberRole role, @Nullable Date joinDate, @Nullable Date muteEndDate, @Nullable ClientSession session), countMembers(Long groupId), countMembers(@Nullable Set<Long> groupIds, @Nullable Set<Long> userIds, @Nullable Set<GroupMemberRole> roles, @Nullable DateRange joinDateRange, @Nullable DateRange muteEndDateRange), findGroupManagersAndOwnerId(Long groupId), findGroupMembers(Long groupId), findGroupMembers(Long groupId, Set<Long> memberIds), findGroupsMembers(@Nullable Set<Long> groupIds, @Nullable Set<Long> userIds, @Nullable Set<@ValidGroupMemberRole GroupMemberRole> roles, @Nullable DateRange joinDateRange, @Nullable DateRange muteEndDateRange, @Nullable Integer page, @Nullable Integer size), findGroupMemberIds(Long groupId), findGroupMemberIds(Set<Long> groupIds), findGroupMemberKeyAndRoleParis(Set<Long> userIds, Long groupId), findGroupMemberRole(Long userId, Long groupId), findMemberIdsByGroupId(Long groupId), findUserJoinedGroupIds(Long userId), findUsersJoinedGroupIds(@Nullable Set<Long> groupIds, @NotEmpty Set<Long> userIds, @Nullable Integer page, @Nullable Integer size), isMemberMuted(Long groupId, Long userId)*

Now I have all the context needed for a thorough comparison. Let me analyze each method systematically.

## deleteAllGroupMembers

Java: Filters by `inIfNotNull(GroupMember.Fields.ID_GROUP_ID, groupIds)` — if `groupIds` is null, the filter is empty (deletes all documents). Accepts a nullable ClientSession.

Go (repository `DeleteByGroupIDs`): Returns early if `len(groupIDs) == 0`, which means an empty/nil slice will NOT delete anything. The Java version with null groupIds would delete ALL members.

Go (service `DeleteAllGroupMembers`): Also returns early on empty slice. Additionally, the Java repository version takes a nullable `ClientSession`, while the Go version doesn't pass the session to the repo call.

## updateGroupMembers

Java: Uses `Filter.newBuilder(1).in(DomainFieldName.ID, keys)` — filters by the composite `_id` field. The update uses `setIfNotNull` for name, role, joinDate. For muteEndDate, compares against `System.currentTimeMillis()`.

Go: Uses `"_id": bson.M{"$in": keys}` — this is correct. Uses `set`/`unset` logic for muteEndDate. However, Go uses `muteEndDate.Before(time.Now())` which is semantically the same as Java's `muteEndDate.getTime() < System.currentTimeMillis()`.

The Go version also has an early return when `len(update) == 0` — the Java version doesn't do this, it would still make the updateMany call (which is a no-op but still hits the database).

## countMembers(Long groupId)

Java: `Filter.newBuilder(1).eq(GroupMember.Fields.ID_GROUP_ID, groupId)` — filters on `_id.gid`.

Go: `bson.M{"_id.gid": groupID}` — correct.

## countMembers(Set<Long> groupIds, Set<Long> userIds, ...)

Java: Multi-parameter count with optional filters. This method exists in the Java repo.

Go: **Missing entirely.** There is no `CountMembersWithFilters` or equivalent method in the Go repository.

## findGroupManagersAndOwnerId

Java: Filters by groupId + role IN (MANAGER, OWNER).

Go: Same logic, correct.

## findGroupMembers(Long groupId)

Java: Filters by `_id.gid`.

Go: Same, correct.

## findGroupMembers(Long groupId, Set<Long> memberIds)

Java: Filters by `_id.gid` AND `_id.uid` IN memberIds.

Go (`FindGroupMembersWithIds`): Same, correct.

## findGroupsMembers

Java: Full implementation with filter on groupIds, userIds, roles, joinDateRange, muteEndDateRange, plus pagination.

Go: **Stub returning nil, nil.** Completely unimplemented.

## findGroupMemberIds(Long groupId)

Java: Filters by `_id.gid`, projects `_id.uid`, returns user IDs.

Go (`FindGroupMemberIDs`): Same, but only handles a single groupId. Correct for the single-group overload.

## findGroupMemberIds(Set<Long> groupIds)

Java: Filters by `_id.gid` IN groupIds, projects `_id.uid`, returns user IDs.

Go: **Missing.** The `FindGroupMemberIDs` method only accepts a single groupID. The `FindMemberIdsByGroupIds` exists but is mapped to a different Java method.

## findGroupMemberKeyAndRoleParis

Java: Filters by `_id.uid` IN userIds AND `_id.gid` = groupId. Projects `role`.

Go: Same filter and projection. However, Java's projection is `include(GroupMember.Fields.ROLE)` which only includes the role field (plus `_id` by default). Go's projection is `bson.M{"role": 1, "_id.uid": 1}` — this includes `_id.uid` explicitly but the Java version doesn't. This is actually a minor difference in what's returned; the Java version would include the full `_id` (both gid and uid) since `_id` is always included, while Go explicitly includes only `_id.uid`. However, the Go version may not return `_id.gid` in the projected result, which means the full key might not be available. This could be a bug depending on how the result is used.

## findGroupMemberRole

Java: Filters by `_id.uid` = userId AND `_id.gid` = groupId. Projects `role`. Returns `GroupMember::getRole`.

Go: Filters by `_id` (full composite key). Projects `role`. Returns `member.Role`. Correct.

## findMemberIdsByGroupId

Java: Filters by `_id.gid` = groupId. Projects `_id.uid`. Returns user IDs.

Go: **Missing from repository.** There is no single-groupId `findMemberIdsByGroupId` equivalent. The `FindMemberIdsByGroupIds` accepts multiple group IDs but is mapped to `queryGroupMemberIds`.

## findUserJoinedGroupIds

Java: Filters by `_id.uid` = userId. Projects `_id.gid`.

Go: Same, correct.

## findUsersJoinedGroupIds

Java: Filters `inIfNotNull(_id.gid, groupIds)` AND `inIfNotNull(_id.uid, userIds)`. Has pagination via `paginateIfNotNull(page, size)`. Projects `_id.gid`.

Go: Has filter logic but **missing pagination**. The comment says "Pagination logic can be added if needed" but it's not implemented.

## isMemberMuted

Java: Filters by `_id.gid`, `_id.uid`, `med > new Date()`. Uses `mongoClient.exists`.

Go: Filters by full `_id` composite key and `med > time.Now()`. Uses `CountDocuments` with limit 1 — functionally equivalent to exists. Correct.

---

Now let me double-check the `FindGroupMemberKeyAndRolePairs` projection more carefully. The Java version uses `QueryOptions.newBuilder(1).include(GroupMember.Fields.ROLE)` — in MongoDB, `_id` is always included in projections. So the Java query returns `_id` (full composite: gid + uid) and `role`. The Go version uses `bson.M{"role": 1, "_id.uid": 1}` which explicitly includes only `_id.uid`, excluding `_id.gid`. Since `_id` is normally fully included, but when you explicitly project sub-fields of `_id`, MongoDB only returns those sub-fields. This means the Go version won't return `_id.gid` in the document, which could be a functional difference.

Actually, let me reconsider: in MongoDB, when you use `{"_id.uid": 1}`, it only includes the `uid` field within `_id`, not the full `_id`. The Java code uses `include(GroupMember.Fields.ROLE)` which by default includes the full `_id` document. So the Go code is missing `_id.gid` in the projection. However, since the `groupID` is passed as a parameter and used in the filter, this may not matter practically — but it means the returned `GroupMember` objects won't have their full key populated.

## Method-by-method bugs:

## deleteAllGroupMembers
- [ ] **Behavior difference with nil/empty groupIds**: Java's `deleteAllGroupMembers` with `null` groupIds produces an empty filter (deleting ALL group members in the collection). The Go `DeleteByGroupIDs` returns early with no deletion when `len(groupIDs) == 0`. This changes the behavior when groupIds is nil/empty.
- [ ] **Missing ClientSession support**: Java accepts a `@Nullable ClientSession session` and passes it to `mongoClient.deleteMany`. The Go repository's `DeleteByGroupIDs` doesn't accept or use a session, and the service layer ignores the `session` parameter entirely when calling `DeleteByGroupIDs`.

## updateGroupMembers
- [ ] **Missing ClientSession support**: Java accepts a `@Nullable ClientSession session` parameter and passes it to `mongoClient.updateMany`. The Go version does not accept or pass a session context.

## countMembers (multi-parameter)
- [ ] **Method completely missing**: The Java overload `countMembers(@Nullable Set<Long> groupIds, @Nullable Set<Long> userIds, @Nullable Set<GroupMemberRole> roles, @Nullable DateRange joinDateRange, @Nullable DateRange muteEndDateRange)` has no equivalent in Go. There is no method accepting multiple filter parameters for counting members.

## findGroupsMembers
- [ ] **Method is a stub**: Returns `nil, nil` with a TODO comment. The full implementation is missing — the method should filter by groupIds, userIds, roles, joinDateRange, muteEndDateRange with pagination.

## findGroupMemberIds(Set<Long> groupIds)
- [ ] **Method missing**: The Java overload accepting `Set<Long> groupIds` (multiple group IDs) has no direct equivalent in the Go repository. `FindGroupMemberIDs` only handles a single `groupID`. `FindMemberIdsByGroupIds` exists but is mapped to `queryGroupMemberIds`, a different Java method.

## findGroupMemberKeyAndRoleParis
- [ ] **Projection difference**: Java projects only `role` (with `_id` fully included by default), while Go explicitly projects `{"role": 1, "_id.uid": 1}` which excludes `_id.gid`. This means the returned Go `GroupMember` objects won't have `ID.GroupID` populated, unlike the Java version where the full key is available.

## findMemberIdsByGroupId(Long groupId)
- [ ] **Method missing**: The Java `findMemberIdsByGroupId(Long groupId)` (single groupId parameter) has no equivalent Go repository method. The Go `FindMemberIdsByGroupIds` accepts a slice of group IDs and corresponds to a different Java method (`queryGroupMemberIds`).

## findUsersJoinedGroupIds
- [ ] **Missing pagination**: Java applies `paginateIfNotNull(page, size)` to support pagination, but the Go version has no pagination implementation — the comment says "Pagination logic can be added if needed" but it's not applied to the query options.

## FindGroupMemberIDs (single groupID overload) — correction, this one is actually correct.

Let me also re-examine the `FindGroupMemberIDs` method more carefully against Java's `findGroupMemberIds(Long groupId)`:

Java: `Filter.newBuilder(1).eq(GroupMember.Fields.ID_GROUP_ID, groupId)` → filters on `_id.gid`
Go: `filter := bson.M{"_id.gid": groupID}` → correct
Java: `QueryOptions.newBuilder(1).include(GroupMember.Fields.ID_USER_ID)` → projects `_id.uid`
Go: `options.Find().SetProjection(bson.M{"_id.uid": 1})` → correct

This looks correct.

## findGroupMemberIds
- [ ] **No bugs detected** for the single-groupId overload.

---

Here is my final analysis:

## deleteAllGroupMembers
- [ ] **Behavioral difference with nil/empty groupIds**: Java's `deleteAllGroupMembers(null, session)` produces an empty filter that deletes ALL group members. Go's `DeleteByGroupIDs` returns early with no-op when `len(groupIDs) == 0`, so passing nil/empty never deletes anything.
- [ ] **Missing ClientSession support**: The Java method passes a `@Nullable ClientSession session` to the MongoDB operation. The Go repository and service both ignore the session parameter entirely.

## updateGroupMembers
- [ ] **Missing ClientSession support**: Java accepts a `@Nullable ClientSession session` and passes it to `mongoClient.updateMany`. The Go version has no session parameter.

## countMembers (multi-parameter)
- [ ] **Method completely missing**: The Java overload `countMembers(@Nullable Set<Long> groupIds, @Nullable Set<Long> userIds, @Nullable Set<GroupMemberRole> roles, @Nullable DateRange joinDateRange, @Nullable DateRange muteEndDateRange)` has no equivalent in the Go codebase.

## findGroupsMembers
- [ ] **Method is an unimplemented stub**: Returns `nil, nil` with `// TODO: implement full query logic`. None of the filter logic (groupIds, userIds, roles, dateRange filters) or pagination is implemented.

## findGroupMemberIds(Set<Long> groupIds)
- [ ] **Multi-group overload missing**: Java has an overload `findGroupMemberIds(Set<Long> groupIds)` that queries by multiple group IDs. The Go `FindGroupMemberIDs` only accepts a single `groupID`. The Go `FindMemberIdsByGroupIds` serves a different Java method (`queryGroupMemberIds`).

## findGroupMemberKeyAndRoleParis
- [ ] **Incomplete projection**: Java uses `include(role)` which returns the full `_id` (both gid and uid) plus role. Go projects `{"role": 1, "_id.uid": 1}`, which excludes `_id.gid` from results. The returned `GroupMember` objects will have `ID.GroupID` unset (zero value).

## findMemberIdsByGroupId(Long groupId)
- [ ] **Method missing**: Java's `findMemberIdsByGroupId(Long groupId)` has no direct Go equivalent. The closest Go method `FindMemberIdsByGroupIds` takes a slice and is mapped to the different Java method `queryGroupMemberIds`.

## findUsersJoinedGroupIds
- [ ] **Missing pagination**: Java applies `paginateIfNotNull(page, size)` via `QueryOptions`. The Go version ignores the `page` and `size` parameters entirely — no skip/limit is applied to the find options.

# GroupQuestionRepository.java
*Checked methods: updateQuestion(Long questionId, @Nullable String question, @Nullable Set<String> answers, @Nullable Integer score), findGroupId(Long questionId)*

Now I have the full picture. Let me summarize my findings:

## updateQuestion

**Java**: `updateQuestion(Long questionId, @Nullable String question, @Nullable Set<String> answers, @Nullable Integer score)`
- Filter by `_id = questionId`
- Uses `setIfNotNull` for fields: question (`q`), answers (`ans`), score (`score`)
- Returns `Mono<UpdateResult>` — the raw MongoDB update result

**Go**: `Update(ctx context.Context, questionID int64, question *string, answers []string, score *int) (bool, error)`
- Filter by `_id = questionID` ✓
- Sets fields only when non-nil ✓
- Field names match (`q`, `ans`, `score`) ✓
- Returns `(bool, error)` instead of `UpdateResult` — this is a design difference (returns whether modified > 0), not a bug.

The Go implementation correctly mirrors the Java logic for `updateQuestion`.

## findGroupId

**Java**: `findGroupId(Long questionId)` in `GroupQuestionRepository`
- Queries `groupJoinQuestion` collection
- Filters by `_id = questionId`
- Projects only `gid` field
- Returns the `groupId` as `Mono<Long>`

**Go**: There is **no** `FindGroupId` method on `GroupJoinQuestionRepository`. The method is completely missing from both the interface and the implementation. The `FindGroupId` that exists in Go is on `GroupJoinRequestRepository`, which is a completely different entity (join requests, not join questions).

---

## updateQuestion
- [ ] No bugs found. The Go implementation correctly mirrors the Java `updateQuestion` logic — filter by `_id`, set-if-not-null for `question`, `answers`, and `score` fields, with matching BSON field names (`q`, `ans`, `score`).

## findGroupId
- [ ] **Missing method**: The Go `GroupJoinQuestionRepository` (interface and implementation in `group_join_question_repository.go`) is missing the `FindGroupId` method entirely. The Java `GroupQuestionRepository.findGroupId(Long questionId)` queries the `groupJoinQuestion` collection by `_id`, projects only the `gid` field, and returns the group ID. No equivalent exists in Go — the only `FindGroupId` in the Go codebase is on `GroupJoinRequestRepository` (`group_join_request_repository.go:258`), which queries a different collection (`groupJoinRequest`).
- [ ] **Missing service-level method**: The Java `GroupQuestionService` (line 309) calls `groupQuestionRepository.findGroupId(questionId)`, but the Go `GroupQuestionService` (`group_question_service.go`) has no corresponding method that delegates to a question-repository `FindGroupId`.

# GroupRepository.java
*Checked methods: updateGroups(Set<Long> groupIds, @Nullable Long typeId, @Nullable Long creatorId, @Nullable Long ownerId, @Nullable String name, @Nullable String intro, @Nullable String announcement, @Nullable Integer minimumScore, @Nullable Boolean isActive, @Nullable Date creationDate, @Nullable Date deletionDate, @Nullable Date muteEndDate, @Nullable Date lastUpdatedDate, @Nullable Map<String, Object> userDefinedAttributes, @Nullable ClientSession session), countGroups(@Nullable Set<Long> ids, @Nullable Set<Long> typeIds, @Nullable Set<Long> creatorIds, @Nullable Set<Long> ownerIds, @Nullable Boolean isActive, @Nullable DateRange creationDateRange, @Nullable DateRange deletionDateRange, @Nullable DateRange lastUpdatedDateRange, @Nullable DateRange muteEndDateRange), countOwnedGroups(Long ownerId), countOwnedGroups(Long ownerId, Long groupTypeId), findGroups(@Nullable Set<Long> ids, @Nullable Set<Long> typeIds, @Nullable Set<Long> creatorIds, @Nullable Set<Long> ownerIds, @Nullable Boolean isActive, @Nullable DateRange creationDateRange, @Nullable DateRange deletionDateRange, @Nullable DateRange lastUpdatedDateRange, @Nullable DateRange muteEndDateRange, @Nullable Integer page, @Nullable Integer size), findAllNames()*

Now I have all the information needed. Let me analyze each method systematically.

## updateGroups

**Java** (lines 62-95): Takes `groupIds`, and nullable params: `typeId`, `creatorId`, `ownerId`, `name`, `intro`, `announcement`, `minimumScore`, `isActive`, `creationDate`, `deletionDate`, `muteEndDate`, `lastUpdatedDate`, `userDefinedAttributes`, `session`. Builds a filter with `in(ID, groupIds)`, builds an update with `setIfNotNull` for each field and `setUserDefinedAttributesIfNotEmpty`. Uses `updateMany` with session support.

**Go**: The `UpdateGroup` method at line 138 only updates a **single** group by `groupID`. There is no `UpdateGroups` (plural) method that accepts multiple `groupIDs`, all the nullable fields, `userDefinedAttributes`, or `session` support.

## countGroups

**Java** (lines 110-131): Takes `ids`, `typeIds`, `creatorIds`, `ownerIds`, `isActive`, `creationDateRange`, `deletionDateRange`, `lastUpdatedDateRange`, `muteEndDateRange`. Builds filter with `inIfNotNull`, `eqIfNotNull`, `addBetweenIfNotNull`.

**Go** (group_repository.go lines 207-209): `CountGroups` only takes a `dateRange` parameter and delegates to `CountCreatedGroups`. This is missing all the filter parameters: `ids`, `typeIds`, `creatorIds`, `ownerIds`, `isActive`, `deletionDateRange`, `lastUpdatedDateRange`, `muteEndDateRange`.

## countOwnedGroups(Long ownerId)

**Java** (lines 133-137): Simple filter `eq(OWNER_ID, ownerId)`. No deletion date filter.

**Go** (lines 132-135): Adds `"dd": bson.M{"$exists": false}` filter for deletion date, which the Java version does **not** have.

## countOwnedGroups(Long ownerId, Long groupTypeId)

**Java** (lines 139-144): Filter with `eq(OWNER_ID, ownerId)` AND `eq(TYPE_ID, groupTypeId)`.

**Go**: This overload is completely **missing**. There is no method that accepts both `ownerID` and `groupTypeId`.

## findGroups

**Java** (lines 146-171): Takes `ids`, `typeIds`, `creatorIds`, `ownerIds`, `isActive`, `creationDateRange`, `deletionDateRange`, `lastUpdatedDateRange`, `muteEndDateRange`, `page`, `size`. Builds comprehensive filter + pagination via `paginateIfNotNull`.

**Go**: 
- `FindGroups` (lines 37-54) only takes `groupIDs`, hardcodes `"dd": bson.M{"$exists": false}` (not in Java), and has no pagination.
- `QueryGroups` (lines 57-98) takes `groupIDs`, `name`, `lastUpdatedDate`, `skip`, `limit` — missing `typeIds`, `creatorIds`, `ownerIds`, `isActive`, `creationDateRange`, `deletionDateRange`, `muteEndDateRange` filters. Also hardcodes `"dd": bson.M{"$exists": false}` which the Java `findGroups` does not.

## findAllNames

**Java** (lines 181-185): Uses `QueryOptions.newBuilder(1).include(Group.Fields.NAME)` — projection to only include the `name` field. Calls `mongoClient.findAll`.

**Go** (user_repository.go lines 213-223): `FindAllNames` on the **user** repository calls `FindMany` with no projection, returning full user documents. This is (a) on the wrong repository (user instead of group), and (b) missing the projection to only include the name field.

Now let me check if there's a `FindAllNames` on the group repository — there isn't one in group_repository.go.

---

Here is the bug report:

## updateGroups

- [ ] **Missing method entirely**: There is no `UpdateGroups` method in the Go `group_repository.go` that accepts multiple `groupIDs []int64`, all the nullable update fields (`typeId`, `creatorId`, `ownerId`, `name`, `intro`, `announcement`, `minimumScore`, `isActive`, `creationDate`, `deletionDate`, `muteEndDate`, `lastUpdatedDate`), `userDefinedAttributes`, and an optional `session`. The existing `UpdateGroup` (singular) only updates a single group by ID and takes a generic `bson.M`, losing the field-level null-safety (`setIfNotNull`) and `setUserDefinedAttributesIfNotEmpty` logic from Java.

## countGroups

- [ ] **Missing filter parameters**: The Go `CountGroups` only accepts a `dateRange` and delegates to `CountCreatedGroups`. It is missing all the filter parameters from the Java signature: `ids`, `typeIds`, `creatorIds`, `ownerIds`, `isActive`, `deletionDateRange`, `lastUpdatedDateRange`, `muteEndDateRange`.

## countOwnedGroups(ownerId)

- [ ] **Spurious deletion-date filter**: The Go version adds `"dd": bson.M{"$exists": false}` to filter out deleted groups, but the Java version (`countOwnedGroups(Long ownerId)`, line 133-137) does **not** filter by deletion date. This will return incorrect (lower) counts.

## countOwnedGroups(ownerId, groupTypeId)

- [ ] **Missing overload entirely**: The Go `group_repository.go` has no method corresponding to the Java `countOwnedGroups(Long ownerId, Long groupTypeId)` (lines 139-144) that filters by both `ownerId` and `groupTypeId`.

## findGroups

- [ ] **Missing filter parameters**: Neither `FindGroups` nor `QueryGroups` supports the full Java parameter set. Missing filters: `typeIds`, `creatorIds`, `ownerIds`, `isActive`, `creationDateRange`, `deletionDateRange`, `muteEndDateRange`.
- [ ] **Spurious deletion-date filter**: Both Go methods (`FindGroups` line 41, `QueryGroups` line 73-75) hardcode `"dd": bson.M{"$exists": false}` to filter out deleted groups, but the Java `findGroups` does **not** filter by deletion date. This changes query semantics — deleted groups are excluded when they shouldn't be.
- [ ] **Missing pagination in FindGroups**: The `FindGroups` method has no `page`/`size` pagination parameters. The Java version supports `paginateIfNotNull(page, size)`.
- [ ] **`QueryGroups` uses `$gt` instead of `$gte` for `lastUpdatedDate`**: Line 67 uses `bson.M{"$gt": *lastUpdatedDate}` but the Java version uses `gtIfNotNull` which translates to `$gt`. This specific one is correct but the `QueryGroups` method itself doesn't map to the Java `findGroups` — it maps to `findNotDeletedGroups` with different semantics.

## findAllNames

- [ ] **Missing from group repository entirely**: There is no `FindAllNames` method on the Go `GroupRepository`. The Java `findAllNames()` (lines 181-185) lives on `GroupRepository`, not `UserRepository`.
- [ ] **Missing projection**: The `FindAllNames` on the user repository fetches full documents without a projection. The Java version uses `.include(Group.Fields.NAME)` to only return the `name` field, which is a significant performance optimization.

# GroupVersionRepository.java
*Checked methods: updateVersion(Long groupId, String field), updateVersion(Long groupId, boolean updateMembers, boolean updateBlocklist, boolean joinRequests, boolean joinQuestions), findInvitations(Long groupId)*

Now I have all the information needed to perform the comparison. Let me analyze each method.

## Method 1: `updateVersion(Long groupId, String field)` (Java)

**Java** (lines 61-67): Filters by `_id = groupId`, sets `field` to `new Date()`, calls `updateOne` (no upsert).

**Go** `UpdateVersion` (lines 48-55): Filters by `_id = groupID`, sets `field` to `time.Now()`, calls `UpdateOne` with `SetUpsert(true)`.

**Bug**: The Go code uses `SetUpsert(true)`, but the Java code does NOT use upsert. This changes behavior — if no matching document exists, Java does nothing while Go creates a new partial document.

## Method 2: `updateVersion(Long groupId, boolean updateMembers, boolean updateBlocklist, boolean joinRequests, boolean joinQuestions)` (Java)

**Java** (lines 69-84): Filters by `_id = groupId`, creates a single `Date now` object, then conditionally sets `MEMBERS`, `BLOCKLIST`, `JOIN_REQUESTS`, `JOIN_QUESTIONS` fields each to that same `now` timestamp. Calls `updateOne` (no upsert).

**Go**: This method does NOT exist in the Go code. There is no multi-field conditional update version method. The `UpdateVersion` method only accepts a single `field` string parameter, so there's no way to update multiple version fields atomically with the same timestamp using boolean flags.

**Bug**: The entire multi-field `updateVersion` overload is missing.

## Method 3: `findInvitations(Long groupId)` (Java)

Looking at the Java code, `GroupVersionRepository.findInvitations(Long groupId)` (lines 95-102) filters by `_id = groupId`, projects only the `INVITATIONS` field (`"invt"`), and returns `GroupVersion::getInvitations`.

In the Go code, there is NO `FindInvitations` method on the `GroupVersionRepository`. The `findSpecificVersion` helper at line 170-191 handles `"bl"`, `"jr"`, `"jq"`, `"mbr"` but does NOT handle `"invt"` (invitations).

Looking at the `group_invitation_repository.go` file — `FindInvitations` there is a completely different method that queries the `groupInvitation` collection (not `groupVersion`). It finds group invitation records, not the version timestamp for invitations.

**Bug**: The `findInvitations(Long groupId)` method from `GroupVersionRepository` is missing. The `findSpecificVersion` helper doesn't include a case for `"invt"`, and there's no `FindInvitations` method on the `GroupVersionRepository`.

---

## updateVersion(Long groupId, String field)

- [ ] The Go `UpdateVersion` method uses `SetUpsert(true)` (line 51), but the Java `updateVersion(Long groupId, String field)` does **not** use upsert. This changes behavior: if the document doesn't exist, Java does nothing while Go inserts a new partial document with only the specified field set.

## updateVersion(Long groupId, boolean updateMembers, boolean updateBlocklist, boolean joinRequests, boolean joinQuestions)

- [ ] This entire method is **missing** from the Go `GroupVersionRepository`. The Java version atomically updates multiple version fields (members, blocklist, joinRequests, joinQuestions) using boolean flags, all set to the same `Date now` timestamp. There is no Go method that supports this multi-field conditional update pattern.

## findInvitations(Long groupId)

- [ ] This method is **missing** from the Go `GroupVersionRepository`. The Java version queries the `groupVersion` collection with a projection on the `"invt"` field and returns the invitations timestamp. The Go `findSpecificVersion` helper does not include a case for `"invt"`, and no `FindInvitations` method exists on `GroupVersionRepository`. The `FindInvitations` in `groupInvitationRepository` is an entirely different method querying a different collection.

# GroupBlocklistService.java
*Checked methods: authAndBlockUser(@NotNull Long requesterId, @NotNull Long groupId, @NotNull Long userIdToBlock, @Nullable ClientSession session), unblockUser(@NotNull Long requesterId, @NotNull Long groupId, @NotNull Long userIdToUnblock, @Nullable ClientSession session, boolean updateBlocklistVersion), findBlockedUserIds(@NotNull Long groupId, @NotNull Set<Long> userIds), isBlocked(@NotNull Long groupId, @NotNull Long userId), queryBlockedUsers(@Nullable Set<Long> groupIds, @Nullable Set<Long> userIds, @Nullable DateRange blockDateRange, @Nullable Set<Long> requesterIds, @Nullable Integer page, @Nullable Integer size), updateBlockedUsers(@NotEmpty Set<GroupBlockedUser.@ValidGroupBlockedUserKey Key> keys, @Nullable @PastOrPresent Date blockDate, @Nullable Long requesterId)*

Now I have a complete picture. Let me analyze each method carefully.

## authAndBlockUser

- **Java**: Validates not-null params → checks `requesterId.equals(userIdToBlock)` (self-block check) → calls `isOwnerOrManager` → if not auth, returns error → calls `isGroupMember(groupId, userIdToBlock)` → if isGroupMember: wraps delete member + insert blocked user in a transaction (with `updateVersion` for members AND blocklist), and if session is provided uses it, otherwise creates new transaction → if NOT isGroupMember: inserts blocked user + updates blocklist version only (with error logging). The version update in the member case updates both members and blocklist versions (`true, true, false, false`).
- **Go**: No self-block check (`requesterId == userId`) → checks role via `FindGroupMemberRole` → checks if target is a member → if target is a manager/owner, adds role-based restriction for managers (not in Java) → blocks user → if target was a member, deletes member (ignoring error with `_`) → updates blocklist version. The Go version calls `UpdateBlocklistVersion` in all cases, not `UpdateVersion(groupId, true, true, false, false)` for the member case. Also, the order is wrong: Go blocks first then removes member (Java removes member then blocks, in a transaction).

## unblockUser

- **Java**: Validates not-null → checks `isOwnerOrManager` → if not auth, returns specific error → deletes by key → if `updateBlocklistVersion` is true AND was actually blocked, updates blocklist version (with error logging, not failing on error) → returns whether was blocked.
- **Go (`AuthAndUnblockUser`)**: Checks role → calls `UnblockUser` → always updates blocklist version. `UnblockUser` just calls `Delete` and doesn't return whether the user was actually blocked. `AuthAndUnblockUser` always updates blocklist version regardless of whether the user was actually unblocked. Also, the `UnblockUser` Go method has no `updateBlocklistVersion` parameter — it's unconditional in `AuthAndUnblockUser`.

## findBlockedUserIds

- **Java**: Takes `groupId` and `Set<Long> userIds`, creates keys for each (groupId, userId) pair, calls `findIdsByIds(keys)`, maps result to extract userId from key.
- **Go (`FilterBlockedUserIDs`)**: Takes `groupID` and `[]int64` userIDs. The repo method uses `$in` query on `_id.uid` with groupID filter. The Java version queries by compound key set (exact pairs), while the Go version queries by `{_id.gid: groupID, _id.uid: {$in: userIDs}}`. These are functionally equivalent since gid is fixed and uid is constrained by `$in`.

## isBlocked

- **Java**: Creates key `(groupId, userId)`, calls `existsById(key)`.
- **Go**: Calls `Exists(ctx, groupID, userID)` which does a count query on `{_id.gid, _id.uid}`. Functionally equivalent, though Java likely uses a more efficient `existsById` check.

## queryBlockedUsers

- **Java**: Delegates directly to `groupBlocklistRepository.findBlockedUsers(groupIds, userIds, blockDateRange, requesterIds, page, size)`.
- **Go**: The `QueryBlockedUsers` method in the service only takes `groupID` (single group). There is no service-level method matching the full Java signature `(groupIds, userIds, blockDateRange, requesterIds, page, size)`. The repository has `FindBlockedUsers` with all those params, but the service doesn't expose it.

## updateBlockedUsers

- **Java**: Validates keys not empty → validates each key → validates blockDate is past or present → if all params are null (blockDate and requesterId), returns acknowledged update result immediately → delegates to repo `updateBlockedUsers(keys, blockDate, requesterId)`.
- **Go**: The repository has `UpdateBlockedUsers` but there is **no service-level method** exposed for it at all. The Go service doesn't have an `UpdateBlockedUsers` method.

---

## authAndBlockUser

- [ ] **Missing self-block check**: Java checks `requesterId.equals(userIdToBlock)` and returns `ILLEGAL_ARGUMENT` error "Cannot block oneself". Go has no such check.
- [ ] **Wrong version update for member case**: When the target is a group member, Java calls `groupVersionService.updateVersion(groupId, true, true, false, false)` (updating both members and blocklist versions). Go always calls only `UpdateBlocklistVersion(ctx, groupID)`, missing the member version update.
- [ ] **Wrong operation order**: Java removes the group member first, then inserts the blocked user (in a transaction). Go inserts the blocked user first (line 77), then removes the member (line 85). Without a transaction, if the member deletion fails, the user is both blocked and still a member.
- [ ] **Missing transaction support**: Java wraps the member deletion + blocked user insertion in a MongoDB transaction (with retry). Go performs these as separate, non-transactional operations. If the `DeleteGroupMember` call fails (whose error is also silently discarded with `_`), the data is left in an inconsistent state.
- [ ] **Silently discarding delete-member error**: Go uses `_ = s.groupMemberService.DeleteGroupMember(...)` which swallows any error. In Java, the operation chain is reactive and errors propagate.
- [ ] **Non-member case always updates blocklist version**: In Java's non-member path, the blocklist version update error is caught and logged but doesn't fail the operation (`onErrorResume` → `Mono.empty()`). In Go, `s.groupVersionService.UpdateBlocklistVersion` error is returned directly, which would fail the block operation even though the block itself succeeded.

## unblockUser

- [ ] **Missing return value for whether user was actually blocked**: Java returns `Mono<Boolean>` indicating whether the user was actually unblocked (was previously blocked). Go's `UnblockUser` returns only `error` with no boolean, so callers cannot know if the user was actually blocked before unblocking.
- [ ] **Blocklist version update is unconditional**: In Java, `updateBlocklistVersion` is a parameter — when true, the version is updated only if `wasBlocked` is true. In Go's `AuthAndUnblockUser`, the version is always updated via `s.groupVersionService.UpdateBlocklistVersion` regardless of whether the user was actually blocked.
- [ ] **Missing error-resilient version update**: Java catches and logs errors during blocklist version update (`onErrorResume` → logs and returns empty), so the unblock operation succeeds even if version update fails. Go returns the version update error directly.
- [ ] **`UnblockUser` doesn't map to the Java `unblockUser` correctly**: The Java `unblockUser` includes both auth check logic and the `updateBlocklistVersion` parameter. Go splits these into `AuthAndUnblockUser` (with auth) and `UnblockUser` (without auth), but the parameter `updateBlocklistVersion` is lost — neither method accepts it.

## findBlockedUserIds

- [ ] No bugs found — the Go `FilterBlockedUserIDs` is functionally equivalent to Java's `findBlockedUserIds`.

## isBlocked

- [ ] No bugs found — functionally equivalent.

## queryBlockedUsers

- [ ] **Missing full-parameter query method**: Java has `queryBlockedUsers(@Nullable Set<Long> groupIds, @Nullable Set<Long> userIds, @Nullable DateRange blockDateRange, @Nullable Set<Long> requesterIds, @Nullable Integer page, @Nullable Integer size)`. The Go service's `QueryBlockedUsers` only accepts a single `groupID`. The repository has `FindBlockedUsers` with the full signature, but it is not exposed at the service level.

## updateBlockedUsers

- [ ] **Missing service-level method entirely**: Java exposes `updateBlockedUsers` as a public service method that validates keys, validates blockDate, short-circuits when all update fields are null, and delegates to the repository. Go has no corresponding service method — only the repository has `UpdateBlockedUsers`.

# GroupInvitationService.java
*Checked methods: authAndCreateGroupInvitation(@NotNull Long groupId, @NotNull Long inviterId, @NotNull Long inviteeId, @Nullable String content), createGroupInvitation(@Nullable Long id, @NotNull Long groupId, @NotNull Long inviterId, @NotNull Long inviteeId, @Nullable String content, @Nullable @ValidRequestStatus RequestStatus status, @Nullable @PastOrPresent Date creationDate, @Nullable @PastOrPresent Date responseDate), authAndRecallPendingGroupInvitation(@NotNull Long requesterId, @NotNull Long invitationId), queryUserGroupInvitationsWithVersion(@NotNull Long userId, boolean areSentByUser, @Nullable Date lastUpdatedDate), authAndQueryGroupInvitationsWithVersion(@NotNull Long userId, @NotNull Long groupId, @Nullable Date lastUpdatedDate), queryInvitations(@Nullable Set<Long> ids, @Nullable Set<Long> groupIds, @Nullable Set<Long> inviterIds, @Nullable Set<Long> inviteeIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange, @Nullable Integer page, @Nullable Integer size), countInvitations(@Nullable Set<Long> ids, @Nullable Set<Long> groupIds, @Nullable Set<Long> inviterIds, @Nullable Set<Long> inviteeIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange), deleteInvitations(@Nullable Set<Long> ids), authAndHandleInvitation(@NotNull Long requesterId, @NotNull Long invitationId, @NotNull @ValidResponseAction ResponseAction action, @Nullable String reason), updateInvitations(@NotEmpty Set<Long> invitationIds, @Nullable Long inviterId, @Nullable Long inviteeId, @Nullable String content, @Nullable @ValidRequestStatus RequestStatus status, @Nullable @PastOrPresent Date creationDate, @Nullable @PastOrPresent Date responseDate)*

Now I have all the information needed for a thorough comparison. Let me compile the findings.

## authAndCreateGroupInvitation

- [ ] **Missing `isAllowedToInviteUser` delegation**: The Java version calls `groupMemberService.isAllowedToInviteUser(groupId, inviterId)` which checks: the inviter's role in the group, whether the group is active and not deleted, the group type's invitation strategy, and whether the inviter's role is allowed per that strategy. The Go version reimplements this logic locally but does so **incorrectly** — it calls `groupService.QueryGroupTypeIdIfActiveAndNotDeleted` then `groupMemberService.IsGroupMember` (checking if the **invitee** is a member), then manually evaluates the strategy. The Java `isAllowedToInviteUser` also returns `GROUP_INVITER_NOT_MEMBER` when the inviter has no role at all (not a member), which is a different check from the strategy-based allow list. The Go code uses `requesterRole != nil` as a proxy for membership, but for strategies `ALL` and `ALL_REQUIRING_APPROVAL`, the Go code sets `allowed = true` regardless of whether the inviter is even a member — this is actually **correct** per Java's logic (`isAllowedToInviteUserWithSpecifiedRole` returns OK for ALL strategies regardless of role). However...

- [ ] **Missing `isAllowedToBeInvited` check**: The Java `authAndCreateGroupInvitation` calls `groupMemberService.isAllowedToBeInvited(groupId, inviteeId)` after the inviter permission check. This checks whether the invitee is already a group member AND whether the invitee is on the group blocklist. The Go code only checks `IsGroupMember` (for invitee membership) but **completely omits the blocklist check** (`groupBlocklistService.isBlocked`). This means a blocked user can be invited in the Go version.

- [ ] **Missing `SEND_GROUP_INVITATION_TO_GROUP_NOT_REQUIRING_USERS_APPROVAL` check**: In Java, after the strategy is resolved, there's a check `if (!strategy.requiresApproval())` that returns `SEND_GROUP_INVITATION_TO_GROUP_NOT_REQUIRING_USERS_APPROVAL`. This prevents sending invitations to groups where approval is NOT required (those groups use a different flow). The Go version **inverts** this logic — it auto-accepts when `!strategy.RequiresApproval()` and creates the invitation as `ACCEPTED` status, which is fundamentally different behavior. The Java code **rejects** such invitations in `authAndCreateGroupInvitation`.

- [ ] **Auto-accepting invitations when approval is not required**: Following from above, when `!strategy.RequiresApproval()`, the Go code sets `status = RequestStatusAccepted` and then **adds the invitee as a group member** directly. The Java `authAndCreateGroupInvitation` never does this — it only creates PENDING invitations. This is a major behavioral difference.

- [ ] **Wrong error code for invitee-already-member**: The Go code uses `codes.AddUserToGroupWithSizeLimitReached` when the invitee is already a member, but Java uses `SEND_GROUP_INVITATION_TO_GROUP_MEMBER` (a different, more specific code).

- [ ] **Missing `maxContentLength` validation**: The Java version validates `Validator.maxLength(content, "content", maxContentLength)`. The Go version has no content length validation at all.

- [ ] **Missing null content normalization**: Java normalizes `content == null ? "" : content` and also sets it to `""` in `createGroupInvitation` when `content == null`. The Go version passes content as-is without this normalization.

## createGroupInvitation

- [ ] **Missing field: expirationDate**: The Java `GroupInvitation` constructor has a 9th parameter for `expirationDate` (set to `null` in `authAndCreateGroupInvitation`). The Go `GroupInvitation` struct has no `ExpirationDate` field, so this field is lost.

- [ ] **Missing `maxContentLength` validation**: Java validates `Validator.maxLength(content, "content", maxContentLength)` but Go does not.

- [ ] **Missing `validRequestStatus` validation**: Java calls `DataValidator.validRequestStatus(status)` but Go does not validate the status.

- [ ] **Missing `pastOrPresent` validation for dates**: Java validates both `creationDate` and `responseDate` with `Validator.pastOrPresent()`. Go has no such validation.

- [ ] **Missing content null-to-empty normalization**: Java sets `content = ""` when null, Go does not.

- [ ] **Missing default status handling**: Java defaults `status` to `RequestStatus.PENDING` when null, but the Go version of `createGroupInvitation` is not directly callable (it's inlined into `AuthAndCreateGroupInvitation`), so this validation is effectively missing.

- [ ] **Missing version update error logging**: In Java, version update errors are logged with `LOGGER.error(...)` via `onErrorResume`. In Go, version update errors are silently swallowed with `_ = ...`.

## authAndRecallPendingGroupInvitation

- [ ] **Missing property-based allow/deny check**: Java checks `allowRecallPendingInvitationByOwnerAndManager` and `allowRecallPendingInvitationBySender` properties, and returns `RECALLING_GROUP_INVITATION_IS_DISABLED` if both are false. The Go version has no such property check — recall is always allowed.

- [ ] **Missing expiration check**: Java checks `groupInvitationRepository.isExpired(invitation.getCreationDate().getTime())` and returns `RECALL_NON_PENDING_GROUP_INVITATION` with message "The invitation is under the status EXPIRED" if the invitation has expired. The Go version has no expiration check at all.

- [ ] **Missing user version updates**: Java does NOT update user versions (sent/received group invitations versions) in `authAndRecallPendingGroupInvitation`. The Go version calls `s.userVersionService.UpdateSentGroupInvitationsVersion(ctx, inviterID)` and `s.userVersionService.UpdateReceivedGroupInvitationsVersion(ctx, inviteeID)`, which is extra behavior not present in Java.

- [ ] **Wrong query method depending on allowRecallBySender**: Java conditionally queries different fields: if `allowRecallBySender`, it calls `queryGroupIdAndInviterIdAndInviteeIdAndStatus` (to get inviterId for comparison); otherwise it calls `queryGroupIdAndInviteeIdAndStatus`. The Go version always calls `FindGroupIdAndInviterIdAndInviteeIdAndStatus` regardless of property settings.

## queryUserGroupInvitationsWithVersion

- [ ] **Missing `NO_CONTENT` check for empty results**: Java checks `if (groupInvitations.isEmpty())` and throws `NO_CONTENT` exception. The Go version returns an empty list without this check.

- [ ] **Missing expiration status transformation**: Java calls `groupInvitation2proto(groupInvitation, expireAfterSeconds)` which transforms the status to EXPIRED if the invitation has expired based on `expireAfterSeconds`. The Go version returns raw invitations without any expiration status transformation.

- [ ] **Hardcoded pagination (0, 1000)**: The Go version hardcodes `page=0, size=1000` when querying invitations, while Java queries all invitations by the inviter/invitee ID without pagination limits.

- [ ] **Missing `switchIfEmpty` for version**: Java has `.switchIfEmpty(ResponseExceptionPublisherPool.alreadyUpToUpdate())` after the version `flatMap`, meaning if the version is empty (null), it returns "already up to date". The Go version only checks `lastUpdatedDate != nil && version != nil` — if `version` is nil, it proceeds to query invitations anyway.

## authAndQueryGroupInvitationsWithVersion

- [ ] **Missing `NO_CONTENT` check for empty results**: Java checks `if (groupInvitations.isEmpty())` and throws `NO_CONTENT` exception. The Go version returns an empty list without this check.

- [ ] **Missing expiration status transformation**: Same as `queryUserGroupInvitationsWithVersion` — Java transforms expired invitations' statuses.

- [ ] **Hardcoded pagination (0, 1000)**: Same issue — hardcoded limits vs. no pagination in Java.

- [ ] **Missing `switchIfEmpty` for version**: Same as `queryUserGroupInvitationsWithVersion`.

## queryInvitations

- [ ] **Signature mismatch: Go takes single IDs instead of sets**: Java takes `Set<Long>` for ids, groupIds, inviterIds, inviteeIds, and `Set<RequestStatus>` for statuses. Go takes single `*int64` pointers. This means the Go version cannot query by multiple IDs/statuses, which is a significant functional regression.

- [ ] **Missing filter fields**: Java supports filtering by `ids`, `creationDateRange`, `responseDateRange`, `expirationDateRange`, and `statuses` (as sets). The Go version only supports `groupID`, `inviterID`, `inviteeID`, `status` (all single values) and `lastUpdatedDate` (used as creation date > filter). Missing: `ids`, `statuses` set, `creationDateRange`, `responseDateRange`, `expirationDateRange`.

- [ ] **Missing `isExpiredOrNot` filter**: Java applies `isExpiredOrNot(statuses, Fields.CREATION_DATE, getEntityExpirationDate())` which transforms PENDING status to EXPIRED for expired invitations. The Go repository's `buildFilter` has no such logic.

## countInvitations

- [ ] **Same signature/filter mismatch as queryInvitations**: All the same missing filter fields apply — no `ids`, no `statuses` set, no `creationDateRange`, no `responseDateRange`, no `expirationDateRange`.

## deleteInvitations

- [ ] **Missing null/empty set handling**: Java uses `inIfNotNull` for IDs — if `ids` is null, the filter is empty and ALL documents are deleted. The Go version directly uses `$in: ids` without a null check. If `ids` is nil or empty, behavior differs: an empty Go slice in `$in: []` matches nothing, while Java's null set matches everything.

## authAndHandleInvitation

- [ ] **Missing expiration check**: Java checks `groupInvitationRepository.isExpired(invitation.getCreationDate().getTime())` and returns `UPDATE_NON_PENDING_GROUP_INVITATION` with "The invitation is under the status EXPIRED". The Go version has no expiration check.

- [ ] **Missing transaction for ACCEPT action**: Java wraps the ACCEPT action in a transaction (`inTransaction`) that atomically updates the invitation status and adds the group member, with `DuplicateKeyException` handling (if member already exists, returns result with `false`). The Go version does these two operations without a transaction, risking inconsistent state.

- [ ] **Missing DuplicateKeyException handling for ACCEPT**: In Java, if the member already exists (duplicate key), the error is caught and a result with `addedMember=false` is returned. The Go version has no such handling — a duplicate member error would propagate up.

- [ ] **Missing `maxResponseReasonLength` validation**: Java validates `Validator.maxLength(reason, "reason", maxResponseReasonLength)`. The Go version has no reason length validation.

- [ ] **Unnecessary extra DB query for inviterID**: At line 222, the Go version does `s.invRepo.FindByID(ctx, invitationID)` just to get the `InviterID` for version updates. But the earlier query `FindInviteeIdAndGroupIdAndCreationDateAndStatus` already had access to all the data — it just doesn't return `InviterID`. Java handles this differently: `updatePendingInvitationStatus` updates the group version, and the individual action handlers know the invitation data from the initial query. The extra DB query is wasteful but not a behavioral bug.

- [ ] **Missing user version update error handling**: Java logs errors on version update failures but continues. Go silently discards errors with `_ = ...`.

## updateInvitations

- [ ] **Missing `responseDate` update logic based on status**: Java calls `updateResponseDateBasedOnStatus()` which: if `status` is a "processed by responder" status (ACCEPTED, DECLINED, IGNORED) and `responseDate` is null, sets it to now; if `status` is set but NOT processed by responder, **unsets** the responseDate. The Go version simply sets `responseDate` as-is if provided, with no conditional logic based on the status value.

- [ ] **Missing `notEmpty` validation for invitationIds**: Java validates `Validator.notEmpty(invitationIds, "invitationIds")`. Go takes `[]int64` but has no validation.

- [ ] **Missing `maxContentLength` validation**: Java validates content length. Go does not.

- [ ] **Missing `validRequestStatus` validation**: Java validates the status. Go does not.

- [ ] **Missing `pastOrPresent` date validation**: Java validates creationDate and responseDate. Go does not.

- [ ] **Missing early return for all-null updates**: Java checks `if (Validator.areAllNull(inviterId, inviteeId, content, status, creationDate))` and returns early with `ACKNOWLEDGED_UPDATE_RESULT` if no fields to update (note: `responseDate` is NOT checked here because it's derived from status). Go checks `if len(set) == 0` in the repository, which is similar but also includes `responseDate` in the null check, meaning it won't early-return in the same cases.

## updateStatusIfPending (repository)

- [ ] **Missing `isNotExpired` filter condition**: Java adds `.isNotExpired(Fields.CREATION_DATE, getEntityExpirationDate())` to the filter, ensuring that expired invitations cannot have their status updated. The Go repository filter only checks `_id` and `stat: PENDING`, without checking expiration.

- [ ] **Missing conditional reason field handling**: Java uses `.setIfNotNull(Fields.RESPONSE_DATE, new Date())` — note the `setIfNotNull` with a non-null value means response date is always set. The Go version always sets response date. This is actually correct. But for reason, Java uses `setIfNotNull(Fields.REASON, reason)`, while Go also conditionally sets reason only if not nil — this is correct.

- [ ] **Response date always set vs conditionally set**: Java uses `setIfNotNull(Fields.RESPONSE_DATE, new Date())` — this always sets the response date because `new Date()` is never null. The Go version always sets `rd: responseDate` where `responseDate` is `time.Now()`, which is equivalent. This is correct behavior.

# GroupJoinRequestService.java
*Checked methods: queryJoinRequests(@Nullable Set<Long> ids, @Nullable Set<Long> groupIds, @Nullable Set<Long> requesterIds, @Nullable Set<Long> responderIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange, @Nullable Integer page, @Nullable Integer size), authAndHandleJoinRequest(@NotNull Long requesterId, @NotNull Long joinRequestId, @NotNull @ValidResponseAction ResponseAction action, @Nullable String responseReason)*

Now I have all the information needed to compare the two methods. Let me analyze the differences carefully.

## queryJoinRequests

**Java signature:**
```java
queryJoinRequests(@Nullable Set<Long> ids, @Nullable Set<Long> groupIds, @Nullable Set<Long> requesterIds, @Nullable Set<Long> responderIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange, @Nullable Integer page, @Nullable Integer size)
```

**Go signature:**
```go
QueryJoinRequests(ctx context.Context, groupID *int64, requesterID *int64, responderID *int64, status *po.RequestStatus, creationDate *time.Time, page int, size int)
```

The Java version accepts `Set<Long>` (multiple IDs) for ids, groupIds, requesterIds, responderIds, and `Set<RequestStatus>` for statuses, plus DateRange objects for creation, response, and expiration. The Go version only accepts single `*int64` values and a single `*po.RequestStatus`, and is missing `ids`, `responseDate`, and `expirationDate` parameters entirely. It also passes `nil` for those in the call to `FindRequests`.

## authAndHandleJoinRequest

**Java logic:**
1. Validate input (notNull, validResponseAction, maxLength for reason)
2. Find join request by ID; if not found, return NOT_GROUP_OWNER_OR_MANAGER_TO_UPDATE (security: don't reveal existence)
3. Check if requester is owner or manager; if not, return NOT_GROUP_OWNER_OR_MANAGER_TO_UPDATE
4. Check if status is PENDING; if expired based on creation date, return UPDATE_NON_PENDING with "EXPIRED" message
5. If status is not PENDING, return error
6. Switch on action: ACCEPT (transaction: update status + add member, with DuplicateKeyException handling and retry), IGNORE, DECLINE
7. ACCEPT case uses a transaction wrapping both the status update AND the addGroupMember call

**Go logic:**
1. Find join request; if not found, return UPDATE_NON_PENDING (leaks existence info)
2. Check if status is not PENDING, return error
3. Check role via `QueryGroupMemberRole` instead of `IsOwnerOrManager`
4. Update status (no transaction)
5. If accepted, add group member (no transaction, separate from update)
6. Missing: expiration check, IGNORE action, validation of responseAction, DuplicateKeyException handling, transaction wrapping, userVersionService.updateSentGroupJoinRequestsVersion call

## QueryJoinRequests
- [ ] Missing `ids` parameter - Java accepts `Set<Long> ids` but Go has no equivalent
- [ ] `groupIds` downgraded from `Set<Long>` to single `*int64` - cannot query by multiple group IDs
- [ ] `requesterIds` downgraded from `Set<Long>` to single `*int64` - cannot query by multiple requester IDs
- [ ] `responderIds` downgraded from `Set<Long>` to single `*int64` - cannot query by multiple responder IDs
- [ ] `statuses` downgraded from `Set<RequestStatus>` to single `*po.RequestStatus` - cannot query by multiple statuses
- [ ] Missing `responseDateRange` parameter - Java supports filtering by response date range, Go passes `nil`
- [ ] Missing `expirationDateRange` parameter - Java supports filtering by expiration date range, Go passes `nil`

## AuthAndHandleJoinRequest
- [ ] Missing input validation - Java validates `requesterId`, `joinRequestId`, `action` are not null, validates `responseAction` is valid, and validates `responseReason` max length; Go does none of these
- [ ] Wrong error code on request not found - Java returns `NOT_GROUP_OWNER_OR_MANAGER_TO_UPDATE_GROUP_JOIN_REQUEST` (to avoid leaking request existence to unauthorized users), Go returns `UPDATE_NON_PENDING_GROUP_JOIN_REQUEST`
- [ ] Wrong authorization check order - Java checks authorization (isOwnerOrManager) BEFORE checking request status, Go checks status before authorization. This is a security difference: Java prevents unauthorized users from learning the request's status
- [ ] Missing expiration check - Java checks `isExpired(joinRequest.getCreationDate().getTime())` and returns an error with "EXPIRED" status message if the request has expired; Go has no expiration check
- [ ] Missing IGNORE action - Java handles 3 actions: ACCEPT, IGNORE, DECLINE. Go only handles ACCEPT and implicitly DECLINE (the `status` parameter), with no IGNORE equivalent
- [ ] Missing transaction for ACCEPT action - Java wraps the status update AND addGroupMember in a single database transaction with `inTransaction` and `retryWhen(TRANSACTION_RETRY)`; Go performs them as separate non-transactional operations, risking inconsistent state if the member add fails after the status was already updated
- [ ] Missing DuplicateKeyException handling on ACCEPT - Java catches `DuplicateKeyException` during addGroupMember and returns a result with `false` (member already exists) instead of failing; Go will return an error if the user is already a member
- [ ] Missing `userVersionService.updateSentGroupJoinRequestsVersion` call - Java updates the requester's sent join requests version after handling; Go only updates the group version
- [ ] Wrong method for role check - Java uses `isOwnerOrManager(requesterId, groupId, false)` which is a dedicated boolean check; Go uses `QueryGroupMemberRole` which queries the full role, a different approach that may have different caching/performance characteristics
- [ ] Missing `updateGroupInvitationsVersion` call after successful status update - Java's `updatePendingJoinRequestStatus` calls `groupVersionService.updateGroupInvitationsVersion(groupId)` after a successful update; Go's equivalent `UpdateStatusIfPending` in the repository doesn't do this, and the service method only calls `UpdateJoinRequestsVersion`

# GroupMemberService.java
*Checked methods: addGroupMember(@NotNull Long groupId, @NotNull Long userId, @NotNull @ValidGroupMemberRole GroupMemberRole groupMemberRole, @Nullable String name, @Nullable @PastOrPresent Date joinDate, @Nullable Date muteEndDate, @Nullable ClientSession session), addGroupMembers(@NotNull Long groupId, @NotNull Set<Long> userIds, @NotNull @ValidGroupMemberRole GroupMemberRole groupMemberRole, @Nullable String name, @Nullable @PastOrPresent Date joinDate, @Nullable Date muteEndDate, @Nullable ClientSession session), authAndAddGroupMembers(@NotNull Long requesterId, @NotNull Long groupId, @NotNull Set<Long> userIds, @Nullable @ValidGroupMemberRole GroupMemberRole groupMemberRole, @Nullable String name, @Nullable Date muteEndDate, @Nullable ClientSession session), authAndDeleteGroupMembers(@NotNull Long requesterId, @NotNull Long groupId, @NotNull Set<Long> memberIdsToDelete, @Nullable Long successorId, @Nullable Boolean quitAfterTransfer), deleteGroupMember(@NotNull Long groupId, @NotNull Long memberId, @Nullable ClientSession session, boolean updateGroupMembersVersion), deleteGroupMembers(@NotEmpty Collection<GroupMember.Key> keys, @Nullable ClientSession session, boolean updateGroupMembersVersion), updateGroupMember(@NotNull Long groupId, @NotNull Long memberId, @Nullable String name, @Nullable @ValidGroupMemberRole GroupMemberRole role, @Nullable @PastOrPresent Date joinDate, @Nullable Date muteEndDate, @Nullable ClientSession session, boolean updateGroupMembersVersion), updateGroupMembers(@NotEmpty Set<GroupMember.Key> keys, @Nullable String name, @Nullable @ValidGroupMemberRole GroupMemberRole role, @Nullable @PastOrPresent Date joinDate, @Nullable Date muteEndDate, @Nullable ClientSession session, boolean updateGroupMembersVersion), updateGroupMembers(@NotNull Long groupId, @NotEmpty Set<Long> memberIds, @Nullable String name, @Nullable @ValidGroupMemberRole GroupMemberRole role, @Nullable @PastOrPresent Date joinDate, @Nullable Date muteEndDate, @Nullable ClientSession session, boolean updateGroupMembersVersion), isGroupMember(@NotNull Long groupId, @NotNull Long userId, boolean preferCache), isGroupMember(@NotEmpty Set<Long> groupIds, @NotNull Long userId), isMemberMuted(@NotNull Long groupId, @NotNull Long userId, boolean preferCache), queryGroupMemberRole(@NotNull Long userId, @NotNull Long groupId, boolean preferCache), isOwner(@NotNull Long userId, @NotNull Long groupId, boolean preferCache), isOwnerOrManager(@NotNull Long userId, @NotNull Long groupId, boolean preferCache), queryGroupMembers(@Nullable Set<Long> groupIds, @Nullable Set<Long> userIds, @Nullable Set<@ValidGroupMemberRole GroupMemberRole> roles, @Nullable DateRange joinDateRange, @Nullable DateRange muteEndDateRange, @Nullable Integer page, @Nullable Integer size), countMembers(@Nullable Set<Long> groupIds, @Nullable Set<Long> userIds, @Nullable Set<@ValidGroupMemberRole GroupMemberRole> roles, @Nullable DateRange joinDateRange, @Nullable DateRange muteEndDateRange), deleteGroupMembers(boolean updateGroupMembersVersion), queryGroupMembers(@NotNull Long groupId, boolean preferCache), queryGroupMembers(@NotNull Long groupId, @NotEmpty Set<Long> memberIds, boolean preferCache), authAndUpdateGroupMember(@NotNull Long requesterId, @NotNull Long groupId, @NotNull Long memberId, @Nullable String name, @Nullable @ValidGroupMemberRole GroupMemberRole role, @Nullable Date muteEndDate), deleteAllGroupMembers(@Nullable Set<Long> groupIds, @Nullable ClientSession session, boolean updateMembersVersion)*

This confirms a critical issue: `GetRole()` returns the zero value (`GroupMemberRole_OWNER`) when `Role` is nil (not provided). Now I have a comprehensive picture. Let me also quickly check the `AuthAndAddGroupMembers` signature in the controller to see how it's being called.

Now I have all the information needed. Let me compile the comprehensive bug report.

## addGroupMember

- [ ] **Missing `name` parameter**: The Java `addGroupMember` accepts a `@Nullable String name` parameter and passes it to the `GroupMember` constructor, but the Go `AddGroupMember` at `group_member_service.go:66` does not accept a `name` parameter and does not set `Name` on the `po.GroupMember` struct.

- [ ] **Missing `joinDate` parameter**: The Java version accepts a `@Nullable @PastOrPresent Date joinDate` parameter and defaults it to `new Date()` only when null. The Go version ignores the parameter entirely and always sets `JoinDate` to `time.Now()`, making it impossible to backfill a custom join date.

- [ ] **Missing `session` parameter**: The Java version accepts `@Nullable ClientSession session` and passes it to `groupMemberRepository.insert(groupMember, session)`. The Go version does not accept or pass a MongoDB session, so it cannot participate in a transaction.

- [ ] **Missing group version update**: The Java version calls `groupVersionService.updateMembersVersion(groupId)` after the insert (with error logging on failure). The Go version at `group_member_service.go:89` returns immediately after `AddGroupMember` without updating the group members version.

- [ ] **Missing cache update**: The Java version calls `cacheMember(groupMember.getKey().getGroupId(), groupMember)` after insert. The Go version does not update the member cache.

- [ ] **Uses upsert (`UpdateOne` with `SetUpsert(true)`) instead of insert**: The Go repository `AddGroupMember` at `group_member_repository.go:35` uses `UpdateOne` with upsert, while the Java version uses `groupMemberRepository.insert(groupMember, session)`. This means the Go version silently overwrites existing members instead of failing on duplicate insertion.

- [ ] **Incorrect RBAC check**: The Go version at `group_member_service.go:67-75` performs its own simplified RBAC check (owner/manager only) that doesn't exist in the Java `addGroupMember` method. The Java `addGroupMember` is a low-level method with no auth checks — auth is handled separately by `authAndAddGroupMembers`.

## addGroupMembers

- [ ] **Missing `session` parameter**: Java accepts `@Nullable ClientSession session` and passes it to `insertAllOfSameType`. The Go version does not accept a session.

- [ ] **Missing validation**: Java validates `maxSize(userIds, "userIds", 100)` and `pastOrPresent(joinDate, "joinDate")`. The Go version has no size limit validation or joinDate validation.

- [ ] **Inserts one-by-one instead of bulk**: The Go version at `group_member_service.go:392` calls `s.groupMemberRepo.AddGroupMember` in a loop, while the Java version uses `groupMemberRepository.insertAllOfSameType(groupMembers, session)` for a single bulk insert. This is both a performance issue and a correctness issue — if one insert fails partway through, partial members are added.

- [ ] **Missing error handling on version update**: The Go version at `group_member_service.go:399` uses `_ = s.groupVersionService.UpdateMembersVersion(ctx, groupID)` which silently discards errors. The Java version logs the error but continues.

- [ ] **Missing cache invalidation/updates**: The Go version only deletes from cache, while Java does not have explicit per-member cache ops for bulk add (it updates the version but doesn't invalidate/add to the existing member-level cache in `addGroupMembers`).

## authAndAddGroupMembers

- [ ] **Missing `name` parameter**: Java accepts `@Nullable String name` and passes it to `addGroupMembers`. The Go version does not accept or forward a `name` parameter.

- [ ] **Missing `session` parameter**: Java accepts `@Nullable ClientSession session`. The Go version does not.

- [ ] **Incorrect role handling for `onlyAddRequesterToGroup`**: In Java (line 245-274), when `userIds.contains(requesterId)` (onlyAddRequesterToGroup), there are specific checks: if size > 1 → error, if role is null → default to MEMBER, if role != MEMBER → error, and then role must not be OWNER or GUEST. In Go, this logic is completely different — it checks `isOwnerOrManager` instead of `onlyAddRequesterToGroup`, and the join strategy / invitation strategy branching is simplified incorrectly.

- [ ] **Missing invitation strategy-based authorization**: The Java version has a detailed `switch` on `invitationStrategy` (ALL, OWNER_MANAGER_MEMBER, OWNER_MANAGER, OWNER, plus _REQUIRING_APPROVAL variants) to determine if the requester's role is sufficient. The Go version at `group_member_service.go:440-449` only checks if the requester is owner/manager, completely ignoring the group type's invitation strategy.

- [ ] **Missing join strategy check for self-join**: When `onlyAddRequesterToGroup` is true in Java, the method checks the group type's `joinStrategy` (MEMBERSHIP_REQUEST, INVITATION, QUESTION, JOIN_REQUEST) and returns different error codes. The Go version has a simplified check on `groupType.JoinStrategy` but only handles MEMBERSHIP_REQUEST, and uses incorrect error codes.

- [ ] **Missing `requiresApproval` check**: Java checks `invitationStrategy.requiresApproval()` and returns `ADD_USER_TO_GROUP_REQUIRING_USERS_APPROVAL`. The Go version omits this entirely.

- [ ] **Uses `FilterBlockedUserIDs` instead of `findBlockedUserIds`**: Java collects the blocked user IDs and returns an error listing them. Go filters them out silently with `FilterBlockedUserIDs` and proceeds with only the non-blocked users. This is different behavior — Java rejects the entire operation if *any* user is blocked, while Go silently skips blocked users.

- [ ] **Wrong error code for inactive group**: The Go version returns `ResponseStatusCode_UPDATE_INFO_OF_NONEXISTENT_GROUP` when the group is inactive, while Java returns `ResponseStatusCode_ADD_USER_TO_INACTIVE_GROUP`.

- [ ] **Missing Manager role check for adding MANAGER**: Java checks at line 304 that if `finalGroupMemberRole == GroupMemberRole.MANAGER` and `requesterRole != GroupMemberRole.OWNER`, it returns an error. The Go version omits this check.

- [ ] **Controller passes `GetRole()` value instead of pointer**: At `group_service_controller.go:352`, `createReq.GetRole()` returns `GroupMemberRole` (value), but when the proto field `Role` is nil (not set by client), `GetRole()` returns `GroupMemberRole_OWNER` (zero value). This means "no role specified" is indistinguishable from "OWNER specified", which breaks the auth logic.

## authAndDeleteGroupMembers

- [ ] **Missing successorId null check with non-empty memberIdsToDelete**: Java at line 386-389 checks: if `successorId != null` AND `!memberIdsToDelete.isEmpty()` → error "Cannot transfer the group ownership while removing other members". The Go version does not perform this check.

- [ ] **Different behavior when successorId is provided with empty memberIdsToDelete**: Java at line 390-396 calls `groupService.authAndTransferGroupOwnership` and returns `Set.of(requesterId)`. The Go version does not handle the successorId case at all for the transfer path correctly — it deletes the group via `AuthAndDeleteGroup` when no successorId and owner is quitting, which is completely different from the Java behavior.

- [ ] **Missing "cannot quit with others" validation**: Java at line 403-405 returns error if `quitGroup && memberIdsToDelete.size() > 1`. The Go version at line 487-493 sets `isQuitting = true` but does not return an error for this case.

- [ ] **Missing "owner must specify successor to quit" check**: Java at line 407-409 checks `isOwner(requesterId, groupId, false)` and returns `GROUP_OWNER_QUIT_WITHOUT_SPECIFYING_SUCCESSOR` if the owner tries to quit without a successor. The Go version instead tries to delete the entire group via `AuthAndDeleteGroup`, which is wrong.

- [ ] **Missing role comparison for target members**: Java at line 439-441 checks `keyAndRole.getRole().getNumber() <= requesterRoleNumber` to prevent a manager from removing someone with equal or higher rank. The Go version at line 538 only checks if target is OWNER or MANAGER (hardcoded), not using numeric comparison.

- [ ] **Missing `finalMemberIdsToDelete` computation**: Java at line 427-433 computes the actual set of deletable member IDs (filtering out non-existent ones) and returns it. The Go version just passes the original `userIDs` to deletion without this filtering.

- [ ] **Missing `quitAfterTransfer` null handling**: Java at line 394 uses `quitAfterTransfer != null && quitAfterTransfer`. The Go version accepts `bool` (not `*bool`), so it cannot distinguish "not specified" from "false".

- [ ] **Always updates version even when no members deleted**: Go at line 558 always calls `UpdateMembersVersion`, while Java only updates when `deletedCount > 0`.

## deleteGroupMember

- [ ] **Missing validation of `groupId` and `memberId`**: Java validates `@NotNull` for both parameters. The Go version has no null/validity checks.

- [ ] **Missing version update error handling**: The Go version at line 340-341 calls `UpdateMembersVersion` but returns its error directly, while the Java version catches and logs the error, returning the original `DeleteResult`.

- [ ] **Missing cache invalidation logic parity**: Java invalidates cache only when `deletedCount == keys.size()` (i.e., all requested deletes succeeded). The Go version always deletes from cache regardless of whether the DB deletion actually succeeded.

## deleteGroupMembers(keys, session, updateGroupMembersVersion)

- [ ] **Missing validation**: Java validates `Validator.notEmpty(keys, "keys")` and `DataValidator.validGroupMemberKey(key)` for each key. The Go version has no validation.

- [ ] **Missing session parameter**: Java accepts `@Nullable ClientSession session` and passes it to `groupMemberRepository.deleteByIds(keys, session)`. The Go `DeleteGroupMember` does not pass the session to the repo.

- [ ] **Missing group IDs collection and version update**: Java collects all group IDs from the keys, then calls `groupVersionService.updateMembersVersion(groupIds)` (plural version) to update all affected groups. The Go `DeleteGroupMember` at line 340 only updates a single group's version.

- [ ] **Missing conditional cache invalidation**: Java invalidates cache only when `deletedCount == keys.size()`. The Go `DeleteGroupMember` always deletes from cache.

- [ ] **The `deleteGroupMembers` overload accepting `Collection<GroupMember.Key>` is not implemented**: Java has a dedicated `deleteGroupMembers(Collection<GroupMember.Key> keys, ...)` method that handles bulk deletion across multiple groups. The Go version does not have an equivalent bulk method; `DeleteGroupMember` only handles a single member.

## updateGroupMember

- [ ] **Missing `joinDate` parameter**: Java accepts `@Nullable @PastOrPresent Date joinDate` and passes it through. The Go version does not accept `joinDate`.

- [ ] **Missing `session` parameter**: Java accepts `@Nullable ClientSession session`. The Go version does not.

- [ ] **Missing `updateGroupMembersVersion` parameter**: Java accepts `boolean updateGroupMembersVersion`. The Go version always updates the version, with no option to skip it.

- [ ] **Missing validation**: Java validates `Validator.notNull(memberId, "memberId")`. The Go version has no validation.

## updateGroupMembers(keys, name, role, joinDate, muteEndDate, session, updateGroupMembersVersion)

- [ ] **Missing `session` parameter**: Java accepts `@Nullable ClientSession session`. The Go repository method does not accept a session.

- [ ] **Missing `updateGroupMembersVersion` flag**: Java conditionally updates version. The Go service always updates.

- [ ] **Missing early return when all null**: Java checks `Validator.areAllNull(name, role, joinDate, muteEndDate)` and returns `ACKNOWLEDGED_UPDATE_RESULT` without hitting the DB. The Go version doesn't have this optimization.

- [ ] **Missing cache update**: Java calls `updateMembersCache(keys, name, role, joinDate, muteEndDate)` when `modifiedCount == keys.size()`. The Go version does not update cache on member updates (only deletes).

- [ ] **Missing version update for multiple groups**: Java computes group IDs from keys and calls `updateMembersVersion(groupIds)` for multiple groups. The Go version only handles single-group updates.

- [ ] **Missing error logging on version update failure**: Java logs errors from version updates. The Go version either propagates or silently discards them.

## updateGroupMembers(groupId, memberIds, name, role, joinDate, muteEndDate, session, updateGroupMembersVersion)

- [ ] **This overload is not implemented in Go**: The Java version at line 577-603 takes `groupId` and `memberIds` (Set<Long>), converts them to `GroupMember.Key` set, and delegates to the key-based `updateGroupMembers`. The Go version does not have an equivalent method.

## isGroupMember(groupId, userId, preferCache)

- [ ] **Cache structure mismatch**: Java uses `Cache<Long, Map<GroupMember.Key, GroupMember>>` — a per-group map of all members, allowing O(1) lookup of any member within a cached group. Go uses `TTLCache[string, bool]` — a simple boolean cache per (group, user) pair. This means the Go cache cannot be leveraged for other operations like `queryGroupMembers` or `isMemberMuted` that use the same group-level cache in Java.

- [ ] **Missing `preferCache` parameter**: Java has `boolean preferCache` to allow callers to bypass cache. The Go version always uses cache, with no bypass option.

## isGroupMember(groupIds, userId)

- [ ] **Not implemented**: The Java version accepts `@NotEmpty Set<Long> groupIds` and a single `userId`, creates keys, and calls `groupMemberRepository.existsByIds(keys)` to check membership across multiple groups. The Go version does not have this overload.

## isMemberMuted

- [ ] **Incorrect cache logic**: Java checks the member's `muteEndDate` field from the cached `GroupMember` object: `member.getMuteEndDate() != null`. If the member is in cache, it checks whether `muteEndDate` is non-null (i.e., the member has a mute end date that was set, regardless of whether it's in the past). The Go version caches the boolean result of `IsMemberMuted`, which depends on the current time at the moment of caching. A muted member whose mute expires during the cache TTL would incorrectly show as muted.

- [ ] **Missing `preferCache` parameter**: Java has `boolean preferCache`. Go always uses cache.

- [ ] **Missing cache structure parity**: Java's cache stores `GroupMember` objects (including `muteEndDate`) so it can compute muted status from the actual field. Go caches only the boolean, leading to stale results.

## queryGroupMemberRole

- [ ] **Missing `preferCache` parameter**: Java has `boolean preferCache` and checks the `groupIdToMembersCache` before hitting the DB. The Go version always hits the DB.

- [ ] **Missing cache lookup**: Java checks the per-group member cache map for the role. The Go version at line 137-139 just delegates directly to the repository without any cache.

## isOwner

- [ ] **Missing `preferCache` parameter**: Java passes `preferCache` through to `queryGroupMemberRole`. The Go version at line 109-118 does not accept this parameter.

## isOwnerOrManager

- [ ] **Swapped parameter order**: At `group_member_service.go:121`, `IsOwnerOrManager` calls `s.FindGroupMemberRole(ctx, groupID, userID)` but `FindGroupMemberRole` at line 132 takes `(ctx, groupID, userID)`. This is actually correct for `FindGroupMemberRole`, but the method signature `IsOwnerOrManager(ctx, groupID, userID)` has a different parameter order than `IsOwner(ctx, userID, groupID)` at line 109. This inconsistency is confusing but not necessarily a bug since Go doesn't have named parameters at call sites — however, callers need to be careful about argument order.

## queryGroupMembers(groupIds, userIds, roles, joinDateRange, muteEndDateRange, page, size)

- [ ] **Not implemented — repository stub returns nil**: The Go repository `FindGroupsMembers` at `group_member_repository.go:317-320` is a stub returning `(nil, nil)`. The full query logic with filtering by groupIds, userIds, roles, date ranges, and pagination is not implemented.

- [ ] **Missing validation**: Java validates each role with `DataValidator.validGroupMemberRole(role)`.

- [ ] **Missing cache population**: Java populates the member cache when `isMemberCacheEnabled` and certain conditions are met. The Go version has no equivalent.

## countMembers

- [ ] **Missing parameters — only accepts groupID**: The Go repository `CountMembers` at `group_member_repository.go:208` only takes a `groupID` parameter, while Java takes `(groupIds, userIds, roles, joinDateRange, muteEndDateRange)` for flexible counting with filters.

- [ ] **Missing validation**: Java validates roles with `DataValidator.validGroupMemberRole`.

## deleteGroupMembers(updateGroupMembersVersion)

- [ ] **Not implemented**: The Java version at line 1080 calls `groupMemberRepository.deleteAll()` (deletes ALL group members in the entire collection). The Go version does not have this method.

## queryGroupMembers(groupId, preferCache)

- [ ] **Not fully implemented**: Java returns `Mono<Collection<GroupMember>>` with cache support. The Go version does not have this specific overload — it only has `FindGroupMembers` in the repository which always hits the DB.

- [ ] **Missing cache read/write**: Java checks cache when `preferCache && isMemberCacheEnabled`, and populates cache after DB fetch. The Go version has no cache integration for this method.

## queryGroupMembers(groupId, memberIds, preferCache)

- [ ] **Not implemented as a service method**: Java has a service-level `queryGroupMembers(Long groupId, Set<Long> memberIds, boolean preferCache)` with cache support. The Go version only has the repository method `FindGroupMembersWithIds` which bypasses cache entirely.

- [ ] **Missing cache lookup**: Java checks cache first when `preferCache && isMemberCacheEnabled`. The Go version always goes to the DB.

## authAndUpdateGroupMember

- [ ] **Missing validation**: Java validates `Validator.notNull(memberId, "memberId")` and `DataValidator.validGroupMemberRole(role)`. The Go version has no such validation.

- [ ] **Missing `joinDate` parameter**: Java passes `null` for joinDate in the call to `updateGroupMember`. The Go version at line 602 passes `nil` for `joinDate`, which is correct.

- [ ] **Simplified and incorrect authorization logic**: The Java version has nuanced authorization:
  - If `role != null`: checks group is active, then checks if requester is OWNER
  - If `muteEndDate != null`: checks requester can't mute self, checks group type, checks both requester and target roles with numeric comparison, then checks member info update strategy
  - If only `name`: checks group type, checks if target is member, checks selfInfoUpdatable or member info update strategy
  - If all null: returns ACKNOWLEDGED_UPDATE_RESULT immediately
  
  The Go version at lines 563-612 has a much simpler check: if `requesterID != memberID` → owner/manager required; if updating role → owner required. This misses the group-type-based strategy checks, the self-mute prevention, the role numeric comparison, and the selfInfoUpdatable check.

- [ ] **Missing `checkIfAuthorized` flow for muteEndDate**: Java has a dedicated path when `muteEndDate != null` that prevents self-muting, verifies the group type is active, checks both requester and target roles with numeric comparison (`requesterRole.getNumber() >= targetMemberRole.getNumber()`), and applies the member info update strategy. The Go version at line 582-598 does none of this — it only checks if the requester is owner/manager.

- [ ] **Missing group active/deleted check**: Java at line 1234 checks `groupService.isGroupActiveAndNotDeleted(groupId)` before allowing role updates. The Go version does not check group activity status.

- [ ] **Always updates version unconditionally**: The Go version at line 611 always calls `UpdateMembersVersion`. The Java version passes `false` for `updateGroupMembersVersion` in the call to `updateGroupMember` (line 1335), meaning it does NOT update the version.

- [ ] **Missing member info update strategy check**: Java checks `groupType.getMemberInfoUpdateStrategy()` (OWNER, OWNER_MANAGER, OWNER_MANAGER_MEMBER, ALL) and `groupType.getSelfInfoUpdatable()`. The Go version does not check these group type properties.

## deleteAllGroupMembers

- [ ] **Missing null groupIds handling**: Java at line 1349 handles `groupIds == null` by calling `groupIdToMembersCache.invalidateAll()` (clears entire cache). The Go version at line 349 returns immediately when `len(groupIDs) == 0`, which treats null/empty the same — but in Java, null means "delete ALL group members across ALL groups" while empty should mean "nothing to delete".

- [ ] **Missing DB call for null groupIds**: Java calls `groupMemberRepository.deleteAllGroupMembers(groupIds, session)` which handles the null case (deletes all). The Go version returns early without deleting anything when `groupIDs` is empty.

- [ ] **Missing session parameter**: Java accepts `@Nullable ClientSession session` and passes it to the repository. The Go version accepts a `mongo.SessionContext` but does not pass it through to `DeleteByGroupIDs`.

- [ ] **Missing cache invalidation**: Java invalidates the cache (either all or per-group). The Go version does not invalidate any cache.

- [ ] **Missing conditional version update**: Java at line 1355 checks `updateMembersVersion` flag AND `deletedCount > 0`. The Go version always updates when the flag is true regardless of whether anything was actually deleted.

- [ ] **Version update uses different method**: Java calls `groupVersionService.updateMembersVersion(groupIds)` (which handles null groupIds as "update all"). The Go version iterates and calls `UpdateMembersVersion(ctx, groupID)` per group, which is different from the Java batch update.

# GroupQuestionService.java
*Checked methods: checkGroupQuestionAnswerAndGetScore(@NotNull Long questionId, @NotNull String answer, @Nullable Long groupId), queryGroupJoinQuestions(@Nullable Set<Long> ids, @Nullable Set<Long> groupIds, @Nullable Integer page, @Nullable Integer size, boolean withAnswers), countGroupJoinQuestions(@Nullable Set<Long> ids, @Nullable Set<Long> groupIds), deleteGroupJoinQuestions(@Nullable Set<Long> ids), updateGroupJoinQuestions(@NotEmpty Set<Long> ids, @Nullable Long groupId, @Nullable String question, @Nullable Set<String> answers, @Nullable @Min(0)*

Now I have all the information needed. Let me analyze each method carefully.

## checkGroupQuestionAnswerAndGetScore

**Java** (line 128-139): Validates `questionId` and `answer` are not null, then delegates to `groupQuestionRepository.checkQuestionAnswerAndGetScore(questionId, answer, groupId)`. The repository method queries MongoDB with a filter matching `_id`, `ans` (the answer string), and optionally `gid`.

**Go** (`CheckQuestionAnswerAndGetScore`, line 261-276): Fetches the full question document via `FindByID`, checks `groupID` match, then iterates answers locally. This is a **different approach** (client-side checking vs server-side MongoDB query), but functionally equivalent.

However, there's a subtle difference: The Java repository's `checkQuestionAnswerAndGetScore` returns `Mono<Integer>` — it queries MongoDB with the answer as a filter (`"ans": answer`). The Go version fetches the full document and iterates. This works but means the Go `CheckQuestionAnswerAndGetScore` at the service layer fetches more data than needed. The repository also has a `CheckQuestionAnswerAndGetScore` method that does the server-side query, but the service layer doesn't use it.

Wait — actually looking more carefully: The Go service `CheckQuestionAnswerAndGetScore` doesn't use the repository's `CheckQuestionAnswerAndGetScore` at all. It implements the logic itself by fetching the question and iterating. This is functionally equivalent for correctness, just less efficient. Not a bug per se.

But let me check if the Go repository's `CheckQuestionAnswerAndGetScore` has a bug. In Java, the answer matching is done via MongoDB query `"ans": answer` which checks if the answer string exists in the array. In Go repo, the same is done. So the Go service reimplements what the Go repo already does. Not a functional bug.

## queryGroupJoinQuestions

**Java** (line 350-357): `queryGroupJoinQuestions(ids, groupIds, page, size, withAnswers)` — delegates to `groupQuestionRepository.findQuestions(ids, groupIds, page, size, withAnswers)`. The `withAnswers` parameter is passed to the repo to control whether answers are included in the MongoDB projection.

**Go** (`FindQuestions`, line 282-293): The service calls `s.questionRepo.FindQuestions(ctx, ids, groupIds, page, size)` — note the repo method does NOT accept a `withAnswers` parameter. Then the service strips answers if `!withAnswers`. This is functionally equivalent but fetches answers from DB then discards them. Slightly less efficient but not a behavioral bug.

## countGroupJoinQuestions

**Java** (line 359-363): `countGroupJoinQuestions(ids, groupIds)` → `groupQuestionRepository.countQuestions(ids, groupIds)`.

**Go** (`CountQuestions`, line 278-280): `CountQuestions(ctx, ids, groupIds)` → `s.questionRepo.CountQuestions(ctx, ids, groupIds)`. Functionally equivalent.

## deleteGroupJoinQuestions

**Java** (line 365-367): `deleteGroupJoinQuestions(ids)` → `groupQuestionRepository.deleteByIds(ids)`. Accepts `@Nullable Set<Long> ids`.

**Go**: There is **no `DeleteGroupJoinQuestions` method** in the Go service that takes a set of IDs and delegates to the repo. The only delete method is `AuthAndDeleteQuestion` which takes a single `questionID`. The Go repository also lacks a `DeleteByIds` method — only has `Delete(ctx, questionID int64)`.

This is a **missing method** — the admin controller's `DeleteGroupJoinQuestions(ids []int64)` exists as a stub but the service layer has no corresponding bulk delete implementation.

## updateGroupJoinQuestions

**Java** (line 449-468): 
1. Validates `ids` not null
2. Validates `question` max length against `questionContentLimit`
3. Validates `answers` size range (1 to `maxAnswerCount`)
4. Validates `answers` max length against `answerContentLimit`  
5. Validates `score >= 0`
6. If all params are null/falsy (groupId, question, answers, score), returns ACKNOWLEDGED_UPDATE_RESULT immediately
7. Delegates to `groupQuestionRepository.updateQuestions(ids, groupId, question, answers, score)`

**Go** (`UpdateQuestions`, line 295-322):
1. No validation of `question` content length limit
2. No validation of `answers` size range
3. No validation of `answers` content length limit
4. No validation of `score >= 0`
5. Does NOT short-circuit when all params are nil (it proceeds to call the repo, which also checks for empty updateOps)
6. Has additional logic not in Java: looks up group IDs for version updates after the update, which the Java admin-level method does NOT do (only `authAndUpdateGroupJoinQuestion` does version updates)
7. The Go version calls `groupVersionService.UpdateJoinQuestionsVersion` for affected groups — the Java `updateGroupJoinQuestions` does NOT do this

---

Here is my final analysis:

## checkGroupQuestionAnswerAndGetScore

- [ ] The Go `CheckQuestionAnswerAndGetScore` (line 261-276) reimplements the logic client-side instead of using the repository's `CheckQuestionAnswerAndGetScore` method. While functionally equivalent for the score-checking logic, the Go version returns `(int, error)` where a score of `0` could mean either "answer is wrong" or "question has score 0". The Java version delegates to the repo which returns `Mono<Integer>` (nullable), distinguishing "no match" (null/empty) from "score is 0". The Go service always returns `0` for wrong answers, making it impossible to distinguish a wrong answer from a correct answer worth 0 points.

## deleteGroupJoinQuestions

- [ ] **Missing method**: The Go service has no `DeleteGroupJoinQuestions(ids []int64)` equivalent. The Java method accepts `@Nullable Set<Long> ids` and delegates to `groupQuestionRepository.deleteByIds(ids)`. The Go service only has `AuthAndDeleteQuestion(ctx, requesterID, groupID, questionID)` which deletes a single question by ID. There is no bulk delete-by-IDs method in either the service or repository. The admin controller stub `DeleteGroupJoinQuestions(ids []int64)` has no service-level implementation to call.

## updateGroupJoinQuestions

- [ ] **Missing validations**: The Go `UpdateQuestions` method does not validate `question` against `questionContentLimit`, does not validate `answers` size range (1 to `maxAnswerCount`), does not validate `answers` content against `answerContentLimit`, and does not validate `score >= 0`. The Java version performs all of these validations before proceeding.
- [ ] **Missing short-circuit for empty update**: The Java version returns `ACKNOWLEDGED_UPDATE_RESULT` immediately when all update parameters (`groupId`, `question`, `answers`, `score`) are null/falsy. The Go version does not short-circuit — it proceeds to look up affected questions and call the repo, though the repo does its own empty-check on `updateOps`.
- [ ] **Unintended version update logic**: The Go `UpdateQuestions` method adds extra logic to update group join question versions (`groupVersionService.UpdateJoinQuestionsVersion`) for affected groups after the update. The Java `updateGroupJoinQuestions` method does NOT do this — only the `authAndUpdateGroupJoinQuestion` (single-question, user-facing) method handles version updates. The Go version adds version updates to the admin-level batch update method, which is behavior not present in the Java original.

# GroupService.java
*Checked methods: createGroup(@NotNull Long creatorId, @NotNull Long ownerId, @Nullable String groupName, @Nullable String intro, @Nullable String announcement, @Nullable @Min(value = 0), authAndDeleteGroup(boolean queryGroupMemberIds, @NotNull Long requesterId, @NotNull Long groupId), deleteGroupsAndGroupMembers(@Nullable Set<Long> groupIds, @Nullable Boolean deleteLogically), queryGroups(@Nullable Set<Long> ids, @Nullable Set<Long> typeIds, @Nullable Set<Long> creatorIds, @Nullable Set<Long> ownerIds, @Nullable Boolean isActive, @Nullable DateRange creationDateRange, @Nullable DateRange deletionDateRange, @Nullable DateRange lastUpdatedDateRange, @Nullable DateRange muteEndDateRange, @Nullable Set<Long> memberIds, @Nullable Integer page, @Nullable Integer size), queryGroupTypeIdIfActiveAndNotDeleted(@NotNull Long groupId), authAndTransferGroupOwnership(@NotNull Long requesterId, @NotNull Long groupId, @NotNull Long successorId, boolean quitAfterTransfer, @Nullable ClientSession session), countOwnedGroups(@NotNull Long ownerId), countOwnedGroups(@NotNull Long ownerId, @NotNull Long groupTypeId), countGroups(@Nullable Set<Long> ids, @Nullable Set<Long> typeIds, @Nullable Set<Long> creatorIds, @Nullable Set<Long> ownerIds, @Nullable Boolean isActive, @Nullable DateRange creationDateRange, @Nullable DateRange deletionDateRange, @Nullable DateRange lastUpdatedDateRange, @Nullable DateRange muteEndDateRange, @Nullable Set<Long> memberIds), count()*

Now I have all the information needed. Let me compile the complete bug report.

## createGroup

- [ ] **Missing `ownerId` parameter**: Java accepts both `creatorId` and `ownerId` as separate parameters. Go hardcodes `OwnerID: &creatorID`, ignoring `ownerId`. This means the group owner is always set to the creator, which is incorrect when they differ (e.g., admin-created groups on behalf of another user).
- [ ] **Missing `groupTypeId` parameter**: Java accepts `groupTypeId` and sets it on the `Group` entity. Go does not accept or set `TypeID` at all, so all groups are created without a type.
- [ ] **Missing `announcement` parameter**: Java accepts `announcement` and sets it on the `Group` entity. Go does not accept or set `Announcement`.
- [ ] **Missing `muteEndDate` parameter**: Java accepts `muteEndDate` and sets it on the `Group` entity. Go does not accept or set `MuteEndDate`.
- [ ] **Missing `deletionDate` parameter**: Java accepts `deletionDate` and sets it on the `Group` entity. Go does not accept or set `DeletionDate`.
- [ ] **Missing `creationDate` parameter with fallback**: Java accepts an optional `creationDate` that defaults to `now`. Go always uses `time.Now()` with no override option.
- [ ] **Missing `isActive` parameter and defaulting logic**: Java defaults `isActive` to the `activateGroupWhenCreated` property when null. Go never sets `IsActive` at all, so the field is stored as `nil`/absent in MongoDB instead of an explicit boolean.
- [ ] **Missing `minimumScore` default to 0**: Java sets `minimumScore` to `0` when null (`minimumScore == null ? 0 : minimumScore`). Go only sets it when the pointer is non-nil, leaving it absent in MongoDB when not provided.
- [ ] **Missing `LastUpdatedDate` field assignment**: Java sets `new Date()` (now) as the last updated date on the Group entity. Go never sets `LastUpdatedDate` on the created group.
- [ ] **No transactional wrapping**: Java wraps insert + addGroupMember + version upsert in a single MongoDB transaction (`groupRepository.inTransaction(...).retryWhen(TRANSACTION_RETRY)`). Go executes these sequentially with no transaction. A failure in `AddGroupMember` leaves an orphaned group document (the "rollback" via `DeleteGroup` is a hard delete, not equivalent to atomic transaction rollback).

## authAndDeleteGroup

- [ ] **Missing `queryGroupMemberIds` parameter**: Java has a `boolean queryGroupMemberIds` parameter. Go omits it entirely.
- [ ] **Wrong return type**: Java returns `Mono<Set<Long>>` (group member IDs, or empty set). Go returns only `error`. The `queryGroupMemberIds=true` branch in Java first queries member IDs, then deletes, then returns those IDs. Go cannot return member IDs at all.
- [ ] **Missing conditional branch for `queryGroupMemberIds=true`**: When `queryGroupMemberIds` is true in Java, it calls `groupMemberService.queryGroupMemberIds(groupId, false)` before deletion and returns the member IDs. This entire code path is absent in Go.

## deleteGroupsAndGroupMembers

- [ ] **Missing `deleteLogically` parameter and logical deletion support**: Java accepts `Boolean deleteLogically` (defaults to `deleteGroupLogicallyByDefault` property). When true, it calls `groupRepository.updateGroupsDeletionDate(groupIds, session)` to soft-delete. Go always performs a hard physical delete via `DeleteGroup`, never using the existing `UpdateGroupsDeletionDate` repository method.
- [ ] **No transactional wrapping**: Java wraps the entire operation in `groupRepository.inTransaction(...).retryWhen(TRANSACTION_RETRY)`. Go has no transaction, so partial failures can leave the database in an inconsistent state (e.g., groups deleted but members remain).
- [ ] **Missing message sequence ID cleanup**: Java calls `messageService.deleteGroupMessageSequenceIds(groupIds)`. Go has a TODO comment but does not implement this.
- [ ] **Missing conversation deletion**: Java calls `conversationService.deleteGroupConversations(groupIds, session)`. Go has a TODO comment but does not implement this.
- [ ] **Missing deleted groups counter increment**: Java increments `deletedGroupsCounter.increment(count)` when `count > 0`. Go does not track this metric.
- [ ] **Iterative single-delete instead of batch**: Java calls `groupRepository.deleteByIds(groupIds, session)` as a single batch operation. Go iterates over each group ID calling `s.groupRepo.DeleteGroup(ctx, groupID)` individually, which is less efficient and not atomic.

## queryGroups

- [ ] **Missing filter parameters**: The Go repository `QueryGroups` only supports `groupIDs`, `name`, and `lastUpdatedDate`. It completely lacks support for: `typeIds`, `creatorIds`, `ownerIds`, `isActive`, `creationDateRange`, `deletionDateRange`, `muteEndDateRange` — all of which are queryable in Java.
- [ ] **Missing member-based group filtering**: Java checks if `memberIds` is non-empty and calls `queryGroupIdsFromGroupIdsAndMemberIds(ids, memberIds)` to intersect group IDs with groups containing those members before querying. Go has no equivalent logic.
- [ ] **Wrong method mapped**: Java's `queryGroups` with all filter parameters maps to `groupRepository.findGroups(ids, typeIds, creatorIds, ownerIds, isActive, ..., page, size)`. The Go `AuthAndQueryGroups` delegates to `groupRepo.QueryGroups` which is a different, simplified method. The Go `group_controllers.go` `QueryGroupsWithQuery` is a stub returning `nil` and never calls the service.

## authAndTransferGroupOwnership

- [ ] **Missing `successorId == requesterId` early return**: Java checks `if (successorId.equals(requesterId)) { return Mono.empty(); }` before any database calls. Go does not. When `requesterId == successorId` but requester is not the owner, Java silently succeeds (no-op) while Go returns `ErrNotGroupOwner` — a behavioral difference.
- [ ] **Missing `isAllowedToCreateGroupAndHaveGroupType` permission check**: Java's `checkAndTransferGroupOwnership` queries `queryGroupTypeId(groupId)` then calls `isAllowedToCreateGroupAndHaveGroupType(successorId, groupTypeId)` to verify the successor is allowed to own a group of that type. Go only checks if the successor is a group member, skipping the entire permission/type check.
- [ ] **Wrong order of role updates**: Java demotes/deletes the old owner first, then promotes the successor to OWNER. Go promotes the successor first (`UpdateGroupMemberRole` to OWNER), then demotes/deletes the old owner. This means Go briefly has two owners in the group member table.
- [ ] **Missing group document `oid` update in Java vs. present in Go**: Java's `checkAndTransferGroupOwnership` does not directly update the group document's `ownerId` field — it only updates member roles. Go adds `update := bson.M{"oid": successorID}` + `s.groupRepo.UpdateGroup(ctx, groupID, update)`. While updating the document seems reasonable, this is a behavioral difference from Java where the `oid` on the group document may not be updated during transfer.

## countOwnedGroups

- [ ] **Missing overload with `groupTypeId` parameter**: Java has `countOwnedGroups(@NotNull Long ownerId, @NotNull Long groupTypeId)` which counts groups owned by `ownerId` filtered by `groupTypeId`. Go only has the single-parameter `CountOwnedGroups(ctx, ownerID)`. The Go repository's comments reference the two-parameter Java overload but the implementation does not support it. This is used by `isAllowedCreateGroupWithGroupType` to enforce per-type group ownership limits.

## countGroups

- [ ] **Completely wrong signature and behavior**: Java's `countGroups` accepts 10 filter parameters (`ids`, `typeIds`, `creatorIds`, `ownerIds`, `isActive`, `creationDateRange`, `deletionDateRange`, `lastUpdatedDateRange`, `muteEndDateRange`, `memberIds`) and supports member-based filtering. Go's `CountGroups` only accepts `dateRange *turmsmongo.DateRange` and delegates to `CountCreatedGroups`. It cannot filter by any of the other parameters.
- [ ] **Missing member-based counting**: Java checks if `memberIds` is non-empty and calls `queryGroupIdsFromGroupIdsAndMemberIds` before counting. Go has no equivalent.

## count

No bugs — the Go `Count` delegates to `groupRepo.Count(ctx)` which calls `CountDocuments(ctx, bson.M{})`, matching Java's `groupRepository.countAll()`.

## count (verified)

The `count()` method is correctly ported.

# GroupTypeService.java
*Checked methods: queryGroupTypes(@Nullable Integer page, @Nullable Integer size), addGroupType(@Nullable Long id, @NotNull @NoWhitespace String name, @NotNull @Min(1), queryGroupTypes(@NotNull Collection<Long> groupTypeIds)*

Now I have all the information needed to provide a complete review. Let me also verify the UpdateTypes field pointer type mismatch — the service passes pointers from struct fields but the PO struct uses value types.

## QueryGroupTypes(page, size)

**Java** (`queryGroupTypes(@Nullable Integer page, @Nullable Integer size)`): Delegates to `groupTypeRepository.findAll(page, size)` which applies pagination (skip/limit).

**Go** (`QueryGroupTypes`): Delegates to `groupTypeRepo.FindGroupTypes(ctx, nil, page, size)` but the repository method at line 118 has the comment `"Add skip and limit if page/size exist... (Skipped detailed impl for brevity if not mapped exactly)"` and then **ignores page and size entirely** — no `options.Skip()` / `options.Limit()` is applied to the `Find` call.

## AddGroupType

**Java** (`addGroupType`): Has extensive validation (`Validator.notNull`, `Validator.noWhitespace`, `Validator.min`), generates an ID if `id == null` via `node.nextLargeGapId(ServiceType.GROUP_TYPE)`, creates a `GroupType` with all fields, inserts it, and updates the in-memory cache `idToGroupType`.

**Go** (`AddGroupType`): Takes a pre-built `*po.GroupType` and directly inserts — no validation, no ID generation fallback, no in-memory cache update.

## QueryGroupTypes(groupTypeIds)

**Java** (`queryGroupTypes(@NotNull Collection<Long> groupTypeIds)`): First checks the in-memory cache `idToGroupType`, returns early if all found in cache. If any missing, clears the list and falls back to `groupTypeRepository.findByIds`, updating the cache for each result.

**Go** (`QueryGroupTypesByIds`): Directly queries the database with no cache involvement.

## UpdateTypes field type mismatch

The Go `UpdateGroupTypes` service method passes `&update.Name` (which is `*string` from a `*po.GroupType`), but the repository's `UpdateTypes` expects `*int32` for `groupSizeLimit`, `*constant.GroupInvitationStrategy`, etc. The PO fields are value types (not pointers), so taking their addresses will always produce non-nil pointers, meaning **zero-value fields will be written to the database** even when they shouldn't be updated.

---

## queryGroupTypes(@Nullable Integer page, @Nullable Integer size)

- [ ] **Pagination is completely ignored**: The repository method `FindGroupTypes` accepts `page` and `size` parameters but never applies them as `Skip`/`Limit` options on the MongoDB `Find` call. The comment on line 118 explicitly says the implementation was skipped. This means the Go version always returns **all** group types, while the Java version correctly paginates.

## addGroupType(@Nullable Long id, @NotNull @NoWhitespace String name, @NotNull @Min(1) Integer groupSizeLimit, ...)

- [ ] **Missing all input validation**: The Java version validates that `name` is non-null and has no whitespace, `groupSizeLimit >= 1`, and that all enum and boolean parameters are non-null. The Go version performs zero validation and accepts whatever is passed in.
- [ ] **Missing ID auto-generation when `id` is null/zero**: Java generates a new ID via `node.nextLargeGapId(ServiceType.GROUP_TYPE)` when `id == null`. The Go version has no fallback — if the caller doesn't set `ID` on the `GroupType` struct, it defaults to `0` (which collides with `DEFAULT_GROUP_TYPE_ID`).
- [ ] **Missing in-memory cache update after successful insert**: Java calls `.doOnSuccess(unused -> idToGroupType.put(groupType.getId(), groupType))` to keep the in-memory cache in sync. The Go version has no in-memory cache at all, so every `QueryGroupType`/`QueryGroupTypesByIds` call will always hit the database.
- [ ] **Missing return of the created GroupType**: Java returns `Mono<GroupType>` (the inserted entity). Go returns only `error`, discarding the created entity from the caller.

## queryGroupTypes(@NotNull Collection<Long> groupTypeIds)

- [ ] **Missing in-memory cache lookup with early return optimization**: The Java version first attempts to resolve all group types from the in-memory `idToGroupType` cache, and only queries the database if any are missing. The Go version always queries the database directly via `FindGroupTypes`, bypassing any caching opportunity.
- [ ] **Missing cache population from database results**: Java calls `.doOnNext(type -> idToGroupType.put(type.getId(), type))` to populate the in-memory cache with freshly fetched results. The Go version has no such cache population logic.

## UpdateGroupTypes (bonus: critical pointer type bug)

- [ ] **`UpdateGroupTypes` always writes all fields including zero-values**: The service method takes `&update.Name`, `&update.GroupSizeLimit`, etc. Since `po.GroupType` fields are value types (not pointers), their addresses are always non-nil, so the repository's nil-checks (e.g., `if groupSizeLimit != nil`) will always pass. This means zero-value fields (`""`, `0`, `false`) will be written to the database even when the caller only intended to update specific fields. This will corrupt existing data by overwriting fields with their zero values.

# GroupVersionService.java
*Checked methods: queryGroupJoinRequestsVersion(@NotNull Long groupId), queryGroupInvitationsVersion(@NotNull Long groupId), updateVersion(@NotNull Long groupId, boolean updateMembers, boolean updateBlocklist, boolean joinRequests, boolean joinQuestions), updateMembersVersion(@NotNull Long groupId), updateMembersVersion(@Nullable Set<Long> groupIds), updateMembersVersion(), updateBlocklistVersion(@NotNull Long groupId), updateJoinRequestsVersion(@NotNull Long groupId), updateJoinQuestionsVersion(@NotNull Long groupId), updateSpecificVersion(@NotNull Long groupId, @NotNull String field), updateSpecificVersion(@NotNull String field), updateSpecificVersion(@Nullable Set<Long> groupIds, @NotNull String field), upsert(@NotNull Long groupId, @NotNull Date timestamp), delete(@Nullable Set<Long> groupIds, @Nullable ClientSession session)*

Now I have a complete picture. Let me carefully compare each method.

## queryGroupJoinRequestsVersion(@NotNull Long groupId)

**Java:** Calls `groupVersionRepository.findJoinRequests(groupId)` which does a targeted projection query for just the `jr` field.

**Go:** Calls `s.groupVersionRepo.FindVersion(ctx, groupID)` which fetches the **entire document** (no projection), then extracts `v.JoinRequests`. This is functionally correct but less efficient — it fetches all fields instead of projecting just the one needed. However, the repository also has `FindJoinRequests` which does use projection. The service doesn't use it.

Wait — this is a performance concern, not a behavior bug. Let me verify the logic correctness more carefully.

## updateVersion(@NotNull Long groupId, boolean updateMembers, boolean updateBlocklist, boolean joinRequests, boolean joinQuestions)

**Java:** Builds an update with conditional `$set` using `setIfTrue` — only sets fields where the corresponding boolean is `true`. Does **not** use upsert.

**Go:** The `UpdateVersion` repository method takes only a single `field string` parameter. There is **no** method that accepts multiple boolean flags to conditionally update multiple fields at once. The `updateVersion(groupId, updateMembers, updateBlocklist, joinRequests, joinQuestions)` overload is completely **missing**.

## updateMembersVersion(@NotNull Long groupId)

**Java:** Delegates to `updateSpecificVersion(groupId, GroupVersion.Fields.MEMBERS)` which calls `groupVersionRepository.updateVersion(groupId, field)`. This uses `updateOne` **without** upsert.

**Go:** Calls `s.groupVersionRepo.UpdateMembersVersion(ctx, groupID)` → `r.UpdateVersion(ctx, groupID, "mbr")` which uses `UpdateOne` **with upsert=true**. Different behavior — Go will create a new document if it doesn't exist, Java won't.

## updateMembersVersion(@Nullable Set<Long> groupIds)

**Java:** Delegates to `updateSpecificVersion(groupIds, GroupVersion.Fields.MEMBERS)` which calls `groupVersionRepository.updateVersions(groupIds, field)`. In the repository, when `groupIds` is null, it uses `inIfNotNull` which produces an empty filter (matches all documents). When groupIds is empty set, same behavior.

**Go:** No equivalent method exists in the `GroupVersionService`. The repository has `UpdateVersions(ctx context.Context, groupIDs []int64, field string)` but it's not exposed through the service layer. Also, when `groupIDs` is empty (length 0), the Go code applies the filter to **no documents** (empty filter would match all, but the Go code only enters the `$in` branch when `len(groupIDs) > 0`). When `groupIDs` is empty, it creates `filter := bson.M{}` which matches **all** documents — this is actually correct and matches Java behavior for null groupIds. But wait — in Go, an empty slice has `len == 0`, so it falls through to `filter := bson.M{}` which is an empty filter matching all docs. That matches Java. However the method is **missing from the service**.

## updateMembersVersion()

**Java:** Delegates to `updateSpecificVersion(GroupVersion.Fields.MEMBERS)` which calls `groupVersionRepository.updateVersions(field)` — updates ALL documents. **Missing from Go service**.

## updateSpecificVersion(@NotNull Long groupId, @NotNull String field)

**Java:** Uses `updateOne` without upsert. Returns `result.getModifiedCount() > 0`.

**Go:** The repository `UpdateVersion` uses `SetUpsert(true)`. Different upsert behavior.

## updateSpecificVersion(@NotNull String field)

**Java:** Updates all documents in the collection. **Missing from Go service**.

## updateSpecificVersion(@Nullable Set<Long> groupIds, @NotNull String field)

**Java:** Updates documents matching the given groupIds. When groupIds is null, matches all. **Missing from Go service**.

## upsert(@NotNull Long groupId, @NotNull Date timestamp)

**Java:** Creates `new GroupVersion(groupId, timestamp, timestamp, timestamp, timestamp, timestamp)` — sets 5 fields: members, blocklist, joinRequests, joinQuestions, invitations. The constructor takes 6 args (groupId + 5 timestamps).

**Go Upsert:** Sets 6 fields in the `$set`: `mbr`, `bl`, `jr`, `jq`, `invt`, **and `info`**. The Java version only sets 5 fields (members, blocklist, joinRequests, joinQuestions, invitations). The Go version adds an extra `"info": timestamp` field that doesn't exist in the Java version's `GroupVersion` PO.

## delete(@Nullable Set<Long> groupIds, @Nullable ClientSession session)

**Java:** Uses `inIfNotNull` — when groupIds is null, the filter becomes empty and `deleteMany` deletes **nothing** (actually in Java, `inIfNotNull` with null returns `Filter.newBuilder(0)` which is an empty filter... let me check. Actually, looking at the Java BaseRepository:

```java
public Mono<DeleteResult> deleteByIds(@Nullable Collection<K> ids) {
    Filter filter = Filter.newBuilder(1)
            .inIfNotNull(DomainFieldName.ID, ids);
    return mongoClient.deleteMany(entityClass, filter);
}
```

When `ids` is null, `inIfNotNull` does not add the `$in` condition, so it becomes an empty filter which matches all documents. This would delete all documents.

**Go:** `DeleteByIds` returns `nil` immediately when `len(groupIDs) == 0`. In Go, a nil/empty slice has `len == 0`, so this would **not delete anything**. This is different from Java where null ids → empty filter → delete all.

Wait, but the Java `deleteByIds` is inherited from BaseRepository, and the Java service calls it with potentially null groupIds. Let me reconsider: `inIfNotNull` — when ids is null, it doesn't add the `$in` clause, leaving the filter empty. An empty filter in MongoDB matches all documents. So `deleteMany` with an empty filter would delete ALL group version documents. But the Go code returns nil (does nothing) when the slice is empty/null. This is a behavioral difference but arguably the Go behavior is safer.

Actually, re-reading the Java `Filter.newBuilder(1).inIfNotNull(DomainFieldName.ID, ids)` — when `ids` is null, `inIfNotNull` is documented to not add the condition. So the filter is empty, and `deleteMany("")` deletes all. The Go code with empty slice skips entirely. This is a **behavioral difference**.

---

Now let me compile the final findings:

## updateVersion(groupId, updateMembers, updateBlocklist, joinRequests, joinQuestions)

- [ ] **Missing method entirely.** The Go code has no equivalent of the multi-field conditional update. The Java version uses `setIfTrue` to conditionally set only the fields whose boolean flags are true. The Go `UpdateVersion` only accepts a single field string, not multiple boolean flags.

## updateMembersVersion(@Nullable Set<Long> groupIds)

- [ ] **Missing method from service layer.** Java has `updateMembersVersion(Set<Long> groupIds)` that updates members version for specific groups. No equivalent exists in `GroupVersionService`. The repository has `UpdateVersions` but it's not exposed through the service.

## updateMembersVersion()

- [ ] **Missing method entirely.** Java has a no-arg `updateMembersVersion()` that updates the members version for ALL groups in the collection. No equivalent exists in Go service or repository.

## updateSpecificVersion(@NotNull String field)

- [ ] **Missing method entirely.** Java has `updateSpecificVersion(String field)` that updates a specific field for ALL group versions. No equivalent exists in Go service.

## updateSpecificVersion(@Nullable Set<Long> groupIds, @NotNull String field)

- [ ] **Missing method entirely.** Java has `updateSpecificVersion(Set<Long> groupIds, String field)` that updates a specific field for given groups. No equivalent exists in Go service.

## upsert(groupId, timestamp)

- [ ] **Extra field `info` in Go that doesn't exist in Java's GroupVersion PO.** Java creates `GroupVersion(groupId, timestamp, timestamp, timestamp, timestamp, timestamp)` setting 5 fields (members, blocklist, joinRequests, joinQuestions, invitations). Go sets 6 fields including `"info": timestamp`. The Java `GroupVersion` PO has no `info` field.

## updateSpecificVersion(groupId, field) / updateMembersVersion(groupId) / updateBlocklistVersion(groupId) / updateJoinRequestsVersion(groupId) / updateJoinQuestionsVersion(groupId)

- [ ] **Upsert behavior difference.** The Go `UpdateVersion` repository method uses `SetUpsert(true)`, meaning it will create a new document if none exists. The Java `updateVersion(groupId, field)` uses plain `updateOne` without upsert, so it will silently do nothing if no document exists. This changes behavior: Go creates partial documents (with only one field set), Java never creates documents via these update methods.

## delete(groupIds, session)

- [ ] **Behavioral difference with null/empty IDs.** Java's `deleteByIds` with null groupIds produces an empty filter that deletes ALL documents. Go's `DeleteByIds` with empty slice returns nil and deletes nothing. The Go code would silently succeed where Java would delete everything.

## queryGroupJoinRequestsVersion / queryGroupInvitationsVersion (and similar query methods)

- [ ] **Go fetches entire document instead of using projection.** Java repository methods use `.include(field)` projection to fetch only the needed field. Go service methods call `FindVersion` which fetches the entire document. The Go repository has efficient projection-based methods (`FindJoinRequests`, `FindMembers`, etc.) but the service layer doesn't use them, opting for the full-document `FindVersion` instead.

# MessageController.java
*Checked methods: createMessages(@QueryParam(defaultValue = "true"), queryMessages(@QueryParam(required = false), queryMessages(@QueryParam(required = false), countMessages(@QueryParam(required = false), updateMessages(Set<Long> ids, @RequestBody UpdateMessageDTO updateMessageDTO), deleteMessages(Set<Long> ids, @QueryParam(required = false)*

## CreateMessages
- [ ] Method is completely unimplemented (only contains `// TODO: implement`). Missing all core logic: parsing the `send` query param, calling `messageService.authAndSaveAndSendMessage` with all fields from `CreateMessageDTO` (senderId, senderDeviceType, senderIp conversion via `InetAddressUtil.ipStringToBytes`, id, isGroupMessage, isSystemMessage, text, records, targetId, burnAfter, referenceId, preMessageId), and returning the HTTP response.

## QueryMessages (non-paginated)
- [ ] Method is completely unimplemented. Missing all core logic: accepting filter params (ids, areGroupMessages, areSystemMessages, senderIds, targetIds, deliveryDateStart/End, deletionDateStart/End, recallDateStart/End, size, ascending), constructing `DateRange` objects, calling `messageService.queryMessages` with page=0 and `getPageSize(size)`, and returning `HttpHandlerResult.okIfTruthy`.

## QueryMessages (paginated, "page" endpoint)
- [ ] Method is completely missing from the Go code. The Java version has a second `queryMessages` overload mapped to `@GetMapping("page")` that accepts an additional `page` param, calls `countMessages` for pagination total, then queries with `page`/`size`, and returns `HttpHandlerResult.page(count, completeMessagesFlux)`.

## CountMessages
- [ ] Method is completely unimplemented. Missing all core logic: the conditional branching on `divideBy` (NOOP vs. time-division), calling `countSentMessages`, `countSentMessagesOnAverage`, and `checkAndQueryBetweenDate` for divided statistics, building the `MessageStatisticsDTO`, and returning the result. The Java version has complex conditional logic for deciding which counts to compute based on which date params are non-null.

## UpdateMessages
- [ ] Method is completely unimplemented. Missing all core logic: extracting `senderIp` from `UpdateMessageDTO`, converting it via `InetAddressUtil.ipStringToBytes`, calling `messageService.updateMessages` with (senderId, senderDeviceType, ids, isSystemMessage, text, records, burnAfter, recallDate, ip bytes, null), mapping result to `UpdateResultDTO`, and returning `HttpHandlerResult.okIfTruthy`.

## DeleteMessages
- [ ] Method is completely unimplemented. Missing all core logic: accepting `ids` set and `deleteLogically` query param, calling `messageService.deleteMessages(ids, deleteLogically)`, mapping to `DeleteResultDTO`, and returning `HttpHandlerResult.okIfTruthy`.

# Message.java
*Checked methods: groupId()*

Now I have a clear picture. Let me verify the Java `groupId()` method's behavior and compare it to the Go `FindMessageGroupId` implementation.

**Java `groupId()` method (on Message PO):**
```java
@Nullable
public Long groupId() {
    return isGroupMessage != null && isGroupMessage
            ? targetId
            : null;
}
```

This is a computed accessor on the in-memory `Message` object that returns `targetId` only if `isGroupMessage` is true, otherwise returns `null`.

**Go `FindMessageGroupId` (in repository):**
```go
func (r *MessageRepository) FindMessageGroupId(ctx context.Context, messageID int64) (*int64, error) {
    filter := bson.M{"_id": messageID}
    var result struct {
        TargetID int64 `bson:"tid"`
    }
    err := r.col.FindOne(ctx, filter, options.FindOne().SetProjection(bson.M{"tid": 1})).Decode(&result)
    if err != nil {
        return nil, err
    }
    return &result.TargetID, nil
}
```

This queries MongoDB for the `tid` field without checking `isGroupMessage`. It will return a non-nil group ID even for private (non-group) messages, where the Java version would return `null`.

## groupId()

- [ ] **Missing `isGroupMessage` check**: The Go `FindMessageGroupId` method unconditionally returns the `targetId` (as `GroupID`) without verifying that the message is actually a group message. The Java `groupId()` returns `targetId` only when `isGroupMessage != null && isGroupMessage` is `true`, and returns `null` otherwise. The Go version should either project the `gm` field and conditionally return `nil` when it's not a group message, or at minimum document that callers must ensure the message is a group message before calling.

- [ ] **Missing `ErrNoDocuments` handling**: When the message is not found in MongoDB, `FindMessageGroupId` returns the raw `mongo.ErrNoDocuments` error to the caller. In contrast, other similar methods in the same repository (e.g., `FindMessageSenderIDAndTargetIDAndIsGroupMessage`) also lack this handling, but the `FindByID` methods in the group join request repository pattern handle it by returning `(nil, nil)`. Inconsistent error handling — though this may be a design choice, the Java `groupId()` would simply return `null` on a missing message context rather than throw.

# MessageRepository.java
*Checked methods: updateMessages(Set<Long> messageIds, @Nullable Boolean isSystemMessage, @Nullable Integer senderIp, @Nullable byte[] senderIpV6, @Nullable Date recallDate, @Nullable String text, @Nullable List<byte[]> records, @Nullable Integer burnAfter, @Nullable ClientSession session), updateMessagesDeletionDate(@Nullable Set<Long> messageIds), existsBySenderIdAndTargetId(Long senderId, Long targetId), countMessages(@Nullable Set<Long> messageIds, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages, @Nullable Set<Long> senderIds, @Nullable Set<Long> targetIds, @Nullable DateRange deliveryDateRange, @Nullable DateRange deletionDateRange), countUsersWhoSentMessage(@Nullable DateRange dateRange, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages), countGroupsThatSentMessages(@Nullable DateRange dateRange), countSentMessages(@Nullable DateRange dateRange, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages), findDeliveryDate(Long messageId), findExpiredMessageIds(Date expirationDate), findMessageGroupId(Long messageId), findMessageSenderIdAndTargetIdAndIsGroupMessage(Long messageId), findMessages(@Nullable Collection<Long> messageIds, @Nullable Collection<byte[]> conversationIds, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages, @Nullable Set<Long> senderIds, @Nullable Set<Long> targetIds, @Nullable DateRange deliveryDateRange, @Nullable DateRange deletionDateRange, @Nullable DateRange recallDateRange, @Nullable Integer page, @Nullable Integer size, @Nullable Boolean ascending), findIsGroupMessageAndTargetId(Long messageId, Long senderId), findIsGroupMessageAndTargetIdAndDeliveryDate(Long messageId, Long senderId), getGroupConversationId(long groupId), getPrivateConversationId(long id1, long id2)*

Now I have all the information needed. Let me carefully compare each method.

## UpdateMessages

- [ ] **Missing `modificationDate` field**: The Java code always sets `Message.Fields.MODIFICATION_DATE` to `new Date()` (unconditionally, not `setIfNotNull`). The Go code completely omits setting `md` (modificationDate) in the `UpdateMessages` method.
- [ ] **Missing `ClientSession` parameter**: The Java method accepts a `@Nullable ClientSession session` parameter and passes it to `mongoClient.updateMany(session, ...)`. The Go method does not accept or use a session/transaction parameter.
- [ ] **Unconditional filter issue**: The Java code always applies `Filter.newBuilder(1).in(DomainFieldName.ID, messageIds)` regardless of whether `messageIds` is null. The Go code returns nil (does nothing) if `messageIDs` is empty. In Java, `messageIds` is `Set<Long>` (non-nullable), so it's always expected to have IDs. However, the Go code also has an incorrect early return path — it checks `len(set) == 0` before even building the filter, which means if all optional parameters happen to be nil, it returns without error even when there are valid message IDs (though in practice, `modificationDate` would have been set unconditionally in Java, preventing this). With the missing `modificationDate`, this becomes a real bug — if only `modificationDate` should be set and nothing else, Go does nothing.

## UpdateMessagesDeletionDate

- [ ] **Missing `inIfNotNull` null-check semantics**: The Java version uses `filter.inIfNotNull(DomainFieldName.ID, messageIds)`, meaning if `messageIds` is null, the filter has no ID constraint and updates ALL messages. The Go version requires `messageIDs` as a non-nil `[]int64` and returns nil if empty, which changes behavior when messageIDs should be nil (i.e., "clear deletion date for all messages").
- [ ] **Unconditional `new Date()` vs parameterized deletion date**: The Java version always sets `DELETION_DATE` to `new Date()` (current time). The Go version accepts a `deletionDate *time.Time` parameter and sets it to `nil` when the parameter is nil, which is semantically different — Java always sets a deletion date, Go can clear it. However, this may be an intentional API design difference — the Go version takes the date as parameter rather than always using `time.Now()`.

## ExistsBySenderIDAndTargetID

- [ ] **Missing `isGroupMessage = false` filter**: The Java code adds `.eq(Message.Fields.IS_GROUP_MESSAGE, false)` to the filter, ensuring only non-group (private) messages are matched. The Go code's filter is `bson.M{"sid": senderID, "tid": targetID}` — it's missing `"gm": false`. This means the Go version will incorrectly return true for group messages where the sender sent to the target's group.

## CountMessages

- [ ] **Missing `messageIds` filter parameter**: The Java method accepts `@Nullable Set<Long> messageIds` and applies `.inIfNotNull(DomainFieldName.ID, messageIds)`. The Go method `CountMessages` does not accept a `messageIDs` parameter at all.
- [ ] **Missing `areSystemMessages` filter parameter**: The Java method accepts `@Nullable Boolean areSystemMessages` and applies `.eqIfNotNull(Message.Fields.IS_SYSTEM_MESSAGE, areSystemMessages)`. The Go method does not accept or filter by `areSystemMessages`.
- [ ] **Missing `deletionDateRange` filter parameter**: The Java method accepts `@Nullable DateRange deletionDateRange` and applies `.addBetweenIfNotNull(Message.Fields.DELETION_DATE, deletionDateRange)`. The Go method does not accept or filter by a deletion date range.
- [ ] **Date range uses `$gt`/`$lt` instead of `$gte`/`$lte`**: The Java code uses `addBetweenIfNotNull` which (per the `ToBson` in Go's own `DateRange`) maps to `$gte`/`$lte`. The Go `CountMessages` uses `$gt`/`$lt` (strict inequality), which will exclude messages delivered exactly at the boundary timestamps.

## CountUsersWhoSentMessage

- [ ] **Date range uses `$gt`/`$lt` instead of `$gte`/`$lte`**: Same as above — uses strict inequality instead of inclusive bounds, differing from the Java `addBetweenIfNotNull` behavior.
- [ ] **`Distinct` returns `[]interface{}` values which are counted by length**: The Java code uses `mongoClient.countDistinct(...)` which performs a server-side `count` of distinct values. The Go code uses `r.col.Distinct(ctx, "sid", filter)` which returns the actual distinct values, then counts them client-side with `int64(len(uniqueSenderIDs))`. This is functionally equivalent but may have performance implications for large datasets — not a correctness bug per se, but worth noting the approach differs.

## CountGroupsThatSentMessages

- [ ] **Date range uses `$gt`/`$lt` instead of `$gte`/`$lte`**: Same boundary inequality issue as other methods.

## CountSentMessages

- [ ] **Date range uses `$gt`/`$lt` instead of `$gte`/`$lte`**: Same boundary inequality issue.

## FindMessages

- [ ] **Missing `DateRange.NULL` special handling for `deletionDateRange`**: The Java code has special logic: `if (deletionDateRange == DateRange.NULL) { filter.eq(Message.Fields.DELETION_DATE, null); }`. This means when `DateRange.NULL` is passed, it filters for messages where deletion date is null (i.e., non-deleted messages). The Go code has no equivalent — it only checks if the range is non-nil and applies `$gt`/`$lt`, so there's no way to filter for non-deleted messages.
- [ ] **Missing `DateRange.NULL` special handling for `recallDateRange`**: Same issue — Java has `if (recallDateRange == DateRange.NULL) { filter.eq(Message.Fields.RECALL_DATE, null); }` to filter for non-recalled messages. The Go code lacks this entirely.
- [ ] **Date range uses `$gt`/`$lt` instead of `$gte`/`$lte`**: The Java `addBetweenIfNotNull` uses inclusive bounds. The Go code uses strict `$gt`/`$lt` comparisons.
- [ ] **Pagination uses `$gt`/`$lt` instead of `$gte`/`$lte`**: The skip calculation `int64(*page * *size)` is correct for offset-based pagination, matching the Java `paginateIfNotNull` behavior.
- [ ] **Missing `include`/field projection**: The Java code does not use `include` in `QueryOptions`, so it returns full message documents. The Go code also does not apply projection, so this is consistent — no bug here.

## FindMessageGroupId

- [ ] **Missing `isGroupMessage` projection field**: The Java code includes both `Message.Fields.TARGET_ID` and `Message.Fields.IS_GROUP_MESSAGE` in the projection. The Go code only projects `tid` (targetID) and omits `gm` (isGroupMessage). The return type is `*int64` (just the targetID) instead of `*po.Message` (which would carry both fields). This changes the method's contract — the Java version returns a `Message` PO with both `targetId` and `isGroupMessage` so the caller can verify it's a group message, while the Go version only returns the targetID.

## GetGroupConversationID

- [ ] **Byte order is reversed (little-endian vs big-endian)**: The Java `LongUtil.toBytes(groupId)` uses big-endian order (most significant byte first): `(byte) (groupId >> 56)` at index 0, down to `(byte) groupId` at index 7. The Go code iterates from `i = 7` down to `0`, placing the least significant byte at index 7 and working up. However, looking more carefully at the Go loop: `b[i] = byte(groupID); groupID >>= 8` — at i=7, it stores `byte(groupID)` (least significant byte) then shifts. At i=0, it stores the most significant byte. This produces big-endian output, matching Java. No bug here.

## GetPrivateConversationID

- [ ] **Byte order**: Same analysis as above — the Go code produces big-endian output matching Java's `toPrivateConversationId`. The logic correctly sorts min/max IDs. No bug here.

# MessageService.java
*Checked methods: isMessageRecipientOrSender(@NotNull Long messageId, @NotNull Long userId), authAndQueryCompleteMessages(Long requesterId, @Nullable Collection<Long> messageIds, @NotNull Boolean areGroupMessages, @Nullable Boolean areSystemMessages, @Nullable Set<Long> fromIds, @Nullable DateRange deliveryDateRange, @Nullable Integer maxCount, boolean ascending, boolean withTotal), queryMessage(@NotNull Long messageId), queryMessages(@Nullable Collection<Long> messageIds, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages, @Nullable Set<Long> senderIds, @Nullable Set<Long> targetIds, @Nullable DateRange deliveryDateRange, @Nullable DateRange deletionDateRange, @Nullable DateRange recallDateRange, @Nullable Integer page, @Nullable Integer size, @Nullable Boolean ascending), saveMessage(@Nullable Long messageId, @NotNull Long senderId, @Nullable byte[] senderIp, @NotNull Long targetId, @NotNull Boolean isGroupMessage, @NotNull Boolean isSystemMessage, @Nullable String text, @Nullable List<byte[]> records, @Nullable @Min(0), queryExpiredMessageIds(@NotNull Integer retentionPeriodHours), deleteExpiredMessages(@NotNull Integer retentionPeriodHours), deleteMessages(@Nullable Set<Long> messageIds, @Nullable Boolean deleteLogically), updateMessages(@Nullable Long senderId, @Nullable DeviceType senderDeviceType, @NotEmpty Set<Long> messageIds, @Nullable Boolean isSystemMessage, @Nullable String text, @Nullable List<byte[]> records, @Nullable @Min(0), hasPrivateMessage(Long senderId, Long targetId), countMessages(@Nullable Set<Long> messageIds, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages, @Nullable Set<Long> senderIds, @Nullable Set<Long> targetIds, @Nullable DateRange deliveryDateRange, @Nullable DateRange deletionDateRange), countUsersWhoSentMessage(@Nullable DateRange dateRange, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages), countGroupsThatSentMessages(@Nullable DateRange dateRange), countSentMessages(@Nullable DateRange dateRange, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages), countSentMessagesOnAverage(@Nullable DateRange dateRange, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages), authAndUpdateMessage(@NotNull Long senderId, @Nullable DeviceType senderDeviceType, @NotNull Long messageId, @Nullable String text, @Nullable List<byte[]> records, @Nullable @PastOrPresent Date recallDate), queryMessageRecipients(@NotNull Long messageId), authAndSaveMessage(boolean queryRecipientIds, @Nullable Boolean persist, @Nullable Long messageId, @NotNull Long senderId, @Nullable byte[] senderIp, @NotNull Long targetId, @NotNull Boolean isGroupMessage, @NotNull Boolean isSystemMessage, @Nullable String text, @Nullable List<byte[]> records, @Nullable @Min(0), saveMessage(boolean queryRecipientIds, @Nullable Boolean persist, @Nullable Long messageId, @NotNull Long senderId, @Nullable byte[] senderIp, @NotNull Long targetId, @NotNull Boolean isGroupMessage, @NotNull Boolean isSystemMessage, @Nullable String text, @Nullable List<byte[]> records, @Nullable @Min(0), authAndCloneAndSaveMessage(boolean queryRecipientIds, @NotNull Long requesterId, @Nullable byte[] requesterIp, @NotNull Long referenceId, @NotNull Boolean isGroupMessage, @NotNull Boolean isSystemMessage, @NotNull Long targetId), cloneAndSaveMessage(boolean queryRecipientIds, @NotNull Long senderId, @Nullable byte[] senderIp, @NotNull Long referenceId, @NotNull Boolean isGroupMessage, @NotNull Boolean isSystemMessage, @NotNull Long targetId), authAndSaveAndSendMessage(boolean send, @Nullable Boolean persist, @Nullable Long senderId, @Nullable DeviceType senderDeviceType, @Nullable byte[] senderIp, @Nullable Long messageId, @NotNull Boolean isGroupMessage, @NotNull Boolean isSystemMessage, @Nullable String text, @Nullable List<byte[]> records, @NotNull Long targetId, @Nullable @Min(0), saveAndSendMessage(boolean send, @Nullable Boolean persist, @Nullable Long senderId, @Nullable DeviceType senderDeviceType, @Nullable byte[] senderIp, @Nullable Long messageId, @NotNull Boolean isGroupMessage, @NotNull Boolean isSystemMessage, @Nullable String text, @Nullable List<byte[]> records, @NotNull Long targetId, @Nullable @Min(0), saveAndSendMessage(@Nullable Boolean persist, @Nullable Long senderId, @Nullable DeviceType senderDeviceType, @Nullable byte[] senderIp, @Nullable Long messageId, @NotNull Boolean isGroupMessage, @NotNull Boolean isSystemMessage, @Nullable String text, @Nullable List<byte[]> records, @NotNull Long targetId, @Nullable @Min(0), deleteGroupMessageSequenceIds(Set<Long> groupIds), deletePrivateMessageSequenceIds(Set<Long> userIds), fetchGroupMessageSequenceId(Long groupId), fetchPrivateMessageSequenceId(Long userId1, Long userId2)*

Now I have a thorough understanding of both codebases. Let me compile the comparison.

## isMessageRecipientOrSender

- [ ] **Missing sentMessageCache lookup**: The Java version first checks `sentMessageCache` for the message before querying the repository. The Go version has no cache lookup at all.
- [ ] **Missing senderId check for group messages**: In Java, when the message is a group message, the method checks if the user is a group member via `groupMemberService.isGroupMember(targetId, userId, false)`. The Go version does this but also includes a separate `msg.SenderID == userID` early return check at the top. However, in Java for group messages, the sender is NOT automatically considered a recipient - only membership is checked. The Go code returns `true` early if `msg.SenderID == userID` even for group messages, which is correct because the sender of a group message should be a valid recipient/sender. Actually this matches the Java private method `isMessageRecipientOrSender(userId, isGroupMessage, targetId, senderId)` which returns `targetId.equals(userId) || senderId.equals(userId)` for private messages and checks membership for group messages. So the sender check is only for private messages in Java. The Go code checks `msg.SenderID == userID` first regardless of message type, which means for group messages a sender would be considered a recipient even if they left the group. In Java, for group messages, only `groupMemberService.isGroupMember(targetId, userId, false)` is called, which would return false if the sender left the group.

## authAndQueryCompleteMessages

- [ ] **Missing areGroupMessages null check**: The Java version returns an error if `areGroupMessages` is null. The Go version accepts `*bool` for `isGroupMessage` but does not reject nil.
- [ ] **Missing maxCount defaulting/clamping logic**: Java defaults `maxCount` based on `withTotal` flag (using `defaultAvailableMessagesNumberWithTotal` or `defaultAvailableMessagesNumberWithoutTotal`), and clamps it with max values. The Go version uses `size` directly with no defaulting or clamping.
- [ ] **Missing withTotal parameter entirely**: The Go version does not have a `withTotal` parameter.
- [ ] **Missing messageIds parameter**: The Java version accepts `messageIds` but the Go version does not.
- [ ] **Missing areSystemMessages filtering**: The Go version accepts `areSystemMessages` but never passes it to the repository query. The repository `QueryMessages` method does not support `areSystemMessages` filtering.
- [ ] **Missing fromIds parameter**: Java passes `fromIds` to the repository query, but the Go version passes `senderIDs` instead, which is a different semantic. In Java, `fromIds` are used as `targetIds` for group messages (group IDs) and as `senderIds` for private messages.
- [ ] **Missing auth logic for group messages**: Java verifies the requester is a member of queried groups (either by fetching all joined groups or by checking membership in `fromIds`). The Go version has partial auth logic that checks `IsGroupMember` for each targetID, but does not implement the full logic with `queryUserJoinedGroupIds` when `fromIds` is empty.
- [ ] **Missing private message auth**: For private messages, Java sets `targetIds = Set.of(requesterId)` to ensure only messages addressed to the requester are returned. The Go version passes `nil` for targetIDs in private message queries, allowing retrieval of any private messages.
- [ ] **Missing deletionDateRange and recallDateRange filtering**: Java passes `DateRange.NULL` for `deletionDateRange` and `recalledMessageQueryDateRange` for `recallDateRange`. The Go version does not support these filters.

## queryMessage

- [ ] This method appears to match correctly.

## queryMessages

- [ ] **Missing conversation ID logic**: The Java version has extensive conversation ID computation logic based on `useConversationId`, converting sender/target IDs to conversation IDs for efficient querying. The Go version has no conversation ID logic at all.
- [ ] **Missing areSystemMessages parameter**: The Go version does not accept or pass `areSystemMessages` to the repository.
- [ ] **Missing deletionDateRange and recallDateRange parameters**: The Go version does not support `deletionDateRange` or `recallDateRange` filtering.
- [ ] **Missing page parameter**: The Go version does not support pagination via `page`.
- [ ] **Missing messageIds filtering**: The Go `QueryMessages` repository method does not filter by `messageIds`.

## saveMessage

- [ ] **Missing senderIp processing**: Java processes `senderIp` into IPv4 (`senderIpV4`) or IPv6 (`senderIpV6`) using `InetAddressUtil`. The Go version does not accept or process sender IP at all.
- [ ] **Missing isSystemMessage field**: The Go version does not set `IsSystemMessage` on the message. Java sets it.
- [ ] **Missing referenceId field**: The Go version does not set `ReferenceID` on the message. Java sets it.
- [ ] **Missing conversationId computation**: Java computes `conversationId` using `MessageRepository.getGroupConversationId()` or `getPrivateConversationId()` based on `useConversationId`. The Go version does not compute or set `ConversationID`.
- [ ] **Missing conditional sequence ID generation**: Java only generates sequence IDs when `useSequenceIdForGroupConversation` or `useSequenceIdForPrivateConversation` is true. The Go version always generates a sequence ID.
- [ ] **Missing persistRecord check**: Java nulls out `records` if `!persistRecord && !isSystemMessage`. The Go version always persists records.
- [ ] **Missing persistPreMessageId check**: Java nulls out `preMessageId` if `!persistPreMessageId`. The Go version always includes `preMessageId`.
- [ ] **Missing persistSenderIp check**: Java nulls out `senderIp` if `!persistSenderIp`. The Go version has no sender IP at all.
- [ ] **Missing updateReadDateAfterMessageSent logic**: Java optionally upserts conversation read date after saving a message. The Go version has no such logic.
- [ ] **Missing text validation (maxTextLimit)**: Java validates `text` length against `maxTextLimit`. The Go version has no validation.
- [ ] **Missing records validation (maxRecordsSize)**: Java validates total records size against `maxRecordsSize`. The Go version has no validation.
- [ ] **Missing burnAfter validation**: Java validates `burnAfter >= 0`. The Go version has no validation.
- [ ] **Missing messageId null -> generation**: Java generates a new ID via `node.nextLargeGapId()` if `messageId` is null. The Go version always generates a new ID.
- [ ] **Missing timeType handling**: Java checks `timeType == TimeType.LOCAL_SERVER_TIME` before overriding deliveryDate. The Go version always uses current time if deliveryDate is nil.

## queryExpiredMessageIds

- [ ] This method appears to match functionally (simple delegation to repo).

## deleteExpiredMessages

- [ ] **Missing plugin manager hook**: Java invokes `ExpiredMessageDeletionNotifier` plugin extensions to filter which messages to delete. The Go version deletes all expired messages without plugin notification.
- [ ] **Missing fetch-messages-before-delete step**: Java first fetches the full message objects for the plugin hook. The Go version only uses IDs.

## deleteMessages

- [ ] **Missing deleteLogically null -> default behavior**: Java defaults `deleteLogically` to `deleteMessageLogicallyByDefault` when null. The Go version requires explicit non-nil to trigger logical deletion.

## updateMessages

- [ ] **Missing senderIp parameter and processing**: Java accepts `senderIp` and processes it into IPv4/IPv6. The Go version does not accept sender IP.
- [ ] **Missing text validation (maxTextLimit)**: Java validates text length. The Go version does not.
- [ ] **Missing burnAfter validation**: Java validates `burnAfter >= 0`. The Go version does not.
- [ ] **Missing recallDate validation**: Java validates `recallDate` is past or present. The Go version does not.
- [ ] **Missing records validation (maxRecordsSize)**: Java validates records size. The Go version does not.
- [ ] **Missing early return when all update fields are null**: Java returns `ACKNOWLEDGED_UPDATE_RESULT` if all update fields (isSystemMessage, text, records, burnAfter, recallDate, senderIp) are null. The Go version only checks if the built `set` map is empty.
- [ ] **Missing recall notification logic**: When `recallDate` is set, Java re-fetches messages and sends recall notification messages to recipients. The Go version has no such notification logic.
- [ ] **Missing senderDeviceType parameter usage**: The Go version accepts `senderDeviceType` but never uses it.

## hasPrivateMessage

- [ ] This method appears to match functionally.

## countMessages

- [ ] **Missing messageIds parameter**: Java accepts `messageIds` and passes it to the repository. The Go version does not support `messageIds` filtering.
- [ ] **Missing areSystemMessages parameter**: Java accepts and passes `areSystemMessages`. The Go version does not.
- [ ] **Missing deletionDateRange parameter**: Java supports `deletionDateRange`. The Go version does not.

## countUsersWhoSentMessage

- [ ] **This method appears to match functionally** (delegates to repository with date range and filters).

## countGroupsThatSentMessages

- [ ] This method appears to match functionally.

## countSentMessages

- [ ] This method appears to match functionally.

## countSentMessagesOnAverage

- [ ] **Incorrect return value when totalUsers is 0**: Java returns `Long.MAX_VALUE` when totalUsers is 0 (and totalMessages > 0), indicating "infinite" average. The Go version returns 0.

## authAndUpdateMessage

- [ ] **Missing allowEditMessageBySender check**: Java checks if message editing is enabled before allowing updates. The Go version does not.
- [ ] **Missing allowRecallMessage check**: Java checks if recall is enabled before allowing recalls. The Go version does not.
- [ ] **Missing group type check for edit**: Java verifies the group type allows message editing for group messages. The Go version does not.
- [ ] **Missing recall duration check**: Java checks `availableRecallDurationMillis` to ensure the message is within the recall window. The Go version has no time-based check.
- [ ] **Missing group existence check for recall**: Java verifies the group still exists and is active for group message recalls. The Go version does not.
- [ ] **Missing sender verification from DB**: Java uses `findIsGroupMessageAndTargetId` (which also checks senderId) to verify the sender. The Go version loads the full message and checks `msg.SenderID != senderID`.
- [ ] **Missing senderDeviceType usage**: The Go version accepts `senderDeviceType` but never uses it.
- [ ] **Missing recall notification sending**: Java sends recall notification messages (with `BuiltinSystemMessageType.RECALL_MESSAGE_BYTES`) to recipients after a successful recall. The Go version has no such notification logic.
- [ ] **Missing sentMessageCache lookup**: Java checks the cache first for the message. The Go version always queries the DB.

## queryMessageRecipients

- [ ] **Missing sender exclusion for group messages**: For group messages, the Java version calls `queryGroupMemberIds(groupId, true)` which likely returns all members. However, the Go version calls `FindGroupMemberIDs` which may or may not include the sender. The semantic difference depends on the underlying implementation, but the Java code does not explicitly exclude the sender in this specific method.
- [ ] **Missing isGroupMessage field check**: Java uses `messageRepository.findMessageGroupId` which returns an object with `isGroupMessage` and `groupId`. The Go version loads the full message via `QueryMessage`, then checks `IsGroupMessage`. However, the Go version's `FindMessageGroupId` in the repo returns only `TargetID` without `IsGroupMessage`. The service uses `QueryMessage` which loads the full message, so this is functionally equivalent but less efficient.

## authAndSaveMessage

- [ ] **Missing user permission check (userService.isAllowedToSendMessageToTarget)**: Java calls `userService.isAllowedToSendMessageToTarget()` to verify the sender has permission. The Go version uses a simpler `auth()` method that only checks friendship/group membership, not the full permission system.
- [ ] **Missing validation**: Java validates text length, records length, and deliveryDate. The Go version has no validation.
- [ ] **Missing persist flag handling**: Java respects a `persist` flag (defaulting to `persistMessage` property). The Go version always persists.
- [ ] **Missing queryRecipientIds logic**: Java optionally queries recipient IDs and returns them with the message. The Go version returns only the message.
- [ ] **Missing sentMessageCache caching**: Java caches the sent message. The Go version does not cache.
- [ ] **Missing referenceId field**: Java passes `referenceId` to `saveMessage`. The Go version does not set `ReferenceID`.
- [ ] **Missing isSystemMessage field**: The Go version does not set `IsSystemMessage`.
- [ ] **Missing senderIp handling**: The Go version does not accept or process sender IP.

## saveMessage (the overloaded version with queryRecipientIds/persist)

- [ ] **This is a completely different method in Go**: The Java version is a service-level method that delegates to `saveMessage0`, handling `queryRecipientIds`, `persist` flag, and returning `MessageAndRecipientIds`. The Go version is a simple save that always persists and returns the message. All the `saveMessage0` logic (recipient ID querying, persist flag, caching) is missing.

## authAndCloneAndSaveMessage

- [ ] **Missing referenceId handling**: Java passes `referenceId` to `authAndSaveMessage` and generates a new message ID. The Go version calls `CloneAndSaveMessage` which does not set `ReferenceID`.
- [ ] **Missing isSystemMessage, text, records, burnAfter from original message**: Java passes the original message's `text`, `records`, and `burnAfter` to `authAndSaveMessage`. The Go version passes these via `CloneAndSaveMessage` but does not pass `isSystemMessage` or handle `burnAfter` and `deliveryDate` correctly from the reference message.
- [ ] **Missing switchIfEmpty error**: Java returns `ERROR_NOT_MESSAGE_RECIPIENT_OR_SENDER_TO_FORWARD_MESSAGE` when the reference message is not found (switchIfEmpty). The Go version will propagate the error from `QueryMessage` which is a different error type.

## cloneAndSaveMessage

- [ ] **Missing senderIp parameter**: Java passes `senderIp` from the caller. The Go version does not accept sender IP.
- [ ] **Missing new messageId generation**: Java generates a new message ID via `node.nextLargeGapId()`. The Go version delegates to `SaveMessage` which generates its own ID, so this is functionally equivalent.
- [ ] **Missing deliveryDate from original message**: Java passes `message.getDeliveryDate()` from the reference message. The Go version passes `nil` for deliveryDate, meaning a new delivery date is generated.
- [ ] **Missing isSystemMessage handling**: The Go version accepts `isSystemMessage` but never uses it (SaveMessage doesn't set IsSystemMessage).
- [ ] **Missing referenceId**: Java passes `referenceId` to `saveMessage`. The Go version does not set `ReferenceID` on the cloned message.

## authAndSaveAndSendMessage

- [ ] **Missing auth parameter differentiation**: Java's `saveAndSendMessage0` handles both auth and non-auth paths. The Go version always authenticates (via `AuthAndSaveMessage`).
- [ ] **Missing referenceId handling**: Java routes to `authAndCloneAndSaveMessage` when `referenceId` is present. The Go version does not support `referenceId`.
- [ ] **Missing validation**: Java validates isGroupMessage, isSystemMessage, targetId, burnAfter, text/records, text length, and records length. The Go version has no validation.
- [ ] **Missing senderId null handling**: Java defaults senderId to `ADMIN_REQUESTER_ID` for system messages. The Go version does not handle null senderId.
- [ ] **Missing sentMessageCounter**: Java increments a metrics counter. The Go version does not.
- [ ] **Missing sentMessageCache caching**: Java caches the sent message. The Go version does not.
- [ ] **Missing notification to recipients**: Java calls `sendMessage()` which builds a `TurmsNotification` and sends it via `outboundMessageManager.forwardNotification()`. The Go version uses a simpler `outboundDelivery.Deliver()` which has different semantics.
- [ ] **Missing senderDeviceType usage**: Java uses senderDeviceType for session exclusion in notifications. The Go version does not use it.
- [ ] **Missing preMessageId parameter**: The Go version does not accept `preMessageId`.

## saveAndSendMessage (with `send` parameter)

- [ ] **Missing referenceId handling**: Java routes to `cloneAndSaveMessage` when `referenceId` is present. The Go version does not support `referenceId`.
- [ ] **Missing queryRecipientIds logic**: Java returns `MessageAndRecipientIds` with recipient IDs. The Go version returns just the message.
- [ ] **Missing persist flag handling**: Java respects the `persist` flag. The Go version conditionally persists based on the `persist` bool parameter but does not use a property-based default.
- [ ] **Missing isSystemMessage and isGroupMessage validation and field assignment**: The Go version accepts these but `SaveMessage` does not set `IsSystemMessage`.
- [ ] **Missing sentMessageCounter**: No metrics counter increment.
- [ ] **Missing sentMessageCache caching**: No cache update.
- [ ] **Missing notification to multiple recipients**: Java builds a notification and sends to all recipients via `outboundMessageManager`. The Go version uses a simple single-target `Deliver()`.
- [ ] **Missing senderDeviceType for notification**: Java uses senderDeviceType for session exclusion. The Go version does not use it.
- [ ] **Missing preMessageId parameter**: The Go version does not accept `preMessageId`.

## saveAndSendMessage (without `send` parameter - the 3rd overload)

- [ ] **This overload is completely missing in Go**: The Java version (lines 1493-1558) always sends (hardcodes `send=true` equivalent), validates, and calls `saveMessage` (the service-level version with queryRecipientIds). The Go version has no equivalent.

## deleteGroupMessageSequenceIds

- [ ] **No-op implementation**: The Java version executes Redis `HDEL` commands against `KEY_GROUP_MESSAGE_SEQUENCE_ID_BUFFER`. The Go version returns `nil` without doing anything.

## deletePrivateMessageSequenceIds

- [ ] **No-op implementation**: The Java version executes a Lua script against Redis. The Go version returns `nil` without doing anything.

## fetchGroupMessageSequenceId

- [ ] **Incorrectly increments instead of fetching**: The Go version calls `s.seqGen.NextGroupMessageSequenceId()` which likely increments the sequence. Java's `fetchGroupMessageSequenceId` uses `HINCR` which also increments. However, the Java version returns `Mono.empty()` when `redisClientManager == null`, while the Go version will still try to generate a sequence. The comment in Go even notes "usually a fetch, here we might accidentally increment" - indicating the author was aware this might be wrong.
- [ ] **Missing nil check**: Java returns empty/nil when Redis is not configured. The Go version always attempts the operation.

## fetchPrivateMessageSequenceId

- [ ] **Missing userId2 parameter usage**: Java uses both `userId1` and `userId2` to compute the sequence ID key (ordering them so the smaller ID is first). The Go version only uses `userId1`, ignoring `userId2` entirely.
- [ ] **Incorrect sequence key**: Java uses an XOR-based key (`userId1 ^ userId2`) and a Lua script. The Go version uses only `userId1`.

Here is the consolidated bug report:

## IsMessageRecipientOrSender

- [ ] **Sender check bypasses group membership for group messages**: The Go code returns `true` immediately if `msg.SenderID == userID` (line 314), even for group messages. In Java, for group messages, only `groupMemberService.isGroupMember(targetId, userId, false)` is called — a sender who left the group would NOT be considered a recipient/sender. The Go code incorrectly grants access to former group members who were the original sender.
- [ ] **Missing sentMessageCache lookup**: Java first checks the in-memory `sentMessageCache` before querying MongoDB. The Go version always queries the database.

## AuthAndQueryCompleteMessages

- [ ] **Missing areGroupMessages null rejection**: Java returns an error if `areGroupMessages` is null. Go accepts nil without error.
- [ ] **Missing maxCount defaulting/clamping**: Java defaults `maxCount` based on `withTotal` and clamps it with configurable max values. Go passes `size` directly with no defaulting or clamping.
- [ ] **Missing withTotal parameter**: Go does not have a `withTotal` parameter at all.
- [ ] **Missing messageIds parameter**: Go does not accept `messageIds`.
- [ ] **Missing areSystemMessages filtering**: Go accepts `areSystemMessages` but never passes it to the repository query.
- [ ] **Missing group auth logic (queryUserJoinedGroupIds)**: When `fromIds` is empty and querying group messages, Java fetches all joined group IDs for the requester. Go does not implement this path.
- [ ] **Missing group membership verification for fromIds**: Java verifies the requester is a member of each group in `fromIds` via `findExistentMemberGroupIds`. Go checks `IsGroupMember` per targetID but with different error handling and logic.
- [ ] **Missing private message requester constraint**: For private messages, Java sets `targetIds = Set.of(requesterId)` to restrict results. Go passes `nil` for targetIDs, returning any private messages regardless of requester involvement.
- [ ] **Missing deletionDateRange and recallDateRange**: Java passes `DateRange.NULL` for deletion and `recalledMessageQueryDateRange` for recall. Go supports neither.

## QueryMessages (service-level)

- [ ] **Missing conversation ID computation**: Java has extensive logic to compute conversation IDs from sender/target IDs when `useConversationId` is true. Go has no conversation ID logic.
- [ ] **Missing areSystemMessages parameter**: Go does not accept or filter by `areSystemMessages`.
- [ ] **Missing deletionDateRange and recallDateRange**: Go does not support these filters.
- [ ] **Missing page parameter**: Go does not support pagination.
- [ ] **Missing messageIds filtering**: Go does not filter by `messageIds`.

## QueryMessages (repository-level)

- [ ] **Missing areSystemMessages filter in QueryMessages**: The `QueryMessages` repo method does not accept `areSystemMessages`, unlike Java's `findMessages`.
- [ ] **Missing messageIds, conversationIds, deletionDateRange, recallDateRange, page**: The simplified `QueryMessages` repo method lacks many parameters that Java's `findMessages` supports.

## SaveMessage

- [ ] **Missing senderIp processing**: Java processes sender IP into IPv4/IPv6. Go does not accept sender IP.
- [ ] **Missing IsSystemMessage field assignment**: Go does not set `IsSystemMessage` on the message. Java always sets it.
- [ ] **Missing ReferenceID field assignment**: Go does not set `ReferenceID`. Java sets it from the parameter.
- [ ] **Missing ConversationID computation**: Java computes and sets conversation ID based on `useConversationId`. Go does not set `ConversationID`.
- [ ] **Always generates sequence ID**: Java only generates sequence IDs when configured. Go always generates one.
- [ ] **Missing persistRecord check**: Java nulls records if `!persistRecord && !isSystemMessage`. Go always persists records.
- [ ] **Missing persistPreMessageId check**: Java nulls preMessageId if not configured. Go always includes it.
- [ ] **Missing updateReadDateAfterMessageSent logic**: Java optionally upserts conversation read date. Go has no such logic.
- [ ] **Missing text, records, burnAfter validation**: Java validates `maxTextLimit`, `maxRecordsSize`, and `burnAfter >= 0`. Go has no validation.
- [ ] **Missing timeType handling for deliveryDate**: Java checks `timeType == LOCAL_SERVER_TIME` before overriding. Go always overrides if nil.

## DeleteExpiredMessages

- [ ] **Missing plugin manager hook**: Java invokes `ExpiredMessageDeletionNotifier` extensions. Go deletes all expired messages without plugin notification.
- [ ] **Missing fetch-messages-before-delete**: Java fetches full message objects for the plugin hook. Go only uses IDs.

## DeleteMessages

- [ ] **Missing default for null deleteLogically**: Java defaults to `deleteMessageLogicallyByDefault` when null. Go requires explicit non-nil to trigger logical deletion.

## UpdateMessages

- [ ] **Missing senderIp parameter and processing**: Java accepts and processes sender IP. Go does not.
- [ ] **Missing text, records, burnAfter, recallDate validation**: Java validates all update fields. Go has no validation.
- [ ] **Missing recall notification logic**: When `recallDate` is set, Java re-fetches messages and sends recall notifications to recipients. Go has no notification logic.
- [ ] **Missing early return for all-null update fields**: Java returns immediately if all fields are null. Go builds and executes an empty update (though the repo returns early for empty `set`).

## CountMessages

- [ ] **Missing messageIds parameter**: Java accepts and passes `messageIds`. Go does not.
- [ ] **Missing areSystemMessages parameter**: Java accepts `areSystemMessages`. Go does not.
- [ ] **Missing deletionDateRange parameter**: Java supports `deletionDateRange`. Go does not.

## CountSentMessagesOnAverage

- [ ] **Incorrect return when totalUsers is 0**: Java returns `Long.MAX_VALUE` when `totalUsers == 0` (and `totalMessages > 0`). Go returns `0`.

## AuthAndUpdateMessage

- [ ] **Missing allowEditMessageBySender check**: Java verifies editing is enabled. Go does not.
- [ ] **Missing allowRecallMessage check**: Java verifies recall is enabled. Go does not.
- [ ] **Missing group type message editability check**: Java checks `groupType.getMessageEditable()` for group messages. Go does not.
- [ ] **Missing recall duration check**: Java checks `availableRecallDurationMillis`. Go has no time-based check.
- [ ] **Missing group existence check for recall**: Java verifies the group is active for group message recalls. Go does not.
- [ ] **Missing recall notification sending**: Java sends recall notification messages with `BuiltinSystemMessageType.RECALL_MESSAGE_BYTES`. Go has no notification.
- [ ] **Missing sentMessageCache lookup**: Java checks cache first. Go always queries DB.

## QueryMessageRecipients

- [ ] **Missing proper isGroupMessage-based dispatch from repository**: Java uses `findMessageGroupId` which returns `isGroupMessage` and `groupId` together. Go loads the full message, which works but is less efficient. Functionally equivalent.

## AuthAndSaveMessage

- [ ] **Missing full permission check**: Java calls `userService.isAllowedToSendMessageToTarget()`. Go uses a simpler `auth()` checking only friendship/group membership.
- [ ] **Missing validation**: Java validates text length, records size, deliveryDate. Go has no validation.
- [ ] **Missing persist flag handling**: Java respects `persist` flag (defaulting to property). Go always persists.
- [ ] **Missing queryRecipientIds logic**: Java returns `MessageAndRecipientIds` with recipient IDs. Go returns only the message.
- [ ] **Missing sentMessageCache caching**: Java caches sent messages. Go does not.
- [ ] **Missing ReferenceID, IsSystemMessage, senderIp fields**: Go does not set these on the message.

## SaveMessage (overloaded with queryRecipientIds/persist)

- [ ] **Missing saveMessage0 delegation logic**: Java's `saveMessage(queryRecipientIds, persist, ...)` delegates to `saveMessage0` which handles recipient querying, persist flag, caching, and returns `MessageAndRecipientIds`. Go's version is a simple save returning only the message.
- [ ] **Missing queryRecipientIds handling**: Java queries group member IDs (excluding sender) for group messages. Go does not query recipients.
- [ ] **Missing persist flag with property default**: Java defaults `persist` to `persistMessage` property. Go always persists.

## AuthAndCloneAndSaveMessage

- [ ] **Missing referenceId on cloned message**: Java passes `referenceId` to `authAndSaveMessage`. Go's `CloneAndSaveMessage` does not set `ReferenceID`.
- [ ] **Missing deliveryDate from original message**: Java passes `null` for deliveryDate in `authAndSaveMessage`. Go passes `nil` via `CloneAndSaveMessage`, which is equivalent.
- [ ] **Missing switchIfEmpty error**: Java returns `NOT_MESSAGE_RECIPIENT_OR_SENDER_TO_FORWARD_MESSAGE` when the reference message is not found. Go propagates a generic DB error.

## CloneAndSaveMessage

- [ ] **Missing deliveryDate from reference message**: Java passes `message.getDeliveryDate()`. Go passes `nil`, causing a new delivery date to be generated.
- [ ] **Missing referenceId on cloned message**: Java passes `referenceId`. Go does not set `ReferenceID`.
- [ ] **Missing isSystemMessage handling**: Go accepts but never uses `isSystemMessage`.
- [ ] **Missing senderIp parameter**: Java passes sender IP. Go does not accept it.

## AuthAndSaveAndSendMessage

- [ ] **Missing referenceId routing**: Java routes to `authAndCloneAndSaveMessage` when `referenceId` is present. Go does not support `referenceId`.
- [ ] **Missing validation**: Java validates isGroupMessage, isSystemMessage, targetId, burnAfter, text/records, text length, records length. Go has no validation.
- [ ] **Missing senderId null handling**: Java defaults to `ADMIN_REQUESTER_ID` for system messages. Go does not handle null senderId.
- [ ] **Missing sentMessageCounter**: Java increments a metrics counter. Go does not.
- [ ] **Missing sentMessageCache caching**: Java caches. Go does not.
- [ ] **Missing proper notification to recipients**: Java builds a `TurmsNotification` with `CreateMessageRequest` and sends via `outboundMessageManager.forwardNotification()`. Go uses a simpler `outboundDelivery.Deliver()`.
- [ ] **Missing senderDeviceType usage**: Java uses it for session exclusion. Go does not use it.

## SaveAndSendMessage (with `send` parameter)

- [ ] **Missing referenceId routing**: Java routes to `cloneAndSaveMessage` when `referenceId` is present. Go does not support `referenceId`.
- [ ] **Missing queryRecipientIds and proper persist handling**: Java returns `MessageAndRecipientIds`. Go returns only the message.
- [ ] **Missing isSystemMessage field in persisted message**: Go's `SaveMessage` does not set `IsSystemMessage`.
- [ ] **Missing sentMessageCounter**: No metrics counter.
- [ ] **Missing sentMessageCache caching**: No cache update.
- [ ] **Missing proper multi-recipient notification**: Go uses single-target `Deliver()` instead of Java's `forwardNotification()` to all recipients.

## SaveAndSendMessage (3rd overload, always sends)

- [ ] **Entire method missing**: Java has a 3rd overload (lines 1493-1558) that hardcodes `queryRecipientIds=true`, always sends the message, validates, handles null senderId, and calls the service-level `saveMessage`. Go has no equivalent.

## DeleteGroupMessageSequenceIds

- [ ] **No-op implementation**: Java executes Redis `HDEL` commands. Go returns `nil` without doing anything.

## DeletePrivateMessageSequenceIds

- [ ] **No-op implementation**: Java executes a Lua script against Redis. Go returns `nil` without doing anything.

## FetchGroupMessageSequenceId

- [ ] **Potential incorrect behavior when Redis unavailable**: Java returns `Mono.empty()` when `redisClientManager == null`. Go will still attempt generation, potentially erroring differently.
- [ ] **Comment indicates author uncertainty**: The Go code has a comment "usually a fetch, here we might accidentally increment" suggesting the implementation may be incorrect.

## FetchPrivateMessageSequenceId

- [ ] **Missing userId2 in sequence key computation**: Java orders userId1/userId2 (smaller first) and uses both to compute the sequence key via a Lua script. Go only uses `userId1`, completely ignoring `userId2`.
- [ ] **Incorrect sequence key**: Java uses `userId1 ^ userId2` (XOR) for key distribution and ordered keys for the Lua script. Go uses only `userId1`.

# StatisticsService.java
*Checked methods: countOnlineUsers()*

Now I understand the full picture. Let me compare the two implementations:

**Java `countOnlineUsers()` (StatisticsService, line 90-98):**
- Creates an RPC request (`CountOnlineUsersRequest`)
- Sends it to **all other cluster members** via `node.getRpcService().requestResponsesFromOtherMembers(request, true)`
- On `MEMBER_NOT_FOUND` error, falls back to `Mono.just(0)`
- Sums all responses with `MathFlux.sumInt(responses)`
- Returns a `Mono<Integer>` (async result)

**Go `CountOnlineUsers()` (SessionService, line 200-202):**
- Only counts **local** sessions via `s.shardedMap.CountOnlineUsers()` which iterates local shards and counts entries
- Returns a synchronous `int`
- Does **not** send any RPC to other cluster members

## countOnlineUsers

- [ ] **Missing cluster-wide aggregation**: The Java `countOnlineUsers()` sends RPC requests (`CountOnlineUsersRequest`) to all other cluster members via `requestResponsesFromOtherMembers()` and sums the responses from every node in the cluster. The Go version only counts local sessions via `s.shardedMap.CountOnlineUsers()`, completely missing the distributed/cluster-wide count. It never communicates with other nodes.
- [ ] **Missing RPC error handling for MEMBER_NOT_FOUND**: The Java code uses `onErrorResume` to catch `RpcErrorCode.MEMBER_NOT_FOUND` and fallback to `Mono.just(0)`, ensuring resilience when cluster members are unavailable. The Go version has no such error handling because it performs no RPC at all.
- [ ] **Missing RPC request type and dispatch**: The Java code relies on `CountOnlineUsersRequest` being handled by gateway nodes. The Go code has no corresponding RPC handler registration for counting online users across the cluster (only `SetUserOfflineRequest` is registered in the constructor at service.go:79).

# StorageServiceController.java
*Checked methods: handleDeleteResourceRequest(), handleQueryResourceUploadInfoRequest(), handleQueryResourceDownloadInfoRequest(), handleUpdateMessageAttachmentInfoRequest(), handleQueryMessageAttachmentInfosRequest()*

Now I have all the information needed for a thorough comparison. Let me compile the findings.

## HandleDeleteResourceRequest

- [ ] **Missing `customAttributes` parameter**: Java passes `customAttributes` (from `request.getCustomAttributesList()`) to `storageService.deleteResource()`. The Go controller does not pass any equivalent of `customAttributes`/`Extra` from the request to the service. The Go `DeleteResource` service method also doesn't accept a `customAttributes` parameter, so this is missing at both the controller and service layers.

- [ ] **Loss of `resourceIdNum` vs `resourceIdStr` distinction**: Java passes both `resourceIdNum` (nullable Long) and `resourceIdStr` (nullable String) separately to the service, which uses them differently depending on the resource type (e.g., `deleteGroupProfilePicture` requires `resourceIdNum`). The Go controller converts `IdNum` to a string via `strconv.FormatInt`, losing the numeric identity. The Java service uses `resourceIdNum` as a Long directly for GROUP_PROFILE_PICTURE and MESSAGE_ATTACHMENT scenarios.

## HandleQueryResourceUploadInfoRequest

- [ ] **Missing `resourceIdNum` parameter**: Java passes `resourceIdNum` to `storageService.queryResourceUploadInfo()`, where it's critical for routing logic — when `resourceType == MESSAGE_ATTACHMENT`, the service uses `resourceIdNum` to determine whether to query in private conversations, group conversations, or general. The Go controller completely omits `resourceIdNum` from the upload request and does not pass it to the service.

- [ ] **Missing `customAttributes` parameter**: Java passes `customAttributes` to the service. The Go controller does not pass any equivalent from the request's `Extra` map.

- [ ] **Incorrect response format**: Java's `queryResourceUploadInfo` returns `Map<String, String>` which is converted via `ProtoModelConvertor.toList()` into a flat `List<String>` (key-value pairs interleaved: key1, value1, key2, value2...). The Go controller assumes the service returns a single URL string and wraps it in `[]string{url}`. The Java version returns multiple key-value pairs from a map, not a single string.

- [ ] **Hardcoded `maxSize` of 0**: The Go controller passes `0` for `maxSize` to `QueryResourceUploadInfo`, but this parameter doesn't exist in the Java version at all — the Java method signature has `(requesterId, resourceType, resourceIdNum, resourceName, resourceMediaType, customAttributes)`. This is a fabricated parameter.

## HandleQueryResourceDownloadInfoRequest

- [ ] **Missing `customAttributes` parameter**: Java passes `customAttributes` to the service. The Go controller does not extract or pass `Extra` from the request.

- [ ] **Incorrect response format**: Same as upload info — Java's `queryResourceDownloadInfo` returns `Map<String, String>`, converted via `ProtoModelConvertor.toList()` to a flat list of interleaved key-value strings. The Go controller assumes a single URL string is returned and wraps it as `[]string{url}`, producing incorrect output format.

- [ ] **Loss of `resourceIdNum` vs `resourceIdStr` distinction**: Java passes both `resourceIdNum` and `resourceIdStr` as separate nullable parameters to the service. The Go controller converts `IdNum` to a string, losing the original numeric type distinction needed for USER_PROFILE_PICTURE and GROUP_PROFILE_PICTURE download routes.

## HandleUpdateMessageAttachmentInfoRequest

- [ ] **Behavioral difference on no-operation case**: When none of the share/unshare fields are set, Java returns `Mono.empty()` which is then handled by `.thenReturn(RequestHandlerResult.OK)` — the result is still a success notification. In Go, when none of the conditions match, `err` remains `nil` (its zero value), and the function still returns a success notification. This appears functionally equivalent, so no bug here.

## HandleQueryMessageAttachmentInfosRequest

- [ ] **Missing `name` field in proto conversion**: The Java `storageResourceInfo2proto` method does NOT set the `name` field on the proto `StorageResourceInfo` (it only sets `idNum`, `idStr`, `mediaType`, `uploaderId`, `creationDate`). However, the Go controller sets `Name: &name` unconditionally for every `StorageResourceInfo` proto, always including the name field in the response. This diverges from the Java behavior where the name is omitted from the proto output.

- [ ] **Always setting `name` and `mediaType` as non-nil pointers**: The Go code creates `name` and `mediaType` as always-set `*string` pointers (even when empty string), meaning the proto will always include these fields. In Java, the `storageResourceInfo2proto` only sets `mediaType` (always present via builder) but never sets `name`. This means the Go version always sends the name field while Java never does.

- [ ] **`CreationDate` always set in proto**: The Go code conditionally sets `CreationDate` only when `!info.CreationDate.IsZero()`. The Java code always calls `.setCreationDate(info.creationDate().getTime())` unconditionally. If a BO has a zero/epoch creation date, Java would still set it (value 0) while Go would omit it from the proto.

## HandleDeleteResourceRequest

- [ ] **Missing `customAttributes` parameter**: Java passes `customAttributes` (from `request.getCustomAttributesList()`) to `storageService.deleteResource()`. The Go controller does not pass any equivalent of `customAttributes`/`Extra` from the request to the service. The Go `DeleteResource` service method also doesn't accept a `customAttributes` parameter, so this is missing at both the controller and service layers.

- [ ] **Loss of `resourceIdNum` vs `resourceIdStr` distinction**: Java passes both `resourceIdNum` (nullable Long) and `resourceIdStr` (nullable String) separately to the service, which uses them differently depending on the resource type (e.g., `deleteGroupProfilePicture` requires `resourceIdNum`). The Go controller converts `IdNum` to a string via `strconv.FormatInt`, losing the numeric identity. The Java service uses `resourceIdNum` as a Long directly for GROUP_PROFILE_PICTURE and MESSAGE_ATTACHMENT scenarios.

## HandleQueryResourceUploadInfoRequest

- [ ] **Missing `resourceIdNum` parameter**: Java passes `resourceIdNum` to `storageService.queryResourceUploadInfo()`, where it's critical for routing logic — when `resourceType == MESSAGE_ATTACHMENT`, the service uses `resourceIdNum` to determine whether to query in private conversations, group conversations, or general. The Go controller completely omits `resourceIdNum` from the upload request and does not pass it to the service.

- [ ] **Missing `customAttributes` parameter**: Java passes `customAttributes` to the service. The Go controller does not pass any equivalent from the request's `Extra` map.

- [ ] **Incorrect response format**: Java's `queryResourceUploadInfo` returns `Map<String, String>` which is converted via `ProtoModelConvertor.toList()` into a flat `List<String>` (key-value pairs interleaved: key1, value1, key2, value2...). The Go controller assumes the service returns a single URL string and wraps it in `[]string{url}`. The Java version returns multiple key-value pairs from a map, not a single string.

- [ ] **Hardcoded `maxSize` of 0**: The Go controller passes `0` for `maxSize` to `QueryResourceUploadInfo`, but this parameter doesn't exist in the Java version at all — the Java method signature has `(requesterId, resourceType, resourceIdNum, resourceName, resourceMediaType, customAttributes)`. This is a fabricated parameter.

## HandleQueryResourceDownloadInfoRequest

- [ ] **Missing `customAttributes` parameter**: Java passes `customAttributes` to the service. The Go controller does not extract or pass `Extra` from the request.

- [ ] **Incorrect response format**: Same as upload info — Java's `queryResourceDownloadInfo` returns `Map<String, String>`, converted via `ProtoModelConvertor.toList()` to a flat list of interleaved key-value strings. The Go controller assumes a single URL string is returned and wraps it as `[]string{url}`, producing incorrect output format.

- [ ] **Loss of `resourceIdNum` vs `resourceIdStr` distinction**: Java passes both `resourceIdNum` and `resourceIdStr` as separate nullable parameters to the service. The Go controller converts `IdNum` to a string, losing the original numeric type distinction needed for USER_PROFILE_PICTURE and GROUP_PROFILE_PICTURE download routes.

## HandleUpdateMessageAttachmentInfoRequest

No bugs found — the logic correctly mirrors the Java if/else chain for share/unshare operations.

## HandleQueryMessageAttachmentInfosRequest

- [ ] **Extra `name` field set in proto conversion that Java omits**: The Java `storageResourceInfo2proto` method never sets the `name` field on the proto output — it only sets `idNum`, `idStr`, `mediaType`, `uploaderId`, and `creationDate`. The Go controller always sets `Name: &name` for every proto `StorageResourceInfo`, including an extra field that Java does not send.

- [ ] **`CreationDate` conditionally omitted vs always set**: The Go code only sets `protoInfo.CreationDate` when `!info.CreationDate.IsZero()`, while Java always calls `.setCreationDate(info.creationDate().getTime())` unconditionally. If a zero/epoch creation date occurs, Go omits the field from the proto while Java would set it to `0`.

# StorageService.java
*Checked methods: deleteResource(Long requesterId, StorageResourceType resourceType, @Nullable Long resourceIdNum, @Nullable String resourceIdStr, List<Value> customAttributes), queryResourceUploadInfo(Long requesterId, StorageResourceType resourceType, @Nullable Long resourceIdNum, @Nullable String resourceName, @Nullable String resourceMediaType, List<Value> customAttributes), queryResourceDownloadInfo(Long requesterId, StorageResourceType resourceType, @Nullable Long resourceIdNum, @Nullable String resourceIdStr, List<Value> customAttributes), shareMessageAttachmentWithUser(Long requesterId, @Nullable Long messageAttachmentIdNum, @Nullable String messageAttachmentIdStr, Long userIdToShareWith), shareMessageAttachmentWithGroup(Long requesterId, @Nullable Long messageAttachmentIdNum, @Nullable String messageAttachmentIdStr, Long groupIdToShareWith), unshareMessageAttachmentWithUser(Long requesterId, @Nullable Long messageAttachmentIdNum, @Nullable String messageAttachmentIdStr, Long userIdToUnshareWith), unshareMessageAttachmentWithGroup(Long requesterId, @Nullable Long messageAttachmentIdNum, @Nullable String messageAttachmentIdStr, Long groupIdToUnshareWith), queryMessageAttachmentInfosUploadedByRequester(Long requesterId, @Nullable DateRange creationDateRange), queryMessageAttachmentInfosInPrivateConversations(Long requesterId, @Nullable Set<Long> userIds, @Nullable DateRange creationDateRange, @Nullable Boolean areSharedByRequester), queryMessageAttachmentInfosInGroupConversations(Long requesterId, @Nullable Set<Long> groupIds, @Nullable Set<Long> userIds, @Nullable DateRange creationDateRange)*

Now I have a complete picture. Let me compare each method systematically.

## deleteResource

**Java logic:**
1. Validates `requesterId` not null
2. Validates `resourceType` not null
3. Checks `resourceType != UNRECOGNIZED`
4. Switch on `resourceType`:
   - `USER_PROFILE_PICTURE` → delegates to `provider.deleteUserProfilePicture(requesterId, customAttributes)`
   - `GROUP_PROFILE_PICTURE` → validates `resourceIdNum != null`, then delegates to `provider.deleteGroupProfilePicture(requesterId, resourceIdNum, customAttributes)`
   - `MESSAGE_ATTACHMENT` → delegates to `provider.deleteMessageAttachment(requesterId, resourceIdNum, resourceIdStr, customAttributes)`
   - default → error

**Go logic:**
1. Checks `resourceType == 0` (equivalent to UNRECOGNIZED)
2. Directly calls `s.provider.DeleteResource(ctx, resourceType, resourceIDStr)`
3. **Missing**: No validation of `requesterID`
4. **Missing**: No switch dispatching by resource type - just passes to a generic `DeleteResource`
5. **Missing**: The `resourceIdNum` parameter is entirely dropped
6. **Missing**: `customAttributes` parameter is entirely dropped
7. **Missing**: For GROUP_PROFILE_PICTURE, no null check on resourceIdNum

## queryResourceUploadInfo

**Java logic:**
1. Validates `requesterId`, `resourceType`, `customAttributes` not null
2. Checks `resourceType != UNRECOGNIZED`
3. Parses `resourceMediaType` into `MediaType`
4. Switch on `resourceType`:
   - `USER_PROFILE_PICTURE` → `provider.queryUserProfilePictureUploadInfo(requesterId, resourceName, mediaType, customAttributes)`
   - `GROUP_PROFILE_PICTURE` → validates `resourceIdNum != null`, then `provider.queryGroupProfilePictureUploadInfo(requesterId, resourceIdNum, resourceName, mediaType, customAttributes)`
   - `MESSAGE_ATTACHMENT` → complex branching on `resourceIdNum`:
     - null → `provider.queryMessageAttachmentUploadInfo(requesterId, resourceName, mediaType, customAttributes)`
     - negative → `provider.queryMessageAttachmentUploadInfoInGroupConversation(requesterId, -resourceIdNum, resourceName, mediaType, customAttributes)`
     - positive → `provider.queryMessageAttachmentUploadInfoInPrivateConversation(requesterId, resourceIdNum, resourceName, mediaType, customAttributes)`
   - default → error

**Go logic:**
1. Checks `resourceType == 0`
2. Returns `s.provider.GetPresignedUploadURL(ctx, resourceType, resourceName, contentType, maxSize)`
3. **Missing**: No validation of `requesterID`
4. **Missing**: No switch on resource type at all - completely flattened to a single generic call
5. **Missing**: The `resourceIdNum` parameter and its tri-state branching logic (null/negative/positive for MESSAGE_ATTACHMENT) is entirely lost
6. **Missing**: `customAttributes` parameter is entirely dropped
7. **Missing**: GROUP_PROFILE_PICTURE null check on resourceIdNum
8. **Missing**: The Go interface passes `maxSize int64` which doesn't exist in the Java signature at all

## queryResourceDownloadInfo

**Java logic:**
1. Validates `requesterId`, `resourceType`, `customAttributes` not null
2. Checks `resourceType != UNRECOGNIZED`
3. Switch on `resourceType`:
   - `USER_PROFILE_PICTURE` → validates `resourceIdNum != null`, then `provider.queryUserProfilePictureDownloadInfo(requesterId, resourceIdNum, customAttributes)`
   - `GROUP_PROFILE_PICTURE` → validates `resourceIdNum != null`, then `provider.queryGroupProfilePictureDownloadInfo(requesterId, resourceIdNum, customAttributes)`
   - `MESSAGE_ATTACHMENT` → `provider.queryMessageAttachmentDownloadInfo(requesterId, resourceIdNum, resourceIdStr, customAttributes)`
   - default → error

**Go logic:**
1. Checks `resourceType == 0`
2. Returns `s.provider.GetPresignedDownloadURL(ctx, resourceType, resourceIDStr)`
3. **Missing**: No validation of `requesterID`
4. **Missing**: No switch on resource type - flattened to single generic call
5. **Missing**: `resourceIdNum` parameter dropped entirely
6. **Missing**: `customAttributes` parameter dropped entirely
7. **Missing**: USER_PROFILE_PICTURE and GROUP_PROFILE_PICTURE null check on resourceIdNum

## shareMessageAttachmentWithUser

**Java logic:**
1. Validates `requesterId` not null
2. Validates `userIdToShareWith` not null
3. Delegates to provider via plugin manager

**Go logic:**
1. Directly delegates to provider with no validation
2. **Missing**: No validation that `requesterID != 0` (equivalent to null check)
3. **Missing**: No validation that `userIDToShareWith != 0`

## shareMessageAttachmentWithGroup

**Java logic:**
1. Validates `requesterId` not null
2. Validates `groupIdToShareWith` not null
3. Delegates to provider

**Go logic:**
1. Directly delegates with no validation
2. **Missing**: No validation of `requesterID != 0`
3. **Missing**: No validation of `groupIDToShareWith != 0`

## unshareMessageAttachmentWithUser

**Java logic:**
1. Validates `requesterId` not null
2. Validates `userIdToUnshareWith` not null
3. Delegates to provider

**Go logic:**
1. Directly delegates with no validation
2. **Missing**: No validation of `requesterID != 0`
3. **Missing**: No validation of `userIDToUnshareWith != 0`

## unshareMessageAttachmentWithGroup

**Java logic:**
1. Validates `requesterId` not null
2. Validates `groupIdToUnshareWith` not null
3. Delegates to provider

**Go logic:**
1. Directly delegates with no validation
2. **Missing**: No validation of `requesterID != 0`
3. **Missing**: No validation of `groupIDToUnshareWith != 0`

## queryMessageAttachmentInfosUploadedByRequester

**Java logic:**
1. Validates `requesterId` not null
2. Passes `DateRange` object to provider

**Go logic:**
1. Directly delegates with no validation
2. **Missing**: No validation of `requesterID != 0`

## queryMessageAttachmentInfosInPrivateConversations

**Java logic:**
1. Validates `requesterId` not null
2. Passes `userIds` (Set), `DateRange`, `areSharedByRequester` to provider

**Go logic:**
1. Directly delegates with no validation
2. **Missing**: No validation of `requesterID != 0`

## queryMessageAttachmentInfosInGroupConversations

**Java logic:**
1. Validates `requesterId` not null
2. Passes `groupIds` (Set), `userIds` (Set), `DateRange` to provider

**Go logic:**
1. Directly delegates with no validation
2. **Missing**: No validation of `requesterID != 0`

---

## deleteResource

- [ ] **Missing `requesterID` validation**: The Java version validates `requesterId` is not null with `Validator.notNull(requesterId, "requesterId")`. The Go version has no equivalent check.
- [ ] **Missing resource-type-based dispatch logic**: The Java version performs a switch on `resourceType` and routes to completely different provider methods (`deleteUserProfilePicture`, `deleteGroupProfilePicture`, `deleteMessageAttachment`) with different parameter signatures. The Go version collapses this into a single generic `DeleteResource(ctx, resourceType, resourceIDStr)` call, losing all type-specific dispatch.
- [ ] **Missing `resourceIdNum` parameter**: The Java version accepts `resourceIdNum` and passes it to providers. For `GROUP_PROFILE_PICTURE`, it validates `resourceIdNum != null` before delegating. The Go version drops this parameter entirely.
- [ ] **Missing `customAttributes` parameter**: The Java version passes `customAttributes` (a `List<Value>`) to every provider method. The Go version drops it entirely.
- [ ] **Missing GROUP_PROFILE_PICTURE null-check on resourceIdNum**: Java validates `resourceIdNum != null` before calling `deleteGroupProfilePicture`. Go has no equivalent.

## queryResourceUploadInfo

- [ ] **Missing `requesterID` validation**: Java validates `requesterId` is not null. Go has no equivalent check.
- [ ] **Missing resource-type-based dispatch logic**: Java routes to `queryUserProfilePictureUploadInfo`, `queryGroupProfilePictureUploadInfo`, or one of three `MESSAGE_ATTACHMENT` methods depending on resource type. Go collapses all of this into a single generic `GetPresignedUploadURL` call.
- [ ] **Missing `resourceIdNum` tri-state branching for MESSAGE_ATTACHMENT**: Java has critical branching: `resourceIdNum == null` → `queryMessageAttachmentUploadInfo`; `resourceIdNum < 0` → `queryMessageAttachmentUploadInfoInGroupConversation` (with `-resourceIdNum`); `resourceIdNum > 0` → `queryMessageAttachmentUploadInfoInPrivateConversation`. This entire routing logic is absent in Go.
- [ ] **Missing `customAttributes` parameter**: Java passes `customAttributes` to all provider methods. Go drops it entirely.
- [ ] **Missing GROUP_PROFILE_PICTURE null-check on resourceIdNum**: Java validates `resourceIdNum != null` before calling `queryGroupProfilePictureUploadInfo`. Go has no equivalent.
- [ ] **Extra `maxSize int64` parameter not present in Java**: The Go service method accepts a `maxSize int64` parameter and passes it to `GetPresignedUploadURL`. This parameter does not exist in the Java `queryResourceUploadInfo` signature.
- [ ] **Missing media type parsing/validation**: Java parses `resourceMediaType` string into a `MediaType` object with try/catch for invalid values. Go just passes `contentType` as a raw string with no parsing or validation.

## queryResourceDownloadInfo

- [ ] **Missing `requesterID` validation**: Java validates `requesterId` is not null. Go has no equivalent check.
- [ ] **Missing resource-type-based dispatch logic**: Java routes to `queryUserProfilePictureDownloadInfo`, `queryGroupProfilePictureDownloadInfo`, or `queryMessageAttachmentDownloadInfo`. Go collapses all to a single generic `GetPresignedDownloadURL`.
- [ ] **Missing `resourceIdNum` parameter**: Java passes `resourceIdNum` to providers. Go drops it entirely.
- [ ] **Missing `customAttributes` parameter**: Java passes `customAttributes` to all provider methods. Go drops it entirely.
- [ ] **Missing USER_PROFILE_PICTURE null-check on resourceIdNum**: Java validates `resourceIdNum != null` before calling `queryUserProfilePictureDownloadInfo`. Go has no equivalent.
- [ ] **Missing GROUP_PROFILE_PICTURE null-check on resourceIdNum**: Java validates `resourceIdNum != null` before calling `queryGroupProfilePictureDownloadInfo`. Go has no equivalent.

## shareMessageAttachmentWithUser

- [ ] **Missing `requesterID` validation**: Java validates `requesterId` is not null. Go has no validation on `requesterID`.
- [ ] **Missing `userIDToShareWith` validation**: Java validates `userIdToShareWith` is not null. Go has no validation on `userIDToShareWith`.

## shareMessageAttachmentWithGroup

- [ ] **Missing `requesterID` validation**: Java validates `requesterId` is not null. Go has no validation on `requesterID`.
- [ ] **Missing `groupIDToShareWith` validation**: Java validates `groupIdToShareWith` is not null. Go has no validation on `groupIDToShareWith`.

## unshareMessageAttachmentWithUser

- [ ] **Missing `requesterID` validation**: Java validates `requesterId` is not null. Go has no validation on `requesterID`.
- [ ] **Missing `userIDToUnshareWith` validation**: Java validates `userIdToUnshareWith` is not null. Go has no validation on `userIDToUnshareWith`.

## unshareMessageAttachmentWithGroup

- [ ] **Missing `requesterID` validation**: Java validates `requesterId` is not null. Go has no validation on `requesterID`.
- [ ] **Missing `groupIDToUnshareWith` validation**: Java validates `groupIdToUnshareWith` is not null. Go has no validation on `groupIDToUnshareWith`.

## queryMessageAttachmentInfosUploadedByRequester

- [ ] **Missing `requesterID` validation**: Java validates `requesterId` is not null. Go has no validation on `requesterID`.

## queryMessageAttachmentInfosInPrivateConversations

- [ ] **Missing `requesterID` validation**: Java validates `requesterId` is not null. Go has no validation on `requesterID`.

## queryMessageAttachmentInfosInGroupConversations

- [ ] **Missing `requesterID` validation**: Java validates `requesterId` is not null. Go has no validation on `requesterID`.

# UserController.java
*Checked methods: addUser(@RequestBody AddUserDTO addUserDTO), queryUsers(@QueryParam(required = false), queryUsers(@QueryParam(required = false), countUsers(@QueryParam(required = false), updateUser(Set<Long> ids, @RequestBody UpdateUserDTO updateUserDTO), deleteUsers(Set<Long> ids, @QueryParam(required = false)*

Now I have all the information needed to do a thorough comparison. Let me compile my findings.

## addUser

- [ ] **Missing password encoding**: The Java code calls `passwordManager.encodeUserPassword(rawPassword)` to encode the password before storing, and stores it as `byte[]`. The Go code stores the raw password string directly without any encoding/hashing.
- [ ] **Missing default value for `name`**: Java sets `name = name == null ? "" : name` (defaults empty string when null). The Go code does not default null/empty name to `""`.
- [ ] **Missing default value for `intro`**: Java sets `intro = intro == null ? "" : intro` (defaults empty string when null). The Go code does not default null/empty intro to `""`.
- [ ] **Missing default value for `profileAccessStrategy`**: Java sets `profileAccessStrategy = ProfileAccessStrategy.ALL` when null. The Go code does not set a default for `profileAccess`.
- [ ] **Missing default value for `roleId`**: Java sets `roleId = DEFAULT_USER_ROLE_ID` when null. The Go code does not set a default for `permissionGroupID`.
- [ ] **Missing default value for `isActive`**: Java sets `isActive = activateUserWhenAdded` (a configurable property) when null. The Go `AddUser` takes `isActive bool` (not a pointer), so it cannot represent null and defaults to `false` instead of the configurable `activateUserWhenAdded` value.
- [ ] **Missing `LastUpdatedDate` field**: Java constructs User with `now` as the `LastUpdatedDate` (9th constructor arg). The Go code does not set `LastUpdatedDate` on the user struct.
- [ ] **Missing `DeletionDate` field**: Java constructs User with `null` as `DeletionDate` (10th constructor arg). This is likely fine since Go zero-value would be nil, but the Java also explicitly constructs with positional args — the Go struct should ensure `DeletionDate` is nil.
- [ ] **Missing relationship group creation**: Java calls `userRelationshipGroupService.createRelationshipGroup(finalId, 0, "", now, session)` after inserting the user. The Go code only inserts the user without creating a default relationship group.
- [ ] **Missing user version upsert**: Java calls `userVersionService.upsertEmptyUserVersion(user.getId(), date, session)` after inserting. The Go code does not upsert user version.
- [ ] **Missing metrics counter increment**: Java increments `registeredUsersCounter` on success. The Go code does not increment any metrics counter.
- [ ] **Missing transaction support**: Java wraps the insert + relationship group + version upsert in a transaction. The Go code does a simple insert without transaction.
- [ ] **Missing Elasticsearch sync**: Java syncs user doc to Elasticsearch after creation. The Go code does not interact with Elasticsearch.

## queryUsers (non-paginated, `@GetMapping`)

- [ ] **Missing filter parameters**: Java `queryUsers` accepts `ids`, `registrationDateStart/End`, `deletionDateStart/End`, `isActive`, and `size`. The Go `QueryUsers` only accepts `userIDs []int64` — it's missing registration date range, deletion date range, isActive, and pagination (page/size) parameters entirely.
- [ ] **Missing `queryDeletedRecords` parameter**: Java passes `queryDeletedRecords=true` (always true from controller). The Go code has no such filter, so deleted users would never be excluded or included based on this flag — it simply doesn't filter on deletion date.
- [ ] **Empty userIDs returns empty too early**: The Go code returns `[]*po.User{}` when `len(userIDs) == 0`. In Java, `userIds` is nullable/optional, and when null, it queries all users (with `inIfNotNull` skipping the filter). The Go code should support nil/empty userIDs by not filtering on `_id`.

## queryUsers (paginated, `@GetMapping("page")`)

- [ ] **Same missing filters as non-paginated queryUsers**: The Go `QueryUsers` method is the same method used for both — it lacks registration date range, deletion date range, isActive filter, page, size, and queryDeletedRecords parameters.
- [ ] **Missing count query**: Java's paginated endpoint calls `userService.countUsers(ids, ...)` for the total count, then calls `userService.queryUsers(...)` with pagination. The Go code has no equivalent paginated query that combines count + filtered results.

## countUsers (`@GetMapping("count")`)

- [ ] **Missing parameters entirely**: Java's `countUsers` in the controller accepts `registeredStartDate/End`, `deletedStartDate/End`, `sentMessageStartDate/End`, and `divideBy`. The Go `CountUsers` only accepts `activeOnly bool`. The Go `CountRegisteredUsers` and `CountDeletedUsers` exist but the controller-level `countUsers` endpoint logic (which combines counts with `UserStatisticsDTO` and supports `DivideBy`) is completely missing.
- [ ] **CountUsers service method is oversimplified**: The Go `CountUsers(ctx, activeOnly)` only filters by `isActive`. The Java `countUsers(userIds, registrationDateRange, deletionDateRange, isActive)` supports filtering by user IDs, registration date range, deletion date range, and isActive. The Go method is missing all parameters except isActive.
- [ ] **CountRegisteredUsers uses wrong field name**: The Go repository `CountRegisteredUsers` filters on field `"cd"` (creation date?), but the Java code filters on `User.Fields.REGISTRATION_DATE` which maps to `"rd"` based on the User PO field naming convention used elsewhere in the Go code.

## updateUser / updateUsers (`@PutMapping`)

- [ ] **Missing password encoding**: Java calls `passwordManager.encodeUserPassword(rawPassword)` to encode passwords before updating. The Go code accepts a raw `bson.M` update and applies it directly without any password encoding.
- [ ] **Missing validation**: Java validates password length, name length, intro length, profile picture length, profile access strategy, and registration date. The Go code does no validation.
- [ ] **Missing "all falsy" early return**: Java checks if all update fields are null/empty and returns `ACKNOWLEDGED_UPDATE_RESULT` without hitting the database. The Go code doesn't do this check.
- [ ] **Missing session disconnection on inactivation**: Java's private `updateUsers` method disconnects user sessions when `isActive` is set to `false`. The Go code does not handle this.
- [ ] **Missing Elasticsearch sync**: Java syncs user name changes to Elasticsearch. The Go code does not.
- [ ] **Missing user-defined attributes handling**: Java passes `userDefinedAttributes` to the repository update. The Go code uses a generic `bson.M` but doesn't handle user-defined attributes specifically.
- [ ] **The controller passes `null` for userDefinedAttributes**: In Java, the controller calls `updateUsers(ids, ..., null)` for the last parameter (userDefinedAttributes). This is handled correctly since the Go version uses `bson.M` directly, but there's no translation layer from the DTO fields to the proper BSON update document — the caller must construct the correct BSON manually.

## deleteUsers (`@DeleteMapping`)

- [ ] **Missing logical deletion support**: Java's `deleteUsers` accepts a `deleteLogically` parameter. When true (or when `deleteUserLogically` config is true), it updates the deletion date instead of physically deleting. When false, it physically deletes. The Go code always physically deletes via `DeleteMany` — there is no logical deletion path.
- [ ] **Missing cascade deletion of related data**: Java's physical deletion deletes relationships, relationship groups, user settings, private conversations, group member conversations, conversation settings, user versions, and message sequence IDs. The Go code only deletes the user document itself.
- [ ] **Missing session disconnection**: Java disconnects active sessions for deleted/inactivated users. The Go code does not disconnect sessions.
- [ ] **Missing metrics counter increment**: Java increments `deletedUsersCounter` on physical deletion success. The Go code does not.
- [ ] **Missing transaction support**: Java wraps the physical deletion and all cascading deletes in a transaction. The Go code does a simple delete without transaction.
- [ ] **Missing Elasticsearch sync**: Java deletes user docs from Elasticsearch. The Go code does not.
- [ ] **`deleteLogically` parameter is completely ignored**: The Go `DeleteUsers` method signature `(ctx, userIDs)` does not accept a `deleteLogically` parameter at all, so this behavior cannot be controlled by the caller.

# UserOnlineInfoController.java
*Checked methods: countOnlineUsers(boolean countByNodes), queryUserSessions(Set<Long> ids, boolean returnNonExistingUsers), queryUserStatuses(Set<Long> ids, boolean returnNonExistingUsers), queryUsersNearby(Long userId, @QueryParam(required = false), queryUserLocations(Set<Long> ids, @QueryParam(required = false), updateUserOnlineStatus(Set<Long> ids, @QueryParam(required = false)*

Now I have a comprehensive understanding of all the code. Let me compile the review.

## countOnlineUsers

**Java**: `countOnlineUsers(boolean countByNodes)` — when `countByNodes=false`, calls `statisticsService.countOnlineUsers()` which sends RPC to all gateway nodes and sums the counts. When `countByNodes=true`, calls `statisticsService.countOnlineUsersByNodes()` which returns a `Map<String, Integer>` (nodeId→count), sums them, and returns `OnlineUserCountDTO(sum, nodeIdToUserCount)`.

**Go**: `CountOnlineUsers()` in `user_controllers.go` is a stub (`// TODO: implement`). The `StatisticsService` in Go only has `CountLocalOnlineUsers()` which counts local sessions — it does NOT aggregate from all gateway nodes via RPC. There is no `countOnlineUsersByNodes` equivalent at all.

## queryUserSessions

**Java**: `SessionService.queryUserSessions(Set<Long> userIds)` does a complex multi-step process:
1. Fetches `UserSessionsStatus` for each user via `userStatusService.getUserSessionsStatus(userId)`
2. Separates offline users (empty device map) into `offlineUserSessions` list with `UserStatus.OFFLINE`
3. Builds a `nodeIdToUserIds` map for online users
4. Sends RPC `QueryUserSessionsRequest` to each node to get detailed session info (with `id`, `version`, `deviceType`, `deviceDetails`, `loginDate`, `loginIp`, etc.)
5. Merges results from multiple nodes (a user can have sessions on different nodes)

**Go**: `sessionService.QueryUserSessions` in `session_service.go` only fetches `UserSessionsStatus` from Redis, and builds a simple `UserSessionsInfo` with just `DeviceType` in each session. It does NOT:
- Send RPC to gateway nodes to get detailed session info (id, version, loginDate, location, IP, etc.)
- Handle the offline user case properly (Java explicitly creates entries with `UserStatus.OFFLINE` and empty sessions)
- Merge sessions across multiple nodes

The Go `UserSessionInfo` struct only has `DeviceType`, whereas the Java `UserSessionInfo` record has `id`, `version`, `deviceType`, `deviceDetails`, `loginDate`, `loginLocation`, `lastHeartbeatRequestDate`, `lastRequestDate`, `isSessionOpen`, `ipBytes`, `ip`.

Also, the Go `UserSessionsInfo` (in the `onlineuser` package) is missing the `Status` field that the Java version has. Wait, it does have `Status` in the `onlineuser` package. But the `session/bo` package's `UserSessionsInfo` is missing `Status`.

## queryUserStatuses

**Java**: For each userId, calls `userStatusService.getUserSessionsStatus(userId)`. If `!returnNonExistingUsers` and the status is OFFLINE, checks `userService.checkIfUserExists(userId, false)` and filters out non-existing users.

**Go**: Controller is a stub (`// TODO: implement`).

## queryUsersNearby

**Java**: Calls `nearbyUserService.queryNearbyUsers(userId, deviceType, null, null, maxCount, maxDistance, withCoordinates, withDistance, withUserInfo)`. The `deviceType` is a required query param.

**Go**: The `QueryNearbyUsers` Go method signature requires `longitude *float32, latitude *float32` — but the Java controller passes `null` for both longitude and latitude, which means it gets the user's current location first from the session. The Go implementation returns empty if longitude/latitude are nil, so it won't work the same way as Java when called from the controller.

## queryUserLocations

**Java**: For each userId, calls `sessionLocationService.getUserLocation(userId, deviceType)` and maps coordinates to `UserLocationDTO(userId, deviceType, coordinates.getX().doubleValue(), coordinates.getY().doubleValue())`.

**Go**: Controller is a stub (`// TODO: implement`). The `SessionLocationService.GetUserLocation` returns `(*protocol.UserLocation, error)`.

## updateUserOnlineStatus

**Java**: If `onlineStatus == OFFLINE`, calls `sessionService.disconnect(ids, ...)` or `sessionService.disconnect(ids, deviceTypes, ...)`. Otherwise calls `userStatusService.updateOnlineUsersStatus(ids, onlineStatus)`.

**Go**: Controller is a stub (`// TODO: implement`). The Go `SessionService.Disconnect` methods are also stubs that just return `true, nil` without actually doing anything.

Here is the complete bug report:

---

## countOnlineUsers

- [ ] **Controller is unimplemented**: `UserOnlineInfoController.CountOnlineUsers()` in `user_controllers.go:39` is a stub (`// TODO: implement`) with no logic at all.

- [ ] **Missing RPC-based cluster-wide count**: The Java `statisticsService.countOnlineUsers()` sends RPC to all gateway nodes and sums the results. The Go `StatisticsService` only has `CountLocalOnlineUsers()` which counts local sessions — it does not aggregate counts from other gateway nodes via RPC.

- [ ] **Missing `countOnlineUsersByNodes` implementation**: Java supports `countByNodes=true` which returns a `Map<String, Integer>` (nodeId→count). There is no equivalent in the Go code. The Java controller sums all values and returns `OnlineUserCountDTO(sum, nodeIdToUserCount)`, but neither the per-node map nor the `OnlineUserCountDTO` struct exist in Go.

## queryUserSessions

- [ ] **Controller is unimplemented**: `UserOnlineInfoController.QueryUserSessions()` in `user_controllers.go:44` is a stub (`// TODO: implement`) with no logic at all.

- [ ] **Missing RPC to gateway nodes for detailed session info**: Java's `SessionService.queryUserSessions()` builds a `nodeIdToUserIds` map and sends RPC `QueryUserSessionsRequest` to each gateway node to retrieve detailed session data (id, version, deviceType, deviceDetails, loginDate, loginLocation, lastHeartbeatRequestDate, lastRequestDate, isSessionOpen, IP). The Go `sessionService.QueryUserSessions()` only reads from Redis status data and never contacts gateway nodes via RPC.

- [ ] **Missing session detail fields**: The Java `UserSessionInfo` record has fields: `id`, `version`, `deviceType`, `deviceDetails`, `loginDate`, `loginLocation`, `lastHeartbeatRequestDate`, `lastRequestDate`, `isSessionOpen`, `ipBytes`, `ip`. The Go `onlineuser.UserSessionInfo` struct only has `DeviceType`. The Go `session/bo.UserSessionInfo` has `ID`, `Version`, `DeviceType`, `LoginDate`, `Location` but is not used by `sessionService.QueryUserSessions()`.

- [ ] **Missing offline user handling**: Java explicitly creates `UserSessionsInfo` entries with `UserStatus.OFFLINE` and empty sessions list for offline users. The Go `sessionService.QueryUserSessions()` skips users that return errors from `FetchUserSessionsStatus` (line 71: `continue`) and does not add offline user entries, so offline users are silently dropped instead of being returned with OFFLINE status.

- [ ] **Missing cross-node session merging**: Java has a `mergeUserSessions()` method that merges sessions for users who have sessions across multiple gateway nodes. The Go implementation has no such merging logic.

- [ ] **Missing `returnNonExistingUsers` filtering**: Java's controller filters out non-existing users when `returnNonExistingUsers=false` by calling `userService.checkIfUserExists()` for offline entries. The Go controller and service have no such filtering logic.

## queryUserStatuses

- [ ] **Controller is unimplemented**: `UserOnlineInfoController.QueryUserStatuses()` in `user_controllers.go:49` is a stub (`// TODO: implement`) with no logic at all.

- [ ] **Missing `returnNonExistingUsers` filtering**: Java calls `userService.checkIfUserExists(userId, false)` for offline users when `returnNonExistingUsers=false` to filter out non-existent users. The Go code has no equivalent logic anywhere.

## queryUsersNearby

- [ ] **Controller is unimplemented**: `UserOnlineInfoController.QueryUsersNearby()` in `user_controllers.go:54` is a stub (`// TODO: implement`) with no logic at all.

- [ ] **Go service returns empty when longitude/latitude are nil**: The Java controller passes `null` for longitude and latitude to `nearbyUserService.queryNearbyUsers()`, which then fetches the user's current location from their session. The Go `QueryNearbyUsers()` returns `[]*NearbyUser{}` immediately when `longitude == nil || latitude == nil` (line 44-46 of `nearby_user_service.go`), so it cannot work when called from the controller pattern like Java does.

## queryUserLocations

- [ ] **Controller is unimplemented**: `UserOnlineInfoController.QueryUserLocations()` in `user_controllers.go:59` is a stub (`// TODO: implement`) with no logic at all.

- [ ] **Missing `deviceType` parameter handling**: Java's `queryUserLocations` accepts an optional `DeviceType deviceType` parameter and passes it to `sessionLocationService.getUserLocation(userId, deviceType)`. The Go `SessionLocationService.GetUserLocation` also takes a `deviceType` parameter, but the controller stub does not wire it.

- [ ] **Missing UserLocationDTO construction**: Java maps the result into `UserLocationDTO(userId, deviceType, coordinates.getX().doubleValue(), coordinates.getY().doubleValue())` for each user. This DTO does not exist in Go.

## updateUserOnlineStatus

- [ ] **Controller is unimplemented**: `UserOnlineInfoController.UpdateUserOnlineStatus()` in `user_controllers.go:64` is a stub (`// TODO: implement`) with no logic at all.

- [ ] **Disconnect methods are no-ops**: The Go `sessionService.Disconnect()`, `DisconnectWithDeviceTypes()`, `DisconnectMultipleUsers()`, and `DisconnectMultipleUsersWithDeviceTypes()` all return `true, nil` without performing any actual disconnection (lines 43-62 of `session_service.go`). Java's equivalent methods send RPC requests to gateway nodes to actually close sessions.

- [ ] **Missing `updateOnlineUsersStatus` bulk method**: Java calls `userStatusService.updateOnlineUsersStatus(ids, onlineStatus)` which updates the status for multiple users (iterating and calling `updateOnlineUserStatusIfPresent` per user). The Go `UserStatusService` only has `UpdateStatus` for a single user, and no bulk `UpdateOnlineUsersStatus` equivalent.

# UserRoleController.java
*Checked methods: addUserRole(@RequestBody AddUserRoleDTO addUserRoleDTO), queryUserRoles(@QueryParam(required = false), queryUserRoleGroups(int page, @QueryParam(required = false), updateUserRole(Set<Long> ids, @RequestBody UpdateUserRoleDTO updateUserRoleDTO), deleteUserRole(Set<Long> ids)*

Now I have a comprehensive understanding of both codebases. Let me compile the review.

## addUserRole

- [ ] **Controller is a no-op stub**: The Go `UserRoleController.AddUserRole()` method at `user_controllers.go:74` is entirely unimplemented (`// TODO: implement`), with no parameters, no service call, and no return value. The Java version accepts `AddUserRoleDTO`, normalizes null `creatableGroupTypeIds` to `Collections.emptySet()` and null `groupTypeIdToLimit` to `Collections.emptyMap()`, then calls `userRoleService.addUserRole(...)` with all 6 fields.
- [ ] **Service passes DTO directly without normalization**: The Go `UserRoleService.AddUserRole()` at `user_role_service.go:29` accepts a `*po.UserRole` and directly inserts it. The Java version normalizes null `creatableGroupTypeIds` to `emptySet()` and null `groupTypeIdToLimit` to `emptyMap()` before constructing the `UserRole`. The Go service does not perform this null-to-empty normalization.
- [ ] **Missing ID auto-generation**: The Java `UserRoleService.addUserRole()` auto-generates an ID via `node.nextLargeGapId(ServiceType.USER_ROLE)` when `groupId == null`. The Go service has no equivalent ID generation logic.
- [ ] **Missing in-memory cache update**: The Java version updates `idToRole.put(groupId, userRole)` after constructing the role. The Go version has no equivalent in-memory cache (`idToRole`).
- [ ] **Missing null/not-null validation**: The Java version validates that `creatableGroupTypeIds`, `ownedGroupLimit`, `ownedGroupLimitForEachGroupType`, and `groupTypeIdToLimit` are not null. The Go service performs no validation.

## queryUserRoles

- [ ] **Controller is a no-op stub**: The Go `UserRoleController.QueryUserRoles()` at `user_controllers.go:79` is entirely unimplemented. The Java version accepts an optional `size` parameter, applies `getPageSize(size)` for default sizing, calls `userRoleService.queryUserRoles(0, size)`, and wraps the result with `HttpHandlerResult.okIfTruthy()`.
- [ ] **Service uses generic filter instead of pagination parameters**: The Go `UserRoleService.QueryUserRoles()` at `user_role_service.go:23` takes a generic `bson.M` filter, whereas the Java version takes explicit `(Integer page, Integer size)` pagination parameters. The Go service has no pagination support, meaning it will return all roles with no limit, unlike the Java version which defaults page=0 and applies a page size.

## queryUserRoleGroups

- [ ] **Controller is a no-op stub**: The Go `UserRoleController.QueryUserRoleGroups()` at `user_controllers.go:84` is entirely unimplemented. The Java version accepts `page` and `size` parameters, applies `getPageSize(size)`, calls both `countUserRoles()` and `queryUserRoles(page, size)`, and wraps them into a paginated response via `HttpHandlerResult.page(count, userRoleFlux)`.
- [ ] **Missing count + pagination orchestration**: The Java controller combines a count query and a paginated data query to produce a `PaginationDTO<UserRole>` response. The Go code has no equivalent logic anywhere — neither in the controller nor in the service layer.

## updateUserRole

- [ ] **Controller is a no-op stub**: The Go `UserRoleController.UpdateUserRole()` at `user_controllers.go:89` is entirely unimplemented. The Java version accepts `Set<Long> ids` and `UpdateUserRoleDTO`, calls `userRoleService.updateUserRoles(...)` with all 6 parameters, maps the result to `UpdateResultDTO`, and wraps it with `HttpHandlerResult.okIfTruthy()`.
- [ ] **Service is a no-op placeholder**: The Go `UserRoleService.UpdateUserRoles()` at `user_role_service.go:33` returns `nil` unconditionally — it does not perform any update operation. The comment explicitly states this is a placeholder: "Not implementing complex update parsing... We'll leave it as a placeholder."
- [ ] **Missing empty-update short-circuit**: The Java version checks if all update fields are null and returns `ACKNOWLEDGED_UPDATE_RESULT` early without hitting the database. The Go version has no such optimization (nor any update logic at all).
- [ ] **Missing in-memory cache invalidation**: The Java version invalidates cache entries (`idToRole.remove(groupId)`) for all updated IDs via `.doOnNext()`. The Go version has no equivalent cache invalidation.
- [ ] **Missing validation for non-empty ids**: The Java version validates that `groupIds` is not empty (`Validator.notEmpty(groupIds, "groupIds")`). The Go version performs no validation.

## deleteUserRole

- [ ] **Controller is a no-op stub**: The Go `UserRoleController.DeleteUserRole()` at `user_controllers.go:94` is entirely unimplemented. The Java version accepts `Set<Long> ids`, calls `userRoleService.deleteUserRoles(ids)`, maps the result to `DeleteResultDTO`, and wraps it with `HttpHandlerResult.okIfTruthy()`.
- [ ] **Missing default role deletion protection**: The Java `UserRoleService.deleteUserRoles()` checks if `groupIds` contains `DEFAULT_USER_ROLE_ID` and rejects the deletion with an error (`"The default user role cannot be deleted"`). The Go `DeleteUserRoles()` at `user_role_service.go:43` performs no such check and will happily delete the default role.
- [ ] **Missing null-ids fallback to delete-by-not-id**: When `groupIds` is null, the Java version calls `userRoleRepository.deleteByNotIds(Set.of(DEFAULT_USER_ROLE_ID))` to delete all roles except the default. The Go version passes the `filter` directly with no null-handling or fallback logic.
- [ ] **Missing in-memory cache invalidation**: The Java version invalidates cache entries (`idToRole.remove(groupId)`) for all deleted IDs via `.doOnNext()`. The Go version has no equivalent cache invalidation.

# UserFriendRequestController.java
*Checked methods: createFriendRequest(@RequestBody AddFriendRequestDTO addFriendRequestDTO), queryFriendRequests(@QueryParam(required = false), queryFriendRequests(@QueryParam(required = false), updateFriendRequests(Set<Long> ids, @RequestBody UpdateFriendRequestDTO updateFriendRequestDTO), deleteFriendRequests(@QueryParam(required = false)*

Now I have a comprehensive understanding of both implementations. Let me compile my findings.

## createFriendRequest

- [ ] **Missing `Validator.notNull(content, "content")` validation**: The Java code validates that `content` is not null (`Validator.notNull(content, "content")`), but the Go code only validates `MaxLength` on content and never checks that `content` is not empty/null. In Go, `content` is a `string` (not `*string`), so it will default to `""` — but the Java code explicitly validates it's not null, meaning it expects a non-null value. The Go code skips this `NotNull` check.
- [ ] **Missing `DataValidator.validRequestStatus(status)` validation**: The Java code calls `DataValidator.validRequestStatus(status)` to validate the status enum value. The Go code does not perform any equivalent request status validity check.
- [ ] **Missing `Validator.pastOrPresent(creationDate, "creationDate")` and `Validator.pastOrPresent(responseDate, "responseDate")` validations**: The Java code validates both `creationDate` and `responseDate` with `pastOrPresent`. The Go `CreateFriendRequest` service method does not perform any `PastOrPresent` validation on these dates.
- [ ] **Incorrect `responseDate` logic for new records**: The Java code uses `getResponseDateBasedOnStatusForNewRecord(now, status, responseDate)`, which: (a) if `isProcessedByResponder(status)` (i.e., ACCEPTED/DECLINED/IGNORED) and `responseDate` is null → sets to `now`; (b) if `isProcessedByResponder(status)` and `responseDate` is not null → keeps `responseDate`; (c) if NOT `isProcessedByResponder` (PENDING/CANCELED/EXPIRED or null) → sets to `null`. The Go code at lines 119-123 only handles case (a) partially: it only sets `responseDate = &now` when status is not PENDING and responseDate is nil. It does NOT set responseDate to nil when status is PENDING/CANCELED (case c), meaning it preserves a provided responseDate even when the status is PENDING — which contradicts Java behavior that would null it out.

## queryFriendRequests

- [ ] **Missing `expirationDateStart`/`expirationDateEnd` filtering in repository**: The Java `findFriendRequests` uses `getCreationDateRange(creationDateRange, expirationDateRange)` to merge the expiration date range into the creation date filter, and also applies `isExpiredOrNot` logic to handle expired status filtering. The Go `countOrFind` method at `user_friend_request_repository.go:230-265` completely ignores `expirationDateStart` and `expirationDateEnd` — these parameters are accepted but never used in the filter construction.
- [ ] **Missing `isExpiredOrNot` / expired request transformation logic**: The Java `findFriendRequests` calls `isExpiredOrNot(statuses, creationDateFieldName, expirationDate)` which adds additional filters for expired/non-expired requests, and `findExpirableDocs` which transforms pending-but-expired requests to have `EXPIRED` status. The Go code has none of this logic — it does not handle request expiration at all in queries.

## updateFriendRequests

- [ ] **Missing `responseDate` update based on status**: The Java repository's `updateFriendRequests` calls `updateResponseDateBasedOnStatus(UserFriendRequest.Fields.RESPONSE_DATE, update, status, new Date())`, which: (a) if status is ACCEPTED/DECLINED/IGNORED → sets `responseDate` (to `new Date()` if null); (b) if status is null → skips; (c) if status is PENDING/CANCELED/EXPIRED → unsets `responseDate`. The Go repository's `UpdateFriendRequests` at `user_friend_request_repository.go:106-132` does not handle `responseDate` at all — it only updates `requesterID`, `recipientID`, `content`, `status`, `reason`, and `creationDate`. The `responseDate` parameter passed to the service method is completely ignored and never forwarded to the repository.
- [ ] **Missing `DataValidator.validRequestStatus(status)` validation**: The Java code validates the request status. The Go code does not perform equivalent validation.
- [ ] **Different null-safety check for requesterId/recipientId**: The Java code uses `Validator.shouldTrue(requesterId == null || !requesterId.equals(recipientId), ...)` which only checks inequality when BOTH are non-null. The Go code at line 265 checks `requesterID != nil && recipientID != nil && *requesterID == *recipientID`, which is actually equivalent behavior — this is NOT a bug.

## deleteFriendRequests

- [ ] **Different behavior when `ids` is null/empty**: The Java `deleteByIds` uses `inIfNotNull` which, when `ids` is null or empty, does NOT add any `_id` filter to the query. This means it would delete ALL documents in the collection (empty filter = match all). The Go `DeleteByIds` at `user_friend_request_repository.go:221-228` returns `nil` immediately when `len(ids) == 0`, which means it does NOTHING. This is a behavioral difference — when called with nil/empty ids from the controller (which has `@QueryParam(required = false)`), Java deletes all friend requests, while Go silently does nothing.

# UserRelationshipController.java
*Checked methods: addRelationship(@RequestBody AddRelationshipDTO addRelationshipDTO), queryRelationships(@QueryParam(required = false), queryRelationships(@QueryParam(required = false), updateRelationships(List<UserRelationship.Key> keys, @RequestBody UpdateRelationshipDTO updateRelationshipDTO), deleteRelationships(List<UserRelationship.Key> keys)*

Now I have all the information needed for a thorough comparison. Let me analyze the differences systematically.

## addRelationship

**Java** (`UserRelationshipController.java:82-95`): Calls `userRelationshipService.upsertOneSidedRelationship(ownerId, relatedUserId, name, blockDate, DEFAULT_RELATIONSHIP_GROUP_INDEX, null, establishmentDate, false, null)` — passes `name` as the 3rd parameter, `blockDate` as 4th, `DEFAULT_RELATIONSHIP_GROUP_INDEX` (0) as the group index, `null` as deleteGroupIndex, `establishmentDate` as 7th, `false` as `upsert`, and `null` session.

**Go** (`user_controllers.go:129-131`): The controller method `AddRelationship()` is a **stub** with `// TODO: implement`. No logic is implemented.

## queryRelationships (non-paginated)

**Java** (`UserRelationshipController.java:98-119`): Calls `userRelationshipService.queryRelationships(...)` with pagination `0` and `size`, then transforms results via `relationship2dto` which optionally queries group indexes, then returns via `HttpHandlerResult.okIfTruthy`.

**Go** (`user_controllers.go:134-136`): The controller method `QueryRelationships()` is a **stub** with `// TODO: implement`. The service layer `QueryRelationships` exists and appears correct.

## queryRelationships (paginated)

**Java** (`UserRelationshipController.java:121-146`): Separate `@GetMapping("page")` endpoint with pagination that counts first, queries with page offset, converts to DTOs, and returns via `HttpHandlerResult.page(count, dtoFlux)`.

**Go**: No paginated query method exists in the controller at all.

## updateRelationships

**Java** (`UserRelationshipController.java:148-160`): Takes `List<UserRelationship.Key> keys` (where Key = ownerId + relatedUserId pairs) and `UpdateRelationshipDTO` (name, blockDate, establishmentDate). Calls `userRelationshipService.updateUserOneSidedRelationships(CollectionUtil.newSet(keys), name, blockDate, establishmentDate)`.

**Go** (`user_controllers.go:139-141`): The controller method `UpdateRelationships()` is a **stub** with `// TODO: implement`. 

Also, the Go service method `UpdateUserOneSidedRelationships` has a different signature: it takes `(userID int64, relatedUserIDs []int64, ...)` — a single ownerID with multiple relatedUserIDs, whereas the Java version takes `Set<UserRelationship.Key>` which allows **multiple different ownerId/relatedUserId pairs** across different owners. The Go service signature **cannot handle keys with different ownerIds** in a single call.

## deleteRelationships

**Java** (`UserRelationshipController.java:162-170`): Takes `List<UserRelationship.Key> keys`, calls `userRelationshipService.deleteOneSidedRelationships(CollectionUtil.newSet(keys))` which operates on keys (ownerId + relatedUserId pairs, potentially from multiple owners).

**Go** (`user_controllers.go:144-146`): The controller method `DeleteRelationships()` is a **stub** with `// TODO: implement`.

Also, the Go service method `DeleteOneSidedRelationships` takes `(ownerID int64, relatedUserIDs []int64, session)` — a single ownerID with multiple relatedUserIDs, whereas the Java version operates on `Set<Key>` with potentially multiple different owners.

---

## addRelationship
- [ ] Controller method `AddRelationship()` is completely unimplemented (stub with `// TODO: implement`). The Java version calls `userRelationshipService.upsertOneSidedRelationship(ownerId, relatedUserId, name, blockDate, DEFAULT_RELATIONSHIP_GROUP_INDEX, null, establishmentDate, false, null)` with specific parameter ordering. The Go `UpsertOneSidedRelationship` service method exists but has a **different parameter order** compared to Java: Go takes `(ownerID, relatedUserID, blockDate, groupIndex, establishmentDate, name, session)` while Java takes `(ownerId, relatedUserId, name, blockDate, newGroupIndex, deleteGroupIndex, establishmentDate, upsert, session)`. The controller stub would need to correctly map `addRelationshipDTO.name` → `name` param and `addRelationshipDTO.blockDate` → `blockDate` param with the correct Go parameter positions.

## queryRelationships
- [ ] Controller method `QueryRelationships()` is completely unimplemented (stub with `// TODO: implement`). Missing: parameter parsing (ownerIds, relatedUserIds, groupIndexes, isBlocked, establishmentDateStart, establishmentDateEnd, size, withGroupIndexes), calling the service layer, converting results to DTOs (including optional group index lookup via `relationship2dto` logic), and returning the response.
- [ ] The `withGroupIndexes` flag and the `relationship2dto` helper logic (which optionally queries `userRelationshipGroupService.queryGroupIndexes` for each relationship to populate `groupIndexes` in the DTO) has no Go implementation anywhere in the controller or service layer.

## queryRelationships (page)
- [ ] No paginated query endpoint exists in the Go controller at all. The Java version has a separate `@GetMapping("page")` endpoint that calls `countRelationships` for total count and returns `PaginationDTO`. This is missing entirely.

## updateRelationships
- [ ] Controller method `UpdateRelationships()` is completely unimplemented (stub with `// TODO: implement`).
- [ ] **Signature mismatch in service layer**: The Go `UpdateUserOneSidedRelationships` method signature `(userID int64, relatedUserIDs []int64, blockDate, groupIndex, deleteGroupIndex, name, lastUpdatedDate)` differs fundamentally from the Java version which takes `Set<UserRelationship.Key>` (ownerId + relatedUserId pairs from **potentially multiple different owners**). The Go version is limited to a single owner. The controller receives a list of keys that can span multiple owners — the Go service cannot handle this correctly without iterating per-owner.
- [ ] The Go service method `UpdateUserOneSidedRelationships` takes extra parameters (`groupIndex`, `deleteGroupIndex`, `lastUpdatedDate`) that the Java controller does not pass. The Java controller only passes `name`, `blockDate`, and `establishmentDate`. The Go service signature's additional parameters suggest it may be a conflated version of multiple Java methods, which could lead to incorrect behavior if the controller passes nil/null values inappropriately or if the service applies group index changes that should not occur.
- [ ] The Java service method's `establishmentDate` parameter maps to the `lastUpdatedDate` parameter in Go, but the naming mismatch suggests potential semantic confusion. The Java controller passes `updateRelationshipDTO.establishmentDate()` directly as the `establishmentDate` param to `updateUserOneSidedRelationships`, which then passes it to the repository. In Go, this field would need to map to `lastUpdatedDate` — this naming is misleading and could cause bugs if someone maps it to the wrong parameter.

## deleteRelationships
- [ ] Controller method `DeleteRelationships()` is completely unimplemented (stub with `// TODO: implement`).
- [ ] **Signature mismatch in service layer**: The Go `DeleteOneSidedRelationships` method takes `(ownerID int64, relatedUserIDs []int64, session)` — a single ownerID — whereas the Java version takes `Set<UserRelationship.Key>` with potentially multiple different owners. The controller receives `List<UserRelationship.Key>` keys that can span multiple owners. The Go service cannot handle this correctly for multi-owner deletion scenarios without being called multiple times (once per distinct owner).
- [ ] **Missing cache invalidation**: The Java `deleteOneSidedRelationships` service method calls `invalidateRelationshipCache(keys)` after successful deletion. The Go version calls `s.invalidMemberCache(ownerID, relatedUserID)` per related user, which is functionally similar but the ordering/transaction semantics differ — Java invalidates cache in a `doOnSuccess` after the version update, while Go invalidates before the version update.
- [ ] **Missing error handling for version update**: The Java version wraps `updateRelationshipsVersion` in `.onErrorResume()` to log and swallow errors so a version update failure doesn't fail the entire operation. The Go version returns the error from `UpdateRelationshipsVersion` directly, which means a version update failure will propagate up and fail the delete operation — a behavioral difference from Java.

# UserRelationshipGroupController.java
*Checked methods: addRelationshipGroup(@RequestBody AddRelationshipGroupDTO addRelationshipGroupDTO), deleteRelationshipGroups(@QueryParam(required = false), updateRelationshipGroups(List<UserRelationshipGroup.Key> keys, @RequestBody UpdateRelationshipGroupDTO updateRelationshipGroupDTO), queryRelationshipGroups(@QueryParam(required = false), queryRelationshipGroups(@QueryParam(required = false)*

Now I have all the information needed for a thorough comparison. Let me compile the findings.

## addRelationshipGroup

**Java flow:** Controller passes `ownerId`, `index` (nullable Integer), `name`, `creationDate` (nullable Date), `null` (session) to `service.createRelationshipGroup()`. Service uses `RandomUtil.nextPositiveInt()` for random index when null.

**Go flow:** Controller is a `// TODO: implement` stub. The service method `CreateRelationshipGroup` exists and is mostly correct, but uses `int32(time.Now().UnixNano())` instead of a proper positive random int — this produces a large integer that can be negative (int32 overflow) instead of a positive random value. Also, the controller stub calls no service method at all.

## deleteRelationshipGroups

**Java flow:** Controller checks `CollectionUtil.isEmpty(keys)` — if empty, calls `service.deleteRelationshipGroups()` (no-arg, deletes ALL). Otherwise calls `service.deleteRelationshipGroups(CollectionUtil.newSet(keys))` which validates keys and calls `repo.deleteByIds()`.

**Go flow:** Controller is a `// TODO: implement` stub. The service has `DeleteRelationshipGroups` requiring `ownerID` and `groupIndexes`, which is a different signature — the Java version deletes by composite Key set, while the Go version deletes by single ownerID + groupIndexes. There's no "delete all" path in the service for the controller's use case.

## updateRelationshipGroups

**Java flow:** Controller passes `keys` (as Set), `name`, `creationDate` to service. Service validates keys, validates each key with `DataValidator.validRelationshipGroupKey()`, then calls `repo.updateRelationshipGroups()` which uses `$set` with `setIfNotNull` for both `name` and `creationDate`.

**Go flow:** Controller is a `// TODO: implement` stub. The service `UpdateRelationshipGroups` accepts `newName *string` and `creationDate *time.Time` but only updates name — **`creationDate` is never applied to the update query**, even though it's validated and checked for nil. The repository `UpdateRelationshipGroups` only accepts `newName` and doesn't handle `creationDate` at all. Also missing: per-key validation (`DataValidator.validRelationshipGroupKey`), and returning `UpdateResultDTO` (returns `nil` instead of acknowledged result when both fields are null).

## queryRelationshipGroups (non-paginated)

**Java flow:** Controller passes `ownerIds`, `indexes`, `names`, `DateRange.of(creationDateStart, creationDateEnd)`, `0` (page), `size` to service. Repository filters by `ownerIds`, `indexes`, `names`, AND `creationDateRange` (date range filter).

**Go flow:** Controller is a `// TODO: implement` stub. The service `QueryRelationshipGroups` only accepts `ownerIDs` and `groupIndexes` — **missing `names` and `creationDateRange` filter parameters entirely**. The repository `FindRelationshipGroups` also only filters by `ownerIDs` and `groupIndexes`, missing `names` and `creationDateRange` filters.

## queryRelationshipGroups (paginated - "page" endpoint)

**Java flow:** Calls `countRelationshipGroups` with all filters, then `queryRelationshipGroups` with all filters plus pagination. Returns `PaginationDTO`.

**Go flow:** Controller is a `// TODO: implement` stub. The service `CountRelationshipGroups` only accepts `ownerIDs` — **missing `indexes`, `names`, `creationDateRange` filter parameters**. Same missing filters in `QueryRelationshipGroups`.

Here is the full bug report:

---

## addRelationshipGroup

- [ ] **Controller is unimplemented stub**: The Go controller method `AddRelationshipGroup()` at `user_controllers.go:154` is an empty `// TODO: implement` stub. It does not call the service layer at all, making the endpoint non-functional.

- [ ] **Random group index generation is incorrect**: In `user_relationship_group_service.go:80`, when `groupIndex` is nil, Go uses `int32(time.Now().UnixNano())` which can produce negative values due to int32 overflow. The Java version uses `RandomUtil.nextPositiveInt()` which generates a positive random integer. The Go version should use a proper positive random int generator.

## deleteRelationshipGroups

- [ ] **Controller is unimplemented stub**: The Go controller method `DeleteRelationshipGroups()` at `user_controllers.go:159` is an empty `// TODO: implement` stub. It does not call the service layer at all, making the endpoint non-functional.

- [ ] **Missing "delete all" path**: The Java controller branches: if `keys` is empty, it calls `deleteRelationshipGroups()` (no-arg) which deletes ALL relationship groups via `repo.deleteAll()`. If keys are present, it calls `deleteRelationshipGroups(Set<Key>)`. The Go service has no equivalent of the no-arg delete-all path for this controller endpoint. The service's `DeleteRelationshipGroups` requires a single `ownerID` and `groupIndexes`, which doesn't match the controller's Key-set-based deletion semantics.

- [ ] **Signature mismatch for key-based deletion**: Java's `deleteRelationshipGroups(Set<UserRelationshipGroup.Key> keys)` accepts a set of composite keys (ownerId + groupIndex pairs) and calls `repo.deleteByIds(keys)`. The Go service `DeleteRelationshipGroups` takes a single `ownerID int64` and `groupIndexes []int32`, which is a fundamentally different API that cannot handle keys from different owners in a single call.

## updateRelationshipGroups

- [ ] **Controller is unimplemented stub**: The Go controller method `UpdateRelationshipGroups()` at `user_controllers.go:164` is an empty `// TODO: implement` stub. It does not call the service layer at all, making the endpoint non-functional.

- [ ] **creationDate is never applied to the database update**: In `user_relationship_group_service.go:240-277`, the service accepts `creationDate *time.Time` and validates it, but never passes it to the repository. The repository method `UpdateRelationshipGroups` (at `user_relationship_group_repository.go:130`) only accepts and sets `newName`. The Java repository uses `setIfNotNull(UserRelationshipGroup.Fields.CREATION_DATE, creationDate)` to also update the creation date when provided. The Go version silently ignores the `creationDate` parameter.

- [ ] **Missing per-key validation**: Java's service iterates over each key and calls `DataValidator.validRelationshipGroupKey(key)` before proceeding. The Go service only calls `validator.NotEmpty(keys, "keys")` but does not validate individual key fields (e.g., that ownerId and groupIndex are valid).

- [ ] **Missing acknowledged result return for no-op case**: In Java, when both `name` and `creationDate` are null, the service returns `OperationResultPublisherPool.ACKNOWLEDGED_UPDATE_RESULT` (a successful acknowledged update result). The Go service returns `nil` (line 256), which is a different semantic — callers cannot distinguish a successful no-op from an error.

## queryRelationshipGroups (non-paginated, GET endpoint)

- [ ] **Controller is unimplemented stub**: The Go controller method `QueryRelationshipGroups()` at `user_controllers.go:169` is an empty `// TODO: implement` stub. It does not call the service layer at all, making the endpoint non-functional.

- [ ] **Missing `names` filter parameter**: The Java controller accepts `@QueryParam Set<String> names` and passes it to the service and repository. The Go service `QueryRelationshipGroups` (at `user_relationship_group_service.go:562`) and repository `FindRelationshipGroups` (at `user_relationship_group_repository.go:202`) do not accept or filter by names at all.

- [ ] **Missing `creationDateRange` filter parameter**: The Java controller accepts `creationDateStart` and `creationDateEnd` query params, combines them into a `DateRange`, and passes it to the service and repository. The Go service and repository have no creationDate range filtering capability.

- [ ] **Missing page parameter for non-paginated query**: The Java non-paginated query passes `page=0` and `size` to the service. The Go service `QueryRelationshipGroups` does accept `page *int` and `size *int`, but without the controller implementation, this cannot be called correctly.

## queryRelationshipGroups (paginated, GET "page" endpoint)

- [ ] **Missing paginated query endpoint entirely**: The Java controller has a separate `@GetMapping("page")` method that performs a count query followed by a paginated data query, returning `PaginationDTO`. The Go controller has no equivalent method — the second `queryRelationshipGroups` overload is completely absent.

- [ ] **`CountRelationshipGroups` missing filter parameters**: The Java `countRelationshipGroups` accepts `ownerIds`, `indexes`, `names`, and `creationDateRange`. The Go `CountRelationshipGroups` at `user_relationship_group_service.go:551` only passes `ownerIDs` to the repository, missing `indexes`, `names`, and `creationDateRange` filters.

# AddUserDTO.java
*Checked methods: AddUserDTO(Long id, @SensitiveProperty(SensitiveProperty.Access.ALLOW_DESERIALIZATION), toString()*

Now I have the full picture. Let me summarize the analysis of the three specific methods/features being checked:

## AddUserDTO(Long id, ...)

**Java fields**: `id`, `password`, `name`, `intro`, `profilePicture`, `profileAccessStrategy`, `roleId`, `registrationDate`, `isActive`

**Go fields**: `ID`, `Password`, `Name`, `Intro`, `ProfilePicture`, `ProfileAccess`, `PermissionGroupID`, `RegistrationDate`, `IsActive`

The Java field `roleId` is mapped to `PermissionGroupID` in Go with JSON tag `"permissionGroupId"`. This is a **renamed field** — the JSON tag doesn't match the Java record component name `roleId`. Java records use the component name for Jackson deserialization by default, so the JSON key should be `"roleId"`, not `"permissionGroupId"`.

## @SensitiveProperty(ALLOW_DESERIALIZATION)

In Java, `@SensitiveProperty(ALLOW_DESERIALIZATION)` means:
- **Deserialization**: ALLOWED (password can be read from JSON input)
- **Serialization**: DENIED (password is NOT written to JSON output)

This is implemented via the `AnnotationIntrospector` in `JsonCodecPool.java` (lines 74-98). During serialization, fields annotated with `@SensitiveProperty` where `value != ALLOW_SERIALIZATION` are ignored. So `password` with `ALLOW_DESERIALIZATION` will be **ignored during serialization** (not written to JSON responses).

In Go, the field is `Password *string json:"password,omitempty"` — this means the password CAN be serialized in JSON responses when non-nil, which leaks the password.

## toString()

Java's `toString()` explicitly masks `password` with `SecurityValueConst.SENSITIVE_VALUE` (`"***"`). The Go `AddUserDTO` has no `String()` method at all.

Here are the findings:

## AddUserDTO

- [ ] **Field renamed from `roleId` to `permissionGroupId`**: The Java record has field `Long roleId` (Jackson will use `"roleId"` as the JSON key). The Go struct maps this to `PermissionGroupID *int64` with JSON tag `"permissionGroupId"`. The JSON key mismatch means API requests using `"roleId"` (as the Java server expects) will not populate this field in Go.

## @SensitiveProperty(ALLOW_DESERIALIZATION)

- [ ] **Password not excluded from JSON serialization**: The Java `@SensitiveProperty(ALLOW_DESERIALIZATION)` prevents the `password` field from being included in JSON serialization (responses). The Go field `Password *string json:"password,omitempty"` allows the password to be serialized in JSON responses when non-nil, leaking the password to clients. The Go code needs a mechanism to prevent password serialization (e.g., custom `MarshalJSON`, or using a separate write-only type).

## toString()

- [ ] **Missing `String()` / `fmt.Stringer` implementation with password masking**: The Java `toString()` masks the `password` field with `SecurityValueConst.SENSITIVE_VALUE` (`"***"`). The Go `AddUserDTO` has no `String()` method, meaning `fmt.Sprintf("%v", dto)` or `%+v` will print the raw password in plaintext, which is a security risk in logging/debugging scenarios.
