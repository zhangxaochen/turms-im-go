
### PolicyDeserializer.java
Checked methods: parse(Map<String, Object> map)
## Bugs in Go `Parse` method

The Go implementation at `policy.go:75-78` is a **stub that returns an empty `Policy` with zero logic**. Comparing to the Java `parse(Map<String, Object> map)` method (lines 80-105), the following core logic is entirely missing:

### 1. **Missing: Extraction and validation of `statements` field**
The Java code extracts `map.get(Policy.Fields.statements)`, validates it's a `List<?>`, and throws `IllegalPolicyException` if not. The Go code does none of this.

### 2. **Missing: Iteration over statement objects**
The Java code iterates each `statementObj` in the list, validates each is a `Map<?, ?>`, and throws an error if not. The Go code does none of this.

### 3. **Missing: Parsing `effect` from each statement**
The Java code calls `parseEffect(statementMap.get(PolicyStatement.Fields.effect))` which validates the value is a string matching `"ALLOW"` or `"DENY"`, throwing on invalid input. The Go code does none of this.

### 4. **Missing: Parsing `actions` from each statement**
The Java code calls `parseActions(...)` which handles:
- A single string `"*"` mapping to `PolicyStatementAction.ALL`, or a specific action string (`"CREATE"`, `"DELETE"`, `"UPDATE"`, `"QUERY"`)
- An array of action strings
- Error on invalid types or unknown action strings

The Go code does none of this.

### 5. **Missing: Parsing `resources` from each statement**
The Java code calls `parseResources(...)` which handles:
- A single string `"*"` returning all resources
- A specific resource name string parsed via `PolicyStatementResource.valueOf()`
- An array of resource strings (with early return on `"*"`)
- Error on invalid types or unknown resource strings

The Go code does none of this.

### 6. **Missing: Construction of `PolicyStatement` and `Policy`**
The Java code constructs `new PolicyStatement(effect, actions, resources)` for each statement and collects them into a list, then returns `new Policy(statements)`. The Go code returns `&Policy{}` (empty, no statements).

### 7. **Missing: Error handling / propagation**
The Java code throws `IllegalPolicyException` with descriptive messages for every malformed input case. The Go code returns `nil` error unconditionally.

**Summary:** The Go `Parse` method is an unimplemented stub (`// Dummy parser pending full implementation`) that returns an empty `Policy` and `nil` error regardless of input, missing all parsing, validation, and error-handling logic from the Java original.

### ServiceAvailabilityHandler.java
Checked methods: channelRegistered(ChannelHandlerContext ctx), exceptionCaught(ChannelHandlerContext ctx, Throwable cause)
Now I have a complete picture. Here are the bugs:

## Bugs

### 1. `HandleConnection` — Availability check is hardcoded to `true` instead of calling `serverStatusManager`

**Java (line 58-59):**
```java
ServiceAvailability serviceAvailability = serverStatusManager.getServiceAvailability();
if (serviceAvailability.isAvailable()) {
```

**Go (line 31-32):**
```go
// serviceAvailability := h.serverStatusManager.GetServiceAvailability() // Not fully migrated
available := true
```

The `ServerStatusManager` interface already has a `GetServiceAvailability()` method defined (in `client_request_dispatcher.go:64`), and the `ServiceAvailability` type has an `IsAvailable()` method (in `service_availability.go:38-40`). The real call is commented out and replaced with a hardcoded `true`, meaning the server will **never reject connections when unavailable** (e.g., during shutdown or before startup).

### 2. `HandleException` — Missing `OutOfDirectMemoryError` handling

**Java (lines 92-94):**
```java
} else if (cause instanceof OutOfDirectMemoryError) {
    ctx.close();
}
```

**Go:** The entire `else if` branch for `OutOfDirectMemoryError` is absent. In the Java version, when the server runs out of direct memory, the connection is explicitly closed. The Go version silently ignores this error condition entirely.

### 3. `HandleException` — Missing `ctx.fireExceptionCaught(cause)` propagation

**Java (line 95):**
```java
ctx.fireExceptionCaught(cause);
```

**Go:** The Java code **always** propagates the exception downstream via `ctx.fireExceptionCaught(cause)` — this call is unconditional (outside the `if/else if` blocks). The Go version has no equivalent; exceptions are never forwarded to downstream handlers in the pipeline.

### NetConnection.java
Checked methods: getAddress(), send(ByteBuf buffer), close(CloseReason closeReason), close(), switchToUdp(), tryNotifyClientToRecover()
## Bugs Found

### 1. `SwitchToUdp()` uses wrong close status — **Behavioral Bug**

**Java** (`switchToUdp`):
```java
public void switchToUdp() {
    close(CloseReason.get(SessionCloseStatus.SWITCH));
}
```

**Go** (`SwitchToUdp`):
```go
func (b *BaseNetConnection) SwitchToUdp() {
    b.CloseWithReason(NewCloseReason(constant.SessionCloseStatus_SERVER_CLOSED))
}
```

The Java version passes `SessionCloseStatus.SWITCH`, but the Go version passes `SessionCloseStatus_SERVER_CLOSED`. This is a **critical difference**: with the wrong status, `CloseWithReason` will set `isSwitchingToUdp = false` instead of `true`, since `isSwitchingToUdp` is only set to `true` when the status equals `SWITCH`. This completely breaks the UDP switching logic.

---

### 2. `GetAddress()` and `Send()` are declared in the interface but have no implementation on `BaseNetConnection` — **Structural Difference**

In Java, `getAddress()` and `send(ByteBuf buffer)` are **abstract methods** on the class, meaning concrete subclasses must implement them. In Go, `GetAddress()` and `Send()` are declared in the `NetConnection` interface, but `BaseNetConnection` does **not** provide stub or forwarding implementations. This means any concrete type embedding `BaseNetConnection` must independently satisfy the full interface — this is architecturally fine for Go, but note that `BaseNetConnection` itself does not satisfy the `NetConnection` interface (unlike the Java abstract class which partially implements it).

---

### 3. `TryNotifyClientToRecover()` — actual UDP dispatch is commented out / missing

**Java**:
```java
UdpRequestDispatcher.instance.sendSignal(udpAddress, UdpNotificationType.OPEN_CONNECTION);
```

**Go**:
```go
// Mock logic: UdpRequestDispatcher.instance.sendSignal
// udpdispatcher.SendSignal(b.udpAddress, udp.OPEN_CONNECTION)
```

The actual signal-sending call is entirely missing (commented out). The method only sets `isConnectionRecovering = true` without ever notifying the client. This is a **missing core logic** bug — the method has no real effect.

---

### Summary

| Method | Bug |
|---|---|
| `SwitchToUdp()` | Uses `SERVER_CLOSED` instead of `SWITCH` — sets `isSwitchingToUdp` to `false`, breaking the entire UDP switch flow |
| `TryNotifyClientToRecover()` | The `sendSignal` call is commented out / unimplemented — the method silently does nothing besides setting a flag |
| `GetAddress()` / `Send()` | No `BaseNetConnection` implementation (acceptable Go pattern, but departs from Java's abstract class design) |

### ExtendedHAProxyMessageReader.java
Checked methods: channelRead(ChannelHandlerContext ctx, Object msg)
## Bugs found in `Read` (mapped from `channelRead`)

1. **Missing PROXY protocol parsing logic entirely**  
   The Java version receives an `HAProxyMessage` and extracts `sourceAddress` and `sourcePort` from it via `proxyMessage.sourceAddress()` and `proxyMessage.sourcePort()`. It creates an unresolved `InetSocketAddress` from those values and passes that to the callback. The Go version skips all of this — it never reads or parses the PROXY protocol header (v1 or v2) from the wire. It immediately calls the callback with `conn.RemoteAddr()`, which is the direct TCP remote address, not the address reported by the proxy.

2. **Missing fallback logic to `ctx.channel().remoteAddress()`**  
   In Java, when `remoteAddress` is `null` (i.e., the PROXY header had no valid source address/port), the code falls back to the channel's actual remote address: `ctx.channel().remoteAddress()`. The Go code unconditionally uses `conn.RemoteAddr()` — while this happens to match the fallback path, it does so because the primary path (parsing the PROXY header) is entirely missing.

3. **Missing `proxyMessage.release()` equivalent (resource cleanup)**  
   The Java code calls `proxyMessage.release()` in a `finally` block to release the reference-counted HAProxy message. The Go version has no equivalent cleanup of the parsed proxy message buffer.

4. **Missing self-removal from pipeline**  
   The Java code removes `this` handler from the pipeline after processing the proxy message (`ctx.channel().pipeline().remove(this)`), because the PROXY protocol header is only sent once at connection start. The Go version does not have any equivalent "read once and detach" semantics.

5. **Missing `ctx.read()` re-trigger**  
   After processing the proxy message and removing itself, the Java code calls `ctx.read()` to resume reading subsequent data on the channel. The Go version has no equivalent — the `Read` method does not trigger continued reading of the connection stream after the proxy header is consumed.

6. **Missing non-HAProxy message passthrough**  
   The Java `channelRead` has an `else` branch: if `msg` is not an `HAProxyMessage`, it calls `super.channelRead(ctx, msg)` to pass the message up the pipeline. The Go version only handles the proxy case and has no equivalent passthrough for non-proxy data.

### HAProxyUtil.java
Checked methods: addProxyProtocolHandlers(ChannelPipeline pipeline, Consumer<InetSocketAddress> onRemoteAddressConfirmed), addProxyProtocolDetectorHandler(ChannelPipeline pipeline, Consumer<InetSocketAddress> onRemoteAddressConfirmed)
## Bugs

1. **`AddProxyProtocolHandlers` is an empty stub with no implementation.** The Java version adds two handlers to a pipeline: an `HAProxyMessageDecoder` (named `ProxyProtocolDecoder`) placed first, followed by an `ExtendedHAProxyMessageReader` (named `ProxyProtocolReader`) placed after the decoder. The Go function body is a comment stating "Pending implementation" — it performs zero work and ignores the callback entirely.

2. **`AddProxyProtocolDetectorHandler` is an empty stub with no implementation.** The Java version adds an `ExtendedHAProxyMessageDetector` handler (named `ProxyProtocolDecoder`) as the first handler in the pipeline. The Go function body is a comment stating "Pending implementation" — it performs zero work and ignores the callback entirely.

3. **`AddProxyProtocolHandlers` and `AddProxyProtocolDetectorHandler` are missing the pipeline parameter.** The Java signatures accept a `ChannelPipeline pipeline` argument that identifies which connection/channel to attach handlers to. The Go versions take only a `callback func(net.Addr)` with no equivalent of a pipeline or connection, making it structurally impossible to replicate the Java behavior even if implemented.

4. **`ExtendedHAProxyMessageReader.Read` does not actually parse PROXY protocol headers.** The Java `ExtendedHAProxyMessageReader` (which is `channelRead` in a Netty `ChannelInboundHandlerAdapter`) receives a decoded `HAProxyMessage` and extracts the source address from it (the address sent by the proxy, not the direct TCP remote address). The Go version simply calls `r.OnRemoteAddressConfirmed(conn.RemoteAddr())`, which returns the direct connection's remote address — this is the address *of the proxy*, not the *client behind the proxy*. This defeats the entire purpose of PROXY protocol parsing.

5. **Missing `ExtendedHAProxyMessageDetector` type.** The Java code references a class `ExtendedHAProxyMessageDetector` used by `addProxyProtocolDetectorHandler`. This type does not exist anywhere in the Go file. The detector is a distinct handler from the reader — it detects/validates the PROXY protocol header without fully decoding it, serving a different role in the pipeline.

### TcpConnection.java
Checked methods: getAddress(), send(ByteBuf buffer), close(CloseReason closeReason), close()
Now I have all the context needed. Here is the comparison:

## Bugs in `close(CloseReason closeReason)` → `CloseWithReason`

### 1. Missing guard: no check for `isConnected()` or connection disposed state

**Java** (`TcpConnection.java:71-73`):
```java
if (!isConnected() || connection.isDisposed()) {
    return;
}
```
The Java version returns early if the connection is already disconnected or disposed.

**Go** (`tcp_server.go:38-42`): No such guard exists. `CloseWithReason` proceeds unconditionally, which means:
- It will re-set state on an already-closed connection
- It will attempt to close a `net.Conn` that may already be closed, risking a double-close panic or error

### 2. Missing close notification send with retry logic

**Java** (`TcpConnection.java:75-82`): Before closing the connection, it sends a close notification to the client:
```java
Mono<Void> mono = connection.sendObject(NotificationFactory.createBuffer(closeReason))
        .then()
        .doOnError(...)
        .retryWhen(RETRY_SEND_CLOSE_NOTIFICATION);
```
This includes retry logic (up to 2 retries with 3s backoff) and error logging for non-disconnected-client errors.

**Go** (`tcp_server.go:39`): Only has a comment `// Pending logic to send notification before closing`. The notification is **never sent**. The client receives no close notification before the connection is torn down.

### 3. Missing `closeTimeout`-based scheduling logic

**Java** (`TcpConnection.java:83-90`): Has two timeout-based paths:
- **If `closeTimeout` is zero**: calls `close()` via `doFinally` immediately after the notification send completes.
- **If `closeTimeout` is positive (non-negative)**: waits for `connection.onTerminate()` with a `timeout(closeTimeout)`, falling back to `close()` on timeout or completion.

**Go** (`tcp_server.go:41`): Always calls `c.conn.Close()` immediately after `BaseNetConnection.CloseWithReason(reason)`. There is no timeout-based waiting for the notification to be acknowledged or for a graceful termination.

### 4. Missing error handling/logging in `close(CloseReason)`

**Java** (`TcpConnection.java:91-97`): Subscribes with an error handler that logs failures (filtering out disconnected-client errors):
```java
mono.subscribe(null, t -> {
    if (!ThrowableUtil.isDisconnectedClientError(t)) {
        LOGGER.error("Failed to send the close notification after (2) attempts", t);
    }
});
```

**Go**: No error logging at all in `CloseWithReason`.

## Bugs in `close()` → `Close()`

### 5. Missing error handling/logging in `close()`

**Java** (`TcpConnection.java:101-112`): Wraps `connection.dispose()` in a try-catch, logging errors (excluding disconnected-client errors) with the remote host address:
```java
try {
    connection.dispose();
} catch (Exception e) {
    if (!ThrowableUtil.isDisconnectedClientError(e)) {
        LOGGER.error("Failed to close the TCP connection: " + getAddress().getAddress().getHostAddress(), e);
    }
}
```

**Go** (`tcp_server.go:45-48`): Silently returns any error from `c.conn.Close()` with no logging.

## Bug in `SwitchToUdp()` (BaseNetConnection)

### 6. Wrong close status in `SwitchToUdp()`

**Java** (`NetConnection.java`):
```java
public void switchToUdp() {
    close(CloseReason.get(SessionCloseStatus.SWITCH));
}
```

**Go** (`net_connection.go`):
```go
func (b *BaseNetConnection) SwitchToUdp() {
    b.CloseWithReason(NewCloseReason(constant.SessionCloseStatus_SERVER_CLOSED))
}
```

Uses `SessionCloseStatus_SERVER_CLOSED` instead of `SessionCloseStatus_SWITCH`. This is incorrect — `switchToUdp` should set `isSwitchingToUdp = true`, but with `SERVER_CLOSED`, `isSwitchingToUdp` will always be `false`.

### TcpServerFactory.java
Checked methods: create(TcpProperties tcpProperties, BlocklistService blocklistService, ServerStatusManager serverStatusManager, SessionService sessionService, ConnectionListener connectionListener, int maxFrameLength)
## Bugs in Go port of `TcpServerFactory.create(...)`

The Go implementation at `tcp_server.go:54-60` is a **completely empty stub** — the `CreateWithArgs` method has no body logic whatsoever. Here are the specific missing pieces compared to the Java `create(...)` method:

### 1. **Missing: Entire TCP server creation and configuration** (Critical)
The Java code configures a `TcpServer` with host, port, socket options (`CONNECT_TIMEOUT_MILLIS`, `SO_REUSEADDR`, `SO_BACKLOG`, `SO_LINGER`, `TCP_NODELAY`), wiretap, event loop threads, and metrics. None of this exists in Go.

### 2. **Missing: `ServiceAvailabilityHandler` instantiation** (Critical)
Java creates a `ServiceAvailabilityHandler` with `blocklistService`, `serverStatusManager`, and `sessionService` and adds it as the first pipeline handler. Go does nothing with these three parameters.

### 3. **Missing: Channel pipeline / codec setup** (Critical)
Java adds:
- `serviceAvailabilityHandler` (first in pipeline)
- `varintLengthBasedFrameDecoder` (inbound, before reactive bridge, using `maxFrameLength`)
- `varintLengthFieldPrepender` (outbound)
- `protobufFrameEncoder` (outbound)

Go ignores `maxFrameLength` entirely and has no codec pipeline.

### 4. **Missing: Proxy protocol handling** (Critical)
Java handles three `RemoteAddressSource.ProxyProtocolMode` cases:
- `REQUIRED` — adds HAProxy proxy protocol handlers with blocklist IP check
- `OPTIONAL` — adds proxy protocol detector with blocklist IP check
- Default — uses `channel.remoteAddress()` directly

Go does none of this and does not use `blocklistService` at all.

### 5. **Missing: Connection handling / `remoteAddressSink` logic** (Critical)
Java creates a `Sinks.One<InetSocketAddress>` to asynchronously resolve the remote address, then in `.handle()` sets `autoRead(true)` and calls `connectionListener.onAdded(connection, remoteAddress, in.receive(), out, connection.onDispose())`. The Go method takes `connectionListener any` but never invokes it.

### 6. **Missing: SSL/TLS configuration** (Critical)
Java checks `tcpProperties.getSsl().isEnabled()` and configures SSL via `SslUtil.configureSslContextSpec(...)`. Go ignores SSL entirely.

### 7. **Missing: Server bind with error handling** (Critical)
Java calls `server.bind().block()` and wraps failures in a `BindException` with a descriptive message including host and port. Go does no binding at all.

### 8. **Missing: Parameter types are `any` instead of concrete types** (Design)
The Go method signature uses `any` for all parameters instead of concrete types (e.g., `TcpProperties`, `BlocklistService`, etc.), making it impossible to access any fields or call any methods on them.

### TcpUserSessionAssembler.java
Checked methods: getHost(), getPort()
## Bugs

### 1. `getHost()` / `GetHost()` — Missing disabled-server guard (throws `FeatureDisabledException`)

**Java** (`getHost()`, line 112-117):
```java
public String getHost() {
    if (server == null) {
        throw new FeatureDisabledException("TCP server is disabled");
    }
    return host;
}
```

**Go** (`GetHost()`, line 77-79):
```go
func (a *TcpUserSessionAssembler) GetHost() string {
    return a.Host
}
```

The Go version unconditionally returns `a.Host` without checking whether the TCP server is disabled (i.e., `Server == nil`). In Java, when the server is disabled, `host` is `null` and `server` is `null`, and calling `getHost()` throws a `FeatureDisabledException` rather than silently returning `null`. The Go code will return the zero-value `""` instead of signaling an error.

### 2. `getPort()` / `GetPort()` — Missing disabled-server guard (throws `FeatureDisabledException`)

**Java** (`getPort()`, line 119-124):
```java
public int getPort() {
    if (server == null) {
        throw new FeatureDisabledException("TCP server is disabled");
    }
    return port;
}
```

**Go** (`GetPort()`, line 82-84):
```go
func (a *TcpUserSessionAssembler) GetPort() int {
    return a.Port
}
```

Same issue as `GetHost()`. The Go version returns `a.Port` unconditionally. In Java, when the server is disabled, `port` is set to `-1` and calling `getPort()` throws a `FeatureDisabledException`. The Go code will silently return the zero-value `0` (or `-1` if explicitly set, but the constructor `NewTcpUserSessionAssembler` doesn't set `Port` to `-1` for the disabled case) instead of returning an error.

### 3. Constructor `NewTcpUserSessionAssembler` — Missing disabled-case field initialization

**Java** constructor (lines 100-104):
```java
} else {
    server = null;
    host = null;
    port = -1;
}
```

**Go** constructor (lines 70-74):
```go
func NewTcpUserSessionAssembler() *TcpUserSessionAssembler {
    return &TcpUserSessionAssembler{
        Enabled: false,
    }
}
```

When the server is disabled, the Java code explicitly sets `host = null` and `port = -1`. The Go constructor only sets `Enabled: false` and leaves `Host` as `""` (zero value) and `Port` as `0` (zero value). Even if the guard from bugs 1/2 were added, the sentinel values differ: Java uses `null`/`-1`, Go defaults to `""`/`0`.

### UdpRequestDispatcher.java
Checked methods: sendSignal(InetSocketAddress address, UdpNotificationType signal), get(ResponseStatusCode code), get(UdpNotificationType type)
Now I have the full picture. Let me compare the specific methods requested.

## Bug Report

### 1. `sendSignal` — Silent drop on full channel vs. `tryEmitNext` failure handling

**Java** (`sendSignal`):
```java
public void sendSignal(InetSocketAddress address, UdpNotificationType signal) {
    if (notificationSink != null) {
        notificationSink.tryEmitNext(new UdpNotification(address, signal));
    }
}
```

**Go** (`SendSignal`):
```go
func (d *UdpRequestDispatcher) SendSignal(address net.Addr, signal UdpNotificationType) {
    if d.notificationSink != nil {
        select {
        case d.notificationSink <- UdpNotification{
            RecipientAddress: address,
            Type:             signal,
        }:
        default:
            // Handle sink full
        }
    }
}
```

**Bug:** The Java version uses `tryEmitNext`, which will attempt to emit and returns a result code (e.g., `FAIL_TERMINATED`, `FAIL_OVERFLOW`) that the caller *could* inspect (though here it's ignored). The Go version silently drops the notification when the buffered channel is full via `default`. This is a reasonable semantic approximation but is technically a behavioral difference — Java's `Sinks.many().unicast().onBackpressureBuffer()` has an unbounded buffer by default, meaning `tryEmitNext` will essentially never fail due to backpressure (only if the sink is terminated). The Go version uses a bounded channel (`cap: 1024`), meaning it can silently drop notifications under load, while the Java version does not. **This is a behavioral difference in backpressure handling — the Java sink is unbounded, the Go channel is bounded at 1024.**

---

### 2. `get(ResponseStatusCode code)` — Missing entirely

The Java `UdpSignalResponseBufferPool.get(ResponseStatusCode code)` method:
- Lazily caches `ByteBuf` responses per status code
- Returns `EMPTY_BUFFER` for `OK`
- Returns a 2-byte buffer with `code.getBusinessCode()` as a short for all other codes
- Uses double-checked locking for thread safety

**The Go code has no equivalent of this method.** There is no `UdpSignalResponseBufferPool` or any function that serializes a `ResponseStatusCode` into bytes for UDP response packets.

---

### 3. `get(UdpNotificationType type)` — Missing entirely

The Java `UdpSignalResponseBufferPool.get(UdpNotificationType type)` method:
- Returns a pre-cached 1-byte `ByteBuf` containing `type.ordinal() + 1`
- Initialized in a static block for all enum values

**The Go code has no equivalent of this method.** There is no function that serializes a `UdpNotificationType` into a byte (ordinal + 1) for sending in UDP notification packets.

---

### 4. `UdpNotificationType` values mismatch

**Java:** `UdpNotificationType` has a single value: `OPEN_CONNECTION`

**Go:** `UdpNotificationType` has values: `HeartbeatNotification` (0), `GoOfflineNotification` (1)

These are completely different enum values with different ordinals, which would produce different on-the-wire bytes even if the buffer pool were implemented.

---

### 5. `UdpRequestType` values mismatch

**Java:** `UdpRequestType` has values: `HEARTBEAT`, `GO_OFFLINE` (these are parsed from incoming packets)

**Go:** `UdpRequestType` has values: `HeartbeatRequest` (0), `GoOfflineRequest` (1)

The Go code has `ParseUdpRequestType(number int)` which does a simple cast with **no validation** — the Java `UdpRequestType.parse(byte)` returns `null` for invalid values, while the Go version will happily create any `UdpRequestType` from any integer, producing an invalid enum value.

---

### Summary of Bugs

| # | Method/Area | Bug |
|---|---|---|
| 1 | `sendSignal` | Bounded channel (1024) vs Java's unbounded sink — notifications can be silently dropped |
| 2 | `get(ResponseStatusCode)` | **Missing entirely** — no buffer pool or serialization for response status codes |
| 3 | `get(UdpNotificationType)` | **Missing entirely** — no buffer pool or serialization for notification types (ordinal+1 byte) |
| 4 | `UdpNotificationType` | Wrong enum values: Go has `HeartbeatNotification`/`GoOfflineNotification`; Java has only `OPEN_CONNECTION` |
| 5 | `ParseUdpRequestType` | No validation of invalid input — returns an invalid enum instead of `nil`/error |

### UdpRequestType.java
Checked methods: parse(int number), getNumber()
## Bugs Found

### Bug 1: `parse(int number)` — Missing bounds checking and null/invalid handling

**Java:**
```java
public static UdpRequestType parse(int number) {
    int index = number - 1;
    if (index > -1 && index < ALL.length) {
        return ALL[index];
    }
    return null;
}
```

The Java version:
1. Subtracts 1 from `number` to get the ordinal index (numbers are 1-based: HEARTBEAT=1, GO_OFFLINE=2)
2. Validates the index is in bounds (`> -1 && < length`)
3. Returns `null` for any invalid number (out-of-range)

**Go:**
```go
func ParseUdpRequestType(number int) UdpRequestType {
    return UdpRequestType(number)
}
```

The Go version:
1. Does **not** subtract 1 — it treats the number as a 0-based value directly
2. Does **not** validate bounds — any integer is blindly cast to `UdpRequestType`
3. Cannot return a "not found" sentinel (the Java returns `null`)

**Impact:** Calling `ParseUdpRequestType(1)` returns `GoOfflineRequest` (value 1) instead of `HeartbeatRequest` (which should map from number 1). `ParseUdpRequestType(0)` returns `HeartbeatRequest`, but Java's `parse(0)` would return `null` because `index = 0 - 1 = -1` fails the bounds check.

---

### Bug 2: `getNumber()` — Off-by-one; returns 0-based value instead of 1-based

**Java:**
```java
public int getNumber() {
    return this.ordinal() + 1;
}
```
Returns `ordinal + 1`, so HEARTBEAT=1, GO_OFFLINE=2.

**Go:**
```go
func (t UdpRequestType) GetNumber() int {
    return int(t)
}
```
Returns the raw `int(t)` value, which (due to `iota` starting at 0) gives HeartbeatRequest=0, GoOfflineRequest=1.

**Impact:** `HeartbeatRequest.GetNumber()` returns `0` instead of `1`. `GoOfflineRequest.GetNumber()` returns `1` instead of `2`. All number values are off by 1.

---

### Bug 3: Coupled inconsistency between `parse` and `getNumber`

Because both methods are off by 1 (but in opposite directions), the round-trip invariant is also broken:

- **Java:** `parse(x.getNumber()) == x` holds (e.g., `HEARTBEAT.getNumber()=1`, `parse(1)=HEARTBEAT`)
- **Go:** `ParseUdpRequestType(x.GetNumber()) == x` does hold trivially only because both use the same wrong basis, but the values exchanged over the wire (which the Java code defines as 1-based) will be wrong — any protocol message carrying these numbers will be incompatible with the Java server/client.

### HttpForwardedHeaderHandler.java
Checked methods: apply(ConnectionInfo connectionInfo, HttpRequest request), getAddress(), send(ByteBuf buffer), close(CloseReason closeReason), close()
Now I have a complete picture. The user asked me to compare the Java `HttpForwardedHeaderHandler.apply()` and `NetConnection` methods (`getAddress`, `send`, `close(CloseReason)`, `close()`) against the Go refactor code in `ws_server.go`. The `ws_server.go` file is clearly a stub, but the actual implementations are in `net_connection.go` (base class) and `tcp_server.go` (TCP implementation). There is no WebSocket-specific connection implementation in Go — only a stub factory.

Here is my analysis:

## Bugs Found

### 1. `HttpForwardedHeaderHandler.Apply()` — Complete Missing Logic
**File:** `ws_server.go:11-13`  
**Java ref:** `HttpForwardedHeaderHandler.java:63-70`

The Go `Apply` method is an empty stub returning `nil`. The Java version:
1. Reads the `Forwarded` header; if present, delegates to `parseForwardedInfo()` which parses `for`, `proto`, and `host` directives using regex.
2. Otherwise delegates to `parseXForwardedInfo()` which parses `X-Forwarded-For`, `X-Forwarded-Proto`, `X-Forwarded-Host`, and `X-Forwarded-Port` headers.
3. Validates `isForwardedIpRequired` — throws `IllegalArgumentException` when the required forwarded IP is missing.

**None of this logic is present.** The struct also lacks the `isForwardedIpRequired` field.

---

### 2. `WebSocketConnection` — Entirely Missing
**Java ref:** `WebSocketConnection.java`

There is **no `WebSocketConnection` struct** in the Go codebase. The `ws_server.go` file only contains stubs for `HttpForwardedHeaderHandler` and `WebSocketServerFactory`. The Go code has a TCP implementation (`TcpConnection` in `tcp_server.go`) but no WebSocket counterpart implementing `GetAddress()`, `Send()`, `CloseWithReason()`, or `Close()`.

The Java `WebSocketConnection` has distinct behavior from TCP:
- `send()` wraps data in a `BinaryWebSocketFrame` (Go TCP just does raw `conn.Write`)
- `close(CloseReason)` sends a `BinaryWebSocketFrame` with notification data before closing
- `close()` calls `out.sendClose()` with a WebSocket close status code (Go TCP calls `conn.Close()`)

---

### 3. `BaseNetConnection.CloseWithReason()` — Missing Connection Disposal
**File:** `net_connection.go:50-57`  
**Java ref:** `NetConnection.java:72-77` (base) + `TcpConnection.java:70-98` / `WebSocketConnection.java:76-106` (concrete)

The Java `close(CloseReason)` is split into two levels:
- **Base `NetConnection.close(CloseReason)`:** Sets `isConnected=false`, `isConnectionRecovering=false`, `isSwitchingToUdp` based on status.
- **Concrete `TcpConnection.close(CloseReason)`:** Checks `!isConnected() || connection.isDisposed()` guard, calls `super.close(closeReason)`, sends a close notification buffer with retry logic (2 retries, 3s backoff), and then disposes the connection (with optional timeout).

The Go `BaseNetConnection.CloseWithReason()` only sets the three boolean flags. The `TcpConnection.CloseWithReason()` checks `IsConnected()`, calls the base method, then does `conn.Close()` — but is missing:
- The `connection.isDisposed()` check
- Sending the close notification (`NotificationFactory.createBuffer(closeReason)`) before closing
- The retry logic (2 retries with 3-second backoff)
- The `closeTimeout` handling (zero vs negative vs positive timeout paths)

---

### 4. `BaseNetConnection.Close()` / `TcpConnection.Close()` — Missing Connection Disposal Semantics
**File:** `net_connection.go:60-67`, `tcp_server.go:48-54`  
**Java ref:** `TcpConnection.java:101-112`, `WebSocketConnection.java:109-119`

**Java `TcpConnection.close()`:** Calls `connection.dispose()` with error handling for disconnected-client errors.  
**Java `WebSocketConnection.close()`:** Sends a WebSocket close frame (`out.sendClose(WebSocketCloseStatus.NORMAL_CLOSURE.code(), null)`) with error handling.

**Go `BaseNetConnection.Close()`:** Only sets three booleans to false. Does NOT close any underlying connection.  
**Go `TcpConnection.Close()`:** Calls base `Close()` then `conn.Close()`. This is closer but still differs:
- Missing the `isDisposed()` guard
- Missing error classification (checking for disconnected-client errors before logging)

---

### 5. `TcpConnection.Send()` — Missing Varint Length Prefix
**File:** `tcp_server.go:32-35`  
**Java ref:** `TcpConnection.java:61-64`

Java `send()` calls `connection.sendObject(buffer)` where `buffer` is a `ByteBuf` that the pipeline frames. The Go version does a raw `conn.Write(buffer)` without any framing protocol (no Varint length prefix), which means the receiver cannot delineate message boundaries.

---

### 6. `BaseNetConnection.CloseWithReason()` / `Close()` — Missing Volatile Semantics
**File:** `net_connection.go:35-41`  
**Java ref:** `NetConnection.java:44-46`

Java uses `volatile` for `isConnected`, `isSwitchingToUdp`, and `isConnectionRecovering`, providing visibility guarantees across threads without full locking. Go uses `sync.RWMutex` which is a valid alternative — **this is not a bug** per se, but worth noting as a design difference. The Go approach is actually more conservative and correct for general use.

---

### 7. `WebSocketServerFactory.Create()` — Empty Stub
**File:** `ws_server.go:19-20`  
**Java ref:** `WebSocketServerFactory` class

The factory method that should create a WebSocket server (with WebSocket properties, blocklist service, session service, connection listener, max frame payload length) is an empty stub with no server setup, no route handlers, and no connection lifecycle management.

### WebSocketServerFactory.java
Checked methods: create(WebSocketProperties webSocketProperties, BlocklistService blocklistService, ServerStatusManager serverStatusManager, SessionService sessionService, ConnectionListener connectionListener, int maxFramePayloadLength), sendNotificationToLocalClients(TracingContext tracingContext, ByteBuf notificationData, Set<Long> recipientIds, Set<UserSessionId> excludedUserSessionIds, @Nullable DeviceType excludedDeviceType), countLocalOnlineUsers(), handleServiceRequest(UserSession session, ServiceRequest serviceRequest), deleteSessions(@QueryParam(required = false), handleDeleteSessionRequest(UserSessionWrapper sessionWrapper), handleCreateSessionRequest(UserSessionWrapper sessionWrapper, CreateSessionRequest createSessionRequest)
The Go file is a single stub file with no real implementation. Here is the review:

## Bugs Found

### 1. `Create` method is a complete empty stub
The Java `create()` method contains extensive server setup logic that is entirely missing from the Go `Create()` method:
- **Missing:** `ServiceAvailabilityHandler` initialization with `blocklistService`, `serverStatusManager`, `sessionService`
- **Missing:** `WebsocketServerSpec` construction with `maxFramePayloadLength`
- **Missing:** `HttpServer.create()` with host/port binding
- **Missing:** Socket options: `CONNECT_TIMEOUT_MILLIS`, `SO_REUSEADDR`, `SO_BACKLOG`, `SO_LINGER`, `TCP_NODELAY`
- **Missing:** Proxy protocol support configuration (`REQUIRED`→`ON`, `OPTIONAL`→`AUTO`, `DISABLED`→`OFF`)
- **Missing:** Event loop configuration via `LoopResourcesFactory.createForServer`
- **Missing:** Metrics recorder setup (`TurmsMicrometerChannelMetricsRecorder`)
- **Missing:** HTTP request handler registration (`handleHttpRequest`)
- **Missing:** Channel pipeline initialization (`doOnChannelInit` adding `serviceAvailabilityHandler`)
- **Missing:** Forwarded header handler based on `RemoteAddressSourceHttpHeaderMode` (`REQUIRED`/`OPTIONAL`)
- **Missing:** SSL configuration when `ssl.isEnabled()`
- **Missing:** Server bind with error handling (`BindException` on failure)

### 2. `handleHttpRequest` is not implemented
The entire HTTP request handling logic is missing, which includes:
- **Missing:** CORS preflight request handling (returning OK with `Access-Control-Allow-Origin: *`, `Allow-Methods: *`, `Allow-Headers: *`, `Max-Age: 7200`, then `Mono.never()`)
- **Missing:** Handshake request validation (`validateHandshakeRequest`)
- **Missing:** IP blocklist check via `blocklistService.isIpBlocked()`
- **Missing:** WebSocket upgrade with `sendWebsocket`
- **Missing:** Inbound frame processing (aggregating frames, filtering for `BinaryWebSocketFrame` only)
- **Missing:** Close status handling (`receiveCloseStatus`)
- **Missing:** Remote address resolution fallback (`remoteAddress` → `connection.channel().remoteAddress()`)
- **Missing:** Connection listener invocation (`connectionListener.onAdded`)

### 3. `validateHandshakeRequest` is not implemented
Missing validation checks:
- **Missing:** HTTP method must be `GET`
- **Missing:** `Upgrade: websocket` header check
- **Missing:** `Connection: upgrade` header check
- **Missing:** `Sec-WebSocket-Key` header presence check

### 4. Methods listed in the review request don't exist in either file
The following methods were requested for review but do **not** exist in the Java `WebSocketServerFactory.java`:
- `sendNotificationToLocalClients`
- `countLocalOnlineUsers`
- `handleServiceRequest`
- `deleteSessions`
- `handleDeleteSessionRequest`
- `handleCreateSessionRequest`

These methods belong to other Java classes (likely in the `session` or `common` packages), not `WebSocketServerFactory`. No comparison is possible since the Go file does not reference them either.

### UserPermissionInfo.java
Checked methods: UserPermissionInfo(...), setCloseIdleSessionAfterSeconds(int closeIdleSessionAfterSeconds), setClientHeartbeatIntervalSeconds(int clientHeartbeatIntervalSeconds), destroy(), estimatedSize(), next(), addSessionIfAbsent(int version, Set<TurmsRequest.KindCase> permissions, DeviceType loggingInDeviceType, Map<String, String> deviceDetails, @Nullable Location location), closeSession(@NotNull DeviceType deviceType, @NotNull CloseReason closeReason), pushSessionNotification(DeviceType deviceType, String serverId), getSession(@NotNull DeviceType deviceType), countSessions(), getLoggedInDeviceTypes(), findPassword(Long userId), isActiveAndNotDeleted(Long userId), verifyAndGrant(UserLoginInfo userLoginInfo), verifyAndGrant(UserLoginInfo userLoginInfo), verifyAndGrant(UserLoginInfo userLoginInfo), verifyAndGrant(UserLoginInfo userLoginInfo), verifyAndGrant(UserLoginInfo userLoginInfo), updateGlobalProperties(TurmsProperties properties)
The Java file is a simple `record` with only 2 fields and a handful of static constants plus one convenience constructor. It does **not** contain any of the methods listed in your review request (`setCloseIdleSessionAfterSeconds`, `destroy`, `estimatedSize`, `addSessionIfAbsent`, `closeSession`, `pushSessionNotification`, `getSession`, `countSessions`, `getLoggedInDeviceTypes`, `findPassword`, `isActiveAndNotDeleted`, `verifyAndGrant`, `updateGlobalProperties`, etc.). Those methods belong to a completely different class — likely `UserSession` or `UserSessionManager` in the Java codebase.

The Java `UserPermissionInfo` is strictly:
1. A `record` with `ResponseStatusCode authenticationCode` and `Set<TurmsRequest.KindCase> permissions`
2. Static constants (`GRANTED_WITH_ALL_PERMISSIONS`, `LOGIN_AUTHENTICATION_FAILED`, etc.)
3. A convenience constructor `UserPermissionInfo(ResponseStatusCode authenticationCode)` that passes `Collections.emptySet()` for permissions

The Go code correctly mirrors this structure with:
1. A `struct` with `AuthenticationCode` and `Permissions` fields
2. A `NewUserPermissionInfo` constructor that maps to the two-arg Java constructor

**Bugs found:**

1. **Missing convenience constructor**: The Go code does not implement the single-argument constructor `UserPermissionInfo(ResponseStatusCode authenticationCode)` which defaults permissions to `Collections.emptySet()`. This constructor is used by the static constants `LOGIN_AUTHENTICATION_FAILED` and `LOGGING_IN_USER_NOT_ACTIVE`.

2. **Missing static constants**: The Go code does not define the following static constants that exist in the Java version:
   - `GRANTED_WITH_ALL_PERMISSIONS`
   - `GRANTED_WITH_ALL_PERMISSIONS_MONO`
   - `LOGIN_AUTHENTICATION_FAILED`
   - `LOGIN_AUTHENTICATION_FAILED_MONO`
   - `LOGGING_IN_USER_NOT_ACTIVE_MONO`

3. **Permissions type mismatch**: The Go code uses `map[interface{}]struct{}` for the permissions set, while Java uses `Set<TurmsRequest.KindCase>`. The use of `interface{}` loses type safety compared to the Java generic `Set<TurmsRequest.KindCase>`.

**Note**: The remaining methods listed in the review request (`setCloseIdleSessionAfterSeconds`, `destroy`, `estimatedSize`, `addSessionIfAbsent`, `closeSession`, `pushSessionNotification`, `getSession`, `countSessions`, `getLoggedInDeviceTypes`, `findPassword`, `isActiveAndNotDeleted`, `verifyAndGrant`, `updateGlobalProperties`) do **not** exist in the Java `UserPermissionInfo.java` file. They likely belong to a different class (e.g., `UserSession` or a session manager) and would need to be reviewed against their respective source files.

### SessionIdentityAccessManager.java
Checked methods: verifyAndGrant(int version, Long userId, @Nullable String password, DeviceType deviceType, @Nullable Map<String, String> deviceDetails, @Nullable UserStatus userStatus, @Nullable Location location, @Nullable String ip), destroy(), handleHeartbeatUpdateRequest(UserSession session), handleLoginRequest(int version, @NotNull ByteArrayWrapper ip, @NotNull Long userId, @Nullable String password, @NotNull DeviceType deviceType, @Nullable Map<String, String> deviceDetails, @Nullable UserStatus userStatus, @Nullable Location location, @Nullable String ipStr), closeLocalSessions(@NotNull List<byte[]> ips, @NotNull CloseReason closeReason), closeLocalSessions(@NotNull byte[] ip, @NotNull CloseReason closeReason), closeLocalSession(@NotNull Long userId, @NotNull @ValidDeviceType DeviceType deviceType, @NotNull SessionCloseStatus closeStatus), closeLocalSession(@NotNull Long userId, @NotNull @ValidDeviceType DeviceType deviceType, @NotNull CloseReason closeReason), closeLocalSession(@NotNull Long userId, @NotEmpty Set<@ValidDeviceType DeviceType> deviceTypes, @NotNull CloseReason closeReason), closeLocalSessions(@NotNull Set<Long> userIds, @NotNull CloseReason closeReason), authAndCloseLocalSession(@NotNull Long userId, @NotNull DeviceType deviceType, @NotNull CloseReason closeReason, int sessionId), closeAllLocalSessions(@NotNull CloseReason closeReason), closeLocalSession(Long userId, SessionCloseStatus closeStatus), closeLocalSession(Long userId, CloseReason closeReason), getSessions(Set<Long> userIds), authAndUpdateHeartbeatTimestamp(long userId, @NotNull @ValidDeviceType DeviceType deviceType, int sessionId), tryRegisterOnlineUser(int version, @NotNull Set<TurmsRequest.KindCase> permissions, @NotNull ByteArrayWrapper ip, @NotNull Long userId, @NotNull DeviceType deviceType, @Nullable Map<String, String> deviceDetails, @Nullable UserStatus userStatus, @Nullable Location location), getUserSessionsManager(@NotNull Long userId), getLocalUserSession(@NotNull Long userId, @NotNull DeviceType deviceType), getLocalUserSession(ByteArrayWrapper ip), countLocalOnlineUsers(), onSessionEstablished(@NotNull UserSessionsManager userSessionsManager, @NotNull @ValidDeviceType DeviceType deviceType), addOnSessionClosedListeners(Consumer<UserSession> onSessionClosed), invokeGoOnlineHandlers(@NotNull UserSessionsManager userSessionsManager, @NotNull UserSession userSession)
Now I have both files. Let me compare each ported method systematically.

## Bugs Found

### 1. `Destroy()` — Missing core logic
**Java** (`SessionService.java:214-223`): Calls `heartbeatManager.destroy()`, then `closeAllLocalSessions(closeReason)` with error mapping.
**Go** (`service.go:129-131`): Returns `nil` unconditionally. No heartbeat manager destruction, no session closing, no error handling.

### 2. `HandleHeartbeatUpdateRequest()` — Missing core logic
**Java** (`SessionService.java:225-227`): Calls `session.setLastHeartbeatRequestTimestampToNow()`.
**Go** (`service.go:134-135`): Empty method body. Does not update the heartbeat timestamp.

### 3. `HandleLoginRequest()` — Missing core logic
**Java** (`SessionService.java:229-270`):
1. Validates `version != 1` → returns error
2. Checks `userSimultaneousLoginService.isForbiddenDeviceType(deviceType)` → returns error
3. Calls `sessionAuthenticationManager.verifyAndGrant(...)` then `flatMap`s to either `tryRegisterOnlineUser` on success or error on non-OK status

**Go** (`service.go:138-140`): Returns `nil, nil` unconditionally. No version check, no forbidden device type check, no authentication, no user registration.

### 4. `CloseLocalSessions()` (by `List<byte[]> ips`) — Missing core logic
**Java** (`SessionService.java:275-296`): Validates params, iterates IPs, calls `closeLocalSessions(ip, closeReason)` for each, aggregates counts with `AtomicInteger`, uses `Mono.whenDelayError`.
**Go** (`service.go:144-146`): Returns `nil` unconditionally. No iteration, no closing, no error aggregation.

### 5. `CloseLocalSessions()` (by `byte[] ip`) — Missing core logic
**Java** (`SessionService.java:301-333`): Looks up sessions from `ipToSessions`, iterates and closes each session using `closeLocalSession(userId, ALL_AVAILABLE_DEVICE_TYPES_SET, closeReason)`, aggregates counts.
**Go** (`service.go:144-146`): Same stub — returns `nil`. No IP-based session lookup, no closing logic.

### 6. `CloseLocalSession()` (by userId + deviceType + SessionCloseStatus) — Missing core logic
**Java** (`SessionService.java:338-350`): Validates deviceType, wraps in singleton set, delegates to `closeLocalSession(userId, Set<DeviceType>, CloseReason)`.
**Go** (`service.go:153-155`): Returns `nil` unconditionally. No delegation.

### 7. `CloseLocalSession()` (by userId + deviceType + CloseReason) — Missing core logic
**Java** (`SessionService.java:355-365`): Validates deviceType, wraps in singleton set, delegates.
**Go** (`service.go:153-155`): Returns `nil` unconditionally.

### 8. `CloseLocalSession()` (by userId + Set<DeviceType> + CloseReason) — Missing core logic
**Java** (`SessionService.java:371-385`): Validates userId, looks up `UserSessionsManager`, delegates to private `closeLocalSessions(userId, deviceTypes, closeReason, manager)`.
**Go** (`service.go:153-155`): Returns `nil` unconditionally.

### 9. `CloseLocalSessions()` (by `Set<Long> userIds`) — Missing core logic
**Java** (`SessionService.java:390-411`): Validates params, iterates user IDs, calls `closeLocalSession(userId, closeReason)` for each (which uses `ALL_AVAILABLE_DEVICE_TYPES_SET`), aggregates counts with `whenDelayError`.
**Go** (`service.go:158-160`): Returns `nil` unconditionally.

### 10. `AuthAndCloseLocalSession()` — Missing core logic
**Java** (`SessionService.java:416-440`): Validates params, looks up manager, checks session exists and `session.getId() == sessionId`, then closes matching sessions.
**Go** (`service.go:163-165`): Returns `nil` unconditionally. No session ID verification, no closing.

### 11. `CloseAllLocalSessions()` — Missing core logic
**Java** (`SessionService.java:445-464`): Validates, iterates all entries in `userIdToSessionsManager`, gets `loggedInDeviceTypes` for each, calls `closeLocalSession` for each, aggregates counts.
**Go** (`service.go:168-170`): Returns `nil` unconditionally.

### 12. `CloseLocalSession()` (by userId + SessionCloseStatus) — Missing core logic
**Java** (`SessionService.java:469-473`): Delegates to `closeLocalSession(userId, ALL_AVAILABLE_DEVICE_TYPES_SET, CloseReason.get(closeStatus))`.
**Go** (`service.go:153-155`): Returns `nil` unconditionally.

### 13. `CloseLocalSession()` (by userId + CloseReason) — Missing core logic
**Java** (`SessionService.java:479-483`): Delegates to `closeLocalSession(userId, ALL_AVAILABLE_DEVICE_TYPES_SET, closeReason)`.
**Go** (`service.go:153-155`): Returns `nil` unconditionally.

### 14. `GetSessions()` — Missing core logic
**Java** (`SessionService.java:546-552`): Iterates userIds, builds `UserSessionsInfo` for each (with session details: ID, version, deviceType, deviceDetails, loginDate, loginLocation, heartbeat/request timestamps, IP, etc.).
**Go** (`service.go:173-175`): Returns `nil`. No session info gathering.

### 15. `AuthAndUpdateHeartbeatTimestamp()` — Missing core logic
**Java** (`SessionService.java:587-605`):
1. Validates deviceType
2. Looks up `UserSessionsManager` by userId
3. Gets session by deviceType
4. Checks `session.getId() == sessionId` AND `!session.getConnection().isConnectionRecovering()`
5. If all pass, calls `session.setLastHeartbeatRequestTimestampToNow()` and returns session

**Go** (`service.go:178-180`): Returns `nil` unconditionally. No session lookup, no ID/auth check, no heartbeat update.

### 16. `TryRegisterOnlineUser()` — Missing core logic
**Java** (`SessionService.java:612-767`): This is a very complex method (~155 lines) that:
1. Validates IP, deviceType, userStatus (not UNRECOGNIZED/OFFLINE), location bounds
2. Fetches user sessions status from Redis via `userStatusService.fetchUserSessionsStatus`
3. Closes local sessions if registered on other nodes
4. Handles offline user case → `addOnlineDeviceIfAbsent`
5. Handles online user with device conflict: closed local session recovery path, simultaneous login conflict declination
6. Handles conflicted device types → `closeSessionsWithConflictedDeviceTypes` then `addOnlineDeviceIfAbsent`

**Go** (`service.go:183-185`): Returns `nil, nil` unconditionally. None of the logic is implemented.

### 17. `GetUserSessionsManager()` — Missing core logic
**Java** (`SessionService.java:770-773`): Validates userId, returns `userIdToSessionsManager.get(userId)`.
**Go** (`service.go:188-190`): Returns `nil` unconditionally. No map lookup.

### 18. `GetLocalUserSession()` (by userId + deviceType) — Missing core logic
**Java** (`SessionService.java:776-783`): Validates params, looks up manager, returns `manager.getSession(deviceType)`.
**Go** (`service.go:194-196`): Returns `nil` unconditionally.

### 19. `GetLocalUserSession()` (by `ByteArrayWrapper ip`) — Missing core logic
**Java** (`SessionService.java:786-788`): Returns `ipToSessions.get(ip)` — a queue of sessions.
**Go** (`service.go:194-196`): Returns `nil`. No IP-based lookup. Also, the Go method signature doesn't support IP-based lookup (only takes userId + deviceType).

### 20. `CountLocalOnlineUsers()` — Correctly implemented
**Java** (`SessionService.java:790-792`): Returns `userIdToSessionsManager.size()`.
**Go** (`service.go:124-126`): Delegates to `s.shardedMap.CountOnlineUsers()`. This appears functionally correct.

### 21. `OnSessionEstablished()` — Missing core logic
**Java** (`SessionService.java:858-865`):
1. Increments `loggedInUsersCounter` (metrics)
2. If `notifyClientsOfSessionInfoAfterConnected` is true, calls `userSessionsManager.pushSessionNotification(deviceType, serverId)`

**Go** (`service.go:199-200`): Empty method body. No metrics counter increment, no session notification push.

### 22. `AddOnSessionClosedListeners()` — Missing core logic
**Java** (`SessionService.java:974-976`): Adds the listener to `onSessionClosedListeners` list.
**Go** (`service.go:203-204`): Empty method body. The listener is not stored anywhere.

### 23. `InvokeGoOnlineHandlers()` — Missing core logic
**Java** (`SessionService.java:992-999`): Invokes plugin extension points for `UserOnlineStatusChangeHandler.goOnline`.
**Go** (`service.go:207-208`): Empty method body. No plugin invocation.

### 24. `VerifyAndGrant()` (SessionIdentityAccessManager) — Missing core logic
**Java** (`SessionIdentityAccessManager.java:107-149`):
1. Checks `userId.equals(AdminConst.ADMIN_REQUESTER_ID)` → returns `LOGIN_AUTHENTICATION_FAILED`
2. If `!enableIdentityAccessManagement` → returns `GRANTED_WITH_ALL_PERMISSIONS`
3. Creates `UserLoginInfo`, calls `sessionIdentityAccessManagementSupport.verifyAndGrant(userLoginInfo)`
4. If plugins with `UserAuthenticator` are running, invokes them sequentially with `switchIfEmpty` fallback to default handler

**Go** (`identity_access_manager.go`): Only the interface and stub implementations exist. The `verifyAndGrant` orchestration logic (admin check, enablement check, plugin invocation chain) is completely absent from all implementations. Each concrete type's `VerifyAndGrant` either returns `nil, nil` (Http, Jwt, Ldap) or a hardcoded value (Noop, Password). The `PasswordSessionIdentityAccessManager` returns `SERVER_INTERNAL_ERROR` instead of performing password verification.

### 25. Missing `ipToSessions` map in Go `SessionService`
**Java**: Maintains a `ConcurrentHashMap<ByteArrayWrapper, ConcurrentLinkedQueue<UserSession>> ipToSessions` used for IP-based session lookup and cleanup.
**Go**: The `SessionService` struct has no `ipToSessions` equivalent. This means IP-based session lookups (`getLocalUserSession(ByteArrayWrapper ip)`) and IP-based session cleanup in `closeLocalSessions(byte[] ip)` cannot function.

### 26. Missing `onSessionClosedListeners` in Go `SessionService`
**Java**: Maintains a `List<Consumer<UserSession>> onSessionClosedListeners` with a `notifyOnSessionClosedListeners` method that is called when sessions are closed.
**Go**: No equivalent field or notification mechanism exists.

### 27. Missing private `closeLocalSessions(userId, deviceTypes, closeReason, manager)` method
**Java** (`SessionService.java:491-543`): This is the core session closing logic that:
1. Calls `userStatusService.removeStatusByUserIdAndDeviceTypes` to remove Redis status
2. Iterates device types, closes sessions via `manager.closeSession`
3. Removes session from `ipToSessions`
4. Removes user location if location service is enabled
5. Notifies session closed listeners
6. Calls `removeSessionsManagerIfEmpty` (which also invokes `goOffline` plugin handlers)

**Go**: This entire private method and all its logic are absent.

### 28. Missing `removeSessionsManagerIfEmpty` logic
**Java** (`SessionService.java:958-970`): Removes the manager from the map if no sessions remain, then invokes `UserOnlineStatusChangeHandler.goOffline` plugin extension points.
**Go**: No equivalent logic exists. `ShardedUserSessionsMap.RemoveIfEmpty` is called in `UnregisterSession` but does not trigger plugin handlers.

### 29. Missing `addOnlineDeviceIfAbsent` method
**Java** (`SessionService.java:867-956`): Critical private method (~90 lines) that:
1. Filters device details per configuration
2. Calls `userStatusService.addOnlineDeviceIfAbsent` to register in Redis
3. Creates/retrieves `UserSessionsManager` and adds session
4. Handles the case where session already exists (close and retry)
5. Registers session in `ipToSessions`
6. Upserts user location if location service is enabled

**Go**: No equivalent method exists.

### 30. Missing `closeSessionsWithConflictedDeviceTypes` method
**Java** (`SessionService.java:794-856`): Handles multi-device conflict resolution by sending RPC requests to other nodes to disconnect conflicting sessions, with handling for dead nodes.
**Go**: No equivalent method exists.

### 31. `GetLocalUserSession` — Signature doesn't support IP-based overload
**Java** has two overloads: `getLocalUserSession(Long userId, DeviceType deviceType)` and `getLocalUserSession(ByteArrayWrapper ip)`.
**Go** (`service.go:194-196`): Only one method exists with `(userId, deviceType)` signature. The IP-based overload is listed in the `@MappedFrom` comment but is not implemented as a separate method.

### UserService.java
Checked methods: authenticate(@NotNull Long userId, @Nullable String rawPassword), isActiveAndNotDeleted(@NotNull Long userId)
## Bugs in Go refactor of `UserService`

### 1. Missing struct fields / dependencies

**Java:**
```java
private final UserRepository userRepository;
private final PasswordManager passwordManager;
@Getter
private final boolean enabled;

public UserService(UserRepository userRepository, PasswordManager passwordManager) {
    this.userRepository = userRepository;
    this.passwordManager = passwordManager;
    enabled = userRepository.isEnabled();
}
```

**Go:**
```go
type UserService struct {
}
```

The Go `UserService` struct is **empty** — it has no `userRepository`, `passwordManager`, or `enabled` fields. None of the constructor initialization logic is present.

### 2. `Authenticate` is a stub with no real logic

**Java `authenticate`:**
1. Validates `userId` is not null
2. Calls `userRepository.findPassword(userId)` to look up the user's stored password
3. Calls `passwordManager.matchesUserPassword(rawPassword, user.getPassword())` to compare passwords
4. Returns `false` if user doesn't exist (via `.defaultIfEmpty(false)`)

**Go `Authenticate`:**
```go
func (s *UserService) Authenticate(ctx context.Context, userID int64, rawPassword string) (bool, error) {
    // Stub implementation
    return false, nil
}
```

Missing:
- No call to a repository to find the user's stored password
- No password comparison via a `PasswordManager` equivalent
- Always returns `false, nil` regardless of input

### 3. Method `IsActiveAndNotDeleted` is completely missing

**Java `isActiveAndNotDeleted`:**
```java
public Mono<Boolean> isActiveAndNotDeleted(@NotNull Long userId) {
    Validator.notNull(userId, "userId");
    return userRepository.isActiveAndNotDeleted(userId);
}
```

The Go file has **no** `IsActiveAndNotDeleted` method at all.

### Summary

| Bug | Severity |
|-----|----------|
| `UserService` struct missing `userRepository`, `passwordManager`, `enabled` fields | Critical |
| No constructor/initialization logic | Critical |
| `Authenticate` is a no-op stub — no password lookup or comparison | Critical |
| `IsActiveAndNotDeleted` method entirely absent | Critical |

### UserSimultaneousLoginService.java
Checked methods: getConflictedDeviceTypes(@NotNull @ValidDeviceType DeviceType deviceType), isForbiddenDeviceType(DeviceType deviceType), shouldDisconnectLoggingInDeviceIfConflicts()
## Bugs Found

### 1. `GetConflictedDeviceTypes` — Missing all core logic

**Java behavior:** Looks up `deviceType` in the `deviceTypeToExclusiveDeviceTypes` map (populated in `updateProperties` from the configured `SimultaneousLoginStrategy`) and returns the set of conflicted device types. Returns `null` (Go: `nil`) if the key is absent.

**Go code:** Returns hardcoded `nil` with no fields and no logic. The entire `deviceTypeToExclusiveDeviceTypes` map field, its initialization via `newExclusiveDeviceFromStrategy()`, and the lookup are missing.

---

### 2. `IsForbiddenDeviceType` — Missing all core logic

**Java behavior:** Checks `forbiddenDeviceTypes.contains(deviceType)`, where `forbiddenDeviceTypes` is populated in `updateProperties` from `newForbiddenDeviceTypesFromStrategy()`. The forbidden set depends on the `SimultaneousLoginStrategy`, `allowDeviceTypeUnknownLogin`, and `allowDeviceTypeOthersLogin`.

**Go code:** Returns hardcoded `false`. The `forbiddenDeviceTypes` field, the `allowDeviceTypeUnknownLogin` / `allowDeviceTypeOthersLogin` fields, and `newForbiddenDeviceTypesFromStrategy()` logic are all missing.

---

### 3. `ShouldDisconnectLoggingInDeviceIfConflicts` — Missing all core logic

**Java behavior:** Returns `loginConflictStrategy == LoginConflictStrategy.DISCONNECT_LOGGING_IN_DEVICE`.

**Go code:** Returns hardcoded `false`. The `loginConflictStrategy` field and the comparison are missing.

---

### 4. Missing struct fields and initialization

The Go struct `UserSimultaneousLoginService` is completely empty. The following fields present in Java are absent:
- `deviceTypeToExclusiveDeviceTypes` (`Map<DeviceType, Set<DeviceType>>`)
- `forbiddenDeviceTypes` (`Set<DeviceType>`)
- `allowDeviceTypeUnknownLogin` (`bool`)
- `allowDeviceTypeOthersLogin` (`bool`)
- `loginConflictStrategy` (`LoginConflictStrategy`)

The `updateProperties` method (and its helper methods `newExclusiveDeviceFromStrategy`, `newForbiddenDeviceTypesFromStrategy`, `addConflictedDeviceTypes`, `addDeviceTypeConflictedWithAllTypes`) is entirely unimplemented, meaning none of the three target methods can ever return correct results.

### ServiceAddressManager.java
Checked methods: getWsAddress(), getTcpAddress(), getUdpAddress()
## Bugs

### 1. Missing struct fields for `wsAddress`, `tcpAddress`, `udpAddress`

The Java class has three `@Nullable String` fields that store the resolved addresses:

```java
@Nullable private String wsAddress;
@Nullable private String tcpAddress;
@Nullable private String udpAddress;
```

The Go struct is completely empty:

```go
type ServiceAddressManager struct {
}
```

It should have fields like:
```go
wsAddress  string
tcpAddress string
udpAddress string
```

### 2. Getter methods are stubs — they always return empty string instead of the stored field values

`GetWsAddress()`, `GetTcpAddress()`, and `GetUdpAddress()` all return the hardcoded empty string `""` instead of returning their respective struct fields. The Java versions return `wsAddress`, `tcpAddress`, `udpAddress`.

### 3. Missing `gatewayApiDiscoveryProperties` field

The Java class stores a `DiscoveryProperties` reference used to detect config changes:

```java
private DiscoveryProperties gatewayApiDiscoveryProperties;
```

This field is absent from the Go struct. Without it, the `areAddressPropertiesChange` optimization cannot be implemented.

### 4. Missing `updateCustomAddresses` method — the entire logic that populates the three address fields

This is the core logic that actually sets `wsAddress`, `tcpAddress`, and `udpAddress`. The Java implementation:

- Checks `areAddressPropertiesChange` to short-circuit if nothing changed.
- For **WS**: when `webSocketProperties.isEnabled()`, it calls `queryHost(...)` and builds the address as `"ws://" + host + ":port"` (or `"wss://"` if `adminHttpProperties.getSsl().isEnabled()`).
- For **TCP**: when `tcpProperties.isEnabled()`, it calls `queryHost(...)` and builds `host + ":port"` (plain, no protocol prefix).
- For **UDP**: when `udpProperties.isEnabled()`, it calls `queryHost(...)` and builds `host + ":port"` (plain, no protocol prefix).
- Port is only appended when `attachPortToHost` is `true`.
- Uses `Mono.whenDelayError` to run all updates concurrently.

None of this logic exists in the Go file.

### 5. Missing `areAddressPropertiesChange` helper

The Java class has a private method that compares the current and new discovery properties (advertise strategy, advertise host, attach-port-to-host flag). This is absent.

### 6. Missing `getAdminAddressProperties` override

The Java class overrides `getAdminAddressProperties` to return `properties.getGateway().getAdminApi().getAddress()`. This is absent.

### 7. Missing constructor / initialization — no `BaseServiceAddressManager` embedding

The Java class extends `BaseServiceAddressManager` and passes `AdminHttpProperties`, `IpDetector`, and `TurmsPropertiesManager` to its constructor. The Go struct has no embedding of an equivalent base type and no constructor.

### 8. `GetWsAddress`/`GetTcpAddress`/`GetUdpAddress` return types should be nullable

The Java methods return `@Nullable String`, meaning they can return `null`. The Go methods return `string` (zero value `""`), which doesn't distinguish "no address" from "empty address". They should likely return `*string` or have a second `bool` return to indicate presence.

### LdapClient.java
Checked methods: isConnected(), connect(), bind(boolean useFastBind, String dn, String password), search(String baseDn, Scope scope, DerefAliases derefAliases, int sizeLimit, int timeLimit, boolean typeOnly, List<String> attributes, String filter), modify(String dn, List<ModifyOperationChange> changes)
## Bugs Found

### 1. **`isConnected()` — Method completely missing**
The Go code has no `isConnected()` method. The Java version checks `connection != null && !connection.isDisposed()`. There is no equivalent in Go.

### 2. **`LdapClient` struct — All fields missing**
The Go struct `LdapClient` is empty. The Java version has essential fields:
- `host`, `port`, `sslProperties` (connection config)
- `messageId` (atomic counter starting at 1 per RFC 4511)
- `messageIdToRequestContext` (pending request map)
- `connection` (current channel)
- `connectionMono` (reactive connection holder)

### 3. **`Connect()` — Entire implementation is a stub returning `nil`**
Missing all core logic from the Java version:
- Connection deduplication via `CONNECTION_MONO_UPDATER.compareAndSet`
- TCP client creation with host/port
- SSL/TLS configuration when `sslProperties` is enabled
- Adding LDAP message encoder/decoder handlers
- Error-handling subscription on `receiveObject()` that disposes the connection on error
- Assigning the `connection` field after successful connect

### 4. **`Connect()` — Wrong return type**
Returns `error` but the Java version returns `Mono<ChannelOperations>` (a reactive async result). The Go version provides no way to access the established connection.

### 5. **`Bind()` — Entire implementation is a stub returning `nil`**
Missing all core logic from the Java version:
- Sending a `BindRequest` with DN and password
- Passing `REQUEST_CONTROLS_FAST_BIND` when `useFastBind` is true
- Checking `response.isSuccess()` → return `true`
- Checking `resultCode == INVALID_CREDENTIALS` → return `false`
- Throwing `LdapException` for other error result codes
- The return type should distinguish between success (`true`), invalid credentials (`false`), and error, but Go returns only `error`

### 6. **`Bind()` — Wrong return type**
Returns `error` but the Java version returns `Mono<Boolean>`. The method cannot communicate the three-way result (authenticated / invalid credentials / error) that the Java version supports.

### 7. **`search()` — Method completely missing**
The `search(String baseDn, Scope scope, DerefAliases derefAliases, int sizeLimit, int timeLimit, boolean typeOnly, List<String> attributes, String filter)` method is entirely absent from the Go code.

### 8. **`Modify()` — Entire implementation is a stub returning `nil`**
Missing all core logic from the Java version:
- Early return for empty `changes` list (`Mono.empty()`)
- Validation loop: checking that `ADD` operations have at least one attribute value, throwing `LdapException(INVALID_ATTRIBUTE_SYNTAX, ...)` if not
- Sending a `ModifyRequest` and returning the response

### 9. **`Modify()` — Wrong parameter type for `changes`**
Uses `[]any` instead of a typed slice (e.g., `[]ModifyOperationChange`), losing all type safety for the changes parameter.

### BerBuffer.java
Checked methods: skipTag(), skipTagAndLength(), skipTagAndLengthAndValue(), readTag(), peekAndCheckTag(int tag), skipLength(), skipLengthAndValue(), writeLength(int length), readLength(), tryReadLengthIfReadable(), beginSequence(), beginSequence(int tag), endSequence(), writeBoolean(boolean value), writeBoolean(int tag, boolean value), readBoolean(), writeInteger(int value), writeInteger(int tag, int value), readInteger(), readIntWithTag(int tag), writeOctetString(String value), writeOctetString(byte[] value), writeOctetString(int tag, byte[] value), writeOctetString(byte[] value, int start, int length), writeOctetString(int tag, byte[] value, int start, int length), writeOctetString(int tag, String value), writeOctetStrings(List<String> values), readOctetString(), readOctetStringWithTag(int tag), readOctetStringWithLength(int length), writeEnumeration(int value), readEnumeration(), getBytes(), skipBytes(int length), close(), refCnt(), retain(), retain(int increment), touch(), touch(Object hint), release(), release(int decrement), isReadable(int length), isReadable(), isReadableWithEnd(int end), readerIndex()
## Bugs

The Go file `ber_buffer.go` contains **no actual implementation** — every method is an empty stub with zero/hardcoded return values. The `BerBuffer` struct has **no fields** at all. Comparing against the Java source:

### 1. Missing struct fields

- `sequenceLengthWriterIndexes []int` — used by `beginSequence`/`endSequence` for tracking nested sequence lengths
- `currentSequenceLengthIndex int` — tracks current nesting depth
- `buffer` (the underlying `ByteBuf` equivalent) — the core data buffer that every method reads from/writes to

Without these fields, none of the methods can function.

### 2. `skipTag()` — empty stub
Java advances the reader index by 1 byte. Go does nothing.

### 3. `skipTagAndLength()` — empty stub
Java skips 1 byte then calls `readLength()`. Go does nothing.

### 4. `skipTagAndLengthAndValue()` — empty stub
Java skips 1 byte, reads the length, then skips that many bytes. Go does nothing.

### 5. `readTag()` — returns hardcoded `0`
Java returns `buffer.readByte()`. Go returns `0`.

### 6. `peekAndCheckTag()` — wrong signature and empty stub
Java returns `bool`, Go returns nothing (`void`). Java checks `buffer.isReadable() && buffer.getByte(buffer.readerIndex()) == tag`. Go does nothing.

### 7. `skipLength()` — empty stub
Java calls `readLength()` and discards the result. Go does nothing.

### 8. `skipLengthAndValue()` — empty stub
Java calls `readLength()` then `buffer.skipBytes(length)`. Go does nothing.

### 9. `writeLength()` — empty stub
Java writes 1–5 bytes encoding the BER length. Go does nothing.

### 10. `readLength()` — returns hardcoded `0`
Java reads a BER-encoded length (1–5 bytes) with validation for indefinite length, overflow, and insufficient data. Go returns `0`.

### 11. `tryReadLengthIfReadable()` — returns hardcoded `0` instead of `-1`
Java returns `-1` when not readable. Go returns `0`, which is a valid length value — callers cannot distinguish "not readable" from "length is 0".

### 12. `beginSequence()` — empty stub
Java writes a tag, reserves 3 bytes for the length, records the writer index, and increments the nesting counter. Go does nothing.

### 13. `beginSequenceWithTag()` — empty stub
Same as above but with a custom tag. Go does nothing.

### 14. `endSequence()` — empty stub
Java decrements the nesting index, calculates value length, writes a 3-byte length at the reserved position, and restores the writer index. Go does nothing.

### 15. `writeBoolean()` — empty stub
Java writes tag + length(1) + value (0xFF or 0x00). Go does nothing.

### 16. `writeBooleanWithTag()` — empty stub
Same as above with custom tag. Go does nothing.

### 17. `readBoolean()` — returns hardcoded `false`
Java reads tag, validates it, reads length, validates it, reads value byte and returns `!= 0`. Go returns `false`.

### 18. `writeInteger()` — empty stub
Java calls `writeInteger(TAG_INTEGER, value)`. Go does nothing.

### 19. `writeIntegerWithTag()` — empty stub
Java writes tag + BER-encoded integer (1–4 bytes) with sign handling. Go does nothing.

### 20. `readInteger()` — returns hardcoded `0`
Java calls `readIntWithTag(TAG_INTEGER)`. Go returns `0`.

### 21. `readIntWithTag()` — returns hardcoded `0`
Java reads tag, validates, reads length, reads bytes with sign extension. Go returns `0`.

### 22. `writeOctetString(string)` — empty stub
Java delegates to `writeOctetString(TAG_OCTET_STRING, value)`. Go does nothing.

### 23. `writeOctetStringBytes()` — empty stub
Java writes tag + length + raw bytes. Go does nothing.

### 24. `writeOctetStringBytesWithTag()` — empty stub
Same as above with custom tag. Go does nothing.

### 25. `writeOctetStringBytesRange()` — empty stub
Java writes tag + length + byte slice subset. Go does nothing.

### 26. `writeOctetStringBytesRangeWithTag()` — empty stub
Same as above with custom tag. Go does nothing.

### 27. `writeOctetStringWithTag()` — empty stub
Java writes tag, reserves 3 bytes for length, writes UTF-8 encoded string, backfills the length. Go does nothing.

### 28. `writeOctetStrings()` — empty stub
Java iterates the list calling `writeOctetString` for each. Go does nothing.

### 29. `readOctetString()` — returns hardcoded `""`
Java delegates to `readOctetStringWithTag(TAG_OCTET_STRING)`. Go returns `""`.

### 30. `readOctetStringWithTag()` — returns hardcoded `""`
Java reads tag, validates, reads length, reads UTF-8 chars. Go returns `""`.

### 31. `readOctetStringWithLength()` — returns hardcoded `""`
Java reads `length` bytes as UTF-8. Go returns `""`.

### 32. `writeEnumeration()` — empty stub
Java delegates to `writeInteger(TAG_ENUMERATED, value)`. Go does nothing.

### 33. `readEnumeration()` — returns hardcoded `0`
Java delegates to `readIntWithTag(TAG_ENUMERATED)`. Go returns `0`.

### 34. `getBytes()` — returns `nil`
Java returns a copy of the entire buffer contents. Go returns `nil`.

### 35. `skipBytes()` — empty stub
Java calls `buffer.skipBytes(length)`. Go does nothing.

### 36. Missing `close()` method
Java implements `Closeable.close()` calling `buffer.release()`. Go has no equivalent.

### 37. `refCnt()` — returns hardcoded `0`
Java delegates to `buffer.refCnt()`. Go returns `0`.

### 38. `retain()` — empty stub
Java calls `buffer.retain()`. Go does nothing.

### 39. `retainIncrement()` — empty stub
Java calls `buffer.retain(increment)`. Go does nothing.

### 40. `touch()` — empty stub
Java calls `buffer.touch()`. Go does nothing.

### 41. `touchWithHint()` — empty stub
Java calls `buffer.touch(hint)`. Go does nothing.

### 42. `release()` — returns hardcoded `false`
Java calls `buffer.release()`. Go returns `false`.

### 43. `releaseDecrement()` — returns hardcoded `false`
Java calls `buffer.release(decrement)`. Go returns `false`.

### 44. `isReadableLen()` — returns hardcoded `false`
Java calls `buffer.isReadable(length)`. Go returns `false`.

### 45. `isReadable()` — returns hardcoded `false`
Java calls `buffer.isReadable()`. Go returns `false`.

### 46. `isReadableWithEnd()` — returns hardcoded `false`
Java checks `buffer.readerIndex() < end`. Go returns `false`.

### 47. `readerIndex()` — returns hardcoded `0`
Java calls `buffer.readerIndex()`. Go returns `0`.

### Summary

Every single method is an unimplemented stub. The struct has no fields. None of the Java logic has been ported. This file is a scaffold with correct method names/signatures but zero functional code.

### Attribute.java
Checked methods: isEmpty(), decode(BerBuffer buffer)
## Bugs in `Attribute` Go implementation

### 1. Missing struct fields
**Java** (`Attribute.java:36-37`): The class has two fields:
```java
private final String type;
private final List<String> values;
```
**Go** (`elements.go:7-8`): The struct has **no fields at all**:
```go
type Attribute struct {
}
```
The `type string` and `values []string` fields are completely missing.

### 2. `isEmpty()` not implemented
**Java** (`Attribute.java:39-41`): Returns `values.isEmpty()`.
**Go**: There is no `IsEmpty()` method on `Attribute`. The method is entirely absent.

### 3. `Decode()` is an empty stub
**Java** (`Attribute.java:43-61`): The `decode` method performs significant logic:
1. Skips tag and length (`buffer.skipTagAndLength()`)
2. Reads an octet string into `type`
3. Reads a tag and checks if it equals `TAG_SET | FORM_CONSTRUCTED`
4. If the tag doesn't match: skips length/value, returns `Attribute` with empty values list
5. If the tag matches: reads length, computes end position, loops reading octet strings into `values` until `isReadableWithEnd(end)`
6. Returns a new `Attribute` with the decoded `type` and `values`

**Go** (`elements.go:11-12`): The method body is completely empty:
```go
func (a *Attribute) Decode(buffer *asn1.BerBuffer) {
}
```
None of the above logic is present — no tag/length skipping, no octet string reading, no tag validation, no loop, no field assignments.

### 4. `Decode` signature mismatch — should be a constructor-like factory function
**Java**: `decode` is a `static` factory method that returns a **new** `Attribute` instance with the decoded `type` and `values`. The Java code is:
```java
public static Attribute decode(BerBuffer buffer) { ... return new Attribute(type, values); }
```
**Go**: `Decode` is a method on `*Attribute` with no return value. Even if implemented, this design cannot work correctly because there's no way to populate the struct fields (they don't exist) or return a new instance. The Go version should either be a function `func DecodeAttribute(buffer *asn1.BerBuffer) *Attribute` or the struct needs fields and the method needs to populate them.

### LdapResult.java
Checked methods: isSuccess()
## Bugs in `LdapResult` Go port

### Bug 1: `LdapResult` struct is empty — missing all fields

**Java** (`LdapResult.java:37-40`):
```java
private final int resultCode;
private final String matchedDn;
private final String diagnosticMessage;
private final List<String> referrals;
```

**Go** (`elements.go:30-31`):
```go
type LdapResult struct {
}
```

The struct has zero fields. `resultCode`, `matchedDn`, `diagnosticMessage`, and `referrals` are all missing.

### Bug 2: `IsSuccess()` always returns `false` instead of checking `resultCode`

**Java** (`LdapResult.java:80-82`):
```java
public boolean isSuccess() {
    return resultCode == ResultCodeConst.SUCCESS;
}
```

**Go** (`elements.go:34-36`):
```go
func (r *LdapResult) IsSuccess() bool {
    return false
}
```

The method is a hardcoded stub that always returns `false`. It never checks a `resultCode` field against a success constant. This means it will report every LDAP operation as failed, even successful ones.

### Bug 3: `decodeResult()` static factory method is entirely missing

**Java** (`LdapResult.java:53-78`): The `decodeResult(BerBuffer)` method parses raw BER-encoded data to construct an `LdapResult` — reading `resultCode` (as enumeration), `matchedDn`, `diagnosticMessage`, and conditionally parsing referrals.

**Go**: No equivalent `DecodeResult` function exists anywhere in the file.

### Control.java
Checked methods: decode(BerBuffer buffer)
## Bugs found in Go `Control.Decode` compared to Java `Control.decode`

### 1. **Stub implementation — entire method body is missing**

The Go `Control.Decode` (line 44-45) is an empty stub:
```go
func (c *Control) Decode(buffer *asn1.BerBuffer) {
}
```

The Java `decode` method contains all of the following logic that is entirely absent:

### 2. **Missing struct fields**

The Go `Control` struct (line 40-41) has no fields at all:
```go
type Control struct {
}
```

The Java class has two fields that are never assigned:
- `oid` (String) — the control OID
- `criticality` (boolean) — whether the control is critical

### 3. **Missing return type — returns nothing instead of `[]Control`**

The Java method is `static` and returns `List<Control>`. The Go method is a receiver method on `*Control` and returns nothing. It cannot serve the same purpose — callers in Java get a list of decoded controls; Go callers get nothing.

### 4. **Missing decode logic (enumerated)**

The following logic from Java is completely absent:

| Step | Java Logic | Go Status |
|------|-----------|-----------|
| 1 | Check `buffer.isReadable()` and peek for `LdapTagConst.CONTROLS` tag; return empty list if absent | Missing |
| 2 | Read outer tag (`buffer.readTag()`) and length (`buffer.readLength()`) | Missing |
| 3 | Calculate end position (`readerIndex() + length`) | Missing |
| 4 | Loop `while (buffer.isReadableWithEnd(end))` over each control element | Missing |
| 5 | `buffer.skipTagAndLength()` for each control sequence | Missing |
| 6 | `buffer.readOctetString()` to read the OID | Missing |
| 7 | Peek for `TAG_BOOLEAN` → read criticality, then optionally skip octet string value | Missing |
| 8 | Else peek for `TAG_OCTET_STRING` → skip tag+length+value | Missing |
| 9 | Add `new Control(oid, criticality)` to result list if `oid != null` | Missing |
| 10 | Return the list of controls | Missing |

### BindRequest.java
Checked methods: estimateSize(), writeTo(BerBuffer buffer)
## Bugs in Go `BindRequest`

### 1. Missing fields `dn` and `password`
**Java** (`BindRequest.java:34-35`):
```java
private final String dn;
private final String password;
```
**Go** (`elements.go:49-50`):
```go
type BindRequest struct {
}
```
The Go struct has no fields at all. It needs `Dn string` and `Password string`.

### 2. `EstimateSize()` returns a hardcoded `0` instead of the actual size computation
**Java** (`BindRequest.java:43-45`):
```java
public int estimateSize() {
    return dn.length() + password.length() + 16;
}
```
**Go** (`elements.go:53-55`):
```go
func (r *BindRequest) EstimateSize() int {
    return 0
}
```
Should be `return len(r.Dn) + len(r.Password) + 16`.

### 3. `WriteTo()` is an empty stub — all encoding logic is missing
**Java** (`BindRequest.java:48-54`):
```java
public void writeTo(BerBuffer buffer) {
    buffer.beginSequence(LdapTagConst.BIND_REQUEST)
            .writeInteger(LDAP_VERSION3)
            .writeOctetString(dn)
            .writeOctetString(Asn1IdConst.TAG_CLASS_CONTEXT, password)
            .endSequence();
}
```
**Go** (`elements.go:58-59`):
```go
func (r *BindRequest) WriteTo(buffer *asn1.BerBuffer) {
}
```
The method body is entirely empty. It should:
1. Call `buffer.BeginSequence(LdapTagConst.BIND_REQUEST)`
2. Write integer `3` (LDAP version 3)
3. Write octet string for `dn`
4. Write octet string for `password` with `TAG_CLASS_CONTEXT`
5. Call `buffer.EndSequence()`

### BindResponse.java
Checked methods: decode(BerBuffer buffer)
Now I have all the information needed for a thorough comparison.

## Bugs in `BindResponse.Decode` (Go) vs `BindResponse.decode` (Java)

### 1. `BindResponse` struct is empty — missing all fields

**Java:** `BindResponse extends LdapResult` and inherits fields: `resultCode` (int), `matchedDn` (String), `diagnosticMessage` (String), `referrals` (List<String>).

**Go:** `BindResponse` struct is completely empty:
```go
type BindResponse struct {
}
```

Missing fields: `ResultCode`, `MatchedDn`, `DiagnosticMessage`, `Referrals`.

### 2. `BindResponse.Decode` method body is empty — no decoding logic at all

**Java** `decode(BerBuffer buffer)`:
```java
buffer.skipTagAndLength();
LdapResult result = LdapResult.decodeResult(buffer);
return new BindResponse(
        result.getResultCode(),
        result.getMatchedDn(),
        result.getDiagnosticMessage(),
        result.getReferrals());
```

**Go** `Decode(buffer *asn1.BerBuffer)`:
```go
func (r *BindResponse) Decode(buffer *asn1.BerBuffer) {
}
```

Missing logic:
- **`buffer.skipTagAndLength()`** — not called
- **`LdapResult.decodeResult(buffer)`** equivalent — not performed. This means the following are all missing:
  - `buffer.readEnumeration()` for result code
  - `buffer.readOctetString()` for matched DN
  - `buffer.readOctetString()` for diagnostic message
  - Referral parsing loop (check REFERRAL tag, read length, iterate reading octet strings)
- **Returning/constructing a new `BindResponse`** with the decoded values — the Go method returns nothing and assigns nothing to the receiver

### 3. `Decode` does not return a value — signature mismatch

**Java:** `decode` returns a new `BindResponse` instance.

**Go:** `Decode` has no return value (`void`-like); it operates on the receiver `*BindResponse`, but since the struct has no fields and the body is empty, it is a complete no-op.

### 4. `BindResponse` lacks a static DECODER singleton

**Java** defines `public static final BindResponse DECODER` used as the decoder instance. **Go** has no equivalent.

### ModifyRequest.java
Checked methods: estimateSize(), writeTo(BerBuffer buffer)
## Bugs in Go `ModifyRequest` vs Java `ModifyRequest`

### `ModifyRequest` struct (line 72-73)

**Missing fields**: The Go struct is completely empty. It is missing:
- `DN string` — maps to Java's `dn` field
- `Changes []ModifyOperationChange` — maps to Java's `changes` field

---

### `EstimateSize()` (lines 76-78)

**Bug: Returns hardcoded 0 instead of the actual estimate.**

Java logic:
```java
return dn.length() + changes.size() * 32;
```

Go returns `return 0` — no computation at all. Should return `len(r.DN) + len(r.Changes) * 32`.

---

### `WriteTo()` (lines 81-82)

**Bug: Method body is empty — the entire BER/LDAP serialization logic is missing.**

The Java implementation performs this sequence:
1. `buffer.beginSequence(LdapTagConst.MODIFY_REQUEST)` — outer sequence with modify request tag
2. `buffer.writeOctetString(dn)` — write the DN
3. `buffer.beginSequence()` — begin changes sequence
4. For each change:
   - `buffer.beginSequence()` — begin single change sequence
   - `buffer.writeEnumeration(type.getValue())` — write operation type (add/delete/replace)
   - `buffer.beginSequence()` — begin attribute sequence
   - `buffer.writeOctetString(attribute.getType())` — write attribute type
   - `buffer.beginSequence(Asn1IdConst.TAG_SEQUENCE | Asn1IdConst.FORM_CONSTRUCTED)` — begin values sequence with explicit constructed SEQUENCE tag
   - For each attribute value: `buffer.writeOctetString(value)`
   - `buffer.endSequence()` — end values
   - `buffer.endSequence()` — end attribute
   - `buffer.endSequence()` — end change
5. `buffer.endSequence()` — end changes
6. `buffer.endSequence()` — end outer

The Go version does none of this.

---

### Summary

The Go `ModifyRequest` is an unimplemented stub. It has no fields, no `EstimateSize` logic, and no `WriteTo` serialization logic. Both methods are empty no-ops compared to the fully functional Java implementation.

### ModifyResponse.java
Checked methods: decode(BerBuffer buffer)
## Bugs in `ModifyResponse.Decode`

### Bug 1: Method body is entirely unimplemented (stub)

**Java** (`ModifyResponse.decode`):
```java
public ModifyResponse decode(BerBuffer buffer) {
    buffer.skipTagAndLength();
    LdapResult result = LdapResult.decodeResult(buffer);
    return new ModifyResponse(
            result.getResultCode(),
            result.getMatchedDn(),
            result.getDiagnosticMessage(),
            result.getReferrals());
}
```

**Go** (`ModifyResponse.Decode` at line 90):
```go
func (r *ModifyResponse) Decode(buffer *asn1.BerBuffer) {
}
```

The Go method is an empty stub with no logic at all. It is missing:

1. **`buffer.skipTagAndLength()`** — the tag-and-length skip is never performed.
2. **`LdapResult.decodeResult(buffer)`** — the LDAP result is never decoded.
3. **Returning a new `ModifyResponse`** — the Java method returns a new `ModifyResponse` populated with `resultCode`, `matchedDn`, `diagnosticMessage`, and `referrals` from the decoded `LdapResult`. The Go method returns nothing (no return value, no field assignments).

### Bug 2: `ModifyResponse` struct has no fields

**Java** (`ModifyResponse` inherits from `LdapResult`):
```java
public class ModifyResponse extends LdapResult implements ProtocolOperation<ModifyResponse>
```
`LdapResult` holds `resultCode`, `matchedDn`, `diagnosticMessage`, and `referrals`.

**Go** (`ModifyResponse` at line 86):
```go
type ModifyResponse struct {
}
```

The struct is empty. Even if `Decode` were implemented, there are no fields to populate. The Go type does not embed or contain `LdapResult` fields (`ResultCode`, `MatchedDN`, `DiagnosticMessage`, `Referrals`).

### Bug 3: `LdapResult` struct is also an empty stub

**Go** (`LdapResult` at line 30):
```go
type LdapResult struct {
}
```

The `LdapResult` type that should be decoded from the buffer (via `LdapResult.decodeResult`) has no fields and no `DecodeResult` method, so there is no supporting infrastructure to implement the Java logic even partially.

### Filter.java
Checked methods: write(BerBuffer buffer, String filter)
## Bugs

### 1. `Filter.Write` is an empty stub — zero logic ported

**Java** (`Filter.java:57-60`): The `write` method converts the filter string to bytes and calls `writeFilter(buffer, filterBytes, filterBytes.length)`, which triggers a complex recursive parsing engine (~540 lines) including:
- `writeFilter` — main recursive filter parser handling `(`, `)`, spaces, nested filter expressions
- `writeFilterSet` — handles `&` (AND), `|` (OR), `!` (NOT) filter sets with `beginSequence`/`endSequence`
- `writeFilterInSet` — parses individual filters inside a set, including NOT-filter single-child validation
- `writeFilter` (4-arg overload) — parses filter type (`<`, `>`, `~`, `:`, `=`), attribute descriptions, extensible match, substring, present, equality
- `writeSubstringFilter` — handles `*`-based substring filters (initial/any/final)
- `writeExtensibleMatchFilter` — handles extensible match filters with `:dn`, matching rules, attribute types
- `unescapeFilterValue` — handles `\XX` hex-escaped filter values
- `findUnescaped` — finds unescaped `*` in filter values
- `findClosingParenIndex` — finds matching closing parenthesis respecting escapes
- `validate` / `validateAttributeDescription` — validates OID and attribute description characters

**Go** (`elements.go:99-100`): The method body is completely empty:
```go
func (f *Filter) Write(buffer *asn1.BerBuffer, filter string) {
}
```

None of the Java logic has been ported. The method does nothing — no parsing, no BER encoding, no validation, no error handling. It silently accepts any input and produces no output, which will cause LDAP search operations to fail entirely.

### ServiceAvailabilityHandler.java
Checked methods: channelRegistered(ChannelHandlerContext ctx), exceptionCaught(ChannelHandlerContext ctx, Throwable cause)
Here are the bugs found:

## Bugs in `channel_handler.go` vs Java `ServiceAvailabilityHandler`

### 1. `HandleException` — Missing `OutOfDirectMemoryError` handling

**Java** (line 92-94):
```java
} else if (cause instanceof OutOfDirectMemoryError) {
    ctx.close();
}
```

**Go** (`channel_handler.go`): The `else if` branch for `OutOfDirectMemoryError` is completely absent. In the Java version, when an `OutOfDirectMemoryError` occurs, the channel is closed. The Go version only handles `CorruptedFrameException` and does nothing for out-of-memory errors, meaning the connection will remain open when it should be closed.

### 2. `HandleException` — Missing `ctx.fireExceptionCaught(cause)` propagation

**Java** (line 95):
```java
ctx.fireExceptionCaught(cause);
```

**Go**: After handling the exception (whether it was a corrupted frame, out-of-memory, or neither), the Java version **always** propagates the exception up the pipeline via `ctx.fireExceptionCaught(cause)`. The Go version does not propagate the error at all — it simply returns silently after the `if` block. This means downstream handlers will never see the exception.

### 3. `HandleException` — Missing `ctx.close()` for corrupted frame path

**Java** (lines 80-86): After blocking the IP and users for a corrupted frame, the Java code does **not** explicitly close the context here, but the exception is propagated (see bug #2). More critically, note that in the Java version, after handling `CorruptedFrameException`, execution continues to line 95 (`ctx.fireExceptionCaught(cause)`), which typically results in the channel being closed by upstream handlers.

**Go**: After handling the corrupted frame, the method just returns — no error propagation and no connection closure. The connection associated with a corrupted frame attacker stays open.

### 4. `HandleConnection` — Non-TCP address silently allowed through

**Java** (lines 60-67): If `socketAddress` is **not** an `InetSocketAddress`, the code falls through to `ctx.fireChannelRegistered()` (line 68), meaning the connection is allowed.

**Go** (lines 35-41): If `addr` is not a `*net.TCPAddr` (e.g., a UDP address), the code also falls through to `return true` (line 41). This matches Java behavior — **not a bug**, included for completeness of the comparison.

### 5. `GetLocalUserSession` — Return type difference may hide nil-check gap

**Java** (line 86-91): Returns `Queue<UserSession>` which can be `null`, and the code checks `if (sessions != null)` before iterating.

**Go** (line 56-58): Returns `[]*session.UserSession` (a slice). A nil slice is iterable in Go with a `for range` loop (it just won't iterate), so this is not a crash risk. However, the Go interface declares the return as a non-nil slice with no error return, which is a reasonable adaptation. **Not a functional bug**, but worth noting the semantic difference.

### NetConnection.java
Checked methods: getAddress(), send(ByteBuf buffer), close(CloseReason closeReason), close(), switchToUdp(), tryNotifyClientToRecover()
Now I have all the information needed. Let me compare the two implementations method by method.

## Bugs Found

### 1. `getAddress()` / `GetAddress()` — Missing Implementation (Interface Only)

The Java version declares `getAddress()` as an abstract method returning `InetSocketAddress`. The Go version declares it in the interface as `GetAddress() net.Addr` — this is a valid abstract/interface declaration and relies on implementations to provide it. **Not a bug in itself**, but note the return type changed from `InetSocketAddress` to `net.Addr` (broader type), which could allow non-IP addresses that the Java version would not.

### 2. `send(ByteBuf buffer)` / `Send(ctx context.Context, buffer []byte) error` — Missing Implementation (Interface Only)

Same situation — abstract/interface method only. The Go version adds a `context.Context` parameter and uses `[]byte` instead of `ByteBuf`. **Not a bug in itself** as these are reasonable Go adaptations.

### 3. `close(CloseReason closeReason)` / `CloseWithReason(reason CloseReason) error` — **POTENTIAL DEADLOCK**

`CloseWithReason` acquires `b.mu` (a `sync.Mutex`), and `SwitchToUdp` calls `CloseWithReason`. But `TryNotifyClientToRecover` also acquires `b.mu`. This is fine as long as none of these methods call each other in a way that causes re-entrant locking.

However, `SwitchToUdp()` calls `CloseWithReason()` which acquires `b.mu` — this is correct since `SwitchToUdp` does **not** hold the lock itself. **No deadlock bug here.**

**The logic is correct**: sets `isConnected = false`, `isConnectionRecovering = false`, and `isSwitchingToUdp = (reason.Status == SWITCH)`. This matches the Java version exactly.

### 4. `close()` / `Close() error` — **Correct**

Sets all three flags to `false`, matching Java exactly.

### 5. `switchToUdp()` / `SwitchToUdp()` — **Correct**

Calls `CloseWithReason` with `SWITCH` status, matching Java's delegation to `close(CloseReason.get(SessionCloseStatus.SWITCH))`.

### 6. `tryNotifyClientToRecover()` / `TryNotifyClientToRecover()` — **UDP Signal Send is Commented Out**

The Go version at line 79-80 has the actual UDP notification logic **commented out**:

```go
// Mock logic: UdpRequestDispatcher.instance.sendSignal
// udpdispatcher.SendSignal(b.udpAddress, udp.OPEN_CONNECTION)
```

In the Java version, this is the core purpose of the method — it sends a UDP signal to the client:

```java
UdpRequestDispatcher.instance.sendSignal(udpAddress, UdpNotificationType.OPEN_CONNECTION);
```

**This is a functional bug**: `TryNotifyClientToRecover` sets `isConnectionRecovering = true` but **never actually sends the UDP signal**. The method is a no-op for its primary purpose.

---

## Summary

| Method | Status | Issue |
|--------|--------|-------|
| `GetAddress()` | OK | Interface declaration, matches intent |
| `Send()` | OK | Interface declaration, matches intent |
| `CloseWithReason()` | OK | Logic matches Java |
| `Close()` | OK | Logic matches Java |
| `SwitchToUdp()` | OK | Correctly delegates to `CloseWithReason` |
| `TryNotifyClientToRecover()` | **BUG** | UDP signal send is commented out; sets `isConnectionRecovering = true` without actually notifying the client |

### ExtendedHAProxyMessageReader.java
Checked methods: channelRead(ChannelHandlerContext ctx, Object msg)
## Bugs

### 1. Missing PROXY protocol header parsing — entire `channelRead` logic is unimplemented

The Java `channelRead` method performs HAProxy PROXY protocol parsing:
- Checks if `msg` is an `HAProxyMessage` (i.e., a parsed PROXY protocol header)
- Extracts `sourceAddress` and `sourcePort` from the proxy message
- Creates an unresolved `InetSocketAddress` from those values if they are valid (`sourceAddress != null && sourcePort > 0`)
- Falls back to the channel's actual `remoteAddress` if parsing yields no valid address
- Calls `proxyMessage.release()` in a `finally` block (resource cleanup)
- Removes itself from the pipeline (`ctx.channel().pipeline().remove(this)`) so it only processes one PROXY header
- Calls `ctx.read()` to resume reading after handler removal

The Go `Read` method does **none** of this. It skips PROXY protocol header parsing entirely and directly calls `r.OnRemoteAddressConfirmed(conn.RemoteAddr())`. This means:
- **No PROXY protocol v1/v2 header is ever parsed.** The comment says "Pending implementation" but the method is already wired up and will be called in production.
- The callback always receives the direct TCP remote address (the proxy's IP), never the actual client IP from the PROXY header. This is the core purpose of this class — to extract the real client IP from HAProxy's PROXY protocol.
- No fallback logic exists (Java: if `remoteAddress == null`, use `ctx.channel().remoteAddress()`).
- No equivalent of handler self-removal from the pipeline.
- No equivalent of `ctx.read()` to resume reading.

### 2. Missing non-HAProxyMessage pass-through

The Java code has an `else` branch: when `msg` is **not** an `HAProxyMessage`, it delegates to `super.channelRead(ctx, msg)`, passing the message along the pipeline unchanged. The Go implementation has no equivalent pass-through mechanism — any non-PROXY-protocol data is simply ignored since `Read()` is called explicitly for PROXY handling.

### 3. Missing resource cleanup (equivalent of `proxyMessage.release()`)

The Java code calls `proxyMessage.release()` in a `finally` block to free the reference-counted message. The Go version has no cleanup of the parsed PROXY message, since it never parses one.

### HAProxyUtil.java
Checked methods: addProxyProtocolHandlers(ChannelPipeline pipeline, Consumer<InetSocketAddress> onRemoteAddressConfirmed), addProxyProtocolDetectorHandler(ChannelPipeline pipeline, Consumer<InetSocketAddress> onRemoteAddressConfirmed)
## Bugs Found

### 1. `AddProxyProtocolHandlers` — Entire method body is a stub with no implementation

**Java** (`HAProxyUtil.java:43-50`):
```java
public static void addProxyProtocolHandlers(
        ChannelPipeline pipeline,
        Consumer<InetSocketAddress> onRemoteAddressConfirmed) {
    pipeline.addFirst(NettyPipeline.ProxyProtocolDecoder, new HAProxyMessageDecoder())
            .addAfter(NettyPipeline.ProxyProtocolDecoder,
                    NettyPipeline.ProxyProtocolReader,
                    new ExtendedHAProxyMessageReader(onRemoteAddressConfirmed));
}
```

**Go** (`haproxy.go:34-36`):
```go
func AddProxyProtocolHandlers(callback func(net.Addr)) {
    // Pending implementation: Integrate with Go's net package or custom pipeline
}
```

- The function takes only a `callback` parameter; the `pipeline` parameter is missing entirely. The Java version takes a `ChannelPipeline` and adds two handlers to it in a specific order.
- The function body is empty — it does not add any PROXY protocol decoder or reader. The Java version adds `HAProxyMessageDecoder` first, then `ExtendedHAProxyMessageReader` after it.
- No PROXY protocol header parsing or decoding will occur when this function is called.

### 2. `AddProxyProtocolDetectorHandler` — Entire method body is a stub with no implementation

**Java** (`HAProxyUtil.java:52-57`):
```java
public static void addProxyProtocolDetectorHandler(
        ChannelPipeline pipeline,
        Consumer<InetSocketAddress> onRemoteAddressConfirmed) {
    pipeline.addFirst(NettyPipeline.ProxyProtocolDecoder,
            new ExtendedHAProxyMessageDetector(onRemoteAddressConfirmed));
}
```

**Go** (`haproxy.go:39-41`):
```go
func AddProxyProtocolDetectorHandler(callback func(net.Addr)) {
    // Pending implementation: Integrate with Go's net package or custom pipeline
}
```

- Same as above: the `pipeline` parameter is missing, and the function body is empty.
- The Java version adds an `ExtendedHAProxyMessageDetector` handler at the first position in the pipeline. The Go version does nothing.

### 3. Missing `ExtendedHAProxyMessageDetector` type

The Java code references `ExtendedHAProxyMessageDetector` (used in `addProxyProtocolDetectorHandler`), but the Go file only defines `ExtendedHAProxyMessageReader`. There is no `ExtendedHAProxyMessageDetector` struct or any equivalent detection-only handler.

### 4. `ExtendedHAProxyMessageReader.Read` — Stub implementation with incorrect behavior

**Java behavior**: The `ExtendedHAProxyMessageReader` (a Netty `ChannelInboundHandlerAdapter`) receives a decoded `HAProxyMessage` object from the pipeline, extracts the source address from the PROXY protocol header, and passes it to the callback.

**Go** (`haproxy.go:21-28`):
```go
func (r *ExtendedHAProxyMessageReader) Read(conn net.Conn) error {
    if r.OnRemoteAddressConfirmed != nil {
        r.OnRemoteAddressConfirmed(conn.RemoteAddr())
    }
    return nil
}
```

- It passes `conn.RemoteAddr()` (the direct TCP connection's remote address) instead of parsing PROXY protocol headers to extract the *real* client address. This defeats the entire purpose of the PROXY protocol, which is to communicate the original client IP through a proxy/load balancer.
- The comment on line 22 acknowledges this: "Pending implementation: Read initial bytes to parse PROXY v1/v2 headers."

### 5. Callback type mismatch: `func(net.Addr)` vs `Consumer<InetSocketAddress>`

The Java code uses `Consumer<InetSocketAddress>`, which is typed to `InetSocketAddress` (an IP-based socket address). The Go code uses `func(net.Addr)`, which is a broader interface that could include non-IP addresses (e.g., Unix domain sockets). Downstream code expecting an IP-based address may need type assertions that wouldn't be necessary in the Java version.

### SearchRequest.java
Checked methods: estimateSize(), writeTo(BerBuffer buffer)
## Bugs in `SearchRequest.estimateSize()` and `SearchRequest.writeTo()`

### 1. `EstimateSize()` returns wrong value
- **Java:** `return 128;`
- **Go:** `return 0;`
- **Bug:** The estimate is hardcoded to 0 instead of 128.

### 2. `WriteTo()` is an empty stub — all core logic missing
- **Java `writeTo`** performs the full LDAP search request serialization:
  1. `buffer.beginSequence(LdapTagConst.SEARCH_REQUEST)` — opens a tagged sequence with the search request tag
  2. `.writeOctetString(baseDn)` — writes the base DN
  3. `.writeEnumeration(scope)` — writes the scope enum
  4. `.writeEnumeration(derefAliases)` — writes the deref aliases enum
  5. `.writeInteger(sizeLimit)` — writes the size limit
  6. `.writeInteger(timeLimit)` — writes the time limit
  7. `.writeBoolean(typesOnly)` — writes the types-only flag
  8. `Filter.write(buffer, filter)` — writes the search filter
  9. `buffer.beginSequence().writeOctetStrings(attributes).endSequence()` — writes the attribute list in a sub-sequence
  10. `buffer.endSequence()` — closes the outer sequence

- **Go `WriteTo`:** method body is completely empty — no serialization logic at all.

### 3. `SearchRequest` struct has no fields
- **Java:** Has fields `baseDn`, `scope`, `derefAliases`, `sizeLimit`, `timeLimit`, `typesOnly`, `attributes`, `filter` that are used by both methods.
- **Go:** `SearchRequest` is an empty struct `{}` with no fields, so even if the methods had logic, there would be no data to serialize.

### TcpConnection.java
Checked methods: getAddress(), send(ByteBuf buffer), close(CloseReason closeReason), close()
## Bugs in Go Port of `TcpConnection`

### 1. `close(CloseReason)` — Missing check for `connection.isDisposed()`

**Java (line 71):**
```java
if (!isConnected() || connection.isDisposed()) {
    return;
}
```

**Go:**
```go
if !c.IsConnected() {
    return nil
}
```

The Go version only checks `IsConnected()` but omits the `connection.isDisposed()` equivalent. In Java, even if `isConnected` were somehow `true`, the method returns early if the underlying Netty connection is already disposed. The Go version has no corresponding check for whether `c.conn` is already closed/broken before proceeding.

---

### 2. `close(CloseReason)` — Missing notification send before closing

**Java (lines 75-97):** The method sends a close notification buffer (`NotificationFactory.createBuffer(closeReason)`) to the client **before** closing the connection, with retry logic (`RETRY_SEND_CLOSE_NOTIFICATION`), error handling, and a timeout-based close strategy.

**Go (lines 38-45):** The code has a comment `// Pending logic to send notification before closing, similar to Java` and immediately calls `c.conn.Close()`. The entire notification-sending pipeline is missing:
- No call to send a close notification to the client
- No retry logic
- No error handling for notification send failures
- No timeout handling (the `closeTimeout`-based `doFinally` / `onTerminate` / `timeout` chain)

---

### 3. `close(CloseReason)` — Missing `closeTimeout` branching logic

**Java (lines 83-90):** Two distinct timeout strategies:
- `closeTimeout.isZero()` → close immediately after send completes (`doFinally`)
- `!closeTimeout.isNegative()` → wait for `onTerminate` with timeout, then close

**Go:** No timeout handling at all. The connection is closed immediately regardless of `closeTimeout` value.

---

### 4. `close()` — Missing error classification and logging

**Java (lines 101-112):** Catches exceptions from `connection.dispose()` and:
- Checks `ThrowableUtil.isDisconnectedClientError(e)` — only logs if it's **not** a disconnected-client error
- Logs the remote host address from `getAddress().getAddress().getHostAddress()`

**Go (lines 48-54):** Returns the raw error from `c.conn.Close()` without any classification or conditional logging. There is no distinction between expected disconnect errors and unexpected errors.

---

### 5. `close()` — Premature base state clearing before actual connection close

**Go (lines 49-53):**
```go
if !c.IsConnected() {
    return nil
}
c.BaseNetConnection.Close()
return c.conn.Close()
```

**Java (lines 101-112):** `close()` in `TcpConnection` only calls `connection.dispose()`. The base `NetConnection.close()` (which sets `isConnected = false`) is **not** called from `TcpConnection.close()` — it's only called from `TcpConnection.close(CloseReason)` via `super.close(closeReason)`.

In the Java version, `TcpConnection.close()` does **not** call `super.close()`. It only disposes the connection. The `isConnected` flag is left unchanged. The Go version incorrectly calls `c.BaseNetConnection.Close()`, which sets `isConnected = false`, `isConnectionRecovering = false`, `isSwitchingToUdp = false` — this changes state that the Java version intentionally preserves.

---

### 6. `GetAddress()` — Returns `net.Addr` instead of `*net.TCPAddr` (type precision)

**Java (line 57):** Returns `InetSocketAddress` (the concrete, most-specific address type).

**Go (line 27):** Returns `net.Addr` (the interface). This is a minor precision loss — `c.conn.RemoteAddr()` will typically return a `*net.TCPAddr` at runtime, but the return type doesn't enforce it, losing the ability to access TCP-specific fields without a type assertion.

### TcpServerFactory.java
Checked methods: create(TcpProperties tcpProperties, BlocklistService blocklistService, ServerStatusManager serverStatusManager, SessionService sessionService, ConnectionListener connectionListener, int maxFrameLength)
## Bugs

### 1. `CreateWithArgs` is an empty stub — the entire `create(...)` method body is missing

The Go method `CreateWithArgs` (line 65) has an empty body `{}`. The Java `create(...)` method (lines 60–171) contains substantial logic:

- **ServiceAvailabilityHandler creation** (lines 67–70): Instantiates a handler with blocklistService, serverStatusManager, and sessionService. Completely absent from Go.
- **Host/port extraction** from tcpProperties (lines 71–72).
- **Proxy protocol mode handling** (lines 73–75): Extracts `RemoteAddressSourceProxyProtocolMode` from properties.
- **remoteAddressSink** (line 77): Creates a `Sinks.One<InetSocketAddress>` for resolving the remote address asynchronously.
- **TcpServer configuration** (lines 78–155): The entire pipeline configuration including:
  - Socket options: `CONNECT_TIMEOUT_MILLIS`, `SO_REUSEADDR`, `SO_BACKLOG`, `SO_LINGER`, `TCP_NODELAY` (lines 81–88)
  - `wiretap` (line 89)
  - `runOn` with custom loop resources (line 90)
  - `metrics` with `TurmsMicrometerChannelMetricsRecorder` (lines 91–93)
  - **`doOnChannelInit` handler** (lines 95–133):
    - `serviceAvailabilityHandler` added to pipeline (line 97)
    - `varintLengthBasedFrameDecoder` inbound handler (lines 99–101)
    - **Proxy protocol handlers** (lines 102–122): REQUIRED vs OPTIONAL vs no proxy mode, with IP blocklist checking against `remoteAddressSink`
    - `varintLengthFieldPrepender` outbound handler (lines 125–126)
    - `protobufFrameEncoder` outbound handler (lines 131–132)
  - **`handle` callback** (lines 135–155):
    - `setAutoRead(true)` on the channel (line 148)
    - `remoteAddressSink.asMono().flatMap(...)` that delegates to `connectionListener.onAdded(...)`
- **SSL/TLS configuration** (lines 156–160): Conditional SSL setup via `SslUtil.configureSslContextSpec`.
- **Server binding** (lines 161–170): Binds the server and wraps failures in a `BindException`.

All of this is missing from the Go implementation.

### 2. All parameters to `CreateWithArgs` are typed as `any` instead of concrete types

The Java signature uses:
- `TcpProperties tcpProperties`
- `BlocklistService blocklistService`
- `ServerStatusManager serverStatusManager`
- `SessionService sessionService`
- `ConnectionListener connectionListener`

The Go version uses `any` for all five (line 65), losing all type safety and making it impossible to call any methods on these parameters.

### 3. `CreateWithArgs` does not return a value

The Java method returns `DisposableServer` (the bound server instance). The Go method returns nothing (`void` equivalent). The caller has no way to obtain the running server handle.

### SearchResult.java
Checked methods: decode(BerBuffer buffer), isComplete(), shouldLogHeartbeatRequest(), shouldLogRequest(TurmsRequest.KindCase requestType), shouldLogNotification(TurmsRequest.KindCase requestType), log(@Nullable Integer sessionId, @Nullable Long userId, @Nullable DeviceType deviceType, @Nullable Integer version, String ip, long requestId, TurmsRequest.KindCase requestType, int requestSize, long requestTime, TurmsNotification response, long processingTime), log(@Nullable Integer sessionId, @Nullable Long userId, @Nullable DeviceType deviceType, @Nullable Integer version, String ip, long requestId, TurmsRequest.KindCase requestType, int requestSize, long requestTime, int responseCode, long processingTime), log(@Nullable Integer sessionId, @Nullable Long userId, @Nullable DeviceType deviceType, @Nullable Integer version, String ip, long requestId, String requestType, int requestSize, long requestTime, int responseCode, @Nullable String responseDataType, int responseSize, long processingTime), log(SimpleTurmsNotification notification, int notificationBytes, int recipientCount, int onlineRecipientCount)
The user asked me to review specific methods. Most of the methods listed (`shouldLogHeartbeatRequest`, `shouldLogRequest`, `shouldLogNotification`, `log(...)`) are not part of the `SearchResult` Java file at all — they appear to belong to a different class (likely an access logging class). The only methods from the requested list that exist in `SearchResult.java` are `decode(BerBuffer buffer)` and `isComplete()`.

Let me focus on what exists in both files for `SearchResult`:

## Bugs in Go `SearchResult` Implementation

### 1. `decode(BerBuffer buffer)` — Entire method body is missing

The Go `Decode` method (`elements.go:122-123`) is a no-op stub with an empty body. The Java `decode` method (`SearchResult.java:61-109`) contains significant logic:

- **Missing tag read and switch**: The Java code reads a tag via `buffer.readTag()`, skips length, then switches on the tag value (`SEARCH_RESULT_ENTRY` vs `SEARCH_RESULT_DONE` vs default throwing `LdapException`).
- **Missing `SEARCH_RESULT_ENTRY` branch**: In Java, this branch:
  - Creates an `ArrayList` for result entries (capacity 8) and attributes (capacity 16)
  - Reads an object name via `buffer.readOctetString()`
  - Skips a tag, reads a length, computes an end index, and loops reading `Attribute.decode(buffer)` while readable
  - Decodes controls via `Control.decode(buffer)`
  - Creates a `SearchResultEntry` with the object name, attributes, and controls
  - Returns a new `SearchResult` with resultCode=0, null matchedDn/diagnosticMessage/referrals, the result entries, and an overridden `isComplete()` that returns `false`
- **Missing `SEARCH_RESULT_DONE` branch**: In Java, this branch:
  - Checks if `entries == null` and throws `LdapException` with `ResultCode.PROTOCOL_ERROR` if so
  - Calls `LdapResult.decodeResult(buffer)` to get result fields
  - Returns a new `SearchResult` with the decoded result fields and the existing `entries`
- **Missing default branch**: In Java, throws `LdapException` for unexpected tags

### 2. `isComplete()` — Always returns `false` instead of default `true`

The Go `IsComplete()` method (`elements.go:126-128`) always returns `false`. In Java (`SearchResult.java:42`), `SearchResult` extends `LdapResult`, and `isComplete()` is not overridden in the base class — its default implementation returns `true`. The Java code only returns `false` via an anonymous subclass override in the `SEARCH_RESULT_ENTRY` branch of `decode`. The Go implementation hardcodes `false` universally, which is incorrect for the `SEARCH_RESULT_DONE` case (where `isComplete()` should return `true`).

### 3. Missing struct fields

The Go `SearchResult` struct (`elements.go:118-119`) is empty. The Java class has:
- `entries` field (`List<SearchResultEntry>`) — used in both `decode` branches
- Inherited fields from `LdapResult`: `resultCode`, `matchedDn`, `diagnosticMessage`, `referrals`

The remaining methods listed in the request (`shouldLogHeartbeatRequest`, `shouldLogRequest`, `shouldLogNotification`, all `log(...)` overloads, and `log(SimpleTurmsNotification...)`) do not exist in the `SearchResult.java` file and therefore cannot be compared against the Go file. Those methods belong to a different Java class (likely an access/log handler) and were not included in the files provided for review.

### TcpUserSessionAssembler.java
Checked methods: getHost(), getPort()
## Bugs in `getHost()` and `getPort()`

### 1. `GetHost()` / `GetPort()` — Guard condition checks wrong field

**Java** checks `server == null`:
```java
public String getHost() {
    if (server == null) {
        throw new FeatureDisabledException("TCP server is disabled");
    }
    return host;
}
```

**Go** checks `!a.Enabled`:
```go
func (a *TcpUserSessionAssembler) GetHost() string {
    if !a.Enabled {
        panic("TCP server is disabled")
    }
    return a.Host
}
```

While these are functionally equivalent in the Java code (when `enabled == false`, `server` is set to `null`), the Go code is actually a **reasonable semantic equivalent** — not a true bug. The `Enabled` field mirrors the Java `enabled` field, and when disabled, `Server` (the Go analog of `server`) would also be nil. This is acceptable.

### 2. `GetHost()` / `GetPort()` — Error handling uses `panic` instead of returning an error

**Java** throws `FeatureDisabledException`:
```java
throw new FeatureDisabledException("TCP server is disabled");
```

**Go** uses `panic`:
```go
panic("TCP server is disabled")
```

In Go, `panic` is not idiomatic for expected/control-flow errors like a disabled feature. The Java code uses a checked/declared exception type (`FeatureDisabledException`) that callers can handle. The Go port should return an error (e.g., `("", error)`) or use a domain-specific error type to allow callers to handle the "disabled" case gracefully. `panic` will crash the goroutine if unrecovered, unlike the Java exception which is catchable.

### 3. `GetHost()` — Default `Host` is empty string instead of `null`

**Java** sets `host = null` when disabled:
```java
host = null;
```

**Go** sets `Host = ""` when disabled:
```go
Host: "",
```

While Go doesn't have nullable strings, this is a minor semantic difference. The `""` value is never actually returned because the guard condition fires first, so this is **not a functional bug** — but worth noting as a deviation.

### 4. `GetPort()` — Default `Port` is `-1` (matches Java)

This matches correctly.

---

**Summary of real bugs:**

| # | Method | Bug | Severity |
|---|--------|-----|----------|
| 1 | `GetHost()`, `GetPort()` | Uses `panic` instead of returning a Go `error`, making the "disabled" case unrecoverable for callers and diverging from the Java exception-handling pattern | Medium |

### UdpRequestDispatcher.java
Checked methods: sendSignal(InetSocketAddress address, UdpNotificationType signal), get(ResponseStatusCode code), get(UdpNotificationType type)
Now I have all the information needed. Let me analyze the three methods:

## Bugs Found

### 1. `sendSignal` — Non-blocking send silently drops notifications

**Java** (`sendSignal`, line 141-145):
```java
public void sendSignal(InetSocketAddress address, UdpNotificationType signal) {
    if (notificationSink != null) {
        notificationSink.tryEmitNext(new UdpNotification(address, signal));
    }
}
```

**Go** (`SendSignal`, line 89-100):
```go
func (d *UdpRequestDispatcher) SendSignal(address net.Addr, signal UdpNotificationType) {
    if d.notificationSink != nil {
        select {
        case d.notificationSink <- UdpNotification{...}:
        default:
            // Handle sink full
        }
    }
}
```

**Bug**: The Java `tryEmitNext` on a `Sinks.many().unicast().onBackpressureBuffer()` will buffer unboundedly — it will never drop a notification. The Go version uses a buffered channel with capacity 1024 and a `select/default` pattern that **silently drops the notification** when the channel is full. The Java version guarantees delivery; the Go version does not.

---

### 2. `get(ResponseStatusCode code)` — Missing caching/pooling, and uses wrong numeric value

**Java** (`get(ResponseStatusCode code)` in `UdpSignalResponseBufferPool`, lines 50-69):
```java
public static ByteBuf get(ResponseStatusCode code) {
    // ... double-checked locking cache lookup ...
    if (code == ResponseStatusCode.OK) {
        buf = Unpooled.EMPTY_BUFFER;  // returns empty buffer
    } else {
        buf = Unpooled.unreleasableBuffer(Unpooled.directBuffer(Short.BYTES)
                .writeShort(code.getBusinessCode()));  // writes getBusinessCode()
    }
    // ... caches and returns buf ...
}
```

**Go** (`GetBufferFromStatusCode`, lines 127-135):
```go
func (d *UdpRequestDispatcher) GetBufferFromStatusCode(code constant.ResponseStatusCode) []byte {
    if code == constant.ResponseStatusCode_OK {
        return []byte{}
    }
    val := uint16(code)
    return []byte{byte(val >> 8), byte(val)}
}
```

**Bugs**:
- **Missing caching**: The Java version uses a `FastEnumMap` cache with double-checked locking so each status code's buffer is allocated only once and reused. The Go version allocates a new `[]byte` slice on every call.
- **Wrong value written**: The Java version writes `code.getBusinessCode()` (a specific business code property), not the enum ordinal. The Go version uses `uint16(code)` which is the **enum ordinal**, not the business code. These are different values — `getBusinessCode()` is a separate field on the Java enum.

---

### 3. `get(UdpNotificationType type)` — Wrong enum constants and incorrect ordinal mapping

**Java** (`UdpNotificationType` enum, line 23-25):
```java
public enum UdpNotificationType {
    OPEN_CONNECTION  // ordinal = 0
}
```

**Java** (`get(UdpNotificationType type)` in `UdpSignalResponseBufferPool`, static init):
```java
// Pre-populates: type.ordinal() + 1
// OPEN_CONNECTION -> byte value 1 (0 + 1)
```

**Go** (`UdpNotificationType`, lines 13-18):
```go
type UdpNotificationType byte
const (
    HeartbeatNotification UdpNotificationType = iota  // 0
    GoOfflineNotification                             // 1
)
```

**Go** (`GetBufferFromNotificationType`, lines 138-141):
```go
func (d *UdpRequestDispatcher) GetBufferFromNotificationType(t UdpNotificationType) []byte {
    return []byte{byte(t) + 1}
}
```

**Bugs**:
- **Wrong enum values**: The Java `UdpNotificationType` has a single constant `OPEN_CONNECTION` (ordinal 0). The Go version has `HeartbeatNotification` (0) and `GoOfflineNotification` (1), which do not exist in the Java source. These constants were fabricated — they are not a faithful port.
- **Missing pre-populated pool**: The Java version pre-populates a `FastEnumMap` pool in a static initializer so the buffer is allocated once and reused. The Go version allocates a new `[]byte` on every call with no caching.
- **Method should be static/stateless**: In Java, this is a `static` method on `UdpSignalResponseBufferPool` with a static cache. In Go, it's a method on the `UdpRequestDispatcher` receiver, but the `UdpRequestDispatcher` struct has no buffer pool field — so no caching mechanism exists at all.

### UdpRequestType.java
Checked methods: parse(int number), getNumber()
## Bugs in `ParseUdpRequestType` and `GetNumber`

### Bug 1: `ParseUdpRequestType` — Missing bounds checking and 1-based offset

**Java logic** (`parse(int number)`):
```java
int index = number - 1;
if (index > -1 && index < ALL.length) {
    return ALL[index];
}
return null;
```
- Converts the 1-based `number` to a 0-based index (`number - 1`).
- Bounds-checks that the index is within the valid enum range `[0, ALL.length)`.
- Returns `null` for out-of-range values.

**Go code** (`ParseUdpRequestType`):
```go
func ParseUdpRequestType(number int) UdpRequestType {
    return UdpRequestType(number)
}
```
- Performs a direct cast with **no bounds checking** — any integer value is accepted, including negative numbers and values exceeding the valid enum range.
- Does **not apply the `number - 1` offset**. In Java, ordinal 0 (`HEARTBEAT`) maps to number 1. The Go `UdpRequestType` uses `iota` (0-based), so to match Java semantics, the function should compute `UdpRequestType(number - 1)`.
- Returns a zero-value `UdpRequestType` instead of a nil/error for invalid input, since Go doesn't have nullable return types here. At minimum it should validate bounds.

### Bug 2: `GetNumber` — Missing 1-based offset

**Java logic** (`getNumber()`):
```java
return this.ordinal() + 1;
```
- Returns the 1-based number (ordinal + 1), so `HEARTBEAT` returns `1` and `GO_OFFLINE` returns `2`.

**Go code** (`GetNumber()`):
```go
func (t UdpRequestType) GetNumber() int {
    return int(t)
}
```
- Returns the raw 0-based `int(t)` value directly, so `HeartbeatRequest` returns `0` and `GoOfflineRequest` returns `1`.
- This is **off by 1** compared to the Java implementation.
