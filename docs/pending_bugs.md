
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
