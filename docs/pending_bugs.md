
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

- [ ] **Missing `timestamp` parameter**: The Java method signature is `tryAcquireToken(ByteArrayWrapper ip, long timestamp)` and passes `timestamp` to `bucket.tryAcquire(timestamp)` for refill calculation. The Go version `TryAcquireToken(ip string)` drops the `timestamp` parameter entirely and uses `rate.Limiter.Allow()` which uses `time.Now()` internally. This means the Go version cannot honor caller-provided timestamps for refill computation, changing behavior when timestamps are externally controlled or batched.

- [ ] **Different rate limiting algorithm**: Java uses a custom `TokenBucket` with explicit token counting, CAS-based refill logic, configurable `capacity`, `tokensPerPeriod`, and `refillIntervalNanos`. Go uses `golang.org/x/time/rate.Limiter` which is a token bucket with a continuous refill model (not discrete period-based). This produces different throttling behavior: Java refills tokens in discrete batches per period (e.g., 10 tokens every 1 second), while Go's `rate.Limiter` adds tokens continuously at a steady rate.

- [ ] **Shared context vs. static configuration**: In Java, all `TokenBucket` instances share a single `TokenBucketContext` that can be dynamically updated at runtime (via `propertiesManager.addGlobalPropertiesChangeListener`). The Go version stores `Limit` and `Burst` as struct fields at construction time, and new limiters created after an update would use the updated values, but **existing per-IP limiters are never updated** when `Limit`/`Burst` fields change. Java's shared context means all buckets immediately reflect new rate limiting settings.

- [ ] **Unlimited condition is wrong**: Go returns `true` (unlimited) when `t.Burst <= 0 || t.Limit == 0`. In Java, whether the bucket is "unlimited" depends on the `TokenBucketContext` configuration — specifically, if `refillIntervalNanos <= 0`, the bucket returns `false` when empty (no refill). When `capacity` and `tokensPerPeriod` allow unlimited access, it's the initial token count that matters. The Go condition `t.Limit == 0` treating zero limit as "unlimited" is backwards — `rate.Limit == 0` means zero rate (no tokens), not infinite rate. And `Burst <= 0` as unlimited is also incorrect since `rate.NewLimiter` with burst 0 would allow zero requests.

- [ ] **Cleanup logic is fundamentally different**: Java's cleanup iterates entries and only removes those that are both idle for 30+ minutes **and** have tokens >= initial tokens (i.e., fully replenished). Go's cleanup does a full map reset every 10 minutes, wiping **all** entries including actively-used ones. This means active IPs lose their rate limiting state and start fresh every 10 minutes, which is a behavioral difference.

- [ ] **Missing session-closed listener**: Java registers a listener on `SessionService` that removes an IP's token bucket when a session closes (if tokens are replenished). The Go version has no equivalent mechanism for cleaning up on session close.

# NotificationFactory.java
*Checked methods: init(TurmsPropertiesManager propertiesManager), create(ResponseStatusCode code, long requestId), create(ResponseStatusCode code, @Nullable String reason, long requestId), create(ThrowableInfo info, long requestId), createBuffer(CloseReason closeReason), sessionClosed(long requestId)*

Now I have a complete picture. Let me compile the bugs.

## init (NewNotificationFactory)

- [ ] **Missing dynamic config updates**: The Java `init` method registers a `notifyAndAddLocalPropertiesChangeListener` that dynamically updates `returnReasonForServerError` whenever properties change at runtime. The Go `NewNotificationFactory` captures a static `*config.GatewayProperties` snapshot at construction time and never updates it. If the `ReturnReasonForServerError` property changes at runtime, the Go version will not reflect the change.

## Create (create(ResponseStatusCode code, long requestId))

- [ ] **Missing default reason from status code**: The Java `create(code, requestId)` calls `trySetReason(builder, code, code.getReason())`, which passes the status code's built-in default reason (e.g., "ok" for OK, "The client request is invalid" for INVALID_REQUEST). The Go `Create(requestID, code)` delegates to `CreateWithReason(requestID, code, "")`, passing an empty string instead of the status code's default reason. Since `trySetReason` returns early when `reason == ""`, **no reason is ever set** for any status code in this method. In Java, non-server-error codes would always have their default reason included.

## CreateWithReason (create(ResponseStatusCode code, @Nullable String reason, long requestId))

- [ ] **Empty-string vs nil/null semantics mismatch**: The Java version uses `@Nullable String reason` where `null` triggers the fallback to `code.getReason()`: `reason == null ? code.getReason() : reason`. The Go version uses `reason string` (empty string `""` as zero value). When a caller passes an empty reason, the Go code treats it like Java's `null` (no reason set). But the Java version would still set `code.getReason()` as the reason even when the explicit reason is non-null but empty. More critically, when the Java `reason` parameter is null, it falls back to `code.getReason()` (the default reason for that status code). The Go version has no such fallback — it just uses the empty string directly.

## CreateFromError (create(ThrowableInfo info, long requestId))

- [ ] **Wrong default error code for non-TurmsError errors**: In Java, `create(ThrowableInfo info, long requestId)` always extracts `info.code()` from the `ThrowableInfo` record, which already contains the correct `ResponseStatusCode` (resolved by `ThrowableInfo.get(Throwable)`). The Go version defaults to `ResponseStatusCode_SERVER_INTERNAL_ERROR` for non-`TurmsError` errors and falls back to the generic error message. While this is architecturally different (Go uses `error` interface vs Java's `ThrowableInfo` record), it means any custom error types with specific status codes (like the Java equivalents of `RECORD_CONTAINS_DUPLICATE_KEY`, `RESOURCE_NOT_FOUND`, etc.) will all map to `SERVER_INTERNAL_ERROR` instead of their proper codes.
- [ ] **Missing fallback to code's default reason**: The Java version passes `info.reason()` to `trySetReason`, which may be `null` — in which case `trySetReason` returns without setting a reason. However, the Java version could have a non-null reason from `ThrowableInfo`. The Go version sets `reason = err.Error()` for non-TurmsError errors, which is a reasonable but different behavior (Java would use the throwable's message via `ThrowableInfo.get()`). For TurmsError cases, `te.Message` is used which maps to `info.reason()`, which is correct.

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
- [ ] **Missing IP assignment**: The Java version assigns both `this.connection = connection` and `this.ip = ip`, but the Go version at `connection.go:92-94` assigns `s.Conn = connection` but never assigns the `ip` parameter to `s.IP`. The `ip string` parameter is received but discarded.

## setLastHeartbeatRequestTimestampToNow()
- [ ] **Missing nanosecond timestamp tracking**: The Java version updates both `lastHeartbeatRequestTimestampMillis` (via `System.currentTimeMillis()`) and `lastHeartbeatRequestTimestampNanos` (via `System.nanoTime()`). The Go version at `connection.go:41-43` only stores a millisecond timestamp in `lastHeartbeat`. The nanosecond timestamp (`lastHeartbeatRequestTimestampNanos`) is not tracked at all.

## setLastRequestTimestampToNow()
- [ ] **Missing nanosecond timestamp tracking**: Same as above. The Java version updates both `lastRequestTimestampMillis` and `lastRequestTimestampNanos`. The Go version at `connection.go:52-54` only stores a millisecond timestamp in `lastRequest`. The nanosecond timestamp (`lastRequestTimestampNanos`) is not tracked at all.

## close(@NotNull CloseReason closeReason)
- [ ] **Missing `isSessionOpen` state tracking**: The Java version maintains a separate `isSessionOpen` volatile boolean that is set to `false` on close, and the method returns `true` only if the session was previously open (i.e., the first close succeeds). The Go version at `connection.go:125-129` has no `isSessionOpen` equivalent — it only checks `s.Conn != nil` and does not track whether the session has already been closed.
- [ ] **Missing return value**: The Java version returns `boolean` indicating whether the session was actually closed (was open). The Go version returns nothing (`void`).
- [ ] **Missing close-reason propagation**: The Java version passes `closeReason` to `connection.close(closeReason)`. The Go version ignores the `closeReason` parameter entirely and calls `s.Conn.Close()` with no arguments.
- [ ] **Missing warning log when connection is null**: The Java version logs a warning `"The connection is missing for the user session: {}"` when `isSessionOpen` is true but `connection == null`. The Go version silently does nothing when `Conn` is nil.

## isOpen()
- [ ] **Wrong semantics**: The Java version at line 175-177 returns the `isSessionOpen` volatile boolean, which tracks whether the session is open (independently of whether a connection exists — a session can be open with UDP heartbeats even without a connection). The Go version at `connection.go:63-65` returns `s.Conn != nil`, which checks for the presence of a connection, not session openness. These are semantically different: after `close()` is called, Java returns `false` (session closed) while Go would still return `true` if the connection object hasn't been nulled out.

## isConnected()
- [ ] **Missing `connection.isConnected()` check**: The Java version at line 179-181 returns `connection != null && connection.isConnected()` — it checks both that the connection exists AND that it is connected. The Go version at `connection.go:97-99` only checks `s.Conn != nil`, without calling any `IsActive()` or equivalent method on the connection.

## supportsSwitchingToUdp()
- No bugs. The Go version at `connection.go:102-104` correctly checks `s.DeviceType != protocol.DeviceType_BROWSER`, matching the Java logic `deviceType != DeviceType.BROWSER`.

## sendNotification(ByteBuf byteBuf)
- [ ] **Method is completely missing**: The Java version has `sendNotification(ByteBuf byteBuf)` that calls `notificationConsumer.apply(byteBuf, TracingContext.NOOP)`. There is no corresponding method on `UserSession` in the Go code. The `sendNotification` in `router.go` is a completely different method on the `Router` struct that creates a new notification from scratch, rather than forwarding a pre-built ByteBuf via a consumer function.

## sendNotification(ByteBuf byteBuf, TracingContext tracingContext)
- [ ] **Method is completely missing on UserSession**: The Java version stores a `BiFunction<ByteBuf, TracingContext, Mono<Void>> notificationConsumer` field and uses it in `sendNotification`. The Go version has no `notificationConsumer` field and no `SendNotification` method on `UserSession`. The `sendNotification` in `router.go:135-140` is a `Router` method that creates notifications via a factory — it is architecturally different from the Java version which forwards pre-built ByteBuf notifications from turms-service servers.

## acquireDeleteSessionRequestLoggingLock()
- No bugs. The Go version at `connection.go:112-114` correctly uses `atomic.CompareAndSwapUint32(&s.isDeleteSessionLockAcquired, 0, 1)`, matching the Java `AtomicIntegerFieldUpdater` compareAndSet logic.

## hasPermission(TurmsRequest.KindCase requestType)
- [ ] **Changed behavior on nil map**: The Java version uses a `Set<TurmsRequest.KindCase>` and calls `permissions.contains(requestType)`, which would throw `NullPointerException` if `permissions` were null (it's never null due to constructor). The Go version at `connection.go:117-122` adds a nil-check for `s.Permissions` and returns `false`. While this is arguably more defensive, the behavior differs — Java would crash on null permissions, Go silently returns false. More importantly, the Go `map[any]bool` type is less type-safe than Java's `Set<TurmsRequest.KindCase>`, but this is an idiomatic Go trade-off rather than a bug.

## toString()
- [ ] **Stub implementation — missing all fields**: The Java version at lines 206-225 returns a string with `id`, `version`, `userId`, `deviceType`, `loginDate`, `loginLocation`, `isSessionOpen`, and `connection`. The Go version at `connection.go:107-109` returns the hardcoded literal `"UserSession{}"`, omitting all field values entirely.

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
