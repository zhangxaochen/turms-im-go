
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
